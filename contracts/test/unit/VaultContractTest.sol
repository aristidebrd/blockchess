// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import "../BaseTest.sol";
import "../../common/interfaces/IVaultContract.sol";
import "../../src/VaultContract.sol";

contract VaultContractTest is BaseTest {
    VaultContract public vaultContract;

    uint256 constant GAME_ID = 123;

    // Events to test
    event StakeDeposited(
        uint256 indexed gameId,
        address indexed player,
        uint256 amount,
        uint256 newTotal
    );

    event GameEndedInVault(
        uint256 indexed gameId,
        IVaultContract.GameResult result,
        uint256 totalStakes,
        uint256 endedAt
    );

    event RewardsClaimed(
        uint256 indexed gameId,
        address indexed player,
        uint256 amount
    );

    function setUp() public {
        // Deploy vault contract with ALICE as Go backend address
        vaultContract = new VaultContract(ALICE, address(0x1234)); // Mock game contract address

        // Fund test addresses
        fundAddress(ALICE, 10 ether);
        fundAddress(BOB, 10 ether);
        fundAddress(CHARLIE, 10 ether);
        fundAddress(DAVID, 10 ether);
    }

    // ============ STAKING TESTS ============

    function test_Stake_ShouldSucceed() public {
        vm.expectEmit(true, true, false, true);
        emit StakeDeposited(GAME_ID, BOB, FIXED_STAKE, FIXED_STAKE);

        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        IVaultContract.PlayerStakeInfo memory playerInfo = vaultContract
            .getPlayerStakeInfo(GAME_ID, BOB);
        assertEq(
            playerInfo.totalStaked,
            FIXED_STAKE,
            "Should have correct stake amount"
        );
        assertEq(playerInfo.stakeCount, 1, "Should have 1 stake");
        assertFalse(playerInfo.hasClaimed, "Should not have claimed yet");

        assertEq(
            vaultContract.getPlayerStake(GAME_ID, BOB),
            FIXED_STAKE,
            "Should return correct player stake"
        );
        assertEq(
            vaultContract.getTotalGameStakes(GAME_ID),
            FIXED_STAKE,
            "Should have correct total stakes"
        );
    }

    function test_Stake_MultipleStakes_ShouldAccumulate() public {
        // First stake
        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        // Second stake
        vm.expectEmit(true, true, false, true);
        emit StakeDeposited(GAME_ID, BOB, FIXED_STAKE, FIXED_STAKE * 2);

        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        IVaultContract.PlayerStakeInfo memory playerInfo = vaultContract
            .getPlayerStakeInfo(GAME_ID, BOB);
        assertEq(
            playerInfo.totalStaked,
            FIXED_STAKE * 2,
            "Should have accumulated stakes"
        );
        assertEq(playerInfo.stakeCount, 2, "Should have 2 stakes");

        assertEq(
            vaultContract.getPlayerStake(GAME_ID, BOB),
            FIXED_STAKE * 2,
            "Should return accumulated stake"
        );
    }

    function test_Stake_MultiplePlayersInSameGame_ShouldSucceed() public {
        // BOB stakes
        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        // CHARLIE stakes
        vm.prank(CHARLIE);
        vaultContract.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        assertEq(
            vaultContract.getPlayerStake(GAME_ID, BOB),
            FIXED_STAKE,
            "Bob should have correct stake"
        );
        assertEq(
            vaultContract.getPlayerStake(GAME_ID, CHARLIE),
            FIXED_STAKE,
            "Charlie should have correct stake"
        );
        assertEq(
            vaultContract.getTotalGameStakes(GAME_ID),
            FIXED_STAKE * 2,
            "Total should be sum of all stakes"
        );

        IVaultContract.GameVaultInfo memory gameInfo = vaultContract
            .getGameVaultInfo(GAME_ID);
        assertEq(
            gameInfo.totalStakes,
            FIXED_STAKE * 2,
            "Game info should show correct total"
        );
        assertEq(gameInfo.playerCount, 2, "Should have 2 players");
    }

    function test_Stake_ShouldRevertIncorrectAmount() public {
        vm.expectRevert("Sent value must match fixed stake amount");
        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE * 2}(GAME_ID, FIXED_STAKE); // Sent 2x but specified 1x
    }

    function test_Stake_ShouldRevertZeroAmount() public {
        vm.expectRevert("Fixed stake amount must be greater than 0");
        vm.prank(BOB);
        vaultContract.stake{value: 0}(GAME_ID, 0);
    }

    function test_Stake_ShouldRevertGameEnded() public {
        // End game first
        vm.prank(ALICE); // Go backend
        vaultContract.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);

        vm.expectRevert("Game has ended, cannot stake");
        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);
    }

    // ============ GAME ENDING TESTS ============

    function test_EndGame_WhiteWins_ShouldSucceed() public {
        // Add some stakes first
        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        vm.prank(CHARLIE);
        vaultContract.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        vm.expectEmit(true, false, false, true);
        emit GameEndedInVault(
            GAME_ID,
            IVaultContract.GameResult.WhiteWins,
            FIXED_STAKE * 2,
            block.timestamp
        );

        vm.prank(ALICE); // Go backend
        vaultContract.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);

        IVaultContract.GameVaultInfo memory gameInfo = vaultContract
            .getGameVaultInfo(GAME_ID);
        assertTrue(gameInfo.gameEnded, "Game should be marked as ended");
        assertEq(
            uint8(gameInfo.result),
            uint8(IVaultContract.GameResult.WhiteWins),
            "Result should be white wins"
        );
        assertEq(gameInfo.endedAt, block.timestamp, "End time should be set");

        assertTrue(
            vaultContract.isGameEnded(GAME_ID),
            "isGameEnded should return true"
        );
    }

    function test_EndGame_ShouldRevertUnauthorized() public {
        vm.expectRevert("Only authorized backend can end games");
        vm.prank(BOB); // Not the Go backend
        vaultContract.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);
    }

    function test_EndGame_ShouldRevertAlreadyEnded() public {
        vm.prank(ALICE);
        vaultContract.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);

        vm.expectRevert("Game already ended");
        vm.prank(ALICE);
        vaultContract.endGame(GAME_ID, IVaultContract.GameResult.BlackWins);
    }

    function test_EndGame_AllResults_ShouldWork() public {
        uint256[] memory gameIds = new uint256[](3);
        gameIds[0] = 1;
        gameIds[1] = 2;
        gameIds[2] = 3;

        IVaultContract.GameResult[]
            memory results = new IVaultContract.GameResult[](3);
        results[0] = IVaultContract.GameResult.WhiteWins;
        results[1] = IVaultContract.GameResult.BlackWins;
        results[2] = IVaultContract.GameResult.Draw;

        for (uint256 i = 0; i < gameIds.length; i++) {
            vm.prank(ALICE);
            vaultContract.endGame(gameIds[i], results[i]);

            IVaultContract.GameVaultInfo memory gameInfo = vaultContract
                .getGameVaultInfo(gameIds[i]);
            assertEq(
                uint8(gameInfo.result),
                uint8(results[i]),
                "Result should match"
            );
        }
    }

    // ============ REWARD CLAIMING TESTS ============

    function test_ClaimRewards_ShouldSucceed() public {
        uint256 stakeAmount = 1 ether;

        // BOB stakes
        vm.prank(BOB);
        vaultContract.stake{value: stakeAmount}(GAME_ID, stakeAmount);

        // End game (this will determine rewards based on game contract)
        vm.prank(ALICE);
        vaultContract.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);

        uint256 balanceBefore = BOB.balance;

        vm.expectEmit(true, true, false, true);
        emit RewardsClaimed(GAME_ID, BOB, stakeAmount); // Assuming BOB gets his stake back

        vm.prank(BOB);
        vaultContract.claimRewards(GAME_ID);

        uint256 balanceAfter = BOB.balance;
        assertTrue(
            balanceAfter > balanceBefore,
            "BOB should receive some rewards"
        );

        IVaultContract.PlayerStakeInfo memory playerInfo = vaultContract
            .getPlayerStakeInfo(GAME_ID, BOB);
        assertTrue(playerInfo.hasClaimed, "Should be marked as claimed");
    }

    function test_ClaimRewards_ShouldRevertGameNotEnded() public {
        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        vm.expectRevert("Game has not ended yet");
        vm.prank(BOB);
        vaultContract.claimRewards(GAME_ID);
    }

    function test_ClaimRewards_ShouldRevertAlreadyClaimed() public {
        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        vm.prank(ALICE);
        vaultContract.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);

        vm.prank(BOB);
        vaultContract.claimRewards(GAME_ID);

        vm.expectRevert("Rewards already claimed");
        vm.prank(BOB);
        vaultContract.claimRewards(GAME_ID);
    }

    function test_ClaimRewards_ShouldRevertNoStake() public {
        vm.prank(ALICE);
        vaultContract.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);

        vm.expectRevert("No stake in this game");
        vm.prank(BOB);
        vaultContract.claimRewards(GAME_ID);
    }

    // ============ VIEW FUNCTION TESTS ============

    function test_GetGameVaultInfo_EmptyGame_ShouldReturnDefaults()
        public
        view
    {
        IVaultContract.GameVaultInfo memory gameInfo = vaultContract
            .getGameVaultInfo(999);
        assertEq(gameInfo.gameId, 0, "Should return empty game info");
        assertEq(gameInfo.totalStakes, 0, "Should have 0 total stakes");
        assertEq(gameInfo.playerCount, 0, "Should have 0 players");
        assertFalse(gameInfo.gameEnded, "Should not be ended");
    }

    function test_GetPlayerStakeInfo_NoStake_ShouldReturnEmpty() public view {
        IVaultContract.PlayerStakeInfo memory playerInfo = vaultContract
            .getPlayerStakeInfo(GAME_ID, BOB);
        assertEq(playerInfo.totalStaked, 0, "Should have 0 stake");
        assertEq(playerInfo.stakeCount, 0, "Should have 0 stake count");
        assertFalse(playerInfo.hasClaimed, "Should not have claimed");
    }

    function test_CanClaimRewards_ShouldWork() public {
        // Before staking
        assertFalse(
            vaultContract.canClaimRewards(GAME_ID, BOB),
            "Should not be able to claim before staking"
        );

        // After staking but before game ends
        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);
        assertFalse(
            vaultContract.canClaimRewards(GAME_ID, BOB),
            "Should not be able to claim before game ends"
        );

        // After game ends
        vm.prank(ALICE);
        vaultContract.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);
        assertTrue(
            vaultContract.canClaimRewards(GAME_ID, BOB),
            "Should be able to claim after game ends"
        );

        // After claiming
        vm.prank(BOB);
        vaultContract.claimRewards(GAME_ID);
        assertFalse(
            vaultContract.canClaimRewards(GAME_ID, BOB),
            "Should not be able to claim twice"
        );
    }

    // ============ INTEGRATION TESTS ============

    function test_Integration_MultipleGamesAndChains() public {
        uint256 gameId1 = 100;
        uint256 gameId2 = 200;

        // Game 1: BOB and CHARLIE
        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE}(gameId1, FIXED_STAKE);

        vm.prank(CHARLIE);
        vaultContract.stake{value: FIXED_STAKE}(gameId1, FIXED_STAKE);

        // Game 2: BOB and DAVID
        vm.prank(BOB);
        vaultContract.stake{value: FIXED_STAKE}(gameId2, FIXED_STAKE);

        vm.prank(DAVID);
        vaultContract.stake{value: FIXED_STAKE}(gameId2, FIXED_STAKE);

        // Verify game states
        assertEq(
            vaultContract.getTotalGameStakes(gameId1),
            FIXED_STAKE * 2,
            "Game 1 should have correct total"
        );
        assertEq(
            vaultContract.getTotalGameStakes(gameId2),
            FIXED_STAKE * 2,
            "Game 2 should have correct total"
        );

        // BOB should have stakes in both games
        assertEq(
            vaultContract.getPlayerStake(gameId1, BOB),
            FIXED_STAKE,
            "BOB should have stake in game 1"
        );
        assertEq(
            vaultContract.getPlayerStake(gameId2, BOB),
            FIXED_STAKE,
            "BOB should have stake in game 2"
        );

        // End games
        vm.prank(ALICE);
        vaultContract.endGame(gameId1, IVaultContract.GameResult.WhiteWins);

        vm.prank(ALICE);
        vaultContract.endGame(gameId2, IVaultContract.GameResult.Draw);

        // Verify end states
        assertTrue(
            vaultContract.isGameEnded(gameId1),
            "Game 1 should be ended"
        );
        assertTrue(
            vaultContract.isGameEnded(gameId2),
            "Game 2 should be ended"
        );

        // BOB should be able to claim from both games
        assertTrue(
            vaultContract.canClaimRewards(gameId1, BOB),
            "BOB should be able to claim from game 1"
        );
        assertTrue(
            vaultContract.canClaimRewards(gameId2, BOB),
            "BOB should be able to claim from game 2"
        );
    }

    function test_Integration_CompleteVaultFlow() public {
        uint256 bobStake = 2 ether;
        uint256 charlieStake = 3 ether;

        // 1. Players stake
        vm.prank(BOB);
        vaultContract.stake{value: bobStake}(GAME_ID, bobStake);

        vm.prank(CHARLIE);
        vaultContract.stake{value: charlieStake}(GAME_ID, charlieStake);

        // 2. Verify staking state
        assertEq(
            vaultContract.getTotalGameStakes(GAME_ID),
            bobStake + charlieStake,
            "Total stakes should be correct"
        );

        IVaultContract.GameVaultInfo memory gameInfo = vaultContract
            .getGameVaultInfo(GAME_ID);
        assertEq(
            gameInfo.totalStakes,
            bobStake + charlieStake,
            "Game info should show correct total"
        );
        assertEq(gameInfo.playerCount, 2, "Should have 2 players");
        assertFalse(gameInfo.gameEnded, "Game should not be ended yet");

        // 3. End game
        vm.prank(ALICE);
        vaultContract.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);

        // 4. Verify end state
        gameInfo = vaultContract.getGameVaultInfo(GAME_ID);
        assertTrue(gameInfo.gameEnded, "Game should be ended");
        assertEq(
            uint8(gameInfo.result),
            uint8(IVaultContract.GameResult.WhiteWins),
            "Result should be correct"
        );

        // 5. Claim rewards
        uint256 bobBalanceBefore = BOB.balance;
        uint256 charlieBalanceBefore = CHARLIE.balance;

        vm.prank(BOB);
        vaultContract.claimRewards(GAME_ID);

        vm.prank(CHARLIE);
        vaultContract.claimRewards(GAME_ID);

        // 6. Verify rewards distributed
        assertTrue(
            BOB.balance != bobBalanceBefore,
            "BOB balance should change"
        );
        assertTrue(
            CHARLIE.balance != charlieBalanceBefore,
            "CHARLIE balance should change"
        );

        // 7. Verify claim status
        IVaultContract.PlayerStakeInfo memory bobInfo = vaultContract
            .getPlayerStakeInfo(GAME_ID, BOB);
        IVaultContract.PlayerStakeInfo memory charlieInfo = vaultContract
            .getPlayerStakeInfo(GAME_ID, CHARLIE);

        assertTrue(bobInfo.hasClaimed, "BOB should have claimed");
        assertTrue(charlieInfo.hasClaimed, "CHARLIE should have claimed");
    }
}
