import { useState } from 'react';
import { AlertTriangle, ClipboardCheck, CalendarDays, Plus, FileText } from 'lucide-react';
import { StatCard } from '../components/statcard';
import { ChartWidget } from '../components/chartwidget';
import { DataTable, type Column } from '../components/datatable';
import { StatusBadge } from '../components/statusbadge';
import { FilterBar } from '../components/filterbar';
import ModalTambahPemeriksaan, { type FormDataTambah, type PasienOption } from '../components/Modaltambah';
import ModalEditPemeriksaan, { type PemeriksaanData } from '../components/modalubah';
import ModalHapusPemeriksaan from '../components/modalhapus';

// === Types ===
interface BalitaUrgent {
  [key: string]: unknown;
  no: number;
  id: string;
  nama: string;
  ibu: string;
  usia: string;
  bb: string;
  tb: string;
  statusGizi: string;
  verifikasi: string;
  prioritas: string;
}

// === Dummy Data ===
const BALITA_URGENT: BalitaUrgent[] = [
  { no: 1, id: 'b1', nama: 'Danu Saputra', ibu: 'Linda Sari', usia: '38 Bln', bb: '12.1', tb: '84.2', statusGizi: 'Stunting', verifikasi: 'Pending', prioritas: 'Tinggi' },
  { no: 2, id: 'b2', nama: 'Arka Mahendra', ibu: 'Rina Marlina', usia: '24 Bln', bb: '10.2', tb: '76.0', statusGizi: 'Gizi Kurang', verifikasi: 'Pending', prioritas: 'Tinggi' },
  { no: 3, id: 'b3', nama: 'Nabila Putri', ibu: 'Dewi Anggraini', usia: '18 Bln', bb: '8.5', tb: '72.1', statusGizi: 'Gizi Kurang', verifikasi: 'Verified', prioritas: 'Sedang' },
  { no: 4, id: 'b4', nama: 'Rizki Aditya', ibu: 'Yuni Rahayu', usia: '30 Bln', bb: '11.0', tb: '80.5', statusGizi: 'Stunting', verifikasi: 'Pending', prioritas: 'Tinggi' },
  { no: 5, id: 'b5', nama: 'Sari Indah', ibu: 'Mega Lestari', usia: '12 Bln', bb: '7.8', tb: '68.0', statusGizi: 'Gizi Kurang', verifikasi: 'Verified', prioritas: 'Sedang' },
  { no: 6, id: 'b6', nama: 'Bima Prasetyo', ibu: 'Ani Wulandari', usia: '42 Bln', bb: '13.5', tb: '90.0', statusGizi: 'Stunting', verifikasi: 'Pending', prioritas: 'Tinggi' },
];

// Daftar pasien untuk step-1 modal tambah (bisa diganti data dari API)
const PASIEN_LIST: PasienOption[] = BALITA_URGENT.map(b => ({
  id: b.id,
  nama: b.nama,
  namaIbu: b.ibu,
  usia: b.usia,
  nik: `3273${b.id.padEnd(12, '0')}`,
  statusGizi: b.statusGizi,
}));

const DISTRIBUSI_GIZI = [
  { label: 'Gizi Baik', value: 68, color: '#10b981' },
  { label: 'Gizi Kurang', value: 18, color: '#f59e0b' },
  { label: 'Stunting', value: 9, color: '#ef4444' },
  { label: 'Gizi Lebih', value: 5, color: '#3b82f6' },
];

export function BidanSection() {
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('Semua');

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
  const openTambah = () => setModalTambah(true);

  const openEdit = (row: BalitaUrgent) => {
    setEditTarget({
      nama: row.nama,
      data: {
        id: row.id,
        beratBadan: row.bb,
        tinggiBadan: row.tb,
        lingkarKepala: '',
        terakhirDiperbarui: {
          nama: 'Bidan Sri Lestari',
          inisial: 'SL',
          tanggal: '12 Mei 2024, 10:45 WIB',
        },
      },
    });
    setModalEdit(true);
  };

  const openHapus = (row: BalitaUrgent) => {
    setHapusTarget({ id: row.id, nama: row.nama, tanggal: '12 Mei 2024' });
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

  // === Tabel dengan kolom Aksi ===
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

  // === Filtered data ===
  const filteredData = BALITA_URGENT.filter(row => {
    const matchSearch =
      !search ||
      row.nama.toLowerCase().includes(search.toLowerCase()) ||
      row.ibu.toLowerCase().includes(search.toLowerCase());
    const matchStatus = statusFilter === 'Semua' || row.statusGizi === statusFilter;
    return matchSearch && matchStatus;
  });

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
        <button
          onClick={() => openTambah()}
          className="inline-flex items-center gap-2 bg-primary hover:bg-primary-600 text-white rounded-xl px-5 py-2.5 text-sm font-medium transition-colors"
        >
          <Plus size={16} />
          Tambah Data Pengukuran
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
        data={filteredData}
        pageSize={5}
        onExport={() => { }}
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