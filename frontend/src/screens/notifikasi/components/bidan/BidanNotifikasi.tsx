import { useState } from 'react';
import { CalendarClock, AlertTriangle, Bell, Plus } from 'lucide-react';
import NotifTimeline from '../NotifTimeline';
import type { NotifGroup } from '../types';
import { ModalVerifikasiBidan } from '../../../../components/verifikasi/ModalVerifikasiBidan';
import type { VerifikasiTarget } from '../../../../components/verifikasi/VerifikasiPanel';

// Data verifikasi pending dari Kader Posyandu — dalam real app ini dari API
const PENDING_VERIFIKASI: VerifikasiTarget[] = [
  {
    id: 'bv1',
    nama: 'Ahmad Zaelani',
    inisial: 'AZ',
    warnaBg: 'bg-red-500',
    warnaText: 'text-white',
    usia: '36 Bln',
    bb: '10.2',
    tb: '78.5',
    petugas: 'Siti Aminah · Posyandu Melati 01',
    statusGizi: 'Stunting',
  },
  {
    id: 'bv2',
    nama: 'Rina Pertiwi',
    inisial: 'RP',
    warnaBg: 'bg-amber-500',
    warnaText: 'text-white',
    usia: '18 Bln',
    bb: '7.8',
    tb: '70.0',
    petugas: 'Dewi Rahayu · Posyandu Anggrek 02',
    statusGizi: 'Gizi Kurang',
  },
];

const NOTIF_GROUPS: NotifGroup[] = [
  {
    groupLabel: 'Hari Ini',
    items: [
      {
        id: 'b1',
        title: 'Risiko Stunting Berat: An. Ahmad Zaelani (NIK: 3324...)',
        description:
          'Status gizi menurun drastis dalam 2 bulan. Segera buat rujukan faskes atau jadwal kunjungan rumah.',
        time: '09:45 AM',
        category: 'urgent',
        tags: [
          { label: 'Mendesak', color: '#ef4444', bg: '#fee2e2' },
          { label: 'Poli Anak', color: '#4b5563', bg: '#f1f5f9' },
        ],
      },
      {
        id: 'b2',
        title: 'Hasil Lab Keluar: Ny. Siti Aminah',
        description:
          'Hasil pemeriksaan Hb terbaru telah tersedia. Kadar Hb: 9.5 g/dL (Anemia Ringan). Perlu tindak lanjut konseling gizi.',
        time: '08:15 AM',
        category: 'lab',
        tags: [{ label: 'Lab Hasil', color: '#059669', bg: '#d1fae5' }],
      },
    ],
  },
  {
    groupLabel: 'Kemarin',
    items: [
      {
        id: 'b3',
        title: 'Gagal Timbang: Wilayah RW 05',
        description:
          '15 Anak tidak hadir pada penimbangan Posyandu Melati kemarin. Perlu koordinasi dengan Kader Wilayah untuk sweeping.',
        time: 'Kemarin, 16:20',
        category: 'schedule',
        tags: [
          { label: 'Penjadwalan', color: '#3b82f6', bg: '#eff6ff' },
          { label: 'Posyandu Melati', color: '#4b5563', bg: '#f1f5f9' },
        ],
      },
      {
        id: 'b4',
        title: 'Laporan Bulanan Dinkes Ready',
        description:
          'Rekapitulasi data gizi periode September 2023 telah selesai digenerate. Silakan tinjau sebelum dikirim ke Dinkes.',
        time: 'Kemarin, 14:05',
        category: 'report',
        tags: [{ label: 'Laporan', color: '#6b7280', bg: '#f1f5f9' }],
      },
    ],
  },
  {
    groupLabel: 'Pekan Lalu',
    items: [
      {
        id: 'b5',
        title: 'Vaksinasi Polio: Wilayah Selesai',
        description: '',
        time: '4 Sept 2023',
        category: 'success',
        tags: [],
      },
      {
        id: 'b6',
        title: 'Data EPPGBM Diperbarui',
        description: '',
        time: '2 Sept 2023',
        category: 'info',
        tags: [],
      },
    ],
  },
];

export default function BidanNotifikasi(): JSX.Element {
  const [verifikasiTarget, setVerifikasiTarget] = useState<VerifikasiTarget | null>(null);
  const [verifiedIds, setVerifiedIds] = useState<string[]>([]);

  const pendingList = PENDING_VERIFIKASI.filter(v => !verifiedIds.includes(v.id));
  const pendingCount = pendingList.length;

  const handleSetuju = (id: string, _catatan: string) => {
    setVerifiedIds(prev => [...prev, id]);
    setVerifikasiTarget(null);
  };

  const handleTolak = (_id: string, _catatan: string) => {
    setVerifikasiTarget(null);
  };

  return (
    <div>
      {/* Summary Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-5 mb-10">
        {/* Jadwal Kontrol */}
        <div className="bg-primary rounded-3xl p-6 text-white flex flex-col justify-between min-h-[150px] relative overflow-hidden shadow-sm">
          <div className="flex items-center justify-between">
            <CalendarClock className="w-7 h-7 text-white/90" />
            <span className="bg-white/20 backdrop-blur-sm text-white text-[10px] font-bold px-3 py-1 rounded uppercase tracking-widest">
              Prioritas
            </span>
          </div>
          <div className="mt-4">
            <p className="text-4xl font-bold font-headline leading-none mb-1">12</p>
            <p className="text-sm text-white/85 font-medium">Jadwal Kontrol Menunggu</p>
          </div>
          <div className="absolute -bottom-4 -right-4 w-24 h-24 bg-white/10 rounded-full" />
        </div>

        {/* Anak Berisiko */}
        <div className="bg-orange-50 rounded-3xl p-6 border border-orange-100 flex flex-col justify-between min-h-[150px] shadow-sm">
          <div className="flex items-center justify-between">
            <AlertTriangle className="w-7 h-7 text-orange-500" />
            <span className="text-[10px] font-bold text-orange-500 bg-orange-100 px-3 py-1 rounded uppercase tracking-widest">
              Perhatian
            </span>
          </div>
          <div className="mt-4">
            <p className="text-4xl font-bold font-headline text-neutral-900 leading-none mb-1">08</p>
            <p className="text-sm text-neutral-500 font-medium">Anak Berisiko Stunting</p>
          </div>
        </div>

        {/* Rujukan Mendesak */}
        <div className="bg-white rounded-3xl p-6 border border-neutral-100 flex flex-col justify-between min-h-[150px] shadow-sm">
          <div className="flex items-start justify-between">
            <Bell className="w-7 h-7 text-red-500" />
            <span className="w-2.5 h-2.5 bg-red-500 rounded-full mt-1"></span>
          </div>
          <div className="mt-4">
            <p className="text-4xl font-bold font-headline text-red-500 leading-none mb-1">05</p>
            <p className="text-sm font-semibold text-red-500">Rujukan Baru Mendesak</p>
          </div>
        </div>
      </div>

      {/* ── Antrian Verifikasi dari Kader ── */}
      {pendingCount > 0 && (
        <div className="mb-10">
          <div className="flex items-center gap-4 mb-5">
            <span className="text-xs font-bold text-neutral-400 uppercase tracking-widest whitespace-nowrap">
              Antrian Verifikasi Data Kader
            </span>
            <div className="h-px bg-neutral-200 flex-1" />
            <span className="text-[10px] font-bold bg-amber-500 text-white px-2.5 py-0.5 rounded-full whitespace-nowrap">
              {pendingCount} Pending
            </span>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {pendingList.map(item => (
              <div
                key={item.id}
                className="bg-white border border-amber-200 rounded-2xl p-5 shadow-sm flex items-start gap-4 hover:shadow-md transition-all"
              >
                <div className={`w-11 h-11 rounded-full flex items-center justify-center flex-shrink-0 text-sm font-bold ${item.warnaBg} ${item.warnaText}`}>
                  {item.inisial}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-bold text-neutral-800 text-sm">{item.nama}</p>
                  <p className="text-xs text-neutral-500 mt-0.5">
                    {item.usia} · BB {item.bb} kg · TB {item.tb} cm
                  </p>
                  <p className="text-[11px] text-neutral-400 mt-1">Oleh: {item.petugas}</p>
                  <span className={`inline-block mt-2 text-[10px] font-bold px-2 py-0.5 rounded-full ${item.statusGizi === 'Stunting' ? 'bg-red-100 text-red-700' :
                      item.statusGizi === 'Gizi Kurang' ? 'bg-amber-100 text-amber-700' :
                        'bg-emerald-100 text-emerald-700'
                    }`}>
                    {item.statusGizi}
                  </span>
                </div>
                <button
                  onClick={() => setVerifikasiTarget(item)}
                  className="flex-shrink-0 bg-primary hover:bg-primary-600 text-white text-xs font-semibold px-3 py-2 rounded-xl transition-colors"
                >
                  Verifikasi
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Timeline */}
      <NotifTimeline groups={NOTIF_GROUPS} compactLastGroup />

      {/* Edukasi Banner */}
      <div className="mt-10 bg-primary rounded-3xl p-8 text-white relative overflow-hidden">
        <div className="relative z-10 max-w-md">
          <h3 className="text-xl font-bold font-headline mb-2">Meningkatkan Kesadaran Gizi Keluarga</h3>
          <p className="text-sm text-white/85 mb-5 leading-relaxed">
            Pelajari teknik komunikasi efektif untuk memberikan edukasi MP-ASI kepada ibu muda di wilayah binaan Anda.
          </p>
          <button className="flex items-center gap-2 bg-white text-primary font-semibold text-sm px-5 py-2.5 rounded-xl hover:bg-neutral-50 transition-colors">
            Buka Materi Edukasi
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M18 13v6a2 2 0 01-2 2H5a2 2 0 01-2-2V8a2 2 0 012-2h6M15 3h6m0 0v6m0-6L10 14" />
            </svg>
          </button>
        </div>
        <div className="absolute right-0 top-0 bottom-0 w-64 opacity-20 overflow-hidden rounded-r-3xl">
          <img
            src="https://images.unsplash.com/photo-1576091160550-2173dba999ef?q=80&w=600&auto=format&fit=crop"
            alt=""
            className="w-full h-full object-cover"
          />
        </div>
        <div className="absolute right-0 top-0 bottom-0 w-64 bg-gradient-to-r from-primary to-transparent rounded-r-3xl" />
      </div>

      {/* FAB */}
      <button className="fixed bottom-8 right-8 w-14 h-14 bg-primary hover:bg-primary-700 text-white rounded-full shadow-lg flex items-center justify-center transition-transform hover:scale-105 z-50">
        <Plus className="w-6 h-6" />
      </button>

      {/* Modal Verifikasi */}
      {verifikasiTarget && (
        <ModalVerifikasiBidan
          target={verifikasiTarget}
          onClose={() => setVerifikasiTarget(null)}
          onSetuju={handleSetuju}
          onTolak={handleTolak}
        />
      )}
    </div>
  );
}
