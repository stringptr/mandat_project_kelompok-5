package userAccount

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	auditlogDomain "github.com/stringptr/SiGizi/backend/internal/domain/auditlog"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/hash"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Service struct {
	userAccountRepo userAccountDomain.Repo
	authRepo        authDomain.Repo
	auditLogRepo    auditlogDomain.Repo
}

func NewService(userAccountRepo userAccountDomain.Repo, authRepo authDomain.Repo, auditLogRepo auditlogDomain.Repo) *Service {
	return &Service{userAccountRepo: userAccountRepo, authRepo: authRepo, auditLogRepo: auditLogRepo}
}

func (s *Service) logAudit(ctx context.Context, endpoint string, tipeAktivitas model.TipeAktivitas, berhasil bool, tableName string, recordID string, detail string) {
	tipeAktor := model.TipeAktor_User
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		tipeAktor = model.TipeAktor_Anonymous
	}
	s.auditLogRepo.Log(ctx, &model.AuditLog{
		TipeAktor:     &tipeAktor,
		TipeAktivitas: &tipeAktivitas,
		Berhasil:      &berhasil,
		Endpoint:      &endpoint,
		TableName:     &tableName,
		RecordID:      &recordID,
		Detail:        &detail,
	})
}

func (s *Service) GetAllUsers(ctx context.Context, req *userAccountDomain.GetAllUsersRequest) (*userAccountDomain.UserListData, error) {
	users, total, err := s.userAccountRepo.GetAllPaginated(ctx, req.Page, req.PerPage, req.Q, req.Role, req.StatusVerifikasi)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar pengguna: %w", err)
	}

	items := make([]userAccountDomain.UserListItem, len(users))
	for i, u := range users {
		roles, _ := s.authRepo.GetRoles(ctx, u.IDUser)
		items[i] = userAccountDomain.UserListItem{
			IDUser:           u.IDUser,
			Nama:             u.Nama,
			NIK:              u.Nik,
			Email:            u.Email,
			NoHp:             u.NoHp,
			JenisKelamin:     string(u.JenisKelamin),
			StatusVerifikasi: string(u.StatusVerifikasi),
			Roles:            roles,
			IDLokasi:         u.IDLokasi,
			CreatedAt:        u.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        u.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &userAccountDomain.UserListData{
		Users:     items,
		TotalData: total,
		Page:      req.Page,
		PerPage:   req.PerPage,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, idUser int32) (*userAccountDomain.UserDetailResponse, error) {
	u, err := s.userAccountRepo.GetByID(ctx, idUser)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil detail pengguna: %w", err)
	}
	if u == nil {
		return nil, nil
	}

	roles, _ := s.authRepo.GetRoles(ctx, u.IDUser)

	return &userAccountDomain.UserDetailResponse{
		IDUser:           u.IDUser,
		Email:            u.Email,
		NoHp:             u.NoHp,
		Nama:             u.Nama,
		NIK:              u.Nik,
		JenisKelamin:     string(u.JenisKelamin),
		TanggalLahir:     u.TanggalLahir,
		StatusVerifikasi: string(u.StatusVerifikasi),
		IDLokasi:         u.IDLokasi,
		IDPendidikan:     u.IDPendidikan,
		IDPekerjaan:      u.IDPekerjaan,
		IDPendapatan:     u.IDPendapatan,
		JumlahTanggungan: u.JumlahTanggungan,
		Roles:            roles,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}, nil
}

func (s *Service) UpdateUser(ctx context.Context, idUser int32, req *userAccountDomain.UpdateUserRequest) error {
	u, err := s.userAccountRepo.GetByID(ctx, idUser)
	if err != nil {
		return fmt.Errorf("gagal mengambil data pengguna: %w", err)
	}
	if u == nil {
		return userAccountDomain.ErrNotFound
	}

	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.NoHp != nil {
		u.NoHp = *req.NoHp
	}
	if req.Nama != nil {
		u.Nama = *req.Nama
	}
	if req.NIK != nil {
		u.Nik = *req.NIK
	}
	if req.JenisKelamin != nil {
		u.JenisKelamin = model.JenisKelamin(*req.JenisKelamin)
	}
	if req.TanggalLahir != nil {
		u.TanggalLahir = *req.TanggalLahir
	}
	if req.IDLokasi != nil {
		u.IDLokasi = *req.IDLokasi
	}
	if req.IDPendidikan != nil {
		u.IDPendidikan = req.IDPendidikan
	}
	if req.IDPekerjaan != nil {
		u.IDPekerjaan = req.IDPekerjaan
	}
	if req.IDPendapatan != nil {
		u.IDPendapatan = req.IDPendapatan
	}
	if req.JumlahTanggungan != nil {
		u.JumlahTanggungan = req.JumlahTanggungan
	}

	return s.userAccountRepo.Update(ctx, u)
}

func (s *Service) CreateUser(ctx context.Context, req *userAccountDomain.CreateUserRequest) (*userAccountDomain.CreateUserResponse, error) {
	hashedPassword, err := hash.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("gagal mengenkripsi password: %w", err)
	}

	dataModel := model.UserAccount{}
	copier.Copy(&dataModel, &req)
	dataModel.Password = hashedPassword
	dataModel.StatusVerifikasi = model.StatusVerifikasi_Aktif

	idUser, err := s.authRepo.CreateUser(ctx, &dataModel)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat pengguna: %w", err)
	}

	wilayahKerja := int32(0)
	if req.WilayahKerja != nil {
		wilayahKerja = *req.WilayahKerja
	}
	err = s.authRepo.CreateRoleRecord(ctx, idUser, req.Role, req.NoStr, wilayahKerja, req.NoSk)
	if err != nil {
		return nil, fmt.Errorf("gagal menetapkan role pengguna: %w", err)
	}

	s.logAudit(ctx, "POST /v1/users", model.TipeAktivitas_DataInsert, true, "user_account", strconv.Itoa(int(idUser)), "Pengguna berhasil dibuat oleh SuperAdmin")
	return &userAccountDomain.CreateUserResponse{IDUser: idUser}, nil
}

func (s *Service) UpdateUserRole(ctx context.Context, idUser int32, req *userAccountDomain.UpdateUserRoleRequest) error {
	u, err := s.userAccountRepo.GetByID(ctx, idUser)
	if err != nil {
		return fmt.Errorf("gagal mengambil data pengguna: %w", err)
	}
	if u == nil {
		return userAccountDomain.ErrNotFound
	}

	err = s.authRepo.DeleteRoleRecords(ctx, idUser)
	if err != nil {
		return fmt.Errorf("gagal menghapus role lama: %w", err)
	}

	wilayahKerja := int32(0)
	if req.WilayahKerja != nil {
		wilayahKerja = *req.WilayahKerja
	}
	err = s.authRepo.CreateRoleRecord(ctx, idUser, req.Role, req.NoStr, wilayahKerja, req.NoSk)
	if err != nil {
		return fmt.Errorf("gagal menetapkan role baru: %w", err)
	}

	s.logAudit(ctx, "PATCH /v1/users/"+strconv.Itoa(int(idUser))+"/role", model.TipeAktivitas_DataUpdate, true, "user_account", strconv.Itoa(int(idUser)), "Role pengguna diubah menjadi "+req.Role)
	return nil
}

func (s *Service) GetAuditLogs(ctx context.Context, filter *userAccountDomain.AuditLogFilter) (*userAccountDomain.AuditLogListData, error) {
	domainFilter := &auditlogDomain.AuditLogFilter{
		Page:    filter.Page,
		PerPage: filter.PerPage,
	}
	logs, total, err := s.auditLogRepo.GetAll(ctx, domainFilter)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil audit log: %w", err)
	}

	items := make([]userAccountDomain.AuditLogItem, len(logs))
	for i, l := range logs {
		tipeAktor := ""
		if l.TipeAktor != nil {
			tipeAktor = string(*l.TipeAktor)
		}
		tipeAktivitas := ""
		if l.TipeAktivitas != nil {
			tipeAktivitas = string(*l.TipeAktivitas)
		}
		endpoint := ""
		if l.Endpoint != nil {
			endpoint = *l.Endpoint
		}
		tableName := ""
		if l.TableName != nil {
			tableName = *l.TableName
		}
		recordID := ""
		if l.RecordID != nil {
			recordID = *l.RecordID
		}
		detail := ""
		if l.Detail != nil {
			detail = *l.Detail
		}
		ipAddress := ""
		if l.IPAddress != nil {
			ipAddress = *l.IPAddress
		}
		userAgent := ""
		if l.UserAgent != nil {
			userAgent = *l.UserAgent
		}

		items[i] = userAccountDomain.AuditLogItem{
			IDLog:          l.IDLog,
			TipeAktor:      tipeAktor,
			IDUser:         l.IDUser,
			TipeAktivitas:  tipeAktivitas,
			Berhasil:       l.Berhasil,
			Endpoint:       endpoint,
			TableName:      tableName,
			RecordID:       recordID,
			Detail:         detail,
			IPAddress:      ipAddress,
			UserAgent:      userAgent,
			WaktuAktivitas: l.WaktuAktivitas.Format(time.RFC3339),
		}
	}

	return &userAccountDomain.AuditLogListData{
		AuditLogs: items,
		TotalData: total,
		Page:      filter.Page,
		PerPage:   filter.PerPage,
	}, nil
}
