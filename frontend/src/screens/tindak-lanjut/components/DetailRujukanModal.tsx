import { X, MapPin, Clock, FileText, User, Calendar, AlertTriangle, CheckCircle2 } from 'lucide-react';
import type { Rujukan } from './RujukanAktif';

interface DetailRujukanModalProps {
    isOpen: boolean;
    onClose: () => void;
    rujukan: Rujukan | null;
}

export default function DetailRujukanModal({ isOpen, onClose, rujukan }: DetailRujukanModalProps): JSX.Element | null {
    if (!isOpen || !rujukan) return null;

    const getDeadlineIcon = () => {
        switch (rujukan.statusDeadline) {
            case 'terlambat': return <AlertTriangle className="w-4 h-4 text-red-500" />;
            case 'mendekati': return <Clock className="w-4 h-4 text-amber-500" />;
            default: return <CheckCircle2 className="w-4 h-4 text-green-500" />;
        }
    };

    const getDeadlineClass = () => {
        switch (rujukan.statusDeadline) {
            case 'terlambat': return 'text-red-700 bg-red-50 border-red-200';
            case 'mendekati': return 'text-amber-700 bg-amber-50 border-amber-200';
            default: return 'text-green-700 bg-green-50 border-green-200';
        }
    };

    const formatDate = (dateStr: string) => {
        if (!dateStr) return '-';
        return new Date(dateStr).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' });
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4" onClick={onClose}>
            <div className="bg-white rounded-2xl w-full max-w-md shadow-xl overflow-hidden" onClick={(e) => e.stopPropagation()}>
                {/* Header */}
                <div className="px-6 py-4 border-b border-neutral-100 flex items-center justify-between bg-neutral-50/50">
                    <h3 className="font-bold text-lg text-neutral-800">
                        {rujukan.jenisTindakan === 'Tindak Lanjut' ? 'Detail Tindak Lanjut' : 'Detail Rujukan'}
                    </h3>
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
                            <span className="inline-block mt-1 text-[10px] font-bold px-2 py-0.5 rounded uppercase bg-neutral-100 text-neutral-600">
                                {rujukan.jenisTindakan}
                            </span>
                        </div>
                    </div>

                    <div className="space-y-4">
                        {/* Deadline Status */}
                        {rujukan.statusDeadline && (
                            <div className={`rounded-xl p-3 border ${getDeadlineClass()}`}>
                                <div className="flex items-center gap-2">
                                    {getDeadlineIcon()}
                                    <div>
                                        <p className="text-xs font-semibold">Deadline: {formatDate(rujukan.tanggalDeadline)}</p>
                                        <p className="text-[10px] opacity-75">
                                            {rujukan.statusDeadline === 'terlambat' ? 'Rujukan ini sudah melewati batas waktu' :
                                             rujukan.statusDeadline === 'mendekati' ? 'Rujukan ini mendekati batas waktu' :
                                             'Rujukan masih dalam batas waktu'}
                                        </p>
                                    </div>
                                </div>
                            </div>
                        )}

                        <div className="bg-[#f4f7fc] rounded-xl p-4">
                            <h5 className="text-xs font-bold text-neutral-500 uppercase tracking-wider mb-3">
                                {rujukan.jenisTindakan === 'Tindak Lanjut' ? 'Informasi Tindak Lanjut' : 'Informasi Rujukan'}
                            </h5>

                            <div className="space-y-3">
                                {rujukan.jenisTindakan === 'Rujukan' && (
                                    <div className="flex items-start gap-3">
                                        <MapPin className="w-4 h-4 text-[#3b82f6] shrink-0 mt-0.5" />
                                        <div>
                                            <p className="text-xs font-semibold text-neutral-500 mb-0.5">Faskes Tujuan</p>
                                            <p className="text-sm font-medium text-neutral-800">{rujukan.faskes || '-'}</p>
                                        </div>
                                    </div>
                                )}
                                <div className="flex items-start gap-3">
                                    <Calendar className="w-4 h-4 text-[#10b981] shrink-0 mt-0.5" />
                                    <div>
                                        <p className="text-xs font-semibold text-neutral-500 mb-0.5">Tanggal Diajukan</p>
                                        <p className="text-sm font-medium text-neutral-800">{formatDate(rujukan.tanggalRujukan)}</p>
                                    </div>
                                </div>
                                <div className="flex items-start gap-3">
                                    <Clock className="w-4 h-4 text-[#eab308] shrink-0 mt-0.5" />
                                    <div>
                                        <p className="text-xs font-semibold text-neutral-500 mb-0.5">Status Saat Ini</p>
                                        <p className="text-sm font-medium text-neutral-800">{rujukan.status || '-'}</p>
                                    </div>
                                </div>
                                {rujukan.alasanRujukan && (
                                    <div className="flex items-start gap-3">
                                        <FileText className="w-4 h-4 text-[#8b5cf6] shrink-0 mt-0.5" />
                                        <div>
                                            <p className="text-xs font-semibold text-neutral-500 mb-0.5">
                                                {rujukan.jenisTindakan === 'Tindak Lanjut' ? 'Catatan Medis' : 'Alasan Rujukan'}
                                            </p>
                                            <p className="text-sm text-neutral-700 whitespace-pre-wrap">{rujukan.alasanRujukan}</p>
                                        </div>
                                    </div>
                                )}
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
