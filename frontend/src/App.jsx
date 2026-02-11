import { Navigate, Route, Routes } from 'react-router-dom';
import Navbar from './components/Navbar';
import Footer from './components/Footer';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ClubsPage from './pages/ClubsPage';
import BookingPage from './pages/BookingPage';
import ClubDetailPage from './pages/ClubDetailPage';

export default function App() {
  return (
    <div className="app">
      <Navbar />
      <main className="main">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/clubs" element={<ClubsPage />} />
          <Route path="/clubs/:clubId" element={<ClubDetailPage />} />
          <Route path="/booking" element={<BookingPage />} />
          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </main>
      <Footer />
    </div>
  );
}
