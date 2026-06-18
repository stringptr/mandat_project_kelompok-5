/**
 * ModalVerifikasiBidan — triggered from Notifikasi or any "Lakukan Verifikasi" CTA.
 * Full-screen modal that wraps VerifikasiPanel with patient context.
 */
import { X } from 'lucide-react';
import { VerifikasiPanel, type VerifikasiTarget } from './VerifikasiPanel';

interface ModalVerifikasiBidanProps {
  target: VerifikasiTarget;
  onClose: () => void;
  onSetuju: (id: string, catatan: string) => void;
  onTolak: (id: string, catatan: string) => void;
}

export function ModalVerifikasiBidan({
  target,
  onClose,
  onSetuju,
  onTolak,
}: ModalVerifikasiBidanProps): JSX.Element {
  const handleSetuju = (id: string, catatan: string) => {
    onSetuju(id, catatan);
    onClose();
  };
  const handleTolak = (id: string, catatan: string) => {
    onTolak(id, catatan);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md font-body overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-neutral-100">
          <div>
            <h2 className="text-base font-bold text-neutral-900 font-headline">Lakukan Verifikasi</h2>
            <p className="text-xs text-neutral-500 mt-0.5">Data pengukuran dari Kader Posyandu</p>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg transition-colors"
          >
            <X size={18} />
          </button>
        </div>

        {/* Panel */}
        <div className="p-6">
          <VerifikasiPanel
            target={target}
            onSetuju={handleSetuju}
            onTolak={handleTolak}
          />
        </div>
      </div>
    </div>
  );
}
