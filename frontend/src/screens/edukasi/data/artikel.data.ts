import type { Role } from '../../../App';

export interface Artikel {
  id: string;
  judul: string;
  ringkasan: string;
  konten: string;
  kategori: string;
  penulis: string;
  rolePenulis: Role;
  tanggal: string;
  waktuBaca: string;
  gambar: string;
  status: 'published' | 'pending' | 'rejected';
  featured?: boolean;
}

export type KategoriArtikel =
  | 'Semua'
  | 'Gizi'
  | 'Imunisasi'
  | 'Kehamilan'
  | 'Penyakit'
  | 'Tumbuh Kembang';

export const KATEGORI_LIST: KategoriArtikel[] = [
  'Semua',
  'Gizi',
  'Imunisasi',
  'Kehamilan',
  'Penyakit',
  'Tumbuh Kembang',
];

export const WAKTU_BACA_OPTIONS = [
  '2 Menit Baca',
  '3 Menit Baca',
  '4 Menit Baca',
  '5 Menit Baca',
  '6 Menit Baca',
  '8 Menit Baca',
  '10 Menit Baca',
];

export const DUMMY_ARTIKEL: Artikel[] = [
  {
    id: 'a1',
    judul: 'Nutrisi Esensial Selama Trimester Pertama Kehamilan',
    ringkasan:
      'Pelajari jenis makanan super yang mendukung perkembangan otak janin dan menjaga kesehatan ibu di awal masa kehamilan.',
    konten:
      'Trimester pertama kehamilan adalah masa yang sangat kritis bagi perkembangan janin. Pada periode ini, organ-organ vital mulai terbentuk dan membutuhkan nutrisi yang tepat. Asam folat sangat penting untuk mencegah cacat tabung saraf. Konsumsilah sayuran hijau, kacang-kacangan, dan buah-buahan segar setiap hari. Zat besi dari daging merah, hati, dan bayam membantu mencegah anemia. Kalsium dari susu dan produk olahannya mendukung pembentukan tulang bayi. Protein dari telur, ikan, dan tahu mendukung pertumbuhan sel-sel baru. Hindari makanan mentah, alkohol, dan batasi kafein selama kehamilan.',
    kategori: 'Gizi Ibu',
    penulis: 'dr. Sarah Amalia',
    rolePenulis: 'Dinas Kesehatan',
    tanggal: '12 Okt 2023',
    waktuBaca: '5 Menit Baca',
    gambar:
      'https://images.unsplash.com/photo-1490645935967-10de6ba17061?auto=format&fit=crop&w=800&q=80',
    status: 'published',
    featured: true,
  },
  {
    id: 'a2',
    judul: 'Resep MPASI 6 Bulan: Tinggi Zat Besi',
    ringkasan:
      'Panduan praktis mengolah bahan lokal menjadi menu bergizi si kecil. Mudah dibuat dan disukai bayi.',
    konten:
      'MPASI (Makanan Pendamping ASI) mulai diberikan saat bayi berusia 6 bulan. Pada usia ini, bayi membutuhkan zat besi tambahan karena cadangan zat besi dari lahir mulai menipis. Resep puree hati ayam dengan kentang manis adalah pilihan yang kaya zat besi. Haluskan 50g hati ayam yang telah dimasak dengan 100g kentang manis kukus. Tambahkan ASI atau susu formula untuk mengatur kekentalan. Resep lainnya adalah bubur kacang merah yang kaya protein dan zat besi nabati. Sajikan dalam porsi kecil 2-3 kali sehari.',
    kategori: 'Nutrisi Anak',
    penulis: 'Bidan Sri Lestari',
    rolePenulis: 'Bidan',
    tanggal: '10 Okt 2023',
    waktuBaca: '4 Menit Baca',
    gambar:
      'https://images.unsplash.com/photo-1567620905732-2d1ec7ab7445?auto=format&fit=crop&w=800&q=80',
    status: 'published',
    featured: true,
  },
  {
    id: 'a3',
    judul: 'Menjaga Kesehatan Mental Ibu Pasca Melahirkan',
    ringkasan:
      'Baby blues dan postpartum depression adalah kondisi nyata yang perlu dikenali dan ditangani dengan tepat.',
    konten:
      'Setelah melahirkan, tubuh dan pikiran ibu mengalami perubahan besar. Baby blues adalah kondisi normal yang terjadi pada 50-80% ibu baru, ditandai dengan mudah menangis, cemas, dan lelah pada minggu pertama. Namun jika berlangsung lebih dari 2 minggu, bisa menjadi postpartum depression yang memerlukan penanganan profesional. Tanda-tandanya meliputi kesedihan mendalam, sulit bonding dengan bayi, dan pikiran untuk menyakiti diri sendiri. Dukungan keluarga sangat penting. Istirahat yang cukup, nutrisi seimbang, dan berbicara dengan orang yang dipercaya dapat membantu pemulihan.',
    kategori: 'Parenting',
    penulis: 'dr. Rizki Pratama',
    rolePenulis: 'Dinas Kesehatan',
    tanggal: '09 Okt 2023',
    waktuBaca: '6 Menit Baca',
    gambar:
      'https://images.unsplash.com/photo-1555252333-9f8e92e65df9?auto=format&fit=crop&w=800&q=80',
    status: 'published',
    featured: true,
  },
  {
    id: 'a4',
    judul: 'Dampak Kekurangan Vitamin A pada Pertumbuhan Balita',
    ringkasan:
      'Kekurangan vitamin A masih menjadi isu krusial. Ketahui cara mencegahnya dengan bahan-bahan lokal yang tersedia.',
    konten:
      'Vitamin A merupakan nutrisi esensial yang berperan dalam penglihatan, imunitas, dan pertumbuhan sel. Kekurangan vitamin A pada balita dapat menyebabkan gangguan penglihatan hingga kebutaan, menurunkan daya tahan tubuh, dan menghambat pertumbuhan. Sumber vitamin A yang baik meliputi wortel, labu kuning, ubi jalar merah, bayam, dan hati ayam. Pemberian kapsul vitamin A dosis tinggi setiap 6 bulan sekali untuk anak 6-59 bulan sangat dianjurkan. Program suplementasi ini tersedia gratis di Posyandu.',
    kategori: 'Nutrisi Anak',
    penulis: 'Bidan Dian Pertiwi',
    rolePenulis: 'Bidan',
    tanggal: '12 Okt 2023',
    waktuBaca: '4 Menit Baca',
    gambar:
      'https://images.unsplash.com/photo-1512621776951-a57141f2eefd?auto=format&fit=crop&w=800&q=80',
    status: 'published',
  },
  {
    id: 'a5',
    judul: 'ASI Eksklusif: Kunci Utama Imunitas Bayi di 1000 HPK',
    ringkasan:
      'Mengapa 6 bulan pertama begitu krusial? Temukan manfaat jangka panjang ASI bagi kecerdasan anak.',
    konten:
      '1000 Hari Pertama Kehidupan (HPK), dimulai dari kehamilan hingga anak berusia 2 tahun, adalah periode emas yang menentukan kualitas hidup seseorang. ASI eksklusif selama 6 bulan pertama memberikan semua nutrisi yang dibutuhkan bayi, mengandung antibodi alami yang melindungi dari infeksi, mendukung perkembangan otak optimal, mengurangi risiko obesitas dan diabetes di kemudian hari, serta mempererat ikatan emosional ibu dan bayi. Setelah 6 bulan, teruskan pemberian ASI bersama MPASI hingga anak berusia 2 tahun.',
    kategori: 'Gizi Ibu',
    penulis: 'Oleh Dinas Kota',
    rolePenulis: 'Dinas Kesehatan',
    tanggal: '10 Okt 2023',
    waktuBaca: '5 Menit Baca',
    gambar:
      'https://images.unsplash.com/photo-1584515933487-779824d29309?auto=format&fit=crop&w=800&q=80',
    status: 'published',
  },
  {
    id: 'a6',
    judul: 'Stimulasi Psikomotorik Anak Melalui Permainan Tradisional',
    ringkasan:
      'Membangun kecerdasan emosional sekaligus melatih koordinasi fisik anak dengan cara yang menyenangkan.',
    konten:
      'Permainan tradisional seperti congklak, engklek, dan gobak sodor bukan sekadar hiburan. Permainan ini melatih koordinasi motorik, kemampuan berpikir strategis, interaksi sosial, dan kecerdasan emosional anak. Engklek melatih keseimbangan dan koordinasi tubuh. Congklak melatih kemampuan berhitung dan strategi. Bermain kelompok meningkatkan kemampuan bersosialisasi. Alokasikan waktu 30-60 menit setiap hari untuk bermain bebas bersama teman sebaya. Batasi waktu layar (gadget) untuk mendukung perkembangan yang optimal.',
    kategori: 'Parenting',
    penulis: 'Oleh Psikolog Anak',
    rolePenulis: 'Dinas Kesehatan',
    tanggal: '09 Okt 2023',
    waktuBaca: '5 Menit Baca',
    gambar:
      'https://images.unsplash.com/photo-1503454537195-1dcabb73ffb9?auto=format&fit=crop&w=800&q=80',
    status: 'published',
  },
  {
    id: 'a7',
    judul: 'Cuci Tangan Pakai Sabun: Pencegahan Diare pada Balita',
    ringkasan:
      'Perilaku hidup bersih dan sehat dimulai dari kebiasaan sederhana yang bisa menyelamatkan ribuan nyawa.',
    konten:
      'Diare masih menjadi penyebab kematian balita terbesar di Indonesia. Cuci tangan pakai sabun dengan air mengalir adalah cara paling efektif mencegah penularan penyakit. Lakukan cuci tangan sebelum menyiapkan makanan, sebelum makan, setelah dari toilet, setelah mengganti popok bayi, dan setelah memegang hewan. Enam langkah cuci tangan yang benar memerlukan waktu minimal 20 detik. Ajarkan kebiasaan ini kepada anak sejak dini agar menjadi habit seumur hidup.',
    kategori: 'Sanitasi & Lingkungan',
    penulis: 'Tim Kesehatan Dinkes',
    rolePenulis: 'Dinas Kesehatan',
    tanggal: '08 Okt 2023',
    waktuBaca: '3 Menit Baca',
    gambar:
      'https://images.unsplash.com/photo-1584515933487-779824d29309?auto=format&fit=crop&w=800&q=80',
    status: 'published',
  },
  {
    id: 'a8',
    judul: 'Panduan Pemberian Tablet Tambah Darah untuk Remaja Putri',
    ringkasan:
      'Anemia pada remaja putri berisiko menurunkan prestasi belajar dan berdampak pada kehamilan di masa depan.',
    konten:
      'Anemia defisiensi besi sangat umum terjadi pada remaja putri akibat menstruasi dan pola makan yang kurang baik. Tablet Tambah Darah (TTD) mengandung zat besi dan asam folat yang diberikan 1 tablet per minggu sepanjang tahun untuk siswi SMP dan SMA. Minum TTD bersamaan dengan vitamin C (air jeruk atau jus tomat) meningkatkan penyerapan. Hindari minum TTD bersamaan dengan teh, kopi, atau susu. Konsumsi TTD secara rutin mencegah anemia, meningkatkan konsentrasi belajar, dan mempersiapkan remaja menjadi calon ibu yang sehat.',
    kategori: 'Gizi Ibu',
    penulis: 'Bidan Sri Lestari',
    rolePenulis: 'Bidan',
    tanggal: '05 Okt 2023',
    waktuBaca: '4 Menit Baca',
    gambar:
      'https://images.unsplash.com/photo-1576091160550-2173dba999ef?auto=format&fit=crop&w=800&q=80',
    status: 'pending',
  },
];
