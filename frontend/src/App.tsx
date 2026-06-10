import { Routes, Route, useNavigate } from "react-router-dom";
import { AuthProvider, useAuth } from "./context/AuthContext";

import { Sidebar } from "./components/sidebar";
import {Header} from "./components/header";
import { Footer } from "./components/footer";
import { GuestDashboard } from "./screens/guest/GuestDashboard";

// Screens
import LoginPage from "./screens/login/LoginPage";
import RegisterPage from "./screens/register/RegisterPage";
import Dashboard from "./screens/dashboard/dashboard";
import Monitoring from "./screens/monitoring/monitoring";
import TindakLanjut from "./screens/tindak-lanjut/tindak-lanjut";
import Edukasi from "./screens/edukasi/edukasi";
import UserManagement from "./screens/user-management/user-management";
import Notifikasi from "./screens/notifikasi/notifikasi";
import JadwalImunisasi from "./screens/jadwal-imunisasi/jadwal-imunisasi";

export type Role =
  | "Ibu/Wali"
  | "Bidan"
  | "Dinas Kesehatan"
  | "Kader Posyandu";

function AppShell(): JSX.Element {
  const { isLoggedIn, user } = useAuth();
  const navigate = useNavigate();

  const currentRole: Role =
    isLoggedIn && user ? user.role : "Kader Posyandu";

  const goLogin = () => navigate("/login");

  return (
    <Routes>
      {/* Auth Pages */}
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />

      {/* Main Layout */}
      <Route
        path="/*"
        element={
          <div className="flex h-screen bg-neutral-50">
            <Sidebar
              currentRole={currentRole}
              onLoginClick={goLogin}
            />

            <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
              <Header
                currentRole={isLoggedIn ? currentRole : undefined}
                onLoginClick={goLogin}
              />

              <main className="flex-1 overflow-y-auto p-8">
                <Routes>
                  <Route
                    path="/"
                    element={
                      isLoggedIn
                        ? <Dashboard currentRole={currentRole} />
                        : <GuestDashboard onLoginClick={goLogin} />
                    }
                  />
                  <Route
                    path="/monitoring"
                    element={<Monitoring currentRole={currentRole} />}
                  />
                  <Route
                    path="/tindak-lanjut"
                    element={<TindakLanjut currentRole={currentRole} />}
                  />
                  <Route
                    path="/edukasi"
                    element={<Edukasi currentRole={currentRole} />}
                  />
                  <Route
                    path="/user-management"
                    element={<UserManagement />}
                  />
                  <Route
                    path="/notifikasi"
                    element={<Notifikasi role={currentRole} />}
                  />
                  <Route
                    path="/jadwal-imunisasi"
                    element={<JadwalImunisasi currentRole={currentRole} />}
                  />
                </Routes>
              </main>

              <Footer />
            </div>
          </div>
        }
      />
    </Routes>
  );
}

export default function App(): JSX.Element {
  return (
    <AuthProvider>
      <AppShell />
    </AuthProvider>
  );
}