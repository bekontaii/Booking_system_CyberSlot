import { Link } from 'react-router-dom';

export default function ClubCard({ club }) {
  return (
    <Link to={`/clubs/${club.id}`} className="club-card">
      <div className="club-card-media">
        <img src={club.image} alt={`${club.name} club`} />
      </div>
      <div className="club-card-body">
        <h3>{club.name}</h3>
        <p className="club-card-city">{club.city}</p>
        <p className="club-card-address">{club.address}</p>
      </div>
    </Link>
  );
}
