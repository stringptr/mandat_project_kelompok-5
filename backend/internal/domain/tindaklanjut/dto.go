package tindaklanjut

import "github.com/stringptr/SiGizi/backend/internal/pagination"

type GetStatusTindakLanjutRequest struct {
	Page    int `query:"page" minimum:"1" default:"1"`
	PerPage int `query:"per_page" minimum:"1" maximum:"100" default:"20"`
}

type GetPasienTindakLanjutRequest struct {
	Page    int    `query:"page" minimum:"1" default:"1"`
	PerPage int    `query:"per_page" minimum:"1" maximum:"100" default:"20"`
}

type PasienTindakLanjutItem struct {
	IDPasien           int32  `json:"id_pasien"`
	NamaPasien         string `json:"nama_pasien"`
	StatusGizi         string `json:"status_gizi"`
	StatusPasien       string `json:"status_pasien"`
	TanggalPemeriksaan string `json:"tanggal_pemeriksaan"`
}

type PasienTindakLanjutData struct {
	Pasien []PasienTindakLanjutItem `json:"pasien"`
	Meta   pagination.Meta          `json:"meta"`
}

type MonitoringTerakhir struct {
	StatusGizi    string `json:"status_gizi"`
	StatusStunting string `json:"status_stunting"`
	Catatan       string `json:"catatan"`
}

type RiwayatPemeriksaanItem struct {
	Tanggal    string  `json:"tanggal"`
	BeratBadan float64 `json:"berat_badan"`
	TinggiBadan float64 `json:"tinggi_badan"`
}

type DetailPasienTindakLanjut struct {
	IDPasien                int32                    `json:"id_pasien"`
	NamaPasien              string                   `json:"nama_pasien"`
	Usia                    string                   `json:"usia"`
	HasilMonitoringTerakhir *MonitoringTerakhir      `json:"hasil_monitoring_terakhir"`
	RiwayatPemeriksaan      []RiwayatPemeriksaanItem `json:"riwayat_pemeriksaan"`
}

type CreateTindakLanjutRequest struct {
	IDHasilPemeriksaan int32   `json:"id_hasil_pemeriksaan"`
	JenisTindakan      string  `json:"jenis_tindakan" enum:"Rujukan,Kontrol Ulang"`
	CatatanMedis       string  `json:"catatan_medis,omitempty" maxLength:"1000"`
	Rekomendasi        string  `json:"rekomendasi,omitempty" maxLength:"1000"`
	JadwalKontrol      string  `json:"jadwal_kontrol,omitempty"`
	AlasanRujukan      string  `json:"alasan_rujukan,omitempty" maxLength:"1000"`
	IDFaskes           *int32  `json:"id_faskes,omitempty"`
}

type CreateTindakLanjutResponse struct {
	IDTindakLanjut int32  `json:"id_tindak_lanjut"`
	IDRujukan      *int32 `json:"id_rujukan"`
	StatusPasien   string `json:"status_pasien"`
}

type UpdateStatusRujukanRequest struct {
	StatusRujukan string `json:"status_rujukan" enum:"Diajukan,Diproses,Diterima,Ditolak,Selesai"`
}

type UpdateStatusRujukanResponse struct {
	IDRujukan    int32  `json:"id_rujukan"`
	StatusRujukan string `json:"status_rujukan"`
}

type StatusTindakLanjutItem struct {
	IDPasien        int32  `json:"id_pasien"`
	NamaPasien      string `json:"nama_pasien"`
	StatusPasien    string `json:"status_pasien"`
	StatusRujukan   string `json:"status_rujukan"`
	TanggalRujukan  string `json:"tanggal_rujukan"`
	TanggalDeadline string `json:"tanggal_deadline"`
	StatusDeadline  string `json:"status_deadline"`
	Faskes          string `json:"faskes"`
	AlasanRujukan   string `json:"alasan_rujukan"`
	JenisTindakan   string `json:"jenis_tindakan"`
}

type StatusTindakLanjutData struct {
	Pasien []StatusTindakLanjutItem `json:"pasien"`
	Meta   pagination.Meta          `json:"meta"`
}

type LaporanTindakLanjutItem struct {
	Wilayah            string `json:"wilayah"`
	JumlahPasienDirujuk int32  `json:"jumlah_pasien_dirujuk"`
	JumlahPasienDiterima int32 `json:"jumlah_pasien_diterima"`
	JumlahPasienDiproses int32 `json:"jumlah_pasien_diproses"`
}

type LaporanTindakLanjutData struct {
	Laporan   []LaporanTindakLanjutItem `json:"laporan"`
	TotalData int                        `json:"total_data"`
}

type DetailTindakLanjutPasien struct {
	IDTindakLanjut int32   `json:"id_tindak_lanjut"`
	StatusPasien   string  `json:"status_pasien"`
	CatatanMedis   *string `json:"catatan_medis"`
	Rekomendasi    *string `json:"rekomendasi"`
	JadwalKontrol  *string `json:"jadwal_kontrol"`
	StatusRujukan  *string `json:"status_rujukan"`
	NamaFaskes     *string `json:"nama_faskes"`
}
