package game

import (
	"blockchess/contracts-bindings/gamecontract"
	"blockchess/contracts-bindings/gamefactory"
	"blockchess/contracts-bindings/permit2"
	"blockchess/contracts-bindings/vaultcontract"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// BlockchainService interface defines blockchain operations
type BlockchainService interface {
	CreateGame(stakeAmount *big.Int) (*GameCreationResult, error)
	RecordMove(gameID string, player common.Address, team uint8) error
	EndGame(gameID string, result GameResult) error
	IsConnected() bool
	GetGameInfo(gameID string) (*GameInfo, error)
	StakeOnVote(gameID string, player common.Address, stakeAmount *big.Int) error
	RequestPlayerApproval(playerAddress common.Address, gameID string, approvalAmount *big.Int) (map[string]*ApprovalTransactionData, error)
	TransferRewardsCrossChain(gameID string, amount *big.Int, destinationChainId *big.Int, recipient common.Address, useFastTransfer bool, maxFee *big.Int) error
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
	client                  bind.ContractBackend
	gameContract            *gamecontract.Gamecontract
	gameContractAddr        common.Address
	vaultContract           *vaultcontract.Vaultcontract
	vaultContractAddr       common.Address
	gameFactoryContract     *gamefactory.Gamefactory
	gameFactoryContractAddr common.Address
	auth                    *bind.TransactOpts
	privateKey              *ecdsa.PrivateKey
	defaultStakeUSDC        *big.Int
	permits                 map[common.Address]*PermitData // Store Permit2 signatures with data
}

// BlockchainConfig holds blockchain configuration
type BlockchainConfig struct {
	RPCUrl               string
	PrivateKey           string
	GameContractAddress  string
	VaultContractAddress string
	DefaultStakeUSDC     string
}

// USDC contract address on Base Sepolia
const USDCContractAddress = "0x036CbD53842c5426634e7929541eC2318f3dCF7e"

// EIP-712 domain separator for USDC permit
const USDCPermitTypehash = "0x6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c9"
const USDCDomainSeparator = "0x" // This will be computed dynamically

// Permit2 contract address (same on all chains)
var Permit2ContractAddressByChain = map[uint32]string{
	31337: "0x000000000022D473030F116dDEE9F6B43aC78BA3",
}

// Permit2 signature data
type Permit2Data struct {
	Permitted struct {
		Token  common.Address `json:"token"`
		Amount *big.Int       `json:"amount"`
	} `json:"permitted"`
	Spender  common.Address `json:"spender"`
	Nonce    *big.Int       `json:"nonce"`
	Deadline *big.Int       `json:"deadline"`
}

// ApprovalTransactionData represents the data needed for a player to sign an approval transaction
type ApprovalTransactionData struct {
	To       string `json:"to"`       // Contract address to approve
	Data     string `json:"data"`     // Transaction data (encoded function call)
	Value    string `json:"value"`    // Always "0x0" for ERC20 approve
	GasLimit string `json:"gasLimit"` // Estimated gas limit
	GasPrice string `json:"gasPrice"` // Current gas price
}

// GameCreationResult holds the result of creating a game
type GameCreationResult struct {
	GameID              *big.Int
	GameContractAddress common.Address
	TransactionHash     common.Hash
	BlockNumber         uint64
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
	stakeFloat, err := strconv.ParseFloat(config.DefaultStakeUSDC, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse default stake amount: %v", err)
	}
	stakeUSDC := new(big.Float).Mul(big.NewFloat(stakeFloat), big.NewFloat(1e6))
	defaultStakeUSDC, _ := stakeUSDC.Int(nil)

	// Create game contract instance
	gameContractAddr := common.HexToAddress(config.GameContractAddress)
	gameContract, err := gamecontract.NewGamecontract(gameContractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create game contract instance: %v", err)
	}

	// Create vault contract instance
	vaultContractAddr := common.HexToAddress(config.VaultContractAddress)
	vaultContract, err := vaultcontract.NewVaultcontract(vaultContractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault contract instance: %v", err)
	}

	// Create game factory contract instance (using same address as game contract for now)
	gameFactoryContractAddr := gameContractAddr // TODO: Update with actual GameFactory address
	gameFactoryContract, err := gamefactory.NewGamefactory(gameFactoryContractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create game factory contract instance: %v", err)
	}

	service := &EthereumBlockchainService{
		client:                  client,
		gameContract:            gameContract,
		gameContractAddr:        gameContractAddr,
		vaultContract:           vaultContract,
		vaultContractAddr:       vaultContractAddr,
		gameFactoryContract:     gameFactoryContract,
		gameFactoryContractAddr: gameFactoryContractAddr,
		auth:                    auth,
		privateKey:              privateKey,
		defaultStakeUSDC:        defaultStakeUSDC,
		permits:                 make(map[common.Address]*PermitData),
	}

	log.Printf("Blockchain service initialized successfully")
	log.Printf("Game Contract: %s", gameContractAddr.Hex())
	log.Printf("Vault Contract: %s", vaultContractAddr.Hex())
	log.Printf("Default Stake: %s USDC", config.DefaultStakeUSDC)

	return service, nil
}

// CreateGame creates a new game on the blockchain and returns decoded result
func (s *EthereumBlockchainService) CreateGame(stakeAmount *big.Int) (*GameCreationResult, error) {
	// Use provided stake amount or default
	if stakeAmount == nil {
		stakeAmount = s.defaultStakeUSDC
	}

	// Get fresh nonce
	nonce, err := s.client.PendingNonceAt(context.Background(), s.auth.From)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}
	s.auth.Nonce = big.NewInt(int64(nonce))

	// Create game on contract
	tx, err := s.gameFactoryContract.CreateGame(s.auth, stakeAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to create game on-chain: %v", err)
	}

	log.Printf("Game creation transaction sent! Hash: %s", tx.Hash().Hex())

	// Wait for transaction to be mined and decode the result
	result, err := s.decodeGameCreationTransaction(tx.Hash())
	if err != nil {
		log.Printf("Warning: Failed to decode game creation result: %v", err)
		return nil, err
	}

	log.Printf("Game created successfully!")
	log.Printf("  Game ID: %s", result.GameID.String())
	log.Printf("  Game Contract: %s", result.GameContractAddress.Hex())
	log.Printf("  Transaction: %s", result.TransactionHash.Hex())
	log.Printf("  Block: %d", result.BlockNumber)

	return result, nil
}

// decodeGameCreationTransaction waits for a transaction to be mined and decodes the GameCreated event
func (s *EthereumBlockchainService) decodeGameCreationTransaction(txHash common.Hash) (*GameCreationResult, error) {
	// Wait for transaction to be mined
	receipt, err := s.waitForTransaction(txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	// Check if transaction was successful
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed")
	}

	// Parse GameCreated event from logs
	for _, log := range receipt.Logs {
		// Check if this log is from our GameFactory contract
		if log.Address != s.gameFactoryContractAddr {
			continue
		}

		// Parse the GameCreated event
		event, err := s.gameFactoryContract.ParseGameCreated(*log)
		if err != nil {
			continue // Not a GameCreated event, skip
		}

		// Return the decoded result
		return &GameCreationResult{
			GameID:              event.GameId,
			GameContractAddress: event.GameContract,
			TransactionHash:     txHash,
			BlockNumber:         receipt.BlockNumber.Uint64(),
		}, nil
	}

	return nil, fmt.Errorf("GameCreated event not found in transaction logs")
}

// waitForTransaction waits for a transaction to be mined and returns the receipt
func (s *EthereumBlockchainService) waitForTransaction(txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Cast client to ethclient.Client to access TransactionReceipt
	ethClient, ok := s.client.(*ethclient.Client)
	if !ok {
		return nil, fmt.Errorf("client is not an ethclient.Client")
	}

	// Poll for the transaction receipt
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for transaction to be mined")
		default:
			receipt, err := ethClient.TransactionReceipt(context.Background(), txHash)
			if err == nil {
				return receipt, nil
			}

			// If error is "not found", continue polling
			if strings.Contains(err.Error(), "not found") {
				time.Sleep(1 * time.Second)
				continue
			}

			// Other errors are actual failures
			return nil, err
		}
	}
}

// GetGameCreationResult is a helper method to decode any GameFactory transaction
func (s *EthereumBlockchainService) GetGameCreationResult(txHash common.Hash) (*GameCreationResult, error) {
	return s.decodeGameCreationTransaction(txHash)
}

// DecodeGameFactoryTransaction is a standalone method to decode GameFactory createGame transactions
// This method can be called with any transaction hash to extract gameId and gameContractAddress
func (s *EthereumBlockchainService) DecodeGameFactoryTransaction(txHashStr string) (*GameCreationResult, error) {
	// Parse transaction hash
	txHash := common.HexToHash(txHashStr)

	// Cast client to ethclient.Client to access TransactionReceipt
	ethClient, ok := s.client.(*ethclient.Client)
	if !ok {
		return nil, fmt.Errorf("client is not an ethclient.Client")
	}

	// Get transaction receipt
	receipt, err := ethClient.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	// Check if transaction was successful
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed")
	}

	// Parse GameCreated event from logs
	for _, log := range receipt.Logs {
		// Check if this log is from our GameFactory contract
		if log.Address != s.gameFactoryContractAddr {
			continue
		}

		// Parse the GameCreated event
		event, err := s.gameFactoryContract.ParseGameCreated(*log)
		if err != nil {
			continue // Not a GameCreated event, skip
		}

		// Return the decoded result
		return &GameCreationResult{
			GameID:              event.GameId,
			GameContractAddress: event.GameContract,
			TransactionHash:     txHash,
			BlockNumber:         receipt.BlockNumber.Uint64(),
		}, nil
	}

	return nil, fmt.Errorf("GameCreated event not found in transaction logs")
}

// DecodeTransactionByHash is a convenience method that takes a string hash
func DecodeGameFactoryTransactionByHash(client *ethclient.Client, gameFactoryAddr common.Address, txHashStr string) (*GameCreationResult, error) {
	// Parse transaction hash
	txHash := common.HexToHash(txHashStr)

	// Get transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	// Check if transaction was successful
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed")
	}

	// Create GameFactory contract instance for parsing
	gameFactory, err := gamefactory.NewGamefactory(gameFactoryAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create GameFactory instance: %v", err)
	}

	// Parse GameCreated event from logs
	for _, log := range receipt.Logs {
		// Check if this log is from the GameFactory contract
		if log.Address != gameFactoryAddr {
			continue
		}

		// Parse the GameCreated event
		event, err := gameFactory.ParseGameCreated(*log)
		if err != nil {
			continue // Not a GameCreated event, skip
		}

		// Return the decoded result
		return &GameCreationResult{
			GameID:              event.GameId,
			GameContractAddress: event.GameContract,
			TransactionHash:     txHash,
			BlockNumber:         receipt.BlockNumber.Uint64(),
		}, nil
	}

	return nil, fmt.Errorf("GameCreated event not found in transaction logs")
}

// RecordMove records a move on the blockchain
func (s *EthereumBlockchainService) RecordMove(gameID string, player common.Address, team uint8) error {
	// Get fresh nonce
	nonce, err := s.client.PendingNonceAt(context.Background(), s.auth.From)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}
	s.auth.Nonce = big.NewInt(int64(nonce))

	ethclient, ok := s.client.(*ethclient.Client)
	if !ok {
		return fmt.Errorf("client is not an ethclient.Client")
	}
	chainID, err := ethclient.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %v", err)
	}
	chainIDUint32 := uint32(chainID.Uint64())

	// Record move on contract
	tx, err := s.gameContract.AddVote(s.auth, player, chainIDUint32, team)
	if err != nil {
		return fmt.Errorf("failed to record move on-chain: %v", err)
	}

	log.Printf("Move recorded for game %s, player %s! Transaction: %s", gameID, player.Hex(), tx.Hash().Hex())
	return nil
}

// EndGame ends a game on the blockchain
func (s *EthereumBlockchainService) EndGame(gameID string, result GameResult) error {
	// Get fresh nonce
	nonce, err := s.client.PendingNonceAt(context.Background(), s.auth.From)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}
	s.auth.Nonce = big.NewInt(int64(nonce))

	// End game on contract
	tx, err := s.gameContract.EndGame(s.auth, uint8(result))
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
	// Get game info from contract
	gameInfo, err := s.gameContract.GetGameInfo(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get game info: %v", err)
	}

	// Check if game is active based on state
	isActive := gameInfo.State == 1 // Assuming 1 = Active state

	return &GameInfo{
		GameID:           gameInfo.GameId,
		FixedStakeAmount: big.NewInt(0), // Not available in new contract
		TotalWhiteStakes: big.NewInt(0), // Not available in new contract
		TotalBlackStakes: big.NewInt(0), // Not available in new contract
		WhitePlayerCount: big.NewInt(0), // Not available in new contract
		BlackPlayerCount: big.NewInt(0), // Not available in new contract
		CreatedAt:        big.NewInt(0), // Not available in new contract
		EndedAt:          big.NewInt(0), // Not available in new contract
		IsActive:         isActive,
		Result:           GameResult(gameInfo.Result),
	}, nil
}

// StakeOnVote stakes USDC from player's approved allowance
func (s *EthereumBlockchainService) StakeOnVote(gameID string, player common.Address, stakeAmount *big.Int) error {
	gameIDInt, err := s.parseGameID(gameID)
	if err != nil {
		return fmt.Errorf("invalid game ID: %v", err)
	}

	// Use default stake amount if not provided
	if stakeAmount == nil {
		stakeAmount = s.defaultStakeUSDC
	}

	// Get fresh nonce for backend transaction
	nonce, err := s.client.PendingNonceAt(context.Background(), s.auth.From)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}
	s.auth.Nonce = big.NewInt(int64(nonce))

	// Call stakeOnBehalfOf instead of stake
	tx, err := s.vaultContract.Stake(s.auth, player, big.NewInt(int64(gameIDInt)), stakeAmount)
	if err != nil {
		return fmt.Errorf("failed to stake on behalf of player: %v", err)
	}

	log.Printf("Stake deducted from player %s for game %s, amount %s USDC! Transaction: %s",
		player.Hex(), gameID, stakeAmount.String(), tx.Hash().Hex())
	return nil
}

// parseGameID extracts the numeric ID from a game ID string
func (s *EthereumBlockchainService) parseGameID(gameID string) (uint64, error) {
	if strings.HasPrefix(gameID, "game-") {
		idStr := strings.TrimPrefix(gameID, "game-")
		return strconv.ParseUint(idStr, 10, 64)
	}
	return strconv.ParseUint(gameID, 10, 64)
}

// GenerateUSDCApprovalTransaction generates transaction data for the player to sign
// This allows the vault contract to spend USDC on behalf of the player
func (s *EthereumBlockchainService) GenerateUSDCApprovalTransaction(playerAddress common.Address, approvalAmount *big.Int) (*ApprovalTransactionData, error) {
	// ERC20 approve function signature: approve(address spender, uint256 amount)
	// Function selector: 0x095ea7b3

	// Encode the approve function call
	// approve(vaultContractAddress, approvalAmount)
	approveData := make([]byte, 4+32+32) // 4 bytes selector + 32 bytes address + 32 bytes amount

	// Function selector for approve(address,uint256)
	copy(approveData[0:4], []byte{0x09, 0x5e, 0xa7, 0xb3})

	// Vault contract address (padded to 32 bytes)
	vaultAddressBytes := s.vaultContractAddr.Bytes()
	copy(approveData[4+32-len(vaultAddressBytes):4+32], vaultAddressBytes)

	// Approval amount (padded to 32 bytes)
	amountBytes := approvalAmount.Bytes()
	copy(approveData[4+32+32-len(amountBytes):4+32+32], amountBytes)

	// Get current gas price
	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	// Estimate gas for the approve transaction
	gasLimit := uint64(60000) // Standard gas limit for ERC20 approve

	return &ApprovalTransactionData{
		To:       USDCContractAddress,
		Data:     "0x" + hex.EncodeToString(approveData),
		Value:    "0x0",
		GasLimit: "0x" + strconv.FormatUint(gasLimit, 16),
		GasPrice: "0x" + gasPrice.Text(16),
	}, nil
}

// GenerateGameParticipationTransaction generates transaction data for game participation approval
func (s *EthereumBlockchainService) GenerateGameParticipationTransaction(playerAddress common.Address, gameID string, maxStakeAmount *big.Int) (*ApprovalTransactionData, error) {
	gameIDInt, err := s.parseGameID(gameID)
	if err != nil {
		return nil, fmt.Errorf("invalid game ID: %v", err)
	}

	// approveGameParticipation(uint256 gameId, uint256 maxStakeAmount)
	// You'll need to get the function selector from your contract ABI
	// For now, I'll use a placeholder - you'll need to replace this with the actual selector
	approveGameData := make([]byte, 4+32+32) // 4 bytes selector + 32 bytes gameId + 32 bytes amount

	// Function selector for approveGameParticipation(uint256,uint256)
	// You'll need to calculate this from your contract ABI
	copy(approveGameData[0:4], []byte{0x00, 0x00, 0x00, 0x00}) // Replace with actual selector

	// Game ID (padded to 32 bytes)
	gameIDBytes := big.NewInt(int64(gameIDInt)).Bytes()
	copy(approveGameData[4+32-len(gameIDBytes):4+32], gameIDBytes)

	// Max stake amount (padded to 32 bytes)
	amountBytes := maxStakeAmount.Bytes()
	copy(approveGameData[4+32+32-len(amountBytes):4+32+32], amountBytes)

	// Get current gas price
	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	// Estimate gas for the game participation transaction
	gasLimit := uint64(100000) // Estimated gas limit

	return &ApprovalTransactionData{
		To:       s.vaultContractAddr.Hex(),
		Data:     "0x" + hex.EncodeToString(approveGameData),
		Value:    "0x0",
		GasLimit: "0x" + strconv.FormatUint(gasLimit, 16),
		GasPrice: "0x" + gasPrice.Text(16),
	}, nil
}

// RequestPlayerApproval is a simple function to generate both approval transactions
// The frontend will use this to get transaction data for MetaMask signing
func (s *EthereumBlockchainService) RequestPlayerApproval(playerAddress common.Address, gameID string, approvalAmount *big.Int) (map[string]*ApprovalTransactionData, error) {
	log.Printf("Generating approval transactions for player %s, game %s, amount %s USDC",
		playerAddress.Hex(), gameID, approvalAmount.String())

	// Generate USDC approval transaction
	usdcApproval, err := s.GenerateUSDCApprovalTransaction(playerAddress, approvalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to generate USDC approval: %v", err)
	}

	// Generate game participation transaction
	gameApproval, err := s.GenerateGameParticipationTransaction(playerAddress, gameID, approvalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to generate game approval: %v", err)
	}

	return map[string]*ApprovalTransactionData{
		"usdcApproval": usdcApproval,
		"gameApproval": gameApproval,
	}, nil
}

// Generate Permit2 signature data for allowance-based permits
func (s *EthereumBlockchainService) GeneratePermit2SignatureData(playerAddress common.Address) (*apitypes.TypedData, error) {
	// Use timestamp-based nonce for uniqueness
	nonce := uint64(time.Now().UnixNano())

	// Set expiration (1 hour from now)
	expiration := uint64(time.Now().Add(1 * time.Hour).Unix())

	// Set signature deadline (1 hour from now)
	sigDeadline := uint64(time.Now().Add(1 * time.Hour).Unix())

	// Max uint160 for unlimited approval (Permit2 uses uint160 for amounts)
	maxAmount := "1461501637330902918203684832716283019655932542975" // 2^160 - 1

	// Store the permit data for later use during execution
	maxAmountBig := new(big.Int)
	maxAmountBig.SetString(maxAmount, 10)

	s.permits[playerAddress] = &PermitData{
		Nonce:       big.NewInt(int64(nonce)),
		Expiration:  big.NewInt(int64(expiration)),
		SigDeadline: big.NewInt(int64(sigDeadline)),
		Amount:      maxAmountBig,
		Signature:   "", // Will be set when signature is received
	}

	ethclient, ok := s.client.(*ethclient.Client)
	if !ok {
		return nil, fmt.Errorf("client is not an ethclient.Client")
	}
	chainID, err := ethclient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}
	chainIDUint32 := uint32(chainID.Uint64())

	// Create EIP-712 typed data for Permit2 allowance-based permit
	typedData := &apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"PermitSingle": []apitypes.Type{
				{Name: "details", Type: "PermitDetails"},
				{Name: "spender", Type: "address"},
				{Name: "sigDeadline", Type: "uint256"},
			},
			"PermitDetails": []apitypes.Type{
				{Name: "token", Type: "address"},
				{Name: "amount", Type: "uint160"},
				{Name: "expiration", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
			},
		},
		PrimaryType: "PermitSingle",
		Domain: apitypes.TypedDataDomain{
			Name:              "Permit2",
			ChainId:           (*math.HexOrDecimal256)(chainID),
			VerifyingContract: Permit2ContractAddressByChain[chainIDUint32],
		},
		Message: apitypes.TypedDataMessage{
			"details": map[string]interface{}{
				"token":      USDCContractAddress,
				"amount":     maxAmount,
				"expiration": expiration,
				"nonce":      nonce,
			},
			"spender":     s.vaultContractAddr.Hex(),
			"sigDeadline": sigDeadline,
		},
	}

	return typedData, nil
}

// Request Permit2 signature (no amount needed!)
func (s *EthereumBlockchainService) RequestPermit2(playerAddress common.Address) (*apitypes.TypedData, error) {
	return s.GeneratePermit2SignatureData(playerAddress)
}

// Add this method after getUSDCContract
func (s *EthereumBlockchainService) getPermit2Contract() (*permit2.Permit2, error) {
	ethclient, ok := s.client.(*ethclient.Client)
	if !ok {
		return nil, fmt.Errorf("client is not an ethclient.Client")
	}
	chainID, err := ethclient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}
	chainIDUint32 := uint32(chainID.Uint64())
	permit2ContractAddress := Permit2ContractAddressByChain[chainIDUint32]
	return permit2.NewPermit2(common.HexToAddress(permit2ContractAddress), s.client)
}

// StorePermit2Signature stores a signed Permit2 signature for later use
func (s *EthereumBlockchainService) StorePermit2Signature(playerAddress common.Address, signature string) error {
	// Get the existing permit data
	permitData, exists := s.permits[playerAddress]
	if !exists {
		return fmt.Errorf("no permit data found for player %s - must call GeneratePermit2SignatureData first", playerAddress.Hex())
	}

	// Update the signature
	permitData.Signature = signature

	log.Printf("Stored Permit2 signature for player %s", playerAddress.Hex())
	return nil
}

// PermitData stores the permit data used for signature generation
type PermitData struct {
	Nonce       *big.Int
	Expiration  *big.Int
	SigDeadline *big.Int
	Amount      *big.Int
	Signature   string
}

// ExecutePermit2 executes a Permit2 allowance using the stored signature
func (s *EthereumBlockchainService) ExecutePermit2(playerAddress common.Address) error {
	// Get the stored permit data
	permitData, exists := s.permits[playerAddress]
	if !exists {
		return fmt.Errorf("no permit signature found for player %s", playerAddress.Hex())
	}

	// Get Permit2 contract
	permit2Contract, err := s.getPermit2Contract()
	if err != nil {
		return fmt.Errorf("failed to get Permit2 contract: %v", err)
	}

	// Convert hex signature to bytes
	signatureBytes, err := hex.DecodeString(strings.TrimPrefix(permitData.Signature, "0x"))
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}

	// Create permit details using the stored values
	permitDetails := permit2.IAllowanceTransferPermitDetails{
		Token:      common.HexToAddress(USDCContractAddress),
		Amount:     permitData.Amount,
		Expiration: permitData.Expiration,
		Nonce:      permitData.Nonce,
	}

	// Create permit single struct
	permitSingle := permit2.IAllowanceTransferPermitSingle{
		Details:     permitDetails,
		Spender:     s.vaultContractAddr,
		SigDeadline: permitData.SigDeadline,
	}

	// Execute the permit using Permit0 (single permit function)
	tx, err := permit2Contract.Permit0(s.auth, playerAddress, permitSingle, signatureBytes)
	if err != nil {
		return fmt.Errorf("failed to execute permit: %v", err)
	}

	log.Printf("Permit2 executed for player %s, tx hash: %s", playerAddress.Hex(), tx.Hash().Hex())
	return nil
}

func (s *EthereumBlockchainService) TransferRewardsCrossChain(gameID string, amount *big.Int, destinationChainId *big.Int, recipient common.Address, useFastTransfer bool, maxFee *big.Int) error {
	gameIDInt, err := s.parseGameID(gameID)
	if err != nil {
		return fmt.Errorf("invalid game ID: %v", err)
	}

	tx, err := s.vaultContract.TransferRewardsCrossChain(s.auth, big.NewInt(int64(gameIDInt)), amount, destinationChainId, recipient, useFastTransfer, maxFee)
	if err != nil {
		return fmt.Errorf("failed to transfer rewards cross chain: %v", err)
	}

	log.Printf("Rewards transferred cross chain for game %s, tx hash: %s", gameID, tx.Hash().Hex())
	return nil
}
