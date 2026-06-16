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
	JTIKey     key = "access_token_jti"
)

func GetAccessTokenJTI(ctx context.Context) string {
	jti, _ := ctx.Value(JTIKey).(string)
	return jti
}

func WithAccessTokenJTI(ctx context.Context, jti string) context.Context {
	return context.WithValue(ctx, JTIKey, jti)
}

func GetAccessClaim(ctx context.Context) *jwtutils.Claim {
	claims, _ := ctx.Value(AccessKey).(*jwtutils.Claim)
	return claims
}

func GetRefreshToken(ctx context.Context) uuid.UUID {
	claimsRaw := ctx.Value(RefreshKey)
	if claims, ok := claimsRaw.(uuid.UUID); ok {
		return claims
	}
	return uuid.Nil
}

func IsPetugas(roles []string) bool {
	for _, r := range roles {
		switch r {
		case "ADMIN", "BIDAN", "KADER", "DINKES", "SUPER_ADMIN":
			return true
		}
	}
	return false
}

func ReadCookie(ctx huma.Context, name string) (*http.Cookie, error) {
	c := ctx.Header("Cookie")
	if c == "" {
		return nil, http.ErrNoCookie
	}
	req := http.Request{Header: http.Header{"Cookie": {c}}}
	return req.Cookie(name)
}
