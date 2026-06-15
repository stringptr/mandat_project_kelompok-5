import { useState, useEffect, useMemo } from 'react';
import { Pencil, Trash2, MoreHorizontal, MapPin, ChevronLeft, ChevronRight, X, Check } from 'lucide-react';
import { StatusBadge } from '../monitoring/components/statusbadge';
import { apiGet, apiPatch } from '../../lib/api';

// ─── Types ────────────────────────────────────────────────────────────────────

type UserStatus = 'verified' | 'pending' | 'rejected';

interface UserData {
  id: number;
  name: string;
  email: string;
  initials: string;
  avatarColor: string;
  role: string;
  roleBackend: string;
  wilayah: string;
  status: UserStatus;
}

interface EditForm {
  name: string;
  email: string;
  role: string;
  wilayah: string;
  status: UserStatus;
}

interface BackendUser {
  id_user: number;
  nama: string;
  nik: string;
  email: string;
  no_hp: string;
  jenis_kelamin: string;
  status_verifikasi: string;
  roles: string[];
  id_lokasi: number;
  created_at: string;
  updated_at: string;
}

interface UsersResponse {
  users: BackendUser[];
  total_data: number;
  page: number;
  per_page: number;
}

const AVATAR_COLORS = [
  'bg-teal-100 text-teal-700',
  'bg-blue-100 text-blue-700',
  'bg-purple-100 text-purple-700',
  'bg-rose-100 text-rose-700',
  'bg-amber-100 text-amber-700',
  'bg-emerald-100 text-emerald-700',
  'bg-indigo-100 text-indigo-700',
  'bg-pink-100 text-pink-700',
  'bg-cyan-100 text-cyan-700',
  'bg-lime-100 text-lime-700',
  'bg-orange-100 text-orange-700',
  'bg-sky-100 text-sky-700',
];

function getInitials(name: string): string {
  return name.split(' ').map((w) => w[0]).join('').toUpperCase().slice(0, 2);
}

function mapBackendRole(roles: string[]): string {
  if (roles.includes('SUPER_ADMIN')) return 'Super Admin';
  if (roles.includes('DINKES')) return 'Admin Dinkes';
  if (roles.includes('BIDAN')) return 'Bidan';
  if (roles.includes('KADER')) return 'Kader';
  if (roles.includes('IBU_HAMIL')) return 'Ibu Hamil';
  if (roles.includes('PASIEN')) return 'Pasien';
  return roles[0] || 'Unknown';
}

function mapStatus(status: string): UserStatus {
  if (status === 'Aktif') return 'verified';
  if (status === 'Ditolak') return 'rejected';
  return 'pending';
}

function mapUser(b: BackendUser, idx: number): UserData {
  return {
    id: b.id_user,
    name: b.nama,
    email: b.email,
    initials: getInitials(b.nama),
    avatarColor: AVATAR_COLORS[idx % AVATAR_COLORS.length],
    role: mapBackendRole(b.roles),
    roleBackend: b.roles[0] || '',
    wilayah: '',
    status: mapStatus(b.status_verifikasi),
  };
}

const PAGE_SIZE = 10;

// ─── Edit Modal ───────────────────────────────────────────────────────────────

interface EditModalProps {
  user: UserData;
  onSave: (updated: EditForm) => void;
  onClose: () => void;
}

function EditModal({ user, onSave, onClose }: EditModalProps): JSX.Element {
  const [form, setForm] = useState<EditForm>({
    name: user.name,
    email: user.email,
    role: user.role,
    wilayah: user.wilayah,
    status: user.status,
  });

  const set = (field: keyof EditForm, value: string) =>
    setForm((prev) => ({ ...prev, [field]: value }));

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(form);
  };

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm"
      onClick={onClose}
    >
      <div
        className="bg-white rounded-2xl shadow-xl w-full max-w-md mx-4 overflow-hidden"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between px-6 py-5 border-b border-neutral-100">
          <div className="flex items-center gap-3">
            <div className={`w-9 h-9 rounded-full flex items-center justify-center text-xs font-bold shrink-0 ${user.avatarColor}`}>
              {user.initials}
            </div>
            <div>
              <p className="font-bold text-neutral-800 text-sm leading-tight">Edit Pengguna</p>
              <p className="text-xs text-neutral-400 mt-0.5">{user.email}</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-neutral-100 transition-colors text-neutral-400"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="px-6 py-5 space-y-4">
          <div>
            <label className="block text-xs font-semibold text-neutral-500 mb-1.5 uppercase tracking-wide">Nama Lengkap</label>
            <input type="text" value={form.name} onChange={(e) => set('name', e.target.value)} required
              className="w-full px-3.5 py-2.5 text-sm border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-neutral-800" />
          </div>
          <div>
            <label className="block text-xs font-semibold text-neutral-500 mb-1.5 uppercase tracking-wide">Email</label>
            <input type="email" value={form.email} onChange={(e) => set('email', e.target.value)} required
              className="w-full px-3.5 py-2.5 text-sm border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-neutral-800" />
          </div>
          <div>
            <label className="block text-xs font-semibold text-neutral-500 mb-1.5 uppercase tracking-wide">Status</label>
            <div className="flex gap-3">
              {(['verified', 'pending', 'rejected'] as UserStatus[]).map((s) => (
                <button key={s} type="button" onClick={() => set('status', s)}
                  className={`flex-1 flex items-center justify-center gap-2 py-2.5 rounded-xl border text-xs font-bold transition-all ${form.status === s ? s === 'verified' ? 'bg-emerald-50 border-emerald-400 text-emerald-700' : s === 'rejected' ? 'bg-red-50 border-red-400 text-red-700' : 'bg-orange-50 border-orange-400 text-orange-700' : 'border-neutral-200 text-neutral-400 hover:bg-neutral-50'}`}>
                  {form.status === s && <Check className="w-3.5 h-3.5" />}
                  {s === 'verified' ? 'Aktif' : s === 'rejected' ? 'Ditolak' : 'Pending'}
                </button>
              ))}
            </div>
          </div>
          <div className="flex gap-3 pt-2">
            <button type="button" onClick={onClose}
              className="flex-1 py-2.5 rounded-xl border border-neutral-200 text-sm font-semibold text-neutral-600 hover:bg-neutral-50 transition-colors">Batal</button>
            <button type="submit"
              className="flex-1 py-2.5 rounded-xl bg-primary hover:bg-primary-700 text-white text-sm font-bold transition-colors">Simpan Perubahan</button>
          </div>
        </form>
      </div>
    </div>
  );
}

// ─── Main Component ───────────────────────────────────────────────────────────

export default function UserManagement(): JSX.Element {
  const [users, setUsers] = useState<UserData[]>([]);
  const [page, setPage] = useState(1);
  const [menuOpen, setMenuOpen] = useState<number | null>(null);
  const [editingUser, setEditingUser] = useState<UserData | null>(null);
  const [totalData, setTotalData] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUsers = async () => {
      setLoading(true);
      try {
        const res = await apiGet<UsersResponse>(`/users?page=${page}&per_page=${PAGE_SIZE}`);
        setUsers(res.users.map((u, i) => mapUser(u, i)));
        setTotalData(res.total_data);
      } catch {
        setUsers([]);
      } finally {
        setLoading(false);
      }
    };
    fetchUsers();
  }, [page]);

  const filtered = useMemo(() => users, [users]);

  const totalPages = Math.max(1, Math.ceil(totalData / PAGE_SIZE));
  const safePage = Math.min(page, totalPages);

  const openEdit = (user: UserData) => {
    setMenuOpen(null);
    setEditingUser(user);
  };

  const handleSave = async (updated: EditForm) => {
    if (!editingUser) return;
    try {
      const body: Record<string, string | null> = {};
      if (updated.name !== editingUser.name) body.nama = updated.name;
      if (updated.email !== editingUser.email) body.email = updated.email;

      if (Object.keys(body).length > 0) {
        await apiPatch(`/users/${editingUser.id}`, body);
      }

      const newStatus = updated.status === 'verified' ? 'Aktif' : updated.status === 'rejected' ? 'Ditolak' : 'Pending';
      if (newStatus === 'Aktif' || newStatus === 'Ditolak') {
        await apiPatch(`/users/${editingUser.id}/verification`, { status: newStatus });
      }

      setUsers((prev) =>
        prev.map((u) => u.id === editingUser.id ? { ...u, ...updated } : u),
      );
    } catch {
      // silent
    }
    setEditingUser(null);
  };

  const handleVerify = async (id: number) => {
    try {
      await apiPatch(`/users/${id}/verification`, { status: 'Aktif' });
      setUsers((prev) => prev.map((u) => (u.id === id ? { ...u, status: 'verified' as UserStatus } : u)));
    } catch {
      // silent
    }
  };

  return (
    <div className="w-full max-w-5xl mx-auto font-body text-neutral-800">
      <div className="bg-white rounded-2xl shadow-sm border border-neutral-100 overflow-hidden">
        <div className="px-6 pt-6 pb-4">
          <h2 className="text-base font-bold text-neutral-800 font-headline">Daftar Pengguna Aktif</h2>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-t border-neutral-100">
                <th className="px-6 py-3 text-left text-[11px] font-bold text-neutral-400 uppercase tracking-widest">Nama Lengkap</th>
                <th className="px-4 py-3 text-left text-[11px] font-bold text-neutral-400 uppercase tracking-widest">Peran</th>
                <th className="px-4 py-3 text-left text-[11px] font-bold text-neutral-400 uppercase tracking-widest">Wilayah Kerja</th>
                <th className="px-4 py-3 text-left text-[11px] font-bold text-neutral-400 uppercase tracking-widest">Status</th>
                <th className="px-4 py-3 text-right text-[11px] font-bold text-neutral-400 uppercase tracking-widest pr-6">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-neutral-100">
              {loading ? (
                <tr>
                  <td colSpan={5} className="py-16 text-center text-neutral-400 text-sm">Memuat data...</td>
                </tr>
              ) : filtered.length === 0 ? (
                <tr>
                  <td colSpan={5} className="py-16 text-center text-neutral-400 text-sm">Tidak ada pengguna.</td>
                </tr>
              ) : (
                filtered.map((user) => (
                  <tr key={user.id} className="hover:bg-neutral-50 transition-colors" onClick={() => setMenuOpen(null)}>
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <div className={`w-9 h-9 rounded-full flex items-center justify-center text-xs font-bold shrink-0 ${user.avatarColor}`}>
                          {user.initials}
                        </div>
                        <div className="min-w-0">
                          <p className="font-semibold text-neutral-800 leading-tight truncate">{user.name}</p>
                          <p className="text-[12px] text-neutral-400 mt-0.5 truncate">{user.email}</p>
                        </div>
                      </div>
                    </td>
                    <td className="px-4 py-4">
                      <span className="text-primary font-semibold text-sm">{user.role}</span>
                    </td>
                    <td className="px-4 py-4">
                      <span className="flex items-center gap-1 text-neutral-600 text-sm">
                        <MapPin className="w-3.5 h-3.5 text-neutral-400 shrink-0" />
                        {user.wilayah || '-'}
                      </span>
                    </td>
                    <td className="px-4 py-4">
                      <StatusBadge variant={user.status} />
                    </td>
                    <td className="px-4 py-4 pr-6">
                      <div className="flex items-center justify-end gap-1">
                        {user.status === 'pending' ? (
                          <>
                            <button onClick={(e) => { e.stopPropagation(); handleVerify(user.id); }}
                              className="bg-primary hover:bg-primary-700 text-white text-xs font-bold px-4 py-1.5 rounded-lg transition-colors">Verify</button>
                            <div className="relative ml-1">
                              <button onClick={(e) => { e.stopPropagation(); setMenuOpen(menuOpen === user.id ? null : user.id); }}
                                className="w-7 h-7 flex items-center justify-center rounded-lg hover:bg-neutral-100 transition-colors text-neutral-400">
                                <MoreHorizontal className="w-4 h-4" />
                              </button>
                              {menuOpen === user.id && (
                                <div className="absolute right-0 top-8 w-36 bg-white rounded-xl shadow-lg border border-neutral-100 z-20 overflow-hidden">
                                  <button className="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-neutral-700 hover:bg-neutral-50 transition-colors"
                                    onClick={(e) => { e.stopPropagation(); openEdit(user); }}>
                                    <Pencil className="w-3.5 h-3.5 text-neutral-400" /> Edit
                                  </button>
                                  <button className="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-red-500 hover:bg-red-50 transition-colors"
                                    onClick={(e) => { e.stopPropagation(); setMenuOpen(null); }}>
                                    <Trash2 className="w-3.5 h-3.5" /> Hapus
                                  </button>
                                </div>
                              )}
                            </div>
                          </>
                        ) : (
                          <>
                            <button onClick={(e) => { e.stopPropagation(); openEdit(user); }}
                              className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-neutral-100 transition-colors text-neutral-400 hover:text-primary">
                              <Pencil className="w-4 h-4" />
                            </button>
                            <button onClick={(e) => e.stopPropagation()}
                              className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-red-50 transition-colors text-neutral-400 hover:text-red-500">
                              <Trash2 className="w-4 h-4" />
                            </button>
                          </>
                        )}
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        <div className="flex items-center justify-between px-6 py-4 border-t border-neutral-100">
          <span className="text-xs text-neutral-400">
            Menampilkan {totalData === 0 ? 0 : (safePage - 1) * PAGE_SIZE + 1}–
            {Math.min(safePage * PAGE_SIZE, totalData)} dari {totalData.toLocaleString('id-ID')} pengguna
          </span>

          <div className="flex items-center gap-1">
            <button onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={safePage === 1}
              className="w-8 h-8 flex items-center justify-center rounded-lg border border-neutral-200 text-neutral-500 hover:bg-neutral-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors">
              <ChevronLeft className="w-4 h-4" />
            </button>

            {Array.from({ length: totalPages }, (_, i) => i + 1)
              .filter((p) => p === 1 || p === totalPages || Math.abs(p - safePage) <= 1)
              .reduce<(number | '...')[]>((acc, p, idx, arr) => {
                if (idx > 0 && (p as number) - (arr[idx - 1] as number) > 1) acc.push('...');
                acc.push(p);
                return acc;
              }, [])
              .map((p, idx) =>
                p === '...' ? (
                  <span key={`dots-${idx}`} className="w-8 h-8 flex items-center justify-center text-neutral-400 text-sm">…</span>
                ) : (
                  <button key={p} onClick={() => setPage(p as number)}
                    className={`w-8 h-8 flex items-center justify-center rounded-lg text-sm font-medium transition-colors ${safePage === p ? 'bg-primary text-white font-bold' : 'border border-neutral-200 text-neutral-600 hover:bg-neutral-50'}`}>
                    {p}
                  </button>
                ),
              )}

            <button onClick={() => setPage((p) => Math.min(totalPages, p + 1))} disabled={safePage === totalPages}
              className="w-8 h-8 flex items-center justify-center rounded-lg border border-neutral-200 text-neutral-500 hover:bg-neutral-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors">
              <ChevronRight className="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      {editingUser && (
        <EditModal user={editingUser} onSave={handleSave} onClose={() => setEditingUser(null)} />
      )}
    </div>
  );
}
