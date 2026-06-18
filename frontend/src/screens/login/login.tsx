import { useState } from 'react';
import { Eye, EyeOff, Mail, Lock, X, User } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';
import { ApiError } from '../../lib/api';

interface LoginModalProps {
  onClose: () => void;
}

export function LoginModal({ onClose }: LoginModalProps): JSX.Element {
  const { login } = useAuth();
  const [email, setEmail] = useState('');
  const [nik, setNik] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await login({ email: email.trim(), nik: nik.trim(), password });
      onClose();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.detail || 'Email, NIK, atau password salah. Coba lagi.');
      } else {
        setError('Terjadi kesalahan. Silakan coba lagi.');
      }
    } finally {
      setLoading(false);
    }
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

          {/* NIK */}
          <div>
            <label className="text-xs font-bold text-neutral-500 uppercase tracking-wide font-body block mb-1.5">
              NIK
            </label>
            <div className="relative">
              <User size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-neutral-400" />
              <input
                type="text"
                value={nik}
                onChange={(e) => setNik(e.target.value)}
                placeholder="16 digit NIK"
                required
                minLength={16}
                maxLength={16}
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
