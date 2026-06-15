package imunisasi

import "errors"

var (
	ErrImunisasiNotFound     = errors.New("data jadwal imunisasi tidak ditemukan")
	ErrPasienNotFound        = errors.New("data pasien tidak ditemukan")
	ErrImunisasiNotCreated   = errors.New("gagal membuat jadwal imunisasi")
	ErrImunisasiNotUpdated   = errors.New("gagal memperbarui jadwal imunisasi")
	ErrImunisasiNotDeleted   = errors.New("gagal menghapus jadwal imunisasi")
	ErrImunisasiNotRealisasi = errors.New("gagal mencatat realisasi imunisasi")
)
