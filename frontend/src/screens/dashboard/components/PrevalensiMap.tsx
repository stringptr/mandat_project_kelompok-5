interface MapItem {
  nama: string;
  prevalensi: number;
  jumlahKasus: number;
  level: string;
}

interface PrevalensiMapProps {
  data: MapItem[];
}

// Simplified geographic positions for 35 kab/kota in Jawa Tengah
// Each: [row (0-2), col (0-14), nama]
const KABUPATEN_POS: [number, number, string][] = [
  // North coast (row 0)
  [0, 0, 'Kab. Brebes'],
  [0, 1, 'Kota Tegal'],
  [0, 2, 'Kab. Tegal'],
  [0, 3, 'Kab. Pemalang'],
  [0, 4, 'Kab. Pekalongan'],
  [0, 5, 'Kota Pekalongan'],
  [0, 6, 'Kab. Batang'],
  [0, 7, 'Kab. Kendal'],
  [0, 8, 'Kota Semarang'],
  [0, 9, 'Kab. Demak'],
  [0, 10, 'Kab. Jepara'],
  [0, 11, 'Kab. Kudus'],
  [0, 12, 'Kab. Pati'],
  [0, 13, 'Kab. Rembang'],
  [0, 14, 'Kab. Blora'],
  // Central (row 1)
  [1, 0, 'Kab. Cilacap'],
  [1, 1, 'Kab. Banyumas'],
  [1, 2, 'Kab. Purbalingga'],
  [1, 3, 'Kab. Banjarnegara'],
  [1, 4, 'Kab. Wonosobo'],
  [1, 5, 'Kab. Temanggung'],
  [1, 6, 'Kota Salatiga'],
  [1, 7, 'Kab. Semarang'],
  [1, 8, 'Kab. Grobogan'],
  [1, 9, 'Kab. Sragen'],
  // South coast (row 2)
  [2, 1, 'Kab. Kebumen'],
  [2, 2, 'Kab. Purworejo'],
  [2, 3, 'Kab. Magelang'],
  [2, 4, 'Kota Magelang'],
  [2, 5, 'Kab. Boyolali'],
  [2, 6, 'Kab. Klaten'],
  [2, 7, 'Kab. Sukoharjo'],
  [2, 8, 'Kota Surakarta'],
  [2, 9, 'Kab. Karanganyar'],
  [2, 10, 'Kab. Wonogiri'],
];

function findData(name: string, data: MapItem[]): MapItem | undefined {
  return data.find(
    (d) =>
      d.nama.toLowerCase().includes(name.replace('Kab. ', '').replace('Kota ', '').toLowerCase()) ||
      name.toLowerCase().includes(d.nama.toLowerCase().replace('kab. ', '').replace('kota ', ''))
  );
}

function shortName(name: string): string {
  return name.replace('Kab. ', '').replace('Kota ', '');
}

export function PrevalensiMap({ data }: PrevalensiMapProps): JSX.Element {
  const cellW = 78;
  const cellH = 52;
  const gap = 4;
  const svgW = 15 * (cellW + gap);
  const svgH = 4 * (cellH + gap);

  function getFill(kab: MapItem | undefined) {
    if (!kab) return { bg: 'fill-neutral-100', border: 'stroke-neutral-200', text: 'fill-neutral-400', val: 'fill-neutral-400' };
    if (kab.level === 'tinggi') return { bg: 'fill-red-100', border: 'stroke-red-400', text: 'fill-red-700', val: 'fill-red-700' };
    if (kab.level === 'sedang') return { bg: 'fill-orange-100', border: 'stroke-orange-400', text: 'fill-orange-700', val: 'fill-orange-700' };
    return { bg: 'fill-green-100', border: 'stroke-green-400', text: 'fill-green-700', val: 'fill-green-700' };
  }

  return (
    <div className="relative">
      {/* Legend */}
      <div className="flex items-center gap-4 mb-4 flex-wrap">
        {[
          { label: 'Tinggi (>20%)', color: '#dc2626' },
          { label: 'Sedang (10-20%)', color: '#f97316' },
          { label: 'Rendah (<10%)', color: '#22c55e' },
          { label: 'Tidak Ada Data', color: '#e5e7eb' },
        ].map((l) => (
          <div key={l.label} className="flex items-center gap-1.5">
            <span className="w-3 h-3 rounded-sm flex-shrink-0" style={{ background: l.color }} />
            <span className="text-[10px] text-neutral-600 font-body">{l.label}</span>
          </div>
        ))}
      </div>

      {/* SVG Map */}
      <div className="bg-neutral-50 rounded-2xl border border-neutral-100 p-4 relative overflow-x-auto">
        <svg viewBox={`0 0 ${svgW} ${svgH + 20}`} className="w-full" style={{ maxWidth: svgW, minWidth: 700 }}>
          {/* Title */}
          <text x={svgW / 2} y={16} textAnchor="middle" className="fill-neutral-700" style={{ fontSize: 14, fontWeight: 700, fontFamily: 'system-ui' }}>
            PETA PREVALENSI STUNTING JAWA TENGAH
          </text>

          {KABUPATEN_POS.map(([row, col, name]) => {
            const kab = findData(name, data);
            const x = col * (cellW + gap);
            const y = (row + 0.7) * (cellH + gap);
            const s = getFill(kab);

            return (
              <g key={name}>
                <title>{name}{kab ? `\nPrevalensi: ${kab.prevalensi.toFixed(1)}%\nKasus: ${kab.jumlahKasus}\nLevel: ${kab.level}` : '\nTidak ada data'}</title>
                <rect x={x} y={y} width={cellW} height={cellH} rx={6} className={`${s.bg} ${s.border} transition-all hover:opacity-80`} strokeWidth={1.5} />
                <text x={x + cellW / 2} y={y + cellH / 2 - 5} textAnchor="middle" className={s.text} style={{ fontSize: 8, fontWeight: 700, fontFamily: 'system-ui' }}>
                  {shortName(name)}
                </text>
                {kab ? (
                  <text x={x + cellW / 2} y={y + cellH / 2 + 10} textAnchor="middle" className={s.val} style={{ fontSize: 10, fontWeight: 800, fontFamily: 'system-ui' }}>
                    {kab.prevalensi.toFixed(1)}%
                  </text>
                ) : (
                  <text x={x + cellW / 2} y={y + cellH / 2 + 9} textAnchor="middle" className={s.val} style={{ fontSize: 8, fontFamily: 'system-ui' }}>
                    N/A
                  </text>
                )}
              </g>
            );
          })}
        </svg>

        {/* Floating tooltip */}
        {data.length > 0 && (
          <div className="absolute top-4 right-4 bg-white rounded-xl shadow-lg border border-neutral-100 p-3 min-w-44 z-10">
            <p className="text-xs font-bold text-neutral-800 font-body mb-1">Ringkasan</p>
            <p className="text-xs text-neutral-500 font-body">
              Wilayah dengan data: <span className="font-semibold text-neutral-700">{Math.min(data.length, KABUPATEN_POS.length)}/{KABUPATEN_POS.length}</span>
            </p>
            <p className="text-xs text-neutral-500 font-body">
              Rata-rata: <span className="font-semibold text-neutral-700">
                {data.length > 0
                  ? (data.reduce((s, d) => s + d.prevalensi, 0) / data.length).toFixed(1)
                  : 0}%
              </span>
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
