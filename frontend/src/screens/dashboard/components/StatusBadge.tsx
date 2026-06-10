type StatusGizi = 'Normal' | 'Risiko Stunting' | 'Gizi Kurang' | 'Stunting' | 'Underweight' | 'Perlu Perhatian';

const BADGE_MAP: Record<StatusGizi, string> = {
  Normal: 'bg-emerald-100 text-emerald-700',
  'Risiko Stunting': 'bg-red-100 text-red-700',
  'Gizi Kurang': 'bg-amber-100 text-amber-700',
  Stunting: 'bg-red-200 text-red-800',
  Underweight: 'bg-orange-100 text-orange-700',
  'Perlu Perhatian': 'bg-amber-100 text-amber-700',
};

interface StatusBadgeProps {
  status: StatusGizi;
  size?: 'sm' | 'xs';
}

export function StatusBadge({ status, size = 'sm' }: StatusBadgeProps): JSX.Element {
  return (
    <span
      className={`font-bold rounded-full font-body whitespace-nowrap ${BADGE_MAP[status]} ${
        size === 'xs' ? 'text-[10px] px-2 py-0.5' : 'text-xs px-2.5 py-1'
      }`}
    >
      {status}
    </span>
  );
}
