import { useState } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import {
    LayoutDashboard,
    BarChart3,
    ClipboardList,
    GraduationCap,
    Users,
    Bell,
    Plus,
    ChevronLeft,
    ChevronRight,
    Syringe,
} from 'lucide-react';import type { Role } from '../App';
import { useAuth } from '../context/AuthContext';

interface NavItem {
    path: string;
    label: string;
    icon: React.ReactNode;
}

interface SidebarProps {
    currentRole: Role;
    onLoginClick?: () => void;
}

const NAV_ITEMS: NavItem[] = [
    { path: '/', label: 'Dashboard', icon: <LayoutDashboard size={20} /> },
    { path: '/monitoring', label: 'Monitoring', icon: <BarChart3 size={20} /> },
    { path: '/tindak-lanjut', label: 'Tindak Lanjut', icon: <ClipboardList size={20} /> },
    { path: '/jadwal-imunisasi', label: 'Jadwal Imunisasi', icon: <Syringe size={20} /> },
    { path: '/edukasi', label: 'Edukasi', icon: <GraduationCap size={20} /> },
    { path: '/user-management', label: 'User Management', icon: <Users size={20} /> },
    { path: '/notifikasi', label: 'Notifikasi', icon: <Bell size={20} /> },
];

const ROLE_PROFILES = {
    'Ibu/Wali': {
        name: 'Ny. Rina Marlina',
        roleName: 'Ibu / Wali',
        subtitle: 'Orang Tua Arka',
        avatarUrl: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
    },
    'Bidan': {
        name: 'Bidan Sri Lestari',
        roleName: 'Bidan Desa',
        subtitle: 'Wilayah Binaan Melati',
        avatarUrl: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
    },
    'Dinas Kesehatan': {
        name: 'Dr. Budi Hermawan',
        roleName: 'Dinas Kesehatan',
        subtitle: 'Dinkes Jawa Tengah',
        avatarUrl: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
    },
    'Kader Posyandu': {
        name: 'Siti Aminah',
        roleName: 'Kader Posyandu',
        subtitle: 'Wilayah Posyandu Melati',
        avatarUrl: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
    },
};

const getFilteredNavItems = (role: Role) => {
    return NAV_ITEMS.filter((item) => {
        if (role === 'Ibu/Wali') {
            return ['/', '/monitoring', '/jadwal-imunisasi', '/edukasi', '/notifikasi'].includes(item.path);
        }
        if (role === 'Dinas Kesehatan') {
            return ['/', '/monitoring', '/jadwal-imunisasi', '/edukasi', '/user-management', '/notifikasi'].includes(item.path);
        }
        if (role === 'Kader Posyandu') {
            return ['/', '/monitoring', '/tindak-lanjut', '/jadwal-imunisasi', '/edukasi', '/notifikasi'].includes(item.path);
        }
        // Bidan sees all
        return true;
    });
};

export function Sidebar({ currentRole }: SidebarProps): JSX.Element {
    const location = useLocation();
    const [isCollapsed, setIsCollapsed] = useState<boolean>(false);
    const { isLoggedIn, user } = useAuth();

    const profile = isLoggedIn && user ? ROLE_PROFILES[user.role] : null;
    // Logged-in → filtered by role. Guest → only Dashboard & Edukasi
    const filteredItems = isLoggedIn && user ? getFilteredNavItems(user.role) : NAV_ITEMS.filter(item => ['/', '/edukasi'].includes(item.path));

    return (
        <aside className={`flex flex-col h-screen bg-white border-r border-neutral-100 transition-all duration-300 ${isCollapsed ? 'w-20' : 'w-64'}`}>
            {/* Logo */}
            <div className="p-6 flex items-center gap-3">
                <div className="w-12 h-12 flex-shrink-0">
                    <img src="/logo-sigizi.svg" alt="SiGizi" className="w-full h-full object-contain" />
                </div>
                {!isCollapsed && (
                    <div className="flex flex-col">
                        <span className="text-lg font-bold text-primary leading-tight font-headline">SiGizi</span>
                        <span className="text-xs text-primary-600 font-medium font-body">JAWA TENGAH</span>
                    </div>
                )}
            </div>

            {/* User Profile - hanya tampil saat login */}
            {isLoggedIn && profile && (
                <div className="px-4 mb-4">
                    <div className={`flex items-center gap-3 p-3 rounded-xl bg-neutral-50 ${isCollapsed ? 'justify-center' : ''}`}>
                        <div className="w-10 h-10 rounded-full overflow-hidden flex items-center justify-center flex-shrink-0 bg-neutral-200">
                            <img src={profile.avatarUrl} alt={profile.name} className="w-full h-full object-cover" />
                        </div>
                        {!isCollapsed && (
                            <div className="flex flex-col min-w-0">
                                <span className="text-sm font-semibold text-neutral-800 truncate font-body leading-tight">{profile.roleName}</span>
                                <span className="text-[11px] text-neutral-500 truncate font-body mt-0.5">{profile.subtitle}</span>
                            </div>
                        )}
                    </div>
                </div>
            )}

            {/* Navigation */}
            <nav className="flex-1 px-3 space-y-1 overflow-y-auto">
                {filteredItems.map((item: NavItem) => {
                    const isActive: boolean = location.pathname === item.path;
                    return (
                        <NavLink
                            key={item.path}
                            to={item.path}
                            className={`
                flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 font-body text-sm
                ${isActive ? 'bg-primary-50 text-primary font-semibold' : 'text-neutral-600 hover:bg-neutral-50 hover:text-neutral-800'}
                ${isCollapsed ? 'justify-center' : ''}
              `}
                        >
                            <span className={`flex-shrink-0 ${isActive ? 'text-primary' : 'text-neutral-400'}`}>
                                {item.icon}
                            </span>
                            {!isCollapsed && <span>{item.label}</span>}
                        </NavLink>
                    );
                })}
            </nav>

            {/* Bottom Actions */}
            <div className="p-4 space-y-3">
                {/* Logged-in: tambah data baru */}
                {isLoggedIn && currentRole !== 'Ibu/Wali' && currentRole !== 'Dinas Kesehatan' && (
                    <button className={`w-full flex items-center justify-center gap-2 bg-primary hover:bg-primary-600 text-white rounded-xl py-3 px-4 transition-colors font-body text-sm font-medium ${isCollapsed ? 'px-2' : ''}`}>
                        <Plus size={18} />
                        {!isCollapsed && <span>Tambah Data Baru</span>}
                    </button>
                )}

                <button
                    onClick={() => setIsCollapsed(!isCollapsed)}
                    className="w-full flex items-center justify-center p-2 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-50 rounded-lg transition-colors"
                >
                    {isCollapsed ? <ChevronRight size={20} /> : <ChevronLeft size={20} />}
                </button>
            </div>
        </aside>
    );
}