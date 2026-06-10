interface StatCardProps {
  label: string;
  value: string | number;
  delta?: string;
  deltaPositif?: boolean;
  variant?: 'primary' | 'danger' | 'blue' | 'white';
  icon?: React.ReactNode;
  badge?: string;
  sub?: React.ReactNode;
}

export function StatCard({
  label,
  value,
  delta,
  deltaPositif = true,
  variant = 'white',
  icon,
  badge,
  sub,
}: StatCardProps): JSX.Element {
  const bg: Record<string, string> = {
    primary: 'bg-primary text-white',
    danger: 'bg-white border border-neutral-100',
    blue: 'bg-blue-600 text-white',
    white: 'bg-white border border-neutral-100',
  };

  const labelColor = variant === 'primary' || variant === 'blue' ? 'text-white/70' : 'text-neutral-500';
  const valueColor = variant === 'primary' || variant === 'blue'
    ? 'text-white'
    : variant === 'danger'
    ? 'text-red-500'
    : 'text-neutral-900';

  return (
    <div className={`rounded-2xl p-5 relative overflow-hidden ${bg[variant]}`}>
      {badge && (
        <span className="absolute top-4 right-4 text-[10px] font-bold uppercase tracking-wide bg-red-500 text-white px-2 py-0.5 rounded-full">
          {badge}
        </span>
      )}
      {icon && (
        <div className={`mb-3 ${variant === 'primary' || variant === 'blue' ? 'text-white/80' : 'text-primary'}`}>
          {icon}
        </div>
      )}
      <p className={`text-xs font-semibold uppercase tracking-wide mb-1 font-body ${labelColor}`}>{label}</p>
      <p className={`text-4xl font-bold font-headline leading-none mb-2 ${valueColor}`}>{value}</p>
      {delta && (
        <p className={`text-xs font-medium font-body flex items-center gap-1 ${
          variant === 'primary' || variant === 'blue' ? 'text-white/70' : deltaPositif ? 'text-emerald-600' : 'text-red-500'
        }`}>
          <span>{deltaPositif ? '↑' : '↓'}</span>
          {delta}
        </p>
      )}
      {sub}

      {/* decorative circle for primary/blue */}
      {(variant === 'primary' || variant === 'blue') && (
        <>
          <div className="absolute -bottom-6 -right-6 w-24 h-24 rounded-full bg-white/10" />
          <div className="absolute -bottom-10 right-8 w-16 h-16 rounded-full bg-white/10" />
        </>
      )}
    </div>
  );
}
