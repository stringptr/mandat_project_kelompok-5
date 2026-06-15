const API = 'http://backend:8080/v1';

const users = [
  {
    email: 'dinkes@sigizi.id',
    password: 'password123',
    no_hp: '081234567892',
    nama: 'Dr. Budi Hermawan',
    nik: '3234567890123456',
    jenis_kelamin: 'Laki-Laki',
    tanggal_lahir: '1985-03-10T00:00:00Z',
    id_lokasi: 4,
    role: 'Dinkes',
  },
  {
    email: 'bidan@sigizi.id',
    password: 'password123',
    no_hp: '081234567890',
    nama: 'Bidan Sri Lestari',
    nik: '1234567890123456',
    jenis_kelamin: 'Perempuan',
    tanggal_lahir: '1990-01-15T00:00:00Z',
    id_lokasi: 4,
    role: 'Bidan',
    no_str: 'STR-001',
    wilayah_kerja: 4,
  },
  {
    email: 'kader@sigizi.id',
    password: 'password123',
    no_hp: '081234567891',
    nama: 'Siti Aminah',
    nik: '2234567890123456',
    jenis_kelamin: 'Perempuan',
    tanggal_lahir: '1995-06-20T00:00:00Z',
    id_lokasi: 4,
    role: 'Kader',
    no_sk: 'SK-001',
  },
  {
    email: 'ibu@sigizi.id',
    password: 'password123',
    no_hp: '081234567893',
    nama: 'Ny. Rina Marlina',
    nik: '4234567890123456',
    jenis_kelamin: 'Perempuan',
    tanggal_lahir: '1998-12-25T00:00:00Z',
    id_lokasi: 4,
    role: 'Ibu Hamil',
  },
];

async function main() {
  for (const user of users) {
    try {
      const res = await fetch(`${API}/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user),
      });
      const data = await res.json();
      console.log(`${user.email}: ${res.status} - ${data.detail || data.title}`);
    } catch (err) {
      console.log(`${user.email}: ERROR - ${err.message}`);
    }
  }
}

main();
