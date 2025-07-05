#!/bin/bash

# Create backend directory
mkdir -p backend/contracts/gamecontract
mkdir -p backend/contracts/vaultcontract

# Extract ABI from compiled contracts and generate Go bindings
echo "Generating GameContract bindings..."
# Extract ABI and bytecode from the compiled JSON
jq '.abi' contracts/out/GameContract.sol/GameContract.json > /tmp/gamecontract_abi.json
jq -r '.bytecode.object' contracts/out/GameContract.sol/GameContract.json > /tmp/gamecontract_bin.txt
abigen --abi /tmp/gamecontract_abi.json --bin /tmp/gamecontract_bin.txt --pkg gamecontract --out backend/contracts/gamecontract/gamecontract.go

echo "Generating VaultContract bindings..."
# Extract ABI and bytecode from the compiled JSON
jq '.abi' contracts/out/VaultContract.sol/VaultContract.json > /tmp/vaultcontract_abi.json
jq -r '.bytecode.object' contracts/out/VaultContract.sol/VaultContract.json > /tmp/vaultcontract_bin.txt
abigen --abi /tmp/vaultcontract_abi.json --bin /tmp/vaultcontract_bin.txt --pkg vaultcontract --out backend/contracts/vaultcontract/vaultcontract.go

echo "🔧 Generating Permit2 contract bindings..."
mkdir -p backend/contracts/permit2
abigen --abi contracts/abi/Permit2.abi \
       --pkg permit2 \
       --type Permit2 \
       --out backend/contracts/permit2/permit2.go

echo "🔧 Generating USDC contract bindings..."
mkdir -p backend/contracts/usdc
abigen --abi contracts/abi/USDC.abi \
       --pkg usdc \
       --type USDC \
       --out backend/contracts/usdc/usdc.go

# Clean up temporary files
rm -f /tmp/gamecontract_abi.json /tmp/vaultcontract_abi.json /tmp/gamecontract_bin.txt /tmp/vaultcontract_bin.txt

echo "Go bindings generated successfully!"
echo "Files created:"
echo "- backend/contracts/gamecontract/gamecontract.go"
echo "- backend/contracts/vaultcontract/vaultcontract.go"   