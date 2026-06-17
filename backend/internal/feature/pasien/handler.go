package pasien

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	pasienDomain "github.com/stringptr/SiGizi/backend/internal/domain/pasien"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler struct {
	Service pasienDomain.Service
}

func NewHandler(service pasienDomain.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) DaftarIbuHamil(ctx context.Context, input *httputils.APIRequestInput[*pasienDomain.DaftarIbuHamilRequest]) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.DaftarIbuHamil(ctx, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewCreatedOutput[any](nil), nil
}

func (h *Handler) DaftarAnak(ctx context.Context, input *httputils.APIRequestInput[*pasienDomain.DaftarAnakRequest]) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.DaftarAnak(ctx, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewCreatedOutput[any](nil), nil
}

func (h *Handler) GetAll(ctx context.Context, input *pasienDomain.GetAllPasienRequest) (*httputils.APIResponseOutput[*pasienDomain.PasienListData], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	var res *pasienDomain.PasienListData
	var err *errorutils.Error
	if httputils.IsPetugas(claims.Roles) {
		res, err = h.Service.GetAll(ctx, input)
	} else {
		res, err = h.Service.GetAllByUser(ctx, claims.IDUser, input)
	}
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(res), nil
}

func (h *Handler) Search(ctx context.Context, input *pasienDomain.SearchPasienRequest) (*httputils.APIResponseOutput[*pasienDomain.SearchPasienResponseData], error) {
	res, err := h.Service.Search(ctx, input)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetByID(ctx context.Context, input *struct {
	IDPasien int32 `path:"id" minimum:"1"`
},
) (*httputils.APIResponseOutput[*pasienDomain.PasienDetailResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	if !httputils.IsPetugas(claims.Roles) {
		isOwner, err := h.Service.IsOwnPasien(ctx, input.IDPasien, claims.IDUser)
		if err != nil {
			return nil, errorutils.ToHumaError(err)
		}
		if !isOwner {
			return nil, huma.Error404NotFound("Data pasien tidak ditemukan.")
		}
	}

	res, err := h.Service.GetByID(ctx, input.IDPasien)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(res), nil
}

func (h *Handler) Update(ctx context.Context, input *pasienDomain.UpdatePasienInput) (*httputils.APIResponseOutput[*pasienDomain.PasienDetailResponse], error) {
	input.Body.IDPasien = input.IDPasien
	res, err := h.Service.Update(ctx, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(res), nil
}

func (h *Handler) Delete(ctx context.Context, input *struct {
	IDPasien int32 `path:"id" minimum:"1"`
},
) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.Delete(ctx, input.IDPasien)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput[any](nil), nil
}
