import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import './index.css';

import '@rainbow-me/rainbowkit/styles.css';
import { getDefaultConfig, RainbowKitProvider, connectorsForWallets } from '@rainbow-me/rainbowkit';
import { injectedWallet } from '@rainbow-me/rainbowkit/wallets';
import { WagmiProvider, createConfig } from 'wagmi';
import {
  arbitrumSepolia,
  avalancheFuji,
  baseSepolia,
  sepolia,
  lineaSepolia,
  optimismSepolia,
  polygonAmoy
} from 'wagmi/chains';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { http } from 'wagmi';
import { OnlinePlayersProvider } from './contexts/OnlinePlayersContext.tsx';

// Determine environment
const isDev = import.meta.env.MODE === 'development';

// Local Foundry chain definition (no explicit type annotation)
const localChainBaseSepolia = {
  id: 84532,
  name: 'Anvil (Foundry - Base Sepolia)',
  network: 'foundry',
  nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
  rpcUrls: { default: { http: ['http://127.0.0.1:8545'] } },
};

// Local Foundry chain definition (no explicit type annotation)
const localChainOptimismSepolia = {
  id: 84532,
  name: 'Anvil (Foundry - Optimism Sepolia)',
  network: 'foundry',
  nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
  rpcUrls: { default: { http: ['http://127.0.0.1:8546'] } },
};

const sonicTestnet = {
  id: 64165,
  name: 'Sonic Testnet',
  network: 'sonic-testnet',
  nativeCurrency: { name: 'Sonic', symbol: 'S', decimals: 18 },
  rpcUrls: { default: { http: ['https://rpc.testnet.soniclabs.com'] } },
  blockExplorers: { default: { name: 'Sonic Explorer', url: 'https://testnet.soniclabs.com' } },
};

const unichainSepolia = {
  id: 1301,
  name: 'Unichain Sepolia',
  network: 'unichain-sepolia',
  nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
  rpcUrls: { default: { http: ['https://sepolia.unichain.org'] } },
  blockExplorers: { default: { name: 'Unichain Explorer', url: 'https://sepolia.uniscan.xyz' } },
};

const worldChainSepolia = {
  id: 4801,
  name: 'World Chain Sepolia',
  network: 'world-chain-sepolia',
  nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
  rpcUrls: { default: { http: ['https://worldchain-sepolia.gateway.tenderly.co'] } },
  blockExplorers: { default: { name: 'World Chain Explorer', url: 'https://sepolia.worldscan.org' } },
};

// Define dev and prod chains as readonly tuples
const devChains = [
  localChainBaseSepolia,
  localChainOptimismSepolia,
] as const;
const prodChains = [
  // Testnets
  sepolia,                // Ethereum Sepolia
  arbitrumSepolia,        // Arbitrum Sepolia
  avalancheFuji,          // Avalanche Fuji
  baseSepolia,            // Base Sepolia
  lineaSepolia,           // Linea Sepolia
  optimismSepolia,        // OP Sepolia
  polygonAmoy,            // Polygon PoS Amoy
  sonicTestnet,           // Sonic Testnet
  unichainSepolia,        // Unichain Sepolia
  worldChainSepolia,      // World Chain Sepolia
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
      // Testnets - using Alchemy endpoints
      [sepolia.id]: http(`https://eth-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [arbitrumSepolia.id]: http(`https://arb-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [avalancheFuji.id]: http(`https://avax-fuji.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [baseSepolia.id]: http(`https://base-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [lineaSepolia.id]: http(`https://linea-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [optimismSepolia.id]: http(`https://opt-sepolia.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [polygonAmoy.id]: http(`https://polygon-amoy.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [sonicTestnet.id]: http(`https://sonic-blaze.g.alchemy.com/v2/Kv5V9_Sv55ZosXpoVga3IfXYO6JV9gNJ`),
      [unichainSepolia.id]: http(`https://eth-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7`),
      [worldChainSepolia.id]: http(`https://eth-mainnet.g.alchemy.com/v2/22r8dairb21cjlw7`),
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
