import { TrendingUp, TrendingDown } from 'lucide-react';
import type { ReactNode } from 'react';

export interface StatCardProps {
  title: string;
  value: string | number;
  subtitle?: string;
  icon?: ReactNode;
  trend?: {
    direction: 'up' | 'down';
    value: string;
    label?: string;
  };
  variant?: 'gradient' | 'outline' | 'icon';
  color?: 'primary' | 'blue' | 'amber' | 'red' | 'emerald' | 'purple';
}

const COLOR_MAP = {
  primary: {
    gradient: 'from-primary to-[#0c6b49]',
    iconBg: 'bg-primary-50',
    iconText: 'text-primary',
    trendBg: 'bg-white/10 border-white/10',
    outlineBorder: 'border-primary-100',
  },
  blue: {
    gradient: 'from-blue-600 to-blue-800',
    iconBg: 'bg-blue-50',
    iconText: 'text-blue-600',
    trendBg: 'bg-white/10 border-white/10',
    outlineBorder: 'border-blue-100',
  },
  amber: {
    gradient: 'from-amber-500 to-amber-700',
    iconBg: 'bg-amber-50',
    iconText: 'text-amber-600',
    trendBg: 'bg-white/10 border-white/10',
    outlineBorder: 'border-amber-100',
  },
  red: {
    gradient: 'from-red-500 to-red-700',
    iconBg: 'bg-red-50',
    iconText: 'text-red-600',
    trendBg: 'bg-white/10 border-white/10',
    outlineBorder: 'border-red-100',
  },
  emerald: {
    gradient: 'from-emerald-500 to-emerald-700',
    iconBg: 'bg-emerald-50',
    iconText: 'text-emerald-600',
    trendBg: 'bg-white/10 border-white/10',
    outlineBorder: 'border-emerald-100',
  },
  purple: {
    gradient: 'from-purple-500 to-purple-700',
    iconBg: 'bg-purple-50',
    iconText: 'text-purple-600',
    trendBg: 'bg-white/10 border-white/10',
    outlineBorder: 'border-purple-100',
  },
};

export function StatCard({
  title,
  value,
  subtitle,
  icon,
  trend,
  variant = 'outline',
  color = 'primary',
}: StatCardProps) {
  const colors = COLOR_MAP[color];

  if (variant === 'gradient') {
    return (
      <div
        className={`bg-gradient-to-r ${colors.gradient} rounded-2xl p-6 text-white shadow-sm relative overflow-hidden flex flex-col justify-between min-h-[140px]`}
      >
        <div className="z-10">
          <span className="text-xs font-semibold tracking-wider opacity-85 uppercase font-headline">
            {title}
          </span>
          <div className="flex items-baseline mt-2 gap-2">
            <span className="text-4xl font-bold font-headline">{value}</span>
            {subtitle && <span className="text-sm opacity-90">{subtitle}</span>}
          </div>
        </div>
        {trend && (
          <div className="mt-4 z-10">
            <span
              className={`inline-flex items-center gap-1 ${colors.trendBg} px-2.5 py-1 rounded-full text-xs font-medium border`}
            >
              {trend.direction === 'up' ? (
                <TrendingUp size={14} className="text-emerald-300" />
              ) : (
                <TrendingDown size={14} className="text-emerald-300" />
              )}
              {trend.value}
              {trend.label && <span className="ml-0.5">{trend.label}</span>}
            </span>
          </div>
        )}
        <div className="absolute right-0 bottom-0 opacity-15 pointer-events-none">
          <svg width="180" height="100" viewBox="0 0 180 100" fill="none">
            <path d="M10 80 Q 45 10, 80 50 T 170 30" stroke="white" strokeWidth="4" fill="none" />
            <path d="M10 90 Q 45 40, 80 70 T 170 50" stroke="white" strokeWidth="2" fill="none" />
          </svg>
        </div>
      </div>
    );
  }

  if (variant === 'icon') {
    return (
      <div className="bg-white border border-neutral-100 rounded-2xl p-6 shadow-sm flex items-start gap-4 min-h-[140px]">
        {icon && (
          <div
            className={`w-12 h-12 rounded-xl ${colors.iconBg} ${colors.iconText} flex items-center justify-center flex-shrink-0`}
          >
            {icon}
          </div>
        )}
        <div className="flex flex-col justify-between flex-1 min-h-[92px]">
          <div>
            <span className="text-xs font-bold tracking-wider text-neutral-400 uppercase font-headline">
              {title}
            </span>
            <div className="text-3xl font-bold text-neutral-800 mt-1 font-headline">{value}</div>
          </div>
          {trend && (
            <div className="mt-2">
              <span
                className={`inline-flex items-center gap-1 text-xs font-medium ${
                  trend.direction === 'up' ? 'text-emerald-600' : 'text-red-500'
                }`}
              >
                {trend.direction === 'up' ? <TrendingUp size={14} /> : <TrendingDown size={14} />}
                {trend.value}
                {trend.label && <span className="text-neutral-400 ml-1">{trend.label}</span>}
              </span>
            </div>
          )}
          {subtitle && !trend && (
            <span className="text-xs text-neutral-500 mt-2">{subtitle}</span>
          )}
        </div>
      </div>
    );
  }

  // outline
  return (
    <div
      className={`bg-white border ${colors.outlineBorder} rounded-2xl p-6 shadow-sm flex flex-col justify-between min-h-[140px]`}
    >
      <div>
        <span className="text-xs font-bold tracking-wider text-neutral-400 uppercase font-headline">
          {title}
        </span>
        <div className="text-4xl font-bold text-neutral-800 mt-2 font-headline">{value}</div>
      </div>
      {trend && (
        <div className="mt-4">
          <span
            className={`inline-flex items-center gap-1 text-xs font-medium ${
              trend.direction === 'up' ? 'text-emerald-600' : 'text-red-500'
            }`}
          >
            {trend.direction === 'up' ? <TrendingUp size={14} /> : <TrendingDown size={14} />}
            {trend.value}
            {trend.label && <span className="text-neutral-400 ml-1">{trend.label}</span>}
          </span>
        </div>
      )}
      {subtitle && !trend && <span className="text-xs text-neutral-500 mt-3">{subtitle}</span>}
    </div>
  );
}
