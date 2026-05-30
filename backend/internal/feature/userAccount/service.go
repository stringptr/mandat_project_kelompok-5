package userAccount

import (
	"context"

	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Service struct {
	Repository userAccountDomain.Repo
}

func NewService(Repository userAccountDomain.Repo) *Service {
	return &Service{Repository: Repository}
}

func (s *Service) GetAll(ctx context.Context) ([]*model.UserAccount, error) {
	userAccounts, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return userAccounts, nil
}
