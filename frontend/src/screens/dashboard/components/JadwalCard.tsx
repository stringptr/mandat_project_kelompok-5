import type { JadwalItem } from '../data/dashboard.data';

interface JadwalCardProps {
  item: JadwalItem;
}

export function JadwalCard({ item }: JadwalCardProps): JSX.Element {
  return (
    <div className="flex items-start gap-3 py-2.5 border-b border-neutral-50 last:border-0">
      {/* Date chip */}
      <div className="flex flex-col items-center bg-primary-50 rounded-xl px-3 py-2 flex-shrink-0 min-w-12 text-center">
        <span className="text-[10px] font-bold text-primary uppercase tracking-wide">{item.bulan}</span>
        <span className="text-lg font-bold text-primary font-headline leading-tight">{item.tanggal}</span>
      </div>

      {/* Info */}
      <div className="min-w-0">
        <p className="text-sm font-semibold text-neutral-800 font-body">{item.judul}</p>
        <p className="text-xs text-neutral-500 font-body mt-0.5">{item.lokasi} · {item.jam}</p>
      </div>
    </div>
  );
}
