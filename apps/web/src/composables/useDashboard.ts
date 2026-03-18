import { dashboardService } from '@/services/dashboardService'
import type { DashboardSummary } from '@/types/dashboard'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const DASHBOARD_STALE_TIME = 5 * 60_000 // 5 minutos

export function useDashboardSummary() {
  const query = useQuery({
    queryKey: ['dashboard', 'summary'],
    queryFn: () => dashboardService.getSummary().then((r) => r.data),
    staleTime: DASHBOARD_STALE_TIME,
  })

  const summary = computed<DashboardSummary | undefined>(() => query.data.value)
  const totalPets = computed(() => summary.value?.total_pets ?? 0)
  const healthyPets = computed(() => summary.value?.healthy_pets ?? 0)
  const pendingTasks = computed(() => summary.value?.pending_tasks ?? 0)
  const overdueTasks = computed(() => summary.value?.overdue_tasks ?? 0)

  return {
    summary,
    totalPets,
    healthyPets,
    pendingTasks,
    overdueTasks,
    isLoading: query.isLoading,
    isFetching: query.isFetching,
    isError: query.isError,
    error: query.error,
    refresh: query.refetch,
  }
}
