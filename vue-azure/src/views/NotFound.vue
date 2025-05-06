<template>
  <div class="not-found">
    <div class="not-found-content">
      <div class="error-code">404</div>
      <h1>Page Not Found</h1>
      <p>The page you are looking for doesn't exist or has been moved.</p>
      <div class="actions">
        <router-link to="/" class="btn btn-primary">Go Home</router-link>
        <router-link v-if="authStore.isAuthenticated" to="/tasks" class="btn btn-secondary">My Tasks</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { useAuthStore } from '../stores/auth'

export default {
  name: 'NotFoundView',
  setup() {
    const authStore = useAuthStore()
    
    return {
      authStore
    }
  }
}
</script>

<style scoped lang="scss">
.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 70vh;
  text-align: center;
}

.not-found-content {
  max-width: 500px;
  padding: 2rem;
}

.error-code {
  font-size: 8rem;
  font-weight: bold;
  color: var(--primary);
  line-height: 1;
  margin-bottom: 1rem;
  text-shadow: 2px 2px 0 var(--secondary);
}

h1 {
  font-size: 2rem;
  color: var(--dark);
  margin-bottom: 1rem;
}

p {
  color: var(--gray);
  margin-bottom: 2rem;
  font-size: 1.1rem;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 1rem;
}
</style>