package artikel

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler interface {
	GetAllPublished(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*ArtikelListData], error)
	GetByID(ctx context.Context, input *struct {
		IDArtikel int32 `path:"id" minimum:"1"`
	}) (*httputils.APIResponseOutput[*ArtikelDetail], error)
	Create(ctx context.Context, input *httputils.APIRequestInput[*CreateArtikelRequest]) (*httputils.APIResponseOutput[*CreateArtikelResponse], error)
	Update(ctx context.Context, input *UpdateArtikelInput) (*httputils.APIResponseOutput[*ArtikelDetail], error)
	Delete(ctx context.Context, input *struct {
		IDArtikel int32 `path:"id" minimum:"1"`
	}) (*httputils.APIResponseOutput[any], error)
	Review(ctx context.Context, input *struct {
		IDArtikel int32 `path:"id" minimum:"1"`
		Body      *ReviewArtikelRequest
	}) (*httputils.APIResponseOutput[*ReviewArtikelResponse], error)
	GetPending(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*ArtikelPendingData], error)
}
