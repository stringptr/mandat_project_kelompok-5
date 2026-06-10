package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jinzhu/copier"
	"github.com/stringptr/SiGizi/backend/internal/config"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	bannedipDomain "github.com/stringptr/SiGizi/backend/internal/domain/bannedip"
	jwtblacklistDomain "github.com/stringptr/SiGizi/backend/internal/domain/jwtblacklist"
	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	userSessionDomain "github.com/stringptr/SiGizi/backend/internal/domain/userSession"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/hash"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

type Service struct {
	authRepo        authDomain.Repo
	sessionRepo     userSessionDomain.Repo
	userRepo        userAccountDomain.Repo
	jwt             jwtutils.JWT
	authCfg         *config.AuthConfig
	restrictAuthCfg *config.RestrictAuthConfig
	banRepo         bannedipDomain.Repo
	blacklistRepo   jwtblacklistDomain.Repo
}

func NewService(
	authRepo authDomain.Repo,
	sessionRepo userSessionDomain.Repo,
	userRepo userAccountDomain.Repo,
	jwt jwtutils.JWT,
	authCfg *config.AuthConfig,
	restrictAuthCfg *config.RestrictAuthConfig,
	banRepo bannedipDomain.Repo,
	blacklistRepo jwtblacklistDomain.Repo,
) *Service {
	return &Service{
		authRepo:        authRepo,
		sessionRepo:     sessionRepo,
		userRepo:        userRepo,
		jwt:             jwt,
		authCfg:         authCfg,
		restrictAuthCfg: restrictAuthCfg,
		banRepo:         banRepo,
		blacklistRepo:   blacklistRepo,
	}
}

func (s *Service) Register(ctx context.Context, dataDTO *authDomain.RegisterRequest, ip string) *errorutils.Error {
	if ip != "" {
		info, _ := s.banRepo.GetBanInfo(ctx, ip)
		if info != nil && time.Now().Before(info.ExpiresAt) {
			remaining := time.Until(info.ExpiresAt)
			errs := []*httputils.ErrorItem{{
				ID:      "ERR-LOCK-01",
				Message: "Percobaan gagal dalam 3 kali. Akun dikunci untuk sementara. Silahkan coba lagi dalam 15 menit.",
			}}
			return &errorutils.Error{Status: http.StatusForbidden, Message: fmt.Sprintf("Akses ditolak. Terlalu banyak percobaan. Silahkan coba lagi dalam %d menit %d detik.", int(remaining.Minutes()), int(remaining.Seconds())%60), Errors: errs}
		}
	}

	existingUser, err := s.userRepo.GetByNIKEmail(ctx, dataDTO.NIK, dataDTO.Email)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if existingUser != nil {
		if ip != "" {
			s.banRepo.IncrementAttempt(ctx, ip, s.restrictAuthCfg.MaxAttempt, s.restrictAuthCfg.Duration)
		}
		if existingUser.Email == dataDTO.Email && existingUser.Nik == dataDTO.NIK && existingUser.Password == dataDTO.Password {
			return &errorutils.Error{Status: http.StatusConflict, Message: "Akun dengan Email, NIK, dan Password tersebut sudah terdaftar."}
		}
		if err != nil {
			return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
		}
		return &errorutils.Error{Status: http.StatusConflict, Message: "Tidak dapat mendaftar dengan identitas yang diberikan."}
	}

	hashedPassword, err := hash.Hash(dataDTO.Password)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	dataModel := model.UserAccount{}
	copier.Copy(&dataModel, &dataDTO)
	dataModel.Password = hashedPassword

	idUser, err := s.authRepo.CreateUser(ctx, &dataModel)
	if err != nil {
		if isDBConnectionError(err) {
			errs := []*httputils.ErrorItem{{ID: "ERR-SYS-01", Message: "Layanan tidak tersedia. Silakan coba lagi."}}
			return &errorutils.Error{Status: http.StatusServiceUnavailable, Message: "Layanan tidak tersedia. Silakan coba lagi.", Errors: errs}
		}
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	if dataDTO.Role != "" {
		wilayahKerja := int32(0)
		if dataDTO.WilayahKerja != nil {
			wilayahKerja = *dataDTO.WilayahKerja
		}
		err = s.authRepo.CreateRoleRecord(ctx, idUser, dataDTO.Role, dataDTO.NoStr, wilayahKerja, dataDTO.NoSk)
		if err != nil {
			if isDBConnectionError(err) {
				errs := []*httputils.ErrorItem{{ID: "ERR-SYS-01", Message: "Layanan tidak tersedia. Silakan coba lagi."}}
				return &errorutils.Error{Status: http.StatusServiceUnavailable, Message: "Layanan tidak tersedia. Silakan coba lagi.", Errors: errs}
			}
			return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
		}
	}

	return nil
}

func (s *Service) Login(ctx context.Context, req *authDomain.LoginRequest, ip string) (*authDomain.AuthResponse, *errorutils.Error) {
	var errs []*httputils.ErrorItem

	if ip != "" {
		info, _ := s.banRepo.GetBanInfo(ctx, ip)
		if info != nil && time.Now().Before(info.ExpiresAt) {
			remaining := time.Until(info.ExpiresAt)
			errs := []*httputils.ErrorItem{{
				ID:      "ERR-LOCK-01",
				Message: "Percobaan login gagal dalam 3 kali. Akun dikunci untuk sementara. Silahkan coba lagi dalam 15 menit.",
			}}
			return nil, &errorutils.Error{Status: http.StatusForbidden, Message: fmt.Sprintf("Akses ditolak. Terlalu banyak percobaan login. Silahkan coba lagi dalam %d menit %d detik.", int(remaining.Minutes()), int(remaining.Seconds())%60), Errors: errs}
		}
	}

	user, err := s.authRepo.GetByEmailNIK(ctx, req.Email, req.NIK)
	if err != nil {
		if isDBConnectionError(err) {
			errs = append(errs, &httputils.ErrorItem{ID: "ERR-SYS-01", Message: "Layanan tidak tersedia. Silakan coba lagi."})
			return nil, &errorutils.Error{Status: http.StatusServiceUnavailable, Message: "Layanan tidak tersedia. Silakan coba lagi.", Errors: errs}
		}
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pengecekan akun. Silahkan dicoba kembali.", Errors: nil}
	}
	if user == nil {
		if ip != "" {
			s.banRepo.IncrementAttempt(ctx, ip, s.restrictAuthCfg.MaxAttempt, s.restrictAuthCfg.Duration)
		}
		errs = append(errs, &httputils.ErrorItem{ID: "ERR-AUTH-02", Message: "Email, NIK, atau Password tidak valid"})
		return nil, &errorutils.Error{Status: http.StatusUnauthorized, Message: "Email, NIK, atau Password tidak valid", Errors: errs}
	}

	ok, err := hash.VerifyHash(req.Password, user.Password)
	if err != nil {
		if ip != "" {
			s.banRepo.IncrementAttempt(ctx, ip, s.restrictAuthCfg.MaxAttempt, s.restrictAuthCfg.Duration)
		}
		errs = append(errs, &httputils.ErrorItem{ID: "ERR-AUTH-02", Message: "Email, NIK, atau Password tidak valid"})
		return nil, &errorutils.Error{Status: http.StatusUnauthorized, Message: "Email, NIK, atau Password tidak valid", Errors: errs}
	}
	if !ok {
		if ip != "" {
			s.banRepo.IncrementAttempt(ctx, ip, s.restrictAuthCfg.MaxAttempt, s.restrictAuthCfg.Duration)
		}
		errs = append(errs, &httputils.ErrorItem{ID: "ERR-AUTH-02", Message: "Email, NIK, atau Password tidak valid"})
		return nil, &errorutils.Error{Status: http.StatusUnauthorized, Message: "Email, NIK, atau Password tidak valid", Errors: errs}
	}

	if user.StatusVerifikasi == model.StatusVerifikasi_Pending {
		return nil, &errorutils.Error{Status: http.StatusUnauthorized, Message: "Akun sedang dalam proses verifikasi. Silahkan dicek secara berkala."}
	}

	s.banRepo.ClearAttempts(ctx, ip)

	roles, err := s.authRepo.GetRoles(ctx, user.IDUser)
	if err != nil {
		if isDBConnectionError(err) {
			errs := []*httputils.ErrorItem{{ID: "ERR-SYS-01", Message: "Layanan tidak tersedia. Silakan coba lagi."}}
			return nil, &errorutils.Error{Status: http.StatusServiceUnavailable, Message: "Layanan tidak tersedia. Silakan coba lagi.", Errors: errs}
		}
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pengecekan akun. Silahkan dicoba kembali."}
	}

	newUUID, err := uuid.NewV7()
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	refreshTokenExpireDate := time.Now().Add(s.authCfg.RefreshTokenTTL)
	session := &model.UserSession{
		IDSession:     newUUID,
		IDUser:        user.IDUser,
		StatusSession: model.StatusSession_Aktif,
		IPAddress:     &ip,
		ExpiredAt:     refreshTokenExpireDate,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err = s.sessionRepo.Create(ctx, session)
	if err != nil {
		if isDBConnectionError(err) {
			errs := []*httputils.ErrorItem{{ID: "ERR-SYS-01", Message: "Layanan tidak tersedia. Silakan coba lagi."}}
			return nil, &errorutils.Error{Status: http.StatusServiceUnavailable, Message: "Layanan tidak tersedia. Silakan coba lagi.", Errors: errs}
		}
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	claim := jwtutils.Claim{
		IDUser: user.IDUser,
		Roles:  roles,
		Email:  user.Email,
		NIK:    user.Nik,
	}

	accessToken, err := s.jwt.EncodeWithTTL(claim, s.authCfg.AccessTokenTTL)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	return &authDomain.AuthResponse{
		AccessToken:           accessToken,
		RefreshToken:          session.IDSession,
		AccessTokenExpiresIn:  int64(s.authCfg.AccessTokenTTL.Seconds()),
		RefreshTokenExpiresIn: int64(s.authCfg.RefreshTokenTTL.Seconds()),
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken uuid.UUID, ip string) (*authDomain.AuthResponse, *errorutils.Error) {
	session, err := s.sessionRepo.GetByID(ctx, refreshToken)
	if err != nil {
		if isDBConnectionError(err) {
			errs := []*httputils.ErrorItem{{ID: "ERR-SYS-01", Message: "Layanan tidak tersedia. Silakan coba lagi."}}
			return nil, &errorutils.Error{Status: http.StatusServiceUnavailable, Message: "Layanan tidak tersedia. Silakan coba lagi.", Errors: errs}
		}
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pengecekan sesi login. Silahkan login ulang."}
	}
	if session == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Sesi login tidak dapat ditemukan. Silahkan login ulang."}
	}

	if session.StatusSession != model.StatusSession_Aktif {
		return nil, &errorutils.Error{Status: http.StatusForbidden, Message: "Sesi login telah kadaluwarsa. Silahkan login ulang."}
	}

	session.StatusSession = model.StatusSession_Dicabut
	err = s.sessionRepo.Update(ctx, session)
	if err != nil {
		if isDBConnectionError(err) {
			errs := []*httputils.ErrorItem{{ID: "ERR-SYS-01", Message: "Layanan tidak tersedia. Silakan coba lagi."}}
			return nil, &errorutils.Error{Status: http.StatusServiceUnavailable, Message: "Layanan tidak tersedia. Silakan coba lagi.", Errors: errs}
		}
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pembaruan sesi. Silahkan login ulang."}
	}

	user, err := s.userRepo.GetByID(ctx, session.IDUser)
	if err != nil {
		if isDBConnectionError(err) {
			errs := []*httputils.ErrorItem{{ID: "ERR-SYS-01", Message: "Layanan tidak tersedia. Silakan coba lagi."}}
			return nil, &errorutils.Error{Status: http.StatusServiceUnavailable, Message: "Layanan tidak tersedia. Silakan coba lagi.", Errors: errs}
		}
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pengecekan akun. Silahkan login ulang."}
	}
	if user == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Akun tidak dapat ditemukan. Silahkan login ulang."}
	}

	roles, err := s.authRepo.GetRoles(ctx, user.IDUser)
	if err != nil {
		if isDBConnectionError(err) {
			errs := []*httputils.ErrorItem{{ID: "ERR-SYS-01", Message: "Layanan tidak tersedia. Silakan coba lagi."}}
			return nil, &errorutils.Error{Status: http.StatusServiceUnavailable, Message: "Layanan tidak tersedia. Silakan coba lagi.", Errors: errs}
		}
		return nil, &errorutils.Error{Status: http.StatusForbidden, Message: "Terjadi kesalahan dalam pengecekan akun. Silahkan login ulang."}
	}

	claim := jwtutils.Claim{
		IDUser: user.IDUser,
		Roles:  roles,
		Email:  user.Email,
		NIK:    user.Nik,
	}

	newUUID, err := uuid.NewV7()
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pembaruan sesi. Silahkan login ulang."}
	}

	refreshTokenExpireDate := time.Now().Add(s.authCfg.RefreshTokenTTL)
	session = &model.UserSession{
		IDSession:     newUUID,
		IDUser:        user.IDUser,
		StatusSession: model.StatusSession_Aktif,
		IPAddress:     &ip,
		ExpiredAt:     refreshTokenExpireDate,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err = s.sessionRepo.Create(ctx, session)
	if err != nil {
		if isDBConnectionError(err) {
			errs := []*httputils.ErrorItem{{ID: "ERR-SYS-01", Message: "Layanan tidak tersedia. Silakan coba lagi."}}
			return nil, &errorutils.Error{Status: http.StatusServiceUnavailable, Message: "Layanan tidak tersedia. Silakan coba lagi.", Errors: errs}
		}
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pembaruan sesi. Silahkan login ulang."}
	}

	accessToken, err := s.jwt.EncodeWithTTL(claim, s.authCfg.AccessTokenTTL)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pembaruan sesi. Silahkan login ulang."}
	}

	return &authDomain.AuthResponse{
		AccessToken:           accessToken,
		RefreshToken:          session.IDSession,
		AccessTokenExpiresIn:  int64(s.authCfg.AccessTokenTTL.Seconds()),
		RefreshTokenExpiresIn: int64(s.authCfg.RefreshTokenTTL.Seconds()),
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken uuid.UUID, accessTokenJTI string) *errorutils.Error {
	session, err := s.sessionRepo.GetByID(ctx, refreshToken)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan."}
	}
	if session == nil {
		return &errorutils.Error{Status: http.StatusNotFound, Message: "Sesi login tidak dapat ditemukan."}
	}
	session.StatusSession = model.StatusSession_Dicabut

	err = s.sessionRepo.Update(ctx, session)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan."}
	}

	if accessTokenJTI != "" {
		s.blacklistRepo.Blacklist(ctx, accessTokenJTI, session.IDUser, "logout", s.authCfg.AccessTokenTTL)
	}

	return nil
}

func (s *Service) VerifyUser(ctx context.Context, req *authDomain.VerifyUserRequest) *errorutils.Error {
	user, err := s.userRepo.GetByID(ctx, req.IDUser)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if user == nil {
		return &errorutils.Error{Status: http.StatusNotFound, Message: "Akun tidak ditemukan."}
	}

	var newStatus model.StatusVerifikasi
	switch req.Status {
	case "Aktif":
		newStatus = model.StatusVerifikasi_Aktif
	case "Ditolak":
		newStatus = model.StatusVerifikasi_Ditolak
	default:
		return &errorutils.Error{Status: http.StatusBadRequest, Message: "Status yang dimasukkan tidak valid."}
	}

	err = s.userRepo.UpdateStatusVerifikasi(ctx, req.IDUser, newStatus)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	return nil
}

func isDBConnectionError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "no such host") ||
		strings.Contains(err.Error(), "network is unreachable") ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, context.DeadlineExceeded)
}
