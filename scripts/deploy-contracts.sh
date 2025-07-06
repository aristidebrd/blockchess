#!/bin/bash

# Set environment variables
export GO_BACKEND_ADDRESS=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
export RPC_URL_BASE=http://127.0.0.1:8545
export RPC_URL_OPTIMISM=http://127.0.0.1:8546
export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
export USDC_TOKEN_ADDRESS=0x036CbD53842c5426634e7929541eC2318f3dCF7e

echo "🚀 Deploying BlockChess Contracts to Local Anvil Instances"
echo "=========================================================="

# Check if Base Anvil is running (port 8545)
echo "🔍 Checking Base Anvil (port 8545)..."
if ! curl -s -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  $RPC_URL_BASE > /dev/null 2>&1; then
    echo "❌ Base Anvil is not running on port 8545"
    echo "   Please start it first with: anvil --port 8545 --chain-id 31337"
    exit 1
fi

# Check if Optimism Anvil is running (port 8546)
echo "🔍 Checking Optimism Anvil (port 8546)..."
if ! curl -s -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  $RPC_URL_OPTIMISM > /dev/null 2>&1; then
    echo "❌ Optimism Anvil is not running on port 8546"
    echo "   Please start it first with: anvil --port 8546 --chain-id 31338"
    exit 1
fi

echo "✅ Both Anvil instances are running"
echo "🔧 Environment:"
echo "Backend Address: $GO_BACKEND_ADDRESS"
echo "Base RPC URL: $RPC_URL_BASE"
echo "Optimism RPC URL: $RPC_URL_OPTIMISM"

# Deploy contracts to Base Anvil (Chain ID 31337)
echo ""
echo "📦 Deploying contracts to Base Anvil..."
echo "======================================="
forge script contracts/script/DeployVaultMultiChain.s.sol:DeployVaultMultiChain --rpc-url $RPC_URL_BASE --broadcast
forge script contracts/script/DeployGameFactory.s.sol:DeployGameFactory --rpc-url $RPC_URL_BASE --private-key $PRIVATE_KEY --broadcast

if [ $? -ne 0 ]; then
    echo "❌ Base Anvil deployment failed"
    exit 1
fi

# Deploy contracts to Optimism Anvil (Chain ID 31338)
echo ""
echo "📦 Deploying contracts to Optimism Anvil..."
echo "==========================================="
forge script contracts/script/DeployVaultMultiChain.s.sol:DeployVaultMultiChain --rpc-url $RPC_URL_OPTIMISM --broadcast
forge script contracts/script/DeployGameFactory.s.sol:DeployGameFactory --rpc-url $RPC_URL_OPTIMISM --private-key $PRIVATE_KEY --broadcast

if [ $? -ne 0 ]; then
    echo "❌ Optimism Anvil deployment failed"
    exit 1
fi

echo "✅ Contracts deployed successfully to both chains!"

# Extract contract addresses from Base Anvil deployment (Chain ID 31337)
echo ""
echo "📋 Extracting contract addresses..."
echo "=================================="

BASE_RUN=$(find broadcast/DeployVaultMultiChain.s.sol/84532 -name "run-latest.json" 2>/dev/null)
OPTIMISM_RUN=$(find broadcast/DeployVaultMultiChain.s.sol/11155420 -name "run-latest.json" 2>/dev/null)

# Base Anvil addresses
if [ -f "$BASE_RUN" ]; then
    BASE_GAME_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "GameContract") | .contractAddress' "$BASE_RUN")
    BASE_VAULT_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "VaultContract") | .contractAddress' "$BASE_RUN")
    
    echo "📋 Base Anvil Contract Addresses (Chain ID 84532):"
    echo "GameContract: $BASE_GAME_CONTRACT"
    echo "VaultContract: $BASE_VAULT_CONTRACT"
else
    echo "❌ Could not find Base Anvil deployment artifacts"
    exit 1
fi

# Optimism Anvil addresses
if [ -f "$OPTIMISM_RUN" ]; then
    OPTIMISM_GAME_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "GameContract") | .contractAddress' "$OPTIMISM_RUN")
    OPTIMISM_VAULT_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "VaultContract") | .contractAddress' "$OPTIMISM_RUN")
    
    echo "📋 Optimism Anvil Contract Addresses (Chain ID 11155420):"
    echo "GameContract: $OPTIMISM_GAME_CONTRACT"
    echo "VaultContract: $OPTIMISM_VAULT_CONTRACT"
else
    echo "❌ Could not find Optimism Anvil deployment artifacts"
    exit 1
fi

# GameFactory addresses (only deployed on specific chains)
BASE_FACTORY_RUN=$(find broadcast/DeployGameFactory.s.sol/84532 -name "run-latest.json" 2>/dev/null)
OPTIMISM_FACTORY_RUN=$(find broadcast/DeployGameFactory.s.sol/11155420 -name "run-latest.json" 2>/dev/null)

if [ -f "$BASE_FACTORY_RUN" ]; then
    BASE_GAME_FACTORY=$(jq -r '.transactions[] | select(.contractName == "GameFactory") | .contractAddress' "$BASE_FACTORY_RUN")
    echo "GameFactory (Base): $BASE_GAME_FACTORY"
fi

if [ -f "$OPTIMISM_FACTORY_RUN" ]; then
    OPTIMISM_GAME_FACTORY=$(jq -r '.transactions[] | select(.contractName == "GameFactory") | .contractAddress' "$OPTIMISM_FACTORY_RUN")
    echo "GameFactory (Optimism): $OPTIMISM_GAME_FACTORY"
fi

# Update environment files
echo ""
echo "📝 Updating environment files..."
echo "==============================="

# Update .env with Base Anvil as primary (for backwards compatibility)
cat > .env << EOF
# Local Development Environment - Base Anvil Primary
BASE_RPC_URL=http://127.0.0.1:8545
BASE_CHAIN_ID=84532
OPTIMISM_RPC_URL=http://127.0.0.1:8546
OPTIMISM_CHAIN_ID=11155420

# Base Anvil Contract Addresses (Chain ID 31337)
GAME_CONTRACT_ADDRESS=$BASE_GAME_CONTRACT
VAULT_CONTRACT_ADDRESS=$BASE_VAULT_CONTRACT
GAME_FACTORY_ADDRESS=$BASE_GAME_FACTORY

# Optimism Anvil Contract Addresses (Chain ID 31338)
ANVIL_OPTIMISM_SEPOLIA_GAME_CONTRACT_ADDRESS=$OPTIMISM_GAME_CONTRACT
ANVIL_OPTIMISM_SEPOLIA_VAULT_CONTRACT_ADDRESS=$OPTIMISM_VAULT_CONTRACT
ANVIL_OPTIMISM_SEPOLIA_GAME_FACTORY_ADDRESS=$OPTIMISM_GAME_FACTORY

# Backend configuration
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
GO_BACKEND_ADDRESS=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
DEFAULT_STAKE_USDC=0.01

# WebSocket configuration
WS_PORT=8080
HTTP_PORT=8081

# Multi-chain configuration
ANVIL_BASE_SEPOLIA_VAULT_ADDRESS=$BASE_VAULT_CONTRACT
ANVIL_OPTIMISM_SEPOLIA_VAULT_ADDRESS=$OPTIMISM_VAULT_CONTRACT

# Permit2 addresses
OP_SEPOLIA_PERMIT2_ADDRESS=0x000000000022D473030F116dDEE9F6B43aC78BA3
ETHEREUM_SEPOLIA_PERMIT2_ADDRESS=0x000000000022d473030f116ddee9f6b43ac78ba3
AVALANCHE_FUJI_PERMIT2_ADDRESS=0x000000000022D473030F116dDEE9F6B43aC78BA3
ARBITRUM_SEPOLIA_PERMIT2_ADDRESS=0x000000000022D473030F116dDEE9F6B43aC78BA3
UNICHAIN_SEPOLIA_PERMIT2_ADDRESS=0x0000000000000000000000000000000000000000
BASE_SEPOLIA_PERMIT2_ADDRESS=0x000000000022D473030F116dDEE9F6B43aC78BA3
POLYGON_AMOY_PERMIT2_ADDRESS=0x000000000022D473030F116dDEE9F6B43aC78BA3
LINEA_SEPOLIA_PERMIT2_ADDRESS=0x0000000000000000000000000000000000000000
SONIC_TESTNET_PERMIT2_ADDRESS=0x0000000000000000000000000000000000000000
WORLD_CHAIN_SEPOLIA_PERMIT2_ADDRESS=0x000000000022D473030F116dDEE9F6B43aC78BA3
CODEX_TESTNET_PERMIT2_ADDRESS=0x0000000000000000000000000000000000000000
EOF

echo "✅ Environment files updated!"

# Test contract deployments
echo ""
echo "🧪 Testing contract deployments..."
echo "================================="

# Test Base Anvil contracts
echo "Testing Base Anvil contracts..."
if cast call $BASE_GAME_CONTRACT "authorizedBackend()" --rpc-url $RPC_URL_BASE | grep -q "0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266"; then
    echo "✅ Base GameContract is working correctly!"
else
    echo "❌ Base GameContract test failed"
fi

if cast call $BASE_VAULT_CONTRACT "getAuthorizedBackend()" --rpc-url $RPC_URL_BASE | grep -q "0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266"; then
    echo "✅ Base VaultContract is working correctly!"
else
    echo "❌ Base VaultContract test failed"
fi

# Test Optimism Anvil contracts
echo "Testing Optimism Anvil contracts..."
if cast call $OPTIMISM_GAME_CONTRACT "authorizedBackend()" --rpc-url $RPC_URL_OPTIMISM | grep -q "0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266"; then
    echo "✅ Optimism GameContract is working correctly!"
else
    echo "❌ Optimism GameContract test failed"
fi

if cast call $OPTIMISM_VAULT_CONTRACT "getAuthorizedBackend()" --rpc-url $RPC_URL_OPTIMISM | grep -q "0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266"; then
    echo "✅ Optimism VaultContract is working correctly!"
else
    echo "❌ Optimism VaultContract test failed"
fi

echo ""
echo "🎉 Multi-Chain Deployment Complete!"
echo "==================================="
echo "📋 Base Anvil (Chain ID 84532, Port 8545):"
echo "   GameContract: $BASE_GAME_CONTRACT"
echo "   VaultContract: $BASE_VAULT_CONTRACT"
echo "   GameFactory: $BASE_GAME_FACTORY"
echo ""
echo "📋 Optimism Anvil (Chain ID 11155420, Port 8546):"
echo "   GameContract: $OPTIMISM_GAME_CONTRACT"
echo "   VaultContract: $OPTIMISM_VAULT_CONTRACT"
echo "   GameFactory: $OPTIMISM_GAME_FACTORY"
echo ""
echo "🚀 Next steps:"
echo "1. Make sure MetaMask is connected to one of the Anvil instances:"
echo "   - Base Anvil: RPC http://127.0.0.1:8545, Chain ID 84532"
echo "   - Optimism Anvil: RPC http://127.0.0.1:8546, Chain ID 11155420"
echo "2. Start backend: go run main.go"
echo "3. Start frontend: npm run dev"
echo ""
echo "💡 The .env file has been configured with Base Anvil as primary"
echo "   You can switch to Optimism Anvil by changing RPC_URL and CHAIN_ID" 
