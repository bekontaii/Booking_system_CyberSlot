import { Link } from 'react-router-dom';
import { useEffect, useRef, useState } from 'react';
import HeroSection from '../components/HeroSection';
import CTASection from '../components/CTASection';

export default function HomePage() {
  const clubSectionRef = useRef(null);
  const [clubVisible, setClubVisible] = useState(false);

  useEffect(() => {
    const section = clubSectionRef.current;
    if (!section) {
      return undefined;
    }

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setClubVisible(true);
            observer.unobserve(entry.target);
          }
        });
      },
      { threshold: 0.2 }
    );

    observer.observe(section);
    return () => observer.disconnect();
  }, []);

  return (
    <div className="home">
      <HeroSection
        title="Public Computer Club Booking Platform in Kazakhstan"
        description="Make it easy for gamers to book a computer in their favorite gaming club anytime, anywhere. Our platform allows players to choose a city, select a computer club, pick a specific PC, view real-time availability and tariffs, and reserve gaming time online in just a few clicks. Share your booking link across the web, attract new players, reduce manual work for administrators, and increase club occupancy with a modern online booking experience."
        image="https://senet.cloud/wp-content/uploads/2023/03/outside-booking-applied-1.png"
        buttonText="Find a computer"
        buttonHref="/booking"
      />

      <HeroSection
        title="Online Computer Reservation"
        description="Reserve any computer or gaming zone in advance with full transparency. Our system shows which PCs are available, which are already booked, and what configurations they offer so every gamer gets exactly the setup they need. Clubs gain a centralized booking system, automated reservations, and better control over schedules, while players enjoy fast, reliable, and hassle-free booking."
        image="https://senet.cloud/wp-content/uploads/2023/03/frame4381.svg"
        reverse
      />

      <section
        ref={clubSectionRef}
        className={`club-choice ${clubVisible ? 'is-visible' : ''}`}
      >
        <div className="container">
          <h2>Choose your club</h2>
          <div className="club-choice-grid">
            <article className="club-choice-card" style={{ '--delay': '0ms' }}>
              <img
                src="https://avatars.mds.yandex.net/get-altay/10494556/2a0000018e8a1eb0862111c81245a581c238/L_height"
                alt="Standard gaming zone"
              />
              <h3>Standard</h3>
              <p>
                Comfortable gaming PCs for everyday gaming, study, and relaxation.
                Perfect for popular online games and casual sessions.
              </p>
              <Link className="club-choice-link" to="/booking">
                Choose
              </Link>
            </article>
            <article className="club-choice-card" style={{ '--delay': '120ms' }}>
              <img
                src="https://dynamic-media-cdn.tripadvisor.com/media/photo-o/2c/8b/31/36/open-space-a-place-where.jpg?w=900&h=500&s=1"
                alt="Pro gaming zone"
              />
              <h3>Pro</h3>
              <p>
                High-performance PCs for esports and demanding games.
                High FPS, professional peripherals, and maximum comfort.
              </p>
              <Link className="club-choice-link" to="/booking">
                Choose
              </Link>
            </article>
            <article className="club-choice-card" style={{ '--delay': '240ms' }}>
              <img
                src="https://playroom.lv/ru/wp-content/uploads/2023/10/385815420_122101846040063251_2293891381055673392_n-1-725x444.jpg"
                alt="VIP gaming zone"
              />
              <h3>VIP</h3>
              <p>
                Premium gaming zones with top-tier configurations and private atmosphere.
                Maximum performance, comfort, and exclusive gaming experience.
              </p>
              <Link className="club-choice-link" to="/booking">
                Choose
              </Link>
            </article>
          </div>
        </div>
      </section>

      <CTASection />
    </div>
  );
}
