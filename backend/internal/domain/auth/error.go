package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrSessionNotFound    = errors.New("session tidak ditemukan")
	ErrSessionExpired     = errors.New("session sudah kadaluwarsa")
	ErrSessionInactive    = errors.New("session tidak aktif")
)
