import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    isAuthenticated: false,
    loading: false,
    error: null
  }),

  getters: {
    getUser: (state) => state.user,
    getIsAuthenticated: (state) => state.isAuthenticated,
    getLoading: (state) => state.loading,
    getError: (state) => state.error
  },

  actions: {
    async getMicrosoftLoginUrl() {
      try {
        this.loading = true
        this.error = null

        const response = await axios.get(`${import.meta.env.VITE_API_URL}/auth/microsoft`)
        return response.data.login_url
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to get login URL'
        throw error
      } finally {
        this.loading = false
      }
    },

    async getCurrentUser() {
        try {
            const response = await axios.get("/auth/me", {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem("token")}`, // Include the JWT token
                },
            });
    
            // console.log("User Info:", response.data);
            return response.data;
        } catch (error) {
            console.error("Failed to fetch user info:", error.response?.data || error.message);
            throw error;
        }
    },

      async handleLoginCallback(code, user) {
      try {
        this.loading = true
        this.error = null

        // This would typically be handled by the backend
        // For this example, we'll simulate a successful login
        // this.user = {
        //   id: 1,
        //   name: 'Test User',
        //   email: 'test@example.com'
        // }

        this.isAuthenticated = true

        // Get the access token from the backend
        const parsedCode = typeof code === 'string' ? JSON.parse(code) : code;
        // Get the user info from the backend
        const parsedUser = typeof user === 'string' ? JSON.parse(user) : user;
        console.log("access token", parsedCode.access_token);
        // Assign user data to the store
        this.user = {
          id: parsedUser.id,
          name: parsedUser.name,
          email: parsedUser.email
        }

        // Store auth token in localStorage
        localStorage.setItem('access_token', parsedCode.access_token)
        localStorage.setItem('user_id', parsedUser.id)
        localStorage.setItem('user_name', parsedUser.name)

        return this.user
      } catch (error) {
        this.error = error.response?.data?.message || 'Login failed'
        throw error
      } finally {
        this.loading = false
      }
    },

    async logout() {
      try {
        this.loading = true
        this.error = null

        await axios.post(`${import.meta.env.VITE_API_URL}/auth/signout`)

        // Clear user data
        this.user = null
        this.isAuthenticated = false

        // Remove auth token from localStorage
        localStorage.removeItem('access_token')
      } catch (error) {
        this.error = error.response?.data?.message || 'Logout failed'
        throw error
      } finally {
        this.loading = false
      }
    },

    async checkAuth() {
      // Check if user is already authenticated
      const token = localStorage.getItem('access_token')

      if (token) {
        try {
          // In a real app, you would validate the token with the server
          // For this example, we'll just set the user as authenticated
          // this.user = {
          //   id: 1,
          //   name: 'Test User',
          //   email: 'test@example.com'
          // }
          this.getCurrentUser().then(user => {
            // console.log("Current User:", user);
            this.user = user; // Set the user data in the store
          });
          this.isAuthenticated = true
        } catch (error) {
          // If token validation fails, clear user data
          this.user = null
          this.isAuthenticated = false
          localStorage.removeItem('access_token')
        }
      }
    }
  }
})
