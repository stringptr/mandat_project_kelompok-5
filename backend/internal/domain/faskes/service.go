package faskes

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service interface {
	GetFaskes(ctx context.Context, req *GetFaskesRequest) ([]*FaskesItem, *errorutils.Error)
}
