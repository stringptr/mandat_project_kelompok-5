import { useState } from 'react';
import { ChevronLeft, ChevronRight, Download } from 'lucide-react';

export interface Column<T> {
  header: string;
  accessor: keyof T | string;
  render?: (row: T) => React.ReactNode;
  className?: string;
}

interface DataTableProps<T> {
  title: string;
  columns: Column<T>[];
  data: T[];
  pageSize?: number;
  exportLabel?: string;
  onExport?: () => void;
  emptyMessage?: string;
}

export function DataTable<T extends Record<string, unknown>>({
  title,
  columns,
  data,
  pageSize = 5,
  exportLabel = 'Export ke Excel',
  onExport,
  emptyMessage = 'Tidak ada data',
}: DataTableProps<T>) {
  const [currentPage, setCurrentPage] = useState(1);
  const totalPages = Math.max(1, Math.ceil(data.length / pageSize));
  const startIdx = (currentPage - 1) * pageSize;
  const paginatedData = data.slice(startIdx, startIdx + pageSize);

  const getCellValue = (row: T, accessor: keyof T | string): React.ReactNode => {
    const keys = (accessor as string).split('.');
    let val: unknown = row;
    for (const k of keys) {
      val = (val as Record<string, unknown>)?.[k];
    }
    return val as React.ReactNode;
  };

  return (
    <div className="bg-white rounded-2xl border border-neutral-100 shadow-sm p-6 space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-base font-bold text-neutral-800 font-headline">{title}</h3>
        {onExport && (
          <button
            onClick={onExport}
            className="flex items-center gap-1.5 text-xs text-primary font-bold hover:text-primary-600 transition-colors"
          >
            <Download size={14} />
            {exportLabel}
          </button>
        )}
      </div>

      <div className="overflow-x-auto border border-neutral-100 rounded-xl">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-neutral-50/75 border-b border-neutral-100 text-neutral-400 text-xs font-bold tracking-wider">
              {columns.map((col, i) => (
                <th key={i} className={`px-5 py-4 ${col.className || ''}`}>
                  {col.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-neutral-100 text-sm text-neutral-700">
            {paginatedData.length === 0 ? (
              <tr>
                <td colSpan={columns.length} className="px-5 py-12 text-center text-neutral-400 font-medium">
                  {emptyMessage}
                </td>
              </tr>
            ) : (
              paginatedData.map((row, rowIdx) => (
                <tr key={rowIdx} className="hover:bg-neutral-50/40 transition-colors">
                  {columns.map((col, colIdx) => (
                    <td key={colIdx} className={`px-5 py-4 ${col.className || ''}`}>
                      {col.render ? col.render(row) : getCellValue(row, col.accessor)}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {data.length > pageSize && (
        <div className="flex flex-col sm:flex-row justify-between items-center pt-2 gap-4 text-sm text-neutral-500">
          <span>
            Menampilkan {startIdx + 1}-{Math.min(startIdx + pageSize, data.length)} dari {data.length} data
          </span>
          <div className="flex items-center gap-1">
            <button
              onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
              disabled={currentPage === 1}
              className="flex items-center justify-center p-2 rounded-lg border border-neutral-200 text-neutral-400 hover:bg-neutral-50 transition-colors cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed"
            >
              <ChevronLeft size={16} />
            </button>
            {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
              <button
                key={page}
                onClick={() => setCurrentPage(page)}
                className={`w-8 h-8 rounded-lg flex items-center justify-center font-bold text-xs transition-colors ${
                  page === currentPage
                    ? 'bg-primary text-white'
                    : 'border border-neutral-200 text-neutral-600 hover:bg-neutral-50 cursor-pointer'
                }`}
              >
                {page}
              </button>
            ))}
            <button
              onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
              disabled={currentPage === totalPages}
              className="flex items-center justify-center p-2 rounded-lg border border-neutral-200 text-neutral-400 hover:bg-neutral-50 transition-colors cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed"
            >
              <ChevronRight size={16} />
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
