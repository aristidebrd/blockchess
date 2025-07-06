// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import "forge-std/Script.sol";
import "../src/GameFactory.sol";

contract DeployGameFactory is Script {
    // Base Sepolia testnet configuration
    uint256 constant BASE_SEPOLIA_CHAIN_ID = 84532;

    // Configuration - Replace with your actual backend address
    address constant AUTHORIZED_BACKEND =
        0x1234567890123456789012345678901234567890; // Replace with actual address

    function run() external {
        // Get the deployer's private key from environment
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);

        console.log("Deploying GameFactory to Base Sepolia...");
        console.log("Deployer address:", deployer);
        console.log("Authorized backend:", AUTHORIZED_BACKEND);

        // Start broadcasting transactions
        vm.startBroadcast(deployerPrivateKey);

        // Deploy GameFactory
        GameFactory gameFactory = new GameFactory(AUTHORIZED_BACKEND);

        // Stop broadcasting
        vm.stopBroadcast();

        console.log("GameFactory deployed at:", address(gameFactory));
        console.log("Deployment successful!");

        // Verify the deployment
        console.log("Verifying deployment...");
        console.log(
            "Authorized backend from contract:",
            gameFactory.getAuthorizedBackend()
        );
        console.log("Next game ID:", gameFactory.getNextGameId());
        console.log("Total games count:", gameFactory.getTotalGamesCount());
    }

    function deployToBaseSepolia() external {
        require(
            block.chainid == BASE_SEPOLIA_CHAIN_ID,
            "Must be on Base Sepolia"
        );
        this.run();
    }
}
