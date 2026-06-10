package pasien

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler interface {
	DaftarIbuHamil(ctx context.Context, input *httputils.APIRequestInput[*DaftarIbuHamilRequest]) (*httputils.APIResponseOutput[any], error)
	DaftarAnak(ctx context.Context, input *httputils.APIRequestInput[*DaftarAnakRequest]) (*httputils.APIResponseOutput[any], error)
	GetAll(ctx context.Context, input *httputils.APIRequestInput[*GetAllPasienRequest]) (*httputils.APIResponseOutput[*PasienListData], error)
	Search(ctx context.Context, input *httputils.APIRequestInput[*SearchPasienRequest]) (*httputils.APIResponseOutput[[]*PasienListItem], error)
	GetByID(ctx context.Context, input *struct{ IDPasien int32 `path:"id" minimum:"1"` }) (*httputils.APIResponseOutput[*PasienDetailResponse], error)
	Update(ctx context.Context, input *httputils.APIRequestInput[*UpdatePasienRequest]) (*httputils.APIResponseOutput[*PasienDetailResponse], error)
	Delete(ctx context.Context, input *struct{ IDPasien int32 `path:"id" minimum:"1"` }) (*httputils.APIResponseOutput[any], error)
}
