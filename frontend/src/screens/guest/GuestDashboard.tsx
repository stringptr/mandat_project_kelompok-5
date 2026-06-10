/**
 * GuestDashboard — landing page untuk pengguna yang belum login.
 * Menampilkan:
 *  - Hero banner program SiGizi
 *  - Statistik publik Jawa Tengah
 *  - Preview artikel edukasi terbaru
 *  - CTA login per jenis pengguna
 */
import { Users, Baby, TrendingDown, BookOpen, ChevronRight, Clock, User } from 'lucide-react';
import { DUMMY_ARTIKEL } from '../edukasi/data/artikel.data';
import { KABUPATEN_DATA, TARGET_WILAYAH } from '../dashboard/data/dashboard.data';
import { ProgressBar } from '../dashboard/components/ProgressBar';
import { TrendChart } from '../dashboard/components/TrendChart';
import { TREN_NUTRISI } from '../dashboard/data/dashboard.data';

interface GuestDashboardProps {
  onLoginClick: () => void;
}

const PUBLIC_STATS = [
  { icon: <Users size={22} />, label: 'Total Pasien Terdaftar', value: '12,482', delta: '+13% bulan ini', color: 'bg-primary text-white' },
  { icon: <Baby size={22} />, label: 'Balita Dipantau', value: '8,934', delta: 'Aktif di 35 kab/kota', color: 'bg-blue-600 text-white' },
  { icon: <TrendingDown size={22} />, label: 'Kasus Stunting', value: '342', delta: '↓ Turun 8% dari tahun lalu', color: 'bg-white border border-neutral-100 text-neutral-800' },
  { icon: <BookOpen size={22} />, label: 'Artikel Edukasi', value: '120+', delta: 'Diperbarui setiap minggu', color: 'bg-white border border-neutral-100 text-neutral-800' },
];

const PREVIEW_ARTIKEL = DUMMY_ARTIKEL.filter((a) => a.status === 'published').slice(0, 3);

export function GuestDashboard({ onLoginClick }: GuestDashboardProps): JSX.Element {
  return (
    <div className="space-y-8 font-body text-neutral-800">

      {/* ── Hero banner ─────────────────────────────────────────────────── */}
      <div className="relative bg-primary rounded-3xl overflow-hidden px-8 py-10">
        {/* Decorative circles */}
        <div className="absolute -top-10 -right-10 w-48 h-48 rounded-full bg-white/10" />
        <div className="absolute bottom-0 right-24 w-32 h-32 rounded-full bg-white/10" />
        <div className="absolute top-6 right-48 w-16 h-16 rounded-full bg-white/5" />

        <div className="relative z-10 max-w-2xl">
          <span className="inline-block text-[10px] font-bold uppercase tracking-widest bg-white/20 text-white px-3 py-1 rounded-full mb-4">
            Sistem Informasi Gizi — Jawa Tengah
          </span>
          <h2 className="text-3xl font-bold text-white font-headline leading-tight mb-3">
            Wujudkan Generasi<br />
            <span className="text-emerald-300">Bebas Stunting</span> Bersama
          </h2>
          <p className="text-white/80 text-sm font-body leading-relaxed mb-6 max-w-lg">
            Platform monitoring gizi ibu dan anak berbasis komunitas. Pantau tumbuh kembang, akses edukasi kesehatan, dan kolaborasi lintas peran untuk Indonesia sehat.
          </p>
          <div className="flex items-center gap-3 flex-wrap">
            <button className="flex items-center gap-2 text-white/80 hover:text-white text-sm font-semibold transition-colors">
              Pelajari Program <ChevronRight size={14} />
            </button>
          </div>
        </div>
      </div>

      {/* ── Statistik publik ────────────────────────────────────────────── */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-base font-bold text-neutral-800 font-headline">Statistik Kesehatan Jawa Tengah</h3>
          <span className="text-xs text-neutral-400 font-body">Diperbarui: Februari 2024</span>
        </div>
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
          {PUBLIC_STATS.map((s) => (
            <div key={s.label} className={`rounded-2xl p-5 relative overflow-hidden ${s.color}`}>
              <div className={`mb-3 ${s.color.includes('bg-white') ? 'text-primary' : 'text-white/80'}`}>
                {s.icon}
              </div>
              <p className={`text-[10px] font-semibold uppercase tracking-wide mb-1 font-body ${s.color.includes('bg-white') ? 'text-neutral-500' : 'text-white/70'}`}>
                {s.label}
              </p>
              <p className={`text-3xl font-bold font-headline leading-none mb-1 ${s.color.includes('bg-white') ? 'text-neutral-900' : 'text-white'}`}>
                {s.value}
              </p>
              <p className={`text-xs font-body ${s.color.includes('bg-white') ? 'text-neutral-500' : 'text-white/70'}`}>
                {s.delta}
              </p>
              {!s.color.includes('bg-white') && (
                <div className="absolute -bottom-4 -right-4 w-20 h-20 rounded-full bg-white/10" />
              )}
            </div>
          ))}
        </div>
      </div>

      {/* ── Capaian program + Tren ──────────────────────────────────────── */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {/* Capaian program wilayah */}
        <div className="bg-white rounded-2xl p-6 border border-neutral-100">
          <div className="flex items-center gap-2 mb-5">
            <span className="text-lg">🎯</span>
            <div>
              <h4 className="text-sm font-bold text-neutral-800 font-headline">Capaian Program Wilayah</h4>
              <p className="text-xs text-neutral-500 font-body mt-0.5">Target nasional program gizi 2024</p>
            </div>
          </div>
          <div className="space-y-4">
            {TARGET_WILAYAH.map((t) => (
              <ProgressBar key={t.label} label={t.label} persen={t.persen} color={t.color} />
            ))}
          </div>
          <div className="mt-5 pt-4 border-t border-neutral-50 grid grid-cols-2 gap-3">
            {KABUPATEN_DATA.slice(0, 4).map((k) => (
              <div key={k.nama} className="flex items-center justify-between text-xs font-body">
                <span className="text-neutral-600 truncate mr-2">{k.nama.replace('Kab. ', '').replace('Kota ', '')}</span>
                <span className={`font-bold flex-shrink-0 ${k.level === 'tinggi' ? 'text-red-500' : k.level === 'sedang' ? 'text-orange-500' : 'text-emerald-600'}`}>
                  {k.prevalensi}%
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Tren nutrisi */}
        <div className="bg-blue-50 rounded-2xl p-6 border border-blue-100 flex flex-col">
          <div className="mb-1">
            <h4 className="text-sm font-bold text-neutral-800 font-headline">Tren Perbaikan Gizi Regional</h4>
            <p className="text-xs text-neutral-500 font-body mt-0.5">Indeks gizi Jawa Tengah — Okt 2023 s/d Apr 2024</p>
          </div>
          <div className="flex-1 mt-4">
            <TrendChart
              data={TREN_NUTRISI}
              height={120}
              color="#1d4ed8"
              fillColor="rgba(29,78,216,0.08)"
            />
          </div>
          <div className="mt-4 pt-4 border-t border-blue-100 flex items-center justify-between">
            <div>
              <p className="text-xs text-neutral-500 font-body">Tingkat keberhasilan intervensi</p>
              <p className="text-xl font-bold text-blue-700 font-headline">86.4%</p>
            </div>
            <div className="text-right">
              <p className="text-xs text-neutral-500 font-body">Kenaikan indeks</p>
              <p className="text-xl font-bold text-emerald-600 font-headline">+22 poin</p>
            </div>
          </div>
        </div>
      </div>

      {/* ── Artikel edukasi preview ─────────────────────────────────────── */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <div>
            <h3 className="text-base font-bold text-neutral-800 font-headline">Edukasi Kesehatan Terbaru</h3>
            <p className="text-xs text-neutral-500 font-body mt-0.5">Informasi gizi untuk keluarga Indonesia</p>
          </div>
          <button
            onClick={onLoginClick}
            className="text-xs text-primary font-semibold hover:text-primary-600 font-body flex items-center gap-1"
          >
            Lihat Semua <ChevronRight size={12} />
          </button>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          {PREVIEW_ARTIKEL.map((a) => (
            <div key={a.id} className="bg-white rounded-2xl overflow-hidden border border-neutral-100 hover:shadow-md transition-shadow group cursor-pointer" onClick={onLoginClick}>
              <div className="relative h-40 overflow-hidden">
                <img src={a.gambar} alt={a.judul} className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105" />
                <div className="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent" />
                <span className="absolute bottom-3 left-3 text-[10px] font-bold uppercase tracking-wide bg-primary text-white px-2.5 py-1 rounded-full">
                  {a.kategori}
                </span>
              </div>
              <div className="p-4">
                <h4 className="text-sm font-bold text-neutral-800 font-headline leading-snug mb-1.5 line-clamp-2">
                  {a.judul}
                </h4>
                <p className="text-xs text-neutral-500 font-body line-clamp-2 mb-3">{a.ringkasan}</p>
                <div className="flex items-center gap-3 text-neutral-400 text-xs">
                  <span className="flex items-center gap-1"><User size={11} />{a.penulis}</span>
                  <span className="flex items-center gap-1"><Clock size={11} />{a.waktuBaca}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

    </div>
  );
}
