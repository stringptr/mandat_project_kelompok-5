import { Baby, Ruler, Weight, Calendar, Heart, Utensils } from 'lucide-react';
import { StatCard } from '../components/statcard';
import { ChartWidget } from '../components/chartwidget';
import { DataTable, type Column } from '../components/datatable';
import { StatusBadge } from '../components/statusbadge';

// === Dummy Data ===
interface Pengukuran {
  [key: string]: unknown;
  tanggal: string;
  bb: string;
  tb: string;
  statusGizi: string;
  catatan: string;
}

const RIWAYAT_PENGUKURAN: Pengukuran[] = [
  { tanggal: '12 Mei 2024', bb: '13.2', tb: '82.0', statusGizi: 'Gizi Baik', catatan: 'Pertumbuhan normal' },
  { tanggal: '14 Apr 2024', bb: '12.8', tb: '80.5', statusGizi: 'Gizi Baik', catatan: 'Pertumbuhan normal' },
  { tanggal: '10 Mar 2024', bb: '12.3', tb: '79.1', statusGizi: 'Gizi Baik', catatan: 'BB naik 0.4 kg' },
  { tanggal: '12 Feb 2024', bb: '11.9', tb: '77.8', statusGizi: 'Gizi Kurang', catatan: 'Perlu pemantauan' },
  { tanggal: '15 Jan 2024', bb: '11.5', tb: '76.2', statusGizi: 'Gizi Kurang', catatan: 'BB di bawah standar' },
  { tanggal: '11 Des 2023', bb: '11.1', tb: '74.9', statusGizi: 'Gizi Kurang', catatan: 'Rujukan PMT' },
];

const GROWTH_CHART_DATA = [
  { label: 'Des', value: 11.1, color: '#ef4444' },
  { label: 'Jan', value: 11.5, color: '#f59e0b' },
  { label: 'Feb', value: 11.9, color: '#f59e0b' },
  { label: 'Mar', value: 12.3, color: '#10b981' },
  { label: 'Apr', value: 12.8, color: '#10b981' },
  { label: 'Mei', value: 13.2, color: '#095c3e' },
];

const TABLE_COLUMNS: Column<Pengukuran>[] = [
  { header: 'TANGGAL', accessor: 'tanggal', className: 'font-medium text-neutral-500' },
  { header: 'BB (KG)', accessor: 'bb', className: 'font-semibold text-neutral-700' },
  { header: 'TB (CM)', accessor: 'tb', className: 'font-semibold text-neutral-700' },
  {
    header: 'STATUS GIZI',
    accessor: 'statusGizi',
    render: (row) => {
      const variant = row.statusGizi === 'Gizi Baik' ? 'gizi-baik' : 'gizi-kurang';
      return <StatusBadge variant={variant} label={row.statusGizi} />;
    },
  },
  { header: 'CATATAN', accessor: 'catatan', className: 'text-neutral-500' },
];

export function IbuWaliSection() {
  return (
    <div className="space-y-6">
      {/* Greeting */}
      <div className="bg-gradient-to-r from-primary-50 to-emerald-50 rounded-2xl p-6 border border-primary-100">
        <div className="flex items-center gap-4">
          <div className="w-14 h-14 bg-primary/10 rounded-2xl flex items-center justify-center">
            <Baby size={28} className="text-primary" />
          </div>
          <div>
            <h2 className="text-lg font-bold text-neutral-800 font-headline">
              Data Anak: Arka Mahendra
            </h2>
            <p className="text-sm text-neutral-500 mt-0.5">
              Laki-laki • Usia 24 Bulan • Tanggal Lahir: 12 Mei 2022
            </p>
          </div>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Berat Badan"
          value="13.2"
          subtitle="kg"
          icon={<Weight size={22} />}
          variant="icon"
          color="primary"
          trend={{ direction: 'up', value: '+0.4 kg', label: 'dari bulan lalu' }}
        />
        <StatCard
          title="Tinggi Badan"
          value="82.0"
          subtitle="cm"
          icon={<Ruler size={22} />}
          variant="icon"
          color="blue"
          trend={{ direction: 'up', value: '+1.5 cm', label: 'dari bulan lalu' }}
        />
        <StatCard
          title="Status Gizi"
          value="Baik"
          icon={<Heart size={22} />}
          variant="icon"
          color="emerald"
          subtitle="Z-Score: -0.8 SD"
        />
        <StatCard
          title="Jadwal Posyandu"
          value="14 Jun"
          icon={<Calendar size={22} />}
          variant="icon"
          color="purple"
          subtitle="Posyandu Melati, 09:00"
        />
      </div>

      {/* Charts row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <ChartWidget
          title="Grafik Berat Badan"
          subtitle="6 bulan terakhir (kg)"
          type="bar-vertical"
          data={GROWTH_CHART_DATA}
        />

        {/* Nutrisi card */}
        <div className="bg-white rounded-2xl border border-neutral-100 shadow-sm p-6 space-y-4">
          <h3 className="text-base font-bold text-neutral-800 font-headline">
            Rekomendasi Nutrisi
          </h3>
          <div className="space-y-3">
            {[
              { icon: <Utensils size={16} />, title: 'Protein Hewani', desc: 'Tambah konsumsi telur & ikan 2x sehari', color: 'text-amber-600 bg-amber-50' },
              { icon: <Heart size={16} />, title: 'Sayur & Buah', desc: 'Minimal 3 porsi sayur/buah per hari', color: 'text-emerald-600 bg-emerald-50' },
              { icon: <Baby size={16} />, title: 'Susu/ASI', desc: 'Lanjutkan konsumsi susu 2 gelas/hari', color: 'text-blue-600 bg-blue-50' },
            ].map((rec, i) => (
              <div key={i} className="flex items-start gap-3 p-3 rounded-xl bg-neutral-50/50">
                <div className={`w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0 ${rec.color}`}>
                  {rec.icon}
                </div>
                <div>
                  <span className="text-sm font-semibold text-neutral-700">{rec.title}</span>
                  <p className="text-xs text-neutral-500 mt-0.5">{rec.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Table */}
      <DataTable
        title="Riwayat Pengukuran"
        columns={TABLE_COLUMNS}
        data={RIWAYAT_PENGUKURAN}
        pageSize={5}
        onExport={() => { }}
      />
    </div>
  );
}
