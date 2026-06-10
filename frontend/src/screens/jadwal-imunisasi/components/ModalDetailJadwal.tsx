import { X, Calendar, User, Clock } from 'lucide-react';
import type { JadwalImunisasi } from '../data/imunisasi.data';

interface ModalDetailJadwalProps {
  jadwal: JadwalImunisasi;
  onClose: () => void;
}

export function ModalDetailJadwal({ jadwal, onClose }: ModalDetailJadwalProps): JSX.Element {
  return (
    <div className="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-sm font-body">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-neutral-100">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 bg-primary/10 rounded-xl flex items-center justify-center">
              <Calendar size={18} className="text-primary" />
            </div>
            <div>
              <h2 className="text-sm font-bold text-neutral-900 font-headline">Detail Jadwal Imunisasi</h2>
              <p className="text-[10px] text-neutral-400 mt-0.5">Read-only · Dinas Kesehatan</p>
            </div>
          </div>
          <button onClick={onClose} className="p-1.5 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg transition-colors">
            <X size={18} />
          </button>
        </div>

        {/* Content */}
        <div className="px-6 py-5 space-y-4">
          {/* Pasien */}
          <div className="flex items-center gap-3 bg-neutral-50 rounded-xl p-3.5">
            <div className="w-10 h-10 bg-primary-50 rounded-full flex items-center justify-center flex-shrink-0">
              <User size={18} className="text-primary" />
            </div>
            <div>
              <p className="text-sm font-bold text-neutral-800">{jadwal.namaAnak}</p>
              <p className="text-xs text-neutral-500 font-body">{jadwal.idPasien}</p>
            </div>
          </div>

          {/* Vaksin info */}
          <div className="space-y-3">
            <Row label="Nama Vaksin" value={jadwal.namaVaksin} />
            <Row label="Dosis" value={jadwal.dosis} />
            <Row label="Tanggal Jadwal" value={jadwal.tanggalJadwal} icon={<Clock size={13} className="text-neutral-400" />} />
            <Row
              label="Tanggal Realisasi"
              value={jadwal.tanggalRealisasi ?? '—'}
              icon={<Calendar size={13} className="text-neutral-400" />}
            />
          </div>

          {/* Status */}
          <div className="flex items-center justify-between bg-neutral-50 rounded-xl px-4 py-3">
            <span className="text-xs font-semibold text-neutral-500 uppercase tracking-wide">Status</span>
            <span className={`text-xs font-bold px-3 py-1 rounded-full ${jadwal.status === 'SUDAH' ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-600'}`}>
              {jadwal.status}
            </span>
          </div>

          {/* Catatan */}
          {jadwal.catatan && (
            <div>
              <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-wide mb-1.5">Catatan</p>
              <p className="text-xs text-neutral-600 bg-neutral-50 rounded-xl p-3 leading-relaxed">{jadwal.catatan}</p>
            </div>
          )}
        </div>

        <div className="px-6 py-4 border-t border-neutral-100">
          <button onClick={onClose} className="w-full py-2.5 bg-neutral-100 hover:bg-neutral-200 text-neutral-700 rounded-xl text-sm font-semibold transition-colors">
            Tutup
          </button>
        </div>
      </div>
    </div>
  );
}

function Row({ label, value, icon }: { label: string; value: string; icon?: React.ReactNode }): JSX.Element {
  return (
    <div className="flex items-center justify-between text-sm">
      <span className="text-xs font-semibold text-neutral-500 flex items-center gap-1.5">{icon}{label}</span>
      <span className="font-semibold text-neutral-800 text-right max-w-[55%]">{value}</span>
    </div>
  );
}
