import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type {
  DashboardStats, DistribusiGiziItem, TrenStuntingItem, StuntingWilayahItem,
} from '../types/api';
import type { Artikel } from '../screens/edukasi/data/artikel.data';
import type { JadwalImunisasi } from '../screens/jadwal-imunisasi/data/imunisasi.data';
import type { Rujukan } from '../screens/tindak-lanjut/components/RujukanAktif';

interface DashboardSlice {
  dashboardStats: DashboardStats | null;
  dashboardLoading: boolean;
  distribusiGizi: DistribusiGiziItem[];
  trenStunting: TrenStuntingItem[];
  stuntingPerWilayah: StuntingWilayahItem[];
  kehadiranBulanan: TrenStuntingItem[];
  jadwalTerdekat: Record<string, unknown>[];
  aktivitas: Record<string, unknown>[];
  imunisasiPersen: number;
  setDashboardStats: (s: DashboardStats) => void;
  setDashboardLoading: (v: boolean) => void;
  setDistribusiGizi: (d: DistribusiGiziItem[]) => void;
  setTrenStunting: (t: TrenStuntingItem[]) => void;
  setStuntingPerWilayah: (w: StuntingWilayahItem[]) => void;
  setKehadiranBulanan: (t: TrenStuntingItem[]) => void;
  setJadwalTerdekat: (j: Record<string, unknown>[]) => void;
  setAktivitas: (a: Record<string, unknown>[]) => void;
  setImunisasiPersen: (n: number) => void;
}

interface ArtikelSlice {
  artikelList: Artikel[];
  artikelLoading: boolean;
  setArtikelList: (list: Artikel[]) => void;
  setArtikelLoading: (v: boolean) => void;
}

interface ImunisasiSlice {
  imunisasiList: JadwalImunisasi[];
  imunisasiLoading: boolean;
  setImunisasiList: (list: JadwalImunisasi[]) => void;
  setImunisasiLoading: (v: boolean) => void;
}

interface TindakLanjutSlice {
  rujukanList: Rujukan[];
  rujukanLoading: boolean;
  setRujukanList: (list: Rujukan[]) => void;
  setRujukanLoading: (v: boolean) => void;
}

interface GlobalSlice {
  appInitialized: boolean;
}

type AppStore = DashboardSlice & ArtikelSlice & ImunisasiSlice & TindakLanjutSlice & GlobalSlice;

type PersistedState = Omit<
  AppStore,
  'dashboardLoading' | 'artikelLoading' | 'imunisasiLoading' | 'rujukanLoading'
>;

export const useAppStore = create<AppStore>()(
  persist(
    (set) => ({
      // ── Dashboard ──
      dashboardStats: null,
      dashboardLoading: false,
      distribusiGizi: [],
      trenStunting: [],
      stuntingPerWilayah: [],
      kehadiranBulanan: [],
      jadwalTerdekat: [],
      aktivitas: [],
      imunisasiPersen: 0,
      setDashboardStats: (s) => set({ dashboardStats: s }),
      setDashboardLoading: (v) => set({ dashboardLoading: v }),
      setDistribusiGizi: (d) => set({ distribusiGizi: d }),
      setTrenStunting: (t) => set({ trenStunting: t }),
      setStuntingPerWilayah: (w) => set({ stuntingPerWilayah: w }),
      setKehadiranBulanan: (t) => set({ kehadiranBulanan: t }),
      setJadwalTerdekat: (j) => set({ jadwalTerdekat: j }),
      setAktivitas: (a) => set({ aktivitas: a }),
      setImunisasiPersen: (n) => set({ imunisasiPersen: n }),

      // ── Artikel ──
      artikelList: [],
      artikelLoading: false,
      setArtikelList: (list) => set({ artikelList: list }),
      setArtikelLoading: (v) => set({ artikelLoading: v }),

      // ── Imunisasi ──
      imunisasiList: [],
      imunisasiLoading: false,
      setImunisasiList: (list) => set({ imunisasiList: list }),
      setImunisasiLoading: (v) => set({ imunisasiLoading: v }),

      // ── Tindak Lanjut ──
      rujukanList: [],
      rujukanLoading: false,
      setRujukanList: (list) => set({ rujukanList: list }),
      setRujukanLoading: (v) => set({ rujukanLoading: v }),

      // ── Global Init ──
      appInitialized: false,
    }),
    {
      name: 'sigizi-store',
      partialize: (state): PersistedState => {
        const { dashboardLoading, artikelLoading, imunisasiLoading, rujukanLoading, ...persisted } = state;
        return persisted;
      },
    }
  )
);
