package lokasi

import (
	"context"
	"net/http"

	lokasiDomain "github.com/stringptr/SiGizi/backend/internal/domain/lokasi"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service struct {
	repo lokasiDomain.Repo
}

func NewService(repo lokasiDomain.Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetLokasi(ctx context.Context, req *lokasiDomain.GetLokasiRequest) ([]*lokasiDomain.LokasiItem, *errorutils.Error) {
	if req.Tipe == "" {
		return nil, &errorutils.Error{
			Status:  http.StatusUnprocessableEntity,
			Message: "Parameter tipe wajib diisi.",
		}
	}

	items, err := s.repo.GetByTipeAndParent(ctx, req.Tipe, req.BagianDari)
	if err != nil {
		return nil, &errorutils.Error{
			Status:  http.StatusInternalServerError,
			Message: "Gagal mengambil data lokasi.",
		}
	}

	if items == nil {
		items = []*lokasiDomain.LokasiItem{}
	}

	return items, nil
}
