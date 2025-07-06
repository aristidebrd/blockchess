#!/bin/bash

# credit_usdc.sh - Script to credit USDC to Anvil accounts using CreditUSDC.sol on multiple chains

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration - Multi-chain Anvil setup
ANVIL_PORT_BASE=8545
ANVIL_PORT_OPTIMISM=8546
ANVIL_RPC_URL_BASE="http://127.0.0.1:$ANVIL_PORT_BASE"
ANVIL_RPC_URL_OPTIMISM="http://127.0.0.1:$ANVIL_PORT_OPTIMISM"

# Real network RPC URLs for forking
BASE_SEPOLIA_RPC_URL="http://127.0.0.1:$ANVIL_PORT_BASE"
OPTIMISM_SEPOLIA_RPC_URL="http://127.0.0.1:$ANVIL_PORT_OPTIMISM"

# USDC contract addresses for each chain
USDC_CONTRACT_BASE="0x036CbD53842c5426634e7929541eC2318f3dCF7e"
USDC_CONTRACT_OPTIMISM="0x5fd84259d66Cd46123540766Be93DFE6D43130D7"

# Chain IDs
CHAIN_ID_BASE=84532
CHAIN_ID_OPTIMISM=11155420

# Script contract
SCRIPT_CONTRACT_BASE="script/CreditUSDC_base.sol:CreditRealUSDC"
SCRIPT_CONTRACT_OPTIMISM="script/CreditUSDC_op.sol:CreditRealUSDC"

# Default Anvil accounts
ACCOUNT1="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
ACCOUNT2="0x70997970C51812dc3A010C7d01b50e0d17dc79C8"

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if Anvil is running on a specific port
check_anvil_running() {
    local rpc_url=$1
    if curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        $rpc_url >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Function to wait for Anvil to be ready
wait_for_anvil() {
    local rpc_url=$1
    local chain_name=$2
    local max_attempts=30
    local attempt=0
    
    print_info "Waiting for $chain_name Anvil to be ready..."
    
    while [ $attempt -lt $max_attempts ]; do
        if check_anvil_running $rpc_url; then
            print_success "$chain_name Anvil is ready!"
            return 0
        fi
        
        sleep 2
        attempt=$((attempt + 1))
        echo -n "."
    done
    
    print_error "$chain_name Anvil failed to start after $max_attempts attempts"
    return 1
}

# Function to start Base Sepolia Anvil
start_base_anvil() {
    print_info "Starting Base Sepolia Anvil fork..."
    
    # Check if Anvil is already running
    if check_anvil_running $ANVIL_RPC_URL_BASE; then
        print_warning "Base Anvil is already running on port $ANVIL_PORT_BASE"
        return 0
    fi
    
    # Start Anvil in background
    anvil --fork-url $BASE_SEPOLIA_RPC_URL --port $ANVIL_PORT_BASE --chain-id $CHAIN_ID_BASE > anvil_base.log 2>&1 &
    BASE_ANVIL_PID=$!
    
    # Wait for Anvil to be ready
    if wait_for_anvil $ANVIL_RPC_URL_BASE "Base Sepolia"; then
        print_success "Base Sepolia Anvil started successfully (PID: $BASE_ANVIL_PID)"
        echo $BASE_ANVIL_PID > anvil_base.pid
        return 0
    else
        print_error "Failed to start Base Sepolia Anvil"
        return 1
    fi
}

# Function to start Optimism Sepolia Anvil
start_optimism_anvil() {
    print_info "Starting Optimism Sepolia Anvil fork..."
    
    # Check if Anvil is already running
    if check_anvil_running $ANVIL_RPC_URL_OPTIMISM; then
        print_warning "Optimism Anvil is already running on port $ANVIL_PORT_OPTIMISM"
        return 0
    fi
    
    # Start Anvil in background
    anvil --fork-url $OPTIMISM_SEPOLIA_RPC_URL --port $ANVIL_PORT_OPTIMISM --chain-id $CHAIN_ID_OPTIMISM > anvil_optimism.log 2>&1 &
    OPTIMISM_ANVIL_PID=$!
    
    # Wait for Anvil to be ready
    if wait_for_anvil $ANVIL_RPC_URL_OPTIMISM "Optimism Sepolia"; then
        print_success "Optimism Sepolia Anvil started successfully (PID: $OPTIMISM_ANVIL_PID)"
        echo $OPTIMISM_ANVIL_PID > anvil_optimism.pid
        return 0
    else
        print_error "Failed to start Optimism Sepolia Anvil"
        return 1
    fi
}

# Function to start both Anvil instances
start_both_anvils() {
    print_info "Starting both Anvil instances..."
    
    local base_success=0
    local optimism_success=0
    
    # Start Base Anvil
    if start_base_anvil; then
        base_success=1
    fi
    
    # Start Optimism Anvil
    if start_optimism_anvil; then
        optimism_success=1
    fi
    
    if [ $base_success -eq 1 ] && [ $optimism_success -eq 1 ]; then
        print_success "Both Anvil instances started successfully!"
        return 0
    else
        print_error "Failed to start one or both Anvil instances"
        return 1
    fi
}

# Function to stop Anvil instances
stop_anvils() {
    print_info "Stopping Anvil instances..."
    
    # Stop Base Anvil
    if [ -f anvil_base.pid ]; then
        local base_pid=$(cat anvil_base.pid)
        print_info "Stopping Base Anvil (PID: $base_pid)..."
        kill $base_pid 2>/dev/null || true
        rm -f anvil_base.pid
        print_success "Base Anvil stopped"
    fi
    
    # Stop Optimism Anvil
    if [ -f anvil_optimism.pid ]; then
        local optimism_pid=$(cat anvil_optimism.pid)
        print_info "Stopping Optimism Anvil (PID: $optimism_pid)..."
        kill $optimism_pid 2>/dev/null || true
        rm -f anvil_optimism.pid
        print_success "Optimism Anvil stopped"
    fi
}

# Function to deploy CreditUSDC contract on a specific chain
deploy_credit_usdc_chain() {
    local chain_name=$1
    local rpc_url=$2
    local usdc_contract=$3
    
    print_info "Deploying CreditUSDC contract on $chain_name..."
    
    cd contracts
    
    # Set environment variable for the specific USDC contract
    export USDC_CONTRACT_ADDRESS=$usdc_contract

    if [ "$chain_name" == "Base Sepolia" ]; then
        SCRIPT_CONTRACT=$SCRIPT_CONTRACT_BASE
    elif [ "$chain_name" == "Optimism Sepolia" ]; then
        SCRIPT_CONTRACT=$SCRIPT_CONTRACT_OPTIMISM
    fi

    # Run the forge script
    if forge script $SCRIPT_CONTRACT \
        --rpc-url $rpc_url \
        -vvv; then
        print_success "CreditUSDC contract deployed successfully on $chain_name!"
        cd ..
        return 0
    else
        print_error "Failed to deploy CreditUSDC contract on $chain_name"
        cd ..
        return 1
    fi


}

# Function to deploy CreditUSDC contracts on both chains
deploy_credit_usdc_both() {
    print_info "Deploying CreditUSDC contracts on both chains..."
    
    local base_success=0
    local optimism_success=0
    
    # Deploy on Base Sepolia
    if deploy_credit_usdc_chain "Base Sepolia" $ANVIL_RPC_URL_BASE $USDC_CONTRACT_BASE; then
        base_success=1
    fi
    
    # Deploy on Optimism Sepolia
    if deploy_credit_usdc_chain "Optimism Sepolia" $ANVIL_RPC_URL_OPTIMISM $USDC_CONTRACT_OPTIMISM; then
        optimism_success=1
    fi
    
    if [ $base_success -eq 1 ] && [ $optimism_success -eq 1 ]; then
        print_success "CreditUSDC contracts deployed successfully on both chains!"
        return 0
    else
        print_error "Failed to deploy CreditUSDC contracts on one or both chains"
        return 1
    fi
}

# Function to check USDC balances on a specific chain
check_usdc_balances_chain() {
    local chain_name=$1
    local rpc_url=$2
    local usdc_contract=$3
    
    print_info "Checking USDC balances on $chain_name..."
    
    echo ""
    echo "$chain_name USDC Balances:"
    echo "$(printf '=%.0s' {1..30})"
    
    # Check balance for account1
    local balance1_hex=$(cast call $usdc_contract \
        "balanceOf(address)" $ACCOUNT1 \
        --rpc-url $rpc_url 2>/dev/null || echo "0x0")
    
    # Clean up hex value and convert to decimal
    local balance1_clean=$(echo ${balance1_hex} | sed 's/^0x0*/0x/' | sed 's/^0x$/0x0/')
    local balance1_dec=$((${balance1_clean}))
    
    # Convert from wei to USDC (6 decimals)
    local balance1_formatted=$(echo "scale=6; $balance1_dec / 1000000" | bc -l 2>/dev/null || echo "0")
    
    # Check balance for account2
    local balance2_hex=$(cast call $usdc_contract \
        "balanceOf(address)" $ACCOUNT2 \
        --rpc-url $rpc_url 2>/dev/null || echo "0x0")
    
    # Clean up hex value and convert to decimal
    local balance2_clean=$(echo ${balance2_hex} | sed 's/^0x0*/0x/' | sed 's/^0x$/0x0/')
    local balance2_dec=$((${balance2_clean}))
    
    # Convert from wei to USDC (6 decimals)
    local balance2_formatted=$(echo "scale=6; $balance2_dec / 1000000" | bc -l 2>/dev/null || echo "0")
    
    echo "Account 1 ($ACCOUNT1): $balance1_formatted USDC"
    echo "Account 2 ($ACCOUNT2): $balance2_formatted USDC"
    echo ""
}

# Function to check USDC balances on both chains
check_usdc_balances_both() {
    print_info "Checking USDC balances on both chains..."
    
    # Check Base Sepolia balances
    if check_anvil_running $ANVIL_RPC_URL_BASE; then
        check_usdc_balances_chain "Base Sepolia" $ANVIL_RPC_URL_BASE $USDC_CONTRACT_BASE
    else
        print_warning "Base Sepolia Anvil is not running - skipping balance check"
    fi
    
    # Check Optimism Sepolia balances
    if check_anvil_running $ANVIL_RPC_URL_OPTIMISM; then
        check_usdc_balances_chain "Optimism Sepolia" $ANVIL_RPC_URL_OPTIMISM $USDC_CONTRACT_OPTIMISM
    else
        print_warning "Optimism Sepolia Anvil is not running - skipping balance check"
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTION] [CHAIN]"
    echo ""
    echo "Options:"
    echo "  start [chain]     Start Anvil(s) and deploy CreditUSDC contract(s)"
    echo "  stop              Stop all Anvil instances"
    echo "  deploy [chain]    Deploy CreditUSDC contract(s) (assumes Anvil is running)"
    echo "  balance [chain]   Check USDC balances"
    echo "  status            Check if Anvil instances are running"
    echo "  help              Show this help message"
    echo ""
    echo "Chains:"
    echo "  base              Base Sepolia only"
    echo "  optimism          Optimism Sepolia only"
    echo "  both              Both chains (default)"
    echo ""
    echo "Examples:"
    echo "  $0 start          # Start both Anvil instances and credit USDC"
    echo "  $0 start base     # Start only Base Sepolia Anvil"
    echo "  $0 balance        # Check USDC balances on both chains"
    echo "  $0 balance base   # Check USDC balances on Base Sepolia only"
    echo "  $0 stop           # Stop all Anvil instances"
}

# Function to check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."
    
    local missing_deps=()
    
    if ! command_exists anvil; then
        missing_deps+=("anvil (install Foundry)")
    fi
    
    if ! command_exists forge; then
        missing_deps+=("forge (install Foundry)")
    fi
    
    if ! command_exists cast; then
        missing_deps+=("cast (install Foundry)")
    fi
    
    if ! command_exists bc; then
        missing_deps+=("bc (install with: sudo apt-get install bc)")
    fi
    
    if ! command_exists curl; then
        missing_deps+=("curl")
    fi
    
    if [ ${#missing_deps[@]} -gt 0 ]; then
        print_error "Missing dependencies:"
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        echo ""
        echo "Please install missing dependencies and try again."
        return 1
    fi
    
    print_success "All prerequisites met!"
    return 0
}

# Function to show status
show_status() {
    echo ""
    echo "🔍 Anvil Status Check"
    echo "===================="
    
    # Check Base Sepolia Anvil
    if check_anvil_running $ANVIL_RPC_URL_BASE; then
        print_success "Base Sepolia Anvil is running on port $ANVIL_PORT_BASE"
    else
        print_info "Base Sepolia Anvil is not running"
    fi
    
    # Check Optimism Sepolia Anvil
    if check_anvil_running $ANVIL_RPC_URL_OPTIMISM; then
        print_success "Optimism Sepolia Anvil is running on port $ANVIL_PORT_OPTIMISM"
    else
        print_info "Optimism Sepolia Anvil is not running"
    fi
    
    # Show balances if any Anvil is running
    if check_anvil_running $ANVIL_RPC_URL_BASE || check_anvil_running $ANVIL_RPC_URL_OPTIMISM; then
        check_usdc_balances_both
    fi
}

# Main function
main() {
    local action="${1:-start}"
    local chain="${2:-both}"
    
    case "$action" in
        "start")
            if ! check_prerequisites; then
                exit 1
            fi
            
            case "$chain" in
                "base")
                    if start_base_anvil; then
                        sleep 2
                        if deploy_credit_usdc_chain "Base Sepolia" $ANVIL_RPC_URL_BASE $USDC_CONTRACT_BASE; then
                            check_usdc_balances_chain "Base Sepolia" $ANVIL_RPC_URL_BASE $USDC_CONTRACT_BASE
                            print_success "Base Sepolia USDC crediting completed successfully!"
                        else
                            print_error "Failed to credit USDC on Base Sepolia"
                            exit 1
                        fi
                    else
                        print_error "Failed to start Base Sepolia Anvil"
                        exit 1
                    fi
                    ;;
                "optimism")
                    if start_optimism_anvil; then
                        sleep 2
                        if deploy_credit_usdc_chain "Optimism Sepolia" $ANVIL_RPC_URL_OPTIMISM $USDC_CONTRACT_OPTIMISM; then
                            check_usdc_balances_chain "Optimism Sepolia" $ANVIL_RPC_URL_OPTIMISM $USDC_CONTRACT_OPTIMISM
                            print_success "Optimism Sepolia USDC crediting completed successfully!"
                        else
                            print_error "Failed to credit USDC on Optimism Sepolia"
                            exit 1
                        fi
                    else
                        print_error "Failed to start Optimism Sepolia Anvil"
                        exit 1
                    fi
                    ;;
                "both"|*)
                    if start_both_anvils; then
                        sleep 2
                        if deploy_credit_usdc_both; then
                            check_usdc_balances_both
                            print_success "Multi-chain USDC crediting completed successfully!"
                            print_info "Both Anvil instances are running in the background."
                            print_info "Use '$0 stop' to stop them."
                        else
                            print_error "Failed to credit USDC on one or both chains"
                            exit 1
                        fi
                    else
                        print_error "Failed to start Anvil instances"
                        exit 1
                    fi
                    ;;
            esac
            ;;
        "stop")
            stop_anvils
            ;;
        "deploy")
            if ! check_prerequisites; then
                exit 1
            fi
            
            case "$chain" in
                "base")
                    if ! check_anvil_running $ANVIL_RPC_URL_BASE; then
                        print_error "Base Sepolia Anvil is not running. Start it first with '$0 start base'"
                        exit 1
                    fi
                    deploy_credit_usdc_chain "Base Sepolia" $ANVIL_RPC_URL_BASE $USDC_CONTRACT_BASE
                    check_usdc_balances_chain "Base Sepolia" $ANVIL_RPC_URL_BASE $USDC_CONTRACT_BASE
                    ;;
                "optimism")
                    if ! check_anvil_running $ANVIL_RPC_URL_OPTIMISM; then
                        print_error "Optimism Sepolia Anvil is not running. Start it first with '$0 start optimism'"
                        exit 1
                    fi
                    deploy_credit_usdc_chain "Optimism Sepolia" $ANVIL_RPC_URL_OPTIMISM $USDC_CONTRACT_OPTIMISM
                    check_usdc_balances_chain "Optimism Sepolia" $ANVIL_RPC_URL_OPTIMISM $USDC_CONTRACT_OPTIMISM
                    ;;
                "both"|*)
                    if ! check_anvil_running $ANVIL_RPC_URL_BASE && ! check_anvil_running $ANVIL_RPC_URL_OPTIMISM; then
                        print_error "No Anvil instances are running. Start them first with '$0 start'"
                        exit 1
                    fi
                    deploy_credit_usdc_both
                    check_usdc_balances_both
                    ;;
            esac
            ;;
        "balance")
            case "$chain" in
                "base")
                    if ! check_anvil_running $ANVIL_RPC_URL_BASE; then
                        print_error "Base Sepolia Anvil is not running"
                        exit 1
                    fi
                    check_usdc_balances_chain "Base Sepolia" $ANVIL_RPC_URL_BASE $USDC_CONTRACT_BASE
                    ;;
                "optimism")
                    if ! check_anvil_running $ANVIL_RPC_URL_OPTIMISM; then
                        print_error "Optimism Sepolia Anvil is not running"
                        exit 1
                    fi
                    check_usdc_balances_chain "Optimism Sepolia" $ANVIL_RPC_URL_OPTIMISM $USDC_CONTRACT_OPTIMISM
                    ;;
                "both"|*)
                    check_usdc_balances_both
                    ;;
            esac
            ;;
        "status")
            show_status
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        *)
            print_error "Unknown option: $action"
            show_usage
            exit 1
            ;;
    esac
}

# Cleanup function
cleanup() {
    print_info "Cleaning up..."
    # Don't automatically stop Anvil on script exit unless explicitly requested
}

# Set up trap for cleanup
trap cleanup EXIT

# Run main function
main "$@" 