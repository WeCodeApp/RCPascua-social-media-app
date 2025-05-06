<template>
  <div class="login-page">
    <div class="login-container">
      <h1>Login to Social Media App</h1>
      <p>Sign in with your Microsoft account to access your tasks</p>
      
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <p>Loading...</p>
      </div>
      
      <div v-else-if="error" class="error-message">
        <p>{{ error }}</p>
        <button @click="resetError" class="btn btn-secondary">Try Again</button>
      </div>
      
      <button v-else @click="login" class="btn btn-microsoft">
        <span class="icon">M</span>
        Sign in with Microsoft
      </button>
    </div>
  </div>
</template>

<script>
import {ref} from 'vue'
import {useRouter} from 'vue-router'
import {useAuthStore} from '../stores/auth'

export default {
  name: 'LoginView',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const loading = ref(false)
    const error = ref(null)

    // Check if we have a code parameter in the URL (OAuth callback)
    const handleCallback = async () => {
      console.log("handle callback")
      const urlParams = new URLSearchParams(window.location.search)
      const code = urlParams.get('token')
      const user = urlParams.get('user')
      console.log("code", code)
      console.log("user", user)

      if (code) {
        try {
          loading.value = true
          await authStore.handleLoginCallback(code, user)
          router.push({ name: 'tasks' })
        } catch (err) {
          error.value = 'Authentication failed. Please try again.'
          console.error('Login error:', err)
        } finally {
          loading.value = false
        }
      }
    }

    // Call handleCallback on component mount
    handleCallback()
    
    const login = async () => {
      try {
        loading.value = true
        error.value = null
        window.location.href = await authStore.getMicrosoftLoginUrl()
      } catch (err) {
        error.value = 'Failed to get login URL. Please try again.'
        console.error('Login URL error:', err)
      } finally {
        loading.value = false
      }
    }
    
    const resetError = () => {
      error.value = null
    }
    
    return {
      loading,
      error,
      login,
      resetError
    }
  }
}
</script>

<style scoped lang="scss">
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}

.login-container {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  padding: 2rem;
  width: 100%;
  max-width: 500px;
  text-align: center;
  
  h1 {
    color: var(--primary);
    margin-bottom: 1rem;
  }
  
  p {
    color: var(--gray);
    margin-bottom: 2rem;
  }
}

.btn-microsoft {
  background-color: #2f2f2f;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  padding: 0.75rem;
  font-size: 1rem;
  
  &:hover {
    background-color: #1f1f1f;
  }
  
  .icon {
    background-color: white;
    color: #2f2f2f;
    width: 24px;
    height: 24px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 10px;
    font-weight: bold;
  }
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem 0;
  
  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid rgba(0, 0, 0, 0.1);
    border-radius: 50%;
    border-top-color: var(--primary);
    animation: spin 1s ease-in-out infinite;
    margin-bottom: 1rem;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
}

.error-message {
  background-color: rgba(207, 8, 34, 0.1);
  border-left: 4px solid var(--accent);
  padding: 1rem;
  margin-bottom: 1.5rem;
  text-align: left;
  
  p {
    color: var(--accent);
    margin-bottom: 1rem;
  }
}
</style>