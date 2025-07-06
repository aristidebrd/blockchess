// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {IVaultContract} from "../common/interfaces/IVaultContract.sol";

// USDC Token Interface
interface IERC20 {
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
    function allowance(
        address owner,
        address spender
    ) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
}

// Circle CCTP V2 TokenMessenger Interface
interface ITokenMessengerV2 {
    function depositForBurn(
        uint256 amount,
        uint32 destinationDomain,
        bytes32 mintRecipient,
        address burnToken,
        bytes32 destinationCaller,
        uint256 maxFee,
        uint32 minFinalityThreshold
    ) external returns (uint64 nonce);

    function depositForBurnWithHook(
        uint256 amount,
        uint32 destinationDomain,
        bytes32 mintRecipient,
        address burnToken,
        bytes32 destinationCaller,
        uint256 maxFee,
        uint32 minFinalityThreshold,
        bytes calldata hookData
    ) external returns (uint64 nonce);
}

// Circle CCTP V2 MessageTransmitter Interface
interface IMessageTransmitterV2 {
    function receiveMessage(
        bytes calldata message,
        bytes calldata attestation
    ) external returns (bool success);
}

contract VaultContract is IVaultContract {
    address public immutable AUTHORIZED_BACKEND;
    address public immutable USDC_CONTRACT_ADDRESS;
    address public immutable TOKEN_MESSENGER_V2;
    address public immutable MESSAGE_TRANSMITTER_V2;

    uint256 public totalStakes;

    // Chain configuration - using the one from interface

    mapping(uint256 => ChainConfig) public chainConfigs;

    // Events - using the ones from interface

    modifier onlyAuthorizedBackend() {
        require(
            msg.sender == AUTHORIZED_BACKEND,
            "Only authorized backend can transfer rewards"
        );
        _;
    }

    constructor(
        address _authorizedBackend,
        address _usdcContractAddress,
        address _tokenMessengerV2,
        address _messageTransmitterV2
    ) {
        require(
            _authorizedBackend != address(0),
            "Authorized backend cannot be zero address"
        );
        require(
            _usdcContractAddress != address(0),
            "USDC contract address cannot be zero address"
        );
        require(
            _tokenMessengerV2 != address(0),
            "TokenMessenger V2 cannot be zero address"
        );
        require(
            _messageTransmitterV2 != address(0),
            "MessageTransmitter V2 cannot be zero address"
        );

        AUTHORIZED_BACKEND = _authorizedBackend;
        USDC_CONTRACT_ADDRESS = _usdcContractAddress;
        TOKEN_MESSENGER_V2 = _tokenMessengerV2;
        MESSAGE_TRANSMITTER_V2 = _messageTransmitterV2;

        _initializeChainConfigs();
    }

    function _initializeChainConfigs() private {
        // Ethereum Sepolia (Domain 0)
        chainConfigs[11155111] = ChainConfig({
            domainId: 0,
            usdcAddress: 0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238,
            tokenMessenger: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitter: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            isSupported: true
        });

        // Avalanche Fuji (Domain 1)
        chainConfigs[43113] = ChainConfig({
            domainId: 1,
            usdcAddress: 0x5425890298aed601595a70AB815c96711a31Bc65,
            tokenMessenger: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitter: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            isSupported: true
        });

        // OP Sepolia (Domain 2)
        chainConfigs[11155420] = ChainConfig({
            domainId: 2,
            usdcAddress: 0x5fd84259d66Cd46123540766Be93DFE6D43130D7,
            tokenMessenger: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitter: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            isSupported: true
        });

        // Arbitrum Sepolia (Domain 3)
        chainConfigs[421614] = ChainConfig({
            domainId: 3,
            usdcAddress: 0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d,
            tokenMessenger: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitter: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            isSupported: true
        });

        // Base Sepolia (Domain 6)
        chainConfigs[84532] = ChainConfig({
            domainId: 6,
            usdcAddress: 0x036CbD53842c5426634e7929541eC2318f3dCF7e,
            tokenMessenger: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitter: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            isSupported: true
        });

        // Polygon PoS Amoy (Domain 7)
        chainConfigs[80002] = ChainConfig({
            domainId: 7,
            usdcAddress: 0x41E94Eb019C0762f9Bfcf9Fb1E58725BfB0e7582,
            tokenMessenger: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitter: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            isSupported: true
        });

        // Unichain Sepolia (Domain 10)
        chainConfigs[1301] = ChainConfig({
            domainId: 10,
            usdcAddress: 0x31d0220469e10c4E71834a79b1f276d740d3768F,
            tokenMessenger: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitter: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            isSupported: true
        });

        // Linea Sepolia (Domain 11)
        chainConfigs[59141] = ChainConfig({
            domainId: 11,
            usdcAddress: 0xFEce4462D57bD51A6A552365A011b95f0E16d9B7,
            tokenMessenger: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitter: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            isSupported: true
        });

        // Codex Testnet (Domain 12)
        chainConfigs[325000] = ChainConfig({
            domainId: 12,
            usdcAddress: 0x6d7f141b6819C2c9CC2f818e6ad549E7Ca090F8f,
            tokenMessenger: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitter: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            isSupported: true
        });

        // Sonic Testnet (Domain 13)
        chainConfigs[64165] = ChainConfig({
            domainId: 13,
            usdcAddress: 0xA4879Fed32Ecbef99399e5cbC247E533421C4eC6,
            tokenMessenger: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitter: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            isSupported: true
        });

        // World Chain Sepolia (Domain 14)
        chainConfigs[4801] = ChainConfig({
            domainId: 14,
            usdcAddress: 0x66145f38cBAC35Ca6F1Dfb4914dF98F1614aeA88,
            tokenMessenger: 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA,
            messageTransmitter: 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275,
            isSupported: true
        });
    }

    function stake(
        address playerAddress,
        uint256 gameId,
        uint256 amount
    ) external override {
        require(amount > 0, "Stake amount must be greater than 0");
        require(playerAddress != address(0), "Player address cannot be zero");

        // Transfer USDC from player to this contract
        require(
            IERC20(USDC_CONTRACT_ADDRESS).transferFrom(
                playerAddress,
                address(this),
                amount
            ),
            "USDC transfer failed"
        );

        // Update vault total stakes
        totalStakes += amount;

        emit StakeDeposited(playerAddress, gameId, amount, totalStakes);
    }

    function transferRewardsCrossChain(
        uint256 gameId,
        uint256 amount,
        uint256 destinationChainId,
        address recipient,
        bool useFastTransfer,
        uint256 maxFee
    ) external onlyAuthorizedBackend {
        require(amount > 0, "Reward amount must be greater than 0");
        require(totalStakes >= amount, "Insufficient total stakes");
        require(recipient != address(0), "Recipient cannot be zero address");

        ChainConfig memory destConfig = chainConfigs[destinationChainId];
        require(destConfig.isSupported, "Destination chain not supported");

        // Approve TokenMessenger to spend USDC
        require(
            IERC20(USDC_CONTRACT_ADDRESS).approve(TOKEN_MESSENGER_V2, amount),
            "USDC approval failed"
        );

        // Convert recipient address to bytes32
        bytes32 mintRecipient = bytes32(uint256(uint160(recipient)));

        // Set finality threshold: 1000 for fast transfer, 2000 for standard
        uint32 minFinalityThreshold = useFastTransfer ? 1000 : 2000;

        // Initiate cross-chain transfer via CCTP
        uint64 nonce = ITokenMessengerV2(TOKEN_MESSENGER_V2).depositForBurn(
            amount,
            destConfig.domainId,
            mintRecipient,
            USDC_CONTRACT_ADDRESS,
            bytes32(0), // Any destination caller
            maxFee,
            minFinalityThreshold
        );

        // Update total stakes
        totalStakes -= amount;

        emit CrossChainTransferInitiated(
            gameId,
            amount,
            destConfig.domainId,
            nonce
        );
    }

    function getTotalStakes() external view override returns (uint256) {
        return totalStakes;
    }

    // Additional view functions for multi-chain deployment info
    function getUsdcContractAddress() external view returns (address) {
        return USDC_CONTRACT_ADDRESS;
    }

    function getAuthorizedBackend() external view returns (address) {
        return AUTHORIZED_BACKEND;
    }

    function getTokenMessengerV2() external view returns (address) {
        return TOKEN_MESSENGER_V2;
    }

    function getMessageTransmitterV2() external view returns (address) {
        return MESSAGE_TRANSMITTER_V2;
    }

    function getChainConfig(
        uint256 chainId
    ) external view returns (ChainConfig memory) {
        return chainConfigs[chainId];
    }

    function isChainSupported(uint256 chainId) external view returns (bool) {
        return chainConfigs[chainId].isSupported;
    }
}
