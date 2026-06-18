import type { NotifItem } from '../components/types';
import type { Role } from '../../../App';

// Preview items untuk header bell dropdown — 4 item terbaru per role
export const NOTIF_PREVIEW: Record<Role, NotifItem[]> = {
  'Kader Posyandu': [
    {
      id: 'k1',
      title: 'Rujukan Baru: An. Ahmad Zaelani',
      description: 'Pasien terindikasi stunting berat. Segera lakukan verifikasi data rujukan.',
      time: '09:45 AM',
      category: 'urgent',
      tags: [{ label: 'Mendesak', color: '#ef4444', bg: '#fee2e2' }],
    },
    {
      id: 'k2',
      title: 'Verifikasi Selesai',
      description: 'Data pemantauan gizi Ibu Hamil periode Agustus telah diverifikasi oleh Bidan.',
      time: '08:15 AM',
      category: 'success',
      tags: [{ label: 'Selesai', color: '#059669', bg: '#d1fae5' }],
    },
    {
      id: 'k3',
      title: 'Jadwal Ulang: Kunjungan Rumah',
      description: 'Kunjungan rumah balita RW 05 dijadwalkan ulang ke 15 September.',
      time: 'Kemarin, 16:20',
      category: 'schedule',
      tags: [{ label: 'Penjadwalan', color: '#3b82f6', bg: '#eff6ff' }],
    },
    {
      id: 'k4',
      title: 'Laporan Bulanan Disetujui',
      description: 'Rekapitulasi gizi bulanan tingkat kelurahan telah disetujui.',
      time: 'Kemarin, 14:05',
      category: 'report',
      tags: [],
    },
  ],

  'Bidan': [
    {
      id: 'b1',
      title: 'Risiko Stunting Berat: An. Ahmad Zaelani',
      description: 'Status gizi menurun drastis dalam 2 bulan. Segera buat rujukan atau jadwal kunjungan rumah.',
      time: '09:45 AM',
      category: 'urgent',
      tags: [{ label: 'Mendesak', color: '#ef4444', bg: '#fee2e2' }],
      actionLabel: 'Lakukan Verifikasi',
    },
    {
      id: 'b2',
      title: 'Hasil Lab Keluar: Ny. Siti Aminah',
      description: 'Hb: 9.5 g/dL (Anemia Ringan). Perlu tindak lanjut konseling gizi.',
      time: '08:15 AM',
      category: 'lab',
      tags: [{ label: 'Lab Hasil', color: '#059669', bg: '#d1fae5' }],
      actionLabel: 'Lihat Detail',
    },
    {
      id: 'b3',
      title: 'Gagal Timbang: Wilayah RW 05',
      description: '15 Anak tidak hadir penimbangan Posyandu Melati. Perlu sweeping.',
      time: 'Kemarin, 16:20',
      category: 'schedule',
      tags: [{ label: 'Penjadwalan', color: '#3b82f6', bg: '#eff6ff' }],
    },
    {
      id: 'b4',
      title: 'Laporan Bulanan Dinkes Ready',
      description: 'Rekapitulasi data gizi September 2023 siap ditinjau sebelum dikirim ke Dinkes.',
      time: 'Kemarin, 14:05',
      category: 'report',
      tags: [{ label: 'Laporan', color: '#6b7280', bg: '#f1f5f9' }],
      actionLabel: 'Review Laporan',
    },
  ],

  'Dinas Kesehatan': [
    {
      id: 'd1',
      title: 'Lonjakan Stunting: Kec. Banyumanik',
      description: 'Kenaikan 12% kasus stunting baru dalam 30 hari. Perlu intervensi lintas sektor.',
      time: '10:00 AM',
      category: 'urgent',
      tags: [{ label: 'Mendesak', color: '#ef4444', bg: '#fee2e2' }],
      actionLabel: 'Lihat Analitik',
    },
    {
      id: 'd2',
      title: 'Laporan Bulanan Puskesmas Siap',
      description: '14 Puskesmas mengunggah laporan gizi September. Menunggu review Dinkes.',
      time: '08:30 AM',
      category: 'report',
      tags: [{ label: 'Laporan', color: '#6b7280', bg: '#f1f5f9' }],
      actionLabel: 'Review Laporan',
    },
    {
      id: 'd3',
      title: 'Stok PMT Wilayah Utara Kritis',
      description: '3 Puskesmas stok PMT < 10%. Perlu distribusi darurat dalam 48 jam.',
      time: 'Kemarin, 15:45',
      category: 'urgent',
      tags: [{ label: 'Stok PMT', color: '#ef4444', bg: '#fee2e2' }],
      actionLabel: 'Koordinasi Distribusi',
    },
    {
      id: 'd4',
      title: 'Capaian Imunisasi Agustus',
      description: 'Coverage imunisasi dasar lengkap: 87.3% dari target.',
      time: 'Kemarin, 11:20',
      category: 'success',
      tags: [{ label: 'Imunisasi', color: '#059669', bg: '#d1fae5' }],
    },
  ],

  'Ibu/Wali': [
    {
      id: 'iw1',
      title: 'Peringatan Stunting: Ananda Ahmad',
      description: 'Terindikasi stunting berat. Mohon segera kunjungi RSUD untuk verifikasi.',
      time: '09:45 AM',
      category: 'urgent',
      tags: [{ label: 'Mendesak', color: '#ef4444', bg: '#fee2e2' }],
      actionLabel: 'Lihat Instruksi',
    },
    {
      id: 'iw2',
      title: 'Verifikasi Gizi Selesai',
      description: 'Data pemantauan gizi periode Agustus telah diverifikasi oleh Bidan Desa.',
      time: '08:15 AM',
      category: 'success',
      tags: [{ label: 'Selesai', color: '#059669', bg: '#d1fae5' }],
    },
    {
      id: 'iw3',
      title: 'Pengingat Imunisasi Besok',
      description: 'Imunisasi DPT-HB-Hib 3 di Posyandu Melati. Bawa buku KIA.',
      time: 'Kemarin, 16:20',
      category: 'schedule',
      tags: [{ label: 'Posyandu Melati', color: '#3b82f6', bg: '#eff6ff' }],
    },
    {
      id: 'iw4',
      title: 'Edukasi MPASI Baru',
      description: 'Artikel edukasi MPASI baru telah tersedia untuk Anda.',
      time: '2 Sept 2023',
      category: 'info',
      tags: [],
    },
  ],
};

// Warna & icon per kategori untuk preview
export const CATEGORY_DOT: Record<string, string> = {
  urgent: 'bg-red-500',
  success: 'bg-emerald-500',
  schedule: 'bg-blue-500',
  lab: 'bg-emerald-500',
  report: 'bg-neutral-400',
  info: 'bg-blue-400',
  warning: 'bg-amber-500',
};
