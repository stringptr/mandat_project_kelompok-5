import { useState, useEffect } from 'react';
import { CalendarClock, AlertTriangle, Bell } from 'lucide-react';
import type { BackendNotifikasi } from '../useNotifikasi';
import type { NotifGroup } from '../types';
import NotifTimeline from '../NotifTimeline';
import { apiGet } from '../../../../lib/api';

interface BidanDashboard {
  statistik: {
    jadwal_kontrol: number;
    risiko_stunting: number;
    rujukan_mendesak: number;
  };
  notifikasi_risiko_stunting: Array<{
    id_pasien: number;
    nama_pasien: string;
    status_gizi: string;
    status_stunting: string;
    tanggal_monitoring: string;
  }>;
  jadwal_monitoring: Array<{
    id_pasien: number;
    nama_pasien: string;
    jadwal_kontrol: string;
    status: string;
  }>;
  rujukan_mendesak: Array<{
    id_rujukan: number;
    nama_pasien: string;
    status_rujukan: string;
    tanggal_rujukan: string;
  }>;
  laporan_bulanan: {
    bulan: string;
    jumlah_pasien_monitoring: number;
    jumlah_pasien_dirujuk: number;
  } | null;
}

interface BidanNotifikasiProps {
  notifikasi: BackendNotifikasi[];
  markAllRead: () => Promise<void>;
  toNotifGroups: (list: BackendNotifikasi[]) => NotifGroup[];
}

export default function BidanNotifikasi({ notifikasi, markAllRead, toNotifGroups }: BidanNotifikasiProps): JSX.Element {
  const [dashboard, setDashboard] = useState<BidanDashboard | null>(null);

  useEffect(() => {
    apiGet<BidanDashboard>('/notifikasi/bidan')
      .then(setDashboard)
      .catch(() => console.error('Gagal memuat dashboard bidan'));
  }, []);

  const groups = toNotifGroups(notifikasi);

  return (
    <div className="space-y-6">
      {/* Summary Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-amber-50 flex items-center justify-center">
              <AlertTriangle className="w-5 h-5 text-amber-600" />
            </div>
            <div>
              <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest">Risiko Stunting</p>
              <p className="text-2xl font-bold text-neutral-800">{dashboard?.statistik.risiko_stunting ?? '-'}</p>
            </div>
          </div>
        </div>
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center">
              <CalendarClock className="w-5 h-5 text-blue-600" />
            </div>
            <div>
              <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest">Jadwal Kontrol</p>
              <p className="text-2xl font-bold text-neutral-800">{dashboard?.statistik.jadwal_kontrol ?? '-'}</p>
            </div>
          </div>
        </div>
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-red-50 flex items-center justify-center">
              <Bell className="w-5 h-5 text-red-600" />
            </div>
            <div>
              <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest">Rujukan Mendesak</p>
              <p className="text-2xl font-bold text-neutral-800">{dashboard?.statistik.rujukan_mendesak ?? '-'}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Notification Timeline */}
      <div>
        <div className="flex items-center justify-between mb-5">
          <h3 className="text-lg font-bold text-neutral-800 font-headline">Notifikasi Terbaru</h3>
          {notifikasi.some((n) => !n.status_baca) && (
            <button onClick={markAllRead} className="text-sm font-semibold text-primary hover:text-primary-700 transition-colors">
              Tandai semua dibaca
            </button>
          )}
        </div>
        {groups.length === 0 ? (
          <div className="text-center py-12 text-neutral-400 text-sm">Belum ada notifikasi</div>
        ) : (
          <NotifTimeline groups={groups} />
        )}
      </div>
    </div>
  );
}
