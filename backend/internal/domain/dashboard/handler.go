package dashboard

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler interface {
	GetDashboardStats(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*DashboardStatsResponse], error)
	GetDistribusiGizi(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*DistribusiGiziResponse], error)
	GetTrenStunting(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*TrenResponse], error)
	GetStuntingPerWilayah(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*StuntingWilayahResponse], error)
	GetKehadiranBulanan(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*TrenResponse], error)
	GetPublicStats(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*PublicStatsResponse], error)
	GetRiwayat(ctx context.Context, input *RiwayatInput) (*httputils.APIResponseOutput[*RiwayatResponse], error)
	GetTumbuhKembang(ctx context.Context, input *TumbuhKembangInput) (*httputils.APIResponseOutput[*TumbuhKembangResponse], error)
	GetJadwalTerdekat(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*JadwalTerdekatResponse], error)
	GetIbuHamilStats(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*IbuHamilStatsResponse], error)
	GetIbuHamilPerWilayah(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*IbuHamilWilayahResponse], error)
	GetSemuaPemeriksaan(ctx context.Context, input *GetAllPemeriksaanRequest) (*httputils.APIResponseOutput[*PemeriksaanListResponse], error)
}
