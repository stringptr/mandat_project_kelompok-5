package pemeriksaan

import "github.com/stringptr/SiGizi/backend/internal/pagination"

type GetPendingPemeriksaanRequest struct {
	Page    int `query:"page" minimum:"1" default:"1"`
	PerPage int `query:"per_page" minimum:"1" maximum:"100" default:"20"`
}

type CreatePemeriksaanRequest struct {
	IDJadwalImunisasi int32   `json:"id_jadwal_imunisasi" minimum:"1"`
	BeratBadan        float64 `json:"berat_badan" minimum:"0"`
	TinggiBadan       float64 `json:"tinggi_badan" minimum:"0"`
	LingkarKepala     float64 `json:"lingkar_kepala" minimum:"0"`
	TekananDarah      string  `json:"tekanan_darah" minLength:"1" maxLength:"20"`
	Catatan           string  `json:"catatan,omitempty" maxLength:"1000"`
}

type CreatePemeriksaanResponse struct {
	IDHasilPemeriksaan int32  `json:"id_hasil_pemeriksaan"`
	StatusStunting     string `json:"status_stunting"`
	StatusGizi         string `json:"status_gizi"`
	CreatedAt          string `json:"created_at"`
}

type DetailPemeriksaanResponse struct {
	IDHasilPemeriksaan int32              `json:"id_hasil_pemeriksaan"`
	Pasien             PasienInfo         `json:"pasien"`
	Antropometri       AntropometriData   `json:"antropometri"`
	StatusKesehatan    StatusKesehatanData `json:"status_kesehatan"`
	Catatan            string             `json:"catatan"`
}

type PasienInfo struct {
	IDPasien int32  `json:"id_pasien"`
	Nama     string `json:"nama"`
}

type AntropometriData struct {
	BeratBadan   float64 `json:"berat_badan"`
	TinggiBadan  float64 `json:"tinggi_badan"`
	LingkarKepala float64 `json:"lingkar_kepala"`
	TekananDarah  string  `json:"tekanan_darah"`
}

type StatusKesehatanData struct {
	StatusStunting string `json:"status_stunting"`
	StatusGizi     string `json:"status_gizi"`
}

type UpdatePemeriksaanRequest struct {
	BeratBadan    *float64 `json:"berat_badan,omitempty" minimum:"0"`
	TinggiBadan   *float64 `json:"tinggi_badan,omitempty" minimum:"0"`
	LingkarKepala *float64 `json:"lingkar_kepala,omitempty" minimum:"0"`
	TekananDarah  *string  `json:"tekanan_darah,omitempty" minLength:"1" maxLength:"20"`
	Catatan       *string  `json:"catatan,omitempty" maxLength:"1000"`
}

type UpdatePemeriksaanInput struct {
	IDHasilPemeriksaan int32 `path:"id" minimum:"1"`
	Body               *UpdatePemeriksaanRequest
}

type UpdatePemeriksaanResponse struct {
	IDHasilPemeriksaan int32  `json:"id_hasil_pemeriksaan"`
	StatusGiziBaru     string `json:"status_gizi_baru"`
	UpdatedAt          string `json:"updated_at"`
}

type VerifyPemeriksaanResponse struct {
	IDHasilPemeriksaan int32  `json:"id_hasil_pemeriksaan"`
	DiverifikasiOleh   int32  `json:"diverifikasi_oleh"`
	StatusVerifikasi   string `json:"status_verifikasi"`
}

type GetAllPemeriksaanRequest struct {
	Page    int    `query:"page" minimum:"1" default:"1"`
	PerPage int    `query:"per_page" minimum:"1" maximum:"100" default:"20"`
	Q       string `query:"q" maxLength:"255"`
}

type PemeriksaanListItem struct {
	IDHasilPemeriksaan int32  `json:"id_hasil_pemeriksaan"`
	NamaPasien         string `json:"nama_pasien"`
	DiinputOleh        string `json:"diinput_oleh"`
	StatusStunting     string `json:"status_stunting"`
	StatusGizi         string `json:"status_gizi"`
	StatusVerifikasi   string `json:"status_verifikasi"`
	TanggalInput       string `json:"tanggal_input"`
}

type PemeriksaanListData struct {
	Pemeriksaan []PemeriksaanListItem `json:"pemeriksaan"`
	TotalData   int                   `json:"total_data"`
	Page        int                   `json:"page"`
	PerPage     int                   `json:"per_page"`
}

type PendingPemeriksaanItem struct {
	IDHasilPemeriksaan int32  `json:"id_hasil_pemeriksaan"`
	NamaPasien         string `json:"nama_pasien"`
	DiinputOleh        string `json:"diinput_oleh"`
	TanggalInput       string `json:"tanggal_input"`
}

type PendingPemeriksaanData struct {
	PemeriksaanPending []PendingPemeriksaanItem `json:"pemeriksaan_pending"`
	Meta               pagination.Meta          `json:"meta"`
}
