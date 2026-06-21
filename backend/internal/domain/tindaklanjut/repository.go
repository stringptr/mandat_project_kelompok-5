package tindaklanjut

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Repo interface {
	GetPasienByID(ctx context.Context, idPasien int32) (*model.Pasien, error)
	GetPasienTindakLanjut(ctx context.Context, page int, perPage int) ([]*PasienTindakLanjutJoinRow, int, error)
	GetDetailPasienByID(ctx context.Context, idPasien int32) (*DetailPasienJoinRow, error)
	GetRiwayatPemeriksaan(ctx context.Context, idPasien int32) ([]*RiwayatPemeriksaanJoinRow, error)
	GetTindakLanjutByHasilPemeriksaan(ctx context.Context, idHasilPemeriksaan int32) (*model.TindakLanjut, error)
	CreateTindakLanjut(ctx context.Context, data *model.TindakLanjut) error
	CreateRujukan(ctx context.Context, data *model.Rujukan) error
	UpdateStatusRujukan(ctx context.Context, idRujukan int32, statusRujukan model.StatusRujukan) (*model.Rujukan, error)
	GetStatusTindakLanjut(ctx context.Context, page int, perPage int) ([]*StatusTindakLanjutJoinRow, int, error)
	GetLaporanTindakLanjut(ctx context.Context) ([]*LaporanTindakLanjutJoinRow, error)
	GetDetailTindakLanjutByID(ctx context.Context, idTindakLanjut int32) (*DetailTindakLanjutJoinRow, error)
	GetPasienIDByHasilPemeriksaanID(ctx context.Context, idHasilPemeriksaan int32) (*int32, error)
}

type PasienTindakLanjutJoinRow struct {
	IDPasien           int32
	NamaPasien         string
	StatusGizi         string
	StatusPasien       string
	TanggalPemeriksaan string
}

type DetailPasienJoinRow struct {
	IDPasien   int32
	NamaPasien string
	Usia       string
	StatusGizi string
	StatusStunting string
	Catatan    string
}

type RiwayatPemeriksaanJoinRow struct {
	Tanggal    string
	BeratBadan float64
	TinggiBadan float64
}

type StatusTindakLanjutJoinRow struct {
	IDPasien       int32
	NamaPasien     string
	StatusPasien   string
	StatusRujukan  string
	TanggalRujukan string
}

type LaporanTindakLanjutJoinRow struct {
	Wilayah              string
	JumlahPasienDirujuk  int32
	JumlahPasienDiterima int32
	JumlahPasienDiproses int32
}

type DetailTindakLanjutJoinRow struct {
	IDTindakLanjut int32
	StatusPasien   string
	CatatanMedis   *string
	Rekomendasi    *string
	JadwalKontrol  *string
	StatusRujukan  *string
	NamaFaskes     *string
}
