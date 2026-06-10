package userAccount

import (
	"context"
	"fmt"
	"time"

	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Service struct {
	userAccountRepo userAccountDomain.Repo
	authRepo        authDomain.Repo
}

func NewService(userAccountRepo userAccountDomain.Repo, authRepo authDomain.Repo) *Service {
	return &Service{userAccountRepo: userAccountRepo, authRepo: authRepo}
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
