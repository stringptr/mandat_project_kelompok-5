package dashboard

import (
	"context"
	"net/http"

	dashboardDomain "github.com/stringptr/SiGizi/backend/internal/domain/dashboard"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/pagination"
)

type Service struct {
	Repo dashboardDomain.Repo
}

func NewService(repo dashboardDomain.Repo) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetDashboardStats(ctx context.Context) (*dashboardDomain.DashboardStatsResponse, *errorutils.Error) {
	row, err := s.Repo.GetDashboardStats(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	return &dashboardDomain.DashboardStatsResponse{
		TotalPasien:       row.TotalPasien,
		PerluVerifikasi:   row.PerluVerifikasi,
		TindakLanjut:      row.TindakLanjut,
		KasusStunting:     row.KasusStunting,
		JadwalPosyandu:    row.JadwalPosyandu,
		TotalBalita:       row.TotalBalita,
		CakupanPersentase: row.CakupanPersentase,
	}, nil
}

func (s *Service) GetDistribusiGizi(ctx context.Context) (*dashboardDomain.DistribusiGiziResponse, *errorutils.Error) {
	rows, err := s.Repo.GetDistribusiGizi(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]dashboardDomain.DistribusiGiziItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dashboardDomain.DistribusiGiziItem{
			StatusGizi: r.StatusGizi,
			Jumlah:     r.Jumlah,
		})
	}

	return &dashboardDomain.DistribusiGiziResponse{Distribusi: items}, nil
}

func (s *Service) GetTrenStunting(ctx context.Context) (*dashboardDomain.TrenResponse, *errorutils.Error) {
	rows, err := s.Repo.GetTrenStunting(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]dashboardDomain.TrenItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dashboardDomain.TrenItem{
			Bulan:  r.Bulan,
			Jumlah: r.Jumlah,
		})
	}

	return &dashboardDomain.TrenResponse{Tren: items}, nil
}

func (s *Service) GetStuntingPerWilayah(ctx context.Context) (*dashboardDomain.StuntingWilayahResponse, *errorutils.Error) {
	rows, err := s.Repo.GetStuntingPerWilayah(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]dashboardDomain.StuntingWilayahItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dashboardDomain.StuntingWilayahItem{
			NamaWilayah: r.NamaWilayah,
			Prevalensi:  r.Prevalensi,
			JumlahKasus: r.JumlahKasus,
			TotalBalita: r.TotalBalita,
			Level:       r.Level,
		})
	}

	return &dashboardDomain.StuntingWilayahResponse{Wilayah: items}, nil
}

func (s *Service) GetKehadiranBulanan(ctx context.Context) (*dashboardDomain.TrenResponse, *errorutils.Error) {
	rows, err := s.Repo.GetKehadiranBulanan(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]dashboardDomain.TrenItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dashboardDomain.TrenItem{
			Bulan:  r.Bulan,
			Jumlah: r.Jumlah,
		})
	}

	return &dashboardDomain.TrenResponse{Tren: items}, nil
}

func (s *Service) GetPublicStats(ctx context.Context) (*dashboardDomain.PublicStatsResponse, *errorutils.Error) {
	row, err := s.Repo.GetPublicStats(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	return &dashboardDomain.PublicStatsResponse{
		TotalPasien:    row.TotalPasien,
		BalitaDipantau: row.BalitaDipantau,
		KasusStunting:  row.KasusStunting,
		TotalArtikel:   row.TotalArtikel,
	}, nil
}

func (s *Service) GetRiwayat(ctx context.Context, idPasien int32) (*dashboardDomain.RiwayatResponse, *errorutils.Error) {
	rows, err := s.Repo.GetRiwayat(ctx, idPasien)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if len(rows) == 0 {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Riwayat pemeriksaan tidak ditemukan."}
	}

	items := make([]dashboardDomain.RiwayatItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dashboardDomain.RiwayatItem{
			Tanggal:     r.Tanggal,
			BeratBadan:  r.BeratBadan,
			TinggiBadan: r.TinggiBadan,
			StatusGizi:  r.StatusGizi,
			Catatan:     r.Catatan,
			Petugas:     r.Petugas,
		})
	}

	return &dashboardDomain.RiwayatResponse{Riwayat: items}, nil
}

func (s *Service) GetTumbuhKembang(ctx context.Context, idPasien int32) (*dashboardDomain.TumbuhKembangResponse, *errorutils.Error) {
	rows, err := s.Repo.GetTumbuhKembang(ctx, idPasien)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if len(rows) == 0 {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data tumbuh kembang tidak ditemukan."}
	}

	items := make([]dashboardDomain.TumbuhKembangItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dashboardDomain.TumbuhKembangItem{
			Bulan:      r.Bulan,
			BeratBadan: r.BeratBadan,
			TinggiBadan: r.TinggiBadan,
		})
	}

	return &dashboardDomain.TumbuhKembangResponse{Data: items}, nil
}

func (s *Service) GetJadwalTerdekat(ctx context.Context) (*dashboardDomain.JadwalTerdekatResponse, *errorutils.Error) {
	rows, err := s.Repo.GetJadwalTerdekat(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]dashboardDomain.JadwalTerdekatItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dashboardDomain.JadwalTerdekatItem{
			ID:            r.ID,
			NamaVaksin:    r.NamaVaksin,
			TanggalJadwal: r.TanggalJadwal,
			NamaPasien:    r.NamaPasien,
		})
	}

	return &dashboardDomain.JadwalTerdekatResponse{Jadwal: items}, nil
}

func (s *Service) GetIbuHamilStats(ctx context.Context) (*dashboardDomain.IbuHamilStatsResponse, *errorutils.Error) {
	row, err := s.Repo.GetIbuHamilStats(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	return &dashboardDomain.IbuHamilStatsResponse{
		TotalIbuHamil: row.TotalIbuHamil,
		Trimester1:    row.Trimester1,
		Trimester2:    row.Trimester2,
		Trimester3:    row.Trimester3,
		Melahirkan:    row.Melahirkan,
		Nifas:         row.Nifas,
		Keguguran:     row.Keguguran,
	}, nil
}

func (s *Service) GetIbuHamilPerWilayah(ctx context.Context) (*dashboardDomain.IbuHamilWilayahResponse, *errorutils.Error) {
	rows, err := s.Repo.GetIbuHamilPerWilayah(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]dashboardDomain.IbuHamilWilayahItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dashboardDomain.IbuHamilWilayahItem{
			NamaWilayah:   r.NamaWilayah,
			TotalIbuHamil: r.TotalIbuHamil,
			Trimester1:    r.Trimester1,
			Trimester2:    r.Trimester2,
			Trimester3:    r.Trimester3,
			Melahirkan:    r.Melahirkan,
			Nifas:         r.Nifas,
			Keguguran:     r.Keguguran,
		})
	}

	return &dashboardDomain.IbuHamilWilayahResponse{Wilayah: items}, nil
}

func (s *Service) GetSemuaPemeriksaan(ctx context.Context, req *dashboardDomain.GetAllPemeriksaanRequest) (*dashboardDomain.PemeriksaanListResponse, *errorutils.Error) {
	page := req.Page
	perPage := req.PerPage
	if page < 1 { page = 1 }
	if perPage < 1 || perPage > 50 { perPage = 20 }

	idBidan := req.IDBidan
	idPosyandu := req.IDPosyandu

	if req.IDKader > 0 {
		pid, err := s.Repo.GetPosyanduByKaderID(ctx, req.IDKader)
		if err == nil && pid > 0 {
			idPosyandu = pid
		}
	}

	rows, total, err := s.Repo.GetSemuaPemeriksaan(ctx, page, perPage, idBidan, idPosyandu)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]dashboardDomain.PemeriksaanItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dashboardDomain.PemeriksaanItem{
			IDHasilPemeriksaan: r.IDHasilPemeriksaan,
			IDJadwalImunisasi:  r.IDJadwalImunisasi,
			NamaVaksin:         r.NamaVaksin,
			NamaPasien:         r.NamaPasien,
			BeratBadan:         r.BeratBadan,
			TinggiBadan:        r.TinggiBadan,
			LingkarKepala:      r.LingkarKepala,
			TekananDarah:       r.TekananDarah,
			StatusStunting:     r.StatusStunting,
			StatusGizi:         r.StatusGizi,
			Catatan:            r.Catatan,
			Tanggal:            r.Tanggal,
			Petugas:            r.Petugas,
		})
	}

	lastPage := (total + perPage - 1) / perPage
	if lastPage < 1 { lastPage = 1 }
	return &dashboardDomain.PemeriksaanListResponse{
		Pemeriksaan: items,
		Meta:        pagination.Meta{CurrentPage: int32(page), PerPage: int32(perPage), Total: int32(total), LastPage: int32(lastPage)},
	}, nil
}
