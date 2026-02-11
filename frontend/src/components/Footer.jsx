import { Link } from 'react-router-dom';

export default function Footer() {
  return (
    <footer className="footer">
      <div className="container footer-grid">
        <div className="footer-column">
          <h4>CyberSlot</h4>
          <p>
            CyberSlot is an online platform that helps gamers book computers in gaming and esports clubs across Kazakhstan.
            Fast, convenient, and transparent.
          </p>
        </div>
        <div className="footer-column">
          <h4>Company</h4>
          <ul>
            <li>
              <Link to="/">About us</Link>
            </li>
            <li>
              <Link to="/">Contacts</Link>
            </li>
            <li>
              <Link to="/">For clubs</Link>
            </li>
          </ul>
        </div>
        <div className="footer-column">
          <h4>Legal</h4>
          <ul>
            <li>
              <Link to="/">Privacy Policy</Link>
            </li>
            <li>
              <Link to="/">User Agreement</Link>
            </li>
          </ul>
        </div>
      </div>
      <div className="footer-bottom">
        <div className="container footer-bottom-inner">
          <span>(c) 2026 CyberSlot.</span>
          <div className="footer-socials">
            <Link to="/" aria-label="Instagram">
              <img
                src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQWEbfjuyn9VmJ8GgZ59dV9x8cTzp3fMvD_lQ&s"
                alt="Instagram"
              />
            </Link>
            <Link to="/" aria-label="Facebook">
              <img
                src="https://wallpapers.com/images/thumbnail/facebook-icon-blackand-white-lqsy5va0pj1y4v8g.webp"
                alt="Facebook"
              />
            </Link>
            <Link to="/" aria-label="YouTube">
              <img
                src="https://img.freepik.com/premium-vector/free-vector-youtube-icon-logo-black-white_901408-456.jpg"
                alt="YouTube"
              />
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
