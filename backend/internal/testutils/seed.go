package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stringptr/SiGizi/backend/internal/hash"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type AuthSeedIDs struct {
	VerifiedUserID   int32
	UnverifiedUserID int32
	AdminUserID      int32
	RegularUserID    int32
}

func SeedAuthData(t *testing.T, pool *pgxpool.Pool) *AuthSeedIDs {
	t.Helper()
	ctx := context.Background()

	_, err := pool.Exec(ctx, `INSERT INTO lokasi (id_lokasi, nama_lokasi, tipe_lokasi) VALUES (1, 'Test Location', 'Kelurahan') ON CONFLICT (id_lokasi) DO NOTHING`)
	if err != nil {
		t.Fatalf("failed to seed lokasi: %v", err)
	}

	hashedPass, err := hash.Hash("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	now := time.Now().Truncate(time.Microsecond)

	var verifiedUserID int32
	err = pool.QueryRow(ctx, `
		INSERT INTO user_account (email, password, no_hp, status_verifikasi, nama, nik, jenis_kelamin, tanggal_lahir, id_lokasi, akun_ke, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1, $10, $10)
		RETURNING id_user`,
		"verified@example.com", hashedPass, "081111111111", string(model.StatusVerifikasi_Aktif),
		"Verified User", "3201020304050001", string(model.JenisKelamin_Perempuan),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), int32(1), now,
	).Scan(&verifiedUserID)
	if err != nil {
		t.Fatalf("failed to seed verified user: %v", err)
	}

	_, err = pool.Exec(ctx,
		"INSERT INTO bidan (id_user, no_str, wilayah_kerja, created_at, updated_at) VALUES ($1, '12345/STR/2024', 1, $2, $2)",
		verifiedUserID, now)
	if err != nil {
		t.Fatalf("failed to seed bidan role: %v", err)
	}

	var unverifiedUserID int32
	err = pool.QueryRow(ctx, `
		INSERT INTO user_account (email, password, no_hp, status_verifikasi, nama, nik, jenis_kelamin, tanggal_lahir, id_lokasi, akun_ke, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1, $10, $10)
		RETURNING id_user`,
		"unverified@example.com", hashedPass, "081222222222", string(model.StatusVerifikasi_Pending),
		"Unverified User", "3201020304050002", string(model.JenisKelamin_LakiLaki),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), int32(1), now,
	).Scan(&unverifiedUserID)
	if err != nil {
		t.Fatalf("failed to seed unverified user: %v", err)
	}

	var adminUserID int32
	err = pool.QueryRow(ctx, `
		INSERT INTO user_account (email, password, no_hp, status_verifikasi, nama, nik, jenis_kelamin, tanggal_lahir, id_lokasi, akun_ke, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1, $10, $10)
		RETURNING id_user`,
		"admin@example.com", hashedPass, "081333333333", string(model.StatusVerifikasi_Aktif),
		"Admin User", "3201020304050003", string(model.JenisKelamin_LakiLaki),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), int32(1), now,
	).Scan(&adminUserID)
	if err != nil {
		t.Fatalf("failed to seed admin user: %v", err)
	}

	_, err = pool.Exec(ctx,
		"INSERT INTO dinas_kesehatan (id_user, created_at, updated_at) VALUES ($1, $2, $2)",
		adminUserID, now)
	if err != nil {
		t.Fatalf("failed to seed dinkes role: %v", err)
	}

	var regularUserID int32
	err = pool.QueryRow(ctx, `
		INSERT INTO user_account (email, password, no_hp, status_verifikasi, nama, nik, jenis_kelamin, tanggal_lahir, id_lokasi, akun_ke, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1, $10, $10)
		RETURNING id_user`,
		"user@example.com", hashedPass, "081444444444", string(model.StatusVerifikasi_Aktif),
		"Regular User", "3201020304050004", string(model.JenisKelamin_Perempuan),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), int32(1), now,
	).Scan(&regularUserID)
	if err != nil {
		t.Fatalf("failed to seed regular user: %v", err)
	}

	return &AuthSeedIDs{
		VerifiedUserID:   verifiedUserID,
		UnverifiedUserID: unverifiedUserID,
		AdminUserID:      adminUserID,
		RegularUserID:    regularUserID,
	}
}

type PasienSeedIDs struct {
	PosyanduID int32
	IbuHamilID int32
	AnakID     int32
}

func SeedPosyandu(t *testing.T, pool *pgxpool.Pool, idLokasi, idBidan int32) int32 {
	t.Helper()
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)

	var idPosyandu int32
	err := pool.QueryRow(ctx, `
		INSERT INTO posyandu (nama_posyandu, id_lokasi, id_bidan, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id_posyandu`,
		"Posyandu Test", idLokasi, idBidan, now,
	).Scan(&idPosyandu)
	if err != nil {
		t.Fatalf("failed to seed posyandu: %v", err)
	}
	return idPosyandu
}

func SeedPasienIbuHamil(t *testing.T, pool *pgxpool.Pool, idUser, idPosyandu int32) int32 {
	t.Helper()
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)

	_, err := pool.Exec(ctx,
		"INSERT INTO pasien (id_pasien, id_posyandu, created_at, updated_at) VALUES ($1, $2, $3, $3)",
		idUser, idPosyandu, now)
	if err != nil {
		t.Fatalf("failed to seed pasien: %v", err)
	}

	var idIbuHamil int32
	nowDate := now.Format("2006-01-02")
	err = pool.QueryRow(ctx, `
		INSERT INTO ibu_hamil (id_pasien, hamil_ke, bulan_mulai_hamil, hpht, status_kehamilan, created_at, updated_at)
		VALUES ($1, 1, $2::date, $3::date, $4, $5, $5)
		RETURNING id_ibu_hamil`,
		idUser, nowDate, nowDate, "Trimester 1", now,
	).Scan(&idIbuHamil)
	if err != nil {
		t.Fatalf("failed to seed ibu_hamil: %v", err)
	}

	return idUser
}

func SeedPasienAnak(t *testing.T, pool *pgxpool.Pool, idUser, idPosyandu, idWali int32) int32 {
	t.Helper()
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)

	_, err := pool.Exec(ctx,
		"INSERT INTO pasien (id_pasien, id_posyandu, created_at, updated_at) VALUES ($1, $2, $3, $3)",
		idUser, idPosyandu, now)
	if err != nil {
		t.Fatalf("failed to seed pasien: %v", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO anak (id_pasien, id_wali, nama_anak, berat_lahir, panjang_lahir, hubungan_dengan_wali, created_at, updated_at)
		VALUES ($1, $2, 'Anak Test', 3.0, 50.0, 'Kandung', $3, $3)`,
		idUser, idWali, now)
	if err != nil {
		t.Fatalf("failed to seed anak: %v", err)
	}

	return idUser
}

func SeedActiveSession(t *testing.T, pool *pgxpool.Pool, userID int32) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)
	sessionID := uuid.New()

	_, err := pool.Exec(ctx, `
		INSERT INTO user_session (id_session, id_user, status_session, ip_address, expired_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)`,
		sessionID, userID, string(model.StatusSession_Aktif), "127.0.0.1",
		now.Add(7*24*time.Hour), now,
	)
	if err != nil {
		t.Fatalf("failed to seed active session: %v", err)
	}

	return sessionID
}
