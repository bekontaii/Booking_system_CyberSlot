import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { api } from '../../services/api';

export const createBooking = createAsyncThunk('booking/createBooking', async (payload, { rejectWithValue }) => {
  const result = await api.createBooking(payload);
  if (!result.ok) {
    return rejectWithValue(result.data.error || 'Booking failed');
  }
  return result.data;
});

const bookingSlice = createSlice({
  name: 'booking',
  initialState: {
    status: 'idle',
    error: '',
    lastBooking: null
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(createBooking.pending, (state) => {
        state.status = 'loading';
        state.error = '';
      })
      .addCase(createBooking.fulfilled, (state, action) => {
        state.status = 'succeeded';
        state.lastBooking = action.payload;
      })
      .addCase(createBooking.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.payload || 'Booking failed';
      });
  }
});

export default bookingSlice.reducer;
