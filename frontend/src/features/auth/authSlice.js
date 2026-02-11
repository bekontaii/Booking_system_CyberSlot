import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { api } from '../../services/api';

const tokenFromStorage = localStorage.getItem('jwt_token');

export const login = createAsyncThunk('auth/login', async (payload, { rejectWithValue }) => {
  const result = await api.login(payload);
  if (!result.ok) {
    return rejectWithValue(result.data.error || 'Login failed');
  }
  return result.data.token;
});

export const register = createAsyncThunk('auth/register', async (payload, { rejectWithValue }) => {
  const result = await api.register(payload);
  if (!result.ok) {
    return rejectWithValue(result.data.error || 'Registration failed');
  }
  return true;
});

export const logout = createAsyncThunk('auth/logout', async () => {
  await api.logout();
  return true;
});

const authSlice = createSlice({
  name: 'auth',
  initialState: {
    token: tokenFromStorage || '',
    status: 'idle',
    error: ''
  },
  reducers: {
    clearAuth(state) {
      state.token = '';
      state.error = '';
      state.status = 'idle';
      localStorage.removeItem('jwt_token');
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(login.pending, (state) => {
        state.status = 'loading';
        state.error = '';
      })
      .addCase(login.fulfilled, (state, action) => {
        state.status = 'succeeded';
        state.token = action.payload || '';
        localStorage.setItem('jwt_token', state.token);
      })
      .addCase(login.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.payload || 'Login failed';
      })
      .addCase(register.pending, (state) => {
        state.status = 'loading';
        state.error = '';
      })
      .addCase(register.fulfilled, (state) => {
        state.status = 'succeeded';
      })
      .addCase(register.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.payload || 'Registration failed';
      })
      .addCase(logout.fulfilled, (state) => {
        state.token = '';
        state.status = 'idle';
        state.error = '';
        localStorage.removeItem('jwt_token');
      });
  }
});

export const { clearAuth } = authSlice.actions;
export default authSlice.reducer;
