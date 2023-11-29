import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

// Request interceptor
// Note: Authentication is handled via HTTP-only cookies (_yun_session)
// This interceptor is kept for future extensibility (e.g., API keys, CSRF tokens)
api.interceptors.request.use(
  (config) => {
    // Could add custom headers here if needed
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Only redirect to login if:
    // 1. It's a 401 error
    // 2. NOT from the login endpoint
    // 3. NOT from the profile endpoint (used to check auth on page load)
    // 4. NOT already on the login or setup page
    const isLoginEndpoint = error.config?.url === 'public/action/login'
    const isProfileEndpoint = error.config?.url === '/action/profile'
    const isOnPublicPage = window.location.pathname === '/login' ||
                           window.location.pathname === '/setup'

    if (error.response?.status === 401 &&
        !isLoginEndpoint &&
        !isProfileEndpoint &&
        !isOnPublicPage) {
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

//export default api

// API endpoints
export const healthAPI = {
  check: () => api.get('public/_health'),
  ready: () => api.get('public/_ready'),
}

export const authAPI = {
  login: (data) => api.post('public/action/login', data),
  logout: () => api.post('public/action/logout'),
  getProfile: () => api.get('/action/profile'),
}

export const setupAPI = {
  install: (data) => api.post('public/action/setup', data),
  checkInstalled: () => api.get('public/action/setup/status'),
}

// API endpoints for users C
export const userAPI = {
  getUsers: () => api.get('/users'),
  createUser: (data) => api.post('/users', data),
  updateUser: (id, data) => api.put(`/users/${id}`, data),
  deleteUser: (id) => api.delete(`/users/${id}`),
}
