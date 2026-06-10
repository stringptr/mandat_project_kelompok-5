import { X, Trash2, AlertTriangle } from 'lucide-react';
import type { Artikel } from '../data/artikel.data';

interface ModalHapusArtikelProps {
  artikel: Artikel;
  onClose: () => void;
  onHapus: () => void;
}

export function ModalHapusArtikel({ artikel, onClose, onHapus }: ModalHapusArtikelProps): JSX.Element {
  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md">
        <div className="flex items-center justify-between px-6 py-4 border-b border-neutral-100">
          <h2 className="text-lg font-bold text-neutral-800 font-headline">Hapus Artikel</h2>
          <button onClick={onClose} className="p-2 hover:bg-neutral-100 rounded-xl transition-colors">
            <X size={20} className="text-neutral-500" />
          </button>
        </div>

        <div className="p-6">
          <div className="flex items-start gap-4 mb-6">
            <div className="w-12 h-12 bg-red-100 rounded-xl flex items-center justify-center flex-shrink-0">
              <AlertTriangle size={24} className="text-red-500" />
            </div>
            <div>
              <p className="text-sm font-semibold text-neutral-800 mb-1 font-body">
                Apakah Anda yakin ingin menghapus artikel ini?
              </p>
              <p className="text-sm text-neutral-600 font-body leading-relaxed">
                Artikel{' '}
                <span className="font-semibold text-neutral-800">"{artikel.judul}"</span>{' '}
                akan dihapus secara permanen dan tidak dapat dipulihkan.
              </p>
            </div>
          </div>

          <div className="bg-neutral-50 rounded-xl p-3 mb-6 flex items-center gap-3">
            <img
              src={artikel.gambar}
              alt=""
              className="w-14 h-14 rounded-lg object-cover flex-shrink-0"
            />
            <div className="min-w-0">
              <p className="text-sm font-semibold text-neutral-700 truncate font-body">{artikel.judul}</p>
              <p className="text-xs text-neutral-500 font-body">
                {artikel.kategori} · {artikel.tanggal}
              </p>
            </div>
          </div>

          <div className="flex gap-3">
            <button
              onClick={onClose}
              className="flex-1 py-3 border border-neutral-200 text-neutral-600 rounded-xl text-sm font-semibold hover:bg-neutral-50 transition-colors font-body"
            >
              Batal
            </button>
            <button
              onClick={onHapus}
              className="flex-1 py-3 bg-red-500 hover:bg-red-600 text-white rounded-xl text-sm font-semibold transition-colors flex items-center justify-center gap-2 font-body"
            >
              <Trash2 size={16} />
              Hapus Artikel
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
