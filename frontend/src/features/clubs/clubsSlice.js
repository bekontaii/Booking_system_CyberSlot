import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { api } from '../../services/api';

export const fetchClubs = createAsyncThunk('clubs/fetchClubs', async (_, { rejectWithValue }) => {
  const result = await api.getClubs();
  if (!result.ok) {
    return rejectWithValue(result.data.error || 'Failed to load clubs');
  }
  return Array.isArray(result.data) ? result.data : [];
});

const clubsSlice = createSlice({
  name: 'clubs',
  initialState: {
    items: [],
    status: 'idle',
    error: ''
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchClubs.pending, (state) => {
        state.status = 'loading';
        state.error = '';
      })
      .addCase(fetchClubs.fulfilled, (state, action) => {
        state.status = 'succeeded';
        state.items = action.payload;
      })
      .addCase(fetchClubs.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.payload || 'Failed to load clubs';
      });
  }
});

export default clubsSlice.reducer;
