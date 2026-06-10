import {
  BriefcaseMedical,
  CheckCircle2,
  CalendarClock,
  FlaskConical,
  FileText,
  Bell,
  AlertTriangle,
  ArrowRight,
} from 'lucide-react';
import type { NotifItem } from './types';

interface NotifCardProps {
  item: NotifItem;
  compact?: boolean;
}

const CATEGORY_STYLES: Record<
  string,
  { borderColor: string; iconBg: string; iconColor: string; Icon: React.FC<{ className?: string }> }
> = {
  urgent: {
    borderColor: 'border-l-red-500',
    iconBg: 'bg-red-50',
    iconColor: 'text-red-500',
    Icon: ({ className }) => <BriefcaseMedical className={className} />,
  },
  success: {
    borderColor: 'border-l-emerald-500',
    iconBg: 'bg-emerald-50',
    iconColor: 'text-emerald-600',
    Icon: ({ className }) => <CheckCircle2 className={className} />,
  },
  schedule: {
    borderColor: 'border-l-blue-500',
    iconBg: 'bg-blue-50',
    iconColor: 'text-blue-500',
    Icon: ({ className }) => <CalendarClock className={className} />,
  },
  lab: {
    borderColor: 'border-l-emerald-500',
    iconBg: 'bg-emerald-50',
    iconColor: 'text-emerald-600',
    Icon: ({ className }) => <FlaskConical className={className} />,
  },
  report: {
    borderColor: 'border-l-neutral-300',
    iconBg: 'bg-neutral-100',
    iconColor: 'text-neutral-400',
    Icon: ({ className }) => <FileText className={className} />,
  },
  info: {
    borderColor: 'border-l-blue-400',
    iconBg: 'bg-blue-50',
    iconColor: 'text-blue-500',
    Icon: ({ className }) => <Bell className={className} />,
  },
  warning: {
    borderColor: 'border-l-amber-400',
    iconBg: 'bg-amber-50',
    iconColor: 'text-amber-500',
    Icon: ({ className }) => <AlertTriangle className={className} />,
  },
};

export default function NotifCard({ item, compact = false }: NotifCardProps): JSX.Element {
  const style = CATEGORY_STYLES[item.category] ?? CATEGORY_STYLES.info;
  const { Icon } = style;

  if (compact) {
    return (
      <div className="bg-neutral-50 rounded-2xl p-4 flex gap-3 border border-neutral-100 hover:bg-white hover:shadow-sm transition-all">
        <div className={`w-10 h-10 rounded-full bg-white flex items-center justify-center shrink-0 shadow-sm ${style.iconColor}`}>
          <Icon className="w-4.5 h-4.5" />
        </div>
        <div className="flex-1 min-w-0">
          <p className="font-semibold text-neutral-800 text-sm truncate">{item.title}</p>
          <p className="text-xs text-neutral-400 mt-0.5">{item.time}</p>
        </div>
      </div>
    );
  }

  return (
    <div
      className={`bg-white rounded-2xl p-5 shadow-sm border-l-4 ${style.borderColor} border border-neutral-100 flex gap-4 hover:shadow-md transition-all`}
    >
      <div className={`w-12 h-12 rounded-2xl ${style.iconBg} ${style.iconColor} flex items-center justify-center shrink-0`}>
        <Icon className="w-6 h-6" />
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2 mb-1.5">
          <h4 className="font-bold text-neutral-800 text-[15px] leading-snug">{item.title}</h4>
          <span className="text-[11px] text-neutral-400 shrink-0 whitespace-nowrap">{item.time}</span>
        </div>
        <p className="text-sm text-neutral-600 leading-relaxed mb-3">{item.description}</p>
        <div className="flex items-center justify-between gap-2 flex-wrap">
          <div className="flex gap-2 flex-wrap">
            {item.tags.map((tag) => (
              <span
                key={tag.label}
                style={{ backgroundColor: tag.bg, color: tag.color }}
                className="text-[10px] font-bold px-2.5 py-1 rounded uppercase tracking-wide"
              >
                {tag.label}
              </span>
            ))}
          </div>
          {item.actionLabel && (
            <button className="flex items-center gap-1 text-xs font-semibold text-primary hover:text-primary-700 transition-colors shrink-0">
              {item.actionLabel} <ArrowRight className="w-3.5 h-3.5" />
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
