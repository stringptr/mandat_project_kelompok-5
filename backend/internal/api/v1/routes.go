package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/stringptr/SiGizi/backend/internal/config"
	"github.com/stringptr/SiGizi/backend/internal/feature/auth"
	"github.com/stringptr/SiGizi/backend/internal/feature/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/feature/userSession"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(api huma.API, r chi.Router, pool *pgxpool.Pool, cfg *config.Config) {
	userAccountRepo := userAccount.NewRepo(pool)
	userAccountService := userAccount.NewService(userAccountRepo)
	userAccountHandler := userAccount.NewHandler(userAccountService)

	jwtUtil := jwtutils.New(cfg.AuthConfig.JWTSecret)
	userSessionRepo := userSession.NewRepo(pool)

	authRepo := auth.NewRepo(pool)
	authService := auth.NewService(authRepo, userSessionRepo, userAccountRepo, jwtUtil, &cfg.AuthConfig)
	authHandler := auth.NewHandler(authService)

	huma.Post(api, "/auth/register", authHandler.Register)

	authAccess := huma.NewGroup(api, "")
	authRefresh := huma.NewGroup(api, "")
	userGroup := huma.NewGroup(authAccess, "")
	adminGroup := huma.NewGroup(userGroup, "")
	bidanGroup := huma.NewGroup(adminGroup, "")
	kaderGroup := huma.NewGroup(adminGroup, "")
	dinkesGroup := huma.NewGroup(userGroup, "")

	authAccess.UseMiddleware(middleware.AuthAccessMiddleware(api, &jwtUtil))
	authRefresh.UseMiddleware(middleware.AuthRefreshMiddleware(api))
	userGroup.UseMiddleware(middleware.RequireRole(authAccess, "USER"))
	adminGroup.UseMiddleware(middleware.RequireRole(authAccess, "ADMIN"))
	bidanGroup.UseMiddleware(middleware.RequireRole(authAccess, "BIDAN"))
	kaderGroup.UseMiddleware(middleware.RequireRole(authAccess, "KADER"))
	dinkesGroup.UseMiddleware(middleware.RequireRole(authAccess, "DINKES"))
}
