import { Users, AlertTriangle, Target, Download, Calendar, ChevronRight } from 'lucide-react';
import { INTERVENSI_TERBARU, KABUPATEN_DATA, TREN_NUTRISI } from '../data/dashboard.data';
import { PrevalensiMap } from '../components/PrevalensiMap';
import { TrendChart } from '../components/TrendChart';

export function DinkesSection(): JSX.Element {
  return (
    <div className="space-y-5">
      {/* Sub-header */}
      <div className="flex items-start justify-between flex-wrap gap-3">
        <div>
          <h3 className="text-xl font-bold text-neutral-900 font-headline">Statistik Kesehatan Wilayah</h3>
          <p className="text-xs text-neutral-500 font-body mt-1">
            Monitoring gizi maternal dan anak periode Oktober 2023 – Februari 2024. Data diperbarui secara real-time dari posyandu.
          </p>
        </div>
        <div className="flex items-center gap-2 flex-shrink-0">
          <button className="flex items-center gap-1.5 bg-primary-50 border border-primary-100 text-primary text-xs font-semibold px-3 py-2 rounded-xl hover:bg-primary-100 transition-colors font-body">
            <Calendar size={13} /> Februari 2024
          </button>
        </div>
      </div>

      {/* Top stat cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        {/* Total pasien */}
        <div className="bg-primary rounded-2xl p-5 relative overflow-hidden">
          <div className="absolute -bottom-6 -right-6 w-28 h-28 rounded-full bg-white/10" />
          <Users size={22} className="text-white/80 mb-3" />
          <p className="text-white/70 text-[10px] font-semibold uppercase tracking-wide font-body">Total Pasien Terdaftar</p>
          <p className="text-4xl font-bold text-white font-headline leading-none mt-1">12,482</p>
          <p className="text-white/70 text-xs font-body mt-2 flex items-center gap-1">
            <span>↑</span> +13% dari bulan lalu
          </p>
        </div>

        {/* Kasus gizi buruk */}
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 relative">
          <span className="absolute top-4 right-4 text-[10px] font-bold bg-red-500 text-white px-2 py-0.5 rounded-full">MASALAH</span>
          <AlertTriangle size={22} className="text-red-500 mb-3" />
          <p className="text-neutral-500 text-[10px] font-semibold uppercase tracking-wide font-body">Kasus Gizi Buruk (Stunting)</p>
          <p className="text-4xl font-bold text-red-500 font-headline leading-none mt-1">342</p>
          <p className="text-red-400 text-xs font-body mt-2 flex items-center gap-1">
            <span>↑</span> +2K Peningkatan di-verifikasi
          </p>
        </div>

        {/* Capaian pembangunan */}
        <div className="bg-blue-600 rounded-2xl p-5 relative overflow-hidden">
          <div className="absolute -bottom-6 -right-6 w-28 h-28 rounded-full bg-white/10" />
          <Target size={22} className="text-white/80 mb-3" />
          <p className="text-white/70 text-[10px] font-semibold uppercase tracking-wide font-body">Capaian Pembangunan</p>
          <p className="text-4xl font-bold text-white font-headline leading-none mt-1">94.2%</p>
          <p className="text-white/70 text-xs font-body mt-2 flex items-center gap-1">
            <span>🎯</span> Target keluaran: 95%
          </p>
        </div>
      </div>

      {/* Map + Tren chart */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* Map */}
        <div className="lg:col-span-2 bg-white rounded-2xl p-5 border border-neutral-100">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h4 className="text-sm font-bold text-neutral-800 font-headline">Sebaran Prevalensi Gizi Buruk</h4>
              <p className="text-xs text-neutral-500 font-body mt-0.5">Data spasial per wilayah kerja Puskesmas</p>
            </div>
            <div className="flex items-center gap-2">
              <button className="text-xs text-primary border border-primary rounded-lg px-3 py-1.5 font-semibold font-body hover:bg-primary-50 transition-colors">
                Peta Wilayah
              </button>
              <button className="text-xs text-neutral-500 border border-neutral-200 rounded-lg px-3 py-1.5 font-semibold font-body hover:bg-neutral-50 transition-colors">
                Tabel Ranking
              </button>
            </div>
          </div>
          <PrevalensiMap data={KABUPATEN_DATA} />
        </div>

        {/* Tren */}
        <div className="bg-blue-50 rounded-2xl p-5 border border-blue-100 flex flex-col">
          <h4 className="text-sm font-bold text-neutral-800 font-headline mb-1">Tren Nutrisi Bulanan</h4>
          <p className="text-xs text-neutral-500 font-body mb-4">Indeks perbaikan gizi regional</p>
          <div className="flex-1">
            <TrendChart
              data={TREN_NUTRISI}
              height={130}
              color="#1d4ed8"
              fillColor="rgba(29,78,216,0.08)"
            />
          </div>
          <div className="mt-4 pt-4 border-t border-blue-100">
            <div className="flex items-center justify-between mb-1">
              <span className="text-xs text-neutral-600 font-body">Tingkat Kesembuhan</span>
              <span className="text-sm font-bold text-blue-700 font-headline">86.4%</span>
            </div>
            <p className="text-[10px] text-neutral-500 font-body leading-relaxed">
              Program intervensi selama 4 bulan terakhir menunjukkan hasil gizi anak semakin membaik.
            </p>
          </div>
        </div>
      </div>

      {/* Intervensi terbaru table */}
      <div className="bg-white rounded-2xl border border-neutral-100 overflow-hidden">
        <div className="flex items-center justify-between px-5 py-4 border-b border-neutral-50">
          <h4 className="text-sm font-bold text-neutral-800 font-headline">Catatan Intervensi Terbaru</h4>
          <button className="text-xs text-primary font-semibold hover:text-primary-600 font-body flex items-center gap-1">
            Lihat Semua History <ChevronRight size={12} />
          </button>
        </div>

        {/* Table header */}
        <div className="grid grid-cols-5 px-5 py-2.5 bg-neutral-50 border-b border-neutral-100">
          {['Posyandu / Wilayah', 'Tindakan', 'Status Pasien', 'Waktu', 'Progres'].map((h) => (
            <p key={h} className="text-[10px] font-bold text-neutral-400 uppercase tracking-wide font-body">{h}</p>
          ))}
        </div>

        {/* Rows */}
        {INTERVENSI_TERBARU.map((item) => (
          <div
            key={item.id}
            className="grid grid-cols-5 px-5 py-3.5 border-b border-neutral-50 last:border-0 hover:bg-neutral-50/60 transition-colors items-center"
          >
            <div>
              <p className="text-sm font-semibold text-neutral-800 font-body">{item.posyandu}</p>
              <p className="text-xs text-neutral-400 font-body">{item.kecamatan}</p>
            </div>
            <div>
              <span className={`text-[10px] font-bold px-2 py-1 rounded-full ${item.tindakanColor}`}>
                {item.tindakan}
              </span>
            </div>
            <div className="flex items-center gap-1.5">
              <span className={`w-2 h-2 rounded-full flex-shrink-0 ${
                item.statusPasien === 'Normal' ? 'bg-emerald-500' :
                item.statusPasien === 'Underweight' ? 'bg-orange-500' : 'bg-red-500'
              }`} />
              <span className="text-xs font-medium text-neutral-700 font-body">{item.statusPasien}</span>
            </div>
            <p className="text-xs text-neutral-500 font-body">{item.waktu}</p>
            <p className={`text-sm font-bold font-headline ${item.progresPositif ? 'text-emerald-600' : 'text-red-500'}`}>
              {item.progres}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
