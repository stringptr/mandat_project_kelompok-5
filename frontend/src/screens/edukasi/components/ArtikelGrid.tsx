import type { Artikel } from '../data/artikel.data';
import { ArtikelCard } from './ArtikelCard';

interface ArtikelGridProps {
  artikelList: Artikel[];
  onRead: (a: Artikel) => void;
  onEdit?: (a: Artikel) => void;
  onHapus?: (a: Artikel) => void;
  onVerifikasi?: (a: Artikel) => void;
  canEdit?: (a: Artikel) => boolean;
  showVerifikasi?: boolean;
}

export function ArtikelGrid({
  artikelList,
  onRead,
  onEdit,
  onHapus,
  onVerifikasi,
  canEdit,
  showVerifikasi = false,
}: ArtikelGridProps): JSX.Element {
  if (artikelList.length === 0) return <></>;

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
      {artikelList.map((artikel) => (
        <ArtikelCard
          key={artikel.id}
          artikel={artikel}
          onRead={onRead}
          onEdit={onEdit}
          onHapus={onHapus}
          onVerifikasi={onVerifikasi}
          canEdit={canEdit ? canEdit(artikel) : false}
          showVerifikasi={showVerifikasi}
        />
      ))}
    </div>
  );
}
