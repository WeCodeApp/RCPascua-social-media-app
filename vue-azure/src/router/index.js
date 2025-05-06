import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

// Views
// const Home = () => import('../views/Home.vue')
const Login = () => import('../views/Login.vue')
const PostList = () => import('../views/PostList.vue')
const TaskList = () => import('../views/TaskList.vue')
const NotFound = () => import('../views/NotFound.vue')

const routes = [
  {
    path: '/',
    name: 'home',
    component: Login
  },
  {
    path: '/login',
    name: 'login',
    component: Login,
    meta: { requiresGuest: true }
  },
  {
    path: '/tasks',
    name: 'tasks',
    component: TaskList,
    meta: { requiresAuth: true }
  },
  {
    path: '/posts',
    name: 'posts',
    component: PostList,
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Navigation guards
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const requiresGuest = to.matched.some(record => record.meta.requiresGuest)

  //!localStorage.getItem('access_token') !authStore.isAuthenticated
  if (requiresAuth && !localStorage.getItem('access_token')) {
    next({ name: 'login' })
  } else if (requiresGuest && authStore.isAuthenticated) {
    next({ name: 'posts' })
  } else {
    next()
  }
})

export default router
