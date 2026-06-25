import { useEffect, useState } from 'react';
import { Users, AlertTriangle, Target, Calendar, Bell, FileText } from 'lucide-react';
import { apiGet } from '../../../lib/api';
import { useNotification } from '../../../context/NotificationContext';
import { useAppStore } from '../../../store/useAppStore';
import { TrendChart } from '../components/TrendChart';
import { PrevalensiMap } from '../components/PrevalensiMap';
import type { DashboardStats, DistribusiGiziItem, TrenStuntingItem, StuntingWilayahItem } from '../../../types/api';

interface NotifItem {
  id_notifikasi: number;
  judul: string;
  pesan: string | null;
  tipe_notifikasi: string;
  status_baca: boolean;
  tanggal_kirim: string;
}

interface LaporanItem {
  wilayah: string;
  jumlah_pasien_dirujuk: number;
  jumlah_pasien_diterima: number;
  jumlah_pasien_diproses: number;
}

const TIPE_DOT: Record<string, string> = {
  Pemeriksaan: 'bg-blue-500',
  Imunisasi: 'bg-green-500',
  Rujukan: 'bg-orange-500',
  Edukasi: 'bg-purple-500',
  Pengingat: 'bg-yellow-500',
};

export function DinkesSection(): JSX.Element {
  const notify = useNotification();

  const {
    dashboardStats,
    distribusiGizi,
    trenStunting,
    stuntingPerWilayah,
    kehadiranBulanan,
    setDashboardStats,
    setDistribusiGizi,
    setTrenStunting,
    setStuntingPerWilayah,
    setKehadiranBulanan,
  } = useAppStore();

  const filteredWilayah = stuntingPerWilayah.filter((w) => w.nama_wilayah !== 'Tidak Diketahui');

  useEffect(() => {
    if (dashboardStats !== null && filteredWilayah.length > 0) return;

    Promise.allSettled([
      apiGet<DashboardStats>('/dashboard/stats'),
      apiGet<{ distribusi: DistribusiGiziItem[] }>('/dashboard/distribusi-gizi'),
      apiGet<{ tren: TrenStuntingItem[] }>('/dashboard/tren-stunting'),
      apiGet<{ wilayah: StuntingWilayahItem[] }>('/dashboard/stunting-per-wilayah'),
      apiGet<{ tren: TrenStuntingItem[] }>('/dashboard/kehadiran-bulanan'),
    ]).then(([stats, distrib, tren, wilayah, kehadiran]) => {
      if (stats.status === 'fulfilled') setDashboardStats(stats.value);
      else notify.error('Gagal memuat data statistik dashboard.');

      if (distrib.status === 'fulfilled') setDistribusiGizi(distrib.value.distribusi ?? []);
      if (tren.status === 'fulfilled') setTrenStunting(tren.value.tren ?? []);
      if (wilayah.status === 'fulfilled') setStuntingPerWilayah((wilayah.value.wilayah ?? []).filter((w: StuntingWilayahItem) => w.nama_wilayah !== 'Tidak Diketahui'));
      if (kehadiran.status === 'fulfilled') setKehadiranBulanan(kehadiran.value.tren ?? []);
    });
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const totalPasien = dashboardStats ? dashboardStats.total_pasien.toLocaleString('id-ID') : '...';
  const kasusStunting = dashboardStats ? dashboardStats.kasus_stunting.toLocaleString('id-ID') : '...';
  const capaian = dashboardStats?.cakupan_persentase != null ? `${dashboardStats.cakupan_persentase}%` : '...';

  const trenChartData = trenStunting.map((t) => ({ bulan: t.bulan, nilai: t.jumlah }));
  const kehadiranChartData = kehadiranBulanan.map((t) => ({ bulan: t.bulan, nilai: t.jumlah }));

  const totalDistribusi = distribusiGizi.reduce((s, d) => s + d.jumlah, 0) || 1;

  const [notifikasi, setNotifikasi] = useState<NotifItem[]>([]);
  const [laporan, setLaporan] = useState<LaporanItem[]>([]);

  useEffect(() => {
    apiGet<{ notifikasi: NotifItem[] }>('/notifikasi?per_page=5')
      .then((res) => setNotifikasi(res.notifikasi ?? []))
      .catch(() => {});
    apiGet<{ laporan: LaporanItem[] }>('/laporan/tindak-lanjut')
      .then((res) => setLaporan(res.laporan ?? []))
      .catch(() => {});
  }, []);

  const levelColor = (level: string) => {
    if (level === 'tinggi') return 'text-red-600 bg-red-50';
    if (level === 'sedang') return 'text-amber-600 bg-amber-50';
    return 'text-emerald-600 bg-emerald-50';
  };

  return (
    <div className="space-y-5">
      <div className="flex items-start justify-between flex-wrap gap-3">
        <div>
          <h3 className="text-xl font-bold text-neutral-900 font-headline">Statistik Kesehatan Wilayah</h3>
          <p className="text-xs text-neutral-500 font-body mt-1">
            Monitoring gizi maternal dan anak. Data diperbarui secara real-time dari posyandu.
          </p>
        </div>
        <div className="flex items-center gap-2 flex-shrink-0">
          <button className="flex items-center gap-1.5 bg-primary-50 border border-primary-100 text-primary text-xs font-semibold px-3 py-2 rounded-xl hover:bg-primary-100 transition-colors font-body">
            <Calendar size={13} /> Real-Time
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="bg-primary rounded-2xl p-5 relative overflow-hidden">
          <div className="absolute -bottom-6 -right-6 w-28 h-28 rounded-full bg-white/10" />
          <Users size={22} className="text-white/80 mb-3" />
          <p className="text-white/70 text-[10px] font-semibold uppercase tracking-wide font-body">Total Pasien Terdaftar</p>
          <p className="text-4xl font-bold text-white font-headline leading-none mt-1">{totalPasien}</p>
        </div>
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 relative">
          <span className="absolute top-4 right-4 text-[10px] font-bold bg-red-500 text-white px-2 py-0.5 rounded-full">STUNTING</span>
          <AlertTriangle size={22} className="text-red-500 mb-3" />
          <p className="text-neutral-500 text-[10px] font-semibold uppercase tracking-wide font-body">Kasus Stunting Aktif</p>
          <p className="text-4xl font-bold text-red-500 font-headline leading-none mt-1">{kasusStunting}</p>
        </div>
        <div className="bg-blue-600 rounded-2xl p-5 relative overflow-hidden">
          <div className="absolute -bottom-6 -right-6 w-28 h-28 rounded-full bg-white/10" />
          <Target size={22} className="text-white/80 mb-3" />
          <p className="text-white/70 text-[10px] font-semibold uppercase tracking-wide font-body">Capaian Pemantauan</p>
          <p className="text-4xl font-bold text-white font-headline leading-none mt-1">{capaian}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="bg-white rounded-2xl p-5 border border-neutral-100">
          <h4 className="text-sm font-bold text-neutral-800 font-headline mb-3">Tren Stunting</h4>
          {trenChartData.length >= 2 ? (
            <TrendChart data={trenChartData} color="#ef4444" fillColor="rgba(239,68,68,0.08)" />
          ) : (
            <p className="text-sm text-neutral-400 text-center py-8">Belum ada data tren stunting</p>
          )}
        </div>

        <div className="bg-white rounded-2xl p-5 border border-neutral-100">
          <h4 className="text-sm font-bold text-neutral-800 font-headline mb-3">Kehadiran Posyandu Bulanan</h4>
          {kehadiranChartData.length >= 2 ? (
            <TrendChart data={kehadiranChartData} color="#095c3e" fillColor="rgba(9,92,62,0.08)" />
          ) : (
            <p className="text-sm text-neutral-400 text-center py-8">Belum ada data kehadiran</p>
          )}
        </div>
      </div>

      <div className="bg-white rounded-2xl p-5 border border-neutral-100">
        <h4 className="text-sm font-bold text-neutral-800 font-headline mb-4">Prevalensi Stunting Per Wilayah</h4>
        {filteredWilayah.length === 0 ? (
          <p className="text-sm text-neutral-400 text-center py-8">Belum ada data wilayah</p>
        ) : (
          <PrevalensiMap
            data={filteredWilayah.map((w) => ({
              nama: w.nama_wilayah,
              prevalensi: w.prevalensi,
              jumlahKasus: w.jumlah_kasus,
              level: w.level,
            }))}
          />
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="bg-white rounded-2xl p-5 border border-neutral-100">
          <h4 className="text-sm font-bold text-neutral-800 font-headline mb-4">Distribusi Status Gizi</h4>
          {distribusiGizi.length === 0 ? (
            <p className="text-sm text-neutral-400 text-center py-8">Belum ada data distribusi</p>
          ) : (
            <div className="space-y-3">
              {distribusiGizi.map((d) => {
                const pct = Math.round((d.jumlah / totalDistribusi) * 100);
                const isStunting = d.status_gizi.toLowerCase().includes('stunting');
                const isKurang = d.status_gizi.toLowerCase().includes('kurang');
                const barColor = isStunting ? 'bg-red-500' : isKurang ? 'bg-amber-500' : 'bg-emerald-500';
                return (
                  <div key={d.status_gizi}>
                    <div className="flex justify-between text-xs mb-1">
                      <span className="font-medium text-neutral-700">{d.status_gizi}</span>
                      <span className="text-neutral-500">{d.jumlah} ({pct}%)</span>
                    </div>
                    <div className="w-full bg-neutral-100 rounded-full h-2">
                      <div
                        className={`h-2 rounded-full ${barColor} transition-all duration-500`}
                        style={{ width: `${pct}%` }}
                      />
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>

        <div className="bg-white rounded-2xl p-5 border border-neutral-100">
          <h4 className="text-sm font-bold text-neutral-800 font-headline mb-4">Stunting Per Wilayah</h4>
          {filteredWilayah.length === 0 ? (
            <p className="text-sm text-neutral-400 text-center py-8">Belum ada data wilayah</p>
          ) : (
            <div className="space-y-2 max-h-48 overflow-y-auto pr-1">
              {filteredWilayah.map((w) => (
                <div key={w.nama_wilayah} className="flex items-center justify-between py-2 border-b border-neutral-50 last:border-0">
                  <div>
                    <p className="text-sm font-semibold text-neutral-800">{w.nama_wilayah}</p>
                    <p className="text-xs text-neutral-500">{w.jumlah_kasus} kasus dari {w.total_balita} balita</p>
                  </div>
                  <span className={`text-xs font-bold px-2 py-1 rounded-full ${levelColor(w.level)}`}>
                    {w.prevalensi.toFixed(1)}%
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="bg-white rounded-2xl p-5 border border-neutral-100">
          <div className="flex items-center gap-2 mb-3">
            <Bell size={16} className="text-amber-500" />
            <h4 className="text-sm font-bold text-neutral-800 font-headline">Notifikasi Terbaru</h4>
          </div>
          {notifikasi.length === 0 ? (
            <p className="text-sm text-neutral-400 text-center py-6">Belum ada notifikasi</p>
          ) : (
            <div className="space-y-2 max-h-52 overflow-y-auto">
              {notifikasi.map((n) => (
                <div key={n.id_notifikasi} className="flex items-start gap-3 py-2.5 border-b border-neutral-50 last:border-0">
                  <span className={`w-2 h-2 rounded-full mt-1.5 flex-shrink-0 ${TIPE_DOT[n.tipe_notifikasi] ?? 'bg-neutral-300'}`} />
                  <div className="min-w-0 flex-1">
                    <p className="text-xs font-semibold text-neutral-800 truncate">{n.judul}</p>
                    <p className="text-xs text-neutral-500 mt-0.5 line-clamp-1">{n.pesan || ''}</p>
                  </div>
                  <span className="text-[10px] text-neutral-400 flex-shrink-0">{new Date(n.tanggal_kirim).toLocaleDateString('id-ID')}</span>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="bg-white rounded-2xl p-5 border border-neutral-100">
          <div className="flex items-center gap-2 mb-3">
            <FileText size={16} className="text-emerald-500" />
            <h4 className="text-sm font-bold text-neutral-800 font-headline">Laporan Tindak Lanjut</h4>
          </div>
          {laporan.length === 0 ? (
            <p className="text-sm text-neutral-400 text-center py-6">Belum ada data laporan</p>
          ) : (
            <div className="space-y-1 max-h-52 overflow-y-auto">
              <div className="grid grid-cols-4 text-[10px] font-bold text-neutral-400 uppercase tracking-wide px-1 pb-2 border-b border-neutral-100">
                <span>Wilayah</span>
                <span className="text-center">Dirujuk</span>
                <span className="text-center">Diterima</span>
                <span className="text-center">Diproses</span>
              </div>
              {laporan.map((l) => (
                <div key={l.wilayah} className="grid grid-cols-4 text-xs py-2 border-b border-neutral-50 last:border-0 items-center hover:bg-neutral-50 rounded px-1">
                  <span className="font-medium text-neutral-800 truncate">{l.wilayah}</span>
                  <span className="text-center text-amber-600 font-medium">{l.jumlah_pasien_dirujuk}</span>
                  <span className="text-center text-emerald-600 font-medium">{l.jumlah_pasien_diterima}</span>
                  <span className="text-center text-blue-600 font-medium">{l.jumlah_pasien_diproses}</span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
