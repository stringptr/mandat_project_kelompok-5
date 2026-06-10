/**
 * IbuWaliSection
 * Personal view: status anak, jadwal imunisasi, riwayat pemeriksaan, tren tumbuh kembang.
 */
import { Baby, Calendar, ChevronRight, TrendingUp, Stethoscope } from 'lucide-react';
import { RIWAYAT_ANAK, JADWAL_TERDEKAT, TUMBUH_KEMBANG } from '../data/dashboard.data';
import { TrendChart } from '../components/TrendChart';
import { JadwalCard } from '../components/JadwalCard';

export function IbuWaliSection(): JSX.Element {
  const anak = { nama: 'Arka Ramadhan', usia: '14 bulan', bb: '10.2 kg', tb: '78 cm', status: 'Normal' };

  return (
    <div className="space-y-5">
      {/* Hero: status anak */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* Anak card */}
        <div className="lg:col-span-2 bg-primary rounded-2xl p-6 relative overflow-hidden">
          <div className="absolute -bottom-6 -right-6 w-32 h-32 rounded-full bg-white/10" />
          <div className="absolute -top-4 right-16 w-20 h-20 rounded-full bg-white/10" />
          <div className="flex items-start justify-between relative z-10">
            <div>
              <p className="text-white/70 text-xs font-semibold uppercase tracking-wide font-body mb-1">Status Anak Anda</p>
              <h3 className="text-white text-2xl font-bold font-headline">{anak.nama}</h3>
              <p className="text-white/80 text-sm font-body mt-1">{anak.usia} · {anak.bb} · {anak.tb}</p>
            </div>
            <div className="bg-emerald-400 text-white text-xs font-bold px-3 py-1.5 rounded-full">
              {anak.status}
            </div>
          </div>
          <div className="grid grid-cols-2 gap-3 mt-5 relative z-10">
            <div className="bg-white/15 rounded-xl p-3">
              <p className="text-white/70 text-[10px] font-semibold uppercase font-body">Berat Badan</p>
              <p className="text-white text-xl font-bold font-headline">{anak.bb}</p>
              <p className="text-white/60 text-xs font-body">+0.4 kg bulan ini</p>
            </div>
            <div className="bg-white/15 rounded-xl p-3">
              <p className="text-white/70 text-[10px] font-semibold uppercase font-body">Tinggi Badan</p>
              <p className="text-white text-xl font-bold font-headline">{anak.tb}</p>
              <p className="text-white/60 text-xs font-body">+2 cm bulan ini</p>
            </div>
          </div>
        </div>

        {/* Jadwal imunisasi next */}
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 flex flex-col">
          <div className="flex items-center gap-2 mb-4">
            <Calendar size={16} className="text-primary" />
            <p className="text-sm font-bold text-neutral-800 font-headline">Jadwal Imunisasi</p>
          </div>
          <div className="flex-1 space-y-1">
            {JADWAL_TERDEKAT.slice(0, 2).map((j) => (
              <JadwalCard key={j.id} item={j} />
            ))}
          </div>
          <button className="mt-4 w-full text-center text-xs text-primary font-semibold hover:text-primary-600 font-body flex items-center justify-center gap-1">
            Lihat Semua Jadwal <ChevronRight size={13} />
          </button>
        </div>
      </div>

      {/* Tren tumbuh kembang + riwayat */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {/* Tren */}
        <div className="bg-white rounded-2xl p-5 border border-neutral-100">
          <div className="flex items-center gap-2 mb-1">
            <TrendingUp size={16} className="text-primary" />
            <p className="text-sm font-bold text-neutral-800 font-headline">Tren Berat Badan</p>
          </div>
          <p className="text-xs text-neutral-400 font-body mb-4">6 bulan terakhir (kg)</p>
          <TrendChart
            data={TUMBUH_KEMBANG.map((d) => ({ bulan: d.bulan, nilai: d.bb }))}
            height={90}
          />
          <div className="mt-3 pt-3 border-t border-neutral-50 flex items-center justify-between">
            <span className="text-xs text-neutral-500 font-body">Target WHO usia 14 bln</span>
            <span className="text-xs font-bold text-primary font-body">9.8 – 11.5 kg ✓</span>
          </div>
        </div>

        {/* Riwayat pemeriksaan */}
        <div className="bg-white rounded-2xl p-5 border border-neutral-100">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-2">
              <Stethoscope size={16} className="text-primary" />
              <p className="text-sm font-bold text-neutral-800 font-headline">Riwayat Pemeriksaan</p>
            </div>
            <button className="text-xs text-primary font-semibold hover:text-primary-600 font-body">
              Lihat Semua
            </button>
          </div>
          <div className="space-y-0">
            {RIWAYAT_ANAK.map((r) => (
              <div key={r.id} className="flex items-start gap-3 py-3 border-b border-neutral-50 last:border-0">
                <div className={`w-2 h-2 rounded-full mt-1.5 flex-shrink-0 ${r.status === 'Normal' ? 'bg-emerald-500' : 'bg-amber-400'}`} />
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between gap-2">
                    <p className="text-sm font-semibold text-neutral-800 font-body truncate">{r.jenis}</p>
                    <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full flex-shrink-0 ${
                      r.status === 'Normal' ? 'bg-emerald-50 text-emerald-700' : 'bg-amber-50 text-amber-700'
                    }`}>
                      {r.status}
                    </span>
                  </div>
                  <p className="text-xs text-neutral-500 font-body mt-0.5">{r.hasil}</p>
                  <p className="text-[10px] text-neutral-400 font-body mt-0.5">{r.tanggal} · {r.petugas}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Tips edukasi */}
      <div className="bg-primary-50 border border-primary-100 rounded-2xl p-5 flex items-start gap-4">
        <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center flex-shrink-0">
          <Baby size={20} className="text-white" />
        </div>
        <div className="flex-1">
          <p className="text-sm font-bold text-primary font-headline">Tips Gizi Minggu Ini</p>
          <p className="text-xs text-neutral-600 font-body mt-1 leading-relaxed">
            Arka memasuki usia 14 bulan — waktu tepat memperkenalkan variasi makanan keluarga. Pastikan setiap makan mengandung sumber protein, karbohidrat, dan sayuran.
          </p>
        </div>
        <button className="text-xs text-primary font-semibold hover:text-primary-600 font-body flex-shrink-0">
          Baca Artikel →
        </button>
      </div>
    </div>
  );
}
