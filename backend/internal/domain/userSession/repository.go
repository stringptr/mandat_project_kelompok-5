package userSession

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Repo interface {
	GetByID(ctx context.Context, sessionID uuid.UUID) (*model.UserSession, error)
	Create(ctx context.Context, sessionModel *model.UserSession) error
	Update(ctx context.Context, sessionModel *model.UserSession) error
}
