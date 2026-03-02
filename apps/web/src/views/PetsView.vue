<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { usePetStore } from '@/stores/pets'
import type { PetPayload } from '@/types/pet'

const store = usePetStore()

const showForm = ref(false)
const form = ref<PetPayload>({ name: '', species: '', breed: '', age: 0, owner: '' })

onMounted(() => store.fetchPets())

async function handleCreate() {
  await store.createPet(form.value)
  showForm.value = false
  form.value = { name: '', species: '', breed: '', age: 0, owner: '' }
}

async function handleDelete(id: number) {
  if (confirm('Eliminar mascota?')) await store.deletePet(id)
}
</script>

<template>
  <div class="pets-view">
    <header class="pets-header">
      <h1>Mascotas</h1>
      <button class="btn-primary" @click="showForm = !showForm">
        {{ showForm ? 'Cancelar' : '+ Agregar' }}
      </button>
    </header>

    <form v-if="showForm" class="pet-form" @submit.prevent="handleCreate">
      <input v-model="form.name" placeholder="Nombre *" required />
      <input v-model="form.species" placeholder="Especie * (dog, cat...)" required />
      <input v-model="form.breed" placeholder="Raza" />
      <input v-model.number="form.age" type="number" placeholder="Edad" min="0" />
      <input v-model="form.owner" placeholder="Dueño" />
      <button type="submit" class="btn-primary">Guardar</button>
    </form>

    <div v-if="store.loading" class="status">Cargando...</div>
    <div v-else-if="store.error" class="status error">Error: {{ store.error }}</div>

    <ul v-else class="pet-list">
      <li v-for="pet in store.pets" :key="pet.id" class="pet-card">
        <RouterLink :to="`/pets/${pet.id}`" class="pet-name">
          {{ pet.name }}
        </RouterLink>
        <span class="pet-meta">{{ pet.species }} · {{ pet.breed }}</span>
        <span class="pet-meta">Dueño: {{ pet.owner || '—' }}</span>
        <button class="btn-danger" @click="handleDelete(pet.id)">Eliminar</button>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.pets-view { max-width: 800px; margin: 0 auto; padding: 2rem 1rem; }
.pets-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.pet-form { display: flex; flex-wrap: wrap; gap: 0.5rem; margin-bottom: 2rem; padding: 1rem; background: #f9f9f9; border-radius: 8px; }
.pet-form input { flex: 1 1 180px; padding: 0.5rem 0.75rem; border: 1px solid #ddd; border-radius: 6px; font-size: 0.9rem; }
.pet-list { list-style: none; padding: 0; display: flex; flex-direction: column; gap: 0.75rem; }
.pet-card { display: flex; align-items: center; gap: 1rem; padding: 1rem 1.25rem; background: #fff; border: 1px solid #e8e8e8; border-radius: 10px; }
.pet-name { font-weight: 600; font-size: 1.05rem; text-decoration: none; color: #2c3e50; flex: 1; }
.pet-name:hover { color: #42b883; }
.pet-meta { font-size: 0.85rem; color: #888; }
.btn-primary { padding: 0.5rem 1.25rem; background: #42b883; color: #fff; border: none; border-radius: 6px; cursor: pointer; font-weight: 600; }
.btn-primary:hover { background: #369a6e; }
.btn-danger { padding: 0.4rem 0.9rem; background: #e74c3c; color: #fff; border: none; border-radius: 6px; cursor: pointer; font-size: 0.85rem; margin-left: auto; }
.btn-danger:hover { background: #c0392b; }
.status { text-align: center; padding: 2rem; color: #888; }
.error { color: #e74c3c; }
</style>
