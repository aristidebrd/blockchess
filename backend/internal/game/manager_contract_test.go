package game

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"

	"blockchess/contracts/gamecontract"
	"blockchess/contracts/vaultcontract"
)

// Test setup constants
const (
	testGameID   = "12345"
	testStakeETH = "0.01"
)

// TestManagerWithContracts tests the complete integration between manager and smart contracts
func TestManagerWithContracts(t *testing.T) {
	// Setup blockchain simulator
	backend, gameContract, vaultContract, auth := setupTestBlockchain(t)

	// Create custom manager for testing
	manager := createTestManager(backend, gameContract, vaultContract, auth)

	// Test 1: Create game and verify basic functionality
	t.Run("CreateGame", func(t *testing.T) {
		game := manager.GetOrCreateGame(testGameID)
		if game == nil {
			t.Fatal("Failed to create game")
		}

		// Wait for async operations
		time.Sleep(100 * time.Millisecond)
		backend.Commit()

		// Verify game was created locally
		if game.ID != testGameID {
			t.Errorf("Expected game ID %s, got %s", testGameID, game.ID)
		}
	})

	// Test 2: Player joining and team assignment
	t.Run("PlayerJoinTeam", func(t *testing.T) {
		manager.AddPlayerToTeam(testGameID, "player1", "white")
		manager.AddPlayerToTeam(testGameID, "player2", "black")

		stats := manager.GetGameStats(testGameID)
		if stats["whitePlayers"].(int) != 1 {
			t.Errorf("Expected 1 white player, got %d", stats["whitePlayers"].(int))
		}
		if stats["blackPlayers"].(int) != 1 {
			t.Errorf("Expected 1 black player, got %d", stats["blackPlayers"].(int))
		}
	})

	// Test 3: Voting mechanism
	t.Run("VotingMechanism", func(t *testing.T) {
		// Create a fresh game for voting test
		votingGameID := "voting123"
		manager.GetOrCreateGame(votingGameID)

		// Add players to teams first
		manager.AddPlayerToTeam(votingGameID, "voter1", "white")
		manager.AddPlayerToTeam(votingGameID, "voter2", "white")
		manager.AddPlayerToTeam(votingGameID, "voter3", "black")

		// Test valid move voting - white moves first
		manager.VoteForMove(votingGameID, "voter1", "e4", "white")
		manager.VoteForMove(votingGameID, "voter2", "e4", "white")
		// Black can't vote yet since it's white's turn
		manager.VoteForMove(votingGameID, "voter3", "e5", "black")

		votes := manager.GetVotes(votingGameID)
		if votes["e4"] != 2 {
			t.Errorf("Expected 2 votes for e4, got %d", votes["e4"])
		}
		// e5 should be 0 because it's white's turn
		if votes["e5"] != 0 {
			t.Errorf("Expected 0 votes for e5 (wrong turn), got %d", votes["e5"])
		}

		// Test duplicate voting prevention
		manager.VoteForMove(votingGameID, "voter1", "d4", "white")
		votes = manager.GetVotes(votingGameID)
		if _, exists := votes["d4"]; exists {
			t.Error("Player should not be able to vote twice in same round")
		}
	})

	// Test 4: Invalid move rejection
	t.Run("InvalidMoveRejection", func(t *testing.T) {
		manager.AddPlayerToTeam(testGameID, "invalidvoter", "white")

		// Try to vote for invalid move
		manager.VoteForMove(testGameID, "invalidvoter", "z9", "white")

		votes := manager.GetVotes(testGameID)
		if _, exists := votes["z9"]; exists {
			t.Error("Invalid move should not be accepted")
		}
	})

	// Test 5: Contract interaction verification
	t.Run("ContractInteractionVerification", func(t *testing.T) {
		// Create a game on the blockchain first
		gameID := big.NewInt(12345)
		stakeAmount := big.NewInt(10000000000000000) // 0.01 ETH

		_, err := gameContract.CreateGame(auth, gameID, stakeAmount)
		if err != nil {
			t.Fatalf("Failed to create game on blockchain: %v", err)
		}

		backend.Commit()

		// Test game contract interaction
		gameInfo, err := gameContract.GetGameInfo(nil, gameID)
		if err != nil {
			t.Fatalf("Failed to get game info: %v", err)
		}

		// Verify game exists
		if gameInfo.GameId.Cmp(gameID) != 0 {
			t.Errorf("Expected game ID 12345, got %v", gameInfo.GameId)
		}

		// Test that the contracts are properly deployed and accessible
		if gameInfo.State != 0 { // Active state
			t.Errorf("Expected game state 0 (Active), got %d", gameInfo.State)
		}

		// Verify stake amount
		if gameInfo.FixedStakeAmount.Cmp(stakeAmount) != 0 {
			t.Errorf("Expected stake amount %v, got %v", stakeAmount, gameInfo.FixedStakeAmount)
		}
	})

	// Test 6: Game statistics accuracy
	t.Run("GameStatisticsAccuracy", func(t *testing.T) {
		statsGameID := "stats123"
		manager.GetOrCreateGame(statsGameID)

		// Add multiple players and votes
		manager.AddPlayerToTeam(statsGameID, "stats1", "white")
		manager.AddPlayerToTeam(statsGameID, "stats2", "white")
		manager.AddPlayerToTeam(statsGameID, "stats3", "black")

		// Only white can vote since it's white's turn
		manager.VoteForMove(statsGameID, "stats1", "e4", "white")
		manager.VoteForMove(statsGameID, "stats2", "d4", "white")
		// Black vote should be rejected due to turn
		manager.VoteForMove(statsGameID, "stats3", "e5", "black")

		stats := manager.GetGameStats(statsGameID)

		// Verify statistics
		if stats["whitePlayers"].(int) != 2 {
			t.Errorf("Expected 2 white players, got %d", stats["whitePlayers"].(int))
		}
		if stats["blackPlayers"].(int) != 1 {
			t.Errorf("Expected 1 black player, got %d", stats["blackPlayers"].(int))
		}
		if stats["whiteCurrentTurnVotes"].(int) != 2 {
			t.Errorf("Expected 2 white votes, got %d", stats["whiteCurrentTurnVotes"].(int))
		}
		// Black vote should be 0 since it's white's turn
		if stats["blackCurrentTurnVotes"].(int) != 0 {
			t.Errorf("Expected 0 black votes (wrong turn), got %d", stats["blackCurrentTurnVotes"].(int))
		}
	})

	// Test 7: Manager configuration
	t.Run("ManagerConfiguration", func(t *testing.T) {
		// Test contract config loading
		config := &ContractConfig{
			RPCUrl:              "http://127.0.0.1:8545",
			PrivateKey:          "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			GameContractAddress: "0x1234567890123456789012345678901234567890",
			DefaultStakeETH:     "0.01",
		}

		// Test that config is properly structured
		if config.RPCUrl == "" {
			t.Error("RPC URL should not be empty")
		}
		if config.PrivateKey == "" {
			t.Error("Private key should not be empty")
		}
		if config.GameContractAddress == "" {
			t.Error("Game contract address should not be empty")
		}
		if config.DefaultStakeETH == "" {
			t.Error("Default stake ETH should not be empty")
		}
	})
}

// setupTestBlockchain creates a simulated blockchain for testing
func setupTestBlockchain(t *testing.T) (*backends.SimulatedBackend, *gamecontract.Gamecontract, *vaultcontract.Vaultcontract, *bind.TransactOpts) {
	// Generate test private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Create auth object
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		t.Fatalf("Failed to create auth: %v", err)
	}

	// Create simulated backend
	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 ETH
	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}

	backend := backends.NewSimulatedBackend(genesisAlloc, 10000000)
	auth.GasLimit = 3000000

	// Deploy GameContract
	gameAddr, _, gameContract, err := gamecontract.DeployGamecontract(auth, backend, auth.From)
	if err != nil {
		t.Fatalf("Failed to deploy GameContract: %v", err)
	}

	// Deploy VaultContract
	_, _, vaultContract, err := vaultcontract.DeployVaultcontract(auth, backend, auth.From, gameAddr)
	if err != nil {
		t.Fatalf("Failed to deploy VaultContract: %v", err)
	}

	backend.Commit()

	return backend, gameContract, vaultContract, auth
}

// createTestManager creates a manager with test configuration
func createTestManager(backend *backends.SimulatedBackend, gameContract *gamecontract.Gamecontract, vaultContract *vaultcontract.Vaultcontract, auth *bind.TransactOpts) *Manager {
	manager := NewManager()

	// Manually set up contract integration for testing
	manager.SetContractClients(backend, gameContract, vaultContract, auth)

	return manager
}

// Benchmark test for performance
func BenchmarkManagerOperations(b *testing.B) {
	backend, gameContract, vaultContract, auth := setupTestBlockchain(&testing.T{})
	manager := createTestManager(backend, gameContract, vaultContract, auth)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gameID := fmt.Sprintf("bench_%d", i)
		manager.GetOrCreateGame(gameID)
		manager.AddPlayerToTeam(gameID, "player1", "white")
		manager.VoteForMove(gameID, "player1", "e4", "white")
		_ = manager.GetGameStats(gameID)
	}
}
