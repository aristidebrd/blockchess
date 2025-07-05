// ========== TYPES ==========
export type PieceType = 'king' | 'queen' | 'rook' | 'bishop' | 'knight' | 'pawn';
export type PieceColor = 'white' | 'black';
export type Square = string; // e.g., 'a1', 'e4'

export interface ChessPiece {
  type: PieceType;
  color: PieceColor;
  position: Square;
  hasMoved?: boolean;
  hasBeenChecked?: boolean; // For kings: true if this king has ever been in check during the game
}

export interface Move {
  from: Square;
  to: Square;
  piece: ChessPiece;
  notation: string;
  votes: number;
}

export interface GameState {
  board: (ChessPiece | null)[][];
  currentPlayer: PieceColor;
  moveNumber: number;
  gameStatus: 'playing' | 'ended';
  lastMove?: Move;
  gameStartTime?: Date;
}



// Consolidated game information interface
export interface GameInfo {
  // Basic game metadata
  id: string;
  name?: string;
  status: 'waiting' | 'active' | 'completed' | 'ended';
  createdAt?: Date;
  gameStartTime?: Date;

  // Game state
  currentMove: number;
  currentTurn: 'white' | 'black';
  turnStartTime?: number;
  turnTimeLimit?: number;
  timeRemaining: number;
  boardState?: (ChessPiece | null)[][];
  lastExecutedMove?: { from: string, to: string, piece: ChessPiece };
  validMoves?: string[]; // Valid moves for the current position from backend

  // Players and teams
  whitePlayers: number;
  blackPlayers: number;
  spectators: number;
  whitePlayer?: string;
  blackPlayer?: string;

  // Voting system
  proposedMoves: { moveId: string, from: string, to: string, votes: number }[];
  whiteCurrentTurnVotes: number;
  blackCurrentTurnVotes: number;
  whiteTeamTotalVotes: number;
  blackTeamTotalVotes: number;
  playerVotedThisRound: Record<string, boolean>;
  playerTotalVotes: Record<string, number>;

  // Economic system
  totalPot: number;
  whitePot: number;
  blackPot: number;

  // Game end information (for ended games)
  winner?: 'white' | 'black' | 'draw';
  endReason?: string;
  endedAt?: number; // Unix timestamp

  // Player statistics per team (only for ended games)
  whiteTeamPlayers?: PlayerStats[];
  blackTeamPlayers?: PlayerStats[];

  // Check and checkmate status
  isInCheck?: boolean;
  isCheckmate?: boolean;
}





export interface PlayerStats {
  walletAddress: string;
  totalVotes: number;
  totalSpent: number;
}

export interface GameEndInfo {
  gameId: string;
  winner: 'white' | 'black' | 'draw';
  reason: string;
  whitePlayers: number;
  blackPlayers: number;
  whiteTeamTotalVotes: number;
  blackTeamTotalVotes: number;
  totalPot: number;
  whitePot: number;
  blackPot: number;
  currentMove: number;
  playerVotes: number;
  whiteTeamPlayers?: PlayerStats[];
  blackTeamPlayers?: PlayerStats[];
}

// ========== CONSTANTS ==========
export const initialBoard: (ChessPiece | null)[][] = [
  // Row 8 (black pieces)
  [
    { type: 'rook', color: 'black', position: 'a8', hasMoved: false },
    { type: 'knight', color: 'black', position: 'b8', hasMoved: false },
    { type: 'bishop', color: 'black', position: 'c8', hasMoved: false },
    { type: 'queen', color: 'black', position: 'd8', hasMoved: false },
    { type: 'king', color: 'black', position: 'e8', hasMoved: false, hasBeenChecked: false },
    { type: 'bishop', color: 'black', position: 'f8', hasMoved: false },
    { type: 'knight', color: 'black', position: 'g8', hasMoved: false },
    { type: 'rook', color: 'black', position: 'h8', hasMoved: false }
  ],
  // Row 7 (black pawns)
  Array(8).fill(null).map((_, i) => ({
    type: 'pawn' as PieceType,
    color: 'black' as PieceColor,
    position: `${String.fromCharCode(97 + i)}7` as Square,
    hasMoved: false
  })),
  // Rows 6-3 (empty)
  ...Array(4).fill(null).map(() => Array(8).fill(null)),
  // Row 2 (white pawns)
  Array(8).fill(null).map((_, i) => ({
    type: 'pawn' as PieceType,
    color: 'white' as PieceColor,
    position: `${String.fromCharCode(97 + i)}2` as Square,
    hasMoved: false
  })),
  // Row 1 (white pieces)
  [
    { type: 'rook', color: 'white', position: 'a1', hasMoved: false },
    { type: 'knight', color: 'white', position: 'b1', hasMoved: false },
    { type: 'bishop', color: 'white', position: 'c1', hasMoved: false },
    { type: 'queen', color: 'white', position: 'd1', hasMoved: false },
    { type: 'king', color: 'white', position: 'e1', hasMoved: false, hasBeenChecked: false },
    { type: 'bishop', color: 'white', position: 'f1', hasMoved: false },
    { type: 'knight', color: 'white', position: 'g1', hasMoved: false },
    { type: 'rook', color: 'white', position: 'h1', hasMoved: false }
  ]
];

export const pieceSymbols: Record<PieceColor, Record<PieceType, string>> = {
  white: {
    king: '♔',
    queen: '♕',
    rook: '♖',
    bishop: '♗',
    knight: '♘',
    pawn: '♙'
  },
  black: {
    king: '♚',
    queen: '♛',
    rook: '♜',
    bishop: '♝',
    knight: '♞',
    pawn: '♟︎'
  }
};

// ========== UTILITY FUNCTIONS ==========
export function squareToCoords(square: Square): [number, number] {
  const file = square.charCodeAt(0) - 97; // a=0, b=1, etc.
  const rank = 8 - parseInt(square[1]); // 8=0, 7=1, etc. (flipped for array indexing)
  return [rank, file];
}

export function coordsToSquare(row: number, col: number): Square {
  const file = String.fromCharCode(97 + col); // a, b, c, etc.
  const rank = (8 - row).toString(); // 8, 7, 6, etc.
  return `${file}${rank}` as Square;
}




