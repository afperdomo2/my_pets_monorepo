<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { useGetPet, useUpdatePet, useDeletePet } from '@/composables/usePets'
import { petSchema } from '@/schemas/pet'

const route = useRoute()
const router = useRouter()

const id = String(route.params.id)
const editing = ref(false)
const submitError = ref<string | null>(null)

const { data: pet, isLoading, isError } = useGetPet(id)
const updatePet = useUpdatePet()
const deletePet = useDeletePet()

const { defineField, handleSubmit, errors, resetForm } = useForm({
  validationSchema: toTypedSchema(petSchema),
  initialValues: { name: '', species: '', breed: '', age: undefined, owner: '' },
})

const [name, nameAttrs] = defineField('name')
const [species, speciesAttrs] = defineField('species')
const [breed, breedAttrs] = defineField('breed')
const [age, ageAttrs] = defineField('age')
const [owner, ownerAttrs] = defineField('owner')

// Populate form when pet data loads
watch(pet, (found) => {
  if (found) {
    resetForm({
      values: {
        name: found.name,
        species: found.species,
        breed: found.breed,
        age: found.age,
        owner: found.owner,
      },
    })
  }
}, { immediate: true })

const handleUpdate = handleSubmit(async (values) => {
  submitError.value = null
  try {
    await updatePet.mutateAsync({ id, payload: values })
    editing.value = false
  } catch (e) {
    submitError.value = e instanceof Error ? e.message : 'Error al actualizar mascota'
  }
})

async function handleDelete() {
  if (confirm('Eliminar mascota?')) {
    await deletePet.mutateAsync(id)
    router.push('/pets')
  }
}
</script>

<template>
  <div class="pet-detail">
    <RouterLink to="/pets" class="back">&larr; Volver</RouterLink>

    <div v-if="isLoading" class="status">Cargando...</div>

    <div v-else-if="isError || !pet" class="status">Mascota no encontrada.</div>

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

        <div class="actions">
          <button type="submit" class="btn-primary" :disabled="updatePet.isPending.value">Guardar</button>
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
.pet-form { display: flex; flex-direction: column; gap: 0.5rem; margin-top: 1.5rem; }
.pet-form h2 { margin: 0 0 0.5rem; }
.form-field { display: flex; flex-direction: column; gap: 0.25rem; }
.pet-form input,
.pet-form select { padding: 0.5rem 0.75rem; border: 1px solid #ddd; border-radius: 6px; font-size: 0.95rem; background: #fff; }
.actions { display: flex; gap: 0.75rem; margin-top: 1rem; }
.btn-primary { padding: 0.5rem 1.25rem; background: #42b883; color: #fff; border: none; border-radius: 6px; cursor: pointer; font-weight: 600; }
.btn-primary:hover { background: #369a6e; }
.btn-secondary { padding: 0.5rem 1.25rem; background: #eee; color: #333; border: none; border-radius: 6px; cursor: pointer; }
.btn-danger { padding: 0.5rem 1.25rem; background: #e74c3c; color: #fff; border: none; border-radius: 6px; cursor: pointer; }
.status { margin-top: 2rem; color: #888; }
.field-error { font-size: 0.75rem; color: #dc2626; margin: 0; }
.form-error { font-size: 0.85rem; color: #dc2626; background: #fef2f2; border: 1px solid #fecaca; border-radius: 6px; padding: 0.4rem 0.75rem; margin: 0; }
</style>
