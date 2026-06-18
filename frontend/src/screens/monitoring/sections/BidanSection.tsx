import { useState, useEffect, useCallback } from 'react';
import { Plus, Search } from 'lucide-react';
import { useNotification } from '../../../context/NotificationContext';
import { apiGet, apiPost, apiPut, apiPatch, apiDelete } from '../../../lib/api';
import type { PasienListData } from '../../../types/entities';
import ModalTambahPemeriksaan, { type FormDataTambah, type PasienOption } from '../components/modaltambah';
import ModalEditPemeriksaan, { type PemeriksaanData } from '../components/modalubah';
import ModalHapusPemeriksaan from '../components/modalhapus';
import { ModalVerifikasiBidan } from '../../../components/verifikasi/ModalVerifikasiBidan';
import type { VerifikasiTarget } from '../../../components/verifikasi/VerifikasiPanel';

interface PemeriksaanRow {
  no: number;
  id: string;
  nama: string;
  usia: string;
  statusGizi: string;
  verifikasi: string;
}

export function BidanSection() {
  const notify = useNotification();
  const [search, setSearch] = useState('');
  const [data, setData] = useState<PemeriksaanRow[]>([]);
  const [pasienOptions, setPasienOptions] = useState<PasienOption[]>([]);

  const [modalTambah, setModalTambah] = useState(false);
  const [modalEdit, setModalEdit] = useState(false);
  const [editTarget, setEditTarget] = useState<{ nama: string; data: PemeriksaanData } | null>(null);
  const [modalHapus, setModalHapus] = useState(false);
  const [hapusTarget, setHapusTarget] = useState<{ id: string; nama: string; tanggal: string } | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [verifikasiTarget, setVerifikasiTarget] = useState<VerifikasiTarget | null>(null);

  const loadPasienOptions = useCallback(async () => {
    if (pasienOptions.length > 0) return;
    try {
      const res = await apiGet<PasienListData>('/monitoring/pasien?page=1&per_page=50');
      const list = (res.pasien ?? []).map(p => ({
        id: String(p.id_pasien), nama: p.nama, namaIbu: '', usia: p.umur, nik: p.nik, statusGizi: '',
      }));
      setPasienOptions(list);
    } catch { notify.error('Gagal memuat daftar pasien.'); }
  }, [pasienOptions.length, notify]);

  useEffect(() => {
    apiGet<{ pemeriksaan_pending?: { id_hasil_pemeriksaan: number; nama_pasien: string }[] }>('/monitoring/pemeriksaan/pending')
      .then((res) => {
        const list = (res.pemeriksaan_pending ?? []).map((item, idx) => ({
          no: idx + 1,
          id: String(item.id_hasil_pemeriksaan),
          nama: item.nama_pasien,
          usia: '',
          statusGizi: '',
          verifikasi: 'Pending',
        }));
        setData(list);
      })
      .catch(() => console.error('Gagal memuat pemeriksaan pending'));
  }, []);

  const filteredData = data.filter(row =>
    !search || row.nama.toLowerCase().includes(search.toLowerCase())
  );

  const openEdit = (row: PemeriksaanRow) => {
    setEditTarget({
      nama: row.nama,
      data: { id: row.id, beratBadan: '', tinggiBadan: '', lingkarKepala: '', terakhirDiperbarui: { nama: '', inisial: '', tanggal: '' } },
    });
    setModalEdit(true);
  };

  const openHapus = (row: PemeriksaanRow) => {
    setHapusTarget({ id: row.id, nama: row.nama, tanggal: '' });
    setModalHapus(true);
  };

  const openVerifikasi = (row: PemeriksaanRow) => {
    setVerifikasiTarget({
      id: row.id, nama: row.nama, inisial: row.nama.split(' ').map(w => w[0]).slice(0, 2).join(''),
      warnaBg: 'bg-primary', warnaText: 'text-white', usia: row.usia, bb: '', tb: '',
      petugas: '', statusGizi: row.statusGizi,
    });
  };

  const handleVerifikasiSetuju = async (id: string, _catatan: string) => {
    try {
      await apiPatch('/monitoring/pemeriksaan/' + id + '/verify');
      setData(prev => prev.map(b => b.id === id ? { ...b, verifikasi: 'Verified' } : b));
      setVerifikasiTarget(null);
      notify.success('Data pemeriksaan berhasil diverifikasi.');
    } catch { notify.error('Gagal memverifikasi.'); }
  };

  const handleVerifikasiTolak = async (id: string, _catatan: string) => {
    try {
      await apiDelete('/monitoring/pemeriksaan/' + id);
      setData(prev => prev.filter(b => b.id !== id));
      setVerifikasiTarget(null);
      notify.warn('Data pemeriksaan ditolak.');
    } catch { notify.error('Gagal menolak.'); }
  };

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
      apiGet<{ pemeriksaan_pending?: { id_hasil_pemeriksaan: number; nama_pasien: string }[] }>('/monitoring/pemeriksaan/pending')
        .then((res) => {
          const list = (res.pemeriksaan_pending ?? []).map((item, idx) => ({
            no: idx + 1,
            id: String(item.id_hasil_pemeriksaan),
            nama: item.nama_pasien,
            usia: '',
            statusGizi: '',
            verifikasi: 'Pending',
          }));
          setData(list);
        })
        .catch(() => {});
    } catch { notify.error('Gagal menyimpan.'); }
  };

  const handleEditSubmit = async (id: string, dataForm: Omit<PemeriksaanData, 'id' | 'terakhirDiperbarui'>) => {
    const idNum = parseInt(String(id).replace(/[^0-9]/g, ''), 10) || 0;
    if (!idNum) { notify.error('ID pemeriksaan tidak valid.'); return; }
    if (!dataForm.beratBadan || !dataForm.tinggiBadan) {
      notify.warn('Berat badan dan tinggi badan wajib diisi.');
      return;
    }
    try {
      await apiPut('/monitoring/pemeriksaan/' + idNum, {
        berat_badan: parseFloat(dataForm.beratBadan),
        tinggi_badan: parseFloat(dataForm.tinggiBadan),
        lingkar_kepala: dataForm.lingkarKepala ? parseFloat(dataForm.lingkarKepala) : undefined,
      });
      notify.success('Data berhasil diperbarui.');
      setModalEdit(false);
      setEditTarget(null);
    } catch { notify.error('Gagal memperbarui.'); }
  };

  const handleHapusConfirm = async () => {
    if (!hapusTarget) return;
    setIsDeleting(true);
    const idNum = parseInt(String(hapusTarget.id).replace(/[^0-9]/g, ''), 10) || 0;
    try {
      await apiDelete('/monitoring/pemeriksaan/' + idNum);
      setData(prev => prev.filter(b => b.id !== hapusTarget.id));
      setModalHapus(false);
      setHapusTarget(null);
      notify.success('Data berhasil dihapus.');
    } catch { notify.error('Gagal menghapus.'); }
    finally { setIsDeleting(false); }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap gap-3">
        <button onClick={() => { loadPasienOptions(); setModalTambah(true); }}
          className="inline-flex items-center gap-2 bg-primary hover:bg-primary-600 text-white rounded-xl px-5 py-2.5 text-sm font-medium transition-colors">
          <Plus size={16} /> Tambah Data Pengukuran
        </button>
      </div>

      <div className="bg-white rounded-2xl border border-neutral-100 overflow-hidden">
        <div className="flex items-center gap-3 px-5 py-4 border-b border-neutral-50">
          <div className="relative flex-1 max-w-xs">
            <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-400" />
            <input type="text" placeholder="Cari nama pasien..." value={search} onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-9 pr-4 py-2 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-700 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200" />
          </div>
        </div>

        <div className="grid grid-cols-5 px-5 py-2.5 bg-neutral-50 border-b border-neutral-100 text-[10px] font-bold text-neutral-400 uppercase tracking-wide">
          <p>NO</p>
          <p className="col-span-2">NAMA</p>
          <p>STATUS</p>
          <p>AKSI</p>
        </div>

        {filteredData.length === 0 ? (
          <div className="py-12 text-center text-neutral-400"><p className="text-sm">Tidak ada data pemeriksaan</p></div>
        ) : (
          filteredData.map((row) => (
            <div key={row.id} className="grid grid-cols-5 px-5 py-3 border-b border-neutral-50 last:border-0 hover:bg-neutral-50/60 items-center">
              <p className="text-sm text-neutral-400">{row.no}</p>
              <p className="col-span-2 text-sm font-bold text-primary">{row.nama}</p>
              <span className={`text-xs font-bold px-2 py-1 rounded-full w-fit ${row.verifikasi === 'Verified' ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-700'}`}>
                {row.verifikasi === 'Verified' ? 'Terverifikasi' : 'Pending'}
              </span>
              <div className="flex items-center gap-2">
                {row.verifikasi === 'Pending' && (
                  <button onClick={() => openVerifikasi(row)}
                    className="text-xs text-blue-700 hover:text-blue-800 bg-blue-50 hover:bg-blue-100 px-2 py-1.5 rounded-lg font-medium">Verifikasi</button>
                )}
                <button onClick={() => openEdit(row)} className="text-xs text-emerald-600 hover:text-emerald-700 font-medium">Edit</button>
                <button onClick={() => openHapus(row)} className="text-xs text-red-600 hover:text-red-700 font-medium">Hapus</button>
              </div>
            </div>
          ))
        )}
      </div>

      <ModalTambahPemeriksaan isOpen={modalTambah} onClose={() => setModalTambah(false)} pasienList={pasienOptions} onSubmit={handleTambahSubmit} />
      <ModalEditPemeriksaan isOpen={modalEdit} onClose={() => setModalEdit(false)} namaAnak={editTarget?.nama ?? ''} data={editTarget?.data ?? null} onSubmit={handleEditSubmit} />
      <ModalHapusPemeriksaan isOpen={modalHapus} onClose={() => { setModalHapus(false); setHapusTarget(null); }} namaAnak={hapusTarget?.nama ?? ''} tanggalPemeriksaan={hapusTarget?.tanggal ?? ''} onConfirm={handleHapusConfirm} isLoading={isDeleting} />
      {verifikasiTarget && (
        <ModalVerifikasiBidan target={verifikasiTarget} onClose={() => setVerifikasiTarget(null)} onSetuju={handleVerifikasiSetuju} onTolak={handleVerifikasiTolak} />
      )}
    </div>
  );
}
