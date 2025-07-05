#!/bin/bash

# Create backend directory
mkdir -p backend/contracts

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

# Clean up temporary files
rm -f /tmp/gamecontract_abi.json /tmp/vaultcontract_abi.json /tmp/gamecontract_bin.txt /tmp/vaultcontract_bin.txt

echo "Go bindings generated successfully!"
echo "Files created:"
echo "- backend/contracts/gamecontract.go"
echo "- backend/contracts/vaultcontract.go"
