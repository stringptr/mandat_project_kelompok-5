import { X, MapPin, Clock, FileText, User } from 'lucide-react';
import type { Rujukan } from './RujukanAktif';

interface DetailRujukanModalProps {
    isOpen: boolean;
    onClose: () => void;
    rujukan: Rujukan | null;
}

export default function DetailRujukanModal({ isOpen, onClose, rujukan }: DetailRujukanModalProps): JSX.Element | null {
    if (!isOpen || !rujukan) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4">
            <div className="bg-white rounded-2xl w-full max-w-md shadow-xl overflow-hidden" onClick={(e) => e.stopPropagation()}>
                {/* Header */}
                <div className="px-6 py-4 border-b border-neutral-100 flex items-center justify-between bg-neutral-50/50">
                    <h3 className="font-bold text-lg text-neutral-800 font-headline">Detail Rujukan</h3>
                    <button 
                        onClick={onClose}
                        className="p-2 text-neutral-400 hover:text-neutral-700 hover:bg-neutral-100 rounded-full transition-colors cursor-pointer"
                    >
                        <X className="w-5 h-5" />
                    </button>
                </div>

                {/* Content */}
                <div className="p-6">
                    <div className="flex items-center gap-4 mb-6">
                        <div className="w-14 h-14 bg-[#f4f7fc] rounded-full flex items-center justify-center shrink-0 border border-neutral-100">
                            <User className="w-7 h-7 text-neutral-400" />
                        </div>
                        <div>
                            <h4 className="font-bold text-neutral-800 text-lg">{rujukan.patientName}</h4>
                            <p className="text-sm text-neutral-500">{rujukan.patientAge}</p>
                            <span className="inline-block mt-1 text-[10px] font-bold px-2 py-0.5 rounded uppercase bg-neutral-100 text-neutral-600">
                                ID: #RJK-{rujukan.id.substring(0,6).toUpperCase()}
                            </span>
                        </div>
                    </div>

                    <div className="space-y-4">
                        <div className="bg-[#f4f7fc] rounded-xl p-4">
                            <h5 className="text-xs font-bold text-neutral-500 uppercase tracking-wider mb-3">Informasi Rujukan</h5>
                            
                            <div className="space-y-3">
                                <div className="flex items-start gap-3">
                                    <MapPin className="w-4 h-4 text-[#3b82f6] shrink-0 mt-0.5" />
                                    <div>
                                        <p className="text-xs font-semibold text-neutral-500 mb-0.5">Faskes / Tujuan</p>
                                        <p className="text-sm font-medium text-neutral-800">{rujukan.faskes}</p>
                                    </div>
                                </div>
                                <div className="flex items-start gap-3">
                                    <Clock className="w-4 h-4 text-[#eab308] shrink-0 mt-0.5" />
                                    <div>
                                        <p className="text-xs font-semibold text-neutral-500 mb-0.5">Status Saat Ini</p>
                                        <p className="text-sm font-medium text-neutral-800">{rujukan.status}</p>
                                    </div>
                                </div>
                                <div className="flex items-start gap-3">
                                    <FileText className="w-4 h-4 text-[#10b981] shrink-0 mt-0.5" />
                                    <div>
                                        <p className="text-xs font-semibold text-neutral-500 mb-0.5">Tingkat Urgensi</p>
                                        <p className="text-sm font-medium text-neutral-800">{rujukan.urgency}</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Footer */}
                <div className="px-6 py-4 border-t border-neutral-100 bg-neutral-50 flex justify-end">
                    <button 
                        onClick={onClose}
                        className="px-6 py-2 bg-white border border-neutral-200 hover:bg-neutral-50 text-neutral-700 text-sm font-semibold rounded-xl transition-colors cursor-pointer"
                    >
                        Tutup
                    </button>
                </div>
            </div>
        </div>
    );
}
