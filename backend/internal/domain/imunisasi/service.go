package imunisasi

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service interface {
	GetAll(ctx context.Context) (*ImunisasiListData, *errorutils.Error)
	GetAllByUser(ctx context.Context, idUser int32) (*ImunisasiListData, *errorutils.Error)
	GetByID(ctx context.Context, idImunisasi int32) (*ImunisasiDetail, *errorutils.Error)
	Create(ctx context.Context, req *CreateImunisasiRequest) (*CreateImunisasiResponse, *errorutils.Error)
	Update(ctx context.Context, req *UpdateImunisasiRequest, idImunisasi int32) (*UpdateImunisasiResponse, *errorutils.Error)
	Delete(ctx context.Context, idImunisasi int32) *errorutils.Error
	Realisasi(ctx context.Context, idImunisasi int32, req *RealisasiRequest) (*RealisasiResponse, *errorutils.Error)
	GetByPasienID(ctx context.Context, idPasien int32) (*RiwayatImunisasiResponse, *errorutils.Error)
	GetStatistik(ctx context.Context) (*StatistikImunisasi, *errorutils.Error)
	IsOwnPasien(ctx context.Context, idPasien int32, idUser int32) (bool, *errorutils.Error)
}
