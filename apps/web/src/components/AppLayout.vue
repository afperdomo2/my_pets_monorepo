<script setup lang="ts">
import AppSidebar from './AppSidebar.vue'
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="app-layout">
    <AppSidebar />
    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
  background: var(--color-bg);
}

.app-main {
  flex: 1;
  margin-left: var(--sidebar-width);
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
</style>
