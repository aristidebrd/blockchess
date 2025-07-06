// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {IGameContract} from "./IGameContract.sol";

interface IGameFactory {
    // Events
    event GameCreated(
        uint256 indexed gameId,
        address indexed gameContract,
        uint256 fixedStakeAmount,
        uint256 createdAt
    );

    // Structs
    struct Game {
        uint256 gameId;
        address gameContract;
        uint256 fixedStakeAmount;
        uint256 createdAt;
    }

    struct Games {
        Game[] games;
    }

    // Core Functions
    function createGame(
        uint256 fixedStakeAmount
    ) external returns (uint256 gameId, address gameContract);

    function endGame(uint256 gameId, IGameContract.GameResult result) external;

    function addVote(
        uint256 gameId,
        address player,
        uint32 chainId,
        IGameContract.Team team
    ) external;

    // View Functions
    function getGameContractAddress(
        uint256 gameId
    ) external view returns (address);

    function getGame(uint256 gameId) external view returns (Game memory);

    function getAllGames() external view returns (Games memory);

    function getActiveGames() external view returns (Games memory);

    function gameExists(uint256 gameId) external view returns (bool);

    // Game Contract Proxy Functions
    function getGameResult(
        uint256 gameId
    ) external view returns (IGameContract.GameResult);

    function getGameStatus(uint256 gameId) external view returns (bool);

    function getPlayerVoteCounts(
        uint256 gameId
    )
        external
        view
        returns (
            address[] memory players,
            uint256[] memory voteCounts,
            uint32[] memory chainIds,
            uint256[] memory teams
        );
}
