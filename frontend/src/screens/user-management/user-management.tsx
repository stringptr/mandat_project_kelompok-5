import { useState, useEffect } from 'react';
import { Pencil, Trash2, MoreHorizontal, MapPin, ChevronLeft, ChevronRight, X, Check, Loader2, Eye } from 'lucide-react';
import { useNotification } from '../../context/NotificationContext';
import { StatusBadge } from '../monitoring/components/statusbadge';
import { apiGet, apiPatch, apiDelete } from '../../lib/api';

// ─── Types ────────────────────────────────────────────────────────────────────

type UserStatus = 'verified' | 'pending' | 'rejected';

interface UserData {
  id: number;
  name: string;
  email: string;
  nik: string;
  initials: string;
  avatarColor: string;
  role: string;
  roleBackend: string;
  roles: string[];
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
  nama_lokasi: string;
  created_at: string;
  updated_at: string;
}

interface UsersResponse {
  users: BackendUser[];
  meta: {
    current_page: number;
    per_page: number;
    total: number;
    last_page: number;
  };
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

const DISPLAY_ROLES = ['SUPER_ADMIN', 'DINKES', 'BIDAN', 'KADER', 'IBU_HAMIL', 'ANAK', 'PASIEN'];

const ROLE_LABEL: Record<string, string> = {
  SUPER_ADMIN: 'Super Admin',
  DINKES: 'Dinkes',
  BIDAN: 'Bidan',
  KADER: 'Kader',
  IBU_HAMIL: 'Ibu Hamil',
  ANAK: 'Anak',
  PASIEN: 'Pasien',
};

const ROLE_BADGE: Record<string, string> = {
  SUPER_ADMIN: 'bg-purple-100 text-purple-700',
  DINKES: 'bg-indigo-100 text-indigo-700',
  BIDAN: 'bg-emerald-100 text-emerald-700',
  KADER: 'bg-amber-100 text-amber-700',
  IBU_HAMIL: 'bg-pink-100 text-pink-700',
  ANAK: 'bg-sky-100 text-sky-700',
  PASIEN: 'bg-blue-100 text-blue-700',
};

function getUserRoleBadges(roles: string[]): string[] {
  return roles.filter((r) => DISPLAY_ROLES.includes(r));
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
    nik: b.nik,
    initials: getInitials(b.nama),
    avatarColor: AVATAR_COLORS[idx % AVATAR_COLORS.length],
    role: mapBackendRole(b.roles),
    roleBackend: b.roles[0] || '',
    roles: b.roles,
    wilayah: b.nama_lokasi || '',
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

// ─── Detail Modal ──────────────────────────────────────────────────────────────

interface UserDetail {
  id_user: number;
  email: string;
  no_hp: string;
  nama: string;
  nik: string;
  jenis_kelamin: string;
  tanggal_lahir: string;
  status_verifikasi: string;
  id_lokasi: number;
  roles: string[];
  created_at: string;
}

interface PasienDetailData {
  id_pasien: number;
  nama: string;
  nik: string;
  email: string;
  no_hp: string;
  jenis_kelamin: string;
  tanggal_lahir: string;
  nama_posyandu: string;
  jenis_pasien: string;
  data_ibu_hamil?: { id_ibu_hamil: number; hamil_ke: number; bulan_mulai_hamil: string; hpht: string; status_kehamilan: string } | null;
  data_anak?: { nama_anak: string; berat_lahir: number; panjang_lahir: number; hubungan_dengan_wali: string; nama_wali: string } | null;
}

interface DetailModalProps {
  userId: number;
  userName: string;
  userRoles: string[];
  onClose: () => void;
}

function DetailModal({ userId, userName, userRoles, onClose }: DetailModalProps): JSX.Element {
  const [loading, setLoading] = useState(true);
  const [userDetail, setUserDetail] = useState<UserDetail | null>(null);
  const [pasienDetail, setPasienDetail] = useState<PasienDetailData | null>(null);

  useEffect(() => {
    setLoading(true);
    const isPasien = userRoles.some((r) => ['PASIEN', 'IBU_HAMIL', 'ANAK'].includes(r));

    Promise.all([
      apiGet<UserDetail>(`/users/${userId}`).catch(() => null),
      isPasien ? apiGet<PasienDetailData>(`/monitoring/pasien/${userId}`).catch(() => null) : Promise.resolve(null),
    ]).then(([u, p]) => {
      setUserDetail(u);
      setPasienDetail(p);
    }).finally(() => setLoading(false));
  }, [userId, userRoles]);

  const initials = userName.split(' ').map((w) => w[0]).join('').toUpperCase().slice(0, 2);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm" onClick={onClose}>
      <div className="bg-white rounded-2xl shadow-xl w-full max-w-lg mx-4 max-h-[85vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
        <div className="flex items-center justify-between px-6 py-5 border-b border-neutral-100 sticky top-0 bg-white z-10">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-primary-100 text-primary text-sm font-bold flex items-center justify-center shrink-0">{initials}</div>
            <div>
              <p className="font-bold text-neutral-800 text-sm">{userName}</p>
              <p className="text-xs text-neutral-400">{userDetail?.email || '-'}</p>
            </div>
          </div>
          <button onClick={onClose} className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-neutral-100 text-neutral-400"><X size={16} /></button>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-12"><Loader2 size={24} className="animate-spin text-primary" /></div>
        ) : (
          <div className="px-6 py-5 space-y-4">
            {/* Basic Info */}
            <div>
              <p className="text-xs font-semibold text-neutral-400 uppercase tracking-wide mb-2">Informasi Dasar</p>
              <div className="grid grid-cols-2 gap-3 text-sm">
                <div><span className="text-neutral-400">NIK</span><p className="font-medium text-neutral-800">{userDetail?.nik || '-'}</p></div>
                <div><span className="text-neutral-400">No. HP</span><p className="font-medium text-neutral-800">{userDetail?.no_hp || '-'}</p></div>
                <div><span className="text-neutral-400">Jenis Kelamin</span><p className="font-medium text-neutral-800">{userDetail?.jenis_kelamin || '-'}</p></div>
                <div><span className="text-neutral-400">Tgl Lahir</span><p className="font-medium text-neutral-800">{userDetail?.tanggal_lahir ? new Date(userDetail.tanggal_lahir).toLocaleDateString('id-ID') : '-'}</p></div>
                <div><span className="text-neutral-400">Status</span><span className={`inline-block text-xs font-bold px-2 py-0.5 rounded-full ml-1 ${userDetail?.status_verifikasi === 'Aktif' ? 'bg-emerald-100 text-emerald-700' : userDetail?.status_verifikasi === 'Ditolak' ? 'bg-red-100 text-red-700' : 'bg-amber-100 text-amber-700'}`}>{userDetail?.status_verifikasi || '-'}</span></div>
                <div><span className="text-neutral-400">Terdaftar</span><p className="font-medium text-neutral-800">{userDetail?.created_at ? new Date(userDetail.created_at).toLocaleDateString('id-ID') : '-'}</p></div>
              </div>
            </div>

            {/* Roles */}
            <div>
              <p className="text-xs font-semibold text-neutral-400 uppercase tracking-wide mb-2">Role</p>
              <div className="flex flex-wrap gap-1">
                {userRoles.map((r) => (
                  <span key={r} className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${ROLE_BADGE[r] ?? 'bg-neutral-100 text-neutral-600'}`}>{ROLE_LABEL[r] ?? r}</span>
                ))}
              </div>
            </div>

            {/* Pasien Info */}
            {pasienDetail && (
              <div>
                <p className="text-xs font-semibold text-neutral-400 uppercase tracking-wide mb-2">Detail Pasien</p>
                <div className="bg-neutral-50 rounded-xl p-4 space-y-3 text-sm">
                  <div className="grid grid-cols-2 gap-3">
                    <div><span className="text-neutral-400">Posyandu</span><p className="font-medium text-neutral-800">{pasienDetail.nama_posyandu || '-'}</p></div>
                    <div><span className="text-neutral-400">Jenis Pasien</span><span className={`inline-block text-xs font-bold px-2 py-0.5 rounded-full ml-1 ${ROLE_BADGE[pasienDetail.jenis_pasien] ?? 'bg-neutral-100 text-neutral-600'}`}>{pasienDetail.jenis_pasien || '-'}</span></div>
                  </div>

                  {pasienDetail.data_ibu_hamil && (
                    <div className="border-t border-neutral-200 pt-3 mt-3">
                      <p className="text-xs font-semibold text-pink-600 mb-2">Ibu Hamil</p>
                      <div className="grid grid-cols-2 gap-3">
                        <div><span className="text-neutral-400">Hamil Ke</span><p className="font-medium text-neutral-800">{pasienDetail.data_ibu_hamil.hamil_ke}</p></div>
                        <div><span className="text-neutral-400">Status</span><span className="inline-block text-xs font-bold px-2 py-0.5 rounded-full bg-pink-100 text-pink-700 ml-1">{pasienDetail.data_ibu_hamil.status_kehamilan}</span></div>
                        <div><span className="text-neutral-400">HPHT</span><p className="font-medium text-neutral-800">{pasienDetail.data_ibu_hamil.hpht ? new Date(pasienDetail.data_ibu_hamil.hpht).toLocaleDateString('id-ID') : '-'}</p></div>
                        <div><span className="text-neutral-400">Mulai Hamil</span><p className="font-medium text-neutral-800">{pasienDetail.data_ibu_hamil.bulan_mulai_hamil ? new Date(pasienDetail.data_ibu_hamil.bulan_mulai_hamil).toLocaleDateString('id-ID') : '-'}</p></div>
                      </div>
                    </div>
                  )}

                  {pasienDetail.data_anak && (
                    <div className="border-t border-neutral-200 pt-3 mt-3">
                      <p className="text-xs font-semibold text-sky-600 mb-2">Anak</p>
                      <div className="grid grid-cols-2 gap-3">
                        <div><span className="text-neutral-400">Nama Anak</span><p className="font-medium text-neutral-800">{pasienDetail.data_anak.nama_anak}</p></div>
                        <div><span className="text-neutral-400">Wali</span><p className="font-medium text-neutral-800">{pasienDetail.data_anak.nama_wali || '-'}</p></div>
                        <div><span className="text-neutral-400">BB Lahir</span><p className="font-medium text-neutral-800">{pasienDetail.data_anak.berat_lahir} kg</p></div>
                        <div><span className="text-neutral-400">PB Lahir</span><p className="font-medium text-neutral-800">{pasienDetail.data_anak.panjang_lahir} cm</p></div>
                      </div>
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

// ─── Main Component ───────────────────────────────────────────────────────────

export default function UserManagement(): JSX.Element {
  const notify = useNotification();
  const [users, setUsers] = useState<UserData[]>([]);
  const [page, setPage] = useState(1);
  const [menuOpen, setMenuOpen] = useState<number | null>(null);
  const [editingUser, setEditingUser] = useState<UserData | null>(null);
  const [deletingId, setDeletingId] = useState<number | null>(null);
  const [totalData, setTotalData] = useState(0);
  const [loading, setLoading] = useState(true);
  const [detailUser, setDetailUser] = useState<{ id: number; name: string; roles: string[] } | null>(null);
  const [roleFilter, setRoleFilter] = useState('');

  useEffect(() => {
    const fetchUsers = async () => {
      setLoading(true);
      try {
        const params = new URLSearchParams({ page: String(page), per_page: String(PAGE_SIZE) });
        if (roleFilter) params.set('role', roleFilter);
        const res = await apiGet<UsersResponse>(`/users?${params.toString()}`);
        setUsers(res.users.map((u, i) => mapUser(u, i)));
        setTotalData(res.meta.total);
      } catch {
        setUsers([]);
      } finally {
        setLoading(false);
      }
    };
    fetchUsers();
  }, [page, roleFilter]);

  const totalPages = Math.max(1, Math.ceil(totalData / PAGE_SIZE));
  const safePage = Math.min(page, totalPages);

  const ROLE_FILTERS = [
    { value: '', label: 'Semua' },
    { value: 'Bidan', label: 'Bidan', color: 'bg-emerald-100 text-emerald-700 border-emerald-300' },
    { value: 'Kader', label: 'Kader', color: 'bg-amber-100 text-amber-700 border-amber-300' },
    { value: 'Dinkes', label: 'Dinkes', color: 'bg-indigo-100 text-indigo-700 border-indigo-300' },
    { value: 'Ibu Hamil', label: 'Ibu Hamil', color: 'bg-pink-100 text-pink-700 border-pink-300' },
    { value: 'Anak', label: 'Anak', color: 'bg-sky-100 text-sky-700 border-sky-300' },
    { value: 'Pasien', label: 'Pasien', color: 'bg-blue-100 text-blue-700 border-blue-300' },
  ];

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
      notify.success('Data pengguna berhasil diperbarui');
    } catch {
      notify.error('Gagal memperbarui data pengguna.');
    }
    setEditingUser(null);
  };

  const handleVerify = async (id: number) => {
    try {
      await apiPatch(`/users/${id}/verification`, { status: 'Aktif' });
      setUsers((prev) => prev.map((u) => (u.id === id ? { ...u, status: 'verified' as UserStatus } : u)));
      notify.success('Pengguna berhasil diverifikasi');
    } catch {
      notify.error('Gagal memverifikasi pengguna.');
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('Hapus pengguna ini? Tindakan ini tidak dapat dibatalkan.')) return;
    setDeletingId(id);
    try {
      await apiDelete(`/users/${id}`);
      setUsers((prev) => prev.filter((u) => u.id !== id));
      setTotalData((prev) => Math.max(0, prev - 1));
      notify.success('Pengguna berhasil dihapus');
    } catch {
      notify.error('Gagal menghapus pengguna.');
    } finally {
      setDeletingId(null);
      setMenuOpen(null);
    }
  };

  return (
    <div className="w-full max-w-5xl mx-auto font-body text-neutral-800">
      <div className="bg-white rounded-2xl shadow-sm border border-neutral-100 overflow-hidden">
        <div className="px-6 pt-6 pb-4">
          <h2 className="text-base font-bold text-neutral-800 font-headline">Daftar Pengguna Aktif</h2>
        </div>

        <div className="px-6 pb-4 flex flex-wrap gap-2">
          {ROLE_FILTERS.map((f) => (
            <button key={f.value} onClick={() => { setRoleFilter(f.value); setPage(1); }}
              className={`text-xs font-bold px-3 py-1.5 rounded-full border transition-colors ${roleFilter === f.value ? f.color + ' border' : 'border-neutral-200 text-neutral-500 hover:bg-neutral-100'}`}>
              {f.label}
            </button>
          ))}
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-t border-neutral-100">
                <th className="px-6 py-3 text-left text-[11px] font-bold text-neutral-400 uppercase tracking-widest">Nama Lengkap</th>
                <th className="px-4 py-3 text-left text-[11px] font-bold text-neutral-400 uppercase tracking-widest">NIK</th>
                <th className="px-4 py-3 text-left text-[11px] font-bold text-neutral-400 uppercase tracking-widest">Peran</th>
                <th className="px-4 py-3 text-left text-[11px] font-bold text-neutral-400 uppercase tracking-widest">Wilayah Kerja</th>
                <th className="px-4 py-3 text-left text-[11px] font-bold text-neutral-400 uppercase tracking-widest">Status</th>
                <th className="px-4 py-3 text-right text-[11px] font-bold text-neutral-400 uppercase tracking-widest pr-6">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-neutral-100">
                {users.length === 0 && !loading ? (
                  <tr>
                    <td colSpan={6} className="py-16 text-center text-neutral-400 text-sm">Tidak ada pengguna.</td>
                  </tr>
                ) : (
                  users.map((user) => (
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
                        <span className="text-neutral-700 text-sm font-mono">{user.nik}</span>
                      </td>
                      <td className="px-4 py-4">
                        <div className="flex flex-wrap gap-1">
                          {getUserRoleBadges(user.roles).map((r) => (
                            <span key={r} className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${ROLE_BADGE[r] ?? 'bg-neutral-100 text-neutral-600'}`}>
                              {ROLE_LABEL[r] ?? r}
                            </span>
                          ))}
                        </div>
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
                                    onClick={(e) => { e.stopPropagation(); handleDelete(user.id); }}>
                                    <Trash2 className="w-3.5 h-3.5" /> Hapus
                                  </button>
        </div>
                              )}
                            </div>
                          </>
                        ) : (
                          <>
                            <button onClick={(e) => { e.stopPropagation(); setDetailUser({ id: user.id, name: user.name, roles: user.roles }); }}
                              className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-blue-50 transition-colors text-neutral-400 hover:text-blue-500">
                              <Eye size={15} />
                            </button>
                            <button onClick={(e) => { e.stopPropagation(); openEdit(user); }}
                              className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-neutral-100 transition-colors text-neutral-400 hover:text-primary">
                              <Pencil className="w-4 h-4" />
                            </button>
                            <button onClick={(e) => { e.stopPropagation(); handleDelete(user.id); }}
                              disabled={deletingId === user.id}
                              className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-red-50 transition-colors text-neutral-400 hover:text-red-500 disabled:opacity-40">
                              {deletingId === user.id
                                ? <Loader2 className="w-4 h-4 animate-spin" />
                                : <Trash2 className="w-4 h-4" />}
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
      {detailUser && (
        <DetailModal userId={detailUser.id} userName={detailUser.name} userRoles={detailUser.roles} onClose={() => setDetailUser(null)} />
      )}
    </div>
  );
}
