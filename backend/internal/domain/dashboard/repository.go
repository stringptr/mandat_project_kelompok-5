package dashboard

import "context"

type DashboardStatsRow struct {
	TotalPasien       int32
	PerluVerifikasi   int32
	TindakLanjut      int32
	KasusStunting     int32
	JadwalPosyandu    int32
	TotalBalita       int32
	CakupanPersentase float64
}

type DistribusiGiziRow struct {
	StatusGizi string
	Jumlah     int32
}

type TrenRow struct {
	Bulan  string
	Jumlah int32
}

type StuntingWilayahRow struct {
	NamaWilayah string
	Prevalensi  float64
	JumlahKasus int32
	TotalBalita int32
	Level       string
}

type PublicStatsRow struct {
	TotalPasien    int32
	BalitaDipantau int32
	KasusStunting  int32
	TotalArtikel   int32
}

type RiwayatRow struct {
	Tanggal     string
	BeratBadan  float64
	TinggiBadan float64
	StatusGizi  string
	Catatan     *string
	Petugas     string
}

type TumbuhKembangRow struct {
	Bulan       string
	BeratBadan  float64
	TinggiBadan float64
}

type JadwalTerdekatRow struct {
	ID            int32
	NamaVaksin    string
	TanggalJadwal string
	NamaPasien    string
}

type IbuHamilStatsRow struct {
	TotalIbuHamil int32
	Trimester1    int32
	Trimester2    int32
	Trimester3    int32
	Melahirkan    int32
	Nifas         int32
	Keguguran     int32
}

type IbuHamilWilayahRow struct {
	NamaWilayah   string
	TotalIbuHamil int32
	Trimester1    int32
	Trimester2    int32
	Trimester3    int32
	Melahirkan    int32
	Nifas         int32
	Keguguran     int32
}

type PemeriksaanRow struct {
	IDHasilPemeriksaan int32
	IDJadwalImunisasi  int32
	NamaVaksin         string
	NamaPasien         string
	BeratBadan         float64
	TinggiBadan        float64
	LingkarKepala      float64
	TekananDarah       string
	StatusStunting     string
	StatusGizi         string
	Catatan            *string
	Tanggal            string
	Petugas            string
}

// MVRefresher is a minimal interface for asynchronously refreshing materialized views.
// It is injected into write-heavy services (pemeriksaan, tindaklanjut) so that dashboard
// statistics stay up-to-date after data mutations without blocking the HTTP response.
type MVRefresher interface {
	RefreshMaterializedViews(ctx context.Context) error
}

type Repo interface {
	GetDashboardStats(ctx context.Context) (*DashboardStatsRow, error)
	GetDistribusiGizi(ctx context.Context) ([]DistribusiGiziRow, error)
	GetTrenStunting(ctx context.Context) ([]TrenRow, error)
	GetStuntingPerWilayah(ctx context.Context) ([]StuntingWilayahRow, error)
	GetKehadiranBulanan(ctx context.Context) ([]TrenRow, error)
	GetPublicStats(ctx context.Context) (*PublicStatsRow, error)
	GetRiwayat(ctx context.Context, idPasien int32) ([]RiwayatRow, error)
	GetTumbuhKembang(ctx context.Context, idPasien int32) ([]TumbuhKembangRow, error)
	GetJadwalTerdekat(ctx context.Context) ([]JadwalTerdekatRow, error)
	GetIbuHamilStats(ctx context.Context) (*IbuHamilStatsRow, error)
	GetIbuHamilPerWilayah(ctx context.Context) ([]IbuHamilWilayahRow, error)
	GetSemuaPemeriksaan(ctx context.Context, page, perPage int, idBidan, idPosyandu int32) ([]PemeriksaanRow, int, error)
	GetPosyanduByKaderID(ctx context.Context, idKader int32) (int32, error)
	MVRefresher
}
