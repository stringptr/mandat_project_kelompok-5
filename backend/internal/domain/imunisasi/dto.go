package imunisasi

type GetAllImunisasiRequest struct {
	Page    int    `query:"page" minimum:"1" default:"1"`
	PerPage int    `query:"per_page" minimum:"1" maximum:"100" default:"20"`
	Q       string `query:"q" maxLength:"255"`
}

type CreateImunisasiRequest struct {
	IDPasien      int32  `json:"id_pasien" minimum:"1"`
	NamaVaksin    string `json:"nama_vaksin" minLength:"1" maxLength:"100"`
	TanggalJadwal string `json:"tanggal_jadwal" format:"date"`
}

type CreateImunisasiResponse struct {
	IDImunisasi     int32  `json:"id_imunisasi"`
	StatusImunisasi string `json:"status_imunisasi"`
}

type ImunisasiListItem struct {
	IDImunisasi     int32  `json:"id_imunisasi"`
	NamaPasien      string `json:"nama_pasien"`
	NamaVaksin      string `json:"nama_vaksin"`
	TanggalJadwal   string `json:"tanggal_jadwal"`
	StatusImunisasi string `json:"status_imunisasi"`
}

type ImunisasiListData struct {
	Jadwal    []ImunisasiListItem `json:"jadwal"`
	TotalData int                  `json:"total_data"`
}

type ImunisasiDetail struct {
	IDImunisasi      int32   `json:"id_imunisasi"`
	IDPasien         int32   `json:"id_pasien"`
	NamaPasien       string  `json:"nama_pasien"`
	NamaVaksin       string  `json:"nama_vaksin"`
	TanggalJadwal    string  `json:"tanggal_jadwal"`
	TanggalRealisasi *string `json:"tanggal_realisasi"`
	StatusImunisasi  string  `json:"status_imunisasi"`
}

type UpdateImunisasiRequest struct {
	IDPasien      *int32  `json:"id_pasien,omitempty" minimum:"1"`
	NamaVaksin    *string `json:"nama_vaksin,omitempty" minLength:"1" maxLength:"100"`
	TanggalJadwal *string `json:"tanggal_jadwal,omitempty" format:"date"`
}

type UpdateImunisasiInput struct {
	IDImunisasi int32 `path:"id" minimum:"1"`
	Body        *UpdateImunisasiRequest
}

type UpdateImunisasiResponse struct {
	IDImunisasi int32  `json:"id_imunisasi"`
	UpdatedAt   string `json:"updated_at"`
}

type RealisasiRequest struct {
	TanggalRealisasi string `json:"tanggal_realisasi" format:"date"`
}

type RealisasiResponse struct {
	IDImunisasi      int32  `json:"id_imunisasi"`
	StatusImunisasi  string `json:"status_imunisasi"`
	TanggalRealisasi string `json:"tanggal_realisasi"`
}

type RiwayatImunisasiItem struct {
	IDImunisasi      int32   `json:"id_imunisasi"`
	NamaVaksin       string  `json:"nama_vaksin"`
	TanggalJadwal    string  `json:"tanggal_jadwal"`
	TanggalRealisasi *string `json:"tanggal_realisasi"`
	StatusImunisasi  string  `json:"status_imunisasi"`
}

type RiwayatImunisasiResponse struct {
	IDPasien         int32                  `json:"id_pasien"`
	RiwayatImunisasi []RiwayatImunisasiItem `json:"riwayat_imunisasi"`
}

type StatistikImunisasi struct {
	TotalTargetImunisasi  int32   `json:"total_target_imunisasi"`
	TotalTerealisasi      int32   `json:"total_terealisasi"`
	CakupanPersentase     float64 `json:"cakupan_persentase"`
	VaksinTerbanyak       string  `json:"vaksin_terbanyak"`
}
