export type StatusImunisasi = 'SUDAH' | 'BELUM';

export interface JadwalImunisasi {
  id: string;
  idPasien: string;
  namaAnak: string;
  namaVaksin: string;
  dosis: string;
  tanggalJadwal: string;
  tanggalRealisasi: string | null;
  status: StatusImunisasi;
  catatan?: string;
}

export const VAKSIN_OPTIONS = [
  'BCG (Bacillus Calmette-Guérin)',
  'Hepatitis B (HB)',
  'Polio (IPV)',
  'DPT-HB-Hib',
  'Campak-Rubela (MR)',
  'Pneumokokus (PCV)',
  'Rotavirus',
  'Japanese Encephalitis (JE)',
  'HPV',
  'Influenza',
];

export const DUMMY_JADWAL: JadwalImunisasi[] = [
  { id: '1', idPasien: '#PST-09221', namaAnak: 'Arka Ramadhan',   namaVaksin: 'BCG (Bacillus Calmette-Guérin)', dosis: 'Primary Dose',   tanggalJadwal: '12 Okt 2023', tanggalRealisasi: '12 Okt 2023', status: 'SUDAH' },
  { id: '2', idPasien: '#PST-09225', namaAnak: 'Rayyan Putra',    namaVaksin: 'Hepatitis B (HB)',               dosis: 'Second Dosage',  tanggalJadwal: '15 Okt 2023', tanggalRealisasi: null,           status: 'BELUM' },
  { id: '3', idPasien: '#PST-09228', namaAnak: 'Siti Rahayu',     namaVaksin: 'Polio (IPV)',                    dosis: 'Booster Routine',tanggalJadwal: '18 Okt 2023', tanggalRealisasi: '19 Okt 2023', status: 'SUDAH' },
  { id: '4', idPasien: '#PST-09230', namaAnak: 'Budi Santoso',    namaVaksin: 'DPT-HB-Hib',                    dosis: 'Third Dosage',   tanggalJadwal: '20 Okt 2023', tanggalRealisasi: null,           status: 'BELUM' },
  { id: '5', idPasien: '#PST-09234', namaAnak: 'Melati Dewi',     namaVaksin: 'Campak-Rubela (MR)',             dosis: 'Primary Dose',   tanggalJadwal: '22 Okt 2023', tanggalRealisasi: '22 Okt 2023', status: 'SUDAH' },
  { id: '6', idPasien: '#PST-09237', namaAnak: 'Ahmad Fauzi',     namaVaksin: 'Pneumokokus (PCV)',              dosis: 'First Dosage',   tanggalJadwal: '25 Okt 2023', tanggalRealisasi: null,           status: 'BELUM' },
  { id: '7', idPasien: '#PST-09240', namaAnak: 'Intan Permata',   namaVaksin: 'Rotavirus',                      dosis: 'Second Dosage',  tanggalJadwal: '28 Okt 2023', tanggalRealisasi: '28 Okt 2023', status: 'SUDAH' },
  { id: '8', idPasien: '#PST-09243', namaAnak: 'Galih Pratama',   namaVaksin: 'DPT-HB-Hib',                    dosis: 'Booster',        tanggalJadwal: '30 Okt 2023', tanggalRealisasi: null,           status: 'BELUM' },
  { id: '9', idPasien: '#PST-09246', namaAnak: 'Citra Lestari',   namaVaksin: 'Polio (IPV)',                    dosis: 'Third Dosage',   tanggalJadwal: '01 Nov 2023', tanggalRealisasi: '01 Nov 2023', status: 'SUDAH' },
  { id: '10', idPasien: '#PST-09249', namaAnak: 'Rizky Hidayat',  namaVaksin: 'Hepatitis B (HB)',               dosis: 'Third Dosage',   tanggalJadwal: '03 Nov 2023', tanggalRealisasi: null,           status: 'BELUM', catatan: 'Kondisi anak sedang kurang sehat' },
  { id: '11', idPasien: '#PST-09252', namaAnak: 'Putri Cahaya',   namaVaksin: 'BCG (Bacillus Calmette-Guérin)', dosis: 'Primary Dose',   tanggalJadwal: '05 Nov 2023', tanggalRealisasi: '05 Nov 2023', status: 'SUDAH' },
  { id: '12', idPasien: '#PST-09255', namaAnak: 'Dimas Satria',   namaVaksin: 'Campak-Rubela (MR)',             dosis: 'Second Dosage',  tanggalJadwal: '08 Nov 2023', tanggalRealisasi: null,           status: 'BELUM' },
];
