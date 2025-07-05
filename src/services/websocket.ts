// WebSocket service for real-time game communication
import { GameInfo } from '../utils/chess';

export interface VoteUpdate {
    type: 'vote_update';
    gameId: string;
    votes: Record<string, number>;
    board?: string[][];
    whitePlayers?: number;
    blackPlayers?: number;
    whiteVotes?: number;
    blackVotes?: number;
    whiteTeamTotalVotes?: number;
    blackTeamTotalVotes?: number;
    totalPot?: number;
    whitePot?: number;
    blackPot?: number;
    currentTurn?: string;
    playerVotedThisRound?: Record<string, boolean>;
    playerTotalVotes?: Record<string, number>;
}

export interface MoveResult {
    type: 'move_result';
    gameId: string;
    move: string;
    votes?: Record<string, number>;
    whitePlayers?: number;
    blackPlayers?: number;
    whiteVotes?: number;
    blackVotes?: number;
    whiteTeamTotalVotes?: number;
    blackTeamTotalVotes?: number;
    totalPot?: number;
    whitePot?: number;
    blackPot?: number;
    currentTurn?: string;
    currentMove?: number;
    playerVotedThisRound?: Record<string, boolean>;
    playerTotalVotes?: Record<string, number>;
    board?: string[][];
}

export interface TimerTick {
    type: 'timer_tick';
    gameId: string;
    secondsLeft: number;
}

export interface MatchFound {
    type: 'match_found';
    gameId: string;
    players: string[];
    assignedSide: string;
}

export interface GamesList {
    type: 'games_list';
    gamesList: GameInfo[];
}

export interface GamesListUpdate {
    type: 'games_list_update';
    gamesList: GameInfo[];
}

export interface TotalNumberOfPlayers {
    type: 'number_of_players';
    totalConnections: number;
}

export interface ClientConnected {
    type: 'client_connected';
    clientId: string;
}

export interface ErrorMessage {
    type: 'error';
    error: string;
}

export interface PlayerStatus {
    type: 'player_status';
    gameId: string;
    walletAddress: string;
    team: string; // 'white', 'black', or '' if not in game
}

export interface ValidMovesResponse {
    type: 'valid_moves_response';
    gameId: string;
    validMoves: string[];
}

export interface Permit2Data {
    type: 'permit2_data';
    permit2Data: any; // The EIP-712 typed data structure
}

export type ServerMessage = VoteUpdate | MoveResult | TimerTick | MatchFound | GamesList | GamesListUpdate | TotalNumberOfPlayers | ClientConnected | ErrorMessage | PlayerStatus | ValidMovesResponse | Permit2Data;

export interface ClientMessage {
    type: 'join_game' | 'vote_move' | 'join_team' | 'watch_game' | 'join_matchmaking' | 'leave_matchmaking' | 'request_games_list' | 'request_filtered_games_list' | 'check_player_status' | 'get_valid_moves' | 'request_permit2' | 'submit_permit2_signature';
    gameId?: string;
    move?: string;
    team?: 'white' | 'black';
    playerId?: string;
    walletAddress?: string;
    filter?: string;
    signature?: string;
    typedData?: any;
}

export class WebSocketService {
    private ws: WebSocket | null = null;
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectDelay = 1000;
    private listeners: Map<string, Set<(data: any) => void>> = new Map();
    private currentGameId: string | null = null;
    private playerIdGetter: (() => string) | null = null;

    constructor() {
        this.connect();
    }

    public setPlayerIdGetter(getter: () => string) {
        this.playerIdGetter = getter;
    }

    private connect() {
        try {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${protocol}//${window.location.host}/ws`;

            console.log('Connecting to WebSocket:', wsUrl);
            this.ws = new WebSocket(wsUrl);

            this.ws.onopen = () => {
                console.log('WebSocket connected');
                this.reconnectAttempts = 0;

                // Rejoin current game if any
                if (this.currentGameId) {
                    this.joinGame(this.currentGameId);
                }
            };

            this.ws.onmessage = (event) => {
                try {
                    const message: ServerMessage = JSON.parse(event.data);
                    this.handleMessage(message);
                } catch (error) {
                    console.error('Error parsing WebSocket message:', error);
                }
            };

            this.ws.onclose = () => {
                console.log('WebSocket disconnected');
                this.ws = null;
                this.scheduleReconnect();
            };

            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
            };
        } catch (error) {
            console.error('Failed to create WebSocket connection:', error);
            this.scheduleReconnect();
        }
    }

    private scheduleReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts);
            console.log(`Reconnecting in ${delay}ms...`);

            setTimeout(() => {
                this.reconnectAttempts++;
                this.connect();
            }, delay);
        }
    }

    private handleMessage(data: any) {
        console.log('📨 WebSocket message received:', data);
        console.log('📨 Message type:', data.type);

        // Special logging for matchmaking messages
        if (data.type === 'match_found') {
            console.log('🎯 MATCH_FOUND MESSAGE RECEIVED!');
            console.log('🎯 Game ID:', data.gameId);
            console.log('🎯 Players:', data.players);
            console.log('🎯 Assigned Side:', data.assignedSide);
        }

        // Special logging for client_connected messages
        if (data.type === 'client_connected') {
            console.log('🔗 CLIENT_CONNECTED MESSAGE RECEIVED!');
            console.log('🔗 Client ID:', data.clientId);
        }

        // Special logging for error messages
        if (data.type === 'error') {
            console.log('❌ ERROR MESSAGE RECEIVED!');
            console.log('❌ Error:', data.error);
        }

        // Use the existing listeners pattern
        const listeners = this.listeners.get(data.type);
        if (listeners) {
            console.log(`📨 Found ${listeners.size} listeners for type: ${data.type}`);
            listeners.forEach(listener => listener(data));
        } else {
            console.log(`📨 No listeners found for type: ${data.type}`);
        }
    }

    private send(message: ClientMessage) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        } else {
            console.error('WebSocket is not connected');
        }
    }

    public joinGame(gameId: string) {
        this.currentGameId = gameId;
        this.send({
            type: 'join_game',
            gameId
        });
    }

    public voteMove(gameId: string, move: string) {
        this.send({
            type: 'vote_move',
            gameId: gameId,
            move,
            playerId: this.playerIdGetter?.() || ''
        });
    }

    public joinTeam(gameId: string, team: 'white' | 'black') {
        this.send({
            type: 'join_team',
            gameId,
            team,
            playerId: this.playerIdGetter?.() || ''
        });
    }

    public watchGame(gameId: string) {
        this.send({ type: 'watch_game', gameId });
    }

    public joinMatchmaking(walletAddress: string) {
        this.send({
            type: 'join_matchmaking',
            walletAddress: walletAddress
        });
    }

    public leaveMatchmaking() {
        this.send({ type: 'leave_matchmaking' });
    }

    public requestGamesList() {
        this.send({ type: 'request_games_list' });
    }

    public requestFilteredGamesList(filter: 'active' | 'ended' | 'all') {
        this.send({ type: 'request_filtered_games_list', filter });
    }

    public checkPlayerStatus(gameId: string, walletAddress: string) {
        this.send({
            type: 'check_player_status',
            gameId,
            walletAddress
        });
    }

    // Request valid moves for a game
    getValidMoves(gameId: string): void {
        this.send({
            type: 'get_valid_moves',
            gameId: gameId
        });
    }

    public requestPermit2(walletAddress: string) {
        this.send({
            type: 'request_permit2',
            walletAddress
        });
    }

    public submitPermit2Signature(walletAddress: string, signature: string) {
        this.send({
            type: 'submit_permit2_signature',
            walletAddress,
            signature
        });
    }

    public on(type: string, callback: (data: any) => void): () => void {
        if (!this.listeners.has(type)) {
            this.listeners.set(type, new Set());
        }
        this.listeners.get(type)!.add(callback);

        // Return unsubscribe function
        return () => {
            const listeners = this.listeners.get(type);
            if (listeners) {
                listeners.delete(callback);
            }
        };
    }

    // Disconnect from WebSocket
    disconnect(): void {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }
}

// Export singleton instance
export const wsService = new WebSocketService(); 
