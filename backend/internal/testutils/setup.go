package testutils

import (
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/go-chi/chi/v5"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/middleware"
)

type Groups struct {
	AuthAccess  *huma.Group
	AuthRefresh *huma.Group
	UserGroup   *huma.Group
	AdminGroup  *huma.Group
	BidanGroup  *huma.Group
	KaderGroup  *huma.Group
	DinkesGroup *huma.Group
	NonAuth     *huma.Group
	PublicGroup *huma.Group
}

func SetupRouter(t *testing.T) (http.Handler, huma.API, *jwtutils.JWT, *NoopBlacklistRepo) {
	t.Helper()
	testHandler, api := humatest.New(t, huma.DefaultConfig("Test", "1.0.0"))
	jwtUtil := jwtutils.New("test-secret-key")
	blacklistRepo := &NoopBlacklistRepo{}

	api.UseMiddleware(middleware.RealIPMiddleware())
	api.UseMiddleware(middleware.AccessTokenMiddleware(api, &jwtUtil, blacklistRepo))

	return testHandler, api, &jwtUtil, blacklistRepo
}

func SetupChiRouter(t *testing.T) (http.Handler, huma.API, *jwtutils.JWT, *NoopBlacklistRepo) {
	t.Helper()
	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("Test", "1.0.0"))
	jwtUtil := jwtutils.New("test-secret-key")
	blacklistRepo := &NoopBlacklistRepo{}

	api.UseMiddleware(middleware.RealIPMiddleware())
	api.UseMiddleware(middleware.AccessTokenMiddleware(api, &jwtUtil, blacklistRepo))

	return router, api, &jwtUtil, blacklistRepo
}

func CreateGroups(api huma.API, jwtUtil *jwtutils.JWT, blacklistRepo *NoopBlacklistRepo) Groups {
	authAccess := huma.NewGroup(api, "")
	authRefresh := huma.NewGroup(api, "")
	userGroup := huma.NewGroup(authAccess, "")
	adminGroup := huma.NewGroup(authAccess, "")
	bidanGroup := huma.NewGroup(authAccess, "")
	kaderGroup := huma.NewGroup(authAccess, "")
	dinkesGroup := huma.NewGroup(authAccess, "")
	publicGroup := huma.NewGroup(api, "")
	nonAuth := huma.NewGroup(api, "")

	authAccess.UseMiddleware(middleware.AuthAccessMiddleware(api, jwtUtil, blacklistRepo))
	authRefresh.UseMiddleware(middleware.AuthRefreshMiddleware(api, jwtUtil))
	userGroup.UseMiddleware(middleware.RequireRole(authAccess, "USER"))
	adminGroup.UseMiddleware(middleware.RequireRole(authAccess, "ADMIN"))
	bidanGroup.UseMiddleware(middleware.RequireRole(authAccess, "BIDAN"))
	kaderGroup.UseMiddleware(middleware.RequireRole(authAccess, "KADER"))
	dinkesGroup.UseMiddleware(middleware.RequireRole(authAccess, "DINKES"))
	nonAuth.UseMiddleware(middleware.NonAuthenticatedOnlyMiddleware(api, jwtUtil, blacklistRepo))

	return Groups{
		AuthAccess:  authAccess,
		AuthRefresh: authRefresh,
		UserGroup:   userGroup,
		AdminGroup:  adminGroup,
		BidanGroup:  bidanGroup,
		KaderGroup:  kaderGroup,
		DinkesGroup: dinkesGroup,
		NonAuth:     nonAuth,
		PublicGroup: publicGroup,
	}
}
