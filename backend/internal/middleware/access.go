package middleware

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

func AuthAccessMiddleware(api huma.API, jwt *jwtutils.JWT) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		cookie, err := httputils.ReadCookie(ctx, "access_token")
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("no access token provided"))
			return
		}

		token := cookie.Value
		claims, err := jwt.Decode(token)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("invalid token: %w", err))
			return
		}

		newCtx := context.WithValue(ctx.Context(), httputils.AccessKey, claims)
		next(huma.WithContext(ctx, newCtx))
	}
}

func RequireRole(api huma.API, roles ...string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		claims := httputils.GetAccessClaim(ctx.Context())
		if claims == nil {
			huma.WriteErr(api, ctx, http.StatusForbidden, "Forbidden", fmt.Errorf("no claims in context"))
			return
		}

		for _, required := range roles {
			if slices.Contains(claims.Roles, required) {
				next(ctx)
				return
			}
		}
		huma.WriteErr(api, ctx, http.StatusForbidden, "Forbidden", fmt.Errorf("insufficient role"))
	}
}
