// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

interface IVaultContract {
    // Enums
    enum GameResult {
        Ongoing,
        WhiteWins,
        BlackWins,
        Draw
    }

    // Structs
    struct GameVaultInfo {
        uint256 gameId;
        uint256 totalStakes;
        uint256 playerCount;
        GameResult result;
        bool gameEnded;
        uint256 endedAt;
    }

    struct PlayerStakeInfo {
        uint256 totalStaked;
        uint256 stakeCount;
        bool hasClaimed;
    }

    // Events
    event StakeDeposited(
        uint256 indexed gameId,
        address indexed player,
        uint256 amount,
        uint256 newTotal
    );

    event GameEndedInVault(
        uint256 indexed gameId,
        GameResult result,
        uint256 totalStakes,
        uint256 endedAt
    );

    event RewardsClaimed(
        uint256 indexed gameId,
        address indexed player,
        uint256 amount
    );

    // Core Functions
    function stake(uint256 gameId, uint256 fixedStakeAmount) external;

    function endGame(uint256 gameId, GameResult result) external;

    // View Functions
    function getGameVaultInfo(
        uint256 gameId
    ) external view returns (GameVaultInfo memory);

    function getPlayerStakeInfo(
        uint256 gameId,
        address player
    ) external view returns (PlayerStakeInfo memory);

    function getPlayerStake(
        uint256 gameId,
        address player
    ) external view returns (uint256);

    function getTotalGameStakes(uint256 gameId) external view returns (uint256);

    function isGameEnded(uint256 gameId) external view returns (bool);
}
