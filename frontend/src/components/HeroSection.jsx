import { Link } from 'react-router-dom';

const logoUrl = '/logo.jpg';

export default function HeroSection({
  title,
  description,
  image,
  reverse,
  buttonText,
  buttonHref,
}) {
  const content = description.split('\n\n').map((text) => text.trim()).filter(Boolean);

  return (
    <section className={`hero-section ${reverse ? 'reverse' : ''}`}>
      <div className="hero-content">
        <h1>{title}</h1>
        {content.map((paragraph, index) => (
          <p key={index}>{paragraph}</p>
        ))}
        {buttonText && buttonHref && (
          <Link className="btn primary" to={buttonHref}>
            {buttonText}
          </Link>
        )}
      </div>
      <div className="hero-visual">
        <div className="hero-accent" aria-hidden="true" />
        {reverse && (
          <div className="hero-brand">
            <img src={logoUrl} alt="CyberSlot" />
            <span>CyberSlot</span>
          </div>
        )}
        <img className="hero-image" src={image} alt={title} loading="lazy" />
      </div>
    </section>
  );
}
