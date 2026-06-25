package pasien

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Repo interface {
	GetUserByID(ctx context.Context, idUser int32) (*model.UserAccount, error)
	GetPasienByUserID(ctx context.Context, idUser int32) (*model.Pasien, error)
	GetPosyanduByID(ctx context.Context, idPosyandu int32) (*model.Posyandu, error)
	GetIbuHamilByPasienID(ctx context.Context, idPasien int32) ([]*model.IbuHamil, error)
	GetAnakByPasienID(ctx context.Context, idPasien int32) (*model.Anak, error)

	CreatePasien(ctx context.Context, data *model.Pasien) error
	CreateIbuHamil(ctx context.Context, data *model.IbuHamil) error
	CreateAnak(ctx context.Context, data *model.Anak) error

	GetAllPaginated(ctx context.Context, page int, perPage int, q string, idPosyandu int32) ([]*PasienJoinRow, int, error)
	GetAllPaginatedByUser(ctx context.Context, page int, perPage int, q string, idUser int32) ([]*PasienJoinRow, int, error)
	Search(ctx context.Context, q string, page int, perPage int) ([]*PasienJoinRow, int, error)
	GetDetailByID(ctx context.Context, idPasien int32) (*PasienDetailJoinRow, error)
	CheckPasienOwnership(ctx context.Context, idPasien int32, idUser int32) (bool, error)

	UpdatePasien(ctx context.Context, data *model.Pasien) error
	UpdateIbuHamil(ctx context.Context, data *model.IbuHamil) error
	UpdateAnak(ctx context.Context, data *model.Anak) error

	DeletePasien(ctx context.Context, idPasien int32) error
}

type PasienJoinRow struct {
	IDPasien     int32
	Nama         string
	NIK          string
	JenisKelamin string
	TanggalLahir string
	NamaPosyandu string
	JenisPasien  string
	StatusKehamilan *string
}

type PasienDetailJoinRow struct {
	IDPasien         int32
	Nama             string
	NIK              string
	Email            string
	NoHp             string
	JenisKelamin     string
	TanggalLahir     string
	IDLokasi         int32
	NamaPosyandu     string
	IDPosyandu       int32
	JenisPasien      string
	CreatedAt        string
	UpdatedAt        string
	IDIbuHamil       *int32
	HamilKe          *int32
	BulanMulaiHamil  *string
	Hpht             *string
	StatusKehamilan  *string
	NamaAnak         *string
	BeratLahir       *float64
	PanjangLahir     *float64
	HubunganDenganWali *string
	NamaWali         *string
}
