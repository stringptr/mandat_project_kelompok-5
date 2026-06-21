// ── Pasien ────────────────────────────────────────────────────────────────────

export interface PasienListItem {
  id_pasien: number;
  nama: string;
  nik: string;
  jenis_kelamin: string;
  umur: string;
  nama_posyandu: string;
  jenis_pasien: string;
  status_kehamilan?: string | null;
}

export interface PaginationMeta {
  current_page: number;
  per_page: number;
  total: number;
  last_page: number;
}

export interface PasienListData {
  pasien: PasienListItem[];
  meta: PaginationMeta;
}

export interface IbuHamilData {
  id_ibu_hamil: number;
  hamil_ke: number;
  bulan_mulai_hamil: string;
  hpht: string;
  status_kehamilan: string;
}

export interface AnakData {
  nama_anak: string;
  berat_lahir: number;
  panjang_lahir: number;
  hubungan_dengan_wali: string;
  nama_wali: string;
}

export interface PasienDetail {
  id_pasien: number;
  nama: string;
  nik: string;
  email: string;
  no_hp: string;
  jenis_kelamin: string;
  tanggal_lahir: string;
  id_lokasi: number;
  nama_posyandu: string;
  id_posyandu: number;
  jenis_pasien: string;
  data_ibu_hamil?: IbuHamilData | null;
  data_anak?: AnakData | null;
  created_at: string;
  updated_at: string;
}

export interface PasienSearchResult {
  id_pasien: number;
  nama: string;
  nik: string;
  jenis_kelamin: string;
  umur: string;
  nama_posyandu: string;
  jenis_pasien: string;
  status_kehamilan?: string | null;
}

export interface PasienRiwayatItem {
  tanggal: string;
  berat_badan: number;
  tinggi_badan: number;
  status_gizi: string;
  catatan?: string;
  petugas: string;
}

// ── Pemeriksaan ────────────────────────────────────────────────────────────────

export interface CreatePemeriksaanRequest {
  id_jadwal_imunisasi: number;
  berat_badan: number;
  tinggi_badan: number;
  lingkar_kepala: number;
  tekanan_darah: string;
  catatan?: string;
}

export interface CreatePemeriksaanResponse {
  id_hasil_pemeriksaan: number;
  status_stunting: string;
  status_gizi: string;
  created_at: string;
}

export interface UpdatePemeriksaanRequest {
  berat_badan?: number;
  tinggi_badan?: number;
  lingkar_kepala?: number;
  tekanan_darah?: string;
  catatan?: string;
}

export interface UpdatePemeriksaanResponse {
  id_hasil_pemeriksaan: number;
  status_gizi_baru: string;
  updated_at: string;
}

export interface PasienInfo {
  id_pasien: number;
  nama: string;
}

export interface AntropometriData {
  berat_badan: number;
  tinggi_badan: number;
  lingkar_kepala: number;
  tekanan_darah: string;
}

export interface StatusKesehatanData {
  status_stunting: string;
  status_gizi: string;
}

export interface DetailPemeriksaanResponse {
  id_hasil_pemeriksaan: number;
  pasien: PasienInfo;
  antropometri: AntropometriData;
  status_kesehatan: StatusKesehatanData;
  catatan: string;
}

export interface PemeriksaanPendingItem {
  id_hasil_pemeriksaan: number;
  nama_pasien: string;
  diinput_oleh: string;
  tanggal_input: string;
}

export interface PendingPemeriksaanData {
  pemeriksaan_pending: PemeriksaanPendingItem[];
  meta: PaginationMeta;
}

export interface VerifyPemeriksaanResponse {
  id_hasil_pemeriksaan: number;
  diverifikasi_oleh: number;
  status_verifikasi: string;
}

// ── Imunisasi ──────────────────────────────────────────────────────────────────

export interface ImunisasiListItem {
  id_imunisasi: number;
  nama_pasien: string;
  nama_vaksin: string;
  tanggal_jadwal: string;
  status_imunisasi: string;
}

export interface ImunisasiListData {
  jadwal: ImunisasiListItem[];
  meta: PaginationMeta;
}

export interface CreateImunisasiRequest {
  id_pasien: number;
  nama_vaksin: string;
  tanggal_jadwal: string;
}

export interface CreateImunisasiResponse {
  id_imunisasi: number;
  status_imunisasi: string;
}

export interface ImunisasiDetail {
  id_imunisasi: number;
  id_pasien: number;
  nama_pasien: string;
  nama_vaksin: string;
  tanggal_jadwal: string;
  tanggal_realisasi?: string | null;
  status_imunisasi: string;
}

export interface UpdateImunisasiRequest {
  id_pasien?: number;
  nama_vaksin?: string;
  tanggal_jadwal?: string;
}

export interface UpdateImunisasiResponse {
  id_imunisasi: number;
  updated_at: string;
}

export interface RealisasiRequest {
  tanggal_realisasi: string;
}

export interface RealisasiResponse {
  id_imunisasi: number;
  status_imunisasi: string;
  tanggal_realisasi: string;
}

export interface RiwayatImunisasiItem {
  id_imunisasi: number;
  nama_vaksin: string;
  tanggal_jadwal: string;
  tanggal_realisasi?: string | null;
  status_imunisasi: string;
}

export interface RiwayatImunisasiResponse {
  id_pasien: number;
  riwayat_imunisasi: RiwayatImunisasiItem[];
}

export interface StatistikImunisasi {
  total_target_imunisasi: number;
  total_terealisasi: number;
  cakupan_persentase: number;
  vaksin_terbanyak: string;
}

// ── Artikel ────────────────────────────────────────────────────────────────────

export interface ArtikelListItem {
  id_artikel: number;
  judul: string;
  kategori: string;
  ringkasan: string;
  nama_penulis: string;
  tanggal_publish: string;
}

export interface ArtikelListData {
  artikel: ArtikelListItem[];
  meta: PaginationMeta;
}

export interface CreateArtikelRequest {
  judul: string;
  isi_artikel: string;
  kategori?: string;
}

export interface CreateArtikelResponse {
  id_artikel: number;
  status_artikel: string;
}

export interface ArtikelDetail {
  id_artikel: number;
  judul: string;
  isi_artikel: string;
  kategori: string;
  nama_penulis: string;
  nama_verifikator?: string | null;
  tanggal_publish?: string | null;
  created_at: string;
  updated_at: string;
}

export interface UpdateArtikelRequest {
  judul?: string;
  isi_artikel?: string;
  kategori?: string;
}

export interface ReviewArtikelRequest {
  aksi: 'setujui' | 'tolak';
  catatan_review?: string;
}

export interface ReviewArtikelResponse {
  id_artikel: number;
  status_artikel: string;
  tanggal_publish?: string | null;
}

export interface ArtikelPendingItem {
  id_artikel: number;
  judul: string;
  nama_penulis: string;
  created_at: string;
  status_artikel: string;
}

export interface ArtikelPendingData {
  artikel: ArtikelPendingItem[];
  meta: PaginationMeta;
}

// ── Tindak Lanjut ──────────────────────────────────────────────────────────────

export interface PasienTindakLanjutItem {
  id_pasien: number;
  nama_pasien: string;
  status_gizi: string;
  status_pasien: string;
  tanggal_pemeriksaan: string;
}

export interface PasienTindakLanjutData {
  pasien: PasienTindakLanjutItem[];
  meta: PaginationMeta;
}

export interface CreateTindakLanjutRequest {
  id_hasil_pemeriksaan: number;
  jenis_tindakan: 'Rujukan' | 'Kontrol Ulang';
  catatan_medis?: string;
  rekomendasi?: string;
  jadwal_kontrol?: string;
  alasan_rujukan?: string;
  id_faskes?: number;
}

export interface CreateTindakLanjutResponse {
  id_tindak_lanjut: number;
  id_rujukan?: number | null;
  status_pasien: string;
}

export interface UpdateStatusRujukanRequest {
  status_rujukan: 'Diajukan' | 'Diproses' | 'Diterima' | 'Ditolak' | 'Selesai';
}

export interface UpdateStatusRujukanResponse {
  id_rujukan: number;
  status_rujukan: string;
}

export interface StatusTindakLanjutItem {
  id_pasien: number;
  nama_pasien: string;
  status_pasien: string;
  status_rujukan: string;
  tanggal_rujukan: string;
}

export interface StatusTindakLanjutData {
  pasien: StatusTindakLanjutItem[];
  meta: PaginationMeta;
}

export interface MonitoringTerakhir {
  status_gizi: string;
  status_stunting: string;
  catatan: string;
}

export interface RiwayatPemeriksaanItem {
  tanggal: string;
  berat_badan: number;
  tinggi_badan: number;
}

export interface DetailPasienTindakLanjut {
  id_pasien: number;
  nama_pasien: string;
  usia: string;
  hasil_monitoring_terakhir?: MonitoringTerakhir | null;
  riwayat_pemeriksaan: RiwayatPemeriksaanItem[];
}

export interface LaporanTindakLanjutItem {
  wilayah: string;
  jumlah_pasien_dirujuk: number;
  jumlah_pasien_diterima: number;
  jumlah_pasien_diproses: number;
}

export interface LaporanTindakLanjutData {
  laporan: LaporanTindakLanjutItem[];
  total_data: number;
}

export interface DetailTindakLanjutPasien {
  id_tindak_lanjut: number;
  status_pasien: string;
  catatan_medis?: string | null;
  rekomendasi?: string | null;
  jadwal_kontrol?: string | null;
  status_rujukan?: string | null;
  nama_faskes?: string | null;
}

// ── Lokasi ─────────────────────────────────────────────────────────────────────

export interface LokasiItem {
  id_lokasi: number;
  nama_lokasi: string;
  tipe_lokasi: string;
  bagian_dari?: number | null;
}

// ── Notifikasi ─────────────────────────────────────────────────────────────────

export interface BackendNotifikasi {
  id: string;
  title: string;
  message: string;
  is_read: boolean;
  created_at: string;
  notif_type?: string;
}

export interface NotifikasiResponse {
  notifikasi: BackendNotifikasi[];
  meta: PaginationMeta;
}

export interface BidanDashboard {
  statistik: {
    risiko_stunting: number;
    jadwal_ulang: number;
    rujukan_mendesak: number;
  };
  risiko_stunting: { nama_pasien: string; status_gizi: string }[];
  jadwal_monitoring: { nama_pasien: string; tanggal: string }[];
  rujukan_mendesak: { nama_pasien: string; alasan: string }[];
  laporan_bulanan: { bulan: string; jumlah: number }[];
}

export interface StatistikNotifikasi {
  jadwal_ulang: number;
  rujukan_mendesak: number;
  risiko_stunting: number;
  notifikasi_belum_dibaca: number;
}

export interface AktivitasItem {
  id: string;
  title: string;
  message: string;
  created_at: string;
}

export interface AktivitasResponse {
  hari_ini: AktivitasItem[];
  kemarin: AktivitasItem[];
}
