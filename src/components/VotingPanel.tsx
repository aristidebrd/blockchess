import React, { useState, useEffect } from 'react';
import { Vote, TrendingUp, Users, ChevronRight, Check, X } from 'lucide-react';
import { GameInfo, Move } from '../utils/chess';

interface PendingVote {
  moveId: string;
  move: Move;
  betAmount: string;
}

interface VotingPanelProps {
  moves: Move[];
  onVote: (moveId: string) => Promise<boolean>;
  isVotingEnabled: boolean;
  userVotes: string[];
  playerRole: 'white' | 'black';
  backendGameState: GameInfo;
}

const VotingPanel: React.FC<VotingPanelProps> = ({
  moves,
  onVote,
  isVotingEnabled,
  userVotes,
  playerRole,
  backendGameState,
}) => {
  const [pendingVote, setPendingVote] = useState<PendingVote | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);

  const totalVotes = moves.reduce((sum, move) => sum + move.votes, 0);
  const hasVotedThisRound = userVotes.length > 0;

  const handleSelectMove = (move: Move) => {
    if (!isVotingEnabled || hasVotedThisRound || isProcessing) return;

    const moveId = `${move.from}-${move.to}`;
    setPendingVote({
      moveId,
      move,
      betAmount: '0.01'
    });
  };

  const handleConfirmVote = async () => {
    if (!pendingVote) return;

    setIsProcessing(true);
    try {
      const success = await onVote(pendingVote.moveId);
      if (success) {
        setPendingVote(null);
      }
    } finally {
      setIsProcessing(false);
    }
  };

  const handleCancelVote = () => {
    setPendingVote(null);
  };

  const getMoveId = (move: Move) => `${move.from}-${move.to}`;

  // Automatically open confirm popup for newly created moves (with 0 votes)
  useEffect(() => {
    if (moves.length > 0 && !pendingVote && !hasVotedThisRound && isVotingEnabled) {
      // Find the most recently added move (should be the last one with 0 votes)
      const newMoves = moves.filter(move => move.votes === 0);
      if (newMoves.length > 0) {
        const latestMove = newMoves[newMoves.length - 1];
        const moveId = `${latestMove.from}-${latestMove.to}`;
        setPendingVote({
          moveId,
          move: latestMove,
          betAmount: '0.01'
        });
      }
    }
  }, [moves.length]); // Only depend on moves.length to avoid infinite loops

  return (
    <div className="bg-gradient-to-br from-gray-800 to-gray-900 rounded-xl p-6 shadow-xl">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-3">
          <Vote className="w-6 h-6 text-blue-400" />
          <h3 className="text-xl font-bold text-white">
            {playerRole === 'white' ? 'White' : 'Black'} Team - Vote for Next Move
          </h3>
        </div>
        <div className="flex items-center space-x-2 text-gray-300">
          <Users className="w-4 h-4" />
          <span className="text-sm">{playerRole === 'white' ? backendGameState.whiteCurrentTurnVotes : backendGameState.blackCurrentTurnVotes} votes</span>
        </div>
      </div>

      {/* Voting Status */}
      {hasVotedThisRound && (
        <div className="bg-green-900/30 border border-green-500/50 rounded-lg p-4 mb-6">
          <div className="flex items-center space-x-2 text-green-400">
            <Check className="w-5 h-5" />
            <span className="font-medium">You have voted and staked this round</span>
          </div>
          <p className="text-green-300 text-sm mt-1">
            Your 0.01 ETH stake is locked in the vault. Wait for the next round to vote again.
          </p>
        </div>
      )}

      {/* Pending Vote Confirmation */}
      {pendingVote && (
        <div className="bg-blue-900/30 border border-blue-500/50 rounded-lg p-4 mb-6">
          <h4 className="text-blue-400 font-semibold mb-3">Confirm Your Vote & Stake</h4>
          <div className="space-y-2 text-sm mb-4">
            <div className="flex justify-between">
              <span className="text-gray-300">Move:</span>
              <span className="text-white font-bold">{pendingVote.move.from} → {pendingVote.move.to}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-300">Stake Amount:</span>
              <span className="text-yellow-400 font-bold">{pendingVote.betAmount} ETH</span>
            </div>
          </div>
          <div className="bg-yellow-900/30 border border-yellow-500/50 rounded-lg p-3 mb-4">
            <div className="flex items-center space-x-2 text-yellow-400">
              <div className="w-4 h-4 rounded-full bg-yellow-400 flex items-center justify-center">
                <span className="text-xs text-black font-bold">!</span>
              </div>
              <span className="text-sm font-medium">Your stake will be locked in the vault until the game ends</span>
            </div>
          </div>
          <div className="flex space-x-3">
            <button
              onClick={handleConfirmVote}
              disabled={isProcessing}
              className="flex-1 bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded-lg transition-all duration-200 flex items-center justify-center space-x-2 disabled:opacity-50"
            >
              {isProcessing ? (
                <>
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  <span>Staking...</span>
                </>
              ) : (
                <>
                  <Check className="w-4 h-4" />
                  <span>Vote & Stake</span>
                </>
              )}
            </button>
            <button
              onClick={handleCancelVote}
              disabled={isProcessing}
              className="flex-1 bg-gray-600 hover:bg-gray-700 text-white font-medium py-2 px-4 rounded-lg transition-all duration-200 flex items-center justify-center space-x-2"
            >
              <X className="w-4 h-4" />
              <span>Cancel</span>
            </button>
          </div>
        </div>
      )}

      {/* Move options */}
      <div className="space-y-3">
        {moves.map((move) => {
          const moveId = getMoveId(move);
          const percentage = totalVotes > 0 ? (move.votes / totalVotes) * 100 : 0;
          const hasVoted = userVotes.includes(moveId);
          const isSelected = pendingVote?.moveId === moveId;

          return (
            <div
              key={moveId}
              className={`relative overflow-hidden rounded-lg border transition-all duration-200 cursor-pointer ${hasVoted
                ? 'border-green-500 bg-green-900/20'
                : isSelected
                  ? 'border-blue-500 bg-blue-900/20'
                  : hasVotedThisRound || !isVotingEnabled
                    ? 'border-gray-600 bg-gray-800/30 opacity-50 cursor-not-allowed'
                    : 'border-gray-600 bg-gray-800/50 hover:border-blue-500'
                }`}
              onClick={() => handleSelectMove(move)}
            >
              {/* Progress bar background */}
              <div
                className="absolute inset-0 bg-gradient-to-r from-blue-500/20 to-purple-500/20 transition-all duration-500"
                style={{ width: `${percentage}%` }}
              />

              <div className="relative p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <div className="text-2xl">
                      {move.piece.color === 'white' ? '♔' : '♚'}
                    </div>
                    <div>
                      <div className="font-bold text-white text-lg">
                        {move.notation}
                      </div>
                      <div className="text-sm text-gray-400">
                        {move.from} → {move.to}
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center space-x-4">
                    <div className="text-right">
                      <div className="text-lg font-bold text-white">
                        {move.votes}
                      </div>
                      <div className="text-xs text-gray-400">
                        {percentage.toFixed(1)}%
                      </div>
                    </div>

                    <div className="flex items-center space-x-2">
                      {hasVoted ? (
                        <div className="flex items-center space-x-2 text-green-400">
                          <Check className="w-4 h-4" />
                          <span className="text-sm font-medium">Voted</span>
                        </div>
                      ) : isSelected ? (
                        <div className="flex items-center space-x-2 text-blue-400">
                          <div className="w-2 h-2 bg-blue-400 rounded-full animate-pulse" />
                          <span className="text-sm font-medium">Selected</span>
                        </div>
                      ) : (
                        <ChevronRight className="w-4 h-4 text-gray-400" />
                      )}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default VotingPanel;
