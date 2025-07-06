// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {IGameContract} from "../common/interfaces/IGameContract.sol";

/**
 * @title GameContract
 * @dev Individual game contract deployed by GameFactory
 * @dev Only the factory can modify game state
 */
contract GameContract is IGameContract {
    address public immutable FACTORY;

    uint256 private gameId;
    uint256 private fixedStakeAmount;
    GameState private gameState;
    GameResult private gameResult;

    // Player data
    mapping(address => UserInfo) private users;
    address[] private playerAddresses;

    modifier onlyFactory() {
        require(msg.sender == FACTORY, "Only factory can call this function");
        _;
    }

    constructor(address _factory) {
        FACTORY = _factory;
        gameState = GameState.Active;
        gameResult = GameResult.Draw;
    }

    // ============ FACTORY-ONLY FUNCTIONS ============

    function initialize(
        uint256 _gameId,
        uint256 _fixedStakeAmount
    ) external override onlyFactory {
        gameId = _gameId;
        fixedStakeAmount = _fixedStakeAmount;
        emit GameCreated(_gameId, _fixedStakeAmount, block.timestamp);
    }

    function addVote(
        address player,
        uint32 chainId,
        Team team
    ) external override onlyFactory {
        if (users[player].totalVotes == 0) {
            // Add new user
            users[player] = UserInfo({
                chainId: chainId,
                totalVotes: 1,
                team: team
            });
            playerAddresses.push(player);
        } else {
            users[player].totalVotes++;
        }

        emit Vote(gameId, player, team, chainId);
    }

    function endGame(GameResult result) external override onlyFactory {
        require(gameState == GameState.Active, "Game is not active");

        gameState = GameState.Finished;
        gameResult = result;

        emit GameFinished(gameId, result, block.timestamp);
    }

    // ============ VIEW FUNCTIONS ============

    function getGameInfo() external view override returns (GameInfo memory) {
        return GameInfo({gameId: gameId, state: gameState, result: gameResult});
    }

    function getGameStatus() external view override returns (GameState) {
        return gameState;
    }

    function getGameResult() external view override returns (GameResult) {
        return gameResult;
    }

    function getPlayerVoteCounts()
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
        uint256 length = playerAddresses.length;
        players = new address[](length);
        voteCounts = new uint256[](length);
        chainIds = new uint32[](length);
        teams = new uint256[](length);

        for (uint256 i = 0; i < length; i++) {
            address player = playerAddresses[i];
            players[i] = player;
            voteCounts[i] = users[player].totalVotes;
            chainIds[i] = users[player].chainId;
            teams[i] = uint256(users[player].team);
        }

        return (players, voteCounts, chainIds, teams);
    }

    function getGameId() external view override returns (uint256) {
        return gameId;
    }

    function userExists(
        address userAddress
    ) external view override returns (bool) {
        return users[userAddress].totalVotes > 0;
    }
}
