//go:build mailtrap_integration

package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	gomail "github.com/go-mail/mail/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/stringptr/SiGizi/backend/internal/config"
	"github.com/stringptr/SiGizi/backend/internal/feature/auditlog"
	"github.com/stringptr/SiGizi/backend/internal/feature/bannedip"
	"github.com/stringptr/SiGizi/backend/internal/feature/jwtblacklist"
	"github.com/stringptr/SiGizi/backend/internal/feature/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/feature/userSession"
	natsutil "github.com/stringptr/SiGizi/backend/internal/infrastructure/nats"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/middleware"
	"github.com/stringptr/SiGizi/backend/internal/testutils"
)

type timedSender struct {
	dialer *gomail.Dialer
	from   string
	name   string
}

func (s *timedSender) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", s.from, s.name)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return s.dialer.DialAndSend(m)
}

type mailtrapTestFixture struct {
	handler       http.Handler
	jwtUtil       *jwtutils.JWT
	pool          *pgxpool.Pool
	blacklistRepo *jwtblacklist.Repo
	banRepo       *bannedip.Repo
	bannedKV      jetstream.KeyValue
	blacklistKV   jetstream.KeyValue
}

func setupMailtrapTest(t *testing.T) *mailtrapTestFixture {
	t.Helper()

	pool := testutils.NewTestDB(t)
	natsConn := testutils.NewTestNATS(t)

	ctx := context.Background()

	resetKV := func(bucket string, ttl time.Duration) jetstream.KeyValue {
		_ = natsConn.JetStream().DeleteKeyValue(ctx, bucket)
		kv, err := natsConn.CreateKeyValue(ctx, bucket, ttl)
		if err != nil {
			t.Fatalf("failed to create KV bucket %s: %v", bucket, err)
		}
		return kv
	}

	bannedKV := resetKV("banned_ips", 15*time.Minute)
	blacklistKV := resetKV("jwt_blacklist", 30*time.Minute)

	jwtUtil := jwtutils.New("test-secret-key")

	authCfg := &config.AuthConfig{
		JWTSecret:       "test-secret-key",
		AccessTokenTTL:  30 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}
	restrictCfg := &config.RestrictAuthConfig{
		MaxAttempt: 3,
		Duration:   1 * time.Hour,
	}

	host := getEnv("MAIL_HOST", "sandbox.smtp.mailtrap.io")
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	fromEmail := getEnv("MAIL_FROM_EMAIL", "noreply@sigizi.com")
	fromName := getEnv("MAIL_FROM_NAME", "SiGizi")

	d := gomail.NewDialer(host, 2525, username, password)
	d.Timeout = 15 * time.Second
	mailSender := &timedSender{dialer: d, from: fromEmail, name: fromName}

	authRepo := NewRepo(pool)
	userSessionRepo := userSession.NewRepo(pool)
	userAccountRepo := userAccount.NewRepo(pool)
	banRepo := bannedip.NewRepo(natsutil.NewKV(bannedKV))
	br := jwtblacklist.NewRepo(natsutil.NewKV(blacklistKV))
	auditLogRepo := auditlog.NewRepo(pool)

	svc := NewService(authRepo, userSessionRepo, userAccountRepo, jwtUtil, authCfg, restrictCfg, banRepo, br, auditLogRepo, mailSender)
	h := NewHandler(svc, &jwtUtil)

	testHandler, api := humatest.New(t, huma.DefaultConfig("Test", "1.0.0"))

	api.UseMiddleware(middleware.RealIPMiddleware())
	api.UseMiddleware(middleware.AccessTokenMiddleware(api, &jwtUtil, br))

	authAccess := huma.NewGroup(api, "")
	authRefresh := huma.NewGroup(api, "")
	userGroup := huma.NewGroup(authAccess, "")
	adminGroup := huma.NewGroup(authAccess, "")

	authAccess.UseMiddleware(middleware.AuthAccessMiddleware(api, &jwtUtil, br))
	authRefresh.UseMiddleware(middleware.AuthRefreshMiddleware(api, &jwtUtil))
	userGroup.UseMiddleware(middleware.RequireRole(authAccess, "USER"))
	adminGroup.UseMiddleware(middleware.RequireRole(authAccess, "ADMIN"))
	nonAuth := huma.NewGroup(api, "")
	nonAuth.UseMiddleware(middleware.NonAuthenticatedOnlyMiddleware(api, &jwtUtil, br))

	huma.Post(nonAuth, "/auth/register", h.Register)
	huma.Post(nonAuth, "/auth/login", h.Login)
	huma.Post(authRefresh, "/auth/refresh", h.Refresh)
	huma.Post(authRefresh, "/auth/logout", h.Logout)
	huma.Get(userGroup, "/auth/me", h.Me)
	huma.Patch(adminGroup, "/users/{id_user}/verification", h.VerifyUser)

	return &mailtrapTestFixture{
		handler:       testHandler,
		jwtUtil:       &jwtUtil,
		pool:          pool,
		blacklistRepo: br,
		banRepo:       banRepo,
		bannedKV:      bannedKV,
		blacklistKV:   blacklistKV,
	}
}

func (f *mailtrapTestFixture) cleanup(t *testing.T) {
	t.Helper()
	testutils.TruncateAuthTables(t, f.pool)
	ctx := context.Background()
	_ = f.bannedKV.Delete(ctx, "127.0.0.1")
}

func (f *mailtrapTestFixture) seed(t *testing.T) *testutils.AuthSeedIDs {
	t.Helper()
	return testutils.SeedAuthData(t, f.pool)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func TestAuthVerifyUserSendsActivationEmailViaMailTrap(t *testing.T) {
	apiToken := os.Getenv("MAILTRAP_API_TOKEN")
	if apiToken == "" {
		t.Skip("MAILTRAP_API_TOKEN is not set — skipping MailTrap integration test")
	}
	accountID := os.Getenv("MAILTRAP_ACCOUNT_ID")
	if accountID == "" {
		t.Skip("MAILTRAP_ACCOUNT_ID is not set — skipping MailTrap integration test")
	}
	inboxID := os.Getenv("MAILTRAP_INBOX_ID")
	if inboxID == "" {
		t.Skip("MAILTRAP_INBOX_ID is not set — skipping MailTrap integration test")
	}

	mtClient := testutils.NewMailTrapClient(apiToken, accountID, inboxID)

	initialMsgs := mtClient.GetMessages(t)
	initialCount := len(initialMsgs)

	f := setupMailtrapTest(t)
	defer f.cleanup(t)
	ids := f.seed(t)

	body := map[string]any{"status": "Aktif"}
	path := "/users/" + fmt.Sprint(ids.UnverifiedUserID) + "/verification"
	resp := testutils.DoRequest(f.handler, http.MethodPatch, path, body, testutils.AccessCookie(f.jwtUtil, ids.AdminUserID, []string{"ADMIN"}))
	respBody := testutils.ReadBody(resp)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", resp.StatusCode, string(respBody))
	}

	msg := mtClient.WaitForNewMessage(t, initialCount, 30*time.Second)
	defer mtClient.DeleteMessage(t, msg.ID)

	htmlBody := mtClient.GetHTMLBody(t, msg)

	t.Logf("HTML body (first 500 chars): %s", truncate(htmlBody, 500))

	if msg.Subject != "Akun Anda telah diaktifkan - SiGizi" {
		t.Errorf("expected Subject %q, got %q", "Akun Anda telah diaktifkan - SiGizi", msg.Subject)
	}
	if !strings.Contains(htmlBody, "Unverified User") {
		t.Errorf("expected HTML body to contain %q", "Unverified User")
	}

	t.Logf("MailTrap test passed — message ID %d with subject %q", msg.ID, msg.Subject)
}
