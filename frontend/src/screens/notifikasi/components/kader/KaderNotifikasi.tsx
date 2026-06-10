import { CalendarClock, AlertTriangle, Settings, CheckCircle2 } from 'lucide-react';
import NotifTimeline from '../NotifTimeline';
import type { NotifGroup } from '../types';

// Data yang sudah diverifikasi oleh Bidan — akan kader bisa lihat statusnya
const DATA_TERVERIFIKASI = [
  {
    id: 'dv1',
    nama: 'Ny. Siti Aminah',
    inisial: 'SA',
    warnaBg: 'bg-emerald-500',
    warnaText: 'text-white',
    usia: '28 Minggu (Hamil)',
    bb: '62',
    tb: '158',
    statusGizi: 'Normal',
    verifikasiOleh: 'Bidan Sri Lestari',
    tanggal: '10 Jun 2024, 09:30',
    catatan: 'Data valid, pertambahan berat badan sesuai standar.',
  },
  {
    id: 'dv2',
    nama: 'An. Nabila Putri',
    inisial: 'NP',
    warnaBg: 'bg-amber-500',
    warnaText: 'text-white',
    usia: '18 Bln',
    bb: '8.5',
    tb: '72.1',
    statusGizi: 'Gizi Kurang',
    verifikasiOleh: 'Bidan Sri Lestari',
    tanggal: '09 Jun 2024, 14:15',
    catatan: 'Perlu intervensi PMT selama 3 bulan. Jadwalkan kunjungan rumah.',
  },
];

const NOTIF_GROUPS: NotifGroup[] = [
  {
    groupLabel: 'Hari Ini',
    items: [
      {
        id: 'k1',
        title: 'Rujukan Baru: An. Ahmad Zaelani',
        description:
          'Pasien terindikasi stunting berat. Segera lakukan verifikasi data rujukan ke RSUD setempat untuk penanganan lanjut.',
        time: '09:45 AM',
        category: 'urgent',
        tags: [
          { label: 'Mendesak', color: '#ef4444', bg: '#fee2e2' },
          { label: 'Poli Anak', color: '#4b5563', bg: '#f1f5f9' },
        ],
      },
      {
        id: 'k2',
        title: 'Verifikasi Selesai: An. Nabila Putri',
        description:
          'Data pemantauan gizi telah diverifikasi oleh Bidan Sri Lestari. Perlu intervensi PMT selama 3 bulan.',
        time: '08:15 AM',
        category: 'success',
        tags: [{ label: 'Terverifikasi', color: '#059669', bg: '#d1fae5' }],
      },
    ],
  },
  {
    groupLabel: 'Kemarin',
    items: [
      {
        id: 'k3',
        title: 'Jadwal Ulang: Kunjungan Rumah',
        description:
          'Kunjungan rumah untuk pemantauan balita RW 05 telah dijadwalkan ulang ke tanggal 15 September oleh Kader Wilayah.',
        time: 'Kemarin, 16:20',
        category: 'schedule',
        tags: [
          { label: 'Penjadwalan', color: '#3b82f6', bg: '#eff6ff' },
          { label: 'Posyandu Melati', color: '#4b5563', bg: '#f1f5f9' },
        ],
      },
      {
        id: 'k4',
        title: 'Verifikasi Selesai: Ny. Siti Aminah',
        description: 'Laporan rekapitulasi data kehamilan telah disetujui oleh bidan.',
        time: 'Kemarin, 14:05',
        category: 'report',
        tags: [{ label: 'Terverifikasi', color: '#059669', bg: '#d1fae5' }],
      },
    ],
  },
  {
    groupLabel: 'Pekan Lalu',
    items: [
      {
        id: 'k5',
        title: 'Edukasi Parenting Selesai',
        description: '',
        time: '4 Sept 2023',
        category: 'success',
        tags: [],
      },
      {
        id: 'k6',
        title: 'Stok PMT Diperbarui',
        description: '',
        time: '2 Sept 2023',
        category: 'info',
        tags: [],
      },
    ],
  },
];

export default function KaderNotifikasi(): JSX.Element {
  return (
    <div>
      {/* Summary Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-5 mb-10">
        {/* Prioritas */}
        <div className="bg-primary rounded-3xl p-6 text-white flex flex-col justify-between min-h-[160px] relative overflow-hidden shadow-sm">
          <div className="flex items-center justify-between">
            <CalendarClock className="w-8 h-8 text-white/90" />
            <span className="bg-white/20 backdrop-blur-sm text-white text-[10px] font-bold px-3 py-1 rounded uppercase tracking-widest">
              Prioritas
            </span>
          </div>
          <div className="mt-4">
            <p className="text-5xl font-bold font-headline leading-none mb-1">12</p>
            <p className="text-sm text-white/85 font-medium">Jadwal Ulang Menunggu</p>
          </div>
          <div className="absolute -bottom-6 -right-6 w-28 h-28 bg-white/10 rounded-full" />
          <div className="absolute -bottom-10 -right-10 w-40 h-40 bg-white/5 rounded-full" />
        </div>

        {/* Rujukan Mendesak */}
        <div className="bg-white rounded-3xl p-6 border border-neutral-100 border-b-4 border-b-red-500 flex flex-col justify-between min-h-[160px] shadow-sm">
          <div className="flex items-start justify-between">
            <AlertTriangle className="w-8 h-8 text-red-500" />
          </div>
          <div className="mt-4">
            <p className="text-5xl font-bold font-headline text-neutral-900 leading-none mb-1">05</p>
            <p className="text-sm text-neutral-500 font-medium">Rujukan Baru Mendesak</p>
          </div>
        </div>
      </div>

      {/* ── Data Diverifikasi Bidan ── */}
      <div className="mb-10">
        <div className="flex items-center gap-4 mb-5">
          <span className="text-xs font-bold text-neutral-400 uppercase tracking-widest whitespace-nowrap">
            Data Diverifikasi Bidan
          </span>
          <div className="h-px bg-neutral-200 flex-1" />
          <span className="text-[10px] font-bold bg-emerald-500 text-white px-2.5 py-0.5 rounded-full whitespace-nowrap">
            {DATA_TERVERIFIKASI.length} Selesai
          </span>
        </div>
        <div className="space-y-3">
          {DATA_TERVERIFIKASI.map(item => (
            <div
              key={item.id}
              className="bg-white border border-emerald-100 rounded-2xl p-5 shadow-sm flex items-start gap-4 hover:shadow-md transition-all"
            >
              {/* Avatar */}
              <div className={`w-11 h-11 rounded-full flex items-center justify-center flex-shrink-0 text-sm font-bold ${item.warnaBg} ${item.warnaText}`}>
                {item.inisial}
              </div>
              {/* Info */}
              <div className="flex-1 min-w-0">
                <div className="flex items-start justify-between gap-2">
                  <div>
                    <p className="font-bold text-neutral-800 text-sm">{item.nama}</p>
                    <p className="text-xs text-neutral-500 mt-0.5">
                      {item.usia} · BB {item.bb} kg · TB {item.tb} cm
                    </p>
                  </div>
                  <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full flex-shrink-0 ${
                    item.statusGizi === 'Stunting' ? 'bg-red-100 text-red-700' :
                    item.statusGizi === 'Gizi Kurang' ? 'bg-amber-100 text-amber-700' :
                    'bg-emerald-100 text-emerald-700'
                  }`}>
                    {item.statusGizi}
                  </span>
                </div>
                {item.catatan && (
                  <p className="text-xs text-neutral-600 mt-2 leading-relaxed italic border-l-2 border-emerald-300 pl-2">
                    "{item.catatan}"
                  </p>
                )}
                <div className="flex items-center gap-1.5 mt-2">
                  <CheckCircle2 className="w-3.5 h-3.5 text-emerald-500 flex-shrink-0" />
                  <p className="text-[11px] text-neutral-400">
                    Diverifikasi oleh <span className="font-semibold text-emerald-700">{item.verifikasiOleh}</span> · {item.tanggal}
                  </p>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Timeline */}
      <NotifTimeline groups={NOTIF_GROUPS} compactLastGroup />

      {/* FAB */}
      <button className="fixed bottom-8 right-8 w-14 h-14 bg-primary hover:bg-primary-700 text-white rounded-full shadow-lg flex items-center justify-center transition-transform hover:scale-105 z-50">
        <Settings className="w-6 h-6" />
      </button>
    </div>
  );
}
