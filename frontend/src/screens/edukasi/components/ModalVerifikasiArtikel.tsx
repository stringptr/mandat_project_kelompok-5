import { X, CheckCircle, XCircle, User, Clock } from 'lucide-react';
import type { Artikel } from '../data/artikel.data';
import { KategoriBadge } from './KategoriBadge';

interface ModalVerifikasiArtikelProps {
  artikel: Artikel;
  onClose: () => void;
  onVerifikasi: (id: string, action: 'approve' | 'reject') => void;
}

export function ModalVerifikasiArtikel({
  artikel,
  onClose,
  onVerifikasi,
}: ModalVerifikasiArtikelProps): JSX.Element {
  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-xl max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between px-6 py-4 border-b border-neutral-100 sticky top-0 bg-white z-10">
          <div>
            <h2 className="text-lg font-bold text-neutral-800 font-headline">Verifikasi Artikel</h2>
            <p className="text-xs text-neutral-500 mt-0.5 font-body">
              Tinjau dan setujui atau tolak artikel dari Bidan
            </p>
          </div>
          <button onClick={onClose} className="p-2 hover:bg-neutral-100 rounded-xl transition-colors">
            <X size={20} className="text-neutral-500" />
          </button>
        </div>

        <div className="p-6 space-y-4">
          <div className="bg-gradient-to-br from-amber-600 to-amber-800 rounded-xl p-4">
            <div className="flex items-start justify-between">
              <KategoriBadge kategori={artikel.kategori} variant="solid" />
              <span className="bg-amber-500 text-white px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-wide">
                Menunggu Verifikasi
              </span>
            </div>
          </div>

          <div>
            <h3 className="text-base font-bold text-neutral-800 font-headline mb-2">{artikel.judul}</h3>
            <div className="flex items-center gap-3 text-xs text-neutral-500 font-body mb-3">
              <span className="flex items-center gap-1">
                <User size={12} /> {artikel.penulis}
              </span>
              <span className="flex items-center gap-1">
                <Clock size={12} /> {artikel.waktuBaca}
              </span>
            </div>
            <div className="bg-neutral-50 rounded-xl p-4 border border-neutral-100">
              <p className="text-xs font-semibold text-neutral-500 uppercase tracking-wide mb-2">Ringkasan</p>
              <p className="text-sm text-neutral-700 font-body leading-relaxed">{artikel.ringkasan}</p>
            </div>
          </div>

          <div>
            <p className="text-xs font-semibold text-neutral-500 uppercase tracking-wide mb-2">Konten Artikel</p>
            <div className="bg-neutral-50 rounded-xl p-4 border border-neutral-100 max-h-48 overflow-y-auto">
              <p className="text-sm text-neutral-700 font-body leading-relaxed">{artikel.konten}</p>
            </div>
          </div>

          <div className="bg-amber-50 border border-amber-200 rounded-xl p-4">
            <p className="text-sm font-semibold text-amber-800 mb-1 font-body">Keputusan Verifikasi</p>
            <p className="text-xs text-amber-700 font-body mb-4 leading-relaxed">
              Artikel ini dikirim oleh <strong>{artikel.penulis}</strong>. Jika disetujui, artikel akan
              langsung dipublikasikan dan dapat dibaca oleh semua pengguna.
            </p>
            <div className="flex gap-3">
              <button
                onClick={() => onVerifikasi(artikel.id, 'reject')}
                className="flex-1 py-3 bg-white border border-red-200 text-red-600 hover:bg-red-50 rounded-xl text-sm font-semibold transition-colors flex items-center justify-center gap-2 font-body"
              >
                <XCircle size={16} />
                Tolak Artikel
              </button>
              <button
                onClick={() => onVerifikasi(artikel.id, 'approve')}
                className="flex-1 py-3 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-semibold transition-colors flex items-center justify-center gap-2 font-body"
              >
                <CheckCircle size={16} />
                Setujui & Publikasikan
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
