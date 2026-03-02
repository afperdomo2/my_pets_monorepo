<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePetStore } from '@/stores/pets'
import type { PetPayload } from '@/types/pet'

const route = useRoute()
const router = useRouter()
const store = usePetStore()

const id = Number(route.params.id)
const pet = ref(store.pets.find(p => p.id === id))
const editing = ref(false)
const form = ref<PetPayload>({ name: '', species: '', breed: '', age: 0, owner: '' })

onMounted(async () => {
  if (store.pets.length === 0) await store.fetchPets()
  const found = store.pets.find(p => p.id === id)
  if (found) {
    pet.value = found
    form.value = { name: found.name, species: found.species, breed: found.breed, age: found.age, owner: found.owner }
  }
})

async function handleUpdate() {
  await store.updatePet(id, form.value)
  pet.value = store.pets.find(p => p.id === id)
  editing.value = false
}

async function handleDelete() {
  if (confirm('Eliminar mascota?')) {
    await store.deletePet(id)
    router.push('/pets')
  }
}
</script>

<template>
  <div class="pet-detail">
    <RouterLink to="/pets" class="back">&larr; Volver</RouterLink>

    <div v-if="!pet" class="status">Mascota no encontrada.</div>

    <template v-else>
      <div v-if="!editing" class="detail-card">
        <h1>{{ pet.name }}</h1>
        <p><strong>Especie:</strong> {{ pet.species }}</p>
        <p><strong>Raza:</strong> {{ pet.breed || '—' }}</p>
        <p><strong>Edad:</strong> {{ pet.age }} años</p>
        <p><strong>Dueño:</strong> {{ pet.owner || '—' }}</p>
        <div class="actions">
          <button class="btn-primary" @click="editing = true">Editar</button>
          <button class="btn-danger" @click="handleDelete">Eliminar</button>
        </div>
      </div>

      <form v-else class="pet-form" @submit.prevent="handleUpdate">
        <h2>Editar mascota</h2>
        <input v-model="form.name" placeholder="Nombre *" required />
        <input v-model="form.species" placeholder="Especie *" required />
        <input v-model="form.breed" placeholder="Raza" />
        <input v-model.number="form.age" type="number" placeholder="Edad" min="0" />
        <input v-model="form.owner" placeholder="Dueño" />
        <div class="actions">
          <button type="submit" class="btn-primary">Guardar</button>
          <button type="button" class="btn-secondary" @click="editing = false">Cancelar</button>
        </div>
      </form>
    </template>
  </div>
</template>

<style scoped>
.pet-detail { width: 100%; padding: var(--space-8) var(--space-10); }
.back { color: #42b883; text-decoration: none; font-weight: 600; }
.back:hover { text-decoration: underline; }
.detail-card { margin-top: 1.5rem; padding: 1.5rem; background: #fff; border: 1px solid #e8e8e8; border-radius: 10px; }
.detail-card h1 { margin-top: 0; }
.detail-card p { margin: 0.5rem 0; }
.pet-form { display: flex; flex-direction: column; gap: 0.75rem; margin-top: 1.5rem; }
.pet-form h2 { margin: 0; }
.pet-form input { padding: 0.5rem 0.75rem; border: 1px solid #ddd; border-radius: 6px; font-size: 0.95rem; }
.actions { display: flex; gap: 0.75rem; margin-top: 1rem; }
.btn-primary { padding: 0.5rem 1.25rem; background: #42b883; color: #fff; border: none; border-radius: 6px; cursor: pointer; font-weight: 600; }
.btn-primary:hover { background: #369a6e; }
.btn-secondary { padding: 0.5rem 1.25rem; background: #eee; color: #333; border: none; border-radius: 6px; cursor: pointer; }
.btn-danger { padding: 0.5rem 1.25rem; background: #e74c3c; color: #fff; border: none; border-radius: 6px; cursor: pointer; }
.status { margin-top: 2rem; color: #888; }
</style>
