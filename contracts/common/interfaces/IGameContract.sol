// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

interface IGameContract {
    // Enums
    enum GameState {
        Active,
        Ended
    }

    enum GameResult {
        Ongoing,
        WhiteWins,
        BlackWins,
        Draw
    }

    enum Team {
        White,
        Black
    }

    // Structs
    struct GameInfo {
        uint256 gameId;
        GameState state;
        GameResult result;
        uint256 fixedStakeAmount;
        uint256 createdAt;
        uint256 endedAt;
        uint256 totalWhiteStakes;
        uint256 totalBlackStakes;
        uint256 whitePlayerCount;
        uint256 blackPlayerCount;
    }

    struct PlayerInfo {
        Team team;
        uint256 totalStakes;
        uint256 moveCount;
        bool hasJoined;
    }

    // Events
    event GameCreated(
        uint256 indexed gameId,
        uint256 fixedStakeAmount,
        uint256 createdAt
    );

    event PlayerJoinedTeam(
        uint256 indexed gameId,
        address indexed player,
        Team team
    );

    event MoveRecorded(
        uint256 indexed gameId,
        address indexed player,
        uint32 chainId,
        uint256 newMoveCount
    );

    event GameEnded(uint256 indexed gameId, GameResult result, uint256 endedAt);

    // Core Functions
    function createGame(uint256 gameId, uint256 fixedStakeAmount) external;

    function joinTeam(uint256 gameId, Team team) external;

    function recordMove(
        uint256 gameId,
        address player,
        uint32 chainId
    ) external;

    function endGame(uint256 gameId, GameResult result) external;

    // View Functions
    function getGameInfo(
        uint256 gameId
    ) external view returns (GameInfo memory);

    function getPlayerInfo(
        uint256 gameId,
        address player
    ) external view returns (PlayerInfo memory);

    function getPlayerMoveCount(
        uint256 gameId,
        address player,
        uint32 chainId
    ) external view returns (uint256);

    function calculateRewards(
        uint256 gameId,
        address player,
        uint256 playerTotalStakes
    ) external view returns (uint256);

    function isGameActive(uint256 gameId) external view returns (bool);

    function getGameResult(uint256 gameId) external view returns (GameResult);
}
