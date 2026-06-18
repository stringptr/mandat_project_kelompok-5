interface DataPoint {
  bulan: string;
  nilai: number;
}

interface TrendChartProps {
  data: DataPoint[];
  height?: number;
  color?: string;
  fillColor?: string;
  label?: string;
}

export function TrendChart({
  data,
  height = 80,
  color = '#095c3e',
  fillColor = 'rgba(9,92,62,0.08)',
  label,
}: TrendChartProps): JSX.Element {
  if (data.length < 2) return <></>;

  const values = data.map((d) => d.nilai);
  const min = Math.min(...values) - 5;
  const max = Math.max(...values) + 5;
  const range = max - min || 1;

  const W = 260;
  const H = height;
  const PAD = 12;

  const toX = (i: number) => PAD + (i / (data.length - 1)) * (W - PAD * 2);
  const toY = (v: number) => H - PAD - ((v - min) / range) * (H - PAD * 2);

  const points = data.map((d, i) => ({ x: toX(i), y: toY(d.nilai), bulan: d.bulan, nilai: d.nilai }));

  const linePath = points
    .map((p, i) => (i === 0 ? `M ${p.x} ${p.y}` : `L ${p.x} ${p.y}`))
    .join(' ');

  const areaPath = `${linePath} L ${points[points.length - 1].x} ${H} L ${points[0].x} ${H} Z`;

  return (
    <div className="w-full">
      {label && (
        <p className="text-xs font-semibold text-neutral-500 uppercase tracking-wide mb-2 font-body">{label}</p>
      )}
      <svg viewBox={`0 0 ${W} ${H}`} className="w-full" style={{ height }}>
        {/* Area fill */}
        <path d={areaPath} fill={fillColor} />
        {/* Line */}
        <path d={linePath} fill="none" stroke={color} strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
        {/* Dots */}
        {points.map((p, i) => (
          <circle key={i} cx={p.x} cy={p.y} r="3.5" fill="white" stroke={color} strokeWidth="2" />
        ))}
        {/* Last dot highlighted */}
        <circle cx={points[points.length - 1].x} cy={points[points.length - 1].y} r="5" fill={color} />
      </svg>
      {/* X labels */}
      <div className="flex justify-between mt-1 px-3">
        {data.map((d, i) => (
          <span key={i} className="text-[9px] text-neutral-400 font-body">{d.bulan}</span>
        ))}
      </div>
    </div>
  );
}
