import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Search, Bell, LogOut, LogIn } from 'lucide-react';
import type { Role } from '../App';
import { useAuth } from '../context/AuthContext';

interface HeaderProps {
  currentRole?: Role;
  onLoginClick: () => void;
}

const PAGE_TITLES: Record<string, string> = {
  '/': 'Dashboard',
  '/monitoring': 'Monitoring Gizi Ibu dan Anak',
  '/tindak-lanjut': 'Tindak Lanjut',
  '/jadwal-imunisasi': 'Jadwal Imunisasi',
  '/edukasi': 'Edukasi Hub',
  '/user-management': 'User Management',
  '/notifikasi': 'Notifikasi',
};

interface Notification {
  id: number;
  message: string;
  time: string;
  unread: boolean;
}

const NOTIFICATIONS: Notification[] = [
  { id: 1, message: 'Data ibu hamil baru perlu verifikasi', time: '5 menit lalu', unread: true },
  { id: 2, message: 'Laporan bulanan telah diapprove', time: '1 jam lalu', unread: true },
  { id: 3, message: 'Jadwal edukasi besok pukul 09:00', time: '3 jam lalu', unread: false },
];

export function Header({ currentRole, onLoginClick }: HeaderProps): JSX.Element {
  const location = useLocation();
  const navigate = useNavigate();
  const { user, isLoggedIn, logout } = useAuth();

  const [searchQuery, setSearchQuery] = useState('');
  const [showNotifications, setShowNotifications] = useState(false);
  const [showProfileMenu, setShowProfileMenu] = useState(false);

  const pageTitle = PAGE_TITLES[location.pathname] || 'SiGizi';
  const unreadCount = NOTIFICATIONS.filter((n) => n.unread).length;

  const handleLogout = () => {
    setShowProfileMenu(false);
    logout();
  };

  return (
    <header className="flex items-center justify-between px-8 py-4 bg-white border-b border-neutral-100 relative z-30">
      {/* Left: page title + role badge */}
      <div className="flex items-center gap-3">
        <h1 className="text-xl font-bold text-neutral-800 font-headline">{pageTitle}</h1>
        {isLoggedIn && currentRole && (
          <>
            <div className="h-5 w-px bg-neutral-200" />
            <span className="text-xs font-semibold text-primary border border-primary px-2.5 py-1 rounded-lg font-body">
              {currentRole}
            </span>
          </>
        )}
      </div>

      {/* Right */}
      <div className="flex items-center gap-3">
        {/* Search — hide on monitoring */}
        {isLoggedIn && location.pathname !== '/monitoring' && (
          <div className="relative">
            <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-400" />
            <input
              type="text"
              placeholder="Cari data..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-56 pl-9 pr-4 py-2.5 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-700 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 font-body"
            />
          </div>
        )}

        {isLoggedIn ? (
          <>
            {/* Notification bell */}
            <div className="relative">
              <button
                onClick={() => { setShowNotifications(!showNotifications); setShowProfileMenu(false); }}
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
                  {NOTIFICATIONS.map((n) => (
                    <div key={n.id} className={`px-4 py-3 hover:bg-neutral-50 cursor-pointer ${n.unread ? 'bg-primary-50/50' : ''}`}>
                      <p className="text-sm text-neutral-700 font-body">{n.message}</p>
                      <p className="text-xs text-neutral-400 mt-1 font-body">{n.time}</p>
                    </div>
                  ))}
                  <div className="px-4 py-2 border-t border-neutral-100 text-center">
                    <button
                      onClick={() => { navigate('/notifikasi'); setShowNotifications(false); }}
                      className="text-xs text-primary font-medium hover:text-primary-600 font-body"
                    >
                      Lihat Semua
                    </button>
                  </div>
                </div>
              )}
            </div>

            {/* Profile avatar + dropdown */}
            <div className="relative">
              <button
                onClick={() => { setShowProfileMenu(!showProfileMenu); setShowNotifications(false); }}
                className="flex items-center gap-2 p-1 hover:bg-neutral-50 rounded-xl transition-colors"
              >
                <div className="w-8 h-8 rounded-full overflow-hidden bg-neutral-200 flex-shrink-0">
                  <img src={user!.avatarUrl} alt={user!.name} className="w-full h-full object-cover" />
                </div>
                <span className="text-sm font-semibold text-neutral-700 font-body hidden sm:inline">
                  {user!.name}
                </span>
              </button>

              {showProfileMenu && (
                <div className="absolute right-0 top-full mt-2 w-56 bg-white rounded-xl shadow-lg border border-neutral-100 py-2 z-50">
                  {/* User info */}
                  <div className="px-4 py-3 border-b border-neutral-100">
                    <p className="text-sm font-semibold text-neutral-800 font-body">{user!.name}</p>
                    <p className="text-xs text-neutral-500 font-body mt-0.5">{user!.role}</p>
                  </div>
                  {/* Logout */}
                  <button
                    onClick={handleLogout}
                    className="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 transition-colors font-body"
                  >
                    <LogOut size={15} />
                    Keluar
                  </button>
                </div>
              )}
            </div>
          </>
        ) : (
          /* Not logged in → Login button */
          <button
            onClick={onLoginClick}
            className="flex items-center gap-2 bg-primary hover:bg-primary-600 text-white px-4 py-2 rounded-xl text-sm font-semibold font-body transition-colors"
          >
            <LogIn size={16} />
            Masuk
          </button>
        )}
      </div>

      {/* Click-outside dismiss overlay */}
      {(showNotifications || showProfileMenu) && (
        <div
          className="fixed inset-0 z-40"
          onClick={() => { setShowNotifications(false); setShowProfileMenu(false); }}
        />
      )}
    </header>
  );
}
