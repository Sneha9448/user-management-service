import React, { useState } from 'react';
import { useQuery, useMutation } from 'urql';
import { GET_USERS_QUERY } from '../graphql/queries';
import { CREATE_USER_MUTATION, UPDATE_USER_MUTATION, DELETE_USER_MUTATION } from '../graphql/mutations';
import { useAuth } from '../context/AuthContext';
import { Layout } from '../components/Layout';
import {
    Plus,
    Edit2,
    Trash2,
    X,
    Save,
    Search,
    Loader2,
    ShieldAlert,
    MoreVertical,
    Mail
} from 'lucide-react';

export const Dashboard: React.FC = () => {
    const { user, logout } = useAuth();
    const [result, reexecuteQuery] = useQuery({ query: GET_USERS_QUERY });
    const { data, fetching, error } = result;

    const [, createUser] = useMutation(CREATE_USER_MUTATION);
    const [, updateUser] = useMutation(UPDATE_USER_MUTATION);
    const [, deleteUser] = useMutation(DELETE_USER_MUTATION);

    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingUser, setEditingUser] = useState<any>(null);
    const [formData, setFormData] = useState({ name: '', email: '' });
    const [searchTerm, setSearchTerm] = useState('');

    const handleOpenModal = (user?: any) => {
        if (user) {
            setEditingUser(user);
            setFormData({ name: user.name, email: user.email });
        } else {
            setEditingUser(null);
            setFormData({ name: '', email: '' });
        }
        setIsModalOpen(true);
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setEditingUser(null);
        setFormData({ name: '', email: '' });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (editingUser) {
            await updateUser({ id: editingUser.id, ...formData });
        } else {
            await createUser(formData);
        }
        handleCloseModal();
        reexecuteQuery({ requestPolicy: 'network-only' });
    };

    const handleDelete = async (id: string) => {
        if (window.confirm('Are you sure you want to delete this user?')) {
            await deleteUser({ id });
            reexecuteQuery({ requestPolicy: 'network-only' });
        }
    };

    const filteredUsers = (data?.users || []).filter((u: any) =>
        u.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        u.email.toLowerCase().includes(searchTerm.toLowerCase())
    );

    if (fetching) return (
        <div className="flex h-screen items-center justify-center bg-slate-50">
            <Loader2 className="animate-spin h-10 w-10 text-indigo-600" />
        </div>
    );

    const renderContent = () => {
        if (error) {
            return (
                <div className="flex flex-col items-center justify-center py-20 text-center animate-fade-in card p-10">
                    <div className="bg-red-50 p-6 rounded-full mb-6">
                        <ShieldAlert className="h-12 w-12 text-red-600" />
                    </div>
                    <h2 className="text-2xl font-black text-slate-900 mb-2">Access Denied</h2>
                    <p className="text-slate-500 max-w-sm mb-8">
                        {error.message.includes('access denied')
                            ? "Your account doesn't have the required permissions to manage teammates."
                            : error.message}
                    </p>
                    <button onClick={logout} className="btn btn-primary">
                        Back to Login
                    </button>
                </div>
            );
        }

        return (
            <div className="space-y-8">
                {/* Header */}
                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                    <div>
                        <h1 className="text-3xl font-black text-slate-900 tracking-tight">Organization Teammates</h1>
                        <p className="text-slate-500 font-medium">Manage user identities and security access.</p>
                    </div>
                    {user?.role === 'ADMIN' && (
                        <button onClick={() => handleOpenModal()} className="btn btn-primary gap-2 h-11 px-6 shadow-indigo-200 shadow-lg">
                            <Plus className="h-5 w-5" />
                            Add Teammate
                        </button>
                    )}
                </div>

                {/* Filters */}
                <div className="card p-6 flex flex-col md:flex-row gap-4 items-center">
                    <div className="search-wrapper flex-1 w-full">
                        <Search className="search-icon h-5 w-5" />
                        <input
                            type="text"
                            placeholder="Find teammates by name, email or role..."
                            className="input-modern search-input h-11"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                    </div>
                    <div className="flex gap-2 shrink-0">
                        <div className="badge badge-indigo">
                            <span className="font-black mr-1">{filteredUsers.length}</span> Members Total
                        </div>
                    </div>
                </div>

                {/* Table */}
                <div className="card">
                    <div className="table-wrapper">
                        <table>
                            <thead>
                                <tr>
                                    <th>User</th>
                                    <th>Role</th>
                                    <th>Status</th>
                                    <th className="text-right">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredUsers.map((u: any) => (
                                    <tr key={u.id}>
                                        <td>
                                            <div className="flex items-center gap-4">
                                                <div className="h-10 w-10 bg-slate-100 rounded-lg flex items-center justify-center font-black text-slate-400 uppercase">
                                                    {u.name.charAt(0)}
                                                </div>
                                                <div>
                                                    <div className="font-bold text-slate-900">{u.name}</div>
                                                    <div className="text-xs text-slate-500 flex items-center gap-1">
                                                        <Mail className="h-3 w-3" />
                                                        {u.email}
                                                    </div>
                                                </div>
                                            </div>
                                        </td>
                                        <td>
                                            <span className={`badge ${u.role === 'ADMIN' ? 'badge-indigo' : 'badge-success'} uppercase text-[10px]`}>
                                                {u.role || 'USER'}
                                            </span>
                                        </td>
                                        <td>
                                            <div className="flex items-center gap-1.5 font-bold text-emerald-600 text-xs uppercase tracking-wider">
                                                <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse"></span>
                                                Active
                                            </div>
                                        </td>
                                        <td className="text-right">
                                            {user?.role === 'ADMIN' ? (
                                                <div className="flex justify-end gap-1">
                                                    <button onClick={() => handleOpenModal(u)} className="btn btn-ghost p-2 rounded-lg">
                                                        <Edit2 className="h-4 w-4" />
                                                    </button>
                                                    <button onClick={() => handleDelete(u.id)} className="btn btn-ghost p-2 rounded-lg text-red-500 hover:bg-red-50">
                                                        <Trash2 className="h-4 w-4" />
                                                    </button>
                                                </div>
                                            ) : (
                                                <button className="btn btn-ghost p-2 text-slate-300 pointer-events-none">
                                                    <MoreVertical className="h-4 w-4" />
                                                </button>
                                            )}
                                        </td>
                                    </tr>
                                ))}
                                {filteredUsers.length === 0 && (
                                    <tr>
                                        <td colSpan={4} className="py-20 text-center">
                                            <Search className="h-10 w-10 text-slate-200 mx-auto mb-4" />
                                            <p className="text-slate-400 font-bold uppercase tracking-widest text-xs">No Results Found</p>
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        );
    };

    return (
        <Layout>
            {renderContent()}

            {/* Modal */}
            {isModalOpen && (
                <div className="modal-overlay">
                    <div className="card w-full max-w-lg p-10 relative animate-fade-in shadow-2xl" onClick={e => e.stopPropagation()}>
                        <button onClick={handleCloseModal} className="absolute top-6 right-6 text-slate-400 hover:text-slate-900">
                            <X className="h-6 w-6" />
                        </button>

                        <div className="mb-8">
                            <h2 className="text-2xl font-black text-slate-900 tracking-tight">{editingUser ? 'Edit Teammate' : 'Add New Member'}</h2>
                            <p className="text-slate-500">Configure profile details for organization access.</p>
                        </div>

                        <form onSubmit={handleSubmit} className="space-y-6">
                            <div>
                                <label className="label-modern">Full Name</label>
                                <input
                                    type="text"
                                    required
                                    value={formData.name}
                                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                    className="input-modern"
                                    placeholder="Jane Smith"
                                />
                            </div>

                            <div>
                                <label className="label-modern">Work Email</label>
                                <input
                                    type="email"
                                    required
                                    value={formData.email}
                                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                                    className="input-modern"
                                    placeholder="jane@organization.com"
                                />
                            </div>

                            <div className="flex justify-end gap-3 pt-4">
                                <button type="button" onClick={handleCloseModal} className="btn btn-ghost px-6">Cancel</button>
                                <button type="submit" className="btn btn-primary px-8">
                                    <Save className="h-4 w-4 mr-2" />
                                    {editingUser ? 'Save Changes' : 'Invite Member'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </Layout>
    );
};
