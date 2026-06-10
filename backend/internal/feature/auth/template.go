package auth

import "fmt"

func accountActivatedEmail(name string) (subject, body string) {
	subject = "Akun Anda telah diaktifkan - SiGizi"
	body = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: sans-serif; padding: 20px;">
<h2>Selamat, %s!</h2>
<p>Akun Anda telah berhasil diaktifkan oleh admin.</p>
<p>Silakan login menggunakan email/NIK dan password yang telah didaftarkan untuk menggunakan layanan SiGizi.</p>
<hr>
<p style="color: #666; font-size: 12px;">Email ini dikirim secara otomatis. Harap tidak membalas email ini.</p>
</body>
</html>`, name)
	return
}

func accountRejectedEmail(name, reason string) (subject, body string) {
	subject = "Akun Anda ditolak - SiGizi"
	body = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: sans-serif; padding: 20px;">
<h2>Yth. %s</h2>
<p>Mohon maaf, akun Anda tidak dapat diverifikasi.</p>
<p><strong>Alasan:</strong> %s</p>
<p>Silakan hubungi admin untuk informasi lebih lanjut.</p>
<hr>
<p style="color: #666; font-size: 12px;">Email ini dikirim secara otomatis. Harap tidak membalas email ini.</p>
</body>
</html>`, name, reason)
	return
}
