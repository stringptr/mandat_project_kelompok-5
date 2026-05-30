package auth

import (
	"context"
)

type Repo interface {
	GetRoles(ctx context.Context, idUser int32) ([]string, error)
}
