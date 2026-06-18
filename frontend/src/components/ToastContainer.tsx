import { useContext } from 'react';
import { X, CheckCircle, AlertTriangle, AlertCircle, Info } from 'lucide-react';
import { NotificationContext } from '../context/NotificationContext';
import type { Toast, ToastType } from '../context/NotificationContext';

// ── Style maps ──────────────────────────────────────────────────────────

const TYPE_STYLES: Record<
  ToastType,
  { container: string; border: string; icon: string; title: string; message: string; dismiss: string }
> = {
  success: {
    container: 'bg-green-50',
    border: 'border-green-200',
    icon: 'text-green-500',
    title: 'text-green-800',
    message: 'text-green-700',
    dismiss: 'text-green-400 hover:text-green-600 hover:bg-green-100',
  },
  error: {
    container: 'bg-red-50',
    border: 'border-red-200',
    icon: 'text-red-500',
    title: 'text-red-800',
    message: 'text-red-700',
    dismiss: 'text-red-400 hover:text-red-600 hover:bg-red-100',
  },
  warning: {
    container: 'bg-yellow-50',
    border: 'border-yellow-200',
    icon: 'text-yellow-500',
    title: 'text-yellow-800',
    message: 'text-yellow-700',
    dismiss: 'text-yellow-400 hover:text-yellow-600 hover:bg-yellow-100',
  },
  info: {
    container: 'bg-blue-50',
    border: 'border-blue-200',
    icon: 'text-blue-500',
    title: 'text-blue-800',
    message: 'text-blue-700',
    dismiss: 'text-blue-400 hover:text-blue-600 hover:bg-blue-100',
  },
};

const TYPE_ICONS: Record<ToastType, React.ReactNode> = {
  success: <CheckCircle size={20} />,
  error: <AlertCircle size={20} />,
  warning: <AlertTriangle size={20} />,
  info: <Info size={20} />,
};

// ── Single Toast ────────────────────────────────────────────────────────

function ToastItem({ toast, onDismiss }: { toast: Toast; onDismiss: (id: string) => void }): JSX.Element {
  const styles = TYPE_STYLES[toast.type];
  const icon = TYPE_ICONS[toast.type];

  return (
    <div
      role="alert"
      className={`${styles.container} ${styles.border} border rounded-2xl shadow-lg px-5 py-4 w-full max-w-md pointer-events-auto animate-slide-in font-body`}
      style={{ animation: 'slideIn 0.3s ease-out' }}
    >
      <div className="flex items-start gap-3">
        {/* Icon */}
        <span className={`flex-shrink-0 mt-0.5 ${styles.icon}`}>{icon}</span>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2">
            <p className={`text-sm font-bold ${styles.title} font-headline`}>{toast.title}</p>
            <button
              onClick={() => onDismiss(toast.id)}
              className={`flex-shrink-0 p-0.5 rounded-lg transition-colors ${styles.dismiss}`}
              aria-label="Tutup"
            >
              <X size={16} />
            </button>
          </div>

          {toast.message && (
            <p className={`text-xs mt-1 leading-relaxed ${styles.message}`}>{toast.message}</p>
          )}

          {/* 422 validation error breakdown */}
          {toast.errors && toast.errors.length > 0 && (
            <ul className="mt-2 space-y-1">
              {toast.errors.filter(Boolean).map((e, i) => (
                <li key={i} className={`text-[11px] leading-relaxed flex items-start gap-1.5 ${styles.message}`}>
                  <span className="flex-shrink-0 mt-0.5">·</span>
                  <span>
                    {e.location ? `${e.location}: ` : ''}{e.message}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
}

// ── Container ───────────────────────────────────────────────────────────

export default function ToastContainer(): JSX.Element {
  const ctx = useContext(NotificationContext);
  if (!ctx) return <></>;

  const { toasts, notify } = ctx;

  if (toasts.length === 0) return <></>;

  return (
    <div className="fixed top-4 right-4 z-[9999] flex flex-col gap-3 pointer-events-none">
      {toasts.map((toast) => (
        <ToastItem key={toast.id} toast={toast} onDismiss={notify.dismiss} />
      ))}
    </div>
  );
}
