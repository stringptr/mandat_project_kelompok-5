# ====================================================================
# SKENARIO ARTIKEL GA MUNCUL -- Debug & Test
# Base URL: http://localhost:8090/api/v1
# ====================================================================

$base = "http://localhost:8090/api/v1"
$pass = 0; $fail = 0; $total = 0

function Test-Case {
    param([string]$Label, [scriptblock]$Script)
    $global:total++
    try {
        & $Script
        Write-Host "  [PASS] $Label" -ForegroundColor Green
        $global:pass++
    } catch {
        $msg = $_.Exception.Message.Substring(0, [Math]::Min(150, $_.Exception.Message.Length))
        Write-Host "  [FAIL] $Label -- $msg" -ForegroundColor Red
        $global:fail++
    }
}

# ====================================================================
Write-Host "================================================" -ForegroundColor Cyan
Write-Host " SKENARIO 1: Artikel Dipublikasikan (public)" -ForegroundColor Cyan
Write-Host "================================================"

$url1 = "$base/artikel"
$res = Invoke-WebRequest $url1 -UseBasicParsing -Body @{page=1; per_page=100}
$pub = ($res.Content | ConvertFrom-Json).data
Write-Host "  Total artikel publik: $($pub.meta.total)"
Write-Host "  Di halaman 1: $($pub.artikel.Count) artikel"

Test-Case "Ada artikel publik" {
    if ($pub.artikel.Count -eq 0) { throw "ARTIKEL GA MUNCUL!" }
}
Test-Case "Status selalu Dipublikasikan" {
    $bukan = $pub.artikel | Where-Object { $_.status_artikel -ne "Dipublikasikan" }
    if ($bukan) { throw "Ada $($bukan.Count) artikel bukan Dipublikasikan!" }
}
Test-Case "Punya nama penulis" {
    $kosong = $pub.artikel | Where-Object { -not $_.nama_penulis }
    if ($kosong) { throw "$($kosong.Count) artikel tanpa penulis" }
}

# ====================================================================
Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host " SKENARIO 2: Alasan artikel ga muncul di publik" -ForegroundColor Cyan
Write-Host "================================================"

# Login admin
$loginBody = '{"email":"bidan@test.com","password":"password123"}'
$loginRes = Invoke-WebRequest "$base/auth/login" -Method POST -ContentType "application/json" -Body $loginBody -UseBasicParsing -SessionVariable sess
$token = ($loginRes.Content | ConvertFrom-Json).data.access_token

# Ambil SEMUA artikel via admin
$allUrl = "$base/artikel/semua"
$allRes = Invoke-WebRequest $allUrl -UseBasicParsing -WebSession $sess -Body @{page=1; per_page=200}
$all = ($allRes.Content | ConvertFrom-Json).data

Write-Host "  Total artikel di sistem: $($all.meta.total)"
Write-Host "  Artikel publik (Dipublikasikan):  $($pub.meta.total)"
Write-Host "  Artikel TERSEMBUNYI dari publik: $($all.meta.total - $pub.meta.total)"
Write-Host ""
Write-Host "  Alasan ga muncul di halaman publik:"
$all.artikel | Where-Object { $_.status_artikel -ne "Dipublikasikan" } | 
    Group-Object status_artikel | 
    ForEach-Object { Write-Host "    Status '$($_.Name)' = $($_.Count) artikel TERSEMBUNYI" }

# ====================================================================
Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host " SKENARIO 3: Artikel Draft -- ga muncul di mana pun (publik)" -ForegroundColor Cyan
Write-Host "================================================"

# Buat artikel baru via admin (status Draft karena Bidan, bukan Dinkes)
$timeStr = Get-Date -Format "HHmmss"
$newBody = (@{
    judul = "Test Artikel Debug $timeStr"
    isi_artikel = "Artikel test untuk debugging. Seharusnya berstatus Draft."
    kategori = "Debug"
} | ConvertTo-Json)

$createRes = Invoke-WebRequest "$base/artikel" -Method POST -ContentType "application/json" -Body $newBody -UseBasicParsing -WebSession $sess
$created = ($createRes.Content | ConvertFrom-Json).data
$newId = $created.id_artikel
Write-Host "  Artikel baru dibuat: id=$newId status=$($created.status_artikel)"

Test-Case "Status artikel baru = Draft (Bidan bukan Dinkes)" {
    if ($created.status_artikel -ne "Draft") { throw "Status: $($created.status_artikel), bukan Draft" }
}

# Cek: apakah muncul di publik? (seharusnya TIDAK)
$pubCheck = Invoke-WebRequest "$base/artikel" -UseBasicParsing -Body @{page=1; per_page=200}
$pubAfter = ($pubCheck.Content | ConvertFrom-Json).data
$foundPub = $pubAfter.artikel | Where-Object { $_.id_artikel -eq $newId }
Test-Case "Artikel Draft TIDAK muncul di halaman publik" {
    if ($foundPub) { throw "Artikel Draft id=$newId muncul di publik -- SEHARUSNYA TIDAK!" }
}

# Cek: apakah muncul di admin/semua? (seharusnya YA)
$allCheck = Invoke-WebRequest "$base/artikel/semua" -UseBasicParsing -WebSession $sess -Body @{page=1; per_page=200}
$allAfter = ($allCheck.Content | ConvertFrom-Json).data
$foundAll = $allAfter.artikel | Where-Object { $_.id_artikel -eq $newId }
Test-Case "Artikel Draft MUNCUL di /artikel/semua (admin)" {
    if (-not $foundAll) { throw "Artikel Draft id=$newId tidak muncul di admin -- SEHARUSNYA MUNCUL!" }
}

# ====================================================================
Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host " SKENARIO 4: Artikel Ditolak -- permanen ga muncul" -ForegroundColor Cyan
Write-Host "================================================"

# Login Dinkes untuk review
$dinkesBody = '{"email":"admin@dinkes.test","password":"password123"}'
try {
    $dinkesRes = Invoke-WebRequest "$base/auth/login" -Method POST -ContentType "application/json" -Body $dinkesBody -UseBasicParsing -SessionVariable ds
    $dinkesToken = ($dinkesRes.Content | ConvertFrom-Json).data.access_token
    Write-Host "  Login Dinkes OK"

    # Tolak artikel yang barusan dibuat
    $tolak = (@{ aksi = "tolak"; catatan_review = "Artikel tidak sesuai pedoman" } | ConvertTo-Json)
    $rvUrl = "$base/artikel/$newId/review"
    $reviewRes = Invoke-WebRequest $rvUrl -Method PATCH -ContentType "application/json" -Body $tolak -UseBasicParsing -WebSession $ds
    $reviewed = ($reviewRes.Content | ConvertFrom-Json).data
    Write-Host "  Artikel id=$newId ditolak -- status: $($reviewed.status_artikel)"

    Test-Case "Status jadi Ditolak setelah review" {
        if ($reviewed.status_artikel -ne "Ditolak") { throw "Status: $($reviewed.status_artikel)" }
    }

    # Cek: apakah muncul di publik?
    $pubTolakCheck = Invoke-WebRequest "$base/artikel" -UseBasicParsing -Body @{page=1; per_page=200}
    $pubTolakData = ($pubTolakCheck.Content | ConvertFrom-Json).data
    $foundTolak = $pubTolakData.artikel | Where-Object { $_.id_artikel -eq $newId }
    Test-Case "Artikel Ditolak TIDAK muncul di publik" {
        if ($foundTolak) { throw "Artikel Ditolak masih muncul di publik!" }
    }

    # Cek: apakah muncul di pending? (seharusnya TIDAK)
    $pendUrl = "$base/artikel/pending"
    $pendingRes = Invoke-WebRequest $pendUrl -UseBasicParsing -WebSession $ds -Body @{page=1; per_page=200}
    $pending = ($pendingRes.Content | ConvertFrom-Json).data
    $foundPending = $pending.artikel | Where-Object { $_.id_artikel -eq $newId }
    Test-Case "Artikel Ditolak TIDAK muncul di pending Dinkes" {
        if ($foundPending) { throw "Artikel Ditolak masih di pending Dinkes!" }
    }
} catch {
    Write-Host "  [SKIP] Dinkes login/review gagal: $($_.Exception.Message.Substring(0, 100))"
}

# ====================================================================
Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host " SKENARIO 5: Frontend cache -- artikel baru ga muncul" -ForegroundColor Cyan
Write-Host "================================================"

Write-Host "  Frontend ada cache 10 MENIT (CACHE_TTL = 600000ms)"
Write-Host "  GET request ke endpoint yang sama dalam 10 menit return data cache"
Write-Host "  BUKAN dari API call baru."
Write-Host ""
Write-Host "  IMPLIKASI:"
Write-Host "   - Buat artikel baru -> publikasi -> ga muncul di frontend 10 menit"
Write-Host "   - Edit artikel -> perubahan ga keliatan 10 menit"
Write-Host "   - Hapus artikel -> masih muncul di list 10 menit"
Write-Host ""
Write-Host "  SOLUSI:"
Write-Host "   - POST/PATCH/DELETE otomatis clear cache"
Write-Host "   - Atau user refresh halaman setelah 10 menit"
Write-Host "   - Atau kurangi CACHE_TTL di lib/api.ts"

# ====================================================================
Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host " SKENARIO 6: Artikel tidak ada -- 404 & empty page" -ForegroundColor Cyan
Write-Host "================================================"

# Cek artikel publik terakhir
$lastPubId = $pub.artikel[-1].id_artikel
$detUrl = "$base/artikel/$lastPubId"
$detailRes = Invoke-WebRequest $detUrl -UseBasicParsing
$detail = ($detailRes.Content | ConvertFrom-Json).data
Write-Host "  Artikel publik terakhir: id=$lastPubId judul='$($detail.judul)'"

Test-Case "GET /artikel/{id} -- detail tersedia" {
    if (-not $detail) { throw "Detail artikel tidak ditemukan" }
}

# Cek ID yang tidak ada
$missingId = 99999999
try {
    $murl = "$base/artikel/$missingId"
    $missing = Invoke-WebRequest $murl -UseBasicParsing -ErrorAction Stop
    Write-Host "  GET /artikel/$missingId -> HTTP $($missing.StatusCode) (seharusnya 404)"
} catch {
    $ec = $_.Exception.Response.StatusCode.value__
    Write-Host "  GET /artikel/$missingId -> HTTP $ec (tidak ada)"
    Test-Case "Artikel tidak ada return 404" {
        if ($ec -ne 404) { throw "Expected 404, got $ec" }
    }
}

# Cek halaman jauh (seharusnya kosong)
$emptyUrl = "$base/artikel"
$emptyRes = Invoke-WebRequest $emptyUrl -UseBasicParsing -Body @{page=999; per_page=20}
$empty = ($emptyRes.Content | ConvertFrom-Json).data
Write-Host "  Halaman 999: $($empty.artikel.Count) artikel (meta.total=$($empty.meta.total))"

Test-Case "Halaman jauh return array kosong (tidak error)" {
    if ($null -eq $empty.artikel) { throw "Response tidak punya field artikel" }
}

# ====================================================================
Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host " SKENARIO 7: Verifikator null -- artikel belum direview" -ForegroundColor Cyan
Write-Host "================================================"

# Ambil detail beberapa artikel untuk cek verifikator
$sampleArts = $allAfter.artikel | Select-Object -First 5
foreach ($art in $sampleArts) {
    $sid = $art.id_artikel
    $sdUrl = "$base/artikel/$sid"
    try {
        $sdRes = Invoke-WebRequest $sdUrl -UseBasicParsing
        $d = ($sdRes.Content | ConvertFrom-Json).data
        $v = if ($d.nama_verifikator) { $d.nama_verifikator } else { "NULL (belum direview)" }
        Write-Host "  Artikel $sid : $($d.judul.Substring(0, [Math]::Min(35, $d.judul.Length))) verifikator=$v status=$($d.status_artikel)"
    } catch { Write-Host "  Artikel $sid : GAGAL fetch" }
}

# ====================================================================
Write-Host ""
Write-Host "================================================" -ForegroundColor White
Write-Host "        SUMMARY ARTIKEL DEBUG" -ForegroundColor White
Write-Host "================================================" -ForegroundColor White
Write-Host "  Artikel publik (Dipublikasikan) : $($pub.meta.total)" -ForegroundColor Green
Write-Host "  Artikel tersembunyi (non-publik): $($all.meta.total - $pub.meta.total)" -ForegroundColor Yellow
Write-Host "  Cache frontend                 : 10 menit (CACHE_TTL)" -ForegroundColor Yellow
Write-Host ""
Write-Host "  PENYEBAB ARTIKEL GA MUNCUL DI PUBLIK:" -ForegroundColor White
Write-Host "  1. Status bukan Dipublikasikan (Draft/Ditolak/Menunggu/Diarsipkan)" -ForegroundColor Yellow
Write-Host "  2. Cache frontend 10 menit -- data baru belum kerefresh" -ForegroundColor Yellow
Write-Host "  3. Dihapus (hard delete) -- gap di database" -ForegroundColor Yellow
Write-Host "  4. Halaman terlalu jauh -- halaman > last_page return kosong" -ForegroundColor Yellow
Write-Host "  5. Belum direview Dinkes -- tetap Menunggu Verifikasi" -ForegroundColor Yellow
Write-Host ""
$score = if ($total -gt 0) { [math]::Round($pass / $total * 100, 1) } else { 0 }
Write-Host "  Test: $total | Pass: $pass | Fail: $fail | Score: $score%" -ForegroundColor $(if ($fail -gt 0) { "Red" } else { "Green" })
Write-Host "================================================" -ForegroundColor White
