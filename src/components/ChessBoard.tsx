import React, { useState, useEffect } from 'react';
import { ChessPiece, GameInfo, pieceSymbols, coordsToSquare, squareToCoords } from '../utils/chess';

interface ChessBoardProps {
  gameState: GameInfo;
  highlightedSquares?: string[];
  onCreateMove?: (from: string, to: string) => void;
  isInteractive?: boolean;
  playerSide?: 'white' | 'black' | 'spectator' | 'none';
  hasVoted?: boolean;
}

const ChessBoard: React.FC<ChessBoardProps> = ({
  gameState,
  highlightedSquares = [],
  onCreateMove,
  isInteractive = false,
  playerSide = 'white',
  hasVoted = false
}) => {
  const [selectedSquare, setSelectedSquare] = useState<string | null>(null);
  const [possibleMoves, setPossibleMoves] = useState<string[]>([]);
  const [draggedPiece, setDraggedPiece] = useState<string | null>(null);

  // Get valid moves from game state
  const validMoves = gameState.validMoves || [];

  // Clear selection when game state changes (new turn, etc.)
  useEffect(() => {
    setSelectedSquare(null);
    setPossibleMoves([]);
  }, [gameState.currentTurn, gameState.currentMove]);

  const handleSquareClick = (square: string) => {
    if (!isInteractive) return;

    const piece = getPieceAtSquare(square);

    if (selectedSquare) {
      if (selectedSquare === square) {
        // Deselect
        setSelectedSquare(null);
        setPossibleMoves([]);
      } else if (possibleMoves.includes(square)) {
        // Make move
        onCreateMove?.(selectedSquare, square);
        setSelectedSquare(null);
        setPossibleMoves([]);
      } else if (piece && piece.color === gameState.currentTurn) {
        // Select new piece
        setSelectedSquare(square);
        // Filter valid moves for this piece
        const pieceMoves = validMoves.filter(move => move.startsWith(square));
        setPossibleMoves(pieceMoves.map(move => move.substring(2, 4)));
      } else {
        // Invalid move
        setSelectedSquare(null);
        setPossibleMoves([]);
      }
    } else if (piece && piece.color === gameState.currentTurn) {
      // Select piece
      setSelectedSquare(square);
      // Filter valid moves for this piece
      const pieceMoves = validMoves.filter(move => move.startsWith(square));
      const possibleDestinations = pieceMoves.map(move => move.substring(2, 4));

      setPossibleMoves(possibleDestinations);
    }
  };

  const getPieceAtSquare = (square: string): ChessPiece | null => {
    const [row, col] = squareToCoords(square);
    return gameState.boardState ? gameState.boardState[row][col] : null;
  };



  const canDragPiece = (piece: ChessPiece | null): boolean => {
    if (!isInteractive || !piece || hasVoted) return false;
    if (playerSide === 'spectator' || playerSide === 'none') return false;
    return piece.color === gameState.currentTurn && piece.color === playerSide;
  };

  const handleDragStart = (e: React.DragEvent, position: string) => {
    const piece = getPieceAtSquare(position);
    if (!canDragPiece(piece)) {
      e.preventDefault();
      return;
    }

    setDraggedPiece(position);
    setSelectedSquare(position);
    // Filter valid moves for this piece
    const pieceMoves = validMoves.filter(move => move.startsWith(position));
    const possibleDestinations = pieceMoves.map(move => move.substring(2, 4));

    setPossibleMoves(possibleDestinations);
    e.dataTransfer.effectAllowed = 'move';
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };

  const handleDrop = (e: React.DragEvent, targetSquare: string) => {
    e.preventDefault();

    if (!draggedPiece || draggedPiece === targetSquare) {
      setDraggedPiece(null);
      return;
    }

    if (possibleMoves.includes(targetSquare)) {
      onCreateMove?.(draggedPiece, targetSquare);
    }

    setDraggedPiece(null);
    setSelectedSquare(null);
    setPossibleMoves([]);
  };

  const renderSquare = (row: number, col: number) => {
    // Flip coordinates for black players
    const displayRow = playerSide === 'black' ? 7 - row : row;
    const displayCol = playerSide === 'black' ? 7 - col : col;

    const isLight = (displayRow + displayCol) % 2 === 0;
    const square = coordsToSquare(row, col);
    const piece = gameState.boardState ? gameState.boardState[row][col] : null;
    const isHighlighted = highlightedSquares.includes(square);
    const isSelected = selectedSquare === square;
    const isPossibleMove = possibleMoves.includes(square);
    const isLastMove = gameState.lastExecutedMove &&
      (gameState.lastExecutedMove.from === square || gameState.lastExecutedMove.to === square);
    const isValidMove = isPossibleMove && !piece;

    // Check if this square contains a king that is in check or checkmate
    // The king in check/checkmate is the one of the current turn (the player who needs to respond to the threat)
    const isKingInCheck = piece && piece.type === 'king' && piece.color === gameState.currentTurn && gameState.isInCheck && !gameState.isCheckmate;
    const isKingInCheckmate = piece && piece.type === 'king' && piece.color === gameState.currentTurn && gameState.isCheckmate;

    // Debug log for check/checkmate status (only for actual check/checkmate)
    if (isKingInCheck || isKingInCheckmate) {
      console.log(`King at ${square} (${piece?.color}) - Check: ${isKingInCheck}, Checkmate: ${isKingInCheckmate}, CurrentTurn: ${gameState.currentTurn}, GameCheck: ${gameState.isInCheck}, GameCheckmate: ${gameState.isCheckmate}`);
    }

    return (
      <div
        key={`${row}-${col}`}
        className={`
          relative aspect-square flex items-center justify-center text-4xl select-none
          ${isLight ? 'bg-amber-100' : 'bg-amber-700'}
          ${isSelected ? 'ring-4 ring-yellow-400' : ''}
          ${isValidMove ? 'ring-4 ring-green-400' : ''}
          ${isLastMove ? 'ring-4 ring-blue-400' : ''}
          ${isKingInCheckmate ? 'ring-4 ring-red-400' : isKingInCheck ? 'ring-4 ring-orange-400' : ''}
          ${canDragPiece(piece) ? 'cursor-move' : ''}
        `}
        onClick={() => handleSquareClick(square)}
        onDragStart={(e) => handleDragStart(e, square)}
        onDragOver={handleDragOver}
        onDrop={(e) => handleDrop(e, square)}
        draggable={canDragPiece(piece)}
      >
        {/* Square coordinate labels - same visual position for both players */}
        {displayCol === 0 && (
          <div className="absolute left-1 top-1 text-xs font-bold text-gray-600">
            {playerSide === 'black' ? row + 1 : 8 - row}
          </div>
        )}
        {displayRow === 7 && (
          <div className="absolute right-1 bottom-1 text-xs font-bold text-gray-600">
            {String.fromCharCode(97 + (playerSide === 'black' ? 7 - col : col))}
          </div>
        )}

        {/* Chess piece */}
        {piece && (
          <div className={`text-4xl md:text-5xl select-none hover:scale-110 transition-transform duration-150 ${piece.color === 'black' ? 'text-gray-800' : 'text-gray-900'}`}>
            {pieceSymbols[piece.color][piece.type]}
          </div>
        )}

        {/* Possible move indicator */}
        {isPossibleMove && !piece && (
          <div className="w-4 h-4 bg-purple-400 rounded-full opacity-60" />
        )}

        {/* Capture indicator */}
        {isPossibleMove && piece && (
          <div className="absolute inset-0 border-4 border-red-400 rounded-lg opacity-60" />
        )}
      </div>
    );
  };

  // Return empty board if no board state
  if (!gameState.boardState) {
    return (
      <div className="bg-gradient-to-br from-amber-50 to-amber-200 p-4 rounded-xl shadow-2xl">
        <div className="grid grid-cols-8 gap-0 border-4 border-amber-900 rounded-lg overflow-hidden">
          {Array(64).fill(null).map((_, index) => (
            <div key={index} className="aspect-square bg-amber-100" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-gradient-to-br from-amber-50 to-amber-200 p-4 rounded-xl shadow-2xl">
      <div className="grid grid-cols-8 gap-0 border-4 border-amber-900 rounded-lg overflow-hidden">
        {(playerSide === 'black'
          ? gameState.boardState.slice().reverse().map((row, reversedRowIndex) =>
            row.slice().reverse().map((_, reversedColIndex) => {
              const originalRow = 7 - reversedRowIndex;
              const originalCol = 7 - reversedColIndex;
              return renderSquare(originalRow, originalCol);
            })
          ).flat()
          : gameState.boardState.map((row, rowIndex) =>
            row.map((_, colIndex) => renderSquare(rowIndex, colIndex))
          ).flat()
        )}
      </div>

      {isInteractive && selectedSquare && (
        <div className="mt-2 text-center">
          <div className="inline-flex items-center px-3 py-1 bg-green-600 text-white rounded-full text-sm">
            Selected: {selectedSquare}
          </div>
        </div>
      )}
    </div>
  );
};

export default ChessBoard;
