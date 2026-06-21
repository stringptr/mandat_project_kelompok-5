import { useState, useEffect } from 'react';
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

export default function TindakLanjut({ currentRole }: TindakLanjutProps): JSX.Element {
  if (currentRole !== 'Bidan') {
    return <></>;
  }

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
          patientAge: String(r.umur_pasien ?? ''),
          urgency: 'Berjalan' as Rujukan['urgency'],
          faskes: String(r.faskes ?? ''),
          jenisTindakan: 'Rujukan' as Rujukan['jenisTindakan'],
          status: (r.status_rujukan as Rujukan['status']) || 'Diajukan',
          lastWeight: 0,
          lastHeight: 0,
          nutritionStatus: String(r.status_gizi ?? 'Normal'),
        }));
        setRujukanList(mapped);
      })
      .catch(() => notify.error('Gagal memuat data tindak lanjut.'))
      .finally(() => setRujukanLoading(false));
  };

  useEffect(() => {
    fetchRujukan();
  }, [paginator.page]); // eslint-disable-line react-hooks/exhaustive-deps

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
        {/* Form occupies two columns */}
        <div className="lg:col-span-2">
          <FormTindakLanjut
            onSubmit={handleFormSubmit}
            onCancel={() => {}}
          />
        </div>

        <div className="lg:col-span-1 flex flex-col space-y-4">
          {/* Filter dropdown */}
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
