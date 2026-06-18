/**
 * VerifikasiPanel — shared right-side panel for Bidan to verify patient data.
 * Shows patient info, measurement snapshot, catatan field, and Tolak/Setuju buttons.
 * Used in BidanSection (Monitoring) and anywhere else Bidan needs to verify.
 */
import { useState } from 'react';
import { CheckCircle, X, ClipboardCheck, User } from 'lucide-react';

export interface VerifikasiTarget {
  id: string;
  nama: string;
  inisial: string;        // e.g. "RA"
  warnaBg: string;        // tailwind bg, e.g. "bg-primary"
  warnaText: string;      // tailwind text, e.g. "text-white"
  usia: string;
  bb: string;
  tb: string;
  petugas: string;        // kader / posyandu yang mengukur
  statusGizi: string;
}

interface VerifikasiPanelProps {
  target: VerifikasiTarget;
  onSetuju: (id: string, catatan: string) => void;
  onTolak: (id: string, catatan: string) => void;
}

export function VerifikasiPanel({ target, onSetuju, onTolak }: VerifikasiPanelProps): JSX.Element {
  const [catatan, setCatatan] = useState('');

  return (
    <div className="bg-primary-50 border border-primary-100 rounded-2xl p-5 h-full flex flex-col font-body">
      {/* Header */}
      <div className="flex items-center gap-2 mb-4">
        <div className="w-9 h-9 bg-primary rounded-xl flex items-center justify-center flex-shrink-0">
          <ClipboardCheck size={18} className="text-white" />
        </div>
        <div>
          <h3 className="text-sm font-bold text-neutral-800 font-headline">Verifikasi Data</h3>
          <p className="text-[10px] text-neutral-500 mt-0.5">Tinjau & konfirmasi pengukuran</p>
        </div>
      </div>

      {/* Patient row */}
      <div className="flex items-center gap-3 bg-white rounded-xl p-3 mb-4 border border-primary-100">
        <div className={`w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0 text-sm font-bold ${target.warnaBg} ${target.warnaText}`}>
          {target.inisial}
        </div>
        <div className="min-w-0">
          <p className="text-sm font-bold text-neutral-800 truncate">{target.nama}</p>
          <p className="text-xs text-neutral-500">{target.usia} · BB {target.bb} kg · TB {target.tb} cm</p>
        </div>
      </div>

      {/* Petugas info */}
      <div className="flex items-center gap-2 text-xs text-neutral-600 mb-4 bg-white rounded-xl px-3 py-2.5 border border-neutral-100">
        <User size={13} className="text-neutral-400 flex-shrink-0" />
        <span>
          Petugas Pengukur (Kpder):{' '}
          <span className="font-semibold text-neutral-700">{target.petugas}</span>
        </span>
      </div>

      {/* Status gizi badge */}
      <div className="flex items-center gap-2 mb-4">
        <span className={`text-xs font-bold px-2.5 py-1 rounded-full ${
          target.statusGizi === 'Stunting'
            ? 'bg-red-100 text-red-700'
            : target.statusGizi === 'Gizi Kurang'
            ? 'bg-amber-100 text-amber-700'
            : 'bg-emerald-100 text-emerald-700'
        }`}>
          {target.statusGizi}
        </span>
        <span className="text-[10px] text-neutral-400">Status gizi terdeteksi</span>
      </div>

      {/* Catatan verifikator */}
      <div className="flex-1 flex flex-col mb-4">
        <label className="text-[10px] font-bold text-neutral-500 uppercase tracking-wide mb-1.5">
          Catatan Verifikator
        </label>
        <textarea
          value={catatan}
          onChange={(e) => setCatatan(e.target.value)}
          placeholder="Berikan catatan medis atau verifikasi..."
          className="flex-1 w-full px-3 py-2.5 bg-white border border-neutral-200 rounded-xl text-sm text-neutral-800 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 resize-none min-h-24"
        />
      </div>

      {/* Action buttons */}
      <div className="grid grid-cols-2 gap-2 mb-3">
        <button
          onClick={() => onTolak(target.id, catatan)}
          className="flex items-center justify-center gap-1.5 py-2.5 bg-red-50 hover:bg-red-100 text-red-600 border border-red-200 rounded-xl text-sm font-semibold transition-colors"
        >
          <X size={15} /> Tolak
        </button>
        <button
          onClick={() => onSetuju(target.id, catatan)}
          className="flex items-center justify-center gap-1.5 py-2.5 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-semibold transition-colors"
        >
          <CheckCircle size={15} /> Setuju
        </button>
      </div>

      <p className="text-[10px] text-neutral-400 text-center leading-relaxed">
        Data yang diverifikasi akan tercatat dalam Laporan Bulanan Puskesmas Terpadu.
      </p>
    </div>
  );
}
