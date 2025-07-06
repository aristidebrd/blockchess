package game

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// ChainConfig represents configuration for a specific blockchain
type ChainConfig struct {
	RPCUrl               string
	VaultContractAddress string
	Name                 string
}

// Clients holds all blockchain clients for different chains
type Clients struct {
	EthereumSepolia      *ethclient.Client
	AvalancheFuji        *ethclient.Client
	OpSepolia            *ethclient.Client
	ArbitrumSepolia      *ethclient.Client
	BaseSepolia          *ethclient.Client
	PolygonAmoy          *ethclient.Client
	UnichainSepolia      *ethclient.Client
	LineaSepolia         *ethclient.Client
	SonicTestnet         *ethclient.Client
	WorldChainSepolia    *ethclient.Client
	CodexTestnet         *ethclient.Client
	AnvilBaseSepolia     *ethclient.Client
	AnvilOptimismSepolia *ethclient.Client
}

// InitializeClients creates and initializes all blockchain clients
func InitializeClients() (*Clients, error) {
	clients := &Clients{}

	// Initialize each client with error handling
	var err error

	// Ethereum Sepolia
	clients.EthereumSepolia, err = NewClient(11155111)
	if err != nil {
		log.Printf("Warning: Failed to initialize Ethereum Sepolia client: %v", err)
	}

	// Avalanche Fuji
	clients.AvalancheFuji, err = NewClient(43113)
	if err != nil {
		log.Printf("Warning: Failed to initialize Avalanche Fuji client: %v", err)
	}

	// OP Sepolia
	clients.OpSepolia, err = NewClient(11155420)
	if err != nil {
		log.Printf("Warning: Failed to initialize OP Sepolia client: %v", err)
	}

	// Arbitrum Sepolia
	clients.ArbitrumSepolia, err = NewClient(421614)
	if err != nil {
		log.Printf("Warning: Failed to initialize Arbitrum Sepolia client: %v", err)
	}

	// Base Sepolia
	clients.BaseSepolia, err = NewClient(84532)
	if err != nil {
		log.Printf("Warning: Failed to initialize Base Sepolia client: %v", err)
	}

	// Polygon Amoy
	clients.PolygonAmoy, err = NewClient(80002)
	if err != nil {
		log.Printf("Warning: Failed to initialize Polygon Amoy client: %v", err)
	}

	// Unichain Sepolia
	clients.UnichainSepolia, err = NewClient(1301)
	if err != nil {
		log.Printf("Warning: Failed to initialize Unichain Sepolia client: %v", err)
	}

	// Linea Sepolia
	clients.LineaSepolia, err = NewClient(59141)
	if err != nil {
		log.Printf("Warning: Failed to initialize Linea Sepolia client: %v", err)
	}

	// Sonic Testnet
	clients.SonicTestnet, err = NewClient(64165)
	if err != nil {
		log.Printf("Warning: Failed to initialize Sonic Testnet client: %v", err)
	}

	// World Chain Sepolia
	clients.WorldChainSepolia, err = NewClient(4801)
	if err != nil {
		log.Printf("Warning: Failed to initialize World Chain Sepolia client: %v", err)
	}

	// Codex Testnet
	clients.CodexTestnet, err = NewClient(325000)
	if err != nil {
		log.Printf("Warning: Failed to initialize Codex Testnet client: %v", err)
	}

	// Anvil Base Sepolia (local development)
	clients.AnvilBaseSepolia, err = NewClient(31337)
	if err != nil {
		log.Printf("Warning: Failed to initialize Anvil Base Sepolia client: %v", err)
	}

	// Anvil Optimism Sepolia (local development)
	clients.AnvilOptimismSepolia, err = NewClient(31338)
	if err != nil {
		log.Printf("Warning: Failed to initialize Anvil Optimism Sepolia client: %v", err)
	}

	log.Printf("Blockchain clients initialized successfully")
	return clients, nil
}

// GetClientByChainID returns the appropriate client for a given chain ID
func (c *Clients) GetClientByChainID(chainID uint64) (*ethclient.Client, error) {
	switch chainID {
	case 11155111:
		if c.EthereumSepolia == nil {
			return nil, fmt.Errorf("ethereum Sepolia client not initialized")
		}
		return c.EthereumSepolia, nil
	case 43113:
		if c.AvalancheFuji == nil {
			return nil, fmt.Errorf("avalanche Fuji client not initialized")
		}
		return c.AvalancheFuji, nil
	case 11155420:
		if c.OpSepolia == nil {
			return nil, fmt.Errorf("op Sepolia client not initialized")
		}
		return c.OpSepolia, nil
	case 421614:
		if c.ArbitrumSepolia == nil {
			return nil, fmt.Errorf("arbitrum Sepolia client not initialized")
		}
		return c.ArbitrumSepolia, nil
	case 84532:
		if c.BaseSepolia == nil {
			return nil, fmt.Errorf("base Sepolia client not initialized")
		}
		return c.BaseSepolia, nil
	case 80002:
		if c.PolygonAmoy == nil {
			return nil, fmt.Errorf("polygon Amoy client not initialized")
		}
		return c.PolygonAmoy, nil
	case 1301:
		if c.UnichainSepolia == nil {
			return nil, fmt.Errorf("unichain Sepolia client not initialized")
		}
		return c.UnichainSepolia, nil
	case 59141:
		if c.LineaSepolia == nil {
			return nil, fmt.Errorf("linea Sepolia client not initialized")
		}
		return c.LineaSepolia, nil
	case 64165:
		if c.SonicTestnet == nil {
			return nil, fmt.Errorf("sonic Testnet client not initialized")
		}
		return c.SonicTestnet, nil
	case 4801:
		if c.WorldChainSepolia == nil {
			return nil, fmt.Errorf("world Chain Sepolia client not initialized")
		}
		return c.WorldChainSepolia, nil
	case 325000:
		if c.CodexTestnet == nil {
			return nil, fmt.Errorf("codex Testnet client not initialized")
		}
		return c.CodexTestnet, nil
	case 31337:
		if c.AnvilBaseSepolia == nil {
			return nil, fmt.Errorf("anvil Base Sepolia client not initialized")
		}
		return c.AnvilBaseSepolia, nil
	case 31338:
		if c.AnvilOptimismSepolia == nil {
			return nil, fmt.Errorf("anvil Optimism Sepolia client not initialized")
		}
		return c.AnvilOptimismSepolia, nil
	default:
		return nil, fmt.Errorf("unsupported chain ID: %d", chainID)
	}
}

// Close closes all blockchain client connections
func (c *Clients) Close() {
	if c.EthereumSepolia != nil {
		c.EthereumSepolia.Close()
	}
	if c.AvalancheFuji != nil {
		c.AvalancheFuji.Close()
	}
	if c.OpSepolia != nil {
		c.OpSepolia.Close()
	}
	if c.ArbitrumSepolia != nil {
		c.ArbitrumSepolia.Close()
	}
	if c.BaseSepolia != nil {
		c.BaseSepolia.Close()
	}
	if c.PolygonAmoy != nil {
		c.PolygonAmoy.Close()
	}
	if c.UnichainSepolia != nil {
		c.UnichainSepolia.Close()
	}
	if c.LineaSepolia != nil {
		c.LineaSepolia.Close()
	}
	if c.SonicTestnet != nil {
		c.SonicTestnet.Close()
	}
	if c.WorldChainSepolia != nil {
		c.WorldChainSepolia.Close()
	}
	if c.CodexTestnet != nil {
		c.CodexTestnet.Close()
	}
	if c.AnvilBaseSepolia != nil {
		c.AnvilBaseSepolia.Close()
	}
	if c.AnvilOptimismSepolia != nil {
		c.AnvilOptimismSepolia.Close()
	}
}

// LoadChainConfigs loads chain configurations from environment variables
func LoadChainConfigs() map[uint64]ChainConfig {
	// Load environment variables
	godotenv.Load(".env")

	return map[uint64]ChainConfig{
		// Ethereum Sepolia
		11155111: {
			RPCUrl:               "https://eth-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7",
			VaultContractAddress: os.Getenv("ETHEREUM_SEPOLIA_VAULT_ADDRESS"),
			Name:                 "Ethereum Sepolia",
		},
		// Avalanche Fuji
		43113: {
			RPCUrl:               "https://avax-fuji.g.alchemy.com/v2/22r8dairb21cjlw7",
			VaultContractAddress: os.Getenv("AVALANCHE_FUJI_VAULT_ADDRESS"),
			Name:                 "Avalanche Fuji",
		},
		// OP Sepolia
		11155420: {
			RPCUrl:               "https://opt-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7",
			VaultContractAddress: os.Getenv("OP_SEPOLIA_VAULT_ADDRESS"),
			Name:                 "OP Sepolia",
		},
		// Arbitrum Sepolia
		421614: {
			RPCUrl:               "https://arb-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7",
			VaultContractAddress: os.Getenv("ARBITRUM_SEPOLIA_VAULT_ADDRESS"),
			Name:                 "Arbitrum Sepolia",
		},
		// Base Sepolia
		84532: {
			RPCUrl:               "https://base-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7",
			VaultContractAddress: os.Getenv("BASE_SEPOLIA_VAULT_ADDRESS"),
			Name:                 "Base Sepolia",
		},
		// Polygon Amoy
		80002: {
			RPCUrl:               "https://polygon-amoy.g.alchemy.com/v2/22r8dairb21cjlw7",
			VaultContractAddress: os.Getenv("POLYGON_AMOY_VAULT_ADDRESS"),
			Name:                 "Polygon Amoy",
		},
		// Unichain Sepolia
		1301: {
			RPCUrl:               "https://eth-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7", // Placeholder
			VaultContractAddress: os.Getenv("UNICHAIN_SEPOLIA_VAULT_ADDRESS"),
			Name:                 "Unichain Sepolia",
		},
		// Linea Sepolia
		59141: {
			RPCUrl:               "https://linea-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7",
			VaultContractAddress: os.Getenv("LINEA_SEPOLIA_VAULT_ADDRESS"),
			Name:                 "Linea Sepolia",
		},
		// Sonic Testnet
		64165: {
			RPCUrl:               "https://sonic-blaze.g.alchemy.com/v2/Kv5V9_Sv55ZosXpoVga3IfXYO6JV9gNJ",
			VaultContractAddress: os.Getenv("SONIC_TESTNET_VAULT_ADDRESS"),
			Name:                 "Sonic Testnet",
		},
		// World Chain Sepolia
		4801: {
			RPCUrl:               "https://eth-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7", // Placeholder
			VaultContractAddress: os.Getenv("WORLD_CHAIN_SEPOLIA_VAULT_ADDRESS"),
			Name:                 "World Chain Sepolia",
		},
		// Codex Testnet
		325000: {
			RPCUrl:               "https://eth-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7", // Placeholder
			VaultContractAddress: os.Getenv("CODEX_TESTNET_VAULT_ADDRESS"),
			Name:                 "Codex Testnet",
		},
		// Local Anvil (for development)
		31337: {
			RPCUrl:               "http://127.0.0.1:8545",
			VaultContractAddress: os.Getenv("ANVIL_BASE_SEPOLIA_VAULT_ADDRESS"),
			Name:                 "Anvil Local (Base Sepolia)",
		},
		// Local Anvil (for development)
		31338: {
			RPCUrl:               "http://127.0.0.1:8546",
			VaultContractAddress: os.Getenv("ANVIL_OPTIMISM_SEPOLIA_VAULT_ADDRESS"),
			Name:                 "Anvil Local (Optimism Sepolia)",
		},
	}
}

// ChainConfigs holds the loaded chain configurations
var ChainConfigs = LoadChainConfigs()

// GetPrivateKey returns the private key from environment variables
func GetPrivateKey() string {
	godotenv.Load(".env")
	return os.Getenv("PRIVATE_KEY")
}

// GetGameFactoryAddress returns the GameFactory contract address (deployed only on Base Sepolia)
func GetGameFactoryAddress() string {
	godotenv.Load(".env")
	return os.Getenv("GAME_FACTORY_ADDRESS")
}

// GetEnv gets an environment variable with a default value
func GetEnv(key, defaultValue string) string {
	godotenv.Load(".env")
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt gets an environment variable as integer with a default value
func GetEnvInt(key string, defaultValue int) int {
	godotenv.Load(".env")
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvRequired gets a required environment variable
func GetEnvRequired(key string) (string, error) {
	godotenv.Load(".env")
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s not set", key)
	}
	return value, nil
}

// NewClient creates a new Ethereum client for the specified chain ID
func NewClient(chainID uint64) (*ethclient.Client, error) {
	config, exists := ChainConfigs[chainID]
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
	chains := make([]uint64, 0, len(ChainConfigs))
	for chainID := range ChainConfigs {
		chains = append(chains, chainID)
	}
	return chains
}

// ReloadChainConfigs reloads chain configurations from environment variables
func ReloadChainConfigs() {
	ChainConfigs = LoadChainConfigs()
}

// ConfigService handles application configuration
type ConfigService interface {
	LoadBlockchainConfig() (*BlockchainConfig, error)
	GetEnv(key, defaultValue string) string
	GetEnvInt(key string, defaultValue int) int
	GetEnvRequired(key string) (string, error)
}

// DefaultConfigService implements ConfigService
type DefaultConfigService struct{}

// NewConfigService creates a new configuration service
func NewConfigService() ConfigService {
	return &DefaultConfigService{}
}

// LoadBlockchainConfig loads blockchain configuration from environment variables
func (c *DefaultConfigService) LoadBlockchainConfig() (*BlockchainConfig, error) {
	// Load environment variables
	godotenv.Load(".env")

	config := &BlockchainConfig{
		RPCUrl:               GetEnv("RPC_URL", "http://127.0.0.1:8545"),
		PrivateKey:           GetEnv("PRIVATE_KEY", ""),
		GameContractAddress:  GetEnv("GAME_CONTRACT_ADDRESS", ""),
		VaultContractAddress: GetEnv("VAULT_CONTRACT_ADDRESS", ""),
		DefaultStakeUSDC:     GetEnv("DEFAULT_STAKE_USDC", "0.01"),
	}

	// Validate required fields
	if config.PrivateKey == "" {
		return nil, fmt.Errorf("PRIVATE_KEY environment variable not set")
	}
	if config.GameContractAddress == "" {
		return nil, fmt.Errorf("GAME_CONTRACT_ADDRESS environment variable not set")
	}
	if config.VaultContractAddress == "" {
		return nil, fmt.Errorf("VAULT_CONTRACT_ADDRESS environment variable not set")
	}

	return config, nil
}

// GetEnv gets an environment variable with a default value (DefaultConfigService method)
func (c *DefaultConfigService) GetEnv(key, defaultValue string) string {
	return GetEnv(key, defaultValue)
}

// GetEnvInt gets an environment variable as integer with a default value (DefaultConfigService method)
func (c *DefaultConfigService) GetEnvInt(key string, defaultValue int) int {
	return GetEnvInt(key, defaultValue)
}

// GetEnvRequired gets a required environment variable (DefaultConfigService method)
func (c *DefaultConfigService) GetEnvRequired(key string) (string, error) {
	return GetEnvRequired(key)
}
