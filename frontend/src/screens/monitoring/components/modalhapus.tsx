interface ModalHapusPemeriksaanProps {
    isOpen: boolean;
    onClose: () => void;
    namaAnak: string;
    tanggalPemeriksaan: string;
    onConfirm: () => void;
    isLoading?: boolean;
}

export default function ModalHapusPemeriksaan({
    isOpen,
    onClose,
    namaAnak,
    tanggalPemeriksaan,
    onConfirm,
    isLoading = false,
}: ModalHapusPemeriksaanProps) {
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
            <div className="bg-white rounded-2xl w-full max-w-sm shadow-xl overflow-hidden">
                {/* Close button */}
                <div className="flex justify-end px-4 pt-4">
                    <button
                        onClick={onClose}
                        disabled={isLoading}
                        className="w-7 h-7 flex items-center justify-center rounded-lg text-neutral-400 hover:bg-neutral-100 hover:text-neutral-600 transition-colors disabled:opacity-40"
                    >
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                {/* Content */}
                <div className="px-6 pb-6 text-center">
                    <div className="w-16 h-16 bg-red-50 rounded-full flex items-center justify-center mx-auto mb-4">
                        <svg className="w-8 h-8 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.8}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                        </svg>
                    </div>

                    <h3 className="text-base font-semibold text-neutral-800 mb-2">
                        Hapus Data Pemeriksaan?
                    </h3>

                    <p className="text-sm text-neutral-500 leading-relaxed">
                        Apakah Anda yakin ingin menghapus data pemeriksaan{' '}
                        <span className="font-medium text-neutral-700">{namaAnak}</span> tanggal{' '}
                        <span className="font-semibold text-neutral-800">{tanggalPemeriksaan}</span>?
                        Tindakan ini tidak dapat dibatalkan dan akan menghapus data tindak lanjut yang terkait secara permanen dari sistem SiGizi.
                    </p>
                </div>

                {/* Footer */}
                <div className="flex gap-2.5 px-5 pb-5">
                    <button
                        onClick={onClose}
                        disabled={isLoading}
                        className="flex-1 py-2.5 text-sm font-medium text-neutral-700 bg-neutral-100 hover:bg-neutral-200 rounded-xl transition-colors disabled:opacity-40"
                    >
                        Batal
                    </button>
                    <button
                        onClick={onConfirm}
                        disabled={isLoading}
                        className="flex-1 flex items-center justify-center gap-2 py-2.5 text-sm font-medium text-white bg-red-500 hover:bg-red-600 rounded-xl transition-colors disabled:opacity-60"
                    >
                        {isLoading ? (
                            <>
                                <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                                </svg>
                                Menghapus...
                            </>
                        ) : (
                            <>
                                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                </svg>
                                Ya, Hapus
                            </>
                        )}
                    </button>
                </div>
            </div>
        </div>
    );
}