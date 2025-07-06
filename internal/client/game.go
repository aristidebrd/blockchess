package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"blockchess/contracts-bindings/gamefactory"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GameFactory wraps the GameFactory contract instance
type GameFactory struct {
	contract *gamefactory.Gamefactory
	client   *ethclient.Client
	auth     *bind.TransactOpts
}

// NewGameFactory creates a new GameFactory instance connected to Base Sepolia
func NewGameFactory(client *ethclient.Client, privateKey string) (*GameFactory, error) {
	// Get GameFactory contract address from environment
	factoryAddress := GetGameFactoryAddress()
	if factoryAddress == "" {
		return nil, fmt.Errorf("GAME_FACTORY_ADDRESS not set in environment")
	}

	if !common.IsHexAddress(factoryAddress) {
		return nil, fmt.Errorf("invalid GameFactory address: %s", factoryAddress)
	}

	// Create contract instance
	contract, err := gamefactory.NewGamefactory(common.HexToAddress(factoryAddress), client)
	if err != nil {
		return nil, fmt.Errorf("failed to create GameFactory contract instance: %w", err)
	}

	// Validate private key
	if privateKey == "" {
		return nil, fmt.Errorf("private key cannot be empty")
	}

	// Remove 0x prefix if present
	privateKey = strings.TrimPrefix(privateKey, "0x")

	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get chain ID for transaction auth
	chainID, err := client.ChainID(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	return &GameFactory{
		contract: contract,
		client:   client,
		auth:     auth,
	}, nil
}

// CreateGame creates a new game with the specified stake amount and returns the game ID
func (gf *GameFactory) CreateGame(fixedStakeAmount *big.Int) (uint64, error) {
	log.Printf("Creating new game with stake amount: %s USDC", fixedStakeAmount.String())

	// Call the createGame function
	tx, err := gf.contract.CreateGame(gf.auth, fixedStakeAmount)
	if err != nil {
		return 0, fmt.Errorf("failed to create game transaction: %w", err)
	}

	log.Printf("Game creation transaction sent: %s", tx.Hash().Hex())

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), gf.client, tx)
	if err != nil {
		return 0, fmt.Errorf("failed to wait for transaction to be mined: %w", err)
	}

	if receipt.Status != 1 {
		return 0, fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	// Parse the GameCreated event to get the game ID
	for _, vLog := range receipt.Logs {
		event, err := gf.contract.ParseGameCreated(*vLog)
		if err == nil {
			gameID := event.GameId.Uint64()
			log.Printf("Game created successfully with ID: %d", gameID)
			return gameID, nil
		}
	}

	return 0, fmt.Errorf("failed to find GameCreated event in transaction receipt")
}

// AddVote adds a vote for a player to a game
func (gf *GameFactory) AddVote(gameID uint64, playerAddress common.Address, chainID uint32, team uint8) error {
	log.Printf("Adding vote for player %s in game %d on chain %d for team %d",
		playerAddress.Hex(), gameID, chainID, team)

	gameIDBig := new(big.Int).SetUint64(gameID)

	// Call the addVote function
	tx, err := gf.contract.AddVote(gf.auth, gameIDBig, playerAddress, chainID, team)
	if err != nil {
		return fmt.Errorf("failed to add vote transaction: %w", err)
	}

	log.Printf("Add vote transaction sent: %s", tx.Hash().Hex())

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), gf.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	log.Printf("Vote added successfully for player %s in game %d", playerAddress.Hex(), gameID)
	return nil
}

// EndGame ends a game with the specified result
func (gf *GameFactory) EndGame(gameID uint64, result uint8) error {
	log.Printf("Ending game %d with result: %d", gameID, result)

	gameIDBig := new(big.Int).SetUint64(gameID)

	// Call the endGame function
	tx, err := gf.contract.EndGame(gf.auth, gameIDBig, result)
	if err != nil {
		return fmt.Errorf("failed to end game transaction: %w", err)
	}

	log.Printf("End game transaction sent: %s", tx.Hash().Hex())

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), gf.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	log.Printf("Game %d ended successfully with result: %d", gameID, result)
	return nil
}

// GetGameExists checks if a game exists
func (gf *GameFactory) GetGameExists(gameID uint64) (bool, error) {
	gameIDBig := new(big.Int).SetUint64(gameID)
	exists, err := gf.contract.GameExists(nil, gameIDBig)
	if err != nil {
		return false, fmt.Errorf("failed to check if game exists: %w", err)
	}
	return exists, nil
}

// GetGameInfo retrieves game information from the contract
func (gf *GameFactory) GetGameInfo(gameID uint64) (*gamefactory.IGameFactoryGame, error) {
	gameIDBig := new(big.Int).SetUint64(gameID)
	game, err := gf.contract.GetGame(nil, gameIDBig)
	if err != nil {
		return nil, fmt.Errorf("failed to get game info: %w", err)
	}
	return &game, nil
}

// GetPlayerVoteCounts retrieves vote counts for all players in a game
func (gf *GameFactory) GetPlayerVoteCounts(gameID uint64) ([]common.Address, []*big.Int, []uint32, []*big.Int, error) {
	gameIDBig := new(big.Int).SetUint64(gameID)
	result, err := gf.contract.GetPlayerVoteCounts(nil, gameIDBig)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get player vote counts: %w", err)
	}
	return result.Players, result.VoteCounts, result.ChainIds, result.Teams, nil
}

// Team constants for game contract
const (
	TeamWhite = uint8(0)
	TeamBlack = uint8(1)
)

// GameResult constants for game contract
const (
	GameResultWhiteWins = uint8(0)
	GameResultBlackWins = uint8(1)
	GameResultDraw      = uint8(2)
)

// Helper function to convert team string to uint8
func TeamStringToUint8(team string) (uint8, error) {
	switch strings.ToLower(team) {
	case "white":
		return TeamWhite, nil
	case "black":
		return TeamBlack, nil
	default:
		return 0, fmt.Errorf("invalid team: %s", team)
	}
}

// Helper function to convert result string to uint8
func ResultStringToUint8(result string) (uint8, error) {
	switch strings.ToLower(result) {
	case "white":
		return GameResultWhiteWins, nil
	case "black":
		return GameResultBlackWins, nil
	case "draw":
		return GameResultDraw, nil
	default:
		return 0, fmt.Errorf("invalid result: %s", result)
	}
}
