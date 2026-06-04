import { examService } from '@/services/examService'
import type {
  CompleteExamPayload,
  CreateExamPayload,
  Exam,
  ExamResult,
  ExamWithResults,
  ScheduleExamPayload,
  UpdateExamPayload,
} from '@/types/exam'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, nextTick, type Ref } from 'vue'

// ═══════════════════════════════════════════════════════════════════════════════════
// Cache time constants
// ═══════════════════════════════════════════════════════════════════════════════════
const EXAM_LIST_STALE_TIME = 2 * 60_000 // 2 minutos

export function useGetAllExams(page: Ref<number>, perPage: Ref<number>) {
  const queryClient = useQueryClient()
  const queryKey = computed(() => ['exams', 'all', { page: page.value, perPage: perPage.value }])

  const query = useQuery({
    queryKey,
    queryFn: () => examService.getAll(page.value, perPage.value),
    staleTime: EXAM_LIST_STALE_TIME,
  })

  let refreshing = false
  async function refresh() {
    if (refreshing) return
    refreshing = true
    try {
      page.value = 1
      await nextTick()
      await queryClient.invalidateQueries({ queryKey: ['exams', 'all'], refetchType: 'active' })
    } finally {
      refreshing = false
    }
  }

  return { ...query, refresh }
}

export function useGetExamsByPet(petId: Ref<string>, page: Ref<number>, perPage: Ref<number>) {
  const queryClient = useQueryClient()
  const queryKey = computed(() => [
    'exams',
    'pet',
    petId.value,
    { page: page.value, perPage: perPage.value },
  ])

  const query = useQuery({
    queryKey,
    queryFn: () => examService.getByPetId(petId.value, page.value, perPage.value),
    staleTime: EXAM_LIST_STALE_TIME,
    enabled: computed(() => !!petId.value),
  })

  let refreshing = false
  async function refresh() {
    if (refreshing) return
    refreshing = true
    try {
      page.value = 1
      await nextTick()
      await queryClient.invalidateQueries({
        queryKey: ['exams', 'pet', petId.value],
        refetchType: 'active',
      })
    } finally {
      refreshing = false
    }
  }

  return { ...query, refresh }
}

export function useGetExamById(id: Ref<string>) {
  const queryKey = computed(() => ['exams', 'detail', id.value])

  return useQuery({
    queryKey,
    queryFn: () =>
      examService.getById(id.value).then((r) => {
        const raw = r as unknown as { data: Exam; results?: ExamResult[] }
        return { ...raw.data, results: raw.results ?? [] } as ExamWithResults
      }),
    enabled: computed(() => !!id.value),
    staleTime: 0,
  })
}

export function useCreateExam() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: CreateExamPayload) => examService.create(payload).then((r) => r.data),
    onSuccess: (newItem) => {
      queryClient.invalidateQueries({ queryKey: ['exams', 'pet', newItem.pet_id] })
      queryClient.invalidateQueries({ queryKey: ['exams', 'all'] })
    },
  })
}

export function useUpdateExam() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateExamPayload }) =>
      examService.update(id, payload).then((r) => r.data),
    onSuccess: (updatedItem) => {
      queryClient.invalidateQueries({ queryKey: ['exams', 'detail', updatedItem.id] })
      queryClient.invalidateQueries({ queryKey: ['exams', 'pet', updatedItem.pet_id] })
    },
  })
}

export function useScheduleExam() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: ScheduleExamPayload }) =>
      examService.schedule(id, payload).then((r) => r.data),
    onSuccess: (updatedItem) => {
      queryClient.invalidateQueries({ queryKey: ['exams', 'detail', updatedItem.id] })
      queryClient.invalidateQueries({ queryKey: ['exams', 'pet', updatedItem.pet_id] })
    },
  })
}

export function useCompleteExam() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: CompleteExamPayload }) =>
      examService.complete(id, payload).then((r) => r.data),
    onSuccess: (updatedItem) => {
      queryClient.invalidateQueries({ queryKey: ['exams', 'detail', updatedItem.id] })
      queryClient.invalidateQueries({ queryKey: ['exams', 'pet', updatedItem.pet_id] })
    },
  })
}

export function useDeleteExam() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (variables: { id: string; petId: string }) => examService.remove(variables.id),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['exams', 'pet', variables.petId] })
      queryClient.invalidateQueries({ queryKey: ['exams', 'all'] })
    },
  })
}
