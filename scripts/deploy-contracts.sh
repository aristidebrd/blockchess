#!/bin/bash

# Set environment variables
export GO_BACKEND_ADDRESS=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
export RPC_URL_BASE=http://127.0.0.1:8545
export RPC_URL_OPTIMISM=http://127.0.0.1:8546
export RPC_URL_ETHEREUM=http://127.0.0.1:8547
export RPC_URL_AVALANCHE=http://127.0.0.1:8548
export RPC_URL_ARBITRUM=http://127.0.0.1:8549
export RPC_URL_UNICHAIN=http://127.0.0.1:8550
export RPC_URL_BASE_SEPOLIA=http://127.0.0.1:8551
export RPC_URL_POLYGON_AMOY=http://127.0.0.1:8552
export RPC_URL_LINEA_SEPOLIA=http://127.0.0.1:8553
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

# Check if Ethereum Anvil is running (port 8547)
echo "🔍 Checking Ethereum Anvil (port 8547)..."
if ! curl -s -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  $RPC_URL_ETHEREUM > /dev/null 2>&1; then
    echo "❌ Ethereum Anvil is not running on port 8547"
    echo "   Please start it first with: anvil --port 8547 --chain-id 31337"
    exit 1
fi

# Check if Avalanche Anvil is running (port 8548)   
echo "🔍 Checking Avalanche Anvil (port 8548)..."
if ! curl -s -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  $RPC_URL_AVALANCHE > /dev/null 2>&1; then
    echo "❌ Avalanche Anvil is not running on port 8548"
    echo "   Please start it first with: anvil --port 8548 --chain-id 31337"
    exit 1
fi

# Check if Arbitrum Anvil is running (port 8549)    
echo "🔍 Checking Arbitrum Anvil (port 8549)..."
if ! curl -s -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  $RPC_URL_ARBITRUM > /dev/null 2>&1; then
    echo "❌ Arbitrum Anvil is not running on port 8549"
    echo "   Please start it first with: anvil --port 8549 --chain-id 31337"    
    exit 1
fi

echo "✅ Both Anvil instances are running"
echo "🔧 Environment:"
echo "Backend Address: $GO_BACKEND_ADDRESS"
echo "Base RPC URL: $RPC_URL_BASE"
echo "Optimism RPC URL: $RPC_URL_OPTIMISM"
echo "Ethereum RPC URL: $RPC_URL_ETHEREUM"
echo "Avalanche RPC URL: $RPC_URL_AVALANCHE"
echo "Arbitrum RPC URL: $RPC_URL_ARBITRUM"

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

# Deploy contracts to Ethereum Anvil (Chain ID 31337)
echo ""
echo "📦 Deploying contracts to Ethereum Anvil..."
echo "==========================================="
forge script contracts/script/DeployVaultMultiChain.s.sol:DeployVaultMultiChain --rpc-url $RPC_URL_ETHEREUM --broadcast
forge script contracts/script/DeployGameFactory.s.sol:DeployGameFactory --rpc-url $RPC_URL_ETHEREUM --private-key $PRIVATE_KEY --broadcast

if [ $? -ne 0 ]; then
    echo "❌ Ethereum Anvil deployment failed"
    exit 1
fi

# Deploy contracts to Avalanche Anvil (Chain ID 31337)

echo ""
echo "📦 Deploying contracts to Avalanche Anvil..."
echo "==========================================="
forge script contracts/script/DeployVaultMultiChain.s.sol:DeployVaultMultiChain --rpc-url $RPC_URL_AVALANCHE --broadcast
forge script contracts/script/DeployGameFactory.s.sol:DeployGameFactory --rpc-url $RPC_URL_AVALANCHE --private-key $PRIVATE_KEY --broadcast

if [ $? -ne 0 ]; then
    echo "❌ Avalanche Anvil deployment failed"
    exit 1
fi

# Deploy contracts to Arbitrum Anvil (Chain ID 31337)

echo ""
echo "📦 Deploying contracts to Arbitrum Anvil..."
echo "==========================================="
forge script contracts/script/DeployVaultMultiChain.s.sol:DeployVaultMultiChain --rpc-url $RPC_URL_ARBITRUM --broadcast
forge script contracts/script/DeployGameFactory.s.sol:DeployGameFactory --rpc-url $RPC_URL_ARBITRUM --private-key $PRIVATE_KEY --broadcast

if [ $? -ne 0 ]; then
    echo "❌ Arbitrum Anvil deployment failed"
    exit 1
fi

    echo "✅ Contracts deployed successfully to both chains!"

# Extract contract addresses from Base Anvil deployment (Chain ID 31337)
echo ""
echo "📋 Extracting contract addresses..."
echo "=================================="

BASE_RUN=$(find broadcast/DeployVaultMultiChain.s.sol/84532 -name "run-latest.json" 2>/dev/null)
OPTIMISM_RUN=$(find broadcast/DeployVaultMultiChain.s.sol/11155420 -name "run-latest.json" 2>/dev/null)
ETHEREUM_RUN=$(find broadcast/DeployVaultMultiChain.s.sol/11155111 -name "run-latest.json" 2>/dev/null)
AVALANCHE_RUN=$(find broadcast/DeployVaultMultiChain.s.sol/43113 -name "run-latest.json" 2>/dev/null)
ARBITRUM_RUN=$(find broadcast/DeployVaultMultiChain.s.sol/421613 -name "run-latest.json" 2>/dev/null)

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

# Ethereum Anvil addresses
if [ -f "$ETHEREUM_RUN" ]; then
    ETHEREUM_GAME_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "GameContract") | .contractAddress' "$ETHEREUM_RUN")
    ETHEREUM_VAULT_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "VaultContract") | .contractAddress' "$ETHEREUM_RUN")

    echo "📋 Ethereum Anvil Contract Addresses (Chain ID 11155111):"
    echo "GameContract: $ETHEREUM_GAME_CONTRACT"
    echo "VaultContract: $ETHEREUM_VAULT_CONTRACT"
else   
    echo "❌ Could not find Ethereum Anvil deployment artifacts"
    exit 1
fi

# Avalanche Anvil addresses

if [ -f "$AVALANCHE_RUN" ]; then
    AVALANCHE_GAME_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "GameContract") | .contractAddress' "$AVALANCHE_RUN")
    AVALANCHE_VAULT_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "VaultContract") | .contractAddress' "$AVALANCHE_RUN")
    
    echo "📋 Avalanche Anvil Contract Addresses (Chain ID 43113):"
    echo "GameContract: $AVALANCHE_GAME_CONTRACT"
    echo "VaultContract: $AVALANCHE_VAULT_CONTRACT"
else
    echo "❌ Could not find Avalanche Anvil deployment artifacts"
    exit 1
    fi

# Arbitrum Anvil addresses

if [ -f "$ARBITRUM_RUN" ]; then
    ARBITRUM_GAME_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "GameContract") | .contractAddress' "$ARBITRUM_RUN")
    ARBITRUM_VAULT_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "VaultContract") | .contractAddress' "$ARBITRUM_RUN")

    echo "📋 Arbitrum Anvil Contract Addresses (Chain ID 421613):"
    echo "GameContract: $ARBITRUM_GAME_CONTRACT"
    echo "VaultContract: $ARBITRUM_VAULT_CONTRACT"
else
    echo "❌ Could not find Arbitrum Anvil deployment artifacts"
    exit 1
fi
    # GameFactory addresses (only deployed on specific chains)
BASE_FACTORY_RUN=$(find broadcast/DeployGameFactory.s.sol/84532 -name "run-latest.json" 2>/dev/null)
OPTIMISM_FACTORY_RUN=$(find broadcast/DeployGameFactory.s.sol/11155420 -name "run-latest.json" 2>/dev/null)
ETHEREUM_FACTORY_RUN=$(find broadcast/DeployGameFactory.s.sol/11155111 -name "run-latest.json" 2>/dev/null)
AVALANCHE_FACTORY_RUN=$(find broadcast/DeployGameFactory.s.sol/43113 -name "run-latest.json" 2>/dev/null)
ARBITRUM_FACTORY_RUN=$(find broadcast/DeployGameFactory.s.sol/421613 -name "run-latest.json" 2>/dev/null)

if [ -f "$BASE_FACTORY_RUN" ]; then
    BASE_GAME_FACTORY=$(jq -r '.transactions[] | select(.contractName == "GameFactory") | .contractAddress' "$BASE_FACTORY_RUN")
    echo "GameFactory (Base): $BASE_GAME_FACTORY"
fi

if [ -f "$OPTIMISM_FACTORY_RUN" ]; then
    OPTIMISM_GAME_FACTORY=$(jq -r '.transactions[] | select(.contractName == "GameFactory") | .contractAddress' "$OPTIMISM_FACTORY_RUN")
    echo "GameFactory (Optimism): $OPTIMISM_GAME_FACTORY"
fi

if [ -f "$ETHEREUM_FACTORY_RUN" ]; then
    ETHEREUM_GAME_FACTORY=$(jq -r '.transactions[] | select(.contractName == "GameFactory") | .contractAddress' "$ETHEREUM_FACTORY_RUN")
    echo "GameFactory (Ethereum): $ETHEREUM_GAME_FACTORY"
fi

if [ -f "$AVALANCHE_FACTORY_RUN" ]; then
    AVALANCHE_GAME_FACTORY=$(jq -r '.transactions[] | select(.contractName == "GameFactory") | .contractAddress' "$AVALANCHE_FACTORY_RUN")
    echo "GameFactory (Avalanche): $AVALANCHE_GAME_FACTORY"
fi

if [ -f "$ARBITRUM_FACTORY_RUN" ]; then
    ARBITRUM_GAME_FACTORY=$(jq -r '.transactions[] | select(.contractName == "GameFactory") | .contractAddress' "$ARBITRUM_FACTORY_RUN")
    echo "GameFactory (Arbitrum): $ARBITRUM_GAME_FACTORY"
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
ETHEREUM_RPC_URL=http://127.0.0.1:8547
ETHEREUM_CHAIN_ID=11155111
AVALANCHE_RPC_URL=http://127.0.0.1:8548
AVALANCHE_CHAIN_ID=43113
ARBITRUM_RPC_URL=http://127.0.0.1:8549
ARBITRUM_CHAIN_ID=421613

# Base Anvil Contract Addresses (Chain ID 31337)
GAME_CONTRACT_ADDRESS=$BASE_GAME_CONTRACT
VAULT_CONTRACT_ADDRESS=$BASE_VAULT_CONTRACT
GAME_FACTORY_ADDRESS=$BASE_GAME_FACTORY

# Optimism Anvil Contract Addresses (Chain ID 31338)
ANVIL_OPTIMISM_SEPOLIA_GAME_CONTRACT_ADDRESS=$OPTIMISM_GAME_CONTRACT
ANVIL_OPTIMISM_SEPOLIA_VAULT_CONTRACT_ADDRESS=$OPTIMISM_VAULT_CONTRACT
ANVIL_OPTIMISM_SEPOLIA_GAME_FACTORY_ADDRESS=$OPTIMISM_GAME_FACTORY

# Ethereum Anvil Contract Addresses (Chain ID 31337)
ANVIL_ETHEREUM_SEPOLIA_GAME_CONTRACT_ADDRESS=$ETHEREUM_GAME_CONTRACT
ANVIL_ETHEREUM_SEPOLIA_VAULT_CONTRACT_ADDRESS=$ETHEREUM_VAULT_CONTRACT
ANVIL_ETHEREUM_SEPOLIA_GAME_FACTORY_ADDRESS=$ETHEREUM_GAME_FACTORY

# Avalanche Anvil Contract Addresses (Chain ID 31337)
ANVIL_AVALANCHE_SEPOLIA_GAME_CONTRACT_ADDRESS=$AVALANCHE_GAME_CONTRACT
ANVIL_AVALANCHE_SEPOLIA_VAULT_CONTRACT_ADDRESS=$AVALANCHE_VAULT_CONTRACT
ANVIL_AVALANCHE_SEPOLIA_GAME_FACTORY_ADDRESS=$AVALANCHE_GAME_FACTORY

# Arbitrum Anvil Contract Addresses (Chain ID 31337)
ANVIL_ARBITRUM_SEPOLIA_GAME_CONTRACT_ADDRESS=$ARBITRUM_GAME_CONTRACT
ANVIL_ARBITRUM_SEPOLIA_VAULT_CONTRACT_ADDRESS=$ARBITRUM_VAULT_CONTRACT
ANVIL_ARBITRUM_SEPOLIA_GAME_FACTORY_ADDRESS=$ARBITRUM_GAME_FACTORY

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
ANVIL_ETHEREUM_SEPOLIA_VAULT_ADDRESS=$ETHEREUM_VAULT_CONTRACT
ANVIL_AVALANCHE_SEPOLIA_VAULT_ADDRESS=$AVALANCHE_VAULT_CONTRACT
ANVIL_ARBITRUM_SEPOLIA_VAULT_ADDRESS=$ARBITRUM_VAULT_CONTRACT

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
echo "📋 Ethereum Anvil (Chain ID 11155111, Port 8547):"
echo "   GameContract: $ETHEREUM_GAME_CONTRACT"
echo "   VaultContract: $ETHEREUM_VAULT_CONTRACT"
echo "   GameFactory: $ETHEREUM_GAME_FACTORY"
echo ""
echo "📋 Avalanche Anvil (Chain ID 43113, Port 8548):"
echo "   GameContract: $AVALANCHE_GAME_CONTRACT"
echo "   VaultContract: $AVALANCHE_VAULT_CONTRACT"
echo "   GameFactory: $AVALANCHE_GAME_FACTORY"
echo ""
echo "📋 Arbitrum Anvil (Chain ID 421613, Port 8549):"
echo "   GameContract: $ARBITRUM_GAME_CONTRACT"
echo "   VaultContract: $ARBITRUM_VAULT_CONTRACT"
echo "   GameFactory: $ARBITRUM_GAME_FACTORY"
echo ""
echo "🚀 Next steps:"
echo "1. Make sure MetaMask is connected to one of the Anvil instances:"
echo "   - Base Anvil: RPC http://127.0.0.1:8545, Chain ID 84532"
echo "   - Optimism Anvil: RPC http://127.0.0.1:8546, Chain ID 11155420"
echo "   - Ethereum Anvil: RPC http://127.0.0.1:8547, Chain ID 11155111"
echo "   - Avalanche Anvil: RPC http://127.0.0.1:8548, Chain ID 43113"
echo "2. Start backend: go run main.go"
echo "3. Start frontend: npm run dev"
echo ""
echo "💡 The .env file has been configured with Base Anvil as primary"
echo "   You can switch to Optimism Anvil by changing RPC_URL and CHAIN_ID" 
