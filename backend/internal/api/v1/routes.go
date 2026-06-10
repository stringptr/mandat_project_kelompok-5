package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/stringptr/SiGizi/backend/internal/config"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	bannedipDomain "github.com/stringptr/SiGizi/backend/internal/domain/bannedip"
	jwtblacklistDomain "github.com/stringptr/SiGizi/backend/internal/domain/jwtblacklist"
	notificationDomain "github.com/stringptr/SiGizi/backend/internal/domain/notification"
	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	userSessionDomain "github.com/stringptr/SiGizi/backend/internal/domain/userSession"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
)

type Dependency struct {
	AuthConfig          config.AuthConfig
	JWTUtil             jwtutils.JWT
	UserAccountRepo     userAccountDomain.Repo
	UserAccountHandler  userAccountDomain.Handler
	UserSessionRepo     userSessionDomain.Repo
	AuthRepo            authDomain.Repo
	AuthService         authDomain.Service
	AuthHandler         authDomain.Handler
	BanRepo             bannedipDomain.Repo
	BanAuthRepo         bannedipDomain.Repo
	BlacklistRepo       jwtblacklistDomain.Repo
	NotifPublisher      notificationDomain.Publisher
	NotifHandler        notificationDomain.Handler
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
	dinkesGroup.UseMiddleware(middleware.RequireRole(authAccess, "ADMIN_DINKES"))

	nonAuthenticatedOnlyGroup := huma.NewGroup(api, "")
	nonAuthenticatedOnlyGroup.UseMiddleware(middleware.NonAuthenticatedOnlyMiddleware(api, &d.JWTUtil, d.BlacklistRepo))

	huma.Post(nonAuthenticatedOnlyGroup, "/auth/register", d.AuthHandler.Register)
	huma.Post(nonAuthenticatedOnlyGroup, "/auth/login", d.AuthHandler.Login)

	huma.Post(authRefresh, "/auth/refresh", d.AuthHandler.Refresh)
	huma.Post(authRefresh, "/auth/logout", d.AuthHandler.Logout)

	huma.Get(userGroup, "/auth/me", d.AuthHandler.Me)

	huma.Patch(adminGroup, "/users/{id_user}/verification", d.AuthHandler.VerifyUser)

	// User Management routes — ADMIN and DINKES can list all users
	usermgtGroup := huma.NewGroup(authAccess, "")
	usermgtGroup.UseMiddleware(middleware.RequireRole(authAccess, "ADMIN", "ADMIN_DINKES"))

	huma.Get(usermgtGroup, "/users", d.UserAccountHandler.GetAllUsers)
	huma.Get(userGroup, "/users/{id}", d.UserAccountHandler.GetUserByID)
	huma.Patch(userGroup, "/users/{id}", d.UserAccountHandler.UpdateUser)

	// Notifikasi routes
	notifBidanKaderDinkesGroup := huma.NewGroup(authAccess, "")
	notifBidanKaderDinkesGroup.UseMiddleware(middleware.RequireRole(authAccess, "BIDAN", "KADER", "ADMIN_DINKES"))

	huma.Get(authAccess, "/notifikasi", d.NotifHandler.GetNotifikasi)
	huma.Get(authAccess, "/notifikasi/{id}", d.NotifHandler.GetNotifikasiDetail)
	huma.Patch(authAccess, "/notifikasi/{id}/read", d.NotifHandler.MarkRead)
	huma.Patch(authAccess, "/notifikasi/read-all", d.NotifHandler.MarkAllRead)
	huma.Get(bidanGroup, "/notifikasi/bidan", d.NotifHandler.GetBidanDashboard)
	huma.Get(notifBidanKaderDinkesGroup, "/notifikasi/statistik", d.NotifHandler.GetStatistics)
	huma.Get(notifBidanKaderDinkesGroup, "/notifikasi/aktivitas", d.NotifHandler.GetActivity)
}
