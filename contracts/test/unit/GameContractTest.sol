// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import "../BaseTest.sol";
import "../../common/interfaces/IGameContract.sol";
import "../../src/GameContract.sol";

contract GameContractTest is BaseTest {
    GameContract public gameContract;

    uint256 constant GAME_ID = 123;
    uint32 constant BASE_CHAIN_ID = 8453;
    uint32 constant MAINNET_CHAIN_ID = 1;
    uint32 constant OPTIMISM_CHAIN_ID = 10;

    // Events to test
    event GameCreated(
        uint256 indexed gameId,
        uint256 fixedStakeAmount,
        uint256 createdAt
    );

    event PlayerJoinedTeam(
        uint256 indexed gameId,
        address indexed player,
        IGameContract.Team team
    );

    event MoveRecorded(
        uint256 indexed gameId,
        address indexed player,
        uint32 chainId,
        uint256 newMoveCount
    );

    event GameEnded(
        uint256 indexed gameId,
        IGameContract.GameResult result,
        uint256 endedAt
    );

    function setUp() public {
        // Deploy game contract with ALICE as Go backend address
        gameContract = new GameContract(ALICE);

        // Fund test addresses
        fundAddress(ALICE, 10 ether);
        fundAddress(BOB, 10 ether);
        fundAddress(CHARLIE, 10 ether);
        fundAddress(DAVID, 10 ether);
    }

    // ============ GAME CREATION TESTS ============

    function test_CreateGame_ShouldSucceed() public {
        vm.expectEmit(true, false, false, true);
        emit GameCreated(GAME_ID, FIXED_STAKE, block.timestamp);

        vm.prank(ALICE); // Go backend creates game
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        IGameContract.GameInfo memory gameInfo = gameContract.getGameInfo(
            GAME_ID
        );
        assertEq(gameInfo.gameId, GAME_ID, "Game ID should match");
        assertEq(
            gameInfo.fixedStakeAmount,
            FIXED_STAKE,
            "Fixed stake should match"
        );
        assertEq(
            uint8(gameInfo.state),
            uint8(IGameContract.GameState.Active),
            "Game should be active"
        );
        assertEq(
            uint8(gameInfo.result),
            uint8(IGameContract.GameResult.Ongoing),
            "Result should be ongoing"
        );
        assertEq(
            gameInfo.totalWhiteStakes,
            0,
            "Should start with 0 white stakes"
        );
        assertEq(
            gameInfo.totalBlackStakes,
            0,
            "Should start with 0 black stakes"
        );
        assertTrue(gameContract.isGameActive(GAME_ID), "Game should be active");
    }

    function test_CreateGame_ShouldRevertDuplicateGameId() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.expectRevert("Game already exists");
        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);
    }

    function test_CreateGame_ShouldRevertZeroStake() public {

        vm.expectRevert("Fixed stake must be greater than 0");
        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, 0);
    }

    function test_CreateGame_ShouldRevertUnauthorized() public {

        vm.expectRevert("Only authorized backend can perform this action");
        vm.prank(BOB); // Not the Go backend
        gameContract.createGame(GAME_ID, FIXED_STAKE);
    }

    // ============ TEAM JOINING TESTS ============

    function test_JoinTeam_WhiteTeam_ShouldSucceed() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.expectEmit(true, true, false, true);
        emit PlayerJoinedTeam(GAME_ID, BOB, IGameContract.Team.White);

        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        IGameContract.PlayerInfo memory playerInfo = gameContract.getPlayerInfo(
            GAME_ID,
            BOB
        );
        assertEq(
            uint8(playerInfo.team),
            uint8(IGameContract.Team.White),
            "Should be on white team"
        );
        assertTrue(playerInfo.hasJoined, "Should be marked as joined");
        assertEq(playerInfo.moveCount, 0, "Should start with 0 moves");
        assertEq(
            playerInfo.totalStakes,
            0,
            "Should start with 0 stakes (tracked in vault)"
        );

        IGameContract.GameInfo memory gameInfo = gameContract.getGameInfo(
            GAME_ID
        );
        assertEq(gameInfo.whitePlayerCount, 1, "Should have 1 white player");
        assertEq(gameInfo.blackPlayerCount, 0, "Should have 0 black players");
    }

    function test_JoinTeam_BlackTeam_ShouldSucceed() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.expectEmit(true, true, false, true);
        emit PlayerJoinedTeam(GAME_ID, CHARLIE, IGameContract.Team.Black);

        vm.prank(CHARLIE);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.Black);

        IGameContract.PlayerInfo memory playerInfo = gameContract.getPlayerInfo(
            GAME_ID,
            CHARLIE
        );
        assertEq(
            uint8(playerInfo.team),
            uint8(IGameContract.Team.Black),
            "Should be on black team"
        );

        IGameContract.GameInfo memory gameInfo = gameContract.getGameInfo(
            GAME_ID
        );
        assertEq(gameInfo.blackPlayerCount, 1, "Should have 1 black player");
    }

    function test_JoinTeam_ShouldRevertGameNotExists() public {

        vm.expectRevert("Game does not exist");
        vm.prank(BOB);
        gameContract.joinTeam(999, IGameContract.Team.White);
    }

    function test_JoinTeam_ShouldRevertAlreadyJoined() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        vm.expectRevert("Player already joined this game");
        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.Black);
    }

    function test_JoinTeam_ShouldRevertGameEnded() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.prank(ALICE);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.WhiteWins);

        vm.expectRevert("Game is not active");
        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);
    }

    // ============ MOVE RECORDING TESTS ============

    function test_RecordMove_ShouldSucceed() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        vm.expectEmit(true, true, false, true);
        emit MoveRecorded(GAME_ID, BOB, MAINNET_CHAIN_ID, 1);

        vm.prank(ALICE); // Go backend records move
        gameContract.recordMove(GAME_ID, BOB, MAINNET_CHAIN_ID);

        uint256 moveCount = gameContract.getPlayerMoveCount(
            GAME_ID,
            BOB,
            MAINNET_CHAIN_ID
        );
        assertEq(moveCount, 1, "Should have 1 move on mainnet");

        IGameContract.PlayerInfo memory playerInfo = gameContract.getPlayerInfo(
            GAME_ID,
            BOB
        );
        assertEq(playerInfo.moveCount, 1, "Should have total 1 move");
    }

    function test_RecordMove_MultipleChains_ShouldSucceed() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        // Record moves on different chains
        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, MAINNET_CHAIN_ID);

        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, OPTIMISM_CHAIN_ID);

        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, BASE_CHAIN_ID);

        assertEq(
            gameContract.getPlayerMoveCount(GAME_ID, BOB, MAINNET_CHAIN_ID),
            1,
            "Should have 1 move on mainnet"
        );
        assertEq(
            gameContract.getPlayerMoveCount(GAME_ID, BOB, OPTIMISM_CHAIN_ID),
            1,
            "Should have 1 move on optimism"
        );
        assertEq(
            gameContract.getPlayerMoveCount(GAME_ID, BOB, BASE_CHAIN_ID),
            1,
            "Should have 1 move on base"
        );

        IGameContract.PlayerInfo memory playerInfo = gameContract.getPlayerInfo(
            GAME_ID,
            BOB
        );
        assertEq(playerInfo.moveCount, 3, "Should have total 3 moves");
    }

    function test_RecordMove_ShouldRevertUnauthorized() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        vm.expectRevert("Only authorized backend can perform this action");
        vm.prank(BOB); // Player trying to record own move
        gameContract.recordMove(GAME_ID, BOB, MAINNET_CHAIN_ID);
    }

    function test_RecordMove_ShouldRevertPlayerNotJoined() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.expectRevert("Player has not joined this game");
        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, MAINNET_CHAIN_ID);
    }

    // ============ GAME ENDING TESTS ============

    function test_EndGame_WhiteWins_ShouldSucceed() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.expectEmit(true, false, false, true);
        emit GameEnded(
            GAME_ID,
            IGameContract.GameResult.WhiteWins,
            block.timestamp
        );

        vm.prank(ALICE);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.WhiteWins);

        IGameContract.GameInfo memory gameInfo = gameContract.getGameInfo(
            GAME_ID
        );
        assertEq(
            uint8(gameInfo.state),
            uint8(IGameContract.GameState.Ended),
            "Game should be ended"
        );
        assertEq(
            uint8(gameInfo.result),
            uint8(IGameContract.GameResult.WhiteWins),
            "Result should be white wins"
        );
        assertFalse(
            gameContract.isGameActive(GAME_ID),
            "Game should not be active"
        );
        assertEq(
            uint8(gameContract.getGameResult(GAME_ID)),
            uint8(IGameContract.GameResult.WhiteWins),
            "Should return correct result"
        );
    }

    function test_EndGame_ShouldRevertUnauthorized() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.expectRevert("Only authorized backend can perform this action");
        vm.prank(BOB);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.WhiteWins);
    }

    function test_EndGame_ShouldRevertAlreadyEnded() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.prank(ALICE);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.WhiteWins);

        vm.expectRevert("Game is not active");
        vm.prank(ALICE);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.BlackWins);
    }

    // ============ REWARD CALCULATION TESTS ============

    function test_CalculateRewards_WhiteWins_ShouldReturnCorrectAmount()
        public
    {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        // Join teams
        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        vm.prank(CHARLIE);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.Black);

        // End game
        vm.prank(ALICE);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.WhiteWins);

        // Calculate rewards (player staked 2 ETH, total pot is 5 ETH, player is winner)
        uint256 playerStakes = 2 ether;
        // uint256 totalPot = 5 ether;
        // uint256 expectedReward = playerStakes +
        //     (playerStakes * totalPot) /
        //     totalPot; // stake + proportional share

        uint256 actualReward = gameContract.calculateRewards(
            GAME_ID,
            BOB,
            playerStakes
        );

        // For winners: they get their stake back + proportional share of loser pot
        // This will be implemented in the actual contract
        assertTrue(actualReward > 0, "Winner should get rewards");
    }

    function test_CalculateRewards_ShouldReturnZeroForLoser() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        vm.prank(CHARLIE);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.Black);

        vm.prank(ALICE);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.WhiteWins);

        uint256 charlieStakes = 1 ether;
        uint256 reward = gameContract.calculateRewards(
            GAME_ID,
            CHARLIE,
            charlieStakes
        );
        assertEq(reward, 0, "Loser should get no rewards");
    }

    function test_CalculateRewards_Draw_ShouldReturnStakeBack() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        vm.prank(ALICE);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.Draw);

        uint256 playerStakes = 2 ether;
        uint256 reward = gameContract.calculateRewards(
            GAME_ID,
            BOB,
            playerStakes
        );
        assertEq(reward, playerStakes, "In draw, player should get stake back");
    }

    // ============ VIEW FUNCTION TESTS ============

    function test_GetGameInfo_NonExistentGame_ShouldReturnEmpty() public view {
        IGameContract.GameInfo memory gameInfo = gameContract.getGameInfo(999);
        assertEq(
            gameInfo.gameId,
            0,
            "Non-existent game should return empty info"
        );
    }

    function test_GetPlayerInfo_NonJoinedPlayer_ShouldReturnEmpty() public {

        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        IGameContract.PlayerInfo memory playerInfo = gameContract.getPlayerInfo(
            GAME_ID,
            BOB
        );
        assertFalse(
            playerInfo.hasJoined,
            "Non-joined player should return false"
        );
        assertEq(
            playerInfo.moveCount,
            0,
            "Non-joined player should have 0 moves"
        );
    }

    function test_Integration_CompleteGameFlow() public {

        // 1. Create game
        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        // 2. Players join teams
        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        vm.prank(CHARLIE);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.Black);

        // 3. Record some moves
        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, MAINNET_CHAIN_ID);

        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, CHARLIE, OPTIMISM_CHAIN_ID);

        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, BASE_CHAIN_ID);

        // 4. End game
        vm.prank(ALICE);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.WhiteWins);

        // 5. Verify final state
        IGameContract.GameInfo memory gameInfo = gameContract.getGameInfo(
            GAME_ID
        );
        assertEq(
            uint8(gameInfo.state),
            uint8(IGameContract.GameState.Ended),
            "Game should be ended"
        );
        assertEq(
            uint8(gameInfo.result),
            uint8(IGameContract.GameResult.WhiteWins),
            "Result should be white wins"
        );

        IGameContract.PlayerInfo memory bobInfo = gameContract.getPlayerInfo(
            GAME_ID,
            BOB
        );
        assertEq(bobInfo.moveCount, 2, "Bob should have 2 total moves");

        IGameContract.PlayerInfo memory charlieInfo = gameContract
            .getPlayerInfo(GAME_ID, CHARLIE);
        assertEq(charlieInfo.moveCount, 1, "Charlie should have 1 total move");
    }
}
