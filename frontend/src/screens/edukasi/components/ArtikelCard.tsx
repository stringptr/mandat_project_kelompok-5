import { Pencil, Trash2, Clock, User, CheckCircle, AlertCircle } from 'lucide-react';
import type { Artikel } from '../data/artikel.data';
import { KategoriBadge } from './KategoriBadge';

interface ArtikelCardProps {
  artikel: Artikel;
  onRead: (a: Artikel) => void;
  onEdit?: (a: Artikel) => void;
  onHapus?: (a: Artikel) => void;
  onVerifikasi?: (a: Artikel) => void;
  canEdit?: boolean;
  showVerifikasi?: boolean;
}

export function ArtikelCard({
  artikel,
  onRead,
  onEdit,
  onHapus,
  onVerifikasi,
  canEdit = false,
  showVerifikasi = false,
}: ArtikelCardProps): JSX.Element {
  const isPending = artikel.status === 'pending';
  const isRejected = artikel.status === 'rejected';

  return (
    <div
      onClick={() => onRead(artikel)}
      className={`bg-white rounded-2xl overflow-hidden border transition-shadow group cursor-pointer hover:shadow-md
        ${isPending ? 'border-amber-200' : isRejected ? 'border-red-200' : 'border-neutral-100'}
      `}
    >
      {/* Status & action header */}
      <div className="px-4 pt-4 pb-0">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            {isPending && (
              <span className="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-amber-500 text-white">
                Menunggu Verifikasi
              </span>
            )}
            {isRejected && (
              <span className="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-red-500 text-white">
                Ditolak
              </span>
            )}
          </div>
          <div className="flex items-center gap-1">
            {showVerifikasi && isPending && onVerifikasi && (
              <button
                onClick={(e) => { e.stopPropagation(); onVerifikasi(artikel); }}
                title="Verifikasi Artikel"
                className="p-1.5 bg-amber-500 hover:bg-amber-600 rounded-lg text-white transition-colors"
              >
                <CheckCircle size={14} />
              </button>
            )}
            {canEdit && onEdit && (
              <button
                onClick={(e) => { e.stopPropagation(); onEdit(artikel); }}
                title="Edit Artikel"
                className="p-1.5 text-neutral-400 hover:text-primary hover:bg-primary-50 rounded-lg transition-colors"
              >
                <Pencil size={14} />
              </button>
            )}
            {canEdit && onHapus && (
              <button
                onClick={(e) => { e.stopPropagation(); onHapus(artikel); }}
                title="Hapus Artikel"
                className="p-1.5 text-neutral-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
              >
                <Trash2 size={14} />
              </button>
            )}
          </div>
        </div>
      </div>

      {/* Body */}
      <div className="p-4">
        <div className="flex items-center justify-between mb-2">
          <KategoriBadge kategori={artikel.kategori} variant="outline" />
          <span className="text-xs text-neutral-400 font-body">{artikel.tanggal}</span>
        </div>

        <h4 className="font-bold text-neutral-800 text-sm font-headline leading-snug mb-1.5 line-clamp-2">
          {artikel.judul}
        </h4>
        <p className="text-xs text-neutral-500 font-body line-clamp-3 mb-4 leading-relaxed">
          {artikel.ringkasan}
        </p>

        <div className="flex items-center gap-3 text-neutral-400 text-xs border-t border-neutral-100 pt-3">
          <span className="flex items-center gap-1 truncate">
            <User size={11} className="flex-shrink-0" />
            <span className="truncate">{artikel.penulis}</span>
          </span>
          <span className="flex items-center gap-1 flex-shrink-0">
            <Clock size={11} />
            {artikel.waktuBaca}
          </span>
        </div>

        {/* Verifikasi CTA for Dinkes inside card */}
        {showVerifikasi && isPending && onVerifikasi && (
          <button
            onClick={(e) => { e.stopPropagation(); onVerifikasi(artikel); }}
            className="mt-3 w-full flex items-center justify-center gap-2 bg-amber-50 border border-amber-200 text-amber-700 rounded-xl px-3 py-2 text-xs font-semibold hover:bg-amber-100 transition-colors"
          >
            <AlertCircle size={13} />
            Klik untuk Verifikasi
          </button>
        )}
      </div>
    </div>
  );
}
