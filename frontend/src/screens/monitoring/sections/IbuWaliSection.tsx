import { useEffect, useState } from 'react';
import { Baby, Weight, Heart, Calendar, ClipboardList, Loader2, AlertCircle, Syringe, TrendingUp } from 'lucide-react';
import { useAuth } from '../../../context/AuthContext';
import { apiGet } from '../../../lib/api';
import type { PasienDetail, PasienRiwayatItem } from '../../../types/entities';

interface ImunisasiItem {
  nama_vaksin: string;
  tanggal_jadwal: string;
  status_imunisasi: string;
}

interface TumbuhKembangItem {
  bulan: string;
  berat_badan: number;
  tinggi_badan: number;
}

export function IbuWaliSection() {
  const { user } = useAuth();
  const idPasien = user?.idUser;

  const [pasien, setPasien] = useState<PasienDetail | null>(null);
  const [riwayat, setRiwayat] = useState<PasienRiwayatItem[]>([]);
  const [imunisasi, setImunisasi] = useState<ImunisasiItem[]>([]);
  const [tumbuhKembang, setTumbuhKembang] = useState<TumbuhKembangItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!idPasien) { setLoading(false); return; }
    setLoading(true);

    Promise.all([
      apiGet<PasienDetail>(`/monitoring/pasien/${idPasien}`).catch(() => null),
      apiGet<{ riwayat: PasienRiwayatItem[] }>(`/monitoring/pasien/${idPasien}/riwayat-pemeriksaan`).catch(() => ({ riwayat: [] as PasienRiwayatItem[] })),
      apiGet<{ jadwal: ImunisasiItem[] }>(`/imunisasi/pasien/${idPasien}?page=1&per_page=10`).catch(() => ({ jadwal: [] as ImunisasiItem[] })),
      apiGet<{ data: TumbuhKembangItem[] }>(`/monitoring/pasien/${idPasien}/tumbuh-kembang`).catch(() => ({ data: [] as TumbuhKembangItem[] })),
    ]).then(([p, r, im, tk]) => {
      if (p) setPasien(p as PasienDetail);
      setRiwayat(r.riwayat);
      setImunisasi(im.jadwal ?? []);
      setTumbuhKembang(tk.data ?? []);
    }).finally(() => setLoading(false));
  }, [idPasien]);

  const namaAnak = pasien?.data_anak?.nama_anak || pasien?.nama || 'Anak';
  const jenisKelamin = pasien?.jenis_kelamin || '';
  const labelKelamin = jenisKelamin === 'L' || jenisKelamin === 'Laki-laki' ? 'Laki-laki' : jenisKelamin === 'P' || jenisKelamin === 'Perempuan' ? 'Perempuan' : '';

  const usia = pasien?.tanggal_lahir
    ? (() => {
        const lahir = new Date(pasien.tanggal_lahir);
        const now = new Date();
        const bulan = (now.getFullYear() - lahir.getFullYear()) * 12 + (now.getMonth() - lahir.getMonth());
        if (bulan < 24) return `${bulan} bulan`;
        const tahun = Math.floor(bulan / 12);
        const sisaBulan = bulan % 12;
        return sisaBulan > 0 ? `${tahun} tahun ${sisaBulan} bulan` : `${tahun} tahun`;
      })()
    : null;

  const bbTerbaru = riwayat.length > 0 ? riwayat[0].berat_badan : undefined;
  const tbTerbaru = riwayat.length > 0 ? riwayat[0].tinggi_badan : undefined;
  const statusGiziTerbaru = riwayat.length > 0 ? riwayat[0].status_gizi : '-';
  const bbSebelum = riwayat.length > 1 ? riwayat[1].berat_badan : null;
  const tbSebelum = riwayat.length > 1 ? riwayat[1].tinggi_badan : null;
  const selisihBB = bbSebelum !== null && bbTerbaru !== undefined ? (bbTerbaru - bbSebelum).toFixed(1) : null;
  const selisihTB = tbSebelum !== null && tbTerbaru !== undefined ? (tbTerbaru - tbSebelum).toFixed(1) : null;

  if (loading) {
    return (
      <div className="flex items-center justify-center py-16">
        <div className="text-center">
          <Loader2 size={32} className="animate-spin mx-auto mb-3 text-primary" />
          <p className="text-sm text-neutral-500">Memuat data anak...</p>
        </div>
      </div>
    );
  }

  if (!pasien) {
    return (
      <div className="bg-amber-50 border border-amber-200 rounded-2xl p-6 flex items-start gap-4">
        <AlertCircle size={20} className="text-amber-600 flex-shrink-0 mt-0.5" />
        <div>
          <h4 className="text-sm font-bold text-amber-800">Data Pasien Belum Tersedia</h4>
          <p className="text-sm text-amber-700 mt-1">Silakan hubungi bidan atau kader posyandu untuk mendaftarkan data anak Anda.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-5">
      {/* Greeting */}
      <div className="bg-gradient-to-r from-primary-50 to-emerald-50 rounded-2xl p-5 border border-primary-100">
        <div className="flex items-center gap-4">
          <div className="w-12 h-12 bg-primary/10 rounded-2xl flex items-center justify-center">
            <Baby size={24} className="text-primary" />
          </div>
          <div>
            <h2 className="text-lg font-bold text-neutral-800 font-headline">Data Anak: {namaAnak}</h2>
            <p className="text-sm text-neutral-500 mt-0.5">
              {jenisKelamin && `${labelKelamin} • `}Usia: {usia ?? 'Menunggu data lengkap'}
            </p>
          </div>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <div className="bg-white rounded-2xl p-4 border border-neutral-100">
          <Weight size={18} className="text-primary mb-1.5" />
          <p className="text-[10px] font-semibold text-neutral-400 uppercase tracking-wide">Berat Badan</p>
          <p className="text-2xl font-bold text-neutral-800 mt-0.5">{bbTerbaru !== undefined ? bbTerbaru.toFixed(1) : '-'} <span className="text-xs font-normal text-neutral-400">kg</span></p>
          {selisihBB && <p className="text-[10px] text-emerald-500 mt-0.5">{parseFloat(selisihBB) >= 0 ? '+' : ''}{selisihBB} kg</p>}
        </div>
        <div className="bg-white rounded-2xl p-4 border border-neutral-100">
          <Baby size={18} className="text-blue-500 mb-1.5" />
          <p className="text-[10px] font-semibold text-neutral-400 uppercase tracking-wide">Tinggi Badan</p>
          <p className="text-2xl font-bold text-neutral-800 mt-0.5">{tbTerbaru !== undefined ? tbTerbaru.toFixed(1) : '-'} <span className="text-xs font-normal text-neutral-400">cm</span></p>
          {selisihTB && <p className="text-[10px] text-emerald-500 mt-0.5">{parseFloat(selisihTB) >= 0 ? '+' : ''}{selisihTB} cm</p>}
        </div>
        <div className="bg-white rounded-2xl p-4 border border-neutral-100">
          <Heart size={18} className={statusGiziTerbaru === 'Gizi Baik' ? 'text-emerald-500 mb-1.5' : 'text-amber-500 mb-1.5'} />
          <p className="text-[10px] font-semibold text-neutral-400 uppercase tracking-wide">Status Gizi</p>
          <p className={`text-lg font-bold mt-0.5 ${statusGiziTerbaru === 'Gizi Baik' ? 'text-emerald-600' : 'text-amber-600'}`}>{statusGiziTerbaru}</p>
        </div>
        <div className="bg-white rounded-2xl p-4 border border-neutral-100">
          <Calendar size={18} className="text-purple-500 mb-1.5" />
          <p className="text-[10px] font-semibold text-neutral-400 uppercase tracking-wide">Pemeriksaan</p>
          <p className="text-2xl font-bold text-neutral-800 mt-0.5">{riwayat.length}</p>
          <p className="text-[10px] text-neutral-400 mt-0.5">kali</p>
        </div>
      </div>

      {/* Growth Chart */}
      {tumbuhKembang.length >= 2 && (
        <div className="bg-white rounded-2xl border border-neutral-100 p-5">
          <div className="flex items-center gap-2 mb-3">
            <TrendingUp size={16} className="text-emerald-500" />
            <h3 className="text-sm font-bold text-neutral-800 font-headline">Grafik Pertumbuhan</h3>
          </div>
          <div className="h-48">
            {(() => {
              const maxBB = Math.max(...tumbuhKembang.map((d) => d.berat_badan));
              const minBB = Math.min(...tumbuhKembang.map((d) => d.berat_badan));
              const maxTB = Math.max(...tumbuhKembang.map((d) => d.tinggi_badan));
              const minTB = Math.min(...tumbuhKembang.map((d) => d.tinggi_badan));
              const bbRange = maxBB - minBB || 1;
              const tbRange = maxTB - minTB || 1;
              const len = tumbuhKembang.length;
              return (
                <svg viewBox="0 0 600 180" className="w-full h-full" preserveAspectRatio="xMidYMid meet">
                  {/* Grid */}
                  {[0, 0.25, 0.5, 0.75, 1].map((p) => (
                    <line key={`g-${p}`} x1={0} y1={p * 160 + 10} x2={600} y2={p * 160 + 10} stroke="#f3f4f6" strokeWidth={1} />
                  ))}
                  {/* BB line */}
                  <polyline
                    points={tumbuhKembang.map((d, i) => `${(i / (len - 1)) * 560 + 20},${10 + 160 - ((d.berat_badan - minBB) / bbRange) * 150}`).join(' ')}
                    fill="none" stroke="#6366f1" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"
                  />
                  {/* BB dots */}
                  {tumbuhKembang.map((d, i) => (
                    <circle key={`bb-${i}`} cx={(i / (len - 1)) * 560 + 20} cy={10 + 160 - ((d.berat_badan - minBB) / bbRange) * 150} r={3} fill="white" stroke="#6366f1" strokeWidth={2} />
                  ))}
                  {/* TB line */}
                  <polyline
                    points={tumbuhKembang.map((d, i) => `${(i / (len - 1)) * 560 + 20},${10 + 160 - ((d.tinggi_badan - minTB) / tbRange) * 150}`).join(' ')}
                    fill="none" stroke="#10b981" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"
                  />
                  {/* TB dots */}
                  {tumbuhKembang.map((d, i) => (
                    <circle key={`tb-${i}`} cx={(i / (len - 1)) * 560 + 20} cy={10 + 160 - ((d.tinggi_badan - minTB) / tbRange) * 150} r={3} fill="white" stroke="#10b981" strokeWidth={2} />
                  ))}
                  {/* Legend */}
                  <rect x={440} y={8} width={10} height={10} rx={2} fill="#6366f1" />
                  <text x={454} y={17} className="text-[10px] fill-neutral-500" fontFamily="system-ui">BB (kg)</text>
                  <rect x={520} y={8} width={10} height={10} rx={2} fill="#10b981" />
                  <text x={534} y={17} className="text-[10px] fill-neutral-500" fontFamily="system-ui">TB (cm)</text>
                </svg>
              );
            })()}
          </div>
        </div>
      )}

      {/* Imunisasi */}
      {imunisasi.length > 0 && (
        <div className="bg-white rounded-2xl border border-neutral-100 p-5">
          <div className="flex items-center gap-2 mb-3">
            <Syringe size={16} className="text-blue-600" />
            <h3 className="text-sm font-bold text-neutral-800 font-headline">Jadwal Imunisasi</h3>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
            {imunisasi.slice(0, 6).map((item, i) => (
              <div key={i} className="flex items-center gap-2 p-2.5 rounded-lg border border-neutral-100 bg-neutral-50/50">
                <div className={`w-7 h-7 rounded-full flex items-center justify-center shrink-0 ${item.status_imunisasi === 'Sudah' ? 'bg-emerald-100 text-emerald-600' : 'bg-blue-100 text-blue-600'}`}>
                  <Syringe size={12} />
                </div>
                <div className="min-w-0">
                  <p className="text-xs font-semibold text-neutral-700 truncate">{item.nama_vaksin}</p>
                  <p className="text-[10px] text-neutral-500">{item.tanggal_jadwal?.slice(0,10)} · {item.status_imunisasi}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Riwayat Table */}
      <div>
        <div className="flex items-center gap-2 mb-3">
          <ClipboardList size={16} className="text-primary" />
          <h3 className="text-sm font-bold text-neutral-800 font-headline">Riwayat Pengukuran</h3>
          <span className="text-xs text-neutral-400 ml-auto">{riwayat.length} data</span>
        </div>

        <div className="bg-white rounded-2xl border border-neutral-100 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-xs">
              <thead>
                <tr className="bg-neutral-50 border-b border-neutral-100 text-[10px] font-bold text-neutral-400 uppercase tracking-wide">
                  <th className="text-left py-2.5 px-4">Tanggal</th>
                  <th className="text-center py-2.5 px-3">BB (kg)</th>
                  <th className="text-center py-2.5 px-3">TB (cm)</th>
                  <th className="text-center py-2.5 px-3">Status Gizi</th>
                  <th className="text-left py-2.5 px-4">Catatan</th>
                  <th className="text-left py-2.5 px-4">Petugas</th>
                </tr>
              </thead>
              <tbody>
                {riwayat.length === 0 ? (
                  <tr><td colSpan={6} className="py-12 text-center text-neutral-400">Belum ada riwayat pengukuran</td></tr>
                ) : (
                  riwayat.map((row, i) => (
                    <tr key={i} className="border-b border-neutral-50 hover:bg-neutral-50/60">
                      <td className="py-2.5 px-4 text-neutral-500 whitespace-nowrap">{row.tanggal?.slice(0, 10) || '-'}</td>
                      <td className="py-2.5 px-3 text-center text-neutral-700 tabular-nums font-medium">{row.berat_badan?.toFixed(1)}</td>
                      <td className="py-2.5 px-3 text-center text-neutral-700 tabular-nums font-medium">{row.tinggi_badan?.toFixed(1)}</td>
                      <td className="py-2.5 px-3 text-center">
                        <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${row.status_gizi === 'Gizi Baik' ? 'bg-emerald-100 text-emerald-600' : 'bg-amber-100 text-amber-600'}`}>{row.status_gizi}</span>
                      </td>
                      <td className="py-2.5 px-4 text-neutral-500 max-w-40 truncate">{row.catatan || '-'}</td>
                      <td className="py-2.5 px-4 text-neutral-500 whitespace-nowrap">{row.petugas || '-'}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
