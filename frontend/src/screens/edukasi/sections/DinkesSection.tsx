/**
 * DinkesSection
 * Dinas Kesehatan has full control:
 * - Write articles (published immediately)
 * - Edit/delete any article
 * - Verify (approve/reject) pending articles from Bidan
 * - Sees pending articles in the list
 */
import type { Artikel, KategoriArtikel } from '../data/artikel.data';
import { EdukasiHero } from '../components/EdukasiHero';
import { KategoriFilter } from '../components/KategoriFilter';
import { ArtikelGrid } from '../components/ArtikelGrid';

interface DinkesSectionProps {
  featuredArtikel: Artikel[];
  gridArtikel: Artikel[];
  showAll: boolean;
  totalGrid: number;
  pendingCount: number;
  kategoriAktif: KategoriArtikel;
  searchQuery: string;
  onChangeKategori: (k: KategoriArtikel) => void;
  onChangeSearch: (q: string) => void;
  onRead: (a: Artikel) => void;
  onEdit: (a: Artikel) => void;
  onHapus: (a: Artikel) => void;
  onTambah: () => void;
  onVerifikasi: (a: Artikel) => void;
  onBukaAntrianVerifikasi: () => void;
  canEdit: (a: Artikel) => boolean;
  onToggleShowAll: () => void;
}

export function DinkesSection({
  featuredArtikel,
  gridArtikel,
  showAll,
  totalGrid,
  pendingCount,
  kategoriAktif,
  searchQuery,
  onChangeKategori,
  onChangeSearch,
  onRead,
  onEdit,
  onHapus,
  onTambah,
  onVerifikasi,
  onBukaAntrianVerifikasi,
  canEdit,
  onToggleShowAll,
}: DinkesSectionProps): JSX.Element {
  return (
    <div className="space-y-6">
      {/* Toolbar */}
      <div className="flex items-center gap-3 flex-wrap">
        {pendingCount > 0 && (
          <button
            onClick={onBukaAntrianVerifikasi}
            className="flex items-center gap-2 bg-amber-500 hover:bg-amber-600 text-white rounded-xl px-4 py-2.5 text-sm font-semibold transition-colors"
          >
            <span className="w-5 h-5 bg-white text-amber-500 rounded-full text-xs font-bold flex items-center justify-center flex-shrink-0">
              {pendingCount}
            </span>
            Verifikasi Artikel
          </button>
        )}
        <button
          onClick={onTambah}
          className="flex items-center gap-2 bg-primary hover:bg-primary-600 text-white rounded-xl px-5 py-2.5 text-sm font-semibold transition-colors shadow-sm ml-auto"
        >
          <span className="text-lg leading-none">+</span>
          Tambah Artikel Baru
        </button>
      </div>

      {/* Pending summary strip */}
      {pendingCount > 0 && (
        <div className="bg-amber-50 border border-amber-200 rounded-xl px-4 py-3 flex items-center justify-between">
          <div className="flex items-center gap-2 text-sm">
            <span className="text-amber-500">📋</span>
            <span className="text-amber-700 font-medium font-body">
              Terdapat <strong>{pendingCount}</strong> artikel dari Bidan yang menunggu verifikasi Anda.
            </span>
          </div>
          <button
            onClick={onBukaAntrianVerifikasi}
            className="text-xs text-amber-700 font-semibold hover:text-amber-800 underline font-body flex-shrink-0 ml-4"
          >
            Tinjau Sekarang →
          </button>
        </div>
      )}

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
        onVerifikasi={onVerifikasi}
        canEdit={canEdit}
        showVerifikasi
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
