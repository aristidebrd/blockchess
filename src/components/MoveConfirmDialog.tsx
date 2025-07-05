import React from 'react';
import { X, CheckCircle, XCircle } from 'lucide-react';
import { ChessPiece } from '../utils/chess';

interface MoveConfirmDialogProps {
    isOpen: boolean;
    from: string;
    to: string;
    piece: ChessPiece;
    onConfirm: () => void;
    onCancel: () => void;
}

const MoveConfirmDialog: React.FC<MoveConfirmDialogProps> = ({
    isOpen,
    from,
    to,
    piece,
    onConfirm,
    onCancel
}) => {
    if (!isOpen) return null;

    const getPieceSymbol = (piece: ChessPiece) => {
        const symbols = {
            king: piece.color === 'white' ? '♔' : '♚',
            queen: piece.color === 'white' ? '♕' : '♛',
            rook: piece.color === 'white' ? '♖' : '♜',
            bishop: piece.color === 'white' ? '♗' : '♝',
            knight: piece.color === 'white' ? '♘' : '♞',
            pawn: piece.color === 'white' ? '♙' : '♟',
        };
        return symbols[piece.type];
    };

    const getPieceName = (type: string) => {
        return type.charAt(0).toUpperCase() + type.slice(1);
    };

    return (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
            <div className="bg-gradient-to-br from-gray-800 to-gray-900 rounded-2xl p-6 max-w-md w-full border border-gray-700 shadow-2xl animate-fade-in">
                {/* Header */}
                <div className="flex items-center justify-between mb-6">
                    <h3 className="text-2xl font-bold text-white">Confirm Move</h3>
                    <button
                        onClick={onCancel}
                        className="text-gray-400 hover:text-white transition-colors p-1 hover:bg-gray-700 rounded-lg"
                    >
                        <X className="w-6 h-6" />
                    </button>
                </div>

                {/* Move Details */}
                <div className="bg-gradient-to-r from-blue-900/30 to-purple-900/30 rounded-xl p-6 mb-6 border border-blue-500/20">
                    <div className="flex items-center justify-center space-x-4 mb-4">
                        {/* From Square */}
                        <div className="text-center">
                            <div className="text-6xl mb-2">{getPieceSymbol(piece)}</div>
                            <div className="text-white font-bold text-xl uppercase">{from}</div>
                        </div>

                        {/* Arrow */}
                        <div className="text-white text-3xl">→</div>

                        {/* To Square */}
                        <div className="text-center">
                            <div className="w-16 h-16 bg-gray-700/50 rounded-lg border-2 border-dashed border-gray-500 mb-2" />
                            <div className="text-white font-bold text-xl uppercase">{to}</div>
                        </div>
                    </div>

                    <div className="text-center text-gray-300">
                        Move your <span className="text-white font-semibold">{piece.color} {getPieceName(piece.type)}</span> from <span className="text-yellow-400 font-bold">{from.toUpperCase()}</span> to <span className="text-yellow-400 font-bold">{to.toUpperCase()}</span>
                    </div>
                </div>

                {/* Bet Info */}
                <div className="bg-gradient-to-r from-yellow-900/30 to-orange-900/30 rounded-xl p-4 mb-6 border border-yellow-500/20">
                    <div className="flex items-center justify-between">
                        <span className="text-gray-300">Vote Cost:</span>
                        <span className="text-yellow-400 font-bold text-lg">0.01 ETH</span>
                    </div>
                    <p className="text-xs text-gray-400 mt-2">
                        This will place a bet on your proposed move. If it wins, you'll share the pot!
                    </p>
                </div>

                {/* Action Buttons */}
                <div className="grid grid-cols-2 gap-4">
                    <button
                        onClick={onCancel}
                        className="bg-gradient-to-r from-gray-600 to-gray-700 hover:from-gray-500 hover:to-gray-600 text-white font-bold py-3 px-6 rounded-xl transition-all duration-200 flex items-center justify-center space-x-2 shadow-lg hover:shadow-xl transform hover:scale-105"
                    >
                        <XCircle className="w-5 h-5" />
                        <span>Cancel</span>
                    </button>
                    <button
                        onClick={onConfirm}
                        className="bg-gradient-to-r from-green-600 to-green-700 hover:from-green-500 hover:to-green-600 text-white font-bold py-3 px-6 rounded-xl transition-all duration-200 flex items-center justify-center space-x-2 shadow-lg hover:shadow-xl transform hover:scale-105"
                    >
                        <CheckCircle className="w-5 h-5" />
                        <span>Confirm</span>
                    </button>
                </div>

                {/* Note */}
                <p className="text-xs text-gray-500 text-center mt-4">
                    Your move will be added to the voting pool. The most voted move will be executed when the timer runs out.
                </p>
            </div>
        </div>
    );
};

export default MoveConfirmDialog; 
