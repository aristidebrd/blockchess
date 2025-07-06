#!/bin/bash

# GameFactory Deployment Script for Base Sepolia
# This script deploys the GameFactory contract to Base Sepolia testnet

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_PATH="contracts/script/DeployGameFactory.s.sol:DeployGameFactory"
NETWORK="base-sepolia"
RPC_URL_VAR="https://base-sepolia.g.alchemy.com/v2/"
PRIVATE_KEY="22r8dairb21cjlw7" 
BASESCAN_API_KEY="5RE2C7SV8V161RH224YHNN8VCA7ZI58IQK"

echo -e "${BLUE}🚀 GameFactory Deployment Script${NC}"
echo -e "${BLUE}=================================${NC}"
echo ""

# Check if required environment variables are set
if [ -z "$PRIVATE_KEY" ]; then
    echo -e "${RED}❌ Error: PRIVATE_KEY environment variable is not set${NC}"
    echo "Please set your private key: export PRIVATE_KEY=your_private_key_here"
    exit 1
fi

if [ -z "${!RPC_URL_VAR}" ]; then
    echo -e "${RED}❌ Error: ${RPC_URL_VAR} environment variable is not set${NC}"
    echo "Please set the Base Sepolia RPC URL: export ${RPC_URL_VAR}=your_rpc_url_here"
    exit 1
fi

echo -e "${YELLOW}📋 Configuration:${NC}"
echo "Network: ${NETWORK}"
echo "RPC URL: ${!RPC_URL_VAR}"
echo "Script: ${SCRIPT_PATH}"
echo ""

# Function to deploy to a specific network
deploy_to_network() {
    local network=$1
    local rpc_url=$2
    
    echo -e "${BLUE}🔄 Deploying GameFactory to ${network}...${NC}"
    
    # Deploy the contract
    forge script ${SCRIPT_PATH} \
        --rpc-url ${rpc_url} \
        --private-key ${PRIVATE_KEY} \
        --broadcast \
        --verify \
        --etherscan-api-key ${BASESCAN_API_KEY:-""} \
        -vvvv
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Successfully deployed GameFactory to ${network}${NC}"
        echo ""
    else
        echo -e "${RED}❌ Failed to deploy GameFactory to ${network}${NC}"
        return 1
    fi
}

# Build contracts first
echo -e "${BLUE}🔨 Building contracts...${NC}"
forge build

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Build successful${NC}"
echo ""

# Deploy to Base Sepolia
echo -e "${BLUE}🚀 Starting deployment to Base Sepolia...${NC}"
deploy_to_network "base-sepolia" "${!RPC_URL_VAR}"

echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
echo ""
echo -e "${YELLOW}📝 Next steps:${NC}"
echo "1. Update your backend configuration with the deployed GameFactory address"
echo "2. Update the AUTHORIZED_BACKEND address in the deployment script if needed"
echo "3. Test the deployment by creating a game through your backend"
echo ""
echo -e "${BLUE}💡 Useful commands:${NC}"
echo "- Check deployment: forge script ${SCRIPT_PATH} --rpc-url ${!RPC_URL_VAR} --private-key \$PRIVATE_KEY --broadcast --verify"
echo "- Verify contract: forge verify-contract <CONTRACT_ADDRESS> ${SCRIPT_PATH} --chain-id 84532"
echo "" 
