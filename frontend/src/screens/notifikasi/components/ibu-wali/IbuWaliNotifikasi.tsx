import { CalendarDays, CheckCircle2, ArrowRight, ChevronRight } from 'lucide-react';
import NotifTimeline from '../NotifTimeline';
import type { NotifGroup } from '../types';

interface IbuWaliNotifikasiProps {
  // reserved for future use
}

const NOTIF_GROUPS: NotifGroup[] = [
  {
    groupLabel: 'Hari Ini',
    items: [
      {
        id: 'iw1',
        title: 'Peringatan Stunting: Ananda Ahmad Zaelani',
        description:
          'Pasien terindikasi stunting berat berdasarkan data pengukuran terakhir. Mohon segera kunjungi RSUD untuk verifikasi lanjutan.',
        time: '09:45 AM',
        category: 'urgent',
        tags: [
          { label: 'Mendesak', color: '#ef4444', bg: '#fee2e2' },
          { label: 'Poli Anak', color: '#4b5563', bg: '#f1f5f9' },
        ],
        actionLabel: 'Lihat Instruksi',
      },
      {
        id: 'iw2',
        title: 'Verifikasi Gizi Selesai',
        description:
          'Data pemantauan gizi Ibu Hamil untuk periode Agustus telah diverifikasi oleh Bidan Desa.',
        time: '08:15 AM',
        category: 'success',
        tags: [{ label: 'Selesai', color: '#059669', bg: '#d1fae5' }],
        actionLabel: 'Unduh Laporan',
      },
    ],
  },
  {
    groupLabel: 'Kemarin',
    items: [
      {
        id: 'iw3',
        title: 'Pengingat Imunisasi: Besok Jam 09:00',
        description:
          'Jangan lupa membawa buku KIA untuk imunisasi DPT-HB-Hib 3 di Posyandu Melati.',
        time: 'Kemarin, 16:20',
        category: 'schedule',
        tags: [
          { label: 'Posyandu Melati', color: '#3b82f6', bg: '#eff6ff' },
          { label: 'Ananda Rizky', color: '#4b5563', bg: '#f1f5f9' },
        ],
      },
    ],
  },
  {
    groupLabel: 'Pekan Lalu',
    items: [
      {
        id: 'iw4',
        title: 'Jadwal Kunjungan Rumah',
        description: '',
        time: '4 Sept 2023',
        category: 'schedule',
        tags: [],
      },
      {
        id: 'iw5',
        title: 'Edukasi MPASI Baru Tersedia',
        description: '',
        time: '2 Sept 2023',
        category: 'info',
        tags: [],
      },
    ],
  },
];

const EDUKASI_ITEMS = [
  {
    id: 1,
    kategori: 'NUTRISI',
    judul: '5 Menu MPASI Sehat untuk Tumbuh Kembang',
    deskripsi: 'Pelajari cara mengolah bahan lokal menjadi asupan gizi terbaik.',
    image: 'https://images.unsplash.com/photo-1512621776951-a57141f2eefd?q=80&w=600&auto=format&fit=crop',
  },
  {
    id: 2,
    kategori: 'PSIKOLOGI',
    judul: 'Membangun Bonding yang Kuat dengan Si Kecil',
    deskripsi: 'Tips praktis untuk memperkuat ikatan emosional sejak dini.',
    image: 'https://images.unsplash.com/photo-1544126592-807ade215a0b?q=80&w=600&auto=format&fit=crop',
  },
  {
    id: 3,
    kategori: 'MONITORING',
    judul: 'Memahami Grafik Pertumbuhan Anak',
    deskripsi: 'Panduan membaca buku KIA agar lebih cepat tanggap.',
    image: 'https://images.unsplash.com/photo-1504829329064-0010c73e04e4?q=80&w=600&auto=format&fit=crop',
  },
];

export default function IbuWaliNotifikasi(_props: IbuWaliNotifikasiProps): JSX.Element {
  return (
    <div>
      {/* Hero Summary Card */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-5 mb-10">
        {/* Main highlight card */}
        <div className="lg:col-span-2 bg-gradient-to-br from-emerald-500 to-blue-500 rounded-3xl p-7 text-white shadow-sm relative overflow-hidden flex flex-col justify-between min-h-[200px]">
          <div className="absolute top-0 right-0 w-60 h-60 bg-white/5 rounded-full -translate-y-1/2 translate-x-1/4" />
          <div className="relative z-10 flex justify-between items-start mb-3">
            <span className="bg-white/20 backdrop-blur-sm text-white text-[10px] font-bold px-3 py-1.5 rounded-full uppercase tracking-widest">
              Pesan Penting
            </span>
          </div>
          <div className="relative z-10">
            <h2 className="text-xl font-bold font-headline mb-2 leading-snug">
              Rujukan Aktif: RSUD Kota — Ananda Rizky
            </h2>
            <p className="text-white/85 text-sm mb-5 leading-relaxed max-w-sm">
              Berkas rujukan telah diterima oleh faskes tingkat lanjut. Mohon segera lengkapi administrasi penunjang.
            </p>
            <button className="flex items-center gap-2 bg-white text-emerald-700 font-bold text-sm px-5 py-2.5 rounded-xl hover:bg-neutral-50 transition-colors w-fit shadow-sm">
              Lihat Detail Rujukan <ArrowRight className="w-4 h-4" />
            </button>
          </div>
        </div>

        {/* Side info cards */}
        <div className="flex flex-col gap-4">
          <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm flex items-start gap-4 flex-1">
            <div className="w-11 h-11 rounded-xl bg-emerald-50 flex items-center justify-center shrink-0">
              <CalendarDays className="w-5 h-5 text-emerald-600" />
            </div>
            <div>
              <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest mb-0.5">Jadwal Pemeriksaan</p>
              <p className="font-bold text-neutral-800 text-sm">12 Sept 2023, 09:00</p>
              <p className="text-xs text-neutral-500">Posyandu Melati</p>
            </div>
          </div>
          <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm flex items-start gap-4 flex-1">
            <div className="w-11 h-11 rounded-xl bg-emerald-50 flex items-center justify-center shrink-0">
              <CheckCircle2 className="w-5 h-5 text-emerald-600" />
            </div>
            <div>
              <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest mb-0.5">Status Monitoring</p>
              <p className="font-bold text-neutral-800 text-sm">Laporan Agustus Selesai</p>
              <p className="text-xs text-neutral-500">Hasil pemantauan gizi tersedia</p>
            </div>
          </div>
        </div>
      </div>

      {/* Notification Timeline */}
      <div className="mb-10">
        <div className="flex items-center justify-between mb-5">
          <h3 className="text-lg font-bold text-neutral-800 font-headline">Pemberitahuan Terbaru</h3>
          <button className="text-sm font-semibold text-primary hover:text-primary-700 transition-colors">
            Tandai semua dibaca
          </button>
        </div>
        <NotifTimeline groups={NOTIF_GROUPS} compactLastGroup />
      </div>

      {/* Edukasi Section */}
      <div>
        <div className="flex items-center justify-between mb-5">
          <h3 className="text-lg font-bold text-neutral-800 font-headline">Edukasi & Tips Baru</h3>
          <button className="w-8 h-8 rounded-full bg-neutral-100 flex items-center justify-center text-neutral-600 hover:bg-neutral-200 transition-colors">
            <ChevronRight className="w-4 h-4" />
          </button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
          {EDUKASI_ITEMS.map((item) => (
            <div
              key={item.id}
              className="bg-white rounded-2xl overflow-hidden border border-neutral-100 shadow-sm hover:shadow-md transition-shadow group flex flex-col"
            >
              <div className="h-40 overflow-hidden relative">
                <img
                  src={item.image}
                  alt={item.judul}
                  className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                />
                <span className="absolute bottom-2 left-2 bg-white/90 text-neutral-700 text-[10px] font-bold px-2.5 py-1 rounded-full uppercase tracking-wide backdrop-blur-sm">
                  {item.kategori}
                </span>
              </div>
              <div className="p-5 flex flex-col flex-1">
                <h4 className="font-bold text-neutral-800 text-sm mb-1.5 leading-snug font-headline">{item.judul}</h4>
                <p className="text-xs text-neutral-500 flex-1 leading-relaxed">{item.deskripsi}</p>
                <button className="mt-4 w-full py-2.5 rounded-xl border border-neutral-200 text-primary font-semibold text-sm hover:bg-neutral-50 hover:border-primary transition-colors">
                  Buka Edukasi
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
