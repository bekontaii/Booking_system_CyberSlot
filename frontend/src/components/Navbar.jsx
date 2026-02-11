import { Link, NavLink, useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { clearAuth, logout } from '../features/auth/authSlice';

const logoUrl = '/logo.jpg';

export default function Navbar() {
  const token = useSelector((state) => state.auth.token);
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleLogout = async () => {
    if (token) {
      await dispatch(logout());
    }
    dispatch(clearAuth());
    navigate('/login');
  };

  return (
    <header className="navbar">
      <div className="container navbar-inner">
        <Link to="/" className="brand">
          <img src={logoUrl} alt="CyberSlot" className="brand-logo" />
          <span>CyberSlot</span>
        </Link>
        <nav className="nav-links">
          <NavLink to="/" end>
            Home
          </NavLink>
          <NavLink to="/clubs">Clubs</NavLink>
          <NavLink to="/booking">Booking</NavLink>
        </nav>
        <div className="nav-actions">
          {token ? (
            <button className="btn" type="button" onClick={handleLogout}>
              Logout
            </button>
          ) : (
            <>
              <NavLink to="/login">Login</NavLink>
              <NavLink to="/register">Register</NavLink>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
