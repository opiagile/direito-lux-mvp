import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { useAuthStore } from './auth'
import { useProcessStore } from './processes'
import { useUsageStore } from './usage'

interface Invoice {
  id: string
  number: string
  date: string
  period: string
  amount: number
  status: 'paid' | 'pending' | 'overdue' | 'cancelled'
  dueDate: string
  paidAt?: string
  downloadUrl?: string
}

interface Usage {
  processes: { used: number; limit: number }
  users: { used: number; limit: number }
  mcpCommands: { used: number; limit: number }
  aiSummaries: { used: number; limit: number }
  reports: { used: number; limit: number }
  datajudQueries: { used: number; limit: number }
}

interface PaymentMethod {
  id: string
  type: 'credit_card' | 'boleto' | 'pix'
  last4?: string
  brand?: string
  expiryMonth?: number
  expiryYear?: number
  isDefault: boolean
}

interface BillingState {
  invoices: Invoice[]
  currentUsage: Usage | null
  paymentMethod: PaymentMethod | null
  isLoading: boolean
  
  // Actions
  loadBillingData: () => void
  updatePaymentMethod: (method: PaymentMethod) => void
  downloadInvoice: (invoiceId: string) => void
  calculateCurrentUsage: () => Usage
}

// Plan limits based on subscription type
const planLimits = {
  starter: {
    processes: 50,
    users: 2,
    mcpCommands: 0,
    aiSummaries: 10,
    reports: 10,
    datajudQueries: 100
  },
  professional: {
    processes: 200,
    users: 5,
    mcpCommands: 200,
    aiSummaries: 50,
    reports: 100,
    datajudQueries: 500
  },
  business: {
    processes: 500,
    users: 15,
    mcpCommands: 1000,
    aiSummaries: 200,
    reports: 500,
    datajudQueries: 2000
  },
  enterprise: {
    processes: -1, // unlimited
    users: -1,
    mcpCommands: -1,
    aiSummaries: -1,
    reports: -1,
    datajudQueries: 10000
  }
}

export const useBillingStore = create<BillingState>()(
  persist(
    (set, get) => ({
      invoices: [],
      currentUsage: null,
      paymentMethod: null,
      isLoading: false,

      loadBillingData: () => {
        set({ isLoading: true })
        
        try {
          // Get tenant info
          const { tenant } = useAuthStore.getState()
          if (!tenant) {
            console.log('No tenant found in billing load')
            set({ isLoading: false })
            return
          }

          // Generate invoices based on tenant creation date
          const invoices = generateInvoicesForTenant(tenant)
          
          // Get current usage
          const currentUsage = get().calculateCurrentUsage()
          
          // Set default payment method for demo
          const paymentMethod: PaymentMethod = {
            id: 'pm_1',
            type: 'credit_card',
            last4: '4532',
            brand: 'visa',
            expiryMonth: 12,
            expiryYear: 2027,
            isDefault: true
          }

          set({ 
            invoices, 
            currentUsage, 
            paymentMethod,
            isLoading: false 
          })
        } catch (error) {
          console.error('Error loading billing data:', error)
          set({ isLoading: false })
        }
      },

      calculateCurrentUsage: () => {
        const { tenant } = useAuthStore.getState()
        const { processes } = useProcessStore.getState()
        
        if (!tenant) {
          return {
            processes: { used: 0, limit: 0 },
            users: { used: 0, limit: 0 },
            mcpCommands: { used: 0, limit: 0 },
            aiSummaries: { used: 0, limit: 0 },
            reports: { used: 0, limit: 0 },
            datajudQueries: { used: 0, limit: 0 }
          }
        }

        const limits = planLimits[tenant.plan as keyof typeof planLimits] || planLimits.starter
        
        // Get real usage
        const processCount = processes.length
        const userCount = 1 // We'll simulate having multiple users per tenant
        
        // Get usage from store if available, otherwise use demo values
        let aiSummariesUsed = 0
        let reportsUsed = 0
        let datajudUsed = 0
        let mcpUsed = 0
        
        try {
          const usageStore = useUsageStore.getState()
          
          // Get base multiplier for demo data
          const today = new Date()
          const dayOfMonth = today.getDate()
          const baseMultiplier = dayOfMonth / 30 // Progress through month
          
          // Calculate usage with tracked data + demo simulation
          aiSummariesUsed = Math.min(
            usageStore.metrics.aiSummaries + Math.floor(limits.aiSummaries * baseMultiplier * 0.2),
            limits.aiSummaries === -1 ? 999999 : limits.aiSummaries
          )
          
          reportsUsed = Math.min(
            usageStore.metrics.reports + Math.floor(limits.reports * baseMultiplier * 0.15),
            limits.reports === -1 ? 999999 : limits.reports
          )
          
          // DataJud is daily, so use less of the month multiplier
          datajudUsed = Math.min(
            usageStore.metrics.datajudQueries + Math.floor(limits.datajudQueries * 0.3),
            limits.datajudQueries
          )
          
          // MCP commands if available
          mcpUsed = limits.mcpCommands > 0 
            ? Math.min(
                usageStore.metrics.mcpCommands + Math.floor(limits.mcpCommands * baseMultiplier * 0.25),
                limits.mcpCommands === -1 ? 999999 : limits.mcpCommands
              )
            : 0
        } catch (e) {
          console.error('Error getting usage metrics:', e)
          // Use fallback demo values
          const dayOfMonth = new Date().getDate()
          aiSummariesUsed = Math.floor((dayOfMonth / 30) * limits.aiSummaries * 0.3)
          reportsUsed = Math.floor((dayOfMonth / 30) * limits.reports * 0.2)
          datajudUsed = Math.floor(limits.datajudQueries * 0.4)
          mcpUsed = limits.mcpCommands > 0 ? Math.floor((dayOfMonth / 30) * limits.mcpCommands * 0.3) : 0
        }

        return {
          processes: { 
            used: processCount, 
            limit: limits.processes 
          },
          users: { 
            used: userCount, 
            limit: limits.users 
          },
          mcpCommands: { 
            used: mcpUsed, 
            limit: limits.mcpCommands 
          },
          aiSummaries: { 
            used: aiSummariesUsed, 
            limit: limits.aiSummaries 
          },
          reports: { 
            used: reportsUsed, 
            limit: limits.reports 
          },
          datajudQueries: { 
            used: datajudUsed, 
            limit: limits.datajudQueries 
          }
        }
      },

      updatePaymentMethod: (method) => {
        set({ paymentMethod: method })
      },

      downloadInvoice: (invoiceId) => {
        // In a real app, this would trigger a download
        console.log('Downloading invoice:', invoiceId)
      }
    }),
    {
      name: 'billing-storage',
      partialize: (state) => ({
        invoices: state.invoices,
        paymentMethod: state.paymentMethod
      })
    }
  )
)

// Helper function to generate invoices based on tenant
function generateInvoicesForTenant(tenant: any): Invoice[] {
  const invoices: Invoice[] = []
  const planPrices = {
    starter: 99,
    professional: 299,
    business: 699,
    enterprise: 1999
  }

  const price = planPrices[tenant.plan as keyof typeof planPrices] || 99
  
  // Generate last 12 months of invoices
  const today = new Date()
  for (let i = 0; i < 12; i++) {
    const invoiceDate = new Date(today.getFullYear(), today.getMonth() - i, 1)
    const dueDate = new Date(invoiceDate.getFullYear(), invoiceDate.getMonth(), 10)
    const isPaid = invoiceDate < today
    
    invoices.push({
      id: `inv_${i}`,
      number: `INV-${invoiceDate.getFullYear()}-${String(invoiceDate.getMonth() + 1).padStart(2, '0')}`,
      date: invoiceDate.toISOString(),
      period: new Intl.DateTimeFormat('pt-BR', { month: 'long', year: 'numeric' }).format(invoiceDate),
      amount: price,
      status: isPaid ? 'paid' : 'pending',
      dueDate: dueDate.toISOString(),
      paidAt: isPaid ? new Date(dueDate.getTime() - 2 * 24 * 60 * 60 * 1000).toISOString() : undefined,
      downloadUrl: isPaid ? `/api/invoices/${invoiceDate.getFullYear()}/${invoiceDate.getMonth() + 1}` : undefined
    })
  }

  return invoices
}