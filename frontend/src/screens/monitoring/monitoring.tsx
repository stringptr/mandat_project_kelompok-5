import { useAuth } from '../../context/AuthContext';
import type { Role } from '../../App';
import { IbuWaliSection } from './sections/IbuWaliSection';
import { BidanSection } from './sections/BidanSection';
import { DinkesSection } from './sections/DinkesSection';
import { KaderSection } from './sections/KaderSection';

interface MonitoringProps {
    currentRole: Role;
}

const ROLE_DESC: Record<Role, string> = {
    'Ibu/Wali': 'Pantau tumbuh kembang anak Anda dengan mudah',
    'Bidan': 'Ringkasan wilayah binaan dan operasional harian Anda',
    'Dinas Kesehatan': 'Overview regional dan capaian program gizi Jawa Tengah',
    'Kader Posyandu': 'Status operasional Posyandu Melati dan data balita',
};

const SECTION_MAP: Record<Role, React.FC> = {
    'Ibu/Wali': IbuWaliSection,
    'Bidan': BidanSection,
    'Dinas Kesehatan': DinkesSection,
    'Kader Posyandu': KaderSection,
};

export default function Monitoring({ currentRole }: MonitoringProps): JSX.Element {
    const { user } = useAuth();
    const today = new Date().toLocaleDateString('id-ID', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric',
    });

    const description = ROLE_DESC[currentRole];
    const SectionComponent = SECTION_MAP[currentRole];

    if (!description || !SectionComponent) {
        return (
            <div className="p-8">
                <div className="bg-red-50 border border-red-200 rounded-xl p-4">
                    <p className="text-red-600 font-bold">Error</p>
                    <p className="text-red-500">Role tidak dikenal: {JSON.stringify(currentRole)}</p>
                </div>
            </div>
        );
    }

    return (
        <div className="space-y-6 font-body text-neutral-800">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2">
                <div>
                    <h2 className="text-xl font-bold text-neutral-800 font-headline">
                        Selamat Datang, {user?.name || 'Pengguna'}
                    </h2>
                    <p className="text-sm text-neutral-500 mt-0.5">{description}</p>
                </div>
                <span className="text-xs text-neutral-400 bg-white border border-neutral-100 px-4 py-2 rounded-xl font-medium">
                    {today}
                </span>
            </div>

            <SectionComponent />
        </div>
    );
}
