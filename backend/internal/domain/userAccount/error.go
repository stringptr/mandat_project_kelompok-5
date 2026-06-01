package userAccount

import "errors"

var (
	ErrAccountStatusPending  = errors.New("account verification still pending")
	ErrAccountStatusRejected = errors.New("account verification rejected")
	ErrNIKEmailExist         = errors.New("combination of email and NIK already used")
	ErrNotDeleted            = errors.New("fail to delete account")
	ErrNotCreated            = errors.New("fail to create account")
)
