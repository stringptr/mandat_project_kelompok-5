import { createContext, useContext, useState } from 'react';
import type { Role } from '../App';

export interface UserProfile {
  name: string;
  role: Role;
  avatarUrl: string;
}

interface AuthContextValue {
  user: UserProfile | null;
  login: (profile: UserProfile) => void;
  logout: () => void;
  isLoggedIn: boolean;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }): JSX.Element {
  const [user, setUser] = useState<UserProfile | null>(null);

  const login = (profile: UserProfile) => setUser(profile);
  const logout = () => setUser(null);

  return (
    <AuthContext.Provider value={{ user, login, logout, isLoggedIn: user !== null }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used inside AuthProvider');
  return ctx;
}
