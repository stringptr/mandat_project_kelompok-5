import { useState } from 'react';
import { Building2, BarChart3, Landmark, Download, Target } from 'lucide-react';
import { StatCard } from '../components/statcard';
import { ChartWidget } from '../components/chartwidget';
import { DataTable, type Column } from '../components/datatable';
import { ProgressBar } from '../components/progressbar';
import { FilterBar } from '../components/filterbar';

// === Dummy Data ===
interface KabupatenData {
  [key: string]: unknown;
  no: number;
  kabupaten: string;
  totalBalita: number;
  stunting: number;
  prevalensi: string;
  cakupan: string;
  status: string;
}

const KABUPATEN_DATA: KabupatenData[] = [
  { no: 1, kabupaten: 'Kota Semarang', totalBalita: 12500, stunting: 875, prevalensi: '7.0%', cakupan: '94%', status: 'Baik' },
  { no: 2, kabupaten: 'Kab. Demak', totalBalita: 8200, stunting: 984, prevalensi: '12.0%', cakupan: '78%', status: 'Perlu Perhatian' },
  { no: 3, kabupaten: 'Kab. Kendal', totalBalita: 6800, stunting: 748, prevalensi: '11.0%', cakupan: '82%', status: 'Perlu Perhatian' },
  { no: 4, kabupaten: 'Kab. Semarang', totalBalita: 7100, stunting: 497, prevalensi: '7.0%', cakupan: '91%', status: 'Baik' },
  { no: 5, kabupaten: 'Kab. Grobogan', totalBalita: 9400, stunting: 1410, prevalensi: '15.0%', cakupan: '68%', status: 'Kritis' },
  { no: 6, kabupaten: 'Kab. Blora', totalBalita: 5600, stunting: 784, prevalensi: '14.0%', cakupan: '72%', status: 'Kritis' },
  { no: 7, kabupaten: 'Kab. Pati', totalBalita: 7800, stunting: 624, prevalensi: '8.0%', cakupan: '88%', status: 'Baik' },
  { no: 8, kabupaten: 'Kab. Jepara', totalBalita: 6300, stunting: 693, prevalensi: '11.0%', cakupan: '80%', status: 'Perlu Perhatian' },
];

const TREN_STUNTING = [
  { label: 'Des', value: 14.2, color: '#ef4444' },
  { label: 'Jan', value: 13.8, color: '#ef4444' },
  { label: 'Feb', value: 12.5, color: '#f59e0b' },
  { label: 'Mar', value: 11.9, color: '#f59e0b' },
  { label: 'Apr', value: 11.2, color: '#10b981' },
  { label: 'Mei', value: 10.8, color: '#10b981' },
];

const DISTRIBUSI_REGIONAL = [
  { label: 'Baik (<10%)', value: 12, color: '#10b981' },
  { label: 'Perhatian (10-14%)', value: 15, color: '#f59e0b' },
  { label: 'Kritis (≥15%)', value: 8, color: '#ef4444' },
];

const TABLE_COLUMNS: Column<KabupatenData>[] = [
  { header: 'NO', accessor: 'no', className: 'font-medium text-neutral-400 w-12' },
  {
    header: 'KABUPATEN/KOTA',
    accessor: 'kabupaten',
    render: (row) => <span className="font-bold text-primary">{row.kabupaten}</span>,
  },
  {
    header: 'TOTAL BALITA', accessor: 'totalBalita', className: 'font-semibold text-neutral-600',
    render: (row) => <span>{(row.totalBalita as number).toLocaleString('id-ID')}</span>,
  },
  {
    header: 'STUNTING', accessor: 'stunting', className: 'font-semibold text-red-600',
    render: (row) => <span className="text-red-600 font-bold">{(row.stunting as number).toLocaleString('id-ID')}</span>,
  },
  { header: 'PREVALENSI', accessor: 'prevalensi', className: 'font-bold' },
  { header: 'CAKUPAN', accessor: 'cakupan', className: 'font-semibold text-neutral-600' },
  {
    header: 'STATUS',
    accessor: 'status',
    render: (row) => {
      const colorMap: Record<string, string> = {
        Baik: 'bg-emerald-50 text-emerald-700',
        'Perlu Perhatian': 'bg-amber-50 text-amber-700',
        Kritis: 'bg-red-50 text-red-700',
      };
      return (
        <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-bold ${colorMap[row.status as string] || ''}`}>
          {row.status as string}
        </span>
      );
    },
  },
];

export function DinkesSection() {
  const [regionFilter, setRegionFilter] = useState('Semua');

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Total Stunting Provinsi"
          value="6,615"
          subtitle="anak"
          variant="gradient"
          color="red"
          trend={{ direction: 'down', value: 'Turun 8.2%', label: 'dari semester lalu' }}
        />
        <StatCard
          title="Cakupan Pemantauan"
          value="83%"
          icon={<BarChart3 size={22} />}
          variant="icon"
          color="blue"
          trend={{ direction: 'up', value: '+5%', label: 'dari bulan lalu' }}
        />
        <StatCard
          title="Posyandu Aktif"
          value="1,247"
          icon={<Building2 size={22} />}
          variant="icon"
          color="emerald"
          subtitle="dari 1,380 total posyandu"
        />
        <StatCard
          title="Anggaran Terserap"
          value="67%"
          icon={<Landmark size={22} />}
          variant="icon"
          color="purple"
          subtitle="Rp 4.2M dari Rp 6.3M"
        />
      </div>

      {/* Capaian Program */}
      <div className="bg-white rounded-2xl border border-neutral-100 shadow-sm p-6 space-y-5">
        <div className="flex justify-between items-center">
          <div>
            <h3 className="text-base font-bold text-neutral-800 font-headline flex items-center gap-2">
              <Target size={18} className="text-primary" />
              Capaian Program vs Target Nasional
            </h3>
            <p className="text-xs text-neutral-400 mt-0.5">Perbandingan capaian Jawa Tengah dengan target nasional 2024</p>
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-4">
            <ProgressBar label="Penurunan Stunting" value={72} color="primary" />
            <ProgressBar label="Cakupan IMD" value={85} color="emerald" />
            <ProgressBar label="ASI Eksklusif" value={68} color="blue" />
          </div>
          <div className="space-y-4">
            <ProgressBar label="Imunisasi Lengkap" value={91} color="emerald" />
            <ProgressBar label="PMT Coverage" value={78} color="amber" />
            <ProgressBar label="Akses Air Bersih" value={64} color="red" />
          </div>
        </div>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <ChartWidget
          title="Tren Prevalensi Stunting"
          subtitle="6 bulan terakhir (%)"
          type="bar-vertical"
          data={TREN_STUNTING}
        />
        <ChartWidget
          title="Distribusi Kab/Kota"
          subtitle="Berdasarkan tingkat prevalensi"
          type="donut"
          data={DISTRIBUSI_REGIONAL}
        />
      </div>

      {/* Filter + Table */}
      <FilterBar
        filters={[
          {
            id: 'region',
            placeholder: 'Wilayah',
            value: regionFilter,
            onChange: setRegionFilter,
            options: [
              { label: 'Wilayah: Semua', value: 'Semua' },
              { label: 'Semarang Raya', value: 'Semarang' },
              { label: 'Pantura', value: 'Pantura' },
            ],
          },
        ]}
        actions={
          <button className="inline-flex items-center gap-2 bg-primary hover:bg-primary-600 text-white rounded-xl px-5 py-2.5 text-sm font-medium transition-colors">
            <Download size={16} />
            Download Laporan
          </button>
        }
      />

      <DataTable
        title="Ranking Kabupaten/Kota — Prevalensi Stunting"
        columns={TABLE_COLUMNS}
        data={KABUPATEN_DATA}
        pageSize={5}
        onExport={() => { }}
        exportLabel="Export Laporan"
      />
    </div>
  );
}
