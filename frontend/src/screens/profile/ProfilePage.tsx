import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { User, Mail, Phone, ShieldCheck, Save, Loader2 } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';
import { useNotification } from '../../context/NotificationContext';
import { apiPatch, apiGet } from '../../lib/api';

export default function ProfilePage(): JSX.Element {
  const { user, refreshProfile } = useAuth();
  const notify = useNotification();
  const navigate = useNavigate();
  const [saving, setSaving] = useState(false);
  const [loading, setLoading] = useState(true);

  const [form, setForm] = useState({
    nama: '',
    email: '',
    no_hp: '',
  });

  useEffect(() => {
    if (!user) return;
    apiGet<{ no_hp: string }>(`/users/${user.idUser}`).then((res) => {
      setForm({
        nama: user.name || '',
        email: user.email || '',
        no_hp: res.no_hp ?? '',
      });
    }).catch(() => {
      setForm({
        nama: user.name || '',
        email: user.email || '',
        no_hp: '',
      });
    }).finally(() => setLoading(false));
  }, [user]);

  if (!user) {
    return (
      <div className="p-8 text-center">
        <p className="text-neutral-500">Silakan login terlebih dahulu.</p>
      </div>
    );
  }

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.nama.trim()) { notify.warn('Nama tidak boleh kosong.'); return; }
    setSaving(true);
    try {
      const body: Record<string, string> = {
        nama: form.nama,
        email: form.email,
        no_hp: form.no_hp,
      };
      await apiPatch(`/users/${user.idUser}`, body);
      notify.success('Profil berhasil diperbarui.');
      await refreshProfile();
      navigate(-1);
    } catch {
      notify.error('Gagal memperbarui profil.');
    } finally {
      setSaving(false);
    }
  };

  const initials = (user.name || '?').split(' ').map((w) => w[0]).join('').toUpperCase().slice(0, 2);

  if (loading) {
    return (
      <div className="p-8 flex items-center justify-center py-16">
        <Loader2 size={24} className="animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="max-w-lg mx-auto p-8 font-body">
      <div className="bg-white rounded-2xl shadow-sm border border-neutral-100 overflow-hidden">
        <div className="px-6 pt-6 pb-4 border-b border-neutral-100">
          <h2 className="text-base font-bold text-neutral-800 font-headline">Edit Profil</h2>
          <p className="text-xs text-neutral-400 mt-1">Perbarui data diri Anda.</p>
        </div>

        <div className="px-6 py-6 space-y-5">
          <div className="flex items-center gap-4">
            <div className="w-16 h-16 rounded-full bg-primary-100 text-primary text-lg font-bold flex items-center justify-center shrink-0">{initials}</div>
            <div>
              <p className="font-bold text-neutral-800">{user.name}</p>
              <p className="text-sm text-neutral-500">{user.role}</p>
              <p className="text-xs text-neutral-400 mt-0.5">NIK: {user.nik}</p>
            </div>
          </div>

          <form onSubmit={handleSave} className="space-y-4">
            <div>
              <label className="flex items-center gap-1.5 text-xs font-semibold text-neutral-500 mb-1.5">
                <User size={12} /> Nama Lengkap
              </label>
              <input type="text" value={form.nama} onChange={(e) => setForm({ ...form, nama: e.target.value })}
                className="w-full px-3.5 py-2.5 text-sm border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/30 text-neutral-800" />
            </div>
            <div>
              <label className="flex items-center gap-1.5 text-xs font-semibold text-neutral-500 mb-1.5">
                <Mail size={12} /> Email
              </label>
              <input type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })}
                className="w-full px-3.5 py-2.5 text-sm border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/30 text-neutral-800" />
            </div>
            <div>
              <label className="flex items-center gap-1.5 text-xs font-semibold text-neutral-500 mb-1.5">
                <Phone size={12} /> No. Telepon
              </label>
              <input type="tel" value={form.no_hp} onChange={(e) => setForm({ ...form, no_hp: e.target.value })}
                className="w-full px-3.5 py-2.5 text-sm border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/30 text-neutral-800" />
            </div>
            <div className="bg-neutral-50 border border-neutral-100 rounded-xl p-3 flex items-start gap-2">
              <ShieldCheck size={14} className="text-neutral-400 mt-0.5" />
              <div className="text-xs text-neutral-500">
                <span className="font-semibold">NIK tidak dapat diubah.</span> Untuk perubahan data kependudukan, hubungi admin Dinas Kesehatan.
              </div>
            </div>
            <div className="flex gap-3 pt-2">
              <button type="button" onClick={() => navigate(-1)}
                className="flex-1 py-2.5 rounded-xl border border-neutral-200 text-sm font-semibold text-neutral-600 hover:bg-neutral-50 transition-colors">Batal</button>
              <button type="submit" disabled={saving}
                className="flex-1 py-2.5 rounded-xl bg-primary hover:bg-primary-700 text-white text-sm font-bold transition-colors flex items-center justify-center gap-2">
                {saving ? <Loader2 size={16} className="animate-spin" /> : <Save size={15} />}
                {saving ? 'Menyimpan...' : 'Simpan'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
