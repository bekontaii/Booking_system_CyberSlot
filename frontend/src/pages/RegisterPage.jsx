import { useDispatch, useSelector } from 'react-redux';
import { register } from '../features/auth/authSlice';
import { useState } from 'react';

export default function RegisterPage() {
  const dispatch = useDispatch();
  const { status, error } = useSelector((state) => state.auth);
  const [form, setForm] = useState({
    name: '',
    surname: '',
    email: '',
    username: '',
    password: ''
  });
  const [message, setMessage] = useState('');

  const handleChange = (e) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage('');
    const result = await dispatch(register(form));
    if (register.fulfilled.match(result)) {
      setMessage('Registration successful. You can login now.');
    }
  };

  return (
    <div className="container form-page">
      <h2>Register</h2>
      <form className="form" onSubmit={handleSubmit}>
        <label>
          Name
          <input name="name" value={form.name} onChange={handleChange} required />
        </label>
        <label>
          Surname
          <input name="surname" value={form.surname} onChange={handleChange} required />
        </label>
        <label>
          Email
          <input type="email" name="email" value={form.email} onChange={handleChange} required />
        </label>
        <label>
          Username
          <input name="username" value={form.username} onChange={handleChange} required />
        </label>
        <label>
          Password
          <input type="password" name="password" value={form.password} onChange={handleChange} required />
        </label>
        <button className="btn" type="submit" disabled={status === 'loading'}>
          {status === 'loading' ? 'Loading...' : 'Register'}
        </button>
      </form>
      {message && <p className="message success">{message}</p>}
      {error && <p className="message error">{error}</p>}
    </div>
  );
}
