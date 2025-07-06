import React, { useState } from 'react';
import { useAccount, useChainId } from 'wagmi';
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
    const chainId = useChainId();
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

        if (!chainId) {
            setError('Chain ID not available');
            return;
        }

        setIsLoading(true);
        setError(null);

        try {
            // Request permit signature data via websocket
            wsService.requestPermitSignature(address, chainId);

            // Listen for permit signature response
            const unsubscribe = wsService.on('permit_signature', (data: any) => {
                if (data.walletAddress === address && data.chainId === chainId) {
                    setTypedData(data.typedData);
                    setStep('sign-permit');
                    setIsLoading(false);
                    unsubscribe();
                }
            });

            // Listen for permit valid response (already have valid permit)
            const validUnsubscribe = wsService.on('permit_valid', (data: any) => {
                if (data.walletAddress === address && data.chainId === chainId) {
                    console.log('Player already has valid permit, skipping signature');
                    setStep('complete');
                    setIsLoading(false);
                    validUnsubscribe();
                    unsubscribe();

                    // Auto-complete after short delay
                    setTimeout(() => {
                        onApprovalComplete('existing_permit');
                    }, 1000);
                }
            });

            // Listen for errors
            const errorUnsubscribe = wsService.on('error', (data: any) => {
                setError(data.error);
                setIsLoading(false);
                errorUnsubscribe();
                validUnsubscribe();
                unsubscribe();
            });

            // Set timeout to prevent hanging
            setTimeout(() => {
                unsubscribe();
                validUnsubscribe();
                errorUnsubscribe();
                if (isLoading) {
                    setError('Timeout waiting for permit data');
                    setIsLoading(false);
                }
            }, 10000);

        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to generate permit data');
            setIsLoading(false);
        }
    };

    // Submit permit signature
    const submitPermit = async () => {
        if (!typedData || !address || !chainId) {
            setError('Required data not available');
            return;
        }

        setIsLoading(true);
        setError(null);

        try {
            // Request signature from user's wallet
            const signature = await window.ethereum.request({
                method: "eth_signTypedData_v4",
                params: [address, JSON.stringify(typedData)],
            });

            // Send signed permit back to backend
            wsService.submitPermitSignature(address, signature, chainId);

            setStep('complete');

            // Call completion callback
            setTimeout(() => {
                onApprovalComplete(signature);
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
                    <p className="text-gray-300">You can now participate in games without gas fees for votes!</p>
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
                    <h4 className="font-medium text-blue-300 mb-2">Gasless Game Participation</h4>
                    <p className="text-gray-300 text-sm">
                        Permit2 allows the vault contract to spend USDC on your behalf for game votes.
                        <br />
                        <strong>Benefits:</strong>
                        <br />
                        • No gas fees for future votes
                        <br />
                        • Seamless game participation
                        <br />
                        • Industry standard (Uniswap)
                        <br />
                        • Revocable anytime
                    </p>
                </div>

                <div className="bg-gray-700/30 rounded-lg p-3">
                    <div className="flex justify-between text-sm">
                        <span className="text-gray-300">Chain:</span>
                        <span className="text-blue-400 font-medium">{chainId}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                        <span className="text-gray-300">Amount:</span>
                        <span className="text-green-400 font-bold">1 USDC (for games)</span>
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
                            <h4 className="font-medium text-white mb-2">Sign Permit Authorization</h4>
                            <p className="text-gray-300 text-sm mb-4">
                                Sign the permit authorization to allow gasless game participation.
                                This signature will be stored securely for future game votes.
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
                                    <span>Signing...</span>
                                </>
                            ) : (
                                <>
                                    <DollarSign className="w-4 h-4" />
                                    <span>Sign Permit</span>
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
