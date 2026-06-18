import { useState, useEffect } from 'react';
import { BarChart3, Activity } from 'lucide-react';
import NotifTimeline from '../NotifTimeline';
import { useNotifikasi } from '../useNotifikasi';
import { apiGet } from '../../../../lib/api';

interface StatistikNotifikasi {
  jadwal_ulang: number;
  rujukan_mendesak: number;
  risiko_stunting: number;
  notifikasi_belum_dibaca: number;
}

interface AktivitasItem {
  id_notifikasi: number;
  judul: string;
  status: string;
  timestamp: string;
}

interface AktivitasResponse {
  hari_ini: AktivitasItem[];
  kemarin: AktivitasItem[];
}

export default function DinkesNotifikasi(): JSX.Element {
  const { notifikasi, markAllRead, toNotifGroups } = useNotifikasi();
  const [statistik, setStatistik] = useState<StatistikNotifikasi | null>(null);
  const [aktivitas, setAktivitas] = useState<AktivitasResponse | null>(null);

  useEffect(() => {
    apiGet<StatistikNotifikasi>('/notifikasi/statistik').then(setStatistik).catch(() => console.error('Gagal memuat statistik'));
    apiGet<AktivitasResponse>('/notifikasi/aktivitas').then(setAktivitas).catch(() => console.error('Gagal memuat aktivitas'));
  }, []);

  const groups = toNotifGroups(notifikasi);

  return (
    <div className="space-y-8">
      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm">
          <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest mb-1">Jadwal Ulang</p>
          <p className="text-2xl font-bold text-neutral-800">{statistik?.jadwal_ulang ?? '-'}</p>
        </div>
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm">
          <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest mb-1">Rujukan Mendesak</p>
          <p className="text-2xl font-bold text-neutral-800">{statistik?.rujukan_mendesak ?? '-'}</p>
        </div>
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm">
          <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest mb-1">Risiko Stunting</p>
          <p className="text-2xl font-bold text-neutral-800">{statistik?.risiko_stunting ?? '-'}</p>
        </div>
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm">
          <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest mb-1">Belum Dibaca</p>
          <p className="text-2xl font-bold text-neutral-800">{statistik?.notifikasi_belum_dibaca ?? '-'}</p>
        </div>
      </div>

      {/* Aktivitas */}
      {aktivitas && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="bg-white rounded-2xl border border-neutral-100 shadow-sm p-5">
            <div className="flex items-center gap-2 mb-4">
              <BarChart3 className="w-4 h-4 text-primary" />
              <h3 className="font-bold text-sm text-neutral-800">Hari Ini</h3>
            </div>
            {aktivitas.hari_ini.length === 0 ? (
              <p className="text-sm text-neutral-400">Tidak ada aktivitas</p>
            ) : (
              <div className="space-y-3">
                {aktivitas.hari_ini.map((item) => (
                  <div key={item.id_notifikasi} className="flex items-center justify-between">
                    <span className="text-sm text-neutral-700">{item.judul}</span>
                    <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${item.status === 'baru' ? 'bg-blue-100 text-blue-700' : 'bg-neutral-100 text-neutral-500'}`}>
                      {item.status}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </div>
          <div className="bg-white rounded-2xl border border-neutral-100 shadow-sm p-5">
            <div className="flex items-center gap-2 mb-4">
              <Activity className="w-4 h-4 text-primary" />
              <h3 className="font-bold text-sm text-neutral-800">Kemarin</h3>
            </div>
            {aktivitas.kemarin.length === 0 ? (
              <p className="text-sm text-neutral-400">Tidak ada aktivitas</p>
            ) : (
              <div className="space-y-3">
                {aktivitas.kemarin.map((item) => (
                  <div key={item.id_notifikasi} className="text-sm text-neutral-700">{item.judul}</div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}

      {/* Timeline */}
      <div>
        <div className="flex items-center justify-between mb-5">
          <h3 className="text-lg font-bold text-neutral-800 font-headline">Notifikasi</h3>
          {notifikasi.some((n) => !n.status_baca) && (
            <button onClick={markAllRead} className="text-sm font-semibold text-primary hover:text-primary-700 transition-colors">
              Tandai semua dibaca
            </button>
          )}
        </div>
        {groups.length === 0 ? (
          <div className="text-center py-12 text-neutral-400 text-sm">Belum ada notifikasi</div>
        ) : (
          <NotifTimeline groups={groups} compactLastGroup />
        )}
      </div>
    </div>
  );
}
