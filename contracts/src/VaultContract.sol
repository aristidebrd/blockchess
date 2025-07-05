// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import "../common/interfaces/IVaultContract.sol";

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
}

contract VaultContract is IVaultContract {
    address public immutable authorizedBackend;
    address public immutable gameContract;
    address public immutable usdcToken; // USDC contract address

    mapping(uint256 gameId => GameVaultInfo) public gameVaults;
    mapping(uint256 gameId => mapping(address player => PlayerStakeInfo))
        public playerStakes;
    mapping(uint256 gameId => address[]) public gamePlayers;

    // New: Track player approvals for games
    mapping(address player => mapping(uint256 gameId => bool))
        public playerApprovals;

    event PlayerApproved(
        uint256 indexed gameId,
        address indexed player,
        uint256 approvedAmount
    );
    event StakeDeductedFromPlayer(
        uint256 indexed gameId,
        address indexed player,
        uint256 amount
    );

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

    constructor(
        address _authorizedBackend,
        address _gameContract,
        address _usdcToken
    ) {
        require(
            _authorizedBackend != address(0),
            "Authorized backend cannot be zero address"
        );
        require(
            _gameContract != address(0),
            "Game contract cannot be zero address"
        );
        require(_usdcToken != address(0), "USDC token cannot be zero address");

        authorizedBackend = _authorizedBackend;
        gameContract = _gameContract;
        usdcToken = _usdcToken;
    }

    // New: Player approves participation in a game
    function approveGameParticipation(
        uint256 gameId,
        uint256 maxStakeAmount
    ) external {
        require(maxStakeAmount > 0, "Max stake amount must be greater than 0");
        require(!gameVaults[gameId].gameEnded, "Game has ended");

        // Check if player has enough USDC balance
        require(
            IERC20(usdcToken).balanceOf(msg.sender) >= maxStakeAmount,
            "Insufficient USDC balance"
        );

        // Check if player has approved enough USDC to this contract
        require(
            IERC20(usdcToken).allowance(msg.sender, address(this)) >=
                maxStakeAmount,
            "Please approve USDC spending first"
        );

        playerApprovals[msg.sender][gameId] = true;

        emit PlayerApproved(gameId, msg.sender, maxStakeAmount);
    }

    // Modified: Backend can stake on behalf of approved players
    function stakeOnBehalfOf(
        uint256 gameId,
        address player,
        uint256 stakeAmount
    ) external onlyAuthorizedBackend gameNotEnded(gameId) {
        require(stakeAmount > 0, "Stake amount must be greater than 0");
        require(
            playerApprovals[player][gameId],
            "Player has not approved participation"
        );

        // Check allowance
        require(
            IERC20(usdcToken).allowance(player, address(this)) >= stakeAmount,
            "Insufficient USDC allowance"
        );

        // Transfer USDC from player to this contract
        require(
            IERC20(usdcToken).transferFrom(player, address(this), stakeAmount),
            "USDC transfer failed"
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

        PlayerStakeInfo storage playerInfo = playerStakes[gameId][player];
        if (playerInfo.stakeCount == 0) {
            gamePlayers[gameId].push(player);
            gameVaults[gameId].playerCount++;
        }

        playerInfo.totalStaked += stakeAmount;
        playerInfo.stakeCount++;
        gameVaults[gameId].totalStakes += stakeAmount;

        emit StakeDeposited(
            gameId,
            player,
            stakeAmount,
            playerInfo.totalStaked
        );
        emit StakeDeductedFromPlayer(gameId, player, stakeAmount);
    }

    // Keep original stake function for direct player staking
    function stake(
        uint256 gameId,
        uint256 stakeAmount
    ) external gameNotEnded(gameId) {
        require(stakeAmount > 0, "Stake amount must be greater than 0");

        // Transfer USDC from player to this contract
        require(
            IERC20(usdcToken).transferFrom(
                msg.sender,
                address(this),
                stakeAmount
            ),
            "USDC transfer failed"
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

        playerInfo.totalStaked += stakeAmount;
        playerInfo.stakeCount++;
        gameVaults[gameId].totalStakes += stakeAmount;

        emit StakeDeposited(
            gameId,
            msg.sender,
            stakeAmount,
            playerInfo.totalStaked
        );
    }

    // New: Check if player has approved game participation
    function hasApprovedGame(
        address player,
        uint256 gameId
    ) external view returns (bool) {
        return playerApprovals[player][gameId];
    }

    // New: Get player's USDC allowance for this contract
    function getPlayerAllowance(
        address player
    ) external view returns (uint256) {
        return IERC20(usdcToken).allowance(player, address(this));
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
}
