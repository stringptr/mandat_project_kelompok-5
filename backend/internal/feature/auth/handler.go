package auth

import (
	"context"
	"errors"
	"github.com/danielgtaylor/huma/v2"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

type AuthOutput struct {
	Body             httputils.APIResponse[*authDomain.AuthResponse]
	SetAccessCookie  http.Cookie `header:"Set-Cookie"`
	SetRefreshCookie http.Cookie `header:"Set-Cookie"`
}

func (h *Handler) Register(ctx context.Context, input *httputils.APIRequestInput[*authDomain.RegisterRequest]) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.Register(ctx, input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest("Registration Failed", err)
	}

	return &httputils.APIResponseOutput[any]{Body: httputils.Created[any](nil)}, nil
}

func (h *Handler) Me(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*jwtutils.Claim], error) {
	accessCookie := httputils.GetAccessClaim(ctx)
	if accessCookie == nil {
		return nil, huma.Error401Unauthorized("Please login first.", errors.New("no access token"))
	}

	return &httputils.APIResponseOutput[*jwtutils.Claim]{Body: httputils.OK(accessCookie)}, nil
}

func (h *Handler) Login(ctx context.Context, input *httputils.APIRequestInput[*authDomain.LoginRequest]) (*AuthOutput, error) {
	ip := httputils.GetRealIP(ctx)
	resp, err := h.Service.Login(ctx, input.Body, ip)
	if err != nil {
		return nil, huma.Error401Unauthorized(err.Error(), err)
	}

	return &AuthOutput{
		Body: httputils.OK(resp),
		SetAccessCookie: http.Cookie{
			Name:     "access_token",
			Value:    resp.AccessToken,
			Expires:  time.Now().Add(time.Duration(resp.AccessTokenExpiresIn) * time.Second),
			MaxAge:   int(resp.AccessTokenExpiresIn),
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
		SetRefreshCookie: http.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			Expires:  time.Now().Add(time.Duration(resp.RefreshTokenExpiresIn) * time.Second),
			MaxAge:   int(resp.RefreshTokenExpiresIn),
			Path:     "/auth/refresh",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	}, nil
}
