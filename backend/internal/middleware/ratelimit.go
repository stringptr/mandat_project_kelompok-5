package middleware

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	bannedipDomain "github.com/stringptr/SiGizi/backend/internal/domain/bannedip"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

func BannedIPMiddleware(api huma.API, bannedIPRepo bannedipDomain.Repo) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ip := httputils.GetRealIP(ctx.Context())
		if ip == "" {
			next(ctx)
			return
		}

		banned, err := bannedIPRepo.IsBanned(ctx.Context(), ip)
		if err != nil || banned {
			huma.WriteErr(api, ctx, http.StatusForbidden, "Akses ditolak. IP Anda sedang diblokir sementara.", nil)
			return
		}

		next(ctx)
	}
}
