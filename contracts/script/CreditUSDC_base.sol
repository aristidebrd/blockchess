// scripts/CreditRealUSDC.s.sol
pragma solidity ^0.8.0;

import "forge-std/Script.sol";

interface IERC20 {
    function transfer(address to, uint amount) external returns (bool);
    function balanceOf(address account) external view returns (uint);
}

contract CreditRealUSDC is Script {
    address usdcAddress = 0x036CbD53842c5426634e7929541eC2318f3dCF7e; // USDC on Base Sepolia
    address richUSDCAccount = 0x52EdA770E87565ddB61cc1E9011192c5e3D5CbEc; // une "whale" sur Base Sepolia avec du USDC
    address account1 = 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266; // une des adresses Anvil
    address account2 = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8; // une des adresses Anvil

    function run() external {
        // On Anvil fork, we can impersonate any account without needing private keys
        vm.prank(richUSDCAccount);
        IERC20(usdcAddress).transfer(account1, 10000e6); // 1000 USDC

        vm.prank(richUSDCAccount);
        IERC20(usdcAddress).transfer(account2, 10000e6); // 1000 USDC
    }
}
