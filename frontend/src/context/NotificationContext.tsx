import { createContext, useContext, useState, useCallback, useRef } from 'react';
import type { ApiResponse, ApiErrorItem } from '../types/api';

// ── Types ──────────────────────────────────────────────────────────────

export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface Toast {
  id: string;
  type: ToastType;
  title: string;
  message: string;
  errors?: ApiErrorItem[];
  createdAt: number;
}

interface NotificationContextValue {
  toasts: Toast[];
  notify: {
    success: (message: string, title?: string) => void;
    error: (message: string, title?: string) => void;
    warn: (message: string, title?: string) => void;
    info: (message: string, title?: string) => void;
    apiError: (err: unknown, fallbackMessage?: string) => void;
    dismiss: (id: string) => void;
  };
}

// ── Helpers ─────────────────────────────────────────────────────────────

const DEFAULT_DURATION = 5000;

const STATUS_CONFIG: Record<number, { type: ToastType; defaultTitle: string; defaultMessage: string }> = {
  200: { type: 'success', defaultTitle: 'Berhasil', defaultMessage: '' },
  201: { type: 'success', defaultTitle: 'Berhasil', defaultMessage: 'Data berhasil disimpan.' },
  400: { type: 'error', defaultTitle: 'Permintaan Tidak Valid', defaultMessage: 'Permintaan tidak dapat diproses.' },
  401: { type: 'error', defaultTitle: 'Sesi Berakhir', defaultMessage: 'Sesi Anda telah berakhir, silakan login kembali.' },
  403: { type: 'error', defaultTitle: 'Akses Ditolak', defaultMessage: 'Anda tidak memiliki izin untuk mengakses halaman/fitur ini.' },
  404: { type: 'error', defaultTitle: 'Data Tidak Ditemukan', defaultMessage: 'Data tidak dapat ditemukan.' },
  409: { type: 'error', defaultTitle: 'Konflik Data', defaultMessage: '' },
  422: { type: 'error', defaultTitle: 'Validasi Gagal', defaultMessage: 'Mohon periksa kembali data yang Anda masukkan.' },
  500: { type: 'error', defaultTitle: 'Kesalahan Sistem', defaultMessage: 'Terjadi kesalahan pada sistem. Mohon dicoba beberapa saat lagi.' },
  503: { type: 'error', defaultTitle: 'Layanan Tidak Tersedia', defaultMessage: 'Terjadi kesalahan pada sistem. Mohon dicoba beberapa saat lagi.' },
};

function generateId(): string {
  return `toast-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
}

function parseApiError(err: unknown): { status: number; title: string; detail: string; errors: ApiErrorItem[] } {
  // If it's an Axios error with response
  const axiosErr = err as { response?: { data?: Record<string, unknown> } } | undefined | null;
  const responseData = axiosErr?.response?.data;

  if (responseData && typeof responseData.status === 'number') {
    return {
      status: responseData.status as number,
      title: (responseData.title as string) ?? '',
      detail: (responseData.detail as string) ?? '',
      errors: (Array.isArray(responseData.errors) ? responseData.errors : []) as ApiErrorItem[],
    };
  }

  // If it's a raw fetch Response
  if (err instanceof Response) {
    // We can't read body here easily — caller should handle that
    return {
      status: err.status,
      title: '',
      detail: err.statusText,
      errors: [],
    };
  }

  // If it's already an ApiResponse-like object
  const apiObj = err as ApiResponse | undefined | null;
  if (apiObj && typeof apiObj.status === 'number') {
    return {
      status: apiObj.status,
      title: apiObj.title ?? '',
      detail: apiObj.detail ?? '',
      errors: Array.isArray(apiObj.errors) ? apiObj.errors : [],
    };
  }

  // Fallback: treat as network/unknown error
  return {
    status: 0,
    title: 'Kesalahan Jaringan',
    detail: err instanceof Error ? err.message : 'Terjadi kesalahan yang tidak terduga.',
    errors: [],
  };
}

// ── Context ─────────────────────────────────────────────────────────────

export const NotificationContext = createContext<NotificationContextValue | null>(null);

export function NotificationProvider({ children }: { children: React.ReactNode }): JSX.Element {
  const [toasts, setToasts] = useState<Toast[]>([]);
  const timersRef = useRef<Map<string, ReturnType<typeof setTimeout>>>(new Map());

  const scheduleDismiss = useCallback((id: string) => {
    const timer = setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
      timersRef.current.delete(id);
    }, DEFAULT_DURATION);
    timersRef.current.set(id, timer);
  }, []);

  const addToast = useCallback(
    (type: ToastType, title: string, message: string, errors?: ApiErrorItem[]) => {
      const id = generateId();
      const toast: Toast = { id, type, title, message, errors, createdAt: Date.now() };
      setToasts((prev) => [...prev.slice(-4), toast]);
      scheduleDismiss(id);
      return id;
    },
    [scheduleDismiss],
  );

  const dismiss = useCallback((id: string) => {
    const timer = timersRef.current.get(id);
    if (timer) {
      clearTimeout(timer);
      timersRef.current.delete(id);
    }
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const success = useCallback(
    (message: string, title?: string) => {
      addToast('success', title ?? 'Berhasil', message);
    },
    [addToast],
  );

  const error = useCallback(
    (message: string, title?: string) => {
      addToast('error', title ?? 'Kesalahan', message);
    },
    [addToast],
  );

  const warn = useCallback(
    (message: string, title?: string) => {
      addToast('warning', title ?? 'Perhatian', message);
    },
    [addToast],
  );

  const info = useCallback(
    (message: string, title?: string) => {
      addToast('info', title ?? 'Informasi', message);
    },
    [addToast],
  );

  const apiError = useCallback(
    (err: unknown, fallbackMessage?: string) => {
      const parsed = parseApiError(err);
      const status = parsed.status;

      let type: ToastType = 'error';
      let title = parsed.title;
      let message = parsed.detail;
      const errors: ApiErrorItem[] | undefined =
        parsed.errors.length > 0 ? parsed.errors : undefined;

      // If no title/message from API, use defaults
      if (STATUS_CONFIG[status]) {
        const cfg = STATUS_CONFIG[status];
        type = cfg.type;
        title = title || cfg.defaultTitle;
        message = message || cfg.defaultMessage;
      }

      // For status 0 (network error) and unknown, use fallback
      if (status === 0 || !STATUS_CONFIG[status]) {
        title = title || 'Kesalahan Jaringan';
        message = message || fallbackMessage || 'Terjadi kesalahan pada sistem. Mohon dicoba beberapa saat lagi.';
        type = 'error';
      }

      addToast(type, title, message, errors);
    },
    [addToast],
  );

  return (
    <NotificationContext.Provider value={{ toasts, notify: { success, error, warn, info, apiError, dismiss } }}>
      {children}
    </NotificationContext.Provider>
  );
}

export function useNotification(): NotificationContextValue['notify'] {
  const ctx = useContext(NotificationContext);
  if (!ctx) throw new Error('useNotification must be used inside NotificationProvider');
  return ctx.notify;
}
