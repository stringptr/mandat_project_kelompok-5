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
	AuthConfig      config.AuthConfig
	JWTUtil         jwtutils.JWT
	UserAccountRepo userAccountDomain.Repo
	UserSessionRepo userSessionDomain.Repo
	AuthRepo        authDomain.Repo
	AuthService     authDomain.Service
	AuthHandler     authDomain.Handler
	BanRepo         bannedipDomain.Repo
	BanAuthRepo     bannedipDomain.Repo
	BlacklistRepo   jwtblacklistDomain.Repo
	NotifPublisher  notificationDomain.Publisher
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

	huma.Get(userGroup, "/me", d.AuthHandler.Me)
}
