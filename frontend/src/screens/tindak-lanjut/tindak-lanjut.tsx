import { useState, useEffect } from 'react';
import { useAuth } from '../../context/AuthContext';
import { useNotification } from '../../context/NotificationContext';
import { apiGet } from '../../lib/api';
import { useAppStore } from '../../store/useAppStore';
import { usePaginator } from '../../hooks/usePaginator';
import { Paginator } from '../../components/Paginator';
import FormTindakLanjut from './components/FormTindakLanjut';
import RujukanAktif from './components/RujukanAktif';
import DetailRujukanModal from './components/DetailRujukanModal';
import UpdateRujukanModal from './components/UpdateRujukanModal';
import type { Rujukan } from './components/RujukanAktif';
import type { Role } from '../../App';

const PAGE_SIZE = 10;

interface TindakLanjutProps {
  currentRole: Role;
}

interface IbuTindakLanjutItem {
  id_pasien: number;
  nama_pasien: string;
  status_pasien: string;
  status_rujukan: string;
  tanggal_rujukan: string;
  tanggal_deadline: string;
  status_deadline: string;
}

function IbuWaliTindakLanjut(): JSX.Element {
  const { user } = useAuth();
  const notifyIbu = useNotification();
  const [items, setItems] = useState<IbuTindakLanjutItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!user?.idUser) return;
    setLoading(true);
    apiGet<{ pasien: IbuTindakLanjutItem[] }>(`/tindak-lanjut/status?page=1&per_page=100`)
      .then((res) => {
        const filtered = (res.pasien ?? []).filter(
          (item) => Number(item.id_pasien) === Number(user.idUser)
        );
        setItems(filtered);
      })
      .catch(() => notifyIbu.error('Gagal memuat data tindak lanjut.'))
      .finally(() => setLoading(false));
  }, [user?.idUser]);

  if (loading) {
    return (
      <section className="p-8 flex items-center justify-center py-16">
        <div className="w-6 h-6 border-2 border-primary/30 border-t-primary rounded-full animate-spin" />
      </section>
    );
  }

  return (
    <section className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-slate-800">Tindak Lanjut Saya</h1>
        <p className="text-sm text-slate-500 max-w-2xl">
          Riwayat tindak lanjut dan rujukan untuk anak Anda.
        </p>
      </div>

      {items.length === 0 ? (
        <div className="bg-white rounded-2xl border border-neutral-100 p-12 text-center">
          <div className="w-16 h-16 bg-emerald-50 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <p className="text-sm text-neutral-500">Tidak ada tindak lanjut yang perlu dikhawatirkan</p>
        </div>
      ) : (
        <div className="bg-white rounded-2xl border border-neutral-100 overflow-hidden">
          <table className="w-full text-xs">
            <thead>
              <tr className="bg-neutral-50 border-b border-neutral-100 text-[10px] font-bold text-neutral-400 uppercase tracking-wide">
                <th className="text-left py-2.5 px-4">Nama</th>
                <th className="text-left py-2.5 px-4">Status</th>
                <th className="text-left py-2.5 px-4">Rujukan</th>
                <th className="text-left py-2.5 px-4">Tanggal</th>
                <th className="text-left py-2.5 px-4">Deadline</th>
              </tr>
            </thead>
            <tbody>
              {items.map((item) => (
                <tr key={item.id_pasien} className="border-b border-neutral-50 hover:bg-neutral-50/60">
                  <td className="py-2.5 px-4 font-medium text-neutral-800">{item.nama_pasien}</td>
                  <td className="py-2.5 px-4">
                    <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${
                      item.status_pasien === 'Perlu Rujukan' ? 'bg-amber-100 text-amber-700' :
                      'bg-blue-100 text-blue-700'
                    }`}>{item.status_pasien}</span>
                  </td>
                  <td className="py-2.5 px-4">
                    <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${
                      item.status_rujukan === 'Diterima' ? 'bg-emerald-100 text-emerald-700' :
                      item.status_rujukan === 'Ditolak' ? 'bg-red-100 text-red-700' :
                      item.status_rujukan === 'Selesai' ? 'bg-blue-100 text-blue-700' :
                      item.status_rujukan === 'Diproses' ? 'bg-purple-100 text-purple-700' :
                      item.status_rujukan ? 'bg-amber-100 text-amber-700' :
                      ''
                    }`}>{item.status_rujukan || '-'}</span>
                  </td>
                  <td className="py-2.5 px-4 text-neutral-500">{item.tanggal_rujukan?.slice(0, 10) || '-'}</td>
                  <td className="py-2.5 px-4">
                    <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${
                      item.status_deadline === 'terlambat' ? 'bg-red-100 text-red-600' :
                      item.status_deadline === 'mendekati' ? 'bg-amber-100 text-amber-600' :
                      item.status_deadline ? 'bg-emerald-100 text-emerald-600' : ''
                    }`}>{
                      item.status_deadline === 'terlambat' ? 'Terlambat' :
                      item.status_deadline === 'mendekati' ? 'Mendekati' :
                      item.status_deadline || '-'
                    }</span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

function BidanTindakLanjut(): JSX.Element {
  const notify = useNotification();
  const { rujukanList, rujukanLoading, setRujukanList, setRujukanLoading } = useAppStore();
  const [filter, setFilter] = useState<'Semua' | 'Rujukan' | 'Tindak Lanjut'>('Semua');
  const [selectedPatient, setSelectedPatient] = useState<Rujukan | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [updatePatient, setUpdatePatient] = useState<Rujukan | null>(null);
  const [updateModalOpen, setUpdateModalOpen] = useState(false);
  const [totalData, setTotalData] = useState(0);
  const paginator = usePaginator({ totalItems: totalData, pageSize: PAGE_SIZE });

  const filteredData = filter === 'Semua'
    ? rujukanList
    : rujukanList.filter((item) =>
        filter === 'Rujukan' ? item.jenisTindakan === 'Rujukan' : item.jenisTindakan === 'Tindak Lanjut',
      );

  const fetchRujukan = (targetPage?: number) => {
    const p = targetPage ?? paginator.page;
    setRujukanLoading(true);
    apiGet<{ pasien?: Record<string, unknown>[]; meta?: { total: number } }>(
      `/tindak-lanjut/status?page=${p}&per_page=${PAGE_SIZE}`,
    )
      .then((res) => {
        const list = (res.pasien ?? []) as Record<string, unknown>[];
        setTotalData(res.meta?.total ?? 0);
        const mapped: Rujukan[] = list.map((r) => ({
          id: String(r.id_pasien ?? ''),
          patientName: String(r.nama_pasien ?? ''),
          patientAge: '',
          urgency: 'Berjalan' as Rujukan['urgency'],
          faskes: String(r.faskes ?? ''),
          jenisTindakan: (r.jenis_tindakan as Rujukan['jenisTindakan']) || 'Rujukan',
          status: (r.status_rujukan as Rujukan['status']) || 'Diajukan',
          lastWeight: 0,
          lastHeight: 0,
          nutritionStatus: String(r.status_pasien ?? 'Normal'),
          tanggalDeadline: String(r.tanggal_deadline ?? ''),
          statusDeadline: String(r.status_deadline ?? ''),
          tanggalRujukan: String(r.tanggal_rujukan ?? ''),
          alasanRujukan: String(r.alasan_rujukan ?? ''),
        }));
        setRujukanList(mapped);
      })
      .catch(() => notify.error('Gagal memuat data tindak lanjut.'))
      .finally(() => setRujukanLoading(false));
  };

  useEffect(() => {
    fetchRujukan();
  }, [paginator.page]);

  const handleUpdate = (updated: Rujukan) => {
    setRujukanList(rujukanList.map((item) => (item.id === updated.id ? updated : item)));
    notify.success('Data rujukan berhasil diperbarui');
  };

  const handleFormSubmit = () => {
    if (paginator.page !== 1) {
      paginator.setPage(1);
    } else {
      fetchRujukan(1);
    }
  };

  const closeModal = () => { setModalOpen(false); setSelectedPatient(null); };
  const closeUpdateModal = () => { setUpdateModalOpen(false); setUpdatePatient(null); };
  const openUpdate = (item: Rujukan) => { setUpdatePatient(item); setUpdateModalOpen(true); };

  return (
    <section className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-slate-800">Manajemen Tindak Lanjut</h1>
        <p className="text-sm text-slate-500 max-w-2xl">
          Kelola rujukan medis, jadwalkan ulang pemeriksaan, dan catat instruksi untuk memastikan pemantauan gizi yang tepat.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <FormTindakLanjut
            onSubmit={handleFormSubmit}
            onCancel={() => {}}
          />
        </div>

        <div className="lg:col-span-1 flex flex-col space-y-4">
          <div className="mb-4">
            <label className="block text-sm font-medium text-slate-700 mb-1">Filter Rujukan / Tindak Lanjut</label>
            <select
              value={filter}
              onChange={(e) => setFilter(e.target.value as typeof filter)}
              className="w-full bg-gray-100 border border-gray-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
            >
              <option value="Semua">Semua</option>
              <option value="Rujukan">Rujukan</option>
              <option value="Tindak Lanjut">Tindak Lanjut</option>
            </select>
          </div>

          {rujukanList.length === 0 && !rujukanLoading ? (
            <div className="flex-1 flex items-center justify-center py-8">
              <p className="text-sm text-neutral-400">Belum ada data rujukan</p>
            </div>
          ) : (
            <div className="flex-1 overflow-y-auto">
              <RujukanAktif
                data={filteredData}
                onDetailClick={(p) => { setSelectedPatient(p); setModalOpen(true); }}
                onUpdateClick={openUpdate}
              />
            </div>
          )}

          <Paginator
            page={paginator.page}
            totalPages={paginator.totalPages}
            totalItems={totalData}
            pageSize={PAGE_SIZE}
            onPageChange={paginator.setPage}
          />
        </div>
      </div>

      <DetailRujukanModal isOpen={modalOpen} onClose={closeModal} rujukan={selectedPatient} />
      {updatePatient && (
        <UpdateRujukanModal
          isOpen={updateModalOpen}
          onClose={closeUpdateModal}
          rujukan={updatePatient}
          onSave={handleUpdate}
        />
      )}
    </section>
  );
}

interface LaporanItem {
  wilayah: string;
  jumlah_pasien_dirujuk: number;
  jumlah_pasien_diterima: number;
  jumlah_pasien_diproses: number;
}

function DinkesLaporanTindakLanjut(): JSX.Element {
  const notify = useNotification();
  const [laporan, setLaporan] = useState<LaporanItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    apiGet<{ laporan: LaporanItem[]; total_data: number }>('/laporan/tindak-lanjut')
      .then((res) => setLaporan(res.laporan ?? []))
      .catch(() => notify.error('Gagal memuat laporan tindak lanjut.'))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <section className="p-8 flex items-center justify-center py-16">
        <div className="w-6 h-6 border-2 border-primary/30 border-t-primary rounded-full animate-spin" />
      </section>
    );
  }

  const totalDirujuk = laporan.reduce((s, i) => s + i.jumlah_pasien_dirujuk, 0);
  const totalDiterima = laporan.reduce((s, i) => s + i.jumlah_pasien_diterima, 0);
  const totalDiproses = laporan.reduce((s, i) => s + i.jumlah_pasien_diproses, 0);
  const maxVal = Math.max(totalDirujuk, totalDiterima, totalDiproses, 1);

  return (
    <section className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-slate-800">Laporan Tindak Lanjut</h1>
        <p className="text-sm text-slate-500 max-w-2xl">
          Ringkasan rujukan dan tindak lanjut per wilayah.
        </p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
        <div className="bg-white rounded-2xl border border-neutral-100 p-5">
          <p className="text-xs text-neutral-400 uppercase tracking-wide font-semibold">Total Dirujuk</p>
          <p className="text-3xl font-bold text-amber-600 mt-1">{totalDirujuk.toLocaleString('id-ID')}</p>
          <div className="mt-2 w-full h-1.5 bg-neutral-100 rounded-full overflow-hidden">
            <div className="h-full bg-amber-400 rounded-full" style={{ width: `${(totalDirujuk / maxVal) * 100}%` }} />
          </div>
        </div>
        <div className="bg-white rounded-2xl border border-neutral-100 p-5">
          <p className="text-xs text-neutral-400 uppercase tracking-wide font-semibold">Total Diterima</p>
          <p className="text-3xl font-bold text-emerald-600 mt-1">{totalDiterima.toLocaleString('id-ID')}</p>
          <div className="mt-2 w-full h-1.5 bg-neutral-100 rounded-full overflow-hidden">
            <div className="h-full bg-emerald-400 rounded-full" style={{ width: `${(totalDiterima / maxVal) * 100}%` }} />
          </div>
        </div>
        <div className="bg-white rounded-2xl border border-neutral-100 p-5">
          <p className="text-xs text-neutral-400 uppercase tracking-wide font-semibold">Total Diproses</p>
          <p className="text-3xl font-bold text-blue-600 mt-1">{totalDiproses.toLocaleString('id-ID')}</p>
          <div className="mt-2 w-full h-1.5 bg-neutral-100 rounded-full overflow-hidden">
            <div className="h-full bg-blue-400 rounded-full" style={{ width: `${(totalDiproses / maxVal) * 100}%` }} />
          </div>
        </div>
      </div>

      {laporan.length === 0 ? (
        <div className="bg-white rounded-2xl border border-neutral-100 p-12 text-center">
          <p className="text-sm text-neutral-400">Belum ada data laporan tindak lanjut</p>
        </div>
      ) : (
        <div className="bg-white rounded-2xl border border-neutral-100 overflow-hidden">
          <table className="w-full text-xs">
            <thead>
              <tr className="bg-neutral-50 border-b border-neutral-100 text-[10px] font-bold text-neutral-400 uppercase tracking-wide">
                <th className="text-left py-3 px-5">Wilayah</th>
                <th className="text-center py-3 px-4">Dirujuk</th>
                <th className="text-center py-3 px-4">Diterima</th>
                <th className="text-center py-3 px-4">Diproses</th>
                <th className="text-left py-3 px-5">Progress</th>
              </tr>
            </thead>
            <tbody>
              {laporan.map((item) => {
                const rowMax = Math.max(item.jumlah_pasien_dirujuk, item.jumlah_pasien_diterima, item.jumlah_pasien_diproses, 1);
                return (
                  <tr key={item.wilayah} className="border-b border-neutral-50 hover:bg-neutral-50/60">
                    <td className="py-3 px-5 font-semibold text-neutral-800">{item.wilayah}</td>
                    <td className="py-3 px-4 text-center font-medium text-amber-600">{item.jumlah_pasien_dirujuk}</td>
                    <td className="py-3 px-4 text-center font-medium text-emerald-600">{item.jumlah_pasien_diterima}</td>
                    <td className="py-3 px-4 text-center font-medium text-blue-600">{item.jumlah_pasien_diproses}</td>
                    <td className="py-3 px-5">
                      <div className="flex items-center gap-0.5 h-2 rounded-full overflow-hidden bg-neutral-100">
                        <div className="h-full bg-amber-400 rounded-l-full" style={{ width: `${(item.jumlah_pasien_dirujuk / rowMax) * 33}%`, minWidth: item.jumlah_pasien_dirujuk > 0 ? '6px' : '0' }} />
                        <div className="h-full bg-emerald-400" style={{ width: `${(item.jumlah_pasien_diterima / rowMax) * 33}%`, minWidth: item.jumlah_pasien_diterima > 0 ? '6px' : '0' }} />
                        <div className="h-full bg-blue-400 rounded-r-full" style={{ width: `${(item.jumlah_pasien_diproses / rowMax) * 34}%`, minWidth: item.jumlah_pasien_diproses > 0 ? '6px' : '0' }} />
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

export default function TindakLanjut({ currentRole }: TindakLanjutProps): JSX.Element {
  if (currentRole === 'Bidan') {
    return <BidanTindakLanjut />;
  }
  if (currentRole === 'Ibu/Wali') {
    return <IbuWaliTindakLanjut />;
  }
  if (currentRole === 'Dinas Kesehatan') {
    return <DinkesLaporanTindakLanjut />;
  }
  return <></>;
}
