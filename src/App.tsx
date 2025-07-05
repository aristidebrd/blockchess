import { useState, useEffect } from 'react';
import { Crown, Users } from 'lucide-react';
import { useAccount } from 'wagmi';
import WalletConnect from './components/WalletConnect';
import GameLobby from './components/GameLobby';
import GameView from './components/GameView';
import GameNotFound from './components/GameNotFound';
import { GameInfo, ChessPiece } from './utils/chess';
import { wsService } from './services/websocket';
import { gameService } from './services/gameService';

type AppView = 'lobby' | 'game' | 'game-not-found';

function App() {
  const { isConnected, address } = useAccount();

  // Application-level state
  const [currentView, setCurrentView] = useState<AppView>('lobby');
  const [currentGameId, setCurrentGameId] = useState<string | null>(null);
  const [playerRole, setPlayerRole] = useState<'white' | 'black' | 'spectator' | 'none'>('none');
  const [isFromMatchmaking, setIsFromMatchmaking] = useState<boolean>(false);
  const [gameNotFoundId, setGameNotFoundId] = useState<string | null>(null);

  // Games list state
  const [games, setGames] = useState<GameInfo[]>([]);
  const [allGames, setAllGames] = useState<GameInfo[]>([]); // Store complete games list
  const [totalConnections, setTotalConnections] = useState<number>(0);
  const [currentFilter, setCurrentFilter] = useState<'active' | 'ended' | 'all'>('all');

  // Matchmaking state
  const [gamePlayersCount, setGamePlayersCount] = useState<number>(0);
  const [matchmakingStartTime, setMatchmakingStartTime] = useState<Date | null>(null);
  const [isMatchmaking, setIsMatchmaking] = useState<boolean>(false);

  // URL handling for deep linking to games
  useEffect(() => {
    const handleURLChange = () => {
      const urlParams = new URLSearchParams(window.location.search);
      const gameIdFromURL = urlParams.get('game_id');

      if (gameIdFromURL) {
        console.log('Found game_id in URL:', gameIdFromURL);
        setCurrentGameId(gameIdFromURL);
        setCurrentView('game');
        setPlayerRole('none'); // Let user choose team unless coming from matchmaking
        setIsFromMatchmaking(false);
      } else if (window.location.pathname === '/' || window.location.pathname === '') {
        console.log('Navigating to lobby');
        setCurrentView('lobby');
        setCurrentGameId(null);
        setPlayerRole('none');
        setIsFromMatchmaking(false);
      }
    };

    // Handle initial URL
    handleURLChange();

    // Listen for URL changes (back/forward navigation)
    window.addEventListener('popstate', handleURLChange);

    return () => {
      window.removeEventListener('popstate', handleURLChange);
    };
  }, []);

  // Helper function to update URL
  const updateURL = (gameId: string | null) => {
    if (gameId) {
      const newURL = `${window.location.origin}/?game_id=${gameId}`;
      window.history.pushState({ gameId }, '', newURL);
    } else {
      const newURL = window.location.origin;
      window.history.pushState({}, '', newURL);
    }
  };

  // Helper function to convert WebSocket GameInfo to local GameInfo format
  // Filter games based on status
  const filterGames = (gamesList: GameInfo[], filter: 'active' | 'ended' | 'all'): GameInfo[] => {
    if (filter === 'all') {
      return gamesList;
    } else if (filter === 'active') {
      return gamesList.filter(game => game.status === 'active' || !game.status);
    } else if (filter === 'ended') {
      return gamesList.filter(game => game.status === 'ended');
    }
    return gamesList;
  };

  const convertWSGameToLocalFormat = (wsGame: any): GameInfo => {

    // Convert backend board format to frontend format
    const convertBackendBoardToFrontend = (board: string[][]): (ChessPiece | null)[][] => {
      const pieceMap: Record<string, { type: ChessPiece['type'], color: ChessPiece['color'] }> = {
        'P': { type: 'pawn', color: 'white' },
        'R': { type: 'rook', color: 'white' },
        'N': { type: 'knight', color: 'white' },
        'B': { type: 'bishop', color: 'white' },
        'Q': { type: 'queen', color: 'white' },
        'K': { type: 'king', color: 'white' },
        'p': { type: 'pawn', color: 'black' },
        'r': { type: 'rook', color: 'black' },
        'n': { type: 'knight', color: 'black' },
        'b': { type: 'bishop', color: 'black' },
        'q': { type: 'queen', color: 'black' },
        'k': { type: 'king', color: 'black' },
      };

      return board.map((row, rowIndex) =>
        row.map((cell, colIndex) => {
          if (!cell || cell === '') {
            return null;
          }

          const pieceInfo = pieceMap[cell];
          if (!pieceInfo) {
            return null;
          }

          // Convert board coordinates to chess notation
          const file = String.fromCharCode(97 + colIndex); // a, b, c, ...
          const rank = (8 - rowIndex).toString(); // 8, 7, 6, ...

          return {
            type: pieceInfo.type,
            color: pieceInfo.color,
            position: file + rank,
            hasMoved: false
          } as ChessPiece;
        })
      );
    };

    return {
      id: wsGame.gameId,
      name: `Game ${wsGame.gameId}`,
      status: wsGame.status || (wsGame.currentMove === 0 ? 'waiting' : 'active'),
      currentMove: wsGame.currentMove,
      currentTurn: wsGame.currentTurn as 'white' | 'black',
      turnStartTime: Date.now() - ((10 - wsGame.timeLeft) * 1000),
      turnTimeLimit: 10000,
      timeRemaining: wsGame.timeLeft * 1000,
      spectators: wsGame.spectators,
      createdAt: new Date(), // Backend doesn't provide this yet
      whitePlayers: wsGame.whitePlayers,
      blackPlayers: wsGame.blackPlayers,
      proposedMoves: [],
      whiteCurrentTurnVotes: wsGame.whiteCurrentTurnVotes || 0,
      blackCurrentTurnVotes: wsGame.blackCurrentTurnVotes || 0,
      whiteTeamTotalVotes: wsGame.whiteTeamTotalVotes || 0,
      blackTeamTotalVotes: wsGame.blackTeamTotalVotes || 0,
      playerVotedThisRound: wsGame.playerVotedThisRound || {},
      playerTotalVotes: wsGame.playerTotalVotes || {},
      totalPot: wsGame.totalPot,
      whitePot: wsGame.whitePot,
      blackPot: wsGame.blackPot,
      boardState: (() => {
        const backendBoard = wsGame.board || wsGame.Board;
        if (backendBoard) {
          const converted = convertBackendBoardToFrontend(backendBoard);
          return converted;
        }
        return undefined;
      })(),
      winner: wsGame.winner,
      endReason: wsGame.endReason,
      endedAt: wsGame.endedAt,
      // Add player statistics for ended games
      whiteTeamPlayers: wsGame.whiteTeamPlayers || undefined,
      blackTeamPlayers: wsGame.blackTeamPlayers || undefined,
    };
  };

  // Initialize game service when wallet connects
  useEffect(() => {
    if (isConnected) {
      console.log('Starting game service heartbeat');
      gameService.startHeartbeat();
    } else {
      console.log('Stopping game service - wallet disconnected');
      gameService.stop();
    }

    return () => {
      if (process.env.NODE_ENV === 'production') {
        gameService.stop();
      }
    };
  }, [isConnected]);

  // Update game service with wallet address when it changes
  useEffect(() => {
    if (isConnected && address) {
      gameService.setWalletAddress(address);
    } else {
      gameService.setWalletAddress(null);
    }
  }, [isConnected, address]);

  // Set up game not found callback
  useEffect(() => {
    const unsubscribe = gameService.onGameNotFound((gameId: string) => {
      console.log('Game not found:', gameId);
      setGameNotFoundId(gameId);
      setCurrentView('game-not-found');
      setCurrentGameId(null);
      setPlayerRole('none');
      setIsFromMatchmaking(false);
    });

    return unsubscribe;
  }, []);

  // Function to request games list based on current filter
  const requestGamesList = () => {
    // Always request all games, filtering happens on frontend
    wsService.requestGamesList();
  };

  // Function to handle filter changes
  const handleFilterChange = (filter: 'active' | 'ended' | 'all') => {
    console.log('Filter changed to:', filter);
    setCurrentFilter(filter);
    // Filter will be applied automatically by useEffect
  };

  // Fetch games list via WebSocket
  useEffect(() => {
    console.log('Setting up games list WebSocket subscriptions...');
    requestGamesList();

    // Listen for games list responses (initial request)
    const unsubscribeGamesList = wsService.on('games_list', (data) => {
      console.log('Games list received:', data);
      const convertedGames = data.gamesList.map(convertWSGameToLocalFormat);
      setAllGames(convertedGames); // Store complete list, filtering handled by useEffect
    });

    const unsubscribeNumberOfPlayers = wsService.on('number_of_players', (data) => {
      console.log('Number of players received:', data);
      setTotalConnections(data.totalConnections);
    });

    // Listen for games list updates (real-time changes)
    const unsubscribeGamesUpdate = wsService.on('games_list_update', (data) => {
      console.log('Games list update received:', data);

      const convertedGames = data.gamesList.map(convertWSGameToLocalFormat);
      setAllGames(convertedGames); // Always update complete list, filtering handled by useEffect
    });

    // Cleanup subscriptions
    return () => {
      unsubscribeGamesList();
      unsubscribeGamesUpdate();
      unsubscribeNumberOfPlayers();
    };
  }, [currentFilter]);

  // Apply filter whenever allGames or currentFilter changes
  useEffect(() => {
    const filteredGames = filterGames(allGames, currentFilter);
    setGames(filteredGames);
  }, [allGames, currentFilter]);

  // Check player role when in a game
  useEffect(() => {
    if (currentGameId && currentView === 'game') {
      const checkPlayerRole = async () => {
        console.log('Checking player role for game:', currentGameId);
        const roleInfo = await gameService.getPlayerRole(currentGameId);
        console.log('Role info received:', roleInfo);
        if (roleInfo) {
          console.log('Setting player role to:', roleInfo.role);
          setPlayerRole(roleInfo.role);
        }
      };

      checkPlayerRole();
      // Check role every 5 seconds
      const roleInterval = setInterval(checkPlayerRole, 5000);

      return () => clearInterval(roleInterval);
    } else {
      setPlayerRole('none');
    }
  }, [currentGameId, currentView]);

  // ========== MATCHMAKING HANDLERS ==========

  const handleStartMatchmaking = async () => {
    if (!isConnected) {
      alert('Please connect your wallet to start matchmaking');
      return;
    }

    setIsMatchmaking(true);
    setMatchmakingStartTime(new Date());
    setGamePlayersCount(1);

    try {
      const success = await gameService.joinMatchmaking((players: string[], gameId: string, assignedSide?: 'white' | 'black') => {
        console.log('🔥 MATCH FOUND CALLBACK TRIGGERED!');
        console.log('🔥 Players:', players);
        console.log('🔥 Game ID:', gameId);
        console.log('🔥 Assigned Side:', assignedSide);

        // Reset matchmaking state immediately when match is found
        console.log('🔥 Resetting matchmaking state...');
        setIsMatchmaking(false);
        setMatchmakingStartTime(null);
        setGamePlayersCount(2);
        setCurrentGameId(gameId);

        if (assignedSide) {
          console.log('🔥 Calling handleJoinGame with side:', assignedSide);
          handleJoinGame(gameId, assignedSide);
        } else {
          console.log('🔥 No assigned side, calling handleJoinGame without side');
          handleJoinGame(gameId);
        }
      });

      if (!success) {
        alert('Failed to join matchmaking. Please try again.');
        setIsMatchmaking(false);
        setMatchmakingStartTime(null);
        setGamePlayersCount(0);
      }
    } catch (error) {
      console.error('Matchmaking error:', error);
      alert('Error starting matchmaking. Please try again.');
      setIsMatchmaking(false);
      setMatchmakingStartTime(null);
      setGamePlayersCount(0);
    }
  };

  const handleCancelMatchmaking = async () => {
    setIsMatchmaking(false);
    setGamePlayersCount(0);
    setMatchmakingStartTime(null);

    try {
      await gameService.leaveMatchmaking();
    } catch (error) {
      console.error('Error leaving matchmaking:', error);
    }
  };

  // ========== GAME NAVIGATION HANDLERS ==========

  const handleJoinGame = async (gameId: string, side?: 'white' | 'black' | 'spectator') => {
    console.log(`App: Joining game ${gameId} as ${side || 'viewer'}`);

    setCurrentGameId(gameId);
    setCurrentView('game');
    updateURL(gameId); // Update URL to include game_id

    // Reset matchmaking state
    setGamePlayersCount(0);
    setMatchmakingStartTime(null);

    const fromMatchmaking = Boolean(side && (side === 'white' || side === 'black'));
    setIsFromMatchmaking(fromMatchmaking);

    // Join the appropriate team or watch as spectator
    if (side === 'white') {
      console.log('Joining white team');
      await gameService.joinTeam(gameId, 'white');
      setPlayerRole('white');
    } else if (side === 'black') {
      console.log('Joining black team');
      await gameService.joinTeam(gameId, 'black');
      setPlayerRole('black');
    } else if (side === 'spectator') {
      console.log('Joining as spectator');
      await gameService.watchGame(gameId);
      setPlayerRole('spectator');
    } else {
      console.log('No side provided, setting role to none for manual team selection');
      setPlayerRole('none');
      setIsFromMatchmaking(false);
    }
  };

  const handleJoinTeam = async (side: 'white' | 'black') => {
    if (!isConnected || !address) {
      alert('Please connect your wallet to join a team');
      return;
    }

    // If player is already on the requested team (e.g., from matchmaking), don't send another join message
    if (playerRole === side) {
      console.log(`App: Player already on ${side} team, no need to join again`);
      return;
    }

    // If player is on a different team, they shouldn't be able to switch
    if (playerRole !== 'none' && playerRole !== 'spectator') {
      console.log(`App: Player already on ${playerRole} team, cannot switch to ${side}`);
      alert(`You are already on the ${playerRole} team and cannot switch sides.`);
      return;
    }

    if (currentGameId) {
      console.log(`App: Joining ${side} team in game ${currentGameId}`);

      try {
        const result = await gameService.joinTeam(currentGameId, side);

        if (result.success) {
          setPlayerRole(side);
        } else {
          alert(result.error || 'Failed to join team');
        }
      } catch (error) {
        console.error('Error joining team:', error);
        alert('Error joining team. Please try again.');
      }
    }
  };

  const handleBackToLobby = () => {
    console.log('App: Returning to lobby');
    setCurrentView('lobby');
    setCurrentGameId(null);
    setIsFromMatchmaking(false);
    setPlayerRole('none');
    setGameNotFoundId(null);
    updateURL(null); // Update URL to remove game_id
  };

  // ========== RENDER ==========

  // Create a Header component that will be used across all pages
  const Header = ({ showBackButton = false }: { showBackButton?: boolean }) => (
    <header className="border-b border-gray-800 bg-black/20 backdrop-blur-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            {showBackButton ? (
              <button
                onClick={handleBackToLobby}
                className="bg-gradient-to-r from-yellow-400 to-yellow-600 p-2 rounded-lg hover:scale-105 transition-transform"
              >
                <Crown className="w-8 h-8 text-black" />
              </button>
            ) : (
              <div className="bg-gradient-to-r from-yellow-400 to-yellow-600 p-3 rounded-xl">
                <Crown className="w-10 h-10 text-black" />
              </div>
            )}
            <div>
              <h1 className={`font-bold bg-gradient-to-r from-yellow-400 to-yellow-600 bg-clip-text text-transparent ${showBackButton ? 'text-2xl' : 'text-3xl'
                }`}>
                BlockChess
              </h1>
              <p className="text-gray-400 text-sm">
                {showBackButton ? 'Community-driven chess gaming' : 'Join the community-driven chess revolution'}
              </p>
            </div>
          </div>

          <div className="flex items-center space-x-6">
            <div className="flex items-center space-x-2 text-gray-300">
              <Users className="w-5 h-5" />
              <span className="text-sm">{totalConnections} players online</span>
            </div>

            <WalletConnect isCompact={true} />
          </div>
        </div>
      </div>
    </header>
  );

  // Create a Footer component that will be used across all pages
  const Footer = () => (
    <footer className="border-t border-gray-800 bg-black/20 backdrop-blur-sm mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div className="text-center text-gray-400">
          <p className="text-sm">
            Powered by blockchain technology • Fair, transparent, and decentralized
          </p>
        </div>
      </div>
    </footer>
  );

  if (currentView === 'lobby') {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex flex-col">
        <Header showBackButton={false} />

        <div className="flex-1">
          <GameLobby
            onJoinGame={handleJoinGame}
            onStartMatchmaking={handleStartMatchmaking}
            onCancelMatchmaking={handleCancelMatchmaking}
            gamePlayersCount={gamePlayersCount}
            matchmakingStartTime={matchmakingStartTime}
            isMatchmaking={isMatchmaking}
            games={games}
            totalConnections={totalConnections}
            onFilterChange={handleFilterChange}
          />
        </div>

        <Footer />
      </div>
    );
  }

  // Game Not Found View
  if (currentView === 'game-not-found' && gameNotFoundId) {
    return (
      <GameNotFound
        gameId={gameNotFoundId}
        onBackToLobby={handleBackToLobby}
      />
    );
  }

  // Game View
  if (currentGameId) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex flex-col">
        <Header showBackButton={true} />

        <div className="flex-1">
          <GameView
            gameId={currentGameId}
            playerRole={playerRole}
            onJoinTeam={handleJoinTeam}
            isFromMatchmaking={isFromMatchmaking}
          />
        </div>
        <Footer />
      </div>
    );
  }

  // Loading state
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex flex-col">
      <Header showBackButton={false} />

      <div className="flex-1 flex items-center justify-center">
        <div className="text-center text-white">
          <div className="text-xl mb-4">Loading...</div>
          <div className="w-8 h-8 border-2 border-white border-t-transparent rounded-full animate-spin mx-auto"></div>
        </div>
      </div>
      <Footer />
    </div>
  );
}

export default App;
