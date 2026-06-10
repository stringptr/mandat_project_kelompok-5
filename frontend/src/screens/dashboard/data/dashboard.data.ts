export interface AktivitasPasien {
  id: string;
  nama: string;
  tindakan: string;
  waktu: string;
  status: 'Normal' | 'Risiko Stunting' | 'Gizi Kurang' | 'Stunting' | 'Underweight';
}

export interface JadwalItem {
  id: string;
  bulan: string;
  tanggal: number;
  judul: string;
  lokasi: string;
  jam: string;
}

export interface TargetWilayah {
  label: string;
  persen: number;
  color: string;
}

export interface IntervensiItem {
  id: string;
  posyandu: string;
  kecamatan: string;
  tindakan: string;
  tindakanColor: string;
  statusPasien: 'Normal' | 'Underweight' | 'Stunting';
  waktu: string;
  progres: string;
  progresPositif: boolean;
}

// ── Bidan/Kader data ────────────────────────────────────────────────────────

export const AKTIVITAS_TERBARU: AktivitasPasien[] = [
  { id: 'p1', nama: 'Aminah', tindakan: 'Pemeriksaan Kehamilan', waktu: '2 jam yang lalu', status: 'Normal' },
  { id: 'p2', nama: 'Rayyan', tindakan: 'Pemantauan Gizi Bulanan', waktu: '5 jam yang lalu', status: 'Risiko Stunting' },
  { id: 'p3', nama: 'Linda Wati', tindakan: 'Pemberian Vitamin A', waktu: 'Kemarin', status: 'Normal' },
  { id: 'p4', nama: 'Budi Santoso', tindakan: 'Pengukuran BB/TB', waktu: 'Kemarin', status: 'Gizi Kurang' },
  { id: 'p5', nama: 'Sari Dewi', tindakan: 'Imunisasi Dasar', waktu: '3 hari lalu', status: 'Normal' },
];

export const JADWAL_TERDEKAT: JadwalItem[] = [
  { id: 'j1', bulan: 'JUN', tanggal: 14, judul: 'Posyandu Balita', lokasi: 'Balai Desa Melati', jam: '08:00 WIB' },
  { id: 'j2', bulan: 'JUN', tanggal: 16, judul: 'Kunjungan Rumah', lokasi: 'RW 04', jam: '10:00 WIB' },
  { id: 'j3', bulan: 'JUN', tanggal: 20, judul: 'Imunisasi Lanjutan', lokasi: 'Puskesmas Melati', jam: '09:00 WIB' },
];

export const TARGET_WILAYAH: TargetWilayah[] = [
  { label: 'Imunisasi Dasar', persen: 85, color: 'bg-primary' },
  { label: 'Pemeriksaan K4', persen: 62, color: 'bg-blue-500' },
  { label: 'Stunting Free', persen: 94, color: 'bg-primary' },
];

// ── Ibu/Wali data ───────────────────────────────────────────────────────────

export interface RiwayatAnakItem {
  id: string;
  tanggal: string;
  jenis: string;
  hasil: string;
  status: 'Normal' | 'Perlu Perhatian';
  petugas: string;
}

export interface TumbuhKembangItem {
  bulan: string;
  bb: number;
  tb: number;
}

export const RIWAYAT_ANAK: RiwayatAnakItem[] = [
  { id: 'r1', tanggal: '10 Jun 2024', jenis: 'Penimbangan BB/TB', hasil: 'BB: 10.2 kg / TB: 78 cm', status: 'Normal', petugas: 'Bidan Sri' },
  { id: 'r2', tanggal: '14 Mei 2024', jenis: 'Imunisasi DPT-HB', hasil: 'Selesai, reaksi normal', status: 'Normal', petugas: 'dr. Ayu' },
  { id: 'r3', tanggal: '10 Apr 2024', jenis: 'Penimbangan BB/TB', hasil: 'BB: 9.8 kg / TB: 76 cm', status: 'Normal', petugas: 'Bidan Sri' },
  { id: 'r4', tanggal: '10 Mar 2024', jenis: 'Penimbangan BB/TB', hasil: 'BB: 9.3 kg / TB: 74 cm', status: 'Perlu Perhatian', petugas: 'Kader Siti' },
];

export const TUMBUH_KEMBANG: TumbuhKembangItem[] = [
  { bulan: 'Jan', bb: 8.5, tb: 70 },
  { bulan: 'Feb', bb: 8.8, tb: 71 },
  { bulan: 'Mar', bb: 9.3, tb: 74 },
  { bulan: 'Apr', bb: 9.8, tb: 76 },
  { bulan: 'Mei', bb: 10.0, tb: 77 },
  { bulan: 'Jun', bb: 10.2, tb: 78 },
];

// ── Dinkes data ─────────────────────────────────────────────────────────────

export const INTERVENSI_TERBARU: IntervensiItem[] = [
  { id: 'i1', posyandu: 'Posyandu Melati 1', kecamatan: 'Kec. Balonrejo', tindakan: 'Pemantauan', tindakanColor: 'bg-blue-100 text-blue-700', statusPasien: 'Normal', waktu: 'Hari ini, 09:34', progres: '+2.4kg', progresPositif: true },
  { id: 'i2', posyandu: 'Posyandu Cangkiran', kecamatan: 'Kec. Banyumanik', tindakan: 'Rujukan Gizi', tindakanColor: 'bg-red-100 text-red-700', statusPasien: 'Underweight', waktu: 'Kemarin, 14:15', progres: '-0.9kg', progresPositif: false },
  { id: 'i3', posyandu: 'Posyandu Anggrek 4', kecamatan: 'Kec. Tembalang', tindakan: 'Konseling Ibu', tindakanColor: 'bg-violet-100 text-violet-700', statusPasien: 'Normal', waktu: '12 Feb, 10:20', progres: 'Optimal', progresPositif: true },
  { id: 'i4', posyandu: 'Posyandu Mawar 2', kecamatan: 'Kec. Semarang Utara', tindakan: 'Imunisasi', tindakanColor: 'bg-emerald-100 text-emerald-700', statusPasien: 'Normal', waktu: '10 Feb, 08:00', progres: '+1.1kg', progresPositif: true },
];

export interface KabupatenData {
  nama: string;
  prevalensi: number; // persen
  jumlahKasus: number;
  level: 'tinggi' | 'sedang' | 'rendah';
}

export const KABUPATEN_DATA: KabupatenData[] = [
  { nama: 'Kab. Brebes', prevalensi: 32.4, jumlahKasus: 4821, level: 'tinggi' },
  { nama: 'Kab. Pemalang', prevalensi: 28.7, jumlahKasus: 3654, level: 'tinggi' },
  { nama: 'Kab. Rembang', prevalensi: 24.1, jumlahKasus: 2190, level: 'sedang' },
  { nama: 'Kab. Blora', prevalensi: 22.8, jumlahKasus: 1987, level: 'sedang' },
  { nama: 'Kab. Semarang', prevalensi: 14.5, jumlahKasus: 1340, level: 'rendah' },
  { nama: 'Kota Semarang', prevalensi: 10.2, jumlahKasus: 892, level: 'rendah' },
  { nama: 'Kab. Boyolali', prevalensi: 18.3, jumlahKasus: 1654, level: 'sedang' },
  { nama: 'Kab. Wonosobo', prevalensi: 26.5, jumlahKasus: 2876, level: 'tinggi' },
];

export const TREN_NUTRISI = [
  { bulan: 'Okt 23', nilai: 58 },
  { bulan: 'Nov 23', nilai: 62 },
  { bulan: 'Des 23', nilai: 61 },
  { bulan: 'Jan 24', nilai: 67 },
  { bulan: 'Feb 24', nilai: 72 },
  { bulan: 'Mar 24', nilai: 75 },
  { bulan: 'Apr 24', nilai: 80 },
];
