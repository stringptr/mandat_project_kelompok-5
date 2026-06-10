package pasien

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service interface {
	DaftarIbuHamil(ctx context.Context, req *DaftarIbuHamilRequest) *errorutils.Error
	DaftarAnak(ctx context.Context, req *DaftarAnakRequest) *errorutils.Error
	GetAll(ctx context.Context, req *GetAllPasienRequest) (*PasienListData, *errorutils.Error)
	Search(ctx context.Context, req *SearchPasienRequest) ([]*PasienListItem, *errorutils.Error)
	GetByID(ctx context.Context, idPasien int32) (*PasienDetailResponse, *errorutils.Error)
	Update(ctx context.Context, req *UpdatePasienRequest) (*PasienDetailResponse, *errorutils.Error)
	Delete(ctx context.Context, idPasien int32) *errorutils.Error
}
