import { Calendar, AlertCircle } from 'lucide-react';
import { useState } from 'react';

// ── Types ────────────────────────────────────────
interface RujukanPasien {
  id: string;
  nama: string;
  avatarUrl?: string;
  statusRujukan: string; // e.g., 'Diajukan', 'Diproses', 'Diterima'
  statusGizi?: string;
  statusPasisn?: string;
  tanggalPeriksa?: string;
  umur?: string;
  jenisTindakan?: 'Rujukan' | 'Tindak Lanjut';
}

interface RujukanAktifCardsProps {
  data: RujukanPasien[];
}

// ── Helper Functions ───────────────────────────────
const getBadgeClasses = (status: string) => {
  switch (status) {
    case 'Diajukan':
      return 'bg-amber-50 text-amber-700 border-amber-200';
    case 'Diproses':
      return 'bg-sky-50 text-sky-700 border-sky-200';
    case 'Diterima':
      return 'bg-emerald-50 text-emerald-700 border-emerald-200';
    default:
      return 'bg-gray-50 text-gray-600 border-gray-200';
  }
};

const getBadgeDot = (status: string) => {
  switch (status) {
    case 'Diajukan':
      return 'bg-amber-500';
    case 'Diproses':
      return 'bg-sky-500';
    case 'Diterima':
      return 'bg-emerald-500';
    default:
      return 'bg-gray-400';
  }
};

// ── Component ───────────────────────────────────────
export default function RujukanAktifCards({ data }: RujukanAktifCardsProps): JSX.Element {
  const [filter, setFilter] = useState<'All' | 'Rujukan' | 'Tindak Lanjut'>('All');

  const filteredData = data.filter(item => {
    if (filter === 'All') return true;
    return item.jenisTindakan === filter;
  });

  return (
    <div className="bg-white rounded-2xl border border-slate-100 shadow-sm mb-8">
      {/* Header with filter buttons */}
      <div className="p-6 border-b border-slate-100 flex items-center justify-between">
        <div>
          <h2 className="text-lg font-bold text-slate-800">Rujukan Aktif</h2>
          <p className="text-sm text-slate-500 mt-1">Daftar pasien yang sedang menjalani rujukan atau tindak lanjut</p>
        </div>
        <div className="flex gap-2">
          <button
            className={`px-3 py-1 text-xs rounded ${filter === 'All' ? 'bg-emerald-600 text-white' : 'bg-gray-200 text-gray-800'}`}
            onClick={() => setFilter('All')}
          >
            All
          </button>
          <button
            className={`px-3 py-1 text-xs rounded ${filter === 'Rujukan' ? 'bg-emerald-600 text-white' : 'bg-gray-200 text-gray-800'}`}
            onClick={() => setFilter('Rujukan')}
          >
            Rujukan
          </button>
          <button
            className={`px-3 py-1 text-xs rounded ${filter === 'Tindak Lanjut' ? 'bg-emerald-600 text-white' : 'bg-gray-200 text-gray-800'}`}
            onClick={() => setFilter('Tindak Lanjut')}
          >
            Tindak Lanjut
          </button>
        </div>
      </div>

      <div className="p-6">
        <span className="bg-emerald-50 text-emerald-700 text-xs font-semibold px-3 py-1.5 rounded-full border border-emerald-200">
          {filteredData.length} Aktif
        </span>
        {filteredData.length === 0 ? (
          <div className="text-center py-12">
            <div className="w-16 h-16 bg-slate-50 rounded-full flex items-center justify-center mx-auto mb-4">
              <AlertCircle size={28} className="text-slate-400" />
            </div>
            <h3 className="text-base font-semibold text-slate-700 mb-1">Tidak ada pasien yang memerlukan tindak lanjut</h3>
            <p className="text-sm text-slate-500">Semua pasien telah ditangani dengan baik</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
            {filteredData.map(item => (
              <div key={item.id} className="flex items-start gap-4 p-4 rounded-xl border border-slate-100 hover:border-emerald-200 hover:shadow-sm transition-all bg-slate-50/30">
                <div className="w-12 h-12 rounded-full overflow-hidden flex-shrink-0 bg-slate-200 border-2 border-white shadow-sm">
                  {item.avatarUrl ? (
                    <img src={item.avatarUrl} alt={item.nama} className="w-full h-full object-cover" />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center bg-emerald-100 text-emerald-600 font-bold text-sm">
                      {item.nama.charAt(0)}
                    </div>
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between mb-1">
                    <h4 className="text-sm font-bold text-slate-800 truncate">{item.nama}</h4>
                    <span className={`inline-flex items-center gap-1.5 text-[11px] font-semibold px-2.5 py-1 rounded-full border ${getBadgeClasses(item.statusRujukan)}`}
                    >
                      <span className={`w-1.5 h-1.5 rounded-full ${getBadgeDot(item.statusRujukan)}`} />
                      {item.statusRujukan}
                    </span>
                  </div>
                  {item.umur && (
                    <p className="text-xs text-slate-500 mb-2">{item.umur}</p>
                  )}
                  <div className="flex items-center gap-3 mb-3">
                    {item.statusGizi && (
                      <span className="text-xs text-red-600 bg-red-50 px-2 py-0.5 rounded-md font-medium">{item.statusGizi}</span>
                    )}
                    {item.statusPasisn && (
                      <span className="text-xs text-slate-500">{item.statusPasisn}</span>
                    )}
                  </div>
                  {item.tanggalPeriksa && (
                    <div className="flex items-center gap-1.5 text-xs text-slate-400 mb-3">
                      <Calendar size={12} />
                      <span>Pemeriksaan terakhir: {item.tanggalPeriksa}</span>
                    </div>
                  )}
                  {item.jenisTindakan && (
                    <div className="flex items-center gap-1.5 text-xs text-slate-400 mb-3">
                      <span className="text-slate-500">{item.jenisTindakan}</span>
                    </div>
                  )}
                  <div className="flex items-center gap-2">
                    <button className="flex-1 px-3 py-2 text-xs font-medium text-slate-600 bg-white border border-slate-200 rounded-lg hover:bg-slate-50 hover:border-slate-300 transition-colors">
                      Detail
                    </button>
                    <button className="flex-1 px-3 py-2 text-xs font-medium text-white bg-emerald-600 rounded-lg hover:bg-emerald-700 transition-colors shadow-sm">
                      Update
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}