interface ProgressBarProps {
  label: string;
  value: number; // 0–100
  color?: 'primary' | 'blue' | 'amber' | 'red' | 'emerald' | 'purple';
  showPercent?: boolean;
  size?: 'sm' | 'md';
}

const BAR_COLORS = {
  primary: 'bg-primary',
  blue: 'bg-blue-600',
  amber: 'bg-amber-500',
  red: 'bg-red-500',
  emerald: 'bg-emerald-500',
  purple: 'bg-purple-500',
};

const TEXT_COLORS = {
  primary: 'text-primary',
  blue: 'text-blue-600',
  amber: 'text-amber-600',
  red: 'text-red-600',
  emerald: 'text-emerald-600',
  purple: 'text-purple-600',
};

export function ProgressBar({
  label,
  value,
  color = 'primary',
  showPercent = true,
  size = 'sm',
}: ProgressBarProps) {
  const clampedValue = Math.max(0, Math.min(100, value));
  const barHeight = size === 'sm' ? 'h-2' : 'h-3';

  return (
    <div className="space-y-1.5">
      <div className="flex justify-between items-center">
        <span className="text-xs text-neutral-600 font-medium">{label}</span>
        {showPercent && (
          <span className={`text-xs font-bold ${TEXT_COLORS[color]}`}>{clampedValue}%</span>
        )}
      </div>
      <div className={`w-full bg-neutral-100 ${barHeight} rounded-full overflow-hidden`}>
        <div
          className={`${BAR_COLORS[color]} ${barHeight} rounded-full transition-all duration-500 ease-out`}
          style={{ width: `${clampedValue}%` }}
        />
      </div>
    </div>
  );
}
