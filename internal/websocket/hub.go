package websocket

import (
	"blockchess/internal/game"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Message types
const (
	TypeJoinGame                 = "join_game"
	TypeVoteMove                 = "vote_move"
	TypeVoteUpdate               = "vote_update"
	TypeMoveResult               = "move_result"
	TypeTimerTick                = "timer_tick"
	TypeJoinMatchmaking          = "join_matchmaking"
	TypeLeaveMatchmaking         = "leave_matchmaking"
	TypeMatchFound               = "match_found"
	TypeJoinTeam                 = "join_team"
	TypeWatchGame                = "watch_game"
	TypeGamesList                = "games_list"
	TypeRequestGamesList         = "request_games_list"
	TypeGamesListUpdate          = "games_list_update"
	TypeNumberOfPlayers          = "number_of_players"
	TypeClientConnected          = "client_connected"
	TypeGameEnd                  = "game_end"
	TypeRequestFilteredGamesList = "request_filtered_games_list"
	TypeCheckPlayerStatus        = "check_player_status"
	TypePlayerStatus             = "player_status"
	TypeGetValidMoves            = "get_valid_moves"
	TypeValidMovesResponse       = "valid_moves_response"
	TypeError                    = "error"
	TypeRequestPermit2           = "request_permit2"
	TypePermit2Data              = "permit2_data"
	TypeSubmitPermit2Signature   = "submit_permit2_signature"
)

// GameInfo holds summary information about a single game
type GameInfo struct {
	GameID       string     `json:"gameId"`
	WhitePlayers int        `json:"whitePlayers"`
	BlackPlayers int        `json:"blackPlayers"`
	TimeLeft     int        `json:"timeLeft"`
	CurrentMove  int        `json:"currentMove"`
	TotalPot     float64    `json:"totalPot"`
	WhitePot     float64    `json:"whitePot"`
	BlackPot     float64    `json:"blackPot"`
	Spectators   int        `json:"spectators"`
	CurrentTurn  string     `json:"currentTurn"`
	Status       string     `json:"status"`              // "active" or "ended"
	Winner       string     `json:"winner,omitempty"`    // "white", "black", "draw" (only for ended games)
	EndReason    string     `json:"endReason,omitempty"` // "checkmate", "stalemate", etc. (only for ended games)
	CreatedAt    int64      `json:"createdAt"`           // Unix timestamp when game was created
	EndedAt      *int64     `json:"endedAt,omitempty"`   // Unix timestamp when game ended
	Board        [][]string `json:"board,omitempty"`     // Current board state

	// Player statistics per team (only for ended games)
	WhiteTeamPlayers []PlayerStats `json:"whiteTeamPlayers,omitempty"`
	BlackTeamPlayers []PlayerStats `json:"blackTeamPlayers,omitempty"`
}

type Message struct {
	Type             string         `json:"type"`
	GameID           string         `json:"gameId,omitempty"`
	Move             string         `json:"move,omitempty"`
	Votes            map[string]int `json:"votes,omitempty"`
	SecondsLeft      int            `json:"secondsLeft,omitempty"`
	Players          []string       `json:"players,omitempty"`
	AssignedSide     string         `json:"assignedSide,omitempty"`
	Team             string         `json:"team,omitempty"`
	PlayerID         string         `json:"playerId,omitempty"`
	WalletAddress    string         `json:"walletAddress,omitempty"`
	ClientID         string         `json:"clientId,omitempty"`
	Board            [][]string     `json:"board,omitempty"`
	GamesList        []GameInfo     `json:"gamesList,omitempty"`
	TotalConnections int            `json:"totalConnections,omitempty"`
	Filter           string         `json:"filter,omitempty"`     // "active", "ended", or "" for all
	Error            string         `json:"error,omitempty"`      // Error message
	ValidMoves       []string       `json:"validMoves,omitempty"` // List of valid moves in coordinate notation
	Signature        string         `json:"signature,omitempty"`  // Permit2 signature
	TypedData        interface{}    `json:"typedData,omitempty"`
	ChainId          uint32         `json:"chainId,omitempty"` // EIP-712 typed data

	// Game statistics
	WhitePlayers          int             `json:"whitePlayers,omitempty"`
	BlackPlayers          int             `json:"blackPlayers,omitempty"`
	WhiteCurrentTurnVotes int             `json:"whiteCurrentTurnVotes,omitempty"`
	BlackCurrentTurnVotes int             `json:"blackCurrentTurnVotes,omitempty"`
	WhiteTeamTotalVotes   int             `json:"whiteTeamTotalVotes,omitempty"`
	BlackTeamTotalVotes   int             `json:"blackTeamTotalVotes,omitempty"`
	TotalPot              float64         `json:"totalPot,omitempty"`
	WhitePot              float64         `json:"whitePot,omitempty"`
	BlackPot              float64         `json:"blackPot,omitempty"`
	CurrentTurn           string          `json:"currentTurn,omitempty"`
	CurrentMove           int             `json:"currentMove,omitempty"`
	PlayerVotedThisRound  map[string]bool `json:"playerVotedThisRound,omitempty"`
	PlayerTotalVotes      map[string]int  `json:"playerTotalVotes,omitempty"`

	// Game end information
	Winner        string `json:"winner,omitempty"`        // "white", "black", "draw"
	GameEndReason string `json:"gameEndReason,omitempty"` // "checkmate", "stalemate", "draw"
	PlayerVotes   int    `json:"playerVotes,omitempty"`   // Current player's total votes

	// Player statistics per team
	WhiteTeamPlayers []PlayerStats `json:"whiteTeamPlayers,omitempty"`
	BlackTeamPlayers []PlayerStats `json:"blackTeamPlayers,omitempty"`

	// Check and checkmate status
	IsInCheck   bool `json:"isInCheck,omitempty"`
	IsCheckmate bool `json:"isCheckmate,omitempty"`
}

type PlayerStats struct {
	WalletAddress string  `json:"walletAddress"`
	TotalVotes    int     `json:"totalVotes"`
	TotalSpent    float64 `json:"totalSpent"`
}
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan *ClientMessage

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Game manager
	gameManager *game.Manager

	// Game rooms - gameID -> clients
	gameRooms map[string]map[*Client]bool

	// Ended games - gameID -> GameInfo
	endedGames map[string]*GameInfo

	// Matchmaking queue - wallet addresses with their clients
	matchmakingQueue map[string]*Client

	// Client teams - walletAddress -> team
	clientTeams map[string]string

	// Client wallet addresses - client -> wallet address
	clientWallets map[*Client]string
}

func NewHub(gm *game.Manager) *Hub {
	h := &Hub{
		broadcast:        make(chan *ClientMessage),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		clients:          make(map[*Client]bool),
		gameManager:      gm,
		gameRooms:        make(map[string]map[*Client]bool),
		endedGames:       make(map[string]*GameInfo),
		matchmakingQueue: make(map[string]*Client),
		clientTeams:      make(map[string]string),
		clientWallets:    make(map[*Client]string),
	}

	// Set up the move result callback
	gm.SetMoveResultCallback(h.handleMoveResult)

	// Set up the game end callback
	gm.SetGameEndCallback(h.handleGameEnd)

	// Start periodic updates
	go h.startPeriodicUpdates()

	return h
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			// Send client ID to the newly connected client
			clientMsg := &Message{
				Type:     TypeClientConnected,
				ClientID: client.id,
			}
			if data, err := json.Marshal(clientMsg); err == nil {
				log.Printf("Sending client_connected to %s: %s", client.id, string(data))
				select {
				case client.send <- data:
					log.Printf("✅ Successfully sent client_connected to %s", client.id)
				default:
					log.Printf("❌ Failed to send client_connected to %s - channel full", client.id)
				}
			} else {
				log.Printf("❌ Failed to marshal client_connected for %s: %v", client.id, err)
			}

			h.broadcastToAll(&Message{
				Type:             TypeNumberOfPlayers,
				TotalConnections: h.GetTotalConnections(),
			})
			log.Printf("Client registered: %s (Total connections: %d)", client.id, len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Remember if this client was in any game
				wasInGame := false

				// Remove from matchmaking queue
				h.removeFromMatchmaking(client)

				// Remove from game rooms
				for _, room := range h.gameRooms {
					if _, ok := room[client]; ok {
						delete(room, client)
						wasInGame = true
						// Don't delete empty game rooms - games should persist even with no spectators
						// Only delete rooms when games actually end
					}
				}

				// Remove from client teams using wallet address
				if walletAddress, exists := h.clientWallets[client]; exists {
					if h.clientTeams[walletAddress] != "" {
						wasInGame = true
					}
					delete(h.clientTeams, walletAddress)
				}

				// Remove from client wallets
				delete(h.clientWallets, client)

				log.Printf("Client unregistered: %s (Total connections: %d)", client.id, len(h.clients))

				// Broadcast updated total connections count to all remaining clients
				h.broadcastToAll(&Message{
					Type:             TypeNumberOfPlayers,
					TotalConnections: h.GetTotalConnections(),
				})

				// Broadcast updated games list if client was in a game
				if wasInGame {
					h.broadcastGamesListUpdate()
				}
			}

		case clientMessage := <-h.broadcast:
			// Handle incoming message from client
			h.handleClientMessage(clientMessage)
		}
	}
}

func (h *Hub) handleClientMessage(clientMessage *ClientMessage) {
	var msg Message
	if err := json.Unmarshal(clientMessage.data, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	h.handleMessage(&msg, clientMessage.client)
}

func (h *Hub) handleMessage(msg *Message, client *Client) {
	switch msg.Type {
	case TypeJoinGame:
		log.Printf("Player %s joining game: %s", client.id, msg.GameID)

		// Check if game exists before allowing join
		game := h.gameManager.GetGame(msg.GameID)
		if game == nil {
			log.Printf("Game %s does not exist, cannot join", msg.GameID)
			h.sendErrorToClient(client, "Game does not exist. Games can only be created through matchmaking.")
			return
		}

		// Add client to game room
		h.AddClientToGame(client, msg.GameID)

		// Send initial game state to the joining player
		votes := h.gameManager.GetVotes(msg.GameID)
		stats := h.gameManager.GetGameStats(msg.GameID)

		initialMsg := &Message{
			Type:   TypeVoteUpdate,
			GameID: msg.GameID,
			Votes:  votes,
		}

		h.updateStats(stats, initialMsg)

		// Send initial state to the joining client
		if data, err := json.Marshal(initialMsg); err == nil {
			select {
			case client.send <- data:
			default:
			}
		}

		// Broadcast updated games list to all clients
		h.broadcastGamesListUpdate()

	case TypeVoteMove:
		// Use player ID (wallet address) from message if provided, otherwise fall back to client ID
		walletAddress := msg.PlayerID
		if walletAddress == "" {
			walletAddress = client.id
		}

		// Store the wallet address mapping for this client
		h.clientWallets[client] = walletAddress

		log.Printf("Vote for move %s in game %s from wallet %s", msg.Move, msg.GameID, walletAddress)

		// Get player's team from the game manager (authoritative source)
		team := h.gameManager.GetPlayerTeam(msg.GameID, walletAddress)
		if team == "" {
			log.Printf("Player %s has no team, cannot vote", walletAddress)
			h.sendErrorToClient(client, "You must join a team before voting")
			return
		}

		// Attempt to vote
		if err := h.gameManager.VoteForMove(msg.GameID, walletAddress, msg.Move, team, msg.ChainId); err != nil {
			log.Printf("Vote failed for player %s: %v", walletAddress, err)
			h.sendErrorToClient(client, err.Error())
			return
		}

		// Get updated votes and game stats
		votes := h.gameManager.GetVotes(msg.GameID)
		stats := h.gameManager.GetGameStats(msg.GameID)

		// Create vote update message with all stats
		updateMsg := &Message{
			Type:   TypeVoteUpdate,
			GameID: msg.GameID,
			Votes:  votes,
		}

		h.updateStats(stats, updateMsg)

		h.broadcastToGame(msg.GameID, updateMsg)

		// Broadcast updated games list since pot has changed
		h.broadcastGamesListUpdate()

	case TypeJoinMatchmaking:
		walletAddress := msg.WalletAddress
		if walletAddress == "" {
			log.Printf("No wallet address provided for matchmaking from client %s", client.id)
			h.sendErrorToClient(client, "Wallet address is required for matchmaking")
			return
		}

		// Store the wallet address mapping for this client
		h.clientWallets[client] = walletAddress

		log.Printf("Player %s (wallet: %s) joining matchmaking", client.id, walletAddress)
		h.addToMatchmaking(client, walletAddress)

	case TypeLeaveMatchmaking:
		log.Printf("Player %s leaving matchmaking", client.id)
		h.removeFromMatchmaking(client)

	case TypeJoinTeam:
		// Use player ID (wallet address) from message if provided, otherwise fall back to client ID
		walletAddress := msg.PlayerID
		if walletAddress == "" {
			log.Printf("No wallet address provided for team join from client %s", client.id)
			h.sendErrorToClient(client, "Wallet address is required to join a team")
			return
		}

		// Store the wallet address mapping for this client
		h.clientWallets[client] = walletAddress

		log.Printf("Player %s (wallet: %s) joining %s team in game %s", client.id, walletAddress, msg.Team, msg.GameID)

		// First check if player is already in the game
		existingTeam := h.gameManager.GetPlayerTeam(msg.GameID, walletAddress)

		if existingTeam != "" {
			// Player is already in the game - this is a reconnection
			if existingTeam != msg.Team {
				// Player is trying to join a different team than they're already on
				log.Printf("Player %s is already on %s team, cannot join %s team", walletAddress, existingTeam, msg.Team)
				h.sendErrorToClient(client, fmt.Sprintf("You are already on the %s team. Cannot switch teams.", existingTeam))
				return
			}

			// Player is reconnecting to their existing team
			log.Printf("Player %s reconnecting to existing %s team in game %s", walletAddress, existingTeam, msg.GameID)
			h.clientTeams[walletAddress] = existingTeam

			// Add client to the game room if not already there
			h.AddClientToGame(client, msg.GameID)

			// Send initial game state to the reconnecting player
			votes := h.gameManager.GetVotes(msg.GameID)
			stats := h.gameManager.GetGameStats(msg.GameID)

			initialMsg := &Message{
				Type:   TypeVoteUpdate,
				GameID: msg.GameID,
				Votes:  votes,
			}

			h.updateStats(stats, initialMsg)

			// Send initial state to the reconnecting client
			if data, err := json.Marshal(initialMsg); err == nil {
				select {
				case client.send <- data:
				default:
				}
			}

			log.Printf("Player %s successfully reconnected to %s team", walletAddress, existingTeam)
		} else {
			// Player is not in the game yet - attempt to add them to the requested team
			if err := h.gameManager.AddPlayerToTeam(msg.GameID, walletAddress, msg.Team); err != nil {
				log.Printf("Failed to add player %s to team %s: %v", walletAddress, msg.Team, err)
				h.sendErrorToClient(client, err.Error())
				return
			}

			// Success - update local state
			h.clientTeams[walletAddress] = msg.Team
			log.Printf("Successfully added player %s to %s team", walletAddress, msg.Team)
		}

		// Broadcast updated games list since player count may have changed
		h.broadcastGamesListUpdate()

	case TypeWatchGame:
		log.Printf("Player %s watching game %s", client.id, msg.GameID)
		h.AddClientToGame(client, msg.GameID)

		// Broadcast updated games list since spectator count changed
		h.broadcastGamesListUpdate()

	case TypeRequestGamesList:
		log.Printf("🎯 Player %s requesting games list", client.id)
		gamesList := h.collectGamesInfo("all") // Return all games
		log.Printf("🔍 Collected %d games for list request", len(gamesList))
		for i, game := range gamesList {
			log.Printf("🔍 Game %d: %s - Status: %s - Move: %d - HasBoard: %t", i, game.GameID, game.Status, game.CurrentMove, len(game.Board) > 0)
			if len(game.Board) > 0 {
				log.Printf("🔍 Board sample: %v", game.Board[0])
			}
		}

		gamesMsg := &Message{
			Type:             TypeGamesList,
			GamesList:        gamesList,
			TotalConnections: h.GetTotalConnections(),
		}

		// Send games list to requesting client
		if data, err := json.Marshal(gamesMsg); err == nil {
			select {
			case client.send <- data:
			default:
			}
		}

	case TypeRequestFilteredGamesList:
		log.Printf("Player %s requesting filtered games list with filter: %s", client.id, msg.Filter)
		gamesList := h.collectGamesInfo("all") // Always send all games, let frontend filter

		gamesMsg := &Message{
			Type:             TypeGamesList,
			GamesList:        gamesList,
			TotalConnections: h.GetTotalConnections(),
			Filter:           msg.Filter, // Include the requested filter so frontend knows what was requested
		}

		// Send games list to requesting client
		if data, err := json.Marshal(gamesMsg); err == nil {
			select {
			case client.send <- data:
			default:
			}
		}

	case TypeNumberOfPlayers:
		h.broadcastToAll(&Message{
			Type:             TypeNumberOfPlayers,
			TotalConnections: h.GetTotalConnections(),
		})

	case TypeCheckPlayerStatus:
		walletAddress := msg.WalletAddress
		if walletAddress == "" {
			log.Printf("No wallet address provided for player status check from client %s", client.id)
			h.sendErrorToClient(client, "Wallet address is required for player status check")
			return
		}

		// Store the wallet address mapping for this client
		h.clientWallets[client] = walletAddress

		log.Printf("Checking player status for wallet %s in game %s", walletAddress, msg.GameID)

		// Check if player is already in the game
		team := h.gameManager.GetPlayerTeam(msg.GameID, walletAddress)

		statusMsg := &Message{
			Type:          TypePlayerStatus,
			GameID:        msg.GameID,
			WalletAddress: walletAddress,
			Team:          team, // Will be "white", "black", or "" if not in game
		}

		// Send status back to the requesting client
		if data, err := json.Marshal(statusMsg); err == nil {
			select {
			case client.send <- data:
			default:
			}
		}

	case TypeGetValidMoves:
		log.Printf("Player %s requesting valid moves for game %s", client.id, msg.GameID)

		// Check if game exists
		game := h.gameManager.GetGame(msg.GameID)
		if game == nil {
			log.Printf("Game %s does not exist, cannot get valid moves", msg.GameID)
			h.sendErrorToClient(client, "Game does not exist")
			return
		}

		// Get valid moves from the game manager (already in coordinate notation)
		validMoves := h.gameManager.GetValidMoves(msg.GameID)
		if validMoves == nil {
			validMoves = []string{} // Ensure we send an empty array instead of null
		}

		validMovesMsg := &Message{
			Type:       TypeValidMovesResponse,
			GameID:     msg.GameID,
			ValidMoves: validMoves,
		}

		// Send valid moves back to the requesting client
		if data, err := json.Marshal(validMovesMsg); err == nil {
			select {
			case client.send <- data:
			default:
			}
		}

	case TypeRequestPermit2:
		walletAddress := msg.WalletAddress
		chainId := msg.ChainId
		if walletAddress == "" {
			log.Printf("No wallet address provided for permit2 request from client %s", client.id)
			h.sendErrorToClient(client, "Wallet address is required for permit2 request")
			return
		}

		// Store the wallet address mapping for this client
		h.clientWallets[client] = walletAddress

		permit2Data, err := h.gameManager.GeneratePermit2SignatureData(walletAddress, chainId)
		if err != nil {
			log.Printf("Failed to generate permit2 data for wallet %s: %v", walletAddress, err)
			h.sendErrorToClient(client, "Failed to generate permit2 data")
			return
		}

		h.sendPermit2DataToClient(client, permit2Data)

	case TypeSubmitPermit2Signature:
		walletAddress := msg.WalletAddress
		signature := msg.Signature
		chainId := msg.ChainId
		if walletAddress == "" {
			log.Printf("No wallet address provided for permit2 signature from client %s", client.id)
			h.sendErrorToClient(client, "Wallet address is required for permit2 signature")
			return
		}

		if signature == "" {
			log.Printf("No signature provided for permit2 signature from client %s", client.id)
			h.sendErrorToClient(client, "Signature is required for permit2 signature")
			return
		}

		// Store the permit2 signature
		err := h.gameManager.StorePermit2Signature(walletAddress, signature)
		if err != nil {
			log.Printf("Failed to store permit2 signature for wallet %s: %v", walletAddress, err)
			h.sendErrorToClient(client, "Failed to store permit2 signature")
			return
		}

		// Execute the permit2 allowance
		err = h.gameManager.ExecutePermit2(walletAddress, chainId)
		if err != nil {
			log.Printf("Failed to execute permit2 for wallet %s: %v", walletAddress, err)
			h.sendErrorToClient(client, "Failed to execute permit2")
			return
		}

		log.Printf("Successfully executed permit2 for wallet %s", walletAddress)
	}
}

func (h *Hub) addToMatchmaking(client *Client, walletAddress string) {
	// Check if wallet address is already in queue
	if _, exists := h.matchmakingQueue[walletAddress]; exists {
		log.Printf("Wallet %s already in matchmaking queue", walletAddress)
		h.sendErrorToClient(client, "You are already in the matchmaking queue")
		return
	}

	h.matchmakingQueue[walletAddress] = client
	log.Printf("Matchmaking queue size: %d", len(h.matchmakingQueue))

	// Check if we can create a match
	if len(h.matchmakingQueue) >= 2 {
		// Get first 2 wallet addresses
		walletAddresses := make([]string, 0, 2)
		clients := make([]*Client, 0, 2)

		for walletAddr, client := range h.matchmakingQueue {
			walletAddresses = append(walletAddresses, walletAddr)
			clients = append(clients, client)
			if len(walletAddresses) == 2 {
				break
			}
		}

		// Remove matched players from queue
		delete(h.matchmakingQueue, walletAddresses[0])
		delete(h.matchmakingQueue, walletAddresses[1])

		player1Wallet := walletAddresses[0]
		player2Wallet := walletAddresses[1]
		player1Client := clients[0]
		player2Client := clients[1]

		// Create the game
		gameState := h.gameManager.GetOrCreateGame()

		// Create unique game ID
		gameID := gameState.ID
		// Assign sides randomly
		assignedSides := []string{"white", "black"}
		if time.Now().Unix()%2 == 0 {
			assignedSides[0], assignedSides[1] = assignedSides[1], assignedSides[0]
		}

		// Add players to teams using their wallet addresses
		if err := h.gameManager.AddPlayerToTeam(gameID, player1Wallet, assignedSides[0]); err != nil {
			log.Printf("Failed to add player %s to team %s: %v", player1Wallet, assignedSides[0], err)
			h.sendErrorToClient(player1Client, "Failed to join team")
			return
		}

		if err := h.gameManager.AddPlayerToTeam(gameID, player2Wallet, assignedSides[1]); err != nil {
			log.Printf("Failed to add player %s to team %s: %v", player2Wallet, assignedSides[1], err)
			h.sendErrorToClient(player2Client, "Failed to join team")
			return
		}

		// Notify both players
		matchMsg1 := &Message{
			Type:         TypeMatchFound,
			GameID:       gameID,
			Players:      []string{player1Client.id, player2Client.id},
			AssignedSide: assignedSides[0],
		}

		matchMsg2 := &Message{
			Type:         TypeMatchFound,
			GameID:       gameID,
			Players:      []string{player1Client.id, player2Client.id},
			AssignedSide: assignedSides[1],
		}

		// Send match found messages
		if data1, err := json.Marshal(matchMsg1); err == nil {
			log.Printf("Sending match_found to player1 (%s): %s", player1Client.id, string(data1))
			select {
			case player1Client.send <- data1:
				log.Printf("✅ Successfully sent match_found to player1 (%s)", player1Client.id)
			default:
				log.Printf("❌ Failed to send match_found to player1 (%s) - channel full", player1Client.id)
			}
		} else {
			log.Printf("❌ Failed to marshal match_found for player1 (%s): %v", player1Client.id, err)
		}

		if data2, err := json.Marshal(matchMsg2); err == nil {
			log.Printf("Sending match_found to player2 (%s): %s", player2Client.id, string(data2))
			select {
			case player2Client.send <- data2:
				log.Printf("✅ Successfully sent match_found to player2 (%s)", player2Client.id)
			default:
				log.Printf("❌ Failed to send match_found to player2 (%s) - channel full", player2Client.id)
			}
		} else {
			log.Printf("❌ Failed to marshal match_found for player2 (%s): %v", player2Client.id, err)
		}

		// Add both players to the game
		h.AddClientToGame(player1Client, gameID)
		h.AddClientToGame(player2Client, gameID)

		// Set their teams locally using wallet addresses
		h.clientTeams[player1Wallet] = assignedSides[0]
		h.clientTeams[player2Wallet] = assignedSides[1]

		log.Printf("Match created: %s with players %s (%s) and %s (%s)",
			gameID, player1Wallet, assignedSides[0], player2Wallet, assignedSides[1])

		// Broadcast updated games list since a new game was created
		h.broadcastGamesListUpdate()
	}
}

func (h *Hub) removeFromMatchmaking(client *Client) {
	// Find and remove the wallet address associated with this client
	for walletAddr, c := range h.matchmakingQueue {
		if c == client {
			delete(h.matchmakingQueue, walletAddr)
			log.Printf("Removed wallet %s from matchmaking queue", walletAddr)
			break
		}
	}
}

func (h *Hub) broadcastToGame(gameID string, msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	// Send to all clients in the game room
	if room, ok := h.gameRooms[gameID]; ok {
		for client := range room {
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(room, client)
			}
		}
	}
}

func (h *Hub) broadcastToAll(msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) startPeriodicUpdates() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Broadcast timer updates for all active games
		for gameID := range h.gameRooms {
			timeLeft := h.gameManager.GetTimeLeft(gameID)

			h.broadcastToGame(gameID, &Message{
				Type:        TypeTimerTick,
				GameID:      gameID,
				SecondsLeft: timeLeft,
			})
		}
	}
}

func (h *Hub) AddClientToGame(client *Client, gameID string) {
	if _, ok := h.gameRooms[gameID]; !ok {
		h.gameRooms[gameID] = make(map[*Client]bool)
	}
	h.gameRooms[gameID][client] = true
}

// Handle move result from game manager
func (h *Hub) handleMoveResult(gameID, move string) {
	log.Printf("Move executed in game %s: %s", gameID, move)

	// Get updated game statistics after the move
	stats := h.gameManager.GetGameStats(gameID)
	votes := h.gameManager.GetVotes(gameID)

	// Create move result message with updated stats
	moveMsg := &Message{
		Type:   TypeMoveResult,
		GameID: gameID,
		Move:   move,
		Votes:  votes,
	}

	h.updateStats(stats, moveMsg)

	// Broadcast the move result to all clients in the game
	h.broadcastToGame(gameID, moveMsg)

	// Broadcast updated games list since move number has changed
	h.broadcastGamesListUpdate()
}

// Handle game end from game manager
func (h *Hub) handleGameEnd(gameID, winner, reason string, gameStats map[string]any) {
	log.Printf("Game ended: %s, winner: %s, reason: %s", gameID, winner, reason)

	// Create game end message with all the required statistics
	gameEndMsg := &Message{
		Type:          TypeGameEnd,
		GameID:        gameID,
		Winner:        winner,
		GameEndReason: reason,
	}

	// Add all game statistics to the message
	h.updateStats(gameStats, gameEndMsg)

	// Extract and add team player statistics
	if whiteTeamPlayers, ok := gameStats["whiteTeamPlayers"].([]map[string]any); ok {
		log.Printf("White team players found: %d players", len(whiteTeamPlayers))
		gameEndMsg.WhiteTeamPlayers = make([]PlayerStats, len(whiteTeamPlayers))
		for i, player := range whiteTeamPlayers {
			gameEndMsg.WhiteTeamPlayers[i] = PlayerStats{
				WalletAddress: player["walletAddress"].(string),
				TotalVotes:    player["totalVotes"].(int),
				TotalSpent:    player["totalSpent"].(float64),
			}
			log.Printf("White player %d: %s - %d votes, %.3f USDC", i+1,
				player["walletAddress"].(string), player["totalVotes"].(int), player["totalSpent"].(float64))
		}
	} else {
		log.Printf("No white team players found in game stats")
	}

	if blackTeamPlayers, ok := gameStats["blackTeamPlayers"].([]map[string]any); ok {
		log.Printf("Black team players found: %d players", len(blackTeamPlayers))
		gameEndMsg.BlackTeamPlayers = make([]PlayerStats, len(blackTeamPlayers))
		for i, player := range blackTeamPlayers {
			gameEndMsg.BlackTeamPlayers[i] = PlayerStats{
				WalletAddress: player["walletAddress"].(string),
				TotalVotes:    player["totalVotes"].(int),
				TotalSpent:    player["totalSpent"].(float64),
			}
			log.Printf("Black player %d: %s - %d votes, %.3f USDC", i+1,
				player["walletAddress"].(string), player["totalVotes"].(int), player["totalSpent"].(float64))
		}
	} else {
		log.Printf("No black team players found in game stats")
	}

	// Broadcast to all clients in the game with their individual vote counts
	if room, exists := h.gameRooms[gameID]; exists {
		for client := range room {
			// Create a copy of the message for each client with their personal vote count
			clientMsg := *gameEndMsg
			if playerVotes, ok := gameStats["playerTotalVotes"].(map[string]int); ok {
				// Get the wallet address for this client
				if walletAddress, found := h.clientWallets[client]; found {
					if votes, found := playerVotes[walletAddress]; found {
						clientMsg.PlayerVotes = votes
					}
				}
			}

			// Send personalized message to each client
			if data, err := json.Marshal(clientMsg); err == nil {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(room, client)
				}
			}
		}
	}

	// Store the ended game info before cleaning up
	endedGameInfo := &GameInfo{
		GameID:    gameID,
		Status:    "ended",
		Winner:    winner,
		EndReason: reason,
		CreatedAt: h.gameManager.GetGameCreatedAt(gameID),
		EndedAt:   func() *int64 { t := time.Now().Unix(); return &t }(),
	}

	// Extract final game stats
	if whitePlayers, ok := gameStats["whitePlayers"].(int); ok {
		endedGameInfo.WhitePlayers = whitePlayers
	}
	if blackPlayers, ok := gameStats["blackPlayers"].(int); ok {
		endedGameInfo.BlackPlayers = blackPlayers
	}
	if currentMove, ok := gameStats["currentMove"].(int); ok {
		endedGameInfo.CurrentMove = currentMove
	}
	if totalPot, ok := gameStats["totalPot"].(float64); ok {
		endedGameInfo.TotalPot = totalPot
	}
	if whitePot, ok := gameStats["whitePot"].(float64); ok {
		endedGameInfo.WhitePot = whitePot
	}
	if blackPot, ok := gameStats["blackPot"].(float64); ok {
		endedGameInfo.BlackPot = blackPot
	}

	// Extract board data for ended games
	if board, ok := gameStats["board"].([][]string); ok {
		endedGameInfo.Board = board
	}

	// Extract team player statistics
	if whiteTeamPlayers, ok := gameStats["whiteTeamPlayers"].([]map[string]any); ok {
		endedGameInfo.WhiteTeamPlayers = make([]PlayerStats, len(whiteTeamPlayers))
		for i, player := range whiteTeamPlayers {
			endedGameInfo.WhiteTeamPlayers[i] = PlayerStats{
				WalletAddress: player["walletAddress"].(string),
				TotalVotes:    player["totalVotes"].(int),
				TotalSpent:    player["totalSpent"].(float64),
			}
		}
	}

	if blackTeamPlayers, ok := gameStats["blackTeamPlayers"].([]map[string]any); ok {
		endedGameInfo.BlackTeamPlayers = make([]PlayerStats, len(blackTeamPlayers))
		for i, player := range blackTeamPlayers {
			endedGameInfo.BlackTeamPlayers[i] = PlayerStats{
				WalletAddress: player["walletAddress"].(string),
				TotalVotes:    player["totalVotes"].(int),
				TotalSpent:    player["totalSpent"].(float64),
			}
		}
	}

	// Store the ended game
	h.endedGames[gameID] = endedGameInfo

	// Clean up the game room since the game has ended
	delete(h.gameRooms, gameID)

	// Broadcast updated games list since game has ended
	h.broadcastGamesListUpdate()
}

// collectGamesInfo gathers information about games based on filter
// filter can be "active", "ended", or "all" for all games
func (h *Hub) collectGamesInfo(filter string) []GameInfo {
	var gamesList []GameInfo

	// Include active games if filter is "active" or "all"
	if filter == "active" || filter == "all" {
		// Iterate through all games in the game manager (not just game rooms)
		allGameIDs := h.gameManager.GetAllGames()
		log.Printf("🔍 Found %d total games in manager", len(allGameIDs))

		for _, gameID := range allGameIDs {
			// Skip if this game is already ended
			if _, isEnded := h.endedGames[gameID]; isEnded {
				continue
			}

			// Get game statistics from game manager
			stats := h.gameManager.GetGameStats(gameID)
			if stats == nil {
				log.Printf("🔍 No stats for game %s", gameID)
				continue
			}
			log.Printf("🔍 Got stats for game %s - HasBoard: %t", gameID, stats["board"] != nil)

			// Count spectators (clients in room who are not on a team)
			spectators := 0
			if room, exists := h.gameRooms[gameID]; exists {
				for client := range room {
					// Get wallet address for this client
					if walletAddress, found := h.clientWallets[client]; found {
						if team, ok := h.clientTeams[walletAddress]; !ok || team == "" {
							spectators++
						}
					} else {
						// If no wallet address found, count as spectator
						spectators++
					}
				}
			}

			// Create GameInfo struct for active game
			gameInfo := GameInfo{
				GameID:     gameID,
				Status:     "active",
				CreatedAt:  h.gameManager.GetGameCreatedAt(gameID),
				Spectators: spectators,
			}

			// Extract stats from the map
			if whitePlayers, ok := stats["whitePlayers"].(int); ok {
				gameInfo.WhitePlayers = whitePlayers
			}
			if blackPlayers, ok := stats["blackPlayers"].(int); ok {
				gameInfo.BlackPlayers = blackPlayers
			}
			if timeLeft, ok := stats["timeLeft"].(int); ok {
				gameInfo.TimeLeft = timeLeft
			}
			if currentMove, ok := stats["currentMove"].(int); ok {
				gameInfo.CurrentMove = currentMove
			}
			if totalPot, ok := stats["totalPot"].(float64); ok {
				gameInfo.TotalPot = totalPot
			}
			if whitePot, ok := stats["whitePot"].(float64); ok {
				gameInfo.WhitePot = whitePot
			}
			if blackPot, ok := stats["blackPot"].(float64); ok {
				gameInfo.BlackPot = blackPot
			}
			if currentTurn, ok := stats["currentTurn"].(string); ok {
				gameInfo.CurrentTurn = currentTurn
			}
			if board, ok := stats["board"].([][]string); ok {
				gameInfo.Board = board
			}

			gamesList = append(gamesList, gameInfo)
		}
	}

	// Include ended games if filter is "ended" or "all"
	if filter == "ended" || filter == "all" {
		for _, endedGame := range h.endedGames {
			// Count current spectators for ended games (people who might be viewing the final state)
			spectators := 0
			if room, exists := h.gameRooms[endedGame.GameID]; exists {
				for client := range room {
					// Get wallet address for this client
					if walletAddress, found := h.clientWallets[client]; found {
						if team, ok := h.clientTeams[walletAddress]; !ok || team == "" {
							spectators++
						}
					} else {
						// If no wallet address found, count as spectator
						spectators++
					}
				}
			}

			// Create a copy and update spectators
			gameInfo := *endedGame
			gameInfo.Spectators = spectators
			gamesList = append(gamesList, gameInfo)
		}
	}

	return gamesList
}

func (h *Hub) updateStats(stats map[string]any, moveMsg *Message) {
	if stats == nil {
		return
	}

	if whitePlayers, ok := stats["whitePlayers"].(int); ok {
		moveMsg.WhitePlayers = whitePlayers
	}
	if blackPlayers, ok := stats["blackPlayers"].(int); ok {
		moveMsg.BlackPlayers = blackPlayers
	}
	if wv, ok := stats["whiteCurrentTurnVotes"].(int); ok {
		moveMsg.WhiteCurrentTurnVotes = wv
	}
	if bv, ok := stats["blackCurrentTurnVotes"].(int); ok {
		moveMsg.BlackCurrentTurnVotes = bv
	}
	if wttv, ok := stats["whiteTeamTotalVotes"].(int); ok {
		moveMsg.WhiteTeamTotalVotes = wttv
	}
	if bttv, ok := stats["blackTeamTotalVotes"].(int); ok {
		moveMsg.BlackTeamTotalVotes = bttv
	}
	if tp, ok := stats["totalPot"].(float64); ok {
		moveMsg.TotalPot = tp
	}
	if whitePot, ok := stats["whitePot"].(float64); ok {
		moveMsg.WhitePot = whitePot
	}
	if blackPot, ok := stats["blackPot"].(float64); ok {
		moveMsg.BlackPot = blackPot
	}
	if ct, ok := stats["currentTurn"].(string); ok {
		moveMsg.CurrentTurn = ct
	}
	if cm, ok := stats["currentMove"].(int); ok {
		moveMsg.CurrentMove = cm
	}
	if pvr, ok := stats["playerVotedThisRound"].(map[string]bool); ok {
		moveMsg.PlayerVotedThisRound = pvr
	}
	if ptv, ok := stats["playerTotalVotes"].(map[string]int); ok {
		moveMsg.PlayerTotalVotes = ptv
	}
	if board, ok := stats["board"].([][]string); ok {
		moveMsg.Board = board
	}
	if isInCheck, ok := stats["isInCheck"].(bool); ok {
		moveMsg.IsInCheck = isInCheck
	}
	if isCheckmate, ok := stats["isCheckmate"].(bool); ok {
		moveMsg.IsCheckmate = isCheckmate
	}
}

// GetTotalConnections returns the total number of connected clients
func (h *Hub) GetTotalConnections() int {
	return len(h.clients)
}

// broadcastGamesListUpdate sends game list updates to all connected clients
func (h *Hub) broadcastGamesListUpdate() {
	gamesList := h.collectGamesInfo("all") // Return all games by default

	h.broadcastToAll(&Message{
		Type:             TypeGamesListUpdate,
		GamesList:        gamesList,
		TotalConnections: h.GetTotalConnections(),
	})
}

// Send error message to a specific client
func (h *Hub) sendErrorToClient(client *Client, errorMsg string) {
	errorMessage := &Message{
		Type:  TypeError,
		Error: errorMsg,
	}

	if data, err := json.Marshal(errorMessage); err == nil {
		select {
		case client.send <- data:
		default:
			log.Printf("Failed to send error message to client %s: channel full", client.id)
		}
	} else {
		log.Printf("Failed to marshal error message for client %s: %v", client.id, err)
	}
}

// sendPermit2DataToClient sends Permit2 typed data to a specific client
func (h *Hub) sendPermit2DataToClient(client *Client, permit2Data interface{}) {
	data, err := json.Marshal(map[string]interface{}{
		"type":        TypePermit2Data,
		"permit2Data": permit2Data,
	})

	if err != nil {
		log.Printf("Failed to marshal permit2 data for client %s: %v", client.id, err)
		h.sendErrorToClient(client, "Failed to generate permit2 data")
		return
	}

	select {
	case client.send <- data:
		log.Printf("Successfully sent permit2 data to client %s", client.id)
	default:
		log.Printf("Failed to send permit2 data to client %s: channel full", client.id)
	}
}
