import React, { useState, useEffect } from 'react';
import { Play, Users, Clock, Crown, Eye, AlertTriangle, Swords, Trophy, Timer, Copy, ExternalLink, Filter, Vote, Coins } from 'lucide-react';
import { useAccount } from 'wagmi';
import { GameInfo, ChessPiece } from '../utils/chess';
import { wsService } from '../services/websocket';
import PlayerStatistics from './PlayerStatistics';

interface GameLobbyProps {
  onJoinGame: (gameId: string, side?: 'white' | 'black' | 'spectator') => void;
  onStartMatchmaking: () => void;
  onCancelMatchmaking: () => void;
  gamePlayersCount: number;
  matchmakingStartTime: Date | null;
  isMatchmaking: boolean;
  games: GameInfo[];
  totalConnections: number;
  onFilterChange?: (filter: 'active' | 'ended' | 'all') => void;
}

// Compact board preview component
const BoardPreview: React.FC<{ boardState?: (ChessPiece | null)[][], onClick?: () => void, isLarge?: boolean }> = ({ boardState, onClick, isLarge = false }) => {
  if (!boardState) {
    // Show initial board setup if no board state available
    return (
      <div
        className={`grid grid-cols-8 gap-px bg-amber-900 p-1 rounded-lg border-2 border-amber-700 ${onClick ? 'cursor-pointer hover:border-amber-500 transition-colors' : ''} ${isLarge ? 'w-96 h-96' : 'w-32 h-32'}`}
        onClick={onClick}
      >
        {Array.from({ length: 64 }).map((_, index) => {
          const row = Math.floor(index / 8);
          const col = index % 8;
          const isLight = (row + col) % 2 === 0;

          return (
            <div
              key={index}
              className={`aspect-square ${isLight ? 'bg-amber-100' : 'bg-amber-800'} 
                         rounded-sm flex items-center justify-center ${isLarge ? 'text-4xl' : 'text-xs'}`}
            />
          );
        })}
      </div>
    );
  }

  const getPieceSymbol = (piece: ChessPiece) => {
    const symbols = {
      king: piece.color === 'white' ? '♔' : '♚',
      queen: piece.color === 'white' ? '♕' : '♛',
      rook: piece.color === 'white' ? '♖' : '♜',
      bishop: piece.color === 'white' ? '♗' : '♝',
      knight: piece.color === 'white' ? '♘' : '♞',
      pawn: piece.color === 'white' ? '♙' : '♟︎',
    };
    return symbols[piece.type];
  };

  return (
    <div
      className={`grid grid-cols-8 gap-px bg-amber-900 p-1 rounded-lg border-2 border-amber-700 ${onClick ? 'cursor-pointer hover:border-amber-500 transition-colors' : ''} ${isLarge ? 'w-96 h-96' : 'w-32 h-32'}`}
      onClick={onClick}
    >
      {boardState.flat().map((piece, index) => {
        const row = Math.floor(index / 8);
        const col = index % 8;
        const isLight = (row + col) % 2 === 0;

        return (
          <div
            key={index}
            className={`aspect-square ${isLight ? 'bg-amber-100' : 'bg-amber-800'} 
                       rounded-sm flex items-center justify-center ${isLarge ? 'text-4xl' : 'text-xs'} font-bold`}
          >
            {piece && (
              <span className={`drop-shadow-sm ${piece.color === 'black' ? 'text-gray-800' : 'text-gray-900'}`}>
                {getPieceSymbol(piece)}
              </span>
            )}
          </div>
        );
      })}
    </div>
  );
};

const GameLobby: React.FC<GameLobbyProps> = ({ onJoinGame, onStartMatchmaking, onCancelMatchmaking, gamePlayersCount, matchmakingStartTime, isMatchmaking, games, onFilterChange }) => {
  const { isConnected, address } = useAccount();
  const [selectedGameInfo, setSelectedGameInfo] = useState<GameInfo | null>(null);
  const [selectedStatsGame, setSelectedStatsGame] = useState<GameInfo | null>(null);
  const [matchmakingTime, setMatchmakingTime] = useState(0);
  const [copiedGameId, setCopiedGameId] = useState<string | null>(null);
  const [currentFilter, setCurrentFilter] = useState<'active' | 'ended' | 'all'>('all');

  // Track matchmaking time
  useEffect(() => {
    if (matchmakingStartTime) {
      const interval = setInterval(() => {
        const elapsed = Math.floor((Date.now() - matchmakingStartTime.getTime()) / 1000);
        setMatchmakingTime(elapsed);
      }, 1000);

      return () => clearInterval(interval);
    } else {
      setMatchmakingTime(0);
    }
  }, [matchmakingStartTime]);



  const handleMatchmaking = () => {
    if (!isConnected) {
      alert('Please connect your wallet to start matchmaking');
      return;
    }
    onStartMatchmaking();
  };

  const handleFilterChange = (filter: 'active' | 'ended' | 'all') => {
    setCurrentFilter(filter);
    if (onFilterChange) {
      onFilterChange(filter);
    }
  };

  const canJoinAsPlayer = (game: GameInfo, playerAddress?: string): boolean => {
    if (game.status !== 'waiting' && game.gameStartTime) {
      const gameAge = Date.now() - game.gameStartTime.getTime();
      return gameAge <= 7 * 60 * 1000; // 7 minutes in milliseconds
    }
    return game.status === 'waiting';
  };

  const [playerStatusCache, setPlayerStatusCache] = useState<Record<string, { team: string; gameId: string }>>({});

  const checkPlayerStatus = (gameId: string, walletAddress: string) => {
    return new Promise<{ team: string; gameId: string }>((resolve) => {
      const cacheKey = `${gameId}-${walletAddress}`;
      console.log(`🔍 checkPlayerStatus called for ${walletAddress} in game ${gameId}`);
      console.log(`🔍 Cache key: ${cacheKey}`);
      console.log(`🔍 Current cache:`, playerStatusCache);

      // Check cache first
      if (playerStatusCache[cacheKey]) {
        console.log(`🔍 Found in cache:`, playerStatusCache[cacheKey]);
        resolve(playerStatusCache[cacheKey]);
        return;
      }

      console.log(`🔍 Not in cache, sending WebSocket request`);

      // Set up listener for the response
      const unsubscribe = wsService.on('player_status', (data: any) => {
        console.log(`🔍 Received player_status message:`, data);
        if (data.gameId === gameId && data.walletAddress === walletAddress) {
          const statusData = { team: data.team, gameId: data.gameId };
          console.log(`🔍 Status data matches our request:`, statusData);

          // Cache the result
          setPlayerStatusCache(prev => ({
            ...prev,
            [cacheKey]: statusData
          }));

          unsubscribe();
          resolve(statusData);
        }
      });

      // Send the request
      console.log(`🔍 Sending checkPlayerStatus request via WebSocket`);
      wsService.checkPlayerStatus(gameId, walletAddress);

      // Add timeout to prevent hanging
      setTimeout(() => {
        console.log(`🔍 Timeout waiting for player status response`);
        unsubscribe();
        resolve({ team: '', gameId });
      }, 5000);
    });
  };

  // Check player status for all games when connected wallet changes
  useEffect(() => {
    if (isConnected && address && games.length > 0) {
      // Check player status for all active games
      games.forEach(game => {
        if (game.status === 'active' || game.status === 'waiting') {
          const cacheKey = `${game.id}-${address}`;
          // Only check if not already cached
          if (!playerStatusCache[cacheKey]) {
            checkPlayerStatus(game.id, address);
          }
        }
      });
    }
  }, [isConnected, address, games, playerStatusCache, checkPlayerStatus]);

  const getGameAgeMinutes = (game: GameInfo): number => {
    if (!game.gameStartTime) return 0;
    return Math.floor((Date.now() - game.gameStartTime.getTime()) / 60000);
  };

  const formatTimeAgo = (date: Date) => {
    const gameDate = new Date(date);
    const minutes = Math.floor((Date.now() - gameDate.getTime()) / 60000);
    if (minutes < 1) return 'Just now';
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    return `${hours}h ago`;
  };

  const getGameStatusBadge = (game: GameInfo) => {
    const gameAge = getGameAgeMinutes(game);
    const isTooOld = gameAge > 7 && game.status === 'active';

    if (game.status === 'waiting') {
      return (
        <div className="flex items-center space-x-1 bg-yellow-500/20 text-yellow-400 px-2 py-1 rounded-full text-xs font-medium">
          <Clock className="w-3 h-3" />
          <span>Waiting</span>
        </div>
      );
    }

    if (game.status === 'active' && !isTooOld) {
      return (
        <div className="flex items-center space-x-1 bg-green-500/20 text-green-400 px-2 py-1 rounded-full text-xs font-medium">
          <Swords className="w-3 h-3" />
          <span>Active</span>
        </div>
      );
    }

    if (isTooOld) {
      return (
        <div className="flex items-center space-x-1 bg-orange-500/20 text-orange-400 px-2 py-1 rounded-full text-xs font-medium">
          <AlertTriangle className="w-3 h-3" />
          <span>Late Join</span>
        </div>
      );
    }

    if (game.status === 'completed' || game.status === 'ended') {
      return (
        <div className="flex items-center space-x-1 bg-purple-500/20 text-purple-400 px-2 py-1 rounded-full text-xs font-medium">
          <Trophy className="w-3 h-3" />
          <span>Finished</span>
        </div>
      );
    }

    return null;
  };

  const getWinnerBadge = (game: GameInfo) => {
    if (game.status !== 'ended' || !game.winner) return null;

    const winnerColor = game.winner === 'white' ? 'bg-gray-100 text-gray-900' :
      game.winner === 'black' ? 'bg-gray-800 text-white' :
        'bg-yellow-500 text-yellow-900';

    return (
      <div className={`flex items-center space-x-1 px-2 py-1 rounded-full text-xs font-medium ${winnerColor}`}>
        <Crown className="w-3 h-3" />
        <span>{game.winner === 'draw' ? 'Draw' : `${game.winner} wins`}</span>
      </div>
    );
  };

  const openBoardModal = (game: GameInfo) => {
    setSelectedGameInfo(game);
  };

  const closeBoardModal = () => {
    setSelectedGameInfo(null);
  };

  const openStatsModal = (game: GameInfo) => {
    setSelectedStatsGame(game);
  };

  const closeStatsModal = () => {
    setSelectedStatsGame(null);
  };

  // Helper function to copy game URL to clipboard
  const copyGameURL = async (gameId: string) => {
    const gameURL = `${window.location.origin}/?game_id=${gameId}`;
    try {
      await navigator.clipboard.writeText(gameURL);
      setCopiedGameId(gameId);
      setTimeout(() => setCopiedGameId(null), 2000); // Reset after 2 seconds
    } catch (err) {
      console.error('Failed to copy URL:', err);
      // Fallback for older browsers
      const textArea = document.createElement('textarea');
      textArea.value = gameURL;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
      setCopiedGameId(gameId);
      setTimeout(() => setCopiedGameId(null), 2000);
    }
  };

  const handleJoinTeam = async (gameId: string, side: 'white' | 'black') => {
    console.log(`🔘 handleJoinTeam called with gameId: ${gameId}, side: ${side}`);
    console.log(`🔘 isConnected: ${isConnected}, address: ${address}`);

    if (!isConnected || !address) {
      alert('Please connect your wallet to join a team');
      return;
    }

    try {
      console.log(`🔘 Checking player status for ${address} in game ${gameId}`);
      // Check if player is already in the game
      const playerStatus = await checkPlayerStatus(gameId, address);
      console.log(`🔘 Player status result:`, playerStatus);

      if (playerStatus.team && playerStatus.team !== '') {
        if (playerStatus.team === side) {
          // Player is reconnecting to their existing team
          console.log(`🔘 Player ${address} reconnecting to ${playerStatus.team} team in game ${gameId}`);
          onJoinGame(gameId, side);
        } else {
          // Player is trying to join a different team
          console.log(`🔘 Player trying to switch teams from ${playerStatus.team} to ${side}`);
          alert(`You are already on the ${playerStatus.team} team. You cannot switch teams.`);
        }
      } else {
        // Player is not in the game yet, proceed with joining
        console.log(`🔘 Player not in game yet, proceeding with join`);
        onJoinGame(gameId, side);
      }
    } catch (error) {
      console.error('🔘 Error checking player status:', error);
      // Fallback to normal join if status check fails
      console.log(`🔘 Fallback to normal join`);
      onJoinGame(gameId, side);
    }
  };

  const handleWatchGame = (gameId: string) => {
    // Spectating doesn't require wallet connection
    onJoinGame(gameId, 'spectator');
  };

  const handleJoinTeamFromModal = async (gameId: string, side: 'white' | 'black') => {
    if (!isConnected || !address) {
      alert('Please connect your wallet to join a team');
      return;
    }

    try {
      // Check if player is already in the game
      const playerStatus = await checkPlayerStatus(gameId, address);

      if (playerStatus.team && playerStatus.team !== '') {
        if (playerStatus.team === side) {
          // Player is reconnecting to their existing team
          console.log(`Player ${address} reconnecting to ${playerStatus.team} team in game ${gameId}`);
          onJoinGame(gameId, side);
          closeBoardModal();
        } else {
          // Player is trying to join a different team
          alert(`You are already on the ${playerStatus.team} team. You cannot switch teams.`);
        }
      } else {
        // Player is not in the game yet, proceed with joining
        onJoinGame(gameId, side);
        closeBoardModal();
      }
    } catch (error) {
      console.error('Error checking player status:', error);
      // Fallback to normal join if status check fails
      onJoinGame(gameId, side);
      closeBoardModal();
    }
  };

  const handleWatchGameFromModal = (gameId: string) => {
    // Spectating doesn't require wallet connection
    onJoinGame(gameId, 'spectator');
    closeBoardModal();
  };

  return (
    <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">

      {/* Quick Actions */}
      <div className="mb-8">
        <div className="bg-gradient-to-r from-blue-600 to-purple-600 rounded-2xl p-8 text-white shadow-2xl">
          <div className="flex flex-col md:flex-row items-center justify-between">
            <div className="mb-6 md:mb-0">
              <h2 className="text-3xl font-bold mb-3">Ready to Play?</h2>
              <p className="text-blue-100 text-lg">Join a game or start matchmaking to find an opponent</p>
            </div>
            <div className="flex flex-col space-y-3">
              <button
                onClick={handleMatchmaking}
                disabled={isMatchmaking}
                className="bg-white text-blue-600 hover:bg-gray-100 font-bold py-4 px-8 rounded-xl transition-all duration-200 flex items-center space-x-3 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg hover:shadow-xl transform hover:scale-105"
              >
                {isMatchmaking ? (
                  <>
                    <div className="w-6 h-6 border-2 border-blue-600 border-t-transparent rounded-full animate-spin" />
                    <div className="text-left">
                      <div>Finding Opponent...</div>
                      <div className="text-sm font-normal">
                        {gamePlayersCount}/2 players • {Math.floor(matchmakingTime / 60)}:{(matchmakingTime % 60).toString().padStart(2, '0')}
                      </div>
                    </div>
                  </>
                ) : (
                  <>
                    <Play className="w-6 h-6" />
                    <span>Quick Match</span>
                  </>
                )}
              </button>

              {/* Cancel matchmaking button */}
              {isMatchmaking && (
                <button
                  onClick={onCancelMatchmaking}
                  className="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded-lg transition-all duration-200 shadow-lg"
                >
                  Cancel Matchmaking
                </button>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Filter Section */}
      <div className="mb-6">
        <div className="bg-gradient-to-r from-gray-800 to-gray-900 rounded-xl p-4 border border-gray-700">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <Filter className="w-5 h-5 text-gray-400" />
              <span className="text-white font-medium">Filter Games</span>
            </div>
            <div className="flex space-x-2">
              <button
                onClick={() => handleFilterChange('all')}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${currentFilter === 'all'
                  ? 'bg-blue-600 text-white shadow-lg'
                  : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                  }`}
              >
                All Games
              </button>
              <button
                onClick={() => handleFilterChange('active')}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${currentFilter === 'active'
                  ? 'bg-green-600 text-white shadow-lg'
                  : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                  }`}
              >
                Active Games
              </button>
              <button
                onClick={() => handleFilterChange('ended')}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${currentFilter === 'ended'
                  ? 'bg-purple-600 text-white shadow-lg'
                  : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                  }`}
              >
                Ended Games
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Games Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
        {games.map((game) => {
          const canJoin = canJoinAsPlayer(game, address);
          const gameAge = getGameAgeMinutes(game);
          const isTooOld = gameAge > 7 && game.status === 'active';

          // Get cached player status if available
          const playerStatusKey = address ? `${game.id}-${address}` : null;
          const playerStatus = playerStatusKey ? playerStatusCache[playerStatusKey] : null;

          return (
            <div
              key={game.id}
              className="bg-gradient-to-br from-gray-800 via-gray-900 to-black rounded-2xl border border-gray-700 hover:border-gray-600 transition-all duration-300 shadow-xl hover:shadow-2xl transform hover:scale-[1.02] overflow-hidden"
            >
              {/* Header with Status and Time */}
              <div className="p-4 border-b border-gray-700/50 bg-gradient-to-r from-gray-800/50 to-gray-900/50">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    {getGameStatusBadge(game)}
                    {getWinnerBadge(game)}
                    <div className="flex items-center space-x-1 text-gray-400 text-sm">
                      <Timer className="w-3 h-3" />
                      <span>Move {game.currentMove}</span>
                    </div>
                  </div>
                  <div className="flex items-center space-x-3">
                    <button
                      onClick={() => copyGameURL(game.id)}
                      className="flex items-center space-x-1 text-gray-400 hover:text-white text-sm transition-colors duration-200 hover:bg-gray-700 px-2 py-1 rounded"
                      title="Copy game URL"
                    >
                      {copiedGameId === game.id ? (
                        <>
                          <ExternalLink className="w-3 h-3" />
                          <span>Copied!</span>
                        </>
                      ) : (
                        <>
                          <Copy className="w-3 h-3" />
                          <span>Share</span>
                        </>
                      )}
                    </button>
                    <div className="text-gray-400 text-sm">
                      {formatTimeAgo(game.createdAt || new Date())}
                    </div>
                  </div>
                </div>
              </div>

              {/* Main Content */}
              <div className="p-6">
                {/* Board and Game Info Grid */}
                <div className="grid grid-cols-3 gap-4 mb-6">
                  {/* Chess Board Preview */}
                  <div className="col-span-1 flex justify-center">
                    <BoardPreview
                      boardState={game.boardState}
                      onClick={() => openBoardModal(game)}
                    />
                  </div>

                  {/* Teams and Players */}
                  <div className="col-span-2 space-y-4">
                    {/* White Team */}
                    <div className="bg-gradient-to-r from-gray-700/30 to-gray-600/30 rounded-lg p-3 border border-gray-600/50">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-2">
                          <div className="w-4 h-4 bg-white rounded-full shadow-sm" />
                          <span className="text-white font-medium text-sm">White Team</span>
                        </div>
                        <div className="flex items-center space-x-2">
                          <Users className="w-4 h-4 text-gray-400" />
                          <span className="text-white font-bold">{game.whitePlayers}</span>
                        </div>
                      </div>
                      <div className="mt-2 text-xs text-gray-400">
                        {game.whitePot.toFixed(2)} USDC contributed
                      </div>
                    </div>

                    {/* Black Team */}
                    <div className="bg-gradient-to-r from-gray-800/30 to-gray-900/30 rounded-lg p-3 border border-gray-700/50">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-2">
                          <div className="w-4 h-4 bg-gray-600 rounded-full shadow-sm" />
                          <span className="text-gray-300 font-medium text-sm">Black Team</span>
                        </div>
                        <div className="flex items-center space-x-2">
                          <Users className="w-4 h-4 text-gray-400" />
                          <span className="text-gray-300 font-bold">{game.blackPlayers}</span>
                        </div>
                      </div>
                      <div className="mt-2 text-xs text-gray-400">
                        {game.blackPot.toFixed(2)} USDC contributed
                      </div>
                    </div>
                  </div>
                </div>

                {/* Game Stats - Different for active vs ended games */}
                {game.status === 'active' ? (
                  <div className="grid grid-cols-3 gap-3 mb-6">
                    <div className="bg-gradient-to-br from-green-900/30 to-green-800/30 rounded-lg p-3 text-center border border-green-500/20">
                      <div className="text-green-400 font-bold text-lg">{game.totalPot.toFixed(2)}</div>
                      <div className="text-green-300 text-xs">Total Pot (USDC)</div>
                    </div>
                    <div className="bg-gradient-to-br from-blue-900/30 to-blue-800/30 rounded-lg p-3 text-center border border-blue-500/20">
                      <div className="text-blue-400 font-bold text-lg">{game.spectators}</div>
                      <div className="text-blue-300 text-xs">Spectators</div>
                    </div>
                    <div className="bg-gradient-to-br from-purple-900/30 to-purple-800/30 rounded-lg p-3 text-center border border-purple-500/20">
                      <div className="text-purple-400 font-bold text-lg capitalize">{game.currentTurn}</div>
                      <div className="text-purple-300 text-xs">Turn</div>
                    </div>
                  </div>
                ) : (
                  <div></div>
                )}

                {/* Action Buttons */}
                <div className="space-y-3">
                  {game.status === 'active' && !isTooOld && (
                    <div className="grid grid-cols-3 gap-2">
                      <button
                        onClick={() => {
                          console.log('🔘 DEBUG: White button clicked!');
                          handleJoinTeam(game.id, 'white');
                        }}
                        disabled={!isConnected || (playerStatus != null && playerStatus.team !== '' && playerStatus.team !== 'white')}
                        className={`font-medium py-2 px-3 rounded-lg transition-all duration-200 text-sm shadow-lg hover:shadow-xl transform hover:scale-105 
                          ${playerStatus && playerStatus.team === 'white'
                            ? 'bg-gradient-to-r from-green-500 to-green-600 text-white hover:from-green-400 hover:to-green-500'
                            : playerStatus && playerStatus.team !== '' && playerStatus.team !== 'white'
                              ? 'bg-gray-400 text-gray-700 cursor-not-allowed'
                              : 'bg-gradient-to-r from-white to-gray-100 text-black hover:from-gray-100 hover:to-gray-200'}`}
                      >
                        {playerStatus && playerStatus.team === 'white' ? 'Reconnect White' : 'Join White'}
                      </button>
                      <button
                        onClick={() => {
                          console.log('🔘 DEBUG: Black button clicked!');
                          handleJoinTeam(game.id, 'black');
                        }}
                        disabled={!isConnected || (playerStatus != null && playerStatus.team !== '' && playerStatus.team !== 'black')}
                        className={`font-medium py-2 px-3 rounded-lg transition-all duration-200 text-sm shadow-lg hover:shadow-xl transform hover:scale-105 
                          ${playerStatus && playerStatus.team === 'black'
                            ? 'bg-gradient-to-r from-green-500 to-green-600 text-white hover:from-green-400 hover:to-green-500'
                            : playerStatus && playerStatus.team !== '' && playerStatus.team !== 'black'
                              ? 'bg-gray-400 text-gray-700 cursor-not-allowed'
                              : 'bg-gradient-to-r from-gray-700 to-gray-800 text-white hover:from-gray-600 hover:to-gray-700'}`}
                      >
                        {playerStatus && playerStatus.team === 'black' ? 'Reconnect Black' : 'Join Black'}
                      </button>
                      <button
                        onClick={() => {
                          console.log('🔘 DEBUG: Watch button clicked!');
                          handleWatchGame(game.id);
                        }}
                        className="bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 text-white font-medium py-2 px-3 rounded-lg transition-all duration-200 text-sm shadow-lg hover:shadow-xl transform hover:scale-105 flex items-center justify-center space-x-1"
                      >
                        <Eye className="w-3 h-3" />
                        <span>Watch</span>
                      </button>
                    </div>
                  )}

                  {game.status === 'active' && isTooOld && (
                    <button
                      onClick={() => handleWatchGame(game.id)}
                      className="w-full bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 text-white font-medium py-2 px-3 rounded-lg transition-all duration-200 text-sm shadow-lg hover:shadow-xl transform hover:scale-105 flex items-center justify-center space-x-2"
                    >
                      <Eye className="w-4 h-4" />
                      <span>Watch Game</span>
                    </button>
                  )}

                  {game.status === 'waiting' && (
                    <div className="grid grid-cols-2 gap-2">
                      {!game.whitePlayer && (
                        <button
                          onClick={() => handleJoinTeam(game.id, 'white')}
                          disabled={!isConnected || (playerStatus != null && playerStatus.team !== '' && playerStatus.team !== 'white')}
                          className={`font-medium py-3 px-4 rounded-lg transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105 
                            ${playerStatus && playerStatus.team === 'white'
                              ? 'bg-gradient-to-r from-green-500 to-green-600 text-white hover:from-green-400 hover:to-green-500'
                              : playerStatus && playerStatus.team !== '' && playerStatus.team !== 'white'
                                ? 'bg-gray-400 text-gray-700 cursor-not-allowed'
                                : 'bg-gradient-to-r from-white to-gray-100 text-black hover:from-gray-100 hover:to-gray-200'}`}
                        >
                          {playerStatus && playerStatus.team === 'white' ? 'Reconnect White' : 'Join White'}
                        </button>
                      )}
                      {!game.blackPlayer && (
                        <button
                          onClick={() => handleJoinTeam(game.id, 'black')}
                          disabled={!isConnected || (playerStatus != null && playerStatus.team !== '' && playerStatus.team !== 'black')}
                          className={`font-medium py-3 px-4 rounded-lg transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105 
                            ${playerStatus && playerStatus.team === 'black'
                              ? 'bg-gradient-to-r from-green-500 to-green-600 text-white hover:from-green-400 hover:to-green-500'
                              : playerStatus && playerStatus.team !== '' && playerStatus.team !== 'black'
                                ? 'bg-gray-400 text-gray-700 cursor-not-allowed'
                                : 'bg-gradient-to-r from-gray-700 to-gray-800 text-white hover:from-gray-600 hover:to-gray-700'}`}
                        >
                          {playerStatus && playerStatus.team === 'black' ? 'Reconnect Black' : 'Join Black'}
                        </button>
                      )}
                    </div>
                  )}

                  {(game.status === 'completed' || game.status === 'ended') && (
                    <button
                      onClick={() => openStatsModal(game)}
                      className="w-full bg-gradient-to-r from-purple-600 to-purple-700 hover:from-purple-500 hover:to-purple-600 text-white font-medium py-3 px-4 rounded-lg transition-all duration-200 flex items-center justify-center space-x-2 shadow-lg hover:shadow-xl transform hover:scale-105"
                    >
                      <Trophy className="w-4 h-4" />
                      <span>View Game Statistics</span>
                    </button>
                  )}
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {games.length === 0 && (
        <div className="text-center py-16">
          <div className="bg-gradient-to-br from-gray-800 to-gray-900 rounded-2xl p-12 border border-gray-700 shadow-xl">
            <Crown className="w-16 h-16 text-gray-500 mx-auto mb-6" />
            <div className="text-gray-400 text-xl mb-4">No active games found</div>
            <p className="text-gray-500 mb-6">Be the first to start a game! Use Quick Match to find an opponent.</p>
          </div>
        </div>
      )}

      {/* Board Zoom Modal */}
      {selectedGameInfo && (
        <div
          className="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
          onClick={closeBoardModal}
        >
          <div
            className="bg-gradient-to-br from-gray-800 via-gray-900 to-black rounded-2xl p-8 max-w-3xl w-full max-h-screen overflow-auto border border-gray-700 shadow-2xl"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Modal Header */}
            <div className="flex items-center justify-between mb-8">
              <div>
                <h2 className="text-3xl font-bold text-white mb-3">Game Preview</h2>
                <div className="flex items-center space-x-6 text-sm">
                  <div className="flex items-center space-x-2 text-gray-300">
                    <Swords className="w-4 h-4" />
                    <span>Move {selectedGameInfo.currentMove}</span>
                  </div>
                  <div className="flex items-center space-x-2 text-gray-300">
                    <div className={`w-3 h-3 rounded-full ${selectedGameInfo.currentTurn === 'white' ? 'bg-white' : 'bg-gray-600'}`} />
                    <span className="capitalize">{selectedGameInfo.currentTurn} to move</span>
                  </div>
                  <div className="flex items-center space-x-2 text-green-400">
                    <Trophy className="w-4 h-4" />
                    <span>{selectedGameInfo.totalPot.toFixed(2)} USDC pot</span>
                  </div>
                </div>
              </div>
              <button
                onClick={closeBoardModal}
                className="text-gray-400 hover:text-white text-3xl font-bold p-2 hover:bg-gray-700 rounded-lg transition-colors"
              >
                ×
              </button>
            </div>

            {/* Large Board */}
            <div className="flex justify-center mb-8">
              <BoardPreview
                boardState={selectedGameInfo.boardState}
                isLarge={true}
              />
            </div>

            {/* Team Stats */}
            <div className="grid grid-cols-2 gap-6 mb-8">
              <div className="bg-gradient-to-r from-gray-700/30 to-gray-600/30 rounded-xl p-4 border border-gray-600/50">
                <div className="flex items-center space-x-3 mb-3">
                  <div className="w-6 h-6 bg-white rounded-full shadow-sm" />
                  <span className="text-white font-bold text-lg">White Team</span>
                </div>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-300">Players:</span>
                    <span className="text-white font-medium">{selectedGameInfo.whitePlayers}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-300">Contributed:</span>
                    <span className="text-green-400 font-medium">{selectedGameInfo.whitePot.toFixed(2)} USDC</span>
                  </div>
                </div>
              </div>

              <div className="bg-gradient-to-r from-gray-800/30 to-gray-900/30 rounded-xl p-4 border border-gray-700/50">
                <div className="flex items-center space-x-3 mb-3">
                  <div className="w-6 h-6 bg-gray-600 rounded-full shadow-sm" />
                  <span className="text-gray-300 font-bold text-lg">Black Team</span>
                </div>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-300">Players:</span>
                    <span className="text-gray-300 font-medium">{selectedGameInfo.blackPlayers}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-300">Contributed:</span>
                    <span className="text-green-400 font-medium">{selectedGameInfo.blackPot.toFixed(2)} USDC</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Player Statistics for Ended Games */}
            {(selectedGameInfo.status === 'ended' || selectedGameInfo.status === 'completed') &&
              (selectedGameInfo.whiteTeamPlayers || selectedGameInfo.blackTeamPlayers) && (
                <div className="mb-8">
                  <PlayerStatistics
                    whiteTeamPlayers={selectedGameInfo.whiteTeamPlayers}
                    blackTeamPlayers={selectedGameInfo.blackTeamPlayers}
                    darkTheme={true}
                  />
                </div>
              )}

            {/* Game Actions */}
            <div className="flex justify-center space-x-4">
              {(() => {
                const playerStatusKey = address ? `${selectedGameInfo.id}-${address}` : null;
                const playerStatus = playerStatusKey ? playerStatusCache[playerStatusKey] : null;

                return (
                  <>
                    <button
                      onClick={() => handleJoinTeamFromModal(selectedGameInfo.id, 'white')}
                      disabled={!isConnected || (playerStatus != null && playerStatus.team !== '' && playerStatus.team !== 'white')}
                      className={`font-medium py-3 px-6 rounded-xl transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105 
                        ${playerStatus && playerStatus.team === 'white'
                          ? 'bg-gradient-to-r from-green-500 to-green-600 text-white hover:from-green-400 hover:to-green-500'
                          : playerStatus && playerStatus.team !== '' && playerStatus.team !== 'white'
                            ? 'bg-gray-400 text-gray-700 cursor-not-allowed'
                            : 'bg-gradient-to-r from-white to-gray-100 text-black hover:from-gray-100 hover:to-gray-200'}`}
                    >
                      {playerStatus && playerStatus.team === 'white' ? 'Reconnect White' : 'Join White Team'}
                    </button>
                    <button
                      onClick={() => handleJoinTeamFromModal(selectedGameInfo.id, 'black')}
                      disabled={!isConnected || (playerStatus != null && playerStatus.team !== '' && playerStatus.team !== 'black')}
                      className={`font-medium py-3 px-6 rounded-xl transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105 
                        ${playerStatus && playerStatus.team === 'black'
                          ? 'bg-gradient-to-r from-green-500 to-green-600 text-white hover:from-green-400 hover:to-green-500'
                          : playerStatus && playerStatus.team !== '' && playerStatus.team !== 'black'
                            ? 'bg-gray-400 text-gray-700 cursor-not-allowed'
                            : 'bg-gradient-to-r from-gray-700 to-gray-800 text-white hover:from-gray-600 hover:to-gray-700'}`}
                    >
                      {playerStatus && playerStatus.team === 'black' ? 'Reconnect Black' : 'Join Black Team'}
                    </button>
                    <button
                      onClick={() => handleWatchGameFromModal(selectedGameInfo.id)}
                      className="bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 text-white font-medium py-3 px-6 rounded-xl transition-all duration-200 flex items-center space-x-2 shadow-lg hover:shadow-xl transform hover:scale-105"
                    >
                      <Eye className="w-4 h-4" />
                      <span>Watch Game</span>
                    </button>
                  </>
                );
              })()}
            </div>
          </div>
        </div>
      )}

      {/* Statistics Modal for Ended Games */}
      {selectedStatsGame && (
        <div
          className="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
          onClick={closeStatsModal}
        >
          <div
            className="bg-gradient-to-br from-gray-800 via-gray-900 to-black rounded-2xl p-8 max-w-4xl w-full max-h-screen overflow-auto border border-gray-700 shadow-2xl"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Modal Header */}
            <div className="flex items-center justify-between mb-8">
              <div>
                <h2 className="text-3xl font-bold text-white mb-3">Game Statistics</h2>
                <div className="flex items-center space-x-6 text-sm">
                  <div className="flex items-center space-x-2 text-gray-300">
                    <Swords className="w-4 h-4" />
                    <span>Final Move: {selectedStatsGame.currentMove}</span>
                  </div>
                  <div className="flex items-center space-x-2 text-green-400">
                    <Trophy className="w-4 h-4" />
                    <span>{selectedStatsGame.totalPot.toFixed(2)} USDC Total Pot</span>
                  </div>
                  {selectedStatsGame.winner && (
                    <div className="flex items-center space-x-2 text-purple-400">
                      <Crown className="w-4 h-4" />
                      <span className="capitalize">
                        {selectedStatsGame.winner === 'draw' ? 'Draw' : `${selectedStatsGame.winner} Team Won`}
                      </span>
                    </div>
                  )}
                </div>
              </div>
              <button
                onClick={closeStatsModal}
                className="text-gray-400 hover:text-white text-3xl font-bold p-2 hover:bg-gray-700 rounded-lg transition-colors"
              >
                ×
              </button>
            </div>

            {/* Team Statistics and Player Lists */}
            {(selectedStatsGame.whiteTeamPlayers || selectedStatsGame.blackTeamPlayers) && (
              <div className="mb-8">
                <h3 className="font-bold text-white text-center text-lg mb-6">Team Statistics</h3>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* White Team */}
                  {selectedStatsGame.whiteTeamPlayers && selectedStatsGame.whiteTeamPlayers.length > 0 && (
                    <div className="bg-gradient-to-r from-gray-700/30 to-gray-600/30 p-6 rounded-lg border border-gray-600/50">
                      <h4 className="font-bold text-white mb-4 flex items-center gap-2 text-lg">
                        <div className="w-5 h-5 bg-white rounded-full"></div>
                        White Team
                      </h4>

                      {/* Team Summary */}
                      <div className="space-y-2 mb-6 text-sm">
                        <div className="flex justify-between">
                          <span className="text-gray-300">Total Players:</span>
                          <span className="text-white font-medium">{selectedStatsGame.whitePlayers}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-gray-300">Total Contributed:</span>
                          <span className="text-green-400 font-medium">{selectedStatsGame.whitePot.toFixed(3)} USDC</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-gray-300">Total Votes:</span>
                          <span className="text-blue-400 font-medium">
                            {selectedStatsGame.whiteTeamPlayers.reduce((sum, player) => sum + player.totalVotes, 0)}
                          </span>
                        </div>
                      </div>

                      {/* Players List */}
                      <div>
                        <h5 className="font-semibold text-white mb-3">Players:</h5>
                        <div className="space-y-2">
                          {[...selectedStatsGame.whiteTeamPlayers].sort((a, b) => b.totalVotes - a.totalVotes).map((player, index) => (
                            <div key={player.walletAddress} className="flex justify-between items-center text-sm bg-gray-800/50 p-3 rounded border border-gray-700/50">
                              <div className="flex items-center gap-2">
                                <span className="text-white font-semibold bg-gray-600 px-2 py-1 rounded text-xs">#{index + 1}</span>
                                <span className="font-mono text-gray-300">
                                  {player.walletAddress.length <= 10 ? player.walletAddress : `${player.walletAddress.slice(0, 4)}...${player.walletAddress.slice(-4)}`}
                                </span>
                              </div>
                              <div className="text-right">
                                <span className="text-white font-medium">
                                  {player.totalVotes} votes {player.totalSpent.toFixed(3)} USDC
                                </span>
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>
                  )}

                  {/* Black Team */}
                  {selectedStatsGame.blackTeamPlayers && selectedStatsGame.blackTeamPlayers.length > 0 && (
                    <div className="bg-gradient-to-r from-gray-800/30 to-gray-900/30 p-6 rounded-lg border border-gray-700/50">
                      <h4 className="font-bold text-gray-300 mb-4 flex items-center gap-2 text-lg">
                        <div className="w-5 h-5 bg-gray-600 rounded-full"></div>
                        Black Team
                      </h4>

                      {/* Team Summary */}
                      <div className="space-y-2 mb-6 text-sm">
                        <div className="flex justify-between">
                          <span className="text-gray-300">Total Players:</span>
                          <span className="text-gray-300 font-medium">{selectedStatsGame.blackPlayers}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-gray-300">Total Contributed:</span>
                          <span className="text-green-400 font-medium">{selectedStatsGame.blackPot.toFixed(3)} USDC</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-gray-300">Total Votes:</span>
                          <span className="text-blue-400 font-medium">
                            {selectedStatsGame.blackTeamPlayers.reduce((sum, player) => sum + player.totalVotes, 0)}
                          </span>
                        </div>
                      </div>

                      {/* Players List */}
                      <div>
                        <h5 className="font-semibold text-gray-300 mb-3">Players:</h5>
                        <div className="space-y-2">
                          {[...selectedStatsGame.blackTeamPlayers].sort((a, b) => b.totalVotes - a.totalVotes).map((player, index) => (
                            <div key={player.walletAddress} className="flex justify-between items-center text-sm bg-gray-800/50 p-3 rounded border border-gray-700/50">
                              <div className="flex items-center gap-2">
                                <span className="text-gray-300 font-semibold bg-gray-700 px-2 py-1 rounded text-xs">#{index + 1}</span>
                                <span className="font-mono text-gray-300">
                                  {player.walletAddress.length <= 10 ? player.walletAddress : `${player.walletAddress.slice(0, 4)}...${player.walletAddress.slice(-4)}`}
                                </span>
                              </div>
                              <div className="text-right">
                                <span className="text-gray-300 font-medium">
                                  {player.totalVotes} votes {player.totalSpent.toFixed(3)} USDC
                                </span>
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Debug section when no player data */}
            {!(selectedStatsGame.whiteTeamPlayers || selectedStatsGame.blackTeamPlayers) && (
              <div className="mb-8 bg-red-900/20 border border-red-500/50 rounded-lg p-4">
                <h3 className="font-bold text-red-400 text-center text-lg mb-2">No Player Statistics Available</h3>
                <p className="text-red-300 text-center text-sm">
                  Player statistics are not available for this game. This might be because:
                </p>
                <ul className="text-red-300 text-xs mt-2 space-y-1">
                  <li>• The game ended before the new statistics system was implemented</li>
                  <li>• The game data is incomplete</li>
                  <li>• There's an issue with data transmission from the backend</li>
                </ul>
              </div>
            )}

            {/* Action Buttons */}
            <div className="flex justify-center space-x-4">
              <button
                onClick={() => {
                  closeStatsModal();
                  handleWatchGame(selectedStatsGame.id);
                }}
                className="bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 text-white font-medium py-3 px-6 rounded-xl transition-all duration-200 flex items-center space-x-2 shadow-lg hover:shadow-xl transform hover:scale-105"
              >
                <Eye className="w-4 h-4" />
                <span>View Game Board</span>
              </button>
              <button
                onClick={closeStatsModal}
                className="bg-gradient-to-r from-gray-600 to-gray-700 hover:from-gray-500 hover:to-gray-600 text-white font-medium py-3 px-6 rounded-xl transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </main>
  );
};

export default GameLobby;
