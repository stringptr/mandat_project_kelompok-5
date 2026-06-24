# Panduan Penyusunan Paper/Skripsi - Sistem Data Warehouse untuk Monitoring Gizi dan Imunisasi

## 📋 Template Struktur Paper

---

## 1. ABSTRAK

### Komponen yang Harus Ada:
✅ **Latar Belakang Singkat** (2-3 kalimat)
✅ **Permasalahan** (1-2 kalimat)
✅ **Metode Penelitian** (1-2 kalimat)
✅ **Hasil Utama** (2-3 kalimat)
✅ **Kesimpulan** (1 kalimat)

### Contoh Draft:

**Bahasa Indonesia:**
```
Stunting dan malnutrisi pada anak serta Kekurangan Energi Kronis (KEK) pada ibu hamil 
masih menjadi masalah kesehatan serius di Indonesia, namun sistem monitoring yang ada 
masih berbasis pencatatan manual sehingga sulit untuk analisis dan pengambilan keputusan 
berbasis data. Penelitian ini mengembangkan sistem data warehouse dengan arsitektur 
3-tier (Staging-Master-Warehouse) menggunakan skema star schema untuk mendukung 
analisis multidimensional data gizi dan imunisasi. Sistem diimplementasikan menggunakan 
PostgreSQL untuk database warehouse, Apache Hop untuk ETL dengan teknik incremental 
loading dan selective column update, serta Go-Huma untuk backend API. Hasil implementasi 
menunjukkan efisiensi ETL sebesar 99.95% dengan incremental loading, mampu memproses 
hanya 50 dari 100.000 records yang berubah, dan menyediakan 9 dimensi serta 2 tabel 
fakta untuk analisis stunting, imunisasi, dan status gizi. Sistem ini memberikan 
foundation yang solid untuk Business Intelligence dan pengambilan keputusan berbasis 
data dalam program kesehatan masyarakat.
```

**Kata Kunci:** Data Warehouse, ETL, Star Schema, Incremental Loading, Stunting, 
Imunisasi, PostgreSQL, Apache Hop

---

## 2. PENDAHULUAN

### 2.1 Latar Belakang

#### Poin-poin Penting:

**1. Masalah Stunting & Malnutrisi di Indonesia**

📚 **Referensi yang Dibutuhkan (3 paper):**
- Paper tentang prevalensi stunting di Indonesia (WHO/Kemenkes)
- Paper tentang dampak malnutrisi terhadap perkembangan anak
- Paper tentang pentingnya monitoring kesehatan ibu hamil (KEK)

**Contoh Paragraf:**
```
Stunting dan malnutrisi pada anak di bawah lima tahun (balita) masih menjadi 
permasalahan kesehatan publik yang serius di Indonesia. Menurut [REFERENSI 1], 
prevalensi stunting di Indonesia mencapai XX% pada tahun 2023, menempatkan Indonesia 
dalam kategori negara dengan masalah stunting yang tinggi. Kondisi ini berdampak 
pada perkembangan kognitif dan fisik anak jangka panjang [REFERENSI 2]. Selain itu, 
Kekurangan Energi Kronis (KEK) pada ibu hamil juga berkontribusi terhadap kelahiran 
bayi dengan berat badan rendah dan risiko stunting [REFERENSI 3].
```

**2. Sistem Monitoring Saat Ini (Manual & Terfragmentasi)**

📚 **Referensi yang Dibutuhkan (3 paper):**
- Paper tentang tantangan sistem pencatatan kesehatan manual
- Paper tentang pentingnya digitalisasi data kesehatan
- Paper tentang data-driven decision making dalam kesehatan publik

**Contoh Paragraf:**
```
Sistem pencatatan data gizi dan imunisasi di tingkat Posyandu saat ini masih 
didominasi oleh pencatatan manual menggunakan buku register dan Kartu Menuju Sehat 
(KMS). Pendekatan manual ini menimbulkan berbagai kendala seperti inkonsistensi data, 
kesulitan agregasi, dan keterlambatan pelaporan [REFERENSI 4]. Studi oleh [REFERENSI 5] 
menunjukkan bahwa sistem informasi kesehatan yang terdigitalisasi dapat meningkatkan 
akurasi data hingga 95% dan mempercepat proses pengambilan keputusan. Lebih lanjut, 
pendekatan data-driven dalam program kesehatan terbukti efektif untuk early detection 
dan intervensi tepat sasaran [REFERENSI 6].
```

**3. Data Warehouse sebagai Solusi**

📚 **Referensi yang Dibutuhkan (3 paper):**
- Paper fundamental tentang data warehouse (Kimball/Inmon)
- Paper tentang implementasi data warehouse di healthcare
- Paper tentang ETL dan incremental loading techniques

**Contoh Paragraf:**
```
Data warehouse menyediakan pendekatan sistematis untuk mengintegrasikan data dari 
berbagai sumber ke dalam repository terpusat yang dioptimalkan untuk analisis 
[REFERENSI 7]. Berbeda dengan database transaksional (OLTP), data warehouse menggunakan 
skema dimensional (star/snowflake schema) yang memungkinkan query analitis yang kompleks 
dengan performa tinggi [REFERENSI 8]. Implementasi data warehouse di sektor kesehatan 
telah terbukti efektif untuk monitoring epidemiologi, analisis trend penyakit, dan 
evaluasi program kesehatan [REFERENSI 9].
```

### 2.2 Urgensi Penelitian

#### Poin-poin Penting:

**1. Gap Analysis: Sistem Manual vs Kebutuhan Analisis**

- ❌ Data tersebar di berbagai buku register
- ❌ Sulit melakukan analisis trend dan pola
- ❌ Tidak ada single source of truth
- ✅ **Urgensi:** Butuh sistem terintegrasi untuk analisis real-time

**2. Kompleksitas Data Kesehatan**
- Multiple entities: Pasien, Petugas, Posyandu, Lokasi
- Hierarchical data: Lokasi (Provinsi → Kota → Kecamatan → Kelurahan)
- Time-series data: Monitoring berkala dan historis
- ✅ **Urgensi:** Butuh skema yang support analisis multidimensional

**3. Kebutuhan Performance untuk Big Data**
- Potensi jutaan records dari ribuan posyandu
- Query analitis yang kompleks (aggregation, time-series, drill-down)
- ✅ **Urgensi:** Butuh ETL incremental agar tidak full load setiap kali

### 2.3 Research Question

**Rumusan Masalah:**
1. Bagaimana merancang arsitektur data warehouse yang sesuai untuk monitoring gizi dan imunisasi?
2. Bagaimana mengimplementasikan ETL pipeline dengan incremental loading untuk efisiensi proses?
3. Bagaimana mendesain skema star schema yang mendukung analisis multidimensional?
4. Bagaimana performa sistem dalam memproses data transaksional ke warehouse?

### 2.4 Tujuan Penelitian

**Tujuan Utama:**
Mengembangkan sistem data warehouse untuk monitoring gizi dan imunisasi dengan 
arsitektur 3-tier dan ETL incremental loading.

**Tujuan Khusus:**
1. Merancang arsitektur data warehouse 3-tier (Staging → Master → Warehouse)
2. Mengimplementasikan skema star schema dengan 9 dimensi dan 2 tabel fakta
3. Membangun ETL pipeline dengan incremental loading dan selective column update
4. Mengevaluasi performa sistem dalam hal efisiensi ETL dan kecepatan query

### 2.5 Manfaat Penelitian

**Manfaat Teoritis:**
- Kontribusi pada body of knowledge tentang implementasi data warehouse di healthcare
- Best practices untuk ETL incremental loading dalam context Indonesia

**Manfaat Praktis:**
- Sistem monitoring real-time untuk deteksi dini stunting dan malnutrisi
- Foundation untuk Business Intelligence dan data-driven decision making
- Template untuk implementasi di daerah lain

---

## 3. LANDASAN TEORI

### 3.1 Data Warehouse

#### Konsep Fundamental (Referensi: Kimball & Inmon)

**Definisi:**
```
Data warehouse adalah subject-oriented, integrated, time-variant, dan non-volatile 
collection of data yang mendukung proses pengambilan keputusan manajemen (Inmon, 1992).
```

**Karakteristik:**
1. **Subject-Oriented:** Diorganisasi berdasarkan subjek bisnis (pasien, imunisasi)
2. **Integrated:** Data dari berbagai sumber distandarkan dan diintegrasikan
3. **Time-Variant:** Menyimpan data historis dengan timestamp
4. **Non-Volatile:** Data tidak berubah setelah dimuat (read-only untuk users)

**OLTP vs OLAP:**
| Aspek | OLTP (Master DB) | OLAP (Warehouse) |
|-------|------------------|------------------|
| Purpose | Transaction processing | Analytical processing |
| Schema | Normalized (3NF) | Denormalized (Star) |
| Query | Simple, frequent | Complex, aggregated |
| Users | Many concurrent | Few analysts |
| Update | Frequent writes | Batch loads |

### 3.2 Star Schema

**Definisi:**
Skema dimensional dengan satu tabel fakta di tengah dan multiple dimensi di sekitarnya 
membentuk pola bintang.

**Komponen:**
1. **Fact Table:** Berisi measures/metrics (quantitative data)
   - Contoh: FACT_IMUNISASI, FACT_PEMERIKSAAN
   - Contains: Foreign keys to dimensions + numeric measures

2. **Dimension Table:** Berisi attributes (descriptive data)
   - Contoh: DIM_PASIEN, DIM_LOKASI, DIM_WAKTU
   - Contains: Primary key + descriptive attributes

**Keuntungan Star Schema:**
- ✅ Query performance tinggi (minimal joins)
- ✅ Intuitive untuk business users
- ✅ Flexible untuk ad-hoc queries
- ✅ Support OLAP operations (slice, dice, drill-down, roll-up)

### 3.3 ETL (Extract, Transform, Load)

**Definisi:**
Proses ekstraksi data dari sumber, transformasi sesuai business rules, dan loading 
ke target warehouse.

**Tahapan ETL:**

**1. Extract:**
