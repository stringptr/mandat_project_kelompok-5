export type NotifCategory = 'urgent' | 'success' | 'info' | 'schedule' | 'report' | 'lab';

export interface NotifTag {
  label: string;
  color: string;
  bg: string;
}

export interface NotifItem {
  id: string;
  title: string;
  description: string;
  time: string;
  category: NotifCategory;
  tags: NotifTag[];
  actionLabel?: string;
  actionUrl?: string;
  read?: boolean;
}

export interface NotifGroup {
  groupLabel: string;
  items: NotifItem[];
}
