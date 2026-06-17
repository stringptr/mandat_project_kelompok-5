package pemeriksaan

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service interface {
	Create(ctx context.Context, idPetugas int32, req *CreatePemeriksaanRequest) (*CreatePemeriksaanResponse, *errorutils.Error)
	GetByID(ctx context.Context, idHasilPemeriksaan int32) (*DetailPemeriksaanResponse, *errorutils.Error)
	Update(ctx context.Context, req *UpdatePemeriksaanRequest, idHasilPemeriksaan int32) (*UpdatePemeriksaanResponse, *errorutils.Error)
	Delete(ctx context.Context, idHasilPemeriksaan int32) *errorutils.Error
	Verify(ctx context.Context, idHasilPemeriksaan int32, idBidan int32) (*VerifyPemeriksaanResponse, *errorutils.Error)
	GetPending(ctx context.Context, page int, perPage int) (*PendingPemeriksaanData, *errorutils.Error)
	IsOwnPemeriksaan(ctx context.Context, idHasilPemeriksaan int32, idUser int32) (bool, *errorutils.Error)
}
