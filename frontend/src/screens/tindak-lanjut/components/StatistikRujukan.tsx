import { Stethoscope, Clock, CheckCircle2 } from 'lucide-react';

interface StatItem {
    label: string;
    value: string;
    icon: React.ReactNode;
    bgIcon: string;
    trend: string;
    trendColor: string;
}

export default function StatistikRujukan(): JSX.Element {
    const stats: StatItem[] = [
        {
            label: 'Total Pasien Dirujuk',
            value: '24',
            icon: <Stethoscope size={22} className="text-emerald-600" />,
            bgIcon: 'bg-emerald-100',
            trend: '+3 minggu ini',
            trendColor: 'text-emerald-600',
        },
        {
            label: 'Sedang Diproses',
            value: '8',
            icon: <Clock size={22} className="text-sky-600" />,
            bgIcon: 'bg-sky-100',
            trend: '2 selesai hari ini',
            trendColor: 'text-sky-600',
        },
        {
            label: 'Sudah Diterima Faskes',
            value: '16',
            icon: <CheckCircle2 size={22} className="text-emerald-600" />,
            bgIcon: 'bg-emerald-100',
            trend: '85% keberhasilan',
            trendColor: 'text-emerald-600',
        },
    ];

    return (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-5 mb-8">
            {stats.map((stat, idx) => (
                <div key={idx} className="bg-white rounded-2xl p-5 border border-slate-100 shadow-sm hover:shadow-md transition-shadow">
                    <div className="flex items-start justify-between mb-3">
                        <div className={`w-11 h-11 rounded-xl ${stat.bgIcon} flex items-center justify-center`}>
                            {stat.icon}
                        </div>
                        <span className={`text-xs font-medium ${stat.trendColor} bg-slate-50 px-2 py-1 rounded-lg`}>
                            {stat.trend}
                        </span>
                    </div>
                    <h3 className="text-2xl font-bold text-slate-800 mb-1">{stat.value}</h3>
                    <p className="text-sm text-slate-500">{stat.label}</p>
                </div>
            ))}
        </div>
    );
}