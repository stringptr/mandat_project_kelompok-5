import { Search } from 'lucide-react';
import type { KategoriArtikel } from '../data/artikel.data';
import { KATEGORI_LIST } from '../data/artikel.data';

interface KategoriFilterProps {
  aktif: KategoriArtikel;
  onChangeKategori: (k: KategoriArtikel) => void;
  searchQuery: string;
  onChangeSearch: (q: string) => void;
}

export function KategoriFilter({
  aktif,
  onChangeKategori,
  searchQuery,
  onChangeSearch,
}: KategoriFilterProps): JSX.Element {
  return (
    <div className="flex items-center gap-2 flex-wrap">
      {KATEGORI_LIST.map((k) => (
        <button
          key={k}
          onClick={() => onChangeKategori(k)}
          className={`px-4 py-2 rounded-full text-sm font-semibold transition-all ${
            aktif === k
              ? 'bg-primary text-white shadow-sm'
              : 'bg-white border border-neutral-200 text-neutral-600 hover:border-primary hover:text-primary'
          }`}
        >
          {k}
        </button>
      ))}

      <div className="ml-auto relative">
        <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-400" />
        <input
          type="text"
          placeholder="Cari artikel..."
          value={searchQuery}
          onChange={(e) => onChangeSearch(e.target.value)}
          className="pl-9 pr-4 py-2 bg-white border border-neutral-200 rounded-full text-sm text-neutral-700 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 w-52"
        />
      </div>
    </div>
  );
}
