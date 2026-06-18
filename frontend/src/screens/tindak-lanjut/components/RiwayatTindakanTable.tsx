import { Calendar, MapPin, ArrowRight, Eye, User, Filter } from 'lucide-react';

export interface RiwayatTindakanData {
    id: string;
    namaPasien: string;
    jenisTindakan: string;
    tanggal: string;
    faskes: string;
    statusRujukan: 'Diajukan' | 'Diproses' | 'Diterima' | 'Selesai';
}

interface RiwayatTindakanTableProps {
    data: RiwayatTindakanData[];
}

const getBadgeClasses = (status: string) => {
    switch (status) {
        case 'Diajukan': return 'bg-amber-50 text-amber-700 border-amber-200';
        case 'Diproses': return 'bg-sky-50 text-sky-700 border-sky-200';
        case 'Diterima': return 'bg-emerald-50 text-emerald-700 border-emerald-200';
        case 'Selesai': return 'bg-slate-50 text-slate-600 border-slate-200';
        default: return 'bg-gray-50 text-gray-600 border-gray-200';
    }
};

const getBadgeDot = (status: string) => {
    switch (status) {
        case 'Diajukan': return 'bg-amber-500';
        case 'Diproses': return 'bg-sky-500';
        case 'Diterima': return 'bg-emerald-500';
        case 'Selesai': return 'bg-slate-400';
        default: return 'bg-gray-400';
    }
};

export default function RiwayatTindakanTable({ data }: RiwayatTindakanTableProps): JSX.Element {
    return (
        <div className="bg-white rounded-2xl border border-slate-100 shadow-sm">
            <div className="p-6 border-b border-slate-100">
                <div className="flex items-center justify-between">
                    <div>
                        <h2 className="text-lg font-bold text-slate-800">Riwayat Tindakan Pasien</h2>
                        <p className="text-sm text-slate-500 mt-1">Catatan tindak lanjut dan rujukan yang telah dilakukan</p>
                    </div>
                    <button className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-emerald-700 bg-emerald-50 rounded-xl hover:bg-emerald-100 transition-colors border border-emerald-200">
                        <Filter size={16} />
                        Filter
                    </button>
                </div>
            </div>

            <div className="overflow-x-auto">
                <table className="w-full">
                    <thead>
                        <tr className="bg-slate-50/80 border-b border-slate-100">
                            <th className="text-left text-xs font-semibold text-slate-500 uppercase tracking-wider px-6 py-4">Nama Pasien</th>
                            <th className="text-left text-xs font-semibold text-slate-500 uppercase tracking-wider px-6 py-4">Jenis Tindakan</th>
                            <th className="text-left text-xs font-semibold text-slate-500 uppercase tracking-wider px-6 py-4">Tanggal</th>
                            <th className="text-left text-xs font-semibold text-slate-500 uppercase tracking-wider px-6 py-4">Fasilitas Kesehatan</th>
                            <th className="text-left text-xs font-semibold text-slate-500 uppercase tracking-wider px-6 py-4">Status Rujukan</th>
                            <th className="text-left text-xs font-semibold text-slate-500 uppercase tracking-wider px-6 py-4">Aksi</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                        {data.map((item) => (
                            <tr key={item.id} className="hover:bg-slate-50/50 transition-colors">
                                <td className="px-6 py-4">
                                    <div className="flex items-center gap-3">
                                        <div className="w-8 h-8 rounded-full bg-emerald-100 flex items-center justify-center flex-shrink-0">
                                            <User size={14} className="text-emerald-600" />
                                        </div>
                                        <span className="text-sm font-medium text-slate-800">{item.namaPasien}</span>
                                    </div>
                                </td>
                                <td className="px-6 py-4 text-sm text-slate-600">{item.jenisTindakan}</td>
                                <td className="px-6 py-4">
                                    <div className="flex items-center gap-1.5 text-sm text-slate-600">
                                        <Calendar size={14} className="text-slate-400" />
                                        {item.tanggal}
                                    </div>
                                </td>
                                <td className="px-6 py-4">
                                    <div className="flex items-center gap-1.5 text-sm text-slate-600">
                                        <MapPin size={14} className="text-slate-400" />
                                        {item.faskes}
                                    </div>
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`inline-flex items-center gap-1.5 text-xs font-semibold px-2.5 py-1 rounded-full border ${getBadgeClasses(item.statusRujukan)}`}>
                                        <span className={`w-1.5 h-1.5 rounded-full ${getBadgeDot(item.statusRujukan)}`}></span>
                                        {item.statusRujukan}
                                    </span>
                                </td>
                                <td className="px-6 py-4">
                                    <button className="flex items-center gap-1.5 text-xs font-medium text-emerald-600 hover:text-emerald-700 hover:bg-emerald-50 px-3 py-1.5 rounded-lg transition-colors">
                                        <Eye size={14} />
                                        Lihat Detail
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            <div className="p-4 border-t border-slate-100 text-center">
                <button className="inline-flex items-center gap-2 text-sm font-medium text-emerald-600 hover:text-emerald-700 hover:underline transition-colors">
                    Lihat Semua Riwayat Rujukan
                    <ArrowRight size={16} />
                </button>
            </div>
        </div>
    );
}