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
  const normalizedDescription = String(description || '')
    .replace(/\\n\\n/g, ' ')
    .replace(/\n\n/g, ' ')
    .replace(/\s+/g, ' ')
    .trim();

  return (
    <section className={`hero-section ${reverse ? 'reverse' : ''}`}>
      <div className="hero-content">
        <h1>{title}</h1>
        <p>{normalizedDescription}</p>
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
