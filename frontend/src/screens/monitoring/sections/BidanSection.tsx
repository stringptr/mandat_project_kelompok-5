import { useState, useEffect, useCallback } from 'react';
import { Plus, Search, ClipboardList, Users, FileSearch, AlertOctagon, Eye, X, Pencil, Trash2 } from 'lucide-react';
import { useAuth } from '../../../context/AuthContext';
import { useNotification } from '../../../context/NotificationContext';
import { apiGet, apiPost, apiPut, apiPatch, apiDelete } from '../../../lib/api';
import { Paginator } from '../../../components/Paginator';
import type { PasienListData } from '../../../types/entities';
import ModalTambahPemeriksaan, { type FormDataTambah, type PasienOption } from '../components/modaltambah';
import ModalEditPemeriksaan, { type PemeriksaanData } from '../components/modalubah';
import ModalHapusPemeriksaan from '../components/modalhapus';
import { ModalVerifikasiBidan } from '../../../components/verifikasi/ModalVerifikasiBidan';
import type { VerifikasiTarget } from '../../../components/verifikasi/VerifikasiPanel';

interface PemeriksaanItem {
  id_hasil_pemeriksaan: number;
  id_jadwal_imunisasi: number;
  nama_vaksin: string;
  nama_pasien: string;
  berat_badan: number;
  tinggi_badan: number;
  lingkar_kepala: number;
  tekanan_darah: string;
  status_stunting: string;
  status_gizi: string;
  catatan?: string;
  tanggal: string;
  petugas: string;
}

interface PendingRow {
  id: string;
  nama: string;
}

export function BidanSection() {
  const { user } = useAuth();
  const notify = useNotification();
  const [search, setSearch] = useState('');

  // Pemeriksaan table
  const [pemeriksaan, setPemeriksaan] = useState<PemeriksaanItem[]>([]);
  const [pemPage, setPemPage] = useState(1);
  const [pemTotal, setPemTotal] = useState(0);
  const [pemLastPage, setPemLastPage] = useState(1);
  const [pendingCount, setPendingCount] = useState(0);
  const [pendingList, setPendingList] = useState<PendingRow[]>([]);
  const PEM_PER_PAGE = 15;

  // Pasien options for tambah modal
  const [pasienOptions, setPasienOptions] = useState<PasienOption[]>([]);

  // Modal states
  const [modalTambah, setModalTambah] = useState(false);
  const [modalEdit, setModalEdit] = useState(false);
  const [editTarget, setEditTarget] = useState<{ nama: string; data: PemeriksaanData } | null>(null);
  const [modalHapus, setModalHapus] = useState(false);
  const [hapusTarget, setHapusTarget] = useState<{ id: string; nama: string; tanggal: string } | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [verifikasiTarget, setVerifikasiTarget] = useState<VerifikasiTarget | null>(null);
  const [detailTarget, setDetailTarget] = useState<PemeriksaanItem | null>(null);

  // Fetch pemeriksaan list
  const fetchPemeriksaan = useCallback((p: number) => {
    const params: Record<string, string> = { page: String(p), per_page: String(PEM_PER_PAGE) };
    if (user?.idUser) params.id_bidan = String(user.idUser);
    apiGet<{ pemeriksaan: PemeriksaanItem[]; meta: { current_page: number; per_page: number; total: number; last_page: number } }>(
      '/monitoring/semua-pemeriksaan', params
    )
      .then((res) => {
        setPemeriksaan(res.pemeriksaan ?? []);
        setPemTotal(res.meta?.total ?? 0);
        setPemLastPage(res.meta?.last_page ?? 1);
      })
      .catch(() => {});
  }, []);

  // Fetch pending count
  useEffect(() => {
    apiGet<{ pemeriksaan_pending: { id_hasil_pemeriksaan: number; nama_pasien: string }[] }>('/monitoring/pemeriksaan/pending')
      .then((res) => {
        const list = res.pemeriksaan_pending ?? [];
        setPendingCount(list.length);
        setPendingList(list.map((item) => ({ id: String(item.id_hasil_pemeriksaan), nama: item.nama_pasien })));
      })
      .catch(() => {});
  }, []);

  useEffect(() => { fetchPemeriksaan(pemPage); }, [pemPage, fetchPemeriksaan]);

  // Pasien options
  const loadPasienOptions = useCallback(async () => {
    if (pasienOptions.length > 0) return;
    try {
      const res = await apiGet<PasienListData>('/monitoring/pasien?page=1&per_page=50');
      setPasienOptions((res.pasien ?? []).map((p) => ({
        id: String(p.id_pasien), nama: p.nama, namaIbu: '', usia: p.umur, nik: p.nik, statusGizi: '',
      })));
    } catch { notify.error('Gagal memuat daftar pasien.'); }
  }, [pasienOptions.length, notify]);

  // Filtered data
  const filtered = pemeriksaan.filter((row) =>
    !search || row.nama_pasien.toLowerCase().includes(search.toLowerCase())
  );

  // Verif logic
  const openVerifikasi = (row: PemeriksaanItem) => {
    setVerifikasiTarget({
      id: String(row.id_hasil_pemeriksaan), nama: row.nama_pasien,
      inisial: row.nama_pasien.split(' ').map((w) => w[0]).slice(0, 2).join(''),
      warnaBg: 'bg-primary', warnaText: 'text-white', usia: '',
      bb: String(row.berat_badan), tb: String(row.tinggi_badan),
      petugas: row.petugas, statusGizi: row.status_gizi,
    });
  };

  const handleVerifikasiSetuju = async (id: string, _catatan: string) => {
    try {
      await apiPatch('/monitoring/pemeriksaan/' + id + '/verify');
      setVerifikasiTarget(null);
      setPendingCount((c) => c - 1);
      setPendingList((l) => l.filter((p) => p.id !== id));
      fetchPemeriksaan(pemPage);
      notify.success('Data berhasil diverifikasi.');
    } catch { notify.error('Gagal memverifikasi.'); }
  };

  const handleVerifikasiTolak = async (id: string, _catatan: string) => {
    try {
      await apiDelete('/monitoring/pemeriksaan/' + id);
      setVerifikasiTarget(null);
      setPendingCount((c) => c - 1);
      setPendingList((l) => l.filter((p) => p.id !== id));
      fetchPemeriksaan(pemPage);
      notify.warn('Data ditolak.');
    } catch { notify.error('Gagal menolak.'); }
  };

  // CRUD Handlers
  const handleTambahSubmit = async (_pasienId: string, namaAnak: string, _idJadwal: number, dataForm: FormDataTambah) => {
    try {
      await apiPost('/monitoring/pemeriksaan', {
        id_jadwal_imunisasi: _idJadwal,
        berat_badan: parseFloat(dataForm.beratBadan),
        tinggi_badan: parseFloat(dataForm.tinggiBadan),
        lingkar_kepala: dataForm.lingkarKepala ? parseFloat(dataForm.lingkarKepala) : 0,
        tekanan_darah: dataForm.tekananDarah || '-',
        catatan: dataForm.catatanMedis || undefined,
      });
      notify.success(`Data ${namaAnak} berhasil disimpan.`);
      setModalTambah(false);
      fetchPemeriksaan(pemPage);
    } catch { notify.error('Gagal menyimpan.'); }
  };

  const openEdit = (row: PemeriksaanItem) => {
    setEditTarget({
      nama: row.nama_pasien,
      data: {
        id: String(row.id_hasil_pemeriksaan), beratBadan: String(row.berat_badan),
        tinggiBadan: String(row.tinggi_badan), lingkarKepala: String(row.lingkar_kepala),
        terakhirDiperbarui: { nama: row.petugas, inisial: row.petugas?.[0] || '?', tanggal: row.tanggal?.slice(0, 10) },
      },
    });
    setModalEdit(true);
  };

  const handleEditSubmit = async (id: string, dataForm: Omit<PemeriksaanData, 'id' | 'terakhirDiperbarui'>) => {
    const idNum = parseInt(String(id).replace(/[^0-9]/g, ''), 10) || 0;
    if (!idNum) { notify.error('ID pemeriksaan tidak valid.'); return; }
    if (!dataForm.beratBadan || !dataForm.tinggiBadan) { notify.warn('BB dan TB wajib diisi.'); return; }
    try {
      await apiPut('/monitoring/pemeriksaan/' + idNum, {
        berat_badan: parseFloat(dataForm.beratBadan),
        tinggi_badan: parseFloat(dataForm.tinggiBadan),
        lingkar_kepala: dataForm.lingkarKepala ? parseFloat(dataForm.lingkarKepala) : undefined,
      });
      notify.success('Data berhasil diperbarui.');
      setModalEdit(false); setEditTarget(null);
      fetchPemeriksaan(pemPage);
    } catch { notify.error('Gagal memperbarui.'); }
  };

  const openHapus = (row: PemeriksaanItem) => {
    setHapusTarget({ id: String(row.id_hasil_pemeriksaan), nama: row.nama_pasien, tanggal: row.tanggal?.slice(0, 10) });
    setModalHapus(true);
  };

  const handleHapusConfirm = async () => {
    if (!hapusTarget) return;
    setIsDeleting(true);
    try {
      await apiDelete('/monitoring/pemeriksaan/' + hapusTarget.id);
      setModalHapus(false); setHapusTarget(null);
      notify.success('Data berhasil dihapus.');
      fetchPemeriksaan(pemPage);
    } catch { notify.error('Gagal menghapus.'); }
    finally { setIsDeleting(false); }
  };

  return (
    <div className="space-y-6">
      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div className="bg-primary rounded-2xl p-4 relative overflow-hidden">
          <div className="absolute -bottom-4 -right-4 w-20 h-20 rounded-full bg-white/10" />
          <Users size={20} className="text-white/80 mb-1.5" />
          <p className="text-white/70 text-[10px] font-semibold uppercase tracking-wide">Total Pemeriksaan</p>
          <p className="text-2xl font-bold text-white font-headline mt-0.5">{pemTotal.toLocaleString('id-ID')}</p>
        </div>
        <div className="bg-white rounded-2xl p-4 border border-neutral-100">
          <FileSearch size={20} className="text-amber-500 mb-1.5" />
          <p className="text-neutral-400 text-[10px] font-semibold uppercase tracking-wide">Perlu Verifikasi</p>
          <p className="text-2xl font-bold text-amber-500 font-headline mt-0.5">{pendingCount}</p>
        </div>
        <div className="bg-white rounded-2xl p-4 border border-neutral-100">
          <AlertOctagon size={20} className="text-red-500 mb-1.5" />
          <p className="text-neutral-400 text-[10px] font-semibold uppercase tracking-wide">Stunting Berat</p>
          <p className="text-2xl font-bold text-red-500 font-headline mt-0.5">
            {pemeriksaan.filter((p) => p.status_stunting?.includes('Berat')).length}
          </p>
        </div>
      </div>

      {/* Action Bar */}
      <div className="flex items-center gap-3">
        <button onClick={() => { loadPasienOptions(); setModalTambah(true); }}
          className="inline-flex items-center gap-2 bg-primary hover:bg-primary-600 text-white rounded-xl px-5 py-2.5 text-sm font-medium transition-colors">
          <Plus size={16} /> Tambah Data Pengukuran
        </button>
        {pendingList.length > 0 && (
          <span className="text-xs text-amber-600 font-medium">{pendingList.length} pemeriksaan menunggu verifikasi</span>
        )}
      </div>

      {/* Penimbangan Table */}
      <div>
        <div className="flex items-center gap-2 mb-3">
          <ClipboardList size={18} className="text-primary" />
          <h3 className="text-base font-bold text-neutral-800 font-headline">Tabel Penimbangan</h3>
          <span className="text-xs text-neutral-400 ml-auto">{pemTotal.toLocaleString('id-ID')} data</span>
        </div>

        <div className="bg-white rounded-2xl border border-neutral-100 overflow-hidden">
          <div className="flex items-center gap-3 px-5 py-3 border-b border-neutral-50">
            <div className="relative flex-1 max-w-xs">
              <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-400" />
              <input type="text" placeholder="Cari nama pasien..." value={search} onChange={(e) => setSearch(e.target.value)}
                className="w-full pl-9 pr-4 py-1.5 bg-neutral-50 border border-neutral-200 rounded-lg text-xs text-neutral-700 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200" />
            </div>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-xs">
              <thead>
                <tr className="bg-neutral-50 border-b border-neutral-100 text-[10px] font-bold text-neutral-400 uppercase tracking-wide">
                  <th className="text-left py-2.5 px-4 whitespace-nowrap">Nama</th>
                  <th className="text-left py-2.5 px-4 whitespace-nowrap">Vaksin</th>
                  <th className="text-center py-2.5 px-3 whitespace-nowrap">BB</th>
                  <th className="text-center py-2.5 px-3 whitespace-nowrap">TB</th>
                  <th className="text-center py-2.5 px-3 whitespace-nowrap">Stunting</th>
                  <th className="text-center py-2.5 px-3 whitespace-nowrap">Gizi</th>
                  <th className="text-center py-2.5 px-3 whitespace-nowrap">Status</th>
                  <th className="text-center py-2.5 px-3 whitespace-nowrap">Tgl</th>
                  <th className="text-center py-2.5 px-3 whitespace-nowrap">Aksi</th>
                </tr>
              </thead>
              <tbody>
                {filtered.length === 0 ? (
                  <tr><td colSpan={9} className="py-12 text-center text-neutral-400">Tidak ada data penimbangan</td></tr>
                ) : (
                  filtered.map((row) => {
                    const isPending = pendingList.some((p) => p.id === String(row.id_hasil_pemeriksaan));
                    return (
                      <tr key={row.id_hasil_pemeriksaan} className="border-b border-neutral-50 hover:bg-neutral-50/60 transition-colors">
                        <td className="py-2.5 px-4 font-semibold text-neutral-800 whitespace-nowrap max-w-40 truncate">{row.nama_pasien}</td>
                        <td className="py-2.5 px-4 text-blue-600 font-medium whitespace-nowrap max-w-32 truncate">{row.nama_vaksin || '-'}</td>
                        <td className="py-2.5 px-3 text-center text-neutral-700 tabular-nums whitespace-nowrap">{row.berat_badan?.toFixed(1)} kg</td>
                        <td className="py-2.5 px-3 text-center text-neutral-700 tabular-nums whitespace-nowrap">{row.tinggi_badan?.toFixed(1)} cm</td>
                        <td className="py-2.5 px-3 text-center whitespace-nowrap">
                          <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${row.status_stunting?.includes('Stunting') ? 'bg-red-100 text-red-600' : 'bg-green-100 text-green-600'}`}>{row.status_stunting || '-'}</span>
                        </td>
                        <td className="py-2.5 px-3 text-center whitespace-nowrap">
                          <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${row.status_gizi?.includes('Buruk') || row.status_gizi?.includes('Kurang') ? 'bg-amber-100 text-amber-600' : 'bg-emerald-100 text-emerald-600'}`}>{row.status_gizi || '-'}</span>
                        </td>
                        <td className="py-2.5 px-3 text-center whitespace-nowrap">
                          {isPending ? (
                            <button onClick={() => openVerifikasi(row)} className="text-[10px] font-bold text-amber-700 bg-amber-100 hover:bg-amber-200 px-2 py-0.5 rounded-full">Verifikasi</button>
                          ) : (
                            <span className="text-[10px] font-medium text-emerald-700 bg-emerald-100 px-2 py-0.5 rounded-full">✓</span>
                          )}
                        </td>
                        <td className="py-2.5 px-3 text-center text-neutral-500 whitespace-nowrap">{row.tanggal?.slice(0, 10)}</td>
                        <td className="py-2.5 px-3 text-center whitespace-nowrap">
                          <div className="flex items-center justify-center gap-2">
                            <button onClick={() => setDetailTarget(row)} className="p-1.5 hover:bg-primary/10 rounded-lg text-primary transition-colors" title="Detail">
                              <Eye size={14} />
                            </button>
                            <button onClick={() => openEdit(row)} className="p-1.5 hover:bg-blue-50 rounded-lg text-blue-600 transition-colors" title="Edit">
                              <Pencil size={14} />
                            </button>
                            <button onClick={() => openHapus(row)} className="p-1.5 hover:bg-red-50 rounded-lg text-red-500 transition-colors" title="Hapus">
                              <Trash2 size={14} />
                            </button>
                          </div>
                        </td>
                      </tr>
                    );
                  })
                )}
              </tbody>
            </table>
          </div>

          <Paginator page={pemPage} totalPages={pemLastPage} totalItems={pemTotal} pageSize={PEM_PER_PAGE} onPageChange={setPemPage} />
        </div>
      </div>

      {/* Detail Modal */}
      {detailTarget && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4" onClick={() => setDetailTarget(null)}>
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg max-h-[80vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
            <div className="px-6 py-4 border-b border-neutral-100 flex items-center justify-between">
              <h3 className="text-lg font-bold text-neutral-800 font-headline">Detail Pemeriksaan</h3>
              <button onClick={() => setDetailTarget(null)} className="p-2 hover:bg-neutral-100 rounded-xl">
                <X size={18} className="text-neutral-500" />
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <p className="text-xs text-neutral-400 uppercase tracking-wide">Pasien</p>
                <p className="text-sm font-semibold text-neutral-800">{detailTarget.nama_pasien}</p>
              </div>
              <div className="grid grid-cols-3 gap-4">
                <div>
                  <p className="text-xs text-neutral-400 uppercase tracking-wide">BB</p>
                  <p className="text-sm font-semibold text-neutral-800">{detailTarget.berat_badan?.toFixed(1)} kg</p>
                </div>
                <div>
                  <p className="text-xs text-neutral-400 uppercase tracking-wide">TB</p>
                  <p className="text-sm font-semibold text-neutral-800">{detailTarget.tinggi_badan?.toFixed(1)} cm</p>
                </div>
                <div>
                  <p className="text-xs text-neutral-400 uppercase tracking-wide">LK</p>
                  <p className="text-sm font-semibold text-neutral-800">{detailTarget.lingkar_kepala?.toFixed(1)} cm</p>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-xs text-neutral-400 uppercase tracking-wide">Tekanan Darah</p>
                  <p className="text-sm font-semibold text-neutral-800">{detailTarget.tekanan_darah || '-'}</p>
                </div>
                <div>
                  <p className="text-xs text-neutral-400 uppercase tracking-wide">Vaksin</p>
                  <p className="text-sm font-semibold text-blue-600">{detailTarget.nama_vaksin || '-'}</p>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-xs text-neutral-400 uppercase tracking-wide">Status Stunting</p>
                  <span className={`inline-block mt-1 text-xs font-bold px-2 py-0.5 rounded-full ${detailTarget.status_stunting?.includes('Stunting') ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>{detailTarget.status_stunting || '-'}</span>
                </div>
                <div>
                  <p className="text-xs text-neutral-400 uppercase tracking-wide">Status Gizi</p>
                  <span className={`inline-block mt-1 text-xs font-bold px-2 py-0.5 rounded-full ${detailTarget.status_gizi?.includes('Buruk') || detailTarget.status_gizi?.includes('Kurang') ? 'bg-amber-100 text-amber-700' : 'bg-emerald-100 text-emerald-700'}`}>{detailTarget.status_gizi || '-'}</span>
                </div>
              </div>
              <div>
                <p className="text-xs text-neutral-400 uppercase tracking-wide">Catatan</p>
                <p className="text-sm text-neutral-700 bg-neutral-50 rounded-xl p-3 mt-1">{detailTarget.catatan || 'Tidak ada catatan'}</p>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-xs text-neutral-400 uppercase tracking-wide">Petugas</p>
                  <p className="text-sm font-medium text-neutral-800">{detailTarget.petugas || '-'}</p>
                </div>
                <div>
                  <p className="text-xs text-neutral-400 uppercase tracking-wide">Tanggal</p>
                  <p className="text-sm font-medium text-neutral-800">{detailTarget.tanggal?.slice(0, 10)}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Modals */}
      <ModalTambahPemeriksaan isOpen={modalTambah} onClose={() => setModalTambah(false)} pasienList={pasienOptions} onSubmit={handleTambahSubmit} />
      <ModalEditPemeriksaan isOpen={modalEdit} onClose={() => setModalEdit(false)} namaAnak={editTarget?.nama ?? ''} data={editTarget?.data ?? null} onSubmit={handleEditSubmit} />
      <ModalHapusPemeriksaan isOpen={modalHapus} onClose={() => { setModalHapus(false); setHapusTarget(null); }} namaAnak={hapusTarget?.nama ?? ''} tanggalPemeriksaan={hapusTarget?.tanggal ?? ''} onConfirm={handleHapusConfirm} isLoading={isDeleting} />
      {verifikasiTarget && (
        <ModalVerifikasiBidan target={verifikasiTarget} onClose={() => setVerifikasiTarget(null)} onSetuju={handleVerifikasiSetuju} onTolak={handleVerifikasiTolak} />
      )}
    </div>
  );
}
