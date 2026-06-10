package auditlog

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/pgxV5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	. "github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/table"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Log(ctx context.Context, entry *model.AuditLog) error {
	claims := httputils.GetAccessClaim(ctx)
	if claims != nil {
		entry.IDUser = &claims.IDUser
	}
	ip := httputils.GetRealIP(ctx)
	if ip != "" {
		entry.IPAddress = &ip
	}
	entry.WaktuAktivitas = time.Now()

	stmt := AuditLog.INSERT(AuditLog.MutableColumns).
		MODEL(entry)
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}
