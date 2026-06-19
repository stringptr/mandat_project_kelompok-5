package pasien

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	pasienDomain "github.com/stringptr/SiGizi/backend/internal/domain/pasien"
	auditlogDomain "github.com/stringptr/SiGizi/backend/internal/domain/auditlog"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
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
	repo      pasienDomain.Repo
	auditRepo auditlogDomain.Repo
}

func NewService(repo pasienDomain.Repo, auditRepo auditlogDomain.Repo) *Service {
	return &Service{repo: repo, auditRepo: auditRepo}
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

func (s *Service) DaftarIbuHamil(ctx context.Context, req *pasienDomain.DaftarIbuHamilRequest) *errorutils.Error {
	user, err := s.repo.GetUserByID(ctx, req.IDUser)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if user == nil {
		return &errorutils.Error{Status: http.StatusNotFound, Message: "User tidak ditemukan."}
	}

	existingPasien, err := s.repo.GetPasienByUserID(ctx, req.IDUser)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if existingPasien != nil {
		return &errorutils.Error{Status: http.StatusConflict, Message: "User sudah terdaftar sebagai pasien."}
	}

	posyandu, err := s.repo.GetPosyanduByID(ctx, req.IDPosyandu)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if posyandu == nil {
		return &errorutils.Error{Status: http.StatusNotFound, Message: "Posyandu tidak ditemukan."}
	}

	now := time.Now()

	pasienModel := &model.Pasien{
		IDPasien:   req.IDUser,
		IDPosyandu: req.IDPosyandu,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	var bulanMulaiHamil time.Time
	if req.BulanMulaiHamil != "" {
		t, err := time.Parse("2006-01-02", req.BulanMulaiHamil)
		if err != nil {
			return &errorutils.Error{Status: http.StatusBadRequest, Message: "Format bulan_mulai_hamil tidak valid. Gunakan YYYY-MM-DD."}
		}
		bulanMulaiHamil = t
	}

	var hpht time.Time
	if req.Hpht != "" {
		t, err := time.Parse("2006-01-02", req.Hpht)
		if err != nil {
			return &errorutils.Error{Status: http.StatusBadRequest, Message: "Format hpht tidak valid. Gunakan YYYY-MM-DD."}
		}
		hpht = t
	}

	ibuHamilModel := &model.IbuHamil{
		IDPasien:        req.IDUser,
		HamilKe:         req.HamilKe,
		BulanMulaiHamil: bulanMulaiHamil,
		Hpht:            hpht,
		StatusKehamilan: model.StatusKehamilan(req.StatusKehamilan),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = s.repo.CreatePasien(ctx, pasienModel)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal membuat data pasien."}
	}

	err = s.repo.CreateIbuHamil(ctx, ibuHamilModel)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal membuat data ibu hamil."}
	}

	s.logAudit(ctx, "POST /pasien/daftar-ibu-hamil", model.TipeAktivitas_DataInsert, true, "pasien, ibu_hamil", strconv.Itoa(int(req.IDUser)), "Berhasil mendaftarkan ibu hamil")
	return nil
}

func (s *Service) DaftarAnak(ctx context.Context, req *pasienDomain.DaftarAnakRequest) *errorutils.Error {
	user, err := s.repo.GetUserByID(ctx, req.IDUser)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if user == nil {
		return &errorutils.Error{Status: http.StatusNotFound, Message: "User tidak ditemukan."}
	}

	existingPasien, err := s.repo.GetPasienByUserID(ctx, req.IDUser)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if existingPasien != nil {
		return &errorutils.Error{Status: http.StatusConflict, Message: "User sudah terdaftar sebagai pasien."}
	}

	posyandu, err := s.repo.GetPosyanduByID(ctx, req.IDPosyandu)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if posyandu == nil {
		return &errorutils.Error{Status: http.StatusNotFound, Message: "Posyandu tidak ditemukan."}
	}

	now := time.Now()

	pasienModel := &model.Pasien{
		IDPasien:   req.IDUser,
		IDPosyandu: req.IDPosyandu,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	anakModel := &model.Anak{
		IDPasien:           req.IDUser,
		IDIbuHamil:         req.IDIbuHamil,
		IDWali:             req.IDWali,
		NamaAnak:           req.NamaAnak,
		BeratLahir:         decimal.NewFromFloat(req.BeratLahir),
		PanjangLahir:       decimal.NewFromFloat(req.PanjangLahir),
		HubunganDenganWali: model.HubunganDenganWali(req.HubunganDenganWali),
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	err = s.repo.CreatePasien(ctx, pasienModel)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal membuat data pasien."}
	}

	err = s.repo.CreateAnak(ctx, anakModel)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal membuat data anak."}
	}

	s.logAudit(ctx, "POST /pasien/daftar-anak", model.TipeAktivitas_DataInsert, true, "pasien, anak", strconv.Itoa(int(req.IDUser)), "Berhasil mendaftarkan anak")
	return nil
}

func (s *Service) GetAllByUser(ctx context.Context, idUser int32, req *pasienDomain.GetAllPasienRequest) (*pasienDomain.PasienListData, *errorutils.Error) {
	if req == nil {
		req = &pasienDomain.GetAllPasienRequest{}
	}
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	rows, total, err := s.repo.GetAllPaginatedByUser(ctx, page, perPage, req.Q, idUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	items := make([]pasienDomain.PasienListItem, len(rows))
	for i, r := range rows {
		items[i] = pasienDomain.PasienListItem{
			IDPasien:        r.IDPasien,
			Nama:            r.Nama,
			NIK:             r.NIK,
			JenisKelamin:    r.JenisKelamin,
			Umur:            calculateAge(r.TanggalLahir),
			NamaPosyandu:    r.NamaPosyandu,
			JenisPasien:     r.JenisPasien,
			StatusKehamilan: r.StatusKehamilan,
		}
	}

	return &pasienDomain.PasienListData{
		Pasien:    items,
		TotalData: total,
		Page:      page,
		PerPage:   perPage,
	}, nil
}

func (s *Service) IsOwnPasien(ctx context.Context, idPasien int32, idUser int32) (bool, *errorutils.Error) {
	owned, err := s.repo.CheckPasienOwnership(ctx, idPasien, idUser)
	if err != nil {
		return false, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	return owned, nil
}

func (s *Service) GetAll(ctx context.Context, req *pasienDomain.GetAllPasienRequest) (*pasienDomain.PasienListData, *errorutils.Error) {
	if req == nil {
		req = &pasienDomain.GetAllPasienRequest{}
	}
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	rows, total, err := s.repo.GetAllPaginated(ctx, page, perPage, req.Q)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	items := make([]pasienDomain.PasienListItem, len(rows))
	for i, r := range rows {
		items[i] = pasienDomain.PasienListItem{
			IDPasien:        r.IDPasien,
			Nama:            r.Nama,
			NIK:             r.NIK,
			JenisKelamin:    r.JenisKelamin,
			Umur:            calculateAge(r.TanggalLahir),
			NamaPosyandu:    r.NamaPosyandu,
			JenisPasien:     r.JenisPasien,
			StatusKehamilan: r.StatusKehamilan,
		}
	}

	return &pasienDomain.PasienListData{
		Pasien:    items,
		TotalData: total,
		Page:      page,
		PerPage:   perPage,
	}, nil
}

func (s *Service) Search(ctx context.Context, req *pasienDomain.SearchPasienRequest) (*pasienDomain.SearchPasienResponseData, *errorutils.Error) {
	if req == nil {
		req = &pasienDomain.SearchPasienRequest{}
	}
	if req.Q == "" {
		return nil, &errorutils.Error{Status: http.StatusBadRequest, Message: "Parameter pencarian (q) wajib diisi."}
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	rows, total, err := s.repo.Search(ctx, req.Q, page, perPage)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	items := make([]pasienDomain.PasienListItem, len(rows))
	for i, r := range rows {
		items[i] = pasienDomain.PasienListItem{
			IDPasien:        r.IDPasien,
			Nama:            r.Nama,
			NIK:             r.NIK,
			JenisKelamin:    r.JenisKelamin,
			Umur:            calculateAge(r.TanggalLahir),
			NamaPosyandu:    r.NamaPosyandu,
			JenisPasien:     r.JenisPasien,
			StatusKehamilan: r.StatusKehamilan,
		}
	}

	if items == nil {
		items = []pasienDomain.PasienListItem{}
	}

	return &pasienDomain.SearchPasienResponseData{
		Pasien:    items,
		TotalData: total,
		Page:      page,
		PerPage:   perPage,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, idPasien int32, claims *jwtutils.Claim) (*pasienDomain.PasienDetailResponse, *errorutils.Error) {
	row, err := s.repo.GetDetailByID(ctx, idPasien)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if row == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Pasien tidak ditemukan."}
	}

	if claims != nil && !isPetugas(claims.Roles) {
		if !s.isOwnPasien(ctx, idPasien, claims.IDUser) {
			return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Pasien tidak ditemukan."}
		}
	}

	tglLahir, _ := time.Parse("2006-01-02T15:04:05Z", row.TanggalLahir)
	if tglLahir.IsZero() {
		tglLahir, _ = time.Parse("2006-01-02", row.TanggalLahir)
	}

	createdAt, _ := time.Parse("2006-01-02T15:04:05Z", row.CreatedAt)
	updatedAt, _ := time.Parse("2006-01-02T15:04:05Z", row.UpdatedAt)

	resp := &pasienDomain.PasienDetailResponse{
		IDPasien:     row.IDPasien,
		Nama:         row.Nama,
		NIK:          row.NIK,
		Email:        row.Email,
		NoHp:         row.NoHp,
		JenisKelamin: row.JenisKelamin,
		TanggalLahir: tglLahir,
		IDLokasi:     row.IDLokasi,
		NamaPosyandu: row.NamaPosyandu,
		IDPosyandu:   row.IDPosyandu,
		JenisPasien:  row.JenisPasien,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	if row.IDIbuHamil != nil {
		var hamilKe int32
		if row.HamilKe != nil {
			hamilKe = *row.HamilKe
		}
		resp.DataIbuHamil = &pasienDomain.IbuHamilData{
			IDIbuHamil:      *row.IDIbuHamil,
			HamilKe:         hamilKe,
			BulanMulaiHamil: nullString(row.BulanMulaiHamil),
			Hpht:            nullString(row.Hpht),
			StatusKehamilan: nullString(row.StatusKehamilan),
		}
	}

	if row.NamaAnak != nil {
		resp.DataAnak = &pasienDomain.AnakData{
			NamaAnak:           *row.NamaAnak,
			BeratLahir:         nullFloat(row.BeratLahir),
			PanjangLahir:       nullFloat(row.PanjangLahir),
			HubunganDenganWali: nullString(row.HubunganDenganWali),
			NamaWali:           nullString(row.NamaWali),
		}
	}

	return resp, nil
}

func (s *Service) Update(ctx context.Context, req *pasienDomain.UpdatePasienRequest) (*pasienDomain.PasienDetailResponse, *errorutils.Error) {
	existing, err := s.repo.GetDetailByID(ctx, req.IDPasien)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if existing == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Pasien tidak ditemukan."}
	}

	if req.IDPosyandu != nil {
		posyandu, err := s.repo.GetPosyanduByID(ctx, *req.IDPosyandu)
		if err != nil {
			return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
		}
		if posyandu == nil {
			return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Posyandu tidak ditemukan."}
		}
		err = s.repo.UpdatePasien(ctx, &model.Pasien{
			IDPasien:   req.IDPasien,
			IDPosyandu: *req.IDPosyandu,
		})
		if err != nil {
			return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal memperbarui data pasien."}
		}
	}

	if existing.JenisPasien == "Ibu Hamil" && (req.HamilKe != nil || req.BulanMulaiHamil != nil || req.Hpht != nil || req.StatusKehamilan != nil) {
		ibuHamilList, err := s.repo.GetIbuHamilByPasienID(ctx, req.IDPasien)
		if err != nil {
			return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
		}
		if len(ibuHamilList) > 0 {
			ih := ibuHamilList[0]
			if req.HamilKe != nil {
				ih.HamilKe = *req.HamilKe
			}
			if req.BulanMulaiHamil != nil {
				t, parseErr := time.Parse("2006-01-02", *req.BulanMulaiHamil)
				if parseErr != nil {
					return nil, &errorutils.Error{Status: http.StatusBadRequest, Message: "Format bulan_mulai_hamil tidak valid."}
				}
				ih.BulanMulaiHamil = t
			}
			if req.Hpht != nil {
				t, parseErr := time.Parse("2006-01-02", *req.Hpht)
				if parseErr != nil {
					return nil, &errorutils.Error{Status: http.StatusBadRequest, Message: "Format hpht tidak valid."}
				}
				ih.Hpht = t
			}
			if req.StatusKehamilan != nil {
				ih.StatusKehamilan = model.StatusKehamilan(*req.StatusKehamilan)
			}
			err = s.repo.UpdateIbuHamil(ctx, ih)
			if err != nil {
				return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal memperbarui data ibu hamil."}
			}
		}
	}

	if existing.JenisPasien == "Anak" && (req.NamaAnak != nil || req.BeratLahir != nil || req.PanjangLahir != nil || req.HubunganDenganWali != nil || req.IDWali != nil) {
		anak, err := s.repo.GetAnakByPasienID(ctx, req.IDPasien)
		if err != nil {
			return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
		}
		if anak != nil {
			if req.NamaAnak != nil {
				anak.NamaAnak = *req.NamaAnak
			}
			if req.BeratLahir != nil {
				anak.BeratLahir = decimal.NewFromFloat(*req.BeratLahir)
			}
			if req.PanjangLahir != nil {
				anak.PanjangLahir = decimal.NewFromFloat(*req.PanjangLahir)
			}
			if req.HubunganDenganWali != nil {
				anak.HubunganDenganWali = model.HubunganDenganWali(*req.HubunganDenganWali)
			}
			if req.IDWali != nil {
				anak.IDWali = *req.IDWali
			}
			err = s.repo.UpdateAnak(ctx, anak)
			if err != nil {
				return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal memperbarui data anak."}
			}
		}
	}

	s.logAudit(ctx, "PATCH /pasien/"+strconv.Itoa(int(req.IDPasien)), model.TipeAktivitas_DataUpdate, true, "pasien", strconv.Itoa(int(req.IDPasien)), "Berhasil memperbarui data pasien")

	return s.GetByID(ctx, req.IDPasien)
}

func (s *Service) Delete(ctx context.Context, idPasien int32) *errorutils.Error {
	existing, err := s.repo.GetDetailByID(ctx, idPasien)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if existing == nil {
		return &errorutils.Error{Status: http.StatusNotFound, Message: "Pasien tidak ditemukan."}
	}

	err = s.repo.DeletePasien(ctx, idPasien)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Gagal menghapus data pasien."}
	}

	s.logAudit(ctx, "DELETE /pasien/"+strconv.Itoa(int(idPasien)), model.TipeAktivitas_DataDelete, true, "pasien", strconv.Itoa(int(idPasien)), "Berhasil menghapus data pasien")
	return nil
}

func calculateAge(tanggalLahir string) string {
	t, err := time.Parse("2006-01-02", tanggalLahir)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", tanggalLahir)
		if err != nil {
			return ""
		}
	}

	now := time.Now()
	years := now.Year() - t.Year()
	if now.YearDay() < t.YearDay() {
		years--
	}

	if years < 1 {
		months := int(now.Month() - t.Month())
		if now.Day() < t.Day() {
			months--
		}
		if months < 0 {
			months += 12
		}
		return fmt.Sprintf("%d bulan", months)
	}

	return fmt.Sprintf("%d tahun", years)
}

func nullString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (s *Service) isOwnPasien(ctx context.Context, idPasien, idUser int32) bool {
	if idPasien == idUser {
		return true
	}
	anak, err := s.repo.GetAnakByPasienID(ctx, idPasien)
	if err != nil {
		return false
	}
	return anak != nil && anak.IDWali == idUser
}

func nullFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func toInt32Ptr(s *string) *int32 {
	if s == nil {
		return nil
	}
	v, err := strconv.ParseInt(*s, 10, 32)
	if err != nil {
		return nil
	}
	r := int32(v)
	return &r
}
