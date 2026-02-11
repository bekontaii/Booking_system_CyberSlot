import { useDispatch, useSelector } from 'react-redux';
import { login } from '../features/auth/authSlice';
import { useState } from 'react';

export default function LoginPage() {
  const dispatch = useDispatch();
  const { status, error } = useSelector((state) => state.auth);
  const [form, setForm] = useState({ username: '', password: '' });
  const [message, setMessage] = useState('');

  const handleChange = (e) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage('');
    const result = await dispatch(login(form));
    if (login.fulfilled.match(result)) {
      setMessage('Login successful. Token saved.');
    }
  };

  return (
    <div className="container form-page">
      <h2>Login</h2>
      <form className="form" onSubmit={handleSubmit}>
        <label>
          Username
          <input name="username" value={form.username} onChange={handleChange} required />
        </label>
        <label>
          Password
          <input type="password" name="password" value={form.password} onChange={handleChange} required />
        </label>
        <button className="btn" type="submit" disabled={status === 'loading'}>
          {status === 'loading' ? 'Loading...' : 'Login'}
        </button>
      </form>
      {message && <p className="message success">{message}</p>}
      {error && <p className="message error">{error}</p>}
    </div>
  );
}
