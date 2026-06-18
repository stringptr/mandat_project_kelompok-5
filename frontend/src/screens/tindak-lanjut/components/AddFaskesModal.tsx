import { X, Loader2 } from 'lucide-react';
import { useState } from 'react';

interface AddFaskesModalProps {
  isOpen: boolean;
  onClose: () => void;
  onAdd: (newFaskes: string) => void;
}

export default function AddFaskesModal({ isOpen, onClose, onAdd }: AddFaskesModalProps): JSX.Element | null {
  if (!isOpen) return null;

  const [value, setValue] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSave = () => {
    if (!value.trim()) return;
    setIsSubmitting(true);
    // Simulate async save
    setTimeout(() => {
      onAdd(value.trim());
      setIsSubmitting(false);
      onClose();
      setValue('');
    }, 500);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4" onClick={e => e.stopPropagation()}>
      <div className="bg-white rounded-2xl w-full max-w-md shadow-xl overflow-hidden" onClick={e => e.stopPropagation()}>
        {/* Header */}
        <div className="px-6 py-4 border-b border-neutral-100 flex items-center justify-between bg-neutral-50/50">
          <h3 className="font-bold text-lg text-neutral-800">Tambah Faskes Baru</h3>
          <button onClick={onClose} className="p-2 text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100 rounded-full transition-colors">
            <X className="w-5 h-5" />
          </button>
        </div>
        {/* Content */}
        <div className="p-6 space-y-4">
          <input
            type="text"
            placeholder="Nama fasilitas faskes..."
            value={value}
            onChange={e => setValue(e.target.value)}
            className="w-full bg-gray-100 border border-gray-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
          />
        </div>
        {/* Footer */}
        <div className="px-6 py-4 border-t border-neutral-100 bg-neutral-50 flex justify-end space-x-2">
          <button onClick={onClose} disabled={isSubmitting} className="px-4 py-2 bg-white border border-neutral-200 hover:bg-neutral-50 text-neutral-700 rounded-md disabled:opacity-50">
            Batal
          </button>
          <button onClick={handleSave} disabled={isSubmitting} className="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-md flex items-center space-x-2 disabled:opacity-50">
            {isSubmitting && <Loader2 className="w-4 h-4 animate-spin" />}
            <span>Tambah</span>
          </button>
        </div>
      </div>
    </div>
  );
}
