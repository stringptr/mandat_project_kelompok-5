# SiGizi - Sistem Informasi Gizi dan Imunisasi

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26.2-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18.3-blue.svg)](https://www.postgresql.org)
[![Apache Hop](https://img.shields.io/badge/Apache%20Hop-2.17.0-orange.svg)](https://hop.apache.org)

## 📋 Deskripsi Proyek

**SiGizi** adalah sistem informasi terintegrasi untuk mengelola data gizi dan imunisasi anak serta ibu hamil di tingkat Posyandu. Sistem ini dirancang dengan arsitektur modern yang memisahkan **Sistem Informasi Transaksional (OLTP)** dan **Data Warehouse (OLAP)** untuk mendukung operasional harian sekaligus analisis data yang mendalam.

### 🎯 Tujuan Utama

1. **Digitalisasi Data Posyandu** - Menggantikan pencatatan manual dengan sistem digital
2. **Monitoring Real-time** - Memantau status gizi dan imunisasi secara langsung
3. **Business Intelligence** - Analisis data untuk pengambilan keputusan berbasis data
4. **Deteksi Dini** - Identifikasi kasus stunting, gizi buruk, dan KEK pada ibu hamil
5. **Integrasi Data** - Sinkronisasi data dari berbagai sumber ke data warehouse

---

## 🏗️ Arsitektur Sistem

### **Arsitektur 3-Tier: Staging → Master → Warehouse**

```
┌─────────────────────────────────────────────────────────────────────┐
│                          PRESENTATION LAYER                         │
│  ┌──────────────┐         ┌──────────────┐        ┌──────────────┐ │
│  │   Frontend   │ ──────► │   Backend    │ ◄────► │    Master    │ │
│  │   (React)    │         │  (Go/Huma)   │        │  (PostgreSQL)│ │
│  └──────────────┘         └──────────────┘        └──────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                            ETL LAYER                                │
│                      (Apache Hop 2.17.0)                            │
│                                                                     │
│  ┌──────────────┐         ┌──────────────┐        ┌──────────────┐│
│  │   Staging    │────ETL─►│    Master    │────ETL─►│  Warehouse   ││
│  │ (SQL Server) │   #1    │ (PostgreSQL) │   #2   │ (PostgreSQL) ││
│  └──────────────┘         └──────────────┘        └──────────────┘│
│       OLTP                     OLTP                    OLAP        │
│   (External Data)         (Normalized)            (Star Schema)    │
└─────────────────────────────────────────────────────────────────────┘

### **Penjelasan Layer:**

#### 1️⃣ **Staging Database (SQL Server)**
- **Purpose:** Data sementara dari sistem eksternal/legacy
- **Type:** OLTP (Transactional)
- **Schema:** Mirror dari sistem sumber
- **ETL:** `StagingToMaster` - Transformasi & validasi data

#### 2️⃣ **Master Database (PostgreSQL)**
- **Purpose:** Database operasional utama untuk sistem informasi
- **Type:** OLTP (Transactional)
- **Schema:** Normalized (3NF) dengan relasi lengkap
- **Users:** Backend API, Web Application, Mobile App
- **Features:** 
  - CRUD operations untuk semua entitas
  - Authentication & authorization
  - Soft delete mechanism
  - Audit trail (created_at, updated_at)

#### 3️⃣ **Warehouse Database (PostgreSQL)**
- **Purpose:** Data warehouse untuk analisis dan reporting
- **Type:** OLAP (Analytical)
- **Schema:** Star Schema (Fact & Dimension tables)
- **Users:** BI Tools, Analytics, Dashboards
- **Features:**
  - Optimized for complex queries
  - Historical data tracking
  - Aggregated metrics
  - Pre-calculated KPIs

---

## 🛠️ Tech Stack

### **Backend**
- **Language:** Go 1.26.2
- **Framework:** Huma v2.38.0 (OpenAPI-first REST framework)
- **Router:** Chi v5.3.0
- **ORM:** go-jet v2.15.1 (Type-safe SQL builder)
- **Database Driver:** pgx v5.9.2
- **Authentication:** JWT (golang-jwt/jwt v5.3.1)
- **Password Hashing:** Argon2id (golang.org/x/crypto)
- **Development:** Air (hot reload)

### **Frontend**
- **Framework:** React (Details in frontend folder)
- **Build Tool:** Modern JavaScript bundler
- **Development Server:** Integrated with Docker

### **Web Server**
- **Reverse Proxy:** Caddy 2.11.3
- **Features:** Automatic HTTPS, HTTP/2, HTTP/3 (QUIC)

### **ETL (Extract, Transform, Load)**
- **Tool:** Apache Hop 2.17.0
- **Interface:** Web-based GUI (port 8090)
- **Workflows:** Visual pipeline designer
- **Scheduling:** Built-in scheduler untuk automasi

### **Database**
- **Master & Warehouse:** PostgreSQL 18.3-alpine
- **Staging:** Microsoft SQL Server 2025-CU3
- **Migration Tool:** Goose (Go-based migration)

### **Infrastructure**
- **Containerization:** Docker & Docker Compose
- **Orchestration:** Docker Compose profiles
- **Development:** Hot reload untuk Backend & Frontend

---

## 📊 Database Schema

### **Master Database (Normalized - 3NF)**

#### **Core Entities:**

**User Management:**
- `user_account` - Data pengguna (bidan, kader, pasien, dinas kesehatan)
- `bidan` - Data bidan posyandu
- `kader_posyandu` - Data kader
- `dinas_kesehatan` - Data petugas dinas
- `pasien` - Data pasien (anak & ibu hamil)

**Patient Specific:**
- `anak` - Detail data anak (berat/panjang lahir, dll)
- `ibu_hamil` - Detail ibu hamil (HPHT, status kehamilan, dll)

**Health Records:**
- `jadwal_imunisasi` - Jadwal dan realisasi imunisasi
- `hasil_pemeriksaan` - Hasil pemeriksaan kesehatan
- `tindak_lanjut` - Tindak lanjut berdasarkan hasil pemeriksaan

**Master Data:**
- `posyandu` - Data posyandu
- `lokasi` - Hierarchical location (Provinsi → Kota → Kabupaten → Kecamatan → Kelurahan)
- `pendidikan` - Referensi tingkat pendidikan
- `pekerjaan` - Referensi jenis pekerjaan
- `kategori_pendapatan` - Referensi kategori penghasilan

**Key Features:**
- ✅ Foreign Key constraints
- ✅ Soft delete (`is_deleted` flag)
- ✅ Timestamp tracking (`created_at`, `updated_at`)
- ✅ Enum types (PostgreSQL native)
- ✅ Hierarchical data (lokasi)

### **Warehouse Database (Star Schema - Dimensional)**

#### **Dimension Tables:**

- `DIM_WAKTU` - Dimensi waktu (granularity: hari)
  - Tahun, semester, kuartal, bulan, minggu
  - Pre-generated untuk 2020-2030

- `DIM_LOKASI` - Dimensi lokasi (flattened hierarchy)
  - Provinsi, Kota, Kabupaten, Kecamatan, Kelurahan

- `DIM_PASIEN` - Dimensi pasien
  - Demographics: jenis kelamin, usia, pendidikan, pekerjaan
  - Socio-economic: kategori pendapatan, jumlah tanggungan
  - Calculated fields: kategori usia (Balita, Anak-Anak, Remaja, Dewasa)

- `DIM_ANAK` - Outrigger dari DIM_PASIEN
  - Berat lahir, panjang lahir
  - Hubungan dengan wali
  - Link ke ibu hamil (jika ada)

- `DIM_IBU_HAMIL` - Outrigger dari DIM_PASIEN
  - Hamil ke-N, status kehamilan
  - HPHT, bulan mulai hamil, trimester

- `DIM_POSYANDU` - Dimensi posyandu
  - Nama posyandu, wilayah kerja
  - Link ke bidan penanggung jawab

- `DIM_PETUGAS` - Dimensi petugas kesehatan
  - Nama, peran (Bidan, Kader, Dinas Kesehatan)

#### **Fact Tables:**

- `FACT_IMUNISASI` - Fakta imunisasi
  - **Dimensions:** Waktu, Lokasi, Pasien, Posyandu
  - **Measures:** 
    - Flag terlaksana (0/1)
    - Flag terlambat (0/1)
    - Hari keterlambatan
  - **Attributes:** Nama vaksin, status imunisasi

- `FACT_PEMERIKSAAN` - Fakta pemeriksaan kesehatan
  - **Dimensions:** Waktu, Lokasi, Pasien, Posyandu, Petugas, Kehamilan (optional)
  - **Measures:**
    - Berat badan, tinggi badan, lingkar kepala
    - Tekanan darah
    - Flag stunting, gizi buruk, KEK, anemia (0/1)
    - Flag perlu rujukan, risiko tinggi (0/1)
  - **Attributes:** Status stunting, status gizi, tipe pasien

---

## 🔄 ETL Pipelines

### **Pipeline 1: Staging → Master**

**Lokasi:** `etl/StagingToMaster/`

**Workflow:**
1. **00_cleanup.hpl** - Truncate staging tables (opsional)
2. **Data Transformation Pipelines:**
   - Validasi format data
   - Transformasi tipe data
   - Mapping ke schema Master
   - Handle missing values
3. **Data Loading:**
   - INSERT dengan conflict handling
   - Update timestamp metadata

**Features:**
- ✅ Data validation & cleansing
- ✅ Error handling & logging
- ✅ Rollback on failure

### **Pipeline 2: Master → Warehouse** ⭐

**Lokasi:** `etl/MasterToWarehouse/`

**Workflow Phases:**
1. **Phase 1 (Parallel):** Load reference dimensions
   - `dim_waktu.hpl` - Calendar table (2020-2030)
   - `dim_lokasi.hpl` - Location hierarchy (flattened)
   - `dim_petugas.hpl` - Healthcare workers (UNION all types)

2. **Phase 2 (Sequential):** Load dependent dimensions
   - `dim_posyandu.hpl` - Posyandu (depends on dim_petugas)

3. **Phase 3 (Sequential):** Load patient dimensions
   - `dim_pasien.hpl` - Patient demographics & socio-economic

4. **Phase 4 (Parallel):** Load patient outriggers
   - `dim_anak.hpl` - Child specific data
   - `dim_ibu_hamil.hpl` - Pregnant mother specific data

5. **Phase 5 (Parallel):** Load facts
   - `fact_imunisasi.hpl` - Immunization facts
   - `fact_pemeriksaan.hpl` - Health examination facts

**🚀 Key Features:**

#### **1. Incremental Loading** ✅
```sql
WHERE updated_at >= COALESCE(NULLIF('${LAST_ETL_RUN}', ''), '1900-01-01')::TIMESTAMPTZ
  AND is_deleted = false
```
- Hanya proses data yang berubah sejak ETL terakhir
- Filter berdasarkan `updated_at` timestamp
- Soft delete aware (`is_deleted = false`)
- Multi-table incremental (OR logic untuk joined tables)

#### **2. UPSERT Pattern** ✅
```sql
INSERT INTO DIM_TABLE (...) VALUES (...)
ON CONFLICT (primary_key) DO UPDATE SET
  column1 = EXCLUDED.column1,
  ...
WHERE DIM_TABLE.column1 IS DISTINCT FROM EXCLUDED.column1
   OR DIM_TABLE.column2 IS DISTINCT FROM EXCLUDED.column2
```
- Idempotent operations (aman di-run berkali-kali)
- No duplicates (conflict resolution)
- Skip update jika tidak ada perubahan

#### **3. Selective Column Update** ✅ (dim_pasien only)
```sql
ON CONFLICT (...) DO UPDATE SET
  column = CASE 
    WHEN target.column IS DISTINCT FROM source.column 
    THEN source.column 
    ELSE target.column 
  END
```
- **Hanya update kolom yang BENAR-BENAR berubah**
- Hemat write I/O & index maintenance
- NULL-safe comparison dengan `IS DISTINCT FROM`

#### **4. Change Detection** ✅
- WHERE clause untuk skip update jika semua kolom identik
- Efisiensi 99.95% (hanya proses data berubah)

#### **5. Parallel Execution** ✅
- Independent dimensions load secara parallel
- Dependency management dengan JOIN steps

**Performance:**
- Full Load: 100,000 rows → Process 100,000
- Incremental: 100,000 rows → Process 50 (99.95% reduction!)

---

## 🚀 Getting Started

### **Prerequisites**

- Docker Engine 20.10+
- Docker Compose v2.0+
- Git
- Minimal 8GB RAM
- 20GB free disk space

### **Installation**

#### 1. Clone Repository
```bash
git clone https://github.com/stringptr/SiGizi.git
cd SiGizi
```

#### 2. Setup Environment Variables
```bash
cp .env.example .env
# Edit .env sesuai kebutuhan
```

**Konfigurasi `.env` Penting:**
```env
# Database Credentials
MASTER_PASSWORD=your_secure_password
WAREHOUSE_PASSWORD=your_secure_password
STAGING_SA_PASSWORD=YourStrongPassword123!

# Port Configuration
MASTER_HOST_PORT=5432
WAREHOUSE_HOST_PORT=5433
STAGING_HOST_PORT=1433
HOP_HOST_PORT=8090
WEBSERVER_HTTP_HOST_PORT=8080

# Host IP (default localhost)
HOST_IP=127.0.0.1
```

#### 3. Start Services dengan Docker Compose Profiles

**Profile Options:**
- `dev` - All services (development)
- `information_system` - Frontend, Backend, Master DB, Webserver
- `data_management` - ETL, Staging, Master, Warehouse, Hop
- `backend` - Backend service only
- `frontend` - Frontend service only
- `warehouse` - Warehouse DB only
- `master` - Master DB only
- `staging` - Staging DB only
- `hop` - Apache Hop ETL tool only

**Start Full Development Environment:**
```bash
docker compose --profile dev up -d
```

**Start Information System Only:**
```bash
docker compose --profile information_system up -d
```

**Start Data Management (ETL) Only:**
```bash
docker compose --profile data_management up -d
```

#### 4. Run Database Migrations

Migrations akan otomatis jalan saat container start. Untuk manual migration:

```bash
# Master DB
docker compose --profile master up migration-master

# Warehouse DB
docker compose --profile warehouse up migration-warehouse

# Staging DB
docker compose --profile staging up migration-staging
```

#### 5. Access Services

| Service | URL | Credentials |
|---------|-----|-------------|
| **Frontend** | http://localhost:8080 | - |
| **Backend API** | http://localhost:8070 | - |
| **Apache Hop GUI** | http://localhost:8090 | hop/hop |
| **Master DB** | localhost:5432 | postgres/your_password |
| **Warehouse DB** | localhost:5433 | postgres/your_password |
| **Staging DB** | localhost:1433 | sa/your_password |

---

## 📖 Usage Guide

### **Backend API**

**OpenAPI Documentation:**
```
http://localhost:8070/docs
```

**Key Endpoints:**
- `POST /api/v1/auth/login` - Authentication
- `POST /api/v1/auth/refresh` - Refresh token
- `GET /api/v1/users` - List users
- `GET /api/v1/patients` - List patients
- `POST /api/v1/immunizations` - Create immunization record
- `POST /api/v1/examinations` - Create examination record

**Authentication:**
```bash
# Login
curl -X POST http://localhost:8070/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"pass"}'

# Use token
curl -X GET http://localhost:8070/api/v1/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Testing:**
Lihat dokumentasi lengkap di:
- `TestScript-Backend-Complete.md` - Comprehensive test scripts
- `backend/QUICK-START-TESTING.md` - Quick start guide
- `backend/test/` - Integration test files

### **Apache Hop ETL**

**Access GUI:**
1. Open http://localhost:8090
2. Login: `hop` / `hop`
3. Navigate to Project: `imunisasi`

**Run ETL Pipeline:**

**Via GUI:**
1. Open `MasterToWarehouse/RunAll.hwf`
2. Click "Run" button
3. Monitor execution logs

**Via Command Line:**
```bash
# Enter Hop container
docker exec -it sigizi-hop-1 bash

# Run workflow
sh /opt/hop/hop-run.sh \
  -j /files/projects/imunisasi/MasterToWarehouse/RunAll.hwf \
  -r local \
  -p LAST_ETL_RUN="2024-01-01 00:00:00"
```

**ETL Scheduling:**
1. Apache Hop GUI → Tools → Schedule Workflow
2. Set cron expression (e.g., `0 2 * * *` for 2 AM daily)
3. Configure email notifications (optional)

### **Database Management**

**Connect to PostgreSQL (Master):**
```bash
# Via psql
docker exec -it sigizi-master-1 psql -U postgres -d imunisasi

# Via DBeaver/pgAdmin
Host: localhost
Port: 5432
Database: imunisasi
Username: postgres
Password: your_password
```

**Connect to PostgreSQL (Warehouse):**
```bash
# Via psql
docker exec -it sigizi-warehouse-1 psql -U postgres -d imunisasi

# Via DBeaver/pgAdmin
Host: localhost
Port: 5433
Database: imunisasi
Username: postgres
Password: your_password
```

**Connect to SQL Server (Staging):**
```bash
# Via sqlcmd
docker exec -it sigizi-staging-1 /opt/mssql-tools18/bin/sqlcmd \
  -S localhost -U sa -P 'your_password' -C

# Via SSMS/Azure Data Studio
Host: localhost,1433
Username: sa
Password: your_password
```

---

## 🧪 Testing

### **Backend Testing**

**Run Unit Tests:**
```bash
cd backend
go test ./...
```

**Run Integration Tests:**
```bash
# Start test database
docker compose --profile dev up -d

# Run integration tests
cd backend
go test -tags=integration ./test/...
```

**Test Scripts:**
- `TestScript-Backend-Complete.md` - Manual testing guide
- `backend/test/health_test.go` - Health check tests
- `backend/test/monitoring_test.go` - Monitoring tests
- `backend/test/artikel_test.go` - Article CRUD tests
- `backend/test/imunisasi_test.go` - Immunization tests

**Test Coverage:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### **ETL Testing**

**Test Single Pipeline:**
1. Open Apache Hop GUI
2. Navigate to pipeline (e.g., `dim_pasien.hpl`)
3. Click Preview button
4. Set sample rows (e.g., 100)
5. Verify output data

**Test Full Workflow:**
1. Backup warehouse database
2. Truncate target tables
3. Run `RunAll.hwf`
4. Verify record counts
5. Check data quality

**Validation Queries:**
```sql
-- Check record counts
SELECT 'DIM_PASIEN' as table_name, COUNT(*) as count FROM DIM_PASIEN
UNION ALL
SELECT 'FACT_IMUNISASI', COUNT(*) FROM FACT_IMUNISASI
UNION ALL
SELECT 'FACT_PEMERIKSAAN', COUNT(*) FROM FACT_PEMERIKSAAN;

-- Check for nulls in FK
SELECT COUNT(*) FROM FACT_IMUNISASI WHERE id_dim_pasien IS NULL;

-- Verify date ranges
SELECT MIN(tanggal), MAX(tanggal) FROM DIM_WAKTU;
```

---

## 📁 Project Structure

```
SiGizi/
├── backend/                    # Go Backend API
│   ├── cmd/
│   │   ├── api/               # Main API application
│   │   └── jetgen/            # Code generator for go-jet
│   ├── internal/
│   │   ├── api/v1/            # API v1 routes
│   │   ├── config/            # Configuration management
│   │   ├── domain/            # Domain models & interfaces
│   │   ├── feature/           # Feature implementations
│   │   ├── infrastructure/    # Generated DB code (go-jet)
│   │   ├── middleware/        # HTTP middlewares
│   │   ├── hash/              # Password hashing (Argon2)
│   │   ├── jwtutils/          # JWT utilities
│   │   └── httputils/         # HTTP utilities
│   ├── test/                  # Integration tests
│   ├── .air.toml              # Air config (hot reload)
│   ├── go.mod                 # Go dependencies
│   └── Containerfile          # Docker image definition
│
├── frontend/                  # React Frontend
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── Containerfile
│
├── etl/                       # Apache Hop ETL
│   ├── MasterToWarehouse/     # Master → Warehouse ETL
│   │   ├── dimensions/        # Dimension load pipelines
│   │   │   ├── dim_waktu.hpl
│   │   │   ├── dim_lokasi.hpl
│   │   │   ├── dim_pasien.hpl ⭐ (selective column update)
│   │   │   ├── dim_anak.hpl
│   │   │   ├── dim_ibu_hamil.hpl
│   │   │   ├── dim_posyandu.hpl
│   │   │   └── dim_petugas.hpl
│   │   ├── facts/             # Fact load pipelines
│   │   │   ├── fact_imunisasi.hpl
│   │   │   └── fact_pemeriksaan.hpl
│   │   └── RunAll.hwf         # Master workflow
│   │
│   ├── StagingToMaster/       # Staging → Master ETL
│   │   ├── 00_cleanup.hpl
│   │   └── RunAllPipeline.hwf
│   │
│   ├── metadata/              # Hop metadata
│   ├── .env.json              # Environment config for Hop
│   └── project-config.json    # Hop project config
│
├── migrations/                # Database migrations (Goose)
│   ├── master/                # Master DB migrations
│   ├── warehouse/             # Warehouse DB migrations
│   ├── staging/               # Staging DB migrations
│   ├── data/                  # Sample data
│   └── Containerfile          # Migration runner image
│
├── webserver/                 # Caddy configuration
│   └── Caddyfile
│
├── storage/                   # Persistent data (gitignored)
│   ├── master/
│   ├── warehouse/
│   └── staging/
│
├── .github/                   # GitHub workflows
│   └── workflows/
│       ├── backend-unit-test.yml
│       └── backend-integration-test.yml
│
├── compose.yaml               # Docker Compose main config
├── compose.test.yaml          # Docker Compose test config
├── .env.example               # Environment variables template
└── README.md                  # This file
```

---

## 🔐 Security

### **Authentication & Authorization**
- Menggunakan standar **JSON Web Tokens (JWT)** untuk stateless authentication.
- Terdapat role-based access control (RBAC) dengan pemisahan akses untuk Bidan, Kader Posyandu, Dinas Kesehatan, dan Pasien Umum.

### **Password Hashing**
- Seluruh password di-hash secara aman menggunakan algoritma **Argon2id**, standar industri terkini untuk ketahanan terhadap brute-force dan rainbow table attacks.
- Salt unik di-generate secara acak untuk setiap user.

### **Data Protection**
- Soft Delete mechanism (`is_deleted`) untuk mencegah kehilangan data historis secara permanen.
- Penggunaan PostgreSQL Row-Level Security (RLS) jika dibutuhkan di masa mendatang.

---

## 👥 Tim Pengembang
- **Organisasi:** SiGizi Team
- **Kontributor:** [Daftar Kontributor]

## 📄 Lisensi
Proyek ini didistribusikan di bawah lisensi **MIT License**. Lihat file `LICENSE` untuk informasi lebih lanjut.

- **JWT Tokens:** Access token (15 min) + Refresh token (7 days)
- **Password Hashing:** Argon2id with salt
- **Role-Based Access Control:** Bidan, Kader, Dinas Kesehatan, Admin
- **IP Tracking:** Request IP logging untuk audit
- **HTTPS:** Automatic via Caddy (production)

### **Database Security**
- **Network Isolation:** Database containers tidak exposed ke public
- **Strong Passwords:** Enforced via environment variables
- **Soft Delete:** Data tidak benar-benar dihapus (audit trail)
- **Prepared Statements:** SQL injection prevention via go-jet
- **Connection Pooling:** pgx pool untuk resource management

### **Best Practices**
```bash
# JANGAN commit .env file
# JANGAN hardcode credentials
# GUNAKAN .env.example sebagai template
# ROTATE passwords secara berkala
# BACKUP database secara rutin
```

---

## 📊 Data Warehouse Analytics

### **Pre-built Analysis Queries**

#### **1. Monitoring Stunting by Location**
```sql
SELECT 
  dl.nama_kecamatan,
  dl.nama_kelurahan,
  COUNT(*) as total_pemeriksaan,
  SUM(fp.flag_stunting) as jumlah_stunting,
  ROUND(SUM(fp.flag_stunting)::NUMERIC / COUNT(*) * 100, 2) as persentase_stunting
FROM FACT_PEMERIKSAAN fp
JOIN DIM_LOKASI dl ON fp.id_lokasi = dl.id_lokasi
JOIN DIM_WAKTU dw ON fp.id_waktu = dw.id_waktu
WHERE dw.tahun = 2026
  AND fp.tipe_pasien = 'Anak'
GROUP BY dl.nama_kecamatan, dl.nama_kelurahan
ORDER BY persentase_stunting DESC;
```

#### **2. Immunization Coverage Analysis**
```sql
SELECT 
  fi.nama_vaksin,
  COUNT(*) as total_jadwal,
  SUM(fi.flag_terlaksana) as terlaksana,
  SUM(fi.flag_terlambat) as terlambat,
  ROUND(SUM(fi.flag_terlaksana)::NUMERIC / COUNT(*) * 100, 2) as coverage_rate,
  ROUND(AVG(fi.hari_keterlambatan), 1) as avg_hari_keterlambatan
FROM FACT_IMUNISASI fi
JOIN DIM_WAKTU dw ON fi.id_waktu = dw.id_waktu
WHERE dw.tahun = 2026
GROUP BY fi.nama_vaksin
ORDER BY coverage_rate DESC;
```

#### **3. KEK in Pregnant Women Trend**
```sql
SELECT 
  dw.nama_bulan,
  dw.tahun,
  COUNT(*) as total_ibu_hamil,
  SUM(fp.flag_kek) as jumlah_kek,
  ROUND(SUM(fp.flag_kek)::NUMERIC / COUNT(*) * 100, 2) as persentase_kek
FROM FACT_PEMERIKSAAN fp
JOIN DIM_WAKTU dw ON fp.id_waktu = dw.id_waktu
WHERE fp.tipe_pasien = 'Ibu Hamil'
  AND dw.tahun BETWEEN 2025 AND 2026
GROUP BY dw.tahun, dw.bulan, dw.nama_bulan
ORDER BY dw.tahun, dw.bulan;
```

#### **4. Posyandu Performance**
```sql
SELECT 
  dp.nama_posyandu,
  dp.wilayah_kerja,
  COUNT(DISTINCT fp.id_dim_pasien) as jumlah_pasien,
  COUNT(*) as total_pemeriksaan,
  SUM(fp.flag_risiko_tinggi) as kasus_risiko_tinggi,
  ROUND(SUM(fp.flag_risiko_tinggi)::NUMERIC / COUNT(*) * 100, 2) as persentase_risiko
FROM FACT_PEMERIKSAAN fp
JOIN DIM_POSYANDU dp ON fp.id_posyandu = dp.id_posyandu
JOIN DIM_WAKTU dw ON fp.id_waktu = dw.id_waktu
WHERE dw.tahun = 2026
GROUP BY dp.nama_posyandu, dp.wilayah_kerja
ORDER BY total_pemeriksaan DESC;
```

#### **5. Socio-Economic Impact Analysis**
```sql
SELECT 
  dp_dim.kategori_pendapatan,
  dp_dim.pendidikan,
  COUNT(*) as jumlah_kasus,
  SUM(fp.flag_gizi_buruk) as gizi_buruk,
  SUM(fp.flag_stunting) as stunting,
  ROUND(AVG(fp.berat_badan), 2) as avg_berat_badan,
  ROUND(AVG(fp.tinggi_badan), 2) as avg_tinggi_badan
FROM FACT_PEMERIKSAAN fp
JOIN DIM_PASIEN dp_dim ON fp.id_dim_pasien = dp_dim.id_dim_pasien
WHERE fp.tipe_pasien = 'Anak'
GROUP BY dp_dim.kategori_pendapatan, dp_dim.pendidikan
ORDER BY gizi_buruk DESC, stunting DESC;
```

### **BI Tool Integration**

Warehouse siap diintegrasikan dengan:
- **Tableau** - Koneksi PostgreSQL langsung
- **Power BI** - Via PostgreSQL connector
- **Metabase** - Open-source BI tool
- **Apache Superset** - Python-based BI
- **Custom Dashboards** - Via Backend API

---

## 🛠️ Development

### **Backend Development**

**Setup:**
```bash
cd backend
go mod download
go mod tidy
```

**Run with Hot Reload:**
```bash
# Via Docker (recommended)
docker compose --profile backend up

# Local development
go install github.com/cosmtrek/air@latest
air -c .air.toml
```

**Generate DB Code (go-jet):**
```bash
# After schema changes
cd backend
go run cmd/jetgen/main.go
```

**Code Style:**
```bash
# Format code
go fmt ./...

# Lint
golangci-lint run ./...
```

### **Frontend Development**

```bash
cd frontend
npm install
npm run dev
```

### **ETL Development**

**Edit Pipelines:**
1. Open Apache Hop GUI: http://localhost:8090
2. Navigate to pipeline file
3. Edit transformations visually
4. Save changes
5. Test with Preview feature
6. Run workflow to validate

**Best Practices:**
- ✅ Test pipeline dengan sample data terlebih dahulu
- ✅ Backup database sebelum run ETL
- ✅ Monitor logs untuk error detection
- ✅ Validate output data setelah ETL
- ✅ Document pipeline logic di description field
- ✅ Use meaningful transform names
- ✅ Add error handling steps

### **Database Migrations**

**Create New Migration:**
```bash
# Master DB
cd migrations/master
goose create add_new_feature sql

# Warehouse DB
cd migrations/warehouse
goose create add_new_dimension sql
```

**Migration File Structure:**
```sql
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE new_table (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE new_table;
```

**Run Migrations:**
```bash
# Via Docker (automatic on startup)
docker compose --profile dev up -d

# Manual
docker exec -it sigizi-migration-master-1 goose up
docker exec -it sigizi-migration-warehouse-1 goose up
```

**Rollback:**
```bash
docker exec -it sigizi-migration-master-1 goose down
```

---

## 🐛 Troubleshooting

### **Common Issues**

#### **1. Database Connection Failed**
```bash
# Check if containers running
docker ps

# Check container logs
docker logs sigizi-master-1
docker logs sigizi-warehouse-1

# Restart containers
docker compose --profile dev restart
```

#### **2. Port Already in Use**
```bash
# Check what's using the port
# Windows
netstat -ano | findstr :5432

# Kill process or change port in .env
MASTER_HOST_PORT=5435
```

#### **3. ETL Pipeline Fails**
```bash
# Check Hop logs
docker logs sigizi-hop-1

# Check database connectivity in Hop GUI
# Tools → Database Connections → Test Connection

# Verify LAST_ETL_RUN variable
# Should be set or pipeline will full load
```

#### **4. Migration Fails**
```bash
# Check migration status
docker exec -it sigizi-migration-master-1 goose status

# Fix failed migration
docker exec -it sigizi-migration-master-1 goose fix

# Reset database (DANGER: deletes all data)
docker exec -it sigizi-master-1 psql -U postgres -c "DROP DATABASE imunisasi;"
docker exec -it sigizi-master-1 psql -U postgres -c "CREATE DATABASE imunisasi;"
docker compose --profile master up migration-master
```

#### **5. Backend API Not Responding**
```bash
# Check backend logs
docker logs sigizi-backend-1 -f

# Check if database is accessible
docker exec -it sigizi-backend-1 ping master

# Rebuild backend
docker compose --profile backend up --build
```

### **Debug Mode**

**Enable Verbose Logging:**
```env
# .env file
LOG_LEVEL=debug
HOP_LOG_LEVEL=debug
```

**Access Container Shell:**
```bash
# Backend
docker exec -it sigizi-backend-1 sh

# Database
docker exec -it sigizi-master-1 psql -U postgres -d imunisasi

# Hop
docker exec -it sigizi-hop-1 bash
```

---

## 📈 Performance Optimization

### **Database Optimization**

#### **Indexing Strategy:**
```sql
-- Master DB - OLTP indexes
CREATE INDEX idx_pasien_updated_at ON pasien(updated_at) WHERE is_deleted = false;
CREATE INDEX idx_user_account_updated_at ON user_account(updated_at) WHERE is_deleted = false;
CREATE INDEX idx_jadwal_imunisasi_updated_at ON jadwal_imunisasi(updated_at);

-- Warehouse - OLAP indexes
CREATE INDEX idx_fact_pemeriksaan_waktu ON FACT_PEMERIKSAAN(id_waktu);
CREATE INDEX idx_fact_pemeriksaan_lokasi ON FACT_PEMERIKSAAN(id_lokasi);
CREATE INDEX idx_fact_pemeriksaan_pasien ON FACT_PEMERIKSAAN(id_dim_pasien);
CREATE INDEX idx_fact_imunisasi_waktu ON FACT_IMUNISASI(id_waktu);
```

#### **Query Optimization:**
```sql
-- Use EXPLAIN ANALYZE untuk melihat query plan
EXPLAIN ANALYZE
SELECT * FROM FACT_PEMERIKSAAN 
WHERE id_waktu IN (SELECT id_waktu FROM DIM_WAKTU WHERE tahun = 2026);

-- Materialize common aggregations
CREATE MATERIALIZED VIEW mv_monthly_stats AS
SELECT 
  dw.tahun, dw.bulan,
  COUNT(*) as total_pemeriksaan,
  SUM(flag_stunting) as total_stunting
FROM FACT_PEMERIKSAAN fp
JOIN DIM_WAKTU dw ON fp.id_waktu = dw.id_waktu
GROUP BY dw.tahun, dw.bulan;

-- Refresh materialized view (setelah ETL)
REFRESH MATERIALIZED VIEW mv_monthly_stats;
```

### **ETL Optimization**

#### **Batch Processing:**
```
# di Transform settings
Batch Size: 1000
Commit Size: 1000
```

#### **Parallel Execution:**
- Set number of copies untuk intensive transforms
- Use parallel workflow execution
- Optimize JOIN operations

#### **Memory Management:**
```bash
# Increase Hop memory
# Edit docker-compose.yaml
environment:
  HOP_OPTIONS: "-Xms2g -Xmx4g"
```

### **Backend Optimization**

#### **Connection Pooling:**
```go
// internal/config/config.go
pgxConfig.MaxConns = 25
pgxConfig.MinConns = 5
pgxConfig.MaxConnLifetime = time.Hour
pgxConfig.MaxConnIdleTime = 30 * time.Minute
```

#### **Caching:**
```go
// Implement Redis caching for frequent queries
// Cache dimension tables (rarely change)
```

---

## 🔄 CI/CD

### **GitHub Actions**

**Automated Tests:**
- `.github/workflows/backend-unit-test.yml` - Unit tests on push
- `.github/workflows/backend-integration-test.yml` - Integration tests on PR

**Workflow Triggers:**
- Push to `main` or `develop` branch
- Pull request creation/update
- Manual dispatch

### **Deployment**

**Production Deployment:**
```bash
# Build production images
docker compose -f compose.prod.yaml build

# Deploy to production
docker compose -f compose.prod.yaml up -d

# Run migrations
docker compose -f compose.prod.yaml exec migration-master goose up
docker compose -f compose.prod.yaml exec migration-warehouse goose up
```

**Environment-specific Config:**
```
.env.development
.env.staging
.env.production
```

---

## 📚 Documentation

### **Additional Resources**

- **Backend API:** `http://localhost:8070/docs` (OpenAPI/Swagger)
- **Testing Guide:** `TestScript-Backend-Complete.md`
- **Quick Start:** `backend/QUICK-START-TESTING.md`
- **Test Summary:** `TESTING-SUMMARY.md`
- **Automated Testing:** `AUTOMATED-TEST-GUIDE.md`

### **ETL Documentation**

- **Pipeline Descriptions:** Lihat description field di setiap `.hpl` file
- **Workflow Notes:** Lihat notepad di `RunAll.hwf`
- **Transformation Logic:** Comment di JavaScript transform

### **External Documentation**

- [Apache Hop Documentation](https://hop.apache.org/manual/latest/)
- [Go-Jet Documentation](https://github.com/go-jet/jet)
- [Huma Framework](https://huma.rocks/)
- [PostgreSQL 18 Docs](https://www.postgresql.org/docs/18/)

---

## 🤝 Contributing

### **How to Contribute**

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

### **Coding Standards**

**Go:**
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Run `golangci-lint` before commit
- Write tests for new features
- Update OpenAPI spec when adding endpoints

**SQL:**
- Use uppercase for SQL keywords
- Proper indentation (2 spaces)
- Always use table aliases
- Comment complex queries

**ETL:**
- Meaningful transform names
- Add description untuk setiap pipeline
- Error handling di setiap critical step
- Test dengan sample data

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👥 Team

**Developers:**
- Backend: [stringptr](https://github.com/stringptr)
- ETL: Data Engineering Team
- Frontend: Frontend Team

**Contact:**
- Email: support@sigizi.id
- Issues: [GitHub Issues](https://github.com/stringptr/SiGizi/issues)

---

## 🙏 Acknowledgments

- Apache Software Foundation for Apache Hop
- PostgreSQL Global Development Group
- Go Team at Google
- All open-source contributors

---

## 📊 Project Status

**Current Version:** 1.0.0-beta  
**Status:** Active Development  
**Last Updated:** June 2026

### **Roadmap**

- [x] Backend API with JWT authentication
- [x] Master database schema with migrations
- [x] Warehouse star schema design
- [x] ETL Staging → Master
- [x] ETL Master → Warehouse (Incremental + Selective Column Update)
- [x] Docker containerization
- [ ] Frontend UI completion
- [ ] BI Dashboard integration
- [ ] Mobile application
- [ ] Automated ETL scheduling
- [ ] Data quality monitoring
- [ ] Real-time alerts
- [ ] Multi-tenancy support

---

## 📞 Support

Jika mengalami masalah atau memiliki pertanyaan:

1. **Check Documentation:** README ini, test guides, API docs
2. **Check Issues:** [GitHub Issues](https://github.com/stringptr/SiGizi/issues)
3. **Create New Issue:** Dengan detail lengkap (logs, screenshots, steps to reproduce)
4. **Contact Team:** Via email atau project communication channel

---

**⭐ Star this repository if you find it useful!**

**🔗 Links:**
- GitHub: https://github.com/stringptr/SiGizi
- Documentation: [Wiki](https://github.com/stringptr/SiGizi/wiki)
- Issues: [Issue Tracker](https://github.com/stringptr/SiGizi/issues)
