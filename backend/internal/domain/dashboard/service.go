package dashboard

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service interface {
	GetDashboardStats(ctx context.Context) (*DashboardStatsResponse, *errorutils.Error)
	GetDistribusiGizi(ctx context.Context) (*DistribusiGiziResponse, *errorutils.Error)
	GetTrenStunting(ctx context.Context) (*TrenResponse, *errorutils.Error)
	GetStuntingPerWilayah(ctx context.Context) (*StuntingWilayahResponse, *errorutils.Error)
	GetKehadiranBulanan(ctx context.Context) (*TrenResponse, *errorutils.Error)
	GetPublicStats(ctx context.Context) (*PublicStatsResponse, *errorutils.Error)
	GetRiwayat(ctx context.Context, idPasien int32) (*RiwayatResponse, *errorutils.Error)
	GetTumbuhKembang(ctx context.Context, idPasien int32) (*TumbuhKembangResponse, *errorutils.Error)
	GetJadwalTerdekat(ctx context.Context) (*JadwalTerdekatResponse, *errorutils.Error)
	GetIbuHamilStats(ctx context.Context) (*IbuHamilStatsResponse, *errorutils.Error)
	GetIbuHamilPerWilayah(ctx context.Context) (*IbuHamilWilayahResponse, *errorutils.Error)
	GetSemuaPemeriksaan(ctx context.Context, req *GetAllPemeriksaanRequest) (*PemeriksaanListResponse, *errorutils.Error)
}
