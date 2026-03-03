import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'

import App from './App.vue'
import router from './router'
import { useAuthStore } from '@/stores/auth'

const app = createApp(App)
const pinia = createPinia()

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 60_000,       // 1 min — listas
      gcTime: 5 * 60_000,      // 5 min en cache tras desmontarse
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

app.use(pinia)
app.use(router)
app.use(VueQueryPlugin, { queryClient })

// Restore session from cookies before the router processes the first navigation.
// This ensures isAuthenticated is correct when the beforeEach guard runs,
// preventing a redirect to /login on F5/page refresh.
const authStore = useAuthStore()
await authStore.initSession().catch(() => {
  // Silently ignore errors (e.g. no cookies, network down) — the guard will
  // redirect to login if the user is not authenticated.
})

app.mount('#app')

