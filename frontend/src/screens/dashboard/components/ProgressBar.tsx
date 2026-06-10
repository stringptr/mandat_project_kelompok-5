interface ProgressBarProps {
  label: string;
  persen: number;
  color?: string;
}

export function ProgressBar({ label, persen, color = 'bg-primary' }: ProgressBarProps): JSX.Element {
  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between text-xs font-body">
        <span className="font-semibold text-neutral-600 uppercase tracking-wide text-[10px]">{label}</span>
        <span className="font-bold text-primary">{persen}%</span>
      </div>
      <div className="w-full bg-neutral-100 rounded-full h-2">
        <div
          className={`h-2 rounded-full transition-all duration-700 ${color}`}
          style={{ width: `${persen}%` }}
        />
      </div>
    </div>
  );
}
