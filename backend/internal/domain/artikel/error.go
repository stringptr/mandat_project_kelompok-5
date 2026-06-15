package artikel

import "errors"

var (
	ErrArtikelNotFound    = errors.New("artikel tidak ditemukan")
	ErrArtikelNotCreated  = errors.New("gagal membuat artikel")
	ErrArtikelNotUpdated  = errors.New("gagal memperbarui artikel")
	ErrArtikelNotDeleted  = errors.New("gagal menghapus artikel")
	ErrArtikelNotReviewed = errors.New("gagal mereview artikel")
)
