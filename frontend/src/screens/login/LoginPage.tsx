import { useState } from 'react';
import { Eye, EyeOff, Mail, Lock, User } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { useNotification } from '../../context/NotificationContext';
import { ApiError } from '../../lib/api';

export default function LoginPage(): JSX.Element {
  const { login } = useAuth();
  const notify = useNotification();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [nik, setNik] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!email.trim() || !nik.trim() || !password.trim()) {
      notify.warn('Mohon lengkapi semua data form yang wajib diisi sebelum mengirim.');
      return;
    }

    setLoading(true);
    try {
      await login({ email: email.trim(), nik: nik.trim(), password });
      notify.success('Selamat datang kembali!');
      navigate('/');
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.detail || 'Email, NIK, atau password salah. Coba lagi.');
        notify.apiError(err, 'Email, NIK, atau password salah. Silakan coba lagi.');
      } else {
        setError('Terjadi kesalahan. Silakan coba lagi.');
        notify.error('Terjadi kesalahan. Silakan coba lagi.');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-primary-50 via-white to-primary-100 px-4 font-body">
      {/* Logo */}
      <div className="flex flex-col items-center mb-8">
        <img src="/logo-sigizi.svg" alt="SiGizi" className="w-16 h-16 object-contain mb-3" />
        <p className="text-xl font-bold text-primary font-headline">SiGizi</p>
        <p className="text-xs text-neutral-500 mt-0.5 text-center leading-snug">
          Sistem Monitoring Gizi Ibu dan Anak<br />Berbasis Komunitas
        </p>
      </div>

      {/* Card */}
      <div className="bg-white rounded-2xl shadow-xl border border-neutral-100 w-full max-w-sm p-8">
        <h1 className="text-xl font-bold text-neutral-900 font-headline mb-1">Selamat Datang Kembali</h1>
        <p className="text-sm text-neutral-500 mb-6">Silakan masuk ke akun Anda.</p>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="text-xs font-bold text-neutral-500 uppercase tracking-wide block mb-1.5">Email</label>
            <div className="relative">
              <Mail size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-neutral-400" />
              <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="nama@email.com" required
                className="w-full pl-10 pr-4 py-3 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-800 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 focus:border-primary transition-colors" />
            </div>
          </div>

          <div>
            <label className="text-xs font-bold text-neutral-500 uppercase tracking-wide block mb-1.5">NIK</label>
            <div className="relative">
              <User size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-neutral-400" />
              <input type="text" value={nik} onChange={(e) => setNik(e.target.value)} placeholder="16 digit NIK" required minLength={16} maxLength={16}
                className="w-full pl-10 pr-4 py-3 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-800 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 focus:border-primary transition-colors" />
            </div>
          </div>

          <div>
            <label className="text-xs font-bold text-neutral-500 uppercase tracking-wide block mb-1.5">Password</label>
            <div className="relative">
              <Lock size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-neutral-400" />
              <input type={showPassword ? 'text' : 'password'} value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••" required
                className="w-full pl-10 pr-11 py-3 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-800 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 focus:border-primary transition-colors" />
              <button type="button" onClick={() => setShowPassword(!showPassword)} className="absolute right-3.5 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-neutral-600 transition-colors">
                {showPassword ? <EyeOff size={15} /> : <Eye size={15} />}
              </button>
            </div>
          </div>

          {error && <p className="text-xs text-red-500 bg-red-50 border border-red-100 rounded-lg px-3 py-2">{error}</p>}

          <button type="submit" disabled={loading}
            className="w-full py-3 bg-primary hover:bg-primary-600 disabled:bg-primary/60 text-white rounded-xl text-sm font-bold transition-colors mt-2 flex items-center justify-center gap-2">
            {loading ? <><span className="w-4 h-4 border-2 border-white/40 border-t-white rounded-full animate-spin" />Memproses...</> : 'Masuk'}
          </button>
        </form>

        <p className="text-center text-xs text-neutral-500 mt-5">
          Belum memiliki akun?{' '}
          <button onClick={() => navigate('/register')} className="text-primary font-semibold hover:text-primary-600 transition-colors">Daftar</button>
        </p>
      </div>

      <p className="text-xs text-neutral-400 mt-6 text-center">
        © {new Date().getFullYear()} Dinas Kesehatan & Komunitas SiGizi. Seluruh hak dilindungi.
      </p>
    </div>
  );
}
