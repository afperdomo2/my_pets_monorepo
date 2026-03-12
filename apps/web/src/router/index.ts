import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '../components/AppLayout.vue'
import HomeView from '../views/HomeView.vue'
import { useAuthStore } from '@/stores/auth'
import { setupService } from '@/services/setupService'

// Cache setup status so we only call the endpoint once per page load
let setupChecked = false
let needsSetup = false

export function resetSetupCache() {
  setupChecked = false
  needsSetup = false
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Auth routes (no layout, no JWT required)
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
    },

    // Setup route (public — only accessible when no users exist)
    {
      path: '/setup',
      name: 'setup',
      component: () => import('../views/SetupView.vue'),
    },

    // App routes (sidebar layout, JWT required)
    {
      path: '/',
      component: AppLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'home',
          component: HomeView,
        },
        {
          path: 'pets',
          name: 'pets',
          component: () => import('../views/PetsView.vue'),
        },
        {
          path: 'pets/:id',
          name: 'pet-detail',
          component: () => import('../views/PetDetailView.vue'),
        },
        {
          path: 'vaccines',
          name: 'vaccines',
          component: () => import('../views/VaccinesView.vue'),
        },
        {
          path: 'reports',
          name: 'reports',
          component: () => import('../views/ReportsView.vue'),
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('../views/SettingsView.vue'),
        },
        {
          path: 'profile',
          name: 'profile',
          component: () => import('../views/ProfileView.vue'),
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('../views/UsersView.vue'),
          meta: { requiresSystemUser: true },
        },
      ],
    },

    // Catch-all → login
    {
      path: '/:pathMatch(.*)*',
      redirect: '/login',
    },
  ],
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  // Wait for the initial session check to complete before evaluating auth state.
  // This prevents race conditions between app startup and the first navigation,
  // which caused intermittent redirects to /login on F5.
  await authStore.sessionReady

  // Check setup status once per page load
  if (!setupChecked) {
    try {
      const status = await setupService.checkStatus()
      needsSetup = status.needs_setup
    } catch {
      needsSetup = false
    }
    setupChecked = true
  }

  // If system needs setup, redirect everything to /setup (except /setup itself)
  if (needsSetup) {
    if (to.name !== 'setup') {
      return { name: 'setup' }
    }
    return
  }

  // System is initialized — redirect away from /setup
  if (to.name === 'setup') {
    return { name: 'login' }
  }

  if (to.meta['requiresAuth'] && !authStore.isAuthenticated) {
    return { name: 'login' }
  }

  // System-user-only routes: redirect non-admins to home
  if (to.meta['requiresSystemUser'] && !authStore.user?.is_system_user) {
    return { name: 'home' }
  }

  // Redirect already-authenticated users away from login/register
  if ((to.name === 'login' || to.name === 'register') && authStore.isAuthenticated) {
    return { name: 'home' }
  }
})

export default router
