import { useState, useEffect } from 'react';
import { X, Clock, User, Calendar } from 'lucide-react';
import { apiGet } from '../../../lib/api';
import type { Artikel } from '../data/artikel.data';
import { KategoriBadge } from './KategoriBadge';

interface ModalDetailArtikelProps {
  artikel: Artikel;
  onClose: () => void;
}

export function ModalDetailArtikel({ artikel, onClose }: ModalDetailArtikelProps): JSX.Element {
  const [isi, setIsi] = useState<string>('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const idNum = parseInt(String(artikel.id).replace(/[^0-9]/g, ''), 10) || 0;
    if (idNum) {
      apiGet<{ isi_artikel: string }>(`/artikel/${idNum}`)
        .then((res) => setIsi(res.isi_artikel ?? ''))
        .catch(() => setIsi(artikel.ringkasan || ''))
        .finally(() => setLoading(false));
    } else {
      setIsi(artikel.ringkasan || '');
      setLoading(false);
    }
  }, [artikel.id, artikel.ringkasan]);

  const paragraphs = isi
    .split('. ')
    .reduce<string[][]>((acc, s, i) => {
      const idx = Math.floor(i / 3);
      if (!acc[idx]) acc[idx] = [];
      acc[idx].push(s);
      return acc;
    }, [])
    .map((group) => {
      const text = group.join('. ');
      return text.endsWith('.') ? text : text + '.';
    });

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <div className="relative bg-gradient-to-br from-primary-600 to-primary-800 rounded-t-2xl p-6">
          <div className="flex items-start justify-between">
            <KategoriBadge kategori={artikel.kategori} variant="solid" />
            <button onClick={onClose} className="p-2 bg-white/20 hover:bg-white/40 backdrop-blur-sm rounded-xl text-white transition-colors">
              <X size={18} />
            </button>
          </div>
        </div>

        <div className="p-6">
          <h2 className="text-xl font-bold text-neutral-900 font-headline leading-snug mb-3">{artikel.judul}</h2>

          <div className="flex items-center gap-4 text-xs text-neutral-500 font-body mb-4 pb-4 border-b border-neutral-100">
            <span className="flex items-center gap-1.5"><User size={13} className="text-primary" /> {artikel.penulis}</span>
            <span className="flex items-center gap-1.5"><Calendar size={13} className="text-primary" /> {artikel.tanggal}</span>
            <span className="flex items-center gap-1.5"><Clock size={13} className="text-primary" /> {artikel.waktuBaca}</span>
            {artikel.status === 'pending' && (
              <span className="ml-auto px-2.5 py-1 bg-amber-100 text-amber-700 rounded-full text-[10px] font-bold uppercase tracking-wide">
                Menunggu Verifikasi
              </span>
            )}
          </div>

          {loading ? (
            <div className="text-center py-12 text-neutral-400 text-sm">Memuat artikel...</div>
          ) : (
            <div className="text-sm text-neutral-700 font-body leading-relaxed space-y-3">
              {paragraphs.map((p, i) => <p key={i}>{p}</p>)}
            </div>
          )}

          <div className="mt-6 pt-4 border-t border-neutral-100">
            <button onClick={onClose} className="w-full py-3 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-semibold transition-colors font-body">
              Tutup
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
