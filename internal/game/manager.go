package game

import (
	"blockchess/internal/client"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/corentings/chess/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

// Constants
const (
	GameTimerSeconds = 15   // Timer duration in seconds for each turn
	StakeAmount      = 0.01 // 0.01 USDC
)

type Vote struct {
	Move  string
	Count int
}

type GameState struct {
	ID          string
	Votes       map[string]int // move -> vote count
	TimeLeft    int            // seconds
	Game        *chess.Game    // Chess game state from CorentinGS/chess library
	Players     []string       // connected player IDs
	CurrentMove int            // Current move number
	CreatedAt   int64          // Unix timestamp when game was created

	// Team tracking with wallet addresses
	WhitePlayers map[string]bool // walletAddress -> true if on white team
	BlackPlayers map[string]bool // walletAddress -> true if on black team

	// Pot tracking
	TotalPot float64
	WhitePot float64
	BlackPot float64

	// Vote tracking per round
	WhiteVotesThisTurn   int
	BlackVotesThisTurn   int
	PlayerVotedThisRound map[string]bool // Track who has voted this round (by wallet address)

	// Total vote tracking (persistent across all rounds)
	WhiteTeamTotalVotes int
	BlackTeamTotalVotes int
	PlayerTotalVotes    map[string]int // walletAddress -> total votes made throughout the game

	// Blockchain integration
	BlockchainGameID uint64 // Game ID from the smart contract

	mu sync.RWMutex
}

// Manager handles game logic and orchestrates services
type Manager struct {
	games              map[string]*GameState
	mu                 sync.RWMutex
	moveResultCallback func(gameID, move string)
	gameEndCallback    func(gameID, winner, reason string, gameStats map[string]any)

	// Blockchain clients for multi-chain operations
	clients        *client.Clients
	gameFactory    *client.GameFactory
	vaultManager   *client.VaultManager
	permit2Manager *client.Permit2Manager

	// Player chain ID mapping - walletAddress -> chainID
	playerChainIDs map[string]uint32
	chainIDMutex   sync.RWMutex

	// Player permit signatures - walletAddress -> permit signature data
	playerPermits map[string]*client.PermitSignatureData
	permitMutex   sync.RWMutex
}

func NewGamesManager(clients *client.Clients) *Manager {
	// Initialize GameFactory with Base Sepolia client
	var gameFactory *client.GameFactory
	baseSepoliaChainID := uint64(84532)
	baseSepoliaClient, err := clients.GetClientByChainID(baseSepoliaChainID)
	if err != nil {
		log.Printf("Warning: Failed to get Base Sepolia client: %v", err)
	} else {
		privateKey := clients.GetPrivateKey()
		if privateKey == "" {
			log.Printf("Warning: No private key available for GameFactory")
		} else {
			gameFactory, err = client.NewGameFactory(baseSepoliaClient, privateKey)
			if err != nil {
				log.Printf("Warning: Failed to initialize GameFactory: %v", err)
			}
		}
	}

	// Initialize VaultManager for all chains
	vaultManager, err := client.NewVaultManager(clients)
	if err != nil {
		log.Printf("Warning: Failed to initialize VaultManager: %v", err)
	} else {
		availableChains := vaultManager.GetAvailableChains()
		log.Printf("VaultManager initialized with %d chains: %v", len(availableChains), availableChains)
	}

	// Initialize Permit2Manager for all chains
	permit2Manager, err := client.NewPermit2Manager(clients)
	if err != nil {
		log.Printf("Warning: Failed to initialize Permit2Manager: %v", err)
	} else {
		availableChains := permit2Manager.GetAvailableChains()
		log.Printf("Permit2Manager initialized with %d chains: %v", len(availableChains), availableChains)
	}

	return &Manager{
		games:          make(map[string]*GameState),
		clients:        clients,
		gameFactory:    gameFactory,
		vaultManager:   vaultManager,
		permit2Manager: permit2Manager,
		playerChainIDs: make(map[string]uint32),
		playerPermits:  make(map[string]*client.PermitSignatureData),
	}
}

// SetMoveResultCallback sets the callback for broadcasting move results
func (m *Manager) SetMoveResultCallback(callback func(gameID, move string)) {
	m.moveResultCallback = callback
}

// SetGameEndCallback sets the callback for broadcasting game end
func (m *Manager) SetGameEndCallback(callback func(gameID, winner, reason string, gameStats map[string]any)) {
	m.gameEndCallback = callback
}

// SetPlayerChainID sets the chain ID for a player (called during matchmaking)
func (m *Manager) SetPlayerChainID(walletAddress string, chainID uint32) {
	m.chainIDMutex.Lock()
	defer m.chainIDMutex.Unlock()
	m.playerChainIDs[walletAddress] = chainID
	log.Printf("Set chain ID %d for player %s", chainID, walletAddress)
}

// GetPlayerChainID gets the chain ID for a player
func (m *Manager) GetPlayerChainID(walletAddress string) uint32 {
	m.chainIDMutex.RLock()
	defer m.chainIDMutex.RUnlock()
	return m.playerChainIDs[walletAddress]
}

// StorePlayerPermit stores a permit signature for a player
func (m *Manager) StorePlayerPermit(walletAddress string, permitData *client.PermitSignatureData) {
	m.permitMutex.Lock()
	defer m.permitMutex.Unlock()
	m.playerPermits[walletAddress] = permitData
	log.Printf("Stored permit signature for player %s on chain %d", walletAddress, permitData.ChainID)
}

// GetPlayerPermit retrieves a permit signature for a player
func (m *Manager) GetPlayerPermit(walletAddress string) *client.PermitSignatureData {
	m.permitMutex.RLock()
	defer m.permitMutex.RUnlock()
	return m.playerPermits[walletAddress]
}

// CreatePermitForPlayer creates a permit signature request for a player
func (m *Manager) CreatePermitForPlayer(walletAddress string, chainID uint32) (*client.PermitSignatureData, interface{}, error) {
	if m.permit2Manager == nil {
		return nil, nil, fmt.Errorf("Permit2 manager not available")
	}

	// Get the Permit2 client for the player's chain
	permit2Client, err := m.permit2Manager.GetPermit2Client(uint64(chainID))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get Permit2 client for chain %d: %w", chainID, err)
	}

	// Get the vault address for the player's chain
	vaultAddress := client.GetVaultAddress(uint64(chainID))
	if vaultAddress == "" {
		return nil, nil, fmt.Errorf("no vault address configured for chain %d", chainID)
	}

	// Create permit signature data
	owner := common.HexToAddress(walletAddress)
	vault := common.HexToAddress(vaultAddress)

	permitData, typedData, err := permit2Client.CreateGameStakePermit(owner, vault)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create permit: %w", err)
	}

	log.Printf("Created permit for player %s on chain %d (vault: %s)", walletAddress, chainID, vaultAddress)
	return permitData, typedData, nil
}

// HasValidPermit checks if a player has a valid permit signature
func (m *Manager) HasValidPermit(walletAddress string, chainID uint32) bool {
	permitData := m.GetPlayerPermit(walletAddress)
	if permitData == nil {
		return false
	}

	// Check if the permit is for the correct chain
	if permitData.ChainID != uint64(chainID) {
		return false
	}

	// Check if the permit signature is present
	if permitData.Signature == "" {
		return false
	}

	// Check if the permit hasn't expired (signature deadline)
	now := big.NewInt(time.Now().Unix())
	if permitData.SigDeadline.Cmp(now) <= 0 {
		log.Printf("Permit expired for player %s on chain %d", walletAddress, chainID)
		return false
	}

	return true
}

// EnsurePlayerPermit ensures a player has a valid permit for their chain
func (m *Manager) EnsurePlayerPermit(walletAddress string, chainID uint32) error {
	if m.HasValidPermit(walletAddress, chainID) {
		return nil // Already has valid permit
	}

	// Check if we have permit data but it's invalid
	permitData := m.GetPlayerPermit(walletAddress)
	if permitData != nil {
		if permitData.ChainID != uint64(chainID) {
			return fmt.Errorf("permit is for chain %d, but player is on chain %d", permitData.ChainID, chainID)
		}
		if permitData.Signature == "" {
			return fmt.Errorf("permit signature is missing - please sign the permit")
		}
		now := big.NewInt(time.Now().Unix())
		if permitData.SigDeadline.Cmp(now) <= 0 {
			return fmt.Errorf("permit has expired - please request a new permit")
		}
	}

	return fmt.Errorf("no permit found - please sign permit before voting")
}

// GetOrCreatePlayerPermit gets existing permit or creates a new one if needed
func (m *Manager) GetOrCreatePlayerPermit(walletAddress string, chainID uint32) (*client.PermitSignatureData, interface{}, error) {
	// Check if we already have a valid permit
	if m.HasValidPermit(walletAddress, chainID) {
		permitData := m.GetPlayerPermit(walletAddress)
		return permitData, nil, nil
	}

	// Create new permit
	return m.CreatePermitForPlayer(walletAddress, chainID)
}

// GetOrCreateGame gets an existing game or creates a new one
func (m *Manager) GetOrCreateGame() *GameState {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate a unique game ID
	gameID := uuid.New().String()

	game := &GameState{
		ID:                   gameID,
		Votes:                make(map[string]int),
		TimeLeft:             GameTimerSeconds, // Use constant for timer duration
		Game:                 chess.NewGame(),
		Players:              make([]string, 0),
		CurrentMove:          1,
		CreatedAt:            time.Now().Unix(),
		WhitePlayers:         make(map[string]bool),
		BlackPlayers:         make(map[string]bool),
		TotalPot:             0,
		WhitePot:             0,
		BlackPot:             0,
		WhiteVotesThisTurn:   0,
		BlackVotesThisTurn:   0,
		WhiteTeamTotalVotes:  0,
		BlackTeamTotalVotes:  0,
		PlayerVotedThisRound: make(map[string]bool),
		PlayerTotalVotes:     make(map[string]int),
	}

	// Create blockchain game if GameFactory is available
	if m.gameFactory != nil {
		// Use fixed stake amount of 0.01 USDC (converted to wei: 0.01 * 10^6)
		stakeAmount := new(big.Int).SetInt64(10000) // 0.01 USDC in USDC's 6 decimal places
		blockchainGameID, err := m.gameFactory.CreateGame(stakeAmount)
		if err != nil {
			log.Fatalf("Failed to create game contract: %v", err)
		} else {
			game.BlockchainGameID = blockchainGameID
			log.Printf("Created game contract with ID: %d for local game: %s", blockchainGameID, gameID)
		}
	} else {
		log.Fatal("No game factory available")
	}

	m.games[game.ID] = game

	// Start game timer
	go m.runGameTimer(game)

	return game
}

// GetGame retrieves an existing game without creating it
func (m *Manager) GetGame(gameID string) *GameState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if game, exists := m.games[gameID]; exists {
		return game
	}

	return nil
}

// runGameTimer manages the game timer and move execution
func (m *Manager) runGameTimer(game *GameState) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		game.mu.Lock()
		game.TimeLeft--

		// Check if we should execute the move early (when only 1 player on current team)
		shouldExecuteEarly := false
		if game.TimeLeft <= GameTimerSeconds-1 { // After 1 second (started at GameTimerSeconds)
			currentTurn := game.Game.Position().Turn()
			var currentTeamPlayerCount int
			if currentTurn == chess.White {
				currentTeamPlayerCount = len(game.WhitePlayers)
			} else {
				currentTeamPlayerCount = len(game.BlackPlayers)
			}

			// If only 1 player on current team AND there's at least one vote, execute after 1 second
			if currentTeamPlayerCount == 1 && len(game.Votes) > 0 {
				shouldExecuteEarly = true
			}
		}

		if game.TimeLeft <= 0 || shouldExecuteEarly {
			// Time's up or early execution - execute the move with most votes
			bestMove := m.getBestMove(game)

			// If no votes were cast, the current team forfeits the game.
			if bestMove == "skip" {
				currentTurn := game.Game.Position().Turn()
				game.Game.Resign(currentTurn) // This sets the game outcome
				log.Printf("Game %s: Team %s forfeited due to inactivity.", game.ID, currentTurn.String())

				// Get final stats before unlocking
				gameStats := m.getGameStatsUnsafe(game, true)
				game.mu.Unlock()

				// Handle game end logic (blockchain update, broadcast)
				m.handleGameEnd(game.ID, gameStats)
				log.Printf("Game %s timer stopped due to forfeit.", game.ID)
				return // Stop the timer goroutine for this game
			}

			// Apply the move to the board if it's not "skip"
			m.applyMoveToBoard(game, bestMove)

			// Check if the game ended after this move
			gameEnded := m.checkGameEnd(game)

			// Reset for next turn immediately to prevent race conditions
			game.Votes = make(map[string]int)
			game.TimeLeft = GameTimerSeconds // Reset to constant value
			game.CurrentMove++

			// Reset round vote tracking
			game.WhiteVotesThisTurn = 0
			game.BlackVotesThisTurn = 0
			game.PlayerVotedThisRound = make(map[string]bool)

			if shouldExecuteEarly {
				log.Printf("Game %s: Move executed early due to single player on team", game.ID)
			}

			game.mu.Unlock()

			// Broadcast move result after reset
			m.BroadcastMoveResult(game.ID, bestMove)

			// Handle game end if applicable
			if gameEnded {
				log.Printf("Game %s ended after move %s!", game.ID, bestMove)

				// Get final stats
				game.mu.RLock()
				gameStats := m.getGameStatsUnsafe(game, true)
				game.mu.RUnlock()

				// Broadcast game end and handle blockchain update
				m.handleGameEnd(game.ID, gameStats)
				log.Printf("Game %s timer stopped", game.ID)
				return
			}
		} else {
			game.mu.Unlock()
		}
	}
}

// handleGameEnd processes game end logic including blockchain updates
func (m *Manager) handleGameEnd(gameID string, gameStats map[string]any) {
	// Determine winner and reason from game state
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return
	}

	game.mu.RLock()
	outcome := game.Game.Outcome()
	method := game.Game.Method()
	blockchainGameID := game.BlockchainGameID
	game.mu.RUnlock()

	if outcome == chess.NoOutcome {
		return
	}

	// Map chess outcomes to winner and reason
	var winner, reason string

	switch outcome {
	case chess.WhiteWon:
		winner = "white"
		if method == chess.Checkmate {
			reason = "checkmate"
		} else {
			reason = "resignation"
		}
	case chess.BlackWon:
		winner = "black"
		if method == chess.Checkmate {
			reason = "checkmate"
		} else {
			reason = "resignation"
		}
	case chess.Draw:
		winner = "draw"
		switch method {
		case chess.Stalemate:
			reason = "stalemate"
		case chess.InsufficientMaterial:
			reason = "insufficient_material"
		case chess.ThreefoldRepetition:
			reason = "threefold_repetition"
		case chess.FiftyMoveRule:
			reason = "fifty_move_rule"
		default:
			reason = "draw"
		}
	}

	// Distribute rewards to winners before ending the game
	if m.vaultManager != nil && blockchainGameID != 0 {
		m.distributeRewards(gameID, winner, gameStats)
	}

	// End the blockchain game if available
	if m.gameFactory != nil && blockchainGameID != 0 {
		result, err := client.ResultStringToUint8(winner)
		if err != nil {
			log.Printf("Warning: Failed to convert result '%s' to uint8: %v", winner, err)
		} else {
			err = m.gameFactory.EndGame(blockchainGameID, result)
			if err != nil {
				log.Printf("Warning: Failed to end blockchain game %d: %v", blockchainGameID, err)
			} else {
				log.Printf("Successfully ended blockchain game %d with result: %s", blockchainGameID, winner)
			}
		}
	}

	// Broadcast game end
	if m.gameEndCallback != nil {
		m.gameEndCallback(gameID, winner, reason, gameStats)
	}
}

// getBestMove returns the move with the most votes
func (m *Manager) getBestMove(game *GameState) string {
	var bestMove string
	var maxVotes int

	for move, votes := range game.Votes {
		if votes > maxVotes {
			maxVotes = votes
			bestMove = move
		}
	}

	if bestMove == "" {
		return "skip"
	}

	return bestMove
}

// VoteForMove allows a player to vote for a move
func (m *Manager) VoteForMove(gameID, walletAddress, move string, team string, chainId uint32) error {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game not found: %s", gameID)
	}

	game.mu.Lock()
	defer game.mu.Unlock()

	// Validate wallet address
	if walletAddress == "" {
		return fmt.Errorf("wallet address cannot be empty")
	}

	// Check if player is on the team they're trying to vote for
	switch team {
	case "white":
		if !game.WhitePlayers[walletAddress] {
			return fmt.Errorf("player not on white team")
		}
	case "black":
		if !game.BlackPlayers[walletAddress] {
			return fmt.Errorf("player not on black team")
		}
	default:
		return fmt.Errorf("invalid team: %s", team)
	}

	// Check if player already voted this round
	if game.PlayerVotedThisRound[walletAddress] {
		return fmt.Errorf("player already voted this round")
	}

	// Validate the move
	if !m.isValidMove(game, move) {
		return fmt.Errorf("invalid move: %s", move)
	}

	// MANDATORY: Ensure player has valid permit before allowing vote
	if chainId != 0 {
		err := m.EnsurePlayerPermit(walletAddress, chainId)
		if err != nil {
			return fmt.Errorf("permit required for voting: %w", err)
		}
	}

	// Stake USDC to vault contract on player's chain using stored permit
	if m.vaultManager != nil && chainId != 0 {
		vault, err := m.vaultManager.GetVault(uint64(chainId))
		if err != nil {
			log.Printf("Warning: Failed to get vault for chain %d: %v", chainId, err)
		} else {
			playerAddress := common.HexToAddress(walletAddress)
			stakeAmount := new(big.Int).SetInt64(10000) // 0.01 USDC in USDC's 6 decimal places

			// Get the stored permit (we already validated it exists above)
			permitData := m.GetPlayerPermit(walletAddress)
			log.Printf("Using stored permit for player %s on chain %d", walletAddress, chainId)
			err = vault.StakeWithPermit(playerAddress, game.BlockchainGameID, stakeAmount, permitData)
			if err != nil {
				log.Printf("Error: Failed to stake with permit to vault on chain %d: %v", chainId, err)
				return fmt.Errorf("staking with permit failed: %w", err)
			} else {
				log.Printf("Successfully staked 0.01 USDC for player %s on chain %d using Permit2", walletAddress, chainId)
			}
		}
	}

	// Add vote to blockchain if available
	if m.gameFactory != nil && game.BlockchainGameID != 0 {
		teamUint8, err := client.TeamStringToUint8(team)
		if err != nil {
			log.Printf("Warning: Failed to convert team '%s' to uint8: %v", team, err)
		} else {
			playerAddress := common.HexToAddress(walletAddress)
			err = m.gameFactory.AddVote(game.BlockchainGameID, playerAddress, chainId, teamUint8)
			if err != nil {
				log.Printf("Warning: Failed to add vote to blockchain: %v", err)
				// Continue with local vote even if blockchain fails
			} else {
				log.Printf("Successfully added vote to blockchain for player %s", walletAddress)
			}
		}
	}

	// Record vote locally
	if _, exists := game.Votes[move]; !exists {
		game.Votes[move] = 0
	}
	game.Votes[move]++
	game.PlayerVotedThisRound[walletAddress] = true
	game.PlayerTotalVotes[walletAddress]++

	// Update team vote count and pot
	switch team {
	case "white":
		game.WhiteVotesThisTurn++
		game.WhiteTeamTotalVotes++
		game.WhitePot += 0.01
		game.TotalPot += 0.01
	case "black":
		game.BlackVotesThisTurn++
		game.BlackTeamTotalVotes++
		game.BlackPot += 0.01
		game.TotalPot += 0.01
	}

	log.Printf("Player %s voted for move %s in team %s (game %s)", walletAddress, move, team, gameID)
	return nil
}

// isValidMove checks if a move is valid in the current game state
func (m *Manager) isValidMove(game *GameState, moveStr string) bool {
	if len(moveStr) < 2 {
		return false
	}

	// Handle coordinate notation (e.g., "e2e4")
	if len(moveStr) == 4 {
		from := parseSquare(moveStr[:2])
		to := parseSquare(moveStr[2:])

		if from == chess.NoSquare || to == chess.NoSquare {
			return false
		}

		validMoves := game.Game.ValidMoves()
		for _, move := range validMoves {
			if move.S1() == from && move.S2() == to {
				return true
			}
		}
		return false
	}

	// Handle algebraic notation
	gameCopy := game.Game.Clone()
	err := gameCopy.PushNotationMove(moveStr, chess.AlgebraicNotation{}, nil)
	return err == nil
}

// parseSquare converts square notation like "e2" into chess.Square
func parseSquare(square string) chess.Square {
	if len(square) != 2 {
		return chess.NoSquare
	}

	file := square[0]
	rank := square[1]

	if file < 'a' || file > 'h' || rank < '1' || rank > '8' {
		return chess.NoSquare
	}

	fileIndex := int(file - 'a')
	rankIndex := int(rank - '1')
	squareIndex := rankIndex*8 + fileIndex

	return chess.Square(squareIndex)
}

// GetValidMoves returns all valid moves for the current position in coordinate notation
func (m *Manager) GetValidMoves(gameID string) []string {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return nil
	}

	game.mu.RLock()
	defer game.mu.RUnlock()

	moves := game.Game.ValidMoves()
	moveStrings := make([]string, len(moves))
	for i, move := range moves {
		// Convert to coordinate notation (e.g., "e2e4" instead of "e4")
		from := move.S1().String()
		to := move.S2().String()
		moveStrings[i] = from + to
	}

	return moveStrings
}

// GetVotes returns a copy of the current votes
func (m *Manager) GetVotes(gameID string) map[string]int {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return nil
	}

	game.mu.RLock()
	defer game.mu.RUnlock()

	votes := make(map[string]int)
	for k, v := range game.Votes {
		votes[k] = v
	}

	return votes
}

// GetTimeLeft returns the time left in the current turn
func (m *Manager) GetTimeLeft(gameID string) int {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return 0
	}

	game.mu.RLock()
	defer game.mu.RUnlock()

	return game.TimeLeft
}

// BroadcastMoveResult broadcasts the result of a move
func (m *Manager) BroadcastMoveResult(gameID, move string) {
	if m.moveResultCallback != nil {
		m.moveResultCallback(gameID, move)
	}
}

// AddPlayerToTeam adds a player to a team with wallet address validation
func (m *Manager) AddPlayerToTeam(gameID, walletAddress, team string) error {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game not found: %s", gameID)
	}

	// Validate wallet address
	if walletAddress == "" {
		return fmt.Errorf("wallet address cannot be empty")
	}

	if len(walletAddress) != 42 || !strings.HasPrefix(walletAddress, "0x") {
		return fmt.Errorf("invalid wallet address format")
	}

	game.mu.Lock()
	defer game.mu.Unlock()

	// Check if player is already on any team
	if game.WhitePlayers[walletAddress] {
		if team == "white" {
			return fmt.Errorf("player already on white team")
		}
		return fmt.Errorf("player cannot join black team - already on white team")
	}

	if game.BlackPlayers[walletAddress] {
		if team == "black" {
			return fmt.Errorf("player already on black team")
		}
		return fmt.Errorf("player cannot join white team - already on black team")
	}

	// Add to requested team
	switch team {
	case "white":
		game.WhitePlayers[walletAddress] = true
		log.Printf("Player %s joined white team in game %s", walletAddress, gameID)
	case "black":
		game.BlackPlayers[walletAddress] = true
		log.Printf("Player %s joined black team in game %s", walletAddress, gameID)
	default:
		return fmt.Errorf("invalid team: %s", team)
	}

	return nil
}

// GetGameStats returns game statistics
func (m *Manager) GetGameStats(gameID string) map[string]any {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return nil
	}

	game.mu.RLock()
	defer game.mu.RUnlock()

	return m.getGameStatsUnsafe(game, false)
}

// getGameStatsUnsafe returns game statistics without locking (caller must hold the lock)
func (m *Manager) getGameStatsUnsafe(game *GameState, gameEnded bool) map[string]any {
	// Convert chess board to string array
	board := make([][]string, 8)
	pieceMap := map[string]string{
		"♜": "r", "♞": "n", "♝": "b", "♛": "q", "♚": "k", "♟": "p", // Black pieces
		"♖": "R", "♘": "N", "♗": "B", "♕": "Q", "♔": "K", "♙": "P", // White pieces
	}

	for i := range 8 {
		board[i] = make([]string, 8)
		for j := range 8 {
			rank := 7 - i // Convert row to rank (row 0 = rank 8)
			file := j     // file stays the same (col 0 = file a)
			square := chess.Square(rank*8 + file)
			piece := game.Game.Position().Board().Piece(square)
			if piece == chess.NoPiece {
				board[i][j] = ""
			} else {
				unicodeSymbol := piece.String()
				if simpleSymbol, exists := pieceMap[unicodeSymbol]; exists {
					board[i][j] = simpleSymbol
				} else {
					board[i][j] = ""
				}
			}
		}
	}

	// Get current turn from chess game
	currentTurn := "white"
	if game.Game.Position().Turn() == chess.Black {
		currentTurn = "black"
	}

	currentMove := game.CurrentMove
	if gameEnded {
		currentMove = currentMove - 1
	}

	// Detect check and checkmate status
	isInCheck := false
	isCheckmate := false

	if !gameEnded {
		// Use the chess library's built-in position status to detect check and checkmate
		position := game.Game.Position()
		outcome := game.Game.Outcome()
		method := game.Game.Method()

		// Check for checkmate first
		if outcome != chess.NoOutcome && method == chess.Checkmate {
			isCheckmate = true
			isInCheck = true // checkmate implies check
		} else {
			// For ongoing games, detect check using a more reliable method
			// The chess library doesn't have a direct InCheck() method, but we can detect it
			// by checking if the king is under attack

			currentPlayer := position.Turn()

			// Find the king of the current player
			var kingSquare chess.Square = chess.NoSquare
			for sq := chess.A1; sq <= chess.H8; sq++ {
				piece := position.Board().Piece(sq)
				if piece != chess.NoPiece {
					pieceType := piece.Type()
					pieceColor := piece.Color()
					if pieceType == chess.King && pieceColor == currentPlayer {
						kingSquare = sq
						break
					}
				}
			}

			// If we found the king, check if it's under attack
			if kingSquare != chess.NoSquare {
				// Get all valid moves for the current position
				validMoves := game.Game.ValidMoves()

				// A simple but effective heuristic: if the king has to move or if there are very few legal moves
				// and some of them involve the king, it's likely in check
				kingMoves := 0
				totalMoves := len(validMoves)

				for _, move := range validMoves {
					if move.S1() == kingSquare {
						kingMoves++
					}
				}

				// Check heuristics (more conservative approach):
				// 1. If there are very few total moves and multiple king moves, likely in check
				// 2. If king has many escape moves relative to total moves, likely in check
				if totalMoves > 0 {
					kingMoveRatio := float64(kingMoves) / float64(totalMoves)

					// More conservative check detection:
					// - If more than 40% of moves are king moves AND there are multiple king moves, likely in check
					// - OR if there are very few total moves (3 or less) and at least 2 king moves
					if (kingMoveRatio > 0.4 && kingMoves >= 2) || (totalMoves <= 3 && kingMoves >= 2) {
						isInCheck = true
						log.Printf("Check detected for %s: totalMoves=%d, kingMoves=%d, ratio=%.2f", currentPlayer.String(), totalMoves, kingMoves, kingMoveRatio)
					} else {
						log.Printf("No check for %s: totalMoves=%d, kingMoves=%d, ratio=%.2f", currentPlayer.String(), totalMoves, kingMoves, kingMoveRatio)
					}
				}
			}
		}
	}

	// Collect white team player details
	whiteTeamPlayers := make([]map[string]any, 0)
	for walletAddress := range game.WhitePlayers {
		votes := game.PlayerTotalVotes[walletAddress]
		spent := float64(votes) * 0.01 // Each vote costs 0.01 USDC
		whiteTeamPlayers = append(whiteTeamPlayers, map[string]any{
			"walletAddress": walletAddress,
			"totalVotes":    votes,
			"totalSpent":    spent,
		})
		log.Printf("Collecting white player: %s - %d votes, %.3f USDC", walletAddress, votes, spent)
	}
	log.Printf("Total white team players collected: %d", len(whiteTeamPlayers))

	// Collect black team player details
	blackTeamPlayers := make([]map[string]any, 0)
	for walletAddress := range game.BlackPlayers {
		votes := game.PlayerTotalVotes[walletAddress]
		spent := float64(votes) * 0.01 // Each vote costs 0.01 USDC
		blackTeamPlayers = append(blackTeamPlayers, map[string]any{
			"walletAddress": walletAddress,
			"totalVotes":    votes,
			"totalSpent":    spent,
		})
		log.Printf("Collecting black player: %s - %d votes, %.3f USDC", walletAddress, votes, spent)
	}
	log.Printf("Total black team players collected: %d", len(blackTeamPlayers))

	return map[string]any{
		"whitePlayers":          len(game.WhitePlayers),
		"blackPlayers":          len(game.BlackPlayers),
		"whiteCurrentTurnVotes": game.WhiteVotesThisTurn,
		"blackCurrentTurnVotes": game.BlackVotesThisTurn,
		"whiteTeamTotalVotes":   game.WhiteTeamTotalVotes,
		"blackTeamTotalVotes":   game.BlackTeamTotalVotes,
		"totalPot":              game.TotalPot,
		"whitePot":              game.WhitePot,
		"blackPot":              game.BlackPot,
		"currentTurn":           currentTurn,
		"timeLeft":              game.TimeLeft,
		"currentMove":           currentMove,
		"playerVotedThisRound":  game.PlayerVotedThisRound,
		"playerTotalVotes":      game.PlayerTotalVotes,
		"board":                 board,
		"whiteTeamPlayers":      whiteTeamPlayers,
		"blackTeamPlayers":      blackTeamPlayers,
		"isInCheck":             isInCheck,
		"isCheckmate":           isCheckmate,
	}
}

// HasPlayerVoted checks if a specific player has voted in the current round
func (m *Manager) HasPlayerVoted(gameID, walletAddress string) bool {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return false
	}

	game.mu.RLock()
	defer game.mu.RUnlock()

	return game.PlayerVotedThisRound[walletAddress]
}

// GetPlayerTotalVotes returns the total number of votes a player has made
func (m *Manager) GetPlayerTotalVotes(gameID, walletAddress string) int {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return 0
	}

	game.mu.RLock()
	defer game.mu.RUnlock()

	return game.PlayerTotalVotes[walletAddress]
}

// GetPlayerTeam returns the team a player is on
func (m *Manager) GetPlayerTeam(gameID, walletAddress string) string {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return ""
	}

	game.mu.RLock()
	defer game.mu.RUnlock()

	if game.WhitePlayers[walletAddress] {
		return "white"
	}
	if game.BlackPlayers[walletAddress] {
		return "black"
	}
	return ""
}

// GetGameCreatedAt returns the creation timestamp of a game
func (m *Manager) GetGameCreatedAt(gameID string) int64 {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return 0
	}

	game.mu.RLock()
	defer game.mu.RUnlock()

	return game.CreatedAt
}

// applyMoveToBoard applies a move to the chess board
func (m *Manager) applyMoveToBoard(game *GameState, move string) {
	if len(move) < 2 {
		log.Printf("Invalid move format: %s", move)
		return
	}

	log.Printf("Applying move %s to game %s", move, game.ID)

	// Handle coordinate notation (e.g., "e2e4")
	if len(move) == 4 {
		from := parseSquare(move[:2])
		to := parseSquare(move[2:])

		if from == chess.NoSquare || to == chess.NoSquare {
			log.Printf("Invalid square notation in move: %s", move)
			return
		}

		validMoves := game.Game.ValidMoves()
		for _, validMove := range validMoves {
			if validMove.S1() == from && validMove.S2() == to {
				if err := game.Game.Move(&validMove, nil); err != nil {
					log.Printf("Error applying move %s: %v", move, err)
					return
				}
				log.Printf("Successfully applied coordinate move %s", move)
				return
			}
		}
		log.Printf("No valid move found for %s", move)
		return
	}

	// Handle algebraic notation (e.g., "e4", "Nf3")
	if err := game.Game.PushNotationMove(move, chess.AlgebraicNotation{}, nil); err != nil {
		log.Printf("Error applying algebraic move %s: %v", move, err)
		return
	}
	log.Printf("Successfully applied algebraic move %s", move)
}

// checkGameEnd checks if the game has ended
func (m *Manager) checkGameEnd(game *GameState) bool {
	outcome := game.Game.Outcome()
	return outcome != chess.NoOutcome
}

// GetAllGames returns a list of all game IDs
func (m *Manager) GetAllGames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	gameIDs := make([]string, 0, len(m.games))
	for gameID := range m.games {
		gameIDs = append(gameIDs, gameID)
	}
	return gameIDs
}

// distributeRewards distributes rewards to winning players using multicall approach
func (m *Manager) distributeRewards(gameID, winner string, gameStats map[string]any) {
	if winner == "draw" {
		log.Printf("Game %s ended in draw, no rewards to distribute", gameID)
		return
	}

	// Step 1: Gather all rewards from participating vaults to Base Sepolia
	err := m.gatherRewards(gameID, gameStats)
	if err != nil {
		log.Printf("Warning: Failed to gather rewards for game %s: %v", gameID, err)
		return
	}

	// Step 2: Calculate and distribute rewards from total pot
	m.distributeRewardsFromTotalPot(gameID, winner, gameStats)
}

// gatherRewards sends all vault rewards to the central Base Sepolia vault
func (m *Manager) gatherRewards(gameID string, gameStats map[string]any) error {
	// Get all players from both teams to determine which chains were involved
	var allPlayers []map[string]any

	if whiteTeamPlayers, ok := gameStats["whiteTeamPlayers"].([]map[string]any); ok {
		allPlayers = append(allPlayers, whiteTeamPlayers...)
	}
	if blackTeamPlayers, ok := gameStats["blackTeamPlayers"].([]map[string]any); ok {
		allPlayers = append(allPlayers, blackTeamPlayers...)
	}

	if len(allPlayers) == 0 {
		return fmt.Errorf("no players found in game stats")
	}

	// Get unique chain IDs from all players
	involvedChains := make(map[uint64]bool)
	for _, player := range allPlayers {
		if walletAddress, ok := player["walletAddress"].(string); ok {
			playerChainID := m.GetPlayerChainID(walletAddress)
			if playerChainID != 0 {
				involvedChains[uint64(playerChainID)] = true
			}
		}
	}

	// Base Sepolia chain ID (destination for all rewards)
	baseSepoliaChainID := uint64(84532)

	// Remove Base Sepolia from involved chains since we're gathering TO it
	delete(involvedChains, baseSepoliaChainID)

	if len(involvedChains) == 0 {
		log.Printf("No external chains involved in game %s, all rewards already on Base Sepolia", gameID)
		return nil
	}

	// Get blockchain game ID
	m.mu.RLock()
	game, exists := m.games[gameID]
	gameIDUint := uint64(0)
	if exists {
		gameIDUint = game.BlockchainGameID
	}
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	// Get total pot to transfer
	totalPot := 0.0
	if tp, ok := gameStats["totalPot"].(float64); ok {
		totalPot = tp
	}

	if totalPot <= 0 {
		return fmt.Errorf("no pot to gather for game %s", gameID)
	}

	// Calculate amount per chain (distribute the gathering equally)
	totalPotWei := new(big.Int).SetInt64(int64(totalPot * 1000000)) // Convert to USDC wei
	amountPerChain := new(big.Int).Div(totalPotWei, big.NewInt(int64(len(involvedChains))))

	log.Printf("Gathering %s USDC total pot from %d chains to Base Sepolia (%s per chain)",
		totalPotWei.String(), len(involvedChains), amountPerChain.String())

	// Get Base Sepolia vault address as recipient
	baseVaultAddress := client.GetVaultAddress(baseSepoliaChainID)
	if baseVaultAddress == "" {
		return fmt.Errorf("no Base Sepolia vault address configured")
	}
	recipient := common.HexToAddress(baseVaultAddress)

	// Transfer from each involved chain to Base Sepolia
	for chainID := range involvedChains {
		vault, err := m.vaultManager.GetVault(chainID)
		if err != nil {
			log.Printf("Warning: Failed to get vault for chain %d: %v", chainID, err)
			continue
		}

		// Transfer rewards to Base Sepolia vault
		useFastTransfer := false // Use standard transfer for lower fees
		maxFee := big.NewInt(0)  // Let the contract determine the fee

		err = vault.TransferRewards(gameIDUint, amountPerChain, baseSepoliaChainID, recipient, useFastTransfer, maxFee)
		if err != nil {
			log.Printf("Warning: Failed to gather rewards from chain %d: %v", chainID, err)
			continue
		}

		log.Printf("Successfully gathered %s USDC from chain %d to Base Sepolia vault",
			amountPerChain.String(), chainID)
	}

	return nil
}

// distributeRewardsFromTotalPot distributes rewards from the total pot using multicall
func (m *Manager) distributeRewardsFromTotalPot(gameID, winner string, gameStats map[string]any) {
	// Get the winning team players
	var winningPlayers []map[string]any

	switch winner {
	case "white":
		if whiteTeamPlayers, ok := gameStats["whiteTeamPlayers"].([]map[string]any); ok {
			winningPlayers = whiteTeamPlayers
		}
	case "black":
		if blackTeamPlayers, ok := gameStats["blackTeamPlayers"].([]map[string]any); ok {
			winningPlayers = blackTeamPlayers
		}
	}

	if len(winningPlayers) == 0 {
		log.Printf("No winning players for game %s", gameID)
		return
	}

	// Get total pot (not just losing team pot)
	totalPot := 0.0
	if tp, ok := gameStats["totalPot"].(float64); ok {
		totalPot = tp
	}

	if totalPot <= 0 {
		log.Printf("No total pot for game %s", gameID)
		return
	}

	// Calculate total votes from winning team
	totalWinningVotes := 0
	for _, player := range winningPlayers {
		if votes, ok := player["totalVotes"].(int); ok {
			totalWinningVotes += votes
		}
	}

	if totalWinningVotes == 0 {
		log.Printf("No votes from winning team for game %s", gameID)
		return
	}

	// Convert total pot to USDC wei (6 decimal places)
	totalPotWei := new(big.Int).SetInt64(int64(totalPot * 1000000)) // Convert to USDC wei

	log.Printf("Distributing %s USDC from total pot to %d winning players based on %d total votes",
		totalPotWei.String(), len(winningPlayers), totalWinningVotes)

	// Prepare multicall data for all reward transfers
	var rewardTransfers []RewardTransfer

	// Calculate each player's share
	for _, player := range winningPlayers {
		walletAddress, ok := player["walletAddress"].(string)
		if !ok {
			continue
		}

		playerVotes, ok := player["totalVotes"].(int)
		if !ok || playerVotes <= 0 {
			continue
		}

		// Calculate player's proportional share of the total pot
		playerShare := new(big.Int).Mul(totalPotWei, big.NewInt(int64(playerVotes)))
		playerShare.Div(playerShare, big.NewInt(int64(totalWinningVotes)))

		if playerShare.Cmp(big.NewInt(0)) <= 0 {
			continue
		}

		// Get player's chain ID for destination
		playerChainID := m.GetPlayerChainID(walletAddress)
		if playerChainID == 0 {
			log.Printf("Warning: No chain ID found for player %s, skipping reward", walletAddress)
			continue
		}

		rewardTransfers = append(rewardTransfers, RewardTransfer{
			Recipient:        common.HexToAddress(walletAddress),
			Amount:           playerShare,
			DestinationChain: uint64(playerChainID),
		})

		log.Printf("Prepared reward transfer: %s USDC to %s on chain %d",
			playerShare.String(), walletAddress, playerChainID)
	}

	if len(rewardTransfers) == 0 {
		log.Printf("No valid reward transfers for game %s", gameID)
		return
	}

	// Execute multicall reward distribution
	err := m.executeMulticallRewards(gameID, rewardTransfers)
	if err != nil {
		log.Printf("Warning: Failed to execute multicall rewards for game %s: %v", gameID, err)
	}
}

// RewardTransfer represents a single reward transfer
type RewardTransfer struct {
	Recipient        common.Address
	Amount           *big.Int
	DestinationChain uint64
}

// executeMulticallRewards executes multiple reward transfers using multicall
func (m *Manager) executeMulticallRewards(gameID string, transfers []RewardTransfer) error {
	// Get Base Sepolia vault (where all rewards are now gathered)
	baseSepoliaChainID := uint64(84532)
	baseVault, err := m.vaultManager.GetVault(baseSepoliaChainID)
	if err != nil {
		return fmt.Errorf("failed to get Base Sepolia vault: %w", err)
	}

	// Get blockchain game ID
	m.mu.RLock()
	game, exists := m.games[gameID]
	gameIDUint := uint64(0)
	if exists {
		gameIDUint = game.BlockchainGameID
	}
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	log.Printf("Executing multicall reward distribution for %d transfers", len(transfers))

	for i, transfer := range transfers {
		log.Printf("Executing transfer %d/%d: %s USDC to %s on chain %d",
			i+1, len(transfers), transfer.Amount.String(), transfer.Recipient.Hex(), transfer.DestinationChain)

		useFastTransfer := false // Use standard transfer for lower fees
		maxFee := big.NewInt(0)  // Let the contract determine the fee

		err := baseVault.TransferRewards(gameIDUint, transfer.Amount, transfer.DestinationChain, transfer.Recipient, useFastTransfer, maxFee)
		if err != nil {
			log.Printf("Warning: Failed to transfer reward %d: %v", i+1, err)
			continue
		}

		log.Printf("Successfully transferred reward %d: %s USDC to %s on chain %d",
			i+1, transfer.Amount.String(), transfer.Recipient.Hex(), transfer.DestinationChain)
	}

	log.Printf("Completed multicall reward distribution for game %s", gameID)
	return nil
}
