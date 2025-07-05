package game

import (
	"blockchess/contracts/gamecontract"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockchainService interface defines blockchain operations
type BlockchainService interface {
	CreateGame(gameID string, stakeAmount *big.Int) error
	RecordMove(gameID string, player common.Address, chainID uint32) error
	EndGame(gameID string, result GameResult) error
	IsConnected() bool
	GetGameInfo(gameID string) (*GameInfo, error)
}

// GameResult represents the outcome of a game
type GameResult uint8

const (
	GameResultOngoing GameResult = iota
	GameResultWhiteWins
	GameResultBlackWins
	GameResultDraw
)

// GameInfo represents game information from the blockchain
type GameInfo struct {
	GameID           *big.Int
	FixedStakeAmount *big.Int
	TotalWhiteStakes *big.Int
	TotalBlackStakes *big.Int
	WhitePlayerCount *big.Int
	BlackPlayerCount *big.Int
	CreatedAt        *big.Int
	EndedAt          *big.Int
	IsActive         bool
	Result           GameResult
}

// EthereumBlockchainService implements BlockchainService for Ethereum
type EthereumBlockchainService struct {
	client           bind.ContractBackend
	gameContract     *gamecontract.Gamecontract
	gameContractAddr common.Address
	auth             *bind.TransactOpts
	privateKey       *ecdsa.PrivateKey
	defaultStakeWei  *big.Int
}

// BlockchainConfig holds blockchain configuration
type BlockchainConfig struct {
	RPCUrl              string
	PrivateKey          string
	GameContractAddress string
	DefaultStakeETH     string
}

// NewEthereumBlockchainService creates a new Ethereum blockchain service
func NewEthereumBlockchainService(config *BlockchainConfig) (*EthereumBlockchainService, error) {
	// Connect to Ethereum client
	client, err := ethclient.Dial(config.RPCUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	// Parse private key
	privateKeyHex := strings.TrimPrefix(config.PrivateKey, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// Get chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	// Set gas parameters
	auth.GasLimit = uint64(3000000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}
	auth.GasPrice = gasPrice

	// Parse default stake amount
	stakeFloat, err := strconv.ParseFloat(config.DefaultStakeETH, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse default stake amount: %v", err)
	}
	stakeWei := new(big.Float).Mul(big.NewFloat(stakeFloat), big.NewFloat(1e18))
	defaultStakeWei, _ := stakeWei.Int(nil)

	// Create game contract instance
	gameContractAddr := common.HexToAddress(config.GameContractAddress)
	gameContract, err := gamecontract.NewGamecontract(gameContractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create game contract instance: %v", err)
	}

	service := &EthereumBlockchainService{
		client:           client,
		gameContract:     gameContract,
		gameContractAddr: gameContractAddr,
		auth:             auth,
		privateKey:       privateKey,
		defaultStakeWei:  defaultStakeWei,
	}

	log.Printf("Blockchain service initialized successfully")
	log.Printf("Game Contract: %s", gameContractAddr.Hex())
	log.Printf("Default Stake: %s Wei (%s ETH)", defaultStakeWei.String(), config.DefaultStakeETH)

	return service, nil
}

// CreateGame creates a new game on the blockchain
func (s *EthereumBlockchainService) CreateGame(gameID string, stakeAmount *big.Int) error {
	gameIDInt, err := s.parseGameID(gameID)
	if err != nil {
		return fmt.Errorf("invalid game ID: %v", err)
	}

	// Use provided stake amount or default
	if stakeAmount == nil {
		stakeAmount = s.defaultStakeWei
	}

	// Get fresh nonce
	nonce, err := s.client.PendingNonceAt(context.Background(), s.auth.From)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}
	s.auth.Nonce = big.NewInt(int64(nonce))

	// Create game on contract
	tx, err := s.gameContract.CreateGame(s.auth, big.NewInt(int64(gameIDInt)), stakeAmount)
	if err != nil {
		return fmt.Errorf("failed to create game on-chain: %v", err)
	}

	log.Printf("Game %s created on-chain! Transaction: %s", gameID, tx.Hash().Hex())
	return nil
}

// RecordMove records a move on the blockchain
func (s *EthereumBlockchainService) RecordMove(gameID string, player common.Address, chainID uint32) error {
	gameIDInt, err := s.parseGameID(gameID)
	if err != nil {
		return fmt.Errorf("invalid game ID: %v", err)
	}

	// Get fresh nonce
	nonce, err := s.client.PendingNonceAt(context.Background(), s.auth.From)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}
	s.auth.Nonce = big.NewInt(int64(nonce))

	// Record move on contract
	tx, err := s.gameContract.RecordMove(s.auth, big.NewInt(int64(gameIDInt)), player, chainID)
	if err != nil {
		return fmt.Errorf("failed to record move on-chain: %v", err)
	}

	log.Printf("Move recorded for game %s, player %s! Transaction: %s", gameID, player.Hex(), tx.Hash().Hex())
	return nil
}

// EndGame ends a game on the blockchain
func (s *EthereumBlockchainService) EndGame(gameID string, result GameResult) error {
	gameIDInt, err := s.parseGameID(gameID)
	if err != nil {
		return fmt.Errorf("invalid game ID: %v", err)
	}

	// Get fresh nonce
	nonce, err := s.client.PendingNonceAt(context.Background(), s.auth.From)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}
	s.auth.Nonce = big.NewInt(int64(nonce))

	// End game on contract
	tx, err := s.gameContract.EndGame(s.auth, big.NewInt(int64(gameIDInt)), uint8(result))
	if err != nil {
		return fmt.Errorf("failed to end game on-chain: %v", err)
	}

	log.Printf("Game %s ended on-chain with result %d! Transaction: %s", gameID, result, tx.Hash().Hex())
	return nil
}

// IsConnected checks if the blockchain service is connected
func (s *EthereumBlockchainService) IsConnected() bool {
	if s.client == nil {
		return false
	}

	// Try to get the latest block number as a connectivity test
	_, err := s.client.PendingNonceAt(context.Background(), s.auth.From)
	return err == nil
}

// GetGameInfo retrieves game information from the blockchain
func (s *EthereumBlockchainService) GetGameInfo(gameID string) (*GameInfo, error) {
	gameIDInt, err := s.parseGameID(gameID)
	if err != nil {
		return nil, fmt.Errorf("invalid game ID: %v", err)
	}

	// Get game info from contract
	gameInfo, err := s.gameContract.GetGameInfo(nil, big.NewInt(int64(gameIDInt)))
	if err != nil {
		return nil, fmt.Errorf("failed to get game info: %v", err)
	}

	// Check if game is active
	isActive, err := s.gameContract.IsGameActive(nil, big.NewInt(int64(gameIDInt)))
	if err != nil {
		return nil, fmt.Errorf("failed to check game status: %v", err)
	}

	return &GameInfo{
		GameID:           gameInfo.GameId,
		FixedStakeAmount: gameInfo.FixedStakeAmount,
		TotalWhiteStakes: gameInfo.TotalWhiteStakes,
		TotalBlackStakes: gameInfo.TotalBlackStakes,
		WhitePlayerCount: gameInfo.WhitePlayerCount,
		BlackPlayerCount: gameInfo.BlackPlayerCount,
		CreatedAt:        gameInfo.CreatedAt,
		EndedAt:          gameInfo.EndedAt,
		IsActive:         isActive,
		Result:           GameResult(gameInfo.Result),
	}, nil
}

// parseGameID extracts the numeric ID from a game ID string
func (s *EthereumBlockchainService) parseGameID(gameID string) (uint64, error) {
	if strings.HasPrefix(gameID, "game-") {
		idStr := strings.TrimPrefix(gameID, "game-")
		return strconv.ParseUint(idStr, 10, 64)
	}
	return strconv.ParseUint(gameID, 10, 64)
}

// MockBlockchainService implements BlockchainService for testing
type MockBlockchainService struct {
	connected bool
	games     map[string]*GameInfo
}

// NewMockBlockchainService creates a new mock blockchain service
func NewMockBlockchainService() *MockBlockchainService {
	return &MockBlockchainService{
		connected: true,
		games:     make(map[string]*GameInfo),
	}
}

func (m *MockBlockchainService) CreateGame(gameID string, stakeAmount *big.Int) error {
	if !m.connected {
		return fmt.Errorf("blockchain service not connected")
	}

	gameIDInt, _ := strconv.ParseUint(strings.TrimPrefix(gameID, "game-"), 10, 64)
	m.games[gameID] = &GameInfo{
		GameID:           big.NewInt(int64(gameIDInt)),
		FixedStakeAmount: stakeAmount,
		IsActive:         true,
		Result:           GameResultOngoing,
	}

	log.Printf("Mock: Game %s created", gameID)
	return nil
}

func (m *MockBlockchainService) RecordMove(gameID string, player common.Address, chainID uint32) error {
	if !m.connected {
		return fmt.Errorf("blockchain service not connected")
	}
	log.Printf("Mock: Move recorded for game %s, player %s", gameID, player.Hex())
	return nil
}

func (m *MockBlockchainService) EndGame(gameID string, result GameResult) error {
	if !m.connected {
		return fmt.Errorf("blockchain service not connected")
	}

	if game, exists := m.games[gameID]; exists {
		game.IsActive = false
		game.Result = result
	}

	log.Printf("Mock: Game %s ended with result %d", gameID, result)
	return nil
}

func (m *MockBlockchainService) IsConnected() bool {
	return m.connected
}

func (m *MockBlockchainService) GetGameInfo(gameID string) (*GameInfo, error) {
	if !m.connected {
		return nil, fmt.Errorf("blockchain service not connected")
	}

	if game, exists := m.games[gameID]; exists {
		return game, nil
	}

	return nil, fmt.Errorf("game not found")
}
