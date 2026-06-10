/**
 * IbuWaliSection
 * Read-only layout. Shows hero + category filter + article grid.
 * No write/edit/delete actions.
 */
import type { Artikel, KategoriArtikel } from '../data/artikel.data';
import { EdukasiHero } from '../components/EdukasiHero';
import { KategoriFilter } from '../components/KategoriFilter';
import { ArtikelGrid } from '../components/ArtikelGrid';

interface IbuWaliSectionProps {
  featuredArtikel: Artikel[];
  gridArtikel: Artikel[];
  showAll: boolean;
  totalGrid: number;
  kategoriAktif: KategoriArtikel;
  searchQuery: string;
  onChangeKategori: (k: KategoriArtikel) => void;
  onChangeSearch: (q: string) => void;
  onRead: (a: Artikel) => void;
  onToggleShowAll: () => void;
}

export function IbuWaliSection({
  featuredArtikel,
  gridArtikel,
  showAll,
  totalGrid,
  kategoriAktif,
  searchQuery,
  onChangeKategori,
  onChangeSearch,
  onRead,
  onToggleShowAll,
}: IbuWaliSectionProps): JSX.Element {
  return (
    <div className="space-y-6">
      {featuredArtikel.length > 0 && (
        <EdukasiHero featured={featuredArtikel} onRead={onRead} />
      )}

      <KategoriFilter
        aktif={kategoriAktif}
        onChangeKategori={onChangeKategori}
        searchQuery={searchQuery}
        onChangeSearch={onChangeSearch}
      />

      <ArtikelGrid artikelList={gridArtikel} onRead={onRead} />

      {totalGrid > 6 && (
        <div className="text-center">
          <button
            onClick={onToggleShowAll}
            className="text-sm text-primary font-semibold hover:text-primary-600 font-body"
          >
            {showAll ? 'Tampilkan Lebih Sedikit ∧' : 'Lihat Materi Lainnya ∨'}
          </button>
        </div>
      )}
    </div>
  );
}
