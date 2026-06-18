import { useEffect } from 'react';
import { Users, FileSearch, AlertOctagon } from 'lucide-react';
import { apiGet } from '../../../lib/api';
import type { Role } from '../../../App';
import type { DashboardStats } from '../../../types/api';
import { StatCard } from '../components/StatCard';
import { useAppStore } from '../../../store/useAppStore';

interface BidanKaderSectionProps {
  currentRole: Role;
}

export function BidanKaderSection({ currentRole }: BidanKaderSectionProps): JSX.Element {
  const isBidan = currentRole === 'Bidan';

  const {
    dashboardStats,
    aktivitas,
    jadwalTerdekat,
    imunisasiPersen,
    setDashboardStats,
    setAktivitas,
    setJadwalTerdekat,
    setImunisasiPersen,
  } = useAppStore();

  useEffect(() => {
    if (dashboardStats !== null) return;

    Promise.allSettled([
      apiGet<DashboardStats>('/dashboard/stats'),
      apiGet<{ aktivitas: Record<string, unknown>[] }>('/dashboard/aktivitas'),
      apiGet<{ jadwal: Record<string, unknown>[] }>('/dashboard/jadwal-terdekat'),
      apiGet<Record<string, unknown>>('/imunisasi/statistik'),
    ]).then(([stats, akt, jadwal, imun]) => {
      if (stats.status === 'fulfilled') setDashboardStats(stats.value);
      if (akt.status === 'fulfilled') setAktivitas((akt.value.aktivitas ?? []).slice(0, 5));
      if (jadwal.status === 'fulfilled') setJadwalTerdekat((jadwal.value.jadwal ?? []).slice(0, 3));
      if (imun.status === 'fulfilled') setImunisasiPersen(Math.round((imun.value.cakupan_persentase as number) ?? 0));
    });
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const totalPasien = dashboardStats ? String(dashboardStats.total_pasien) : '...';
  const perluVerifikasi = dashboardStats ? String(dashboardStats.perlu_verifikasi) : '...';
  const tindakLanjut = dashboardStats ? String(dashboardStats.tindak_lanjut) : '...';

  return (
    <div className="space-y-5">
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <StatCard
          variant="primary"
          label="Total Pasien Aktif"
          value={totalPasien}
          icon={<Users size={22} />}
        />
        <StatCard
          variant="white"
          label="Perlu Verifikasi"
          value={perluVerifikasi}
          badge={dashboardStats ? undefined : '...'}
          icon={<FileSearch size={22} className="text-primary" />}
        />
        <StatCard
          variant="danger"
          label="Tindak Lanjut"
          value={tindakLanjut}
          icon={<AlertOctagon size={22} className="text-red-500" />}
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <div className="lg:col-span-2 bg-white rounded-2xl p-5 border border-neutral-100">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-bold text-neutral-800 font-headline">Aktivitas Terbaru</h3>
          </div>
          {aktivitas.length === 0 ? (
            <p className="text-sm text-neutral-400 text-center py-4">Belum ada aktivitas</p>
          ) : (
            aktivitas.map((item) => (
              <div key={String(item.id)} className="flex items-center gap-3 py-2 border-b border-neutral-50 last:border-0">
                <span className="text-sm font-medium text-neutral-700">{String(item.nama_pasien ?? '')}</span>
                <span className="text-xs text-neutral-500">{String(item.tindakan ?? '')}</span>
                <span className="text-xs text-neutral-400 ml-auto">{String(item.waktu ?? '')}</span>
              </div>
            ))
          )}
        </div>

        <div className="flex flex-col gap-4">
          <div className="bg-blue-50 rounded-2xl p-5 border border-blue-100">
            <div className="flex items-center gap-2 mb-4">
              <span className="text-base">🎯</span>
              <h3 className="text-sm font-bold text-neutral-800 font-headline">Target Wilayah</h3>
            </div>
            <>
              <p className="text-sm text-neutral-500 mb-2">Cakupan imunisasi: <span className="font-bold text-blue-700">{imunisasiPersen}%</span></p>
              <div className="w-full bg-blue-100 rounded-full h-2">
                <div
                  className="bg-blue-500 h-2 rounded-full transition-all duration-700"
                  style={{ width: `${Math.min(imunisasiPersen, 100)}%` }}
                />
              </div>
            </>
          </div>

          <div className="bg-white rounded-2xl p-5 border border-neutral-100 flex-1">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-bold text-neutral-800 font-headline">Jadwal Terdekat</h3>
            </div>
            {jadwalTerdekat.length === 0 ? (
              <p className="text-sm text-neutral-400 text-center py-4">Belum ada jadwal</p>
            ) : (
              jadwalTerdekat.map((j) => (
                <div key={String(j.id)} className="text-sm py-1.5 border-b border-neutral-50 last:border-0">
                  <p className="font-medium text-neutral-800">{String(j.nama_vaksin ?? '')}</p>
                  <p className="text-xs text-neutral-500">{String(j.tanggal_jadwal ?? '')} · {String(j.nama_pasien ?? '')}</p>
                </div>
              ))
            )}
          </div>
        </div>
      </div>

      {isBidan && (
        <div className="bg-primary-50 border border-primary-100 rounded-2xl p-5">
          <p className="text-sm font-bold text-primary font-headline">Status Bidan</p>
          <p className="text-xs text-neutral-600 mt-1">
            {dashboardStats
              ? `${dashboardStats.perlu_verifikasi} pemeriksaan menunggu verifikasi Anda.`
              : 'Data dashboard bidan akan tampil setelah ada aktivitas pemeriksaan.'}
          </p>
        </div>
      )}
    </div>
  );
}
