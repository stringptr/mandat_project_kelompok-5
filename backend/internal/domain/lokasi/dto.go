package lokasi

type GetLokasiRequest struct {
	Tipe       string `query:"tipe" enum:"Provinsi,Kabupaten,Kota,Kecamatan,Kelurahan"`
	BagianDari int32 `query:"bagian_dari"`
}

type LokasiItem struct {
	IDLokasi   int32  `json:"id_lokasi"`
	NamaLokasi string `json:"nama_lokasi"`
	TipeLokasi string `json:"tipe_lokasi"`
	BagianDari *int32 `json:"bagian_dari"`
}
