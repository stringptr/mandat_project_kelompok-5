package middleware

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stringptr/SiGizi/backend/internal/domain/bannedip"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

func RealIPMiddleware() func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ip := ctx.RemoteAddr()
		next(huma.WithContext(ctx, httputils.WithRealIP(ctx.Context(), ip)))
	}
}

func IPBanMiddleware(api huma.API, banRepo bannedip.Repo) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ip := httputils.GetRealIP(ctx.Context())
		if ip == "" {
			next(ctx)
			return
		}

		banned, err := banRepo.IsBanned(ctx.Context(), ip)
		if err != nil {
			next(ctx)
			return
		}
		if banned {
			huma.WriteErr(api, ctx, http.StatusForbidden, "Akses ditolak. IP Anda diblokir sementara.", nil)
			return
		}

		next(ctx)
	}
}
