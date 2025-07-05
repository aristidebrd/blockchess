import React from 'react';
import { Crown, AlertTriangle, Home, Search } from 'lucide-react';

interface GameNotFoundProps {
    gameId: string;
    onBackToLobby: () => void;
}

const GameNotFound: React.FC<GameNotFoundProps> = ({ gameId, onBackToLobby }) => {
    return (
        <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex items-center justify-center px-4">
            <div className="max-w-md w-full text-center">
                {/* Error Icon */}
                <div className="mb-8">
                    <div className="mx-auto w-24 h-24 bg-red-500/20 rounded-full flex items-center justify-center mb-4">
                        <AlertTriangle className="w-12 h-12 text-red-400" />
                    </div>
                    <h1 className="text-6xl font-bold text-white mb-2">404</h1>
                    <h2 className="text-2xl font-semibold text-gray-300 mb-4">Game Not Found</h2>
                </div>

                {/* Error Message */}
                <div className="bg-black/30 backdrop-blur-sm rounded-xl p-6 mb-8 border border-gray-700">
                    <p className="text-gray-300 mb-4">
                        The game <span className="font-mono text-yellow-400 bg-gray-800 px-2 py-1 rounded">{gameId}</span> doesn't exist or has ended.
                    </p>
                    <p className="text-gray-400 text-sm">
                        Games can only be created through matchmaking. You cannot join a game by entering a random game ID.
                    </p>
                </div>

                {/* Action Buttons */}
                <div className="space-y-4">
                    <button
                        onClick={onBackToLobby}
                        className="w-full bg-gradient-to-r from-yellow-400 to-yellow-600 text-black font-semibold py-3 px-6 rounded-xl hover:from-yellow-300 hover:to-yellow-500 transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105 flex items-center justify-center space-x-2"
                    >
                        <Home className="w-5 h-5" />
                        <span>Back to Lobby</span>
                    </button>
                </div>

                {/* Help Text */}
                <div className="mt-8 p-4 bg-blue-500/10 rounded-lg border border-blue-500/20">
                    <div className="flex items-center justify-center space-x-2 mb-2">
                        <Crown className="w-5 h-5 text-blue-400" />
                        <span className="text-blue-300 font-medium">How to Play</span>
                    </div>
                    <p className="text-blue-200 text-sm">
                        Start a new game through matchmaking or join an existing active game from the lobby.
                    </p>
                </div>
            </div>
        </div>
    );
};

export default GameNotFound; 
