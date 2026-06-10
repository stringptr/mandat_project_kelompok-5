package pasien

import "errors"

var (
	ErrUserNotFound      = errors.New("user tidak ditemukan")
	ErrPasienNotFound    = errors.New("pasien tidak ditemukan")
	ErrAlreadyPasien     = errors.New("user sudah terdaftar sebagai pasien")
	ErrNotCreated        = errors.New("gagal membuat data pasien")
	ErrNotUpdated        = errors.New("gagal memperbarui data pasien")
	ErrNotDeleted        = errors.New("gagal menghapus data pasien")
	ErrPosyanduNotFound  = errors.New("posyandu tidak ditemukan")
)
