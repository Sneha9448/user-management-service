import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import {
    LayoutDashboard,
    Users,
    Settings,
    LogOut,
    Menu,
    X,
    Bell,
    ChevronDown,
    Shield
} from 'lucide-react';

interface LayoutProps {
    children: React.ReactNode;
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
    const { user, logout } = useAuth();
    const [isSidebarOpen, setIsSidebarOpen] = useState(true);

    const navItems = [
        { icon: LayoutDashboard, label: 'Dashboard', path: '/' },
        { icon: Users, label: 'Teammates', path: '/' },
        { icon: Shield, label: 'Roles', path: '#' },
        { icon: Settings, label: 'Settings', path: '#' },
    ];

    return (
        <div className="app-container">
            {/* Sidebar */}
            <aside className={`sidebar transition-all duration-300 ${isSidebarOpen ? 'w-[260px]' : 'w-20'}`}>
                <div className="p-6 flex items-center gap-3 border-bottom border-light mb-4">
                    <div className="bg-indigo-600 rounded-lg p-2">
                        <Shield className="text-white h-6 w-6" />
                    </div>
                    {isSidebarOpen && <span className="font-black text-lg tracking-tighter text-slate-900 uppercase">User Management</span>}
                </div>

                <nav className="flex-1 px-4 space-y-2">
                    {navItems.map((item) => (
                        <button
                            key={item.label}
                            className={`btn btn-ghost w-full !justify-start gap-4 ${isSidebarOpen ? '' : 'px-2'}`}
                        >
                            <item.icon className="h-5 w-5" />
                            {isSidebarOpen && <span>{item.label}</span>}
                        </button>
                    ))}
                </nav>

                <div className="p-4 border-t border-light">
                    <button
                        onClick={logout}
                        className={`btn btn-ghost w-full !justify-start gap-4 text-red-600 hover:bg-red-50 ${isSidebarOpen ? '' : 'px-2'}`}
                    >
                        <LogOut className="h-5 w-5" />
                        {isSidebarOpen && <span>Sign Out</span>}
                    </button>
                </div>
            </aside>

            {/* Main Content Area */}
            <div className="flex-1 flex flex-col min-h-screen">
                {/* Navbar */}
                <header className="h-16 bg-white border-b border-light flex items-center justify-between px-8 sticky top-0 z-40">
                    <button
                        onClick={() => setIsSidebarOpen(!isSidebarOpen)}
                        className="p-2 hover:bg-slate-100 rounded-lg transition-colors"
                    >
                        <Menu className="h-5 w-5 text-slate-500" />
                    </button>

                    <div className="flex items-center gap-6">
                        <button className="relative p-2 text-slate-400 hover:text-slate-600">
                            <Bell className="h-5 w-5" />
                            <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full border-2 border-white"></span>
                        </button>

                        <div className="h-8 w-px bg-slate-200"></div>

                        <div className="flex items-center gap-3 cursor-pointer group">
                            <div className="h-9 w-9 bg-indigo-100 rounded-full flex items-center justify-center font-bold text-indigo-700">
                                {user?.name.charAt(0)}
                            </div>
                            <div className="hidden sm:block">
                                <div className="text-sm font-bold text-slate-900">{user?.name}</div>
                                <div className="text-[10px] uppercase tracking-wider font-black text-slate-400">{user?.role}</div>
                            </div>
                            <ChevronDown className="h-4 w-4 text-slate-400 group-hover:text-slate-600" />
                        </div>
                    </div>
                </header>

                <main className="main-content animate-fade-in">
                    {children}
                </main>
            </div>
        </div>
    );
};
