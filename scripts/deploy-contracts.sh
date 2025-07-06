#!/bin/bash

# Set environment variables
export GO_BACKEND_ADDRESS=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
export RPC_URL=http://127.0.0.1:8545
export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
export USDC_TOKEN_ADDRESS=0x036CbD53842c5426634e7929541eC2318f3dCF7e

echo "🚀 Deploying BlockChess Contracts to Local Anvil"
echo "================================================"

# Check if Anvil is running
if ! curl -s -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  http://127.0.0.1:8545 > /dev/null 2>&1; then
    echo "❌ Anvil is not running. Please start it first with: npm run anvil"
    exit 1
fi

echo "✅ Anvil is running with Chain ID 84532"
echo "🔧 Environment:"
echo "Backend Address: $GO_BACKEND_ADDRESS"
echo "RPC URL: $RPC_URL"

# Deploy contracts
echo "📦 Deploying contracts..."
forge script contracts/script/DeployVaultSystem.s.sol --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast
forge script contracts/script/DeployVaultMultiChain.s.sol:DeployVaultMultiChain --rpc-url $RPC_URL --broadcast

if [ $? -eq 0 ]; then
    echo "✅ Contracts deployed successfully!"
    
    # Extract contract addresses from the latest deployment
    LATEST_RUN=$(find broadcast/DeployVaultSystem.s.sol/84532 -name "run-latest.json" 2>/dev/null)
    
    if [ -f "$LATEST_RUN" ]; then
        GAME_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "GameContract") | .contractAddress' "$LATEST_RUN")
        VAULT_CONTRACT=$(jq -r '.transactions[] | select(.contractName == "VaultContract") | .contractAddress' "$LATEST_RUN")
        
        echo "📋 Contract Addresses:"
        echo "GameContract: $GAME_CONTRACT"
        echo "VaultContract: $VAULT_CONTRACT"
        
        # Update environment files
        echo "📝 Updating environment files..."
        
        # Update contracts/.env
        cat > backend/.env << EOF
# Local Anvil Configuration for Development
RPC_URL=http://127.0.0.1:8545
CHAIN_ID=84532

# Use first Anvil private key
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
GO_BACKEND_ADDRESS=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266

# Deployed contract addresses
GAME_CONTRACT_ADDRESS=$GAME_CONTRACT
VAULT_CONTRACT_ADDRESS=$VAULT_CONTRACT
EOF
        
        # Update .env
        cat > .env << EOF
# Local Development Environment
RPC_URL=http://127.0.0.1:8545
CHAIN_ID=84532

# Deployed contract addresses
GAME_CONTRACT_ADDRESS=$GAME_CONTRACT
VAULT_CONTRACT_ADDRESS=$VAULT_CONTRACT

# Backend configuration
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
GO_BACKEND_ADDRESS=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
DEFAULT_STAKE_ETH=0.01

# WebSocket configuration
WS_PORT=8080
HTTP_PORT=8081

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
        
        # Test contract deployment
        echo "🧪 Testing contract deployment..."
        if cast call $GAME_CONTRACT "authorizedBackend()" --rpc-url $RPC_URL | grep -q "0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266"; then
            echo "✅ GameContract is working correctly!"
        else
            echo "❌ GameContract test failed"
        fi
        
        if cast call $VAULT_CONTRACT "authorizedBackend()" --rpc-url $RPC_URL | grep -q "0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266"; then
            echo "✅ VaultContract is working correctly!"
        else
            echo "❌ VaultContract test failed"
        fi
        
        echo ""
        echo "🎉 Deployment Complete!"
        echo "======================="
        echo "GameContract: $GAME_CONTRACT"
        echo "VaultContract: $VAULT_CONTRACT"
        echo ""
        echo "Next steps:"
        echo "1. Make sure MetaMask is connected to Anvil Local (Chain ID: 84532)"
        echo "2. Start backend: cd ../backend && go run cmd/server/main.go"
        echo "3. Start frontend: cd .. && npm run dev"
        
    else
        echo "❌ Could not find deployment artifacts"
    fi
else
    echo "❌ Deployment failed"
    exit 1
fi 
