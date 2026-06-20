package imunisasi

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler interface {
	GetAll(ctx context.Context, input *GetAllImunisasiRequest) (*httputils.APIResponseOutput[*ImunisasiListData], error)
	GetByID(ctx context.Context, input *struct {
		IDImunisasi int32 `path:"id" minimum:"1"`
	}) (*httputils.APIResponseOutput[*ImunisasiDetail], error)
	Create(ctx context.Context, input *httputils.APIRequestInput[*CreateImunisasiRequest]) (*httputils.APIResponseOutput[*CreateImunisasiResponse], error)
	Update(ctx context.Context, input *UpdateImunisasiInput) (*httputils.APIResponseOutput[*UpdateImunisasiResponse], error)
	Delete(ctx context.Context, input *struct {
		IDImunisasi int32 `path:"id" minimum:"1"`
	}) (*httputils.APIResponseOutput[any], error)
	Realisasi(ctx context.Context, input *struct {
		IDImunisasi int32  `path:"id" minimum:"1"`
		Body        *RealisasiRequest
	}) (*httputils.APIResponseOutput[*RealisasiResponse], error)
	GetByPasienID(ctx context.Context, input *struct {
		IDPasien int32 `path:"id_pasien" minimum:"1"`
	}) (*httputils.APIResponseOutput[*RiwayatImunisasiResponse], error)
	GetStatistik(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*StatistikImunisasi], error)
}
