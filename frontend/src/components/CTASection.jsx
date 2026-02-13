import { Link } from 'react-router-dom';

export default function CTASection() {
  return (
    <section className="cta">
      <div className="container cta-inner">
        <h2>Book a gaming PC in just a few minutes</h2>
        <p>
          Find available computers in gaming clubs across Kazakhstan, choose the right configuration, and reserve online
        </p>
        <Link className="btn primary" to="/booking">
          Find a computer
        </Link>
      </div>
    </section>
  );
}

