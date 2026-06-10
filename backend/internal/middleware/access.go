package middleware

import (
	"context"
	"net/http"
	"slices"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stringptr/SiGizi/backend/internal/domain/jwtblacklist"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

func AuthAccessMiddleware(api huma.API, jwt *jwtutils.JWT, blacklistRepo jwtblacklist.Repo) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		cookie, err := httputils.ReadCookie(ctx, "access_token")
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Silahkan login terlebih dahulu.", nil)
			return
		}

		token := cookie.Value
		claims, err := jwt.Decode(token)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusInternalServerError, "Terjadi kesalahan. Silahkan dicoba kembali.", nil)
			return
		}

		if claims.ID != "" {
			blacklisted, err := blacklistRepo.IsBlacklisted(ctx.Context(), claims.ID)
			if err == nil && blacklisted {
				huma.WriteErr(api, ctx, http.StatusUnauthorized, "Sesi telah berakhir. Silahkan login ulang.", nil)
				return
			}
		}

		newCtx := context.WithValue(ctx.Context(), httputils.AccessKey, claims)
		next(huma.WithContext(ctx, newCtx))
	}
}

func AuthRequiredAccessMiddleware(api huma.API, jwt *jwtutils.JWT, blacklistRepo jwtblacklist.Repo) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		claims := httputils.GetAccessClaim(ctx.Context())

		if claims == nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Silahkan login terlebih dahulu.", nil)
			return
		}

		next(ctx)
	}
}

func NotLoggedInRequiredMiddleware(api huma.API, jwt *jwtutils.JWT, blacklistRepo jwtblacklist.Repo) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		claims := httputils.GetAccessClaim(ctx.Context())

		if claims == nil {
			next(ctx)
		}

		huma.WriteErr(api, ctx, http.StatusUnauthorized, "Anda sudah loginl. Silahkan logout terlebih dahulu.", nil)
	}
}

func RequireRole(api huma.API, roles ...string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		claims := httputils.GetAccessClaim(ctx.Context())
		if claims == nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Silahkan login terlebih dahulu.", nil)
			return
		}

		for _, required := range roles {
			if slices.Contains(claims.Roles, required) {
				next(ctx)
				return
			}
		}
		huma.WriteErr(api, ctx, http.StatusForbidden, "Tidak mempunyai akses untuk halaman ini.", nil)
	}
}
