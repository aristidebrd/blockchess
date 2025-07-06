package game

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// EndGameProcessor handles the end game logic and reward distribution
type EndGameProcessor struct {
	blockchainService BlockchainService
}

// NewEndGameProcessor creates a new end game processor
func NewEndGameProcessor(blockchainService BlockchainService) *EndGameProcessor {
	return &EndGameProcessor{
		blockchainService: blockchainService,
	}
}

// ProcessGameEnd handles the complete game ending process including reward distribution
func (e *EndGameProcessor) ProcessGameEnd(gameID, winner string, game *GameState) error {
	log.Printf("Processing game end for game %s, winner: %s", gameID, winner)

	// End game on blockchain
	if e.blockchainService != nil && e.blockchainService.IsConnected() {
		var blockchainResult GameResult
		switch winner {
		case "white":
			blockchainResult = GameResultWhiteWins
		case "black":
			blockchainResult = GameResultBlackWins
		case "draw":
			blockchainResult = GameResultDraw
		default:
			blockchainResult = GameResultDraw
		}

		if err := e.blockchainService.EndGame(gameID, blockchainResult); err != nil {
			log.Printf("Error ending game %s on blockchain: %v", gameID, err)
			return err
		}

		// Distribute rewards to winning team
		if err := e.distributeRewards(gameID, winner, game); err != nil {
			log.Printf("Error distributing rewards for game %s: %v", gameID, err)
			return err
		}
	}

	return nil
}

// distributeRewards handles cross-chain reward distribution to players
func (e *EndGameProcessor) distributeRewards(gameID, winner string, game *GameState) error {
	if winner == "draw" {
		log.Printf("Game %s ended in draw, no rewards to distribute", gameID)
		return nil
	}

	// Get winning team players
	var winningPlayers map[string]bool
	if winner == "white" {
		winningPlayers = game.WhitePlayers
	} else {
		winningPlayers = game.BlackPlayers
	}

	if len(winningPlayers) == 0 {
		log.Printf("No players in winning team for game %s", gameID)
		return nil
	}

	// Calculate total votes by winning team
	totalWinningVotes := 0
	for playerAddress := range winningPlayers {
		totalWinningVotes += game.PlayerTotalVotes[playerAddress]
	}

	if totalWinningVotes == 0 {
		log.Printf("No votes from winning team in game %s", gameID)
		return nil
	}

	// Convert total pot to USDC wei (6 decimals)
	totalPotWei := big.NewInt(int64(game.TotalPot * 1e6))

	// Distribute rewards proportionally based on votes
	for playerAddress := range winningPlayers {
		playerVotes := game.PlayerTotalVotes[playerAddress]
		if playerVotes == 0 {
			continue
		}

		// Calculate player's share based on their votes
		playerShare := float64(playerVotes) / float64(totalWinningVotes)
		playerRewardWei := new(big.Int).Mul(totalPotWei, big.NewInt(int64(playerShare*1e6)))
		playerRewardWei = new(big.Int).Div(playerRewardWei, big.NewInt(1e6))

		// Transfer rewards cross-chain to Base Sepolia (chain ID 84532)
		if err := e.blockchainService.TransferRewardsCrossChain(
			gameID,
			playerRewardWei,
			big.NewInt(84532), // Base Sepolia chain ID
			common.HexToAddress(playerAddress),
			false,                           // Use standard transfer (not fast)
			big.NewInt(1000000000000000000), // Max fee: 1 ETH worth
		); err != nil {
			log.Printf("Error transferring rewards to player %s: %v", playerAddress, err)
			continue
		}

		log.Printf("Transferred %s USDC to player %s (votes: %d, share: %.2f%%)",
			playerRewardWei.String(), playerAddress, playerVotes, playerShare*100)
	}

	return nil
}

// updatePlayerBalances updates player balances after game end (legacy function)
func (e *EndGameProcessor) updatePlayerBalances(gameID string, users []User) error {
	// This is a placeholder for the legacy balance update logic
	// You may want to implement this based on your user management system
	for _, user := range users {
		// user.updateBalance(user.balance + amount)
		log.Printf("Would update balance for user %s in game %s", user.ID, gameID)
	}
	return nil
}

// User represents a game user (placeholder struct)
type User struct {
	ID      string
	Balance float64
}
