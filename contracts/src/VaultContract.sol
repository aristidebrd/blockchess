// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import "../common/interfaces/IVaultContract.sol";

contract VaultContract is IVaultContract {
    address public immutable authorizedBackend;
    address public immutable gameContract;

    mapping(uint256 gameId => GameVaultInfo) public gameVaults;
    mapping(uint256 gameId => mapping(address player => PlayerStakeInfo))
        public playerStakes;
    mapping(uint256 gameId => address[]) public gamePlayers;

    modifier onlyAuthorizedBackend() {
        require(
            msg.sender == authorizedBackend,
            "Only authorized backend can end games"
        );
        _;
    }

    modifier gameNotEnded(uint256 gameId) {
        require(!gameVaults[gameId].gameEnded, "Game has ended, cannot stake");
        _;
    }

    modifier gameEnded(uint256 gameId) {
        require(gameVaults[gameId].gameEnded, "Game has not ended yet");
        _;
    }

    constructor(address _authorizedBackend, address _gameContract) {
        require(
            _authorizedBackend != address(0),
            "Authorized backend cannot be zero address"
        );
        require(
            _gameContract != address(0),
            "Game contract cannot be zero address"
        );
        authorizedBackend = _authorizedBackend;
        gameContract = _gameContract;
    }

    function stake(
        uint256 gameId,
        uint256 fixedStakeAmount
    ) external payable gameNotEnded(gameId) {
        require(
            fixedStakeAmount > 0,
            "Fixed stake amount must be greater than 0"
        );
        require(
            msg.value == fixedStakeAmount,
            "Sent value must match fixed stake amount"
        );

        if (gameVaults[gameId].gameId == 0) {
            gameVaults[gameId] = GameVaultInfo({
                gameId: gameId,
                totalStakes: 0,
                playerCount: 0,
                result: GameResult.Ongoing,
                gameEnded: false,
                endedAt: 0
            });
        }

        PlayerStakeInfo storage playerInfo = playerStakes[gameId][msg.sender];
        if (playerInfo.stakeCount == 0) {
            gamePlayers[gameId].push(msg.sender);
            gameVaults[gameId].playerCount++;
        }

        playerInfo.totalStaked += fixedStakeAmount;
        playerInfo.stakeCount++;

        gameVaults[gameId].totalStakes += fixedStakeAmount;

        emit StakeDeposited(
            gameId,
            msg.sender,
            fixedStakeAmount,
            playerInfo.totalStaked
        );
    }

    function endGame(
        uint256 gameId,
        GameResult result
    ) external onlyAuthorizedBackend {
        require(
            result != GameResult.Ongoing,
            "Cannot end game with ongoing result"
        );
        require(!gameVaults[gameId].gameEnded, "Game already ended");

        if (gameVaults[gameId].gameId == 0) {
            gameVaults[gameId] = GameVaultInfo({
                gameId: gameId,
                totalStakes: 0,
                playerCount: 0,
                result: result,
                gameEnded: true,
                endedAt: block.timestamp
            });
        } else {
            gameVaults[gameId].result = result;
            gameVaults[gameId].gameEnded = true;
            gameVaults[gameId].endedAt = block.timestamp;
        }

        emit GameEndedInVault(
            gameId,
            result,
            gameVaults[gameId].totalStakes,
            block.timestamp
        );
    }

    function claimRewards(uint256 gameId) external gameEnded(gameId) {
        PlayerStakeInfo storage playerInfo = playerStakes[gameId][msg.sender];
        require(playerInfo.totalStaked > 0, "No stake in this game");
        require(!playerInfo.hasClaimed, "Rewards already claimed");

        uint256 rewardAmount = _calculatePlayerRewards(
            gameId,
            msg.sender,
            playerInfo.totalStaked
        );

        playerInfo.hasClaimed = true;

        if (rewardAmount > 0) {
            (bool success, ) = payable(msg.sender).call{value: rewardAmount}(
                ""
            );
            require(success, "Failed to transfer rewards");
        }

        emit RewardsClaimed(gameId, msg.sender, rewardAmount);
    }

    function getGameVaultInfo(
        uint256 gameId
    ) external view returns (GameVaultInfo memory) {
        return gameVaults[gameId];
    }

    function getPlayerStakeInfo(
        uint256 gameId,
        address player
    ) external view returns (PlayerStakeInfo memory) {
        return playerStakes[gameId][player];
    }

    function getPlayerStake(
        uint256 gameId,
        address player
    ) external view returns (uint256) {
        return playerStakes[gameId][player].totalStaked;
    }

    function getTotalGameStakes(
        uint256 gameId
    ) external view returns (uint256) {
        return gameVaults[gameId].totalStakes;
    }

    function isGameEnded(uint256 gameId) external view returns (bool) {
        return gameVaults[gameId].gameEnded;
    }

    function canClaimRewards(
        uint256 gameId,
        address player
    ) external view returns (bool) {
        PlayerStakeInfo memory playerInfo = playerStakes[gameId][player];
        return
            gameVaults[gameId].gameEnded &&
            playerInfo.totalStaked > 0 &&
            !playerInfo.hasClaimed;
    }

    function _calculatePlayerRewards(
        uint256 gameId,
        address /* player */,
        uint256 playerTotalStakes
    ) internal view returns (uint256) {
        GameVaultInfo memory vaultInfo = gameVaults[gameId];

        if (vaultInfo.result == GameResult.Draw) {
            return playerTotalStakes;
        }

        // For wins/losses, return stake back for now (simplified)
        // In production, this would query the game contract for team info
        return playerTotalStakes;
    }
}
