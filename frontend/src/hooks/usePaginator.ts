import { useState, useMemo, useCallback } from 'react';

interface UsePaginatorOptions {
  totalItems: number;
  pageSize?: number;
  initialPage?: number;
}

interface UsePaginatorReturn {
  page: number;
  setPage: (p: number) => void;
  totalPages: number;
  pageSize: number;
  offset: number;
  from: number;
  to: number;
}

export function usePaginator({ totalItems, pageSize = 10, initialPage = 1 }: UsePaginatorOptions): UsePaginatorReturn {
  const [page, setPage] = useState(initialPage);

  const totalPages = useMemo(() => Math.max(1, Math.ceil(totalItems / pageSize)), [totalItems, pageSize]);

  const safeSetPage = useCallback((p: number) => {
    setPage(Math.max(1, Math.min(p, totalPages)));
  }, [totalPages]);

  const offset = (page - 1) * pageSize;
  const from = totalItems === 0 ? 0 : offset + 1;
  const to = Math.min(page * pageSize, totalItems);

  return { page, setPage: safeSetPage, totalPages, pageSize, offset, from, to };
}
