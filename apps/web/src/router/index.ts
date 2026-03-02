import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '../components/AppLayout.vue'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Auth routes (sin layout)
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

    // App routes (con sidebar layout)
    {
      path: '/',
      component: AppLayout,
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
      ],
    },

    // Redirect raíz a login
    {
      path: '/:pathMatch(.*)*',
      redirect: '/login',
    },
  ],
})

export default router
