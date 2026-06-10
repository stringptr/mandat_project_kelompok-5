/**
 * Edukasi — orchestrator
 *
 * Holds shared state (artikel list, filters, modal targets) and
 * delegates all rendering to the appropriate role section.
 *
 * Role → Section map
 *   Ibu/Wali       → IbuWaliSection  (read-only)
 *   Kader Posyandu → KaderSection    (read-only)
 *   Bidan          → BidanSection    (write own, pending verification)
 *   Dinas Kesehatan→ DinkesSection   (full CRUD + verifikasi)
 */
import { useState } from 'react';
import type { Role } from '../../App';
import type { Artikel, KategoriArtikel } from './data/artikel.data';
import { DUMMY_ARTIKEL } from './data/artikel.data';

// Sections
import { IbuWaliSection } from './sections/IbuWaliSection';
import { KaderSection } from './sections/KaderSection';
import { BidanSection } from './sections/BidanSection';
import { DinkesSection } from './sections/DinkesSection';

// Shared modals
import { ModalDetailArtikel } from './components/ModalDetailArtikel';
import { ModalTambahArtikel } from './components/ModalTambahArtikel';
import { ModalEditArtikel } from './components/ModalEditArtikel';
import { ModalHapusArtikel } from './components/ModalHapusArtikel';
import { ModalVerifikasiArtikel } from './components/ModalVerifikasiArtikel';

interface EdukasiProps {
  currentRole: Role;
}

export default function Edukasi({ currentRole }: EdukasiProps): JSX.Element {
  // ── State ──────────────────────────────────────────────────────────────
  const [artikelList, setArtikelList] = useState<Artikel[]>(DUMMY_ARTIKEL);
  const [kategoriAktif, setKategoriAktif] = useState<KategoriArtikel>('Semua');
  const [searchQuery, setSearchQuery] = useState('');
  const [showAll, setShowAll] = useState(false);

  // Modal targets
  const [modalTambah, setModalTambah] = useState(false);
  const [modalEdit, setModalEdit] = useState<Artikel | null>(null);
  const [modalHapus, setModalHapus] = useState<Artikel | null>(null);
  const [modalDetail, setModalDetail] = useState<Artikel | null>(null);
  const [modalVerifikasi, setModalVerifikasi] = useState<Artikel | null>(null);

  // ── Derived role flags ─────────────────────────────────────────────────
  const isDinkes = currentRole === 'Dinas Kesehatan';
  const isBidan = currentRole === 'Bidan';

  // ── Permission helpers ─────────────────────────────────────────────────
  const canEditArtikel = (a: Artikel): boolean => {
    if (isDinkes) return true;
    if (isBidan && a.rolePenulis === 'Bidan') return true;
    return false;
  };

  // ── Derived data ───────────────────────────────────────────────────────
  // Dinkes sees all statuses; everyone else only sees published
  const visibleArtikel = artikelList.filter((a) =>
    isDinkes ? true : a.status === 'published'
  );

  const filteredArtikel = visibleArtikel.filter((a) => {
    const matchKat = kategoriAktif === 'Semua' || a.kategori === kategoriAktif;
    const matchQ =
      searchQuery === '' ||
      a.judul.toLowerCase().includes(searchQuery.toLowerCase()) ||
      a.ringkasan.toLowerCase().includes(searchQuery.toLowerCase());
    return matchKat && matchQ;
  });

  const featuredArtikel = filteredArtikel.filter(
    (a) => a.featured && a.status === 'published'
  );
  const allGridArtikel = filteredArtikel.filter(
    (a) => !a.featured || a.status !== 'published'
  );
  const gridArtikel = showAll ? allGridArtikel : allGridArtikel.slice(0, 6);

  const pendingCount = artikelList.filter((a) => a.status === 'pending').length;

  // ── Handlers ───────────────────────────────────────────────────────────
  const handleTambah = (data: Omit<Artikel, 'id' | 'tanggal' | 'status' | 'rolePenulis'>) => {
    const newArtikel: Artikel = {
      ...data,
      id: `a${Date.now()}`,
      tanggal: new Date().toLocaleDateString('id-ID', {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
      }),
      status: isDinkes ? 'published' : 'pending',
      rolePenulis: currentRole,
    };
    setArtikelList((prev) => [newArtikel, ...prev]);
    setModalTambah(false);
  };

  const handleEdit = (data: Artikel) => {
    setArtikelList((prev) => prev.map((a) => (a.id === data.id ? data : a)));
    setModalEdit(null);
  };

  const handleHapus = (id: string) => {
    setArtikelList((prev) => prev.filter((a) => a.id !== id));
    setModalHapus(null);
  };

  const handleVerifikasi = (id: string, action: 'approve' | 'reject') => {
    setArtikelList((prev) =>
      prev.map((a) =>
        a.id === id
          ? { ...a, status: action === 'approve' ? 'published' : 'rejected' }
          : a
      )
    );
    setModalVerifikasi(null);
  };

  // Open first pending article in verifikasi modal
  const bukaAntrianVerifikasi = () => {
    const pending = artikelList.find((a) => a.status === 'pending');
    if (pending) setModalVerifikasi(pending);
  };

  // ── Shared section props ───────────────────────────────────────────────
  const commonProps = {
    featuredArtikel,
    gridArtikel,
    showAll,
    totalGrid: allGridArtikel.length,
    kategoriAktif,
    searchQuery,
    onChangeKategori: setKategoriAktif,
    onChangeSearch: setSearchQuery,
    onRead: setModalDetail,
    onToggleShowAll: () => setShowAll((v) => !v),
  };

  // ── Page header (shared across all roles) ─────────────────────────────
  return (
    <div className="space-y-6 font-body text-neutral-800">
      {/* Page heading */}
      <div>
        <p className="text-xs font-semibold text-primary uppercase tracking-widest mb-1 font-body">
          Literasi Kesehatan Komunitas
        </p>
        <h2 className="text-2xl font-bold text-neutral-900 font-headline leading-tight">
          Wujudkan Generasi{' '}
          <span className="text-primary">Bebas Stunting</span>{' '}
          Melalui Edukasi.
        </h2>
      </div>

      {/* Role-specific section */}
      {currentRole === 'Ibu/Wali' && (
        <IbuWaliSection {...commonProps} />
      )}

      {currentRole === 'Kader Posyandu' && (
        <KaderSection {...commonProps} />
      )}

      {currentRole === 'Bidan' && (
        <BidanSection
          {...commonProps}
          onEdit={(a) => canEditArtikel(a) && setModalEdit(a)}
          onHapus={(a) => canEditArtikel(a) && setModalHapus(a)}
          onTambah={() => setModalTambah(true)}
          canEdit={canEditArtikel}
        />
      )}

      {currentRole === 'Dinas Kesehatan' && (
        <DinkesSection
          {...commonProps}
          pendingCount={pendingCount}
          onEdit={(a) => canEditArtikel(a) && setModalEdit(a)}
          onHapus={(a) => canEditArtikel(a) && setModalHapus(a)}
          onTambah={() => setModalTambah(true)}
          onVerifikasi={setModalVerifikasi}
          onBukaAntrianVerifikasi={bukaAntrianVerifikasi}
          canEdit={canEditArtikel}
        />
      )}

      {/* Empty state */}
      {filteredArtikel.length === 0 && (
        <div className="text-center py-16 text-neutral-400">
          <div className="text-5xl mb-3">📚</div>
          <p className="font-semibold text-neutral-500">Tidak ada artikel ditemukan</p>
          <p className="text-sm mt-1">Coba ubah filter atau kata kunci pencarian</p>
        </div>
      )}

      {/* ── Shared modals ─────────────────────────────────────────────── */}
      {modalDetail && (
        <ModalDetailArtikel artikel={modalDetail} onClose={() => setModalDetail(null)} />
      )}

      {modalTambah && (
        <ModalTambahArtikel
          currentRole={currentRole}
          onClose={() => setModalTambah(false)}
          onSubmit={handleTambah}
        />
      )}

      {modalEdit && (
        <ModalEditArtikel
          artikel={modalEdit}
          onClose={() => setModalEdit(null)}
          onSubmit={handleEdit}
        />
      )}

      {modalHapus && (
        <ModalHapusArtikel
          artikel={modalHapus}
          onClose={() => setModalHapus(null)}
          onHapus={() => handleHapus(modalHapus.id)}
        />
      )}

      {modalVerifikasi && isDinkes && (
        <ModalVerifikasiArtikel
          artikel={modalVerifikasi}
          onClose={() => setModalVerifikasi(null)}
          onVerifikasi={handleVerifikasi}
        />
      )}
    </div>
  );
}
