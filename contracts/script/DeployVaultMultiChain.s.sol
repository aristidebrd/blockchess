// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import "forge-std/Script.sol";
import "../src/VaultContract.sol";

contract DeployVaultMultiChain is Script {
    // Network configurations
    struct NetworkConfig {
        string name;
        uint256 chainId;
        address usdcAddress;
        address tokenMessengerV2;
        address messageTransmitterV2;
        uint32 domainId;
    }

    mapping(uint256 => NetworkConfig) public networkConfigs;

    function setUp() public {
        _initializeNetworkConfigs();
    }

    function _initializeNetworkConfigs() private {
        // Ethereum Sepolia
        networkConfigs[11155111] = NetworkConfig({
            name: "Ethereum Sepolia",
            chainId: 11155111,
            usdcAddress: 0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238,
            tokenMessengerV2: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitterV2: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            domainId: 0
        });

        // Avalanche Fuji
        networkConfigs[43113] = NetworkConfig({
            name: "Avalanche Fuji",
            chainId: 43113,
            usdcAddress: 0x5425890298aed601595a70AB815c96711a31Bc65,
            tokenMessengerV2: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitterV2: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            domainId: 1
        });

        // OP Sepolia
        networkConfigs[11155420] = NetworkConfig({
            name: "OP Sepolia",
            chainId: 11155420,
            usdcAddress: 0x5fd84259d66Cd46123540766Be93DFE6D43130D7,
            tokenMessengerV2: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitterV2: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            domainId: 2
        });

        // Arbitrum Sepolia
        networkConfigs[421614] = NetworkConfig({
            name: "Arbitrum Sepolia",
            chainId: 421614,
            usdcAddress: 0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d,
            tokenMessengerV2: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitterV2: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            domainId: 3
        });

        // Base Sepolia
        networkConfigs[84532] = NetworkConfig({
            name: "Base Sepolia",
            chainId: 84532,
            usdcAddress: 0x036CbD53842c5426634e7929541eC2318f3dCF7e,
            tokenMessengerV2: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitterV2: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            domainId: 6
        });

        // Polygon PoS Amoy
        networkConfigs[80002] = NetworkConfig({
            name: "Polygon PoS Amoy",
            chainId: 80002,
            usdcAddress: 0x41E94Eb019C0762f9Bfcf9Fb1E58725BfB0e7582,
            tokenMessengerV2: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitterV2: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            domainId: 7
        });

        // Unichain Sepolia
        networkConfigs[1301] = NetworkConfig({
            name: "Unichain Sepolia",
            chainId: 1301,
            usdcAddress: 0x31d0220469e10c4E71834a79b1f276d740d3768F,
            tokenMessengerV2: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitterV2: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            domainId: 10
        });

        // Linea Sepolia
        networkConfigs[59141] = NetworkConfig({
            name: "Linea Sepolia",
            chainId: 59141,
            usdcAddress: 0xFEce4462D57bD51A6A552365A011b95f0E16d9B7,
            tokenMessengerV2: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitterV2: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            domainId: 11
        });

        // Codex Testnet
        networkConfigs[325000] = NetworkConfig({
            name: "Codex Testnet",
            chainId: 325000,
            usdcAddress: 0x6d7f141b6819C2c9CC2f818e6ad549E7Ca090F8f,
            tokenMessengerV2: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitterV2: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            domainId: 12
        });

        // Sonic Testnet
        networkConfigs[64165] = NetworkConfig({
            name: "Sonic Testnet",
            chainId: 64165,
            usdcAddress: 0xA4879Fed32Ecbef99399e5cbC247E533421C4eC6,
            tokenMessengerV2: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitterV2: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            domainId: 13
        });

        // World Chain Sepolia
        networkConfigs[4801] = NetworkConfig({
            name: "World Chain Sepolia",
            chainId: 4801,
            usdcAddress: 0x66145f38cBAC35Ca6F1Dfb4914dF98F1614aeA88,
            tokenMessengerV2: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitterV2: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            domainId: 14
        });
    }

    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address authorizedBackend = vm.envAddress("GO_BACKEND_ADDRESS");

        uint256 currentChainId = block.chainid;
        NetworkConfig memory config = networkConfigs[currentChainId];

        require(config.chainId != 0, "Unsupported chain");

        console.log("Deploying VaultContractWithCCTP on", config.name);
        console.log("Chain ID:", config.chainId);
        console.log("USDC Address:", config.usdcAddress);
        console.log("TokenMessenger V2:", config.tokenMessengerV2);
        console.log("MessageTransmitter V2:", config.messageTransmitterV2);
        console.log("Domain ID:", config.domainId);
        console.log("Authorized Backend:", authorizedBackend);

        vm.startBroadcast(deployerPrivateKey);

        VaultContract vault = new VaultContract(
            authorizedBackend,
            config.usdcAddress,
            config.tokenMessengerV2,
            config.messageTransmitterV2
        );

        vm.stopBroadcast();

        console.log("VaultContract deployed at:", address(vault));

        // Verify deployment
        console.log("Verifying deployment...");
        require(
            vault.getAuthorizedBackend() == authorizedBackend,
            "Backend mismatch"
        );
        require(
            vault.getUsdcContractAddress() == config.usdcAddress,
            "USDC address mismatch"
        );
        require(
            vault.getTokenMessengerV2() == config.tokenMessengerV2,
            "TokenMessenger mismatch"
        );
        require(
            vault.getMessageTransmitterV2() == config.messageTransmitterV2,
            "MessageTransmitter mismatch"
        );
        require(
            vault.isChainSupported(config.chainId),
            "Chain not supported in config"
        );

        console.log("Deployment verified successfully!");

        // Output deployment info for documentation
        _outputDeploymentInfo(config.name, address(vault), config);
    }

    function deployToSpecificChain(
        uint256 chainId,
        address authorizedBackend
    ) external {
        NetworkConfig memory config = networkConfigs[chainId];
        require(config.chainId != 0, "Unsupported chain");

        console.log("Deploying VaultContract on", config.name);

        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        VaultContract vault = new VaultContract(
            authorizedBackend,
            config.usdcAddress,
            config.tokenMessengerV2,
            config.messageTransmitterV2
        );

        vm.stopBroadcast();

        console.log("VaultContract deployed at:", address(vault));
        _outputDeploymentInfo(config.name, address(vault), config);
    }

    function _outputDeploymentInfo(
        string memory networkName,
        address vaultAddress,
        NetworkConfig memory config
    ) private view {
        console.log("\n=== DEPLOYMENT SUMMARY ===");
        console.log("Network:", networkName);
        console.log("Chain ID:", config.chainId);
        console.log("Domain ID:", config.domainId);
        console.log("Vault Address:", vaultAddress);
        console.log("USDC Address:", config.usdcAddress);
        console.log("TokenMessenger V2:", config.tokenMessengerV2);
        console.log("MessageTransmitter V2:", config.messageTransmitterV2);
        console.log("========================\n");
    }

    // Helper function to get network configuration
    function getNetworkConfig(
        uint256 chainId
    ) external view returns (NetworkConfig memory) {
        return networkConfigs[chainId];
    }

    // Function to list all supported networks
    function listSupportedNetworks() external view {
        console.log("=== SUPPORTED TESTNET NETWORKS ===");

        uint256[] memory chainIds = new uint256[](11);
        chainIds[0] = 11155111; // Ethereum Sepolia
        chainIds[1] = 43113; // Avalanche Fuji
        chainIds[2] = 11155420; // OP Sepolia
        chainIds[3] = 421614; // Arbitrum Sepolia
        chainIds[4] = 84532; // Base Sepolia
        chainIds[5] = 80002; // Polygon PoS Amoy
        chainIds[6] = 1301; // Unichain Sepolia
        chainIds[7] = 59141; // Linea Sepolia
        chainIds[8] = 325000; // Codex Testnet
        chainIds[9] = 64165; // Sonic Testnet
        chainIds[10] = 4801; // World Chain Sepolia

        for (uint256 i = 0; i < chainIds.length; i++) {
            NetworkConfig memory config = networkConfigs[chainIds[i]];
            console.log(
                string(
                    abi.encodePacked(
                        config.name,
                        " (Chain ID: ",
                        vm.toString(config.chainId),
                        ", Domain: ",
                        vm.toString(config.domainId),
                        ")"
                    )
                )
            );
        }
        console.log("==================================");
    }
}
