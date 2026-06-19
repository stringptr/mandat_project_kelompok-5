package pemeriksaan

import (
	"context"
	"time"

	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Repo interface {
	GetAllPaginated(ctx context.Context, page, perPage int, q string) ([]*GetAllJoinRow, int, error)
	GetByID(ctx context.Context, idHasilPemeriksaan int32) (*model.HasilPemeriksaan, error)
	GetJadwalImunisasiByID(ctx context.Context, idJadwalImunisasi int32) (*model.JadwalImunisasi, error)
	GetUserByID(ctx context.Context, idUser int32) (*model.UserAccount, error)
	GetNamaPetugasByID(ctx context.Context, idUser int32) (string, error)
	Create(ctx context.Context, data *model.HasilPemeriksaan) error
	Update(ctx context.Context, data *model.HasilPemeriksaan) error
	Delete(ctx context.Context, idHasilPemeriksaan int32) error
	GetPendingVerification(ctx context.Context, page int, perPage int) ([]*PendingJoinRow, int, error)
	GetDetailJoinByID(ctx context.Context, idHasilPemeriksaan int32) (*DetailJoinRow, error)
	CheckPemeriksaanOwnership(ctx context.Context, idHasilPemeriksaan int32, idUser int32) (bool, error)
}

type GetAllJoinRow struct {
	IDHasilPemeriksaan int32
	NamaPasien         string
	DiinputOleh        string
	StatusStunting     string
	StatusGizi         string
	TanggalInput       time.Time
}

type PendingJoinRow struct {
	IDHasilPemeriksaan int32
	NamaPasien         string
	DiinputOleh        string
	TanggalInput       time.Time
}

type DetailJoinRow struct {
	IDHasilPemeriksaan int32
	IDPasien           int32
	NamaPasien         string
	BeratBadan         float64
	TinggiBadan        float64
	LingkarKepala      float64
	TekananDarah       string
	StatusStunting     string
	StatusGizi         string
	Catatan            *string
}
