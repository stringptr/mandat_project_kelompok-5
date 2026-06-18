import { useState } from 'react';
import { X, Calendar } from 'lucide-react';
import { useNotification } from '../../../context/NotificationContext';
import { apiPut, apiPatch } from '../../../lib/api';
import type { JadwalImunisasi } from '../data/imunisasi.data';
import { VAKSIN_OPTIONS } from '../data/imunisasi.data';

interface ModalUpdateJadwalProps {
  jadwal: JadwalImunisasi;
  onClose: () => void;
  onSimpan: (updated: JadwalImunisasi) => void;
}

export function ModalUpdateJadwal({ jadwal, onClose, onSimpan }: ModalUpdateJadwalProps): JSX.Element {
  const notify = useNotification();
  const [namaVaksin, setNamaVaksin] = useState(jadwal.namaVaksin);
  const [tanggalRealisasi, setTanggalRealisasi] = useState(
    jadwal.tanggalRealisasi
      ? new Date().toISOString().split('T')[0]
      : new Date().toISOString().split('T')[0]
  );
  const [catatan, setCatatan] = useState(jadwal.catatan ?? '');

  const handleSimpan = async () => {
    if (!namaVaksin) {
      notify.warn('Mohon lengkapi semua data form yang wajib diisi sebelum mengirim.');
      return;
    }
    const formatted = new Date(tanggalRealisasi).toLocaleDateString('id-ID', {
      day: '2-digit', month: 'short', year: 'numeric',
    });
    const idNum = parseInt(String(jadwal.id).replace(/[^0-9]/g, ''), 10) || 0;
    try {
      if (namaVaksin !== jadwal.namaVaksin) {
        await apiPut('/imunisasi/' + idNum, { nama_vaksin: namaVaksin });
      }
      await apiPatch('/imunisasi/' + idNum + '/realisasi', {
        tanggal_realisasi: tanggalRealisasi,
      });
      onSimpan({
        ...jadwal,
        namaVaksin,
        status: 'SUDAH',
        tanggalRealisasi: formatted,
        catatan: catatan || undefined,
      });
    } catch {
      notify.error('Gagal memperbarui jadwal. Silakan coba lagi.');
    }
  };

  return (
    <div className="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md font-body">
        {/* Header */}
        <div className="flex items-center gap-3 px-6 py-5 border-b border-neutral-100">
          <div className="w-10 h-10 bg-primary/10 rounded-xl flex items-center justify-center flex-shrink-0">
            <Calendar size={20} className="text-primary" />
          </div>
          <div className="flex-1">
            <h2 className="text-base font-bold text-neutral-900 font-headline">Update Data Imunisasi</h2>
            <p className="text-xs text-neutral-500 mt-0.5">Perbarui catatan riwayat kesehatan anak</p>
          </div>
          <button onClick={onClose} className="p-1.5 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg transition-colors">
            <X size={18} />
          </button>
        </div>

        <div className="px-6 py-5 space-y-4">
          {/* ID Pasien + Nama Vaksin */}
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="text-[10px] font-bold text-neutral-400 uppercase tracking-wide block mb-1.5">ID Pasien</label>
              <div className="flex items-center gap-2 px-3 py-2.5 bg-neutral-50 border border-neutral-200 rounded-xl">
                <span className="text-neutral-400 text-sm">🔒</span>
                <span className="text-sm font-semibold text-neutral-700">{jadwal.idPasien}</span>
              </div>
              <p className="text-[10px] text-neutral-400 mt-1">Identitas pasien terenkripsi secara otomatis</p>
            </div>
            <div>
              <label className="text-[10px] font-bold text-neutral-400 uppercase tracking-wide block mb-1.5">Nama Vaksin</label>
              <select
                value={namaVaksin}
                onChange={(e) => setNamaVaksin(e.target.value)}
                className="w-full px-3 py-2.5 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-800 focus:outline-none focus:ring-2 focus:ring-primary-200 focus:border-primary"
              >
                {VAKSIN_OPTIONS.map((v) => <option key={v} value={v}>{v.split(' (')[0]}</option>)}
              </select>
            </div>
          </div>

          {/* Progress monitoring */}
          <div className="bg-neutral-50 rounded-xl border border-neutral-100 p-4">
            <p className="text-[10px] font-bold text-primary uppercase tracking-widest mb-3 flex items-center gap-1.5">
              <Calendar size={11} /> Progress Monitoring
            </p>
            <div className="grid grid-cols-2 gap-4 mb-3">
              <div>
                <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-wide mb-1">Status Realisasi</p>
                <div className="flex items-center gap-2">
                  <span className="text-sm font-bold text-neutral-800">{jadwal.status}</span>
                  <span className={`w-5 h-5 rounded-full flex items-center justify-center ${jadwal.status === 'SUDAH' ? 'text-emerald-500' : 'text-amber-500'}`}>
                    {jadwal.status === 'SUDAH' ? '✓' : '○'}
                  </span>
                </div>
              </div>
              <div>
                <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-wide mb-1">Tanggal Jadwal</p>
                <p className="text-sm font-bold text-neutral-800">
                  {new Date().toLocaleDateString('en-US', { month: '2-digit', day: '2-digit', year: 'numeric' }).replace(/\//g, '/')}
                </p>
              </div>
            </div>
            <div>
              <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-wide mb-1.5">Tanggal Realisasi</p>
              <div className="flex items-center gap-2">
                <input
                  type="date"
                  value={tanggalRealisasi}
                  onChange={(e) => setTanggalRealisasi(e.target.value)}
                  className="flex-1 px-3 py-2 bg-white border border-neutral-200 rounded-xl text-sm text-neutral-800 focus:outline-none focus:ring-2 focus:ring-primary-200"
                />
                <Calendar size={18} className="text-primary flex-shrink-0" />
              </div>
              <p className="text-[10px] text-neutral-400 mt-1">Catatan: Sesuaikan dengan tanggal kunjungan aktual di Posyandu</p>
            </div>
          </div>

          {/* Catatan tambahan */}
          <div>
            <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest mb-1.5">Catatan Kesehatan Tambahan (Opsional)</p>
            <textarea
              value={catatan}
              onChange={(e) => setCatatan(e.target.value)}
              placeholder="Contoh: Kondisi bayi sedikit demam pasca imunisasi, disarankan pemberian perasetamol..."
              rows={3}
              className="w-full px-4 py-3 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-800 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 resize-none"
            />
          </div>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-end gap-3 px-6 py-4 border-t border-neutral-100">
          <button onClick={onClose} className="px-5 py-2.5 text-sm font-semibold text-neutral-600 hover:text-neutral-800 transition-colors">
            Batal
          </button>
          <button
            onClick={handleSimpan}
            className="px-6 py-2.5 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-bold transition-colors"
          >
            Simpan Perubahan
          </button>
        </div>
      </div>
    </div>
  );
}
