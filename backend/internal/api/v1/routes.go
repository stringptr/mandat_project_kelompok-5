package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/stringptr/SiGizi/backend/internal/config"
	artikelDomain "github.com/stringptr/SiGizi/backend/internal/domain/artikel"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	dashboardDomain "github.com/stringptr/SiGizi/backend/internal/domain/dashboard"
	faskesDomain "github.com/stringptr/SiGizi/backend/internal/domain/faskes"
	imunisasiDomain "github.com/stringptr/SiGizi/backend/internal/domain/imunisasi"
	tindaklanjutDomain "github.com/stringptr/SiGizi/backend/internal/domain/tindaklanjut"
	jwtblacklistDomain "github.com/stringptr/SiGizi/backend/internal/domain/jwtblacklist"
	lokasiDomain "github.com/stringptr/SiGizi/backend/internal/domain/lokasi"
	notificationDomain "github.com/stringptr/SiGizi/backend/internal/domain/notification"
	pasienDomain "github.com/stringptr/SiGizi/backend/internal/domain/pasien"
	pemeriksaanDomain "github.com/stringptr/SiGizi/backend/internal/domain/pemeriksaan"
	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
)

type Dependency struct {
	AuthConfig          config.AuthConfig
	JWTUtil             jwtutils.JWT
	UserAccountHandler  userAccountDomain.Handler
	AuthHandler         authDomain.Handler
	BlacklistRepo       jwtblacklistDomain.Repo
	NotifPublisher      notificationDomain.Publisher
	NotifHandler        notificationDomain.Handler
	PasienHandler       pasienDomain.Handler
	LokasiHandler       lokasiDomain.Handler
	FaskesHandler       faskesDomain.Handler
	PemeriksaanHandler  pemeriksaanDomain.Handler
	ImunisasiHandler    imunisasiDomain.Handler
	ArtikelHandler      artikelDomain.Handler
	TindakLanjutHandler tindaklanjutDomain.Handler
	DashboardHandler    dashboardDomain.Handler
}

func RegisterRoutes(api huma.API, r chi.Router, d *Dependency) {
	publicGroup := huma.NewGroup(api, "")
	authAccess := huma.NewGroup(api, "")
	authRefresh := huma.NewGroup(api, "")
	userGroup := huma.NewGroup(authAccess, "")
	publicGroup = huma.NewGroup(api, "")
	adminGroup := huma.NewGroup(userGroup, "")
	bidanGroup := huma.NewGroup(adminGroup, "")
	kaderGroup := huma.NewGroup(adminGroup, "")
	dinkesGroup := huma.NewGroup(userGroup, "")
	superAdminGroup := huma.NewGroup(adminGroup, "")

	authAccess.UseMiddleware(middleware.AuthAccessMiddleware(api, &d.JWTUtil, d.BlacklistRepo))
	authRefresh.UseMiddleware(middleware.AuthRefreshMiddleware(api, &d.JWTUtil))
	userGroup.UseMiddleware(middleware.RequireRole(authAccess, "USER"))
	adminGroup.UseMiddleware(middleware.RequireRole(authAccess, "ADMIN"))
	bidanGroup.UseMiddleware(middleware.RequireRole(authAccess, "BIDAN"))
	kaderGroup.UseMiddleware(middleware.RequireRole(authAccess, "KADER"))
	dinkesGroup.UseMiddleware(middleware.RequireRole(authAccess, "DINKES"))
	superAdminGroup.UseMiddleware(middleware.RequireRole(authAccess, "SUPER_ADMIN"))

	nonAuthenticatedOnlyGroup := huma.NewGroup(api, "")
	nonAuthenticatedOnlyGroup.UseMiddleware(middleware.NonAuthenticatedOnlyMiddleware(api, &d.JWTUtil, d.BlacklistRepo))

	huma.Post(nonAuthenticatedOnlyGroup, "/auth/register", d.AuthHandler.Register)
	huma.Post(nonAuthenticatedOnlyGroup, "/auth/login", d.AuthHandler.Login)

	huma.Get(publicGroup, "/lokasi", d.LokasiHandler.GetLokasi)

	huma.Get(userGroup, "/faskes", d.FaskesHandler.GetFaskes)

	huma.Post(authRefresh, "/auth/refresh", d.AuthHandler.Refresh)
	huma.Post(nonAuthenticatedOnlyGroup, "/auth/logout", d.AuthHandler.Logout)

	huma.Get(userGroup, "/auth/me", d.AuthHandler.Me)

	huma.Patch(adminGroup, "/users/{id_user}/verification", d.AuthHandler.VerifyUser)

	huma.Get(adminGroup, "/users", d.UserAccountHandler.GetAllUsers)
	huma.Get(userGroup, "/users/{id}", d.UserAccountHandler.GetUserByID)
	huma.Patch(userGroup, "/users/{id}", d.UserAccountHandler.UpdateUser)

	huma.Post(superAdminGroup, "/users", d.UserAccountHandler.CreateUser)
	huma.Patch(superAdminGroup, "/users/{id}/role", d.UserAccountHandler.UpdateUserRole)
	huma.Get(superAdminGroup, "/admin/audit-logs", d.UserAccountHandler.GetAuditLogs)
	huma.Delete(superAdminGroup, "/users/{id}", d.UserAccountHandler.DeleteUser)

	huma.Get(bidanGroup, "/notifikasi/bidan", d.NotifHandler.GetBidanDashboard)
	huma.Get(adminGroup, "/notifikasi/aktivitas", d.NotifHandler.GetActivity)
	huma.Get(userGroup, "/notifikasi", d.NotifHandler.GetNotifikasi)
	huma.Get(userGroup, "/notifikasi/{id}", d.NotifHandler.GetNotifikasiDetail)
	huma.Patch(userGroup, "/notifikasi/{id}/read", d.NotifHandler.MarkRead)
	huma.Patch(userGroup, "/notifikasi/read-all", d.NotifHandler.MarkAllRead)

	huma.Post(adminGroup, "/pasien/ibu-hamil", d.PasienHandler.DaftarIbuHamil)
	huma.Post(adminGroup, "/pasien/anak", d.PasienHandler.DaftarAnak)
	huma.Get(userGroup, "/monitoring/pasien", d.PasienHandler.GetAll)
	huma.Get(adminGroup, "/monitoring/pasien/search", d.PasienHandler.Search)
	huma.Get(userGroup, "/monitoring/pasien/{id}", d.PasienHandler.GetByID)
	huma.Patch(adminGroup, "/pasien/{id}", d.PasienHandler.Update)
	huma.Delete(bidanGroup, "/pasien/{id}", d.PasienHandler.Delete)

	huma.Get(adminGroup, "/monitoring/pemeriksaan", d.PemeriksaanHandler.GetAll)
	huma.Post(adminGroup, "/monitoring/pemeriksaan", d.PemeriksaanHandler.Create)
	huma.Get(userGroup, "/monitoring/pemeriksaan/{id}", d.PemeriksaanHandler.GetByID)
	huma.Put(adminGroup, "/monitoring/pemeriksaan/{id}", d.PemeriksaanHandler.Update)
	huma.Delete(adminGroup, "/monitoring/pemeriksaan/{id}", d.PemeriksaanHandler.Delete)
	huma.Patch(bidanGroup, "/monitoring/pemeriksaan/{id}/verify", d.PemeriksaanHandler.Verify)
	huma.Get(bidanGroup, "/monitoring/pemeriksaan/pending", d.PemeriksaanHandler.GetPending)

	huma.Get(userGroup, "/imunisasi", d.ImunisasiHandler.GetAll)
	huma.Get(userGroup, "/imunisasi/{id}", d.ImunisasiHandler.GetByID)
	huma.Post(adminGroup, "/imunisasi", d.ImunisasiHandler.Create)
	huma.Put(adminGroup, "/imunisasi/{id}", d.ImunisasiHandler.Update)
	huma.Delete(adminGroup, "/imunisasi/{id}", d.ImunisasiHandler.Delete)
	huma.Patch(adminGroup, "/imunisasi/{id}/realisasi", d.ImunisasiHandler.Realisasi)
	huma.Get(userGroup, "/imunisasi/pasien/{id_pasien}", d.ImunisasiHandler.GetByPasienID)
	huma.Get(adminGroup, "/imunisasi/statistik", d.ImunisasiHandler.GetStatistik)

	huma.Get(publicGroup, "/artikel", d.ArtikelHandler.GetAllPublished)
	huma.Get(publicGroup, "/artikel/{id}", d.ArtikelHandler.GetByID)
	huma.Get(adminGroup, "/artikel/semua", d.ArtikelHandler.GetAll)
	huma.Post(adminGroup, "/artikel", d.ArtikelHandler.Create)
	huma.Patch(adminGroup, "/artikel/{id}", d.ArtikelHandler.Update)
	huma.Delete(dinkesGroup, "/artikel/{id}", d.ArtikelHandler.Delete)
	huma.Patch(dinkesGroup, "/artikel/{id}/review", d.ArtikelHandler.Review)
	huma.Get(dinkesGroup, "/artikel/pending", d.ArtikelHandler.GetPending)

	huma.Get(authAccess, "/tindak-lanjut/status", d.TindakLanjutHandler.GetStatusTindakLanjut)
	huma.Get(bidanGroup, "/tindak-lanjut/pasien", d.TindakLanjutHandler.GetPasienTindakLanjut)
	huma.Get(bidanGroup, "/tindak-lanjut/pasien/{id}", d.TindakLanjutHandler.GetDetailPasienByID)
	huma.Post(bidanGroup, "/tindak-lanjut", d.TindakLanjutHandler.CreateTindakLanjut)
	huma.Patch(bidanGroup, "/rujukan/{id}/status", d.TindakLanjutHandler.UpdateStatusRujukan)
	huma.Get(dinkesGroup, "/laporan/tindak-lanjut", d.TindakLanjutHandler.GetLaporanTindakLanjut)
	huma.Get(userGroup, "/tindak-lanjut/{id}", d.TindakLanjutHandler.GetDetailTindakLanjutByID)

	// Dashboard
	huma.Get(userGroup, "/dashboard/stats", d.DashboardHandler.GetDashboardStats)
	huma.Get(userGroup, "/dashboard/distribusi-gizi", d.DashboardHandler.GetDistribusiGizi)
	huma.Get(userGroup, "/dashboard/tren-stunting", d.DashboardHandler.GetTrenStunting)
	huma.Get(userGroup, "/dashboard/stunting-per-wilayah", d.DashboardHandler.GetStuntingPerWilayah)
	huma.Get(userGroup, "/dashboard/kehadiran-bulanan", d.DashboardHandler.GetKehadiranBulanan)
	huma.Get(userGroup, "/dashboard/jadwal-terdekat", d.DashboardHandler.GetJadwalTerdekat)
	huma.Get(publicGroup, "/stats", d.DashboardHandler.GetPublicStats)
	huma.Get(userGroup, "/monitoring/pasien/{id}/riwayat-pemeriksaan", d.DashboardHandler.GetRiwayat)
	huma.Get(userGroup, "/monitoring/pasien/{id}/tumbuh-kembang", d.DashboardHandler.GetTumbuhKembang)

	// Ibu Hamil
	huma.Get(userGroup, "/dashboard/ibu-hamil-stats", d.DashboardHandler.GetIbuHamilStats)
	huma.Get(userGroup, "/dashboard/ibu-hamil-per-wilayah", d.DashboardHandler.GetIbuHamilPerWilayah)
	huma.Get(adminGroup, "/monitoring/semua-pemeriksaan", d.DashboardHandler.GetSemuaPemeriksaan)
}
