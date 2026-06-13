package auditlog

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type AuditLogFilter struct {
	Page    int
	PerPage int
}

type Repo interface {
	Log(ctx context.Context, entry *model.AuditLog) error
	GetAll(ctx context.Context, filter *AuditLogFilter) ([]*model.AuditLog, int, error)
}
