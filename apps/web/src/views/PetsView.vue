<script setup lang="ts">
import { ref } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { useGetPets, useCreatePet, useDeletePet } from '@/composables/usePets'
import { petSchema } from '@/schemas/pet'

const { data: pets, isLoading, isError, error } = useGetPets()
const createPet = useCreatePet()
const deletePet = useDeletePet()

const showForm = ref(false)
const submitError = ref<string | null>(null)

const { defineField, handleSubmit, errors, resetForm } = useForm({
  validationSchema: toTypedSchema(petSchema),
  initialValues: { name: '', species: '', breed: '', age: undefined, owner: '' },
})

const [name, nameAttrs] = defineField('name')
const [species, speciesAttrs] = defineField('species')
const [breed, breedAttrs] = defineField('breed')
const [age, ageAttrs] = defineField('age')
const [owner, ownerAttrs] = defineField('owner')

const handleCreate = handleSubmit(async (values) => {
  submitError.value = null
  try {
    await createPet.mutateAsync(values)
    showForm.value = false
    resetForm()
  } catch (e) {
    submitError.value = e instanceof Error ? e.message : 'Error al crear mascota'
  }
})

function handleCancel() {
  showForm.value = false
  resetForm()
  submitError.value = null
}

async function handleDelete(id: string) {
  if (confirm('Eliminar mascota?')) await deletePet.mutateAsync(id)
}
</script>

<template>
  <div class="pets-view">
    <header class="pets-header">
      <h1>Mascotas</h1>
      <button class="btn-primary" @click="showForm ? handleCancel() : (showForm = true)">
        {{ showForm ? 'Cancelar' : '+ Agregar' }}
      </button>
    </header>

    <form v-if="showForm" class="pet-form" @submit.prevent="handleCreate">
      <div class="form-field">
        <input v-model="name" v-bind="nameAttrs" placeholder="Nombre *" />
        <p v-if="errors.name" class="field-error">{{ errors.name }}</p>
      </div>

      <div class="form-field">
        <select v-model="species" v-bind="speciesAttrs">
          <option value="" disabled>Especie *</option>
          <option value="dog">Perro</option>
          <option value="cat">Gato</option>
          <option value="bird">Ave</option>
          <option value="rabbit">Conejo</option>
          <option value="fish">Pez</option>
          <option value="other">Otro</option>
        </select>
        <p v-if="errors.species" class="field-error">{{ errors.species }}</p>
      </div>

      <div class="form-field">
        <input v-model="breed" v-bind="breedAttrs" placeholder="Raza" />
      </div>

      <div class="form-field">
        <input
          v-model.number="age"
          v-bind="ageAttrs"
          type="number"
          placeholder="Edad"
          min="0"
          max="100"
        />
        <p v-if="errors.age" class="field-error">{{ errors.age }}</p>
      </div>

      <div class="form-field">
        <input v-model="owner" v-bind="ownerAttrs" placeholder="Dueño" />
      </div>

      <p v-if="submitError" class="form-error">{{ submitError }}</p>

      <button type="submit" class="btn-primary" :disabled="createPet.isPending.value">Guardar</button>
    </form>

    <div v-if="isLoading" class="status">Cargando...</div>

    <div v-else-if="isError" class="status">Error: {{ error?.message }}</div>

    <ul v-else class="pet-list">
      <li v-for="pet in pets" :key="pet.id" class="pet-card">
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
.pets-view { width: 100%; padding: var(--space-8) var(--space-10); }
.pets-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.pet-form { display: flex; flex-wrap: wrap; gap: 0.5rem; margin-bottom: 2rem; padding: 1rem; background: #f9f9f9; border-radius: 8px; }
.form-field { flex: 1 1 180px; display: flex; flex-direction: column; gap: 0.25rem; }
.pet-form input,
.pet-form select { width: 100%; padding: 0.5rem 0.75rem; border: 1px solid #ddd; border-radius: 6px; font-size: 0.9rem; background: #fff; }
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
.field-error { font-size: 0.75rem; color: #dc2626; margin: 0; }
.form-error { flex: 1 1 100%; font-size: 0.85rem; color: #dc2626; background: #fef2f2; border: 1px solid #fecaca; border-radius: 6px; padding: 0.4rem 0.75rem; margin: 0; }
</style>
