import type { BackendNotifikasi } from '../useNotifikasi';
import type { NotifGroup } from '../types';
import NotifTimeline from '../NotifTimeline';

interface KaderNotifikasiProps {
  notifikasi: BackendNotifikasi[];
  markAllRead: () => Promise<void>;
  toNotifGroups: (list: BackendNotifikasi[]) => NotifGroup[];
  onMarkRead?: (id: string) => void;
}

export default function KaderNotifikasi({ notifikasi, markAllRead, toNotifGroups, onMarkRead }: KaderNotifikasiProps): JSX.Element {
  const groups = toNotifGroups(notifikasi);
  const unreadCount = notifikasi.filter((n) => !n.status_baca).length;

  return (
    <div className="space-y-8">
      {/* Summary */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm">
          <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest mb-1">Total Notifikasi</p>
          <p className="text-2xl font-bold text-neutral-800">{notifikasi.length}</p>
        </div>
        <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm">
          <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest mb-1">Belum Dibaca</p>
          <p className="text-2xl font-bold text-neutral-800">{unreadCount}</p>
        </div>
      </div>

      {/* Timeline */}
      <div>
        <div className="flex items-center justify-between mb-5">
          <h3 className="text-lg font-bold text-neutral-800 font-headline">Notifikasi Terbaru</h3>
          {unreadCount > 0 && (
            <button onClick={markAllRead} className="text-sm font-semibold text-primary hover:text-primary-700 transition-colors">
              Tandai semua dibaca
            </button>
          )}
        </div>
        {groups.length === 0 ? (
          <div className="text-center py-12 text-neutral-400 text-sm">Belum ada notifikasi</div>
        ) : (
          <NotifTimeline groups={groups} onMarkRead={onMarkRead} />
        )}
      </div>
    </div>
  );
}
