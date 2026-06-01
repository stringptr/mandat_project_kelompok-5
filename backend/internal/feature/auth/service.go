package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jinzhu/copier"
	"github.com/stringptr/SiGizi/backend/internal/config"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	userSessionDomain "github.com/stringptr/SiGizi/backend/internal/domain/userSession"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/hash"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

type Service struct {
	authRepo    authDomain.Repo
	sessionRepo userSessionDomain.Repo
	userRepo    userAccountDomain.Repo
	jwt         jwtutils.JWT
	cfg         *config.AuthConfig
}

func NewService(
	authRepo authDomain.Repo,
	sessionRepo userSessionDomain.Repo,
	userRepo userAccountDomain.Repo,
	jwt jwtutils.JWT,
	cfg *config.AuthConfig,
) *Service {
	return &Service{
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		jwt:         jwt,
		cfg:         cfg,
	}
}

func (s *Service) Register(ctx context.Context, dataDTO *authDomain.RegisterRequest) *errorutils.Error {
	existingUser, err := s.userRepo.GetByNIK(ctx, dataDTO.NIK)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	if existingUser != nil {
		if existingUser.Email == dataDTO.Email && existingUser.Nik == dataDTO.NIK && existingUser.Password == dataDTO.Password {
			return &errorutils.Error{Status: http.StatusConflict, Message: "Akun dengan Email, NIK, dan Password tersebut sudah terdaftar."}
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

	err = s.userRepo.Create(ctx, &dataModel)
	if err != nil {
		return &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}
	return nil
}

func (s *Service) Login(ctx context.Context, req *authDomain.LoginRequest, ip string) (*authDomain.AuthResponse, *errorutils.Error) {
	var errs []*httputils.ErrorItem

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pengecekan akun. Silahkan dicoba kembali.", Errors: nil}
	}
	if user == nil {
		errs = append(errs, &httputils.ErrorItem{ID: "ERR-AUTH-02", Message: "Email, NIK, atau Password tidak valid"})
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Email, NIK, atau Password tidak valid", Errors: errs}
	}

	ok, err := hash.VerifyHash(req.Password, user.Password)
	if err != nil {
		errs = append(errs, &httputils.ErrorItem{ID: "ERR-AUTH-02", Message: "Email, NIK, atau Password tidak valid"})
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Email, NIK, atau Password tidak valid", Errors: errs}
	}
	if !ok {
		errs = append(errs, &httputils.ErrorItem{ID: "ERR-AUTH-02", Message: "Email, NIK, atau Password tidak valid"})
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Email, NIK, atau Password tidak valid", Errors: errs}
	}

	if user.StatusVerifikasi == model.StatusVerifikasi_Pending {
		return nil, &errorutils.Error{Status: http.StatusUnauthorized, Message: "Akun sedang dalam proses verifikasi. Silahkan dicek secara berkala."}
	}

	roles, err := s.authRepo.GetRoles(ctx, user.IDUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pengecekan akun. Silahkan dicoba kembali."}
	}

	newUUID, err := uuid.NewV7()
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	refreshTokenExpireDate := time.Now().Add(s.cfg.RefreshTokenTTL)
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
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	claim := jwtutils.Claim{
		IDUser: user.IDUser,
		Roles:  roles,
		Email:  user.Email,
		NIK:    user.Nik,
	}

	accessToken, err := s.jwt.EncodeWithTTL(claim, s.cfg.AccessTokenTTL)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan. Silahkan dicoba kembali."}
	}

	return &authDomain.AuthResponse{
		AccessToken:           accessToken,
		RefreshToken:          session.IDSession,
		AccessTokenExpiresIn:  int64(s.cfg.AccessTokenTTL.Seconds()),
		RefreshTokenExpiresIn: int64(s.cfg.RefreshTokenTTL.Seconds()),
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken uuid.UUID, ip string) (*authDomain.AuthResponse, *errorutils.Error) {
	session, err := s.sessionRepo.GetByID(ctx, refreshToken)
	if err != nil {
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
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pembaruan sesi. Silahkan login ulang."}
	}

	user, err := s.userRepo.GetByID(ctx, session.IDUser)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pengecekan akun. Silahkan login ulang."}
	}
	if user == nil {
		return nil, &errorutils.Error{Status: http.StatusNotFound, Message: "Akun tidak dapat ditemukan. Silahkan login ulang."}
	}

	roles, err := s.authRepo.GetRoles(ctx, user.IDUser)
	if err != nil {
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

	refreshTokenExpireDate := time.Now().Add(s.cfg.RefreshTokenTTL)
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
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pembaruan sesi. Silahkan login ulang."}
	}

	accessToken, err := s.jwt.EncodeWithTTL(claim, s.cfg.AccessTokenTTL)
	if err != nil {
		return nil, &errorutils.Error{Status: http.StatusInternalServerError, Message: "Terjadi kesalahan dalam pembaruan sesi. Silahkan login ulang."}
	}

	return &authDomain.AuthResponse{
		AccessToken:           accessToken,
		RefreshToken:          session.IDSession,
		AccessTokenExpiresIn:  int64(s.cfg.AccessTokenTTL.Seconds()),
		RefreshTokenExpiresIn: int64(s.cfg.RefreshTokenTTL.Seconds()),
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken uuid.UUID) *errorutils.Error {
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
	return nil
}
