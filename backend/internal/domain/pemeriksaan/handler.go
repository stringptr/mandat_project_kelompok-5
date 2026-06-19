package pemeriksaan

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler interface {
	GetAll(ctx context.Context, input *httputils.APIRequestInput[*GetAllPemeriksaanRequest]) (*httputils.APIResponseOutput[*PemeriksaanListData], error)
	Create(ctx context.Context, input *httputils.APIRequestInput[*CreatePemeriksaanRequest]) (*httputils.APIResponseOutput[*CreatePemeriksaanResponse], error)
	GetByID(ctx context.Context, input *struct {
		IDHasilPemeriksaan int32 `path:"id" minimum:"1"`
	}) (*httputils.APIResponseOutput[*DetailPemeriksaanResponse], error)
	Update(ctx context.Context, input *UpdatePemeriksaanInput) (*httputils.APIResponseOutput[*UpdatePemeriksaanResponse], error)
	Delete(ctx context.Context, input *struct {
		IDHasilPemeriksaan int32 `path:"id" minimum:"1"`
	}) (*httputils.APIResponseOutput[any], error)
	Verify(ctx context.Context, input *struct {
		IDHasilPemeriksaan int32 `path:"id" minimum:"1"`
	}) (*httputils.APIResponseOutput[*VerifyPemeriksaanResponse], error)
	GetPending(ctx context.Context, input *GetPendingPemeriksaanRequest) (*httputils.APIResponseOutput[*PendingPemeriksaanData], error)
}
