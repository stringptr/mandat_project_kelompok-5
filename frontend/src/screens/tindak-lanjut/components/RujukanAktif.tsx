
import { User, Clock, CheckCircle, ArrowRight, ClipboardList, Calendar } from 'lucide-react';

export interface Rujukan {
    id: string;
    patientName: string;
    patientAge: string;
    urgency: 'Mendesak' | 'Berjalan' | 'Review';
    faskes: string;
    nutritionStatus: string;
    jenisTindakan: 'Rujukan' | 'Tindak Lanjut';
    status: string;
    lastWeight: number;
    lastHeight: number;
    timeLabel?: string;
}

interface RujukanAktifProps {
    data: Rujukan[];
    onDetailClick: (rujukan: Rujukan) => void;
    onUpdateClick: (rujukan: Rujukan) => void;
}

export default function RujukanAktif({ data, onDetailClick, onUpdateClick }: RujukanAktifProps): JSX.Element {

    const getUrgencyStyles = (urgency: string) => {
        switch (urgency) {
            case 'Mendesak':
                return 'bg-[#fee2e2] text-[#dc2626]';
            case 'Berjalan':
                return 'bg-[#dcfce7] text-[#16a34a]';
            case 'Review':
                return 'bg-[#e0e7ff] text-[#4f46e5]';
            default:
                return 'bg-neutral-100 text-neutral-600';
        }
    };

    const getUrgencyIcon = (urgency: string) => {
        switch (urgency) {
            case 'Mendesak':
                return <ClipboardList className="w-3.5 h-3.5" />;
            case 'Berjalan':
                return <Clock className="w-3.5 h-3.5" />;
            case 'Review':
                return <ClipboardList className="w-3.5 h-3.5" />;
            default:
                return <ClipboardList className="w-3.5 h-3.5" />;
        }
    };

    const getStatusIcon = (urgency: string) => {
        switch (urgency) {
            case 'Mendesak':
                return <Clock className="w-3.5 h-3.5" />;
            case 'Berjalan':
                return <Calendar className="w-3.5 h-3.5" />;
            case 'Review':
                return <CheckCircle className="w-3.5 h-3.5" />;
            default:
                return <Clock className="w-3.5 h-3.5" />;
        }
    };

    return (
        <div className="bg-[#f4f7fc] rounded-2xl p-6 border border-neutral-50 h-fit">
            <div className="flex items-center justify-between mb-6">
                <h2 className="text-lg font-bold text-neutral-800 font-headline">Rujukan Aktif</h2>
                <span className="bg-[#e0e7ff] text-blue-700 text-[11px] font-bold px-3 py-1 rounded-full tracking-wide uppercase">
                    {data.length} Aktif
                </span>
            </div>

            <div className="mb-6">

            </div>

            <div className="space-y-4">
                {data.map((item) => (
                    <div key={item.id} className="bg-white rounded-xl p-4 shadow-sm border border-neutral-100 hover:shadow-md transition-shadow">
                        <div className="flex justify-between items-start mb-3">
                            <div className="flex gap-3">
                                <div className="w-10 h-10 bg-neutral-100 rounded-full flex items-center justify-center shrink-0">
                                    <User className="w-5 h-5 text-neutral-400" />
                                </div>
                                <div>
                                    <h3 className="font-bold text-neutral-800 text-sm">{item.patientName}</h3>
                                    <p className="text-xs text-neutral-500 mt-0.5">{item.patientAge}</p>
                                </div>
                            </div>
                            <span className={`text-[10px] font-bold px-2 py-1 rounded uppercase ${getUrgencyStyles(item.urgency)}`}>
                                {item.urgency}
                            </span>
                        </div>
                        <div className="space-y-2 mb-4 mt-2">
                            <div className="flex items-center gap-2 text-xs text-neutral-600">
                                <div className="w-4 flex justify-center">{getUrgencyIcon(item.urgency)}</div>
                                <span>{item.faskes}</span>
                            </div>
                            <div className="flex items-center gap-2 text-xs text-neutral-600">
                                <div className="w-4 flex justify-center">{getStatusIcon(item.urgency)}</div>
                                <span>{item.status}</span>
                            </div>
                        </div>
                        <div className="flex gap-2">
                            <button
                                onClick={() => onDetailClick(item)}
                                className="flex-1 bg-[#f0f4ff] text-[#3b82f6] hover:bg-[#e0e7ff] text-xs font-semibold py-2 rounded-lg transition-colors cursor-pointer"
                            >
                                Detail
                            </button>
                            <button
                                onClick={() => onUpdateClick(item)}
                                className="flex-1 bg-[#e6f4ea] text-[#095c3e] hover:bg-[#d1fae5] text-xs font-semibold py-2 rounded-lg transition-colors cursor-pointer"
                            >
                                Update
                            </button>
                        </div>
                    </div>
                ))}
            </div>

            <button className="w-full mt-6 flex items-center justify-center gap-2 text-sm text-[#3b82f6] font-semibold hover:text-[#2563eb] transition-colors cursor-pointer">
                Lihat Semua Riwayat
                <ArrowRight className="w-4 h-4" />
            </button>
        </div>
    );
}
