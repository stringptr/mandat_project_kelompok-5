package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

type Handler struct {
	Service authDomain.Service
	jwt     *jwtutils.JWT
}

func NewHandler(service authDomain.Service, jwt *jwtutils.JWT) *Handler {
	return &Handler{Service: service, jwt: jwt}
}

func (h *Handler) Register(ctx context.Context, input *httputils.APIRequestInput[*authDomain.RegisterRequest]) (*httputils.APIResponseOutput[any], error) {
	ip := httputils.GetRealIP(ctx)
	err := h.Service.Register(ctx, input.Body, ip)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewSuccessOutput[any](201, nil, "Register berhasil. Akun sedang diverifikasi. Silahkan dicek secara berkala."), nil
}

func (h *Handler) Me(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*jwtutils.Claim], error) {
	accessCookie := httputils.GetAccessClaim(ctx)
	if accessCookie == nil {
		return nil, huma.Error401Unauthorized("Silahkan login terlebih dahulu.", nil)
	}

	return httputils.NewOKOutput(accessCookie), nil
}

func (h *Handler) Login(ctx context.Context, input *httputils.APIRequestInput[*authDomain.LoginRequest]) (*authDomain.AuthOutput, error) {
	ip := httputils.GetRealIP(ctx)
	res, err := h.Service.Login(ctx, input.Body, ip)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	var returnCookie []http.Cookie
	returnCookie = append(returnCookie,
		http.Cookie{
			Name:     "access_token",
			Value:    res.AccessToken,
			Expires:  time.Now().Add(time.Duration(res.AccessTokenExpiresIn) * time.Second),
			MaxAge:   int(res.AccessTokenExpiresIn),
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
		http.Cookie{
			Name:     "refresh_token",
			Value:    res.RefreshToken.String(),
			Expires:  time.Now().Add(time.Duration(res.RefreshTokenExpiresIn) * time.Second),
			MaxAge:   int(res.RefreshTokenExpiresIn),
			Path:     "api/v1/auth/refresh",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	)

	return &authDomain.AuthOutput{
		Body:      httputils.OK(res),
		SetCookie: returnCookie,
	}, nil
}

func (h *Handler) Refresh(ctx context.Context, input *struct{}) (*authDomain.AuthOutput, error) {
	ip := httputils.GetRealIP(ctx)
	refreshToken := httputils.GetRefreshToken(ctx)

	res, err := h.Service.Refresh(ctx, refreshToken, ip)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	var returnCookie []http.Cookie
	returnCookie = append(returnCookie,
		http.Cookie{
			Name:     "access_token",
			Value:    res.AccessToken,
			Expires:  time.Now().Add(time.Duration(res.AccessTokenExpiresIn) * time.Second),
			MaxAge:   int(res.AccessTokenExpiresIn),
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
		http.Cookie{
			Name:     "refresh_token",
			Value:    res.RefreshToken.String(),
			Expires:  time.Now().Add(time.Duration(res.RefreshTokenExpiresIn) * time.Second),
			MaxAge:   int(res.RefreshTokenExpiresIn),
			Path:     "/api/v1/auth/refresh",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	)

	return &authDomain.AuthOutput{
		Body:      httputils.OK(res),
		SetCookie: returnCookie,
	}, nil
}

func (h *Handler) Logout(ctx context.Context, input *struct{}) (*authDomain.LogoutOutput, error) {
	refreshToken := httputils.GetRefreshToken(ctx)
	jti := httputils.GetAccessTokenJTI(ctx)

	err := h.Service.Logout(ctx, refreshToken, jti)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	var returnCookie []http.Cookie
	returnCookie = append(returnCookie,
		http.Cookie{
			Name:     "access_token",
			Value:    "",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
		http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			Path:     "/api/v1/auth/refresh",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	)
	return &authDomain.LogoutOutput{Body: httputils.Success[any](http.StatusOK, nil, "Logout berhasil."), SetCookie: returnCookie}, nil
}

func (h *Handler) VerifyUser(ctx context.Context, input *authDomain.VerifyUserInput) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.VerifyUser(ctx, &authDomain.VerifyUserRequest{
		IDUser:          input.IDUser,
		Status:          input.Body.Status,
		AlasanPenolakan: input.Body.AlasanPenolakan,
	})
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput[any](nil), nil
}
