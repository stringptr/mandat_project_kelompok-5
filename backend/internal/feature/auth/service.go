package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jinzhu/copier"
	"github.com/stringptr/SiGizi/backend/internal/config"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	userSessionDomain "github.com/stringptr/SiGizi/backend/internal/domain/userSession"
	"github.com/stringptr/SiGizi/backend/internal/hash"
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

func (s *Service) Register(ctx context.Context, dataDTO *authDomain.RegisterRequest) error {
	existingUser, err := s.userRepo.GetByNIK(ctx, dataDTO.NIK)
	if err != nil {
		return err
	}
	if existingUser != nil {
		if existingUser.Email == dataDTO.Email && existingUser.Nik == dataDTO.NIK {
			return fmt.Errorf("tidak bisa daftar dengan data yang diberikan")
		}
	}

	hashedPassword, err := hash.Hash(dataDTO.Password)
	if err != nil {
		return err
	}

	dataModel := model.UserAccount{}
	copier.Copy(&dataModel, &dataDTO)
	dataModel.Password = hashedPassword

	err = s.userRepo.Create(ctx, &dataModel)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Login(ctx context.Context, req *authDomain.LoginRequest, ip string) (*authDomain.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	if user == nil {
		return nil, authDomain.ErrInvalidCredentials
	}

	ok, err := hash.VerifyHash(req.Password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("verify hash: %w", err)
	}
	if !ok {
		return nil, authDomain.ErrInvalidCredentials
	}

	if user.StatusVerifikasi == model.StatusVerifikasi_Pending {
		return nil, userAccountDomain.ErrStatusVerifikasiPending
	}

	roles, err := s.authRepo.GetRoles(ctx, user.IDUser)
	if err != nil {
		return nil, fmt.Errorf("get roles: %w", err)
	}

	newUUID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("UUID generation: %w", err)
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
		return nil, fmt.Errorf("create session: %w", err)
	}

	claim := jwtutils.Claim{
		IDUser: user.IDUser,
		Roles:  roles,
		Email:  user.Email,
		NIK:    user.Nik,
	}

	accessToken, err := s.jwt.EncodeWithTTL(claim, s.cfg.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("encode token: %w", err)
	}

	return &authDomain.AuthResponse{
		AccessToken:           accessToken,
		RefreshToken:          session.IDSession,
		AccessTokenExpiresIn:  int64(s.cfg.AccessTokenTTL.Seconds()),
		RefreshTokenExpiresIn: int64(s.cfg.RefreshTokenTTL.Seconds()),
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken uuid.UUID, ip string) (*authDomain.AuthResponse, error) {
	log.Print(refreshToken)
	session, err := s.sessionRepo.GetByID(ctx, refreshToken)
	if err != nil {
		log.Print(err)
		return nil, authDomain.ErrSessionNotFound
	}
	if session == nil {
		return nil, authDomain.ErrSessionNotFound
	}

	if session.StatusSession != model.StatusSession_Aktif {
		return nil, authDomain.ErrSessionInactive
	}

	session.StatusSession = model.StatusSession_Dicabut
	err = s.sessionRepo.Update(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("invalidate session: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, session.IDUser)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, authDomain.ErrInvalidCredentials
	}

	roles, err := s.authRepo.GetRoles(ctx, user.IDUser)
	if err != nil {
		return nil, fmt.Errorf("get roles: %w", err)
	}

	claim := jwtutils.Claim{
		IDUser: user.IDUser,
		Roles:  roles,
		Email:  user.Email,
		NIK:    user.Nik,
	}

	newUUID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("UUID generation: %w", err)
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
		return nil, fmt.Errorf("create session: %w", err)
	}

	accessToken, err := s.jwt.EncodeWithTTL(claim, s.cfg.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("encode token: %w", err)
	}

	return &authDomain.AuthResponse{
		AccessToken:           accessToken,
		RefreshToken:          session.IDSession,
		AccessTokenExpiresIn:  int64(s.cfg.AccessTokenTTL.Seconds()),
		RefreshTokenExpiresIn: int64(s.cfg.RefreshTokenTTL.Seconds()),
	}, nil
}
