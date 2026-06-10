import { MoreVertical } from 'lucide-react';
import type { AktivitasPasien } from '../data/dashboard.data';
import { StatusBadge } from './StatusBadge';

interface AktivitasRowProps {
  item: AktivitasPasien;
}

export function AktivitasRow({ item }: AktivitasRowProps): JSX.Element {
  return (
    <div className="flex items-center gap-4 py-3 border-b border-neutral-50 last:border-0 hover:bg-neutral-50/60 rounded-xl px-2 -mx-2 transition-colors group">
      {/* Avatar */}
      <div className="w-10 h-10 rounded-full bg-neutral-100 flex items-center justify-center flex-shrink-0">
        <svg className="w-6 h-6 text-neutral-400" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z"/>
        </svg>
      </div>

      {/* Info */}
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-neutral-800 font-body">{item.nama}</p>
        <p className="text-xs text-neutral-500 font-body truncate">
          {item.tindakan} · <span className="text-neutral-400">{item.waktu}</span>
        </p>
      </div>

      {/* Status */}
      <StatusBadge status={item.status} />

      {/* Menu */}
      <button className="p-1 text-neutral-300 hover:text-neutral-500 opacity-0 group-hover:opacity-100 transition-opacity">
        <MoreVertical size={16} />
      </button>
    </div>
  );
}
