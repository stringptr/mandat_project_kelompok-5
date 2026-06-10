import { BarChart3, AlertTriangle, FileText, Bell, Plus } from 'lucide-react';
import NotifTimeline from '../NotifTimeline';
import type { NotifGroup } from '../types';

interface DinkesNotifikasiProps {
  // reserved for future use
}

const NOTIF_GROUPS: NotifGroup[] = [
  {
    groupLabel: 'Hari Ini',
    items: [
      {
        id: 'd1',
        title: 'Lonjakan Stunting: Kecamatan Banyumanik',
        description:
          'Terdeteksi kenaikan 12% kasus stunting baru dalam 30 hari terakhir. Perlu intervensi lintas sektor segera.',
        time: '10:00 AM',
        category: 'urgent',
        tags: [
          { label: 'Mendesak', color: '#ef4444', bg: '#fee2e2' },
          { label: 'Regional', color: '#4b5563', bg: '#f1f5f9' },
        ],
        actionLabel: 'Lihat Analitik',
      },
      {
        id: 'd2',
        title: 'Laporan Bulanan Puskesmas Siap Ditinjau',
        description:
          '14 Puskesmas telah mengunggah laporan rekapitulasi gizi periode September 2023. Menunggu review Dinkes.',
        time: '08:30 AM',
        category: 'report',
        tags: [{ label: 'Laporan', color: '#6b7280', bg: '#f1f5f9' }],
        actionLabel: 'Review Laporan',
      },
    ],
  },
  {
    groupLabel: 'Kemarin',
    items: [
      {
        id: 'd3',
        title: 'Stok PMT Wilayah Utara Kritis',
        description:
          '3 Puskesmas di wilayah utara melaporkan stok PMT tersisa < 10%. Perlu distribusi darurat dalam 48 jam.',
        time: 'Kemarin, 15:45',
        category: 'urgent',
        tags: [
          { label: 'Stok PMT', color: '#ef4444', bg: '#fee2e2' },
          { label: 'Wilayah Utara', color: '#4b5563', bg: '#f1f5f9' },
        ],
        actionLabel: 'Koordinasi Distribusi',
      },
      {
        id: 'd4',
        title: 'Capaian Imunisasi Bulan Agustus',
        description:
          'Rekapitulasi capaian imunisasi dasar lengkap Agustus telah tersedia. Coverage: 87.3% dari target.',
        time: 'Kemarin, 11:20',
        category: 'success',
        tags: [{ label: 'Imunisasi', color: '#059669', bg: '#d1fae5' }],
        actionLabel: 'Lihat Laporan',
      },
    ],
  },
  {
    groupLabel: 'Pekan Lalu',
    items: [
      {
        id: 'd5',
        title: 'Rapat Koordinasi Lintas Sektor',
        description: '',
        time: '4 Sept 2023',
        category: 'info',
        tags: [],
      },
      {
        id: 'd6',
        title: 'Data EPPGBM Jawa Tengah Diperbarui',
        description: '',
        time: '2 Sept 2023',
        category: 'success',
        tags: [],
      },
    ],
  },
];

export default function DinkesNotifikasi(_props: DinkesNotifikasiProps): JSX.Element {
  return (
    <div>
      {/* Summary Cards – 3 cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-5 mb-10">
        {/* Laporan Menunggu */}
        <div className="bg-primary rounded-3xl p-6 text-white flex flex-col justify-between min-h-[150px] relative overflow-hidden shadow-sm">
          <div className="flex items-center justify-between">
            <FileText className="w-7 h-7 text-white/90" />
            <span className="bg-white/20 backdrop-blur-sm text-white text-[10px] font-bold px-3 py-1 rounded uppercase tracking-widest">
              Prioritas
            </span>
          </div>
          <div className="mt-4">
            <p className="text-4xl font-bold font-headline leading-none mb-1">14</p>
            <p className="text-sm text-white/85 font-medium">Laporan Menunggu Review</p>
          </div>
          <div className="absolute -bottom-4 -right-4 w-24 h-24 bg-white/10 rounded-full" />
        </div>

        {/* Kecamatan Kritis */}
        <div className="bg-orange-50 rounded-3xl p-6 border border-orange-100 flex flex-col justify-between min-h-[150px] shadow-sm">
          <div className="flex items-center justify-between">
            <AlertTriangle className="w-7 h-7 text-orange-500" />
            <span className="text-[10px] font-bold text-orange-500 bg-orange-100 px-3 py-1 rounded uppercase tracking-widest">
              Perhatian
            </span>
          </div>
          <div className="mt-4">
            <p className="text-4xl font-bold font-headline text-neutral-900 leading-none mb-1">06</p>
            <p className="text-sm text-neutral-500 font-medium">Kecamatan Risiko Tinggi</p>
          </div>
        </div>

        {/* Dashboard Update */}
        <div className="bg-white rounded-3xl p-6 border border-neutral-100 flex flex-col justify-between min-h-[150px] shadow-sm">
          <div className="flex items-start justify-between">
            <BarChart3 className="w-7 h-7 text-blue-500" />
            <Bell className="w-5 h-5 text-neutral-300" />
          </div>
          <div className="mt-4">
            <p className="text-4xl font-bold font-headline text-blue-500 leading-none mb-1">03</p>
            <p className="text-sm text-neutral-500 font-medium">Dashboard Perlu Update</p>
          </div>
        </div>
      </div>

      {/* Timeline */}
      <NotifTimeline groups={NOTIF_GROUPS} compactLastGroup />

      {/* FAB */}
      <button className="fixed bottom-8 right-8 w-14 h-14 bg-primary hover:bg-primary-700 text-white rounded-full shadow-lg flex items-center justify-center transition-transform hover:scale-105 z-50">
        <Plus className="w-6 h-6" />
      </button>
    </div>
  );
}
