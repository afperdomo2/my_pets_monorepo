import { ref, watch } from 'vue'
import { defineStore } from 'pinia'

const USERS_PER_PAGE_KEY = 'ui:usersPerPage'
const VALID_VALUES = [10, 25, 50] as const
type PerPageValue = (typeof VALID_VALUES)[number]

function readPerPage(): PerPageValue {
  const raw = localStorage.getItem(USERS_PER_PAGE_KEY)
  const parsed = Number(raw)
  return (VALID_VALUES as readonly number[]).includes(parsed) ? (parsed as PerPageValue) : 10
}

export const useUIStore = defineStore('ui', () => {
  const usersPerPage = ref<PerPageValue>(readPerPage())

  watch(usersPerPage, (val) => {
    localStorage.setItem(USERS_PER_PAGE_KEY, String(val))
  })

  return { usersPerPage }
})
