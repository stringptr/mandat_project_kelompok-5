/**
 * BidanKaderSection — shared layout for Bidan & Kader Posyandu
 * Ringkasan monitoring: total pasien, perlu verifikasi, tindak lanjut,
 * aktivitas terbaru, target wilayah, jadwal terdekat.
 */
import { Users, FileSearch, AlertOctagon, Filter, ChevronRight } from 'lucide-react';
import type { Role } from '../../../App';
import { AKTIVITAS_TERBARU, JADWAL_TERDEKAT, TARGET_WILAYAH } from '../data/dashboard.data';
import { StatCard } from '../components/StatCard';
import { AktivitasRow } from '../components/AktivitasRow';
import { JadwalCard } from '../components/JadwalCard';
import { ProgressBar } from '../components/ProgressBar';

interface BidanKaderSectionProps {
  currentRole: Role;
}

export function BidanKaderSection({ currentRole }: BidanKaderSectionProps): JSX.Element {
  const isBidan = currentRole === 'Bidan';

  return (
    <div className="space-y-5">
      {/* Stat row */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        {/* Total pasien */}
        <StatCard
          variant="primary"
          label="Total Pasien Aktif"
          value="128"
          delta="+4% dari bulan lalu"
          deltaPositif
          icon={<Users size={22} />}
        />

        {/* Perlu verifikasi */}
        <StatCard
          variant="white"
          label="Perlu Verifikasi"
          value="12"
          badge="PRIORITY"
          icon={<FileSearch size={22} className="text-primary" />}
          sub={
            <button className="mt-3 text-xs text-primary font-semibold hover:text-primary-600 font-body flex items-center gap-1">
              Cek Sekarang <ChevronRight size={12} />
            </button>
          }
        />

        {/* Tindak lanjut */}
        <StatCard
          variant="danger"
          label="Tindak Lanjut"
          value="5"
          icon={<AlertOctagon size={22} className="text-red-500" />}
          sub={
            <button className="mt-3 text-xs text-red-500 font-semibold hover:text-red-600 font-body flex items-center gap-1">
              Lihat Daftar <ChevronRight size={12} />
            </button>
          }
        />
      </div>

      {/* Main content */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* Aktivitas terbaru */}
        <div className="lg:col-span-2 bg-white rounded-2xl p-5 border border-neutral-100">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-bold text-neutral-800 font-headline">Aktivitas Terbaru</h3>
            <button className="flex items-center gap-1.5 text-xs text-primary font-semibold hover:text-primary-600 font-body">
              <Filter size={13} /> Lihat Semua
            </button>
          </div>
          <div>
            {AKTIVITAS_TERBARU.map((item) => (
              <AktivitasRow key={item.id} item={item} />
            ))}
          </div>
        </div>

        {/* Right column */}
        <div className="flex flex-col gap-4">
          {/* Target wilayah */}
          <div className="bg-blue-50 rounded-2xl p-5 border border-blue-100">
            <div className="flex items-center gap-2 mb-4">
              <span className="text-base">🎯</span>
              <h3 className="text-sm font-bold text-neutral-800 font-headline">Target Wilayah</h3>
            </div>
            <div className="space-y-4">
              {TARGET_WILAYAH.map((t) => (
                <ProgressBar key={t.label} label={t.label} persen={t.persen} color={t.color} />
              ))}
            </div>
          </div>

          {/* Jadwal terdekat */}
          <div className="bg-white rounded-2xl p-5 border border-neutral-100 flex-1">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-bold text-neutral-800 font-headline">Jadwal Terdekat</h3>
            </div>
            <div>
              {JADWAL_TERDEKAT.map((j) => (
                <JadwalCard key={j.id} item={j} />
              ))}
            </div>
            <button className="mt-4 w-full py-2.5 border border-neutral-200 text-neutral-600 rounded-xl text-xs font-semibold hover:bg-neutral-50 transition-colors font-body">
              Atur Jadwal
            </button>
          </div>
        </div>
      </div>

      {/* Bidan-only: quick stats strip */}
      {isBidan && (
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
          {[
            { label: 'Ibu Hamil', value: '34', icon: '🤰', color: 'bg-pink-50 border-pink-100' },
            { label: 'Balita Dipantau', value: '94', icon: '👶', color: 'bg-blue-50 border-blue-100' },
            { label: 'Kunjungan Bulan Ini', value: '67', icon: '🏠', color: 'bg-emerald-50 border-emerald-100' },
            { label: 'Artikel Dikirim', value: '3', icon: '📝', color: 'bg-amber-50 border-amber-100' },
          ].map((s) => (
            <div key={s.label} className={`bg-white rounded-2xl p-4 border ${s.color}`}>
              <span className="text-xl">{s.icon}</span>
              <p className="text-2xl font-bold text-neutral-800 font-headline mt-1">{s.value}</p>
              <p className="text-xs text-neutral-500 font-body mt-0.5">{s.label}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
