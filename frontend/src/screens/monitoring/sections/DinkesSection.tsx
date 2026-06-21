import { useState, useEffect, useCallback } from 'react';
import { Users, AlertTriangle, MapPin, TrendingUp, Search, Heart, ClipboardList } from 'lucide-react';
import { apiGet } from '../../../lib/api';
import { Paginator } from '../../../components/Paginator';

interface WilayahStats {
  nama_wilayah: string;
  total_balita: number;
  jumlah_kasus: number;
  prevalensi: number;
  level: string;
}

interface GlobalStats {
  total_pasien: number;
  kasus_stunting: number;
  total_balita: number;
}

interface IbuHamilStats {
  total_ibu_hamil: number;
  trimester_1: number;
  trimester_2: number;
  trimester_3: number;
  melahirkan: number;
  nifas: number;
  keguguran: number;
}

interface IbuHamilWilayah {
  nama_wilayah: string;
  total_ibu_hamil: number;
  trimester_1: number;
  trimester_2: number;
  trimester_3: number;
  melahirkan: number;
  nifas: number;
  keguguran: number;
}

export function DinkesSection() {
  const [search, setSearch] = useState('');
  const [wilayah, setWilayah] = useState<WilayahStats[]>([]);
  const [stats, setStats] = useState<GlobalStats | null>(null);
  const [ibuHamilStats, setIbuHamilStats] = useState<IbuHamilStats | null>(null);
  const [ibuHamilWilayah, setIbuHamilWilayah] = useState<IbuHamilWilayah[]>([]);
  const [ibuSearch, setIbuSearch] = useState('');
  const [ibuSortBy, setIbuSortBy] = useState<keyof IbuHamilWilayah>('total_ibu_hamil');
  const [ibuSortDir, setIbuSortDir] = useState<'asc' | 'desc'>('desc');
  const [sortBy, setSortBy] = useState<keyof WilayahStats>('prevalensi');
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('desc');

  useEffect(() => {
    apiGet<{ wilayah: WilayahStats[] }>('/dashboard/stunting-per-wilayah')
      .then((res) => setWilayah((res.wilayah ?? []).filter((w) => w.nama_wilayah !== 'Tidak Diketahui')))
      .catch(() => {});
    apiGet<GlobalStats>('/dashboard/stats')
      .then((res) => setStats(res))
      .catch(() => {});
    apiGet<IbuHamilStats>('/dashboard/ibu-hamil-stats')
      .then((res) => setIbuHamilStats(res))
      .catch(() => {});
    apiGet<{ wilayah: IbuHamilWilayah[] }>('/dashboard/ibu-hamil-per-wilayah')
      .then((res) => setIbuHamilWilayah((res.wilayah ?? []).filter((w) => w.nama_wilayah !== 'Tidak Diketahui')))
      .catch(() => {});
  }, []);

  const handleSort = (key: keyof WilayahStats) => {
    if (sortBy === key) setSortDir((d) => (d === 'desc' ? 'asc' : 'desc'));
    else { setSortBy(key); setSortDir('desc'); }
  };

  const filtered = wilayah
    .filter((w) => !search || w.nama_wilayah.toLowerCase().includes(search.toLowerCase()))
    .sort((a, b) => {
      const va = a[sortBy]; const vb = b[sortBy];
      if (typeof va === 'string' && typeof vb === 'string')
        return sortDir === 'asc' ? va.localeCompare(vb) : vb.localeCompare(va);
      return sortDir === 'asc' ? (va as number) - (vb as number) : (vb as number) - (va as number);
    });

  const barData = [...wilayah]
    .filter((w) => w.total_balita > 0)
    .sort((a, b) => b.prevalensi - a.prevalensi)
    .slice(0, 12)
    .map((w) => ({
      label: w.nama_wilayah.replace('Kabupaten ', 'Kab. '),
      value: w.prevalensi,
      color: w.level === 'tinggi' ? '#dc2626' : w.level === 'sedang' ? '#f97316' : '#22c55e',
    }));

  const ibuFiltered = ibuHamilWilayah
    .filter((w) => !ibuSearch || w.nama_wilayah.toLowerCase().includes(ibuSearch.toLowerCase()))
    .sort((a, b) => {
      const va = a[ibuSortBy]; const vb = b[ibuSortBy];
      if (typeof va === 'string' && typeof vb === 'string')
        return ibuSortDir === 'asc' ? va.localeCompare(vb) : vb.localeCompare(va);
      return ibuSortDir === 'asc' ? (va as number) - (vb as number) : (vb as number) - (va as number);
    });

  const ibuBarData = [...ibuHamilWilayah]
    .sort((a, b) => b.total_ibu_hamil - a.total_ibu_hamil)
    .slice(0, 10)
    .map((w) => ({
      label: w.nama_wilayah.replace('Kabupaten ', 'Kab. '),
      value: w.total_ibu_hamil,
    }));
  const maxIbuBar = ibuBarData.length > 0 ? Math.max(...ibuBarData.map((d) => d.value)) : 1;

  const handleIbuSort = (key: keyof IbuHamilWilayah) => {
    if (ibuSortBy === key) setIbuSortDir((d) => (d === 'desc' ? 'asc' : 'desc'));
    else { setIbuSortBy(key); setIbuSortDir('desc'); }
  };

  // Tabel Penimbangan
  interface PemeriksaanItem {
    id_hasil_pemeriksaan: number;
    id_jadwal_imunisasi: number;
    nama_vaksin: string;
    nama_pasien: string;
    berat_badan: number;
    tinggi_badan: number;
    lingkar_kepala: number;
    tekanan_darah: string;
    status_stunting: string;
    status_gizi: string;
    catatan?: string;
    tanggal: string;
    petugas: string;
  }
  const [pemeriksaan, setPemeriksaan] = useState<PemeriksaanItem[]>([]);
  const [pemPage, setPemPage] = useState(1);
  const [pemTotal, setPemTotal] = useState(0);
  const [pemLastPage, setPemLastPage] = useState(1);
  const PEM_PER_PAGE = 15;

  const fetchPemeriksaan = useCallback((p: number) => {
    apiGet<{ pemeriksaan: PemeriksaanItem[]; meta: { current_page: number; per_page: number; total: number; last_page: number } }>('/monitoring/semua-pemeriksaan', { page: String(p), per_page: String(PEM_PER_PAGE) })
      .then((res) => {
        setPemeriksaan(res.pemeriksaan ?? []);
        setPemTotal(res.meta?.total ?? 0);
        setPemLastPage(res.meta?.last_page ?? 1);
      })
      .catch(() => {});
  }, []);

  useEffect(() => { fetchPemeriksaan(pemPage); }, [pemPage, fetchPemeriksaan]);

  return (
    <div className="space-y-5">
      {/* ═══ SECTION: BALITA & STUNTING ═══ */}
      <div>
        <div className="flex items-center gap-2 mb-4">
          <AlertTriangle size={18} className="text-red-500" />
          <h3 className="text-base font-bold text-neutral-800 font-headline">Balita & Stunting</h3>
        </div>

        <div className="space-y-4">
          {/* Stats */}
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
            <div className="bg-primary rounded-2xl p-4 relative overflow-hidden">
              <div className="absolute -bottom-4 -right-4 w-20 h-20 rounded-full bg-white/10" />
              <Users size={20} className="text-white/80 mb-1.5" />
              <p className="text-white/70 text-[10px] font-semibold uppercase tracking-wide">Total Pasien</p>
              <p className="text-2xl font-bold text-white font-headline mt-0.5">{stats ? stats.total_pasien.toLocaleString('id-ID') : '...'}</p>
            </div>
            <div className="bg-white rounded-2xl p-4 border border-neutral-100">
              <AlertTriangle size={20} className="text-red-500 mb-1.5" />
              <p className="text-neutral-400 text-[10px] font-semibold uppercase tracking-wide">Kasus Stunting</p>
              <p className="text-2xl font-bold text-red-500 font-headline mt-0.5">{stats ? stats.kasus_stunting.toLocaleString('id-ID') : '...'}</p>
              {stats && stats.total_balita > 0 && (
                <p className="text-[10px] text-red-400 mt-0.5">{((stats.kasus_stunting / stats.total_balita) * 100).toFixed(1)}% dari balita</p>
              )}
            </div>
            <div className="bg-white rounded-2xl p-4 border border-neutral-100">
              <MapPin size={20} className="text-blue-500 mb-1.5" />
              <p className="text-neutral-400 text-[10px] font-semibold uppercase tracking-wide">Wilayah</p>
              <p className="text-2xl font-bold text-blue-500 font-headline mt-0.5">{filtered.filter((w) => w.total_balita > 0).length}</p>
              <p className="text-[10px] text-neutral-400 mt-0.5">dari 35 kab/kota</p>
            </div>
            <div className="bg-white rounded-2xl p-4 border border-neutral-100">
              <TrendingUp size={20} className="text-emerald-500 mb-1.5" />
              <p className="text-neutral-400 text-[10px] font-semibold uppercase tracking-wide">Prevalensi</p>
              <p className="text-2xl font-bold text-emerald-500 font-headline mt-0.5">
                {filtered.filter((w) => w.total_balita > 0).length > 0
                  ? (filtered.reduce((s, w) => s + w.prevalensi, 0) / filtered.filter((w) => w.total_balita > 0).length).toFixed(1) : '...'}%
              </p>
              <p className="text-[10px] text-neutral-400 mt-0.5">rata-rata</p>
            </div>
          </div>

          {/* Bar Chart + Ringkasan */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
            <div className="lg:col-span-2 bg-white rounded-2xl p-5 border border-neutral-100">
              <h4 className="text-sm font-bold text-neutral-800 font-headline mb-4">Prevalensi Stunting Per Wilayah</h4>
              {barData.length === 0 ? (
                <p className="text-sm text-neutral-400 text-center py-12">Belum ada data</p>
              ) : (
                <div className="space-y-2">
                  {barData.map((d) => (
                    <div key={d.label} className="flex items-center gap-2 text-xs">
                      <span className="w-24 text-right truncate text-neutral-600" title={d.label}>{d.label}</span>
                      <div className="flex-1 h-5 bg-neutral-100 rounded-full overflow-hidden">
                        <div className="h-full rounded-full transition-all duration-500" style={{ width: `${Math.min(d.value, 100)}%`, backgroundColor: d.color }} />
                      </div>
                      <span className="w-10 text-right font-semibold text-neutral-700">{d.value.toFixed(1)}%</span>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div className="bg-white rounded-2xl p-5 border border-neutral-100">
              <h4 className="text-sm font-bold text-neutral-800 font-headline mb-3">Ringkasan</h4>
              <div className="space-y-2.5">
                <Row label="Total Balita" value={filtered.reduce((s, w) => s + w.total_balita, 0).toLocaleString('id-ID')} />
                <Row label="Total Kasus" value={filtered.reduce((s, w) => s + w.jumlah_kasus, 0).toLocaleString('id-ID')} color="text-red-500" />
                <Row label="Tinggi (>20%)" value={String(wilayah.filter((w) => w.level === 'tinggi').length)} color="text-red-500" />
                <Row label="Sedang" value={String(wilayah.filter((w) => w.level === 'sedang').length)} color="text-orange-500" />
                <Row label="Rendah (<10%)" value={String(wilayah.filter((w) => w.level === 'rendah').length)} color="text-green-500" />
              </div>
            </div>
          </div>

          {/* Table */}
          <Table
            search={search} onSearch={setSearch}
            count={filtered.length}
            headers={[
              { label: 'KABUPATEN / KOTA', span: 2, key: 'nama_wilayah' as keyof WilayahStats, sortBy, sortDir, onSort: () => handleSort('nama_wilayah') },
              { label: 'BALITA', key: 'total_balita' as keyof WilayahStats, sortBy, sortDir, onSort: () => handleSort('total_balita') },
              { label: 'KASUS', key: 'jumlah_kasus' as keyof WilayahStats, sortBy, sortDir, onSort: () => handleSort('jumlah_kasus') },
              { label: '%', key: 'prevalensi' as keyof WilayahStats, sortBy, sortDir, onSort: () => handleSort('prevalensi') },
            ]}
            emptyMessage="Tidak ada data"
            rows={filtered.map((w) => ({
              key: w.nama_wilayah,
              bg: w.level === 'tinggi' ? 'bg-red-50/30' : w.level === 'sedang' ? 'bg-orange-50/30' : '',
              cells: [
                { span: 2, content: <><span className={`w-2 h-2 rounded-full shrink-0 ${w.level === 'tinggi' ? 'bg-red-500' : w.level === 'sedang' ? 'bg-orange-500' : 'bg-green-500'}`} /> <span className="text-sm font-semibold text-neutral-800 truncate">{w.nama_wilayah.replace('Kabupaten ', 'Kab. ')}</span></> },
                { content: <span className="text-sm text-neutral-700 tabular-nums">{w.total_balita.toLocaleString('id-ID')}</span> },
                { content: <span className={`text-sm font-semibold tabular-nums ${w.jumlah_kasus > 0 ? 'text-red-600' : 'text-neutral-400'}`}>{w.jumlah_kasus.toLocaleString('id-ID')}</span> },
                { content: <span className={`text-sm font-bold tabular-nums ${w.level === 'tinggi' ? 'text-red-600' : w.level === 'sedang' ? 'text-orange-600' : 'text-green-600'}`}>{w.prevalensi.toFixed(1)}%</span> },
              ],
            }))}
          />
        </div>
      </div>

      {/* ═══ SECTION: IBU HAMIL ═══ */}
      <div>
        <div className="flex items-center gap-2 mb-4 pt-2 border-t-2 border-neutral-100">
          <Heart size={18} className="text-pink-500" />
          <h3 className="text-base font-bold text-neutral-800 font-headline">Monitoring Ibu Hamil</h3>
        </div>

        <div className="space-y-4">
          {/* Stats */}
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
            <div className="bg-pink-500 rounded-2xl p-4 relative overflow-hidden">
              <div className="absolute -bottom-4 -right-4 w-20 h-20 rounded-full bg-white/10" />
              <Heart size={20} className="text-white/80 mb-1.5" />
              <p className="text-white/70 text-[10px] font-semibold uppercase tracking-wide">Total Ibu Hamil</p>
              <p className="text-2xl font-bold text-white font-headline mt-0.5">{ibuHamilStats ? ibuHamilStats.total_ibu_hamil.toLocaleString('id-ID') : '...'}</p>
            </div>
            <div className="bg-white rounded-2xl p-4 border border-neutral-100">
              <span className="absolute top-3 right-3 w-2 h-2 rounded-full bg-purple-500" />
              <p className="text-neutral-400 text-[10px] font-semibold uppercase tracking-wide">Trimester 1</p>
              <p className="text-2xl font-bold text-purple-500 font-headline mt-0.5">{ibuHamilStats ? ibuHamilStats.trimester_1.toLocaleString('id-ID') : '...'}</p>
            </div>
            <div className="bg-white rounded-2xl p-4 border border-neutral-100">
              <span className="absolute top-3 right-3 w-2 h-2 rounded-full bg-blue-500" />
              <p className="text-neutral-400 text-[10px] font-semibold uppercase tracking-wide">Trimester 2</p>
              <p className="text-2xl font-bold text-blue-500 font-headline mt-0.5">{ibuHamilStats ? ibuHamilStats.trimester_2.toLocaleString('id-ID') : '...'}</p>
            </div>
            <div className="bg-white rounded-2xl p-4 border border-neutral-100">
              <span className="absolute top-3 right-3 w-2 h-2 rounded-full bg-indigo-500" />
              <p className="text-neutral-400 text-[10px] font-semibold uppercase tracking-wide">Trimester 3</p>
              <p className="text-2xl font-bold text-indigo-500 font-headline mt-0.5">{ibuHamilStats ? ibuHamilStats.trimester_3.toLocaleString('id-ID') : '...'}</p>
            </div>
          </div>

          {/* Bar Chart + Ringkasan */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
            <div className="lg:col-span-2 bg-white rounded-2xl p-5 border border-neutral-100">
              <h4 className="text-sm font-bold text-neutral-800 font-headline mb-4">Ibu Hamil Per Wilayah</h4>
              {ibuBarData.length === 0 ? (
                <p className="text-sm text-neutral-400 text-center py-12">Belum ada data</p>
              ) : (
                <div className="space-y-2">
                  {ibuBarData.map((d) => (
                    <div key={d.label} className="flex items-center gap-2 text-xs">
                      <span className="w-24 text-right truncate text-neutral-600" title={d.label}>{d.label}</span>
                      <div className="flex-1 h-5 bg-neutral-100 rounded-full overflow-hidden">
                        <div className="h-full rounded-full bg-pink-500 transition-all duration-500" style={{ width: `${(d.value / maxIbuBar) * 100}%` }} />
                      </div>
                      <span className="w-12 text-right font-semibold text-neutral-700">{d.value.toLocaleString('id-ID')}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div className="bg-white rounded-2xl p-5 border border-neutral-100">
              <h4 className="text-sm font-bold text-neutral-800 font-headline mb-3">Ringkasan</h4>
              {ibuHamilStats ? (
                <div className="space-y-2.5">
                  <Row label="Total Ibu Hamil" value={ibuHamilStats.total_ibu_hamil.toLocaleString('id-ID')} color="text-pink-500" />
                  <Row label="Trimester 1" value={ibuHamilStats.trimester_1.toLocaleString('id-ID')} />
                  <Row label="Trimester 2" value={ibuHamilStats.trimester_2.toLocaleString('id-ID')} />
                  <Row label="Trimester 3" value={ibuHamilStats.trimester_3.toLocaleString('id-ID')} />
                  <Row label="Melahirkan" value={ibuHamilStats.melahirkan.toLocaleString('id-ID')} />
                  <Row label="Nifas" value={ibuHamilStats.nifas.toLocaleString('id-ID')} />
                  <Row label="Keguguran" value={ibuHamilStats.keguguran.toLocaleString('id-ID')} />
                </div>
              ) : (
                <p className="text-sm text-neutral-400 text-center py-4">Memuat...</p>
              )}
            </div>
          </div>

          {/* Table */}
          <Table
            search={ibuSearch} onSearch={setIbuSearch}
            count={ibuFiltered.length}
            headers={[
              { label: 'KABUPATEN / KOTA', span: 2, key: 'nama_wilayah' as keyof IbuHamilWilayah, sortBy: ibuSortBy, sortDir: ibuSortDir, onSort: () => handleIbuSort('nama_wilayah') },
              { label: 'TOTAL', key: 'total_ibu_hamil' as keyof IbuHamilWilayah, sortBy: ibuSortBy, sortDir: ibuSortDir, onSort: () => handleIbuSort('total_ibu_hamil') },
              { label: 'T1', key: 'trimester_1' as keyof IbuHamilWilayah, sortBy: ibuSortBy, sortDir: ibuSortDir, onSort: () => handleIbuSort('trimester_1') },
              { label: 'T2+T3', key: 'trimester_2' as keyof IbuHamilWilayah, sortBy: ibuSortBy, sortDir: ibuSortDir, onSort: () => handleIbuSort('trimester_2') },
            ]}
            emptyMessage="Tidak ada data"
            rows={ibuFiltered.map((w) => ({
              key: w.nama_wilayah,
              cells: [
                { span: 2, content: <span className="text-sm font-semibold text-neutral-800 truncate">{w.nama_wilayah.replace('Kabupaten ', 'Kab. ')}</span> },
                { content: <span className="text-sm font-semibold text-pink-600 tabular-nums">{w.total_ibu_hamil.toLocaleString('id-ID')}</span> },
                { content: <span className="text-sm text-purple-600 tabular-nums">{w.trimester_1.toLocaleString('id-ID')}</span> },
                { content: <span className="text-sm text-neutral-500 tabular-nums">{w.trimester_2 + w.trimester_3}</span> },
              ],
            }))}
          />
        </div>
      </div>

      {/* ═══ SECTION: TABEL PENIMBANGAN ═══ */}
      <div>
        <div className="flex items-center gap-2 mb-4 pt-2 border-t-2 border-neutral-100">
          <ClipboardList size={18} className="text-primary" />
          <h3 className="text-base font-bold text-neutral-800 font-headline">Tabel Penimbangan</h3>
        </div>

        <div className="bg-white rounded-2xl border border-neutral-100 overflow-hidden">
          <div className="grid grid-cols-12 px-5 py-2.5 bg-neutral-50 border-b border-neutral-100 text-[10px] font-bold text-neutral-400 uppercase tracking-wide">
            <p className="col-span-2">NAMA</p>
            <p>VAKSIN</p>
            <p>BB (kg)</p>
            <p>TB (cm)</p>
            <p>LK (cm)</p>
            <p>TD</p>
            <p>STUNTING</p>
            <p>GIZI</p>
            <p>CATATAN</p>
            <p>TANGGAL</p>
          </div>

          {pemeriksaan.length === 0 ? (
            <div className="py-12 text-center text-neutral-400 text-sm">Tidak ada data penimbangan</div>
          ) : (
            pemeriksaan.map((row) => (
              <div key={row.id_hasil_pemeriksaan} className="grid grid-cols-12 px-5 py-2.5 border-b border-neutral-50 last:border-0 hover:bg-neutral-50/60 items-center">
                <p className="col-span-2 text-sm font-semibold text-neutral-800 truncate pr-2">{row.nama_pasien}</p>
                <p className="text-xs text-blue-600 font-medium truncate pr-1">{row.nama_vaksin || '-'}</p>
                <p className="text-sm text-neutral-700 tabular-nums">{row.berat_badan?.toFixed(1)}</p>
                <p className="text-sm text-neutral-700 tabular-nums">{row.tinggi_badan?.toFixed(1)}</p>
                <p className="text-sm text-neutral-700 tabular-nums">{row.lingkar_kepala?.toFixed(1)}</p>
                <p className="text-xs text-neutral-500">{row.tekanan_darah || '-'}</p>
                <p className="text-xs">
                  <span className={`px-1.5 py-0.5 rounded-full ${row.status_stunting?.includes('Stunting') ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
                    {row.status_stunting || '-'}
                  </span>
                </p>
                <p className="text-xs">
                  <span className={`px-1.5 py-0.5 rounded-full ${row.status_gizi?.includes('Buruk') || row.status_gizi?.includes('Kurang') ? 'bg-amber-100 text-amber-700' : 'bg-emerald-100 text-emerald-700'}`}>
                    {row.status_gizi || '-'}
                  </span>
                </p>
                <p className="text-xs text-neutral-400 truncate">{row.catatan || '-'}</p>
                <p className="text-xs text-neutral-500">{row.tanggal?.slice(0, 10)}</p>
              </div>
            ))
          )}

          <Paginator page={pemPage} totalPages={pemLastPage} totalItems={pemTotal} pageSize={PEM_PER_PAGE} onPageChange={setPemPage} />
        </div>
      </div>
    </div>
  );
}

/* ── Reusable sub-components ── */

function Row({ label, value, color }: { label: string; value: string; color?: string }) {
  return (
    <div className="flex items-center justify-between text-xs">
      <span className="text-neutral-500">{label}</span>
      <span className={`font-semibold ${color || 'text-neutral-800'}`}>{value}</span>
    </div>
  );
}

interface HeaderDef<K extends string> {
  label: string;
  span?: number;
  key: K;
  sortBy: K;
  sortDir: 'asc' | 'desc';
  onSort: () => void;
}

interface CellDef {
  span?: number;
  content: React.ReactNode;
}

interface RowDef {
  key: string;
  bg?: string;
  cells: CellDef[];
}

function Table<K extends string>({
  search, onSearch, count, headers, emptyMessage, rows,
}: {
  search: string;
  onSearch: (v: string) => void;
  count: number;
  headers: HeaderDef<K>[];
  emptyMessage: string;
  rows: RowDef[];
}) {
  return (
    <div className="bg-white rounded-2xl border border-neutral-100 overflow-hidden">
      <div className="flex items-center gap-3 px-5 py-3 border-b border-neutral-50">
        <div className="relative flex-1 max-w-xs">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-400" />
          <input
            type="text" placeholder="Cari kabupaten/kota..." value={search}
            onChange={(e) => onSearch(e.target.value)}
            className="w-full pl-9 pr-4 py-1.5 bg-neutral-50 border border-neutral-200 rounded-lg text-xs text-neutral-700 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200"
          />
        </div>
        <p className="text-[10px] text-neutral-400 ml-auto">{count} wilayah</p>
      </div>

      <div
        className="grid px-5 py-2 bg-neutral-50 border-b border-neutral-100 text-[10px] font-bold text-neutral-400 uppercase tracking-wide"
        style={{ gridTemplateColumns: headers.map((h) => h.span ? `span ${h.span}` : '1fr').join(' ') }}
      >
        {(() => {
          const totalSpan = headers.reduce((s, h) => s + (h.span || 1), 0);
          const gridCols = `repeat(${totalSpan}, minmax(0, 1fr))`;
          return (
            <div className="grid w-full" style={{ gridTemplateColumns: gridCols, gridColumn: `1 / -1` }}>
              {headers.map((h) => (
                <p
                  key={h.key}
                  className={`text-right cursor-pointer select-none hover:text-neutral-600 ${h.key === headers[0].key ? 'text-left' : ''}`}
                  style={h.span ? { gridColumn: `span ${h.span}` } : undefined}
                  onClick={h.onSort}
                >
                  {h.label} {h.sortBy === h.key ? (h.sortDir === 'asc' ? '▲' : '▼') : ''}
                </p>
              ))}
            </div>
          );
        })()}
      </div>

      {rows.length === 0 ? (
        <div className="py-12 text-center text-neutral-400 text-xs">{emptyMessage}</div>
      ) : (
        rows.map((row) => (
          <div
            key={row.key}
            className={`grid px-5 py-2.5 border-b border-neutral-50 last:border-0 hover:bg-neutral-50/60 items-center ${row.bg || ''}`}
            style={{ gridTemplateColumns: row.cells.map((c) => c.span ? `span ${c.span}` : '1fr').join(' ') }}
          >
            {(() => {
              const totalSpan = row.cells.reduce((s, c) => s + (c.span || 1), 0);
              const gridCols = `repeat(${totalSpan}, minmax(0, 1fr))`;
              return (
                <div className="grid w-full items-center" style={{ gridTemplateColumns: gridCols, gridColumn: `1 / -1` }}>
                  {row.cells.map((cell, i) => (
                    <div
                      key={i}
                      className={`flex items-center gap-2 ${i === 0 ? 'justify-start' : 'justify-end'}`}
                      style={cell.span ? { gridColumn: `span ${cell.span}` } : undefined}
                    >
                      {cell.content}
                    </div>
                  ))}
                </div>
              );
            })()}
          </div>
        ))
      )}
    </div>
  );
}
