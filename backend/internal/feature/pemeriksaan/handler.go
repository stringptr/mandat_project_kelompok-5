package pemeriksaan

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	pemeriksaanDomain "github.com/stringptr/SiGizi/backend/internal/domain/pemeriksaan"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/pagination"
)

type Handler struct {
	Service pemeriksaanDomain.Service
}

func NewHandler(service pemeriksaanDomain.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Create(ctx context.Context, input *httputils.APIRequestInput[*pemeriksaanDomain.CreatePemeriksaanRequest]) (*httputils.APIResponseOutput[*pemeriksaanDomain.CreatePemeriksaanResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	res, err := h.Service.Create(ctx, claims.IDUser, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewCreatedOutput(res), nil
}

func (h *Handler) GetByID(ctx context.Context, input *struct {
	IDHasilPemeriksaan int32 `path:"id" minimum:"1"`
}) (*httputils.APIResponseOutput[*pemeriksaanDomain.DetailPemeriksaanResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	if !httputils.IsPetugas(claims.Roles) {
		isOwner, err := h.Service.IsOwnPemeriksaan(ctx, input.IDHasilPemeriksaan, claims.IDUser)
		if err != nil {
			return nil, errorutils.ToHumaError(err)
		}
		if !isOwner {
			return nil, huma.Error404NotFound("Data pemeriksaan tidak ditemukan.")
		}
	}

	res, err := h.Service.GetByID(ctx, input.IDHasilPemeriksaan)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(res), nil
}

func (h *Handler) Update(ctx context.Context, input *pemeriksaanDomain.UpdatePemeriksaanInput) (*httputils.APIResponseOutput[*pemeriksaanDomain.UpdatePemeriksaanResponse], error) {
	res, err := h.Service.Update(ctx, input.Body, input.IDHasilPemeriksaan)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(res), nil
}

func (h *Handler) Delete(ctx context.Context, input *struct {
	IDHasilPemeriksaan int32 `path:"id" minimum:"1"`
}) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.Delete(ctx, input.IDHasilPemeriksaan)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput[any](nil), nil
}

func (h *Handler) Verify(ctx context.Context, input *struct {
	IDHasilPemeriksaan int32 `path:"id" minimum:"1"`
}) (*httputils.APIResponseOutput[*pemeriksaanDomain.VerifyPemeriksaanResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	res, err := h.Service.Verify(ctx, input.IDHasilPemeriksaan, claims.IDUser)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetPending(ctx context.Context, input *pemeriksaanDomain.GetPendingPemeriksaanRequest) (*httputils.APIResponseOutput[*pemeriksaanDomain.PendingPemeriksaanData], error) {
	if input == nil {
		input = &pemeriksaanDomain.GetPendingPemeriksaanRequest{}
	}
	page := pagination.ValidatePage(input.Page)
	perPage := pagination.ValidatePerPage(input.PerPage)

	res, err := h.Service.GetPending(ctx, page, perPage)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(res), nil
}
