const API_BASE = import.meta.env.VITE_API_URL || import.meta.env.VITE_API_BASE || '';

function buildUrl(path) {
  if (API_BASE) {
    const normalized = API_BASE.endsWith('/') ? API_BASE.slice(0, -1) : API_BASE;
    return `${normalized}${path}`;
  }
  return path;
}

function getToken() {
  return localStorage.getItem('jwt_token') || '';
}

async function request(path, options = {}, requireAuth = false) {
  const headers = options.headers ? { ...options.headers } : {};
  headers['Content-Type'] = 'application/json';

  if (requireAuth) {
    const token = getToken();
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
  }

  const response = await fetch(buildUrl(path), { ...options, headers });
  const data = await response.json().catch(() => ({}));
  return { ok: response.ok, status: response.status, data };
}

export const api = {
  login: (payload) => request('/auth/login', { method: 'POST', body: JSON.stringify(payload) }),
  register: (payload) => request('/auth/register', { method: 'POST', body: JSON.stringify(payload) }),
  logout: () => request('/api/logout', { method: 'POST' }, true),
  getClubs: () => request('/api/clubs', { method: 'GET' }, true),
  getPCsByClub: (clubId) => request(`/api/clubs/${clubId}/pcs`, { method: 'GET' }, true),
  createBooking: (payload) => request('/api/bookings', { method: 'POST', body: JSON.stringify(payload) }, true)
};
