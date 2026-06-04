import { Search, ChevronDown } from 'lucide-react';
import type { ReactNode } from 'react';

export interface FilterOption {
  label: string;
  value: string;
}

export interface FilterDropdown {
  id: string;
  placeholder: string;
  options: FilterOption[];
  value: string;
  onChange: (value: string) => void;
}

interface FilterBarProps {
  searchValue?: string;
  onSearchChange?: (value: string) => void;
  searchPlaceholder?: string;
  filters?: FilterDropdown[];
  actions?: ReactNode;
}

export function FilterBar({
  searchValue,
  onSearchChange,
  searchPlaceholder = 'Cari...',
  filters = [],
  actions,
}: FilterBarProps) {
  return (
    <div className="flex flex-col md:flex-row gap-4 justify-between items-stretch">
      {onSearchChange !== undefined && (
        <div className="relative flex-1">
          <Search size={18} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-neutral-400" />
          <input
            type="text"
            placeholder={searchPlaceholder}
            value={searchValue || ''}
            onChange={(e) => onSearchChange(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 bg-white border border-neutral-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-100 focus:border-primary transition-all font-body text-neutral-700"
          />
        </div>
      )}
      <div className="flex flex-wrap sm:flex-nowrap gap-3">
        {filters.map((filter) => (
          <div key={filter.id} className="relative flex-1 sm:w-48">
            <select
              value={filter.value}
              onChange={(e) => filter.onChange(e.target.value)}
              className="w-full appearance-none bg-white border border-neutral-200 rounded-xl px-4 py-2.5 pr-10 text-sm text-neutral-600 focus:outline-none focus:ring-2 focus:ring-primary-100 focus:border-primary cursor-pointer font-medium"
            >
              {filter.options.map((opt) => (
                <option key={opt.value} value={opt.value}>{opt.label}</option>
              ))}
            </select>
            <span className="absolute right-3.5 top-1/2 -translate-y-1/2 pointer-events-none text-neutral-400">
              <ChevronDown size={16} />
            </span>
          </div>
        ))}
        {actions && <div className="flex gap-2">{actions}</div>}
      </div>
    </div>
  );
}
