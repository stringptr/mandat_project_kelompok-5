import { useEffect } from 'react';
import { Calendar } from 'lucide-react';
import { apiGet } from '../../../lib/api';
import { useAppStore } from '../../../store/useAppStore';
import { useAuth } from '../../../context/AuthContext';

export function IbuWaliSection(): JSX.Element {
  const { user } = useAuth();

  const {
    jadwalTerdekat,
    setJadwalTerdekat,
  } = useAppStore();

  useEffect(() => {
    if (jadwalTerdekat.length > 0) return;

    apiGet<{ jadwal: Record<string, unknown>[] }>('/dashboard/jadwal-terdekat')
      .then((res) => setJadwalTerdekat(res.jadwal.slice(0, 3)))
      .catch(() => console.error('Gagal memuat jadwal'));
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const childName = user?.name ?? 'Anak Anda';

  return (
    <div className="space-y-5">
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <div className="lg:col-span-2 bg-primary rounded-2xl p-6 relative overflow-hidden">
          <div className="absolute -bottom-6 -right-6 w-32 h-32 rounded-full bg-white/10" />
          <div className="absolute -top-4 right-16 w-20 h-20 rounded-full bg-white/10" />
          <div className="flex items-start justify-between relative z-10">
            <div>
              <p className="text-white/70 text-xs font-semibold uppercase tracking-wide font-body mb-1">
                Selamat Datang
              </p>
              <h3 className="text-white text-2xl font-bold font-headline">{childName}</h3>
              <p className="text-white/80 text-sm font-body mt-1">
                Pantau tumbuh kembang anak Anda secara berkala
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-2xl p-5 border border-neutral-100 flex flex-col">
          <div className="flex items-center gap-2 mb-4">
            <Calendar size={16} className="text-primary" />
            <p className="text-sm font-bold text-neutral-800 font-headline">Jadwal Imunisasi</p>
          </div>
          <div className="flex-1">
            {jadwalTerdekat.length === 0 ? (
              <p className="text-sm text-neutral-400 text-center py-4">Belum ada jadwal</p>
            ) : (
              jadwalTerdekat.map((j, i) => (
                <div key={String(j.id ?? i)} className="text-sm py-1.5 border-b border-neutral-50 last:border-0">
                  <p className="font-medium text-neutral-700">{String(j.nama_vaksin ?? '')}</p>
                  <p className="text-xs text-neutral-500">{String(j.tanggal_jadwal ?? '')}</p>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
