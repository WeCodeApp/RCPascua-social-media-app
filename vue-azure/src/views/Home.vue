<template>
  <div class="home">
    <div class="hero">
      <h1>Welcome to Task Manager</h1>
      <p>A simple and efficient way to manage your tasks</p>
      <div class="cta-buttons">
        <router-link v-if="!authStore.isAuthenticated" to="/login" class="btn btn-primary">Login</router-link>
        <router-link v-else to="/tasks" class="btn btn-primary">My Tasks</router-link>
      </div>
    </div>
    
    <div class="features">
      <div class="feature">
        <div class="feature-icon" style="background-color: var(--primary)">
          <i class="fas fa-list"></i>
        </div>
        <h3>Organize Tasks</h3>
        <p>Keep track of all your tasks in one place</p>
      </div>
      
      <div class="feature">
        <div class="feature-icon" style="background-color: var(--secondary)">
          <i class="fas fa-check"></i>
        </div>
        <h3>Track Progress</h3>
        <p>Monitor your progress and stay on top of deadlines</p>
      </div>
      
      <div class="feature">
        <div class="feature-icon" style="background-color: var(--accent)">
          <i class="fas fa-sync"></i>
        </div>
        <h3>Stay Synchronized</h3>
        <p>Access your tasks from anywhere, anytime</p>
      </div>
    </div>
  </div>
</template>

<script>
import { useAuthStore } from '../stores/auth'

export default {
  name: 'HomeView',
  setup() {
    const authStore = useAuthStore()
    
    return {
      authStore
    }
  }
}
</script>

<style scoped lang="scss">
.home {
  text-align: center;
}

.hero {
  padding: 4rem 1rem;
  background: linear-gradient(135deg, var(--primary) 0%, rgb(0.5526315789, 29.2894736842, 93.9473684211) 100%);
  color: white;
  border-radius: 8px;
  margin-bottom: 3rem;
  
  h1 {
    font-size: 2.5rem;
    margin-bottom: 1rem;
  }
  
  p {
    font-size: 1.2rem;
    margin-bottom: 2rem;
    opacity: 0.9;
  }
}

.cta-buttons {
  display: flex;
  justify-content: center;
  gap: 1rem;
  
  .btn {
    padding: 0.75rem 1.5rem;
    font-size: 1.1rem;
  }
}

.features {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 2rem;
  margin: 2rem 0;
}

.feature {
  flex: 1;
  min-width: 250px;
  max-width: 350px;
  padding: 2rem;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s;
  
  &:hover {
    transform: translateY(-5px);
  }
  
  h3 {
    margin: 1rem 0;
    color: var(--primary);
  }
  
  p {
    color: var(--gray);
  }
}

.feature-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  color: white;
  font-size: 1.5rem;
}
</style>