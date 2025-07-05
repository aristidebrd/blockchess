#!/bin/bash

# credit_usdc.sh - Script to credit USDC to Anvil accounts using CreditUSDC.sol

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ANVIL_PORT=8545
ANVIL_RPC_URL="http://127.0.0.1:$ANVIL_PORT"
BASE_SEPOLIA_RPC_URL="https://sepolia.base.org"
USDC_CONTRACT="0x036CbD53842c5426634e7929541eC2318f3dCF7e"
SCRIPT_CONTRACT="script/CreditUSDC.sol:CreditRealUSDC"

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

# Function to check if Anvil is running
check_anvil_running() {
    if curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        $ANVIL_RPC_URL >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Function to wait for Anvil to be ready
wait_for_anvil() {
    local max_attempts=30
    local attempt=0
    
    print_info "Waiting for Anvil to be ready..."
    
    while [ $attempt -lt $max_attempts ]; do
        if check_anvil_running; then
            print_success "Anvil is ready!"
            return 0
        fi
        
        sleep 2
        attempt=$((attempt + 1))
        echo -n "."
    done
    
    print_error "Anvil failed to start after $max_attempts attempts"
    return 1
}

# Function to start Anvil
start_anvil() {
    print_info "Starting Anvil with Base Sepolia fork..."
    
    # Check if Anvil is already running
    if check_anvil_running; then
        print_warning "Anvil is already running on port $ANVIL_PORT"
        return 0
    fi
    
    # Start Anvil in background
    anvil --fork-url $BASE_SEPOLIA_RPC_URL --port $ANVIL_PORT > anvil.log 2>&1 &
    ANVIL_PID=$!
    
    # Wait for Anvil to be ready
    if wait_for_anvil; then
        print_success "Anvil started successfully (PID: $ANVIL_PID)"
        echo $ANVIL_PID > anvil.pid
        return 0
    else
        print_error "Failed to start Anvil"
        return 1
    fi
}

# Function to stop Anvil
stop_anvil() {
    if [ -f anvil.pid ]; then
        local pid=$(cat anvil.pid)
        print_info "Stopping Anvil (PID: $pid)..."
        kill $pid 2>/dev/null || true
        rm -f anvil.pid
        print_success "Anvil stopped"
    fi
}

# Function to deploy CreditUSDC contract
deploy_credit_usdc() {
    print_info "Deploying CreditUSDC contract..."
    
    cd contracts
    
    # Run the forge script (no --broadcast needed for local Anvil with vm.prank)
    if forge script $SCRIPT_CONTRACT \
        --rpc-url $ANVIL_RPC_URL \
        -vvv; then
        print_success "CreditUSDC contract deployed successfully!"
        cd ..
        return 0
    else
        print_error "Failed to deploy CreditUSDC contract"
        cd ..
        return 1
    fi
}

# Function to check USDC balances
check_usdc_balances() {
    print_info "Checking USDC balances..."
    
    echo ""
    echo "USDC Balances:"
    echo "=============="
    
    # Check balance for account1
    local balance1_hex=$(cast call $USDC_CONTRACT \
        "balanceOf(address)" $ACCOUNT1 \
        --rpc-url $ANVIL_RPC_URL 2>/dev/null || echo "0x0")
    
    # Clean up hex value (remove leading zeros) and convert to decimal
    local balance1_clean=$(echo ${balance1_hex} | sed 's/^0x0*/0x/' | sed 's/^0x$/0x0/')
    local balance1_dec=$((${balance1_clean}))
    
    # Convert from wei to USDC (6 decimals)
    local balance1_formatted=$(echo "scale=6; $balance1_dec / 1000000" | bc -l 2>/dev/null || echo "0")
    
    # Check balance for account2
    local balance2_hex=$(cast call $USDC_CONTRACT \
        "balanceOf(address)" $ACCOUNT2 \
        --rpc-url $ANVIL_RPC_URL 2>/dev/null || echo "0x0")
    
    # Clean up hex value (remove leading zeros) and convert to decimal
    local balance2_clean=$(echo ${balance2_hex} | sed 's/^0x0*/0x/' | sed 's/^0x$/0x0/')
    local balance2_dec=$((${balance2_clean}))
    
    # Convert from wei to USDC (6 decimals)
    local balance2_formatted=$(echo "scale=6; $balance2_dec / 1000000" | bc -l 2>/dev/null || echo "0")
    
    echo "Account 1 ($ACCOUNT1): $balance1_formatted USDC"
    echo "Account 2 ($ACCOUNT2): $balance2_formatted USDC"
    echo ""
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  start     Start Anvil and deploy CreditUSDC contract"
    echo "  stop      Stop Anvil"
    echo "  deploy    Deploy CreditUSDC contract (assumes Anvil is running)"
    echo "  balance   Check USDC balances"
    echo "  status    Check if Anvil is running"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 start    # Start Anvil and credit USDC"
    echo "  $0 balance  # Check USDC balances"
    echo "  $0 stop     # Stop Anvil"
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
    
    if [ ! -f "contracts/script/CreditUSDC.sol" ]; then
        print_error "CreditUSDC.sol not found at contracts/script/CreditUSDC.sol"
        return 1
    fi
    
    print_success "All prerequisites met!"
    return 0
}

# Function to show status
show_status() {
    if check_anvil_running; then
        print_success "Anvil is running on port $ANVIL_PORT"
        check_usdc_balances
    else
        print_info "Anvil is not running"
    fi
}

# Main function
main() {
    case "${1:-start}" in
        "start")
            if ! check_prerequisites; then
                exit 1
            fi
            
            if start_anvil; then
                sleep 2  # Give Anvil a moment to fully initialize
                if deploy_credit_usdc; then
                    check_usdc_balances
                    print_success "USDC crediting completed successfully!"
                    print_info "Anvil is running in the background. Use '$0 stop' to stop it."
                else
                    print_error "Failed to credit USDC"
                    exit 1
                fi
            else
                print_error "Failed to start Anvil"
                exit 1
            fi
            ;;
        "stop")
            stop_anvil
            ;;
        "deploy")
            if ! check_prerequisites; then
                exit 1
            fi
            
            if ! check_anvil_running; then
                print_error "Anvil is not running. Start it first with '$0 start' or run 'anvil' manually."
                exit 1
            fi
            
            deploy_credit_usdc
            check_usdc_balances
            ;;
        "balance")
            if ! check_anvil_running; then
                print_error "Anvil is not running"
                exit 1
            fi
            check_usdc_balances
            ;;
        "status")
            show_status
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        *)
            print_error "Unknown option: $1"
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