package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"blockchess/contracts-bindings/permit2"
	"blockchess/contracts-bindings/vaultcontract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Vault wraps a VaultContract instance for a specific chain
type Vault struct {
	contract *vaultcontract.Vaultcontract
	client   *ethclient.Client
	auth     *bind.TransactOpts
	chainID  uint64
}

// VaultManager manages vault contracts across all chains
type VaultManager struct {
	vaults map[uint64]*Vault // chainID -> Vault
}

// NewVault creates a new Vault instance for a specific chain
func NewVault(client *ethclient.Client, privateKey string, chainID uint64) (*Vault, error) {
	// Get vault contract address for this chain
	vaultAddress := GetVaultAddress(chainID)
	if vaultAddress == "" {
		return nil, fmt.Errorf("no vault address configured for chain ID: %d", chainID)
	}

	if !common.IsHexAddress(vaultAddress) {
		return nil, fmt.Errorf("invalid vault address for chain %d: %s", chainID, vaultAddress)
	}

	// Create contract instance
	contract, err := vaultcontract.NewVaultcontract(common.HexToAddress(vaultAddress), client)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault contract instance for chain %d: %w", chainID, err)
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

	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(int64(chainID)))
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	return &Vault{
		contract: contract,
		client:   client,
		auth:     auth,
		chainID:  chainID,
	}, nil
}

// NewVaultManager creates vault instances for all configured chains
func NewVaultManager(clients *Clients) (*VaultManager, error) {
	vaultManager := &VaultManager{
		vaults: make(map[uint64]*Vault),
	}

	privateKey := clients.GetPrivateKey()
	if privateKey == "" {
		return nil, fmt.Errorf("no private key available for vault manager")
	}

	// Create vault instances for all chains that have vault addresses configured
	for chainID := range VaultAddresses {
		client, err := clients.GetClientByChainID(chainID)
		if err != nil {
			log.Printf("Warning: Failed to get client for chain %d: %v", chainID, err)
			continue
		}

		vault, err := NewVault(client, privateKey, chainID)
		if err != nil {
			log.Printf("Warning: Failed to create vault for chain %d: %v", chainID, err)
			continue
		}

		vaultManager.vaults[chainID] = vault
		log.Printf("Created vault for chain %d", chainID)
	}

	return vaultManager, nil
}

// GetVault returns the vault instance for a specific chain ID
func (vm *VaultManager) GetVault(chainID uint64) (*Vault, error) {
	vault, exists := vm.vaults[chainID]
	if !exists {
		return nil, fmt.Errorf("no vault available for chain ID: %d", chainID)
	}
	return vault, nil
}

// GetAvailableChains returns a list of chain IDs that have vault instances
func (vm *VaultManager) GetAvailableChains() []uint64 {
	chains := make([]uint64, 0, len(vm.vaults))
	for chainID := range vm.vaults {
		chains = append(chains, chainID)
	}
	return chains
}

// Stake deposits USDC from a player to the vault contract using Permit2
func (v *Vault) Stake(playerAddress common.Address, gameID uint64, amount *big.Int) error {
	log.Printf("Staking %s USDC for player %s in game %d on chain %d",
		amount.String(), playerAddress.Hex(), gameID, v.chainID)

	gameIDBig := new(big.Int).SetUint64(gameID)

	// Call the stake function (the vault contract will handle the USDC transfer internally)
	tx, err := v.contract.Stake(v.auth, playerAddress, gameIDBig, amount)
	if err != nil {
		return fmt.Errorf("failed to stake transaction: %w", err)
	}

	log.Printf("Stake transaction sent: %s", tx.Hash().Hex())

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), v.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	log.Printf("Successfully staked %s USDC for player %s in game %d on chain %d",
		amount.String(), playerAddress.Hex(), gameID, v.chainID)
	return nil
}

// StakeWithPermit deposits USDC from a player to the vault contract using Permit2 signature
func (v *Vault) StakeWithPermit(playerAddress common.Address, gameID uint64, amount *big.Int, permitData *PermitSignatureData) error {
	log.Printf("Staking %s USDC for player %s in game %d on chain %d using Permit2",
		amount.String(), playerAddress.Hex(), gameID, v.chainID)

	if permitData == nil {
		return fmt.Errorf("permit signature data is required")
	}

	if permitData.Signature == "" {
		return fmt.Errorf("permit signature is required")
	}

	// Verify the permit is for the correct player and amount
	if permitData.Owner != playerAddress {
		return fmt.Errorf("permit owner mismatch: expected %s, got %s", playerAddress.Hex(), permitData.Owner.Hex())
	}

	if permitData.Amount.Cmp(amount) < 0 {
		return fmt.Errorf("permit amount insufficient: need %s, have %s", amount.String(), permitData.Amount.String())
	}

	// Get the vault contract address as the spender
	vaultAddress := GetVaultAddress(v.chainID)
	if vaultAddress == "" {
		return fmt.Errorf("no vault address configured for chain %d", v.chainID)
	}

	expectedSpender := common.HexToAddress(vaultAddress)
	if permitData.Spender != expectedSpender {
		return fmt.Errorf("permit spender mismatch: expected %s, got %s", expectedSpender.Hex(), permitData.Spender.Hex())
	}

	// Check if permit hasn't expired
	now := big.NewInt(time.Now().Unix())
	if permitData.SigDeadline.Cmp(now) <= 0 {
		return fmt.Errorf("permit signature has expired")
	}

	// First, execute the permit to allow the vault to spend USDC
	permit2Address := GetPermit2Address(v.chainID)
	if permit2Address == "" {
		return fmt.Errorf("no Permit2 address configured for chain %d", v.chainID)
	}

	// Create permit2 contract instance
	permit2Contract, err := permit2.NewPermit2(common.HexToAddress(permit2Address), v.client)
	if err != nil {
		return fmt.Errorf("failed to create Permit2 contract instance: %w", err)
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
	sigBytes := common.FromHex(permitData.Signature)

	// Execute permit
	permitTx, err := permit2Contract.Permit0(v.auth, permitData.Owner, permitSingle, sigBytes)
	if err != nil {
		return fmt.Errorf("failed to execute permit: %w", err)
	}

	log.Printf("Permit transaction sent: %s", permitTx.Hash().Hex())

	// Wait for permit transaction to be mined
	permitReceipt, err := bind.WaitMined(context.Background(), v.client, permitTx)
	if err != nil {
		return fmt.Errorf("failed to wait for permit transaction to be mined: %w", err)
	}

	if permitReceipt.Status != 1 {
		return fmt.Errorf("permit transaction failed with status: %d", permitReceipt.Status)
	}

	log.Printf("Permit transaction confirmed: %s", permitTx.Hash().Hex())

	// Now call the stake function
	gameIDBig := new(big.Int).SetUint64(gameID)
	stakeTx, err := v.contract.Stake(v.auth, playerAddress, gameIDBig, amount)
	if err != nil {
		return fmt.Errorf("failed to stake transaction: %w", err)
	}

	log.Printf("Stake transaction sent: %s", stakeTx.Hash().Hex())

	// Wait for stake transaction to be mined
	stakeReceipt, err := bind.WaitMined(context.Background(), v.client, stakeTx)
	if err != nil {
		return fmt.Errorf("failed to wait for stake transaction to be mined: %w", err)
	}

	if stakeReceipt.Status != 1 {
		return fmt.Errorf("stake transaction failed with status: %d", stakeReceipt.Status)
	}

	log.Printf("Successfully staked %s USDC for player %s in game %d on chain %d using Permit2",
		amount.String(), playerAddress.Hex(), gameID, v.chainID)
	return nil
}

// TransferRewards transfers rewards from this vault to the main vault on Base Sepolia
func (v *Vault) TransferRewards(gameID uint64, amount *big.Int, toChain uint64, recipient common.Address, useFastTransfer bool, maxFee *big.Int) error {

	log.Printf("Transferring %s USDC rewards from chain %d to chain %d for recipient %s",
		amount.String(), v.chainID, toChain, recipient.Hex())

	gameIDBig := new(big.Int).SetUint64(gameID)
	toChainIDBig := new(big.Int).SetUint64(toChain)

	// Call the transferRewardsCrossChain function
	tx, err := v.contract.TransferRewardsCrossChain(
		v.auth,
		gameIDBig,
		amount,
		toChainIDBig,
		recipient,
		useFastTransfer,
		maxFee,
	)
	if err != nil {
		return fmt.Errorf("failed to transfer rewards transaction: %w", err)
	}

	log.Printf("Transfer rewards transaction sent: %s", tx.Hash().Hex())

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), v.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	log.Printf("Successfully transferred %s USDC rewards from chain %d to chain %d for recipient %s",
		amount.String(), v.chainID, toChain, recipient.Hex())
	return nil
}

// GetTotalStakes returns the total amount staked in this vault
func (v *Vault) GetTotalStakes() (*big.Int, error) {
	totalStakes, err := v.contract.GetTotalStakes(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get total stakes: %w", err)
	}
	return totalStakes, nil
}
