import { useState, useEffect } from 'react';

export interface PemeriksaanData {
    id: string;
    beratBadan: string;
    tinggiBadan: string;
    lingkarKepala: string;
    terakhirDiperbarui: {
        nama: string;
        tanggal: string;
        inisial: string;
    };
}

interface ModalEditPemeriksaanProps {
    isOpen: boolean;
    onClose: () => void;
    namaAnak: string;
    data: PemeriksaanData | null;
    onSubmit: (id: string, data: Omit<PemeriksaanData, 'id' | 'terakhirDiperbarui'>) => void;
}

function validatePositiveDecimal(value: string): string | null {
    if (!value.trim()) return 'Kolom ini wajib diisi.';
    const num = parseFloat(value);
    if (isNaN(num) || num <= 0) return 'Input harus berupa angka desimal positif.';
    return null;
}

export default function ModalEditPemeriksaan({
    isOpen,
    onClose,
    namaAnak,
    data,
    onSubmit,
}: ModalEditPemeriksaanProps) {
    const [beratBadan, setBeratBadan] = useState('');
    const [tinggiBadan, setTinggiBadan] = useState('');
    const [lingkarKepala, setLingkarKepala] = useState('');
    const [errors, setErrors] = useState<{ beratBadan?: string; tinggiBadan?: string }>({});
    const [submitted, setSubmitted] = useState(false);

    useEffect(() => {
        if (data && isOpen) {
            setBeratBadan(data.beratBadan);
            setTinggiBadan(data.tinggiBadan);
            setLingkarKepala(data.lingkarKepala);
            setErrors({});
            setSubmitted(false);
        }
    }, [data, isOpen]);

    if (!isOpen || !data) return null;

    const handleBeratChange = (v: string) => {
        setBeratBadan(v);
        if (submitted) setErrors(prev => ({ ...prev, beratBadan: validatePositiveDecimal(v) ?? undefined }));
    };

    const handleTinggiChange = (v: string) => {
        setTinggiBadan(v);
        if (submitted) setErrors(prev => ({ ...prev, tinggiBadan: validatePositiveDecimal(v) ?? undefined }));
    };

    const handleSubmit = () => {
        const bbErr = validatePositiveDecimal(beratBadan);
        const tbErr = validatePositiveDecimal(tinggiBadan);
        const newErrors = { beratBadan: bbErr ?? undefined, tinggiBadan: tbErr ?? undefined };
        setErrors(newErrors);
        setSubmitted(true);
        if (!bbErr && !tbErr) {
            onSubmit(data.id, { beratBadan, tinggiBadan, lingkarKepala });
            onClose();
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
            <div className="bg-white rounded-2xl w-full max-w-md shadow-xl overflow-hidden">
                {/* Header */}
                <div className="flex items-start gap-3 px-5 py-4 border-b border-neutral-100">
                    <div className="w-9 h-9 rounded-xl bg-emerald-50 flex items-center justify-center flex-shrink-0">
                        <svg className="w-5 h-5 text-emerald-700" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.8}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                    </div>
                    <div className="flex-1 min-w-0">
                        <p className="text-sm font-semibold text-neutral-800 truncate">Edit Riwayat Pemeriksaan – {namaAnak}</p>
                        <p className="text-xs text-neutral-400 mt-0.5">Monitoring Gizi Berkala</p>
                    </div>
                    <button onClick={onClose} className="w-7 h-7 flex items-center justify-center rounded-lg text-neutral-400 hover:bg-neutral-100 hover:text-neutral-600 transition-colors flex-shrink-0">
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                {/* Body */}
                <div className="px-5 py-4 space-y-4">
                    {/* BB & TB */}
                    <div className="grid grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs text-neutral-500 mb-1.5">
                                Berat Badan (BB) <span className="text-red-500">*</span>
                            </label>
                            <div className="relative">
                                <input
                                    type="text"
                                    value={beratBadan}
                                    onChange={e => handleBeratChange(e.target.value)}
                                    className={`w-full pr-14 pl-3 py-2 text-sm rounded-xl border bg-neutral-50 text-neutral-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500 transition-colors ${errors.beratBadan ? 'border-red-400 bg-red-50' : 'border-neutral-200'
                                        }`}
                                />
                                <span className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-neutral-400">kg</span>
                                {errors.beratBadan && (
                                    <svg className="absolute right-8 top-1/2 -translate-y-1/2 w-4 h-4 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                        <circle cx="12" cy="12" r="10" /><line x1="12" y1="8" x2="12" y2="12" /><line x1="12" y1="16" x2="12.01" y2="16" />
                                    </svg>
                                )}
                            </div>
                            {errors.beratBadan && <p className="text-xs text-red-500 mt-1">{errors.beratBadan}</p>}
                        </div>
                        <div>
                            <label className="block text-xs text-neutral-500 mb-1.5">Tinggi Badan (TB)</label>
                            <div className="relative">
                                <input
                                    type="text"
                                    value={tinggiBadan}
                                    onChange={e => handleTinggiChange(e.target.value)}
                                    className={`w-full pr-12 pl-3 py-2 text-sm rounded-xl border bg-neutral-50 text-neutral-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500 transition-colors ${errors.tinggiBadan ? 'border-red-400 bg-red-50' : 'border-neutral-200'
                                        }`}
                                />
                                <span className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-neutral-400">cm</span>
                            </div>
                            {errors.tinggiBadan && <p className="text-xs text-red-500 mt-1">{errors.tinggiBadan}</p>}
                        </div>
                    </div>

                    {/* Lingkar Kepala */}
                    <div>
                        <label className="block text-xs text-neutral-500 mb-1.5">Lingkar Kepala</label>
                        <div className="relative">
                            <input
                                type="text"
                                value={lingkarKepala}
                                onChange={e => setLingkarKepala(e.target.value)}
                                className="w-full pr-12 pl-3 py-2 text-sm rounded-xl border border-neutral-200 bg-neutral-50 text-neutral-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500"
                            />
                            <span className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-neutral-400">cm</span>
                        </div>
                    </div>

                    {/* Info box */}
                    <div className="flex gap-2.5 bg-blue-50 border border-blue-100 rounded-xl px-3 py-3">
                        <svg className="w-4 h-4 text-blue-500 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <p className="text-xs text-blue-600 leading-relaxed">
                            Data ini akan memperbarui nilai Z-Score secara otomatis. Perubahan akan tercatat di Laporan Bulanan Puskesmas Terpadu.
                        </p>
                    </div>

                    {/* Last updated */}
                    <div className="flex items-center gap-2.5 pt-1 border-t border-neutral-100">
                        <div className="w-7 h-7 rounded-full bg-emerald-100 flex items-center justify-center text-[10px] font-semibold text-emerald-800 flex-shrink-0">
                            {data.terakhirDiperbarui.inisial}
                        </div>
                        <p className="text-xs text-neutral-400">
                            Terakhir diperbarui oleh:{' '}
                            <span className="font-medium text-neutral-600">{data.terakhirDiperbarui.nama}</span>{' '}
                            pada {data.terakhirDiperbarui.tanggal}.
                        </p>
                    </div>
                </div>

                {/* Footer */}
                <div className="flex items-center justify-end gap-3 px-5 py-4 border-t border-neutral-100">
                    <button onClick={onClose} className="text-sm text-neutral-500 hover:text-neutral-700 px-3 py-2 rounded-xl hover:bg-neutral-100 transition-colors">
                        Batal
                    </button>
                    <button
                        onClick={handleSubmit}
                        className="flex items-center gap-2 bg-emerald-700 hover:bg-emerald-800 text-white text-sm font-medium px-4 py-2 rounded-xl transition-colors"
                    >
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
                        </svg>
                        Simpan Perubahan
                    </button>
                </div>
            </div>
        </div>
    );
}