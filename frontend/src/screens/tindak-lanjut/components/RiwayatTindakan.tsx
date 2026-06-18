
import { History, ChevronDown } from 'lucide-react';

export default function RiwayatTindakan(): JSX.Element {
    return (
        <div className="bg-[#f4f7fc] rounded-2xl p-6 border border-neutral-50">
            <div className="flex items-center justify-between mb-6">
                <div className="flex items-center gap-2">
                    <History className="w-5 h-5 text-[#3b82f6]" />
                    <h2 className="text-lg font-bold text-neutral-800 font-headline">Riwayat Tindakan Pasien</h2>
                </div>
                <button className="text-sm text-[#3b82f6] font-semibold hover:text-[#2563eb] transition-colors cursor-pointer">
                    Lihat Semua
                </button>
            </div>

            <div className="bg-white rounded-xl overflow-hidden border border-neutral-100">
                <table className="w-full text-left border-collapse">
                    <thead>
                        <tr className="bg-[#e0e7ff]/40 text-[11px] font-bold text-neutral-500 uppercase tracking-wider">
                            <th className="px-6 py-4 font-semibold">Tanggal</th>
                            <th className="px-6 py-4 font-semibold">Tindakan</th>
                            <th className="px-6 py-4 font-semibold">Faskes/Lokasi</th>
                            <th className="px-6 py-4 font-semibold">Status</th>
                            <th className="px-6 py-4 font-semibold text-right">Aksi</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-neutral-100">
                        <tr className="hover:bg-neutral-50/50 transition-colors">
                            <td className="px-6 py-4 text-sm font-medium text-neutral-800">
                                12 Apr<br/><span className="text-neutral-500 font-normal">2024</span>
                            </td>
                            <td className="px-6 py-4 text-sm text-neutral-700">Konseling<br/>Gizi</td>
                            <td className="px-6 py-4 text-sm text-neutral-700">Posyandu<br/>Melati</td>
                            <td className="px-6 py-4">
                                <span className="inline-flex items-center gap-1.5 bg-[#dcfce7] text-[#16a34a] text-[11px] font-bold px-2.5 py-1 rounded-full">
                                    <span className="w-1.5 h-1.5 bg-[#16a34a] rounded-full"></span> Membaik
                                </span>
                            </td>
                            <td className="px-6 py-4 text-right">
                                <button className="text-neutral-400 hover:text-[#3b82f6] transition-colors cursor-pointer">
                                    <ChevronDown className="w-5 h-5 ml-auto" />
                                </button>
                            </td>
                        </tr>
                        <tr className="hover:bg-neutral-50/50 transition-colors">
                            <td className="px-6 py-4 text-sm font-medium text-neutral-800">
                                05 Mar<br/><span className="text-neutral-500 font-normal">2024</span>
                            </td>
                            <td className="px-6 py-4 text-sm text-neutral-700">Rujukan<br/>Puskesmas</td>
                            <td className="px-6 py-4 text-sm text-neutral-700">Puskesmas<br/>Grogol</td>
                            <td className="px-6 py-4">
                                <span className="inline-flex items-center gap-1.5 bg-[#e0e7ff] text-[#3b82f6] text-[11px] font-bold px-2.5 py-1 rounded-full">
                                    <span className="w-1.5 h-1.5 bg-[#3b82f6] rounded-full"></span> Dalam Pemantauan
                                </span>
                            </td>
                            <td className="px-6 py-4 text-right">
                                <button className="text-neutral-400 hover:text-[#3b82f6] transition-colors cursor-pointer">
                                    <ChevronDown className="w-5 h-5 ml-auto" />
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    );
}
