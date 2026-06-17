package imunisasi

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Repo interface {
	GetAll(ctx context.Context, page int, perPage int, q string) ([]*ImunisasiJoinRow, int, error)
	GetAllByUser(ctx context.Context, idUser int32, page int, perPage int, q string) ([]*ImunisasiJoinRow, int, error)
	GetByID(ctx context.Context, idImunisasi int32) (*model.JadwalImunisasi, error)
	GetDetailJoinByID(ctx context.Context, idImunisasi int32) (*DetailJoinRow, error)
	GetByPasienID(ctx context.Context, idPasien int32) ([]*model.JadwalImunisasi, error)
	GetPasienByID(ctx context.Context, idPasien int32) (*model.Pasien, error)
	GetNamaPasienByID(ctx context.Context, idPasien int32) (string, error)
	CheckPasienOwnership(ctx context.Context, idPasien int32, idUser int32) (bool, error)
	Create(ctx context.Context, data *model.JadwalImunisasi) error
	Update(ctx context.Context, data *model.JadwalImunisasi) error
	Delete(ctx context.Context, idImunisasi int32) error
	UpdateRealisasi(ctx context.Context, idImunisasi int32, tanggalRealisasi string) error
	GetStatistik(ctx context.Context) (*StatistikRow, error)
	GetVaksinTerbanyak(ctx context.Context) (string, error)
}

type ImunisasiJoinRow struct {
	IDImunisasi     int32
	NamaPasien      string
	NamaVaksin      string
	TanggalJadwal   string
	StatusImunisasi string
}

type DetailJoinRow struct {
	IDImunisasi      int32
	IDPasien         int32
	NamaPasien       string
	NamaVaksin       string
	TanggalJadwal    string
	TanggalRealisasi *string
	StatusImunisasi  string
}

type StatistikRow struct {
	TotalTarget    int32
	TotalRealisasi int32
}
