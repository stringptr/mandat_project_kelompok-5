import { CalendarDays, CheckCircle2 } from 'lucide-react';
import NotifTimeline from '../NotifTimeline';
import { useNotifikasi } from '../useNotifikasi';

export default function IbuWaliNotifikasi(): JSX.Element {
  const { loading, notifikasi, markAllRead, toNotifGroups } = useNotifikasi();

  const groups = toNotifGroups(notifikasi);
  const unreadCount = notifikasi.filter((n) => !n.status_baca).length;

  return (
    <div>
      {/* Hero Summary Card */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-5 mb-10">
        <div className="lg:col-span-2 bg-gradient-to-br from-emerald-500 to-blue-500 rounded-3xl p-7 text-white shadow-sm relative overflow-hidden flex flex-col justify-between min-h-[200px]">
          <div className="absolute top-0 right-0 w-60 h-60 bg-white/5 rounded-full -translate-y-1/2 translate-x-1/4" />
          <div className="relative z-10 flex justify-between items-start mb-3">
            <span className="bg-white/20 backdrop-blur-sm text-white text-[10px] font-bold px-3 py-1.5 rounded-full uppercase tracking-widest">
              Notifikasi
            </span>
            {unreadCount > 0 && (
              <span className="bg-white/20 backdrop-blur-sm text-white text-[10px] font-bold px-3 py-1.5 rounded-full">
                {unreadCount} belum dibaca
              </span>
            )}
          </div>
          <div className="relative z-10">
            <h2 className="text-xl font-bold font-headline mb-2 leading-snug">
              Pantau Kesehatan Anak Anda
            </h2>
            <p className="text-white/85 text-sm mb-5 leading-relaxed max-w-sm">
              {loading ? 'Memuat...' : `${notifikasi.length} notifikasi tersedia`}
            </p>
          </div>
        </div>

        <div className="flex flex-col gap-4">
          <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm flex items-start gap-4 flex-1">
            <div className="w-11 h-11 rounded-xl bg-emerald-50 flex items-center justify-center shrink-0">
              <CalendarDays className="w-5 h-5 text-emerald-600" />
            </div>
            <div>
              <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest mb-0.5">Total Notifikasi</p>
              <p className="font-bold text-neutral-800 text-sm">{notifikasi.length}</p>
              <p className="text-xs text-neutral-500">{unreadCount} belum dibaca</p>
            </div>
          </div>
          <div className="bg-white rounded-2xl p-5 border border-neutral-100 shadow-sm flex items-start gap-4 flex-1">
            <div className="w-11 h-11 rounded-xl bg-emerald-50 flex items-center justify-center shrink-0">
              <CheckCircle2 className="w-5 h-5 text-emerald-600" />
            </div>
            <div>
              <p className="text-[10px] font-bold text-neutral-400 uppercase tracking-widest mb-0.5">Terbaca</p>
              <p className="font-bold text-neutral-800 text-sm">{notifikasi.filter((n) => n.status_baca).length}</p>
              <p className="text-xs text-neutral-500">notifikasi sudah dibaca</p>
            </div>
          </div>
        </div>
      </div>

      {/* Notification Timeline */}
      <div className="mb-10">
        <div className="flex items-center justify-between mb-5">
          <h3 className="text-lg font-bold text-neutral-800 font-headline">Pemberitahuan Terbaru</h3>
          {unreadCount > 0 && (
            <button onClick={markAllRead} className="text-sm font-semibold text-primary hover:text-primary-700 transition-colors">
              Tandai semua dibaca
            </button>
          )}
        </div>
        {loading ? (
          <div className="text-center py-12 text-neutral-400 text-sm">Memuat notifikasi...</div>
        ) : groups.length === 0 ? (
          <div className="text-center py-12 text-neutral-400 text-sm">Belum ada notifikasi</div>
        ) : (
          <NotifTimeline groups={groups} compactLastGroup />
        )}
      </div>
    </div>
  );
}
