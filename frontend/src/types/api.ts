export interface ApiErrorItem {
  id: string;
  location: string;
  message: string;
  value: unknown;
}

export interface ApiResponse {
  status: number;
  success: boolean;
  title: string;
  detail: string;
  errors: ApiErrorItem[];
}

export function isApiResponse(obj: unknown): obj is ApiResponse {
  if (!obj || typeof obj !== 'object') return false;
  const o = obj as Record<string, unknown>;
  return (
    typeof o.status === 'number' &&
    typeof o.success === 'boolean' &&
    Array.isArray(o.errors)
  );
}

// ── Dashboard ──────────────────────────────────────────────────────────────────

export interface DashboardStats {
  total_pasien: number;
  perlu_verifikasi: number;
  tindak_lanjut: number;
  kasus_stunting: number;
  jadwal_posyandu: number;
  total_balita?: number;
  cakupan_persentase?: number;
}

export interface DistribusiGiziItem {
  status_gizi: string;
  jumlah: number;
}

export interface DistribusiGiziResponse {
  distribusi: DistribusiGiziItem[];
}

export interface PemeriksaanPendingItem {
  id_hasil_pemeriksaan: number;
  nama_pasien: string;
  diinput_oleh: string;
  tanggal_input: string;
}

export interface PemeriksaanPendingResponse {
  pemeriksaan_pending: PemeriksaanPendingItem[];
  total_pending: number;
}

export interface BelumUkurItem {
  id_pasien: number;
  nama_pasien: string;
  nama_ibu: string;
  usia: string;
  pengukuran_terakhir: string;
  status_terakhir: string;
  alamat: string;
}

export interface BelumUkurResponse {
  pasien: BelumUkurItem[];
}

export interface KehadiranBulananItem {
  bulan: string;
  jumlah: number;
}

export interface KehadiranBulananResponse {
  tren: KehadiranBulananItem[];
}

export interface StuntingWilayahItem {
  nama_wilayah: string;
  prevalensi: number;
  jumlah_kasus: number;
  total_balita: number;
  level: string;
}

export interface StuntingWilayahResponse {
  wilayah: StuntingWilayahItem[];
}

export interface TrenStuntingItem {
  bulan: string;
  jumlah: number;
}

export interface TrenStuntingResponse {
  tren: TrenStuntingItem[];
}

export interface PublicStatsResponse {
  total_pasien: number;
  balita_dipantau: number;
  kasus_stunting: number;
  total_artikel: number;
}

export interface PublicArtikelItem {
  id_artikel: number;
  judul: string;
  ringkasan: string;
  kategori?: string;
  nama_penulis: string;
  tanggal_publish?: string;
  status_artikel: string;
}

export interface PublicArtikelResponse {
  artikel: PublicArtikelItem[];
}
