import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface UsageMetrics {
  // Monthly metrics (reset every month)
  aiSummaries: number
  mcpCommands: number
  reports: number
  
  // Daily metrics (reset every day)
  datajudQueries: number
  
  // Lifetime metrics  
  totalProcesses: number
  totalUsers: number
  
  // Last reset dates
  lastMonthlyReset: string
  lastDailyReset: string
}

interface UsageState {
  metrics: UsageMetrics
  
  // Actions
  incrementUsage: (type: keyof Omit<UsageMetrics, 'lastMonthlyReset' | 'lastDailyReset'>, amount?: number) => void
  resetDailyMetrics: () => void
  resetMonthlyMetrics: () => void
  checkAndResetMetrics: () => void
  getUsageForType: (type: keyof Omit<UsageMetrics, 'lastMonthlyReset' | 'lastDailyReset'>) => number
}

const initialMetrics: UsageMetrics = {
  aiSummaries: 0,
  mcpCommands: 0,
  reports: 0,
  datajudQueries: 0,
  totalProcesses: 0,
  totalUsers: 0,
  lastMonthlyReset: new Date().toISOString(),
  lastDailyReset: new Date().toISOString()
}

export const useUsageStore = create<UsageState>()(
  persist(
    (set, get) => ({
      metrics: initialMetrics,

      incrementUsage: (type, amount = 1) => {
        // Don't check reset here to avoid loops - let getUsageForType handle it
        set((state) => ({
          metrics: {
            ...state.metrics,
            [type]: state.metrics[type] + amount
          }
        }))
      },

      resetDailyMetrics: () => {
        set((state) => ({
          metrics: {
            ...state.metrics,
            datajudQueries: 0,
            lastDailyReset: new Date().toISOString()
          }
        }))
      },

      resetMonthlyMetrics: () => {
        set((state) => ({
          metrics: {
            ...state.metrics,
            aiSummaries: 0,
            mcpCommands: 0,
            reports: 0,
            lastMonthlyReset: new Date().toISOString()
          }
        }))
      },

      checkAndResetMetrics: () => {
        const { metrics } = get()
        const now = new Date()
        const lastDaily = new Date(metrics.lastDailyReset)
        const lastMonthly = new Date(metrics.lastMonthlyReset)

        // Check if we need to reset daily metrics (new day)
        if (now.getDate() !== lastDaily.getDate() || 
            now.getMonth() !== lastDaily.getMonth() || 
            now.getFullYear() !== lastDaily.getFullYear()) {
          get().resetDailyMetrics()
        }

        // Check if we need to reset monthly metrics (new month)
        if (now.getMonth() !== lastMonthly.getMonth() || 
            now.getFullYear() !== lastMonthly.getFullYear()) {
          get().resetMonthlyMetrics()
        }
      },

      getUsageForType: (type) => {
        // Check if we need to reset metrics first
        get().checkAndResetMetrics()
        return get().metrics[type]
      }
    }),
    {
      name: 'usage-tracking',
      version: 1
    }
  )
)

// Helper functions to simulate real usage based on different patterns
export const simulateRealisticUsage = (tenantPlan: string) => {
  const now = new Date()
  const dayOfMonth = now.getDate()
  const hour = now.getHours()
  
  // Simulate different usage patterns based on plan
  const planMultipliers = {
    starter: 0.3,
    professional: 0.6,
    business: 0.8,
    enterprise: 1.0
  }
  
  const multiplier = planMultipliers[tenantPlan as keyof typeof planMultipliers] || 0.3
  
  // Business hours simulation (9-18h = more usage)
  const businessHourMultiplier = hour >= 9 && hour <= 18 ? 1.5 : 0.5
  
  // Simulate progressive usage through the month
  const monthProgressMultiplier = dayOfMonth / 30
  
  // Weekly pattern (weekdays = more usage)
  const dayOfWeek = now.getDay()
  const weekdayMultiplier = dayOfWeek >= 1 && dayOfWeek <= 5 ? 1.2 : 0.7
  
  return {
    multiplier,
    businessHourMultiplier,
    monthProgressMultiplier,
    weekdayMultiplier,
    getTotalMultiplier: () => multiplier * businessHourMultiplier * monthProgressMultiplier * weekdayMultiplier
  }
}