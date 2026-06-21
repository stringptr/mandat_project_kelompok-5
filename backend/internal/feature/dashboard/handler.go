package dashboard

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	dashboardDomain "github.com/stringptr/SiGizi/backend/internal/domain/dashboard"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler struct {
	Service dashboardDomain.Service
}

func NewHandler(service dashboardDomain.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) GetDashboardStats(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*dashboardDomain.DashboardStatsResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetDashboardStats(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}

func (h *Handler) GetDistribusiGizi(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*dashboardDomain.DistribusiGiziResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetDistribusiGizi(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}

func (h *Handler) GetTrenStunting(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*dashboardDomain.TrenResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetTrenStunting(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}

func (h *Handler) GetStuntingPerWilayah(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*dashboardDomain.StuntingWilayahResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetStuntingPerWilayah(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}

func (h *Handler) GetKehadiranBulanan(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*dashboardDomain.TrenResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetKehadiranBulanan(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}

func (h *Handler) GetPublicStats(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*dashboardDomain.PublicStatsResponse], error) {
	data, err := h.Service.GetPublicStats(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}

func (h *Handler) GetRiwayat(ctx context.Context, input *dashboardDomain.RiwayatInput) (*httputils.APIResponseOutput[*dashboardDomain.RiwayatResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetRiwayat(ctx, input.IDPasien)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}

func (h *Handler) GetTumbuhKembang(ctx context.Context, input *dashboardDomain.TumbuhKembangInput) (*httputils.APIResponseOutput[*dashboardDomain.TumbuhKembangResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetTumbuhKembang(ctx, input.IDPasien)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}

func (h *Handler) GetJadwalTerdekat(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*dashboardDomain.JadwalTerdekatResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetJadwalTerdekat(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}

func (h *Handler) GetIbuHamilStats(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*dashboardDomain.IbuHamilStatsResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetIbuHamilStats(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}

func (h *Handler) GetIbuHamilPerWilayah(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*dashboardDomain.IbuHamilWilayahResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetIbuHamilPerWilayah(ctx)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}

func (h *Handler) GetSemuaPemeriksaan(ctx context.Context, input *dashboardDomain.GetAllPemeriksaanRequest) (*httputils.APIResponseOutput[*dashboardDomain.PemeriksaanListResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetSemuaPemeriksaan(ctx, input)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(data), nil
}
