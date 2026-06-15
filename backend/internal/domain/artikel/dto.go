package artikel

import "time"

type CreateArtikelRequest struct {
	Judul     string `json:"judul" minLength:"1" maxLength:"255"`
	IsiArtikel string `json:"isi_artikel" minLength:"1"`
	Kategori  string `json:"kategori,omitempty" maxLength:"100"`
}

type CreateArtikelResponse struct {
	IDArtikel      int32  `json:"id_artikel"`
	StatusArtikel  string `json:"status_artikel"`
}

type ArtikelListItem struct {
	IDArtikel      int32  `json:"id_artikel"`
	Judul          string `json:"judul"`
	Kategori       string `json:"kategori"`
	Ringkasan      string `json:"ringkasan"`
	NamaPenulis    string `json:"nama_penulis"`
	TanggalPublish string `json:"tanggal_publish"`
}

type ArtikelListData struct {
	Artikel   []ArtikelListItem `json:"artikel"`
	TotalData int                `json:"total_data"`
}

type ArtikelDetail struct {
	IDArtikel      int32   `json:"id_artikel"`
	Judul          string  `json:"judul"`
	IsiArtikel     string  `json:"isi_artikel"`
	Kategori       string  `json:"kategori"`
	NamaPenulis    string  `json:"nama_penulis"`
	NamaVerifikator *string `json:"nama_verifikator"`
	TanggalPublish *string `json:"tanggal_publish"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type UpdateArtikelRequest struct {
	Judul      *string `json:"judul,omitempty" minLength:"1" maxLength:"255"`
	IsiArtikel *string `json:"isi_artikel,omitempty" minLength:"1"`
	Kategori   *string `json:"kategori,omitempty" maxLength:"100"`
}

type ArtikelPendingItem struct {
	IDArtikel     int32  `json:"id_artikel"`
	Judul         string `json:"judul"`
	NamaPenulis   string `json:"nama_penulis"`
	CreatedAt     string `json:"created_at"`
	StatusArtikel string `json:"status_artikel"`
}

type ArtikelPendingData struct {
	Artikel   []ArtikelPendingItem `json:"artikel"`
	TotalData int                   `json:"total_data"`
}

type ReviewArtikelRequest struct {
	Aksi          string `json:"aksi" enum:"setujui,tolak"`
	CatatanReview string `json:"catatan_review,omitempty"`
}

type ReviewArtikelResponse struct {
	IDArtikel      int32   `json:"id_artikel"`
	StatusArtikel  string  `json:"status_artikel"`
	TanggalPublish *string `json:"tanggal_publish"`
}

type UpdateArtikelInput struct {
	IDArtikel int32 `path:"id" minimum:"1"`
	Body      *UpdateArtikelRequest
}
