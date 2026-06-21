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

func AccessTokenMiddleware(api huma.API, jwt *jwtutils.JWT, blacklistRepo jwtblacklist.Repo) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		cookie, err := httputils.ReadCookie(ctx, "access_token")
		if err != nil {
			if err == http.ErrNoCookie {
				newCtx := context.WithValue(ctx.Context(), httputils.AccessKey, nil)
				next(huma.WithContext(ctx, newCtx))
			} else {
				huma.WriteErr(api, ctx, http.StatusInternalServerError, "Terjadi kesalahan. Silahkan dicoba kembali.", nil)
			}
			return
		}

		token := cookie.Value
		claims, err := jwt.Decode(token)
		if err != nil {
			newCtx := context.WithValue(ctx.Context(), httputils.AccessKey, nil)
			next(huma.WithContext(ctx, newCtx))
			return
		}

		if claims.ID == "" {
			blacklisted, err := blacklistRepo.IsBlacklisted(ctx.Context(), claims.ID)
			if err == nil && blacklisted {
				huma.WriteErr(api, ctx, http.StatusUnauthorized, "Sesi telah berakhir. Silahkan login ulang.", nil)
				return
			}
			newCtx := context.WithValue(ctx.Context(), httputils.AccessKey, nil)
			next(huma.WithContext(ctx, newCtx))
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

func AuthAccessMiddleware(api huma.API, jwt *jwtutils.JWT, blacklistRepo jwtblacklist.Repo) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		claims := httputils.GetAccessClaim(ctx.Context())

		if claims == nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Silahkan login terlebih dahulu.", nil)
			return
		}

		next(ctx)
	}
}

func NonAuthenticatedOnlyMiddleware(api huma.API, jwt *jwtutils.JWT, blacklistRepo jwtblacklist.Repo) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// Allow re-login — old session will be replaced
		next(ctx)
	}
}

func RequireRole(api huma.API, roles ...string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		claims := httputils.GetAccessClaim(ctx.Context())
		if claims == nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Silahkan login terlebih dahulu.", nil)
			return
		}

		if slices.Contains(claims.Roles, "SUPER_ADMIN") {
			next(ctx)
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
