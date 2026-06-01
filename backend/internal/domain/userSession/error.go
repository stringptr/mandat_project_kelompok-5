package userSession

import "errors"

var (
	ErrNotCreated = errors.New("fail to create session")
	ErrNotUpdated = errors.New("fail to update session")
)
