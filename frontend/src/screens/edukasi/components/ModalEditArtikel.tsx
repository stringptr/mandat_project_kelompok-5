import { useState, useEffect } from 'react';
import { X } from 'lucide-react';
import { useNotification } from '../../../context/NotificationContext';
import { apiGet, apiPatch } from '../../../lib/api';
import type { Artikel, KategoriArtikel } from '../data/artikel.data';
import { KATEGORI_LIST } from '../data/artikel.data';

const KATEGORI_OPTIONS = KATEGORI_LIST.filter((k) => k !== 'Semua') as KategoriArtikel[];

interface ModalEditArtikelProps {
  artikel: Artikel;
  onClose: () => void;
  onSubmit: () => void;
}

export function ModalEditArtikel({ artikel, onClose, onSubmit }: ModalEditArtikelProps): JSX.Element {
  const notify = useNotification();
  const [form, setForm] = useState({ judul: artikel.judul, konten: '', kategori: artikel.kategori as KategoriArtikel });
  const [loading, setLoading] = useState(true);
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    const idNum = parseInt(String(artikel.id).replace(/[^0-9]/g, ''), 10) || 0;
    if (idNum) {
      apiGet<{ isi_artikel: string }>(`/artikel/${idNum}`)
        .then((res) => setForm((f) => ({ ...f, konten: res.isi_artikel ?? '' })))
        .catch(() => {})
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, [artikel.id]);

  const validate = () => {
    const e: Record<string, string> = {};
    if (!form.judul.trim()) e.judul = 'Judul wajib diisi';
    if (!form.konten.trim()) e.konten = 'Konten wajib diisi';
    return e;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const errs = validate();
    if (Object.keys(errs).length > 0) { setErrors(errs); notify.warn('Mohon lengkapi data wajib.'); return; }
    try {
      const idNum = parseInt(String(artikel.id).replace(/[^0-9]/g, ''), 10) || 0;
      await apiPatch('/artikel/' + idNum, { judul: form.judul, isi_artikel: form.konten, kategori: form.kategori });
      notify.success('Artikel berhasil diperbarui!');
      onSubmit();
    } catch {
      notify.error('Gagal memperbarui artikel.');
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between px-6 py-4 border-b border-neutral-100 sticky top-0 bg-white z-10">
          <div><h2 className="text-lg font-bold text-neutral-800 font-headline">Edit Artikel</h2></div>
          <button onClick={onClose} className="p-2 hover:bg-neutral-100 rounded-xl transition-colors"><X size={20} className="text-neutral-500" /></button>
        </div>

        {loading ? (
          <div className="text-center py-12 text-neutral-400 text-sm">Memuat artikel...</div>
        ) : (
          <form onSubmit={handleSubmit} className="p-6 space-y-5">
            <div>
              <label className="text-sm font-semibold text-neutral-700 mb-1.5 block">Judul <span className="text-red-500">*</span></label>
              <input type="text" value={form.judul} onChange={(e) => setForm({ ...form, judul: e.target.value })}
                className={`w-full px-4 py-3 border rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-200 ${errors.judul ? 'border-red-300 bg-red-50' : 'border-neutral-200'}`} />
              {errors.judul && <p className="text-xs text-red-500 mt-1">{errors.judul}</p>}
            </div>

            <div>
              <label className="text-sm font-semibold text-neutral-700 mb-1.5 block">Kategori</label>
              <select value={form.kategori} onChange={(e) => setForm({ ...form, kategori: e.target.value as KategoriArtikel })}
                className="w-full px-4 py-3 border border-neutral-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-200 bg-white">
                {KATEGORI_OPTIONS.map((k) => <option key={k} value={k}>{k}</option>)}
              </select>
            </div>

            <div>
              <label className="text-sm font-semibold text-neutral-700 mb-1.5 block">Konten <span className="text-red-500">*</span></label>
              <textarea value={form.konten} onChange={(e) => setForm({ ...form, konten: e.target.value })}
                rows={10} className={`w-full px-4 py-3 border rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-200 resize-none ${errors.konten ? 'border-red-300 bg-red-50' : 'border-neutral-200'}`} />
              {errors.konten && <p className="text-xs text-red-500 mt-1">{errors.konten}</p>}
            </div>

            <div className="flex gap-3 pt-2">
              <button type="button" onClick={onClose} className="flex-1 py-3 border border-neutral-200 text-neutral-600 rounded-xl text-sm font-semibold hover:bg-neutral-50">Batal</button>
              <button type="submit" className="flex-1 py-3 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-semibold">Simpan Perubahan</button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
