// Updated TindakLanjut page with Form and Active Referrals layout
import { useState } from 'react';
import FormTindakLanjut from './components/FormTindakLanjut';
import RujukanAktif from './components/RujukanAktif';
import DetailRujukanModal from './components/DetailRujukanModal';
import UpdateRujukanModal from './components/UpdateRujukanModal';
import type { Rujukan } from './components/RujukanAktif';
import type { Role } from '../../App';

const INITIAL_RUJUKAN: Rujukan[] = [
  { id: '1', patientName: 'Ananda Revan', patientAge: '24 Bulan', urgency: 'Mendesak', faskes: 'RSUD Kota - Spesialis Anak', jenisTindakan: 'Rujukan', status: 'Diajukan', lastWeight: 12, lastHeight: 85, nutritionStatus: 'Normal' },
  { id: '2', patientName: 'Siti Aminah', patientAge: 'Hamil: 28 Minggu', urgency: 'Berjalan', faskes: 'Puskesmas Melati', jenisTindakan: 'Rujukan', status: 'Diproses', lastWeight: 9, lastHeight: 70, nutritionStatus: 'Underweight' },
  { id: '3', patientName: 'Laila Kirana', patientAge: '6 Bulan', urgency: 'Review', faskes: 'Posyandu Melati', jenisTindakan: 'Tindak Lanjut', status: 'Diterima', lastWeight: 6.5, lastHeight: 65, nutritionStatus: 'Normal' },
];

interface TindakLanjutProps {
  currentRole: Role;
}

export default function TindakLanjut({ currentRole }: TindakLanjutProps): JSX.Element {
  // Show only for Bidan role
  if (currentRole !== 'Bidan') {
    return <></>;
  }

  const [data, setData] = useState<Rujukan[]>(INITIAL_RUJUKAN);
  const [filter, setFilter] = useState<'Semua' | 'Rujukan' | 'Tindak Lanjut'>('Semua');
  const [searchQuery] = useState('');
  const filteredData = filter === 'Semua' ? data : data.filter(item => (filter === 'Rujukan' ? item.jenisTindakan === 'Rujukan' : item.jenisTindakan === 'Tindak Lanjut'));


  // Update handler
  const handleUpdate = (updated: Rujukan) => {
    setData(prev => prev.map(item => (item.id === updated.id ? updated : item)));
  };
  const [selectedPatient, setSelectedPatient] = useState<Rujukan | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [updatePatient, setUpdatePatient] = useState<Rujukan | null>(null);
  const [updateModalOpen, setUpdateModalOpen] = useState(false);


  const closeModal = () => {
    setModalOpen(false);
    setSelectedPatient(null);
  };
  const closeUpdateModal = () => {
    setUpdateModalOpen(false);
    setUpdatePatient(null);
  };
  const openUpdate = (item: Rujukan) => {
    setUpdatePatient(item);
    setUpdateModalOpen(true);
  };
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
            onSubmit={data => console.log('Submitted', data)}
            onCancel={() => console.log('Form cancelled')}
          />
          {/* Live search results */}
          {searchQuery && (
            <div className="mt-4 space-y-2">
              {data.filter(p => p.patientName.toLowerCase().includes(searchQuery.toLowerCase()) || p.id.includes(searchQuery)).map(p => (
                <div key={p.id} className="cursor-pointer p-2 bg-[#f4f7fc] rounded-lg hover:bg-[#e0e7ff]" onClick={() => { setSelectedPatient(p); setModalOpen(true); }}>
                  <span className="font-medium">{p.patientName}</span> - <span className="text-sm text-neutral-600">{p.patientAge}</span>
                </div>
              ))}
            </div>
          )}
        </div>
        <div className="lg:col-span-1 flex flex-col space-y-4">
          {/* Filter dropdown */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-slate-700 mb-1">Filter Rujukan / Tindak Lanjut</label>
            <select
              value={filter}
              onChange={e => setFilter(e.target.value as typeof filter)}
              className="w-full bg-gray-100 border border-gray-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
            >
              <option value="Semua">Semua</option>
              <option value="Rujukan">Rujukan</option>
              <option value="Tindak Lanjut">Tindak Lanjut</option>
            </select>
          </div>
          {/* Scrollable active referrals card */}
          <div className="flex-1 overflow-y-auto">
            <RujukanAktif
              data={filteredData}
              onDetailClick={p => { setSelectedPatient(p); setModalOpen(true); }}
              onUpdateClick={openUpdate}
            />
          </div>
        </div>
      </div>
      {/* Detail modal */}
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