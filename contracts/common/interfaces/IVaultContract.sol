// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

interface IVaultContract {
    // Structs
    struct GameVaultInfo {
        uint256 totalStakes;
    }

    struct ChainConfig {
        uint32 domainId;
        address usdcAddress;
        address tokenMessenger;
        address messageTransmitter;
        bool isSupported;
    }

    // Events
    event StakeDeposited(
        address indexed playerAddress,
        uint256 indexed gameId,
        uint256 amount,
        uint256 newTotal
    );

    event RewardsTransferred(uint256 indexed gameId, uint256 amount);

    event CrossChainTransferInitiated(
        uint256 indexed gameId,
        uint256 amount,
        uint32 destinationDomain,
        uint64 nonce
    );

    // Core Functions
    function stake(
        address playerAddress,
        uint256 gameId,
        uint256 amount
    ) external;

    function transferRewardsCrossChain(
        uint256 gameId,
        uint256 amount,
        uint256 destinationChainId,
        address recipient,
        bool useFastTransfer,
        uint256 maxFee
    ) external;

    // View Functions
    function getTotalStakes() external view returns (uint256);

    function getUsdcContractAddress() external view returns (address);

    function getAuthorizedBackend() external view returns (address);

    function getTokenMessengerV2() external view returns (address);

    function getMessageTransmitterV2() external view returns (address);

    function getChainConfig(
        uint256 chainId
    ) external view returns (ChainConfig memory);

    function isChainSupported(uint256 chainId) external view returns (bool);
}
