package pasien

import "time"

type DaftarIbuHamilRequest struct {
	IDUser           int32  `json:"id_user" minimum:"1"`
	IDPosyandu       int32  `json:"id_posyandu" minimum:"1"`
	HamilKe          int32  `json:"hamil_ke" minimum:"1"`
	BulanMulaiHamil  string `json:"bulan_mulai_hamil" format:"date"`
	Hpht             string `json:"hpht" format:"date"`
	StatusKehamilan  string `json:"status_kehamilan" enum:"Trimester 1,Trimester 2,Trimester 3,Melahirkan,Nifas,Keguguran"`
}

type DaftarAnakRequest struct {
	IDUser             int32   `json:"id_user" minimum:"1"`
	IDPosyandu         int32   `json:"id_posyandu" minimum:"1"`
	IDIbuHamil         *int32  `json:"id_ibu_hamil,omitempty" minimum:"1"`
	IDWali             int32   `json:"id_wali" minimum:"1"`
	NamaAnak           string  `json:"nama_anak" minLength:"1" maxLength:"255"`
	BeratLahir         float64 `json:"berat_lahir" minimum:"0"`
	PanjangLahir       float64 `json:"panjang_lahir" minimum:"0"`
	HubunganDenganWali string  `json:"hubungan_dengan_wali" enum:"Kandung,Tiri,Angkat"`
}

type GetAllPasienRequest struct {
	Page    int    `query:"page" minimum:"1" default:"1"`
	PerPage int    `query:"per_page" minimum:"1" maximum:"100" default:"20"`
	Q       string `query:"q" maxLength:"255"`
	IDUser  int32  `json:"-"`
}

type SearchPasienRequest struct {
	Q string `query:"q" maxLength:"255"`
}

type PasienListItem struct {
	IDPasien       int32  `json:"id_pasien"`
	Nama           string `json:"nama"`
	NIK            string `json:"nik"`
	JenisKelamin   string `json:"jenis_kelamin"`
	Umur           string `json:"umur"`
	NamaPosyandu   string `json:"nama_posyandu"`
	JenisPasien    string `json:"jenis_pasien"`
	StatusKehamilan *string `json:"status_kehamilan,omitempty"`
}

type PasienListData struct {
	Pasien    []PasienListItem `json:"pasien"`
	TotalData int              `json:"total_data"`
	Page      int              `json:"page"`
	PerPage   int              `json:"per_page"`
}

type PasienDetailResponse struct {
	IDPasien         int32               `json:"id_pasien"`
	Nama             string              `json:"nama"`
	NIK              string              `json:"nik"`
	Email            string              `json:"email"`
	NoHp             string              `json:"no_hp"`
	JenisKelamin     string              `json:"jenis_kelamin"`
	TanggalLahir     time.Time           `json:"tanggal_lahir"`
	IDLokasi         int32               `json:"id_lokasi"`
	NamaPosyandu     string              `json:"nama_posyandu"`
	IDPosyandu       int32               `json:"id_posyandu"`
	JenisPasien      string              `json:"jenis_pasien"`
	DataIbuHamil     *IbuHamilData       `json:"data_ibu_hamil,omitempty"`
	DataAnak         *AnakData            `json:"data_anak,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
}

type IbuHamilData struct {
	IDIbuHamil       int32  `json:"id_ibu_hamil"`
	HamilKe          int32  `json:"hamil_ke"`
	BulanMulaiHamil  string `json:"bulan_mulai_hamil"`
	Hpht             string `json:"hpht"`
	StatusKehamilan  string `json:"status_kehamilan"`
}

type AnakData struct {
	NamaAnak           string  `json:"nama_anak"`
	BeratLahir         float64 `json:"berat_lahir"`
	PanjangLahir       float64 `json:"panjang_lahir"`
	HubunganDenganWali string  `json:"hubungan_dengan_wali"`
	NamaWali           string  `json:"nama_wali"`
}

type UpdatePasienRequest struct {
	IDPasien         int32   `json:"-" path:"id" minimum:"1"`
	IDPosyandu       *int32  `json:"id_posyandu,omitempty" minimum:"1"`
	HamilKe          *int32  `json:"hamil_ke,omitempty" minimum:"1"`
	BulanMulaiHamil  *string `json:"bulan_mulai_hamil,omitempty" format:"date"`
	Hpht             *string `json:"hpht,omitempty" format:"date"`
	StatusKehamilan  *string `json:"status_kehamilan,omitempty" enum:"Trimester 1,Trimester 2,Trimester 3,Melahirkan,Nifas,Keguguran"`
	NamaAnak         *string `json:"nama_anak,omitempty" minLength:"1" maxLength:"255"`
	BeratLahir       *float64 `json:"berat_lahir,omitempty" minimum:"0"`
	PanjangLahir     *float64 `json:"panjang_lahir,omitempty" minimum:"0"`
	HubunganDenganWali *string `json:"hubungan_dengan_wali,omitempty" enum:"Kandung,Tiri,Angkat"`
	IDWali           *int32  `json:"id_wali,omitempty" minimum:"1"`
}

type UpdatePasienInput struct {
	IDPasien int32 `path:"id" minimum:"1"`
	Body     *UpdatePasienRequest
}

type DeletePasienRequest struct {
	IDPasien int32 `path:"id" minimum:"1"`
}
