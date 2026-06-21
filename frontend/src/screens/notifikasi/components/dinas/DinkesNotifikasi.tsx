import { Bell, MailOpen, CheckCheck } from 'lucide-react';
import type { BackendNotifikasi } from '../useNotifikasi';
import type { NotifGroup } from '../types';
import NotifTimeline from '../NotifTimeline';

interface DinkesNotifikasiProps {
  notifikasi: BackendNotifikasi[];
  markAllRead: () => Promise<void>;
  toNotifGroups: (list: BackendNotifikasi[]) => NotifGroup[];
}

export default function DinkesNotifikasi({ notifikasi, markAllRead, toNotifGroups }: DinkesNotifikasiProps): JSX.Element {
  const groups = toNotifGroups(notifikasi);
  const total = notifikasi.length;
  const unreadCount = notifikasi.filter((n) => !n.status_baca).length;
  const readCount = total - unreadCount;

  return (
    <div className="space-y-6">
      {/* Summary Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div className="bg-white rounded-2xl p-4 border border-neutral-100 shadow-sm flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center shrink-0">
            <Bell className="w-5 h-5 text-primary" />
          </div>
          <div>
            <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest">Total</p>
            <p className="text-xl font-bold text-neutral-800">{total}</p>
          </div>
        </div>
        <div className="bg-white rounded-2xl p-4 border border-neutral-100 shadow-sm flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center shrink-0">
            <MailOpen className="w-5 h-5 text-blue-600" />
          </div>
          <div>
            <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest">Belum Dibaca</p>
            <p className="text-xl font-bold text-neutral-800">{unreadCount}</p>
          </div>
        </div>
        <div className="bg-white rounded-2xl p-4 border border-neutral-100 shadow-sm flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-emerald-50 flex items-center justify-center shrink-0">
            <CheckCheck className="w-5 h-5 text-emerald-600" />
          </div>
          <div>
            <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest">Terbaca</p>
            <p className="text-xl font-bold text-neutral-800">{readCount}</p>
          </div>
        </div>
      </div>

      {/* Notification Timeline */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-base font-bold text-neutral-800 font-headline">Notifikasi</h3>
          {unreadCount > 0 && (
            <button onClick={markAllRead} className="text-sm font-semibold text-primary hover:text-primary-700 transition-colors">
              Tandai semua dibaca
            </button>
          )}
        </div>
        {groups.length === 0 ? (
          <div className="text-center py-12 text-neutral-400 text-sm">Belum ada notifikasi</div>
        ) : (
          <NotifTimeline groups={groups} />
        )}
      </div>
    </div>
  );
}
