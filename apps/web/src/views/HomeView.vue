<template>
  <main class="home">
    <h1>My Pets</h1>
    <p>Gestiona tus mascotas de forma sencilla.</p>
    <RouterLink to="/pets" class="cta-btn">Ver mascotas</RouterLink>
  </main>
</template>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  gap: 1.5rem;
  text-align: center;
}
h1 {
  font-size: 2.5rem;
  font-weight: 700;
}
.cta-btn {
  padding: 0.75rem 2rem;
  background: #42b883;
  color: #fff;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  font-size: 1rem;
  transition: background 0.2s;
}
.cta-btn:hover {
  background: #369a6e;
}
</style>
