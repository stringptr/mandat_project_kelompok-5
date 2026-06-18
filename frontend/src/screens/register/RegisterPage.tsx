import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { CheckCircle, ShieldCheck, Headset, ArrowLeft, ArrowRight, User, Mail, Lock, Phone, Calendar, MapPin, BookOpen, Briefcase, Users, CreditCard, FileText, Building2 } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';

// ── Types ──────────────────────────────────────────────────────────────────
type RegisterRole = 'Ibu/Wali' | 'Bidan' | 'Kader Posyandu';

type Step = 'role' | 'berkas' | 'berkas2' | 'berkas3' | 'done-pending' | 'done';

interface RoleOption {
  id: RegisterRole;
  icon: string;
  label: string;
  desc: string;
}

const ROLE_OPTIONS: RoleOption[] = [
  { id: 'Ibu/Wali', icon: '🤰', label: 'Ibu/Wali', desc: 'Pantau gizi ibu hamil' },
  { id: 'Bidan', icon: '👩‍⚕️', label: 'Bidan', desc: 'Pengelola data wilayah' },
  { id: 'Kader Posyandu', icon: '🏥', label: 'Kader', desc: 'Pengelola data wilayah' },
];

const KABUPATEN = ['Semarang', 'Sukoharjo', 'Boyolali', 'Klaten', 'Sragen', 'Brebes', 'Pekalongan', 'Wonosobo'];
const KECAMATAN = ['Tembalang', 'Banyumanik', 'Candisari', 'Gajahmungkur', 'Genuk', 'Gunungpati'];
const KELURAHAN = ['Tembalang', 'Sendangmulyo', 'Bulusan', 'Jangli', 'Mangunharjo'];
const PENDIDIKAN = ['SD', 'SMP', 'SMA/SMK', 'Diploma 3', 'Sarjana 1', 'Sarjana 2/3'];
const PEKERJAAN = ['Ibu Rumah Tangga', 'Pegawai Swasta', 'PNS', 'Wiraswasta', 'Dosen', 'Guru', 'Petani'];
const PENDAPATAN = ['< Rp1.000.000', 'Rp1.000.000 – Rp2.500.000', 'Rp2.500.000 – Rp5.000.000', '> Rp5.000.000'];

// ── Left panel ─────────────────────────────────────────────────────────────
function LeftPanel(): JSX.Element {
  return (
    <div className="hidden lg:flex w-[420px] flex-shrink-0 bg-gradient-to-b from-primary-800 to-primary flex-col justify-between p-10 relative overflow-hidden"
      style={{ background: 'linear-gradient(160deg, #0a7a52 0%, #095c3e 60%, #064e3b 100%)' }}>
      {/* Decorative shapes */}
      <div className="absolute top-0 right-0 w-64 h-64 rounded-full bg-white/5 -translate-y-1/3 translate-x-1/3" />
      <div className="absolute bottom-0 left-0 w-48 h-48 rounded-full bg-white/5 translate-y-1/3 -translate-x-1/3" />
      <div className="absolute top-1/2 left-1/2 w-96 h-96 rounded-full bg-white/[0.03] -translate-x-1/2 -translate-y-1/2" />

      {/* Logo */}
      <div className="relative z-10">
        <div className="flex items-center gap-2 mb-2">
          <img src="/logo-sigizi.svg" alt="SiGizi" className="w-8 h-8 object-contain brightness-0 invert" />
          <span className="text-white font-bold text-lg font-headline">SiGizi</span>
        </div>
        <p className="text-white/60 text-xs font-body">Nurturing the future through community health monitoring.</p>
      </div>

      {/* Bottom copy */}
      <div className="relative z-10">
        <h2 className="text-white text-3xl font-bold font-headline leading-tight mb-3">
          Bersama Pantau<br />Tumbuh Kembang<br />Buah Hati
        </h2>
        <p className="text-white/60 text-sm font-body leading-relaxed">
          Platform cerdas untuk monitoring gizi ibu dan anak, menghubungkan keluarga dengan tenaga kesehatan secara real-time
        </p>
      </div>
    </div>
  );
}

// ── Form wrapper ───────────────────────────────────────────────────────────
function FormWrap({ children, onLogin }: { children: React.ReactNode; onLogin: () => void }): JSX.Element {
  return (
    <div className="flex flex-col min-h-full p-10">
      <div className="flex-1 flex flex-col max-w-md w-full mx-auto">
        {children}
      </div>
      <div className="max-w-md w-full mx-auto mt-6 text-center space-y-3">
        <p className="text-xs text-neutral-500 font-body">
          Sudah memiliki akun?{' '}
          <button onClick={onLogin} className="text-primary font-semibold hover:text-primary-600 transition-colors">Masuk di sini</button>
        </p>
        <div className="flex items-center justify-center gap-4 text-[10px] text-neutral-400 uppercase tracking-widest font-body">
          <span className="flex items-center gap-1"><ShieldCheck size={11} /> Keamanan Terenkripsi</span>
          <span className="text-neutral-200">·</span>
          <span className="flex items-center gap-1"><Headset size={11} /> Dukungan Teknis 24/7</span>
        </div>
      </div>
    </div>
  );
}

// ── Input helpers ──────────────────────────────────────────────────────────
function Field({ label, icon, children }: { label: string; icon: React.ReactNode; children: React.ReactNode }): JSX.Element {
  return (
    <div>
      <label className="flex items-center gap-1.5 text-xs font-semibold text-neutral-600 font-body mb-1.5">
        <span className="text-neutral-400">{icon}</span>{label}
      </label>
      {children}
    </div>
  );
}

const inputCls = "w-full px-4 py-2.5 bg-white border border-neutral-200 rounded-xl text-sm font-body text-neutral-800 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-200 focus:border-primary transition-colors";
const selectCls = `${inputCls} appearance-none`;

function FileInput({ label, icon }: { label: string; icon: React.ReactNode }): JSX.Element {
  const [name, setName] = useState('');
  return (
    <Field label={label} icon={icon}>
      <div className="flex items-center gap-2 w-full px-4 py-2.5 bg-white border border-neutral-200 rounded-xl">
        <span className="flex-1 text-sm text-neutral-400 font-body truncate">{name || 'Unggah Foto'}</span>
        <label className="text-xs font-semibold text-primary cursor-pointer hover:text-primary-600 transition-colors">
          Pilih File
          <input type="file" accept="image/*" className="hidden" onChange={(e) => setName(e.target.files?.[0]?.name ?? '')} />
        </label>
      </div>
    </Field>
  );
}

function VerifNotice(): JSX.Element {
  return (
    <div className="flex items-start gap-2 bg-blue-50 border border-blue-200 rounded-xl px-4 py-3 text-xs text-blue-700 font-body leading-relaxed">
      <span className="text-blue-500 flex-shrink-0 mt-0.5">ℹ</span>
      <span>Proses Verifikasi: Data NIK akan diverifikasi oleh sistem kependudukan. Akun Kader memerlukan persetujuan admin wilayah dalam 1×24 jam setelah registrasi.</span>
    </div>
  );
}

// ── Main component ─────────────────────────────────────────────────────────
export default function RegisterPage(): JSX.Element {
  const navigate = useNavigate();
  const { register } = useAuth();
  const [step, setStep] = useState<Step>('role');
  const [selectedRole, setSelectedRole] = useState<RegisterRole>('Ibu/Wali');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Form state
  const [form1, setForm1] = useState({ email: '', password: '', nama: '', nik: '', telepon: '', jenisKelamin: 'Laki-Laki', noSk: '', idPosyandu: '' });
  const [form2, setForm2] = useState({ tanggalLahir: '', provinsi: 'Jawa Tengah', kabupaten: 'Sukoharjo', kecamatan: 'Tembalang', kelurahan: 'Bulusan' });
  const [form3, setForm3] = useState({ pendidikan: 'Sarjana 1', pekerjaan: 'Dosen', pendapatan: 'Rp2.500.000 – Rp5.000.000', tanggungan: '1' });

  const needsMultiStep = selectedRole === 'Ibu/Wali';
  const isPendingRole = selectedRole === 'Bidan' || selectedRole === 'Kader Posyandu';

  const goLogin = () => navigate('/login');

  const handleRegister = async () => {
    setLoading(true);
    setError('');
    try {
      const roleMap: Record<string, string> = {
        'Ibu/Wali': 'Ibu Hamil',
        'Bidan': 'Bidan',
        'Kader Posyandu': 'Kader',
      };
      await register({
        email: form1.email,
        password: form1.password,
        no_hp: form1.telepon,
        nama: form1.nama,
        nik: form1.nik,
        jenis_kelamin: form1.jenisKelamin,
        tanggal_lahir: new Date(`${form2.tanggalLahir || '2000'}-01-01T00:00:00Z`).toISOString(),
        id_lokasi: 1,
        id_pendidikan: form3.pendidikan ? 1 : null,
        id_pekerjaan: form3.pekerjaan ? 1 : null,
        id_pendapatan: form3.pendapatan ? 1 : null,
        jumlah_tanggungan: form3.tanggungan ? parseInt(form3.tanggungan) : null,
        role: roleMap[selectedRole] || 'Ibu Hamil',
        no_sk: form1.noSk || undefined,
        id_posyandu: form1.idPosyandu ? parseInt(form1.idPosyandu, 10) : null,
      });
      if (isPendingRole) {
        setStep('done-pending');
      } else {
        setStep('done');
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Registrasi gagal. Silakan coba lagi.');
    } finally {
      setLoading(false);
    }
  };

  // ── Step: Pilih Role ───────────────────────────────────────────────────
  if (step === 'role') {
    return (
      <div className="min-h-screen flex font-body">
        <LeftPanel />
        <div className="flex-1 bg-white overflow-y-auto">
          <FormWrap onLogin={goLogin}>
            <h1 className="text-2xl font-bold text-neutral-900 font-headline mb-1">Buat Akun Baru</h1>
            <p className="text-sm text-neutral-500 mb-6">Lengkapi data diri Anda untuk memulai pemantauan gizi.</p>

            <div className="mb-6">
              <label className="flex items-center gap-1.5 text-xs font-semibold text-neutral-500 mb-3">
                <Users size={13} /> Daftar Sebagai
              </label>
              <div className="space-y-2">
                {ROLE_OPTIONS.map((r) => (
                  <button key={r.id} type="button" onClick={() => setSelectedRole(r.id)}
                    className={`w-full flex items-center gap-4 p-4 rounded-xl border-2 text-left transition-all ${selectedRole === r.id ? 'border-primary bg-primary-50' : 'border-neutral-200 hover:border-neutral-300 bg-white'}`}>
                    <span className={`w-9 h-9 rounded-full flex items-center justify-center text-base flex-shrink-0 ${selectedRole === r.id ? 'bg-primary-100' : 'bg-neutral-100'}`}>
                      {r.icon}
                    </span>
                    <div>
                      <p className={`text-sm font-bold font-headline ${selectedRole === r.id ? 'text-primary' : 'text-neutral-800'}`}>{r.label}</p>
                      <p className="text-xs text-neutral-500">{r.desc}</p>
                    </div>
                    {selectedRole === r.id && <CheckCircle size={18} className="text-primary ml-auto flex-shrink-0" />}
                  </button>
                ))}
              </div>
            </div>

            <button onClick={() => setStep('berkas')}
              className="w-full py-3 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-bold transition-colors flex items-center justify-center gap-2">
              Lanjut Isi Berkas <ArrowRight size={16} />
            </button>
          </FormWrap>
        </div>
      </div>
    );
  }

  // ── Step: Isi Berkas Hal. 1 ────────────────────────────────────────────
  if (step === 'berkas') {
    return (
      <div className="min-h-screen flex font-body">
        <LeftPanel />
        <div className="flex-1 bg-white overflow-y-auto">
          <FormWrap onLogin={goLogin}>
            <h1 className="text-2xl font-bold text-neutral-900 font-headline mb-1">Daftar {selectedRole}</h1>
            <p className="text-sm text-neutral-500 mb-6">Lengkapi berkas untuk melanjutkan pendaftaran.</p>

            <div className="space-y-4">
              <FileInput label="Kartu Tanda Penduduk (KTP)" icon={<span>🪪</span>} />
              {isPendingRole && <FileInput label="Surat Keterangan Aktif Kerja" icon={<span>📄</span>} />}

              <Field label="NIK (16 digit)" icon={<CreditCard size={13} />}>
                <input type="text" value={form1.nik} onChange={(e) => setForm1({ ...form1, nik: e.target.value })}
                  placeholder="1234567890123456" maxLength={16} className={inputCls} />
              </Field>

              <Field label="Email Aktif" icon={<Mail size={13} />}>
                <input type="email" value={form1.email} onChange={(e) => setForm1({ ...form1, email: e.target.value })}
                  placeholder="dani@gmail.com" className={inputCls} />
              </Field>

              <Field label="Kata Sandi" icon={<Lock size={13} />}>
                <input type="password" value={form1.password} onChange={(e) => setForm1({ ...form1, password: e.target.value })}
                  placeholder="dani1234" className={inputCls} />
              </Field>

              {isPendingRole && selectedRole === 'Kader Posyandu' && (
                <>
                  <Field label="No. Surat Keputusan (SK)" icon={<FileText size={13} />}>
                    <input type="text" value={form1.noSk} onChange={(e) => setForm1({ ...form1, noSk: e.target.value })}
                      placeholder="SK.01/Posyandu/2025" className={inputCls} />
                  </Field>
                  <Field label="ID Posyandu" icon={<Building2 size={13} />}>
                    <input type="number" value={form1.idPosyandu} onChange={(e) => setForm1({ ...form1, idPosyandu: e.target.value })}
                      placeholder="1" className={inputCls} />
                  </Field>
                </>
              )}

              {!isPendingRole && (
                <>
                  <Field label="Nama Lengkap" icon={<User size={13} />}>
                    <input type="text" value={form1.nama} onChange={(e) => setForm1({ ...form1, nama: e.target.value })}
                      placeholder="dani" className={inputCls} />
                  </Field>
                  <Field label="Nomor Telepon" icon={<Phone size={13} />}>
                    <input type="tel" value={form1.telepon} onChange={(e) => setForm1({ ...form1, telepon: e.target.value })}
                      placeholder="62123456789" className={inputCls} />
                  </Field>
                  <Field label="Jenis Kelamin" icon={<User size={13} />}>
                    <select value={form1.jenisKelamin} onChange={(e) => setForm1({ ...form1, jenisKelamin: e.target.value })} className={selectCls}>
                      <option>Laki-Laki</option><option>Perempuan</option>
                    </select>
                  </Field>
                </>
              )}

              <VerifNotice />

              {error && <p className="text-xs text-red-500 bg-red-50 border border-red-100 rounded-lg px-3 py-2">{error}</p>}
            </div>

            <div className="flex gap-3 mt-6">
              <button onClick={() => setStep('role')}
                className="flex-1 py-3 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-bold transition-colors flex items-center justify-center gap-2">
                <ArrowLeft size={15} /> Pilih Role
              </button>
              <button onClick={() => isPendingRole ? handleRegister() : setStep('berkas2')}
                disabled={loading}
                className="flex-1 py-3 bg-primary hover:bg-primary-600 disabled:bg-primary/60 text-white rounded-xl text-sm font-bold transition-colors flex items-center justify-center gap-2">
                {loading ? <><span className="w-4 h-4 border-2 border-white/40 border-t-white rounded-full animate-spin" />Memproses...</> : isPendingRole ? 'Selesaikan Registrasi Awal' : 'Halaman 2'} <ArrowRight size={15} />
              </button>
            </div>
          </FormWrap>
        </div>
      </div>
    );
  }

  // ── Step: Isi Berkas Hal. 2 (Ibu/Wali only) ───────────────────────────
  if (step === 'berkas2') {
    return (
      <div className="min-h-screen flex font-body">
        <LeftPanel />
        <div className="flex-1 bg-white overflow-y-auto">
          <FormWrap onLogin={goLogin}>
            <h1 className="text-2xl font-bold text-neutral-900 font-headline mb-1">Daftar {selectedRole}</h1>
            <p className="text-sm text-neutral-500 mb-6">Lengkapi berkas untuk melanjutkan pendaftaran.</p>

            <div className="space-y-4">
              <Field label="Tanggal Lahir" icon={<Calendar size={13} />}>
                <select value={form2.tanggalLahir} onChange={(e) => setForm2({ ...form2, tanggalLahir: e.target.value })} className={selectCls}>
                  <option value="">01 Januari 2027</option>
                  {Array.from({ length: 50 }, (_, i) => 1970 + i).map((y) => (
                    <option key={y} value={y.toString()}>{y}</option>
                  ))}
                </select>
              </Field>
              <Field label="Provinsi Tempat Tinggal" icon={<MapPin size={13} />}>
                <select value={form2.provinsi} onChange={(e) => setForm2({ ...form2, provinsi: e.target.value })} className={selectCls}>
                  <option>Jawa Tengah</option><option>DKI Jakarta</option><option>Jawa Barat</option><option>Jawa Timur</option>
                </select>
              </Field>
              <Field label="Kabupaten/Kota Tempat Tinggal" icon={<MapPin size={13} />}>
                <select value={form2.kabupaten} onChange={(e) => setForm2({ ...form2, kabupaten: e.target.value })} className={selectCls}>
                  {KABUPATEN.map((k) => <option key={k}>{k}</option>)}
                </select>
              </Field>
              <Field label="Kecamatan Tempat Tinggal" icon={<MapPin size={13} />}>
                <select value={form2.kecamatan} onChange={(e) => setForm2({ ...form2, kecamatan: e.target.value })} className={selectCls}>
                  {KECAMATAN.map((k) => <option key={k}>{k}</option>)}
                </select>
              </Field>
              <Field label="Kelurahan Tempat Tinggal" icon={<MapPin size={13} />}>
                <select value={form2.kelurahan} onChange={(e) => setForm2({ ...form2, kelurahan: e.target.value })} className={selectCls}>
                  {KELURAHAN.map((k) => <option key={k}>{k}</option>)}
                </select>
              </Field>
              <VerifNotice />
            </div>

            <div className="flex gap-3 mt-6">
              <button onClick={() => setStep('berkas')}
                className="flex-1 py-3 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-bold transition-colors flex items-center justify-center gap-2">
                <ArrowLeft size={15} /> Halaman 1
              </button>
              <button onClick={() => needsMultiStep ? setStep('berkas3') : setStep('done')}
                className="flex-1 py-3 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-bold transition-colors flex items-center justify-center gap-2">
                Halaman 3 <ArrowRight size={15} />
              </button>
            </div>
          </FormWrap>
        </div>
      </div>
    );
  }

  // ── Step: Isi Berkas Hal. 3 (Ibu/Wali only) ───────────────────────────
  if (step === 'berkas3') {
    return (
      <div className="min-h-screen flex font-body">
        <LeftPanel />
        <div className="flex-1 bg-white overflow-y-auto">
          <FormWrap onLogin={goLogin}>
            <h1 className="text-2xl font-bold text-neutral-900 font-headline mb-1">Daftar {selectedRole}</h1>
            <p className="text-sm text-neutral-500 mb-6">Lengkapi berkas untuk melanjutkan pendaftaran.</p>

            <div className="space-y-4">
              <Field label="Pendidikan Terakhir" icon={<BookOpen size={13} />}>
                <select value={form3.pendidikan} onChange={(e) => setForm3({ ...form3, pendidikan: e.target.value })} className={selectCls}>
                  {PENDIDIKAN.map((p) => <option key={p}>{p}</option>)}
                </select>
              </Field>
              <Field label="Pekerjaan" icon={<Briefcase size={13} />}>
                <select value={form3.pekerjaan} onChange={(e) => setForm3({ ...form3, pekerjaan: e.target.value })} className={selectCls}>
                  {PEKERJAAN.map((p) => <option key={p}>{p}</option>)}
                </select>
              </Field>
              <Field label="Kategori Pendapatan" icon={<Briefcase size={13} />}>
                <select value={form3.pendapatan} onChange={(e) => setForm3({ ...form3, pendapatan: e.target.value })} className={selectCls}>
                  {PENDAPATAN.map((p) => <option key={p}>{p}</option>)}
                </select>
              </Field>
              <Field label="Jumlah Tanggungan" icon={<Users size={13} />}>
                <input type="number" min="0" max="20" value={form3.tanggungan} onChange={(e) => setForm3({ ...form3, tanggungan: e.target.value })}
                  placeholder="1" className={inputCls} />
              </Field>
              <VerifNotice />

              {error && <p className="text-xs text-red-500 bg-red-50 border border-red-100 rounded-lg px-3 py-2">{error}</p>}
            </div>

            <div className="flex gap-3 mt-6">
              <button onClick={() => setStep('berkas2')}
                className="flex-1 py-3 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-bold transition-colors flex items-center justify-center gap-2">
                <ArrowLeft size={15} /> Halaman 2
              </button>
              <button onClick={handleRegister}
                disabled={loading}
                className="flex-1 py-3 bg-primary hover:bg-primary-600 disabled:bg-primary/60 text-white rounded-xl text-sm font-bold transition-colors flex items-center justify-center gap-2">
                {loading ? <><span className="w-4 h-4 border-2 border-white/40 border-t-white rounded-full animate-spin" />Memproses...</> : 'Selesaikan Registrasi'} <ArrowRight size={15} />
              </button>
            </div>
          </FormWrap>
        </div>
      </div>
    );
  }

  // ── Step: Done Pending (Bidan/Kader) ──────────────────────────────────
  if (step === 'done-pending') {
    return (
      <div className="min-h-screen flex font-body">
        <LeftPanel />
        <div className="flex-1 bg-white overflow-y-auto">
          <FormWrap onLogin={goLogin}>
            <div className="text-center py-8">
              <h1 className="text-2xl font-bold text-neutral-900 font-headline mb-2">Registrasi Awal Selesai</h1>
              <p className="text-sm text-neutral-500 mb-8">Berkas sedang diverifikasi<br />Tunggu pesan untuk langkah selanjutnya</p>
              <div className="w-24 h-24 bg-primary-50 rounded-2xl flex items-center justify-center mx-auto mb-8 border-2 border-primary-100">
                <svg className="w-12 h-12 text-neutral-800" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
                  <circle cx="12" cy="12" r="10" /><line x1="10" y1="15" x2="10" y2="9" /><line x1="14" y1="15" x2="14" y2="9" />
                </svg>
              </div>
              <button onClick={() => navigate('/')}
                className="w-full py-3 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-bold transition-colors flex items-center justify-center gap-2">
                Kembali ke Halaman Awal <ArrowRight size={15} />
              </button>
            </div>
          </FormWrap>
        </div>
      </div>
    );
  }

  // ── Step: Done (Ibu/Wali) ──────────────────────────────────────────────
  return (
    <div className="min-h-screen flex font-body">
      <LeftPanel />
      <div className="flex-1 bg-white overflow-y-auto">
        <FormWrap onLogin={goLogin}>
          <div className="text-center py-8">
            <h1 className="text-2xl font-bold text-neutral-900 font-headline mb-2">Registrasi Selesai</h1>
            <p className="text-sm text-neutral-500 mb-8">Masuk untuk memulai pemantauan gizi.</p>
            <div className="w-24 h-24 bg-primary-50 rounded-2xl flex items-center justify-center mx-auto mb-8 border-2 border-primary-100">
              <svg className="w-12 h-12 text-neutral-800" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
                <circle cx="12" cy="12" r="10" /><polyline points="9 12 11 14 15 10" />
              </svg>
            </div>
            <button onClick={goLogin}
              className="w-full py-3 bg-primary hover:bg-primary-600 text-white rounded-xl text-sm font-bold transition-colors flex items-center justify-center gap-2">
              Masuk <ArrowRight size={15} />
            </button>
          </div>
        </FormWrap>
      </div>
    </div>
  );
}
