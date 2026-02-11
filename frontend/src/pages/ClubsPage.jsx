import ClubCard from '../components/ClubCard';
import ClubMap from '../components/ClubMap';
import { clubs } from '../data/clubs';

export default function ClubsPage() {
  return (
    <div className="container page clubs-page page-fade">
      <div className="clubs-header">
        <h2>Clubs in Astana</h2>
        <p>
          Choose a gaming club, explore zones, and book a computer with the right configuration for your session.
        </p>
      </div>
      <div className="clubs-grid">
        {clubs.map((club) => (
          <ClubCard key={club.id} club={club} />
        ))}
      </div>
      <div className="clubs-map-section">
        <h3>Club locations on the map</h3>
        <ClubMap />
      </div>
    </div>
  );
}
