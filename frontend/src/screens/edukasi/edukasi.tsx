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
import { useState, useEffect, useCallback, useRef } from 'react';
import { useNotification } from '../../context/NotificationContext';
import { apiGet, apiDelete, apiPatch } from '../../lib/api';
import type { ArtikelListData } from '../../types/entities';
import type { Role } from '../../App';
import type { Artikel, KategoriArtikel } from './data/artikel.data';
import { useAppStore } from '../../store/useAppStore';

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
  const notify = useNotification();
  const notifyRef = useRef(notify);
  useEffect(() => { notifyRef.current = notify; });

  const { artikelList, setArtikelList, setArtikelLoading } = useAppStore();

  const fetchArtikel = useCallback((force = false) => {
    // Skip if already loaded and not forced
    if (artikelList.length > 0 && !force) return;
    setArtikelLoading(true);
    apiGet<ArtikelListData>('/artikel')
      .then((res) => {
        const list = res.artikel ?? [];
        const mapped: Artikel[] = list.map((a) => ({
          id: String(a.id_artikel ?? ''),
          judul: String(a.judul ?? ''),
          ringkasan: String(a.ringkasan ?? ''),
          konten: '',
          kategori: (a.kategori as KategoriArtikel) || 'Gizi Ibu',
          penulis: String(a.nama_penulis ?? ''),
          gambar: '',
          waktuBaca: '5 Menit Baca',
          featured: false,
          tanggal: String(a.tanggal_publish ?? ''),
          status: 'published',
          rolePenulis: 'Dinas Kesehatan',
        }));
        setArtikelList(mapped);
      })
      .catch(() => notifyRef.current.error('Gagal memuat artikel. Pastikan backend terhubung.'))
      .finally(() => setArtikelLoading(false));
  }, [artikelList.length, setArtikelList, setArtikelLoading]);

  useEffect(() => {
    fetchArtikel();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps
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
  const handleTambah = (_data: Omit<Artikel, 'id' | 'tanggal' | 'status' | 'rolePenulis'>) => {
    setModalTambah(false);
    notify.success('Artikel berhasil ditambahkan');
    fetchArtikel(true);
  };

  const handleEdit = (_data: Artikel) => {
    setModalEdit(null);
    notify.success('Artikel berhasil diperbarui');
    fetchArtikel(true);
  };

  const handleHapus = async (id: string) => {
    const idNum = parseInt(String(id).replace(/[^0-9]/g, ''), 10) || 0;
    if (!idNum) { notify.error('ID artikel tidak valid.'); return; }
    try {
      await apiDelete('/artikel/' + idNum);
      setModalHapus(null);
      notify.success('Artikel berhasil dihapus');
      fetchArtikel(true);
    } catch {
      notify.error('Gagal menghapus artikel. Silakan coba lagi.');
    }
  };

  const handleVerifikasi = async (id: string, action: 'approve' | 'reject') => {
    const idNum = parseInt(String(id).replace(/[^0-9]/g, ''), 10) || 0;
    if (!idNum) { notify.error('ID artikel tidak valid.'); return; }
    try {
      await apiPatch('/artikel/' + idNum + '/review', {
        aksi: action === 'approve' ? 'setujui' : 'tolak',
      });
      setModalVerifikasi(null);
      notify.success(
        action === 'approve'
          ? 'Artikel berhasil disetujui dan dipublikasikan'
          : 'Artikel berhasil ditolak'
      );
      fetchArtikel(true);
    } catch {
      notify.error('Gagal memverifikasi artikel. Silakan coba lagi.');
    }
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
