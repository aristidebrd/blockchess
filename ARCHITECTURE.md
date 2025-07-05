# BlockChess Architecture - Clean Separation of Concerns

## 🎯 **Architecture Overview**

The BlockChess application now follows a clean architecture with proper separation of concerns across three main layers:

### **1. Application Layer (`App.tsx`)**
**Responsibilities:**
- Wallet connection management
- Games list display & WebSocket subscription  
- Total players count tracking
- Matchmaking logic (start/cancel/callbacks)
- Navigation between lobby/game views
- High-level routing and state management
- **Global header and footer components** (consistent across all pages)

**Key State:**
- `currentView` (lobby | game)
- `games[]` (from WebSocket games list)
- `totalConnections` (platform-wide player count)
- `matchmaking` state (players count, start time, etc.)

**Global UI Components:**
- **Header**: Adapts to context (lobby vs game) with back button functionality
- **Footer**: Consistent across all pages

### **2. Game Engine (`gameService.ts`)**
**Responsibilities:**
- Game state management (timer, turns, voting)
- WebSocket game-specific subscriptions
- Move validation and creation logic
- Player role management within games
- Board state management
- Vote tracking and validation

**Key Features:**
- `initializeGame()` - Sets up game and WebSocket subscriptions
- `voteOnMove()` - Handles voting with validation
- `createMove()` - Validates and creates moves
- `isVotingEnabled()` - Centralized voting logic
- `onGameStateUpdate()` - Subscribe to real-time updates

### **3. Presentation Layer (`GameView.tsx`)**
**Responsibilities:**
- Pure presentation of game state
- User interactions (click handlers)
- UI components (ChessBoard, VotingPanel, Timer)
- Move confirmation dialogs

**Key Features:**
- Consumes game state from gameService
- Handles UI interactions and forwards to gameService
- Clean component with minimal business logic
- **No header/footer** (handled at top level)

---

## 🔄 **Data Flow**

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
```

1. **Backend → GameService**: Real-time game updates via WebSocket
2. **GameService → GameView**: Game state updates via callback subscription  
3. **GameView → GameService**: User actions (votes, moves) via method calls
4. **Backend → App**: Games list updates via WebSocket
5. **App ↔ GameView**: Navigation and high-level state
6. **App**: Provides consistent header/footer across all views

---

## 🏗️ **Key Benefits**

### **1. Separation of Concerns**
- **App.tsx**: Only handles app-level concerns (wallet, games list, navigation, global UI)
- **GameService**: Encapsulates all game logic and state management
- **GameView**: Pure presentation layer with minimal business logic

### **2. Consistent UI**
- **Global header**: Adapts to context (lobby vs game) with proper back navigation
- **Global footer**: Same across all pages for consistent branding
- **Responsive layout**: Flexbox structure ensures proper spacing

### **3. Testability**
- GameService can be unit tested independently
- UI components receive props and call clear methods
- No mixing of UI and business logic

### **4. Maintainability**
- Game logic changes only affect GameService
- UI changes only affect GameView
- App-level features don't interfere with game logic
- Global UI elements managed in one place

### **5. Reusability**
- GameService could be used by different UI frameworks
- GameView could display different game types
- Clear interfaces between layers

---

## 📁 **File Structure**

```
src/
├── App.tsx                 # Application layer (wallet, games list, matchmaking, global UI)
├── components/
│   ├── GameView.tsx        # Game presentation layer (no header/footer)
│   ├── GameLobby.tsx       # Lobby UI
│   ├── ChessBoard.tsx      # Board UI component
│   └── VotingPanel.tsx     # Voting UI component
├── services/
│   ├── gameService.ts      # Game engine (business logic)
│   └── websocket.ts        # WebSocket communication
└── types/
    └── chess.ts            # Type definitions
```

---

## 🔧 **Usage Examples**

### **Global Header Usage (App.tsx)**
```typescript
// Lobby view - larger header with full title
<Header showBackButton={false} />

// Game view - smaller header with back button
<Header showBackButton={true} />

// The header adapts its size and functionality based on context
```

### **Layout Structure**
```typescript
// Consistent layout across all views
<div className="min-h-screen bg-gradient flex flex-col">
  <Header showBackButton={isGameView} />
  <div className="flex-1">
    {/* Main content */}
  </div>
  <Footer />
</div>
```

### **Game Engine Usage (GameView.tsx)**
```typescript
// Initialize game
useEffect(() => {
  gameService.initializeGame(gameId);
  
  const unsubscribe = gameService.onGameStateUpdate((state) => {
    setBackendGameState(state);
    // Update UI based on game state
  });
  
  return () => {
    unsubscribe();
    gameService.cleanupGame();
  };
}, [gameId]);

// Handle voting
const handleVote = async (moveId: string) => {
  const result = await gameService.voteOnMove(moveId, isConnected);
  if (!result.success) {
    alert(result.error);
  }
  return result.success;
};
```

---

## 🚀 **Migration Benefits**

### **Before Refactoring:**
- ❌ App.tsx handled everything (900+ lines)
- ❌ Mixed concerns (UI, game logic, WebSocket, voting)
- ❌ Duplicate headers and inconsistent UI
- ❌ Hard to test and maintain

### **After Refactoring:**
- ✅ Clean separation of concerns
- ✅ App.tsx focused on application concerns (~220 lines)
- ✅ GameService encapsulates all game logic
- ✅ GameView is a pure presentation component
- ✅ **Consistent global header and footer on all pages**
- ✅ Easily testable and maintainable
- ✅ Single source of truth for game state
- ✅ Unified user experience

This architecture makes the codebase much more maintainable, testable, and scalable for future features, with a consistent and professional UI across all pages. 
