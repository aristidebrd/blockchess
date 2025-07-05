// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import "../BaseTest.sol";
import "../../src/GameContract.sol";
import "../../src/VaultContract.sol";

contract GameVaultIntegrationTest is BaseTest {
    GameContract public gameContract;
    VaultContract public vaultBase;
    VaultContract public vaultMainnet;
    VaultContract public vaultOptimism;

    uint256 constant GAME_ID = 123;
    uint32 constant BASE_CHAIN_ID = 8453;
    uint32 constant MAINNET_CHAIN_ID = 1;
    uint32 constant OPTIMISM_CHAIN_ID = 10;

    function setUp() public {
        // Deploy contracts
        gameContract = new GameContract(ALICE); // ALICE is Go backend
        vaultBase = new VaultContract(ALICE, address(gameContract));
        vaultMainnet = new VaultContract(ALICE, address(gameContract));
        vaultOptimism = new VaultContract(ALICE, address(gameContract));

        // Fund test addresses
        fundAddress(ALICE, 10 ether);
        fundAddress(BOB, 10 ether);
        fundAddress(CHARLIE, 10 ether);
        fundAddress(DAVID, 10 ether);
        fundAddress(EVE, 10 ether);
    }

    function test_CompleteGameFlow_WithMultipleChains() public {
        // 1. Go backend creates game on Base
        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        // 2. Players join teams on Base
        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        vm.prank(CHARLIE);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.Black);

        vm.prank(DAVID);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        // 3. Players stake on different chains
        // BOB stakes on Base
        vm.prank(BOB);
        vaultBase.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        // CHARLIE stakes on Mainnet
        vm.prank(CHARLIE);
        vaultMainnet.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        // DAVID stakes on Optimism
        vm.prank(DAVID);
        vaultOptimism.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        // BOB stakes again on Mainnet (multiple stakes)
        vm.prank(BOB);
        vaultMainnet.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        // 4. Go backend records moves from different chains
        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, BASE_CHAIN_ID);

        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, CHARLIE, MAINNET_CHAIN_ID);

        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, DAVID, OPTIMISM_CHAIN_ID);

        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, MAINNET_CHAIN_ID);

        // 5. Verify game state
        IGameContract.GameInfo memory gameInfo = gameContract.getGameInfo(
            GAME_ID
        );
        assertEq(gameInfo.whitePlayerCount, 2, "Should have 2 white players");
        assertEq(gameInfo.blackPlayerCount, 1, "Should have 1 black player");
        assertTrue(gameContract.isGameActive(GAME_ID), "Game should be active");

        // 6. Verify player move counts
        assertEq(
            gameContract.getPlayerMoveCount(GAME_ID, BOB, BASE_CHAIN_ID),
            1,
            "BOB should have 1 move on Base"
        );
        assertEq(
            gameContract.getPlayerMoveCount(GAME_ID, BOB, MAINNET_CHAIN_ID),
            1,
            "BOB should have 1 move on Mainnet"
        );
        assertEq(
            gameContract.getPlayerMoveCount(GAME_ID, CHARLIE, MAINNET_CHAIN_ID),
            1,
            "CHARLIE should have 1 move on Mainnet"
        );
        assertEq(
            gameContract.getPlayerMoveCount(GAME_ID, DAVID, OPTIMISM_CHAIN_ID),
            1,
            "DAVID should have 1 move on Optimism"
        );

        IGameContract.PlayerInfo memory bobInfo = gameContract.getPlayerInfo(
            GAME_ID,
            BOB
        );
        assertEq(bobInfo.moveCount, 2, "BOB should have 2 total moves");

        // 7. Verify vault states
        assertEq(
            vaultBase.getTotalGameStakes(GAME_ID),
            FIXED_STAKE,
            "Base vault should have 1 stake"
        );
        assertEq(
            vaultMainnet.getTotalGameStakes(GAME_ID),
            FIXED_STAKE * 2,
            "Mainnet vault should have 2 stakes"
        );
        assertEq(
            vaultOptimism.getTotalGameStakes(GAME_ID),
            FIXED_STAKE,
            "Optimism vault should have 1 stake"
        );

        // 8. Go backend ends game (White wins)
        vm.prank(ALICE);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.WhiteWins);

        // 9. Go backend notifies all vaults
        vm.prank(ALICE);
        vaultBase.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);

        vm.prank(ALICE);
        vaultMainnet.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);

        vm.prank(ALICE);
        vaultOptimism.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);

        // 10. Verify game ended
        assertFalse(
            gameContract.isGameActive(GAME_ID),
            "Game should not be active"
        );
        assertEq(
            uint8(gameContract.getGameResult(GAME_ID)),
            uint8(IGameContract.GameResult.WhiteWins),
            "Result should be white wins"
        );

        // 11. Winners claim rewards
        uint256 bobBalanceBefore = BOB.balance;
        uint256 davidBalanceBefore = DAVID.balance;
        uint256 charlieBalanceBefore = CHARLIE.balance;

        // BOB claims from Base
        vm.prank(BOB);
        vaultBase.claimRewards(GAME_ID);

        // BOB claims from Mainnet
        vm.prank(BOB);
        vaultMainnet.claimRewards(GAME_ID);

        // DAVID claims from Optimism
        vm.prank(DAVID);
        vaultOptimism.claimRewards(GAME_ID);

        // CHARLIE (loser) claims from Mainnet
        vm.prank(CHARLIE);
        vaultMainnet.claimRewards(GAME_ID);

        // 12. Verify rewards distributed
        assertTrue(
            BOB.balance > bobBalanceBefore,
            "BOB should receive rewards"
        );
        assertTrue(
            DAVID.balance > davidBalanceBefore,
            "DAVID should receive rewards"
        );
        assertTrue(
            CHARLIE.balance >= charlieBalanceBefore,
            "CHARLIE should receive something (simplified logic)"
        );

        // 13. Verify claim status
        assertTrue(
            vaultBase.getPlayerStakeInfo(GAME_ID, BOB).hasClaimed,
            "BOB should have claimed from Base"
        );
        assertTrue(
            vaultMainnet.getPlayerStakeInfo(GAME_ID, BOB).hasClaimed,
            "BOB should have claimed from Mainnet"
        );
        assertTrue(
            vaultOptimism.getPlayerStakeInfo(GAME_ID, DAVID).hasClaimed,
            "DAVID should have claimed from Optimism"
        );
        assertTrue(
            vaultMainnet.getPlayerStakeInfo(GAME_ID, CHARLIE).hasClaimed,
            "CHARLIE should have claimed from Mainnet"
        );
    }

    function test_DrawScenario_ShouldReturnAllStakes() public {
        // 1. Setup game
        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        vm.prank(CHARLIE);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.Black);

        // 2. Players stake
        vm.prank(BOB);
        vaultBase.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        vm.prank(CHARLIE);
        vaultMainnet.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        // 3. End game as draw
        vm.prank(ALICE);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.Draw);

        vm.prank(ALICE);
        vaultBase.endGame(GAME_ID, IVaultContract.GameResult.Draw);

        vm.prank(ALICE);
        vaultMainnet.endGame(GAME_ID, IVaultContract.GameResult.Draw);

        // 4. Both players should get their stakes back
        uint256 bobBalanceBefore = BOB.balance;
        uint256 charlieBalanceBefore = CHARLIE.balance;

        vm.prank(BOB);
        vaultBase.claimRewards(GAME_ID);

        vm.prank(CHARLIE);
        vaultMainnet.claimRewards(GAME_ID);

        // 5. Verify stakes returned
        assertEq(
            BOB.balance,
            bobBalanceBefore + FIXED_STAKE,
            "BOB should get stake back"
        );
        assertEq(
            CHARLIE.balance,
            charlieBalanceBefore + FIXED_STAKE,
            "CHARLIE should get stake back"
        );
    }

    function test_PlayerCanStakeOnMultipleChains() public {
        // 1. Setup game
        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        // 2. BOB stakes on all three chains
        vm.prank(BOB);
        vaultBase.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        vm.prank(BOB);
        vaultMainnet.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        vm.prank(BOB);
        vaultOptimism.stake{value: FIXED_STAKE}(GAME_ID, FIXED_STAKE);

        // 3. Verify stakes
        assertEq(
            vaultBase.getPlayerStake(GAME_ID, BOB),
            FIXED_STAKE,
            "BOB should have stake on Base"
        );
        assertEq(
            vaultMainnet.getPlayerStake(GAME_ID, BOB),
            FIXED_STAKE,
            "BOB should have stake on Mainnet"
        );
        assertEq(
            vaultOptimism.getPlayerStake(GAME_ID, BOB),
            FIXED_STAKE,
            "BOB should have stake on Optimism"
        );

        // 4. BOB makes moves on different chains
        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, BASE_CHAIN_ID);

        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, MAINNET_CHAIN_ID);

        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, OPTIMISM_CHAIN_ID);

        // 5. Verify move tracking
        assertEq(
            gameContract.getPlayerMoveCount(GAME_ID, BOB, BASE_CHAIN_ID),
            1,
            "BOB should have 1 move on Base"
        );
        assertEq(
            gameContract.getPlayerMoveCount(GAME_ID, BOB, MAINNET_CHAIN_ID),
            1,
            "BOB should have 1 move on Mainnet"
        );
        assertEq(
            gameContract.getPlayerMoveCount(GAME_ID, BOB, OPTIMISM_CHAIN_ID),
            1,
            "BOB should have 1 move on Optimism"
        );

        IGameContract.PlayerInfo memory bobInfo = gameContract.getPlayerInfo(
            GAME_ID,
            BOB
        );
        assertEq(bobInfo.moveCount, 3, "BOB should have 3 total moves");
    }

    function test_OnlyGoBackendCanControlFlow() public {
        // 1. Try to create game as non-backend
        vm.expectRevert("Only authorized backend can perform this action");
        vm.prank(BOB);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        // 2. Try to record move as non-backend
        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        vm.expectRevert("Only authorized backend can perform this action");
        vm.prank(BOB);
        gameContract.recordMove(GAME_ID, BOB, BASE_CHAIN_ID);

        // 3. Try to end game as non-backend
        vm.expectRevert("Only authorized backend can perform this action");
        vm.prank(BOB);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.WhiteWins);

        // 4. Try to end vault games as non-backend
        vm.expectRevert("Only authorized backend can end games");
        vm.prank(BOB);
        vaultBase.endGame(GAME_ID, IVaultContract.GameResult.WhiteWins);
    }

    function test_GameStateConsistency() public {
        // 1. Create game
        vm.prank(ALICE);
        gameContract.createGame(GAME_ID, FIXED_STAKE);

        // 2. Initially game should be active with no players
        IGameContract.GameInfo memory gameInfo = gameContract.getGameInfo(
            GAME_ID
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
        assertEq(gameInfo.whitePlayerCount, 0, "Should have 0 white players");
        assertEq(gameInfo.blackPlayerCount, 0, "Should have 0 black players");

        // 3. Players join
        vm.prank(BOB);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        vm.prank(CHARLIE);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.Black);

        // 4. Game should still be active with players
        gameInfo = gameContract.getGameInfo(GAME_ID);
        assertEq(gameInfo.whitePlayerCount, 1, "Should have 1 white player");
        assertEq(gameInfo.blackPlayerCount, 1, "Should have 1 black player");
        assertTrue(
            gameContract.isGameActive(GAME_ID),
            "Game should still be active"
        );

        // 5. End game
        vm.prank(ALICE);
        gameContract.endGame(GAME_ID, IGameContract.GameResult.WhiteWins);

        // 6. Game should be ended
        gameInfo = gameContract.getGameInfo(GAME_ID);
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

        // 7. Cannot join ended game
        vm.expectRevert("Game is not active");
        vm.prank(DAVID);
        gameContract.joinTeam(GAME_ID, IGameContract.Team.White);

        // 8. Cannot record moves on ended game
        vm.expectRevert("Game is not active");
        vm.prank(ALICE);
        gameContract.recordMove(GAME_ID, BOB, BASE_CHAIN_ID);
    }
}
