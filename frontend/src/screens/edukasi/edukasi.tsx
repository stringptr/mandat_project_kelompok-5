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
import { useState, useEffect, useCallback, useRef, useMemo } from 'react';
import { useNotification } from '../../context/NotificationContext';
import { apiGet, apiDelete, apiPatch } from '../../lib/api';
import type { ArtikelListData } from '../../types/entities';
import type { Role } from '../../App';
import type { Artikel, KategoriArtikel } from './data/artikel.data';
import { KATEGORI_LIST } from './data/artikel.data';
import { useAppStore } from '../../store/useAppStore';
import { Paginator } from '../../components/Paginator';

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
  currentRole?: Role;
}

const PAGE_SIZE = 12;

export default function Edukasi({ currentRole }: EdukasiProps): JSX.Element {
  const notify = useNotification();
  const notifyRef = useRef(notify);
  useEffect(() => { notifyRef.current = notify; });

  const { artikelList, setArtikelList, setArtikelLoading } = useAppStore();

  const [page, setPage] = useState(1);
  const isDinkes = currentRole === 'Dinas Kesehatan';
  const isBidan = currentRole === 'Bidan';

  const fetchArtikel = useCallback(() => {
    setArtikelLoading(true);
    const endpoint = isDinkes ? '/artikel/semua' : '/artikel';

    apiGet<ArtikelListData>(endpoint, { page: '1', per_page: '100' })
      .then((res) => {
        const list = res.artikel ?? [];
        const mapped: Artikel[] = list.map((a) => ({
          id: String(a.id_artikel ?? ''),
          judul: String(a.judul ?? ''),
          ringkasan: String(a.ringkasan ?? ''),
          konten: '',
          kategori: (KATEGORI_LIST.includes(a.kategori as KategoriArtikel) ? a.kategori : 'Gizi') as KategoriArtikel,
          penulis: String(a.nama_penulis ?? ''),
          gambar: '',
          waktuBaca: '5 Menit Baca',
          featured: false,
          tanggal: String(a.tanggal_publish ?? ''),
          status: (a as { status_artikel?: string }).status_artikel === 'Menunggu Verifikasi' ? 'pending' : 'published',
          rolePenulis: 'Dinas Kesehatan',
        }));
        setArtikelList(mapped);
      })
      .catch(() => notifyRef.current.error('Gagal memuat artikel. Pastikan backend terhubung.'))
      .finally(() => setArtikelLoading(false));
  }, [isDinkes, setArtikelList, setArtikelLoading]);

  useEffect(() => {
    fetchArtikel();
  }, [fetchArtikel]);

  const [kategoriAktif, setKategoriAktif] = useState<KategoriArtikel>('Semua');
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => { setPage(1); }, [kategoriAktif, searchQuery]);

  // Modal targets
  const [modalTambah, setModalTambah] = useState(false);
  const [modalEdit, setModalEdit] = useState<Artikel | null>(null);
  const [modalHapus, setModalHapus] = useState<Artikel | null>(null);
  const [modalDetail, setModalDetail] = useState<Artikel | null>(null);
  const [modalVerifikasi, setModalVerifikasi] = useState<Artikel | null>(null);

  // ── Derived role flags ─────────────────────────────────────────────────

  // ── Permission helpers ─────────────────────────────────────────────────
  const canEditArtikel = (a: Artikel): boolean => {
    if (isDinkes) return true;
    if (isBidan && a.rolePenulis === 'Bidan') return true;
    return false;
  };

  // ── Derived data ───────────────────────────────────────────────────────
  const filteredArtikel = useMemo(() => {
    const visible = isDinkes ? artikelList : artikelList.filter((a) => a.status === 'published');
    return visible.filter((a) => {
      const matchKat = kategoriAktif === 'Semua' || a.kategori === kategoriAktif;
      const matchQ = searchQuery === '' || a.judul.toLowerCase().includes(searchQuery.toLowerCase()) || a.ringkasan.toLowerCase().includes(searchQuery.toLowerCase());
      return matchKat && matchQ;
    });
  }, [artikelList, isDinkes, kategoriAktif, searchQuery]);

  const totalFiltered = filteredArtikel.length;
  const totalPages = Math.max(1, Math.ceil(totalFiltered / PAGE_SIZE));
  const paginatedArtikel = useMemo(() => {
    const start = (page - 1) * PAGE_SIZE;
    return filteredArtikel.slice(start, start + PAGE_SIZE);
  }, [filteredArtikel, page]);

  const featuredArtikel = paginatedArtikel.filter((a) => a.featured && a.status === 'published');
  const gridArtikel = paginatedArtikel.filter((a) => !a.featured || a.status !== 'published');

  const pendingCount = artikelList.filter((a) => a.status === 'pending').length;

  // ── Handlers ───────────────────────────────────────────────────────────
  const handleTambah = () => {
    setModalTambah(false);
    notify.success('Artikel berhasil ditambahkan');
    setPage(1);
    fetchArtikel();
  };

  const handleEdit = () => {
    setModalEdit(null);
    notify.success('Artikel berhasil diperbarui');
    fetchArtikel();
  };

  const handleHapus = async (id: string) => {
    const idNum = parseInt(String(id).replace(/[^0-9]/g, ''), 10) || 0;
    if (!idNum) { notify.error('ID artikel tidak valid.'); return; }
    try {
      await apiDelete('/artikel/' + idNum);
      setModalHapus(null);
      notify.success('Artikel berhasil dihapus');
      fetchArtikel();
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
      notify.success(action === 'approve' ? 'Artikel berhasil disetujui' : 'Artikel berhasil ditolak');
      fetchArtikel();
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
    showAll: true,
    totalGrid: 0,
    kategoriAktif,
    searchQuery,
    onChangeKategori: setKategoriAktif,
    onChangeSearch: setSearchQuery,
    onRead: setModalDetail,
    onToggleShowAll: () => {},
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

      {/* Guest / unauthenticated — read-only article grid */}
      {!currentRole && <IbuWaliSection {...commonProps} />}

      {/* Empty state */}
      {filteredArtikel.length === 0 && (
        <div className="text-center py-16 text-neutral-400">
          <div className="text-5xl mb-3">📚</div>
          <p className="font-semibold text-neutral-500">Tidak ada artikel ditemukan</p>
          <p className="text-sm mt-1">Coba ubah filter atau kata kunci pencarian</p>
        </div>
      )}

      {/* Paginator */}
      <Paginator
        page={page}
        totalPages={totalPages}
        totalItems={totalFiltered}
        pageSize={PAGE_SIZE}
        onPageChange={setPage}
      />

      {/* ── Shared modals ─────────────────────────────────────────────── */}
      {modalDetail && (
        <ModalDetailArtikel artikel={modalDetail} onClose={() => setModalDetail(null)} />
      )}

      {modalTambah && currentRole && (
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
