package auditlog

import (
	"context"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/pgxV5"
	"github.com/jackc/pgx/v5/pgxpool"
	auditlogDomain "github.com/stringptr/SiGizi/backend/internal/domain/auditlog"
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

func (r *Repo) GetAll(ctx context.Context, filter *auditlogDomain.AuditLogFilter) ([]*model.AuditLog, int, error) {
	offset := (filter.Page - 1) * filter.PerPage

	var total int
	countQuery := "SELECT COUNT(*) FROM audit_log"
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf(`
		SELECT id_log, tipe_aktor, id_user, id_user_session, tipe_aktivitas,
		       berhasil, endpoint, table_name, record_id, old_value, new_value,
		       detail, ip_address, user_agent, waktu_aktivitas
		FROM audit_log
		ORDER BY waktu_aktivitas DESC
		LIMIT $1 OFFSET $2
	`)
	args := []any{filter.PerPage, offset}

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*model.AuditLog
	for rows.Next() {
		var l model.AuditLog
		err := rows.Scan(
			&l.IDLog, &l.TipeAktor, &l.IDUser, &l.IDUserSession, &l.TipeAktivitas,
			&l.Berhasil, &l.Endpoint, &l.TableName, &l.RecordID, &l.OldValue, &l.NewValue,
			&l.Detail, &l.IPAddress, &l.UserAgent, &l.WaktuAktivitas,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, &l)
	}

	return logs, total, nil
}
