import { createContext, useContext, useEffect, useState, useCallback } from 'react';
import { apiGet, apiPost, setOnUnauthorized } from '../lib/api';
import type { Role } from '../App';

export interface UserProfile {
  id_user: number;
  name: string;
  email: string;
  nik: string;
  role: Role;
  roles: string[];
}

interface AuthContextValue {
  user: UserProfile | null;
  login: (email: string, password: string) => Promise<void>;
  register: (data: RegisterFormData) => Promise<void>;
  logout: () => Promise<void>;
  isLoggedIn: boolean;
  isLoading: boolean;
}

export interface LoginRequest {
  email?: string;
  nik?: string;
  password: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  access_token_expires_in: number;
  refresh_token_expires_in: number;
}

export interface JwtClaim {
  id_user: number;
  roles: string[];
  email: string;
  nik: string;
  jti: string;
  exp: number;
  iat: number;
}

export interface RegisterFormData {
  email: string;
  password: string;
  no_hp: string;
  nama: string;
  nik: string;
  jenis_kelamin: string;
  tanggal_lahir: string;
  id_lokasi: number;
  id_pendidikan?: number | null;
  id_pekerjaan?: number | null;
  id_pendapatan?: number | null;
  jumlah_tanggungan?: number | null;
  role: string;
  no_str?: string;
  wilayah_kerja?: number | null;
  no_sk?: string;
  id_posyandu?: number | null;
}

function mapRolesToDisplayRole(roles: string[]): Role {
  if (roles.includes('IBU_HAMIL') || roles.includes('PASIEN')) return 'Ibu/Wali';
  if (roles.includes('BIDAN')) return 'Bidan';
  if (roles.includes('DINKES')) return 'Dinas Kesehatan';
  if (roles.includes('KADER')) return 'Kader Posyandu';
  return 'Kader Posyandu';
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }): JSX.Element {
  const [user, setUser] = useState<UserProfile | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const clearAuth = useCallback(() => {
    setUser(null);
  }, []);

  useEffect(() => {
    setOnUnauthorized(clearAuth);
  }, [clearAuth]);

  useEffect(() => {
    const checkSession = async () => {
      try {
        const claim = await apiGet<JwtClaim>('/auth/me');
        setUser({
          id_user: claim.id_user,
          name: claim.email,
          email: claim.email,
          nik: claim.nik,
          role: mapRolesToDisplayRole(claim.roles),
          roles: claim.roles,
        });
      } catch {
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    };
    checkSession();
  }, []);

  const login = async (email: string, password: string) => {
    const body: LoginRequest = { email, password };
    await apiPost<AuthResponse>('/auth/login', body);
    const claim = await apiGet<JwtClaim>('/auth/me');
    setUser({
      id_user: claim.id_user,
      name: claim.email,
      email: claim.email,
      nik: claim.nik,
      role: mapRolesToDisplayRole(claim.roles),
      roles: claim.roles,
    });
  };

  const register = async (data: RegisterFormData) => {
    await apiPost<null>('/auth/register', data);
  };

  const logout = async () => {
    try {
      await apiPost<null>('/auth/logout');
    } catch {
      // ignore errors on logout
    }
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, register, logout, isLoggedIn: user !== null, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used inside AuthProvider');
  return ctx;
}
