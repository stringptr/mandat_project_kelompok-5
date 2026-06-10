import { useState, useMemo } from 'react';

// ─── Types ────────────────────────────────────────────────────────────────────

export interface PasienOption {
    id: string;
    nama: string;
    namaIbu: string;
    usia: string;
    nik: string;
    statusGizi?: string;
}

interface ModalTambahPemeriksaanProps {
    isOpen: boolean;
    onClose: () => void;
    /** Daftar pasien/balita yang bisa dipilih */
    pasienList: PasienOption[];
    onSubmit: (pasienId: string, namaAnak: string, data: FormDataTambah) => void;
}

export interface FormDataTambah {
    periodeBulan: string;
    beratBadan: string;
    tinggiBadan: string;
    lingkarKepala: string;
    tekananDarah: string;
    catatanMedis: string;
    rekomendasiGizi: string;
    jadwalKontrol: string;
    statusPasien: string;
}

// ─── Constants ────────────────────────────────────────────────────────────────

const PERIODE_OPTIONS = [
    'Bulan ke-24 (Saat ini)',
    'Bulan ke-25',
    'Bulan ke-26',
    'Bulan ke-27',
];

const STATUS_OPTIONS = [
    'Dalam Pemantauan',
    'Normal',
    'Perlu Tindak Lanjut',
    'Prioritas',
];

const EMPTY_FORM: FormDataTambah = {
    periodeBulan: PERIODE_OPTIONS[0],
    beratBadan: '',
    tinggiBadan: '',
    lingkarKepala: '',
    tekananDarah: '',
    catatanMedis: '',
    rekomendasiGizi: '',
    jadwalKontrol: '',
    statusPasien: STATUS_OPTIONS[0],
};

function validatePositiveDecimal(value: string): string | null {
    if (!value.trim()) return 'Kolom ini wajib diisi.';
    const num = parseFloat(value);
    if (isNaN(num) || num <= 0) return 'Input harus berupa angka desimal positif.';
    return null;
}

const STATUS_GIZI_COLOR: Record<string, string> = {
    'Stunting': 'bg-red-50 text-red-600',
    'Gizi Kurang': 'bg-amber-50 text-amber-600',
    'Gizi Baik': 'bg-emerald-50 text-emerald-700',
    'Prioritas': 'bg-red-50 text-red-600',
};

// ─── Main Component ───────────────────────────────────────────────────────────

export default function ModalTambahPemeriksaan({
    isOpen,
    onClose,
    pasienList,
    onSubmit,
}: ModalTambahPemeriksaanProps) {
    // Step 1 = pilih pasien, Step 2 = isi form
    const [step, setStep] = useState<1 | 2>(1);
    const [searchPasien, setSearchPasien] = useState('');
    const [selectedPasien, setSelectedPasien] = useState<PasienOption | null>(null);

    const [form, setForm] = useState<FormDataTambah>(EMPTY_FORM);
    const [errors, setErrors] = useState<Partial<Record<keyof FormDataTambah, string>>>({});
    const [submitted, setSubmitted] = useState(false);

    if (!isOpen) return null;

    // ── Step 1: filter pasien ──
    const filteredPasien = useMemo(() => {
        const q = searchPasien.toLowerCase().trim();
        if (!q) return pasienList;
        return pasienList.filter(
            p =>
                p.nama.toLowerCase().includes(q) ||
                p.namaIbu.toLowerCase().includes(q) ||
                p.nik.includes(q),
        );
    }, [searchPasien, pasienList]);

    const handleSelectPasien = (p: PasienOption) => {
        setSelectedPasien(p);
        setStep(2);
    };

    // ── Step 2: form handlers ──
    const setField = (key: keyof FormDataTambah, value: string) => {
        setForm(prev => ({ ...prev, [key]: value }));
        if (submitted) {
            if (key === 'beratBadan' || key === 'tinggiBadan') {
                setErrors(prev => ({ ...prev, [key]: validatePositiveDecimal(value) ?? undefined }));
            }
        }
    };

    const handleSubmit = () => {
        const bbErr = validatePositiveDecimal(form.beratBadan);
        const tbErr = validatePositiveDecimal(form.tinggiBadan);
        setErrors({ beratBadan: bbErr ?? undefined, tinggiBadan: tbErr ?? undefined });
        setSubmitted(true);
        if (!bbErr && !tbErr && selectedPasien) {
            onSubmit(selectedPasien.id, selectedPasien.nama, form);
            handleClose();
        }
    };

    const handleClose = () => {
        setStep(1);
        setSearchPasien('');
        setSelectedPasien(null);
        setForm(EMPTY_FORM);
        setErrors({});
        setSubmitted(false);
        onClose();
    };

    const handleBack = () => {
        setStep(1);
        setSelectedPasien(null);
        setForm(EMPTY_FORM);
        setErrors({});
        setSubmitted(false);
    };

    // ── Shared header ──
    const modalTitle = step === 1
        ? 'Tambah Data Pemeriksaan'
        : `Tambah Pemeriksaan – ${selectedPasien?.nama}`;
    const modalSub = step === 1
        ? 'Cari dan pilih pasien terlebih dahulu'
        : 'Monitoring Gizi Berkala';

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
            <div className="bg-white rounded-2xl w-full max-w-lg shadow-xl overflow-hidden">

                {/* ── Header ── */}
                <div className="flex items-start gap-3 px-5 py-4 border-b border-neutral-100">
                    <div className="w-9 h-9 rounded-xl bg-blue-50 flex items-center justify-center flex-shrink-0">
                        <svg className="w-5 h-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.8}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
                        </svg>
                    </div>
                    <div className="flex-1 min-w-0">
                        <p className="text-sm font-semibold text-neutral-800 truncate">{modalTitle}</p>
                        <p className="text-xs text-neutral-400 mt-0.5">{modalSub}</p>
                    </div>
                    {/* Step indicator */}
                    <div className="flex items-center gap-1.5 mr-2">
                        <StepDot active={step === 1} done={step > 1} label="1" />
                        <div className="w-4 h-px bg-neutral-200" />
                        <StepDot active={step === 2} done={false} label="2" />
                    </div>
                    <button onClick={handleClose} className="w-7 h-7 flex items-center justify-center rounded-lg text-neutral-400 hover:bg-neutral-100 hover:text-neutral-600 transition-colors flex-shrink-0">
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                {/* ── Step 1: Pilih Pasien ── */}
                {step === 1 && (
                    <>
                        <div className="px-5 py-4">
                            {/* Search */}
                            <div className="relative mb-3">
                                <svg className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-neutral-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                                </svg>
                                <input
                                    type="text"
                                    placeholder="Cari nama balita, nama ibu, atau NIK..."
                                    value={searchPasien}
                                    onChange={e => setSearchPasien(e.target.value)}
                                    className="w-full pl-9 pr-3 py-2 text-sm rounded-xl border border-neutral-200 bg-neutral-50 text-neutral-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500"
                                    autoFocus
                                />
                                {searchPasien && (
                                    <button
                                        onClick={() => setSearchPasien('')}
                                        className="absolute right-3 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-neutral-600"
                                    >
                                        <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                                        </svg>
                                    </button>
                                )}
                            </div>

                            {/* Pasien list */}
                            <div className="max-h-72 overflow-y-auto space-y-1.5 pr-0.5">
                                {filteredPasien.length === 0 ? (
                                    <div className="py-10 text-center text-sm text-neutral-400">
                                        <svg className="w-8 h-8 mx-auto mb-2 text-neutral-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                                        </svg>
                                        Pasien tidak ditemukan
                                    </div>
                                ) : (
                                    filteredPasien.map(p => (
                                        <button
                                            key={p.id}
                                            onClick={() => handleSelectPasien(p)}
                                            className="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-emerald-50 border border-transparent hover:border-emerald-100 transition-all text-left group"
                                        >
                                            {/* Avatar inisial */}
                                            <div className="w-9 h-9 rounded-full bg-emerald-100 flex items-center justify-center text-[11px] font-semibold text-emerald-800 flex-shrink-0">
                                                {p.nama.split(' ').map(n => n[0]).slice(0, 2).join('')}
                                            </div>
                                            <div className="flex-1 min-w-0">
                                                <p className="text-sm font-semibold text-neutral-800 truncate">{p.nama}</p>
                                                <p className="text-xs text-neutral-400 truncate">
                                                    Ibu: {p.namaIbu} · {p.usia}
                                                </p>
                                            </div>
                                            <div className="flex items-center gap-2 flex-shrink-0">
                                                {p.statusGizi && (
                                                    <span className={`text-[10px] font-medium px-2 py-0.5 rounded-full ${STATUS_GIZI_COLOR[p.statusGizi] ?? 'bg-neutral-100 text-neutral-500'}`}>
                                                        {p.statusGizi}
                                                    </span>
                                                )}
                                                <svg className="w-4 h-4 text-neutral-300 group-hover:text-emerald-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                                    <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
                                                </svg>
                                            </div>
                                        </button>
                                    ))
                                )}
                            </div>

                            <p className="text-[11px] text-neutral-400 mt-3 text-center">
                                Menampilkan {filteredPasien.length} dari {pasienList.length} pasien
                            </p>
                        </div>

                        <div className="flex justify-end px-5 py-4 border-t border-neutral-100">
                            <button onClick={handleClose} className="text-sm text-neutral-500 hover:text-neutral-700 px-3 py-2 rounded-xl hover:bg-neutral-100 transition-colors">
                                Batal
                            </button>
                        </div>
                    </>
                )}

                {/* ── Step 2: Form Pengukuran ── */}
                {step === 2 && selectedPasien && (
                    <>
                        {/* Pasien terpilih */}
                        <div className="mx-5 mt-4 flex items-center gap-3 bg-emerald-50 border border-emerald-100 rounded-xl px-3 py-2.5">
                            <div className="w-8 h-8 rounded-full bg-emerald-200 flex items-center justify-center text-[10px] font-semibold text-emerald-900 flex-shrink-0">
                                {selectedPasien.nama.split(' ').map(n => n[0]).slice(0, 2).join('')}
                            </div>
                            <div className="flex-1 min-w-0">
                                <p className="text-sm font-semibold text-emerald-900 truncate">{selectedPasien.nama}</p>
                                <p className="text-xs text-emerald-600">NIK: {selectedPasien.nik} · {selectedPasien.usia}</p>
                            </div>
                            <button
                                onClick={handleBack}
                                className="text-[11px] text-emerald-700 hover:text-emerald-900 font-medium flex items-center gap-1 hover:underline"
                            >
                                <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
                                </svg>
                                Ganti
                            </button>
                        </div>

                        {/* Form body */}
                        <div className="px-5 py-4 space-y-4 max-h-[55vh] overflow-y-auto">
                            {/* Periode */}
                            <div>
                                <label className="block text-xs text-neutral-500 mb-1.5">
                                    Jadwal / Periode Bulan <span className="text-red-500">*</span>
                                </label>
                                <select
                                    value={form.periodeBulan}
                                    onChange={e => setField('periodeBulan', e.target.value)}
                                    className="w-full px-3 py-2 text-sm rounded-xl border border-neutral-200 bg-neutral-50 text-neutral-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500"
                                >
                                    {PERIODE_OPTIONS.map(opt => <option key={opt}>{opt}</option>)}
                                </select>
                            </div>

                            {/* Parameter Fisik */}
                            <div>
                                <SectionLabel>Parameter Fisik</SectionLabel>
                                <div className="grid grid-cols-2 gap-3">
                                    <InputField
                                        label="Berat Badan (BB)" required
                                        value={form.beratBadan} onChange={v => setField('beratBadan', v)}
                                        unit="kg" error={errors.beratBadan} placeholder="0.0"
                                    />
                                    <InputField
                                        label="Tinggi Badan (TB/PB)" required
                                        value={form.tinggiBadan} onChange={v => setField('tinggiBadan', v)}
                                        unit="cm" error={errors.tinggiBadan} placeholder="0.0"
                                    />
                                    <div>
                                        <InputField
                                            label="Lingkar Kepala"
                                            value={form.lingkarKepala} onChange={v => setField('lingkarKepala', v)}
                                            unit="cm" placeholder="0.0"
                                        />
                                        <p className="text-[11px] text-neutral-400 mt-1">Opsional untuk anak &gt; 2 tahun</p>
                                    </div>
                                    <InputField
                                        label="Tekanan Darah"
                                        value={form.tekananDarah} onChange={v => setField('tekananDarah', v)}
                                        placeholder="Contoh: 120/80"
                                    />
                                </div>
                            </div>

                            {/* Verifikasi Medis */}
                            <div>
                                <SectionLabel>Verifikasi Medis</SectionLabel>
                                <div className="grid grid-cols-2 gap-3">
                                    <div>
                                        <label className="block text-xs text-neutral-500 mb-1.5">Catatan Medis</label>
                                        <textarea
                                            value={form.catatanMedis}
                                            onChange={e => setField('catatanMedis', e.target.value)}
                                            rows={3} placeholder="Masukkan observasi medis..."
                                            className="w-full px-3 py-2 text-sm rounded-xl border border-neutral-200 bg-neutral-50 text-neutral-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500 resize-none"
                                        />
                                        <p className="text-[11px] text-neutral-400 text-right">{form.catatanMedis.length}/1000</p>
                                    </div>
                                    <div>
                                        <label className="block text-xs text-neutral-500 mb-1.5">Rekomendasi Gizi</label>
                                        <textarea
                                            value={form.rekomendasiGizi}
                                            onChange={e => setField('rekomendasiGizi', e.target.value)}
                                            rows={3} placeholder="Berikan saran asupan gizi..."
                                            className="w-full px-3 py-2 text-sm rounded-xl border border-neutral-200 bg-neutral-50 text-neutral-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500 resize-none"
                                        />
                                        <p className="text-[11px] text-neutral-400 text-right">{form.rekomendasiGizi.length}/1000</p>
                                    </div>
                                </div>
                                <div className="grid grid-cols-2 gap-3 mt-3">
                                    <div>
                                        <label className="block text-xs text-neutral-500 mb-1.5">Jadwal Kontrol Berikutnya</label>
                                        <input
                                            type="date" value={form.jadwalKontrol}
                                            onChange={e => setField('jadwalKontrol', e.target.value)}
                                            className="w-full px-3 py-2 text-sm rounded-xl border border-neutral-200 bg-neutral-50 text-neutral-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-xs text-neutral-500 mb-1.5">Status Pasien</label>
                                        <select
                                            value={form.statusPasien}
                                            onChange={e => setField('statusPasien', e.target.value)}
                                            className="w-full px-3 py-2 text-sm rounded-xl border border-neutral-200 bg-neutral-50 text-neutral-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500"
                                        >
                                            {STATUS_OPTIONS.map(opt => <option key={opt}>{opt}</option>)}
                                        </select>
                                    </div>
                                </div>
                            </div>
                        </div>

                        {/* Footer */}
                        <div className="flex items-center justify-between px-5 py-4 border-t border-neutral-100">
                            <button
                                onClick={handleBack}
                                className="flex items-center gap-1.5 text-sm text-neutral-500 hover:text-neutral-700 px-3 py-2 rounded-xl hover:bg-neutral-100 transition-colors"
                            >
                                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
                                </svg>
                                Kembali
                            </button>
                            <div className="flex items-center gap-2">
                                <button onClick={handleClose} className="text-sm text-neutral-500 hover:text-neutral-700 px-3 py-2 rounded-xl hover:bg-neutral-100 transition-colors">
                                    Batal
                                </button>
                                <button
                                    onClick={handleSubmit}
                                    className="flex items-center gap-2 bg-emerald-700 hover:bg-emerald-800 text-white text-sm font-medium px-4 py-2 rounded-xl transition-colors"
                                >
                                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                        <path strokeLinecap="round" strokeLinejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                    Simpan &amp; Kirim
                                </button>
                            </div>
                        </div>
                    </>
                )}
            </div>
        </div>
    );
}

// ─── Sub-components ────────────────────────────────────────────────────────────

function StepDot({ active, done, label }: { active: boolean; done: boolean; label: string }) {
    if (done) {
        return (
            <div className="w-5 h-5 rounded-full bg-emerald-600 flex items-center justify-center">
                <svg className="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                </svg>
            </div>
        );
    }
    return (
        <div className={`w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-semibold border ${active ? 'bg-emerald-700 border-emerald-700 text-white' : 'border-neutral-200 text-neutral-400 bg-white'}`}>
            {label}
        </div>
    );
}

function SectionLabel({ children }: { children: React.ReactNode }) {
    return (
        <div className="flex items-center gap-2 mb-3">
            <div className="h-0.5 w-5 bg-emerald-600 rounded-full" />
            <span className="text-[11px] font-semibold text-neutral-400 tracking-widest uppercase">{children}</span>
        </div>
    );
}

function InputField({
    label, value, onChange, unit, error, placeholder, required,
}: {
    label: string; value: string; onChange: (v: string) => void;
    unit?: string; error?: string; placeholder?: string; required?: boolean;
}) {
    return (
        <div>
            <label className="block text-xs text-neutral-500 mb-1.5">
                {label} {required && <span className="text-red-500">*</span>}
            </label>
            <div className="relative">
                <input
                    type="text" value={value}
                    onChange={e => onChange(e.target.value)}
                    placeholder={placeholder}
                    className={`w-full px-3 py-2 text-sm rounded-xl border bg-neutral-50 text-neutral-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500 transition-colors ${error ? 'border-red-400 bg-red-50' : 'border-neutral-200'} ${unit ? 'pr-12' : ''}`}
                />
                {unit && <span className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-neutral-400">{unit}</span>}
                {error && (
                    <svg className="absolute right-8 top-1/2 -translate-y-1/2 w-4 h-4 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <circle cx="12" cy="12" r="10" /><line x1="12" y1="8" x2="12" y2="12" /><line x1="12" y1="16" x2="12.01" y2="16" />
                    </svg>
                )}
            </div>
            {error && <p className="text-xs text-red-500 mt-1">{error}</p>}
        </div>
    );
}