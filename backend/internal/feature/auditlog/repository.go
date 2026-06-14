package auditlog

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/pgxV5"
	"github.com/jackc/pgx/v5/pgxpool"
	auditlogDomain "github.com/stringptr/SiGizi/backend/internal/domain/auditlog"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	. "github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/table"

	. "github.com/go-jet/jet/v2/postgres"
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

func (r *Repo) GetAll(ctx context.Context, filter *auditlogDomain.AuditLogFilter) ([]*model.AuditLog, int, error) {
	offset := (filter.Page - 1) * filter.PerPage

	var count struct {
		Count int32
	}
	countStmt := SELECT(COUNT(STAR).AS("count")).
		FROM(AuditLog)
	err := pgxV5.Query(ctx, countStmt, r.db, &count)
	if err != nil {
		return nil, 0, err
	}

	dataStmt := SELECT(AuditLog.AllColumns).
		FROM(AuditLog).
		ORDER_BY(AuditLog.WaktuAktivitas.DESC()).
		LIMIT(int64(filter.PerPage)).
		OFFSET(int64(offset))

	var logs []*model.AuditLog
	err = pgxV5.Query(ctx, dataStmt, r.db, &logs)
	if err != nil {
		return nil, 0, err
	}

	return logs, int(count.Count), nil
}
