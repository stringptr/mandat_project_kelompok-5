package auditlog

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Repo interface {
	Log(ctx context.Context, entry *model.AuditLog) error
}
