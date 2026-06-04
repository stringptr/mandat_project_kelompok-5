import { useState } from 'react';
import { useLocation } from 'react-router-dom';
import { Search, Bell, ChevronDown } from 'lucide-react';
import type { Role } from '../App';

interface Notification {
  id: number;
  message: string;
  time: string;
  unread: boolean;
}

interface HeaderProps {
  currentRole: Role;
  onChangeRole: (role: Role) => void;
}

const PAGE_TITLES: Record<string, string> = {
  '/': 'Dashboard Overview',
  '/monitoring': 'Monitoring Gizi Ibu dan Anak',
  '/tindak-lanjut': 'Tindak Lanjut',
  '/edukasi': 'Edukasi',
  '/user-management': 'User Management',
  '/notifikasi': 'Notifikasi',
};

const ROLE_PROFILES = {
  'Ibu/Wali': {
    name: 'Ny. Rina Marlina',
    avatarUrl: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
  },
  'Bidan': {
    name: 'Bidan Sri Lestari',
    avatarUrl: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
  },
  'Dinas Kesehatan': {
    name: 'Dr. Budi Hermawan',
    avatarUrl: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
  },
  'Kader Posyandu': {
    name: 'Siti Aminah',
    avatarUrl: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
  },
};

export function Header({ currentRole, onChangeRole }: HeaderProps): JSX.Element {
  const location = useLocation();
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [showNotifications, setShowNotifications] = useState<boolean>(false);

  const pageTitle: string = PAGE_TITLES[location.pathname] || 'SiGizi';
  const profile = ROLE_PROFILES[currentRole];

  const notifications: Notification[] = [
    { id: 1, message: 'Data ibu hamil baru perlu verifikasi', time: '5 menit lalu', unread: true },
    { id: 2, message: 'Laporan bulanan telah diapprove', time: '1 jam lalu', unread: true },
    { id: 3, message: 'Jadwal edukasi besok pukul 09:00', time: '3 jam lalu', unread: false },
  ];

  const unreadCount: number = notifications.filter((n: Notification) => n.unread).length;

  return (
    <header className="flex items-center justify-between px-8 py-4 bg-white border-b border-neutral-100">
      <div className="flex items-center gap-4">
        <h1 className="text-xl font-bold text-neutral-800 font-headline">{pageTitle}</h1>
        <div className="flex items-center gap-3">
          <div className="h-5 w-px bg-neutral-200"></div>
          <div className="relative flex items-center">
            <select 
              value={currentRole}
              onChange={(e) => onChangeRole(e.target.value as Role)}
              className="appearance-none bg-white border border-primary text-primary text-xs font-semibold px-3 py-1.5 pr-8 rounded-lg focus:outline-none cursor-pointer"
            >
              <option value="Ibu/Wali">Role: Ibu/Wali</option>
              <option value="Bidan">Role: Bidan</option>
              <option value="Dinas Kesehatan">Role: Dinas Kesehatan</option>
              <option value="Kader Posyandu">Role: Kader Posyandu</option>
            </select>
            <span className="absolute right-2.5 pointer-events-none text-primary">
              <ChevronDown size={14} />
            </span>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-4">
        {/* Search */}
        {location.pathname !== '/monitoring' && (
          <div className="relative">
            <Search size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-400" />
            <input
              type="text"
              placeholder="Cari data wilayah..."
              value={searchQuery}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value)}
              className="w-64 pl-10 pr-4 py-2.5 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-700 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 font-body"
            />
          </div>
        )}

        {/* Notification */}
        <div className="relative">
          <button
            onClick={() => setShowNotifications(!showNotifications)}
            className="p-2 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-50 rounded-lg transition-colors relative"
          >
            <Bell size={20} />
            {unreadCount > 0 && (
              <span className="absolute top-1 right-1 w-4 h-4 bg-red-500 text-white text-[10px] font-bold rounded-full flex items-center justify-center">
                {unreadCount}
              </span>
            )}
          </button>

          {showNotifications && (
            <div className="absolute right-0 top-full mt-2 w-80 bg-white rounded-xl shadow-lg border border-neutral-100 py-2 z-50">
              <div className="px-4 py-2 border-b border-neutral-100">
                <span className="text-sm font-semibold text-neutral-800 font-headline">Notifikasi</span>
              </div>
              {notifications.map((notif: Notification) => (
                <div key={notif.id} className={`px-4 py-3 hover:bg-neutral-50 cursor-pointer ${notif.unread ? 'bg-primary-50/50' : ''}`}>
                  <p className="text-sm text-neutral-700 font-body">{notif.message}</p>
                  <p className="text-xs text-neutral-400 mt-1 font-body">{notif.time}</p>
                </div>
              ))}
              <div className="px-4 py-2 border-t border-neutral-100 text-center">
                <button className="text-xs text-primary font-medium hover:text-primary-600 font-body">Lihat Semua</button>
              </div>
            </div>
          )}
        </div>

        {/* User Avatar */}
        <button className="flex items-center gap-2 p-1 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-50 rounded-xl transition-colors">
          <div className="w-8 h-8 rounded-full overflow-hidden bg-neutral-200 flex items-center justify-center">
            <img 
              src={profile.avatarUrl} 
              alt={profile.name} 
              className="w-full h-full object-cover"
            />
          </div>
          <span className="text-sm font-semibold text-neutral-700 font-body hidden sm:inline">{profile.name}</span>
        </button>
      </div>
    </header>
  );
}