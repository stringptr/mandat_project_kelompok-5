package auth

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Repo interface {
	GetRoles(ctx context.Context, idUser int32) ([]string, error)
	GetByEmailNIK(ctx context.Context, email string, nik string) (*model.UserAccount, error)
	CreateUser(ctx context.Context, data *model.UserAccount) (int32, error)
	CreateRoleRecord(ctx context.Context, idUser int32, role string, noStr string, wilayahKerja int32, noSk string) error
}
