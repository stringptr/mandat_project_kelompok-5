import type { ReactNode } from 'react';

export interface ChartBar {
  label: string;
  value: number;
  color?: string;
}

interface ChartWidgetProps {
  title: string;
  subtitle?: string;
  type: 'bar-horizontal' | 'bar-vertical' | 'donut';
  data: ChartBar[];
  action?: ReactNode;
}

function BarHorizontalChart({ data }: { data: ChartBar[] }) {
  const maxVal = Math.max(...data.map((d) => d.value), 1);
  return (
    <div className="space-y-3">
      {data.map((item, i) => (
        <div key={i} className="space-y-1">
          <div className="flex justify-between items-center text-xs">
            <span className="text-neutral-600 font-medium">{item.label}</span>
            <span className="text-neutral-800 font-bold">{item.value}</span>
          </div>
          <div className="w-full bg-neutral-100 h-2.5 rounded-full overflow-hidden">
            <div
              className="h-full rounded-full transition-all duration-700 ease-out"
              style={{ width: `${(item.value / maxVal) * 100}%`, backgroundColor: item.color || '#095c3e' }}
            />
          </div>
        </div>
      ))}
    </div>
  );
}

function BarVerticalChart({ data }: { data: ChartBar[] }) {
  const maxVal = Math.max(...data.map((d) => d.value), 1);
  return (
    <div className="flex items-end justify-between gap-3 h-40 px-2">
      {data.map((item, i) => (
        <div key={i} className="flex flex-col items-center gap-1.5 flex-1">
          <span className="text-[11px] font-bold text-neutral-700">{item.value}</span>
          <div className="w-full bg-neutral-100 rounded-t-lg overflow-hidden relative" style={{ height: '100%' }}>
            <div
              className="absolute bottom-0 left-0 right-0 rounded-t-lg transition-all duration-700 ease-out"
              style={{ height: `${(item.value / maxVal) * 100}%`, backgroundColor: item.color || '#095c3e', minHeight: '4px' }}
            />
          </div>
          <span className="text-[10px] text-neutral-500 font-medium text-center leading-tight truncate w-full">
            {item.label}
          </span>
        </div>
      ))}
    </div>
  );
}

function DonutChart({ data }: { data: ChartBar[] }) {
  const total = data.reduce((sum, d) => sum + d.value, 0) || 1;
  const defaultColors = ['#095c3e', '#3b82f6', '#f59e0b', '#ef4444', '#8b5cf6', '#10b981'];
  let cumulativePercent = 0;
  const gradientSegments = data.map((item, i) => {
    const percent = (item.value / total) * 100;
    const start = cumulativePercent;
    cumulativePercent += percent;
    return `${item.color || defaultColors[i % defaultColors.length]} ${start}% ${cumulativePercent}%`;
  });
  const background = `conic-gradient(${gradientSegments.join(', ')})`;

  return (
    <div className="flex items-center gap-6">
      <div className="relative flex-shrink-0">
        <div className="w-32 h-32 rounded-full" style={{ background }} />
        <div className="absolute inset-3 bg-white rounded-full flex items-center justify-center">
          <div className="text-center">
            <div className="text-xl font-bold text-neutral-800 font-headline">{total}</div>
            <div className="text-[10px] text-neutral-400 font-medium">Total</div>
          </div>
        </div>
      </div>
      <div className="space-y-2 flex-1">
        {data.map((item, i) => (
          <div key={i} className="flex items-center gap-2">
            <span className="w-3 h-3 rounded-full flex-shrink-0" style={{ backgroundColor: item.color || defaultColors[i % defaultColors.length] }} />
            <span className="text-xs text-neutral-600 flex-1">{item.label}</span>
            <span className="text-xs font-bold text-neutral-700">{item.value}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

export function ChartWidget({ title, subtitle, type, data, action }: ChartWidgetProps) {
  return (
    <div className="bg-white rounded-2xl border border-neutral-100 shadow-sm p-6 space-y-5">
      <div className="flex justify-between items-start">
        <div>
          <h3 className="text-base font-bold text-neutral-800 font-headline">{title}</h3>
          {subtitle && <p className="text-xs text-neutral-400 mt-0.5">{subtitle}</p>}
        </div>
        {action && <div>{action}</div>}
      </div>
      {type === 'bar-horizontal' && <BarHorizontalChart data={data} />}
      {type === 'bar-vertical' && <BarVerticalChart data={data} />}
      {type === 'donut' && <DonutChart data={data} />}
    </div>
  );
}
