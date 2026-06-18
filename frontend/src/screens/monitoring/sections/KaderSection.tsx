import { useState, useEffect, useCallback } from 'react';
import { Plus, Search } from 'lucide-react';
import { useNotification } from '../../../context/NotificationContext';
import { apiGet, apiPost, apiPut, apiDelete } from '../../../lib/api';
import type { PasienListData } from '../../../types/entities';
import ModalTambahPemeriksaan, { type FormDataTambah, type PasienOption } from '../components/modaltambah';
import ModalEditPemeriksaan, { type PemeriksaanData } from '../components/modalubah';
import ModalHapusPemeriksaan from '../components/modalhapus';

interface PasienRow {
  no: number;
  id: string;
  nama: string;
  umur: string;
  nik: string;
}

export function KaderSection() {
  const notify = useNotification();
  const [search, setSearch] = useState('');
  const [pasienData, setPasienData] = useState<PasienRow[]>([]);
  const [pasienOptions, setPasienOptions] = useState<PasienOption[]>([]);

  const [modalTambah, setModalTambah] = useState(false);
  const [modalEdit, setModalEdit] = useState(false);
  const [editTarget, setEditTarget] = useState<{ nama: string; data: PemeriksaanData } | null>(null);
  const [modalHapus, setModalHapus] = useState(false);
  const [hapusTarget, setHapusTarget] = useState<{ id: string; nama: string; tanggal: string } | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  const loadPasienList = useCallback(async () => {
    if (pasienOptions.length > 0) return;
    try {
      const res = await apiGet<PasienListData>('/monitoring/pasien?page=1&per_page=50');
      const list = (res.pasien ?? []).map(p => ({
        id: String(p.id_pasien),
        nama: p.nama,
        namaIbu: '',
        usia: p.umur,
        nik: p.nik,
        statusGizi: '',
      }));
      setPasienOptions(list);
    } catch {
      notify.error('Gagal memuat daftar pasien.');
    }
  }, [pasienOptions.length, notify]);

  useEffect(() => {
    apiGet<PasienListData>('/monitoring/pasien?page=1&per_page=50')
      .then((res) => {
        const list = (res.pasien ?? []).map((p, idx) => ({
          no: idx + 1,
          id: String(p.id_pasien),
          nama: p.nama,
          umur: p.umur,
          nik: p.nik,
        }));
        setPasienData(list);
      })
      .catch(() => notify.error('Gagal memuat data pasien.'));
  }, []);

  const filteredData = pasienData.filter(row =>
    !search ||
    row.nama.toLowerCase().includes(search.toLowerCase()) ||
    row.nik.toLowerCase().includes(search.toLowerCase())
  );

  const openEdit = (row: PasienRow) => {
    setEditTarget({
      nama: row.nama,
      data: {
        id: row.id,
        beratBadan: '',
        tinggiBadan: '',
        lingkarKepala: '',
        terakhirDiperbarui: { nama: '', inisial: '', tanggal: '' },
      },
    });
    setModalEdit(true);
  };

  const openHapus = (row: PasienRow) => {
    setHapusTarget({ id: row.id, nama: row.nama, tanggal: '' });
    setModalHapus(true);
  };

  const handleTambahSubmit = async (_pasienId: string, namaAnak: string, idJadwalImunisasi: number, data: FormDataTambah) => {
    try {
      await apiPost('/monitoring/pemeriksaan', {
        id_jadwal_imunisasi: idJadwalImunisasi,
        berat_badan: parseFloat(data.beratBadan),
        tinggi_badan: parseFloat(data.tinggiBadan),
        lingkar_kepala: data.lingkarKepala ? parseFloat(data.lingkarKepala) : 0,
        tekanan_darah: data.tekananDarah || '-',
        catatan: data.catatanMedis || undefined,
      });
      notify.success(`Data pemeriksaan ${namaAnak} berhasil disimpan.`);
      setModalTambah(false);
    } catch {
      notify.error('Gagal menyimpan data pemeriksaan.');
    }
  };

  const handleEditSubmit = async (id: string, data: Omit<PemeriksaanData, 'id' | 'terakhirDiperbarui'>) => {
    const idNum = parseInt(String(id).replace(/[^0-9]/g, ''), 10) || 0;
    if (!idNum) { notify.error('ID pemeriksaan tidak valid.'); return; }
    if (!data.beratBadan || !data.tinggiBadan) {
      notify.warn('Berat badan dan tinggi badan wajib diisi.');
      return;
    }
    try {
      await apiPut('/monitoring/pemeriksaan/' + idNum, {
        berat_badan: parseFloat(data.beratBadan),
        tinggi_badan: parseFloat(data.tinggiBadan),
        lingkar_kepala: data.lingkarKepala ? parseFloat(data.lingkarKepala) : undefined,
      });
      notify.success('Data pemeriksaan berhasil diperbarui.');
      setModalEdit(false);
      setEditTarget(null);
    } catch {
      notify.error('Gagal memperbarui data pemeriksaan.');
    }
  };

  const handleHapusConfirm = async () => {
    if (!hapusTarget) return;
    setIsDeleting(true);
    const idNum = parseInt(String(hapusTarget.id).replace(/[^0-9]/g, ''), 10) || 0;
    try {
      await apiDelete('/monitoring/pemeriksaan/' + idNum);
      setPasienData(prev => prev.filter(p => p.id !== hapusTarget.id));
      setModalHapus(false);
      setHapusTarget(null);
      notify.success(`Data berhasil dihapus.`);
    } catch {
      notify.error('Gagal menghapus data.');
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap gap-3">
        <button
          onClick={() => { loadPasienList(); setModalTambah(true); }}
          className="inline-flex items-center gap-2 bg-primary hover:bg-primary-600 text-white rounded-xl px-5 py-2.5 text-sm font-medium transition-colors"
        >
          <Plus size={16} />
          Input Data Pengukuran
        </button>
      </div>

      <div className="bg-white rounded-2xl border border-neutral-100 overflow-hidden">
        <div className="flex items-center gap-3 px-5 py-4 border-b border-neutral-50">
          <div className="relative flex-1 max-w-xs">
            <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-400" />
            <input
              type="text"
              placeholder="Cari nama pasien..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-9 pr-4 py-2 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-700 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200"
            />
          </div>
          <p className="text-xs text-neutral-400 ml-auto">{filteredData.length} pasien</p>
        </div>

        <div className="grid grid-cols-5 px-5 py-2.5 bg-neutral-50 border-b border-neutral-100 text-[10px] font-bold text-neutral-400 uppercase tracking-wide">
          <p>NO</p>
          <p className="col-span-2">NAMA</p>
          <p>UMUR</p>
          <p>NIK</p>
        </div>

        {filteredData.length === 0 ? (
          <div className="py-12 text-center text-neutral-400">
            <p className="text-sm">Tidak ada data pasien</p>
          </div>
        ) : (
          filteredData.map((row) => (
            <div key={row.id} className="grid grid-cols-5 px-5 py-3 border-b border-neutral-50 last:border-0 hover:bg-neutral-50/60 items-center">
              <p className="text-sm text-neutral-400">{row.no}</p>
              <p className="col-span-2 text-sm font-bold text-primary">{row.nama}</p>
              <p className="text-sm text-neutral-600">{row.umur}</p>
              <div className="flex items-center gap-2">
                <p className="text-sm text-neutral-500">{row.nik}</p>
                <button onClick={() => openEdit(row)} className="text-xs text-emerald-600 hover:text-emerald-700 font-medium ml-auto">Edit</button>
                <button onClick={() => openHapus(row)} className="text-xs text-red-600 hover:text-red-700 font-medium">Hapus</button>
              </div>
            </div>
          ))
        )}
      </div>

      <ModalTambahPemeriksaan isOpen={modalTambah} onClose={() => setModalTambah(false)} pasienList={pasienOptions} onSubmit={handleTambahSubmit} />
      <ModalEditPemeriksaan isOpen={modalEdit} onClose={() => setModalEdit(false)} namaAnak={editTarget?.nama ?? ''} data={editTarget?.data ?? null} onSubmit={handleEditSubmit} />
      <ModalHapusPemeriksaan isOpen={modalHapus} onClose={() => { setModalHapus(false); setHapusTarget(null); }} namaAnak={hapusTarget?.nama ?? ''} tanggalPemeriksaan={hapusTarget?.tanggal ?? ''} onConfirm={handleHapusConfirm} isLoading={isDeleting} />
    </div>
  );
}
