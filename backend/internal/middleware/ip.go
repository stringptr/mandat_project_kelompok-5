package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stringptr/SiGizi/backend/internal/domain/bannedip"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

func RealIPMiddleware() func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ip := ctx.RemoteAddr()
		next(huma.WithContext(ctx, httputils.WithRealIP(ctx.Context(), ip)))
	}
}

func IPBanMiddleware(api huma.API, banRepo bannedip.Repo, message string, err *errorutils.Error) func(ctx huma.Context, next func(huma.Context)) {
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
			huma.WriteErr(api, ctx, http.StatusForbidden, message, err)
			return
		}

		next(ctx)
	}
}

func AuthRestrictMiddleware(api huma.API, banRepo bannedip.Repo) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ip := httputils.GetRealIP(ctx.Context())
		info, _ := banRepo.GetBanInfo(ctx.Context(), ip)
		remaining := time.Until(info.ExpiresAt)
		IPBanMiddleware(api, banRepo, "Akses ditolak. Terlalu banyak percobaan login. Silahkan coba lagi dalam %d menit %d detik.", &errorutils.Error{Status: http.StatusForbidden, Message: fmt.Sprintf("Akses ditolak. Terlalu banyak percobaan. Silahkan coba lagi dalam %d menit %d detik.", int(remaining.Minutes()), int(remaining.Seconds())%60)})
	}
}
