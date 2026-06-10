package testutils

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stringptr/SiGizi/backend/internal/config"
)

func NewTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	cfg := config.Load()

	poolCfg, err := pgxpool.ParseConfig(cfg.DBMasterConfig.DSN())
	if err != nil {
		t.Fatalf("failed to parse DB config: %v", err)
	}
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterType(&pgtype.Type{
			Name:  "inet",
			OID:   pgtype.InetOID,
			Codec: &pgtype.TextCodec{},
		})
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		t.Fatalf("failed to connect to test DB: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func TruncateNotifikasiTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), "TRUNCATE TABLE notifikasi CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate notifikasi: %v", err)
	}
}

func TruncateAuthTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()

	tables := []string{
		"audit_log",
		"user_session",
		"bidan",
		"dinas_kesehatan",
		"kader_posyandu",
		"pasien",
		"user_account",
		"lokasi",
	}

	for _, table := range tables {
		_, err := pool.Exec(ctx, "TRUNCATE TABLE "+table+" CASCADE")
		if err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}

	var seqs = []string{
		"audit_log_id_log_seq",
		"lokasi_id_lokasi_seq",
		"user_account_id_user_seq",
	}
	for _, seq := range seqs {
		_, err := pool.Exec(ctx, "ALTER SEQUENCE "+seq+" RESTART WITH 1")
		if err != nil {
			_ = err
		}
	}
}
