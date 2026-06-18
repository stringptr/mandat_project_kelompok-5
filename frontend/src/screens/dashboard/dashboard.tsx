/**
 * Dashboard — orchestrator
 *
 * Routes to the correct section based on currentRole.
 *
 * Role → Section
 *   Ibu/Wali        → IbuWaliSection   (status anak, riwayat, jadwal imunisasi)
 *   Kader Posyandu  → BidanKaderSection (ringkasan monitoring posyandu)
 *   Bidan           → BidanKaderSection (ringkasan monitoring + bidan stats)
 *   Dinas Kesehatan → DinkesSection    (statistik regional, peta, tren)
 */
import type { Role } from '../../App';
import { IbuWaliSection } from './sections/IbuWaliSection';
import { BidanKaderSection } from './sections/BidanKaderSection';
import { DinkesSection } from './sections/DinkesSection';

interface DashboardProps {
  currentRole: Role;
}

export default function Dashboard({ currentRole }: DashboardProps): JSX.Element {
  return (
    <div className="space-y-5 font-body text-neutral-800">
      {/* Role section */}
      {currentRole === 'Ibu/Wali' && <IbuWaliSection />}
      {(currentRole === 'Bidan' || currentRole === 'Kader Posyandu') && (
        <BidanKaderSection currentRole={currentRole} />
      )}
      {currentRole === 'Dinas Kesehatan' && <DinkesSection />}
    </div>
  );
}
