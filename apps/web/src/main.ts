import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAuthStore } from '@/stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// Restore session from cookies before the router processes the first navigation.
// This ensures isAuthenticated is correct when the beforeEach guard runs,
// preventing a redirect to /login on F5/page refresh.
const authStore = useAuthStore()
await authStore.initSession().catch(() => {
  // Silently ignore errors (e.g. no cookies, network down) — the guard will
  // redirect to login if the user is not authenticated.
})

app.mount('#app')
