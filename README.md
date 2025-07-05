# BlockChess - Community-Driven Chess Voting Game

A real-time collaborative chess game where users vote for the next move. Built with React + Vite + TypeScript on the frontend and Go (Golang) as the backend, with WebSocket communication for real-time updates.

## 🏗️ Architecture

```
[Player Browser]
     |
     | UI + WalletConnect
     v
[Frontend: Vite + React + TS]   ←──────→  [Backend: Go HTTP + WebSocket]
     |                                            |
     | WebSocket `/ws`                           | HTTP (serves index.html)
     v                                            v
[WebSocket connection to Go]               [Game logic, timers, voting]
```

## 🚀 Features

- **Real-time voting**: Players vote for chess moves in real-time
- **WebSocket communication**: Live updates for votes, moves, and timer
- **WalletConnect integration**: Connect crypto wallets (frontend-only for now)
- **30-second turns**: Each turn has a 30-second timer
- **Single binary deployment**: Go backend serves both API and frontend
- **Game Lobby**: Browse and join active games

## 🛠️ Tech Stack

### Frontend
- React + TypeScript
- Vite for fast development and building
- WalletConnect for wallet integration
- WebSocket client for real-time updates
- Tailwind CSS for styling

### Backend
- Go (Golang)
- Gorilla WebSocket for WebSocket handling
- Gorilla Mux for HTTP routing
- In-memory game state management

## 📡 WebSocket Protocol

### Client → Server Messages:
```json
{ "type": "join_game", "gameId": "game-123" }
{ "type": "vote_move", "move": "e2e4" }
```

### Server → Client Messages:
```json
{ "type": "vote_update", "votes": { "e2e4": 4, "d2d4": 2 } }
{ "type": "move_result", "move": "e2e4" }
{ "type": "timer_tick", "secondsLeft": 14 }
```

## 🚀 Getting Started

### Prerequisites
- Node.js 18+ and npm
- Go 1.23+
- Foundry (for Anvil local blockchain) - Install from [getfoundry.sh](https://getfoundry.sh/)

### Development Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd blockchess
```

2. **Install frontend dependencies**
```bash
npm install
```

3. **Start local blockchain (for complete test setup)**
```bash
npm run anvil
```
This starts a local Ethereum test network using Foundry's Anvil. Required for wallet connections and testing the full blockchain integration.

4. **Start the Go backend**
```bash
cd backend
go run cmd/server/main.go
```
The backend will run on http://localhost:8080

5. **In another terminal, run frontend development server**
```bash
npm run dev
```
The frontend will run on http://localhost:5173

⚠️ **Important**: In development, always access the app through http://localhost:5173. The Vite dev server will proxy WebSocket connections to the Go backend automatically.

### Production Build

1. **Build the frontend**
```bash
npm run build
```

2. **Copy dist to backend**
```bash
cp -r dist backend/
```

3. **Build the Go backend**
```bash
cd backend
go build -o blockchess-server cmd/server/main.go
```

4. **Run the server**
```bash
./blockchess-server -addr=:8080
```

The entire application (frontend + backend) will be served on http://localhost:8080

## 💰 Credit USDC for Testing

For testing the blockchain integration features, you'll need USDC tokens on your Anvil accounts. We provide a convenient script to credit your accounts with USDC from Base Sepolia.

### Quick Start

```bash
# Start Anvil and credit accounts with 10,000 USDC each
./credit_usdc.sh start

# Check USDC balances
./credit_usdc.sh balance

# Stop Anvil when done
./credit_usdc.sh stop
```

### Available Commands

- `./credit_usdc.sh start` - Start Anvil with Base Sepolia fork and credit accounts
- `./credit_usdc.sh balance` - Check current USDC balances
- `./credit_usdc.sh deploy` - Deploy CreditUSDC contract (if Anvil is already running)
- `./credit_usdc.sh status` - Check if Anvil is running and show balances
- `./credit_usdc.sh stop` - Stop Anvil
- `./credit_usdc.sh help` - Show help message

### What the Script Does

The script automatically:
1. **Starts Anvil** with a Base Sepolia fork at `http://127.0.0.1:8545`
2. **Deploys a credit contract** that transfers USDC from a whale account
3. **Credits 10,000 USDC** to two default Anvil accounts:
   - Account 1: `0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266`
   - Account 2: `0x70997970C51812dc3A010C7d01b50e0d17dc79C8`

### Adding USDC to MetaMask

To see your USDC balance in MetaMask:

1. **Connect MetaMask to Anvil**:
   - Network Name: `Anvil Local`
   - RPC URL: `http://127.0.0.1:8545`
   - Chain ID: `31337`
   - Currency Symbol: `ETH`

2. **Import USDC Token**:
   - Click "Import tokens" in MetaMask
   - **Token Contract Address**: `0x036CbD53842c5426634e7929541eC2318f3dCF7e`
   - **Token Symbol**: `USDC`
   - **Token Decimals**: `6`

3. **Import Anvil Account** (optional):
   - Use the private key: `0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
   - This corresponds to the first Anvil account that gets credited with USDC

### Prerequisites for USDC Script

Make sure you have these tools installed:
- **Foundry** (anvil, forge, cast) - [Install here](https://getfoundry.sh/)
- **bc** - For balance calculations: `sudo apt-get install bc`
- **curl** - Usually pre-installed on most systems

### Troubleshooting

- **Script fails to start**: Make sure Foundry is installed and `anvil` is in your PATH
- **No USDC showing**: Verify you've imported the correct token contract address in MetaMask
- **Connection issues**: Ensure Anvil is running on port 8545 and MetaMask is connected to the right network

## �� Project Structure

```
blockchess/
├── src/                    # Frontend React source
│   ├── components/         # React components
│   ├── services/          # WebSocket service
│   ├── types/             # TypeScript types
│   └── utils/             # Utility functions
├── backend/               # Go backend
│   ├── cmd/server/        # Main server entry point
│   ├── internal/          # Internal packages
│   │   ├── game/         # Game logic
│   │   └── websocket/    # WebSocket handling
│   └── dist/             # Built frontend (production)
├── package.json          # Frontend dependencies
└── vite.config.ts       # Vite configuration
```

## 🎮 How to Play

1. Connect your wallet using WalletConnect
2. Browse available games in the lobby or click "Quick Play"
3. Join a game to enter the chess board
4. Click on the chess board to vote for moves
5. Watch the timer countdown - the most voted move executes when time runs out
6. The game alternates between white and black turns

## 🔧 Configuration

- **Backend port**: Change with `-addr` flag (default: `:8080`)
- **Turn duration**: Modify `TimeLeft` in `backend/internal/game/manager.go` (default: 15 seconds)
- **Game ID**: Currently hardcoded as "game-123" in the frontend

## 🚢 Deployment

The application is designed to run as a single Go binary that serves both the frontend and backend:

1. Build the frontend: `npm run build`
2. Copy dist to backend: `cp -r dist backend/`
3. Build Go binary: `cd backend && go build -o blockchess-server cmd/server/main.go`
4. Deploy the binary and run it on your server

The server will:
- Serve the React frontend on `/`
- Handle WebSocket connections on `/ws`
- Manage game state and broadcasting

## 🤝 Contributing

Feel free to submit issues and pull requests!

## 📄 License

MIT
