import React, { useState } from 'react';
import { useMutation } from 'urql';
import { useNavigate } from 'react-router-dom';
import { REQUEST_OTP_MUTATION, VERIFY_OTP_MUTATION } from '../graphql/mutations';
import { useAuth } from '../context/AuthContext';
import { Mail, Lock, ArrowRight, Loader2, ShieldCheck, User, ShieldAlert, ChevronLeft } from 'lucide-react';

type LoginStep = 'welcome' | 'email' | 'otp';

export const Login: React.FC = () => {
    const [email, setEmail] = useState('');
    const [otp, setOtp] = useState('');
    const [selectedRole, setSelectedRole] = useState<'USER' | 'ADMIN'>('USER');
    const [step, setStep] = useState<LoginStep>('welcome');
    const [error, setError] = useState('');

    const navigate = useNavigate();
    const { login, isAuthenticated } = useAuth();

    React.useEffect(() => {
        if (isAuthenticated) {
            navigate('/', { replace: true });
        }
    }, [isAuthenticated, navigate]);

    const [requestOtpResult, requestOtp] = useMutation(REQUEST_OTP_MUTATION);
    const [verifyOtpResult, verifyOtp] = useMutation(VERIFY_OTP_MUTATION);

    const handleSelectRole = (role: 'USER' | 'ADMIN') => {
        setSelectedRole(role);
        setStep('email');
    };

    const handleRequestOtp = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (!email) {
            setError('Please enter your email');
            return;
        }

        const result = await requestOtp({ email });
        if (result.error) {
            setError(result.error.message);
        } else {
            setStep('otp');
        }
    };

    const handleVerifyOtp = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (!otp) {
            setError('Please enter the OTP');
            return;
        }

        const result = await verifyOtp({ email, otp, role: selectedRole });
        if (result.error) {
            setError(result.error.message);
        } else if (result.data?.verifyOtp) {
            const { token, user } = result.data.verifyOtp;
            login(token, user);
            navigate('/');
        }
    };

    const isLoading = requestOtpResult.fetching || verifyOtpResult.fetching;

    return (
        <div className="flex items-center justify-center min-h-screen bg-slate-50 p-4">
            <div className="card w-full max-w-md p-10 animate-fade-in shadow-2xl bg-white border-slate-200">

                {/* Header Section */}
                <div className="text-center mb-10">
                    <div className="mx-auto w-16 h-16 bg-indigo-600 rounded-2xl flex items-center justify-center mb-6 shadow-indigo-100 shadow-xl">
                        <ShieldCheck className="h-8 w-8 text-white" />
                    </div>
                    <h1 className="text-3xl font-black text-slate-900 tracking-tight">
                        {step === 'welcome' ? 'User Management' : step === 'email' ? 'Identification' : 'Verify Access'}
                    </h1>
                    <p className="text-slate-500 font-medium mt-2">
                        {step === 'welcome'
                            ? 'Select your organization access point'
                            : step === 'email'
                                ? `Logging in as ${selectedRole.toLowerCase()}`
                                : `Authentication code sent to ${email}`}
                    </p>
                </div>

                {error && (
                    <div className="bg-red-50 border border-red-100 text-red-700 p-4 rounded-xl mb-8 text-sm flex items-center gap-3">
                        <ShieldAlert className="h-5 w-5 shrink-0 text-red-500" />
                        <span className="font-bold">{error}</span>
                    </div>
                )}

                {step === 'welcome' ? (
                    <div className="space-y-4">
                        <button
                            onClick={() => handleSelectRole('USER')}
                            className="flex items-center gap-4 p-5 rounded-2xl border border-slate-200 bg-white hover:border-indigo-600 hover:bg-slate-50 transition-all group w-full text-left"
                        >
                            <div className="w-12 h-12 rounded-xl bg-slate-100 flex items-center justify-center group-hover:bg-indigo-100 transition-colors">
                                <User className="w-6 h-6 text-slate-500 group-hover:text-indigo-600" />
                            </div>
                            <div className="flex-1">
                                <h3 className="font-bold text-slate-900">Team Member</h3>
                                <p className="text-xs text-slate-500">Standard organization access</p>
                            </div>
                            <ArrowRight className="w-4 h-4 text-slate-300 group-hover:text-indigo-600 group-hover:translate-x-1 transition-all" />
                        </button>

                        <button
                            onClick={() => handleSelectRole('ADMIN')}
                            className="flex items-center gap-4 p-5 rounded-2xl border border-slate-200 bg-white hover:border-indigo-600 hover:bg-slate-50 transition-all group w-full text-left"
                        >
                            <div className="w-12 h-12 rounded-xl bg-slate-100 flex items-center justify-center group-hover:bg-indigo-100 transition-colors">
                                <ShieldAlert className="w-6 h-6 text-slate-500 group-hover:text-indigo-600" />
                            </div>
                            <div className="flex-1">
                                <h3 className="font-bold text-slate-900">Administrator</h3>
                                <p className="text-xs text-slate-500">Full management permissions</p>
                            </div>
                            <ArrowRight className="w-4 h-4 text-slate-300 group-hover:text-indigo-600 group-hover:translate-x-1 transition-all" />
                        </button>
                    </div>
                ) : (
                    <form onSubmit={step === 'email' ? handleRequestOtp : handleVerifyOtp} className="space-y-6">
                        {step === 'email' ? (
                            <div>
                                <label className="label-modern">Work Email</label>
                                <div className="relative">
                                    <input
                                        type="email"
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                        placeholder="jane@organization.com"
                                        className="input-modern !pl-12 h-12"
                                        autoFocus
                                    />
                                    <Mail className="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
                                </div>
                            </div>
                        ) : (
                            <div>
                                <label className="label-modern">Authentication Code</label>
                                <div className="relative">
                                    <input
                                        type="text"
                                        value={otp}
                                        onChange={(e) => setOtp(e.target.value)}
                                        placeholder="000000"
                                        className="input-modern !pl-12 h-12 tracking-[0.5em] font-black text-center text-lg"
                                        maxLength={6}
                                        autoFocus
                                    />
                                    <Lock className="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
                                </div>
                                <p className="text-[10px] text-slate-400 mt-2 text-center uppercase tracking-widest font-black">Code Valid for 5 Minutes</p>
                            </div>
                        )}

                        <div className="flex flex-col gap-4">
                            <button
                                type="submit"
                                className="btn btn-primary w-full h-12 text-base gap-2"
                                disabled={isLoading}
                            >
                                {isLoading ? (
                                    <Loader2 className="h-5 w-5 animate-spin" />
                                ) : (
                                    <>
                                        {step === 'email' ? 'Continue' : 'Verify Identification'}
                                        <ArrowRight className="h-4 w-4" />
                                    </>
                                )}
                            </button>
                            <button
                                type="button"
                                onClick={() => { setStep('welcome'); setError(''); }}
                                className="btn btn-ghost text-slate-500 hover:text-slate-900 h-10 gap-2"
                            >
                                <ChevronLeft className="h-4 w-4" />
                                Change Role
                            </button>
                        </div>
                    </form>
                )}
            </div>

            <div className="fixed bottom-8 text-slate-400 text-xs font-bold uppercase tracking-widest">
                &copy; 2026 User Management System. v2.0
            </div>
        </div>
    );
};
