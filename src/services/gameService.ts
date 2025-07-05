// Game service to communicate with Go backend
import { wsService } from './websocket';
import { GameInfo, ChessPiece, GameEndInfo } from '../utils/chess';

// New interfaces for game engine
export interface PendingMove {
    from: string;
    to: string;
    piece: ChessPiece;
}

export interface VoteResult {
    success: boolean;
    error?: string;
}

export interface MoveResult {
    success: boolean;
    moveId?: string;
    error?: string;
}

class GameService {
    private playerId: string;
    private walletAddress: string | null = null;
    private matchmakingCallbacks: ((players: string[], gameId: string, assignedSide?: 'white' | 'black') => void)[] = [];
    private currentGameState: GameInfo | null = null;
    private playerRoles: Map<string, 'white' | 'black' | 'spectator' | 'none'> = new Map();
    private gameStateCallbacks: Map<string, (state: GameInfo) => void> = new Map();
    private activeGameCleanups: Map<string, () => void> = new Map();
    private gameEndCallbacks: Map<string, (gameEndInfo: GameEndInfo) => void> = new Map();

    // Game engine state
    private pendingMove: PendingMove | null = null;
    private currentGameId: string | null = null;

    constructor() {
        this.playerId = `player_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        console.log('GameService created with playerId:', this.playerId);

        // Set up the player ID getter for the WebSocket service
        wsService.setPlayerIdGetter(() => this.getPlayerIdentifier());

        // Listen for client connected message to sync player ID with backend
        wsService.on('client_connected', (data) => {
            console.log('🔗 Client connected with ID:', data.clientId);
            this.playerId = data.clientId;
            console.log('🔗 Updated playerId to match backend:', this.playerId);
        });

        // Set up matchmaking listener
        wsService.on('match_found', (data) => {
            console.log('🎯 Match found event received:', data);
            console.log('🎯 Callbacks registered:', this.matchmakingCallbacks.length);

            // Trigger all registered matchmaking callbacks
            this.matchmakingCallbacks.forEach((callback, index) => {
                console.log(`🎯 Triggering callback ${index + 1}:`, data.players, data.gameId, data.assignedSide);
                callback(data.players, data.gameId, data.assignedSide);
            });

            // Clear callbacks since match was found
            console.log('🎯 Clearing matchmaking callbacks');
            this.matchmakingCallbacks = [];
        });
    }

    // Set wallet address for the current player
    setWalletAddress(address: string | null): void {
        this.walletAddress = address;
        console.log('🔐 Wallet address set:', address);
    }

    // Get the current player identifier (wallet address if available, otherwise fallback to generated ID)
    getPlayerIdentifier(): string {
        return this.walletAddress || this.playerId;
    }

    // Check if wallet is connected
    isWalletConnected(): boolean {
        return this.walletAddress !== null;
    }

    // Join matchmaking
    async joinMatchmaking(onMatchFound: (players: string[], gameId: string, assignedSide?: 'white' | 'black') => void): Promise<boolean> {
        // Validate wallet connection
        if (!this.isWalletConnected()) {
            console.error('Cannot join matchmaking without wallet connection');
            return false;
        }

        this.matchmakingCallbacks.push(onMatchFound);

        // Send join matchmaking message with wallet address
        wsService.joinMatchmaking(this.getPlayerIdentifier());

        return true;
    }

    // Leave matchmaking
    async leaveMatchmaking(): Promise<boolean> {
        this.matchmakingCallbacks = [];

        // Send leave matchmaking message
        wsService.leaveMatchmaking();

        return true;
    }

    // Join a team with wallet address validation
    async joinTeam(gameId: string, team: 'white' | 'black'): Promise<{ success: boolean; error?: string }> {
        console.log(`Joining ${team} team in game ${gameId}`);

        // Validate wallet connection
        if (!this.isWalletConnected()) {
            return { success: false, error: 'Please connect your wallet to join a team' };
        }

        // Store the role locally
        this.playerRoles.set(gameId, team);

        // Send join team message to backend with wallet address
        wsService.joinTeam(gameId, team);

        return { success: true };
    }

    // Watch game as spectator
    async watchGame(gameId: string): Promise<boolean> {
        console.log(`Watching game ${gameId}`);
        this.playerRoles.set(gameId, 'spectator');

        wsService.watchGame(gameId);

        return true;
    }

    // Get player role in game
    async getPlayerRole(gameId: string): Promise<{ role: 'white' | 'black' | 'spectator' | 'none' } | null> {
        // First check local cache
        const localRole = this.playerRoles.get(gameId);
        if (localRole && localRole !== 'none') {
            return { role: localRole };
        }

        // If no local role or role is 'none', check with backend if wallet is connected
        if (this.isWalletConnected()) {
            try {
                const playerStatus = await this.checkPlayerStatusWithBackend(gameId, this.getPlayerIdentifier());
                if (playerStatus && playerStatus.team && playerStatus.team !== '') {
                    // Update local cache
                    this.playerRoles.set(gameId, playerStatus.team as 'white' | 'black');
                    return { role: playerStatus.team as 'white' | 'black' };
                }
            } catch (error) {
                console.error('Error checking player status with backend:', error);
            }
        }

        return { role: 'none' };
    }

    // Helper method to check player status with backend
    private checkPlayerStatusWithBackend(gameId: string, walletAddress: string): Promise<{ team: string; gameId: string }> {
        return new Promise((resolve, reject) => {
            // Set up listener for the response
            const unsubscribe = wsService.on('player_status', (data: any) => {
                if (data.gameId === gameId && data.walletAddress === walletAddress) {
                    const statusData = { team: data.team, gameId: data.gameId };
                    unsubscribe();
                    resolve(statusData);
                }
            });

            // Send the request
            wsService.checkPlayerStatus(gameId, walletAddress);

            // Set a timeout to reject if no response in 5 seconds
            setTimeout(() => {
                unsubscribe();
                reject(new Error('Timeout waiting for player status response'));
            }, 5000);
        });
    }

    // Get game state
    async getGameState(): Promise<GameInfo | null> {
        return this.currentGameState;
    }

    // ========== GAME ENGINE METHODS ==========

    // Initialize game and start managing it
    async initializeGame(gameId: string): Promise<void> {
        console.log(`Initializing game: ${gameId}`);
        this.currentGameId = gameId;

        // Clean up any existing game
        if (this.activeGameCleanups.has(gameId)) {
            this.activeGameCleanups.get(gameId)!();
        }

        // Join the game via WebSocket
        wsService.joinGame(gameId);

        // Create initial game state
        const initialState: GameInfo = {
            id: gameId,
            name: '',
            whitePlayer: undefined,
            blackPlayer: undefined,
            status: 'waiting',
            currentMove: 1,
            currentTurn: 'white',
            turnStartTime: Date.now(),
            turnTimeLimit: 10000,
            timeRemaining: 10000,
            spectators: 0,
            createdAt: new Date(),
            gameStartTime: undefined,
            whitePlayers: 0,
            blackPlayers: 0,
            proposedMoves: [],
            whiteCurrentTurnVotes: 0,
            blackCurrentTurnVotes: 0,
            whiteTeamTotalVotes: 0,
            blackTeamTotalVotes: 0,
            playerVotedThisRound: {},
            playerTotalVotes: {},
            totalPot: 0,
            whitePot: 0,
            blackPot: 0,
            boardState: this.getInitialBoardState()
        };

        this.currentGameState = initialState;

        // Set up WebSocket subscriptions for this game
        const unsubVotes = wsService.on('vote_update', (data) => {
            if (this.currentGameState && data.gameId === this.currentGameId) {
                this.currentGameState.proposedMoves = Object.entries(data.votes || {}).map(([move, votes]) => ({
                    moveId: move,
                    from: move.substring(0, 2),
                    to: move.substring(2, 4),
                    votes: votes as number
                }));

                // Update game statistics from backend
                this.updateGameStatsFromBackend(data);
                this.notifyGameStateUpdate();
            }
        });

        const unsubTimer = wsService.on('timer_tick', (data) => {
            if (this.currentGameState && data.gameId === this.currentGameId) {
                this.currentGameState.timeRemaining = data.secondsLeft * 1000;
                this.notifyGameStateUpdate();
            }
        });

        const unsubMove = wsService.on('move_result', (data) => {
            if (this.currentGameState && data.gameId === this.currentGameId) {
                console.log('Move result received for game:', data.gameId);

                // Clear proposed moves for the new turn
                this.currentGameState.proposedMoves = [];

                // Update all statistics from backend
                this.updateGameStatsFromBackend(data);

                // Update move and turn info
                if (data.currentTurn !== undefined) {
                    this.currentGameState.currentTurn = data.currentTurn as 'white' | 'black';
                }
                if (data.currentMove !== undefined) {
                    this.currentGameState.currentMove = data.currentMove;
                }

                // Reset timer for new turn
                this.currentGameState.timeRemaining = 10000;
                this.currentGameState.turnStartTime = Date.now();

                this.notifyGameStateUpdate();
            }
        });

        const unsubGameEnd = wsService.on('game_end', (data) => {
            if (data.gameId === this.currentGameId) {
                console.log('🏁 Game ended:', data);
                console.log('🏁 Player statistics received:', {
                    whiteTeamPlayers: data.whiteTeamPlayers,
                    blackTeamPlayers: data.blackTeamPlayers
                });
                console.log('🏁 Player statistics received:', {
                    whiteTeamPlayers: data.whiteTeamPlayers,
                    blackTeamPlayers: data.blackTeamPlayers
                });

                // Create game end info
                const gameEndInfo: GameEndInfo = {
                    gameId: data.gameId,
                    winner: data.winner as 'white' | 'black' | 'draw',
                    reason: data.gameEndReason,
                    whitePlayers: data.whitePlayers || 0,
                    blackPlayers: data.blackPlayers || 0,
                    whiteTeamTotalVotes: data.whiteTeamTotalVotes || 0,
                    blackTeamTotalVotes: data.blackTeamTotalVotes || 0,
                    totalPot: data.totalPot || 0,
                    whitePot: data.whitePot || 0,
                    blackPot: data.blackPot || 0,
                    currentMove: data.currentMove || 0,
                    playerVotes: data.playerVotes || 0,
                    whiteTeamPlayers: data.whiteTeamPlayers || [],
                    blackTeamPlayers: data.blackTeamPlayers || []
                };

                console.log('🏁 Triggering game end callbacks:', this.gameEndCallbacks.size);

                // Notify all game end callbacks
                this.gameEndCallbacks.forEach((callback, id) => {
                    try {
                        console.log(`🏁 Calling game end callback ${id}`);
                        callback(gameEndInfo);
                    } catch (error) {
                        console.error(`🏁 Error in game end callback ${id}:`, error);
                    }
                });
            }
        });

        // Store cleanup function
        this.activeGameCleanups.set(gameId, () => {
            unsubVotes();
            unsubTimer();
            unsubMove();
            unsubGameEnd();
        });
    }

    // Subscribe to game state updates
    onGameStateUpdate(callback: (state: GameInfo) => void): () => void {
        const callbackId = Math.random().toString(36);
        this.gameStateCallbacks.set(callbackId, callback);

        // Immediately call with current state if available
        if (this.currentGameState) {
            callback(this.currentGameState);
        }

        // Return unsubscribe function
        return () => {
            this.gameStateCallbacks.delete(callbackId);
        };
    }

    // Subscribe to game end events
    onGameEnd(callback: (gameEndInfo: GameEndInfo) => void): () => void {
        const callbackId = Math.random().toString(36);
        this.gameEndCallbacks.set(callbackId, callback);

        // Return unsubscribe function
        return () => {
            this.gameEndCallbacks.delete(callbackId);
        };
    }

    // Vote on a move (main game engine method)
    async voteOnMove(moveId: string, isConnected: boolean): Promise<VoteResult> {
        if (!isConnected || !this.isWalletConnected()) {
            return { success: false, error: 'Please connect your wallet to vote' };
        }

        if (!this.currentGameId || !this.currentGameState) {
            return { success: false, error: 'Game not ready for voting' };
        }

        const playerRole = this.playerRoles.get(this.currentGameId);
        if (this.currentGameState.currentTurn !== playerRole) {
            return { success: false, error: `It's not your team's turn. Currently ${this.currentGameState.currentTurn} team's turn.` };
        }

        if (this.currentGameState.timeRemaining <= 0) {
            return { success: false, error: 'Time is up - cannot vote now' };
        }

        // Check if player has already voted this round
        if (this.hasPlayerVoted(this.currentGameId)) {
            return { success: false, error: 'You have already voted this round' };
        }

        try {
            wsService.voteMove(this.currentGameId, moveId);
            return { success: true };
        } catch (error) {
            console.error('Voting error:', error);
            return { success: false, error: 'Error placing vote. Please try again.' };
        }
    }

    // Create a new move (main game engine method)
    async createMove(from: string, to: string, board: (ChessPiece | null)[][]): Promise<MoveResult> {
        console.log('🎮 CREATE MOVE DEBUG:', {
            currentGameId: this.currentGameId,
            hasCurrentGameState: !!this.currentGameState,
            from,
            to,
            walletConnected: this.isWalletConnected()
        });

        if (!this.isWalletConnected()) {
            return { success: false, error: 'Please connect your wallet to create moves' };
        }

        if (!this.currentGameId || !this.currentGameState) {
            console.log('❌ Game not ready for move creation');
            return { success: false, error: 'Game not ready for move creation' };
        }

        const playerRole = this.playerRoles.get(this.currentGameId);
        console.log('🎮 PLAYER ROLE CHECK:', {
            playerRole,
            currentTurn: this.currentGameState.currentTurn,
            isMyTurn: this.currentGameState.currentTurn === playerRole
        });

        if (this.currentGameState.currentTurn !== playerRole) {
            console.log('❌ Not your turn');
            return { success: false, error: 'Not your turn' };
        }

        console.log('🎮 TIME CHECK:', {
            timeRemaining: this.currentGameState.timeRemaining,
            isTimeUp: this.currentGameState.timeRemaining <= 0
        });

        if (this.currentGameState.timeRemaining <= 0) {
            console.log('❌ Time is up');
            return { success: false, error: 'Time is up - cannot propose moves now' };
        }

        // Check if player has already voted in this round
        const hasVoted = this.hasPlayerVoted(this.currentGameId);
        console.log('🎮 VOTING CHECK:', {
            hasVoted,
            playerIdentifier: this.getPlayerIdentifier(),
            playerVotedThisRound: this.currentGameState.playerVotedThisRound
        });

        if (hasVoted) {
            console.log('❌ Already voted this round');
            return { success: false, error: 'You have already voted this round' };
        }

        // Find the piece at the 'from' position
        const fromFile = from.charCodeAt(0) - 97; // a=0, b=1, etc
        const fromRank = 8 - parseInt(from[1]); // 8=0, 7=1, etc
        const piece = board[fromRank][fromFile];

        console.log('🎮 PIECE CHECK:', {
            fromFile,
            fromRank,
            piece,
            boardPosition: `board[${fromRank}][${fromFile}]`
        });

        if (!piece) {
            console.log('❌ No piece at position');
            return { success: false, error: 'No piece at position: ' + from };
        }

        // Store the pending move
        this.pendingMove = { from, to, piece };
        console.log('✅ Move creation successful, confirmation popup should appear');

        return { success: true, moveId: `${from}${to}` };
    }

    // Confirm pending move
    async confirmPendingMove(): Promise<MoveResult> {
        if (!this.pendingMove || !this.currentGameId) {
            return { success: false, error: 'No pending move to confirm' };
        }

        try {
            // Check if this move already exists in proposed moves
            const moveId = `${this.pendingMove.from}${this.pendingMove.to}`;
            const existingMove = this.currentGameState?.proposedMoves.find(
                move => `${move.from}${move.to}` === moveId
            );

            if (existingMove) {
                // Move already exists, just vote for it
                wsService.voteMove(this.currentGameId, moveId);
            } else {
                // Propose new move (this will automatically vote for it on the backend)
                wsService.voteMove(this.currentGameId, moveId);
            }

            this.pendingMove = null;
            return { success: true, moveId };
        } catch (error) {
            console.error('Error confirming move:', error);
            return { success: false, error: 'Error proposing move. Please try again.' };
        }
    }

    // Cancel pending move
    cancelPendingMove(): void {
        this.pendingMove = null;
    }

    // Get pending move
    getPendingMove(): PendingMove | null {
        return this.pendingMove;
    }

    // Check if voting is enabled for current player
    isVotingEnabled(): boolean {
        if (!this.currentGameState || !this.currentGameId || !this.isWalletConnected()) return false;

        const playerRole = this.playerRoles.get(this.currentGameId);
        return (
            this.currentGameState.currentTurn === playerRole &&
            this.currentGameState.timeRemaining > 0 &&
            !this.hasPlayerVoted(this.currentGameId)
        );
    }

    // Get current player role for the active game
    getCurrentPlayerRole(): 'white' | 'black' | 'spectator' | 'none' {
        if (!this.currentGameId) return 'none';
        return this.playerRoles.get(this.currentGameId) || 'none';
    }

    // Cleanup game
    cleanupGame(): void {
        if (this.currentGameId && this.activeGameCleanups.has(this.currentGameId)) {
            this.activeGameCleanups.get(this.currentGameId)!();
            this.activeGameCleanups.delete(this.currentGameId);
        }
        this.currentGameId = null;
        this.currentGameState = null;
        this.pendingMove = null;
        this.gameStateCallbacks.clear();
    }

    // ========== PRIVATE HELPER METHODS ==========

    private updateGameStatsFromBackend(data: any): void {
        if (!this.currentGameState) return;

        if (data.whitePlayers !== undefined) {
            this.currentGameState.whitePlayers = data.whitePlayers;
        }
        if (data.blackPlayers !== undefined) {
            this.currentGameState.blackPlayers = data.blackPlayers;
        }
        if (data.whiteCurrentTurnVotes !== undefined) {
            this.currentGameState.whiteCurrentTurnVotes = data.whiteCurrentTurnVotes;
        }
        if (data.blackCurrentTurnVotes !== undefined) {
            this.currentGameState.blackCurrentTurnVotes = data.blackCurrentTurnVotes;
        }
        if (data.whiteTeamTotalVotes !== undefined) {
            this.currentGameState.whiteTeamTotalVotes = data.whiteTeamTotalVotes;
        }
        if (data.blackTeamTotalVotes !== undefined) {
            this.currentGameState.blackTeamTotalVotes = data.blackTeamTotalVotes;
        }
        if (data.totalPot !== undefined) {
            this.currentGameState.totalPot = data.totalPot;
        }
        if (data.whitePot !== undefined) {
            this.currentGameState.whitePot = data.whitePot;
        }
        if (data.blackPot !== undefined) {
            this.currentGameState.blackPot = data.blackPot;
        }
        if (data.playerVotedThisRound !== undefined) {
            this.currentGameState.playerVotedThisRound = data.playerVotedThisRound;
        }
        if (data.playerTotalVotes !== undefined) {
            this.currentGameState.playerTotalVotes = data.playerTotalVotes;
        }
        if (data.board !== undefined) {
            this.currentGameState.boardState = this.convertBackendBoardToFrontend(data.board);
        }
    }

    private notifyGameStateUpdate(): void {
        if (this.currentGameState) {
            this.gameStateCallbacks.forEach(callback => callback(this.currentGameState!));
        }
    }

    // ========== PRIVATE HELPER METHODS ==========

    // Get initial board state
    private getInitialBoardState(): (ChessPiece | null)[][] {
        const board: (ChessPiece | null)[][] = Array(8).fill(null).map(() => Array(8).fill(null));

        // Place pieces in starting positions
        const pieceOrder: ('rook' | 'knight' | 'bishop' | 'queen' | 'king' | 'bishop' | 'knight' | 'rook')[] =
            ['rook', 'knight', 'bishop', 'queen', 'king', 'bishop', 'knight', 'rook'];

        // Black pieces
        for (let i = 0; i < 8; i++) {
            board[0][i] = { type: pieceOrder[i], color: 'black', position: String.fromCharCode(97 + i) + '8' };
            board[1][i] = { type: 'pawn', color: 'black', position: String.fromCharCode(97 + i) + '7' };
        }

        // White pieces
        for (let i = 0; i < 8; i++) {
            board[7][i] = { type: pieceOrder[i], color: 'white', position: String.fromCharCode(97 + i) + '1' };
            board[6][i] = { type: 'pawn', color: 'white', position: String.fromCharCode(97 + i) + '2' };
        }

        return board;
    }

    // Start heartbeat
    startHeartbeat() {
        // For now, simulate status updates
        setInterval(() => {
            wsService.requestGamesList();
        }, 5000);
    }

    // Check if the current player has voted in this round
    hasPlayerVoted(gameId: string): boolean {
        if (!this.currentGameState || !this.currentGameState.playerVotedThisRound || gameId !== this.currentGameId) {
            return false;
        }
        return this.currentGameState.playerVotedThisRound[this.getPlayerIdentifier()] || false;
    }

    // Get total votes by current player (across all rounds)
    getPlayerTotalVotes(gameId: string): number {
        if (!this.currentGameState || !this.currentGameState.playerTotalVotes || gameId !== this.currentGameId) {
            return 0;
        }

        return this.currentGameState.playerTotalVotes[this.getPlayerIdentifier()] || 0;
    }

    // Stop all services
    stop() {
        // Cleanup all active games
        this.activeGameCleanups.forEach(cleanup => cleanup());
        this.activeGameCleanups.clear();
        this.gameStateCallbacks.clear();

        wsService.disconnect();
    }

    getPlayerId(): string {
        return this.playerId;
    }

    private convertBackendBoardToFrontend(board: string[][]): (ChessPiece | null)[][] {
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
                };
            })
        );
    }
}

// Export singleton instance
export const gameService = new GameService(); 
