package auth

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type Service interface {
	Register(ctx context.Context, dataDTO *RegisterRequest) error
	Login(ctx context.Context, req *LoginRequest, ip string) (*AuthResponse, error)
	Refresh(ctx context.Context, refreshToken uuid.UUID, ip string) (*AuthResponse, error)
	Logout(ctx context.Context, refreshToken uuid.UUID) error
}
