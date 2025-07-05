import React from 'react';
import { Users, Vote, Coins } from 'lucide-react';
import { PlayerStats } from '../utils/chess';

interface PlayerStatisticsProps {
    whiteTeamPlayers?: PlayerStats[];
    blackTeamPlayers?: PlayerStats[];
    className?: string;
    darkTheme?: boolean;
}

const PlayerStatistics: React.FC<PlayerStatisticsProps> = ({
    whiteTeamPlayers,
    blackTeamPlayers,
    className = '',
    darkTheme = false
}) => {
    const formatWalletAddress = (address: string) => {
        if (address.length <= 10) return address;
        return `${address.slice(0, 6)}...${address.slice(-4)}`;
    };

    const sortPlayersByVotes = (players: PlayerStats[]) => {
        return [...players].sort((a, b) => b.totalVotes - a.totalVotes);
    };

    if (!whiteTeamPlayers?.length && !blackTeamPlayers?.length) {
        return null;
    }

    return (
        <div className={`space-y-4 ${className}`}>
            <h3 className={`font-bold text-center text-lg flex items-center justify-center gap-2 ${darkTheme ? 'text-white' : 'text-gray-800'}`}>
                <Users className="w-5 h-5" />
                Player Statistics
            </h3>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* White Team Players */}
                {whiteTeamPlayers && whiteTeamPlayers.length > 0 && (
                    <div className="bg-amber-50 p-4 rounded-lg border-2 border-amber-200">
                        <h4 className="font-bold text-amber-800 mb-3 flex items-center gap-2">
                            <div className="w-3 h-3 bg-amber-400 rounded"></div>
                            White Team ({whiteTeamPlayers.length} players)
                        </h4>
                        <div className="space-y-2">
                            {sortPlayersByVotes(whiteTeamPlayers).map((player, index) => (
                                <div key={player.walletAddress} className={`flex justify-between items-center text-sm p-2 rounded shadow-sm ${darkTheme ? 'bg-gray-700/50' : 'bg-white'}`}>
                                    <div className="flex items-center gap-2">
                                        <span className="text-amber-600 font-semibold min-w-[24px]">#{index + 1}</span>
                                        <span className={`font-mono ${darkTheme ? 'text-gray-300' : 'text-gray-700'}`}>{formatWalletAddress(player.walletAddress)}</span>
                                    </div>
                                    <div className="text-right">
                                        <div className="flex items-center gap-1 text-amber-700 font-semibold">
                                            <Vote className="w-3 h-3" />
                                            <span>{player.totalVotes}</span>
                                        </div>
                                        <div className="flex items-center gap-1 text-amber-600 text-xs">
                                            <Coins className="w-3 h-3" />
                                            <span>${player.totalSpent.toFixed(3)}</span>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {/* Black Team Players */}
                {blackTeamPlayers && blackTeamPlayers.length > 0 && (
                    <div className="bg-gray-50 p-4 rounded-lg border-2 border-gray-200">
                        <h4 className="font-bold text-gray-800 mb-3 flex items-center gap-2">
                            <div className="w-3 h-3 bg-gray-600 rounded"></div>
                            Black Team ({blackTeamPlayers.length} players)
                        </h4>
                        <div className="space-y-2">
                            {sortPlayersByVotes(blackTeamPlayers).map((player, index) => (
                                <div key={player.walletAddress} className={`flex justify-between items-center text-sm p-2 rounded shadow-sm ${darkTheme ? 'bg-gray-700/50' : 'bg-white'}`}>
                                    <div className="flex items-center gap-2">
                                        <span className="text-gray-600 font-semibold min-w-[24px]">#{index + 1}</span>
                                        <span className={`font-mono ${darkTheme ? 'text-gray-300' : 'text-gray-700'}`}>{formatWalletAddress(player.walletAddress)}</span>
                                    </div>
                                    <div className="text-right">
                                        <div className="flex items-center gap-1 text-gray-700 font-semibold">
                                            <Vote className="w-3 h-3" />
                                            <span>{player.totalVotes}</span>
                                        </div>
                                        <div className="flex items-center gap-1 text-gray-600 text-xs">
                                            <Coins className="w-3 h-3" />
                                            <span>${player.totalSpent.toFixed(3)}</span>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default PlayerStatistics;
