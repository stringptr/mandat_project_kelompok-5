import type { Role } from '../../App';
import { useNotifikasi } from './components/useNotifikasi';
import { Paginator } from '../../components/Paginator';
import KaderNotifikasi from './components/kader/KaderNotifikasi';
import BidanNotifikasi from './components/bidan/BidanNotifikasi';
import DinkesNotifikasi from './components/dinas/DinkesNotifikasi';
import IbuWaliNotifikasi from './components/ibu-wali/IbuWaliNotifikasi';

interface NotifikasiProps {
  role: Role;
}

export default function Notifikasi({ role }: NotifikasiProps): JSX.Element {
  const { notifikasi, meta, fetchNotifikasi, markAllRead, handleCardMarkRead, toNotifGroups } = useNotifikasi();

  return (
    <div className="w-full max-w-5xl mx-auto pb-12 font-body text-neutral-800">
      {role === 'Kader Posyandu' && (
        <KaderNotifikasi
          notifikasi={notifikasi}
          markAllRead={markAllRead}
          toNotifGroups={toNotifGroups}
          onMarkRead={handleCardMarkRead}
        />
      )}
      {role === 'Bidan' && (
        <BidanNotifikasi
          notifikasi={notifikasi}
          markAllRead={markAllRead}
          toNotifGroups={toNotifGroups}
        />
      )}
      {role === 'Dinas Kesehatan' && (
        <DinkesNotifikasi
          notifikasi={notifikasi}
          markAllRead={markAllRead}
          toNotifGroups={toNotifGroups}
        />
      )}
      {role === 'Ibu/Wali' && (
        <IbuWaliNotifikasi
          notifikasi={notifikasi}
          markAllRead={markAllRead}
          toNotifGroups={toNotifGroups}
        />
      )}
      <Paginator
        page={meta.current_page}
        totalPages={meta.last_page}
        totalItems={meta.total}
        pageSize={meta.per_page}
        onPageChange={fetchNotifikasi}
      />
    </div>
  );
}
