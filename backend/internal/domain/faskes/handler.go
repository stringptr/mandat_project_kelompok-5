package faskes

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler interface {
	GetFaskes(ctx context.Context, input *GetFaskesRequest) (*httputils.APIResponseOutput[[]*FaskesItem], error)
}
