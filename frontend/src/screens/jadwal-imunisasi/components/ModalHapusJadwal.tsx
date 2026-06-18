import { AlertTriangle } from 'lucide-react';
import type { JadwalImunisasi } from '../data/imunisasi.data';

interface ModalHapusJadwalProps {
  jadwal: JadwalImunisasi;
  onClose: () => void;
  onHapus: () => void;
  loading?: boolean;
}

export function ModalHapusJadwal({ jadwal, onClose, onHapus, loading = false }: ModalHapusJadwalProps): JSX.Element {
  return (
    <div className="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-sm font-body p-8 text-center">
        {/* Icon */}
        <div className="w-14 h-14 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-5">
          <AlertTriangle size={28} className="text-red-500" />
        </div>

        <h2 className="text-xl font-bold text-neutral-900 font-headline mb-2">Hapus Jadwal Imunisasi?</h2>
        <p className="text-sm text-neutral-500 leading-relaxed mb-6">
          Apakah Anda yakin ingin menghapus jadwal imunisasi ini? Tindakan ini tidak dapat dibatalkan.
        </p>

        {/* Preview */}
        <div className="bg-neutral-50 border border-neutral-100 rounded-xl p-4 mb-6 text-left flex items-center gap-3">
          <div className="w-9 h-9 bg-primary/10 rounded-xl flex items-center justify-center flex-shrink-0">
            <span className="text-sm">📅</span>
          </div>
          <div>
            <p className="text-sm font-bold text-neutral-800">{jadwal.namaVaksin.split(' (')[0]} &amp; {jadwal.idPasien}</p>
            <p className="text-xs text-neutral-500 mt-0.5">Target: {jadwal.tanggalJadwal} · {jadwal.namaAnak}</p>
          </div>
        </div>

        <button
          onClick={onHapus}
          disabled={loading}
          className="w-full flex items-center justify-center gap-2 py-3 bg-red-500 hover:bg-red-600 disabled:bg-red-300 text-white rounded-xl text-sm font-bold transition-colors mb-3"
        >
          {loading ? (
            <>
              <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
              </svg>
              Menghapus...
            </>
          ) : (
            'Ya, Hapus'
          )}
        </button>
        <button
          onClick={onClose}
          disabled={loading}
          className="w-full py-3 bg-neutral-100 hover:bg-neutral-200 disabled:opacity-40 text-neutral-600 rounded-xl text-sm font-semibold transition-colors"
        >
          Batal
        </button>
      </div>
    </div>
  );
}
