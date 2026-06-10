import { useState } from 'react';
import type { KabupatenData } from '../data/dashboard.data';

interface PrevalensiMapProps {
  data: KabupatenData[];
}

// Simple visual representation using a grid of tiles (SVG map abstraction)
export function PrevalensiMap({ data }: PrevalensiMapProps): JSX.Element {
  const [hovered, setHovered] = useState<KabupatenData | null>(null);

  const levelColor: Record<string, string> = {
    tinggi: '#dc2626',
    sedang: '#f97316',
    rendah: '#22c55e',
  };

  const levelBg: Record<string, string> = {
    tinggi: 'bg-red-100 border-red-200',
    sedang: 'bg-orange-100 border-orange-200',
    rendah: 'bg-green-100 border-green-200',
  };

  const levelText: Record<string, string> = {
    tinggi: 'text-red-700',
    sedang: 'text-orange-700',
    rendah: 'text-green-700',
  };

  // Stylized Central Java silhouette via irregular tile grid
  // Row layout that loosely mimics the elongated shape of Jawa Tengah
  const layout: (number | null)[][] = [
    [null, 0, 1, null, null, null],
    [null, 2, 3, 4, 5, null],
    [null, 6, 7, null, null, null],
  ];

  return (
    <div className="relative">
      {/* Legend */}
      <div className="flex items-center gap-4 mb-4">
        {[
          { label: 'Tinggi (>25%)', color: '#dc2626' },
          { label: 'Sedang (15–25%)', color: '#f97316' },
          { label: 'Rendah (<15%)', color: '#22c55e' },
        ].map((l) => (
          <div key={l.label} className="flex items-center gap-1.5">
            <span className="w-3 h-3 rounded-sm flex-shrink-0" style={{ background: l.color }} />
            <span className="text-[10px] text-neutral-600 font-body">{l.label}</span>
          </div>
        ))}
      </div>

      {/* Tile grid map */}
      <div className="bg-neutral-50 rounded-2xl border border-neutral-100 p-6 relative min-h-48">
        <div className="flex flex-col gap-3 items-center">
          {layout.map((row, ri) => (
            <div key={ri} className="flex gap-3">
              {row.map((idx, ci) => {
                if (idx === null) {
                  return <div key={ci} className="w-24 h-16 opacity-0" />;
                }
                const kab = data[idx];
                if (!kab) return null;
                return (
                  <div
                    key={ci}
                    onMouseEnter={() => setHovered(kab)}
                    onMouseLeave={() => setHovered(null)}
                    className={`w-24 h-16 rounded-xl border-2 cursor-pointer transition-all hover:scale-105 hover:shadow-md flex flex-col items-center justify-center text-center relative ${levelBg[kab.level]}`}
                    style={{ borderColor: levelColor[kab.level] }}
                  >
                    <span className={`text-[10px] font-bold font-body leading-tight px-1 ${levelText[kab.level]}`}>
                      {kab.nama.replace('Kab. ', '').replace('Kota ', '')}
                    </span>
                    <span className={`text-sm font-bold font-headline ${levelText[kab.level]}`}>
                      {kab.prevalensi}%
                    </span>
                  </div>
                );
              })}
            </div>
          ))}
        </div>

        {/* Tooltip */}
        {hovered && (
          <div className="absolute top-4 right-4 bg-white rounded-xl shadow-lg border border-neutral-100 p-3 min-w-44 z-10">
            <p className="text-xs font-bold text-neutral-800 font-body">{hovered.nama}</p>
            <p className="text-xs text-neutral-500 font-body mt-1">
              Prevalensi: <span className="font-semibold text-neutral-700">{hovered.prevalensi}%</span>
            </p>
            <p className="text-xs text-neutral-500 font-body">
              Kasus: <span className="font-semibold text-neutral-700">{hovered.jumlahKasus.toLocaleString('id-ID')}</span>
            </p>
            <div className={`mt-2 text-[10px] font-bold uppercase px-2 py-0.5 rounded-full inline-block ${
              hovered.level === 'tinggi' ? 'bg-red-100 text-red-700' :
              hovered.level === 'sedang' ? 'bg-orange-100 text-orange-700' :
              'bg-green-100 text-green-700'
            }`}>
              {hovered.level === 'tinggi' ? 'Tinggi' : hovered.level === 'sedang' ? 'Sedang' : 'Rendah'}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
