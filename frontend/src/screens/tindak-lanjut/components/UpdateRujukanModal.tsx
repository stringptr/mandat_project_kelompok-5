// UpdateRujukanModal with dropdown for faskes from API
import { X, Loader2 } from 'lucide-react';
import { useState, useEffect } from 'react';
import { useNotification } from '../../../context/NotificationContext';
import { apiGet, apiPatch } from '../../../lib/api';
import type { Rujukan } from './RujukanAktif';
import type { FaskesItem } from '../../../types/entities';

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
  const [faskesList, setFaskesList] = useState<FaskesItem[]>([]);
  const [selectedFaskesId, setSelectedFaskesId] = useState<number | string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    apiGet<FaskesItem[]>('/faskes')
      .then((data) => {
        setFaskesList(data);
        const match = data.find((f) => f.nama_faskes === rujukan.faskes);
        if (match) setSelectedFaskesId(match.id_faskes);
      })
      .catch((err) => {
          console.error('Gagal memuat daftar faskes:', err);
          notify.error('Gagal memuat daftar fasilitas kesehatan.');
          setFaskesList([]);
      });
  }, [rujukan.faskes]);

  const handleSave = async () => {
    if (!selectedFaskesId) {
      notify.warn('Mohon pilih fasilitas kesehatan.');
      return;
    }
    setIsSubmitting(true);
    const selectedFaskes = faskesList.find((f) => f.id_faskes === Number(selectedFaskesId));
    const updated: Rujukan = { ...rujukan, status, faskes: selectedFaskes?.nama_faskes ?? rujukan.faskes };
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
          <div>
            <label className="block text-sm font-medium text-neutral-700 mb-1">
              Fasilitas Tujuan (Faskes)
            </label>
            <select
              value={selectedFaskesId}
              onChange={(e) => setSelectedFaskesId(e.target.value)}
              className="w-full bg-gray-100 border border-gray-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
            >
              <option value="">-- Pilih Faskes --</option>
              {faskesList.map((f) => (
                <option key={f.id_faskes} value={f.id_faskes}>
                  {f.nama_faskes} ({f.tipe_faskes})
                </option>
              ))}
            </select>
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
