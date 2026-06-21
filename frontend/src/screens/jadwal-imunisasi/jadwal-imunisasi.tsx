import { useState, useMemo, useEffect, useCallback } from 'react';
import { Search, Plus, Pencil, Trash2, Eye } from 'lucide-react';
import { useNotification } from '../../context/NotificationContext';
import { useAuth } from '../../context/AuthContext';
import { apiGet, apiDelete } from '../../lib/api';
import { useAppStore } from '../../store/useAppStore';
import type { ImunisasiListData, RiwayatImunisasiResponse } from '../../types/entities';
import type { Role } from '../../App';
import type { JadwalImunisasi } from './data/imunisasi.data';
import { ModalTambahJadwal } from './components/ModalTambahJadwal';
import { ModalUpdateJadwal } from './components/ModalUpdateJadwal';
import { ModalHapusJadwal } from './components/ModalHapusJadwal';
import { ModalDetailJadwal } from './components/ModalDetailJadwal';
import { Paginator } from '../../components/Paginator';
import { usePaginator } from '../../hooks/usePaginator';

const PAGE_SIZE = 10;

interface JadwalImunisasiProps {
  currentRole: Role;
}

export default function JadwalImunisasi({ currentRole }: JadwalImunisasiProps): JSX.Element {
  const notify = useNotification();
  const { user } = useAuth();
  const { setImunisasiLoading } = useAppStore();
  const [deleting, setDeleting] = useState(false);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState<'SEMUA' | 'SUDAH' | 'BELUM'>('SEMUA');

  // Server-side paginated data (local state, not zustand)
  const [items, setItems] = useState<JadwalImunisasi[]>([]);
  const [totalData, setTotalData] = useState(0);

  const { page, setPage, totalPages } = usePaginator({ totalItems: totalData, pageSize: PAGE_SIZE });

  const [modalTambah, setModalTambah] = useState(false);
  const [modalUpdate, setModalUpdate] = useState<JadwalImunisasi | null>(null);
  const [modalHapus, setModalHapus] = useState<JadwalImunisasi | null>(null);
  const [modalDetail, setModalDetail] = useState<JadwalImunisasi | null>(null);

  // ── Role flags ─────────────────────────────────────────────────────────────
  const canCRUD = currentRole === 'Bidan' || currentRole === 'Kader Posyandu';
  const isIbuWali = currentRole === 'Ibu/Wali';

  // ── Fetch data from backend ────────────────────────────────────────────────
  const fetchData = useCallback(async (targetPage: number, perPage: number = PAGE_SIZE) => {
    setImunisasiLoading(true);
    try {
      if (isIbuWali && user) {
        const res = await apiGet<RiwayatImunisasiResponse>(`/imunisasi/pasien/${user.idUser}`);
        const list = res.riwayat_imunisasi ?? [];
        const mapped: JadwalImunisasi[] = list.map((item) => ({
          id: String(item.id_imunisasi),
          idPasien: String(res.id_pasien),
          namaAnak: user.name,
          namaVaksin: item.nama_vaksin,
          dosis: 'Primary Dose',
          tanggalJadwal: item.tanggal_jadwal,
          tanggalRealisasi: item.tanggal_realisasi ?? null,
          status: item.status_imunisasi === 'Sudah' ? 'SUDAH' : 'BELUM',
        }));
        setItems(mapped);
        setTotalData(mapped.length);
      } else {
        const res = await apiGet<ImunisasiListData>('/imunisasi', { page: String(targetPage), per_page: String(perPage) });
        const list = res.jadwal ?? [];
        const mapped: JadwalImunisasi[] = list.map((item) => ({
          id: String(item.id_imunisasi),
          idPasien: String(item.id_imunisasi),
          namaAnak: item.nama_pasien,
          namaVaksin: item.nama_vaksin,
          dosis: 'Primary Dose',
          tanggalJadwal: item.tanggal_jadwal,
          tanggalRealisasi: null,
          status: item.status_imunisasi === 'Sudah' ? 'SUDAH' : 'BELUM',
        }));
        setItems(mapped);
        setTotalData(res.meta?.total ?? mapped.length);
      }
    } catch (err) {
      console.error('Gagal memuat jadwal imunisasi:', err);
    } finally {
      setImunisasiLoading(false);
    }
  }, [currentRole, user, isIbuWali, setImunisasiLoading]);

  useEffect(() => {
    fetchData(page, PAGE_SIZE);
  }, [page, fetchData]);

  // ── Derived data ───────────────────────────────────────────────────────────
  const filtered = useMemo(() => {
    return items.filter((j) => {
      const matchSearch =
        search === '' ||
        j.idPasien.toLowerCase().includes(search.toLowerCase()) ||
        j.namaVaksin.toLowerCase().includes(search.toLowerCase()) ||
        j.namaAnak.toLowerCase().includes(search.toLowerCase());
      const matchStatus = statusFilter === 'SEMUA' || j.status === statusFilter;
      return matchSearch && matchStatus;
    });
  }, [items, search, statusFilter]);

  const totalSudah = items.filter((j) => j.status === 'SUDAH').length;
  const totalBelum = items.filter((j) => j.status === 'BELUM').length;

  // ── Handlers ───────────────────────────────────────────────────────────────
  const handleTambah = () => {
    setModalTambah(false);
    setPage(1);
    notify.success('Jadwal imunisasi baru berhasil ditambahkan.');
  };

  const handleUpdate = () => {
    setModalUpdate(null);
    setPage(1);
    notify.success('Data jadwal imunisasi berhasil diperbarui.');
  };

  const handleHapus = async () => {
    if (!modalHapus) return;
    const idNum = parseInt(String(modalHapus.id).replace(/[^0-9]/g, ''), 10);
    if (!idNum) {
      notify.error('ID imunisasi tidak valid.');
      return;
    }

    setDeleting(true);
    try {
      await apiDelete(`/imunisasi/${idNum}`);
      setModalHapus(null);
      if (items.length <= 1 && page > 1) {
        setPage(page - 1);
      } else {
        fetchData(page, PAGE_SIZE);
      }
      notify.success('Jadwal imunisasi berhasil dihapus.');
    } catch {
      notify.error('Gagal menghapus jadwal imunisasi. Silakan coba lagi.');
    } finally {
      setDeleting(false);
    }
  };

  // ── Page heading info per role ─────────────────────────────────────────────
  const firstChildName = items.length > 0 ? items[0].namaAnak : 'Anak Anda';
  const pageDesc = isIbuWali
    ? `Jadwal imunisasi untuk ${firstChildName}. Pantau status dan tanggal pelaksanaan vaksin anak Anda.`
    : canCRUD
    ? 'Monitor dan kelola jadwal imunisasi komunitas. Tambah, edit, atau hapus catatan vaksinasi pasien.'
    : 'Monitor jadwal imunisasi komunitas. Data diperbarui secara real-time dari posyandu.';

  // ── Render ─────────────────────────────────────────────────────────────────
  return (
    <div className="space-y-6 font-body text-neutral-800">
      {/* ── Page heading ──────────────────────────────────────────────────── */}
      <div className="flex items-start justify-between gap-4">
        <div>
          <h2 className="text-xl font-bold text-neutral-900 font-headline">Jadwal Imunisasi</h2>
          <p className="text-xs text-neutral-500 mt-1 max-w-xl">{pageDesc}</p>
        </div>
        {/* Ibu/Wali: banner pasien */}
        {isIbuWali && items.length > 0 && (
          <div className="bg-primary-50 border border-primary-100 rounded-xl px-4 py-2.5 flex items-center gap-2 flex-shrink-0">
            <span className="text-primary text-base">👶</span>
            <div>
              <p className="text-xs font-bold text-primary font-body">{firstChildName}</p>
              <p className="text-[10px] text-neutral-500 font-body">ID: {items[0].idPasien}</p>
            </div>
          </div>
        )}
      </div>

      {/* ── Stat cards ────────────────────────────────────────────────────── */}
      <div className="grid grid-cols-3 gap-4">
        <div className="bg-white rounded-2xl border border-neutral-100 p-5 relative">
          <span className="absolute top-4 right-4 text-[10px] font-bold bg-primary text-white px-2 py-0.5 rounded-full">Total Stats</span>
          <div className="w-10 h-10 bg-primary-50 rounded-xl flex items-center justify-center mb-3">
            <span className="text-primary text-lg">📅</span>
          </div>
          <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-wide">Total Jadwal Imunisasi</p>
          <p className="text-4xl font-bold text-neutral-900 font-headline mt-1">{totalData}</p>
          <div className="h-0.5 bg-primary mt-3 rounded-full" />
        </div>

        <div className="bg-white rounded-2xl border border-neutral-100 p-5 relative">
          <span className="absolute top-4 right-4 text-[10px] font-bold bg-emerald-500 text-white px-2 py-0.5 rounded-full">
            {items.length > 0 ? `${Math.round((totalSudah / items.length) * 100)}%` : '0%'}
          </span>
          <div className="w-10 h-10 bg-emerald-50 rounded-xl flex items-center justify-center mb-3">
            <span className="text-emerald-500 text-xl">✓</span>
          </div>
          <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-wide">Jumlah Status 'Sudah'</p>
          <p className="text-4xl font-bold text-neutral-900 font-headline mt-1">{totalSudah}</p>
          <div className="h-0.5 bg-emerald-500 mt-3 rounded-full" />
        </div>

        <div className="bg-white rounded-2xl border border-neutral-100 p-5 relative">
          <span className="absolute top-4 right-4 text-[10px] font-bold bg-red-100 text-red-600 px-2 py-0.5 rounded-full">{totalBelum} left</span>
          <div className="w-10 h-10 bg-red-50 rounded-xl flex items-center justify-center mb-3">
            <span className="text-red-500 text-xl">📋</span>
          </div>
          <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-wide">Jumlah Status 'Belum'</p>
          <p className="text-4xl font-bold text-neutral-900 font-headline mt-1">{totalBelum}</p>
          <div className="h-0.5 bg-red-400 mt-3 rounded-full" />
        </div>
      </div>

      {/* ── Table ─────────────────────────────────────────────────────────── */}
      <div className="bg-white rounded-2xl border border-neutral-100 overflow-hidden">
        {/* Toolbar */}
        <div className="flex items-center gap-3 px-5 py-4 border-b border-neutral-50 flex-wrap">
          <div className="relative flex-1 min-w-48 max-w-xs">
            <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-400" />
            <input
              type="text"
              placeholder={isIbuWali ? 'Cari nama vaksin...' : 'Vaccine Name or Patient ID'}
              value={search}
              onChange={(e) => { setSearch(e.target.value); setPage(1); }}
              className="w-full pl-9 pr-4 py-2 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-700 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200"
            />
          </div>

          <select
            value={statusFilter}
            onChange={(e) => { setStatusFilter(e.target.value as typeof statusFilter); setPage(1); }}
            className="px-4 py-2 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-700 focus:outline-none focus:ring-2 focus:ring-primary-200"
          >
            <option value="SEMUA">Status Filter</option>
            <option value="SUDAH">Sudah</option>
            <option value="BELUM">Belum</option>
          </select>

          {/* Only Bidan & Kader see the + button */}
          {canCRUD && (
            <button
              onClick={() => setModalTambah(true)}
              className="flex items-center gap-2 bg-primary hover:bg-primary-600 text-white px-4 py-2 rounded-xl text-sm font-bold transition-colors ml-auto"
            >
              <Plus size={16} /> New Record
            </button>
          )}
        </div>

        {/* Table header */}
        <div className={`grid px-5 py-2.5 bg-neutral-50 border-b border-neutral-100 ${isIbuWali ? 'grid-cols-5' : 'grid-cols-6'}`}>
          {(isIbuWali
            ? ['Nama Vaksin', 'Dosis', 'Tanggal Jadwal', 'Tanggal Realisasi', 'Status']
            : ['ID Pasien', 'Nama Vaksin', 'Tanggal Jadwal', 'Tanggal Realisasi', 'Status', 'Aksi']
          ).map((h) => (
            <p key={h} className="text-[10px] font-bold text-neutral-400 uppercase tracking-wide">{h}</p>
          ))}
        </div>

        {filtered.length === 0 ? (
          <div className="py-16 text-center text-neutral-400">
            <p className="text-4xl mb-3">📅</p>
            <p className="text-sm font-semibold text-neutral-500">Tidak ada jadwal ditemukan</p>
            <p className="text-xs mt-1">Coba ubah filter atau kata kunci pencarian</p>
          </div>
        ) : (
          filtered.map((j) => (
            <div
              key={j.id}
              className={`grid px-5 py-4 border-b border-neutral-50 last:border-0 hover:bg-neutral-50/60 transition-colors items-start ${isIbuWali ? 'grid-cols-5' : 'grid-cols-6'}`}
            >
              {/* Ibu/Wali view: nama vaksin first */}
              {isIbuWali ? (
                <>
                  <div>
                    <p className="text-sm font-bold text-neutral-800">{j.namaVaksin.split(' (')[0]}</p>
                    <p className="text-xs text-neutral-400 mt-0.5">{j.namaVaksin.includes('(') ? j.namaVaksin.match(/\(([^)]+)\)/)?.[1] : ''}</p>
                  </div>
                  <p className="text-sm text-neutral-600">{j.dosis}</p>
                  <p className="text-sm text-neutral-600">{j.tanggalJadwal}</p>
                  <p className="text-sm text-neutral-500">
                    {j.tanggalRealisasi ?? <span className="text-neutral-300 italic">Pending...</span>}
                  </p>
                  <div>
                    <span className={`text-xs font-bold px-2.5 py-1 rounded-full ${j.status === 'SUDAH' ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-600'}`}>
                      {j.status}
                    </span>
                  </div>
                </>
              ) : (
                <>
                  {/* Non-Ibu/Wali view: ID first */}
                  <p className="text-sm font-bold text-primary">{j.namaAnak}</p>
                  <div>
                    <p className="text-sm font-bold text-neutral-800">{j.namaVaksin.split(' (')[0]}</p>
                    <p className="text-xs text-neutral-400 mt-0.5">{j.dosis}</p>
                  </div>
                  <p className="text-sm text-neutral-600">{j.tanggalJadwal}</p>
                  <p className="text-sm text-neutral-500">
                    {j.tanggalRealisasi ?? <span className="text-neutral-300 italic">Pending...</span>}
                  </p>
                  <div>
                    <span className={`text-xs font-bold px-2.5 py-1 rounded-full ${j.status === 'SUDAH' ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-600'}`}>
                      {j.status}
                    </span>
                  </div>
                  {/* Aksi column */}
                  <div className="flex items-center gap-2">
                    {canCRUD ? (
                      <>
                        <button
                          onClick={() => setModalUpdate(j)}
                          className="p-1.5 text-neutral-400 hover:text-primary hover:bg-primary-50 rounded-lg transition-colors"
                          title="Edit"
                        >
                          <Pencil size={14} />
                        </button>
                        <button
                          onClick={() => setModalHapus(j)}
                          className="p-1.5 text-neutral-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                          title="Hapus"
                        >
                          <Trash2 size={14} />
                        </button>
                      </>
                    ) : (
                      /* Dinkes: view only */
                      <button
                        onClick={() => setModalDetail(j)}
                        className="p-1.5 text-neutral-400 hover:text-primary hover:bg-primary-50 rounded-lg transition-colors"
                        title="Lihat Detail"
                      >
                        <Eye size={14} />
                      </button>
                    )}
                  </div>
                </>
              )}
            </div>
          ))
        )}

        {/* Pagination */}
        <Paginator
          page={page}
          totalPages={totalPages}
          totalItems={totalData}
          pageSize={PAGE_SIZE}
          onPageChange={setPage}
        />
      </div>

      {/* ── Bottom info cards — hidden for Ibu/Wali ───────────────────────── */}
      {!isIbuWali && (
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div className="bg-neutral-50 rounded-2xl border border-neutral-100 p-5">
            <h4 className="text-sm font-bold text-neutral-800 font-headline mb-1">Upcoming Campaigns</h4>
            <p className="text-xs text-neutral-500 leading-relaxed mb-4">
              Mass Measles-Rubella (MR) campaign scheduled for next month. Ensure all cold chain equipment is validated.
            </p>
            <div className="flex items-center gap-2">
              <div className="flex -space-x-2">
                {['bg-red-400', 'bg-blue-400', 'bg-emerald-400'].map((c, i) => (
                  <div key={i} className={`w-8 h-8 rounded-full border-2 border-white ${c} flex items-center justify-center text-white text-xs font-bold`}>
                    {['A', 'B', 'C'][i]}
                  </div>
                ))}
              </div>
              <span className="text-xs text-neutral-500 font-semibold">+4 Assigned Specialist Team</span>
            </div>
          </div>
          <div className="bg-primary-50 rounded-2xl border border-primary-100 p-5">
            <h4 className="text-sm font-bold text-primary font-headline mb-1">Vaccine Supply Chain</h4>
            <p className="text-xs text-neutral-600 leading-relaxed mb-4">
              Current stock of BCG and Polio vaccines is at 82%. Order replenishment soon to avoid delivery gaps.
            </p>
            <button className="text-xs text-primary font-bold hover:text-primary-600 transition-colors">
              Check Detailed Stock Inventory →
            </button>
          </div>
        </div>
      )}

      {/* ── Ibu/Wali: tips card ────────────────────────────────────────────── */}
      {isIbuWali && (
        <div className="bg-primary-50 border border-primary-100 rounded-2xl p-5 flex items-start gap-4">
          <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center flex-shrink-0">
            <span className="text-white text-lg">💉</span>
          </div>
          <div>
            <p className="text-sm font-bold text-primary font-headline">Jadwal Berikutnya</p>
            <p className="text-xs text-neutral-600 font-body mt-1 leading-relaxed">
              Pastikan Anda hadir tepat waktu di Posyandu untuk imunisasi berikutnya. Bawa buku KIA dan kartu imunisasi anak.
            </p>
          </div>
        </div>
      )}

      {/* ── FAB — only for Bidan & Kader ──────────────────────────────────── */}
      {canCRUD && (
        <button
          onClick={() => setModalTambah(true)}
          className="fixed bottom-8 right-8 w-14 h-14 bg-primary hover:bg-primary-600 text-white rounded-full shadow-lg flex items-center justify-center transition-colors z-20"
          title="Tambah Jadwal"
        >
          <Plus size={24} />
        </button>
      )}

      {/* ── Modals ────────────────────────────────────────────────────────── */}
      {modalTambah && canCRUD && (
        <ModalTambahJadwal onClose={() => setModalTambah(false)} onSimpan={handleTambah} />
      )}
      {modalUpdate && canCRUD && (
        <ModalUpdateJadwal jadwal={modalUpdate} onClose={() => setModalUpdate(null)} onSimpan={handleUpdate} />
      )}
      {modalHapus && canCRUD && (
        <ModalHapusJadwal jadwal={modalHapus} onClose={() => setModalHapus(null)} onHapus={handleHapus} loading={deleting} />
      )}
      {modalDetail && (
        <ModalDetailJadwal jadwal={modalDetail} onClose={() => setModalDetail(null)} />
      )}
    </div>
  );
}
