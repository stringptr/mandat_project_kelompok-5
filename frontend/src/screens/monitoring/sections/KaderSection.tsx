import { useState } from 'react';
import { AlertCircle, Calendar, FileWarning, Plus, Printer } from 'lucide-react';
import { StatCard } from '../components/statcard';
import { ChartWidget } from '../components/chartwidget';
import { DataTable, type Column } from '../components/datatable';
import { StatusBadge } from '../components/statusbadge';
import { FilterBar } from '../components/filterbar';
import ModalTambahPemeriksaan, { type FormDataTambah, type PasienOption } from '../components/modaltambah';
import ModalEditPemeriksaan, { type PemeriksaanData } from '../components/modalubah';
import ModalHapusPemeriksaan from '../components/modalhapus';

// === Types ===
interface BalitaBelumUkur {
  [key: string]: unknown;
  no: number;
  id: string;
  nama: string;
  ibu: string;
  usia: string;
  lastPengukuran: string;
  statusTerakhir: string;
  alamat: string;
  bb?: string;
  tb?: string;
}

// === Dummy Data ===
const BALITA_BELUM_UKUR: BalitaBelumUkur[] = [
  { no: 1, id: 'k1', nama: 'Fajar Ramadhan', ibu: 'Susi Hartati', usia: '14 Bln', lastPengukuran: '10 Apr 2024', statusTerakhir: 'Gizi Baik', alamat: 'RT 03/RW 02', bb: '9.5', tb: '74.0' },
  { no: 2, id: 'k2', nama: 'Citra Dewi', ibu: 'Wati Susilowati', usia: '26 Bln', lastPengukuran: '08 Apr 2024', statusTerakhir: 'Gizi Kurang', alamat: 'RT 01/RW 05', bb: '10.1', tb: '80.2' },
  { no: 3, id: 'k3', nama: 'Ahmad Fauzi', ibu: 'Nur Hasanah', usia: '20 Bln', lastPengukuran: '15 Apr 2024', statusTerakhir: 'Stunting', alamat: 'RT 04/RW 01', bb: '8.8', tb: '72.5' },
  { no: 4, id: 'k4', nama: 'Melati Putri', ibu: 'Endang Purwanti', usia: '10 Bln', lastPengukuran: '22 Apr 2024', statusTerakhir: 'Gizi Baik', alamat: 'RT 02/RW 03', bb: '7.5', tb: '68.0' },
  { no: 5, id: 'k5', nama: 'Galih Pratama', ibu: 'Sri Mulyani', usia: '32 Bln', lastPengukuran: '05 Apr 2024', statusTerakhir: 'Gizi Kurang', alamat: 'RT 05/RW 02', bb: '11.2', tb: '82.0' },
  { no: 6, id: 'k6', nama: 'Intan Permata', ibu: 'Retno Wulan', usia: '16 Bln', lastPengukuran: '18 Apr 2024', statusTerakhir: 'Gizi Baik', alamat: 'RT 01/RW 04', bb: '9.0', tb: '75.5' },
  { no: 7, id: 'k7', nama: 'Bayu Setiawan', ibu: 'Yuliana', usia: '28 Bln', lastPengukuran: '12 Apr 2024', statusTerakhir: 'Gizi Kurang', alamat: 'RT 03/RW 01', bb: '10.8', tb: '79.0' },
];

// Derive pasien list untuk modal tambah
const PASIEN_LIST: PasienOption[] = BALITA_BELUM_UKUR.map(b => ({
  id: b.id,
  nama: b.nama,
  namaIbu: b.ibu,
  usia: b.usia,
  nik: `3273${b.id.padEnd(12, '0')}`,
  statusGizi: b.statusTerakhir,
}));

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

export function KaderSection() {
  const [search, setSearch] = useState('');

  // — Modal Tambah —
  const [modalTambah, setModalTambah] = useState(false);

  // — Modal Edit —
  const [modalEdit, setModalEdit] = useState(false);
  const [editTarget, setEditTarget] = useState<{ nama: string; data: PemeriksaanData } | null>(null);

  // — Modal Hapus —
  const [modalHapus, setModalHapus] = useState(false);
  const [hapusTarget, setHapusTarget] = useState<{ id: string; nama: string; tanggal: string } | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  // === Handlers ===
  const openEdit = (row: BalitaBelumUkur) => {
    setEditTarget({
      nama: row.nama,
      data: {
        id: row.id,
        beratBadan: row.bb ?? '',
        tinggiBadan: row.tb ?? '',
        lingkarKepala: '',
        terakhirDiperbarui: {
          nama: 'Kader Siti Aminah',
          inisial: 'SA',
          tanggal: row.lastPengukuran,
        },
      },
    });
    setModalEdit(true);
  };

  const openHapus = (row: BalitaBelumUkur) => {
    setHapusTarget({ id: row.id, nama: row.nama, tanggal: row.lastPengukuran });
    setModalHapus(true);
  };

  const handleTambahSubmit = (pasienId: string, namaAnak: string, data: FormDataTambah) => {
    console.log('Tambah data pasien:', pasienId, namaAnak, data);
    // TODO: POST ke API
  };

  const handleEditSubmit = (id: string, data: Omit<PemeriksaanData, 'id' | 'terakhirDiperbarui'>) => {
    console.log('Edit id:', id, data);
    // TODO: PUT ke API
  };

  const handleHapusConfirm = async () => {
    if (!hapusTarget) return;
    setIsDeleting(true);
    try {
      await new Promise(r => setTimeout(r, 1000)); // ganti dengan API call
      console.log('Hapus id:', hapusTarget.id);
      setModalHapus(false);
      setHapusTarget(null);
    } finally {
      setIsDeleting(false);
    }
  };

  // === Kolom tabel (dengan aksi) ===
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
    {
      header: 'AKSI',
      accessor: 'id',
      render: (row) => (
        <div className="flex items-center gap-1.5">
          <button
            onClick={() => openEdit(row)}
            className="flex items-center gap-1 text-xs text-emerald-700 hover:text-emerald-800 bg-emerald-50 hover:bg-emerald-100 px-2 py-1.5 rounded-lg transition-colors"
          >
            <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            Edit
          </button>
          <button
            onClick={() => openHapus(row)}
            className="flex items-center gap-1 text-xs text-red-600 hover:text-red-700 bg-red-50 hover:bg-red-100 px-2 py-1.5 rounded-lg transition-colors"
          >
            <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Hapus
          </button>
        </div>
      ),
    },
  ];

  const filteredData = BALITA_BELUM_UKUR.filter(row =>
    !search ||
    row.nama.toLowerCase().includes(search.toLowerCase()) ||
    row.ibu.toLowerCase().includes(search.toLowerCase()),
  );

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
        <button
          onClick={() => setModalTambah(true)}
          className="inline-flex items-center gap-2 bg-primary hover:bg-primary-600 text-white rounded-xl px-5 py-2.5 text-sm font-medium transition-colors"
        >
          <Plus size={16} />
          Input Data Pengukuran
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
        data={filteredData}
        pageSize={5}
        onExport={() => { }}
        exportLabel="Cetak Daftar"
      />

      {/* ── Modals ── */}
      <ModalTambahPemeriksaan
        isOpen={modalTambah}
        onClose={() => setModalTambah(false)}
        pasienList={PASIEN_LIST}
        onSubmit={handleTambahSubmit}
      />

      <ModalEditPemeriksaan
        isOpen={modalEdit}
        onClose={() => setModalEdit(false)}
        namaAnak={editTarget?.nama ?? ''}
        data={editTarget?.data ?? null}
        onSubmit={handleEditSubmit}
      />

      <ModalHapusPemeriksaan
        isOpen={modalHapus}
        onClose={() => {
          setModalHapus(false);
          setHapusTarget(null);
        }}
        namaAnak={hapusTarget?.nama ?? ''}
        tanggalPemeriksaan={hapusTarget?.tanggal ?? ''}
        onConfirm={handleHapusConfirm}
        isLoading={isDeleting}
      />
    </div>
  );
}