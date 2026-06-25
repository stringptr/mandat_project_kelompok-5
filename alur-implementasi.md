# Alur Implementasi — SiGizi

> Tahapan implementasi Sistem Informasi Gizi & Monitoring Imunisasi Anak dari nol hingga production.

---

## 1. Persiapan Environment

Langkah pertama adalah menyiapkan seluruh environment pengembangan. Pastikan Docker Engine versi terbaru sudah terinstal karena seluruh infrastruktur — database, message broker, reverse proxy, ETL, frontend, dan backend — dijalankan melalui Docker Compose. Untuk development, dibutuhkan Go 1.26 sebagai runtime backend, Node.js dengan Bun sebagai package manager frontend, dan Git untuk version control.

Setelah repository di-clone, langkah berikutnya adalah mengkonfigurasi environment variables melalui file `.env` yang diletakkan di root project. File ini mendefinisikan kredensial database (master PostgreSQL, staging SQL Server, warehouse PostgreSQL), port mapping untuk setiap service agar tidak bentrok dengan port lokal, konfigurasi NATS, serta port webserver Caddy. Template tersedia di `.env.example` yang bisa langsung disalin dan disesuaikan. Environment variables penting yang perlu diatur antara lain `MASTER_DB`, `MASTER_USER`, `MASTER_PASSWORD` untuk database utama, `STAGING_SA_PASSWORD` untuk SQL Server, `WEBSERVER_HTTP_HOST_PORT` untuk akses web, serta `HOST_IP` untuk binding address.

Selanjutnya, siapkan struktur direktori yang dibutuhkan Docker volumes — folder `storage/` untuk persistent data PostgreSQL dan NATS, serta folder `migrations/data/` untuk file CSV seed. Verifikasi bahwa Docker Compose tersedia dengan menjalankan `docker compose version`. Jika seluruh prasyarat terpenuhi, environment siap untuk memulai implementasi database.

---

## 2. Implementasi Database

Implementasi database dimulai dengan mendefinisikan skema PostgreSQL melalui file migrasi Goose. Migrasi pertama (`20260522040849_init.sql`) membuat fondasi seluruh sistem: 17 enum type untuk standarisasi nilai kolom seperti `status_verifikasi` (Pending, Aktif, Ditolak), `status_gizi` (Gizi Baik, Gizi Kurang, Gizi Buruk, Risiko Gizi Lebih, Gizi Lebih, Obesitas), `status_stunting` (Normal, Berisiko Stunting, Stunting, Stunting Berat), `status_kehamilan` (Trimester 1 hingga Keguguran), `status_rujukan` (Diajukan hingga Selesai), dan sepuluh enum lainnya yang menjadi constraint validasi di level database.

Setelah enum, dibuat 21 tabel yang dikelompokkan menjadi empat kategori. Kategori master mencakup `lokasi` (hierarki wilayah self-referencing), `user_account` (tabel utama seluruh pengguna dengan soft-delete via `is_deleted`), serta lima tabel role-specific yang memperluas `user_account` melalui foreign key: `dinas_kesehatan`, `bidan`, `kader_posyandu`, `pasien`, dan `fasilitas_kesehatan`. Setiap tabel role memiliki trigger `BEFORE UPDATE` yang otomatis mengisi `updated_at` melalui fungsi `update_updated_at_column()`. Kategori referensi mencakup `pendidikan`, `pekerjaan`, dan `kategori_pendapatan` sebagai data master demografi. Kategori pasien mencakup `ibu_hamil` (melacak riwayat kehamilan per pasien) dan `anak` (data balita dengan relasi ke wali). Kategori transaksi mencakup `jadwal_imunisasi`, `hasil_pemeriksaan` (pencatatan antropometri dan status gizi), `tindak_lanjut` (rujukan atau kontrol ulang), `rujukan` (tracking status rujukan ke faskes), `notifikasi` (sistem notifikasi pengguna), `user_session` (manajemen sesi JWT), dan `audit_log` (jejak audit semua aktivitas).

Sebanyak 25 foreign key constraint didefinisikan untuk menjaga integritas referensial, termasuk cascade delete dari `user_account` ke tabel role, dari `pasien` ke `ibu_hamil` dan `anak`, serta dari `tindak_lanjut` ke `rujukan`. Satu fungsi trigger `update_updated_at_column()` ditulis dalam plpgsql untuk otomatisasi timestamp, dipasang di 15 tabel melalui trigger `BEFORE UPDATE`. Tabel `notifikasi`, `lokasi`, `pendidikan`, `pekerjaan`, `kategori_pendapatan`, dan `audit_log` tidak memerlukan trigger ini karena sifatnya yang append-only atau referensi statis.

Setelah struktur tabel siap, migrasi kedua (`20260522040850_seed_master.sql`) memasukkan data awal dari 20 file CSV. Data mencakup hierarki lokasi (provinsi, kabupaten, kecamatan, kelurahan), data referensi pendidikan dan pekerjaan, data posyandu, bidan, kader, serta data faskes. Migrasi ketiga (`20260522040851_seed_login_accounts.sql`) membuat akun dummy untuk 6 role — Admin, Bidan, Kader, Pasien, Ibu Hamil, dan User — lengkap dengan password hash untuk keperluan testing.

Untuk performa query, 39 index ditambahkan melalui dua file migrasi terpisah. Index prioritas tinggi mencakup hot-path columns seperti `is_deleted` dan `status_verifikasi` di `user_account`, foreign key di `pasien` ke `posyandu`, `jadwal_imunisasi` ke `pasien`, `hasil_pemeriksaan` ke `jadwal_imunisasi` dan `id_petugas_input`, serta compound index untuk `notifikasi` (user + status_baca). Index prioritas menengah mencakup seluruh foreign key yang sering di-JOIN dan kolom yang sering di-filter atau di-sort.

Untuk kebutuhan dashboard yang memerlukan agregasi data besar secara cepat, dibuat 9 regular view di migrasi dashboard, kemudian segera digantikan dengan 11 materialized view di migrasi berikutnya. Materialized view ini menyimpan hasil pre-komputasi sehingga query dashboard tidak perlu melakukan JOIN dan agregasi berat setiap kali dipanggil. Kesebelas materialized view tersebut adalah: `mv_dashboard_stats` (ringkasan dashboard), `mv_dashboard_distribusi_gizi`, `mv_dashboard_tren_stunting`, `mv_dashboard_stunting_per_wilayah`, `mv_dashboard_kehadiran_bulanan`, `mv_dashboard_jadwal_terdekat`, `mv_public_stats` (statistik halaman publik), `mv_riwayat_pemeriksaan` (timeline per pasien), `mv_tumbuh_kembang` (grafik pertumbuhan), `mv_ibu_hamil_stats`, dan `mv_ibu_hamil_per_wilayah`. Karena materialized view tidak memiliki unique index, operasi refresh bersifat blocking dan perlu dijadwalkan secara periodik.

Seluruh migrasi dijalankan otomatis saat container `migration-master` startup melalui Goose. Container ini menunggu health check PostgreSQL master berhasil, lalu mengeksekusi `goose up` dan tetap berjalan untuk keperluan migrasi lanjutan.

---

## 3. Implementasi Backend

Backend dibangun dengan bahasa Go menggunakan framework Huma v2 yang menghasilkan dokumentasi OpenAPI secara otomatis dari kode, dipadukan dengan Chi router untuk performa tinggi. Arsitektur backend mengikuti pola clean architecture dengan tiga layer utama: domain (interface + DTO), feature/repository (data access), dan feature/service (business logic), serta feature/handler (HTTP handler).

Langkah pertama adalah code generation menggunakan Go-Jet. Tool `jetgen` membaca skema database PostgreSQL yang sudah berjalan, lalu menghasilkan type-safe Go struct untuk seluruh tabel (`table/`), model data (`model/`), dan enum (`enum/`). Ini memastikan bahwa setiap query SQL yang ditulis di repository terverifikasi tipe pada waktu kompilasi, mencegah kesalahan seperti typo nama kolom atau ketidakcocokan tipe data.

Layer domain mendefinisikan kontrak melalui interface dan DTO. Sebelas domain package dibuat: `auth`, `userAccount`, `pasien`, `pemeriksaan`, `imunisasi`, `artikel`, `tindaklanjut`, `dashboard`, `lokasi`, `faskes`, dan `notification`. Setiap domain memiliki file `dto.go` yang mendefinisikan struktur request dan response lengkap dengan tag validasi Huma (format, minLength, maxLength, enum, pattern). Interface `Handler`, `Service`, dan `Repo` didefinisikan di masing-masing domain untuk memungkinkan dependency injection dan memudahkan unit testing dengan mock.

Layer repository (`feature/*/repository.go`) mengimplementasikan seluruh akses database. Sebagian besar repository menggunakan Go-Jet query builder untuk operasi standar (INSERT, UPDATE, DELETE, SELECT sederhana), sementara query kompleks dengan banyak JOIN, subquery LATERAL, atau agregasi menggunakan raw SQL via `pgxpool.Pool`. Pola yang menonjol termasuk soft-delete (UPDATE `is_deleted = true, deleted_at = NOW()`) di tabel `user_account`, `pasien`, `ibu_hamil`, `anak`, dan tabel role; hard delete di `jadwal_imunisasi`, `hasil_pemeriksaan`, dan `artikel`; serta transaksi database untuk operasi yang melibatkan banyak tabel seperti delete pasien yang harus menghapus `anak` dan `ibu_hamil` secara atomik. Query kompleks menggunakan `LEFT JOIN LATERAL` untuk mendapatkan record terbaru (seperti hasil pemeriksaan terakhir per pasien), `COUNT(DISTINCT ...) FILTER (WHERE ...)` untuk agregasi pivot, `EXTRACT(YEAR FROM AGE(...))` untuk kalkulasi usia, dan CASE expression untuk kategorisasi deadline.

Layer service (`feature/*/service.go`) berisi business logic. Service auth menangani hashing password (bcrypt), validasi login (email atau NIK), pembuatan JWT access dan refresh token, session management, serta verifikasi akun oleh admin. Service pasien menangani pendaftaran ibu hamil dan anak dengan validasi relasi (wali, posyandu, ibu hamil), serta ownership check — hanya wali atau petugas yang bisa mengakses data pasien tertentu. Service pemeriksaan mengkalkulasi status gizi berdasarkan berat badan, tinggi badan, dan usia (WHO growth standards) serta menentukan status stunting. Service mengelola workflow verifikasi: kader menginput data, bidan memverifikasi. Service dashboard seluruhnya membaca dari materialized view untuk performa instant.

Layer handler (`feature/*/handler.go`) menjadi jembatan antara HTTP request dan service layer. Setiap handler method menerima input yang sudah divalidasi oleh Huma (path params, query params, request body) dan mengembalikan response terstandarisasi melalui `httputils.APIResponseOutput`. Error dari service layer dikonversi menjadi HTTP status code yang sesuai via `errorutils.ToHumaError`. Auth handler memiliki penanganan khusus untuk cookie — saat login dan refresh, dua HttpOnly cookie (`access_token` dan `refresh_token`) di-set dengan flag SameSite=Strict, sementara saat logout cookie dikosongkan.

Middleware menjadi lapisan keamanan utama. `AuthAccessMiddleware` memvalidasi JWT access token dari cookie, mengekstrak claims, dan menyimpannya di context. `AuthRefreshMiddleware` melakukan hal serupa untuk refresh token. `RequireRole` memeriksa apakah user memiliki role yang sesuai (ADMIN, BIDAN, KADER, DINKES, SUPER_ADMIN). Middleware rate limiter membatasi jumlah request per IP, dan middleware IP ban memblokir IP yang terdeteksi melakukan serangan.

Seluruh route didaftarkan di `backend/internal/api/v1/routes.go` menggunakan hierarki group bersarang. Group paling luar (`publicGroup`) tidak memerlukan autentikasi. `nonAuthenticatedOnlyGroup` khusus untuk register dan login (menolak user yang sudah login). `authAccess` dan `authRefresh` memvalidasi token. Di bawah `authAccess`, `userGroup` memberikan akses ke semua user terautentikasi. `adminGroup` membatasi ke role admin (Bidan, Kader, Dinkes). `bidanGroup`, `kaderGroup`, `dinkesGroup`, dan `superAdminGroup` memberikan akses spesifik per role. Total terdapat sekitar 70 endpoint yang tersebar di 11 domain fungsional.

---

## 4. Implementasi Frontend

Frontend dibangun dengan React 19 dan TypeScript 6 di atas Vite 8 untuk development experience yang cepat dengan Hot Module Replacement. Tailwind CSS 4 digunakan untuk styling utility-first, React Router 7 untuk client-side routing, Zustand 5 untuk state management global yang ringan, dan Leaflet untuk komponen peta di dashboard.

Langkah pertama adalah mendefinisikan type entity di `types/entities.ts` yang merepresentasikan seluruh domain model dari backend — `PasienListItem`, `PasienDetail`, `IbuHamilData`, `AnakData`, `PemeriksaanDetail`, `ImunisasiDetail`, `ArtikelDetail`, `TindakLanjutRequest`, `NotifikasiResponse`, dan lainnya. Type ini menjadi kontrak antara frontend dan backend. File `types/api.ts` mendefinisikan struktur response API yang spesifik untuk dashboard dan statistik. HTTP client di `lib/api.ts` mengimplementasikan fetch wrapper dengan interceptors untuk attach credential, auto-refresh token saat access token expired, dan response caching.

Context provider menjadi tulang punggung autentikasi frontend. `AuthContext` menyimpan state user (claims JWT), menyediakan fungsi login, register, logout, dan refresh, serta secara otomatis mendeteksi status autentikasi dari cookie. `NotificationContext` menangani koneksi real-time atau polling untuk notifikasi baru. State global aplikasi yang lebih luas (seperti filter dashboard, preferensi tampilan) dikelola melalui Zustand store di `useAppStore`.

Setiap screen diorganisir berdasarkan role dan fungsionalitas. Screen login dan register menangani form dengan validasi client-side. Screen dashboard memiliki tampilan berbeda per role: Bidan dan Kader melihat statistik monitoring, jadwal posyandu, dan daftar pasien yang perlu ditindaklanjuti; Ibu dan Wali melihat data anak mereka (jadwal imunisasi, riwayat pemeriksaan, grafik tumbuh kembang); Dinkes melihat agregat data per wilayah dan laporan. Screen monitoring menampilkan daftar pasien dengan pencarian dan filter, detail pasien dengan tab (data diri, riwayat pemeriksaan, imunisasi, tindak lanjut). Screen jadwal imunisasi menyediakan kalender dan form penjadwalan. Screen pemeriksaan menampilkan form input antropometri lengkap dengan kalkulator status gizi otomatis. Screen tindak lanjut memungkinkan bidan membuat rujukan atau jadwal kontrol ulang. Screen edukasi adalah mini-CMS untuk artikel kesehatan dengan workflow Draft → Submit → Review → Publish. Screen user management (super admin) menampilkan tabel user dengan filter role dan status verifikasi.

Komponen shared dibangun untuk digunakan ulang di seluruh screen. Header menampilkan navigasi utama yang berubah sesuai role user. Sidebar menampilkan menu kontekstual. Paginator komponen menangani navigasi halaman dengan informasi total data. Toast komponen menampilkan notifikasi sukses/gagal. Modal komponen digunakan untuk konfirmasi aksi (delete, verifikasi, review artikel). Custom hooks seperti `usePaginator` menyederhanakan logika paginasi di berbagai screen.

Routing di `App.tsx` menggunakan React Router dengan nested routes. Route publik (landing page, artikel, login, register) bisa diakses tanpa autentikasi. Route terproteksi memerlukan user login dan role tertentu — misalnya `/dashboard` hanya untuk user terautentikasi, `/user-management` hanya untuk super admin, `/monitoring` untuk bidan dan kader. Guard diimplementasikan melalui wrapper component yang memeriksa AuthContext sebelum me-render children.

---

## 5. Implementasi Message Broker dan Cache

NATS dipilih sebagai message broker untuk sistem notifikasi real-time dan penyimpanan key-value untuk data yang memerlukan akses cepat dengan TTL. NATS berjalan dalam mode JetStream yang menyediakan persistensi pesan dan kemampuan streaming.

Publisher notifikasi diimplementasikan di backend melalui interface `notificationDomain.Publisher` dan implementasi konkret berbasis NATS. Setiap kali terjadi event penting — hasil pemeriksaan baru diinput, jadwal imunisasi dibuat, rujukan diajukan, status rujukan berubah, artikel dipublikasikan — publisher akan mengirim pesan ke subject NATS yang sesuai. Subscriber (listener) di sisi backend kemudian mengonsumsi pesan tersebut, membuat record di tabel `notifikasi`, dan mengirimkan ke user yang relevan melalui mekanisme yang sesuai (bisa via polling dari frontend atau WebSocket). Setiap notifikasi memiliki tipe (`Pemeriksaan`, `Imunisasi`, `Rujukan`, `Edukasi`, `Pengingat`) yang menentukan cara penanganannya di frontend — apakah menampilkan badge, toast, atau redirect ke halaman terkait.

NATS Key-Value store digunakan untuk dua fitur keamanan. Pertama, JWT blacklist: saat user logout atau refresh token, JTI (JWT ID) dari token yang dicabut disimpan di NATS KV dengan TTL sesuai masa berlaku token. Middleware `AuthAccessMiddleware` memeriksa KV store setiap kali memvalidasi token — jika JTI ditemukan dalam blacklist, akses ditolak. Kedua, IP ban: middleware rate limiter melacak jumlah percobaan gagal per IP. Jika melebihi threshold, IP tersebut di-ban dengan menyimpan `BanInfo` (reason, attempts, bannedAt, expiresAt) di NATS KV. Setiap request berikutnya dari IP tersebut akan ditolak sebelum mencapai handler.

NATS juga menyediakan monitoring dashboard melalui NATS NUI (NATS User Interface) yang berjalan di port 31311, memungkinkan developer dan operator memantau jumlah koneksi, throughput pesan, status JetStream, dan isi KV store secara real-time.

---

## 6. Integrasi Sistem

Integrasi seluruh komponen sistem diatur melalui Docker Compose dengan 16 service yang dikelompokkan dalam profile. Profile `dev` menjalankan seluruh service untuk development. Profile `information_system` menjalankan komponen inti: master database, backend, frontend, webserver (Caddy), NATS, dan NUI. Profile `data_management` menjalankan komponen data: master, staging, warehouse, Hop ETL, dan migrasi. Pemisahan profile ini memungkinkan developer hanya menjalankan komponen yang dibutuhkan, menghemat resource.

Caddy berperan sebagai reverse proxy dan API gateway. Dikonfigurasi melalui `Caddyfile`, Caddy menerima seluruh traffic HTTP/HTTPS dan merutekan berdasarkan path: request ke `/api/*` diteruskan ke backend container (port 8080 internal), request lainnya ke frontend container (Vite dev server). Caddy juga menangani SSL/TLS secara otomatis melalui fitur auto-HTTPS, serta menyediakan admin API untuk monitoring dan konfigurasi runtime.

Apache Hop digunakan untuk ETL (Extract, Transform, Load) dari database staging (SQL Server) ke database master (PostgreSQL). Pipeline ETL terdiri dari workflow utama (`AllPipeline.hwf`) yang mengorkestrasi 20 pipeline tingkat tabel (`.hpl` files). Setiap pipeline membaca data dari tabel staging yang sesuai, melakukan transformasi (konversi tipe data, pemetaan enum, deduplikasi), dan menulis ke tabel master. Workflow `RunAllIfEmpty.hwf` secara otomatis mendeteksi apakah tabel master masih kosong dan hanya menjalankan ETL jika diperlukan, berguna untuk initial load. Workflow `TruncateAndRunAll.hwf` untuk full reload. Konfigurasi koneksi database disimpan di `metadata/rdbms/` (Master.json, Staging.json, Warehouse.json) dan environment variables di `.env.json`.

Health check antar service memastikan urutan startup yang benar. Container `migration-master` menunggu `master` healthy sebelum menjalankan Goose. Backend menunggu `migration-master` selesai dan `nats` healthy. Frontend menunggu backend ready. Volume Docker digunakan untuk mount source code (hot reload di development), persistent data (database files), dan package cache (Go modules, Node modules) untuk mempercepat build.

---

## 7. Testing dan Validasi

Testing backend dilakukan dalam dua tingkatan. Unit test menggunakan standard library Go `testing` yang menguji service layer secara terisolasi. Mock repository dibuat menggunakan interface yang sudah didefinisikan di domain layer, memungkinkan pengujian business logic tanpa ketergantungan database. Test case mencakup validasi input (password terlalu pendek, NIK tidak 16 digit), business rule (tidak bisa verifikasi user yang sudah aktif, tidak bisa daftar pasien ke posyandu yang tidak ada), dan edge case (login dengan email dan NIK sekaligus, refresh token expired).

Integration test menggunakan testcontainers untuk menjalankan PostgreSQL instance sementara. Sebelum test berjalan, migrasi Goose dijalankan untuk membuat skema, dan seed data minimal dimasukkan. Test mengirim HTTP request sungguhan ke handler (menggunakan `httptest`) dan memverifikasi response code, response body, serta efek samping di database. Test case integration mencakup flow bisnis end-to-end: register → verifikasi → login → akses resource → CRUD → logout. Setiap domain memiliki file integration test terpisah.

CI/CD diimplementasikan melalui GitHub Actions. Dua workflow didefinisikan: `backend-unit-test.yml` menjalankan unit test setiap push dan pull request ke branch main, `backend-integration-test.yml` menjalankan integration test dengan Docker Compose (profile test) yang menyalakan PostgreSQL, menjalankan migrasi, lalu mengeksekusi test. Kedua workflow harus lulus sebelum PR bisa di-merge.

Selain automated test, validasi manual dilakukan menggunakan script Python `seed_dummy.py` yang menghasilkan data dummy dalam jumlah besar untuk stress testing dan visualisasi dashboard. Script `truncate_and_seed.sql` memungkinkan developer mereset database ke state awal dengan cepat selama development.

---

## 8. Deployment

Deployment production menggunakan Docker Compose dengan file `compose.prod.yaml` yang berbeda dari development. Setiap service memiliki konfigurasi production: resource limit (CPU dan memory), restart policy (`unless-stopped`), dan volume untuk persistent data. Secret seperti password database dan JWT signing key tidak lagi di-hardcode di `.env` melainkan di-inject melalui Docker secrets atau environment variable dari orchestrator.

Backend di-build menggunakan multi-stage Containerfile. Stage pertama menggunakan Go image untuk kompilasi binary statis dengan flags `-ldflags="-s -w"` untuk mengurangi ukuran. Stage kedua menggunakan Alpine minimal sebagai runtime image, hanya berisi binary hasil kompilasi. Ini menghasilkan image yang kecil dan aman. Frontend di-build dengan `bun run build` yang menghasilkan static files di folder `dist/`, kemudian disajikan oleh Caddy sebagai static file server — tidak memerlukan Node.js runtime di production.

Caddy di production menangani auto-HTTPS dengan Let's Encrypt, menyediakan sertifikat SSL/TLS gratis yang diperbarui otomatis. Konfigurasi production Caddyfile mengarahkan domain sebenarnya (misalnya `sigizi.kemkes.go.id`) ke service yang sesuai, dengan security header tambahan seperti HSTS, CSP, dan X-Frame-Options.

Database production memerlukan strategi backup reguler. PostgreSQL menyediakan `pg_dump` untuk full backup dan WAL (Write-Ahead Log) archiving untuk point-in-time recovery. Backup dijadwalkan melalui cron job di host atau menggunakan pg_cron extension. Materialized view perlu di-refresh secara periodik — dapat menggunakan pg_cron untuk menjadwalkan `REFRESH MATERIALIZED VIEW CONCURRENTLY` (setelah menambahkan unique index) setiap jam atau setiap hari tergantung kebutuhan kesegaran data.

Monitoring production mengandalkan health check endpoint di backend (`/health`) dan database health check bawaan Docker. Log dari seluruh container dikumpulkan melalui Docker logging driver dan bisa di-forward ke sistem monitoring terpusat seperti Grafana Loki atau ELK stack. NATS monitoring dashboard menyediakan insight real-time tentang performa message broker.

Untuk scaling, arsitektur microservice memungkinkan scaling horizontal komponen tertentu. Backend bisa di-replikasi (multiple container) di belakang load balancer karena bersifat stateless (session disimpan di database via `user_session`, bukan di memory). Database menggunakan replikasi PostgreSQL (primary-standby) untuk high availability. NATS cluster bisa diperluas untuk meningkatkan throughput message.
