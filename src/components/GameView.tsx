import React, { useState, useEffect } from 'react';
import { useAccount } from 'wagmi';
import ChessBoard from './ChessBoard';
import VotingPanel from './VotingPanel';
import MoveConfirmDialog from './MoveConfirmDialog';
import GameEndDialog from './GameEndDialog';
import { GameState, GameEndInfo, GameInfo } from '../utils/chess';
import { initialBoard } from '../utils/chess';
import { gameService } from '../services/gameService';

interface GameViewProps {
    gameId: string;
    playerRole: 'white' | 'black' | 'spectator' | 'none';
    onJoinTeam: (side: 'white' | 'black') => void;
    isFromMatchmaking: boolean;
    onBackToLobby?: () => void;
}

const GameView: React.FC<GameViewProps> = ({
    gameId,
    playerRole,
    onJoinTeam,
    isFromMatchmaking,
    onBackToLobby
}) => {
    const { isConnected } = useAccount();

    // Local UI state
    const [gameState, setGameState] = useState<GameState>({
        board: initialBoard,
        currentPlayer: 'white',
        moveNumber: 1,
        gameStatus: 'playing',
        lastMove: undefined,
        gameStartTime: new Date()
    });

    const [backendGameState, setBackendGameState] = useState<GameInfo | null>(null);
    const [showMoveConfirm, setShowMoveConfirm] = useState(false);
    const [gameEndInfo, setGameEndInfo] = useState<GameEndInfo | null>(null);

    // Initialize game and subscribe to updates
    useEffect(() => {
        console.log(`GameView: Initializing game ${gameId}`);

        gameService.initializeGame(gameId);

        const unsubscribeGameState = gameService.onGameStateUpdate((state) => {
            console.log(`GameView: Received game state update:`, state);
            setBackendGameState(state);

            // Convert backend board state to frontend format if available
            if (state.boardState) {
                const frontendBoard = state.boardState.map((row, rowIndex) =>
                    row.map((piece, colIndex) => {
                        if (!piece) return null;

                        // Calculate the correct position for this piece
                        const file = String.fromCharCode(97 + colIndex); // a, b, c...
                        const rank = (8 - rowIndex).toString(); // 8, 7, 6...
                        const calculatedPosition = `${file}${rank}`;

                        return {
                            type: piece.type,
                            color: piece.color,
                            position: calculatedPosition, // Use calculated position instead of piece.position
                            hasMoved: piece.hasMoved
                        };
                    })
                );

                // Update local game state to match backend
                setGameState(prev => ({
                    ...prev,
                    board: frontendBoard,
                    currentPlayer: state.currentTurn,
                    moveNumber: state.currentMove,
                    gameStatus: 'playing',
                    lastMove: state.lastExecutedMove ? {
                        from: state.lastExecutedMove.from,
                        to: state.lastExecutedMove.to,
                        piece: state.lastExecutedMove.piece,
                        notation: `${state.lastExecutedMove.from}-${state.lastExecutedMove.to}`,
                        votes: 0
                    } : prev.lastMove
                }));
            }
        });

        const unsubscribeGameEnd = gameService.onGameEnd((gameEndData) => {
            setGameEndInfo(gameEndData);
        });

        // Cleanup on unmount
        return () => {
            unsubscribeGameState();
            unsubscribeGameEnd();
            gameService.cleanupGame();
        };
    }, [gameId]);

    // Handle voting
    const handleVote = async (moveId: string): Promise<boolean> => {
        const result = await gameService.voteOnMove(moveId, isConnected);

        if (!result.success && result.error) {
            alert(result.error);
        }

        return result.success;
    };

    // Handle move creation
    const handleCreateMove = async (from: string, to: string) => {
        if (!backendGameState || !backendGameState.boardState) {
            alert('Game not ready for move creation');
            return;
        }

        const result = await gameService.createMove(from, to, backendGameState.boardState);

        if (result.success) {
            setShowMoveConfirm(true);
        } else if (result.error) {
            if (result.error !== 'Not your turn' && result.error !== 'You have already voted this round') {
                alert(result.error);
            }
        }
    };

    // Handle move confirmation
    const handleConfirmMove = async () => {
        const result = await gameService.confirmPendingMove();

        if (!result.success && result.error) {
            alert(result.error);
        }

        setShowMoveConfirm(false);
    };

    // Handle move cancellation
    const handleCancelMove = () => {
        gameService.cancelPendingMove();
        setShowMoveConfirm(false);
    };

    // Get pending move for confirmation dialog
    const pendingMove = gameService.getPendingMove();

    // Handle game end dialog close
    const handleGameEndClose = () => {
        setGameEndInfo(null);
        // Navigate back to lobby if callback is provided, otherwise navigate to root
        if (onBackToLobby) {
            onBackToLobby();
        } else {
            // Navigate to the main page
            window.location.href = '/';
        }
    };

    return (
        <div>
            {/* Main content */}
            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                {/* Show loading state until we have game state */}
                {!backendGameState ? (
                    <div className="text-center text-white">
                        <div className="text-xl mb-4">Loading game...</div>
                        <div className="w-8 h-8 border-2 border-white border-t-transparent rounded-full animate-spin mx-auto"></div>
                    </div>
                ) : (
                    <>
                        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                            {/* Left column - Chess board and voting panel */}
                            <div className="lg:col-span-2 space-y-6">
                                <ChessBoard
                                    gameState={backendGameState}
                                    onCreateMove={handleCreateMove}
                                    isInteractive={true}
                                    playerSide={playerRole}
                                    hasVoted={gameService.hasPlayerVoted(gameId)}
                                />

                                {/* Voting Panel - Always visible under the chess board */}
                                {playerRole !== 'none' && playerRole !== 'spectator' && backendGameState && backendGameState.currentTurn === playerRole && (
                                    <VotingPanel
                                        moves={backendGameState.proposedMoves.map(move => ({
                                            from: move.from,
                                            to: move.to,
                                            piece: {
                                                type: 'pawn' as const,
                                                color: 'white' as const,
                                                position: move.from
                                            },
                                            notation: `${move.from}-${move.to}`,
                                            votes: move.votes
                                        }))}
                                        onVote={(moveId) => handleVote(moveId)}
                                        isVotingEnabled={gameService.isVotingEnabled()}
                                        userVotes={backendGameState.proposedMoves.map(move => move.moveId)}
                                        playerRole={playerRole as 'white' | 'black'}
                                        backendGameState={backendGameState}
                                    />
                                )}

                                {/* Waiting message */}
                                {playerRole !== 'none' && playerRole !== 'spectator' && backendGameState && backendGameState.currentTurn !== playerRole && (
                                    <div className="bg-gradient-to-r from-gray-700 to-gray-800 rounded-xl p-6 text-center">
                                        <h3 className="text-lg font-bold text-white mb-2">
                                            Waiting for {backendGameState.currentTurn === 'white' ? 'White' : 'Black'} Team
                                        </h3>
                                        <p className="text-gray-300">It's the other team's turn to vote and make their move.</p>
                                    </div>
                                )}
                            </div>

                            {/* Right column - Timer and Stats */}
                            <div className="space-y-6">
                                {backendGameState ? (
                                    <div className="bg-gradient-to-br from-slate-800 to-slate-900 rounded-xl p-6 border border-gray-700">
                                        <div className="flex items-center justify-between mb-3">
                                            <div className="flex items-center">
                                                <div className={`w-2 h-2 rounded-full mr-2 ${backendGameState.currentTurn === 'white' ? 'bg-white' : 'bg-gray-400'}`} />
                                                <div className="text-white font-medium">
                                                    {backendGameState.currentTurn === 'white' ? 'White' : 'Black'} to move
                                                    {backendGameState.isCheckmate && <span className="ml-2 text-red-400 font-bold">CHECKMATE!</span>}
                                                    {backendGameState.isInCheck && !backendGameState.isCheckmate && <span className="ml-2 text-orange-400 font-bold">CHECK!</span>}
                                                </div>
                                            </div>
                                            <div className="text-3xl font-bold text-yellow-400">
                                                {Math.ceil(backendGameState.timeRemaining / 1000)}s
                                            </div>
                                            <div className="text-white font-medium">
                                                Move #{backendGameState.currentMove}
                                            </div>
                                        </div>
                                        <div className="w-full bg-gray-700 rounded-full h-2">
                                            <div
                                                className="bg-yellow-400 h-2 rounded-full transition-all duration-1000"
                                                style={{
                                                    width: `${Math.max(0, (backendGameState.timeRemaining / 30000) * 100)}%`
                                                }}
                                            />
                                        </div>
                                    </div>
                                ) : (
                                    <div className="bg-gradient-to-br from-slate-800 to-slate-900 rounded-xl p-6 border border-gray-700">
                                        <div className="text-center text-gray-400">
                                            Loading game timer...
                                        </div>
                                    </div>
                                )}

                                {/* Betting Info */}
                                <div className="bg-gradient-to-r from-blue-900/50 to-purple-900/50 rounded-xl p-6 border border-blue-500/30">
                                    <div className="flex justify-between items-center text-sm mb-2">
                                        <span className="text-gray-300">Fixed bet per vote:</span>
                                        <span className="text-blue-400 font-bold">0.01 USDC</span>
                                    </div>
                                    {backendGameState && (
                                        <>
                                            <div className="flex justify-between items-center text-sm mb-2">
                                                <span className="text-gray-300">White pot:</span>
                                                <span className="text-white font-bold">{backendGameState.whitePot.toFixed(2)} USDC</span>
                                            </div>
                                            <div className="flex justify-between items-center text-sm mb-2">
                                                <span className="text-gray-300">Black pot:</span>
                                                <span className="text-gray-300 font-bold">{backendGameState.blackPot.toFixed(2)} USDC</span>
                                            </div>
                                        </>
                                    )}
                                    <div className="flex justify-between items-center text-sm">
                                        <span className="text-gray-300">Total pot:</span>
                                        <span className="text-green-400 font-bold">
                                            {backendGameState ? backendGameState.totalPot.toFixed(2) : '0.00'} USDC
                                        </span>
                                    </div>
                                </div>

                                {/* Voting Status */}
                                {playerRole !== 'none' && playerRole !== 'spectator' && backendGameState && (
                                    <div className="bg-gray-800/50 rounded-xl p-6 border border-gray-700">
                                        <div className="flex items-center justify-between text-sm mb-2">
                                            <span className="text-gray-300">Your votes:</span>
                                            <span className="text-blue-400 font-medium">{gameService.getPlayerTotalVotes(gameId)}</span>
                                        </div>
                                        <div className="flex items-center justify-between text-sm">
                                            <span className="text-gray-300">{playerRole === 'white' ? 'White' : 'Black'} team votes:</span>
                                            <span className="text-white font-medium">
                                                {playerRole === 'white' ? backendGameState.whiteTeamTotalVotes : backendGameState.blackTeamTotalVotes}
                                            </span>
                                        </div>
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* Bottom section - Team Status and Messages */}
                        <div className="mt-8">
                            {playerRole === 'spectator' ? (
                                <div className="bg-gradient-to-r from-gray-800 to-gray-900 rounded-xl p-6 text-center">
                                    <h3 className="text-xl font-bold text-white mb-2">Spectating Game</h3>
                                    <p className="text-gray-300">You are watching this game. Join a team to participate in voting!</p>
                                </div>
                            ) : playerRole === 'none' && !isFromMatchmaking ? (
                                <div className="bg-gradient-to-r from-yellow-600 to-orange-600 rounded-xl p-6 text-center">
                                    <h3 className="text-xl font-bold text-white mb-2">Choose Your Team</h3>
                                    <p className="text-gray-100 mb-4">Join the white or black team to participate in voting!</p>
                                    <div className="flex gap-4 justify-center">
                                        <button
                                            onClick={() => onJoinTeam('white')}
                                            className="bg-white text-black hover:bg-gray-200 font-medium py-2 px-6 rounded-lg transition-all duration-200"
                                        >
                                            Join White Team
                                        </button>
                                        <button
                                            onClick={() => onJoinTeam('black')}
                                            className="bg-gray-800 text-white hover:bg-gray-700 font-medium py-2 px-6 rounded-lg transition-all duration-200"
                                        >
                                            Join Black Team
                                        </button>
                                    </div>
                                </div>
                            ) : playerRole === 'none' && isFromMatchmaking ? (
                                <div className="bg-gradient-to-r from-blue-600 to-purple-600 rounded-xl p-6 text-center">
                                    <div className="flex items-center justify-center space-x-3">
                                        <div className="w-6 h-6 border-2 border-white border-t-transparent rounded-full animate-spin" />
                                        <div>
                                            <h3 className="text-xl font-bold text-white mb-1">Determining Your Team...</h3>
                                            <p className="text-blue-100">Please wait while we assign you to your team.</p>
                                        </div>
                                    </div>
                                </div>
                            ) : null}
                        </div>
                    </>
                )}
            </main>

            {/* Move Confirmation Dialog */}
            {pendingMove && (
                <MoveConfirmDialog
                    isOpen={showMoveConfirm}
                    from={pendingMove.from}
                    to={pendingMove.to}
                    piece={pendingMove.piece}
                    onConfirm={handleConfirmMove}
                    onCancel={handleCancelMove}
                />
            )}

            {/* Game End Dialog */}
            {gameEndInfo && (
                <GameEndDialog
                    gameEndInfo={gameEndInfo}
                    returnToLobby={handleGameEndClose}
                />
            )}
        </div>
    );
};

export default GameView; 
