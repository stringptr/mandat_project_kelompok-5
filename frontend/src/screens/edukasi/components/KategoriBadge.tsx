import type { KategoriArtikel } from '../data/artikel.data';

const KATEGORI_COLORS: Record<string, { solid: string; outline: string }> = {
  'Gizi Ibu': {
    solid: 'bg-emerald-600 text-white',
    outline: 'bg-emerald-50 text-emerald-700 border border-emerald-200',
  },
  'Nutrisi Anak': {
    solid: 'bg-blue-600 text-white',
    outline: 'bg-blue-50 text-blue-700 border border-blue-200',
  },
  'Parenting': {
    solid: 'bg-violet-600 text-white',
    outline: 'bg-violet-50 text-violet-700 border border-violet-200',
  },
  'Sanitasi & Lingkungan': {
    solid: 'bg-teal-600 text-white',
    outline: 'bg-teal-50 text-teal-700 border border-teal-200',
  },
  'Semua': {
    solid: 'bg-neutral-700 text-white',
    outline: 'bg-neutral-50 text-neutral-600 border border-neutral-200',
  },
};

interface KategoriBadgeProps {
  kategori: KategoriArtikel;
  variant?: 'solid' | 'outline';
}

export function KategoriBadge({ kategori, variant = 'outline' }: KategoriBadgeProps): JSX.Element {
  const colors = KATEGORI_COLORS[kategori] ?? KATEGORI_COLORS['Semua'];
  return (
    <span
      className={`text-[10px] font-bold uppercase tracking-wider px-2.5 py-1 rounded-full ${colors[variant]}`}
    >
      {kategori}
    </span>
  );
}
