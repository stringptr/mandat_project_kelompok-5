package tindaklanjut

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	tindaklanjutDomain "github.com/stringptr/SiGizi/backend/internal/domain/tindaklanjut"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler struct {
	Service tindaklanjutDomain.Service
}

func NewHandler(service tindaklanjutDomain.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) GetPasienTindakLanjut(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*tindaklanjutDomain.PasienTindakLanjutData], error) {
	res, err := h.Service.GetPasienTindakLanjut(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetDetailPasienByID(ctx context.Context, input *struct {
	IDPasien int32 `path:"id" minimum:"1"`
}) (*httputils.APIResponseOutput[*tindaklanjutDomain.DetailPasienTindakLanjut], error) {
	res, err := h.Service.GetDetailPasienByID(ctx, input.IDPasien)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) CreateTindakLanjut(ctx context.Context, input *httputils.APIRequestInput[*tindaklanjutDomain.CreateTindakLanjutRequest]) (*httputils.APIResponseOutput[*tindaklanjutDomain.CreateTindakLanjutResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	res, err := h.Service.CreateTindakLanjut(ctx, claims.IDUser, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewCreatedOutput(res), nil
}

func (h *Handler) UpdateStatusRujukan(ctx context.Context, input *struct {
	IDRujukan int32 `path:"id" minimum:"1"`
	Body      *tindaklanjutDomain.UpdateStatusRujukanRequest
}) (*httputils.APIResponseOutput[*tindaklanjutDomain.UpdateStatusRujukanResponse], error) {
	res, err := h.Service.UpdateStatusRujukan(ctx, input.IDRujukan, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetStatusTindakLanjut(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*tindaklanjutDomain.StatusTindakLanjutData], error) {
	res, err := h.Service.GetStatusTindakLanjut(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetLaporanTindakLanjut(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*tindaklanjutDomain.LaporanTindakLanjutData], error) {
	res, err := h.Service.GetLaporanTindakLanjut(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetDetailTindakLanjutByID(ctx context.Context, input *struct {
	IDTindakLanjut int32 `path:"id" minimum:"1"`
}) (*httputils.APIResponseOutput[*tindaklanjutDomain.DetailTindakLanjutPasien], error) {
	res, err := h.Service.GetDetailTindakLanjutByID(ctx, input.IDTindakLanjut)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}
