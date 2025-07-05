// ========== TYPES ==========
export type PieceType = 'king' | 'queen' | 'rook' | 'bishop' | 'knight' | 'pawn';
export type PieceColor = 'white' | 'black';
export type Square = string; // e.g., 'a1', 'e4'

export interface ChessPiece {
  type: PieceType;
  color: PieceColor;
  position: Square;
  hasMoved?: boolean;
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
    { type: 'rook', color: 'black', position: 'a8' },
    { type: 'knight', color: 'black', position: 'b8' },
    { type: 'bishop', color: 'black', position: 'c8' },
    { type: 'queen', color: 'black', position: 'd8' },
    { type: 'king', color: 'black', position: 'e8' },
    { type: 'bishop', color: 'black', position: 'f8' },
    { type: 'knight', color: 'black', position: 'g8' },
    { type: 'rook', color: 'black', position: 'h8' }
  ],
  // Row 7 (black pawns)
  Array(8).fill(null).map((_, i) => ({
    type: 'pawn' as PieceType,
    color: 'black' as PieceColor,
    position: `${String.fromCharCode(97 + i)}7` as Square
  })),
  // Rows 6-3 (empty)
  ...Array(4).fill(null).map(() => Array(8).fill(null)),
  // Row 2 (white pawns)
  Array(8).fill(null).map((_, i) => ({
    type: 'pawn' as PieceType,
    color: 'white' as PieceColor,
    position: `${String.fromCharCode(97 + i)}2` as Square
  })),
  // Row 1 (white pieces)
  [
    { type: 'rook', color: 'white', position: 'a1' },
    { type: 'knight', color: 'white', position: 'b1' },
    { type: 'bishop', color: 'white', position: 'c1' },
    { type: 'queen', color: 'white', position: 'd1' },
    { type: 'king', color: 'white', position: 'e1' },
    { type: 'bishop', color: 'white', position: 'f1' },
    { type: 'knight', color: 'white', position: 'g1' },
    { type: 'rook', color: 'white', position: 'h1' }
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
  const file = String.fromCharCode(97 + col);
  const rank = (8 - row).toString();
  return `${file}${rank}` as Square;
}



// ========== MOVE CALCULATION ==========
export function calculatePossibleMoves(
  square: string,
  piece: ChessPiece,
  board: (ChessPiece | null)[][]
): string[] {
  const [row, col] = squareToCoords(square);
  const moves: string[] = [];

  const getPieceAtSquare = (square: string): ChessPiece | null => {
    const [r, c] = squareToCoords(square);
    return board[r][c];
  };

  switch (piece.type) {
    case 'pawn':
      const direction = piece.color === 'white' ? -1 : 1;
      const newRow = row + direction;
      if (newRow >= 0 && newRow < 8) {
        const targetSquare = coordsToSquare(newRow, col);
        if (!getPieceAtSquare(targetSquare)) {
          moves.push(targetSquare);

          // Double move from starting position
          if ((piece.color === 'white' && row === 6) || (piece.color === 'black' && row === 1)) {
            const doubleMove = coordsToSquare(newRow + direction, col);
            if (!getPieceAtSquare(doubleMove)) {
              moves.push(doubleMove);
            }
          }
        }

        // Captures
        for (const captureCol of [col - 1, col + 1]) {
          if (captureCol >= 0 && captureCol < 8) {
            const captureSquare = coordsToSquare(newRow, captureCol);
            const targetPiece = getPieceAtSquare(captureSquare);
            if (targetPiece && targetPiece.color !== piece.color) {
              moves.push(captureSquare);
            }
          }
        }
      }
      break;

    case 'rook':
      const rookDirections = [[0, 1], [0, -1], [1, 0], [-1, 0]];
      for (const [dr, dc] of rookDirections) {
        for (let distance = 1; distance < 8; distance++) {
          const newRow = row + (dr * distance);
          const newCol = col + (dc * distance);

          if (newRow < 0 || newRow >= 8 || newCol < 0 || newCol >= 8) break;

          const targetSquare = coordsToSquare(newRow, newCol);
          const targetPiece = getPieceAtSquare(targetSquare);

          if (!targetPiece) {
            moves.push(targetSquare);
          } else {
            if (targetPiece.color !== piece.color) {
              moves.push(targetSquare);
            }
            break;
          }
        }
      }
      break;

    case 'knight':
      const knightMoves = [[-2, -1], [-2, 1], [-1, -2], [-1, 2], [1, -2], [1, 2], [2, -1], [2, 1]];
      for (const [dr, dc] of knightMoves) {
        const newRow = row + dr;
        const newCol = col + dc;
        if (newRow >= 0 && newRow < 8 && newCol >= 0 && newCol < 8) {
          const targetSquare = coordsToSquare(newRow, newCol);
          const targetPiece = getPieceAtSquare(targetSquare);
          if (!targetPiece || targetPiece.color !== piece.color) {
            moves.push(targetSquare);
          }
        }
      }
      break;

    case 'bishop':
      const bishopDirections = [[1, 1], [1, -1], [-1, 1], [-1, -1]];
      for (const [dr, dc] of bishopDirections) {
        for (let distance = 1; distance < 8; distance++) {
          const newRow = row + (dr * distance);
          const newCol = col + (dc * distance);

          if (newRow < 0 || newRow >= 8 || newCol < 0 || newCol >= 8) break;

          const targetSquare = coordsToSquare(newRow, newCol);
          const targetPiece = getPieceAtSquare(targetSquare);

          if (!targetPiece) {
            moves.push(targetSquare);
          } else {
            if (targetPiece.color !== piece.color) {
              moves.push(targetSquare);
            }
            break;
          }
        }
      }
      break;

    case 'queen':
      const queenDirections = [[0, 1], [0, -1], [1, 0], [-1, 0], [1, 1], [1, -1], [-1, 1], [-1, -1]];
      for (const [dr, dc] of queenDirections) {
        for (let distance = 1; distance < 8; distance++) {
          const newRow = row + (dr * distance);
          const newCol = col + (dc * distance);

          if (newRow < 0 || newRow >= 8 || newCol < 0 || newCol >= 8) break;

          const targetSquare = coordsToSquare(newRow, newCol);
          const targetPiece = getPieceAtSquare(targetSquare);

          if (!targetPiece) {
            moves.push(targetSquare);
          } else {
            if (targetPiece.color !== piece.color) {
              moves.push(targetSquare);
            }
            break;
          }
        }
      }
      break;

    case 'king':
      const kingMoves = [
        [-1, -1], [-1, 0], [-1, 1],
        [0, -1], [0, 1],
        [1, -1], [1, 0], [1, 1]
      ];
      for (const [dr, dc] of kingMoves) {
        const newRow = row + dr;
        const newCol = col + dc;

        if (newRow >= 0 && newRow < 8 && newCol >= 0 && newCol < 8) {
          const targetSquare = coordsToSquare(newRow, newCol);
          const targetPiece = board[newRow][newCol];

          if (!targetPiece || targetPiece.color !== piece.color) {
            moves.push(targetSquare);
          }
        }
      }
      break;

    default:
      // For any unimplemented pieces, allow basic movement
      for (let r = Math.max(0, row - 2); r <= Math.min(7, row + 2); r++) {
        for (let c = Math.max(0, col - 2); c <= Math.min(7, col + 2); c++) {
          if (r !== row || c !== col) {
            const targetSquare = coordsToSquare(r, c);
            const targetPiece = getPieceAtSquare(targetSquare);
            if (!targetPiece || targetPiece.color !== piece.color) {
              moves.push(targetSquare);
            }
          }
        }
      }
  }

  return moves;
}
