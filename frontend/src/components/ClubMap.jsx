import { useEffect, useRef, useState } from 'react';

const MAP_KEY = '11990778-0991-4611-b4a8-09689f7a5c35';
const MAP_CENTER = [71.4491, 51.1694];
const MAP_ZOOM = 12;

const CLUBS = [
  {
    name: 'Top Game',
    address: 'Astana, Dinmukhamed Kunayev Street 23',
    coordinates: [71.4304, 51.1289]
  },
  {
    name: 'BRO',
    address: 'Astana, Kuishi Dina Street 31',
    coordinates: [71.4728, 51.1526]
  },
  {
    name: 'Xan.exe',
    address: 'Astana, Heydar Aliyev Street 3',
    coordinates: [71.4183, 51.0907]
  },
  {
    name: 'Prime Game Hub',
    address: 'Astana, Kerei and Zhanibek Khandar Street 14/2',
    coordinates: [71.4049, 51.1324]
  },
  {
    name: 'Yamato Cyber Club',
    address: 'Astana, Syganak Street 21/1',
    coordinates: [71.4092, 51.1347]
  }
];

let mapglLoader;

const loadMapgl = () => {
  if (mapglLoader) {
    return mapglLoader;
  }

  mapglLoader = new Promise((resolve, reject) => {
    if (window.mapgl && window.mapgl.Map) {
      resolve(window.mapgl);
      return;
    }

    const scriptId = 'mapgl-script';
    const existing = document.getElementById(scriptId);
    if (existing) {
      existing.addEventListener('load', () => resolve(window.mapgl));
      existing.addEventListener('error', () => reject(new Error('MapGL failed to load')));
      return;
    }

    const script = document.createElement('script');
    script.id = scriptId;
    script.src = `https://mapgl.2gis.com/api/js/v1?key=${MAP_KEY}`;
    script.async = true;
    script.onload = () => resolve(window.mapgl);
    script.onerror = () => reject(new Error('MapGL failed to load'));
    document.head.appendChild(script);
  });

  return mapglLoader;
};

export default function ClubMap() {
  const containerRef = useRef(null);
  const mapRef = useRef(null);
  const markersRef = useRef([]);
  const [error, setError] = useState(false);

  useEffect(() => {
    let mounted = true;

    loadMapgl()
      .then((mapgl) => {
        if (!mounted || !containerRef.current) {
          return;
        }

        if (!mapRef.current) {
          mapRef.current = new mapgl.Map(containerRef.current, {
            center: MAP_CENTER,
            zoom: MAP_ZOOM,
            key: MAP_KEY
          });
        }

        markersRef.current.forEach((marker) => marker.destroy && marker.destroy());
        markersRef.current = [];

        CLUBS.forEach((club) => {
          const marker = new mapgl.Marker(mapRef.current, {
            coordinates: club.coordinates
          });
          marker.on('click', () => {
            const popup = new mapgl.Popup(mapRef.current, {
              coordinates: club.coordinates,
              html: `<strong>${club.name}</strong><br/>${club.address}`
            });
            popup.open();
          });
          markersRef.current.push(marker);
        });
      })
      .catch(() => {
        if (mounted) {
          setError(true);
        }
      });

    return () => {
      mounted = false;
      markersRef.current.forEach((marker) => marker.destroy && marker.destroy());
      markersRef.current = [];
      if (mapRef.current) {
        mapRef.current.destroy();
        mapRef.current = null;
      }
    };
  }, []);

  return (
    <div className="map-container" ref={containerRef}>
      {error && (
        <div className="clubs-map-error">
          Map failed to load. Please check your 2GIS API key and network access.
        </div>
      )}
    </div>
  );
}
