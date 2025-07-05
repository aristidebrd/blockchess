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

## 📁 Project Structure

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
- **Turn duration**: Modify `TimeLeft` in `backend/internal/game/manager.go` (default: 30 seconds)
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
