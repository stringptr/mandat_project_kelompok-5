import { useNavigate } from 'react-router-dom';
import { Bell, AlertTriangle, CalendarClock, CheckCircle2, FileText, FlaskConical } from 'lucide-react';
import { apiPatch } from '../../../lib/api';
import type { NotifItem } from './types';

interface NotifCardProps {
  item: NotifItem;
  onMarkRead?: (id: string) => void;
}

const CATEGORY_STYLES: Record<string, { dot: string; color: string; Icon: React.FC<{ className?: string }> }> = {
  urgent:   { dot: 'bg-red-500',   color: 'text-red-600',   Icon: ({ className }) => <AlertTriangle className={className} /> },
  success:  { dot: 'bg-emerald-500', color: 'text-emerald-600', Icon: ({ className }) => <CheckCircle2 className={className} /> },
  schedule: { dot: 'bg-blue-500',  color: 'text-blue-600',  Icon: ({ className }) => <CalendarClock className={className} /> },
  lab:      { dot: 'bg-emerald-500', color: 'text-emerald-600', Icon: ({ className }) => <FlaskConical className={className} /> },
  report:   { dot: 'bg-neutral-400', color: 'text-neutral-500', Icon: ({ className }) => <FileText className={className} /> },
  info:     { dot: 'bg-blue-400',  color: 'text-blue-500',  Icon: ({ className }) => <Bell className={className} /> },
  warning:  { dot: 'bg-amber-400', color: 'text-amber-500', Icon: ({ className }) => <AlertTriangle className={className} /> },
};

export default function NotifCard({ item, onMarkRead }: NotifCardProps): JSX.Element {
  const navigate = useNavigate();
  const style = CATEGORY_STYLES[item.category] ?? CATEGORY_STYLES.info;
  const { Icon } = style;

  const handleClick = async () => {
    const idNotifikasi = parseInt(item.id.replace('n-', ''), 10);
    if (!isNaN(idNotifikasi)) {
      apiPatch(`/notifikasi/${idNotifikasi}/read`).catch(() => {});
      onMarkRead?.(item.id);
    }
    if (item.actionUrl) {
      navigate(item.actionUrl);
    }
  };

  const card = (
    <div className={`flex items-start gap-3 px-4 py-3 rounded-xl transition-colors ${item.read ? '' : 'bg-blue-50/50'} cursor-pointer hover:bg-neutral-100`}>
      <div className={`w-8 h-8 rounded-full flex items-center justify-center shrink-0 mt-0.5 ${item.read ? 'bg-neutral-100' : 'bg-white shadow-sm'}`}>
        <Icon className={`w-4 h-4 ${style.color}`} />
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2">
          <p className={`text-sm leading-snug ${item.read ? 'font-medium text-neutral-600' : 'font-semibold text-neutral-800'}`}>
            {item.title}
          </p>
          <span className="text-[11px] text-neutral-400 shrink-0 whitespace-nowrap mt-0.5">{item.time}</span>
        </div>
        {item.description && (
          <p className={`text-xs mt-0.5 line-clamp-2 ${item.read ? 'text-neutral-400' : 'text-neutral-500'}`}>
            {item.description}
          </p>
        )}
        <div className="flex gap-1.5 mt-1.5 flex-wrap items-center">
          {item.tags.map((tag) => (
            <span
              key={tag.label}
              className="text-[10px] font-medium px-1.5 py-0.5 rounded"
              style={{ backgroundColor: tag.bg, color: tag.color }}
            >
              {tag.label}
            </span>
          ))}
          <span className="text-[10px] text-primary font-medium ml-auto">Lihat →</span>
        </div>
      </div>
    </div>
  );

  return (
    <div onClick={handleClick} className="group">
      {card}
    </div>
  );
}
