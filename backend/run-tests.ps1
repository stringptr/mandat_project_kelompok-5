# PowerShell Script untuk menjalankan automated tests
# Usage: .\run-tests.ps1 [option]
# Options: all, health, monitoring, artikel, tindaklanjut, imunisasi, coverage

param(
    [Parameter(Position=0)]
    [string]$TestType = "all"
)

Write-Host "🧪 SiGizi Backend Automated Test Runner" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Check if test database is running (Docker)
Write-Host "📊 Checking test database..." -ForegroundColor Yellow
$dockerStatus = docker ps --filter "name=sigizi-test-db" --format "{{.Status}}" 2>$null

if (-not $dockerStatus) {
    Write-Host "⚠️  Test database not found!" -ForegroundColor Red
    Write-Host "`nDo you want to start a test database container? (Y/N)" -ForegroundColor Yellow
    $response = Read-Host
    
    if ($response -eq "Y" -or $response -eq "y") {
        Write-Host "`n🚀 Starting test database container..." -ForegroundColor Green
        docker run -d `
            --name sigizi-test-db `
            -e POSTGRES_DB=sigizi_test `
            -e POSTGRES_USER=postgres `
            -e POSTGRES_PASSWORD=postgres `
            -p 5433:5432 `
            postgres:16
        
        Write-Host "⏳ Waiting for database to be ready..." -ForegroundColor Yellow
        Start-Sleep -Seconds 5
    } else {
        Write-Host "`n❌ Test aborted. Please start test database first." -ForegroundColor Red
        exit 1
    }
}

Write-Host "✅ Test database is running`n" -ForegroundColor Green

# Run tests based on type
switch ($TestType.ToLower()) {
    "all" {
        Write-Host "🏃 Running ALL tests..." -ForegroundColor Cyan
        go test ./test/... -v -cover -coverprofile=coverage.out
    }
    "health" {
        Write-Host "🏃 Running Health Check tests..." -ForegroundColor Cyan
        go test ./test/ -v -run TestHealthCheck
    }
    "monitoring" {
        Write-Host "🏃 Running Monitoring tests..." -ForegroundColor Cyan
        go test ./test/ -v -run TestMonitoring
    }
    "artikel" {
        Write-Host "🏃 Running Artikel tests..." -ForegroundColor Cyan
        go test ./test/ -v -run TestArtikel
    }
    "tindaklanjut" {
        Write-Host "🏃 Running Tindak Lanjut tests..." -ForegroundColor Cyan
        go test ./test/ -v -run TestTindakLanjut
    }
    "imunisasi" {
        Write-Host "🏃 Running Imunisasi tests..." -ForegroundColor Cyan
        go test ./test/ -v -run TestImunisasi
    }
    "coverage" {
        Write-Host "🏃 Running tests with coverage report..." -ForegroundColor Cyan
        go test ./test/... -v -cover -coverprofile=coverage.out
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "`n📊 Generating HTML coverage report..." -ForegroundColor Green
            go tool cover -html=coverage.out -o coverage.html
            Write-Host "✅ Coverage report generated: coverage.html" -ForegroundColor Green
            Write-Host "`n🌐 Opening coverage report in browser..." -ForegroundColor Cyan
            Start-Process coverage.html
        }
    }
    "clean" {
        Write-Host "🧹 Cleaning test cache..." -ForegroundColor Cyan
        go clean -testcache
        Write-Host "✅ Test cache cleared" -ForegroundColor Green
    }
    default {
        Write-Host "❌ Unknown test type: $TestType" -ForegroundColor Red
        Write-Host "`nAvailable options:" -ForegroundColor Yellow
        Write-Host "  all          - Run all tests"
        Write-Host "  health       - Run health check tests only"
        Write-Host "  monitoring   - Run monitoring tests only"
        Write-Host "  artikel      - Run artikel tests only"
        Write-Host "  tindaklanjut - Run tindak lanjut tests only"
        Write-Host "  imunisasi    - Run imunisasi tests only"
        Write-Host "  coverage     - Run tests with coverage report"
        Write-Host "  clean        - Clear test cache"
        exit 1
    }
}

# Check test result
if ($LASTEXITCODE -eq 0) {
    Write-Host "`n✅ All tests passed!" -ForegroundColor Green
} else {
    Write-Host "`n❌ Some tests failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}
