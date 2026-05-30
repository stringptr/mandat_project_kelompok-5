package httputils

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

type key string

const (
	AccessKey  key = "access_token"
	RefreshKey key = "refresh_token"
)

func GetAccessClaim(ctx context.Context) *jwtutils.Claim {
	claims, _ := ctx.Value(AccessKey).(*jwtutils.Claim)
	return claims
}

func GetRefreshToken(ctx context.Context) uuid.UUID {
	claims, _ := ctx.Value(RefreshKey).(uuid.UUID)
	return claims
}

func ReadCookie(ctx huma.Context, name string) (*http.Cookie, error) {
	c := ctx.Header("Cookie")
	if c == "" {
		return nil, http.ErrNoCookie
	}
	req := http.Request{Header: http.Header{"Cookie": {c}}}
	return req.Cookie(name)
}
