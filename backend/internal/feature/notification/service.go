package notification

import (
	"context"
	"fmt"
	"net/http"

	notificationDomain "github.com/stringptr/SiGizi/backend/internal/domain/notification"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Service struct {
	Repo notificationDomain.Repo
}

func NewService(repo notificationDomain.Repo) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetNotifikasi(ctx context.Context, idUser int32, input *notificationDomain.NotifikasiListInput) (*notificationDomain.NotifikasiListData, *errorutils.Error) {
	limit := input.PerPage
	offset := (input.Page - 1) * input.PerPage

	notif, err := s.Repo.GetByUserID(ctx, idUser, input.Search, limit, offset)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	total, err := s.Repo.CountByUserID(ctx, idUser, input.Search)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	lastPage := total / input.PerPage
	if total%input.PerPage != 0 {
		lastPage++
	}

	items := make([]notificationDomain.NotifikasiItem, 0, len(notif))
	for _, n := range notif {
		items = append(items, toNotifikasiItem(n))
	}

	return &notificationDomain.NotifikasiListData{
		Notifikasi: items,
		Meta: notificationDomain.Meta{
			CurrentPage: input.Page,
			PerPage:     input.PerPage,
			Total:       total,
			LastPage:    lastPage,
		},
	}, nil
}

func (s *Service) GetNotifikasiDetail(ctx context.Context, idUser int32, idNotifikasi int32) (*notificationDomain.NotifikasiDetail, *errorutils.Error) {
	notif, err := s.Repo.GetByID(ctx, idNotifikasi, idUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if notif == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Notifikasi tidak ditemukan."}
	}

	detail := &notificationDomain.NotifikasiDetail{
		IDNotifikasi:   notif.IDNotifikasi,
		Judul:          notif.Judul,
		Pesan:          notif.Pesan,
		TipeNotifikasi: string(notif.TipeNotifikasi),
		StatusBaca:     notif.StatusBaca,
		TanggalKirim:   notif.TanggalKirim,
	}

	tipeStr := string(notif.TipeNotifikasi)
	switch tipeStr {
	case "Rujukan":
		detail.Aksi = &notificationDomain.AksiItem{
			Label: "Lihat Detail Rujukan",
			URL:   fmt.Sprintf("/tindak-lanjut/%d", notif.IDNotifikasi),
		}
	case "Pemeriksaan":
		detail.Aksi = &notificationDomain.AksiItem{
			Label: "Lihat Detail Pemeriksaan",
			URL:   fmt.Sprintf("/monitoring/pemeriksaan/%d", notif.IDNotifikasi),
		}
	case "Imunisasi":
		detail.Aksi = &notificationDomain.AksiItem{
			Label: "Lihat Jadwal Imunisasi",
			URL:   fmt.Sprintf("/imunisasi/%d", notif.IDNotifikasi),
		}
	case "Edukasi":
		detail.Aksi = &notificationDomain.AksiItem{
			Label: "Baca Artikel",
			URL:   fmt.Sprintf("/artikel/%d", notif.IDNotifikasi),
		}
	}

	return detail, nil
}

func (s *Service) MarkRead(ctx context.Context, idUser int32, idNotifikasi int32) (*notificationDomain.MarkReadResponse, *errorutils.Error) {
	notif, err := s.Repo.GetByID(ctx, idNotifikasi, idUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if notif == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Notifikasi tidak ditemukan."}
	}

	err = s.Repo.MarkRead(ctx, idNotifikasi, idUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	return &notificationDomain.MarkReadResponse{
		IDNotifikasi: idNotifikasi,
		StatusBaca:   true,
	}, nil
}

func (s *Service) MarkAllRead(ctx context.Context, idUser int32) (*notificationDomain.MarkAllReadResponse, *errorutils.Error) {
	count, err := s.Repo.MarkAllRead(ctx, idUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	return &notificationDomain.MarkAllReadResponse{
		JumlahDiperbarui: count,
		Status:           "SEMUA_DIBACA",
	}, nil
}

func (s *Service) GetBidanDashboard(ctx context.Context, idUser int32) (*notificationDomain.BidanNotificationResponse, *errorutils.Error) {
	concreteRepo, ok := s.Repo.(*Repo)
	if !ok {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	bidanID, err := concreteRepo.getBidanIDByUserID(ctx, idUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if bidanID == 0 {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Tidak terdapat data pusat notifikasi bidan."}
	}

	jadwalKontrol, err := concreteRepo.countBidanJadwalKontrol(ctx, bidanID)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	risikoStuntingCount, err := concreteRepo.countBidanRisikoStunting(ctx, bidanID)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	rujukanMendesakCount, err := concreteRepo.countBidanRujukanMendesak(ctx, bidanID)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	risikoStuntingList, err := concreteRepo.getBidanRisikoStuntingList(ctx, bidanID)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	jadwalMonitoringList, err := concreteRepo.getBidanJadwalMonitoringList(ctx, bidanID)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	rujukanMendesakList, err := concreteRepo.getBidanRujukanMendesakList(ctx, bidanID)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	laporanBulanan, err := concreteRepo.getBidanLaporanBulanan(ctx, bidanID)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	risikoStuntingItems := make([]notificationDomain.PasienRisikoItem, 0, len(risikoStuntingList))
	for _, r := range risikoStuntingList {
		risikoStuntingItems = append(risikoStuntingItems, notificationDomain.PasienRisikoItem{
			IDPasien:          r.IDPasien,
			NamaPasien:        r.NamaPasien,
			StatusGizi:        r.StatusGizi,
			StatusStunting:    r.StatusStunting,
			TanggalMonitoring: r.TanggalMonitoring.Format("2006-01-02 15:04:05"),
		})
	}

	jadwalMonitoringItems := make([]notificationDomain.JadwalMonitoringItem, 0, len(jadwalMonitoringList))
	for _, j := range jadwalMonitoringList {
		jadwalMonitoringItems = append(jadwalMonitoringItems, notificationDomain.JadwalMonitoringItem{
			IDPasien:      j.IDPasien,
			NamaPasien:    j.NamaPasien,
			JadwalKontrol: j.JadwalKontrol.Format("2006-01-02 15:04:05"),
			Status:        j.Status,
		})
	}

	rujukanMendesakItems := make([]notificationDomain.RujukanMendesakItem, 0, len(rujukanMendesakList))
	for _, r := range rujukanMendesakList {
		rujukanMendesakItems = append(rujukanMendesakItems, notificationDomain.RujukanMendesakItem{
			IDRujukan:      r.IDRujukan,
			NamaPasien:     r.NamaPasien,
			StatusRujukan:  r.StatusRujukan,
			TanggalRujukan: r.TanggalRujukan.Format("2006-01-02 15:04:05"),
		})
	}

	return &notificationDomain.BidanNotificationResponse{
		Statistik: notificationDomain.NotifikasiBidanStats{
			JadwalKontrol:   jadwalKontrol,
			RisikoStunting:  risikoStuntingCount,
			RujukanMendesak: rujukanMendesakCount,
		},
		RisikoStunting:   risikoStuntingItems,
		JadwalMonitoring: jadwalMonitoringItems,
		RujukanMendesak:  rujukanMendesakItems,
		LaporanBulanan: notificationDomain.LaporanBulanan{
			Bulan:                  laporanBulanan.Bulan,
			JumlahPasienMonitoring: laporanBulanan.JumlahPasienMonitoring,
			JumlahPasienDirujuk:    laporanBulanan.JumlahPasienDirujuk,
		},
	}, nil
}

func (s *Service) GetStatistics(ctx context.Context, idUser int32) (*notificationDomain.NotificationStats, *errorutils.Error) {
	concreteRepo, ok := s.Repo.(*Repo)
	if !ok {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	bidanID, err := concreteRepo.getBidanIDByUserID(ctx, idUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	jadwalKontrol, err := concreteRepo.countBidanJadwalKontrol(ctx, bidanID)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	risikoStunting, err := concreteRepo.countBidanRisikoStunting(ctx, bidanID)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	rujukanMendesak, err := concreteRepo.countBidanRujukanMendesak(ctx, bidanID)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	unread, err := s.Repo.CountUnreadByUserID(ctx, idUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	return &notificationDomain.NotificationStats{
		JadwalUlang:           jadwalKontrol,
		RujukanMendesak:       rujukanMendesak,
		RisikoStunting:        risikoStunting,
		NotifikasiBelumDibaca: unread,
	}, nil
}

func (s *Service) GetActivity(ctx context.Context, idUser int32) (*notificationDomain.NotificationActivity, *errorutils.Error) {
	concreteRepo, ok := s.Repo.(*Repo)
	if !ok {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	today, err := concreteRepo.getTodayNotifikasi(ctx, idUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	yesterday, err := concreteRepo.getYesterdayNotifikasi(ctx, idUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	hariIni := make([]notificationDomain.AktivitasItem, 0, len(today))
	for _, n := range today {
		hariIni = append(hariIni, toAktivitasItem(n))
	}

	kemarin := make([]notificationDomain.AktivitasItem, 0, len(yesterday))
	for _, n := range yesterday {
		kemarin = append(kemarin, toAktivitasItem(n))
	}

	return &notificationDomain.NotificationActivity{
		HariIni: hariIni,
		Kemarin: kemarin,
	}, nil
}

func toNotifikasiItem(n *model.Notifikasi) notificationDomain.NotifikasiItem {
	return notificationDomain.NotifikasiItem{
		IDNotifikasi:   n.IDNotifikasi,
		Judul:          n.Judul,
		Pesan:          n.Pesan,
		TipeNotifikasi: string(n.TipeNotifikasi),
		StatusBaca:     n.StatusBaca,
		TanggalKirim:   n.TanggalKirim,
	}
}

func toAktivitasItem(n *model.Notifikasi) notificationDomain.AktivitasItem {
	status := "terbaca"
	if !n.StatusBaca {
		status = "baru"
	}
	return notificationDomain.AktivitasItem{
		IDNotifikasi: n.IDNotifikasi,
		Judul:        n.Judul,
		Status:       status,
		Timestamp:    n.TanggalKirim,
	}
}
