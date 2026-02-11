import { configureStore } from '@reduxjs/toolkit';
import authReducer from '../features/auth/authSlice';
import clubsReducer from '../features/clubs/clubsSlice';
import bookingReducer from '../features/booking/bookingSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    clubs: clubsReducer,
    booking: bookingReducer
  }
});
