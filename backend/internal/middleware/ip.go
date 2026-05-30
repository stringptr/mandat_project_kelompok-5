package middleware

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

func RealIPMiddleware() func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ip := ctx.RemoteAddr()
		next(huma.WithContext(ctx, httputils.WithRealIP(ctx.Context(), ip)))
	}
}
