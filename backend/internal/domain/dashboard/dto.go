package dashboard

import "github.com/stringptr/SiGizi/backend/internal/pagination"

type DashboardStatsResponse struct {
	TotalPasien       int32   `json:"total_pasien"`
	PerluVerifikasi   int32   `json:"perlu_verifikasi"`
	TindakLanjut      int32   `json:"tindak_lanjut"`
	KasusStunting     int32   `json:"kasus_stunting"`
	JadwalPosyandu    int32   `json:"jadwal_posyandu"`
	TotalBalita       int32   `json:"total_balita,omitempty"`
	CakupanPersentase float64 `json:"cakupan_persentase,omitempty"`
}

type DistribusiGiziItem struct {
	StatusGizi string `json:"status_gizi"`
	Jumlah     int32  `json:"jumlah"`
}

type DistribusiGiziResponse struct {
	Distribusi []DistribusiGiziItem `json:"distribusi"`
}

type TrenItem struct {
	Bulan  string `json:"bulan"`
	Jumlah int32  `json:"jumlah"`
}

type TrenResponse struct {
	Tren []TrenItem `json:"tren"`
}

type StuntingWilayahItem struct {
	NamaWilayah string  `json:"nama_wilayah"`
	Prevalensi  float64 `json:"prevalensi"`
	JumlahKasus int32   `json:"jumlah_kasus"`
	TotalBalita int32   `json:"total_balita"`
	Level       string  `json:"level"`
}

type StuntingWilayahResponse struct {
	Wilayah []StuntingWilayahItem `json:"wilayah"`
}

type PublicStatsResponse struct {
	TotalPasien    int32 `json:"total_pasien"`
	BalitaDipantau int32 `json:"balita_dipantau"`
	KasusStunting  int32 `json:"kasus_stunting"`
	TotalArtikel   int32 `json:"total_artikel"`
}

type RiwayatInput struct {
	IDPasien int32 `path:"id" minimum:"1"`
}

type RiwayatItem struct {
	Tanggal     string  `json:"tanggal"`
	BeratBadan  float64 `json:"berat_badan"`
	TinggiBadan float64 `json:"tinggi_badan"`
	StatusGizi  string  `json:"status_gizi"`
	Catatan     *string `json:"catatan,omitempty"`
	Petugas     string  `json:"petugas"`
}

type RiwayatResponse struct {
	Riwayat []RiwayatItem `json:"riwayat"`
}

type TumbuhKembangInput struct {
	IDPasien int32 `path:"id" minimum:"1"`
}

type TumbuhKembangItem struct {
	Bulan       string  `json:"bulan"`
	BeratBadan  float64 `json:"berat_badan"`
	TinggiBadan float64 `json:"tinggi_badan"`
}

type TumbuhKembangResponse struct {
	Data []TumbuhKembangItem `json:"data"`
}

type JadwalTerdekatItem struct {
	ID            int32  `json:"id"`
	NamaVaksin    string `json:"nama_vaksin"`
	TanggalJadwal string `json:"tanggal_jadwal"`
	NamaPasien    string `json:"nama_pasien"`
}

type JadwalTerdekatResponse struct {
	Jadwal []JadwalTerdekatItem `json:"jadwal"`
}

type PemeriksaanItem struct {
	IDHasilPemeriksaan int32   `json:"id_hasil_pemeriksaan"`
	IDJadwalImunisasi  int32   `json:"id_jadwal_imunisasi"`
	NamaVaksin         string  `json:"nama_vaksin"`
	NamaPasien         string  `json:"nama_pasien"`
	BeratBadan         float64 `json:"berat_badan"`
	TinggiBadan        float64 `json:"tinggi_badan"`
	LingkarKepala      float64 `json:"lingkar_kepala"`
	TekananDarah       string  `json:"tekanan_darah"`
	StatusStunting     string  `json:"status_stunting"`
	StatusGizi         string  `json:"status_gizi"`
	Catatan            *string `json:"catatan,omitempty"`
	Tanggal            string  `json:"tanggal"`
	Petugas            string  `json:"petugas"`
}

type PemeriksaanListResponse struct {
	Pemeriksaan []PemeriksaanItem `json:"pemeriksaan"`
	Meta        pagination.Meta   `json:"meta"`
}

type GetAllPemeriksaanRequest struct {
	Page       int   `query:"page" minimum:"1" default:"1"`
	PerPage    int   `query:"per_page" minimum:"1" maximum:"50" default:"20"`
	IDBidan    int32 `query:"id_bidan" minimum:"0" required:"false"`
	IDPosyandu int32 `query:"id_posyandu" minimum:"0" required:"false"`
	IDKader    int32 `query:"id_kader" minimum:"0" required:"false"`
}

type IbuHamilStatsResponse struct {
	TotalIbuHamil int32 `json:"total_ibu_hamil"`
	Trimester1    int32 `json:"trimester_1"`
	Trimester2    int32 `json:"trimester_2"`
	Trimester3    int32 `json:"trimester_3"`
	Melahirkan    int32 `json:"melahirkan"`
	Nifas         int32 `json:"nifas"`
	Keguguran     int32 `json:"keguguran"`
}

type IbuHamilWilayahItem struct {
	NamaWilayah   string `json:"nama_wilayah"`
	TotalIbuHamil int32  `json:"total_ibu_hamil"`
	Trimester1    int32  `json:"trimester_1"`
	Trimester2    int32  `json:"trimester_2"`
	Trimester3    int32  `json:"trimester_3"`
	Melahirkan    int32  `json:"melahirkan"`
	Nifas         int32  `json:"nifas"`
	Keguguran     int32  `json:"keguguran"`
}

type IbuHamilWilayahResponse struct {
	Wilayah []IbuHamilWilayahItem `json:"wilayah"`
}
