package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

func AuthRefreshMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		refreshCookie, err := httputils.ReadCookie(ctx, "refresh_token")
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("no refresh token provided"))
			return
		}

		refreshToken, err := uuid.FromString(refreshCookie.Value)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("invalid refresh token"))
			return
		}
		newCtx := context.WithValue(ctx.Context(), httputils.RefreshKey, refreshToken)

		next(huma.WithContext(ctx, newCtx))
	}
}
