import type { DashboardSummaryResponse } from '@/types/dashboard'
import { get } from '@/services/http'

export const dashboardService = {
  getSummary(): Promise<DashboardSummaryResponse> {
    return get('/dashboard/summary')
  },
}
