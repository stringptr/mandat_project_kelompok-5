package imunisasi

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

type Service interface {
	GetAll(ctx context.Context, req *GetAllImunisasiRequest) (*ImunisasiListData, *errorutils.Error)
	GetAllByUser(ctx context.Context, idUser int32, req *GetAllImunisasiRequest) (*ImunisasiListData, *errorutils.Error)
	GetByID(ctx context.Context, idImunisasi int32, claims *jwtutils.Claim) (*ImunisasiDetail, *errorutils.Error)
	Create(ctx context.Context, req *CreateImunisasiRequest) (*CreateImunisasiResponse, *errorutils.Error)
	Update(ctx context.Context, req *UpdateImunisasiRequest, idImunisasi int32) (*UpdateImunisasiResponse, *errorutils.Error)
	Delete(ctx context.Context, idImunisasi int32) *errorutils.Error
	Realisasi(ctx context.Context, idImunisasi int32, req *RealisasiRequest) (*RealisasiResponse, *errorutils.Error)
	GetByPasienID(ctx context.Context, idPasien int32, claims *jwtutils.Claim) (*RiwayatImunisasiResponse, *errorutils.Error)
	GetStatistik(ctx context.Context) (*StatistikImunisasi, *errorutils.Error)
	IsOwnPasien(ctx context.Context, idPasien int32, idUser int32) (bool, *errorutils.Error)
}
