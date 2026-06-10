import { useState } from 'react';
import { Eye, EyeOff, Mail, Lock, X } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';
import type { Role } from '../../App';

// Mock credentials — satu akun per role
const MOCK_USERS: Record<string, { password: string; role: Role; name: string; avatarUrl: string }> = {
  'ibu@sigizi.id': {
    password: 'password',
    role: 'Ibu/Wali',
    name: 'Ny. Rina Marlina',
    avatarUrl: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
  },
  'bidan@sigizi.id': {
    password: 'password',
    role: 'Bidan',
    name: 'Bidan Sri Lestari',
    avatarUrl: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
  },
  'dinkes@sigizi.id': {
    password: 'password',
    role: 'Dinas Kesehatan',
    name: 'Dr. Budi Hermawan',
    avatarUrl: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
  },
  'kader@sigizi.id': {
    password: 'password',
    role: 'Kader Posyandu',
    name: 'Siti Aminah',
    avatarUrl: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
  },
};

interface LoginModalProps {
  onClose: () => void;
}

export function LoginModal({ onClose }: LoginModalProps): JSX.Element {
  const { login } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    setTimeout(() => {
      const user = MOCK_USERS[email.toLowerCase().trim()];
      if (user && user.password === password) {
        login({ name: user.name, role: user.role, avatarUrl: user.avatarUrl });
        onClose();
      } else {
        setError('Email atau password salah. Coba lagi.');
      }
      setLoading(false);
    }, 600);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" onClick={onClose} />

      {/* Modal card */}
      <div className="relative bg-white rounded-2xl shadow-2xl w-full max-w-sm p-8 z-10">
        {/* Close */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 p-1.5 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg transition-colors"
        >
          <X size={18} />
        </button>

        {/* Logo */}
        <div className="flex items-center gap-3 mb-6">
          <img src="/logo-sigizi.svg" alt="SiGizi" className="w-10 h-10 object-contain" />
          <div>
            <p className="text-base font-bold text-primary font-headline leading-tight">SiGizi</p>
            <p className="text-[10px] text-neutral-400 font-body">Jawa Tengah</p>
          </div>
        </div>

        <h1 className="text-xl font-bold text-neutral-900 font-headline mb-1">Selamat Datang Kembali</h1>
        <p className="text-sm text-neutral-500 font-body mb-6">Silakan masukkan akun Anda.</p>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Email */}
          <div>
            <label className="text-xs font-bold text-neutral-500 uppercase tracking-wide font-body block mb-1.5">
              Email
            </label>
            <div className="relative">
              <Mail size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-neutral-400" />
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="nama@kantor.com"
                required
                className="w-full pl-10 pr-4 py-3 bg-neutral-50 border border-neutral-200 rounded-xl text-sm font-body text-neutral-800 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 focus:border-primary transition-colors"
              />
            </div>
          </div>

          {/* Password */}
          <div>
            <label className="text-xs font-bold text-neutral-500 uppercase tracking-wide font-body block mb-1.5">
              Password
            </label>
            <div className="relative">
              <Lock size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-neutral-400" />
              <input
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                required
                className="w-full pl-10 pr-11 py-3 bg-neutral-50 border border-neutral-200 rounded-xl text-sm font-body text-neutral-800 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 focus:border-primary transition-colors"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3.5 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-neutral-600 transition-colors"
              >
                {showPassword ? <EyeOff size={15} /> : <Eye size={15} />}
              </button>
            </div>
          </div>

          {/* Error */}
          {error && (
            <p className="text-xs text-red-500 font-body bg-red-50 border border-red-100 rounded-lg px-3 py-2">
              {error}
            </p>
          )}

          {/* Submit */}
          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-primary hover:bg-primary-600 disabled:bg-primary/60 text-white rounded-xl text-sm font-bold font-body transition-colors mt-2 flex items-center justify-center gap-2"
          >
            {loading ? (
              <>
                <span className="w-4 h-4 border-2 border-white/40 border-t-white rounded-full animate-spin" />
                Memproses...
              </>
            ) : (
              'Masuk'
            )}
          </button>
        </form>

        <p className="text-center text-xs text-neutral-500 font-body mt-5">
          Belum memiliki akun?{' '}
          <button className="text-primary font-semibold hover:text-primary-600 transition-colors">
            Daftar
          </button>
        </p>
      </div>
    </div>
  );
}
