package faskes

import (
	"context"
	"net/http"

	faskesDomain "github.com/stringptr/SiGizi/backend/internal/domain/faskes"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service struct {
	repo faskesDomain.Repo
}

func NewService(repo faskesDomain.Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetFaskes(ctx context.Context, req *faskesDomain.GetFaskesRequest) ([]*faskesDomain.FaskesItem, *errorutils.Error) {
	items, err := s.repo.GetAll(ctx, req.Search)
	if err != nil {
		return nil, &errorutils.Error{
			Status:  http.StatusInternalServerError,
			Message: "Gagal mengambil data fasilitas kesehatan.",
		}
	}

	if items == nil {
		items = []*faskesDomain.FaskesItem{}
	}

	return items, nil
}
