package testutils

import (
	"testing"
	"time"

	natsutil "github.com/stringptr/SiGizi/backend/internal/infrastructure/nats"
)

func NewTestNATS(t *testing.T) *natsutil.Conn {
	t.Helper()
	conn, err := natsutil.Connect("nats://nats:4222", "")
	if err != nil {
		t.Fatalf("failed to connect to test NATS: %v", err)
	}
	t.Cleanup(conn.Close)
	return conn
}

func CreateTestNATSBuckets(t *testing.T, natsConn *natsutil.Conn) {
	t.Helper()
	ctx := t.Context()

	_, err := natsConn.CreateKeyValue(ctx, "banned_ips", 15*time.Minute)
	if err != nil {
		t.Fatalf("failed to create banned_ips KV bucket: %v", err)
	}

	_, err = natsConn.CreateKeyValue(ctx, "jwt_blacklist", 30*time.Minute)
	if err != nil {
		t.Fatalf("failed to create jwt_blacklist KV bucket: %v", err)
	}
}
