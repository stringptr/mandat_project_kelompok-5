import { useState } from 'react';
import { AlertTriangle, ClipboardCheck, CalendarDays, Plus, FileText } from 'lucide-react';
import { StatCard } from '../components/statcard';
import { ChartWidget } from '../components/chartwidget';
import { DataTable, type Column } from '../components/datatable';
import { StatusBadge } from '../components/statusbadge';
import { FilterBar } from '../components/filterbar';

// === Dummy Data ===
interface BalitaUrgent {
  [key: string]: unknown;
  no: number;
  nama: string;
  ibu: string;
  usia: string;
  bb: string;
  tb: string;
  statusGizi: string;
  verifikasi: string;
  prioritas: string;
}

const BALITA_URGENT: BalitaUrgent[] = [
  { no: 1, nama: 'Danu Saputra', ibu: 'Linda Sari', usia: '38 Bln', bb: '12.1', tb: '84.2', statusGizi: 'Stunting', verifikasi: 'Pending', prioritas: 'Tinggi' },
  { no: 2, nama: 'Arka Mahendra', ibu: 'Rina Marlina', usia: '24 Bln', bb: '10.2', tb: '76.0', statusGizi: 'Gizi Kurang', verifikasi: 'Pending', prioritas: 'Tinggi' },
  { no: 3, nama: 'Nabila Putri', ibu: 'Dewi Anggraini', usia: '18 Bln', bb: '8.5', tb: '72.1', statusGizi: 'Gizi Kurang', verifikasi: 'Verified', prioritas: 'Sedang' },
  { no: 4, nama: 'Rizki Aditya', ibu: 'Yuni Rahayu', usia: '30 Bln', bb: '11.0', tb: '80.5', statusGizi: 'Stunting', verifikasi: 'Pending', prioritas: 'Tinggi' },
  { no: 5, nama: 'Sari Indah', ibu: 'Mega Lestari', usia: '12 Bln', bb: '7.8', tb: '68.0', statusGizi: 'Gizi Kurang', verifikasi: 'Verified', prioritas: 'Sedang' },
  { no: 6, nama: 'Bima Prasetyo', ibu: 'Ani Wulandari', usia: '42 Bln', bb: '13.5', tb: '90.0', statusGizi: 'Stunting', verifikasi: 'Pending', prioritas: 'Tinggi' },
];

const DISTRIBUSI_GIZI = [
  { label: 'Gizi Baik', value: 68, color: '#10b981' },
  { label: 'Gizi Kurang', value: 18, color: '#f59e0b' },
  { label: 'Stunting', value: 9, color: '#ef4444' },
  { label: 'Gizi Lebih', value: 5, color: '#3b82f6' },
];

const TABLE_COLUMNS: Column<BalitaUrgent>[] = [
  { header: 'NO', accessor: 'no', className: 'font-medium text-neutral-400 w-12' },
  {
    header: 'NAMA BALITA',
    accessor: 'nama',
    render: (row) => (
      <div>
        <div className="font-bold text-primary">{row.nama}</div>
        <div className="text-[11px] text-neutral-400 mt-0.5">Ibu: {row.ibu}</div>
      </div>
    ),
  },
  { header: 'USIA', accessor: 'usia', className: 'font-semibold text-neutral-600' },
  { header: 'BB (KG)', accessor: 'bb', className: 'font-semibold text-neutral-600' },
  { header: 'TB (CM)', accessor: 'tb', className: 'font-semibold text-neutral-600' },
  {
    header: 'STATUS GIZI',
    accessor: 'statusGizi',
    render: (row) => {
      const v = row.statusGizi === 'Stunting' ? 'stunting' : 'gizi-kurang';
      return <StatusBadge variant={v} label={row.statusGizi} />;
    },
  },
  {
    header: 'VERIFIKASI',
    accessor: 'verifikasi',
    render: (row) => {
      const v = row.verifikasi === 'Verified' ? 'verified' : 'pending';
      return <StatusBadge variant={v} />;
    },
  },
  {
    header: 'PRIORITAS',
    accessor: 'prioritas',
    render: (row) => {
      const v = row.prioritas === 'Tinggi' ? 'urgent' : 'info';
      return <StatusBadge variant={v} label={row.prioritas} dot />;
    },
  },
];

export function BidanSection() {
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('Semua');

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Total Balita Dipantau"
          value="156"
          subtitle="anak"
          variant="gradient"
          color="primary"
          trend={{ direction: 'up', value: '+8', label: 'bulan ini' }}
        />
        <StatCard
          title="Kasus Stunting Aktif"
          value="12"
          icon={<AlertTriangle size={22} />}
          variant="icon"
          color="red"
          trend={{ direction: 'down', value: '-3', label: 'dari bulan lalu' }}
        />
        <StatCard
          title="Perlu Verifikasi"
          value="7"
          icon={<ClipboardCheck size={22} />}
          variant="icon"
          color="amber"
          subtitle="data baru belum diverifikasi"
        />
        <StatCard
          title="Jadwal Posyandu"
          value="3"
          icon={<CalendarDays size={22} />}
          variant="icon"
          color="blue"
          subtitle="sesi bulan ini tersisa"
        />
      </div>

      {/* Quick Actions */}
      <div className="flex flex-wrap gap-3">
        <button className="inline-flex items-center gap-2 bg-primary hover:bg-primary-600 text-white rounded-xl px-5 py-2.5 text-sm font-medium transition-colors">
          <Plus size={16} />
          Tambah Data Pengukuran
        </button>
        <button className="inline-flex items-center gap-2 bg-white border border-neutral-200 hover:bg-neutral-50 text-neutral-700 rounded-xl px-5 py-2.5 text-sm font-medium transition-colors">
          <FileText size={16} />
          Buat Laporan Bulanan
        </button>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <ChartWidget
          title="Distribusi Status Gizi"
          subtitle="Wilayah Binaan Melati"
          type="donut"
          data={DISTRIBUSI_GIZI}
        />
        <ChartWidget
          title="Kasus per Bulan"
          subtitle="Stunting & Gizi Kurang - 6 bulan terakhir"
          type="bar-vertical"
          data={[
            { label: 'Des', value: 18, color: '#ef4444' },
            { label: 'Jan', value: 16, color: '#ef4444' },
            { label: 'Feb', value: 15, color: '#f59e0b' },
            { label: 'Mar', value: 14, color: '#f59e0b' },
            { label: 'Apr', value: 13, color: '#10b981' },
            { label: 'Mei', value: 12, color: '#10b981' },
          ]}
        />
      </div>

      {/* Filter + Table */}
      <FilterBar
        searchValue={search}
        onSearchChange={setSearch}
        searchPlaceholder="Cari nama balita atau ibu..."
        filters={[
          {
            id: 'status',
            placeholder: 'Status Gizi',
            value: statusFilter,
            onChange: setStatusFilter,
            options: [
              { label: 'Status: Semua', value: 'Semua' },
              { label: 'Stunting', value: 'Stunting' },
              { label: 'Gizi Kurang', value: 'Gizi Kurang' },
            ],
          },
        ]}
      />

      <DataTable
        title="Balita Perlu Tindak Lanjut"
        columns={TABLE_COLUMNS}
        data={BALITA_URGENT}
        pageSize={5}
        onExport={() => { }}
      />
    </div>
  );
}
