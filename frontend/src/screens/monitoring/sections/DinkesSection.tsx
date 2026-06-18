import { useState, useEffect } from 'react';
import { Search } from 'lucide-react';
import { apiGet } from '../../../lib/api';
import type { PasienListData } from '../../../types/entities';

interface PasienRow {
  no: number;
  id: string;
  nama: string;
  umur: string;
  nik: string;
}

export function DinkesSection() {
  const [search, setSearch] = useState('');
  const [data, setData] = useState<PasienRow[]>([]);

  useEffect(() => {
    apiGet<PasienListData>('/monitoring/pasien?page=1&per_page=100')
      .then((res) => {
        const list = (res.pasien ?? []).map((p, idx) => ({
          no: idx + 1,
          id: String(p.id_pasien),
          nama: p.nama,
          umur: p.umur,
          nik: p.nik,
        }));
        setData(list);
      })
      .catch(() => console.error('Gagal memuat data pasien'));
  }, []);

  const filteredData = data.filter(row =>
    !search ||
    row.nama.toLowerCase().includes(search.toLowerCase()) ||
    row.nik.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-2xl border border-neutral-100 overflow-hidden">
        <div className="flex items-center gap-3 px-5 py-4 border-b border-neutral-50">
          <div className="relative flex-1 max-w-xs">
            <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-400" />
            <input
              type="text"
              placeholder="Cari nama atau NIK..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-9 pr-4 py-2 bg-neutral-50 border border-neutral-200 rounded-xl text-sm text-neutral-700 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200"
            />
          </div>
          <p className="text-xs text-neutral-400 ml-auto">{filteredData.length} pasien</p>
        </div>

        <div className="grid grid-cols-5 px-5 py-2.5 bg-neutral-50 border-b border-neutral-100 text-[10px] font-bold text-neutral-400 uppercase tracking-wide">
          <p>NO</p>
          <p className="col-span-2">NAMA</p>
          <p>UMUR</p>
          <p>NIK</p>
        </div>

        {filteredData.length === 0 ? (
          <div className="py-12 text-center text-neutral-400">
            <p className="text-sm">Tidak ada data pasien</p>
          </div>
        ) : (
          filteredData.map((row) => (
            <div key={row.id} className="grid grid-cols-5 px-5 py-3 border-b border-neutral-50 last:border-0 hover:bg-neutral-50/60 items-center">
              <p className="text-sm text-neutral-400">{row.no}</p>
              <p className="col-span-2 text-sm font-bold text-primary">{row.nama}</p>
              <p className="text-sm text-neutral-600">{row.umur}</p>
              <p className="text-sm text-neutral-500">{row.nik}</p>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
