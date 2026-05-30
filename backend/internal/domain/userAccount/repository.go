package userAccount

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Repo interface {
	GetByID(ctx context.Context, IDUser int32) (*model.UserAccount, error)
	GetByNIK(ctx context.Context, NIK string) (*model.UserAccount, error)
	GetByEmail(ctx context.Context, email string) (*model.UserAccount, error)
	GetAll(ctx context.Context) ([]*model.UserAccount, error)
	Create(ctx context.Context, dataModel *model.UserAccount) error
	UpdateStatusVerifikasiByID(ctx context.Context, IDUser int32, statusVerifikasi bool) error
	DeleteByID(ctx context.Context, IDUser int32) error
}
