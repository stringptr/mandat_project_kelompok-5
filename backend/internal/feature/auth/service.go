package auth

import (
	"context"
	"fmt"
	"github.com/stringptr/SiGizi/backend/internal/config"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	userSessionDomain "github.com/stringptr/SiGizi/backend/internal/domain/userSession"
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
