/**
 * BidanSection
 * Bidan can write articles (sent to Dinkes for verification before publish).
 * Can edit/delete only their own articles.
 * Shows a pending-verification banner.
 */
import type { Artikel, KategoriArtikel } from '../data/artikel.data';
import { EdukasiHero } from '../components/EdukasiHero';
import { KategoriFilter } from '../components/KategoriFilter';
import { ArtikelGrid } from '../components/ArtikelGrid';

interface BidanSectionProps {
  featuredArtikel: Artikel[];
  gridArtikel: Artikel[];
  showAll: boolean;
  totalGrid: number;
  kategoriAktif: KategoriArtikel;
  searchQuery: string;
  onChangeKategori: (k: KategoriArtikel) => void;
  onChangeSearch: (q: string) => void;
  onRead: (a: Artikel) => void;
  onEdit: (a: Artikel) => void;
  onHapus: (a: Artikel) => void;
  onTambah: () => void;
  canEdit: (a: Artikel) => boolean;
  onToggleShowAll: () => void;
}

export function BidanSection({
  featuredArtikel,
  gridArtikel,
  showAll,
  totalGrid,
  kategoriAktif,
  searchQuery,
  onChangeKategori,
  onChangeSearch,
  onRead,
  onEdit,
  onHapus,
  onTambah,
  canEdit,
  onToggleShowAll,
}: BidanSectionProps): JSX.Element {
  return (
    <div className="space-y-6">
      {/* Toolbar */}
      <div className="flex items-center justify-between">
        <div className="bg-amber-50 border border-amber-200 rounded-xl px-4 py-2.5 flex items-center gap-2 text-sm flex-1 mr-4">
          <span className="text-amber-500 flex-shrink-0">⏳</span>
          <span className="text-amber-700 font-medium font-body text-xs">
            Artikel yang Anda buat akan dikirim ke Dinas Kesehatan untuk diverifikasi sebelum dipublikasikan.
          </span>
        </div>
        <button
          onClick={onTambah}
          className="flex items-center gap-2 bg-primary hover:bg-primary-600 text-white rounded-xl px-5 py-2.5 text-sm font-semibold transition-colors shadow-sm flex-shrink-0"
        >
          <span className="text-lg leading-none">+</span>
          Tambah Artikel
        </button>
      </div>

      {featuredArtikel.length > 0 && (
        <EdukasiHero
          featured={featuredArtikel}
          onRead={onRead}
          onEdit={onEdit}
          onHapus={onHapus}
          canEdit={canEdit}
        />
      )}

      <KategoriFilter
        aktif={kategoriAktif}
        onChangeKategori={onChangeKategori}
        searchQuery={searchQuery}
        onChangeSearch={onChangeSearch}
      />

      <ArtikelGrid
        artikelList={gridArtikel}
        onRead={onRead}
        onEdit={onEdit}
        onHapus={onHapus}
        canEdit={canEdit}
      />

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
