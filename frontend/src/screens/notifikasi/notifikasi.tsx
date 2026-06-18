import type { Role } from '../../App';
import KaderNotifikasi from './components/kader/KaderNotifikasi';
import BidanNotifikasi from './components/bidan/BidanNotifikasi';
import DinkesNotifikasi from './components/dinas/DinkesNotifikasi';
import IbuWaliNotifikasi from './components/ibu-wali/IbuWaliNotifikasi';

interface NotifikasiProps {
  role: Role;
}

export default function Notifikasi({ role }: NotifikasiProps): JSX.Element {
  return (
    <div className="w-full max-w-5xl mx-auto pb-12 font-body text-neutral-800">
      {role === 'Kader Posyandu' && <KaderNotifikasi />}
      {role === 'Bidan' && <BidanNotifikasi />}
      {role === 'Dinas Kesehatan' && <DinkesNotifikasi />}
      {role === 'Ibu/Wali' && <IbuWaliNotifikasi />}
    </div>
  );
}
