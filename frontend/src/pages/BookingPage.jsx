import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { createBooking } from '../features/booking/bookingSlice';
import { fetchClubs } from '../features/clubs/clubsSlice';
import { api } from '../services/api';

export default function BookingPage() {
  const dispatch = useDispatch();
  const { items: clubs } = useSelector((state) => state.clubs);
  const bookingState = useSelector((state) => state.booking);
  const [pcs, setPcs] = useState([]);
  const [form, setForm] = useState({
    clubId: '',
    pcId: '',
    userId: '',
    startTime: '',
    endTime: ''
  });
  const [message, setMessage] = useState('');

  useEffect(() => {
    if (!clubs.length) {
      dispatch(fetchClubs());
    }
  }, [dispatch, clubs.length]);

  useEffect(() => {
    const loadPcs = async () => {
      if (!form.clubId) {
        setPcs([]);
        return;
      }
      const result = await api.getPCsByClub(form.clubId);
      if (result.ok && Array.isArray(result.data)) {
        setPcs(result.data);
      } else {
        setPcs([]);
      }
    };

    loadPcs();
  }, [form.clubId]);

  const handleChange = (e) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage('');

    const payload = {
      pc_id: Number(form.pcId),
      user_id: Number(form.userId),
      start_time: new Date(form.startTime).toISOString(),
      end_time: new Date(form.endTime).toISOString()
    };

    const result = await dispatch(createBooking(payload));
    if (createBooking.fulfilled.match(result)) {
      setMessage('Booking created successfully.');
      setForm({ clubId: '', pcId: '', userId: '', startTime: '', endTime: '' });
    }
  };

  return (
    <div className="container page">
      <h2>Create Booking</h2>

      <form className="form" onSubmit={handleSubmit}>
        <label>
          Club
          <select name="clubId" value={form.clubId} onChange={handleChange} required>
            <option value="">Select club</option>
            {clubs.map((club) => (
              <option key={club.id} value={club.id}>
                {club.name}
              </option>
            ))}
          </select>
        </label>

        <label>
          PC
          <select name="pcId" value={form.pcId} onChange={handleChange} required>
            <option value="">Select PC</option>
            {pcs.map((pc) => (
              <option key={pc.id} value={pc.id}>
                {pc.name || `PC #${pc.id}`}
              </option>
            ))}
          </select>
        </label>

        <label>
          User ID
          <input name="userId" value={form.userId} onChange={handleChange} required />
        </label>

        <label>
          Start Time
          <input type="datetime-local" name="startTime" value={form.startTime} onChange={handleChange} required />
        </label>

        <label>
          End Time
          <input type="datetime-local" name="endTime" value={form.endTime} onChange={handleChange} required />
        </label>

        <button className="btn" type="submit" disabled={bookingState.status === 'loading'}>
          {bookingState.status === 'loading' ? 'Creating...' : 'Create Booking'}
        </button>
      </form>

      {message && <p className="message success">{message}</p>}
      {bookingState.error && <p className="message error">{bookingState.error}</p>}
    </div>
  );
}
