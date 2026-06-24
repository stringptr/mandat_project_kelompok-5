#!/bin/bash
# Bash Script untuk menjalankan automated tests
# Usage: ./run-tests.sh [option]
# Options: all, health, monitoring, artikel, tindaklanjut, imunisasi, coverage

TEST_TYPE=${1:-all}

echo "🧪 SiGizi Backend Automated Test Runner"
echo "========================================"
echo ""

# Check if test database is running (Docker)
echo "📊 Checking test database..."
DOCKER_STATUS=$(docker ps --filter "name=sigizi-test-db" --format "{{.Status}}" 2>/dev/null)

if [ -z "$DOCKER_STATUS" ]; then
    echo "⚠️  Test database not found!"
    echo ""
    read -p "Do you want to start a test database container? (Y/N): " response
    
    if [ "$response" = "Y" ] || [ "$response" = "y" ]; then
        echo ""
        echo "🚀 Starting test database container..."
        docker run -d \
            --name sigizi-test-db \
            -e POSTGRES_DB=sigizi_test \
            -e POSTGRES_USER=postgres \
            -e POSTGRES_PASSWORD=postgres \
            -p 5433:5432 \
            postgres:16
        
        echo "⏳ Waiting for database to be ready..."
        sleep 5
    else
        echo ""
        echo "❌ Test aborted. Please start test database first."
        exit 1
    fi
fi

echo "✅ Test database is running"
echo ""

# Run tests based on type
case $TEST_TYPE in
    all)
        echo "🏃 Running ALL tests..."
        go test ./test/... -v -cover -coverprofile=coverage.out
        ;;
    health)
        echo "🏃 Running Health Check tests..."
        go test ./test/ -v -run TestHealthCheck
        ;;
    monitoring)
        echo "🏃 Running Monitoring tests..."
        go test ./test/ -v -run TestMonitoring
        ;;
    artikel)
        echo "🏃 Running Artikel tests..."
        go test ./test/ -v -run TestArtikel
        ;;
    tindaklanjut)
        echo "🏃 Running Tindak Lanjut tests..."
        go test ./test/ -v -run TestTindakLanjut
        ;;
    imunisasi)
        echo "🏃 Running Imunisasi tests..."
        go test ./test/ -v -run TestImunisasi
        ;;
    coverage)
        echo "🏃 Running tests with coverage report..."
        go test ./test/... -v -cover -coverprofile=coverage.out
        
        if [ $? -eq 0 ]; then
            echo ""
            echo "📊 Generating HTML coverage report..."
            go tool cover -html=coverage.out -o coverage.html
            echo "✅ Coverage report generated: coverage.html"
            echo ""
            echo "🌐 Opening coverage report in browser..."
            
            # Open browser based on OS
            if [[ "$OSTYPE" == "darwin"* ]]; then
                open coverage.html
            elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
                xdg-open coverage.html 2>/dev/null || echo "Please open coverage.html manually"
            fi
        fi
        ;;
    clean)
        echo "🧹 Cleaning test cache..."
        go clean -testcache
        echo "✅ Test cache cleared"
        ;;
    *)
        echo "❌ Unknown test type: $TEST_TYPE"
        echo ""
        echo "Available options:"
        echo "  all          - Run all tests"
        echo "  health       - Run health check tests only"
        echo "  monitoring   - Run monitoring tests only"
        echo "  artikel      - Run artikel tests only"
        echo "  tindaklanjut - Run tindak lanjut tests only"
        echo "  imunisasi    - Run imunisasi tests only"
        echo "  coverage     - Run tests with coverage report"
        echo "  clean        - Clear test cache"
        exit 1
        ;;
esac

# Check test result
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ All tests passed!"
else
    echo ""
    echo "❌ Some tests failed!"
    exit 1
fi
