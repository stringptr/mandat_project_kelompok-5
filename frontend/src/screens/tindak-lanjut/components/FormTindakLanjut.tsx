import { useState, useRef, useEffect } from 'react';
import {
    Calendar, ChevronDown, Send, User, Activity, Search, FileEdit,
    X, CheckCircle2, AlertTriangle, Stethoscope
} from 'lucide-react';
import { useNotification } from '../../../context/NotificationContext';
import { apiGet, apiPost } from '../../../lib/api';
import type { FaskesItem } from '../../../types/entities';

// ─── Types ───────────────────────────────────────────────────────────
export interface Pasien {
    id: string;
    nama: string;
    umur: string;
    statusGizi: string;
    statusPasien: string;
    bb: number;
    tb: number;
    urgency: 'Rendah' | 'Sedang' | 'Tinggi' | 'Mendesak';
    avatarUrl?: string;
    namaIbu?: string;
    tanggalPeriksa: string;
}

interface FormTindakLanjutProps {
    onSubmit?: () => void;
    onCancel?: () => void;
}

export interface FormData {
    pasien: Pasien;
    jenisTindakan: string;
    lokasiFaskes: string;
    tanggalTarget: string;
    catatanMedis: string;
}

interface PasienSearchResult {
    id_pasien: number;
    nama: string;
    nik: string;
    jenis_kelamin: string;
    umur: string;
    nama_posyandu: string;
    jenis_pasien: string;
    status_kehamilan: string | null;
}

// ─── Urgency Badge Helper ───────────────────────────────────────────
const getUrgencyClasses = (urgency: string) => {
    switch (urgency) {
        case 'Mendesak': return 'bg-red-50 text-red-700 border-red-200';
        case 'Tinggi': return 'bg-amber-50 text-amber-700 border-amber-200';
        case 'Sedang': return 'bg-sky-50 text-sky-700 border-sky-200';
        case 'Rendah': return 'bg-emerald-50 text-emerald-700 border-emerald-200';
        default: return 'bg-gray-50 text-gray-600 border-gray-200';
    }
};

const getUrgencyIcon = (urgency: string) => {
    switch (urgency) {
        case 'Mendesak': return <AlertTriangle size={14} className="text-red-500" />;
        case 'Tinggi': return <AlertTriangle size={14} className="text-amber-500" />;
        case 'Sedang': return <Stethoscope size={14} className="text-sky-500" />;
        case 'Rendah': return <CheckCircle2 size={14} className="text-emerald-500" />;
        default: return null;
    }
};

// ─── Main Component ────────────────────────────────────────────────
export default function FormTindakLanjut({ onSubmit, onCancel }: FormTindakLanjutProps): JSX.Element {
    const notify = useNotification();
    const [searchQuery, setSearchQuery] = useState('');
    const [showDropdown, setShowDropdown] = useState(false);
    const [selectedPatient, setSelectedPatient] = useState<Pasien | null>(null);
    const [selectedFaskesId, setSelectedFaskesId] = useState<number | string>('');
    const [faskesList, setFaskesList] = useState<FaskesItem[]>([]);
    const [jenisTindakan, setJenisTindakan] = useState('Rujukan');
    const [tanggalTarget, setTanggalTarget] = useState('');
    const [catatanMedis, setCatatanMedis] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [errors, setErrors] = useState<{ [key: string]: string }>({});
    const [searchResults, setSearchResults] = useState<Pasien[]>([]);

    const dropdownRef = useRef<HTMLDivElement>(null);
    const inputRef = useRef<HTMLInputElement>(null);

    // Search pasien via API
    useEffect(() => {
        if (!searchQuery.trim()) {
            setSearchResults([]);
            return;
        }
        const timer = setTimeout(async () => {
            try {
                const res = await apiGet<PasienSearchResult[]>(`/monitoring/pasien/search?q=${encodeURIComponent(searchQuery)}`);
                setSearchResults(res.map((p) => ({
                    id: `PST-${p.id_pasien}`,
                    nama: p.nama,
                    umur: p.umur,
                    statusGizi: '',
                    statusPasien: 'Perlu Rujukan',
                    bb: 0,
                    tb: 0,
                    urgency: 'Sedang' as const,
                    tanggalPeriksa: '',
                })));
            } catch {
                setSearchResults([]);
            }
        }, 300);
        return () => clearTimeout(timer);
    }, [searchQuery]);

    // Fetch faskes list
    useEffect(() => {
        apiGet<FaskesItem[]>('/faskes')
            .then((data) => setFaskesList(Array.isArray(data) ? data : []))
            .catch(() => setFaskesList([]));
    }, []);

    const filteredPatients = searchQuery.trim() ? searchResults : searchResults;

    // Close dropdown on click outside
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setShowDropdown(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleSelectPatient = (patient: Pasien) => {
        setSelectedPatient(patient);
        setSearchQuery(patient.nama);
        setShowDropdown(false);
    };

    const handleClearPatient = () => {
        setSelectedPatient(null);
        setSearchQuery('');
        setShowDropdown(false);
        inputRef.current?.focus();
    };

    const handleSubmit = async () => {
        // Validation
        const newErrors: { [key: string]: string } = {};
        if (!selectedPatient) newErrors['selectedPatient'] = 'ERR-VAL-05'; // Nama pasien harus dipilih
        if (!jenisTindakan) newErrors['jenisTindakan'] = 'ERR-VAL-07';
        if (jenisTindakan === 'Rujukan' && !selectedFaskesId) newErrors['lokasiFaskes'] = 'ERR-VAL-07';
        if (!tanggalTarget) newErrors['tanggalTarget'] = 'ERR-VAL-06';
        if (!catatanMedis.trim() || catatanMedis.trim().length < 10) newErrors['catatanMedis'] = 'ERR-VAL-03';
        setErrors(newErrors);
        if (Object.keys(newErrors).length > 0) {
            notify.warn('Mohon lengkapi semua data form yang wajib diisi sebelum mengirim.');
            return;
        }
        setIsSubmitting(true);
        try {
            await apiPost('/tindak-lanjut', {
                id_hasil_pemeriksaan: 0, // placeholder - backend accepts id_pasien as fallback
                id_pasien: parseInt(String(selectedPatient!.id).replace(/[^0-9]/g, ''), 10) || undefined,
                jenis_tindakan: jenisTindakan === 'Tindak Lanjut' ? 'Kontrol Ulang' : jenisTindakan,
                jadwal_kontrol: tanggalTarget,
                catatan_medis: catatanMedis,
                alasan_rujukan: jenisTindakan === 'Rujukan' ? catatanMedis : undefined,
                rekomendasi: jenisTindakan !== 'Rujukan' ? catatanMedis : undefined,
                id_faskes: selectedFaskesId ? Number(selectedFaskesId) : undefined,
            });
            onSubmit?.();
            setIsSubmitting(false);
            notify.success('Data tindak lanjut berhasil dikirim');
            // Reset form
            setSelectedPatient(null);
            setSearchQuery('');
            setSelectedFaskesId('');
            setTanggalTarget('');
            setCatatanMedis('');
        } catch {
            setIsSubmitting(false);
            notify.error('Gagal mengirim tindak lanjut. Silakan coba lagi.');
        }
    };

    return (
        <div className="bg-white rounded-2xl p-6 shadow-sm border border-slate-100">
            {/* Form Header */}
            <div className="flex items-start gap-4 mb-8">
                <div className="bg-emerald-50 p-3 rounded-xl flex items-center justify-center">
                    <FileEdit className="w-6 h-6 text-emerald-700" />
                </div>
                <div>
                    <h2 className="text-xl font-bold text-slate-800">
                        {selectedPatient ? 'Formulir Tindak Lanjut' : 'Formulir Tindak Lanjut Baru'}
                    </h2>
                    <p className="text-slate-500 text-sm mt-1">
                        {selectedPatient
                            ? 'Detailkan data rujukan atau penjadwalan kontrol'
                            : 'Input data rujukan atau penjadwalan ulang'}
                    </p>
                </div>
            </div>

            {/* ─── Search Pasien ─────────────────────────────────────────── */}
            <div className="mb-6 relative" ref={dropdownRef}>
                <label className="block text-sm font-semibold text-slate-800 mb-2">
                    Cari Pasien (Nama Anak, Nama Ibu, atau ID)
                </label>
                <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                        <Search className="w-4 h-4 text-slate-400" />
                    </div>
                    <input
  ref={inputRef}
  type="text"
  value={searchQuery}
  onChange={(e) => {
    setSearchQuery(e.target.value);
    setShowDropdown(true);
    if (selectedPatient && e.target.value !== selectedPatient.nama) {
      setSelectedPatient(null);
    }
  }}
  onFocus={() => setShowDropdown(true)}
  placeholder="Ketik nama pasien..."
  className={`w-full bg-slate-50 border ${errors['selectedPatient'] ? 'border-red-500' : 'border-slate-200'} focus:border-emerald-300 focus:bg-white focus:ring-2 focus:ring-emerald-100 rounded-xl pl-11 pr-10 py-3 text-sm text-slate-700 outline-none transition-all placeholder:text-slate-400`}
/>
                    {searchQuery && (
                        <button
                            onClick={handleClearPatient}
                            className="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-600"
                        >
                            <X size={16} />
                        </button>
                    )}
                    {errors['selectedPatient'] && (
                        <p className="text-xs text-red-600 mt-1">{errors['selectedPatient']}</p>
                    )}
                </div>

                {/* ─── Dropdown Hasil Pencarian ──────────────────────────── */}
                {showDropdown && filteredPatients.length > 0 && (
                    <div className="absolute z-50 w-full mt-2 bg-white rounded-xl border border-slate-200 shadow-xl overflow-hidden">
                        <div className="px-4 py-2 bg-slate-50 border-b border-slate-100">
                            <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
                                {searchQuery.trim() ? `Hasil Pencarian (${filteredPatients.length})` : 'Pasien Butuh Rujukan'}
                            </span>
                        </div>
                        <div className="max-h-80 overflow-y-auto">
                            {filteredPatients.map((patient) => (
                                <button
                                    key={patient.id}
                                    onClick={() => handleSelectPatient(patient)}
                                    className="w-full text-left px-4 py-3 hover:bg-slate-50 transition-colors border-b border-slate-50 last:border-b-0"
                                >
                                    <div className="flex items-start gap-3">
                                        {/* Avatar */}
                                        <div className="w-10 h-10 rounded-full overflow-hidden flex-shrink-0 bg-slate-200 border border-slate-100">
                                            {patient.avatarUrl ? (
                                                <img src={patient.avatarUrl} alt={patient.nama} className="w-full h-full object-cover" />
                                            ) : (
                                                <div className="w-full h-full flex items-center justify-center bg-emerald-100 text-emerald-600 font-bold text-xs">
                                                    {patient.nama.charAt(0)}
                                                </div>
                                            )}
                                        </div>

                                        {/* Info */}
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-center justify-between mb-0.5">
                                                <h4 className="text-sm font-bold text-slate-800 truncate">{patient.nama}</h4>
                                                <span className={`text-[10px] font-semibold px-2 py-0.5 rounded-full border ${getUrgencyClasses(patient.urgency)}`}>
                                                    {patient.urgency}
                                                </span>
                                            </div>
                                            <p className="text-xs text-slate-500 mb-1">
                                                ID: {patient.id} · {patient.umur}
                                            </p>
                                            <div className="flex items-center gap-2">
                                                <span className="text-xs text-red-600 bg-red-50 px-1.5 py-0.5 rounded font-medium">
                                                    {patient.statusGizi}
                                                </span>
                                                <span className="text-xs text-slate-400">
                                                    BB: {patient.bb}kg | TB: {patient.tb}cm
                                                </span>
                                            </div>
                                        </div>
                                    </div>
                                </button>
                            ))}
                        </div>
                    </div>
                )}

                {/* Empty state dropdown */}
                {showDropdown && searchQuery.trim() && filteredPatients.length === 0 && (
                    <div className="absolute z-50 w-full mt-2 bg-white rounded-xl border border-slate-200 shadow-xl p-6 text-center">
                        <Search size={24} className="text-slate-300 mx-auto mb-2" />
                        <p className="text-sm text-slate-500">Tidak ada pasien ditemukan</p>
                        <p className="text-xs text-slate-400 mt-1">Coba kata kunci lain</p>
                    </div>
                )}
            </div>

            {/* ─── Selected Patient Card ─────────────────────────────────── */}
            {selectedPatient && (
                <div className="mb-8 bg-slate-50 border border-slate-200 rounded-xl p-4">
                    <div className="flex items-center justify-between mb-3">
                        <h3 className="text-sm font-semibold text-slate-500 uppercase tracking-wider">
                            Hasil Monitoring Terakhir
                        </h3>
                        <span className={`inline-flex items-center gap-1 text-[10px] font-semibold px-2 py-0.5 rounded-full border ${getUrgencyClasses(selectedPatient.urgency)}`}>
                            {getUrgencyIcon(selectedPatient.urgency)}
                            {selectedPatient.urgency}
                        </span>
                    </div>

                    <div className="flex flex-col sm:flex-row sm:items-center gap-4">
                        <div className="flex items-center gap-3">
                            <div className="w-12 h-12 rounded-full overflow-hidden bg-slate-200 border-2 border-white shadow-sm">
                                {selectedPatient.avatarUrl ? (
                                    <img src={selectedPatient.avatarUrl} alt={selectedPatient.nama} className="w-full h-full object-cover" />
                                ) : (
                                    <div className="w-full h-full flex items-center justify-center bg-emerald-100 text-emerald-600 font-bold">
                                        {selectedPatient.nama.charAt(0)}
                                    </div>
                                )}
                            </div>
                            <div>
                                <h4 className="font-bold text-slate-800 text-base">{selectedPatient.nama}</h4>
                                <div className="flex items-center gap-3 text-xs text-slate-500 mt-0.5">
                                    <span className="flex items-center gap-1"><User size={12} /> {selectedPatient.umur}</span>
                                    <span className="flex items-center gap-1"><Activity size={12} /> {selectedPatient.statusGizi}</span>
                                </div>
                            </div>
                        </div>

                        <div className="sm:ml-auto sm:text-right sm:border-l sm:border-slate-200 sm:pl-4 pt-3 sm:pt-0 border-t sm:border-t-0 border-slate-200">
                            <p className="text-xs text-slate-500 mb-1">Pemeriksaan: {selectedPatient.tanggalPeriksa}</p>
                            <p className="text-sm font-medium text-slate-800">
                                BB: <span className="font-bold">{selectedPatient.bb} kg</span>
                                <span className="mx-2 text-slate-300">|</span>
                                TB: <span className="font-bold">{selectedPatient.tb} cm</span>
                                <span className="ml-2 font-semibold text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-md">
                                    ({selectedPatient.statusGizi})
                                </span>
                            </p>
                        </div>
                    </div>
                </div>
            )}

            {/* ─── Form Fields ───────────────────────────────────────────── */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                <div>
                    <label className="block text-sm font-semibold text-slate-800 mb-2">Jenis Tindakan</label>
                    <div className="relative">
                        <select
  value={jenisTindakan}
  onChange={(e) => setJenisTindakan(e.target.value)}
  className={`w-full bg-slate-50 border ${errors['jenisTindakan'] ? 'border-red-500' : 'border-slate-200'} focus:border-emerald-300 focus:bg-white focus:ring-2 focus:ring-emerald-100 rounded-xl px-4 py-3 text-sm text-slate-700 outline-none transition-all appearance-none cursor-pointer`}
>
  <option value="Rujukan">Rujukan</option>
  <option value="Tindak Lanjut">Tindak Lanjut</option>
</select>
                        {errors['jenisTindakan'] && (
                            <p className="text-xs text-red-600 mt-1">{errors['jenisTindakan']}</p>
                        )}
                        <ChevronDown className="w-4 h-4 text-slate-500 absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none" />
                    </div>
                </div>

                {jenisTindakan === 'Rujukan' && (
                    <div>
                        <label className="block text-sm font-semibold text-slate-800 mb-2">Lokasi Faskes</label>
                        <div className="relative">
                            <select
                                value={selectedFaskesId}
                                onChange={(e) => setSelectedFaskesId(e.target.value)}
                                className={`w-full bg-slate-50 border ${errors['lokasiFaskes'] ? 'border-red-500' : 'border-slate-200'} focus:border-emerald-300 focus:bg-white focus:ring-2 focus:ring-emerald-100 rounded-xl px-4 py-3 text-sm text-slate-700 outline-none transition-all appearance-none cursor-pointer`}
                            >
                                <option value="">-- Pilih Fasilitas Kesehatan --</option>
                                {faskesList.map((f) => (
                                    <option key={f.id_faskes} value={f.id_faskes}>
                                        {f.nama_faskes} ({f.tipe_faskes})
                                    </option>
                                ))}
                            </select>
                            <ChevronDown className="w-4 h-4 text-slate-500 absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none" />
                        </div>
                        {errors['lokasiFaskes'] && (
                            <p className="text-xs text-red-600 mt-1">{errors['lokasiFaskes']}</p>
                        )}
                    </div>
                )}
            </div>

            <div className="mb-6">
                <label className="block text-sm font-semibold text-slate-800 mb-2">Tanggal Target Pelaksanaan</label>
                <div className="relative md:w-1/2 md:pr-3">
                    <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                        <Calendar className="w-4 h-4 text-slate-400" />
                    </div>
                    <input
  type="date"
  value={tanggalTarget}
  onChange={(e) => setTanggalTarget(e.target.value)}
  className={`w-full bg-slate-50 border ${errors['tanggalTarget'] ? 'border-red-500' : 'border-slate-200'} focus:border-emerald-300 focus:bg-white focus:ring-2 focus:ring-emerald-100 rounded-xl pl-11 pr-4 py-3 text-sm text-slate-700 outline-none transition-all`}
/>
                </div>
                {errors['tanggalTarget'] && (
                    <p className="text-xs text-red-600 mt-1">{errors['tanggalTarget']}</p>
                )}
            </div>

            <div className="mb-8">
                {jenisTindakan === 'Rujukan' ? (
                    <>
                        <label className="block text-sm font-semibold text-slate-800 mb-2">Alasan Rujukan</label>
                        <textarea
  rows={4}
  value={catatanMedis}
  onChange={(e) => setCatatanMedis(e.target.value)}
                           placeholder={selectedPatient
                                ? `Contoh: ${selectedPatient.nama} membutuhkan ${jenisTindakan === 'Rujukan' ? 'rujukan' : 'rekomendasi'} ke fasilitas karena ${selectedPatient.statusGizi.toLowerCase()}...`
                                : 'Tuliskan catatan medik...'}
  className={`w-full bg-slate-50 border ${errors['catatanMedis'] ? 'border-red-500' : 'border-slate-200'} focus:border-emerald-300 focus:bg-white focus:ring-2 focus:ring-emerald-100 rounded-xl px-4 py-3 text-sm text-slate-700 outline-none transition-all resize-none placeholder:text-slate-400`}
/>
                        {errors['catatanMedis'] && (
                            <p className="text-xs text-red-600 mt-1">{errors['catatanMedis']}</p>
                        )}
                    </>
                ) : (
                    <>
                        <label className="block text-sm font-semibold text-slate-800 mb-2">Catatan Medis & Rekomendasi</label>
                        <textarea
                            rows={4}
                            value={catatanMedis}
                            onChange={(e) => setCatatanMedis(e.target.value)}
                            placeholder={selectedPatient
                                ? `Contoh: ${selectedPatient.nama} membutuhkan rekomendasi penanganan karena ${selectedPatient.statusGizi.toLowerCase()}...`
                                : 'Tuliskan catatan medis dan rekomendasi...'}
                            className="w-full bg-slate-50 border border-slate-200 focus:border-emerald-300 focus:bg-white focus:ring-2 focus:ring-emerald-100 rounded-xl px-4 py-3 text-sm text-slate-700 outline-none transition-all resize-none placeholder:text-slate-400" />
                        {errors['catatanMedis'] && (
                            <p className="text-xs text-red-600 mt-1">{errors['catatanMedis']}</p>
                        )}
                    </>
                )}
            </div>

            {/* ─── Buttons ───────────────────────────────────────────────── */}
            <div className="flex items-center gap-4">
                <button
                    onClick={handleSubmit}
                    disabled={!selectedPatient || isSubmitting}
                    className={`bg-emerald-700 hover:bg-emerald-800 text-white font-semibold py-3 px-6 rounded-xl flex items-center justify-center gap-2 transition-colors text-sm
            ${(!selectedPatient || isSubmitting) ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
          `}
                >
                    {isSubmitting ? (
                        <>
                            <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                            Mengirim...
                        </>
                    ) : (
                        <>
                            <Send className="w-4 h-4" />
                            Kirim Tindak Lanjut
                        </>
                    )}
                </button>
                <button
                    onClick={onCancel}
                    className="text-slate-600 font-semibold py-3 px-6 hover:bg-slate-100 rounded-xl transition-colors cursor-pointer text-sm"
                >
                    Batal
                </button>
            </div>
        </div>
    );
}