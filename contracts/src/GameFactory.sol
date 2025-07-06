// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {IGameFactory} from "../common/interfaces/IGameFactory.sol";
import {IGameContract} from "../common/interfaces/IGameContract.sol";
import {GameContract} from "./GameContract.sol";

/**
 * @title GameFactory
 * @dev Factory contract deployed on Base testnet that creates and manages GameContract instances
 * @dev This contract acts as the central registry for all games and provides proxy functions
 */
contract GameFactory is IGameFactory {
    // State variables
    address public immutable AUTHORIZED_BACKEND;
    uint256 private nextGameId;

    // Game storage
    mapping(uint256 => Game) public games;
    mapping(uint256 => address) public gameContracts;
    uint256[] public gameIds;

    // Modifiers
    modifier onlyAuthorizedBackend() {
        require(
            msg.sender == AUTHORIZED_BACKEND,
            "Only authorized backend can perform this action"
        );
        _;
    }

    modifier gameExistsModifier(uint256 gameId) {
        require(games[gameId].gameId != 0, "Game does not exist");
        _;
    }

    constructor(address _authorizedBackend) {
        require(
            _authorizedBackend != address(0),
            "Authorized backend cannot be zero address"
        );
        AUTHORIZED_BACKEND = _authorizedBackend;
        nextGameId = 1; // Start game IDs from 1
    }

    // ============ CORE FUNCTIONS ============

    function createGame(
        uint256 fixedStakeAmount
    )
        external
        override
        onlyAuthorizedBackend
        returns (uint256 gameId, address gameContract)
    {
        require(
            fixedStakeAmount > 0,
            "Fixed stake amount must be greater than 0"
        );

        gameId = nextGameId;
        nextGameId++;

        // Deploy new GameContract instance
        GameContract newGameContract = new GameContract(address(this));
        gameContract = address(newGameContract);

        // Initialize the game in the new contract
        newGameContract.initialize(gameId, fixedStakeAmount);

        // Store game information
        games[gameId] = Game({
            gameId: gameId,
            gameContract: gameContract,
            fixedStakeAmount: fixedStakeAmount,
            createdAt: block.timestamp
        });

        gameContracts[gameId] = gameContract;
        gameIds.push(gameId);

        emit GameCreated(
            gameId,
            gameContract,
            fixedStakeAmount,
            block.timestamp
        );

        return (gameId, gameContract);
    }

    function endGame(
        uint256 gameId,
        IGameContract.GameResult result
    ) external override onlyAuthorizedBackend gameExistsModifier(gameId) {
        address gameContract = gameContracts[gameId];
        require(gameContract != address(0), "Game contract not found");

        // Call endGame on the specific game contract
        IGameContract(gameContract).endGame(result);
    }

    function addVote(
        uint256 gameId,
        address player,
        uint32 chainId,
        IGameContract.Team team
    ) external override onlyAuthorizedBackend gameExistsModifier(gameId) {
        address gameContract = gameContracts[gameId];
        require(gameContract != address(0), "Game contract not found");

        // Call addVote on the specific game contract
        IGameContract(gameContract).addVote(player, chainId, team);
    }

    // ============ VIEW FUNCTIONS ============

    function getGameContractAddress(
        uint256 gameId
    ) external view override returns (address) {
        return gameContracts[gameId];
    }

    function getGame(
        uint256 gameId
    ) external view override returns (Game memory) {
        return games[gameId];
    }

    function getAllGames() external view override returns (Games memory) {
        Game[] memory allGames = new Game[](gameIds.length);

        for (uint256 i = 0; i < gameIds.length; i++) {
            allGames[i] = games[gameIds[i]];
        }

        return Games({games: allGames});
    }

    function getActiveGames() external view override returns (Games memory) {
        // First, count active games
        uint256 activeCount = 0;
        for (uint256 i = 0; i < gameIds.length; i++) {
            uint256 gameId = gameIds[i];
            address gameContract = gameContracts[gameId];
            if (gameContract != address(0)) {
                IGameContract.GameState state = IGameContract(gameContract)
                    .getGameStatus();
                if (state == IGameContract.GameState.Active) {
                    activeCount++;
                }
            }
        }

        // Create array with active games
        Game[] memory activeGames = new Game[](activeCount);
        uint256 index = 0;

        for (uint256 i = 0; i < gameIds.length; i++) {
            uint256 gameId = gameIds[i];
            address gameContract = gameContracts[gameId];
            if (gameContract != address(0)) {
                IGameContract.GameState state = IGameContract(gameContract)
                    .getGameStatus();
                if (state == IGameContract.GameState.Active) {
                    activeGames[index] = games[gameId];
                    index++;
                }
            }
        }

        return Games({games: activeGames});
    }

    function gameExists(uint256 gameId) external view override returns (bool) {
        return games[gameId].gameId != 0;
    }

    // ============ GAME CONTRACT PROXY FUNCTIONS ============

    function getGameResult(
        uint256 gameId
    ) external view override returns (IGameContract.GameResult) {
        address gameContract = gameContracts[gameId];
        require(gameContract != address(0), "Game contract not found");

        return IGameContract(gameContract).getGameResult();
    }

    function getGameStatus(
        uint256 gameId
    ) external view override returns (bool) {
        address gameContract = gameContracts[gameId];
        require(gameContract != address(0), "Game contract not found");

        IGameContract.GameState state = IGameContract(gameContract)
            .getGameStatus();
        return state == IGameContract.GameState.Active;
    }

    function getPlayerVoteCounts(
        uint256 gameId
    )
        external
        view
        override
        returns (
            address[] memory players,
            uint256[] memory voteCounts,
            uint32[] memory chainIds,
            uint256[] memory teams
        )
    {
        address gameContract = gameContracts[gameId];
        require(gameContract != address(0), "Game contract not found");

        (players, voteCounts, chainIds, teams) = IGameContract(gameContract)
            .getPlayerVoteCounts();
        return (players, voteCounts, chainIds, teams);
    }
}
