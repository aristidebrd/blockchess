import React from 'react';
import { Wallet, CheckCircle } from 'lucide-react';
import { ConnectButton } from '@rainbow-me/rainbowkit';
import { useAccount, useBalance } from 'wagmi';
import { formatEther } from 'viem';

interface WalletConnectProps {
  onWalletChange?: (wallet: any | null) => void;
  isCompact?: boolean;
}

const WalletConnect: React.FC<WalletConnectProps> = ({ onWalletChange, isCompact = false }) => {
  const { address, isConnected } = useAccount();
  const { data: balance } = useBalance({
    address: address,
  });

  // Update parent component when wallet state changes (if callback provided)
  React.useEffect(() => {
    if (onWalletChange) {
      if (isConnected && address && balance) {
        onWalletChange({
          address: `${address.slice(0, 6)}...${address.slice(-4)}`,
          balance: parseFloat(formatEther(balance.value)).toFixed(4),
          isConnected: true
        });
      } else {
        onWalletChange(null);
      }
    }
  }, [isConnected, address, balance, onWalletChange]);

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
                          {account.displayBalance && (
                            <span className="text-green-400 font-bold text-sm">{account.displayBalance}</span>
                          )}
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
            <div className="flex justify-between items-center">
              <span className="text-gray-300 text-sm">Balance</span>
              <span className="text-green-400 font-bold">
                {balance ? `${parseFloat(formatEther(balance.value)).toFixed(4)} ETH` : '0.0000 ETH'}
              </span>
            </div>
          </div>

          {/* Betting info */}
          <div className="bg-gradient-to-r from-yellow-900/30 to-orange-900/30 rounded-lg p-4 border border-yellow-500/30">
            <h4 className="text-yellow-400 font-semibold mb-3">Betting Information</h4>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-300">Cost per vote:</span>
                <span className="text-white font-medium">0.01 ETH</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-300">Network:</span>
                <span className="text-blue-400 font-medium">Ethereum</span>
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
