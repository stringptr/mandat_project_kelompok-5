import type { NotifGroup } from './types';
import NotifCard from './NotifCard';

interface NotifTimelineProps {
  groups: NotifGroup[];
  onMarkRead?: (id: string) => void;
}

export default function NotifTimeline({ groups, onMarkRead }: NotifTimelineProps): JSX.Element {
  if (groups.length === 0) {
    return (
      <div className="py-16 text-center text-neutral-400 text-sm">
        Tidak ada notifikasi.
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {groups.map((group) => (
        <section key={group.groupLabel}>
          <div className="flex items-center gap-3 mb-3">
            <span className="text-[10px] font-bold text-neutral-400 uppercase tracking-wider whitespace-nowrap">
              {group.groupLabel}
            </span>
            <div className="h-px bg-neutral-100 flex-1" />
          </div>
          <div className="space-y-1">
            {group.items.map((item) => (
              <NotifCard key={item.id} item={item} onMarkRead={onMarkRead} />
            ))}
          </div>
        </section>
      ))}
    </div>
  );
}
