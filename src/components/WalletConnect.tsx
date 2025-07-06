import React, { useState, useEffect } from 'react';
import { Wallet, CheckCircle, RefreshCw } from 'lucide-react';
import { ConnectButton } from '@rainbow-me/rainbowkit';
import { useAccount, useBalance, useReadContract, useChainId } from 'wagmi';
import { formatEther, formatUnits } from 'viem';

// USDC Contract Addresses by Chain ID
const USDC_ADDRESSES: { [chainId: number]: string } = {
  // Local Anvil chains (using Base Sepolia USDC for both)
  84532: '0x036CbD53842c5426634e7929541eC2318f3dCF7e',   // Base Sepolia / Local Base
  11155420: '0x5fd84259d66Cd46123540766Be93DFE6D43130D7', // OP Sepolia / Local OP
  // Testnet addresses
  11155111: '0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238', // Ethereum Sepolia
  43113: '0x5425890298aed601595a70AB815c96711a31Bc65',    // Avalanche Fuji
  421614: '0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d',   // Arbitrum Sepolia
  80002: '0x41E94Eb019C0762f9Bfcf9Fb1E58725BfB0e7582',    // Polygon Amoy
  1301: '0x31d0220469e10c4E71834a79b1f276d740d3768F',     // Unichain Sepolia
  59141: '0xFEce4462D57bD51A6A552365A011b95f0E16d9B7',    // Linea Sepolia
  325000: '0x6d7f141b6819C2c9CC2f818e6ad549E7Ca090F8f',   // Codex Testnet
  64165: '0xA4879Fed32Ecbef99399e5cbC247E533421C4eC6',    // Sonic Testnet
  4801: '0x66145f38cBAC35Ca6F1Dfb4914dF98F1614aeA88',     // World Chain Sepolia
};

// ERC20 ABI for balanceOf function
const ERC20_ABI = [
  {
    name: 'balanceOf',
    type: 'function',
    inputs: [{ name: 'account', type: 'address' }],
    outputs: [{ name: '', type: 'uint256' }],
    stateMutability: 'view',
  },
] as const;

interface WalletConnectProps {
  onWalletChange?: (wallet: any | null) => void;
  isCompact?: boolean;
}

const WalletConnect: React.FC<WalletConnectProps> = ({ onWalletChange, isCompact = false }) => {
  const { address, isConnected } = useAccount();
  const chainId = useChainId();
  const [usdcBalance, setUsdcBalance] = useState<string>('0.0');
  const [isLoadingUSDC, setIsLoadingUSDC] = useState(false);

  // Get current chain's USDC address
  const usdcAddress = chainId ? USDC_ADDRESSES[chainId] : undefined;

  // Get ETH balance (keep for reference)
  const { data: ethBalance } = useBalance({
    address: address,
  });

  // Get USDC balance using contract read
  const { data: usdcBalanceData, refetch: refetchUSDC } = useReadContract({
    address: usdcAddress as `0x${string}`,
    abi: ERC20_ABI,
    functionName: 'balanceOf',
    args: address ? [address] : undefined,
    query: {
      enabled: !!address && isConnected && !!usdcAddress,
      refetchInterval: 10000, // Refetch every 10 seconds
    },
  });

  // Format USDC balance (6 decimals)
  useEffect(() => {
    if (usdcBalanceData) {
      const formatted = formatUnits(usdcBalanceData, 6);
      setUsdcBalance(parseFloat(formatted).toFixed(6));
    } else {
      setUsdcBalance('0.0');
    }
  }, [usdcBalanceData]);

  // Manual refresh function
  const refreshBalance = async () => {
    setIsLoadingUSDC(true);
    try {
      await refetchUSDC();
    } finally {
      setIsLoadingUSDC(false);
    }
  };

  // Update parent component when wallet state changes
  useEffect(() => {
    if (onWalletChange) {
      if (isConnected && address) {
        onWalletChange({
          address: `${address.slice(0, 6)}...${address.slice(-4)}`,
          balance: usdcBalance,
          ethBalance: ethBalance ? parseFloat(formatEther(ethBalance.value)).toFixed(4) : '0.0',
          isConnected: true
        });
      } else {
        onWalletChange(null);
      }
    }
  }, [isConnected, address, usdcBalance, ethBalance, onWalletChange]);

  if (isCompact) {
    return (
      <div className="flex items-center space-x-3">
        <ConnectButton.Custom>
          {({
            account,
            chain,
            openAccountModal,
            openChainModal,
            openConnectModal,
            authenticationStatus,
            mounted,
          }) => {
            const ready = mounted && authenticationStatus !== 'loading';
            const connected =
              ready &&
              account &&
              chain &&
              (!authenticationStatus ||
                authenticationStatus === 'authenticated');

            return (
              <div
                {...(!ready && {
                  'aria-hidden': true,
                  'style': {
                    opacity: 0,
                    pointerEvents: 'none',
                    userSelect: 'none',
                  },
                })}
              >
                {(() => {
                  if (!connected) {
                    return (
                      <button
                        onClick={openConnectModal}
                        type="button"
                        className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white font-medium py-2 px-4 rounded-lg transition-all duration-200 flex items-center space-x-2"
                      >
                        <Wallet className="w-4 h-4" />
                        <span>Connect Wallet</span>
                      </button>
                    );
                  }

                  if (chain.unsupported) {
                    return (
                      <button
                        onClick={openChainModal}
                        type="button"
                        className="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded-lg transition-all duration-200"
                      >
                        Wrong network
                      </button>
                    );
                  }

                  return (
                    <div className="flex items-center space-x-3">
                      <div className="bg-gray-800/80 rounded-lg px-3 py-2 border border-gray-700">
                        <div className="flex items-center space-x-2">
                          <CheckCircle className="w-4 h-4 text-green-400" />
                          <span className="text-white font-mono text-sm">{account.displayName}</span>
                          <span className="text-green-400 font-bold text-sm">{usdcBalance} USDC</span>
                        </div>
                      </div>
                      <button
                        onClick={openAccountModal}
                        type="button"
                        className="bg-gray-700 hover:bg-gray-600 text-white font-medium py-2 px-3 rounded-lg transition-all duration-200 text-sm"
                      >
                        Disconnect
                      </button>
                    </div>
                  );
                })()}
              </div>
            );
          }}
        </ConnectButton.Custom>
      </div>
    );
  }

  return (
    <div className="bg-gradient-to-br from-gray-800 to-gray-900 rounded-xl p-6 shadow-xl">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-3">
          <Wallet className="w-6 h-6 text-yellow-400" />
          <h3 className="text-xl font-bold text-white">Wallet</h3>
        </div>

        {isConnected && (
          <div className="flex items-center space-x-2">
            <CheckCircle className="w-5 h-5 text-green-400" />
            <span className="text-green-400 text-sm font-medium">Connected</span>
          </div>
        )}
      </div>

      {!isConnected ? (
        // Not connected state
        <div className="space-y-4">
          <div className="text-center py-6">
            <Wallet className="w-16 h-16 text-gray-600 mx-auto mb-4" />
            <p className="text-gray-300 mb-2">Connect your wallet to participate</p>
            <p className="text-sm text-gray-500">
              You need to connect a wallet to place votes and participate in the game
            </p>
          </div>

          <ConnectButton.Custom>
            {({ openConnectModal }) => (
              <button
                onClick={openConnectModal}
                type="button"
                className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white font-bold py-3 px-6 rounded-lg transition-all duration-200 flex items-center justify-center space-x-2"
              >
                <Wallet className="w-5 h-5" />
                <span>Connect Wallet</span>
              </button>
            )}
          </ConnectButton.Custom>
        </div>
      ) : (
        // Connected state
        <div className="space-y-4">
          {/* Wallet info */}
          <div className="bg-gray-800/50 rounded-lg p-4 border border-gray-700">
            <div className="flex justify-between items-center mb-2">
              <span className="text-gray-300 text-sm">Address</span>
              <span className="text-white font-mono text-sm">
                {address ? `${address.slice(0, 6)}...${address.slice(-4)}` : ''}
              </span>
            </div>

            {/* USDC Balance */}
            <div className="flex justify-between items-center mb-2">
              <span className="text-gray-300 text-sm">USDC Balance</span>
              <div className="flex items-center space-x-2">
                <span className="text-green-400 font-bold">
                  {usdcBalance} USDC
                </span>
                <button
                  onClick={refreshBalance}
                  disabled={isLoadingUSDC}
                  className="text-gray-400 hover:text-white transition-colors"
                  title="Refresh balance"
                >
                  <RefreshCw className={`w-4 h-4 ${isLoadingUSDC ? 'animate-spin' : ''}`} />
                </button>
              </div>
            </div>

            {/* ETH Balance (for reference) */}
            <div className="flex justify-between items-center">
              <span className="text-gray-300 text-sm">ETH Balance</span>
              <span className="text-blue-400 font-bold">
                {ethBalance ? `${parseFloat(formatEther(ethBalance.value)).toFixed(4)} ETH` : '0.0 ETH'}
              </span>
            </div>
          </div>

          {/* Betting info */}
          <div className="bg-gradient-to-r from-yellow-900/30 to-orange-900/30 rounded-lg p-4 border border-yellow-500/30">
            <h4 className="text-yellow-400 font-semibold mb-3">Betting Information</h4>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-300">Cost per vote:</span>
                <span className="text-white font-medium">10.0 USDC</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-300">Network:</span>
                <span className="text-blue-400 font-medium">Chain {chainId}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-300">USDC Contract:</span>
                <span className="text-purple-400 font-mono text-xs">
                  {usdcAddress ? `${usdcAddress.slice(0, 6)}...${usdcAddress.slice(-4)}` : 'Not found'}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-300">Sufficient funds:</span>
                <span className={`font-medium ${parseFloat(usdcBalance) >= 10 ? 'text-green-400' : 'text-red-400'}`}>
                  {parseFloat(usdcBalance) >= 10 ? 'Yes' : 'No'}
                </span>
              </div>
            </div>
          </div>

          {/* Account modal button */}
          <ConnectButton.Custom>
            {({ openAccountModal }) => (
              <button
                onClick={openAccountModal}
                type="button"
                className="w-full bg-gray-700 hover:bg-gray-600 text-white font-medium py-2 px-4 rounded-lg transition-all duration-200"
              >
                Manage Wallet
              </button>
            )}
          </ConnectButton.Custom>
        </div>
      )}
    </div>
  );
};

export default WalletConnect;
