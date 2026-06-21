package tindaklanjut

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	tindaklanjutDomain "github.com/stringptr/SiGizi/backend/internal/domain/tindaklanjut"
	auditlogDomain "github.com/stringptr/SiGizi/backend/internal/domain/auditlog"
	notificationDomain "github.com/stringptr/SiGizi/backend/internal/domain/notification"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	"github.com/stringptr/SiGizi/backend/internal/pagination"
)

type Service struct {
	repo          tindaklanjutDomain.Repo
	auditRepo     auditlogDomain.Repo
	notifRepo     notificationDomain.Repo
	notifPublisher notificationDomain.Publisher
}

func NewService(repo tindaklanjutDomain.Repo, auditRepo auditlogDomain.Repo, notifRepo notificationDomain.Repo, notifPublisher notificationDomain.Publisher) *Service {
	return &Service{
		repo:          repo,
		auditRepo:     auditRepo,
		notifRepo:     notifRepo,
		notifPublisher: notifPublisher,
	}
}

func (s *Service) logAudit(ctx context.Context, endpoint string, tipeAktivitas model.TipeAktivitas, berhasil bool, tableName string, recordID string, detail string) {
	tipeAktor := model.TipeAktor_User
	s.auditRepo.Log(ctx, &model.AuditLog{
		TipeAktor:     &tipeAktor,
		TipeAktivitas: &tipeAktivitas,
		Berhasil:      &berhasil,
		Endpoint:      &endpoint,
		TableName:     &tableName,
		RecordID:      &recordID,
		Detail:        &detail,
	})
}

func (s *Service) GetPasienTindakLanjut(ctx context.Context, req *tindaklanjutDomain.GetPasienTindakLanjutRequest) (*tindaklanjutDomain.PasienTindakLanjutData, *errorutils.Error) {
	if req == nil {
		req = &tindaklanjutDomain.GetPasienTindakLanjutRequest{}
	}
	page := pagination.ValidatePage(req.Page)
	perPage := pagination.ValidatePerPage(req.PerPage)

	rows, total, err := s.repo.GetPasienTindakLanjut(ctx, page, perPage)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]tindaklanjutDomain.PasienTindakLanjutItem, len(rows))
	for i, r := range rows {
		items[i] = tindaklanjutDomain.PasienTindakLanjutItem{
			IDPasien:           r.IDPasien,
			NamaPasien:         r.NamaPasien,
			StatusGizi:         r.StatusGizi,
			StatusPasien:       r.StatusPasien,
			TanggalPemeriksaan: r.TanggalPemeriksaan,
		}
	}

	if items == nil {
		items = []tindaklanjutDomain.PasienTindakLanjutItem{}
	}

	return &tindaklanjutDomain.PasienTindakLanjutData{
		Pasien: items,
		Meta:   pagination.NewMeta(int32(page), int32(perPage), int32(total)),
	}, nil
}

func (s *Service) GetDetailPasienByID(ctx context.Context, idPasien int32) (*tindaklanjutDomain.DetailPasienTindakLanjut, *errorutils.Error) {
	pasien, err := s.repo.GetPasienByID(ctx, idPasien)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if pasien == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data pasien tidak ditemukan."}
	}

	detailRow, err := s.repo.GetDetailPasienByID(ctx, idPasien)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if detailRow == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data pasien tidak ditemukan."}
	}

	riwayat, err := s.repo.GetRiwayatPemeriksaan(ctx, idPasien)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	monitoring := &tindaklanjutDomain.MonitoringTerakhir{
		StatusGizi:    detailRow.StatusGizi,
		StatusStunting: detailRow.StatusStunting,
		Catatan:       detailRow.Catatan,
	}

	itemsRiwayat := make([]tindaklanjutDomain.RiwayatPemeriksaanItem, len(riwayat))
	for i, r := range riwayat {
		itemsRiwayat[i] = tindaklanjutDomain.RiwayatPemeriksaanItem{
			Tanggal:    r.Tanggal,
			BeratBadan: r.BeratBadan,
			TinggiBadan: r.TinggiBadan,
		}
	}

	if itemsRiwayat == nil {
		itemsRiwayat = []tindaklanjutDomain.RiwayatPemeriksaanItem{}
	}

	return &tindaklanjutDomain.DetailPasienTindakLanjut{
		IDPasien:                detailRow.IDPasien,
		NamaPasien:              detailRow.NamaPasien,
		Usia:                    detailRow.Usia,
		HasilMonitoringTerakhir: monitoring,
		RiwayatPemeriksaan:      itemsRiwayat,
	}, nil
}

func (s *Service) CreateTindakLanjut(ctx context.Context, idBidan int32, req *tindaklanjutDomain.CreateTindakLanjutRequest) (*tindaklanjutDomain.CreateTindakLanjutResponse, *errorutils.Error) {
	existing, err := s.repo.GetTindakLanjutByHasilPemeriksaan(ctx, req.IDHasilPemeriksaan)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if existing != nil {
		return nil, &errorutils.Error{Status: http.StatusConflict, Message: "Tindak lanjut untuk hasil pemeriksaan ini sudah ada."}
	}

	statusPasien := model.StatusPasien_DalamPemantauan
	if req.JenisTindakan == "Rujukan" {
		statusPasien = model.StatusPasien_PerluRujukan
	}

	now := time.Now()
	var jadwalKontrol time.Time
	if req.JadwalKontrol != "" {
		parsed, err := time.Parse("2006-01-02", req.JadwalKontrol)
		if err == nil {
			jadwalKontrol = parsed
		}
	}

	tindakLanjutData := &model.TindakLanjut{
		IDHasilPemeriksaan: req.IDHasilPemeriksaan,
		IDBidan:            idBidan,
		CatatanMedis:       strPtr(req.CatatanMedis),
		Rekomendasi:        strPtr(req.Rekomendasi),
		JadwalKontrol:      jadwalKontrol,
		StatusPasien:       statusPasien,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	err = s.repo.CreateTindakLanjut(ctx, tindakLanjutData)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal membuat tindak lanjut."}
	}

	var idRujukan *int32

	if req.JenisTindakan == "Rujukan" && req.AlasanRujukan != "" && req.IDFaskes != nil {
		rujukanData := &model.Rujukan{
			IDTindakLanjut: tindakLanjutData.IDTindakLanjut,
			AlasanRujukan:  req.AlasanRujukan,
			TanggalRujukan: now,
			StatusRujukan:  model.StatusRujukan_Diajukan,
			IDFaskes:       *req.IDFaskes,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		err = s.repo.CreateRujukan(ctx, rujukanData)
		if err != nil {
			return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal membuat rujukan."}
		}
		idRujukan = &rujukanData.IDRujukan
	}

	idStr := strconv.Itoa(int(tindakLanjutData.IDTindakLanjut))
	s.logAudit(ctx, "POST /tindak-lanjut", model.TipeAktivitas_DataInsert, true, "tindak_lanjut", idStr, "Berhasil membuat tindak lanjut")

	notifPesan := "Tindak lanjut telah dibuat untuk hasil pemeriksaan Anda."
	if idRujukan != nil {
		notifPesan = "Tindak lanjut dan rujukan telah dibuat untuk hasil pemeriksaan Anda."
	}
	s.notifRepo.Create(ctx, &model.Notifikasi{
		Judul:          "Tindak Lanjut & Rujukan",
		Pesan:          &notifPesan,
		TipeNotifikasi: model.TipeNotifikasi_Rujukan,
		StatusBaca:     false,
		TanggalKirim:   time.Now(),
	})

	return &tindaklanjutDomain.CreateTindakLanjutResponse{
		IDTindakLanjut: tindakLanjutData.IDTindakLanjut,
		IDRujukan:      idRujukan,
		StatusPasien:   string(statusPasien),
	}, nil
}

func (s *Service) UpdateStatusRujukan(ctx context.Context, idRujukan int32, req *tindaklanjutDomain.UpdateStatusRujukanRequest) (*tindaklanjutDomain.UpdateStatusRujukanResponse, *errorutils.Error) {
	statusRujukan := model.StatusRujukan(req.StatusRujukan)

	valid := false
	for _, v := range model.StatusRujukanAllValues {
		if v == statusRujukan {
			valid = true
			break
		}
	}
	if !valid {
		return nil, &errorutils.Error{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("Status rujukan '%s' tidak valid.", req.StatusRujukan),
		}
	}

	result, err := s.repo.UpdateStatusRujukan(ctx, idRujukan, statusRujukan)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal memperbarui status rujukan."}
	}
	if result == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data rujukan tidak ditemukan."}
	}

	idStr := strconv.Itoa(int(idRujukan))
	s.logAudit(ctx, "PATCH /rujukan/"+idStr+"/status", model.TipeAktivitas_DataUpdate, true, "rujukan", idStr, "Berhasil memperbarui status rujukan menjadi: "+req.StatusRujukan)

	return &tindaklanjutDomain.UpdateStatusRujukanResponse{
		IDRujukan:    idRujukan,
		StatusRujukan: req.StatusRujukan,
	}, nil
}

func (s *Service) GetStatusTindakLanjut(ctx context.Context, req *tindaklanjutDomain.GetStatusTindakLanjutRequest) (*tindaklanjutDomain.StatusTindakLanjutData, *errorutils.Error) {
	if req == nil {
		req = &tindaklanjutDomain.GetStatusTindakLanjutRequest{}
	}
	page := pagination.ValidatePage(req.Page)
	perPage := pagination.ValidatePerPage(req.PerPage)

	rows, total, err := s.repo.GetStatusTindakLanjut(ctx, page, perPage)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]tindaklanjutDomain.StatusTindakLanjutItem, len(rows))
	for i, r := range rows {
		items[i] = tindaklanjutDomain.StatusTindakLanjutItem{
			IDPasien:       r.IDPasien,
			NamaPasien:     r.NamaPasien,
			StatusPasien:   r.StatusPasien,
			StatusRujukan:  r.StatusRujukan,
			TanggalRujukan: r.TanggalRujukan,
		}
	}

	if items == nil {
		items = []tindaklanjutDomain.StatusTindakLanjutItem{}
	}

	return &tindaklanjutDomain.StatusTindakLanjutData{
		Pasien: items,
		Meta:   pagination.NewMeta(int32(page), int32(perPage), int32(total)),
	}, nil
}

func (s *Service) GetLaporanTindakLanjut(ctx context.Context) (*tindaklanjutDomain.LaporanTindakLanjutData, *errorutils.Error) {
	rows, err := s.repo.GetLaporanTindakLanjut(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]tindaklanjutDomain.LaporanTindakLanjutItem, len(rows))
	for i, r := range rows {
		items[i] = tindaklanjutDomain.LaporanTindakLanjutItem{
			Wilayah:             r.Wilayah,
			JumlahPasienDirujuk: r.JumlahPasienDirujuk,
			JumlahPasienDiterima: r.JumlahPasienDiterima,
			JumlahPasienDiproses: r.JumlahPasienDiproses,
		}
	}

	if items == nil {
		items = []tindaklanjutDomain.LaporanTindakLanjutItem{}
	}

	return &tindaklanjutDomain.LaporanTindakLanjutData{
		Laporan:   items,
		TotalData: len(items),
	}, nil
}

func (s *Service) GetDetailTindakLanjutByID(ctx context.Context, idTindakLanjut int32) (*tindaklanjutDomain.DetailTindakLanjutPasien, *errorutils.Error) {
	row, err := s.repo.GetDetailTindakLanjutByID(ctx, idTindakLanjut)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if row == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data tindak lanjut pasien tidak ditemukan."}
	}

	return &tindaklanjutDomain.DetailTindakLanjutPasien{
		IDTindakLanjut: row.IDTindakLanjut,
		StatusPasien:   row.StatusPasien,
		CatatanMedis:   row.CatatanMedis,
		Rekomendasi:    row.Rekomendasi,
		JadwalKontrol:  row.JadwalKontrol,
		StatusRujukan:  row.StatusRujukan,
		NamaFaskes:     row.NamaFaskes,
	}, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
