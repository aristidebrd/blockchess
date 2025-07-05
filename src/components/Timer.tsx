import React, { useState, useEffect } from 'react';
import { Clock, AlertTriangle } from 'lucide-react';

interface TimerProps {
  duration: number; // in seconds
  onTimeUp: () => void;
  isActive: boolean;
  onReset?: () => void;
}

const Timer: React.FC<TimerProps> = ({ duration, onTimeUp, isActive, onReset }) => {
  const [timeLeft, setTimeLeft] = useState(duration);
  const [isWarning, setIsWarning] = useState(false);

  useEffect(() => {
    setTimeLeft(duration);
    setIsWarning(false);
  }, [duration, onReset]);

  useEffect(() => {
    if (!isActive) return;

    const interval = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          onTimeUp();
          return 0;
        }

        const newTime = prev - 1;
        setIsWarning(newTime <= 5);
        return newTime;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [isActive, onTimeUp]);

  const minutes = Math.floor(timeLeft / 60);
  const seconds = timeLeft % 60;
  const percentage = (timeLeft / duration) * 100;

  return (
    <div className="bg-gradient-to-br from-gray-800 to-gray-900 rounded-xl p-6 shadow-xl">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center space-x-2">
          {isWarning ? (
            <AlertTriangle className="w-6 h-6 text-red-400 animate-pulse" />
          ) : (
            <Clock className="w-6 h-6 text-blue-400" />
          )}
          <h3 className="text-lg font-semibold text-white">
            Voting Timer
          </h3>
        </div>
        <div className={`text-2xl font-mono font-bold ${isWarning ? 'text-red-400 animate-pulse' : 'text-white'
          }`}>
          {String(minutes).padStart(2, '0')}:{String(seconds).padStart(2, '0')}
        </div>
      </div>

      {/* Progress bar */}
      <div className="w-full bg-gray-700 rounded-full h-3 overflow-hidden">
        <div
          className={`h-full transition-all duration-1000 ease-linear ${isWarning
            ? 'bg-gradient-to-r from-red-500 to-red-600'
            : percentage > 50
              ? 'bg-gradient-to-r from-green-500 to-green-600'
              : percentage > 25
                ? 'bg-gradient-to-r from-yellow-500 to-yellow-600'
                : 'bg-gradient-to-r from-orange-500 to-red-500'
            }`}
          style={{ width: `${percentage}%` }}
        />
      </div>

      {/* Status text */}
      <div className="mt-3 text-center">
        <p className={`text-sm ${isWarning ? 'text-red-300' : 'text-gray-300'}`}>
          {isActive ? (
            isWarning ? 'Time running out!' : 'Voting in progress...'
          ) : (
            'Waiting for next round'
          )}
        </p>
      </div>
    </div>
  );
};

export default Timer;
