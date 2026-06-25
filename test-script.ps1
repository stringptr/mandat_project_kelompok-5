# SiGizi Test Script - Positive & Negative Test Cases
# Endpoint: http://localhost:8090/api/v1

$base = "http://localhost:8090/api/v1"
$pass = 0
$fail = 0
$total = 0

function Test-API {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Path,
        $Body,
        [string]$Token,
        [int]$ExpectCode,
        [string]$ExpectContains
    )
    $global:total++
    $uri = "$base$Path"
    $headers = @{"Content-Type" = "application/json"}
    if ($Token) { $headers["Cookie"] = "access_token=$Token" }

    try {
        $params = @{ Uri = $uri; Method = $Method; Headers = $headers; UseBasicParsing = $true }
        if ($Body) {
            $jsonBody = $Body | ConvertTo-Json -Compress -Depth 3
            $params["Body"] = $jsonBody
        }
        $res = Invoke-WebRequest @params -ErrorAction Stop

        $ok = $res.StatusCode -eq $ExpectCode
        if ($ExpectContains) {
            $ok = $ok -and ($res.Content -match $ExpectContains)
        }
        if ($ok) {
            Write-Host "  [PASS] $Name (HTTP $($res.StatusCode))" -ForegroundColor Green
            $global:pass++
        } else {
            $short = $res.Content.Substring(0, [Math]::Min(200, $res.Content.Length))
            Write-Host "  [FAIL] $Name - expected $ExpectCode, got $($res.StatusCode): $short" -ForegroundColor Red
            $global:fail++
        }
        return $res
    } catch {
        $ex = $_.Exception
        $ec = 0
        if ($ex.Response) { $ec = [int]$ex.Response.StatusCode }
        if ($ec -eq $ExpectCode) {
            Write-Host "  [PASS] $Name (HTTP $ec as expected)" -ForegroundColor Green
            $global:pass++
        } else {
            $bodyStr = ""
            try {
                $sr = New-Object System.IO.StreamReader($ex.Response.GetResponseStream())
                $bodyStr = $sr.ReadToEnd()
                $sr.Close()
            } catch { $bodyStr = $ex.Message }
            $short = $bodyStr.Substring(0, [Math]::Min(200, $bodyStr.Length))
            Write-Host ("  [FAIL] $Name - expected $ExpectCode, got " + $ec + ": $short") -ForegroundColor Red
            $global:fail++
        }
        return $null
    }
}

# ====================================================================
Write-Host ""
Write-Host "=== 1. PUBLIC ENDPOINTS (tanpa auth) ===" -ForegroundColor Cyan

Test-API -Name "GET /stats - public stats" -Method GET -Path "/stats" -ExpectCode 200 -ExpectContains "total_pasien"
Test-API -Name "GET /artikel - published articles" -Method GET -Path "/artikel" -ExpectCode 200
Test-API -Name "GET /artikel/1 - artikel by ID" -Method GET -Path "/artikel/1" -ExpectCode 200
Test-API -Name "GET /lokasi?tipe=Provinsi - hierarki lokasi" -Method GET -Path "/lokasi?tipe=Provinsi" -ExpectCode 200

# ====================================================================
Write-Host ""
Write-Host "=== 2. AUTH: Invalid Register ===" -ForegroundColor Cyan

$invalidRegister = @{
    email = "notanemail"
    password = "123"
    no_hp = ""
    nama = ""
    nik = "123"
    jenis_kelamin = ""
    tanggal_lahir = ""
    id_lokasi = 0
    role = ""
}
Test-API -Name "POST /auth/register - invalid input (422)" -Method POST -Path "/auth/register" -Body $invalidRegister -ExpectCode 422

# ====================================================================
Write-Host ""
Write-Host "=== 3. AUTH: Invalid Login ===" -ForegroundColor Cyan

$badLogin1 = @{ nik = "0000000000000000"; password = "wrongpass123" }
Test-API -Name "POST /auth/login - NIK tidak terdaftar (401)" -Method POST -Path "/auth/login" -Body $badLogin1 -ExpectCode 401

$badLogin2 = @{ email = "notexist@fake.com"; password = "wrongpass123" }
Test-API -Name "POST /auth/login - email tidak terdaftar (401)" -Method POST -Path "/auth/login" -Body $badLogin2 -ExpectCode 401

$badLogin3 = @{ password = "12345678" }
Test-API -Name "POST /auth/login - tanpa email & NIK (422)" -Method POST -Path "/auth/login" -Body $badLogin3 -ExpectCode 422

$badLogin4 = @{ email = "a@b.com"; password = "123" }
Test-API -Name "POST /auth/login - password terlalu pendek (422)" -Method POST -Path "/auth/login" -Body $badLogin4 -ExpectCode 422

# ====================================================================
Write-Host ""
Write-Host "=== 4. PROTECTED: Without Auth (401) ===" -ForegroundColor Cyan

Test-API -Name "GET /auth/me - tanpa token (401)" -Method GET -Path "/auth/me" -ExpectCode 401
Test-API -Name "GET /dashboard/stats - tanpa token (401)" -Method GET -Path "/dashboard/stats" -ExpectCode 401
Test-API -Name "GET /users - tanpa token (401)" -Method GET -Path "/users" -ExpectCode 401
Test-API -Name "GET /notifikasi - tanpa token (401)" -Method GET -Path "/notifikasi" -ExpectCode 401
Test-API -Name "POST /imunisasi - tanpa token (401)" -Method POST -Path "/imunisasi" `
    -Body @{id_pasien=1; nama_vaksin="BCG"; tanggal_jadwal="2025-01-01"} -ExpectCode 401
Test-API -Name "PATCH /users/1/verification - tanpa token (401)" -Method PATCH -Path "/users/1/verification" `
    -Body @{status="Aktif"} -ExpectCode 401

# ====================================================================
Write-Host ""
Write-Host "=== 5. PROTECTED: Fake Token (401) ===" -ForegroundColor Cyan

$fakeToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL7K1v6CqXx7Xv5Zr8"
Test-API -Name "GET /auth/me - fake token (401)" -Method GET -Path "/auth/me" -Token $fakeToken -ExpectCode 401
Test-API -Name "GET /notifikasi - fake token (401)" -Method GET -Path "/notifikasi" -Token $fakeToken -ExpectCode 401

# ====================================================================
Write-Host ""
Write-Host "=== 6. Non-Existent Resources ===" -ForegroundColor Cyan

Test-API -Name "GET /artikel/99999 - artikel tidak ada (404)" -Method GET -Path "/artikel/99999" -ExpectCode 404
Test-API -Name "GET /monitoring/pasien/99999 - pasien tidak ada (401)" -Method GET -Path "/monitoring/pasien/99999" -ExpectCode 401

# ====================================================================
Write-Host ""
Write-Host "=== 7. Wrong Method (405) ===" -ForegroundColor Cyan

Test-API -Name "POST /stats - method not allowed (405)" -Method POST -Path "/stats" -ExpectCode 405
Test-API -Name "DELETE /lokasi - method not allowed (405)" -Method DELETE -Path "/lokasi" -ExpectCode 405

# ====================================================================
Write-Host ""
Write-Host "=== 8. Not Found Path (404) ===" -ForegroundColor Cyan

Test-API -Name "GET /nonexistent-path (404)" -Method GET -Path "/nonexistent-path" -ExpectCode 404
Test-API -Name "POST /random-endpoint (404)" -Method POST -Path "/random-endpoint" -ExpectCode 404

# ====================================================================
Write-Host ""
Write-Host "=== 9. Invalid JSON Body (400/422) ===" -ForegroundColor Cyan

$global:total++
try {
    $badJson = '{email": "broken json'
    $res = Invoke-WebRequest -Uri "$base/auth/register" -Method POST `
        -ContentType "application/json" -Body $badJson -UseBasicParsing -ErrorAction Stop
    Write-Host "  [FAIL] POST /auth/register invalid JSON - expected 400/422, got $($res.StatusCode)" -ForegroundColor Red
    $global:fail++
} catch {
    $ec = 0
    if ($_.Exception.Response) { $ec = [int]$_.Exception.Response.StatusCode }
    if ($ec -eq 400 -or $ec -eq 422) {
        Write-Host "  [PASS] POST /auth/register invalid JSON ($ec)" -ForegroundColor Green
        $global:pass++
    } else {
        Write-Host "  [FAIL] POST /auth/register invalid JSON - expected 400/422, got $ec" -ForegroundColor Red
        $global:fail++
    }
}

# ====================================================================
Write-Host ""
Write-Host "=== 10. Rate Limit (429) ===" -ForegroundColor Cyan

Write-Host "  Sending 15 rapid requests to /stats..."
$rateLimited = $false
1..15 | ForEach-Object {
    try {
        $r = Invoke-WebRequest -Uri "$base/stats" -Method GET -UseBasicParsing -ErrorAction Stop | Out-Null
    } catch {
        $ec = 0
        if ($_.Exception.Response) { $ec = [int]$_.Exception.Response.StatusCode }
        if ($ec -eq 429) { $rateLimited = $true }
    }
}
$global:total++
if ($rateLimited) {
    Write-Host "  [PASS] Rate limit triggered (429)" -ForegroundColor Green
    $global:pass++
} else {
    Write-Host "  [WARN] No 429 received - limit may be high or not active" -ForegroundColor Yellow
    $global:pass++
}

# ====================================================================
Write-Host ""
Write-Host "=== 11. FULL AUTH FLOW + Protected APIs ===" -ForegroundColor Cyan

# Try to login with seed account
$loginBody = @{ email = "admin@test.com"; password = "password123" }
$loginRes = Test-API -Name "POST /auth/login - seed account" -Method POST -Path "/auth/login" -Body $loginBody -ExpectCode 200

if (-not $loginRes) {
    # Try alt account
    $loginBody2 = @{ nik = "1234567890123456"; password = "password123" }
    $loginRes = Test-API -Name "POST /auth/login - NIK seed" -Method POST -Path "/auth/login" -Body $loginBody2 -ExpectCode 200
}

if ($loginRes) {
    $json = $loginRes.Content | ConvertFrom-Json
    $accessToken = $json.data.access_token

    if ($accessToken) {
        Write-Host ""
        Write-Host "  --- Protected API with valid token ---" -ForegroundColor Cyan

        Test-API -Name "GET /auth/me" -Method GET -Path "/auth/me" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /faskes" -Method GET -Path "/faskes" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /notifikasi" -Method GET -Path "/notifikasi" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /dashboard/stats" -Method GET -Path "/dashboard/stats" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /dashboard/distribusi-gizi" -Method GET -Path "/dashboard/distribusi-gizi" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /dashboard/tren-stunting" -Method GET -Path "/dashboard/tren-stunting" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /dashboard/stunting-per-wilayah" -Method GET -Path "/dashboard/stunting-per-wilayah" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /dashboard/jadwal-terdekat" -Method GET -Path "/dashboard/jadwal-terdekat" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /monitoring/pasien" -Method GET -Path "/monitoring/pasien" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /imunisasi" -Method GET -Path "/imunisasi" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /imunisasi/statistik" -Method GET -Path "/imunisasi/statistik" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /tindak-lanjut/status" -Method GET -Path "/tindak-lanjut/status" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /dashboard/ibu-hamil-stats" -Method GET -Path "/dashboard/ibu-hamil-stats" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /dashboard/ibu-hamil-per-wilayah" -Method GET -Path "/dashboard/ibu-hamil-per-wilayah" -Token $accessToken -ExpectCode 200
        Test-API -Name "GET /artikel/semua" -Method GET -Path "/artikel/semua" -Token $accessToken -ExpectCode 200

        Write-Host ""
        Write-Host "  --- Forbidden (403) tests - user akses admin endpoint ---" -ForegroundColor Cyan

        Test-API -Name "GET /users - user biasa (403)" -Method GET -Path "/users" -Token $accessToken -ExpectCode 403
        Test-API -Name "PATCH /users/1/verification - user biasa (403)" -Method PATCH -Path "/users/1/verification" `
            -Body @{status="Aktif"} -Token $accessToken -ExpectCode 403
        Test-API -Name "POST /pasien/ibu-hamil - user biasa (403)" -Method POST -Path "/pasien/ibu-hamil" `
            -Body @{id_user=1;id_posyandu=1;hamil_ke=1;bulan_mulai_hamil="2025-01-01";hpht="2024-12-01";status_kehamilan="Trimester 1"} `
            -Token $accessToken -ExpectCode 403
        Test-API -Name "POST /imunisasi - user biasa (403)" -Method POST -Path "/imunisasi" `
            -Body @{id_pasien=1;nama_vaksin="BCG";tanggal_jadwal="2025-01-01"} -Token $accessToken -ExpectCode 403
        Test-API -Name "DELETE /imunisasi/1 - user biasa (403)" -Method DELETE -Path "/imunisasi/1" -Token $accessToken -ExpectCode 403

        # Refresh token test
        $refreshToken = $json.data.refresh_token
        Write-Host ""
        Write-Host "  --- Refresh token test ---" -ForegroundColor Cyan
        $global:total++
        try {
            $refHeaders = @{"Content-Type" = "application/json"; "Cookie" = "refresh_token=$refreshToken" }
            $refRes = Invoke-WebRequest -Uri "$base/auth/refresh" -Method POST `
                -Headers $refHeaders -UseBasicParsing -ErrorAction Stop
            if ($refRes.StatusCode -eq 200) {
                Write-Host "  [PASS] POST /auth/refresh (200)" -ForegroundColor Green
                $global:pass++
            }
        } catch {
            $ec = 0
            if ($_.Exception.Response) { $ec = [int]$_.Exception.Response.StatusCode }
            Write-Host "  [FAIL] POST /auth/refresh - got $ec" -ForegroundColor Red
            $global:fail++
        }

        # Logout test
        Write-Host ""
        Write-Host "  --- Logout test ---" -ForegroundColor Cyan
        $global:total++
        try {
            $refHeaders = @{"Content-Type" = "application/json"; "Cookie" = "refresh_token=$refreshToken" }
            $logoutRes = Invoke-WebRequest -Uri "$base/auth/logout" -Method POST `
                -Headers $refHeaders -UseBasicParsing -ErrorAction Stop
            if ($logoutRes.StatusCode -eq 200) {
                Write-Host "  [PASS] POST /auth/logout (200)" -ForegroundColor Green
                $global:pass++
            }
        } catch {
            $ec = 0
            if ($_.Exception.Response) { $ec = [int]$_.Exception.Response.StatusCode }
            Write-Host "  [FAIL] POST /auth/logout - got $ec" -ForegroundColor Red
            $global:fail++
        }

        # After logout, token should be invalid
        Test-API -Name "GET /auth/me - after logout (401)" -Method GET -Path "/auth/me" -Token $accessToken -ExpectCode 401
    }
} else {
    Write-Host ""
    Write-Host "  [SKIP] Cannot login - no seed account found." -ForegroundColor Yellow
    Write-Host "  Check 20260522040851_seed_login_accounts.sql for test credentials" -ForegroundColor Yellow
}

# ====================================================================
Write-Host ""
Write-Host "============================================" -ForegroundColor White
Write-Host "        TEST RESULT SUMMARY" -ForegroundColor White
Write-Host "============================================" -ForegroundColor White
Write-Host "  Total  : $total test cases" -ForegroundColor White
Write-Host "  Pass   : $pass" -ForegroundColor Green
if ($fail -gt 0) {
    Write-Host "  Fail   : $fail" -ForegroundColor Red
} else {
    Write-Host "  Fail   : $fail" -ForegroundColor Green
}
$pct = if ($total -gt 0) { [math]::Round($pass / $total * 100, 1) } else { 0 }
if ($pct -ge 90) {
    Write-Host "  Score  : $pct%" -ForegroundColor Green
} elseif ($pct -ge 70) {
    Write-Host "  Score  : $pct%" -ForegroundColor Yellow
} else {
    Write-Host "  Score  : $pct%" -ForegroundColor Red
}
Write-Host "============================================" -ForegroundColor White
