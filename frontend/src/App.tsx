import { useState } from "react";
import { Routes, Route } from "react-router-dom";
import { Sidebar } from "./components/sidebar";
import { Header } from "./components/header";
import { Footer } from "./components/footer";

// Screens
import Dashboard from "./screens/dashboard/dashboard";
import Monitoring from "./screens/monitoring/monitoring";
import TindakLanjut from "./screens/tindak-lanjut/tindak-lanjut";
import Edukasi from "./screens/edukasi/edukasi";
import UserManagement from "./screens/user-management/user-management";
import Notifikasi from "./screens/notifikasi/notifikasi";

export type Role = 'Ibu/Wali' | 'Bidan' | 'Dinas Kesehatan' | 'Kader Posyandu';

export default function App(): JSX.Element {
  const [currentRole, setCurrentRole] = useState<Role>('Kader Posyandu');

  return (
    <div className="flex h-screen bg-neutral-50">
      <Sidebar currentRole={currentRole} />

      <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
        <Header currentRole={currentRole} onChangeRole={setCurrentRole} />

        <main className="flex-1 overflow-y-auto p-8">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/monitoring" element={<Monitoring currentRole={currentRole} />} /> {/* ← FIX */}
            <Route path="/tindak-lanjut" element={<TindakLanjut />} />
            <Route path="/edukasi" element={<Edukasi />} />
            <Route path="/user-management" element={<UserManagement />} />
            <Route path="/notifikasi" element={<Notifikasi />} />
          </Routes>
        </main>

        <Footer />
      </div>
    </div>
  );
}