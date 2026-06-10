import { Pencil, Trash2, Clock, User } from 'lucide-react';
import type { Artikel } from '../data/artikel.data';
import { KategoriBadge } from './KategoriBadge';

interface EdukasiHeroProps {
  featured: Artikel[];
  onRead: (a: Artikel) => void;
  onEdit?: (a: Artikel) => void;
  onHapus?: (a: Artikel) => void;
  canEdit?: (a: Artikel) => boolean;
}

export function EdukasiHero({
  featured,
  onRead,
  onEdit,
  onHapus,
  canEdit,
}: EdukasiHeroProps): JSX.Element {
  const [main, second, third] = featured;

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
      {/* Main featured card */}
      {main && (
        <div
          onClick={() => onRead(main)}
          className="lg:col-span-2 relative rounded-2xl overflow-hidden cursor-pointer group min-h-64"
        >
          <img
            src={main.gambar}
            alt={main.judul}
            className="absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/30 to-transparent" />

          <div className="absolute top-4 left-4">
            <KategoriBadge kategori={main.kategori} variant="solid" />
          </div>

          {canEdit?.(main) && (
            <div className="absolute top-4 right-4 flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
              {onEdit && (
                <button
                  onClick={(e) => { e.stopPropagation(); onEdit(main); }}
                  className="p-2 bg-white/20 hover:bg-white/40 backdrop-blur-sm rounded-lg text-white transition-colors"
                >
                  <Pencil size={14} />
                </button>
              )}
              {onHapus && (
                <button
                  onClick={(e) => { e.stopPropagation(); onHapus(main); }}
                  className="p-2 bg-red-500/70 hover:bg-red-500 backdrop-blur-sm rounded-lg text-white transition-colors"
                >
                  <Trash2 size={14} />
                </button>
              )}
            </div>
          )}

          <div className="absolute bottom-0 left-0 right-0 p-6">
            <h3 className="text-xl font-bold text-white font-headline leading-snug mb-2 line-clamp-2">
              {main.judul}
            </h3>
            <p className="text-white/80 text-sm line-clamp-2 mb-4 font-body">{main.ringkasan}</p>
            <div className="flex items-center gap-3 text-white/70 text-xs">
              <span className="flex items-center gap-1">
                <User size={12} /> {main.penulis}
              </span>
              <span className="flex items-center gap-1">
                <Clock size={12} /> {main.waktuBaca}
              </span>
            </div>
          </div>
        </div>
      )}

      {/* Side featured cards */}
      <div className="flex flex-col gap-4">
        {[second, third].filter(Boolean).map(
          (a) =>
            a && (
              <div
                key={a.id}
                onClick={() => onRead(a)}
                className="flex-1 bg-white rounded-2xl p-4 border border-neutral-100 cursor-pointer hover:shadow-md transition-shadow group relative"
              >
                <div className="flex items-start justify-between mb-2">
                  <KategoriBadge kategori={a.kategori} variant="outline" />
                  {canEdit?.(a) && (
                    <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      {onEdit && (
                        <button
                          onClick={(e) => { e.stopPropagation(); onEdit(a); }}
                          className="p-1.5 text-neutral-400 hover:text-primary hover:bg-primary-50 rounded-lg transition-colors"
                        >
                          <Pencil size={13} />
                        </button>
                      )}
                      {onHapus && (
                        <button
                          onClick={(e) => { e.stopPropagation(); onHapus(a); }}
                          className="p-1.5 text-neutral-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                        >
                          <Trash2 size={13} />
                        </button>
                      )}
                    </div>
                  )}
                </div>
                <h4 className="font-bold text-neutral-800 text-sm font-headline leading-snug mb-1 line-clamp-2">
                  {a.judul}
                </h4>
                <p className="text-xs text-neutral-500 font-body line-clamp-2 mb-3">{a.ringkasan}</p>
                <button className="text-xs text-primary font-semibold hover:text-primary-600 font-body">
                  Baca Selengkapnya →
                </button>
              </div>
            )
        )}
      </div>
    </div>
  );
}
