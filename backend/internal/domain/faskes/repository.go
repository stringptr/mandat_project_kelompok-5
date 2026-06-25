package faskes

import "context"

type Repo interface {
	GetAll(ctx context.Context, search string) ([]*FaskesItem, error)
}
