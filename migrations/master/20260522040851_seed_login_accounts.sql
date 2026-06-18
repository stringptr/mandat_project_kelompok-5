-- +goose Up
-- ============================================================
-- SEED LOGIN ACCOUNTS — semua jenis role untuk user account
-- Password: password123 (Argon2 hashed)
-- ID dimulai dari 200001 (setelah data CSV yang berakhir di 200000)
-- ============================================================

-- ============================================================
-- 1. USER ACCOUNT (6 akun, masing-masing role berbeda)
-- ============================================================
INSERT INTO user_account (id_user, email, password, no_hp, status_verifikasi, nama, nik, jenis_kelamin, tanggal_lahir, id_lokasi, id_pendidikan, id_pekerjaan, id_pendapatan, jumlah_tanggungan, akun_ke, created_at, updated_at) VALUES
(20000001, 'admin@dinkes.test',  'MDEyMzQ1Njc4OWFiY2RlZg$5G22S13CsChMna+Gxobht0Q2qxBS8D3BYXJ/6nM8hhQ', '081234567801', 'Aktif', 'Admin Dinkes',  '1234567890123401', 'Laki-Laki',  '1990-01-01', 4, 6, 4,  4, 2, 1, NOW(), NOW()),
(20000002, 'bidan@test.com',     'MDEyMzQ1Njc4OWFiY2RlZg$5G22S13CsChMna+Gxobht0Q2qxBS8D3BYXJ/6nM8hhQ', '081234567802', 'Aktif', 'Bidan Test',    '1234567890123402', 'Perempuan', '1995-05-15', 4, 6, 5,  4, 0, 1, NOW(), NOW()),
(20000003, 'kader@test.com',     'MDEyMzQ1Njc4OWFiY2RlZg$5G22S13CsChMna+Gxobht0Q2qxBS8D3BYXJ/6nM8hhQ', '081234567803', 'Aktif', 'Kader Test',    '1234567890123403', 'Perempuan', '1998-08-20', 4, 5, 8,  3, 1, 1, NOW(), NOW()),
(20000004, 'pasien@test.com',    'MDEyMzQ1Njc4OWFiY2RlZg$5G22S13CsChMna+Gxobht0Q2qxBS8D3BYXJ/6nM8hhQ', '081234567804', 'Aktif', 'Pasien Test',   '1234567890123404', 'Laki-Laki',  '2000-10-10', 5, 3, 1,  2, 3, 1, NOW(), NOW()),
(20000005, 'ibuhamil@test.com',  'MDEyMzQ1Njc4OWFiY2RlZg$5G22S13CsChMna+Gxobht0Q2qxBS8D3BYXJ/6nM8hhQ', '081234567805', 'Aktif', 'Ibu Hamil Test','1234567890123405', 'Perempuan', '2002-03-20', 5, 4, 11, 3, 0, 1, NOW(), NOW()),
(20000006, 'user@test.com',      'MDEyMzQ1Njc4OWFiY2RlZg$5G22S13CsChMna+Gxobht0Q2qxBS8D3BYXJ/6nM8hhQ', '081234567806', 'Aktif', 'User Test',     '1234567890123406', 'Laki-Laki',  '2005-07-07', 5, 5, 9,  3, 0, 1, NOW(), NOW());

-- ============================================================
-- 2. DINAS KESEHATAN (Role: USER + DINKES + SUPER_ADMIN + ADMIN)
-- ============================================================
INSERT INTO dinas_kesehatan (id_user, created_at, updated_at) VALUES
(20000001, NOW(), NOW());

-- ============================================================
-- 3. BIDAN (Role: USER + BIDAN + ADMIN)
-- ============================================================
INSERT INTO bidan (id_user, no_str, wilayah_kerja, created_at, updated_at) VALUES
(20000002, 'STR-DUMMY-LOGIN-001', 4, NOW(), NOW());

-- ============================================================
-- 4. POSYANDU (id_posyandu 501, setelah 500 dari CSV)
-- ============================================================
INSERT INTO posyandu (id_posyandu, nama_posyandu, id_lokasi, id_bidan, created_at, updated_at) VALUES
(501, 'Posyandu Login Sehat', 4, 20000002, NOW(), NOW());

-- ============================================================
-- 5. KADER POSYANDU (Role: USER + KADER + ADMIN)
-- ============================================================
INSERT INTO kader_posyandu (id_user, no_sk, id_posyandu, created_at, updated_at) VALUES
(20000003, 'SK-DUMMY-LOGIN-001', 501, NOW(), NOW());

-- ============================================================
-- 6. PASIEN (Role: USER + PASIEN)
-- ============================================================
INSERT INTO pasien (id_pasien, id_posyandu, created_at, updated_at) VALUES
(20000004, 501, NOW(), NOW()),
(20000005, 501, NOW(), NOW());

-- ============================================================
-- 7. IBU HAMIL (Role: USER + PASIEN + IBU_HAMIL)
-- ============================================================
INSERT INTO ibu_hamil (id_pasien, hamil_ke, bulan_mulai_hamil, hpht, status_kehamilan, created_at, updated_at) VALUES
(20000005, 1, '2026-01-15', '2025-12-20', 'Trimester 2', NOW(), NOW());

-- ============================================================
-- RESET SEQUENCES
-- ============================================================
SELECT setval('user_account_id_user_seq',       COALESCE((SELECT MAX(id_user)     FROM user_account),      200000));
SELECT setval('posyandu_id_posyandu_seq',        COALESCE((SELECT MAX(id_posyandu) FROM posyandu),           500));
SELECT setval('ibu_hamil_id_ibu_hamil_seq',      COALESCE((SELECT MAX(id_ibu_hamil) FROM ibu_hamil),       48010));

-- ============================================================
-- VERIFICATION — tampilkan role setiap login user
-- ============================================================
SELECT
    u.id_user,
    u.email,
    u.nama,
    CASE WHEN d.id_user IS NOT NULL THEN 'DINKES'     ELSE NULL END AS role_dinkes,
    CASE WHEN b.id_user IS NOT NULL THEN 'BIDAN'      ELSE NULL END AS role_bidan,
    CASE WHEN k.id_user IS NOT NULL THEN 'KADER'      ELSE NULL END AS role_kader,
    CASE WHEN p.id_pasien IS NOT NULL AND ih.id_pasien IS NOT NULL THEN 'PASIEN+IBU_HAMIL'
         WHEN p.id_pasien IS NOT NULL THEN 'PASIEN'
         ELSE NULL END AS role_pasien
FROM user_account u
LEFT JOIN dinas_kesehatan d  ON d.id_user  = u.id_user
LEFT JOIN bidan b            ON b.id_user  = u.id_user
LEFT JOIN kader_posyandu k   ON k.id_user  = u.id_user
LEFT JOIN pasien p           ON p.id_pasien = u.id_user
LEFT JOIN ibu_hamil ih       ON ih.id_pasien = p.id_pasien
WHERE u.id_user >= 200001
ORDER BY u.id_user;

-- +goose Down

-- ============================================================
-- ROLLBACK
-- ============================================================
DELETE FROM ibu_hamil      WHERE id_pasien = 200005;
DELETE FROM pasien         WHERE id_pasien IN (200004, 200005);
DELETE FROM kader_posyandu WHERE id_user   = 200003;
DELETE FROM posyandu       WHERE id_posyandu = 501;
DELETE FROM bidan          WHERE id_user   = 200002;
DELETE FROM dinas_kesehatan WHERE id_user  = 200001;
DELETE FROM user_account   WHERE id_user   >= 200001;
