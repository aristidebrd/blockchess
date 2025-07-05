import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import './index.css';

import '@rainbow-me/rainbowkit/styles.css';
import { getDefaultConfig, RainbowKitProvider, connectorsForWallets } from '@rainbow-me/rainbowkit';
import { injectedWallet } from '@rainbow-me/rainbowkit/wallets';
import { WagmiProvider, createConfig } from 'wagmi';
import { mainnet, polygon, optimism, arbitrum, base } from 'wagmi/chains';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { http } from 'wagmi';
import { OnlinePlayersProvider } from './contexts/OnlinePlayersContext.tsx';

// Determine environment
const isDev = import.meta.env.MODE === 'development';

// Local Foundry chain definition (no explicit type annotation)
const localChain = {
  id: 84532,
  name: 'Anvil (Foundry)',
  network: 'foundry',
  nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
  rpcUrls: { default: { http: ['http://127.0.0.1:8545'] } },
};

// Define dev and prod chains as readonly tuples
const devChains = [
  localChain,
] as const;
const prodChains = [
  mainnet,
  polygon,
  optimism,
  arbitrum,
  base,
] as const;

// WalletConnect Project ID (set via Vite env or fallback)
const projectId = import.meta.env.VITE_WALLETCONNECT_PROJECT_ID ||
  'b66e4e5ce122268694ffabeb5b666b7c';

// Generate wagmi + RainbowKit config per environment
const wagmiRainbowConfig = isDev
  ? (() => {
    // Development: Use injected wallet (MetaMask) only
    const connectors = connectorsForWallets(
      [
        {
          groupName: 'Development',
          wallets: [injectedWallet],
        },
      ],
      {
        appName: 'BlockChess Dev',
        projectId,
      }
    );

    return createConfig({
      connectors,
      chains: devChains,
      transports: Object.fromEntries(
        devChains.map((chain) => [chain.id, http(chain.rpcUrls.default.http[0])])
      ),
    });
  })()
  : getDefaultConfig({
    // Production: Use full RainbowKit configuration with WalletConnect and custom Alchemy RPCs
    appName: 'BlockChess',
    projectId,
    chains: prodChains,
    transports: {
      [mainnet.id]: http(`https://eth-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [polygon.id]: http(`https://polygon-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [optimism.id]: http(`https://opt-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [arbitrum.id]: http(`https://arb-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [base.id]: http(`https://base-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7`),
    },
    ssr: false,
  });

// React Query client
const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <WagmiProvider config={wagmiRainbowConfig}>
      <QueryClientProvider client={queryClient}>
        <RainbowKitProvider>
          <OnlinePlayersProvider>
            <App />
          </OnlinePlayersProvider>
        </RainbowKitProvider>
      </QueryClientProvider>
    </WagmiProvider>
  </StrictMode>
);
