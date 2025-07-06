#!/bin/bash

# Create contracts-bindings directory
mkdir -p contracts-bindings/gamecontract
mkdir -p contracts-bindings/vaultcontract

# Extract ABI from compiled contracts and generate Go bindings
echo "Generating GameContract bindings..."
# Extract ABI and bytecode from the compiled JSON
jq '.abi' contracts/out/GameContract.sol/GameContract.json > /tmp/gamecontract_abi.json
jq -r '.bytecode.object' contracts/out/GameContract.sol/GameContract.json > /tmp/gamecontract_bin.txt
abigen --abi /tmp/gamecontract_abi.json --bin /tmp/gamecontract_bin.txt --pkg gamecontract --out contracts-bindings/gamecontract/gamecontract.go

echo "Generating VaultContract bindings..."
# Extract ABI and bytecode from the compiled JSON
jq '.abi' contracts/out/VaultContract.sol/VaultContract.json > /tmp/vaultcontract_abi.json
jq -r '.bytecode.object' contracts/out/VaultContract.sol/VaultContract.json > /tmp/vaultcontract_bin.txt
abigen --abi /tmp/vaultcontract_abi.json --bin /tmp/vaultcontract_bin.txt --pkg vaultcontract --out contracts-bindings/vaultcontract/vaultcontract.go

echo "Generating FactoryContract bindings..."
# Extract ABI and bytecode from the compiled JSON
jq '.abi' contracts/out/GameFactory.sol/GameFactory.json > /tmp/factorycontract_abi.json
jq -r '.bytecode.object' contracts/out/GameFactory.sol/GameFactory.json > /tmp/factorycontract_bin.txt
abigen --abi /tmp/factorycontract_abi.json --bin /tmp/factorycontract_bin.txt --pkg factorycontract --out contracts-bindings/factorycontract/factorycontract.go


echo "🔧 Generating Permit2 contract bindings..."
mkdir -p contracts-bindings/permit2
abigen --abi contracts/abi/Permit2.abi \
       --pkg permit2 \
       --type Permit2 \
       --out contracts-bindings/permit2/permit2.go

echo "🔧 Generating USDC contract bindings..."
mkdir -p contracts-bindings/usdc
abigen --abi contracts/abi/USDC.abi \
       --pkg usdc \
       --type USDC \
       --out contracts-bindings/usdc/usdc.go

# Clean up temporary files
rm -f /tmp/gamecontract_abi.json /tmp/vaultcontract_abi.json /tmp/gamecontract_bin.txt /tmp/vaultcontract_bin.txt

echo "Go bindings generated successfully!"
echo "Files created:"
echo "- contracts-bindings/gamecontract/gamecontract.go"
echo "- contracts-bindings/vaultcontract/vaultcontract.go"   
