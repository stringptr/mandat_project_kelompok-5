package userAccount

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/hash"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Service struct {
	Repository userAccount.Repo
}

func NewService(Repository userAccount.Repo) *Service {
	return &Service{Repository: Repository}
}

func (s *Service) GetAll(ctx context.Context) ([]*model.UserAccount, error) {
	userAccounts, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return userAccounts, nil
}

func (s *Service) Register(ctx context.Context, dataDTO *userAccount.RegisterRequestDTO) error {
	existingUser, err := s.Repository.GetByNIK(ctx, dataDTO.NIK)
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

	err = s.Repository.Create(ctx, &dataModel)
	if err != nil {
		return err
	}
	return nil
}
