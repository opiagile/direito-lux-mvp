import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { Process } from '@/types'
import { toast } from 'sonner'

// Interface para o store de processos
interface ProcessState {
  processes: Process[]
  isLoading: boolean
  
  // Actions
  addProcess: (process: Omit<Process, 'id' | 'createdAt' | 'updatedAt'>) => void
  updateProcess: (id: string, updates: Partial<Process>) => void
  deleteProcess: (id: string) => void
  toggleMonitoring: (id: string) => void
  getProcessById: (id: string) => Process | undefined
  getProcessesByFilter: (filter: {
    status?: string[]
    priority?: string[]
    court?: string[]
    monitoring?: boolean
    search?: string
  }) => Process[]
  
  // Stats
  getStats: () => {
    total: number
    active: number
    monitoring: number
    highPriority: number
    byStatus: Record<string, number>
    byPriority: Record<string, number>
  }
}

// Dados iniciais para demonstração
const initialProcesses: Process[] = [
  {
    id: '1',
    tenantId: '11111111-1111-1111-1111-111111111111',
    number: '5001234-20.2023.4.03.6109',
    type: 'Ação de Cobrança',
    subject: 'Cobrança de honorários advocatícios - João Silva vs. Empresa ABC',
    court: 'TJSP - 1ª Vara Cível',
    status: 'active',
    parties: [
      { id: '1', name: 'João Silva', document: '123.456.789-00', type: 'person', role: 'plaintiff' },
      { id: '2', name: 'Empresa ABC Ltda', document: '12.345.678/0001-90', type: 'company', role: 'defendant' }
    ],
    movements: [],
    monitoring: true,
    tags: ['cobrança', 'honorários', 'cível'],
    priority: 'high',
    lawyer: 'Dr. Carlos Oliveira',
    estimatedValue: 50000,
    createdAt: '2025-01-15T08:00:00Z',
    updatedAt: '2025-01-20T10:30:00Z',
    lastMovement: '2025-01-20T10:30:00Z'
  },
  {
    id: '2',
    tenantId: '11111111-1111-1111-1111-111111111111',
    number: '5009876-15.2023.4.03.6109',
    type: 'Ação Trabalhista',
    subject: 'Rescisão indireta por falta de pagamento - Pedro Costa vs. Empresa XYZ',
    court: 'TRT - 2ª Região',
    status: 'active',
    parties: [
      { id: '3', name: 'Pedro Costa', document: '987.654.321-00', type: 'person', role: 'plaintiff' },
      { id: '4', name: 'Empresa XYZ Ltda', document: '98.765.432/0001-10', type: 'company', role: 'defendant' }
    ],
    movements: [],
    monitoring: false,
    tags: ['trabalhista', 'rescisão', 'indireta'],
    priority: 'medium',
    lawyer: 'Dra. Ana Paula',
    estimatedValue: 25000,
    createdAt: '2025-01-10T14:30:00Z',
    updatedAt: '2025-01-17T14:15:00Z',
    lastMovement: '2025-01-17T14:15:00Z'
  },
  {
    id: '3',
    tenantId: '11111111-1111-1111-1111-111111111111',
    number: '5005555-30.2023.4.03.6109',
    type: 'Divórcio Consensual',
    subject: 'Divórcio consensual com partilha de bens - Roberto e Sandra Lima',
    court: 'TJSP - Vara de Família',
    status: 'concluded',
    parties: [
      { id: '5', name: 'Roberto Lima', document: '111.222.333-44', type: 'person', role: 'plaintiff' },
      { id: '6', name: 'Sandra Lima', document: '555.666.777-88', type: 'person', role: 'defendant' }
    ],
    movements: [],
    monitoring: true,
    tags: ['família', 'divórcio', 'consensual'],
    priority: 'low',
    lawyer: 'Dr. Carlos Oliveira',
    estimatedValue: 0,
    createdAt: '2024-12-20T09:00:00Z',
    updatedAt: '2025-01-15T09:00:00Z',
    lastMovement: '2025-01-15T09:00:00Z'
  }
]

export const useProcessStore = create<ProcessState>()(
  persist(
    (set, get) => ({
      processes: initialProcesses,
      isLoading: false,

      addProcess: (processData) => {
        const newProcess: Process = {
          ...processData,
          id: `proc_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          lastMovement: new Date().toISOString()
        }

        set((state) => ({
          processes: [...state.processes, newProcess]
        }))

        toast.success('Processo criado com sucesso!')
        return newProcess.id
      },

      updateProcess: (id, updates) => {
        set((state) => {
          const updatedProcesses = state.processes.map((process) =>
            process.id === id
              ? { 
                  ...process, 
                  ...updates, 
                  updatedAt: new Date().toISOString() 
                }
              : process
          )
          
          return {
            ...state,
            processes: [...updatedProcesses] // Force new array reference
          }
        })

        toast.success('Processo atualizado com sucesso!')
      },

      deleteProcess: (id) => {
        const process = get().processes.find(p => p.id === id)
        if (!process) {
          toast.error('Processo não encontrado')
          return
        }

        set((state) => ({
          processes: state.processes.filter((process) => process.id !== id)
        }))

        toast.success(`Processo ${process.number} excluído com sucesso!`)
      },

      toggleMonitoring: (id) => {
        const process = get().processes.find(p => p.id === id)
        if (!process) {
          toast.error('Processo não encontrado')
          return
        }

        const newMonitoringStatus = !process.monitoring

        set((state) => ({
          processes: state.processes.map((p) =>
            p.id === id
              ? { ...p, monitoring: newMonitoringStatus, updatedAt: new Date().toISOString() }
              : p
          )
        }))

        toast.success(
          newMonitoringStatus 
            ? `Processo ${process.number} adicionado ao monitoramento`
            : `Processo ${process.number} removido do monitoramento`
        )
      },

      getProcessById: (id) => {
        return get().processes.find(p => p.id === id)
      },

      getProcessesByFilter: (filter) => {
        const { processes } = get()
        
        return processes.filter(process => {
          // Status filter
          if (filter.status && filter.status.length > 0) {
            if (!filter.status.includes(process.status)) return false
          }

          // Priority filter
          if (filter.priority && filter.priority.length > 0) {
            if (!filter.priority.includes(process.priority)) return false
          }

          // Court filter
          if (filter.court && filter.court.length > 0) {
            if (!filter.court.some(court => process.court.toLowerCase().includes(court.toLowerCase()))) {
              return false
            }
          }

          // Monitoring filter
          if (filter.monitoring !== undefined) {
            if (process.monitoring !== filter.monitoring) return false
          }

          // Search filter
          if (filter.search && filter.search.length > 0) {
            const searchLower = filter.search.toLowerCase()
            return (
              process.number.toLowerCase().includes(searchLower) ||
              process.subject.toLowerCase().includes(searchLower) ||
              process.court.toLowerCase().includes(searchLower) ||
              process.type.toLowerCase().includes(searchLower) ||
              process.lawyer?.toLowerCase().includes(searchLower) ||
              process.tags.some(tag => tag.toLowerCase().includes(searchLower)) ||
              process.parties.some(party => party.name.toLowerCase().includes(searchLower))
            )
          }

          return true
        })
      },

      getStats: () => {
        const { processes } = get()
        
        const stats = {
          total: processes.length,
          active: processes.filter(p => p.status === 'active').length,
          monitoring: processes.filter(p => p.monitoring).length,
          highPriority: processes.filter(p => p.priority === 'high').length,
          byStatus: {} as Record<string, number>,
          byPriority: {} as Record<string, number>
        }

        // Count by status
        processes.forEach(process => {
          stats.byStatus[process.status] = (stats.byStatus[process.status] || 0) + 1
        })

        // Count by priority  
        processes.forEach(process => {
          stats.byPriority[process.priority] = (stats.byPriority[process.priority] || 0) + 1
        })

        return stats
      }
    }),
    {
      name: 'processes-storage',
      partialize: (state) => ({ 
        processes: state.processes 
      }),
    }
  )
)