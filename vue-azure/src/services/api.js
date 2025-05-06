import axios from 'axios'

// Create axios instance with base URL from environment variables
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Add request interceptor to add auth token to requests
apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// Add response interceptor to handle common errors
apiClient.interceptors.response.use(
  response => {
    return response
  },
  error => {
    // Handle 401 Unauthorized errors (token expired, etc.)
    if (error.response && error.response.status === 401) {
      // Clear auth data
      localStorage.removeItem('access_token')
      // Redirect to login page
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default {
  // Auth endpoints
  auth: {
    getMicrosoftLoginUrl() {
      return apiClient.get('/auth/microsoft')
    },
    signOut() {
      return apiClient.post('/auth/signout')
    }
  },

  // Tasks endpoints
  tasks: {
    getAll() {
      return apiClient.get('/tasks')
    },
    getById(id) {
      return apiClient.get(`/tasks/${id}`)
    },
    create(task) {
      return apiClient.post('/tasks', task)
    },
    update(id, task) {
      return apiClient.put(`/tasks/${id}`, task)
    },
    delete(id) {
      return apiClient.delete(`/tasks/${id}`)
    }
  },

  // Posts endpoints
  posts: {
    getAll() {
      return apiClient.get('/posts')
    },
    getById(id) {
      return apiClient.get(`/posts/${id}`)
    },
    create(post) {
      return apiClient.post('/posts', post)
    },
    update(id, post) {
      return apiClient.put(`/posts/${id}`, post)
    },
    delete(id) {
      return apiClient.delete(`/posts/${id}`)
    }
  }

}
