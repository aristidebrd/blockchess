import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useAccount } from 'wagmi';

interface OnlinePlayersContextType {
    onlinePlayersCount: number;
    isCurrentUserOnline: boolean;
}

const OnlinePlayersContext = createContext<OnlinePlayersContextType | undefined>(undefined);

export const useOnlinePlayers = () => {
    const context = useContext(OnlinePlayersContext);
    if (context === undefined) {
        throw new Error('useOnlinePlayers must be used within an OnlinePlayersProvider');
    }
    return context;
};

interface OnlinePlayersProviderProps {
    children: ReactNode;
}

export const OnlinePlayersProvider: React.FC<OnlinePlayersProviderProps> = ({ children }) => {
    const { isConnected, address } = useAccount();
    const [onlinePlayersCount, setOnlinePlayersCount] = useState<number>(0);

    useEffect(() => {
        const updateOnlineCount = () => {
            const totalPlayers = isConnected ? 1 : 0;
            setOnlinePlayersCount(totalPlayers);
        };

        // Update immediately when connection status changes
        updateOnlineCount();

        // Update periodically to simulate players joining/leaving
        const interval = setInterval(updateOnlineCount, 10000); // Every 10 seconds

        return () => clearInterval(interval);
    }, [isConnected, address]);

    // Send heartbeat when user is connected (for future backend integration)
    useEffect(() => {
        if (isConnected && address) {
            const heartbeatInterval = setInterval(() => {
            }, 10000); // Every 10 seconds

            return () => clearInterval(heartbeatInterval);
        }
    }, [isConnected, address]);

    const value = {
        onlinePlayersCount,
        isCurrentUserOnline: isConnected,
    };

    return (
        <OnlinePlayersContext.Provider value={value}>
            {children}
        </OnlinePlayersContext.Provider>
    );
};
