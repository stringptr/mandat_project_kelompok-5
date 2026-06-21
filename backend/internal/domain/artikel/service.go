package artikel

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service interface {
	GetAllPublished(ctx context.Context, req *GetAllPublishedRequest) (*ArtikelListData, *errorutils.Error)
	GetAll(ctx context.Context, req *GetAllPublishedRequest) (*ArtikelListData, *errorutils.Error)
	GetByID(ctx context.Context, idArtikel int32) (*ArtikelDetail, *errorutils.Error)
	Create(ctx context.Context, idPenulis int32, isDinkes bool, req *CreateArtikelRequest) (*CreateArtikelResponse, *errorutils.Error)
	Update(ctx context.Context, idPenulis int32, isDinkes bool, req *UpdateArtikelRequest, idArtikel int32) (*ArtikelDetail, *errorutils.Error)
	Delete(ctx context.Context, idArtikel int32) *errorutils.Error
	Review(ctx context.Context, idVerifikator int32, idArtikel int32, req *ReviewArtikelRequest) (*ReviewArtikelResponse, *errorutils.Error)
	GetPending(ctx context.Context, req *GetPendingRequest) (*ArtikelPendingData, *errorutils.Error)
}
