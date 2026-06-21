package artikel

import (
	"context"
	"slices"

	"github.com/danielgtaylor/huma/v2"
	artikelDomain "github.com/stringptr/SiGizi/backend/internal/domain/artikel"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler struct {
	Service artikelDomain.Service
}

func NewHandler(service artikelDomain.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) GetAllPublished(ctx context.Context, input *artikelDomain.GetAllPublishedRequest) (*httputils.APIResponseOutput[*artikelDomain.ArtikelListData], error) {
	res, err := h.Service.GetAllPublished(ctx, input)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetAll(ctx context.Context, input *artikelDomain.GetAllPublishedRequest) (*httputils.APIResponseOutput[*artikelDomain.ArtikelListData], error) {
	res, err := h.Service.GetAll(ctx, input)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetByID(ctx context.Context, input *struct {
	IDArtikel int32 `path:"id" minimum:"1"`
}) (*httputils.APIResponseOutput[*artikelDomain.ArtikelDetail], error) {
	res, err := h.Service.GetByID(ctx, input.IDArtikel)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) Create(ctx context.Context, input *httputils.APIRequestInput[*artikelDomain.CreateArtikelRequest]) (*httputils.APIResponseOutput[*artikelDomain.CreateArtikelResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	res, err := h.Service.Create(ctx, claims.IDUser, slices.Contains(claims.Roles, "DINKES"), input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewCreatedOutput(res), nil
}

func (h *Handler) Update(ctx context.Context, input *artikelDomain.UpdateArtikelInput) (*httputils.APIResponseOutput[*artikelDomain.ArtikelDetail], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	res, err := h.Service.Update(ctx, claims.IDUser, slices.Contains(claims.Roles, "DINKES"), input.Body, input.IDArtikel)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) Delete(ctx context.Context, input *struct {
	IDArtikel int32 `path:"id" minimum:"1"`
}) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.Delete(ctx, input.IDArtikel)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput[any](nil), nil
}

func (h *Handler) Review(ctx context.Context, input *struct {
	IDArtikel int32                         `path:"id" minimum:"1"`
	Body      *artikelDomain.ReviewArtikelRequest
}) (*httputils.APIResponseOutput[*artikelDomain.ReviewArtikelResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	res, err := h.Service.Review(ctx, claims.IDUser, input.IDArtikel, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetPending(ctx context.Context, input *artikelDomain.GetPendingRequest) (*httputils.APIResponseOutput[*artikelDomain.ArtikelPendingData], error) {
	res, err := h.Service.GetPending(ctx, input)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}
