package imunisasi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	auditlogDomain "github.com/stringptr/SiGizi/backend/internal/domain/auditlog"
	imunisasiDomain "github.com/stringptr/SiGizi/backend/internal/domain/imunisasi"
	notificationDomain "github.com/stringptr/SiGizi/backend/internal/domain/notification"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/stringptr/SiGizi/backend/internal/pagination"
)

func isPetugas(roles []string) bool {
	for _, r := range roles {
		switch r {
		case "ADMIN", "BIDAN", "KADER", "DINKES", "SUPER_ADMIN":
			return true
		}
	}
	return false
}

type Service struct {
	repo          imunisasiDomain.Repo
	auditRepo     auditlogDomain.Repo
	notifRepo     notificationDomain.Repo
	notifPublisher notificationDomain.Publisher
}

func NewService(repo imunisasiDomain.Repo, auditRepo auditlogDomain.Repo, notifRepo notificationDomain.Repo, notifPublisher notificationDomain.Publisher) *Service {
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

func (s *Service) GetAllByUser(ctx context.Context, idUser int32, req *imunisasiDomain.GetAllImunisasiRequest) (*imunisasiDomain.ImunisasiListData, *errorutils.Error) {
	if req == nil {
		req = &imunisasiDomain.GetAllImunisasiRequest{}
	}
	page := pagination.ValidatePage(req.Page)
	perPage := pagination.ValidatePerPage(req.PerPage)

	rows, total, err := s.repo.GetAllByUser(ctx, idUser, page, perPage, req.Q)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]imunisasiDomain.ImunisasiListItem, len(rows))
	for i, r := range rows {
		items[i] = imunisasiDomain.ImunisasiListItem{
			IDImunisasi:     r.IDImunisasi,
			NamaPasien:      r.NamaPasien,
			NamaVaksin:      r.NamaVaksin,
			TanggalJadwal:   r.TanggalJadwal,
			StatusImunisasi: r.StatusImunisasi,
		}
	}

	if items == nil {
		items = []imunisasiDomain.ImunisasiListItem{}
	}

	return &imunisasiDomain.ImunisasiListData{
		Jadwal: items,
		Meta:   pagination.NewMeta(int32(page), int32(perPage), int32(total)),
	}, nil
}

func (s *Service) GetAll(ctx context.Context, req *imunisasiDomain.GetAllImunisasiRequest) (*imunisasiDomain.ImunisasiListData, *errorutils.Error) {
	if req == nil {
		req = &imunisasiDomain.GetAllImunisasiRequest{}
	}
	page := pagination.ValidatePage(req.Page)
	perPage := pagination.ValidatePerPage(req.PerPage)

	rows, total, err := s.repo.GetAll(ctx, page, perPage, req.Q)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]imunisasiDomain.ImunisasiListItem, len(rows))
	for i, r := range rows {
		items[i] = imunisasiDomain.ImunisasiListItem{
			IDImunisasi:     r.IDImunisasi,
			NamaPasien:      r.NamaPasien,
			NamaVaksin:      r.NamaVaksin,
			TanggalJadwal:   r.TanggalJadwal,
			StatusImunisasi: r.StatusImunisasi,
		}
	}

	if items == nil {
		items = []imunisasiDomain.ImunisasiListItem{}
	}

	return &imunisasiDomain.ImunisasiListData{
		Jadwal: items,
		Meta:   pagination.NewMeta(int32(page), int32(perPage), int32(total)),
	}, nil
}

func (s *Service) GetByID(ctx context.Context, idImunisasi int32, claims *jwtutils.Claim) (*imunisasiDomain.ImunisasiDetail, *errorutils.Error) {
	row, err := s.repo.GetDetailJoinByID(ctx, idImunisasi)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if row == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data jadwal imunisasi dengan ID tersebut tidak ditemukan."}
	}

	if claims != nil && !isPetugas(claims.Roles) {
		ok, err := s.isOwnPasien(ctx, row.IDPasien, claims.IDUser)
		if err != nil {
			return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
		}
		if !ok {
			return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data jadwal imunisasi dengan ID tersebut tidak ditemukan."}
		}
	}

	return &imunisasiDomain.ImunisasiDetail{
		IDImunisasi:      row.IDImunisasi,
		IDPasien:         row.IDPasien,
		NamaPasien:       row.NamaPasien,
		NamaVaksin:       row.NamaVaksin,
		TanggalJadwal:    row.TanggalJadwal,
		TanggalRealisasi: row.TanggalRealisasi,
		StatusImunisasi:  row.StatusImunisasi,
	}, nil
}

func (s *Service) Create(ctx context.Context, req *imunisasiDomain.CreateImunisasiRequest) (*imunisasiDomain.CreateImunisasiResponse, *errorutils.Error) {
	pasien, err := s.repo.GetPasienByID(ctx, req.IDPasien)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if pasien == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data pasien yang dirujuk tidak ditemukan."}
	}

	tanggalJadwal, err := time.Parse("2006-01-02", req.TanggalJadwal)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusBadRequest, Message: "Format tanggal_jadwal tidak valid. Gunakan YYYY-MM-DD."}
	}

	now := time.Now()
	modelData := &model.JadwalImunisasi{
		IDPasien:        req.IDPasien,
		NamaVaksin:      req.NamaVaksin,
		TanggalJadwal:   tanggalJadwal,
		StatusImunisasi: model.StatusImunisasi_Belum,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = s.repo.Create(ctx, modelData)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal membuat jadwal imunisasi."}
	}

	idStr := strconv.Itoa(int(modelData.IDImunisasi))
	s.logAudit(ctx, "POST /imunisasi", model.TipeAktivitas_DataInsert, true, "jadwal_imunisasi", idStr, "Berhasil membuat jadwal imunisasi")

	notifPesan := fmt.Sprintf("Jadwal imunisasi %s telah dibuat pada %s.", req.NamaVaksin, req.TanggalJadwal)
	s.notifRepo.Create(ctx, &model.Notifikasi{
		IDUser:         req.IDPasien,
		Judul:          "Jadwal Imunisasi Baru",
		Pesan:          &notifPesan,
		TipeNotifikasi: model.TipeNotifikasi_Imunisasi,
		StatusBaca:     false,
		TanggalKirim:   time.Now(),
	})

	s.notifPublisher.PublishToUser(req.IDPasien, &notificationDomain.Notification{
		Judul: "Jadwal Imunisasi Baru",
		Pesan: notifPesan,
		Tipe:  string(model.TipeNotifikasi_Imunisasi),
	})

	return &imunisasiDomain.CreateImunisasiResponse{
		IDImunisasi:     modelData.IDImunisasi,
		StatusImunisasi: string(model.StatusImunisasi_Belum),
	}, nil
}

func (s *Service) Update(ctx context.Context, req *imunisasiDomain.UpdateImunisasiRequest, idImunisasi int32) (*imunisasiDomain.UpdateImunisasiResponse, *errorutils.Error) {
	existing, err := s.repo.GetByID(ctx, idImunisasi)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if existing == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data jadwal imunisasi dengan ID tersebut tidak ditemukan."}
	}

	if req.IDPasien != nil {
		pasien, err := s.repo.GetPasienByID(ctx, *req.IDPasien)
		if err != nil {
			return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
		}
		if pasien == nil {
			return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data pasien yang dirujuk tidak ditemukan."}
		}
		existing.IDPasien = *req.IDPasien
	}
	if req.NamaVaksin != nil {
		existing.NamaVaksin = *req.NamaVaksin
	}
	if req.TanggalJadwal != nil {
		t, err := time.Parse("2006-01-02", *req.TanggalJadwal)
		if err != nil {
			return nil, &errorutils.Error{Status: http.StatusBadRequest, Message: "Format tanggal_jadwal tidak valid."}
		}
		existing.TanggalJadwal = t
	}

	existing.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, existing)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal memperbarui jadwal imunisasi."}
	}

	idStr := strconv.Itoa(int(idImunisasi))
	s.logAudit(ctx, "PUT /imunisasi/"+idStr, model.TipeAktivitas_DataUpdate, true, "jadwal_imunisasi", idStr, "Berhasil memperbarui jadwal imunisasi")

	return &imunisasiDomain.UpdateImunisasiResponse{
		IDImunisasi: idImunisasi,
		UpdatedAt:   existing.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) Delete(ctx context.Context, idImunisasi int32) *errorutils.Error {
	existing, err := s.repo.GetByID(ctx, idImunisasi)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if existing == nil {
		return &errorutils.Error{Status: http.StatusNotFound, Message: "Data jadwal imunisasi dengan ID tersebut tidak ditemukan."}
	}

	err = s.repo.Delete(ctx, idImunisasi)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal menghapus jadwal imunisasi."}
	}

	idStr := strconv.Itoa(int(idImunisasi))
	s.logAudit(ctx, "DELETE /imunisasi/"+idStr, model.TipeAktivitas_DataDelete, true, "jadwal_imunisasi", idStr, "Berhasil menghapus jadwal imunisasi")

	return nil
}

func (s *Service) Realisasi(ctx context.Context, idImunisasi int32, req *imunisasiDomain.RealisasiRequest) (*imunisasiDomain.RealisasiResponse, *errorutils.Error) {
	existing, err := s.repo.GetByID(ctx, idImunisasi)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if existing == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data jadwal imunisasi dengan ID tersebut tidak ditemukan."}
	}

	err = s.repo.UpdateRealisasi(ctx, idImunisasi, req.TanggalRealisasi)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal mencatat realisasi imunisasi."}
	}

	idStr := strconv.Itoa(int(idImunisasi))
	s.logAudit(ctx, "PATCH /imunisasi/"+idStr+"/realisasi", model.TipeAktivitas_DataUpdate, true, "jadwal_imunisasi", idStr, "Berhasil mencatat realisasi imunisasi")

	notifPesan := fmt.Sprintf("Imunisasi %s telah dilaksanakan pada %s.", existing.NamaVaksin, req.TanggalRealisasi)
	s.notifRepo.Create(ctx, &model.Notifikasi{
		IDUser:         existing.IDPasien,
		Judul:          "Imunisasi Terlaksana",
		Pesan:          &notifPesan,
		TipeNotifikasi: model.TipeNotifikasi_Imunisasi,
		StatusBaca:     false,
		TanggalKirim:   time.Now(),
	})

	s.notifPublisher.PublishToUser(existing.IDPasien, &notificationDomain.Notification{
		Judul: "Imunisasi Terlaksana",
		Pesan: notifPesan,
		Tipe:  string(model.TipeNotifikasi_Imunisasi),
	})

	return &imunisasiDomain.RealisasiResponse{
		IDImunisasi:      idImunisasi,
		StatusImunisasi:  string(model.StatusImunisasi_Sudah),
		TanggalRealisasi: req.TanggalRealisasi,
	}, nil
}

func (s *Service) IsOwnPasien(ctx context.Context, idPasien int32, idUser int32) (bool, *errorutils.Error) {
	owned, err := s.repo.CheckPasienOwnership(ctx, idPasien, idUser)
	if err != nil {
		return false, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	return owned, nil
}

func (s *Service) GetByPasienID(ctx context.Context, idPasien int32, claims *jwtutils.Claim) (*imunisasiDomain.RiwayatImunisasiResponse, *errorutils.Error) {
	pasien, err := s.repo.GetPasienByID(ctx, idPasien)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if pasien == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data pasien dengan ID tersebut tidak ditemukan."}
	}

	if claims != nil && !isPetugas(claims.Roles) {
		ok, err := s.isOwnPasien(ctx, idPasien, claims.IDUser)
		if err != nil {
			return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
		}
		if !ok {
			return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Data pasien dengan ID tersebut tidak ditemukan."}
		}
	}

	jadwal, err := s.repo.GetByPasienID(ctx, idPasien)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	items := make([]imunisasiDomain.RiwayatImunisasiItem, len(jadwal))
	for i, j := range jadwal {
		var tanggalRealisasi *string
		if j.TanggalRealisasi != nil {
			t := j.TanggalRealisasi.Format("2006-01-02")
			tanggalRealisasi = &t
		}

		items[i] = imunisasiDomain.RiwayatImunisasiItem{
			IDImunisasi:      j.IDImunisasi,
			NamaVaksin:       j.NamaVaksin,
			TanggalJadwal:    j.TanggalJadwal.Format("2006-01-02"),
			TanggalRealisasi: tanggalRealisasi,
			StatusImunisasi:  string(j.StatusImunisasi),
		}
	}

	if items == nil {
		items = []imunisasiDomain.RiwayatImunisasiItem{}
	}

	return &imunisasiDomain.RiwayatImunisasiResponse{
		IDPasien:         idPasien,
		RiwayatImunisasi: items,
	}, nil
}

func (s *Service) GetStatistik(ctx context.Context) (*imunisasiDomain.StatistikImunisasi, *errorutils.Error) {
	row, err := s.repo.GetStatistik(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}
	if row == nil {
		return &imunisasiDomain.StatistikImunisasi{}, nil
	}

	vaksinTerbanyak, err := s.repo.GetVaksinTerbanyak(ctx)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Mohon dicoba kembali."}
	}

	var persentase float64
	if row.TotalTarget > 0 {
		persentase = (float64(row.TotalRealisasi) / float64(row.TotalTarget)) * 100
	}

	return &imunisasiDomain.StatistikImunisasi{
		TotalTargetImunisasi: row.TotalTarget,
		TotalTerealisasi:     row.TotalRealisasi,
		CakupanPersentase:    persentase,
		VaksinTerbanyak:      vaksinTerbanyak,
	}, nil
}

func (s *Service) isOwnPasien(ctx context.Context, idPasien, idUser int32) (bool, error) {
	if idPasien == idUser {
		return true, nil
	}
	anak, err := s.repo.GetAnakByPasienID(ctx, idPasien)
	if err != nil {
		return false, err
	}
	if anak != nil && anak.IDWali == idUser {
		return true, nil
	}
	return false, nil
}

