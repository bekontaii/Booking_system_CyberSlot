import { Link, useNavigate, useParams } from 'react-router-dom';
import { clubs } from '../data/clubs';

export default function ClubDetailPage() {
  const { clubId } = useParams();
  const navigate = useNavigate();
  const club = clubs.find((item) => item.id === clubId);

  if (!club) {
    return (
      <div className="container page club-detail page-fade">
        <h2>Club not found</h2>
        <p>The club you are looking for does not exist.</p>
        <button className="btn" type="button" onClick={() => navigate('/clubs')}>
          Back to clubs
        </button>
      </div>
    );
  }

  return (
    <div className="container page club-detail page-fade">
      <div className="club-hero">
        <div className="club-hero-text">
          <h2>{club.name}</h2>
          <p className="club-hero-city">{club.city}</p>
          <p className="club-hero-address">{club.address}</p>
          <Link className="btn" to="/clubs">
            Back to clubs
          </Link>
        </div>
        <div className="club-hero-media">
          <img src={club.image} alt={`${club.name} club`} />
        </div>
      </div>

      <section className="zones">
        <h3>Zones and prices</h3>
        <div className="zones-grid">
          {club.zones.map((zone) => (
            <article key={zone.name} className="zone-card">
              <div className="zone-card-header">
                <h4>{zone.name}</h4>
              </div>
              <ul className="zone-specs">
                {zone.specs.map((spec) => (
                  <li key={spec}>{spec}</li>
                ))}
              </ul>
              <div className="zone-prices">
                {zone.prices.map((price) => (
                  <p key={price}>{price}</p>
                ))}
              </div>
              <Link className="btn primary" to="/booking">
                Book now
              </Link>
            </article>
          ))}
        </div>
      </section>
    </div>
  );
}
