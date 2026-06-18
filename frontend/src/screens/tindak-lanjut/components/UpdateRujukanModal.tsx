// UpdateRujukanModal with autocomplete for faskes
import { X, Loader2 } from 'lucide-react';
import { useState } from 'react';
import { useNotification } from '../../../context/NotificationContext';
import { apiPatch } from '../../../lib/api';
import type { Rujukan } from './RujukanAktif';

// Static list of facilities (could be fetched later)
const FACILITIES = [
  'RSUD Kota - Spesialis Anak',
  'Puskesmas Melati',
  'Posyandu Melati',
];

interface UpdateRujukanModalProps {
  isOpen: boolean;
  onClose: () => void;
  rujukan: Rujukan;
  onSave: (updated: Rujukan) => void;
}

export default function UpdateRujukanModal({
  isOpen,
  onClose,
  rujukan,
  onSave,
}: UpdateRujukanModalProps): JSX.Element | null {
  if (!isOpen) return null;

  const notify = useNotification();
  const [status, setStatus] = useState(rujukan.status);
  const [faskes, setFaskes] = useState(rujukan.faskes);
  const [query, setQuery] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Filter suggestions based on query (case‑insensitive)
  const suggestions = query
    ? FACILITIES.filter((f) => f.toLowerCase().startsWith(query.toLowerCase()))
    : [];

  const handleSave = async () => {
    if (!faskes.trim()) {
      notify.warn('Mohon lengkapi semua data form yang wajib diisi sebelum mengirim.');
      return;
    }
    setIsSubmitting(true);
    const updated: Rujukan = { ...rujukan, status, faskes };
    try {
      const idNum = parseInt(String(rujukan.id).replace(/[^0-9]/g, ''), 10) || 0;
      await apiPatch('/rujukan/' + idNum + '/status', { status_rujukan: status });
      onSave(updated);
      setIsSubmitting(false);
      onClose();
    } catch {
      setIsSubmitting(false);
      notify.error('Gagal memperbarui rujukan. Silakan coba lagi.');
    }
  };

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4"
      onClick={(e) => e.stopPropagation()}
    >
      <div
        className="bg-white rounded-2xl w-full max-w-md shadow-xl overflow-hidden"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="px-6 py-4 border-b border-neutral-100 flex items-center justify-between bg-neutral-50/50">
          <h3 className="font-bold text-lg text-neutral-800">Update Rujukan</h3>
          <button
            onClick={onClose}
            className="p-2 text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100 rounded-full transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>
        {/* Content */}
        <div className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-neutral-700 mb-1">Status</label>
            <select
              value={status}
              onChange={(e) => setStatus(e.target.value)}
              className="w-full bg-gray-100 border border-gray-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
            >
              <option value="Diajukan">Diajukan</option>
              <option value="Diproses">Diproses</option>
              <option value="Diterima">Diterima</option>
              <option value="Ditolak">Ditolak</option>
            </select>
          </div>
          <div className="relative">
            <label className="block text-sm font-medium text-neutral-700 mb-1">
              Fasilitas Tujuan (Faskes)
            </label>
            <input
              type="text"
              value={faskes}
              onChange={(e) => {
                setFaskes(e.target.value);
                setQuery(e.target.value);
              }}
              className="w-full bg-gray-100 border border-gray-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
            />
            {/* Autocomplete dropdown */}
            {suggestions.length > 0 && (
              <ul className="absolute left-0 right-0 z-10 mt-1 bg-white border border-gray-200 rounded-md shadow-lg max-h-48 overflow-y-auto">
                {suggestions.map((s) => (
                  <li
                    key={s}
                    onClick={() => {
                      setFaskes(s);
                      setQuery('');
                    }}
                    className="px-3 py-2 text-sm cursor-pointer hover:bg-gray-100"
                  >
                    {s}
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
        {/* Footer */}
        <div className="px-6 py-4 border-t border-neutral-100 bg-neutral-50 flex justify-end space-x-2">
          <button
            onClick={onClose}
            disabled={isSubmitting}
            className="px-4 py-2 bg-white border border-neutral-200 hover:bg-neutral-50 text-neutral-700 rounded-md transition-colors disabled:opacity-50"
          >
            Batal
          </button>
          <button
            onClick={handleSave}
            disabled={isSubmitting}
            className="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-md flex items-center space-x-2 disabled:opacity-50"
          >
            {isSubmitting && <Loader2 className="w-4 h-4 animate-spin" />}
            <span>Simpan</span>
          </button>
        </div>
      </div>
    </div>
  );
}
