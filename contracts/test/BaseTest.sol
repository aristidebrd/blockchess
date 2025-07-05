// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import "forge-std/Test.sol";

contract BaseTest is Test {
    // Test addresses
    address constant ALICE = address(0x1);
    address constant BOB = address(0x2);
    address constant CHARLIE = address(0x3);
    address constant DAVID = address(0x4);
    address constant EVE = address(0x5);

    // Test values
    uint256 constant FIXED_STAKE = 0.1 ether;
    uint256 constant MAX_PLAYERS_PER_TEAM = 10;

    // Helper function to create a move string
    function createMove(
        string memory from,
        string memory to
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(from, "-", to));
    }

    // Helper function to fund an address
    function fundAddress(address addr, uint256 amount) internal {
        vm.deal(addr, amount);
    }

    // Helper function to expect an event
    function expectEvent(address emitter, bytes4 /* selector */) internal {
        vm.expectEmit(true, true, false, true, emitter);
    }
}
