package userAccount

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Repo interface {
	GetByID(ctx context.Context, IDUser int32) (*model.UserAccount, error)
	GetByNIK(ctx context.Context, NIK string) (*model.UserAccount, error)
	GetByNIKEmail(ctx context.Context, NIK string, email string) (*model.UserAccount, error)
	GetByEmail(ctx context.Context, email string) (*model.UserAccount, error)
	GetAll(ctx context.Context) ([]*model.UserAccount, error)
	GetAllPaginated(ctx context.Context, page int, perPage int, q string, role string, statusVerifikasi string) ([]*model.UserAccount, int, error)
	Create(ctx context.Context, dataModel *model.UserAccount) error
	Update(ctx context.Context, dataModel *model.UserAccount) error
	UpdateStatusVerifikasi(ctx context.Context, IDUser int32, status model.StatusVerifikasi) error
	DeleteByID(ctx context.Context, IDUser int32) error
	GetLokasiNames(ctx context.Context, ids map[int32]bool) (map[int32]string, error)
}
