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
	return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Tidak terdapat data pusat notifikasi bidan."}
}

func (s *Service) GetStatistics(ctx context.Context, idUser int32) (*notificationDomain.NotificationStats, *errorutils.Error) {
	return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data statistik notifikasi tidak ditemukan."}
}

func (s *Service) GetActivity(ctx context.Context, idUser int32) (*notificationDomain.NotificationActivity, *errorutils.Error) {
	return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Tidak terdapat aktivitas notifikasi."}
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
