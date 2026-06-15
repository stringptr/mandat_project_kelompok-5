package pemeriksaan

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	auditlogDomain "github.com/stringptr/SiGizi/backend/internal/domain/auditlog"
	notificationDomain "github.com/stringptr/SiGizi/backend/internal/domain/notification"
	pemeriksaanDomain "github.com/stringptr/SiGizi/backend/internal/domain/pemeriksaan"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Service struct {
	repo         pemeriksaanDomain.Repo
	auditRepo    auditlogDomain.Repo
	notifRepo    notificationDomain.Repo
	notifPublisher notificationDomain.Publisher
}

func NewService(repo pemeriksaanDomain.Repo, auditRepo auditlogDomain.Repo, notifRepo notificationDomain.Repo, notifPublisher notificationDomain.Publisher) *Service {
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

func (s *Service) Create(ctx context.Context, idPetugas int32, req *pemeriksaanDomain.CreatePemeriksaanRequest) (*pemeriksaanDomain.CreatePemeriksaanResponse, *errorutils.Error) {
	jadwal, err := s.repo.GetJadwalImunisasiByID(ctx, req.IDJadwalImunisasi)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if jadwal == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data jadwal imunisasi referensi tidak ditemukan."}
	}

	now := time.Now()
	statusGizi := calculateStatusGizi(req.BeratBadan)
	statusStunting := calculateStatusStunting(req.TinggiBadan)

	modelData := &model.HasilPemeriksaan{
		IDPetugasInput:    idPetugas,
		IDJadwalImunisasi: req.IDJadwalImunisasi,
		BeratBadan:        decimal.NewFromFloat(req.BeratBadan),
		TinggiBadan:       decimal.NewFromFloat(req.TinggiBadan),
		LingkarKepala:     decimal.NewFromFloat(req.LingkarKepala),
		TekananDarah:      req.TekananDarah,
		StatusStunting:    model.StatusStunting(statusStunting),
		StatusGizi:        model.StatusGizi(statusGizi),
		Catatan:           strPtr(req.Catatan),
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	err = s.repo.Create(ctx, modelData)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal membuat data pemeriksaan."}
	}

	idStr := strconv.Itoa(int(modelData.IDHasilPemeriksaan))
	s.logAudit(ctx, "POST /monitoring/pemeriksaan", model.TipeAktivitas_DataInsert, true, "hasil_pemeriksaan", idStr, "Berhasil membuat data pemeriksaan")

	notifErr := s.createNotificationForCreate(ctx, jadwal, idPetugas)
	if notifErr != nil {
		return nil, notifErr
	}

	return &pemeriksaanDomain.CreatePemeriksaanResponse{
		IDHasilPemeriksaan: modelData.IDHasilPemeriksaan,
		StatusStunting:     statusStunting,
		StatusGizi:         statusGizi,
		CreatedAt:          now.Format(time.RFC3339),
	}, nil
}

func (s *Service) GetByID(ctx context.Context, idHasilPemeriksaan int32) (*pemeriksaanDomain.DetailPemeriksaanResponse, *errorutils.Error) {
	row, err := s.repo.GetDetailJoinByID(ctx, idHasilPemeriksaan)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if row == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data pemeriksaan tidak ditemukan."}
	}

	catatan := ""
	if row.Catatan != nil {
		catatan = *row.Catatan
	}

	return &pemeriksaanDomain.DetailPemeriksaanResponse{
		IDHasilPemeriksaan: row.IDHasilPemeriksaan,
		Pasien: pemeriksaanDomain.PasienInfo{
			IDPasien: row.IDPasien,
			Nama:     row.NamaPasien,
		},
		Antropometri: pemeriksaanDomain.AntropometriData{
			BeratBadan:   row.BeratBadan,
			TinggiBadan:  row.TinggiBadan,
			LingkarKepala: row.LingkarKepala,
			TekananDarah:  row.TekananDarah,
		},
		StatusKesehatan: pemeriksaanDomain.StatusKesehatanData{
			StatusStunting: row.StatusStunting,
			StatusGizi:     row.StatusGizi,
		},
		Catatan: catatan,
	}, nil
}

func (s *Service) Update(ctx context.Context, req *pemeriksaanDomain.UpdatePemeriksaanRequest, idHasilPemeriksaan int32) (*pemeriksaanDomain.UpdatePemeriksaanResponse, *errorutils.Error) {
	existing, err := s.repo.GetByID(ctx, idHasilPemeriksaan)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if existing == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data pemeriksaan dengan ID tersebut tidak ditemukan."}
	}

	updatedStatusGizi := string(existing.StatusGizi)

	if req.BeratBadan != nil {
		existing.BeratBadan = decimal.NewFromFloat(*req.BeratBadan)
	}
	if req.TinggiBadan != nil {
		existing.TinggiBadan = decimal.NewFromFloat(*req.TinggiBadan)
	}
	if req.LingkarKepala != nil {
		existing.LingkarKepala = decimal.NewFromFloat(*req.LingkarKepala)
	}
	if req.TekananDarah != nil {
		existing.TekananDarah = *req.TekananDarah
	}
	if req.Catatan != nil {
		existing.Catatan = req.Catatan
	}

	if req.BeratBadan != nil || req.TinggiBadan != nil {
		bb, _ := existing.BeratBadan.Float64()
		tb, _ := existing.TinggiBadan.Float64()
		if req.BeratBadan != nil {
			bb = *req.BeratBadan
		}
		if req.TinggiBadan != nil {
			tb = *req.TinggiBadan
		}
		if req.BeratBadan != nil {
			updatedStatusGizi = calculateStatusGizi(bb)
			existing.StatusGizi = model.StatusGizi(updatedStatusGizi)
		}
		if req.TinggiBadan != nil {
			existing.StatusStunting = model.StatusStunting(calculateStatusStunting(tb))
		}
	}

	existing.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, existing)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal memperbarui data pemeriksaan."}
	}

	idStr := strconv.Itoa(int(idHasilPemeriksaan))
	s.logAudit(ctx, "PUT /monitoring/pemeriksaan/"+idStr, model.TipeAktivitas_DataUpdate, true, "hasil_pemeriksaan", idStr, "Berhasil memperbarui data pemeriksaan")

	return &pemeriksaanDomain.UpdatePemeriksaanResponse{
		IDHasilPemeriksaan: idHasilPemeriksaan,
		StatusGiziBaru:     updatedStatusGizi,
		UpdatedAt:          existing.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) Delete(ctx context.Context, idHasilPemeriksaan int32) *errorutils.Error {
	existing, err := s.repo.GetByID(ctx, idHasilPemeriksaan)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if existing == nil {
		return &errorutils.Error{Status: http.StatusNotFound, Message: "Data pemeriksaan dengan ID tersebut tidak ditemukan."}
	}

	err = s.repo.Delete(ctx, idHasilPemeriksaan)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal menghapus data pemeriksaan."}
	}

	idStr := strconv.Itoa(int(idHasilPemeriksaan))
	s.logAudit(ctx, "DELETE /monitoring/pemeriksaan/"+idStr, model.TipeAktivitas_DataDelete, true, "hasil_pemeriksaan", idStr, "Berhasil menghapus data pemeriksaan")

	return nil
}

func (s *Service) Verify(ctx context.Context, idHasilPemeriksaan int32, idBidan int32) (*pemeriksaanDomain.VerifyPemeriksaanResponse, *errorutils.Error) {
	existing, err := s.repo.GetByID(ctx, idHasilPemeriksaan)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if existing == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data pemeriksaan yang akan diverifikasi tidak ditemukan."}
	}

	existing.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, existing)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal memverifikasi data pemeriksaan."}
	}

	idStr := strconv.Itoa(int(idHasilPemeriksaan))
	s.logAudit(ctx, "PATCH /monitoring/pemeriksaan/"+idStr+"/verify", model.TipeAktivitas_VerifikasiRegistrasi, true, "hasil_pemeriksaan", idStr, "Berhasil memverifikasi data pemeriksaan oleh bidan")

	bidanUser, err := s.repo.GetUserByID(ctx, idBidan)
	bidanNama := "Bidan"
	if err == nil && bidanUser != nil {
		bidanNama = bidanUser.Nama
	}

	_, pubErr := s.notifRepo.Create(ctx, &model.Notifikasi{
		IDUser:          existing.IDPetugasInput,
		Judul:           "Pemeriksaan Terverifikasi",
		Pesan:           strPtr(fmt.Sprintf("Data pemeriksaan telah diverifikasi oleh %s.", bidanNama)),
		TipeNotifikasi:  model.TipeNotifikasi_Pemeriksaan,
		StatusBaca:      false,
		TanggalKirim:    time.Now(),
	})
	if pubErr == nil {
		s.notifPublisher.PublishToUser(existing.IDPetugasInput, &notificationDomain.Notification{
			Judul: "Pemeriksaan Terverifikasi",
			Pesan: fmt.Sprintf("Data pemeriksaan telah diverifikasi oleh %s.", bidanNama),
			Tipe:  string(model.TipeNotifikasi_Pemeriksaan),
		})
	}

	return &pemeriksaanDomain.VerifyPemeriksaanResponse{
		IDHasilPemeriksaan: idHasilPemeriksaan,
		DiverifikasiOleh:   idBidan,
		StatusVerifikasi:   "Aktif",
	}, nil
}

func (s *Service) GetPending(ctx context.Context) (*pemeriksaanDomain.PendingPemeriksaanData, *errorutils.Error) {
	rows, err := s.repo.GetPendingVerification(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]pemeriksaanDomain.PendingPemeriksaanItem, len(rows))
	for i, r := range rows {
		items[i] = pemeriksaanDomain.PendingPemeriksaanItem{
			IDHasilPemeriksaan: r.IDHasilPemeriksaan,
			NamaPasien:         r.NamaPasien,
			DiinputOleh:        r.DiinputOleh,
			TanggalInput:       r.TanggalInput.Format(time.RFC3339),
		}
	}

	if items == nil {
		items = []pemeriksaanDomain.PendingPemeriksaanItem{}
	}

	return &pemeriksaanDomain.PendingPemeriksaanData{
		PemeriksaanPending: items,
		TotalPending:       len(items),
	}, nil
}

func (s *Service) createNotificationForCreate(ctx context.Context, jadwal *model.JadwalImunisasi, idPetugas int32) *errorutils.Error {
	petugas, err := s.repo.GetUserByID(ctx, idPetugas)
	if err != nil || petugas == nil {
		return nil
	}

	notif := &model.Notifikasi{
		IDUser:         idPetugas,
		Judul:          "Pemeriksaan Baru",
		Pesan:          strPtr(fmt.Sprintf("Data pemeriksaan baru berhasil dibuat oleh %s.", petugas.Nama)),
		TipeNotifikasi: model.TipeNotifikasi_Pemeriksaan,
		StatusBaca:     false,
		TanggalKirim:   time.Now(),
	}
	_, errCreate := s.notifRepo.Create(ctx, notif)
	if errCreate != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal membuat notifikasi."}
	}

	s.notifPublisher.PublishToUser(idPetugas, &notificationDomain.Notification{
		Judul: "Pemeriksaan Baru",
		Pesan: fmt.Sprintf("Data pemeriksaan baru berhasil dibuat oleh %s.", petugas.Nama),
		Tipe:  string(model.TipeNotifikasi_Pemeriksaan),
	})

	return nil
}

func calculateStatusGizi(beratBadan float64) string {
	switch {
	case beratBadan < 3.0:
		return "Gizi Buruk"
	case beratBadan < 5.0:
		return "Gizi Kurang"
	case beratBadan < 15.0:
		return "Gizi Baik"
	case beratBadan < 25.0:
		return "Risiko Gizi Lebih"
	case beratBadan < 35.0:
		return "Gizi Lebih"
	default:
		return "Obesitas"
	}
}

func calculateStatusStunting(tinggiBadan float64) string {
	switch {
	case tinggiBadan < 50:
		return "Stunting Berat"
	case tinggiBadan < 70:
		return "Stunting"
	case tinggiBadan < 90:
		return "Berisiko Stunting"
	default:
		return "Normal"
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
