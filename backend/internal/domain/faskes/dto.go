package faskes

type GetFaskesRequest struct {
	Search string `query:"search"`
}

type FaskesItem struct {
	IDFaskes   int32  `json:"id_faskes"`
	NamaFaskes string `json:"nama_faskes"`
	TipeFaskes string `json:"tipe_faskes"`
}
