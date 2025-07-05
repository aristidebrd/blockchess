// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import "../common/interfaces/IGameContract.sol";

/**
 * @title GameContract
 * @dev Game contract deployed only on Base chain that tracks game state, team assignments, and move counts
 * @dev This contract does NOT handle stakes - that's handled by Vault contracts on each chain
 */
contract GameContract is IGameContract {
    // State variables
    address public immutable authorizedBackend;

    // Game storage
    mapping(uint256 gameId => GameInfo) public games;
    mapping(uint256 gameId => mapping(address player => PlayerInfo))
        public players;
    mapping(uint256 gameId => mapping(uint32 chainId => mapping(address player => uint256 moveCount)))
        public moves;

    // Modifiers
    modifier onlyAuthorizedBackend() {
        require(
            msg.sender == authorizedBackend,
            "Only authorized backend can perform this action"
        );
        _;
    }

    modifier gameExists(uint256 gameId) {
        require(games[gameId].gameId != 0, "Game does not exist");
        _;
    }

    modifier gameActive(uint256 gameId) {
        require(games[gameId].state == GameState.Active, "Game is not active");
        _;
    }

    constructor(address _authorizedBackend) {
        require(
            _authorizedBackend != address(0),
            "Authorized backend cannot be zero address"
        );
        authorizedBackend = _authorizedBackend;
    }

    // ============ CORE FUNCTIONS ============

    function createGame(
        uint256 gameId,
        uint256 fixedStakeAmount
    ) external onlyAuthorizedBackend {
        require(games[gameId].gameId == 0, "Game already exists");
        require(fixedStakeAmount > 0, "Fixed stake must be greater than 0");

        games[gameId] = GameInfo({
            gameId: gameId,
            state: GameState.Active,
            result: GameResult.Ongoing,
            fixedStakeAmount: fixedStakeAmount,
            createdAt: block.timestamp,
            endedAt: 0,
            totalWhiteStakes: 0,
            totalBlackStakes: 0,
            whitePlayerCount: 0,
            blackPlayerCount: 0
        });

        emit GameCreated(gameId, fixedStakeAmount, block.timestamp);
    }

    function joinTeam(
        uint256 gameId,
        Team team
    ) external gameExists(gameId) gameActive(gameId) {
        require(
            !players[gameId][msg.sender].hasJoined,
            "Player already joined this game"
        );

        players[gameId][msg.sender] = PlayerInfo({
            team: team,
            totalStakes: 0, // Stakes are tracked in vault contracts
            moveCount: 0,
            hasJoined: true
        });

        if (team == Team.White) {
            games[gameId].whitePlayerCount++;
        } else {
            games[gameId].blackPlayerCount++;
        }

        emit PlayerJoinedTeam(gameId, msg.sender, team);
    }

    function recordMove(
        uint256 gameId,
        address player,
        uint32 chainId
    ) external onlyAuthorizedBackend gameExists(gameId) gameActive(gameId) {
        require(
            players[gameId][player].hasJoined,
            "Player has not joined this game"
        );

        moves[gameId][chainId][player]++;
        players[gameId][player].moveCount++;

        emit MoveRecorded(
            gameId,
            player,
            chainId,
            moves[gameId][chainId][player]
        );
    }

    function endGame(
        uint256 gameId,
        GameResult result
    ) external onlyAuthorizedBackend gameExists(gameId) gameActive(gameId) {
        require(
            result != GameResult.Ongoing,
            "Cannot end game with ongoing result"
        );

        games[gameId].state = GameState.Ended;
        games[gameId].result = result;
        games[gameId].endedAt = block.timestamp;

        emit GameEnded(gameId, result, block.timestamp);
    }

    // ============ VIEW FUNCTIONS ============

    function getGameInfo(
        uint256 gameId
    ) external view returns (GameInfo memory) {
        return games[gameId];
    }

    function getPlayerInfo(
        uint256 gameId,
        address player
    ) external view returns (PlayerInfo memory) {
        return players[gameId][player];
    }

    function getPlayerMoveCount(
        uint256 gameId,
        address player,
        uint32 chainId
    ) external view returns (uint256) {
        return moves[gameId][chainId][player];
    }

    function calculateRewards(
        uint256 gameId,
        address player,
        uint256 playerTotalStakes
    ) external view gameExists(gameId) returns (uint256) {
        GameInfo memory game = games[gameId];
        require(game.state == GameState.Ended, "Game has not ended");

        PlayerInfo memory playerInfo = players[gameId][player];
        require(playerInfo.hasJoined, "Player has not joined this game");

        if (playerTotalStakes == 0) {
            return 0;
        }

        // Handle different game results
        if (game.result == GameResult.Draw) {
            // In draw, everyone gets their stake back
            return playerTotalStakes;
        } else if (
            (game.result == GameResult.WhiteWins &&
                playerInfo.team == Team.White) ||
            (game.result == GameResult.BlackWins &&
                playerInfo.team == Team.Black)
        ) {
            // Winner: gets stake back + proportional share of total pot
            // For now, return stake back (full reward calculation needs total pot info from vault)
            return playerTotalStakes;
        } else {
            // Loser: gets nothing
            return 0;
        }
    }

    function isGameActive(uint256 gameId) external view returns (bool) {
        return games[gameId].state == GameState.Active;
    }

    function getGameResult(uint256 gameId) external view returns (GameResult) {
        return games[gameId].result;
    }

    // ============ HELPER FUNCTIONS ============

    function getPlayerTeam(
        uint256 gameId,
        address player
    ) external view returns (Team) {
        return players[gameId][player].team;
    }

    function hasPlayerJoined(
        uint256 gameId,
        address player
    ) external view returns (bool) {
        return players[gameId][player].hasJoined;
    }

    function getPlayerTotalMoveCount(
        uint256 gameId,
        address player
    ) external view returns (uint256) {
        return players[gameId][player].moveCount;
    }
}
