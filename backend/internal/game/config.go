package game

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

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
	// Try to load .env.local first, then .env
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	config := &BlockchainConfig{
		RPCUrl:               c.GetEnv("RPC_URL", "http://127.0.0.1:8545"),
		PrivateKey:           c.GetEnv("PRIVATE_KEY", ""),
		GameContractAddress:  c.GetEnv("GAME_CONTRACT_ADDRESS", ""),
		VaultContractAddress: c.GetEnv("VAULT_CONTRACT_ADDRESS", ""),
		DefaultStakeETH:      c.GetEnv("DEFAULT_STAKE_ETH", "0.01"),
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

// GetEnv gets an environment variable with a default value
func (c *DefaultConfigService) GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt gets an environment variable as integer with a default value
func (c *DefaultConfigService) GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvRequired gets a required environment variable
func (c *DefaultConfigService) GetEnvRequired(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s not set", key)
	}
	return value, nil
}

// MockConfigService implements ConfigService for testing
type MockConfigService struct {
	env map[string]string
}

// NewMockConfigService creates a new mock configuration service
func NewMockConfigService(env map[string]string) ConfigService {
	return &MockConfigService{env: env}
}

func (m *MockConfigService) LoadBlockchainConfig() (*BlockchainConfig, error) {
	return &BlockchainConfig{
		RPCUrl:               m.GetEnv("RPC_URL", "http://127.0.0.1:8545"),
		PrivateKey:           m.GetEnv("PRIVATE_KEY", "test_private_key"),
		GameContractAddress:  m.GetEnv("GAME_CONTRACT_ADDRESS", "0x1234567890123456789012345678901234567890"),
		VaultContractAddress: m.GetEnv("VAULT_CONTRACT_ADDRESS", "0x1234567890123456789012345678901234567890"),
		DefaultStakeETH:      m.GetEnv("DEFAULT_STAKE_ETH", "0.01"),
	}, nil
}

func (m *MockConfigService) GetEnv(key, defaultValue string) string {
	if value, exists := m.env[key]; exists {
		return value
	}
	return defaultValue
}

func (m *MockConfigService) GetEnvInt(key string, defaultValue int) int {
	if value, exists := m.env[key]; exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (m *MockConfigService) GetEnvRequired(key string) (string, error) {
	if value, exists := m.env[key]; exists {
		return value, nil
	}
	return "", fmt.Errorf("required environment variable %s not set", key)
}
