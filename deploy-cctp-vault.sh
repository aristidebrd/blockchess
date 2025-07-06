#!/bin/bash

# CCTP Vault Multi-Chain Deployment Script
# This script deploys the VaultContractWithCCTP to all supported CCTP V2 testnets

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_PATH="contracts/script/DeployVaultMultiChain.s.sol:DeployVaultMultiChain"

# Network configurations
declare -A NETWORKS=(
    ["ethereum-sepolia"]="ETHEREUM_SEPOLIA_RPC_URL"
    ["avalanche-fuji"]="AVALANCHE_FUJI_RPC_URL"
    ["op-sepolia"]="OP_SEPOLIA_RPC_URL"
    ["arbitrum-sepolia"]="ARBITRUM_SEPOLIA_RPC_URL"
    ["base-sepolia"]="BASE_SEPOLIA_RPC_URL"
    ["polygon-amoy"]="POLYGON_AMOY_RPC_URL"
    ["unichain-sepolia"]="UNICHAIN_SEPOLIA_RPC_URL"
    ["linea-sepolia"]="LINEA_SEPOLIA_RPC_URL"
    ["codex-testnet"]="CODEX_TESTNET_RPC_URL"
    ["sonic-testnet"]="SONIC_TESTNET_RPC_URL"
    ["world-chain-sepolia"]="WORLD_CHAIN_SEPOLIA_RPC_URL"
)

declare -A CHAIN_IDS=(
    ["ethereum-sepolia"]="11155111"
    ["avalanche-fuji"]="43113"
    ["op-sepolia"]="11155420"
    ["arbitrum-sepolia"]="421614"
    ["base-sepolia"]="84532"
    ["polygon-amoy"]="80002"
    ["unichain-sepolia"]="1301"
    ["linea-sepolia"]="59141"
    ["codex-testnet"]="325000"
    ["sonic-testnet"]="64165"
    ["world-chain-sepolia"]="4801"
)

# Functions
print_header() {
    echo -e "${BLUE}"
    echo "========================================"
    echo "  CCTP Vault Multi-Chain Deployment"
    echo "========================================"
    echo -e "${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

check_prerequisites() {
    print_info "Checking prerequisites..."
    
    # Check if forge is installed
    if ! command -v forge &> /dev/null; then
        print_error "Forge is not installed. Please install Foundry: https://getfoundry.sh/"
        exit 1
    fi
    
    # Check if .env file exists
    if [ ! -f .env ]; then
        print_error ".env file not found. Please create one with your configuration."
        exit 1
    fi
    
    # Source .env file
    source .env
    
    # Check required environment variables
    if [ -z "$PRIVATE_KEY" ]; then
        print_error "PRIVATE_KEY not set in .env file"
        exit 1
    fi
    
    if [ -z "$AUTHORIZED_BACKEND" ]; then
        print_error "AUTHORIZED_BACKEND not set in .env file"
        exit 1
    fi
    
    print_success "Prerequisites check passed"
}

deploy_to_network() {
    local network=$1
    local rpc_var=${NETWORKS[$network]}
    local rpc_url=${!rpc_var}
    local chain_id=${CHAIN_IDS[$network]}
    
    if [ -z "$rpc_url" ]; then
        print_warning "RPC URL not set for $network (${rpc_var}), skipping..."
        return 1
    fi
    
    print_info "Deploying to $network (Chain ID: $chain_id)..."
    
    # Deploy contract
    if forge script $SCRIPT_PATH \
        --rpc-url "$rpc_url" \
        --broadcast \
        --verify \
        --etherscan-api-key "${network//-/_}_API_KEY" 2>/dev/null || \
       forge script $SCRIPT_PATH \
        --rpc-url "$rpc_url" \
        --broadcast; then
        
        print_success "Successfully deployed to $network"
        return 0
    else
        print_error "Failed to deploy to $network"
        return 1
    fi
}

list_supported_networks() {
    print_info "Supported networks:"
    for network in "${!NETWORKS[@]}"; do
        echo "  - $network (Chain ID: ${CHAIN_IDS[$network]})"
    done
}

deploy_all() {
    local success_count=0
    local total_count=${#NETWORKS[@]}
    local failed_networks=()
    
    print_info "Starting deployment to all supported networks..."
    
    for network in "${!NETWORKS[@]}"; do
        if deploy_to_network "$network"; then
            ((success_count++))
            sleep 3  # Rate limiting
        else
            failed_networks+=("$network")
        fi
        echo ""
    done
    
    print_info "Deployment Summary:"
    print_success "Successful deployments: $success_count/$total_count"
    
    if [ ${#failed_networks[@]} -gt 0 ]; then
        print_warning "Failed deployments:"
        for network in "${failed_networks[@]}"; do
            echo "  - $network"
        done
    fi
}

deploy_specific() {
    local target_network=$1
    
    if [ -z "$target_network" ]; then
        print_error "Please specify a network to deploy to"
        list_supported_networks
        exit 1
    fi
    
    if [ -z "${NETWORKS[$target_network]}" ]; then
        print_error "Unsupported network: $target_network"
        list_supported_networks
        exit 1
    fi
    
    deploy_to_network "$target_network"
}

show_help() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  all                    Deploy to all supported networks"
    echo "  deploy <network>       Deploy to specific network"
    echo "  list                   List all supported networks"
    echo "  check                  Check prerequisites only"
    echo "  help                   Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 all                           # Deploy to all networks"
    echo "  $0 deploy ethereum-sepolia       # Deploy to Ethereum Sepolia only"
    echo "  $0 deploy base-sepolia           # Deploy to Base Sepolia only"
    echo "  $0 list                          # List supported networks"
    echo ""
    echo "Environment Variables Required:"
    echo "  PRIVATE_KEY                      # Your deployment private key"
    echo "  AUTHORIZED_BACKEND               # Address of authorized backend"
    echo "  <NETWORK>_RPC_URL               # RPC URL for each network"
    echo ""
    echo "Optional Environment Variables:"
    echo "  <NETWORK>_API_KEY               # Etherscan API key for verification"
}

# Main script logic
main() {
    print_header
    
    case "${1:-help}" in
        "all")
            check_prerequisites
            deploy_all
            ;;
        "deploy")
            check_prerequisites
            deploy_specific "$2"
            ;;
        "list")
            list_supported_networks
            ;;
        "check")
            check_prerequisites
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# Trap to handle script interruption
trap 'print_warning "Deployment interrupted by user"; exit 1' INT

# Run main function with all arguments
main "$@" 
