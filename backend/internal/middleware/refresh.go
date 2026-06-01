package middleware

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

func AuthRefreshMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		refreshCookie, err := httputils.ReadCookie(ctx, "refresh_token")
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Silahkan login terlebih dahulu.", nil)
			return
		}

		refreshToken, err := uuid.FromString(refreshCookie.Value)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusInternalServerError, "Terjadi kesalahan. Silahkan dicoba kembali.", nil)
			return
		}
		newCtx := context.WithValue(ctx.Context(), httputils.RefreshKey, refreshToken)

		next(huma.WithContext(ctx, newCtx))
	}
}
