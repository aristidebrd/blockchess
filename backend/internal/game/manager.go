package game

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/notnil/chess"
)

type Vote struct {
	Move  string
	Count int
}

type GameState struct {
	ID          string
	Votes       map[string]int // move -> vote count
	TimeLeft    int            // seconds
	Game        *chess.Game    // Chess game state from notnil/chess library
	Players     []string       // connected player IDs
	CurrentMove int            // Current move number

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

	mu sync.RWMutex
}

type Manager struct {
	games map[string]*GameState
	mu    sync.RWMutex
	// Callback for broadcasting move results
	moveResultCallback func(gameID, move string)
	// Callback for broadcasting game end
	gameEndCallback func(gameID, winner, reason string, gameStats map[string]any)

	// Contract integration
	ethClient interface {
		PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
		NetworkID(ctx context.Context) (*big.Int, error)
		SuggestGasPrice(ctx context.Context) (*big.Int, error)
	}
	gameContractAddr common.Address
	privateKey       *ecdsa.PrivateKey
	auth             *bind.TransactOpts
	defaultStakeWei  *big.Int
}

// ContractConfig holds the configuration for contract integration
type ContractConfig struct {
	RPCUrl              string
	PrivateKey          string
	GameContractAddress string
	DefaultStakeETH     string // Default stake amount in ETH (e.g., "0.01")
}

func NewManager() *Manager {
	return NewManagerWithContracts(nil)
}

func NewManagerWithContracts(config *ContractConfig) *Manager {
	m := &Manager{
		games: make(map[string]*GameState),
	}

	// Initialize contract integration if config is provided
	if config != nil {
		if err := m.initializeContracts(config); err != nil {
			log.Printf("Warning: Failed to initialize contracts: %v", err)
			log.Printf("Continuing without contract integration...")
		}
	}

	return m
}

// LoadContractConfigFromEnv loads contract configuration from environment variables
func LoadContractConfigFromEnv() (*ContractConfig, error) {
	// Load environment variables
	if err := godotenv.Load(".env.local"); err != nil {
		log.Printf("Warning: .env.local file not found, using environment variables")
	}

	config := &ContractConfig{
		RPCUrl:              getEnv("RPC_URL", "http://127.0.0.1:8545"),
		PrivateKey:          getEnv("PRIVATE_KEY", ""),
		GameContractAddress: getEnv("GAME_CONTRACT_ADDRESS", ""),
		DefaultStakeETH:     getEnv("DEFAULT_STAKE_ETH", "0.01"),
	}

	if config.PrivateKey == "" {
		return nil, fmt.Errorf("PRIVATE_KEY environment variable not set")
	}
	if config.GameContractAddress == "" {
		return nil, fmt.Errorf("GAME_CONTRACT_ADDRESS environment variable not set")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (m *Manager) initializeContracts(config *ContractConfig) error {
	// Connect to Ethereum client
	client, err := ethclient.Dial(config.RPCUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	// Get chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %v", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %v", err)
	}

	// Set gas parameters
	auth.GasLimit = uint64(3000000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}
	auth.GasPrice = gasPrice

	// Parse default stake amount
	stakeFloat, err := strconv.ParseFloat(config.DefaultStakeETH, 64)
	if err != nil {
		return fmt.Errorf("failed to parse default stake amount: %v", err)
	}
	// Convert ETH to Wei (1 ETH = 10^18 Wei)
	stakeWei := new(big.Float).Mul(big.NewFloat(stakeFloat), big.NewFloat(1e18))
	defaultStakeWei, _ := stakeWei.Int(nil)

	// Store contract configuration
	m.ethClient = client
	m.gameContractAddr = common.HexToAddress(config.GameContractAddress)
	m.privateKey = privateKey
	m.auth = auth
	m.defaultStakeWei = defaultStakeWei

	log.Printf("Contract integration initialized successfully")
	log.Printf("Game Contract: %s", m.gameContractAddr.Hex())
	log.Printf("Default Stake: %s Wei (%s ETH)", m.defaultStakeWei.String(), config.DefaultStakeETH)

	return nil
}

// SetContractClients sets the contract clients for testing (used by tests)
func (m *Manager) SetContractClients(client interface{}, gameContract interface{}, vaultContract interface{}, auth *bind.TransactOpts) {
	// Handle both real ethclient and simulated backend
	if ethClient, ok := client.(interface {
		PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
		NetworkID(ctx context.Context) (*big.Int, error)
		SuggestGasPrice(ctx context.Context) (*big.Int, error)
	}); ok {
		m.ethClient = ethClient
	}
	m.auth = auth
	m.defaultStakeWei = big.NewInt(10000000000000000) // 0.01 ETH for testing
}

// Set the callback for broadcasting move results
func (m *Manager) SetMoveResultCallback(callback func(gameID, move string)) {
	m.moveResultCallback = callback
}

// Set the callback for broadcasting game end
func (m *Manager) SetGameEndCallback(callback func(gameID, winner, reason string, gameStats map[string]any)) {
	m.gameEndCallback = callback
}

func (m *Manager) GetOrCreateGame(gameID string) *GameState {
	m.mu.Lock()
	defer m.mu.Unlock()

	if game, exists := m.games[gameID]; exists {
		return game
	}

	// Create new game locally
	game := &GameState{
		ID:                   gameID,
		Votes:                make(map[string]int),
		TimeLeft:             10, // 10 seconds per turn
		Game:                 chess.NewGame(),
		Players:              make([]string, 0),
		CurrentMove:          1,
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

	m.games[gameID] = game

	// Create game on-chain if contract integration is available
	if m.ethClient != nil && m.gameContractAddr != (common.Address{}) {
		go m.createGameOnChain(gameID)
	}

	// Start game timer
	go m.runGameTimer(game)

	return game
}

func (m *Manager) runGameTimer(game *GameState) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		game.mu.Lock()
		game.TimeLeft--

		if game.TimeLeft <= 0 {
			// Time's up - execute the move with most votes
			bestMove := m.getBestMove(game)

			// Apply the move to the board if it's not "skip"
			moveApplied := false
			if bestMove != "skip" {
				m.applyMoveToBoard(game, bestMove)
				moveApplied = true
			}

			// Check if the game ended after this move, but don't stop the timer yet
			gameEnded := moveApplied && m.checkGameEnd(game)

			// Reset for next turn immediately to prevent race conditions
			game.Votes = make(map[string]int)
			game.TimeLeft = 10
			game.CurrentMove++

			// Reset round vote tracking
			game.WhiteVotesThisTurn = 0
			game.BlackVotesThisTurn = 0
			game.PlayerVotedThisRound = make(map[string]bool)

			game.mu.Unlock()

			// Broadcast move result after reset - this allows players to see the final move
			m.BroadcastMoveResult(game.ID, bestMove)

			// If the game ended, handle it after broadcasting the move result
			if gameEnded {
				log.Printf("Game %s ended after move %s! Getting final stats...", game.ID, bestMove)

				// Get final stats (need to re-acquire lock briefly)
				game.mu.RLock()
				gameStats := m.getGameStatsUnsafe(game, true)
				game.mu.RUnlock()

				log.Printf("Broadcasting game end for %s", game.ID)
				// Broadcast game end after players have seen the final move
				m.broadcastGameEnd(game.ID, gameStats)
				log.Printf("Game %s timer stopped", game.ID)
				return // Game ended, stop the timer
			}
		} else {
			// Timer is counting down, just unlock and continue
			game.mu.Unlock()
		}
	}
}

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
		// No votes - return a default move or skip
		return "skip"
	}

	return bestMove
}

func (m *Manager) VoteForMove(gameID, walletAddress, move string, team string) error {
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

	// Validate the move using the chess library
	if !m.isValidMove(game, move) {
		return fmt.Errorf("invalid move: %s", move)
	}

	// Check if move already exists, if not initialize it
	if _, exists := game.Votes[move]; !exists {
		game.Votes[move] = 0
	}

	// Record vote
	game.Votes[move]++
	game.PlayerVotedThisRound[walletAddress] = true

	// Increment total vote count for this player
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

// Check if a move is valid in the current game state
func (m *Manager) isValidMove(game *GameState, moveStr string) bool {
	if len(moveStr) < 2 {
		return false
	}

	// If the move is in coordinate notation (e.g., "e2e4"), find the matching Move
	if len(moveStr) == 4 {
		from := parseSquare(moveStr[:2])
		to := parseSquare(moveStr[2:])

		if from == chess.NoSquare || to == chess.NoSquare {
			return false
		}

		// Get all valid moves and find one that matches our from-to squares
		validMoves := game.Game.ValidMoves()
		for _, move := range validMoves {
			if move.S1() == from && move.S2() == to {
				return true
			}
		}

		return false
	}

	// If it's already in algebraic notation, test it directly
	gameCopy := game.Game.Clone()
	err := gameCopy.MoveStr(moveStr)
	return err == nil
}

// Helper function to parse square notation like "e2" into chess.Square
func parseSquare(square string) chess.Square {
	if len(square) != 2 {
		return chess.NoSquare
	}

	file := square[0]
	rank := square[1]

	// Convert file ('a'-'h') to 0-7
	if file < 'a' || file > 'h' {
		return chess.NoSquare
	}
	fileIndex := int(file - 'a')

	// Convert rank ('1'-'8') to 0-7
	if rank < '1' || rank > '8' {
		return chess.NoSquare
	}
	rankIndex := int(rank - '1')

	// Calculate square index: rank * 8 + file
	squareIndex := rankIndex*8 + fileIndex

	return chess.Square(squareIndex)
}

// Get all valid moves for the current position
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
		moveStrings[i] = move.String()
	}

	return moveStrings
}

func (m *Manager) GetVotes(gameID string) map[string]int {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return nil
	}

	game.mu.RLock()
	defer game.mu.RUnlock()

	// Return copy of votes
	votes := make(map[string]int)
	for k, v := range game.Votes {
		votes[k] = v
	}

	return votes
}

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

func (m *Manager) BroadcastMoveResult(gameID, move string) {
	// Call the callback if it's set (will be handled by the hub)
	if m.moveResultCallback != nil {
		m.moveResultCallback(gameID, move)
	}
}

// Add player to team with wallet address validation
func (m *Manager) AddPlayerToTeam(gameID, walletAddress, team string) error {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game not found: %s", gameID)
	}

	// Validate wallet address format (basic check)
	if walletAddress == "" {
		return fmt.Errorf("wallet address cannot be empty")
	}

	// Simple validation that it looks like an Ethereum address
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

// Get game statistics
func (m *Manager) GetGameStats(gameID string) map[string]any {
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return nil
	}

	game.mu.RLock()
	defer game.mu.RUnlock()

	// Convert chess board to string array
	board := make([][]string, 8)
	pieceMap := map[string]string{
		"♜": "r", "♞": "n", "♝": "b", "♛": "q", "♚": "k", "♟": "p", // Black pieces
		"♖": "R", "♘": "N", "♗": "B", "♕": "Q", "♔": "K", "♙": "P", // White pieces
	}

	for i := range 8 {
		board[i] = make([]string, 8)
		for j := range 8 {
			// Chess library uses files a-h (0-7) and ranks 1-8 (0-7)
			// Square mapping: rank * 8 + file
			// Our board: [row][col] where row 0 = rank 8, row 7 = rank 1
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
		"currentMove":           game.CurrentMove,
		"playerVotedThisRound":  game.PlayerVotedThisRound,
		"playerTotalVotes":      game.PlayerTotalVotes,
		"board":                 board,
	}
}

// Check if a specific player has voted in the current round
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

// Get total number of votes a specific player has made throughout the game
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

// Get player's team in a specific game
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

// Apply a move to the board (move format: "e2e4" or "e4")
func (m *Manager) applyMoveToBoard(game *GameState, move string) {
	if len(move) < 2 {
		log.Printf("Invalid move format: %s", move)
		return // Invalid move format
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

		// Find the matching Move object
		validMoves := game.Game.ValidMoves()
		for _, validMove := range validMoves {
			if validMove.S1() == from && validMove.S2() == to {
				if err := game.Game.Move(validMove); err != nil {
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
	if err := game.Game.MoveStr(move); err != nil {
		log.Printf("Error applying algebraic move %s: %v", move, err)
		return
	}
	log.Printf("Successfully applied algebraic move %s", move)
}

// Check if the game has ended (checkmate, stalemate, draw) - returns true if game ended
func (m *Manager) checkGameEnd(game *GameState) bool {
	outcome := game.Game.Outcome()
	return outcome != chess.NoOutcome
}

// Get game statistics without locking (caller must hold the lock)
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
			// Chess library uses files a-h (0-7) and ranks 1-8 (0-7)
			// Square mapping: rank * 8 + file
			// Our board: [row][col] where row 0 = rank 8, row 7 = rank 1
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
	}
}

// Broadcast game end (called after releasing the game lock)
func (m *Manager) broadcastGameEnd(gameID string, gameStats map[string]any) {
	// Re-acquire the game to get the final state for winner/reason detection
	m.mu.RLock()
	game, exists := m.games[gameID]
	m.mu.RUnlock()

	if !exists {
		return
	}

	game.mu.RLock()
	outcome := game.Game.Outcome()
	method := game.Game.Method()
	game.mu.RUnlock()

	if outcome == chess.NoOutcome {
		return // No game end needed
	}

	// Determine winner and reason
	var winner, reason string

	switch outcome {
	case chess.WhiteWon:
		winner = "white"
		if method == chess.Checkmate {
			reason = "checkmate"
		} else {
			reason = "resignation" // or other methods
		}
	case chess.BlackWon:
		winner = "black"
		if method == chess.Checkmate {
			reason = "checkmate"
		} else {
			reason = "resignation" // or other methods
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

	// Broadcast game end if callback is set
	if m.gameEndCallback != nil {
		m.gameEndCallback(gameID, winner, reason, gameStats)
	}
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

// createGameOnChain creates the game on the smart contract
func (m *Manager) createGameOnChain(gameID string) {
	log.Printf("Creating game %s on-chain...", gameID)

	// Convert gameID string to uint256
	gameIDInt, err := strconv.ParseUint(gameID, 10, 64)
	if err != nil {
		log.Printf("Error: gameID %s is not a valid number for on-chain creation", gameID)
		return
	}

	// Get fresh nonce
	nonce, err := m.ethClient.PendingNonceAt(context.Background(), m.auth.From)
	if err != nil {
		log.Printf("Error getting nonce for game creation: %v", err)
		return
	}
	m.auth.Nonce = big.NewInt(int64(nonce))

	// Call the contract's createGame function
	// This is a simplified call - you'll need to import and use the generated bindings
	// For now, using raw transaction call

	// Create the transaction data
	// Method ID for createGame(uint256,uint256) is the first 4 bytes of keccak256("createGame(uint256,uint256)")
	// This would be: 0x60104cef

	// For a proper implementation, you should use the generated Go bindings:
	/*
		gameContract, err := gamecontract.NewGamecontract(m.gameContractAddr, m.ethClient)
		if err != nil {
			log.Printf("Error creating game contract instance: %v", err)
			return
		}

		tx, err := gameContract.CreateGame(m.auth, big.NewInt(int64(gameIDInt)), m.defaultStakeWei)
		if err != nil {
			log.Printf("Error creating game on-chain: %v", err)
			return
		}

		log.Printf("Game %s created on-chain! Transaction: %s", gameID, tx.Hash().Hex())
	*/

	// For now, just log the attempt
	log.Printf("Would create game %d on-chain with stake %s Wei", gameIDInt, m.defaultStakeWei.String())
	log.Printf("Contract address: %s", m.gameContractAddr.Hex())
}
