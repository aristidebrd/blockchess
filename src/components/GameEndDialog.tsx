import React from 'react';
import { Trophy, Users, Vote, Coins, TrendingUp } from 'lucide-react';
import { GameEndInfo } from '../utils/chess';

interface GameEndDialogProps {
    gameEndInfo: GameEndInfo;
    returnToLobby: () => void;
}

const GameEndDialog: React.FC<GameEndDialogProps> = ({ gameEndInfo, returnToLobby }) => {
    const getWinnerText = () => {
        if (gameEndInfo.winner === 'draw') {
            return 'Draw Game!';
        }
        return `${gameEndInfo.winner.charAt(0).toUpperCase() + gameEndInfo.winner.slice(1)} Team Wins!`;
    };

    const getReasonText = () => {
        switch (gameEndInfo.reason) {
            case 'checkmate':
                return 'by Checkmate';
            case 'stalemate':
                return 'by Stalemate';
            case 'insufficient_material':
                return 'by Insufficient Material';
            case 'threefold_repetition':
                return 'by Threefold Repetition';
            case 'fifty_move_rule':
                return 'by Fifty Move Rule';
            default:
                return `by ${gameEndInfo.reason}`;
        }
    };

    const getWinnerColor = () => {
        if (gameEndInfo.winner === 'white') return 'text-amber-600 bg-amber-50';
        if (gameEndInfo.winner === 'black') return 'text-gray-700 bg-gray-100';
        return 'text-blue-600 bg-blue-50';
    };

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full max-h-[90vh] overflow-y-auto">
                {/* Header */}
                <div className={`p-6 text-center rounded-t-2xl ${getWinnerColor()}`}>
                    <Trophy className="w-16 h-16 mx-auto mb-4" />
                    <h2 className="text-2xl font-bold mb-2">
                        {getWinnerText()}
                    </h2>
                    <p className="text-lg opacity-80">
                        {getReasonText()}
                    </p>
                </div>

                {/* Game Summary */}
                <div className="p-6 space-y-6">
                    {/* Move Count */}
                    <div className="text-center">
                        <div className="inline-flex items-center gap-2 px-4 py-2 bg-purple-100 text-purple-700 rounded-lg">
                            <TrendingUp className="w-5 h-5" />
                            <span className="font-semibold">Game lasted {gameEndInfo.currentMove} moves</span>
                        </div>
                    </div>

                    {/* Teams Stats */}
                    <div className="grid grid-cols-2 gap-4">
                        <div className="bg-amber-50 p-4 rounded-lg border-2 border-amber-200">
                            <h3 className="font-bold text-amber-800 mb-3 flex items-center gap-2">
                                <div className="w-4 h-4 bg-amber-400 rounded"></div>
                                White Team
                            </h3>
                            <div className="space-y-2 text-sm">
                                <div className="flex items-center gap-2">
                                    <Users className="w-4 h-4 text-amber-600" />
                                    <span>{gameEndInfo.whitePlayers} players</span>
                                </div>
                                <div className="flex items-center gap-2">
                                    <Vote className="w-4 h-4 text-amber-600" />
                                    <span>{gameEndInfo.whiteTeamTotalVotes} votes</span>
                                </div>
                                <div className="flex items-center gap-2">
                                    <Coins className="w-4 h-4 text-amber-600" />
                                    <span>${gameEndInfo.whitePot.toFixed(2)}</span>
                                </div>
                            </div>
                        </div>

                        <div className="bg-gray-50 p-4 rounded-lg border-2 border-gray-200">
                            <h3 className="font-bold text-gray-800 mb-3 flex items-center gap-2">
                                <div className="w-4 h-4 bg-gray-600 rounded"></div>
                                Black Team
                            </h3>
                            <div className="space-y-2 text-sm">
                                <div className="flex items-center gap-2">
                                    <Users className="w-4 h-4 text-gray-600" />
                                    <span>{gameEndInfo.blackPlayers} players</span>
                                </div>
                                <div className="flex items-center gap-2">
                                    <Vote className="w-4 h-4 text-gray-600" />
                                    <span>{gameEndInfo.blackTeamTotalVotes} votes</span>
                                </div>
                                <div className="flex items-center gap-2">
                                    <Coins className="w-4 h-4 text-gray-600" />
                                    <span>${gameEndInfo.blackPot.toFixed(2)}</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Total Pot */}
                    <div className="bg-gradient-to-r from-yellow-100 to-amber-100 p-4 rounded-lg border-2 border-yellow-300">
                        <h3 className="font-bold text-yellow-800 mb-2 text-center">Total Pot</h3>
                        <div className="text-center">
                            <span className="text-2xl font-bold text-yellow-700">${gameEndInfo.totalPot.toFixed(2)}</span>
                        </div>
                    </div>

                    {/* Your Performance */}
                    <div className="bg-blue-50 p-4 rounded-lg border-2 border-blue-200">
                        <h3 className="font-bold text-blue-800 mb-2 text-center">Your Performance</h3>
                        <div className="text-center">
                            <div className="inline-flex items-center gap-2 px-3 py-1 bg-blue-200 text-blue-800 rounded-full">
                                <Vote className="w-4 h-4" />
                                <span className="font-semibold">{gameEndInfo.playerVotes} votes cast</span>
                            </div>
                        </div>
                    </div>

                    {/* Close Button */}
                    <button
                        onClick={returnToLobby}
                        className="w-full bg-gradient-to-r from-purple-600 to-indigo-600 text-white py-3 rounded-lg font-semibold hover:from-purple-700 hover:to-indigo-700 transition-all duration-200 transform hover:scale-105"
                    >
                        Return to Lobby
                    </button>
                </div>
            </div>
        </div>
    );
};

export default GameEndDialog;
