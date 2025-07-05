// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import "forge-std/Script.sol";
import "../src/GameContract.sol";
import "../src/VaultContract.sol";

contract DeployVaultSystem is Script {
    // Deployment addresses will be loaded from environment
    address public goBackendAddress;

    // Deployed contract addresses
    GameContract public gameContract;
    VaultContract public vaultContract;

    function setUp() public {
        // Load Go backend address from environment
        goBackendAddress = vm.envAddress("GO_BACKEND_ADDRESS");
        require(goBackendAddress != address(0), "GO_BACKEND_ADDRESS not set");
    }

    function run() public {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        // Determine which contracts to deploy based on chain
        uint256 chainId = block.chainid;

        if (chainId == 8453 || chainId == 84532 || chainId == 31337) {
            // Base Chain - Deploy Game Contract + Vault Contract
            deployBaseChain();
        } else if (chainId == 1 || chainId == 10) {
            // Mainnet or Optimism - Deploy Vault Contract only
            deployVaultOnly();
        } else {
            revert("Unsupported chain for deployment");
        }

        vm.stopBroadcast();

        // Log deployment addresses
        logDeployment(chainId);
    }

    function deployBaseChain() internal {
        uint256 chainId = block.chainid;
        if (chainId == 8453) {
            console.log("Deploying to Base Mainnet (8453)");
        } else if (chainId == 84532) {
            console.log("Deploying to Base Sepolia (84532)");
        }

        // Deploy Game Contract
        gameContract = new GameContract(goBackendAddress);
        console.log("GameContract deployed at:", address(gameContract));

        // Deploy Vault Contract
        vaultContract = new VaultContract(
            goBackendAddress,
            address(gameContract)
        );
        console.log("VaultContract deployed at:", address(vaultContract));
    }

    function deployVaultOnly() internal {
        console.log("Deploying Vault Contract only");

        // Get Game Contract address from environment (deployed on Base)
        address gameContractAddress = vm.envAddress("GAME_CONTRACT_ADDRESS");
        require(
            gameContractAddress != address(0),
            "GAME_CONTRACT_ADDRESS not set"
        );

        // Deploy Vault Contract
        vaultContract = new VaultContract(
            goBackendAddress,
            gameContractAddress
        );
        console.log("VaultContract deployed at:", address(vaultContract));
    }

    function logDeployment(uint256 chainId) internal view {
        console.log("\n=== DEPLOYMENT SUMMARY ===");
        console.log("Chain ID:", chainId);
        console.log("Go Backend Address:", goBackendAddress);

        if (address(gameContract) != address(0)) {
            console.log("GameContract:", address(gameContract));
        }

        if (address(vaultContract) != address(0)) {
            console.log("VaultContract:", address(vaultContract));
        }

        console.log("========================\n");
    }

    // Helper function to verify deployment
    function verifyDeployment() public view {
        if (address(gameContract) != address(0)) {
            require(
                gameContract.authorizedBackend() == goBackendAddress,
                "Game contract backend mismatch"
            );
            console.log("GameContract backend verified");
        }

        if (address(vaultContract) != address(0)) {
            require(
                vaultContract.authorizedBackend() == goBackendAddress,
                "Vault contract backend mismatch"
            );
            console.log("VaultContract backend verified");
        }
    }
}
