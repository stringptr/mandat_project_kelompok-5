import { useEffect, useRef } from 'react';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

interface MapItem {
  nama: string;
  prevalensi: number;
  jumlahKasus: number;
  level: string;
}

interface PrevalensiMapProps {
  data: MapItem[];
}

const CENTROIDS: Record<string, [number, number]> = {
  'Kab. Cilacap': [-7.73, 109.01],
  'Kab. Banyumas': [-7.48, 109.29],
  'Kab. Purbalingga': [-7.39, 109.36],
  'Kab. Banjarnegara': [-7.40, 109.70],
  'Kab. Kebumen': [-7.67, 109.66],
  'Kab. Purworejo': [-7.71, 110.01],
  'Kab. Wonosobo': [-7.36, 109.90],
  'Kab. Magelang': [-7.48, 110.22],
  'Kota Magelang': [-7.48, 110.22],
  'Kab. Boyolali': [-7.52, 110.60],
  'Kab. Klaten': [-7.68, 110.63],
  'Kab. Sukoharjo': [-7.68, 110.83],
  'Kota Surakarta': [-7.56, 110.83],
  'Kab. Wonogiri': [-7.82, 110.93],
  'Kab. Karanganyar': [-7.62, 111.04],
  'Kab. Sragen': [-7.43, 111.02],
  'Kab. Grobogan': [-7.02, 110.93],
  'Kab. Blora': [-6.97, 111.42],
  'Kab. Rembang': [-6.71, 111.34],
  'Kab. Pati': [-6.75, 111.04],
  'Kab. Kudus': [-6.80, 110.84],
  'Kab. Jepara': [-6.58, 110.67],
  'Kab. Demak': [-6.89, 110.64],
  'Kota Semarang': [-6.99, 110.42],
  'Kab. Semarang': [-7.20, 110.44],
  'Kota Salatiga': [-7.33, 110.50],
  'Kab. Temanggung': [-7.32, 110.17],
  'Kab. Kendal': [-6.92, 110.20],
  'Kab. Batang': [-6.91, 109.73],
  'Kab. Pekalongan': [-7.05, 109.53],
  'Kota Pekalongan': [-6.88, 109.67],
  'Kab. Pemalang': [-6.90, 109.38],
  'Kab. Tegal': [-6.87, 109.14],
  'Kota Tegal': [-6.87, 109.14],
  'Kab. Brebes': [-6.87, 108.89],
};

function getColor(level: string | undefined): string {
  if (!level) return '#9ca3af';
  if (level === 'tinggi') return '#dc2626';
  if (level === 'sedang') return '#f97316';
  return '#22c55e';
}

function matchDistrict(centroidName: string, data: MapItem[]): MapItem | undefined {
  const clean = (s: string) =>
    s.toLowerCase().replace(/kab\.?\s*/g, '').replace(/kota\s*/g, '').replace(/\s+/g, ' ').trim();
  const cn = clean(centroidName);
  return data.find(
    (d) => clean(d.nama).includes(cn) || cn.includes(clean(d.nama))
  );
}

export function PrevalensiMap({ data }: PrevalensiMapProps): JSX.Element {
  const mapContainer = useRef<HTMLDivElement>(null);
  const mapInstance = useRef<L.Map | null>(null);
  const circlesLayer = useRef<L.LayerGroup | null>(null);

  useEffect(() => {
    if (!mapContainer.current || mapInstance.current) return;

    const map = L.map(mapContainer.current, {
      center: [-7.30, 110.0],
      zoom: 9,
      zoomControl: true,
      scrollWheelZoom: true,
    });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
      maxZoom: 18,
    }).addTo(map);

    const group = L.layerGroup().addTo(map);

    Object.entries(CENTROIDS).forEach(([name, [lat, lng]]) => {
      const match = matchDistrict(name, data);

      const circle = L.circleMarker([lat, lng], {
        radius: match ? 10 + (match.prevalensi / 10) : 6,
        fillColor: getColor(match?.level),
        color: '#ffffff',
        weight: 1.5,
        fillOpacity: 0.85,
      });

      let tooltipContent = `<strong>${name.replace('Kab. ', '').replace('Kota ', '')}</strong>`;
      if (match) {
        tooltipContent += `<br/>Prevalensi: ${match.prevalensi.toFixed(1)}%`;
        tooltipContent += `<br/>Kasus: ${match.jumlahKasus}`;
      } else {
        tooltipContent += '<br/>Tidak ada data';
      }

      circle.bindTooltip(tooltipContent, {
        sticky: true,
        direction: 'top',
        offset: [0, -8],
      });

      circle.addTo(group);
    });

    circlesLayer.current = group;
    mapInstance.current = map;

    return () => {
      map.remove();
      mapInstance.current = null;
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (!circlesLayer.current) return;

    circlesLayer.current.eachLayer((l) => {
      const circle = l as L.CircleMarker;
      const latlng = circle.getLatLng();
      const name = Object.entries(CENTROIDS).find(
        ([, [lat, lng]]) => lat === latlng.lat && lng === latlng.lng
      )?.[0];
      if (!name) return;

      const match = matchDistrict(name, data);
      circle.setStyle({ fillColor: getColor(match?.level) });
      circle.setRadius(match ? 10 + (match.prevalensi / 10) : 6);

      let content = `<strong>${name.replace('Kab. ', '').replace('Kota ', '')}</strong>`;
      if (match) {
        content += `<br/>Prevalensi: ${match.prevalensi.toFixed(1)}%`;
        content += `<br/>Kasus: ${match.jumlahKasus}`;
      } else {
        content += '<br/>Tidak ada data';
      }
      circle.setTooltipContent(content);
    });
  }, [data]);

  const dataCount = data.length;
  const rataRata =
    dataCount > 0
      ? (data.reduce((s, d) => s + d.prevalensi, 0) / dataCount).toFixed(1)
      : '0';

  return (
    <div className="relative">
      <div className="flex items-center gap-4 mb-4 flex-wrap">
        {[
          { label: 'Tinggi (>20%)', color: '#dc2626' },
          { label: 'Sedang (10-20%)', color: '#f97316' },
          { label: 'Rendah (<10%)', color: '#22c55e' },
          { label: 'Tidak Ada Data', color: '#9ca3af' },
        ].map((l) => (
          <div key={l.label} className="flex items-center gap-1.5">
            <span className="w-3 h-3 rounded-full flex-shrink-0" style={{ background: l.color }} />
            <span className="text-[10px] text-neutral-600 font-body">{l.label}</span>
          </div>
        ))}
      </div>

      <div className="relative bg-neutral-50 rounded-2xl border border-neutral-100 overflow-hidden">
        <div ref={mapContainer} className="w-full" style={{ height: 500 }} />

        <div className="absolute bottom-4 right-4 z-[1000] bg-white rounded-xl shadow-lg border border-neutral-100 p-3 min-w-44">
          <p className="text-xs font-bold text-neutral-800 font-body mb-1">Ringkasan</p>
          <p className="text-xs text-neutral-500 font-body">
            Wilayah: <span className="font-semibold text-neutral-700">{dataCount}/35</span>
          </p>
          <p className="text-xs text-neutral-500 font-body">
            Rata-rata: <span className="font-semibold text-neutral-700">{rataRata}%</span>
          </p>
        </div>
      </div>

      <p className="text-[10px] text-neutral-400 mt-2 text-center">
        Sumber peta: OpenStreetMap. Data: SiGizi.
      </p>
    </div>
  );
}
