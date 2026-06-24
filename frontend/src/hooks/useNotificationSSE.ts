import { useEffect, useRef, useCallback } from 'react';

export interface SSENotification {
  id: number;
  user_id: number;
  judul: string;
  pesan: string;
  tipe: string;
  created_at: string;
}

export function useNotificationSSE(
  enabled: boolean,
  onNotification: (notif: SSENotification) => void,
  onUnreadCount: (count: number) => void
) {
  const esRef = useRef<EventSource | null>(null);
  const retryRef = useRef<ReturnType<typeof setTimeout>>();
  const onNotifRef = useRef(onNotification);
  const onCountRef = useRef(onUnreadCount);
  onNotifRef.current = onNotification;
  onCountRef.current = onUnreadCount;

  const connect = useCallback(() => {
    if (!enabled) return;
    esRef.current?.close();

    const es = new EventSource('/api/v1/sse/notification', { withCredentials: true });

    es.addEventListener('notification', (e) => {
      try {
        onNotifRef.current(JSON.parse(e.data));
      } catch { /* ignore */ }
    });
    es.addEventListener('unread_count', (e) => {
      try {
        const { unread_count } = JSON.parse(e.data);
        onCountRef.current(unread_count);
      } catch { /* ignore */ }
    });
    es.addEventListener('connected', () => {});
    es.onerror = () => {
      es.close();
      retryRef.current = setTimeout(connect, 3000);
    };
    esRef.current = es;
  }, [enabled]);

  useEffect(() => {
    connect();
    return () => {
      clearTimeout(retryRef.current);
      esRef.current?.close();
    };
  }, [connect]);
}
