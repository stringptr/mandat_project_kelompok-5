package notification

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Repo interface {
	GetByID(ctx context.Context, idNotifikasi int32, idUser int32) (*model.Notifikasi, error)
	GetByUserID(ctx context.Context, idUser int32, search string, limit int32, offset int32) ([]*model.Notifikasi, error)
	CountByUserID(ctx context.Context, idUser int32, search string) (int32, error)
	CountUnreadByUserID(ctx context.Context, idUser int32) (int32, error)
	MarkRead(ctx context.Context, idNotifikasi int32, idUser int32) error
	MarkAllRead(ctx context.Context, idUser int32) (int32, error)
	Create(ctx context.Context, data *model.Notifikasi) error
}
