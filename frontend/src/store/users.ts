import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { toast } from 'sonner'
import { User, UserRole } from '@/types'
import { useAuthStore } from './auth'
import { useBillingStore } from './billing'

interface UserState {
  users: User[]
  isLoading: boolean
  
  // Actions
  addUser: (userData: Omit<User, 'id' | 'createdAt' | 'updatedAt' | 'tenantId'>) => string | null
  updateUser: (id: string, updates: Partial<User>) => void
  deleteUser: (id: string) => void
  toggleUserStatus: (id: string) => void
  getUserById: (id: string) => User | undefined
  getUsersByRole: (role: UserRole) => User[]
  getActiveUsers: () => User[]
  getUsersForCurrentTenant: () => User[]
  checkUserQuota: () => { canAdd: boolean; used: number; limit: number; message?: string }
  loadInitialUsers: () => void
}

// Initial users with different tenant data for testing
const initialUsers: User[] = [
  {
    id: 'user_1',
    email: 'admin@silvaassociados.com.br',
    name: 'Carlos Silva',
    role: 'admin',
    tenantId: '11111111-1111-1111-1111-111111111111', // Silva & Associados (Starter)
    isActive: true,
    createdAt: '2024-01-15T09:00:00Z',
    updatedAt: '2025-01-20T14:30:00Z'
  },
  {
    id: 'user_2',
    email: 'gerente@silvaassociados.com.br',
    name: 'Ana Paula Santos',
    role: 'manager',
    tenantId: '11111111-1111-1111-1111-111111111111', // Silva & Associados (Starter)
    isActive: true,
    createdAt: '2024-01-16T10:00:00Z',
    updatedAt: '2025-01-19T16:20:00Z'
  },
  // Costa Santos - Professional (up to 5 users)
  {
    id: 'user_3',
    email: 'admin@costasantos.com.br',
    name: 'João Costa',
    role: 'admin',
    tenantId: '22222222-2222-2222-2222-222222222222',
    isActive: true,
    createdAt: '2024-02-01T09:00:00Z',
    updatedAt: '2025-01-20T14:30:00Z'
  },
  {
    id: 'user_4',
    email: 'advogado@costasantos.com.br',
    name: 'Dr. Pedro Santos',
    role: 'lawyer',
    tenantId: '22222222-2222-2222-2222-222222222222',
    isActive: true,
    createdAt: '2024-02-02T09:00:00Z',
    updatedAt: '2025-01-20T14:30:00Z'
  },
  {
    id: 'user_5',
    email: 'assistente@costasantos.com.br',
    name: 'Maria Santos',
    role: 'assistant',
    tenantId: '22222222-2222-2222-2222-222222222222',
    isActive: true,
    createdAt: '2024-02-03T09:00:00Z',
    updatedAt: '2025-01-20T14:30:00Z'
  },
  // Machado Advogados - Business (up to 15 users)
  {
    id: 'user_6',
    email: 'admin@machadoadvogados.com.br',
    name: 'Dr. Roberto Machado',
    role: 'admin',
    tenantId: '33333333-3333-3333-3333-333333333333',
    isActive: true,
    createdAt: '2024-03-01T09:00:00Z',
    updatedAt: '2025-01-20T14:30:00Z'
  }
]

export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      users: [],
      isLoading: false,

      loadInitialUsers: () => {
        const currentUsers = get().users
        if (currentUsers.length === 0) {
          set({ users: initialUsers })
        }
      },

      checkUserQuota: () => {
        const { tenant } = useAuthStore.getState()
        const currentUsers = get().getUsersForCurrentTenant()
        
        if (!tenant) {
          return { canAdd: false, used: 0, limit: 0, message: 'Tenant não encontrado' }
        }

        // Get plan limits
        const planLimits = {
          starter: 2,
          professional: 5,
          business: 15,
          enterprise: -1 // unlimited
        }

        const limit = planLimits[tenant.plan as keyof typeof planLimits] || 2
        const used = currentUsers.filter(u => u.isActive).length

        if (limit === -1) {
          return { canAdd: true, used, limit: -1 }
        }

        const canAdd = used < limit
        let message = undefined

        if (!canAdd) {
          message = `Limite de ${limit} usuários atingido para o plano ${tenant.plan.charAt(0).toUpperCase() + tenant.plan.slice(1)}. Faça upgrade para adicionar mais usuários.`
        }

        return { canAdd, used, limit, message }
      },

      addUser: (userData) => {
        const { tenant } = useAuthStore.getState()
        if (!tenant) {
          toast.error('Erro: Tenant não encontrado')
          return null
        }

        // Check quota before adding
        const quotaCheck = get().checkUserQuota()
        if (!quotaCheck.canAdd) {
          toast.error(quotaCheck.message || 'Limite de usuários atingido', {
            description: 'Acesse a página de Billing para fazer upgrade do seu plano.',
            duration: 6000,
            action: {
              label: 'Fazer Upgrade',
              onClick: () => window.open('/billing', '_blank')
            }
          })
          return null
        }

        const newUser: User = {
          ...userData,
          id: `user_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`,
          tenantId: tenant.id,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        }

        set((state) => ({
          users: [...state.users, newUser]
        }))

        toast.success('Usuário criado com sucesso!')
        return newUser.id
      },

      updateUser: (id, updates) => {
        set((state) => ({
          users: state.users.map(user =>
            user.id === id
              ? { ...user, ...updates, updatedAt: new Date().toISOString() }
              : user
          )
        }))
        toast.success('Usuário atualizado com sucesso!')
      },

      deleteUser: (id) => {
        const userToDelete = get().getUserById(id)
        if (!userToDelete) {
          toast.error('Usuário não encontrado')
          return
        }

        // Prevent deleting the current user
        const { user: currentUser } = useAuthStore.getState()
        if (currentUser?.id === id) {
          toast.error('Você não pode excluir seu próprio usuário')
          return
        }

        // Prevent deleting the last admin
        const currentTenantUsers = get().getUsersForCurrentTenant()
        const activeAdmins = currentTenantUsers.filter(u => u.role === 'admin' && u.isActive && u.id !== id)
        
        if (userToDelete.role === 'admin' && activeAdmins.length === 0) {
          toast.error('Não é possível excluir o último administrador ativo')
          return
        }

        set((state) => ({
          users: state.users.filter(user => user.id !== id)
        }))
        
        toast.success('Usuário excluído com sucesso!')
      },

      toggleUserStatus: (id) => {
        const userToToggle = get().getUserById(id)
        if (!userToToggle) {
          toast.error('Usuário não encontrado')
          return
        }

        // Prevent deactivating the current user
        const { user: currentUser } = useAuthStore.getState()
        if (currentUser?.id === id && userToToggle.isActive) {
          toast.error('Você não pode desativar seu próprio usuário')
          return
        }

        // Prevent deactivating the last admin
        if (userToToggle.role === 'admin' && userToToggle.isActive) {
          const currentTenantUsers = get().getUsersForCurrentTenant()
          const activeAdmins = currentTenantUsers.filter(u => u.role === 'admin' && u.isActive && u.id !== id)
          
          if (activeAdmins.length === 0) {
            toast.error('Não é possível desativar o último administrador ativo')
            return
          }
        }

        const newStatus = !userToToggle.isActive
        get().updateUser(id, { isActive: newStatus })
        
        toast.success(`Usuário ${newStatus ? 'ativado' : 'desativado'} com sucesso!`)
      },

      getUserById: (id) => {
        return get().users.find(user => user.id === id)
      },

      getUsersByRole: (role) => {
        const currentTenantUsers = get().getUsersForCurrentTenant()
        return currentTenantUsers.filter(user => user.role === role)
      },

      getActiveUsers: () => {
        const currentTenantUsers = get().getUsersForCurrentTenant()
        return currentTenantUsers.filter(user => user.isActive)
      },

      getUsersForCurrentTenant: () => {
        const { tenant } = useAuthStore.getState()
        if (!tenant) return []
        
        return get().users.filter(user => user.tenantId === tenant.id)
      }
    }),
    {
      name: 'users-storage',
      version: 1
    }
  )
)

// Load initial users when the store is created
useUserStore.getState().loadInitialUsers()