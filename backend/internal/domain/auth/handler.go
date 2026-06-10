package auth

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

type Handler interface {
	Register(ctx context.Context, input *httputils.APIRequestInput[*RegisterRequest]) (*httputils.APIResponseOutput[any], error)
	Me(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*jwtutils.Claim], error)
	Login(ctx context.Context, input *httputils.APIRequestInput[*LoginRequest]) (*AuthOutput, error)
	Refresh(ctx context.Context, input *struct{}) (*AuthOutput, error)
	Logout(ctx context.Context, input *struct{}) (*LogoutOutput, error)
	VerifyUser(ctx context.Context, input *httputils.APIRequestInput[*VerifyUserRequest]) (*httputils.APIResponseOutput[any], error)
}
