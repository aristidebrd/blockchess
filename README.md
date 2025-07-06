# BlockChess - Community-Driven Chess Voting Game

A real-time collaborative chess game where users vote for the next move. Built with React + Vite + TypeScript on the frontend, Go (Golang) backend, and Solidity smart contracts for blockchain integration with multi-chain USDC support via Circle's CCTP.

## 🏗️ Architecture

### **System Overview**
```
[Player Browser] ←→ [Frontend: React + TS] ←→ [Backend: Go + WebSocket] ←→ [Smart Contracts]
     │                      │                          │                        │
     │                      │                          │                        │
     ▼                      ▼                          ▼                        ▼
[WalletConnect]     [Game UI + Voting]        [Game Logic + Timer]      [USDC Staking + CCTP]
```

### **Application Layers**

#### **1. Application Layer (`App.tsx`)**
- Wallet connection management
- Games list display & WebSocket subscription  
- Total players count tracking
- Matchmaking logic (start/cancel/callbacks)
- Navigation between lobby/game views
- Global header and footer components

#### **2. Game Engine (`gameService.ts`)**
- Game state management (timer, turns, voting)
- WebSocket game-specific subscriptions
- Move validation and creation logic
- Player role management within games
- Board state management
- Vote tracking and validation

#### **3. Presentation Layer (`GameView.tsx`)**
- Pure presentation of game state
- User interactions (click handlers)
- UI components (ChessBoard, VotingPanel, Timer)
- Move confirmation dialogs

#### **4. Blockchain Layer (Smart Contracts)**
- **GameFactory**: Creates and manages game instances
- **GameContract**: Individual game state and voting tracking
- **VaultContract**: USDC staking and cross-chain rewards via CCTP

### **Data Flow**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   WebSocket     │───▶│   GameService   │───▶│   GameView      │
│   Backend       │    │   (Engine)      │    │   (UI)          │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       ▼                       │
         │              ┌─────────────────┐              │
         └─────────────▶│     App.tsx     │◀─────────────┘
                        │ (Application)   │
                        │ + Header/Footer │
                        └─────────────────┘
                                 │
                                 ▼
                        ┌─────────────────┐
                        │ Smart Contracts │
                        │ (Blockchain)    │
                        └─────────────────┘
```

## 🚀 Features

### **Core Gameplay**
- **Real-time voting**: Players vote for chess moves in real-time
- **WebSocket communication**: Live updates for votes, moves, and timer
- **Team-based gameplay**: Join white or black teams
- **15-second turns**: Each turn has a 15-second timer
- **Game lobby**: Browse and join active games
- **Matchmaking**: Quick play system

### **Blockchain Integration**
- **USDC staking**: Stake USDC to join games
- **Multi-chain support**: 11 supported testnets via Circle's CCTP
- **Cross-chain rewards**: Automatic USDC distribution to winners
- **Smart contract voting**: On-chain vote tracking and game results
- **Wallet integration**: Connect via WalletConnect/RainbowKit

### **Supported Networks**
- Ethereum Sepolia
- Base Sepolia
- Optimism Sepolia
- Arbitrum Sepolia
- Avalanche Fuji
- Polygon Amoy
- Unichain Sepolia
- Linea Sepolia
- Sonic Testnet
- World Chain Sepolia
- Codex Testnet

## 🛠️ Tech Stack

### **Frontend**
- React + TypeScript
- Vite for fast development and building
- RainbowKit + Wagmi for wallet integration
- WebSocket client for real-time updates
- Tailwind CSS for styling

### **Backend**
- Go (Golang)
- Gorilla WebSocket for WebSocket handling
- Gorilla Mux for HTTP routing
- In-memory game state management
- Blockchain integration via Go-Ethereum

### **Smart Contracts**
- Solidity ^0.8.30
- Foundry for development and deployment
- Circle's CCTP for cross-chain USDC transfers
- Multi-chain deployment support

### **Development Tools**
- Foundry (forge, anvil, cast)
- Go bindings generation
- Automated deployment scripts
- Local blockchain forking

## 🚀 Getting Started

### **Prerequisites**
- Node.js 18+ and npm
- Go 1.23+
- Foundry (for blockchain features) - [Install here](https://getfoundry.sh/)

### **Development Setup**

1. **Clone the repository**
```bash
git clone <repository-url>
cd blockchess
```

2. **Install frontend dependencies**
```bash
npm install
```

3. **Start local blockchain (optional - for full blockchain testing)**
```bash
# Start Anvil with Base Sepolia fork
./credit_usdc.sh start
```

4. **Deploy smart contracts (optional)**
```bash
# Deploy to local Anvil
./deploy-contracts.sh

# Or deploy to specific testnet
./deploy-game-factory.sh
```

5. **Start the Go backend**
```bash
go run main.go
```

6. **Start frontend development server**
```bash
npm run dev
```

The application will be available at:
- Frontend: http://localhost:5173
- Backend: http://localhost:8080

### **Production Build**

1. **Build the frontend**
```bash
npm run build
```

2. **Build the Go backend**
```bash
go build -o blockchess-server main.go
```

3. **Run the server**
```bash
./blockchess-server -addr=:8080
```

The entire application will be served on http://localhost:8080

## 💰 USDC Testing Setup

For testing blockchain features, you'll need USDC tokens. We provide a convenient script:

### **Quick Start**
```bash
# Start Anvil and credit accounts with 10,000 USDC each
./credit_usdc.sh start

# Check USDC balances
./credit_usdc.sh balance

# Stop Anvil when done
./credit_usdc.sh stop
```

### **Available Commands**
- `./credit_usdc.sh start` - Start Anvil with Base Sepolia fork and credit accounts
- `./credit_usdc.sh balance` - Check current USDC balances
- `./credit_usdc.sh deploy` - Deploy CreditUSDC contract
- `./credit_usdc.sh status` - Check Anvil status
- `./credit_usdc.sh stop` - Stop Anvil

### **MetaMask Setup**
1. **Add Anvil Network**:
   - Network Name: `Anvil Local`
   - RPC URL: `http://127.0.0.1:8545`
   - Chain ID: `31337`
   - Currency Symbol: `ETH`

2. **Import USDC Token**:
   - Contract Address: `0x036CbD53842c5426634e7929541eC2318f3dCF7e`
   - Symbol: `USDC`
   - Decimals: `6`

## 🏗️ Smart Contract Deployment

### **Local Development**
```bash
# Deploy all contracts to local Anvil
./deploy-contracts.sh
```

### **Testnet Deployment**
```bash
# Deploy GameFactory to Base Sepolia
./deploy-game-factory.sh

# Deploy VaultContract to all supported testnets
./deploy-cctp-vault.sh all

# Deploy to specific testnet
./deploy-cctp-vault.sh deploy base-sepolia
```

### **Environment Configuration**
Create a `.env` file:
```env
# Blockchain Configuration
PRIVATE_KEY=your_private_key_here
GO_BACKEND_ADDRESS=your_backend_address_here

# RPC URLs for different networks
BASE_SEPOLIA_RPC_URL=https://sepolia.base.org
ETHEREUM_SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/your_key
# ... other networks

# Contract Addresses (auto-populated by deployment scripts)
GAME_CONTRACT_ADDRESS=0x...
VAULT_CONTRACT_ADDRESS=0x...
```

## 📁 Project Structure

```
blockchess/
├── src/                          # Frontend React source
│   ├── components/               # React components
│   │   ├── GameView.tsx         # Main game interface
│   │   ├── GameLobby.tsx        # Game lobby
│   │   ├── ChessBoard.tsx       # Chess board component
│   │   ├── VotingPanel.tsx      # Voting interface
│   │   └── ApprovalFlow.tsx     # USDC approval flow
│   ├── services/                # Frontend services
│   │   ├── gameService.ts       # Game logic service
│   │   └── websocket.ts         # WebSocket communication
│   ├── contexts/                # React contexts
│   └── utils/                   # Utility functions
├── internal/                    # Go backend source
│   ├── game/                    # Game logic
│   │   ├── manager.go          # Game manager
│   │   ├── blockchain.go       # Blockchain integration
│   │   └── config.go           # Configuration
│   └── websocket/              # WebSocket handling
│       ├── hub.go              # WebSocket hub
│       └── client.go           # Client management
├── contracts/                   # Smart contracts
│   ├── src/                    # Solidity source
│   │   ├── GameFactory.sol     # Game factory contract
│   │   ├── GameContract.sol    # Individual game contract
│   │   └── VaultContract.sol   # USDC vault with CCTP
│   ├── script/                 # Deployment scripts
│   └── common/interfaces/      # Contract interfaces
├── contracts-bindings/         # Go contract bindings
├── main.go                     # Go backend entry point
├── package.json               # Frontend dependencies
├── go.mod                     # Go dependencies
├── foundry.toml              # Foundry configuration
└── deployment scripts/       # Automated deployment
    ├── deploy-contracts.sh
    ├── deploy-game-factory.sh
    ├── deploy-cctp-vault.sh
    └── credit_usdc.sh
```

## 🎮 How to Play

1. **Connect your wallet** using WalletConnect/RainbowKit
2. **Browse games** in the lobby or use "Quick Play"
3. **Stake USDC** to join a game (testnet only)
4. **Join a team** (white or black)
5. **Vote for moves** by clicking on the chess board
6. **Watch the timer** - most voted move executes when time runs out
7. **Win rewards** - USDC distributed to winning team

## 🔧 Configuration

### **Backend Configuration**
- **Port**: Change with `-addr` flag (default: `:8080`)
- **Turn duration**: Modify `GameTimerSeconds` in `internal/game/manager.go`
- **Blockchain**: Configure via `.env` file

### **Smart Contract Configuration**
- **Stake amounts**: Configurable per game
- **Supported chains**: Add new chains in `VaultContract.sol`
- **CCTP domains**: Configured for all supported testnets

## 🚢 Deployment

### **Single Binary Deployment**
```bash
# Build everything
npm run build
go build -o blockchess-server main.go

# Deploy
./blockchess-server -addr=:8080
```

### **Docker Deployment**
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN npm install && npm run build
RUN go build -o blockchess-server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/blockchess-server .
COPY --from=builder /app/dist ./dist
CMD ["./blockchess-server"]
```

## 🧪 Testing

### **Frontend Testing**
```bash
npm run test
```

### **Smart Contract Testing**
```bash
forge test
```

### **Integration Testing**
```bash
# Start full stack locally
./credit_usdc.sh start
./deploy-contracts.sh
go run main.go &
npm run dev
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

MIT License

---

**Built with ❤️ for ETHGlobal Cannes 2025**
