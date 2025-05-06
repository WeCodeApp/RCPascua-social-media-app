<template>
  <div id="app">
    <header class="app-header">
      <div class="container">
        <h1 class="app-title">Social Media App</h1>
        <nav v-if="authStore.isAuthenticated">
          <button class="btn btn-logout" @click="logout">Logout</button>
        </nav>
        <p class="userInfo" v-if="authStore.user?.name">
          Current User: {{ authStore.user?.name }}
        </p>
      </div>
    </header>
    <main class="container">
      <router-view />
    </main>
    <footer class="app-footer">
      <div class="container">
        <!-- <p>Task Manager is a simple and efficient way to manage your tasks.</p> -->
        <p>&copy; {{ new Date().getFullYear() }} Task Manager</p>
      </div>
    </footer>
  </div>
</template>

<script>
import { useAuthStore } from './stores/auth'
import { useRouter } from 'vue-router'
import { onMounted } from 'vue'

export default {
  name: 'App',
  setup() {
    const authStore = useAuthStore()
    const router = useRouter()

    const logout = async () => {
      await authStore.logout()
      // Redirect to login page after logout
      router.push({ name: 'home' })
    }

    // Fetch the current user when the component is mounted
    onMounted(async () => {
      try {
        await authStore.getCurrentUser()
      } catch (error) {
        console.error('Failed to fetch user:', error)
      }
    })

    return {
      authStore,
      logout
    }
  }
}
</script>

<style lang="scss">
:root {
  --primary: #0135aa;
  --secondary: #fcd211;
  --accent: #cf0822;
  --light: #f8f9fa;
  --dark: #343a40;
  --gray: #6c757d;
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Arial', sans-serif;
  line-height: 1.6;
  color: var(--dark);
  background-color: var(--light);
}

.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.app-header {
  background-color: var(--primary);
  color: white;
  padding: 1rem 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

  .container {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .app-title {
    font-size: 1.5rem;
    margin: 0;
  }
}

.app-footer {
  background-color: var(--dark);
  color: white;
  padding: 1rem 0;
  margin-top: 2rem;
  text-align: center;
}

main {
  padding: 2rem 0;
  min-height: calc(100vh - 140px);
}

.btn {
  display: inline-block;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s, transform 0.2s;

  &:hover {
    transform: translateY(-2px);
  }

  &:active {
    transform: translateY(0);
  }
}

.btn-primary {
  background-color: var(--primary);
  color: white;

  &:hover {
    //background-color: darken(#0135aa, 10%);
  }
}

.btn-secondary {
  background-color: var(--secondary);
  color: var(--dark);

  &:hover {
    //background-color: darken(#fcd211, 10%);
  }
}

.btn-accent {
  background-color: var(--accent);
  color: white;

  &:hover {
    //background-color: darken(#cf0822, 10%);
  }
}

.btn-logout {
  background-color: var(--accent);
  color: white;

  &:hover {
    //background-color: darken(#cf0822, 10%);
  }
}
</style>
