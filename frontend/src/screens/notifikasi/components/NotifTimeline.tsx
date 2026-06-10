import type { NotifGroup } from './types';
import NotifCard from './NotifCard';

interface NotifTimelineProps {
  groups: NotifGroup[];
  compactLastGroup?: boolean;
}

export default function NotifTimeline({ groups, compactLastGroup = true }: NotifTimelineProps): JSX.Element {
  const enriched = groups.map((g, idx) => ({ ...g, isLast: idx === groups.length - 1 }));

  if (enriched.length === 0) {
    return (
      <div className="py-16 text-center text-neutral-400 text-sm">
        Tidak ada notifikasi.
      </div>
    );
  }

  return (
    <div className="space-y-10">
      {enriched.map((group) => {
        const isCompact = compactLastGroup && group.isLast;
        return (
          <section key={group.groupLabel}>
            {/* Group label with divider */}
            <div className="flex items-center gap-4 mb-5">
              <span className="text-xs font-bold text-neutral-400 uppercase tracking-widest whitespace-nowrap">
                {group.groupLabel}
              </span>
              <div className="h-px bg-neutral-200 flex-1" />
            </div>

            {isCompact ? (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {group.items.map((item) => (
                  <NotifCard key={item.id} item={item} compact />
                ))}
              </div>
            ) : (
              <div className="space-y-4">
                {group.items.map((item) => (
                  <NotifCard key={item.id} item={item} />
                ))}
              </div>
            )}
          </section>
        );
      })}
    </div>
  );
}
