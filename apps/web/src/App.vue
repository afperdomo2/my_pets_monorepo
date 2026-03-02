<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header class="navbar">
    <span class="brand">My Pets</span>
    <nav>
      <RouterLink to="/">Inicio</RouterLink>
      <RouterLink to="/pets">Mascotas</RouterLink>
    </nav>
  </header>

  <RouterView />
</template>

<style scoped>
.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 2rem;
  background: #2c3e50;
  color: #fff;
}
.brand {
  font-size: 1.25rem;
  font-weight: 700;
  letter-spacing: 0.5px;
}
nav {
  display: flex;
  gap: 1.5rem;
}
nav a {
  color: #a8d8c2;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;
}
nav a:hover,
nav a.router-link-active {
  color: #42b883;
}
</style>
