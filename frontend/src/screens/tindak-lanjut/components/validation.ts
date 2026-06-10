// validation.ts
export interface ValidationResult {
  valid: boolean;
  code?: string;
  message?: string;
}

export const validatePatient = (selected: any): ValidationResult => {
  if (!selected) {
    return { valid: false, code: 'ERR-VAL-05', message: 'Masukkan nama yang valid' };
  }
  return { valid: true };
};

export const validateJenisTindakan = (value: string): ValidationResult => {
  const allowed = ['Rujukan', 'Tindak Lanjut'];
  if (!allowed.includes(value)) {
    return { valid: false, code: 'ERR-VAL-07', message: 'Jenis tindakan tidak valid' };
  }
  return { valid: true };
};

export const validateLokasi = (value: string, jenis: string): ValidationResult => {
  if (jenis === 'Rujukan' && (!value || !value.trim())) {
    return { valid: false, code: 'ERR-VAL-03', message: 'Lokasi rujukan wajib diisi' };
  }
  return { valid: true };
};

export const validateCatatan = (text: string): ValidationResult => {
  const trimmed = text.trim();
  if (!trimmed) {
    return { valid: false, code: 'ERR-VAL-03', message: 'Catatan Medis wajib diisi' };
  }
  if (trimmed.length < 10) {
    return { valid: false, code: 'ERR-VAL-03', message: 'Catatan Medis minimal 10 karakter' };
  }
  if (trimmed.length > 1000) {
    return { valid: false, code: 'ERR-VAL-03', message: 'Catatan Medis maksimal 1000 karakter' };
  }
  return { valid: true };
};

export const validateTanggal = (dateStr: string): ValidationResult => {
  if (!dateStr) {
    return { valid: false, code: 'ERR-VAL-06', message: 'Tanggal target wajib diisi' };
  }
  const date = new Date(dateStr);
  const now = new Date();
  if (isNaN(date.getTime())) {
    return { valid: false, code: 'ERR-VAL-06', message: 'Tanggal tidak valid' };
  }
  if (date <= now) {
    return { valid: false, code: 'ERR-VAL-06', message: 'Tanggal harus di masa depan' };
  }
  return { valid: true };
};
