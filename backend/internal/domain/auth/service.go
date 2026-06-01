package auth

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service interface {
	Register(ctx context.Context, dataDTO *RegisterRequest) *errorutils.Error
	Login(ctx context.Context, req *LoginRequest, ip string) (*AuthResponse, *errorutils.Error)
	Refresh(ctx context.Context, refreshToken uuid.UUID, ip string) (*AuthResponse, *errorutils.Error)
	Logout(ctx context.Context, refreshToken uuid.UUID) *errorutils.Error
}
