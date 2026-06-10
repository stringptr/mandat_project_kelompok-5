package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/stringptr/SiGizi/backend/internal/config"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	jwtblacklistDomain "github.com/stringptr/SiGizi/backend/internal/domain/jwtblacklist"
	notificationDomain "github.com/stringptr/SiGizi/backend/internal/domain/notification"
	pasienDomain "github.com/stringptr/SiGizi/backend/internal/domain/pasien"
	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
)

type Dependency struct {
	AuthConfig         config.AuthConfig
	JWTUtil            jwtutils.JWT
	UserAccountHandler userAccountDomain.Handler
	AuthHandler        authDomain.Handler
	BlacklistRepo      jwtblacklistDomain.Repo
	NotifPublisher     notificationDomain.Publisher
	NotifHandler       notificationDomain.Handler
	PasienHandler      pasienDomain.Handler
}

func RegisterRoutes(api huma.API, r chi.Router, d *Dependency) {
	authAccess := huma.NewGroup(api, "")
	authRefresh := huma.NewGroup(api, "")
	userGroup := huma.NewGroup(authAccess, "")
	adminGroup := huma.NewGroup(userGroup, "")
	bidanGroup := huma.NewGroup(adminGroup, "")
	kaderGroup := huma.NewGroup(adminGroup, "")
	dinkesGroup := huma.NewGroup(userGroup, "")

	authAccess.UseMiddleware(middleware.AuthAccessMiddleware(api, &d.JWTUtil, d.BlacklistRepo))
	authRefresh.UseMiddleware(middleware.AuthRefreshMiddleware(api, &d.JWTUtil))
	userGroup.UseMiddleware(middleware.RequireRole(authAccess, "USER"))
	adminGroup.UseMiddleware(middleware.RequireRole(authAccess, "ADMIN"))
	bidanGroup.UseMiddleware(middleware.RequireRole(authAccess, "BIDAN"))
	kaderGroup.UseMiddleware(middleware.RequireRole(authAccess, "KADER"))
	dinkesGroup.UseMiddleware(middleware.RequireRole(authAccess, "DINKES"))

	nonAuthenticatedOnlyGroup := huma.NewGroup(api, "")
	nonAuthenticatedOnlyGroup.UseMiddleware(middleware.NonAuthenticatedOnlyMiddleware(api, &d.JWTUtil, d.BlacklistRepo))

	huma.Post(nonAuthenticatedOnlyGroup, "/auth/register", d.AuthHandler.Register)
	huma.Post(nonAuthenticatedOnlyGroup, "/auth/login", d.AuthHandler.Login)

	huma.Post(authRefresh, "/auth/refresh", d.AuthHandler.Refresh)
	huma.Post(authRefresh, "/auth/logout", d.AuthHandler.Logout)

	huma.Get(userGroup, "/auth/me", d.AuthHandler.Me)

	huma.Patch(adminGroup, "/users/{id_user}/verification", d.AuthHandler.VerifyUser)

	huma.Get(adminGroup, "/users", d.UserAccountHandler.GetAllUsers)
	huma.Get(userGroup, "/users/{id}", d.UserAccountHandler.GetUserByID)
	huma.Patch(userGroup, "/users/{id}", d.UserAccountHandler.UpdateUser)

	huma.Get(userGroup, "/notifikasi", d.NotifHandler.GetNotifikasi)
	huma.Get(userGroup, "/notifikasi/{id}", d.NotifHandler.GetNotifikasiDetail)
	huma.Patch(userGroup, "/notifikasi/{id}/read", d.NotifHandler.MarkRead)
	huma.Patch(userGroup, "/notifikasi/read-all", d.NotifHandler.MarkAllRead)
	huma.Get(bidanGroup, "/notifikasi/bidan", d.NotifHandler.GetBidanDashboard)
	huma.Get(adminGroup, "/notifikasi/statistik", d.NotifHandler.GetStatistics)
	huma.Get(adminGroup, "/notifikasi/aktivitas", d.NotifHandler.GetActivity)

	huma.Post(adminGroup, "/pasien/ibu-hamil", d.PasienHandler.DaftarIbuHamil)
	huma.Post(adminGroup, "/pasien/anak", d.PasienHandler.DaftarAnak)
	huma.Get(adminGroup, "/monitoring/pasien", d.PasienHandler.GetAll)
	huma.Get(adminGroup, "/monitoring/pasien/search", d.PasienHandler.Search)
	huma.Get(adminGroup, "/monitoring/pasien/{id}", d.PasienHandler.GetByID)
	huma.Patch(bidanGroup, "/pasien/{id}", d.PasienHandler.Update)
	huma.Delete(bidanGroup, "/pasien/{id}", d.PasienHandler.Delete)
}
