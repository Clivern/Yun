import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Login from '@/views/Login.vue'
import Dashboard from '@/views/Dashboard.vue'
import Setup from '@/views/Setup.vue'
import Settings from '@/views/Settings.vue'
import Gateways from '@/views/Gateways.vue'
import Users from '@/views/Users.vue'
import Servers from '@/views/Servers.vue'
import Tools from '@/views/Tools.vue'
import Mcps from '@/views/Mcps.vue'

const routes = [
  {
    path: '/setup',
    name: 'Setup',
    component: Setup,
    meta: { requiresGuest: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresGuest: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/gateways',
    name: 'Gateways',
    component: Gateways,
    meta: { requiresAuth: true }
  },
  {
    path: '/mcps',
    name: 'Mcps',
    component: Mcps,
    meta: { requiresAuth: true }
  },
  {
    path: '/users',
    name: 'Users',
    component: Users,
    meta: { requiresAuth: true }
  },
  {
    path: '/servers',
    name: 'Servers',
    component: Servers,
    meta: { requiresAuth: true }
  },
  {
    path: '/tools',
    name: 'Tools',
    component: Tools,
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings,
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/login'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard for authentication and setup
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Wait for auth initialization to complete before making decisions
  // This prevents redirects during the initial auth check on page refresh
  while (authStore.isInitializing) {
    await new Promise(resolve => setTimeout(resolve, 50))
  }

  const isAuthenticated = authStore.isAuthenticated

  // Check if setup is completed (you can implement this check via API)
  // For now, we'll allow access to setup page if not authenticated

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresGuest && isAuthenticated) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router

