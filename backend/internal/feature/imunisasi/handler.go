package imunisasi

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	imunisasiDomain "github.com/stringptr/SiGizi/backend/internal/domain/imunisasi"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler struct {
	Service imunisasiDomain.Service
}

func NewHandler(service imunisasiDomain.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) GetAll(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*imunisasiDomain.ImunisasiListData], error) {
	res, err := h.Service.GetAll(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetByID(ctx context.Context, input *struct {
	IDImunisasi int32 `path:"id" minimum:"1"`
}) (*httputils.APIResponseOutput[*imunisasiDomain.ImunisasiDetail], error) {
	res, err := h.Service.GetByID(ctx, input.IDImunisasi)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) Create(ctx context.Context, input *httputils.APIRequestInput[*imunisasiDomain.CreateImunisasiRequest]) (*httputils.APIResponseOutput[*imunisasiDomain.CreateImunisasiResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}
	res, err := h.Service.Create(ctx, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewCreatedOutput(res), nil
}

func (h *Handler) Update(ctx context.Context, input *imunisasiDomain.UpdateImunisasiInput) (*httputils.APIResponseOutput[*imunisasiDomain.UpdateImunisasiResponse], error) {
	res, err := h.Service.Update(ctx, input.Body, input.IDImunisasi)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) Delete(ctx context.Context, input *struct {
	IDImunisasi int32 `path:"id" minimum:"1"`
}) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.Delete(ctx, input.IDImunisasi)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput[any](nil), nil
}

func (h *Handler) Realisasi(ctx context.Context, input *struct {
	IDImunisasi int32                      `path:"id" minimum:"1"`
	Body        *imunisasiDomain.RealisasiRequest
}) (*httputils.APIResponseOutput[*imunisasiDomain.RealisasiResponse], error) {
	res, err := h.Service.Realisasi(ctx, input.IDImunisasi, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetByPasienID(ctx context.Context, input *struct {
	IDPasien int32 `path:"id_pasien" minimum:"1"`
}) (*httputils.APIResponseOutput[*imunisasiDomain.RiwayatImunisasiResponse], error) {
	res, err := h.Service.GetByPasienID(ctx, input.IDPasien)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}

func (h *Handler) GetStatistik(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*imunisasiDomain.StatistikImunisasi], error) {
	res, err := h.Service.GetStatistik(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}
