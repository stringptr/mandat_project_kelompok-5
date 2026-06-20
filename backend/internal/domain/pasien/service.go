package pasien

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

type Service interface {
	DaftarIbuHamil(ctx context.Context, req *DaftarIbuHamilRequest) *errorutils.Error
	DaftarAnak(ctx context.Context, req *DaftarAnakRequest) *errorutils.Error
	GetAll(ctx context.Context, req *GetAllPasienRequest) (*PasienListData, *errorutils.Error)
	GetAllByUser(ctx context.Context, idUser int32, req *GetAllPasienRequest) (*PasienListData, *errorutils.Error)
	Search(ctx context.Context, req *SearchPasienRequest) (*SearchPasienResponseData, *errorutils.Error)
	GetByID(ctx context.Context, idPasien int32, claims *jwtutils.Claim) (*PasienDetailResponse, *errorutils.Error)
	Update(ctx context.Context, req *UpdatePasienRequest, claims *jwtutils.Claim) (*PasienDetailResponse, *errorutils.Error)
	Delete(ctx context.Context, idPasien int32) *errorutils.Error
}
