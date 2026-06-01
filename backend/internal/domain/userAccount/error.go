package userAccount

import "errors"

var (
	ErrStatusVerifikasiPending = errors.New("account belum terverifikasi")
	ErrStatusVerifikasiDitolak = errors.New("verifikasi akun ditolak")
)
