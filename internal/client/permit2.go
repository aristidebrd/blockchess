package client

import (
	"blockchess/contracts-bindings/permit2"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// Permit2Client wraps the Permit2 contract instance
type Permit2Client struct {
	contract *permit2.Permit2
	client   *ethclient.Client
	address  common.Address
	chainID  uint64
}

// PermitSignatureData represents the data needed for a Permit2 signature
type PermitSignatureData struct {
	Owner       common.Address `json:"owner"`
	Spender     common.Address `json:"spender"`
	Token       common.Address `json:"token"`
	Amount      *big.Int       `json:"amount"`
	Expiration  *big.Int       `json:"expiration"`
	Nonce       *big.Int       `json:"nonce"`
	SigDeadline *big.Int       `json:"sigDeadline"`
	ChainID     uint64         `json:"chainId"`
	Signature   string         `json:"signature"`
}

// NewPermit2Client creates a new Permit2 client
func NewPermit2Client(client *ethclient.Client, chainID uint64) (*Permit2Client, error) {
	permit2Address := GetPermit2Address(chainID)
	if permit2Address == "" {
		return nil, fmt.Errorf("no Permit2 address configured for chain %d", chainID)
	}

	contractAddress := common.HexToAddress(permit2Address)
	contract, err := permit2.NewPermit2(contractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create Permit2 contract instance: %w", err)
	}

	return &Permit2Client{
		contract: contract,
		client:   client,
		address:  contractAddress,
		chainID:  chainID,
	}, nil
}

// GetNonce gets the current nonce for a user's token-spender pair
func (p *Permit2Client) GetNonce(owner, token, spender common.Address) (*big.Int, error) {
	allowance, err := p.contract.Allowance(&bind.CallOpts{}, owner, token, spender)
	if err != nil {
		return nil, fmt.Errorf("failed to get allowance: %w", err)
	}
	return allowance.Nonce, nil
}

// GetDomainSeparator gets the domain separator for the Permit2 contract
func (p *Permit2Client) GetDomainSeparator() ([32]byte, error) {
	return p.contract.DOMAINSEPARATOR(&bind.CallOpts{})
}

// CreatePermitTypedData creates EIP-712 typed data for a Permit2 signature
func (p *Permit2Client) CreatePermitTypedData(
	owner, spender, token common.Address,
	amount, expiration, nonce, sigDeadline *big.Int,
) (*apitypes.TypedData, error) {
	// Create the typed data
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
				{Name: "expiration", Type: "uint48"},
				{Name: "nonce", Type: "uint48"},
			},
		},
		PrimaryType: "PermitSingle",
		Domain: apitypes.TypedDataDomain{
			Name:              "Permit2",
			ChainId:           (*math.HexOrDecimal256)(big.NewInt(int64(p.chainID))),
			VerifyingContract: p.address.Hex(),
		},
		Message: apitypes.TypedDataMessage{
			"details": apitypes.TypedDataMessage{
				"token":      token.Hex(),
				"amount":     amount.String(),
				"expiration": expiration.String(),
				"nonce":      nonce.String(),
			},
			"spender":     spender.Hex(),
			"sigDeadline": sigDeadline.String(),
		},
	}

	return typedData, nil
}

// CreatePermitSignatureData creates the complete permit signature data for frontend
func (p *Permit2Client) CreatePermitSignatureData(
	owner, spender, token common.Address,
	amount *big.Int,
) (*PermitSignatureData, *apitypes.TypedData, error) {
	// Get current nonce
	nonce, err := p.GetNonce(owner, token, spender)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	// Set expiration to 1 hour from now
	expiration := big.NewInt(time.Now().Add(time.Hour).Unix())

	// Set signature deadline to 30 minutes from now
	sigDeadline := big.NewInt(time.Now().Add(30 * time.Minute).Unix())

	// Create typed data
	typedData, err := p.CreatePermitTypedData(owner, spender, token, amount, expiration, nonce, sigDeadline)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create typed data: %w", err)
	}

	permitData := &PermitSignatureData{
		Owner:       owner,
		Spender:     spender,
		Token:       token,
		Amount:      amount,
		Expiration:  expiration,
		Nonce:       nonce,
		SigDeadline: sigDeadline,
		ChainID:     p.chainID,
	}

	return permitData, typedData, nil
}

// ExecutePermit executes a permit transaction using the provided signature
func (p *Permit2Client) ExecutePermit(
	privateKey string,
	permitData *PermitSignatureData,
	signature string,
) error {
	// Parse private key
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create auth
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(int64(p.chainID)))
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	// Create permit single struct
	permitSingle := permit2.IAllowanceTransferPermitSingle{
		Details: permit2.IAllowanceTransferPermitDetails{
			Token:      permitData.Token,
			Amount:     permitData.Amount,
			Expiration: permitData.Expiration,
			Nonce:      permitData.Nonce,
		},
		Spender:     permitData.Spender,
		SigDeadline: permitData.SigDeadline,
	}

	// Convert signature to bytes
	sigBytes := common.FromHex(signature)

	// Execute permit
	tx, err := p.contract.Permit0(auth, permitData.Owner, permitSingle, sigBytes)
	if err != nil {
		return fmt.Errorf("failed to execute permit: %w", err)
	}

	fmt.Printf("Permit transaction sent: %s\n", tx.Hash().Hex())
	return nil
}

// GetUSDCAddress returns the USDC contract address for the chain
func (p *Permit2Client) GetUSDCAddress() (common.Address, error) {
	// This should be loaded from environment variables or configuration
	// For now, we'll use a placeholder - you should add USDC addresses to your config
	usdcAddresses := map[uint64]string{
		11155111: "0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238", // Ethereum Sepolia
		84532:    "0x036CbD53842c5426634e7929541eC2318f3dCF7e", // Base Sepolia
		43113:    "0x5425890298aed601595a70AB815c96711a31Bc65", // Avalanche Fuji
		11155420: "0x5fd84259d66Cd46123540766Be93DFE6D43130D7", // OP Sepolia
		421614:   "0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d", // Arbitrum Sepolia
		80002:    "0x41e94eb019c0762f9bfcf9fb1e58725bfb0e7582", // Polygon Amoy
		// Add more chains as needed
	}

	address, exists := usdcAddresses[p.chainID]
	if !exists {
		return common.Address{}, fmt.Errorf("USDC address not configured for chain %d", p.chainID)
	}

	return common.HexToAddress(address), nil
}

// CreateGameStakePermit creates a permit for staking USDC in a game
func (p *Permit2Client) CreateGameStakePermit(
	owner common.Address,
	vaultAddress common.Address,
) (*PermitSignatureData, *apitypes.TypedData, error) {
	// Get USDC address
	usdcAddress, err := p.GetUSDCAddress()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get USDC address: %w", err)
	}

	// Set amount to a reasonable amount for game staking (e.g., 1 USDC)
	amount := new(big.Int).SetInt64(1000000) // 1 USDC in 6 decimal places

	return p.CreatePermitSignatureData(owner, vaultAddress, usdcAddress, amount)
}

// Permit2Manager manages Permit2 clients across multiple chains
type Permit2Manager struct {
	clients map[uint64]*Permit2Client
}

// NewPermit2Manager creates a new Permit2 manager
func NewPermit2Manager(clients *Clients) (*Permit2Manager, error) {
	manager := &Permit2Manager{
		clients: make(map[uint64]*Permit2Client),
	}

	// Initialize Permit2 clients for all supported chains
	supportedChains := GetSupportedChains()
	for _, chainID := range supportedChains {
		client, err := clients.GetClientByChainID(chainID)
		if err != nil {
			fmt.Printf("Warning: Failed to get client for chain %d: %v\n", chainID, err)
			continue
		}

		permit2Client, err := NewPermit2Client(client, chainID)
		if err != nil {
			fmt.Printf("Warning: Failed to create Permit2 client for chain %d: %v\n", chainID, err)
			continue
		}

		manager.clients[chainID] = permit2Client
	}

	return manager, nil
}

// GetPermit2Client returns the Permit2 client for a specific chain
func (pm *Permit2Manager) GetPermit2Client(chainID uint64) (*Permit2Client, error) {
	client, exists := pm.clients[chainID]
	if !exists {
		return nil, fmt.Errorf("Permit2 client not available for chain %d", chainID)
	}
	return client, nil
}

// GetAvailableChains returns a list of chains with available Permit2 clients
func (pm *Permit2Manager) GetAvailableChains() []uint64 {
	chains := make([]uint64, 0, len(pm.clients))
	for chainID := range pm.clients {
		chains = append(chains, chainID)
	}
	return chains
}
