import { useEffect, useState } from 'react';
import { Baby, Ruler, Weight, Calendar, Heart, Utensils, Loader2, AlertCircle } from 'lucide-react';
import { useAuth } from '../../../context/AuthContext';
import { apiGet } from '../../../lib/api';
import type { PasienDetail } from '../../../types/entities';
import type { PasienRiwayatItem } from '../../../types/entities';
import { StatCard } from '../components/statcard';
import { ChartWidget } from '../components/chartwidget';
import { DataTable, type Column } from '../components/datatable';
import { StatusBadge } from '../components/statusbadge';

interface TumbuhKembangItem {
  bulan: string;
  berat_badan: number;
  tinggi_badan: number;
}

interface RiwayatPemeriksaanData {
  riwayat: PasienRiwayatItem[];
}

interface TumbuhKembangData {
  data: TumbuhKembangItem[];
}

const CHART_COLORS = ['#ef4444', '#f59e0b', '#10b981', '#095c3e', '#3b82f6', '#8b5cf6'];

export function IbuWaliSection() {
  const { user } = useAuth();
  const idPasien = user?.idUser;

  const [pasien, setPasien] = useState<PasienDetail | null>(null);
  const [riwayat, setRiwayat] = useState<PasienRiwayatItem[]>([]);
  const [tumbuhKembang, setTumbuhKembang] = useState<TumbuhKembangItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!idPasien) return;
    setLoading(true);

    Promise.all([
      apiGet<PasienDetail>(`/monitoring/pasien/${idPasien}`).catch(() => null),
      apiGet<RiwayatPemeriksaanData>(`/monitoring/pasien/${idPasien}/riwayat-pemeriksaan`).catch(() => ({ riwayat: [] })),
      apiGet<TumbuhKembangData>(`/monitoring/pasien/${idPasien}/tumbuh-kembang`).catch(() => ({ data: [] })),
    ]).then(([p, r, tk]) => {
      if (p) setPasien(p);
      setRiwayat(r.riwayat);
      setTumbuhKembang(tk.data);
    }).finally(() => setLoading(false));
  }, [idPasien]);

  const namaAnak = pasien?.data_anak?.nama_anak || pasien?.nama || 'Anak';
  const jenisKelamin = pasien?.jenis_kelamin || '';
  const labelKelamin = jenisKelamin === 'L' || jenisKelamin === 'Laki-laki' ? 'Laki-laki' : jenisKelamin === 'P' || jenisKelamin === 'Perempuan' ? 'Perempuan' : '';

  const usia = pasien?.tanggal_lahir
    ? (() => {
        const lahir = new Date(pasien.tanggal_lahir);
        const now = new Date();
        const bulan = (now.getFullYear() - lahir.getFullYear()) * 12 + (now.getMonth() - lahir.getMonth());
        if (bulan < 24) return `${bulan} bulan`;
        const tahun = Math.floor(bulan / 12);
        const sisaBulan = bulan % 12;
        return sisaBulan > 0 ? `${tahun} tahun ${sisaBulan} bulan` : `${tahun} tahun`;
      })()
    : null;

  const bbTerbaru = riwayat.length > 0 ? String(riwayat[0].berat_badan) : '-';
  const tbTerbaru = riwayat.length > 0 ? String(riwayat[0].tinggi_badan) : '-';
  const statusGiziTerbaru = riwayat.length > 0 ? riwayat[0].status_gizi : '-';
  const bbSebelum = riwayat.length > 1 ? riwayat[1].berat_badan : null;
  const tbSebelum = riwayat.length > 1 ? riwayat[1].tinggi_badan : null;
  const selisihBB = bbSebelum !== null ? (parseFloat(bbTerbaru) - bbSebelum).toFixed(1) : null;
  const selisihTB = tbSebelum !== null ? (parseFloat(tbTerbaru) - tbSebelum).toFixed(1) : null;

  const GROWTH_CHART_DATA = tumbuhKembang.map((item, i) => ({
    label: item.bulan,
    value: item.berat_badan,
    color: CHART_COLORS[i % CHART_COLORS.length],
  }));

  const TABLE_COLUMNS: Column<PasienRiwayatItem>[] = [
    { header: 'TANGGAL', accessor: 'tanggal', className: 'font-medium text-neutral-500' },
    { header: 'BB (KG)', accessor: 'berat_badan', className: 'font-semibold text-neutral-700' },
    { header: 'TB (CM)', accessor: 'tinggi_badan', className: 'font-semibold text-neutral-700' },
    {
      header: 'STATUS GIZI',
      accessor: 'status_gizi',
      render: (row) => {
        const variant = row.status_gizi === 'Gizi Baik' ? 'gizi-baik' : 'gizi-kurang';
        return <StatusBadge variant={variant} label={row.status_gizi} />;
      },
    },
    {
      header: 'CATATAN',
      accessor: 'catatan',
      className: 'text-neutral-500',
      render: (row) => <span>{row.catatan || '-'}</span>,
    },
  ];

  if (loading) {
    return (
      <div className="flex items-center justify-center py-16">
        <div className="text-center">
          <Loader2 size={32} className="animate-spin mx-auto mb-3 text-primary" />
          <p className="text-sm text-neutral-500">Memuat data anak...</p>
        </div>
      </div>
    );
  }

  if (!pasien) {
    return (
      <div className="bg-amber-50 border border-amber-200 rounded-2xl p-6 flex items-start gap-4">
        <AlertCircle size={20} className="text-amber-600 flex-shrink-0 mt-0.5" />
        <div>
          <h4 className="text-sm font-bold text-amber-800">Data Pasien Belum Tersedia</h4>
          <p className="text-sm text-amber-700 mt-1">Silakan hubungi bidan atau kader posyandu untuk mendaftarkan data anak Anda.</p>
        </div>
      </div>
    );
  }

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
              Data Anak: {namaAnak}
            </h2>
            <p className="text-sm text-neutral-500 mt-0.5">
              {jenisKelamin && `${labelKelamin} • `}Usia: {usia ?? 'Menunggu data lengkap'}
            </p>
          </div>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Berat Badan"
          value={bbTerbaru}
          subtitle="kg"
          icon={<Weight size={22} />}
          variant="icon"
          color="primary"
          trend={selisihBB ? { direction: parseFloat(selisihBB) >= 0 ? 'up' : 'down', value: `${selisihBB} kg`, label: 'dari bulan lalu' } : undefined}
        />
        <StatCard
          title="Tinggi Badan"
          value={tbTerbaru}
          subtitle="cm"
          icon={<Ruler size={22} />}
          variant="icon"
          color="blue"
          trend={selisihTB ? { direction: parseFloat(selisihTB) >= 0 ? 'up' : 'down', value: `${selisihTB} cm`, label: 'dari bulan lalu' } : undefined}
        />
        <StatCard
          title="Status Gizi"
          value={statusGiziTerbaru}
          icon={<Heart size={22} />}
          variant="icon"
          color={statusGiziTerbaru === 'Gizi Baik' ? 'emerald' : 'amber'}
        />
        <StatCard
          title="Total Pemeriksaan"
          value={String(riwayat.length)}
          icon={<Calendar size={22} />}
          variant="icon"
          color="purple"
          subtitle="riwayat pengukuran"
        />
      </div>

      {/* Charts row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {GROWTH_CHART_DATA.length > 0 ? (
          <ChartWidget
            title="Grafik Berat Badan"
            subtitle="Riwayat pengukuran (kg)"
            type="bar-vertical"
            data={GROWTH_CHART_DATA}
          />
        ) : (
          <div className="bg-white rounded-2xl border border-neutral-100 shadow-sm p-6 flex items-center justify-center text-neutral-400 text-sm">
            Belum ada data grafik
          </div>
        )}

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
        data={riwayat}
        pageSize={5}
        emptyMessage="Belum ada riwayat pengukuran"
      />
    </div>
  );
}
