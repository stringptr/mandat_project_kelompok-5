import { useState } from 'react';
import { X, Search, CheckCircle, Clock } from 'lucide-react';
import { useNotification } from '../../../context/NotificationContext';
import { apiPost } from '../../../lib/api';
import type { JadwalImunisasi } from '../data/imunisasi.data';
import { VAKSIN_OPTIONS } from '../data/imunisasi.data';

interface ModalTambahJadwalProps {
  onClose: () => void;
  onSimpan: (data: Omit<JadwalImunisasi, 'id'>) => void;
}

export function ModalTambahJadwal({ onClose, onSimpan }: ModalTambahJadwalProps): JSX.Element {
  const notify = useNotification();
  const [idPasien, setIdPasien] = useState('');
  const [namaAnak] = useState('');
  const [namaVaksin, setNamaVaksin] = useState('');
  const [tanggalJadwal, setTanggalJadwal] = useState('');
  const [status, setStatus] = useState<'BELUM' | 'SUDAH'>('BELUM');
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validate = () => {
    const e: Record<string, string> = {};
    if (!idPasien.trim()) e.idPasien = 'ID/Nama pasien wajib diisi';
    if (!namaVaksin) e.namaVaksin = 'Nama vaksin wajib dipilih';
    if (!tanggalJadwal) e.tanggalJadwal = 'Tanggal jadwal wajib diisi';
    return e;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const errs = validate();
    if (Object.keys(errs).length > 0) { setErrors(errs); notify.warn('Mohon lengkapi semua data form yang wajib diisi sebelum mengirim.'); return; }
    try {
      const pid = parseInt(idPasien.replace(/[^0-9]/g, ''), 10) || 0;
      await apiPost('/imunisasi', {
        id_pasien: pid,
        nama_vaksin: namaVaksin,
        tanggal_jadwal: tanggalJadwal,
      });
      onSimpan({
        idPasien: idPasien.startsWith('#') ? idPasien : `#${idPasien}`,
        namaAnak: namaAnak || idPasien,
        namaVaksin,
        dosis: 'Primary Dose',
        tanggalJadwal: new Date(tanggalJadwal).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }),
        tanggalRealisasi: status === 'SUDAH' ? new Date(tanggalJadwal).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }) : null,
        status,
      });
    } catch {
      notify.error('Gagal menambahkan jadwal. Silakan coba lagi.');
    }
  };

  return (
    <div className="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md font-body">
        {/* Header */}
        <div className="flex items-center gap-3 px-6 py-5 border-b border-neutral-100">
          <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center flex-shrink-0">
            <span className="text-white text-lg">📅</span>
          </div>
          <div className="flex-1">
            <h2 className="text-base font-bold text-neutral-900 font-headline">Tambah Jadwal Imunisasi</h2>
            <p className="text-xs text-neutral-500 mt-0.5">Input data id vaksinasi rutin anak</p>
          </div>
          <button onClick={onClose} className="p-1.5 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg transition-colors">
            <X size={18} />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="px-6 py-5 space-y-4">
          {/* ID / Nama Anak */}
          <div>
            <label className="text-xs font-semibold text-neutral-600 block mb-1.5">ID Pasien / Nama Anak</label>
            <div className="relative">
              <Search size={15} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-neutral-400" />
              <input
                type="text"
                value={idPasien}
                onChange={(e) => setIdPasien(e.target.value)}
                placeholder="Cari ID atau Nama Anak..."
                className={`w-full pl-10 pr-4 py-2.5 bg-neutral-50 border rounded-xl text-sm text-neutral-800 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 focus:border-primary transition-colors ${errors.idPasien ? 'border-red-300' : 'border-neutral-200'}`}
              />
            </div>
            <p className="text-[10px] text-neutral-400 mt-1">Contoh: 'PST-000' atau 'Ani Sulistyowati'</p>
            {errors.idPasien && <p className="text-xs text-red-500 mt-1">{errors.idPasien}</p>}
          </div>

          {/* Nama Vaksin */}
          <div>
            <label className="text-xs font-semibold text-neutral-600 block mb-1.5">Nama Vaksin</label>
            <select
              value={namaVaksin}
              onChange={(e) => setNamaVaksin(e.target.value)}
              className={`w-full px-4 py-2.5 bg-neutral-50 border rounded-xl text-sm text-neutral-800 focus:outline-none focus:ring-2 focus:ring-primary-200 focus:border-primary transition-colors ${errors.namaVaksin ? 'border-red-300' : 'border-neutral-200'}`}
            >
              <option value="">Pilih jenis vaksin</option>
              {VAKSIN_OPTIONS.map((v) => <option key={v} value={v}>{v}</option>)}
            </select>
            {errors.namaVaksin && <p className="text-xs text-red-500 mt-1">{errors.namaVaksin}</p>}
          </div>

          {/* Tanggal Jadwal */}
          <div>
            <label className="text-xs font-semibold text-neutral-600 block mb-1.5">Tanggal Jadwal</label>
            <input
              type="date"
              value={tanggalJadwal}
              onChange={(e) => setTanggalJadwal(e.target.value)}
              className={`w-full px-4 py-2.5 bg-neutral-50 border rounded-xl text-sm text-neutral-800 focus:outline-none focus:ring-2 focus:ring-primary-200 focus:border-primary transition-colors ${errors.tanggalJadwal ? 'border-red-300' : 'border-neutral-200'}`}
            />
            {errors.tanggalJadwal && <p className="text-xs text-red-500 mt-1">{errors.tanggalJadwal}</p>}
          </div>

          {/* Status Pelaksanaan */}
          <div>
            <label className="text-xs font-semibold text-neutral-600 block mb-1.5">Status Pelaksanaan</label>
            <div className="grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={() => setStatus('BELUM')}
                className={`flex items-center justify-center gap-2 py-2.5 rounded-xl border-2 text-sm font-semibold transition-all ${status === 'BELUM' ? 'border-primary bg-primary-50 text-primary' : 'border-neutral-200 text-neutral-500 hover:border-neutral-300'}`}
              >
                <Clock size={15} /> Pending
              </button>
              <button
                type="button"
                onClick={() => setStatus('SUDAH')}
                className={`flex items-center justify-center gap-2 py-2.5 rounded-xl border-2 text-sm font-semibold transition-all ${status === 'SUDAH' ? 'border-primary bg-primary-50 text-primary' : 'border-neutral-200 text-neutral-500 hover:border-neutral-300'}`}
              >
                <CheckCircle size={15} /> Done
              </button>
            </div>
          </div>
        </form>

        {/* Footer */}
        <div className="flex items-center justify-end gap-3 px-6 py-4 border-t border-neutral-100">
          <button onClick={onClose} className="px-5 py-2.5 text-sm font-semibold text-neutral-600 hover:text-neutral-800 transition-colors">
            Batal
          </button>
          <button
            onClick={handleSubmit as unknown as React.MouseEventHandler}
            className="px-6 py-2.5 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-bold transition-colors"
          >
            Simpan Jadwal
          </button>
        </div>
      </div>
    </div>
  );
}
