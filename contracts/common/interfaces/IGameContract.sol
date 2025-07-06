// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

interface IGameContract {
    // Enums
    enum GameState {
        Active,
        Finished
    }

    enum GameResult {
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
    }

    struct UserInfo {
        uint32 chainId;
        uint256 totalVotes;
        Team team;
    }

    // Events
    event GameCreated(
        uint256 indexed gameId,
        uint256 fixedStakeAmount,
        uint256 createdAt
    );

    event Vote(
        uint256 indexed gameId,
        address indexed player,
        Team team,
        uint32 chainId
    );

    event GameFinished(
        uint256 indexed gameId,
        GameResult result,
        uint256 finishedAt
    );

    // Core Functions
    function initialize(uint256 gameId, uint256 fixedStakeAmount) external;

    function addVote(address player, uint32 chainId, Team team) external;

    function endGame(GameResult result) external;

    // View Functions
    function getGameInfo() external view returns (GameInfo memory);

    function getGameStatus() external view returns (GameState);

    function getGameResult() external view returns (GameResult);

    function getPlayerVoteCounts()
        external
        view
        returns (
            address[] memory players,
            uint256[] memory voteCounts,
            uint32[] memory chainIds,
            uint256[] memory teams
        );

    function getGameId() external view returns (uint256);

    function userExists(address userAddress) external view returns (bool);
}
