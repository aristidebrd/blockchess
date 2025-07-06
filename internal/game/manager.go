package game

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/corentings/chess/v2"
	"github.com/ethereum/go-ethereum/common"
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

	mu sync.RWMutex
}

// Manager handles game logic and orchestrates services
type Manager struct {
	games              map[string]*GameState
	mu                 sync.RWMutex
	moveResultCallback func(gameID, move string)
	gameEndCallback    func(gameID, winner, reason string, gameStats map[string]any)

	// Services (dependency injection)
	blockchainService BlockchainService
	configService     ConfigService
	endGameProcessor  *EndGameProcessor
}

// NewManager creates a new game manager with default services
func NewManager() *Manager {
	configService := NewConfigService()

	// Try to initialize blockchain service
	var blockchainService BlockchainService
	config, _ := configService.LoadBlockchainConfig()
	blockchainService, _ = NewEthereumBlockchainService(config)

	return &Manager{
		games:             make(map[string]*GameState),
		blockchainService: blockchainService,
		configService:     configService,
	}
}

// NewManagerWithServices creates a new game manager with injected services
func NewManagerWithServices(blockchainService BlockchainService, configService ConfigService) *Manager {
	return &Manager{
		games:             make(map[string]*GameState),
		blockchainService: blockchainService,
		configService:     configService,
		endGameProcessor:  NewEndGameProcessor(blockchainService),
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

// GetOrCreateGame gets an existing game or creates a new one
func (m *Manager) GetOrCreateGame() *GameState {
	m.mu.Lock()
	defer m.mu.Unlock()

	var gameResult *GameCreationResult

	// Create game on blockchain if service is available
	if m.blockchainService != nil && m.blockchainService.IsConnected() {
		gameResult, err := m.createGameOnBlockchain(StakeAmount)
		if err != nil {
			log.Printf("Error creating game on blockchain: %v", err)
		} else {
			log.Printf("Game %s created on-chain! Transaction: %s", gameResult.GameID.String(), gameResult.TransactionHash.Hex())
		}
	}

	game := &GameState{
		ID:                   gameResult.GameID.String(),
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

// createGameOnBlockchain creates a game on the blockchain asynchronously
func (m *Manager) createGameOnBlockchain(stakeAmount float64) (*GameCreationResult, error) {
	stakeAmountBigInt := big.NewInt(int64(stakeAmount * 1e6))
	if gameResult, err := m.blockchainService.CreateGame(stakeAmountBigInt); err != nil {
		log.Printf("Error creating game on blockchain: %v", err)
		return nil, err
	} else {
		return gameResult, nil
	}
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

	// Process game end using the EndGameProcessor
	if m.endGameProcessor != nil {
		go func() {
			if err := m.endGameProcessor.ProcessGameEnd(gameID, winner, game); err != nil {
				log.Printf("Error processing game end for %s: %v", gameID, err)
			}
		}()
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

	// Record vote
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

	// Record move on blockchain if available
	if m.blockchainService != nil && m.blockchainService.IsConnected() {
		go func() {
			player := common.HexToAddress(walletAddress)

			// Record the move
			var teamCode uint8 = 0 // 0 for white, 1 for black
			if team == "black" {
				teamCode = 1
			}
			if err := m.blockchainService.RecordMove(gameID, player, chainId, teamCode); err != nil {
				log.Printf("Error recording move on blockchain: %v", err)
			}

			// Stake USDC from player's approved allowance
			stakeAmount := big.NewInt(10000) // 0.01 USDC (6 decimals)
			if err := m.blockchainService.StakeOnVote(gameID, player, stakeAmount); err != nil {
				log.Printf("Error staking USDC from player %s: %v", walletAddress, err)
			} else {
				log.Printf("Successfully staked %s USDC from player %s in game %s", stakeAmount.String(), walletAddress, gameID)
			}
		}()
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

// GetBlockchainService returns the blockchain service (for testing)
func (m *Manager) GetBlockchainService() BlockchainService {
	return m.blockchainService
}

// SetBlockchainService sets the blockchain service (for testing)
func (m *Manager) SetBlockchainService(service BlockchainService) {
	m.blockchainService = service
}

// GeneratePermit2SignatureData generates Permit2 typed data for signing
func (m *Manager) GeneratePermit2SignatureData(walletAddress string, chainId uint32) (interface{}, error) {
	if m.blockchainService == nil {
		return nil, fmt.Errorf("blockchain service not available")
	}

	playerAddr := common.HexToAddress(walletAddress)

	// Cast to EthereumBlockchainService to access Permit2 method
	ethService, ok := m.blockchainService.(*EthereumBlockchainService)
	if !ok {
		return nil, fmt.Errorf("blockchain service does not support Permit2")
	}

	return ethService.GeneratePermit2SignatureData(playerAddr, chainId)
}

// StorePermit2Signature stores a signed Permit2 signature for a player
func (m *Manager) StorePermit2Signature(walletAddress, signature string) error {
	if m.blockchainService == nil {
		return fmt.Errorf("blockchain service not available")
	}

	playerAddr := common.HexToAddress(walletAddress)

	// Cast to EthereumBlockchainService to access Permit2 method
	ethService, ok := m.blockchainService.(*EthereumBlockchainService)
	if !ok {
		return fmt.Errorf("blockchain service does not support Permit2")
	}

	return ethService.StorePermit2Signature(playerAddr, signature)
}

// ExecutePermit2 executes a Permit2 allowance using the stored signature
func (m *Manager) ExecutePermit2(walletAddress string, chainId uint32) error {
	if m.blockchainService == nil {
		return fmt.Errorf("blockchain service not available")
	}

	playerAddr := common.HexToAddress(walletAddress)

	// Cast to EthereumBlockchainService to access Permit2 method
	ethService, ok := m.blockchainService.(*EthereumBlockchainService)
	if !ok {
		return fmt.Errorf("blockchain service does not support Permit2")
	}

	return ethService.ExecutePermit2(playerAddr, chainId)
}
