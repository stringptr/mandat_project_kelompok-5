"""
Generate Data Dummy Imunisasi (Skala 10.000)
Total ~100.000-200.000 baris untuk semua tabel
"""

import csv
import os
import random
import hashlib
import uuid
from datetime import datetime, timedelta, date
from textwrap import wrap

# ─── Konfigurasi ───────────────────────────────────────────────
OUTPUT_DIR = "migrations/data"
SEED = 42
random.seed(SEED)

TOTAL_USER = 5000
TOTAL_DINKES = 5
TOTAL_BIDAN = 100
TOTAL_KADER = 200
TOTAL_IBU_HAMIL = 1000
TOTAL_ANAK = 2000
TOTAL_PASIEN_DEWASA = 500
TOTAL_ARTIKEL = 20
TOTAL_FASKES = 10
TOTAL_POSYANDU = 50
TOTAL_VAKSIN_PER_ANAK = 10

# ─── Helper ─────────────────────────────────────────────────────
def fast_hash(password):
    """Fast deterministic hash for seed data (not cryptographically secure for prod)."""
    salt = hashlib.md5(str(random.getrandbits(256)).encode()).hexdigest()[:22]
    key = hashlib.sha256((salt + password).encode()).hexdigest()
    return f"{salt}${key}"

def random_date(start, end):
    return start + timedelta(days=random.randint(0, (end - start).days))

def csv_writer(filename, fieldnames, rows):
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    path = os.path.join(OUTPUT_DIR, filename)
    with open(path, "w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=fieldnames)
        w.writeheader()
        w.writerows(rows)
    print(f"  {filename}: {len(rows)} baris")
    return path

# ─── Nama & Data Kependudukan ──────────────────────────────────
NAMA_DEPAN_LAKI = [
    "Agus", "Bambang", "Cahyo", "Dwi", "Eko", "Fajar", "Gunawan", "Hadi",
    "Indra", "Joko", "Kurniawan", "Lukman", "Mulyono", "Nugroho", "Prasetyo",
    "Rudi", "Slamet", "Teguh", "Untung", "Wahyu", "Yudi", "Zainal",
    "Ahmad", "Budi", "Candra", "Dimas", "Edi", "Feri", "Gilang", "Hendra",
]
NAMA_DEPAN_PEREMPUAN = [
    "Ani", "Bunga", "Citra", "Dewi", "Endang", "Fitri", "Gita", "Heni",
    "Indah", "Juwita", "Kartika", "Lestari", "Maya", "Nina", "Oktavia",
    "Putri", "Rina", "Sari", "Titi", "Utami", "Wulan", "Yuni", "Zahra",
    "Aisyah", "Bella", "Cindy", "Dian", "Elsa", "Fany", "Hana",
]
NAMA_BELAKANG = [
    "Susanti", "Wijaya", "Kusuma", "Pratama", "Santoso", "Wibowo", "Hartono",
    "Setiawan", "Ningsih", "Utami", "Handayani", "Saputra", "Purnomo",
    "Hidayat", "Rahayu", "Maryati", "Astuti", "Wulandari", "Siregar",
    "Nasution", "Harahap", "Simanjuntak", "Sembiring", "Ginting", "Tarigan",
]

def random_nama(jk):
    if jk == "Laki-Laki":
        return f"{random.choice(NAMA_DEPAN_LAKI)} {random.choice(NAMA_BELAKANG)}"
    return f"{random.choice(NAMA_DEPAN_PEREMPUAN)} {random.choice(NAMA_BELAKANG)}"

def random_email(nama, idx):
    safe = nama.lower().replace(" ", "").replace("-", "")
    return f"{safe}{idx}@gmail.com"

def random_nik():
    return f"{random.randint(1000000000000000, 9999999999999999)}"

def random_no_hp():
    return f"08{random.randint(100000000, 999999999):09d}"[:13]

def random_password():
    chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%"
    return "".join(random.choices(chars, k=random.randint(8, 12)))

# ─── 1. LOKASI ─────────────────────────────────────────────────
def gen_lokasi():
    prov = {"id": 1, "nama": "Jawa Tengah", "tipe": "Provinsi", "bagian_dari": None}
    kab_list = [
        (2, "Kabupaten Cilacap"), (3, "Kabupaten Banyumas"), (4, "Kabupaten Purbalingga"),
        (5, "Kabupaten Banjarnegara"), (6, "Kabupaten Kebumen"), (7, "Kabupaten Purworejo"),
        (8, "Kabupaten Wonosobo"), (9, "Kabupaten Magelang"), (10, "Kabupaten Boyolali"),
        (11, "Kabupaten Klaten"), (12, "Kabupaten Sukoharjo"), (13, "Kabupaten Wonogiri"),
        (14, "Kabupaten Karanganyar"), (15, "Kabupaten Sragen"), (16, "Kabupaten Grobogan"),
        (17, "Kabupaten Blora"), (18, "Kabupaten Rembang"), (19, "Kabupaten Pati"),
        (20, "Kota Semarang"), (21, "Kota Surakarta"), (22, "Kota Salatiga"),
    ]
    rows = [prov]
    kid = 2
    kec_id = 23
    for kab_id, kab_nama in kab_list:
        rows.append({"id": kab_id, "nama": kab_nama, "tipe": "Kabupaten", "bagian_dari": 1})
        n_kec = random.randint(2, 3)
        for _ in range(n_kec):
            rows.append({"id": kec_id, "nama": f"Kecamatan {kab_nama.split()[-1]} {random.choice(['Timur','Barat','Utara','Selatan','Pusat'])}", "tipe": "Kecamatan", "bagian_dari": kab_id})
            kec_id += 1
    # tambah kelurahan
    kel_id = kec_id
    for kec in [r for r in rows if r["tipe"] == "Kecamatan"]:
        n_kel = random.randint(1, 2)
        for _ in range(n_kel):
            rows.append({"id": kel_id, "nama": f"Kelurahan {kec['nama'].split()[-1]} {random.choice(['Hilir','Hulu','Tengah','Indah','Mulya'])}", "tipe": "Kelurahan", "bagian_dari": kec["id"]})
            kel_id += 1
    total = len(rows)
    rows = rows[:total]
    return rows

# ─── 2. PENDIDIKAN ─────────────────────────────────────────────
def gen_pendidikan():
    data = [
        (1, "Tidak/Belum Sekolah", "Tidak Sekolah", 0),
        (2, "SD/MI Sederajat", "SD", 6),
        (3, "SMP/MTs Sederajat", "SMP", 3),
        (4, "SMA/SMK/MA Sederajat", "SMA", 3),
        (5, "D1-D3", "Diploma", 3),
        (6, "S1/D4", "Sarjana", 4),
        (7, "S2", "Magister", 2),
        (8, "S3", "Doktor", 4),
    ]
    return [{"id_pendidikan": d[0], "nama_pendidikan": d[1], "jenjang": d[2], "lama_tahun": d[3]} for d in data]

# ─── 3. PEKERJAAN ──────────────────────────────────────────────
def gen_pekerjaan():
    data = [
        (1, "Tidak Bekerja", "Domestik"),
        (2, "Ibu Rumah Tangga", "Domestik"),
        (3, "Petani", "Pertanian"),
        (4, "Buruh Harian Lepas", "Informal"),
        (5, "Pedagang Kecil", "Informal"),
        (6, "Karyawan Swasta", "Formal"),
        (7, "PNS/TNI/Polri", "Pemerintah"),
        (8, "Guru/Dosen", "Pendidikan"),
        (9, "Perawat/Bidan", "Kesehatan"),
        (10, "Dokter", "Kesehatan"),
        (11, "Sopir", "Transportasi"),
        (12, "Nelayan", "Pertanian"),
        (13, "Pengrajin", "Informal"),
        (14, "Wiraswasta", "Formal"),
        (15, "Peternak", "Pertanian"),
    ]
    return [{"id_pekerjaan": d[0], "nama_pekerjaan": d[1], "sektor": d[2]} for d in data]

# ─── 4. KATEGORI PENDAPATAN ────────────────────────────────────
def gen_kategori_pendapatan():
    data = [
        (1, "< Rp 500.000", 0, 500000),
        (2, "Rp 500.000 - Rp 1.000.000", 500000, 1000000),
        (3, "Rp 1.000.000 - Rp 2.000.000", 1000000, 2000000),
        (4, "Rp 2.000.000 - Rp 3.000.000", 2000000, 3000000),
        (5, "Rp 3.000.000 - Rp 5.000.000", 3000000, 5000000),
        (6, "> Rp 5.000.000", 5000000, 50000000),
    ]
    return [{"id_pendapatan": d[0], "kategori_pendapatan": d[1], "pendapatan_min": d[2], "pendapatan_max": d[3]} for d in data]

# ─── 5. USER ACCOUNT ──────────────────────────────────────────
def gen_user_account(lokasi_ids):
    rows = []
    nik_set = set()
    email_set = set()
    used_names = set()
    now = datetime(2025, 6, 1)

    def make_user(jk, role_tag, birth_year_min, birth_year_max):
        nonlocal rows
        while True:
            nama = random_nama(jk)
            if nama not in used_names:
                used_names.add(nama)
                break
        idx = len(rows) + 1
        email = random_email(nama, idx)
        while email in email_set:
            email = random_email(nama + str(random.randint(1, 999)), idx)
        email_set.add(email)
        nik = random_nik()
        while nik in nik_set:
            nik = random_nik()
        nik_set.add(nik)
        tgl = date(random.randint(birth_year_min, birth_year_max), random.randint(1, 12), random.randint(1, 28))
        user = {
            "id_user": idx,
            "email": email,
            "password": fast_hash("password123"),
            "no_hp": random_no_hp()[:13],
            "status_verifikasi": "Aktif",
            "nama": nama,
            "nik": nik,
            "jenis_kelamin": jk,
            "tanggal_lahir": tgl.isoformat(),
            "id_lokasi": random.choice(lokasi_ids),
            "id_pendidikan": random.randint(1, 8),
            "id_pekerjaan": random.randint(1, 15),
            "id_pendapatan": random.randint(1, 6),
            "jumlah_tanggungan": random.randint(0, 5),
            "akun_ke": 1,
            "created_at": now.isoformat(),
            "updated_at": now.isoformat(),
        }
        rows.append(user)
        return user["id_user"], user["nama"], user["nik"], user["jenis_kelamin"]

    # Dinkes (5)
    dinkes_ids = []
    for _ in range(TOTAL_DINKES):
        uid, nama, nik, jk = make_user(random.choice(["Laki-Laki", "Perempuan"]), "Dinkes", 1970, 1990)
        dinkes_ids.append(uid)

    # Bidan (100)
    bidan_ids = []
    for _ in range(TOTAL_BIDAN):
        uid, nama, nik, jk = make_user("Perempuan", "Bidan", 1975, 1995)
        bidan_ids.append(uid)

    # Kader (200)
    kader_ids = []
    for _ in range(TOTAL_KADER):
        uid, nama, nik, jk = make_user("Perempuan", "Kader", 1980, 2000)
        kader_ids.append(uid)

    # Ibu hamil (1000)
    ibu_ids = []
    ibu_data = []
    for _ in range(TOTAL_IBU_HAMIL):
        uid, nama, nik, jk = make_user("Perempuan", "Ibu Hamil", 1990, 2008)
        ibu_ids.append(uid)
        ibu_data.append({"nama": nama, "nik": nik})

    # Anak (2000) - anak laki dan perempuan
    anak_ids = []
    anak_wali = []
    for _ in range(TOTAL_ANAK):
        uid, nama, nik, jk = make_user(random.choice(["Laki-Laki", "Perempuan"]), "Anak", 2020, 2025)
        anak_ids.append(uid)
        anak_wali.append({"nama": nama, "nik": nik})

    # Pasien dewasa non-role (500)
    pasien_dewasa_ids = []
    for _ in range(TOTAL_PASIEN_DEWASA):
        uid, nama, nik, jk = make_user(random.choice(["Laki-Laki", "Perempuan"]), "Pasien", 1985, 2005)
        pasien_dewasa_ids.append(uid)

    total_gen = len(rows)
    if total_gen < TOTAL_USER:
        for _ in range(TOTAL_USER - total_gen):
            make_user(random.choice(["Laki-Laki", "Perempuan"]), "Pasien", 1980, 2010)

    return rows, dinkes_ids, bidan_ids, kader_ids, ibu_ids, anak_ids, pasien_dewasa_ids

# ─── 6-8. ROLE TABLES ──────────────────────────────────────────
def gen_dinkes(dinkes_ids):
    return [{"id_user": uid, "created_at": datetime(2025, 6, 1).isoformat(), "updated_at": datetime(2025, 6, 1).isoformat()} for uid in dinkes_ids]

def gen_bidan(bidan_ids, lokasi_ids):
    rows = []
    for i, uid in enumerate(bidan_ids):
        rows.append({
            "id_user": uid,
            "no_str": f"STR-2024-{uid:04d}",
            "wilayah_kerja": random.choice(lokasi_ids),
            "created_at": datetime(2025, 6, 1).isoformat(),
            "updated_at": datetime(2025, 6, 1).isoformat(),
        })
    return rows

def gen_kader(kader_ids, posyandu_ids):
    rows = []
    for i, uid in enumerate(kader_ids):
        rows.append({
            "id_user": uid,
            "no_sk": f"SK-2024-{uid:04d}",
            "id_posyandu": random.choice(posyandu_ids),
            "created_at": datetime(2025, 6, 1).isoformat(),
            "updated_at": datetime(2025, 6, 1).isoformat(),
        })
    return rows

# ─── 7. FASILITAS KESEHATAN ────────────────────────────────────
def gen_faskes(lokasi_ids):
    tipe = ["Faskes Tingkat Pertama", "Faskes Rujukan Tingkat Lanjutan", "Faskes Penunjang"]
    rows = []
    for i in range(1, TOTAL_FASKES + 1):
        rows.append({
            "id_faskes": i,
            "nama_faskes": f"{random.choice(['Puskesmas','RSUD','RS Swasta','Klinik','RS Ibu dan Anak'])} {random.choice(['Sehat','Harapan','Mulya','Asih','Bersinar','Cendana','Bhakti'])}",
            "tipe_faskes": random.choice(tipe),
            "id_lokasi": random.choice(lokasi_ids),
            "created_at": datetime(2025, 6, 1).isoformat(),
            "updated_at": datetime(2025, 6, 1).isoformat(),
        })
    return rows

# ─── 8. POSYANDU ───────────────────────────────────────────────
def gen_posyandu(lokasi_ids, bidan_ids):
    rows = []
    for i in range(1, TOTAL_POSYANDU + 1):
        rows.append({
            "id_posyandu": i,
            "nama_posyandu": f"Posyandu {random.choice(['Melati','Mawar','Anggrek','Kenanga','Flamboyan','Cempaka','Dahlia','Bougenville','Kamboja','Teratai'])} {i}",
            "id_lokasi": random.choice(lokasi_ids),
            "id_bidan": random.choice(bidan_ids),
            "created_at": datetime(2025, 6, 1).isoformat(),
            "updated_at": datetime(2025, 6, 1).isoformat(),
        })
    return rows

# ─── 9. PASIEN ─────────────────────────────────────────────────
def gen_pasien(ibu_ids, anak_ids, pasien_dewasa_ids, posyandu_ids):
    all_pasien = ibu_ids + anak_ids + pasien_dewasa_ids
    rows = []
    for pid in all_pasien:
        rows.append({
            "id_pasien": pid,
            "id_posyandu": random.choice(posyandu_ids),
            "created_at": datetime(2025, 6, 1).isoformat(),
            "updated_at": datetime(2025, 6, 1).isoformat(),
        })
    return rows

# ─── 10. IBU HAMIL ─────────────────────────────────────────────
def gen_ibu_hamil(ibu_ids):
    rows = []
    now = datetime(2025, 6, 1)
    for i, pid in enumerate(ibu_ids):
        hamil_ke = random.randint(1, 4)
        hpht = now - timedelta(days=random.randint(30, 270))
        bulan_mulai = hpht
        status = random.choice(["Trimester 1", "Trimester 2", "Trimester 3", "Melahirkan"])
        rows.append({
            "id_ibu_hamil": i + 1,
            "id_pasien": pid,
            "hamil_ke": hamil_ke,
            "bulan_mulai_hamil": bulan_mulai.date().isoformat(),
            "hpht": hpht.date().isoformat(),
            "status_kehamilan": status,
            "created_at": now.isoformat(),
            "updated_at": now.isoformat(),
        })
    return rows

# ─── 11. ANAK ──────────────────────────────────────────────────
def gen_anak(anak_ids, ibu_hamil_rows, pasien_rows):
    """Setiap anak punya ibu_hamil dan wali (laki dewasa dari pasien_dewasa_ids)"""
    rows = []
    wali_pool = [r["id_pasien"] for r in pasien_rows if r["id_pasien"] not in anak_ids]
    # ambil yang laki-laki sebagai wali (gunakan sisa user yang bukan ibu hamil/anak)
    for i, a_id in enumerate(anak_ids):
        ibu = random.choice(ibu_hamil_rows)
        wali = random.choice(wali_pool) if wali_pool else 1
        bb = round(random.uniform(2.5, 4.2), 2)
        pb = round(random.uniform(45.0, 55.0), 2)
        rows.append({
            "id_pasien": a_id,
            "id_ibu_hamil": ibu["id_ibu_hamil"],
            "id_wali": wali,
            "berat_lahir": bb,
            "panjang_lahir": pb,
            "hubungan_dengan_wali": random.choice(["Kandung", "Tiri", "Angkat"]),
            "created_at": datetime(2025, 6, 1).isoformat(),
            "updated_at": datetime(2025, 6, 1).isoformat(),
        })
    return rows

# ─── 12. JADWAL IMUNISASI ──────────────────────────────────────
def gen_jadwal_imunisasi(anak_ids):
    vaksin_list = [
        ("Hepatitis B", 0),
        ("BCG", 1),
        ("Polio 1", 1),
        ("DPT-HB-Hib 1", 2),
        ("Polio 2", 2),
        ("DPT-HB-Hib 2", 4),
        ("Polio 3", 4),
        ("DPT-HB-Hib 3", 6),
        ("Campak/MR", 9),
        ("Polio 4", 18),
    ]
    rows = []
    idx = 1
    now = date(2025, 6, 1)
    for a_id in anak_ids:
        tgl_lahir_anak = now - timedelta(days=random.randint(30, 1825))  # 0-5 th
        for vaksin, usia_bln in vaksin_list:
            jadwal = tgl_lahir_anak + timedelta(days=usia_bln * 30)
            if jadwal > now:
                # masa depan -> belum
                status = "Belum"
                real = None
            else:
                # sudah lewat -> 80% sudah, 20% belum
                if random.random() < 0.8:
                    status = "Sudah"
                    real = jadwal + timedelta(days=random.randint(0, 14))
                else:
                    status = "Belum"
                    real = None
            rows.append({
                "id_imunisasi": idx,
                "id_pasien": a_id,
                "nama_vaksin": vaksin,
                "tanggal_jadwal": jadwal.isoformat(),
                "tanggal_realisasi": real.isoformat() if real else None,
                "status_imunisasi": status,
                "created_at": datetime(2025, 6, 1).isoformat(),
                "updated_at": datetime(2025, 6, 1).isoformat(),
            })
            idx += 1
    return rows

# ─── 13. HASIL PEMERIKSAAN ─────────────────────────────────────
def gen_hasil_pemeriksaan(jadwal_rows, bidan_ids):
    rows = []
    status_stunting_opts = ["Normal", "Berisiko Stunting", "Stunting", "Stunting Berat"]
    status_gizi_opts = ["Gizi Baik", "Gizi Kurang", "Gizi Buruk", "Risiko Gizi Lebih", "Gizi Lebih", "Obesitas"]

    for jr in jadwal_rows:
        if jr["status_imunisasi"] != "Sudah":
            continue
        # usia anak dalam bulan saat pemeriksaan
        # estimasi
        bb = round(random.uniform(5.0, 25.0), 2)
        tb = round(random.uniform(50.0, 120.0), 2)
        lk = round(random.uniform(30.0, 52.0), 2)
        td = f"{random.randint(90, 130)}/{random.randint(60, 90)}"

        # korelasi: bb rendah -> gizi buruk -> stunting
        if bb < 8.0:
            sg = "Gizi Buruk"
            ss = "Stunting Berat"
        elif bb < 10.0:
            sg = "Gizi Kurang"
            ss = random.choice(["Berisiko Stunting", "Stunting"])
        elif bb > 20.0:
            sg = random.choice(["Gizi Lebih", "Obesitas"])
            ss = "Normal"
        else:
            sg = "Gizi Baik"
            ss = "Normal"

        rows.append({
            "id_hasil_pemeriksaan": len(rows) + 1,
            "id_petugas_input": random.choice(bidan_ids),
            "id_jadwal_imunisasi": jr["id_imunisasi"],
            "berat_badan": bb,
            "tinggi_badan": tb,
            "lingkar_kepala": lk,
            "tekanan_darah": td,
            "status_stunting": ss,
            "status_gizi": sg,
            "catatan": random.choice([
                "Anak dalam kondisi sehat",
                "Perlu peningkatan asupan gizi",
                "Imunisasi lengkap sesuai jadwal",
                "Berat badan kurang, perlu evaluasi",
                "Perkembangan motorik baik",
                None,
                None,
            ]),
            "created_at": datetime(2025, 6, 1).isoformat(),
            "updated_at": datetime(2025, 6, 1).isoformat(),
        })
    return rows

# ─── 14. TINDAK LANJUT ─────────────────────────────────────────
def gen_tindak_lanjut(hasil_rows, bidan_ids):
    rows = []
    catatan_list = [
        "Asupan gizi ditingkatkan dengan makanan bergizi seimbang",
        "Pemberian ASI eksklusif dilanjutkan",
        "Kontrol berat badan secara rutin setiap bulan",
        "Pemberian vitamin A sesuai jadwal",
        "Imunisasi lanjutan sesuai jadwal yang ditentukan",
        "Pantau tumbuh kembang anak secara berkala",
        "Konseling gizi kepada orang tua",
    ]
    rekom_list = [
        "Datang kembali ke posyandu bulan depan",
        "Segera ke puskesmas jika ada keluhan",
        "Rujuk ke dokter spesialis anak",
        "Perbanyak makanan berprotein tinggi",
        "Pantau tinggi dan berat badan setiap minggu",
        "Istirahat cukup dan aktivitas fisik teratur",
    ]
    for hr in hasil_rows:
        bidan = random.choice(bidan_ids)
        jadwal_kontrol = date(2025, 6, 1) + timedelta(days=random.randint(7, 90))
        status = random.choice(["Dalam Pemantauan", "Membaik", "Perlu Rujukan", "Selesai Pemantauan"])
        rows.append({
            "id_tindak_lanjut": len(rows) + 1,
            "id_hasil_pemeriksaan": hr["id_hasil_pemeriksaan"],
            "id_bidan": bidan,
            "catatan_medis": random.choice(catatan_list),
            "rekomendasi": random.choice(rekom_list),
            "jadwal_kontrol": jadwal_kontrol.isoformat(),
            "status_pasien": status,
            "created_at": datetime(2025, 6, 1).isoformat(),
            "updated_at": datetime(2025, 6, 1).isoformat(),
        })
    return rows

# ─── 15. RUJUKAN ───────────────────────────────────────────────
def gen_rujukan(tindak_lanjut_rows, faskes_rows):
    rows = []
    faskes_ids = [r["id_faskes"] for r in faskes_rows]
    alasan_list = [
        "Gizi buruk memerlukan penanganan intensif",
        "Stunting berat membutuhkan intervensi medis lanjutan",
        "Berat badan tidak kunjung naik setelah intervensi",
        "Terdapat kelainan kongenital yang perlu evaluasi",
        "Kebutuhan pemeriksaan laboratorium lanjutan",
        "Riwayat kehamilan risiko tinggi",
    ]
    for tl in tindak_lanjut_rows:
        if tl["status_pasien"] != "Perlu Rujukan":
            continue
        rows.append({
            "id_rujukan": len(rows) + 1,
            "id_tindak_lanjut": tl["id_tindak_lanjut"],
            "alasan_rujukan": random.choice(alasan_list),
            "tanggal_rujukan": date(2025, 6, 1).isoformat(),
            "status_rujukan": random.choice(["Diajukan", "Diproses", "Diterima", "Selesai"]),
            "id_faskes": random.choice(faskes_ids),
            "created_at": datetime(2025, 6, 1).isoformat(),
            "updated_at": datetime(2025, 6, 1).isoformat(),
        })
    return rows

# ─── 16. NOTIFIKASI ────────────────────────────────────────────
def gen_notifikasi(all_user_ids):
    rows = []
    tipe_list = ["Pemeriksaan", "Imunisasi", "Rujukan", "Edukasi", "Pengingat"]
    judul_list = [
        ("Jadwal Imunisasi", "Anak Anda memiliki jadwal imunisasi dalam waktu dekat. Silakan datang ke posyandu terdekat."),
        ("Hasil Pemeriksaan", "Hasil pemeriksaan anak Anda sudah tersedia. Silakan cek di aplikasi."),
        ("Pengingat Kontrol", "Ingat! Jadwal kontrol anak Anda sudah dekat. Jangan lupa datang ke posyandu."),
        ("Edukasi Gizi", "Pastikan anak Anda mendapatkan asupan gizi yang cukup setiap hari."),
        ("Status Rujukan", "Status rujukan anak Anda telah diperbarui. Silakan cek detailnya."),
        ("Imunisasi Lengkap", "Selamat! Imunisasi dasar anak Anda sudah lengkap. Tetap jaga kesehatannya."),
        ("Pemeriksaan Bulanan", "Waktunya pemeriksaan bulanan untuk memantau tumbuh kembang anak."),
    ]
    for i in range(1, 10001):
        judul, pesan = random.choice(judul_list)
        tgl = datetime(2025, 6, 1) - timedelta(days=random.randint(0, 180))
        rows.append({
            "id_notifikasi": i,
            "id_user": random.choice(all_user_ids),
            "judul": judul,
            "pesan": pesan,
            "tipe_notifikasi": random.choice(tipe_list),
            "status_baca": random.choice([True, False]),
            "tanggal_kirim": tgl.isoformat(),
            "is_deleted": False,
            "deleted_at": None,
        })
    return rows

# ─── 17. USER SESSION ──────────────────────────────────────────
def gen_user_session(all_user_ids):
    rows = []
    status_opts = ["AKTIF", "KADALUWARSA", "TERGANTI", "DICABUT"]
    for i in range(1, 20001):
        uid = random.choice(all_user_ids)
        created = datetime(2025, 6, 1) - timedelta(days=random.randint(0, 30))
        expired = created + timedelta(days=7)
        rows.append({
            "id_session": str(uuid.uuid4()),
            "id_user": uid,
            "status_session": random.choice(status_opts),
            "ip_address": f"192.168.{random.randint(0, 255)}.{random.randint(1, 254)}",
            "created_at": created.isoformat(),
            "updated_at": created.isoformat(),
            "expired_at": expired.isoformat(),
        })
    return rows

# ─── 18. AUDIT LOG ─────────────────────────────────────────────
def gen_audit_log(all_user_ids):
    rows = []
    aktivitas_list = ["LOGIN", "REGISTRASI", "DATA_INSERT", "DATA_UPDATE", "DATA_DELETE"]
    detail_list = [
        "Pengguna berhasil login ke sistem",
        "Registrasi pengguna baru berhasil",
        "Data pemeriksaan baru berhasil ditambahkan",
        "Data imunisasi berhasil diperbarui",
        "Data pasien berhasil dihapus",
        "Verifikasi akun pengguna berhasil",
        "Data rujukan berhasil dibuat",
    ]
    for i in range(1, 10001):
        uid = random.choice(all_user_ids)
        rows.append({
            "id_log": i,
            "tipe_aktor": "USER",
            "id_user": uid,
            "id_user_session": None,
            "tipe_aktivitas": random.choice(aktivitas_list),
            "berhasil": random.choice([True, True, True, False]),
            "endpoint": f"/v1/{random.choice(['auth/login','auth/register','imunisasi','pemeriksaan','pasien','rujukan'])}",
            "table_name": None,
            "record_id": None,
            "old_value": None,
            "new_value": None,
            "detail": random.choice(detail_list),
            "ip_address": f"192.168.{random.randint(0, 255)}.{random.randint(1, 254)}",
            "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
            "waktu_aktivitas": (datetime(2025, 6, 1) - timedelta(days=random.randint(0, 180))).isoformat(),
        })
    return rows

# ─── 19. ARTIKEL ───────────────────────────────────────────────
def gen_artikel(dinkes_ids):
    articles = [
        ("Pentingnya Imunisasi Dasar Lengkap untuk Bayi",
         "Imunisasi merupakan salah satu upaya pencegahan penyakit yang paling efektif dan efisien. Dengan memberikan imunisasi dasar lengkap, bayi akan terlindungi dari berbagai penyakit berbahaya seperti tuberkulosis, hepatitis B, polio, difteri, pertusis, tetanus, dan campak. Pemerintah Indonesia telah mewajibkan pemberian imunisasi dasar lengkap sebelum anak berusia satu tahun. Orang tua diharapkan aktif membawa bayinya ke posyandu atau puskesmas terdekat untuk mendapatkan imunisasi sesuai jadwal yang telah ditentukan. Imunisasi tidak hanya melindungi anak secara individu, tetapi juga menciptakan kekebalan kelompok yang melindungi masyarakat secara keseluruhan."),
        ("Mengenal Stunting dan Cara Pencegahannya",
         "Stunting adalah kondisi gagal tumbuh pada anak akibat kekurangan gizi kronis dalam waktu lama, terutama pada 1.000 hari pertama kehidupan. Anak stunting memiliki tinggi badan yang lebih pendek dibandingkan standar usianya. Dampak stunting tidak hanya pada fisik, tetapi juga pada perkembangan otak dan kecerdasan anak. Pencegahan stunting dimulai dari masa kehamilan dengan memastikan ibu hamil mendapatkan asupan gizi yang cukup. Setelah lahir, bayi harus mendapatkan ASI eksklusif selama enam bulan, dilanjutkan dengan MPASI yang bergizi. Pemantauan tumbuh kembang secara rutin di posyandu juga sangat penting untuk mendeteksi dini tanda-tanda stunting."),
        ("Manfaat ASI Eksklusif bagi Tumbuh Kembang Bayi",
         "Air Susu Ibu (ASI) adalah makanan terbaik untuk bayi karena mengandung semua nutrisi yang dibutuhkan dalam jumlah yang tepat. ASI eksklusif diberikan selama enam bulan pertama kehidupan tanpa tambahan makanan atau minuman lain. ASI mengandung antibodi yang melindungi bayi dari infeksi, serta berbagai zat gizi yang mendukung perkembangan otak dan sistem saraf. Ibu yang memberikan ASI eksklusif juga mendapatkan manfaat seperti mempererat ikatan batin dengan bayi dan membantu pemulihan pasca melahirkan. Dukungan keluarga dan lingkungan kerja sangat penting untuk keberhasilan pemberian ASI eksklusif."),
        ("Gizi Seimbang untuk Ibu Hamil",
         "Masa kehamilan adalah periode kritis yang membutuhkan perhatian khusus terhadap asupan gizi. Ibu hamil memerlukan tambahan energi, protein, vitamin, dan mineral untuk mendukung pertumbuhan janin. Asam folat sangat penting pada trimester pertama untuk mencegah cacat tabung saraf pada janin. Zat besi diperlukan untuk mencegah anemia yang sering terjadi pada ibu hamil. Kalsium dan vitamin D mendukung pembentukan tulang dan gigi janin. Ibu hamil disarankan mengonsumsi makanan bergizi seimbang dengan porsi lebih banyak dari biasanya serta minum suplemen sesuai anjuran dokter."),
        ("Perkembangan Motorik Anak Usia 0-2 Tahun",
         "Perkembangan motorik anak terbagi menjadi motorik kasar dan motorik halus. Motorik kasar meliputi gerakan tubuh besar seperti tengkurap, duduk, merangkak, dan berjalan. Motorik halus meliputi gerakan tangan dan jari seperti memegang benda, meraih mainan, dan memasukkan benda ke dalam wadah. Setiap anak memiliki kecepatan perkembangan yang berbeda, namun ada tonggak perkembangan yang dapat dijadikan acuan. Orang tua perlu memberikan stimulasi yang sesuai dengan usia anak untuk mendukung perkembangan motoriknya. Jika ada keterlambatan perkembangan, segera konsultasikan dengan dokter atau bidan."),
        ("Jadwal Imunisasi Anak Usia 0-18 Bulan",
         "Jadwal imunisasi anak di Indonesia mengacu pada rekomendasi IDAI dan Kementerian Kesehatan. Imunisasi Hepatitis B diberikan segera setelah lahir, diikuti BCG dan Polio 1 pada usia 1 bulan. DPT-HB-Hib 1 dan Polio 2 diberikan pada usia 2 bulan, dilanjutkan dengan DPT-HB-Hib 2 dan Polio 3 pada usia 4 bulan. Pada usia 6 bulan, anak mendapatkan DPT-HB-Hib 3 dan Polio 4. Imunisasi Campak/MR diberikan pada usia 9 bulan. Orang tua disarankan untuk mencatat setiap imunisasi yang telah diterima anak dan memastikan tidak ada yang terlewat."),
        ("Peran Posyandu dalam Pemantauan Tumbuh Kembang Anak",
         "Posyandu adalah pusat kegiatan masyarakat di bidang kesehatan yang dikelola oleh kader kesehatan. Posyandu memiliki peran penting dalam memantau tumbuh kembang anak melalui penimbangan berat badan, pengukuran tinggi badan, dan lingkar kepala. Hasil pengukuran dicatat dalam Kartu Menuju Sehat (KMS) untuk memantau status gizi anak secara berkala. Posyandu juga memberikan pelayanan imunisasi, penyuluhan gizi, dan vitamin A. Kader posyandu yang terlatih dapat mendeteksi dini jika ada masalah pertumbuhan pada anak dan merujuk ke puskesmas untuk penanganan lebih lanjut."),
        ("Makanan Pendamping ASI (MPASI) yang Tepat untuk Bayi",
         "Setelah usia 6 bulan, bayi memerlukan makanan pendamping ASI (MPASI) untuk memenuhi kebutuhan gizinya yang semakin meningkat. MPASI pertama harus memiliki tekstur yang lembut dan mudah dicerna, seperti bubur susu atau pure buah. Secara bertahap, tekstur MPASI dapat ditingkatkan menjadi lebih kental dan kasar sesuai dengan kemampuan mengunyah bayi. MPASI harus mengandung karbohidrat, protein hewani dan nabati, lemak, vitamin, dan mineral. Hindari menambahkan gula dan garam berlebihan pada MPASI bayi. Perkenalkan satu jenis makanan baru dalam satu waktu untuk mendeteksi kemungkinan alergi."),
        ("Bahaya Gizi Buruk pada Anak dan Penanganannya",
         "Gizi buruk adalah kondisi medis serius yang terjadi ketika anak tidak mendapatkan asupan nutrisi yang cukup. Gejala gizi buruk meliputi berat badan sangat kurang dari standar, pertumbuhan terhambat, mudah sakit, lesu, dan pembengkakan pada kaki (edema). Gizi buruk dapat menyebabkan kerusakan organ, gangguan perkembangan otak, bahkan kematian jika tidak ditangani dengan tepat. Penanganan gizi buruk memerlukan intervensi medis intensif seperti pemberian makanan terapeutik, suplemen vitamin dan mineral, serta pengobatan komplikasi. Pencegahan tetap menjadi langkah terbaik melalui pemberian ASI eksklusif dan MPASI bergizi."),
        ("Tips Menjaga Kesehatan Ibu Nifas",
         "Masa nifas adalah periode pemulihan setelah melahirkan yang berlangsung sekitar 6 minggu. Selama masa nifas, ibu membutuhkan istirahat yang cukup, asupan makanan bergizi, dan dukungan emosional dari keluarga. Ibu nifas dianjurkan mengonsumsi makanan kaya zat besi untuk menggantikan darah yang hilang saat melahirkan. Menjaga kebersihan area kewanitaan sangat penting untuk mencegah infeksi. Ibu menyusui perlu minum air putih yang cukup untuk mempertahankan produksi ASI. Jangan ragu untuk berkonsultasi dengan bidan atau dokter jika ada keluhan selama masa nifas."),
        ("Pentingnya Vitamin A bagi Kesehatan Anak",
         "Vitamin A memiliki peran penting dalam menjaga kesehatan mata, sistem kekebalan tubuh, dan pertumbuhan sel. Kekurangan vitamin A dapat menyebabkan gangguan penglihatan, terutama rabun senja, serta meningkatkan kerentanan terhadap infeksi. Pemerintah Indonesia melalui Kementerian Kesehatan mengadakan program pemberian kapsul vitamin A setiap bulan Februari dan Agustus. Bayi usia 6-11 bulan mendapatkan kapsul vitamin A dosis 100.000 IU (biru). Anak usia 12-59 bulan mendapatkan kapsul vitamin A dosis 200.000 IU (merah)."),
        ("Cara Mendeteksi Dini Gangguan Tumbuh Kembang Anak",
         "Deteksi dini gangguan tumbuh kembang anak sangat penting untuk intervensi yang tepat waktu. Orang tua dapat memantau perkembangan anak menggunakan buku KIA atau KMS. Beberapa tanda yang perlu diwaspadai antara lain: berat badan tidak naik dalam 2 bulan berturut-turut, anak tidak bisa tengkurap pada usia 6 bulan, tidak bisa duduk tanpa bantuan pada usia 9 bulan, belum bisa berjalan pada usia 18 bulan, dan kesulitan berbicara. Jika menemukan tanda-tanda tersebut, segera bawa anak ke posyandu, puskesmas, atau dokter spesialis anak untuk evaluasi lebih lanjut."),
        ("Manfaat Asi Eksklusif dan Dampaknya bagi Kesehatan Bayi",
         "Pemberian ASI eksklusif selama 6 bulan memberikan perlindungan optimal bagi bayi dari berbagai penyakit infeksi seperti diare, pneumonia, dan infeksi telinga. ASI mengandung kolostrum yang kaya antibodi pada hari-hari pertama setelah melahirkan. Bayi yang mendapat ASI eksklusif cenderung memiliki berat badan yang ideal dan risiko obesitas yang lebih rendah saat dewasa. Selain itu, ASI eksklusif juga mendukung perkembangan sistem pencernaan dan otak bayi. Ibu menyusui perlu menjaga asupan makanannya agar kualitas ASI tetap optimal."),
        ("Program Penanggulangan Stunting di Indonesia",
         "Pemerintah Indonesia telah menetapkan program penanggulangan stunting sebagai prioritas nasional. Program ini melibatkan berbagai sektor termasuk kesehatan, pendidikan, pertanian, dan sosial. Intervensi spesifik dilakukan pada 1.000 hari pertama kehidupan, meliputi pemberian makanan tambahan untuk ibu hamil, promosi ASI eksklusif, MPASI bergizi, dan tatalaksana gizi buruk. Intervensi sensitif meliputi akses air bersih, sanitasi, dan jaminan kesehatan. Peran serta masyarakat melalui posyandu menjadi kunci keberhasilan program ini. Target penurunan stunting menjadi 14% pada tahun 2024 menjadi komitmen bersama."),
        ("Olahraga Ringan untuk Ibu Hamil",
         "Ibu hamil dianjurkan untuk tetap aktif berolahraga selama kehamilan dengan intensitas ringan hingga sedang. Olahraga selama kehamilan membantu menjaga kebugaran, mengontrol berat badan, mengurangi nyeri punggung, dan mempersiapkan proses persalinan. Jenis olahraga yang aman untuk ibu hamil antara lain jalan kaki, berenang, yoga prenatal, dan senam hamil. Namun, ibu hamil dengan kondisi tertentu seperti risiko keguguran atau plasenta previa harus berkonsultasi dengan dokter terlebih dahulu. Hindari olahraga yang membahayakan seperti lompat-lompat, kontak fisik, atau olahraga di tempat terlalu panas."),
        ("Imunisasi Polio dan Upaya Eradikasi Global",
         "Polio adalah penyakit menular yang disebabkan oleh virus polio dan dapat menyebabkan kelumpuhan permanen. Imunisasi polio diberikan sebanyak 4 kali sebelum anak berusia 1 tahun, yaitu pada usia 1, 2, 4, dan 18 bulan. Indonesia telah berkomitmen untuk mendukung program eradikasi polio global. Melalui Pekan Imunisasi Nasional (PIN), pemerintah berupaya mencapai cakupan imunisasi polio yang tinggi di seluruh wilayah. Meskipun Indonesia telah dinyatakan bebas polio, kewaspadaan tetap diperlukan mengingat masih adanya kasus polio di negara tetangga."),
        ("Pentingnya Buku KIA bagi Ibu dan Anak",
         "Buku Kesehatan Ibu dan Anak (KIA) adalah dokumen penting yang mencatat seluruh riwayat kesehatan ibu selama kehamilan, persalinan, dan nifas, serta tumbuh kembang anak hingga usia 5 tahun. Buku KIA berisi jadwal imunisasi, grafik pertumbuhan, tabel perkembangan anak, serta informasi kesehatan yang bermanfaat. Ibu disarankan untuk selalu membawa buku KIA setiap kali memeriksakan kehamilan atau membawa anak ke posyandu. Buku KIA juga menjadi alat komunikasi antara ibu dan tenaga kesehatan dalam memantau kesehatan ibu dan anak secara berkelanjutan."),
        ("Kenali Tanda Bahaya Kehamilan yang Perlu Diwaspadai",
         "Setiap ibu hamil perlu mengetahui tanda-tanda bahaya kehamilan agar dapat segera mendapatkan pertolongan medis. Tanda bahaya kehamilan meliputi perdarahan pervaginam, sakit kepala hebat, penglihatan kabur, demam tinggi, bengkak pada wajah dan tangan, gerakan janin berkurang, dan nyeri perut hebat. Pada trimester ketiga, waspadai pecah ketuban sebelum waktunya dan kontraksi prematur. Jika mengalami salah satu tanda tersebut, segera bawa ibu ke fasilitas kesehatan terdekat untuk mendapatkan penanganan. Deteksi dan penanganan dini dapat mencegah komplikasi serius pada ibu dan janin."),
        ("Kalender Imunisasi Anak Usia Sekolah",
         "Anak usia sekolah juga memerlukan imunisasi lanjutan untuk mempertahankan kekebalan tubuhnya. Imunisasi yang diberikan pada anak usia sekolah meliputi imunisasi DT pada kelas 1 SD, Campak/MR pada kelas 1 SD, dan Td pada kelas 2 dan 5 SD. Program Bulan Imunisasi Anak Sekolah (BIAS) dilaksanakan setiap tahun oleh puskesmas bekerja sama dengan sekolah. Orang tua diharapkan mendukung program ini dengan memberikan izin kepada anak untuk diimunisasi di sekolah. Imunisasi lanjutan penting untuk mencegah wabah penyakit yang dapat mengganggu proses belajar anak."),
        ("Peran Ayah dalam Mendukung Pemberian ASI Eksklusif",
         "Dukungan ayah sangat berperan dalam keberhasilan pemberian ASI eksklusif. Ayah dapat membantu dengan cara mengambil alih tugas rumah tangga, memberikan dukungan emosional, mengingatkan istri untuk makan makanan bergizi, dan menemani saat memeriksakan bayi ke posyandu. Kehadiran ayah yang suportif dapat mengurangi stres ibu menyusui yang berdampak positif pada produksi ASI. Ayah juga dapat belajar tentang teknik menyusui dan cara menyimpan ASI perah agar dapat membantu saat istri beristirahat. Dukungan suami menjadi salah satu faktor kunci keberhasilan ASI eksklusif."),
    ]
    rows = []
    for i, (judul, isi) in enumerate(articles, 1):
        rows.append({
            "id_artikel": i,
            "judul": judul,
            "isi_artikel": isi,
            "kategori": random.choice(["Imunisasi", "Stunting", "Gizi", "Tumbuh Kembang", "Kehamilan"]),
            "status_artikel": random.choice(["Dipublikasikan", "Dipublikasikan", "Dipublikasikan", "Draft"]),
            "id_penulis": random.choice(dinkes_ids),
            "id_verifikator": random.choice(dinkes_ids),
            "tanggal_publish": (date(2025, 6, 1) - timedelta(days=random.randint(0, 180))).isoformat(),
            "created_at": datetime(2025, 6, 1).isoformat(),
            "updated_at": datetime(2025, 6, 1).isoformat(),
        })
    return rows

# ─── MAIN ──────────────────────────────────────────────────────
def main():
    print("=== GENERATE DATA DUMMY IMUNISASI ===")

    print("\n[1/19] Lokasi...")
    lokasi = gen_lokasi()
    lokasi_ids = [r["id"] for r in lokasi]
    csv_writer("lokasi.csv",
        ["id_lokasi","nama_lokasi","tipe_lokasi","bagian_dari"],
        [{"id_lokasi": r["id"], "nama_lokasi": r["nama"], "tipe_lokasi": r["tipe"], "bagian_dari": r["bagian_dari"]} for r in lokasi])

    print("\n[2/19] Pendidikan...")
    csv_writer("pendidikan.csv",
        ["id_pendidikan","nama_pendidikan","jenjang","lama_tahun"],
        gen_pendidikan())

    print("\n[3/19] Pekerjaan...")
    csv_writer("pekerjaan.csv",
        ["id_pekerjaan","nama_pekerjaan","sektor"],
        gen_pekerjaan())

    print("\n[4/19] Kategori Pendapatan...")
    csv_writer("kategori_pendapatan.csv",
        ["id_pendapatan","kategori_pendapatan","pendapatan_min","pendapatan_max"],
        gen_kategori_pendapatan())

    print("\n[5/19] User Account...")
    users, dinkes_ids, bidan_ids, kader_ids, ibu_ids, anak_ids, pasien_dewasa_ids = gen_user_account(lokasi_ids)
    all_user_ids = [u["id_user"] for u in users]
    csv_writer("user_account.csv",
        ["id_user","email","password","no_hp","status_verifikasi","nama","nik","jenis_kelamin","tanggal_lahir","id_lokasi","id_pendidikan","id_pekerjaan","id_pendapatan","jumlah_tanggungan","akun_ke","created_at","updated_at"],
        users)

    print("\n[6/19] Dinas Kesehatan...")
    csv_writer("dinas_kesehatan.csv",
        ["id_user","created_at","updated_at"],
        gen_dinkes(dinkes_ids))

    print("\n[7/19] Fasilitas Kesehatan...")
    faskes = gen_faskes(lokasi_ids)
    csv_writer("fasilitas_kesehatan.csv",
        ["id_faskes","nama_faskes","tipe_faskes","id_lokasi","created_at","updated_at"],
        faskes)

    print("\n[8/19] Posyandu...")
    posyandu = gen_posyandu(lokasi_ids, bidan_ids)
    posyandu_ids = [r["id_posyandu"] for r in posyandu]
    csv_writer("posyandu.csv",
        ["id_posyandu","nama_posyandu","id_lokasi","id_bidan","created_at","updated_at"],
        posyandu)

    print("\n[9/19] Kader Posyandu...")
    csv_writer("kader_posyandu.csv",
        ["id_user","no_sk","id_posyandu","created_at","updated_at"],
        gen_kader(kader_ids, posyandu_ids))

    print("\n[10/19] Pasien...")
    pasien = gen_pasien(ibu_ids, anak_ids, pasien_dewasa_ids, posyandu_ids)
    csv_writer("pasien.csv",
        ["id_pasien","id_posyandu","created_at","updated_at"],
        pasien)

    print("\n[11/19] Ibu Hamil...")
    ibu_hamil = gen_ibu_hamil(ibu_ids)
    csv_writer("ibu_hamil.csv",
        ["id_ibu_hamil","id_pasien","hamil_ke","bulan_mulai_hamil","hpht","status_kehamilan","created_at","updated_at"],
        ibu_hamil)

    print("\n[12/19] Anak...")
    anak = gen_anak(anak_ids, ibu_hamil, pasien)
    csv_writer("anak.csv",
        ["id_pasien","id_ibu_hamil","id_wali","berat_lahir","panjang_lahir","hubungan_dengan_wali","created_at","updated_at"],
        anak)

    print("\n[13/19] Jadwal Imunisasi...")
    jadwal = gen_jadwal_imunisasi(anak_ids)
    csv_writer("jadwal_imunisasi.csv",
        ["id_imunisasi","id_pasien","nama_vaksin","tanggal_jadwal","tanggal_realisasi","status_imunisasi","created_at","updated_at"],
        jadwal)

    print("\n[14/19] Hasil Pemeriksaan...")
    hasil = gen_hasil_pemeriksaan(jadwal, bidan_ids)
    csv_writer("hasil_pemeriksaan.csv",
        ["id_hasil_pemeriksaan","id_petugas_input","id_jadwal_imunisasi","berat_badan","tinggi_badan","lingkar_kepala","tekanan_darah","status_stunting","status_gizi","catatan","created_at","updated_at"],
        hasil)

    print("\n[15/19] Tindak Lanjut...")
    tindak = gen_tindak_lanjut(hasil, bidan_ids)
    csv_writer("tindak_lanjut.csv",
        ["id_tindak_lanjut","id_hasil_pemeriksaan","id_bidan","catatan_medis","rekomendasi","jadwal_kontrol","status_pasien","created_at","updated_at"],
        tindak)

    print("\n[16/19] Rujukan...")
    rujukan = gen_rujukan(tindak, faskes)
    csv_writer("rujukan.csv",
        ["id_rujukan","id_tindak_lanjut","alasan_rujukan","tanggal_rujukan","status_rujukan","id_faskes","created_at","updated_at"],
        rujukan)

    print("\n[17/19] Notifikasi...")
    csv_writer("notifikasi.csv",
        ["id_notifikasi","id_user","judul","pesan","tipe_notifikasi","status_baca","tanggal_kirim","is_deleted","deleted_at"],
        gen_notifikasi(all_user_ids))

    print("\n[18/19] User Session...")
    csv_writer("user_session.csv",
        ["id_session","id_user","status_session","ip_address","created_at","updated_at","expired_at"],
        gen_user_session(all_user_ids))

    print("\n[19/19] Audit Log...")
    csv_writer("audit_log.csv",
        ["id_log","tipe_aktor","id_user","id_user_session","tipe_aktivitas","berhasil","endpoint","table_name","record_id","old_value","new_value","detail","ip_address","user_agent","waktu_aktivitas"],
        gen_audit_log(all_user_ids))

    print("\n=== SELESAI ===")
    print(f"Total user_account: {len(users)}")
    print(f"  - Dinkes: {len(dinkes_ids)}")
    print(f"  - Bidan: {len(bidan_ids)}")
    print(f"  - Kader: {len(kader_ids)}")
    print(f"  - Ibu Hamil: {len(ibu_ids)}")
    print(f"  - Anak: {len(anak_ids)}")
    print(f"  - Pasien Dewasa: {len(pasien_dewasa_ids)}")

if __name__ == "__main__":
    main()
