package client

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// ChainConfig represents configuration for a specific blockchain
type ChainConfig struct {
	RPCUrl               string
	VaultContractAddress string
	Permit2Address       string // Add Permit2 contract address
	Name                 string
	EnvVaultKey          string // Environment variable key for vault address
}

// Clients holds all blockchain clients for different chains
type Clients struct {
	clients    map[uint64]*ethclient.Client
	privateKey string // Store the private key
	mu         sync.RWMutex
}

// NewClients creates a new Clients instance
func NewClients() *Clients {
	return &Clients{
		clients: make(map[uint64]*ethclient.Client),
	}
}

// supportedChains defines all supported blockchain configurations
var supportedChains = map[uint64]ChainConfig{
	// Ethereum Sepolia
	11155111: {
		RPCUrl:      "https://eth-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7",
		Name:        "Ethereum Sepolia",
		EnvVaultKey: "ETHEREUM_SEPOLIA_VAULT_ADDRESS",
	},
	// Avalanche Fuji
	43113: {
		RPCUrl:      "https://avax-fuji.g.alchemy.com/v2/22r8dairb21cjlw7",
		Name:        "Avalanche Fuji",
		EnvVaultKey: "AVALANCHE_FUJI_VAULT_ADDRESS",
	},
	// OP Sepolia
	11155420: {
		RPCUrl:      "https://opt-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7",
		Name:        "OP Sepolia",
		EnvVaultKey: "OP_SEPOLIA_VAULT_ADDRESS",
	},
	// Arbitrum Sepolia
	421614: {
		RPCUrl:      "https://arb-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7",
		Name:        "Arbitrum Sepolia",
		EnvVaultKey: "ARBITRUM_SEPOLIA_VAULT_ADDRESS",
	},
	// Base Sepolia
	84532: {
		RPCUrl:      "https://base-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7",
		Name:        "Base Sepolia",
		EnvVaultKey: "BASE_SEPOLIA_VAULT_ADDRESS",
	},
	// Polygon Amoy
	80002: {
		RPCUrl:      "https://polygon-amoy.g.alchemy.com/v2/22r8dairb21cjlw7",
		Name:        "Polygon Amoy",
		EnvVaultKey: "POLYGON_AMOY_VAULT_ADDRESS",
	},
	// Unichain Sepolia
	1301: {
		RPCUrl:      "https://eth-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7", // Placeholder
		Name:        "Unichain Sepolia",
		EnvVaultKey: "UNICHAIN_SEPOLIA_VAULT_ADDRESS",
	},
	// Linea Sepolia
	59141: {
		RPCUrl:      "https://linea-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7",
		Name:        "Linea Sepolia",
		EnvVaultKey: "LINEA_SEPOLIA_VAULT_ADDRESS",
	},
	// Sonic Testnet
	64165: {
		RPCUrl:      "https://sonic-blaze.g.alchemy.com/v2/Kv5V9_Sv55ZosXpoVga3IfXYO6JV9gNJ",
		Name:        "Sonic Testnet",
		EnvVaultKey: "SONIC_TESTNET_VAULT_ADDRESS",
	},
	// World Chain Sepolia
	4801: {
		RPCUrl:      "https://eth-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7", // Placeholder
		Name:        "World Chain Sepolia",
		EnvVaultKey: "WORLD_CHAIN_SEPOLIA_VAULT_ADDRESS",
	},
	// Codex Testnet
	325000: {
		RPCUrl:      "https://eth-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7", // Placeholder
		Name:        "Codex Testnet",
		EnvVaultKey: "CODEX_TESTNET_VAULT_ADDRESS",
	},
	// Local Anvil (for development)
	31337: {
		RPCUrl:      "http://127.0.0.1:8545",
		Name:        "Anvil Local (Base Sepolia)",
		EnvVaultKey: "ANVIL_BASE_SEPOLIA_VAULT_ADDRESS",
	},
	// Local Anvil (for development)
	31338: {
		RPCUrl:      "http://127.0.0.1:8546",
		Name:        "Anvil Local (Optimism Sepolia)",
		EnvVaultKey: "ANVIL_OPTIMISM_SEPOLIA_VAULT_ADDRESS",
	},
}

// envLoader ensures .env is loaded only once
var envLoader sync.Once

// loadEnv loads environment variables once
func loadEnv() {
	envLoader.Do(func() {
		godotenv.Load(".env")
	})
}

// InitializeClients creates and initializes all blockchain clients
func InitializeClients() (*Clients, error) {
	clients := NewClients()

	// Store the pre-loaded private key
	if PrivateKey == "" {
		log.Fatalf("PRIVATE_KEY not set in environment")
	} else {
		clients.setPrivateKey(PrivateKey)
	}

	for chainID, config := range supportedChains {
		client, err := NewClient(chainID)
		if err != nil {
			log.Printf("Warning: Failed to initialize %s client: %v", config.Name, err)
			continue
		}
		clients.setClient(chainID, client)
	}

	log.Printf("Blockchain clients initialized successfully")
	return clients, nil
}

// setClient safely sets a client for a chain ID
func (c *Clients) setClient(chainID uint64, client *ethclient.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[chainID] = client
}

// setPrivateKey safely sets the private key
func (c *Clients) setPrivateKey(privateKey string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.privateKey = privateKey
}

// GetPrivateKey safely gets the private key
func (c *Clients) GetPrivateKey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.privateKey
}

// GetClientByChainID returns the appropriate client for a given chain ID
func (c *Clients) GetClientByChainID(chainID uint64) (*ethclient.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	client, exists := c.clients[chainID]
	if !exists || client == nil {
		config, configExists := supportedChains[chainID]
		if !configExists {
			return nil, fmt.Errorf("unsupported chain ID: %d", chainID)
		}
		return nil, fmt.Errorf("%s client not initialized", config.Name)
	}
	return client, nil
}

// Close closes all blockchain client connections
func (c *Clients) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for chainID, client := range c.clients {
		if client != nil {
			client.Close()
			delete(c.clients, chainID)
		}
	}
}

// LoadChainConfigs loads chain configurations from environment variables
func LoadChainConfigs() map[uint64]ChainConfig {
	loadEnv()

	configs := make(map[uint64]ChainConfig)
	for chainID, baseConfig := range supportedChains {
		config := baseConfig
		config.VaultContractAddress = os.Getenv(baseConfig.EnvVaultKey)
		config.Permit2Address = GetPermit2Address(chainID)
		configs[chainID] = config
	}

	return configs
}

// ChainConfigs holds the loaded chain configurations
var ChainConfigs = LoadChainConfigs()

// PrivateKey holds the loaded private key
var PrivateKey = LoadPrivateKey()

// GameFactoryAddress holds the loaded GameFactory contract address
var GameFactoryAddress = LoadGameFactoryAddress()

// VaultAddresses holds the loaded vault contract addresses for all chains
var VaultAddresses = LoadVaultAddresses()

// Permit2Addresses holds the loaded Permit2 contract addresses for all chains
var Permit2Addresses = LoadPermit2Addresses()

// LoadPrivateKey loads the private key from environment variables
func LoadPrivateKey() string {
	loadEnv()
	return os.Getenv("PRIVATE_KEY")
}

// LoadGameFactoryAddress loads the GameFactory contract address from environment variables
func LoadGameFactoryAddress() string {
	loadEnv()
	return os.Getenv("GAME_FACTORY_ADDRESS")
}

// LoadVaultAddresses loads all vault contract addresses from environment variables
func LoadVaultAddresses() map[uint64]string {
	loadEnv()

	vaultAddresses := make(map[uint64]string)
	for chainID, config := range supportedChains {
		address := os.Getenv(config.EnvVaultKey)
		if address != "" {
			vaultAddresses[chainID] = address
		}
	}

	return vaultAddresses
}

// LoadPermit2Addresses loads all Permit2 contract addresses from environment variables
func LoadPermit2Addresses() map[uint64]string {
	loadEnv()

	// Permit2 is deployed at the same address on all chains
	permit2Address := os.Getenv("PERMIT2_ADDRESS")
	if permit2Address == "" {
		// Use the canonical Permit2 address if not specified
		permit2Address = "0x000000000022D473030F116dDEE9F6B43aC78BA3"
	}

	permit2Addresses := make(map[uint64]string)
	for chainID := range supportedChains {
		permit2Addresses[chainID] = permit2Address
	}

	return permit2Addresses
}

// GetPrivateKey returns the private key from environment variables
func GetPrivateKey() string {
	return PrivateKey
}

// GetGameFactoryAddress returns the GameFactory contract address (deployed only on Base Sepolia)
func GetGameFactoryAddress() string {
	return GameFactoryAddress
}

// GetVaultAddress returns the vault contract address for a specific chain ID
func GetVaultAddress(chainID uint64) string {
	return VaultAddresses[chainID]
}

// GetPermit2Address returns the Permit2 contract address for a specific chain ID
func GetPermit2Address(chainID uint64) string {
	return Permit2Addresses[chainID]
}

// GetEnv gets an environment variable with a default value
func GetEnv(key, defaultValue string) string {
	loadEnv()
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt gets an environment variable as integer with a default value
func GetEnvInt(key string, defaultValue int) int {
	loadEnv()
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvRequired gets a required environment variable
func GetEnvRequired(key string) (string, error) {
	loadEnv()
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s not set", key)
	}
	return value, nil
}

// NewClient creates a new Ethereum client for the specified chain ID
func NewClient(chainID uint64) (*ethclient.Client, error) {
	config, exists := supportedChains[chainID]
	if !exists {
		return nil, fmt.Errorf("unsupported chain ID: %d", chainID)
	}

	client, err := ethclient.Dial(config.RPCUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", config.Name, err)
	}

	return client, nil
}

// GetChainConfig returns the configuration for a specific chain ID
func GetChainConfig(chainID uint64) (ChainConfig, error) {
	config, exists := ChainConfigs[chainID]
	if !exists {
		return ChainConfig{}, fmt.Errorf("unsupported chain ID: %d", chainID)
	}
	return config, nil
}

// UpdateVaultAddress updates the vault contract address for a specific chain
func UpdateVaultAddress(chainID uint64, vaultAddress string) error {
	config, exists := ChainConfigs[chainID]
	if !exists {
		return fmt.Errorf("unsupported chain ID: %d", chainID)
	}

	config.VaultContractAddress = vaultAddress
	ChainConfigs[chainID] = config
	return nil
}

// GetSupportedChains returns a list of all supported chain IDs
func GetSupportedChains() []uint64 {
	chains := make([]uint64, 0, len(supportedChains))
	for chainID := range supportedChains {
		chains = append(chains, chainID)
	}
	return chains
}

// ReloadChainConfigs reloads chain configurations from environment variables
func ReloadChainConfigs() {
	ChainConfigs = LoadChainConfigs()
}
