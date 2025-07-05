package main

import (
	"blockchess/internal/game"
	"log"
)

func main() {
	// Option 1: Create manager without contract integration (existing behavior)
	manager := game.NewManager()

	// Option 2: Create manager with contract integration
	config, err := game.LoadContractConfigFromEnv()
	if err != nil {
		log.Printf("Failed to load contract config: %v", err)
		log.Printf("Using manager without contract integration...")
		manager = game.NewManager()
	} else {
		manager = game.NewManagerWithContracts(config)
	}

	// Set up callbacks for move results and game end
	manager.SetMoveResultCallback(func(gameID, move string) {
		log.Printf("Move result for game %s: %s", gameID, move)
		// Broadcast to websocket clients, etc.
	})

	manager.SetGameEndCallback(func(gameID, winner, reason string, gameStats map[string]any) {
		log.Printf("Game %s ended: winner=%s, reason=%s", gameID, winner, reason)
		// Handle game end, distribute rewards, etc.
	})

	// Create a game - this will now also create it on-chain if contracts are configured
	gameState := manager.GetOrCreateGame("1")
	log.Printf("Created game: %s", gameState.ID)

	// Your existing websocket/HTTP server code would go here...
	// The game manager now automatically handles contract integration
}
