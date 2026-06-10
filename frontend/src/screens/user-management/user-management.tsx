import { useState, useMemo } from 'react';
import { Pencil, Trash2, MoreHorizontal, MapPin, ChevronLeft, ChevronRight, X, Check } from 'lucide-react';
import { StatusBadge } from '../monitoring/components/statusbadge';

// ─── Types ────────────────────────────────────────────────────────────────────

type UserStatus = 'verified' | 'pending';

interface UserData {
  id: number;
  name: string;
  email: string;
  initials: string;
  avatarColor: string;
  role: string;
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

// ─── Dummy Data ───────────────────────────────────────────────────────────────

const INITIAL_USERS: UserData[] = [
  { id: 1,  name: 'Aisyah Nurhaliza',   email: 'aisyah.n@dinkes.go.id',         initials: 'AN', avatarColor: 'bg-teal-100 text-teal-700',      role: 'Admin Dinkes',       wilayah: 'Kota Bandung',    status: 'verified' },
  { id: 2,  name: 'Budi Pratama',       email: 'budi.kader@puskesmas.id',        initials: 'BP', avatarColor: 'bg-blue-100 text-blue-700',       role: 'Kader Puskesmas',    wilayah: 'Kec. Sukajadi',   status: 'pending'  },
  { id: 3,  name: 'Siti Wahyuni',       email: 'siti.w@gmail.com',               initials: 'SW', avatarColor: 'bg-purple-100 text-purple-700',   role: 'Kader Posyandu',     wilayah: 'Kel. Sukawarna',  status: 'verified' },
  { id: 4,  name: 'Dewi Anggraini',     email: 'dewi.a@puskesmas.id',            initials: 'DA', avatarColor: 'bg-rose-100 text-rose-700',       role: 'Bidan Desa',         wilayah: 'Kec. Cidadap',    status: 'verified' },
  { id: 5,  name: 'Rudi Hartono',       email: 'rudi.h@dinkes.go.id',            initials: 'RH', avatarColor: 'bg-amber-100 text-amber-700',    role: 'Admin Dinkes',       wilayah: 'Kota Semarang',   status: 'verified' },
  { id: 6,  name: 'Mega Lestari',       email: 'mega.lestari@gmail.com',         initials: 'ML', avatarColor: 'bg-emerald-100 text-emerald-700', role: 'Kader Posyandu',     wilayah: 'Kel. Tembalang',  status: 'pending'  },
  { id: 7,  name: 'Ahmad Fauzi',        email: 'ahmad.f@puskesmas.id',           initials: 'AF', avatarColor: 'bg-indigo-100 text-indigo-700',  role: 'Bidan Koordinator',  wilayah: 'Kec. Banyumanik', status: 'verified' },
  { id: 8,  name: 'Yuni Rahayu',        email: 'yuni.r@posyandu.go.id',          initials: 'YR', avatarColor: 'bg-pink-100 text-pink-700',      role: 'Kader Posyandu',     wilayah: 'Kel. Srondol',    status: 'verified' },
  { id: 9,  name: 'Hendra Wijaya',      email: 'hendra.w@dinkes.go.id',          initials: 'HW', avatarColor: 'bg-cyan-100 text-cyan-700',      role: 'Admin Dinkes',       wilayah: 'Kab. Semarang',   status: 'pending'  },
  { id: 10, name: 'Lina Marlina',       email: 'lina.m@gmail.com',               initials: 'LM', avatarColor: 'bg-lime-100 text-lime-700',      role: 'Ibu / Wali',         wilayah: 'Kel. Pedurungan', status: 'verified' },
  { id: 11, name: 'Farida Hanum',       email: 'farida.h@puskesmas.id',          initials: 'FH', avatarColor: 'bg-orange-100 text-orange-700',  role: 'Bidan Desa',         wilayah: 'Kec. Gunungpati', status: 'verified' },
  { id: 12, name: 'Tono Sutrisno',      email: 'tono.s@gmail.com',               initials: 'TS', avatarColor: 'bg-sky-100 text-sky-700',        role: 'Kader Puskesmas',    wilayah: 'Kec. Mijen',      status: 'pending'  },
];

const ROLE_OPTIONS = [
  'Admin Dinkes', 'Bidan Desa', 'Bidan Koordinator',
  'Kader Posyandu', 'Kader Puskesmas', 'Ibu / Wali',
];

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
    /* Backdrop */
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm"
      onClick={onClose}
    >
      {/* Modal card */}
      <div
        className="bg-white rounded-2xl shadow-xl w-full max-w-md mx-4 overflow-hidden"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
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

        {/* Form */}
        <form onSubmit={handleSubmit} className="px-6 py-5 space-y-4">
          {/* Nama */}
          <div>
            <label className="block text-xs font-semibold text-neutral-500 mb-1.5 uppercase tracking-wide">
              Nama Lengkap
            </label>
            <input
              type="text"
              value={form.name}
              onChange={(e) => set('name', e.target.value)}
              required
              className="w-full px-3.5 py-2.5 text-sm border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-neutral-800"
            />
          </div>

          {/* Email */}
          <div>
            <label className="block text-xs font-semibold text-neutral-500 mb-1.5 uppercase tracking-wide">
              Email
            </label>
            <input
              type="email"
              value={form.email}
              onChange={(e) => set('email', e.target.value)}
              required
              className="w-full px-3.5 py-2.5 text-sm border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-neutral-800"
            />
          </div>

          {/* Peran */}
          <div>
            <label className="block text-xs font-semibold text-neutral-500 mb-1.5 uppercase tracking-wide">
              Peran
            </label>
            <select
              value={form.role}
              onChange={(e) => set('role', e.target.value)}
              className="w-full px-3.5 py-2.5 text-sm border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-neutral-800 bg-white"
            >
              {ROLE_OPTIONS.map((r) => (
                <option key={r} value={r}>{r}</option>
              ))}
            </select>
          </div>

          {/* Wilayah */}
          <div>
            <label className="block text-xs font-semibold text-neutral-500 mb-1.5 uppercase tracking-wide">
              Wilayah Kerja
            </label>
            <input
              type="text"
              value={form.wilayah}
              onChange={(e) => set('wilayah', e.target.value)}
              required
              className="w-full px-3.5 py-2.5 text-sm border border-neutral-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary text-neutral-800"
            />
          </div>

          {/* Status */}
          <div>
            <label className="block text-xs font-semibold text-neutral-500 mb-1.5 uppercase tracking-wide">
              Status
            </label>
            <div className="flex gap-3">
              {(['verified', 'pending'] as UserStatus[]).map((s) => (
                <button
                  key={s}
                  type="button"
                  onClick={() => set('status', s)}
                  className={`flex-1 flex items-center justify-center gap-2 py-2.5 rounded-xl border text-xs font-bold transition-all ${
                    form.status === s
                      ? s === 'verified'
                        ? 'bg-emerald-50 border-emerald-400 text-emerald-700'
                        : 'bg-orange-50 border-orange-400 text-orange-700'
                      : 'border-neutral-200 text-neutral-400 hover:bg-neutral-50'
                  }`}
                >
                  {form.status === s && <Check className="w-3.5 h-3.5" />}
                  {s === 'verified' ? 'Verified' : 'Pending'}
                </button>
              ))}
            </div>
          </div>

          {/* Footer buttons */}
          <div className="flex gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 py-2.5 rounded-xl border border-neutral-200 text-sm font-semibold text-neutral-600 hover:bg-neutral-50 transition-colors"
            >
              Batal
            </button>
            <button
              type="submit"
              className="flex-1 py-2.5 rounded-xl bg-primary hover:bg-primary-700 text-white text-sm font-bold transition-colors"
            >
              Simpan Perubahan
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// ─── Main Component ───────────────────────────────────────────────────────────

export default function UserManagement(): JSX.Element {
  const [users, setUsers] = useState<UserData[]>(INITIAL_USERS);
  const [page, setPage] = useState(1);
  const [menuOpen, setMenuOpen] = useState<number | null>(null);
  const [editingUser, setEditingUser] = useState<UserData | null>(null);

  const filtered = useMemo(() => users, [users]);

  const totalPages = Math.max(1, Math.ceil(filtered.length / PAGE_SIZE));
  const safePage = Math.min(page, totalPages);
  const paged = filtered.slice((safePage - 1) * PAGE_SIZE, safePage * PAGE_SIZE);

  const openEdit = (user: UserData) => {
    setMenuOpen(null);
    setEditingUser(user);
  };

  const handleSave = (updated: EditForm) => {
    if (!editingUser) return;
    setUsers((prev) =>
      prev.map((u) =>
        u.id === editingUser.id
          ? { ...u, ...updated }
          : u,
      ),
    );
    setEditingUser(null);
  };

  const handleVerify = (id: number) => {
    setUsers((prev) =>
      prev.map((u) => (u.id === id ? { ...u, status: 'verified' } : u)),
    );
  };

  return (
    <div className="w-full max-w-5xl mx-auto font-body text-neutral-800">
      {/* Table Card */}
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
              {paged.length === 0 ? (
                <tr>
                  <td colSpan={5} className="py-16 text-center text-neutral-400 text-sm">
                    Tidak ada pengguna.
                  </td>
                </tr>
              ) : (
                paged.map((user) => (
                  <tr
                    key={user.id}
                    className="hover:bg-neutral-50 transition-colors"
                    onClick={() => setMenuOpen(null)}
                  >
                    {/* Avatar + Name */}
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

                    {/* Role */}
                    <td className="px-4 py-4">
                      <span className="text-primary font-semibold text-sm">{user.role}</span>
                    </td>

                    {/* Wilayah */}
                    <td className="px-4 py-4">
                      <span className="flex items-center gap-1 text-neutral-600 text-sm">
                        <MapPin className="w-3.5 h-3.5 text-neutral-400 shrink-0" />
                        {user.wilayah}
                      </span>
                    </td>

                    {/* Status */}
                    <td className="px-4 py-4">
                      <StatusBadge variant={user.status} />
                    </td>

                    {/* Actions */}
                    <td className="px-4 py-4 pr-6">
                      <div className="flex items-center justify-end gap-1">
                        {user.status === 'pending' ? (
                          <>
                            <button
                              onClick={(e) => { e.stopPropagation(); handleVerify(user.id); }}
                              className="bg-primary hover:bg-primary-700 text-white text-xs font-bold px-4 py-1.5 rounded-lg transition-colors"
                            >
                              Verify
                            </button>
                            <div className="relative ml-1">
                              <button
                                onClick={(e) => { e.stopPropagation(); setMenuOpen(menuOpen === user.id ? null : user.id); }}
                                className="w-7 h-7 flex items-center justify-center rounded-lg hover:bg-neutral-100 transition-colors text-neutral-400"
                              >
                                <MoreHorizontal className="w-4 h-4" />
                              </button>
                              {menuOpen === user.id && (
                                <div className="absolute right-0 top-8 w-36 bg-white rounded-xl shadow-lg border border-neutral-100 z-20 overflow-hidden">
                                  <button
                                    className="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-neutral-700 hover:bg-neutral-50 transition-colors"
                                    onClick={(e) => { e.stopPropagation(); openEdit(user); }}
                                  >
                                    <Pencil className="w-3.5 h-3.5 text-neutral-400" /> Edit
                                  </button>
                                  <button
                                    className="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-red-500 hover:bg-red-50 transition-colors"
                                    onClick={(e) => { e.stopPropagation(); setMenuOpen(null); }}
                                  >
                                    <Trash2 className="w-3.5 h-3.5" /> Hapus
                                  </button>
                                </div>
                              )}
                            </div>
                          </>
                        ) : (
                          <>
                            <button
                              onClick={(e) => { e.stopPropagation(); openEdit(user); }}
                              className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-neutral-100 transition-colors text-neutral-400 hover:text-primary"
                            >
                              <Pencil className="w-4 h-4" />
                            </button>
                            <button
                              onClick={(e) => e.stopPropagation()}
                              className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-red-50 transition-colors text-neutral-400 hover:text-red-500"
                            >
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

        {/* Pagination Footer */}
        <div className="flex items-center justify-between px-6 py-4 border-t border-neutral-100">
          <span className="text-xs text-neutral-400">
            Menampilkan {filtered.length === 0 ? 0 : (safePage - 1) * PAGE_SIZE + 1}–
            {Math.min(safePage * PAGE_SIZE, filtered.length)} dari {filtered.length.toLocaleString('id-ID')} pengguna
          </span>

          <div className="flex items-center gap-1">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={safePage === 1}
              className="w-8 h-8 flex items-center justify-center rounded-lg border border-neutral-200 text-neutral-500 hover:bg-neutral-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
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
                  <button
                    key={p}
                    onClick={() => setPage(p as number)}
                    className={`w-8 h-8 flex items-center justify-center rounded-lg text-sm font-medium transition-colors ${
                      safePage === p
                        ? 'bg-primary text-white font-bold'
                        : 'border border-neutral-200 text-neutral-600 hover:bg-neutral-50'
                    }`}
                  >
                    {p}
                  </button>
                ),
              )}

            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={safePage === totalPages}
              className="w-8 h-8 flex items-center justify-center rounded-lg border border-neutral-200 text-neutral-500 hover:bg-neutral-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
              <ChevronRight className="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      {/* Edit Modal */}
      {editingUser && (
        <EditModal
          user={editingUser}
          onSave={handleSave}
          onClose={() => setEditingUser(null)}
        />
      )}
    </div>
  );
}
