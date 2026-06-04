import { useState } from 'react';
import { AlertCircle, Calendar, FileWarning, Plus, Printer } from 'lucide-react';
import { StatCard } from '../components/statcard';
import { ChartWidget } from '../components/chartwidget';
import { DataTable, type Column } from '../components/datatable';
import { StatusBadge } from '../components/statusbadge';
import { FilterBar } from '../components/filterbar';

// === Dummy Data ===
interface BalitaBelumUkur {
  [key: string]: unknown;
  no: number;
  nama: string;
  ibu: string;
  usia: string;
  lastPengukuran: string;
  statusTerakhir: string;
  alamat: string;
}

const BALITA_BELUM_UKUR: BalitaBelumUkur[] = [
  { no: 1, nama: 'Fajar Ramadhan', ibu: 'Susi Hartati', usia: '14 Bln', lastPengukuran: '10 Apr 2024', statusTerakhir: 'Gizi Baik', alamat: 'RT 03/RW 02' },
  { no: 2, nama: 'Citra Dewi', ibu: 'Wati Susilowati', usia: '26 Bln', lastPengukuran: '08 Apr 2024', statusTerakhir: 'Gizi Kurang', alamat: 'RT 01/RW 05' },
  { no: 3, nama: 'Ahmad Fauzi', ibu: 'Nur Hasanah', usia: '20 Bln', lastPengukuran: '15 Apr 2024', statusTerakhir: 'Stunting', alamat: 'RT 04/RW 01' },
  { no: 4, nama: 'Melati Putri', ibu: 'Endang Purwanti', usia: '10 Bln', lastPengukuran: '22 Apr 2024', statusTerakhir: 'Gizi Baik', alamat: 'RT 02/RW 03' },
  { no: 5, nama: 'Galih Pratama', ibu: 'Sri Mulyani', usia: '32 Bln', lastPengukuran: '05 Apr 2024', statusTerakhir: 'Gizi Kurang', alamat: 'RT 05/RW 02' },
  { no: 6, nama: 'Intan Permata', ibu: 'Retno Wulan', usia: '16 Bln', lastPengukuran: '18 Apr 2024', statusTerakhir: 'Gizi Baik', alamat: 'RT 01/RW 04' },
  { no: 7, nama: 'Bayu Setiawan', ibu: 'Yuliana', usia: '28 Bln', lastPengukuran: '12 Apr 2024', statusTerakhir: 'Gizi Kurang', alamat: 'RT 03/RW 01' },
];

const KEHADIRAN_BULANAN = [
  { label: 'Des', value: 42, color: '#9ca3af' },
  { label: 'Jan', value: 48, color: '#9ca3af' },
  { label: 'Feb', value: 52, color: '#10b981' },
  { label: 'Mar', value: 45, color: '#9ca3af' },
  { label: 'Apr', value: 55, color: '#10b981' },
  { label: 'Mei', value: 51, color: '#095c3e' },
];

const DISTRIBUSI_GIZI_POSYANDU = [
  { label: 'Gizi Baik', value: 42, color: '#10b981' },
  { label: 'Gizi Kurang', value: 12, color: '#f59e0b' },
  { label: 'Stunting', value: 5, color: '#ef4444' },
  { label: 'Gizi Lebih', value: 3, color: '#3b82f6' },
];

const TABLE_COLUMNS: Column<BalitaBelumUkur>[] = [
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
  { header: 'PENGUKURAN TERAKHIR', accessor: 'lastPengukuran', className: 'text-neutral-500 font-medium' },
  {
    header: 'STATUS TERAKHIR',
    accessor: 'statusTerakhir',
    render: (row) => {
      const variantMap: Record<string, 'gizi-baik' | 'gizi-kurang' | 'stunting'> = {
        'Gizi Baik': 'gizi-baik',
        'Gizi Kurang': 'gizi-kurang',
        'Stunting': 'stunting',
      };
      return <StatusBadge variant={variantMap[row.statusTerakhir as string] || 'normal'} label={row.statusTerakhir as string} />;
    },
  },
  { header: 'ALAMAT', accessor: 'alamat', className: 'text-neutral-500' },
];

export function KaderSection() {
  const [search, setSearch] = useState('');

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Total Balita Posyandu"
          value="62"
          subtitle="anak"
          variant="gradient"
          color="primary"
          trend={{ direction: 'up', value: '+4', label: 'bulan ini' }}
        />
        <StatCard
          title="Belum Diukur Bulan Ini"
          value="17"
          icon={<AlertCircle size={22} />}
          variant="icon"
          color="amber"
          subtitle="dari 62 total balita"
        />
        <StatCard
          title="Jadwal Posyandu"
          value="18 Jun"
          icon={<Calendar size={22} />}
          variant="icon"
          color="blue"
          subtitle="Posyandu Melati, 08:00"
        />
        <StatCard
          title="Laporan Pending"
          value="3"
          icon={<FileWarning size={22} />}
          variant="icon"
          color="red"
          subtitle="belum dikirim ke bidan"
        />
      </div>

      {/* Quick Actions */}
      <div className="flex flex-wrap gap-3">
        <button className="inline-flex items-center gap-2 bg-primary hover:bg-primary-600 text-white rounded-xl px-5 py-2.5 text-sm font-medium transition-colors">
          <Plus size={16} />
          Input Data Pengukuran
        </button>
        <button className="inline-flex items-center gap-2 bg-white border border-neutral-200 hover:bg-neutral-50 text-neutral-700 rounded-xl px-5 py-2.5 text-sm font-medium transition-colors">
          <Printer size={16} />
          Cetak Laporan Posyandu
        </button>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <ChartWidget
          title="Kehadiran Balita"
          subtitle="Per bulan di Posyandu Melati"
          type="bar-vertical"
          data={KEHADIRAN_BULANAN}
        />
        <ChartWidget
          title="Status Gizi Posyandu"
          subtitle="Distribusi saat ini"
          type="donut"
          data={DISTRIBUSI_GIZI_POSYANDU}
        />
      </div>

      {/* Reminder card */}
      <div className="bg-amber-50 border border-amber-200 rounded-2xl p-5 flex items-start gap-4">
        <div className="w-10 h-10 bg-amber-100 rounded-xl flex items-center justify-center flex-shrink-0">
          <AlertCircle size={20} className="text-amber-600" />
        </div>
        <div>
          <h4 className="text-sm font-bold text-amber-800 font-headline">Pengingat</h4>
          <p className="text-sm text-amber-700 mt-1">
            Terdapat <strong>17 balita</strong> yang belum diukur bulan ini. Pastikan semua data pengukuran dilengkapi sebelum jadwal posyandu berikutnya pada <strong>18 Juni 2024</strong>.
          </p>
        </div>
      </div>

      {/* Filter + Table */}
      <FilterBar
        searchValue={search}
        onSearchChange={setSearch}
        searchPlaceholder="Cari nama balita..."
      />

      <DataTable
        title="Balita Belum Diukur Bulan Ini"
        columns={TABLE_COLUMNS}
        data={BALITA_BELUM_UKUR}
        pageSize={5}
        onExport={() => { }}
        exportLabel="Cetak Daftar"
      />
    </div>
  );
}
