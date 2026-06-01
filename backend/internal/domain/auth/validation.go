package auth

import "strings"

func MapValidationError(location, message string) (string, string, bool) {
	switch location {
	case "body.email":
		return "ERR-VAL-01", "Alamat email tidak valid", true
	case "body.password":
		return "ERR-VAL-02", "Password tidak valid atau kurang panjang", true
	case "body.no_hp":
		return "ERR-VAL-04", "Nomor telepon tidak valid", true
	case "body.nama":
		return "ERR-VAL-05", "Nama tidak valid", true
	case "body.tanggal_lahir":
		return "ERR-VAL-06", "Tanggal tidak valid", true
	case "body.jenis_kelamin", "body.id_lokasi", "body.id_pendidikan",
		"body.id_pekerjaan", "body.id_pendapatan", "body.jumlah_tanggungan":
		return "ERR-VAL-07", "Pilihan tidak valid", true
	}
	if strings.Contains(message, "required") {
		return "ERR-VAL-03", "Dokumen tidak lengkap", true
	}
	return "", "", false
}
