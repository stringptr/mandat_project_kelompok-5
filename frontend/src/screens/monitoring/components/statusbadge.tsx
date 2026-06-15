export type BadgeVariant =
  | 'gizi-baik'
  | 'gizi-kurang'
  | 'stunting'
  | 'gizi-lebih'
  | 'pending'
  | 'verified'
  | 'rejected'
  | 'urgent'
  | 'normal'
  | 'info';

interface StatusBadgeProps {
  variant: BadgeVariant;
  label?: string;
  dot?: boolean;
}

const BADGE_STYLES: Record<BadgeVariant, { bg: string; text: string; defaultLabel: string }> = {
  'gizi-baik': { bg: 'bg-emerald-50', text: 'text-emerald-700', defaultLabel: 'Gizi Baik' },
  'gizi-kurang': { bg: 'bg-amber-50', text: 'text-amber-700', defaultLabel: 'Gizi Kurang' },
  'stunting': { bg: 'bg-red-50', text: 'text-red-700', defaultLabel: 'Stunting' },
  'gizi-lebih': { bg: 'bg-blue-50', text: 'text-blue-700', defaultLabel: 'Gizi Lebih' },
  'pending': { bg: 'bg-orange-50', text: 'text-orange-700', defaultLabel: 'Pending' },
  'verified': { bg: 'bg-emerald-50', text: 'text-emerald-700', defaultLabel: 'Verified' },
  'rejected': { bg: 'bg-red-50', text: 'text-red-700', defaultLabel: 'Ditolak' },
  'urgent': { bg: 'bg-red-50', text: 'text-red-700', defaultLabel: 'Urgent' },
  'normal': { bg: 'bg-neutral-100', text: 'text-neutral-600', defaultLabel: 'Normal' },
  'info': { bg: 'bg-blue-50', text: 'text-blue-700', defaultLabel: 'Info' },
};

const DOT_COLORS: Record<BadgeVariant, string> = {
  'gizi-baik': 'bg-emerald-500',
  'gizi-kurang': 'bg-amber-500',
  'stunting': 'bg-red-500',
  'gizi-lebih': 'bg-blue-500',
  'pending': 'bg-orange-500',
  'verified': 'bg-emerald-500',
  'rejected': 'bg-red-500',
  'urgent': 'bg-red-500',
  'normal': 'bg-neutral-400',
  'info': 'bg-blue-500',
};

export function StatusBadge({ variant, label, dot = false }: StatusBadgeProps) {
  const style = BADGE_STYLES[variant];
  const displayLabel = label || style.defaultLabel;

  return (
    <span
      className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold ${style.bg} ${style.text}`}
    >
      {dot && <span className={`w-1.5 h-1.5 rounded-full ${DOT_COLORS[variant]}`} />}
      {displayLabel}
    </span>
  );
}
