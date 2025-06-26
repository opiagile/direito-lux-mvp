import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { User, Tenant } from '@/types'

// Auth Store
interface AuthState {
  user: User | null
  tenant: Tenant | null
  isAuthenticated: boolean
  isLoading: boolean
  token: string | null
  setUser: (user: User | null) => void
  setTenant: (tenant: Tenant | null) => void
  setToken: (token: string | null) => void
  setLoading: (loading: boolean) => void
  login: (user: User, tenant: Tenant, token: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      tenant: null,
      isAuthenticated: false,
      isLoading: false,
      token: null,
      setUser: (user) => set({ user, isAuthenticated: !!user }),
      setTenant: (tenant) => set({ tenant }),
      setToken: (token) => set({ token }),
      setLoading: (isLoading) => set({ isLoading }),
      login: (user, tenant, token) => 
        set({ user, tenant, token, isAuthenticated: true, isLoading: false }),
      logout: () => 
        set({ user: null, tenant: null, token: null, isAuthenticated: false, isLoading: false }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ 
        user: state.user, 
        tenant: state.tenant, 
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)

// Helper hooks
export const useAuth = () => useAuthStore()

// Selectors
export const selectUser = () => useAuthStore(state => state.user)
export const selectTenant = () => useAuthStore(state => state.tenant)
export const selectIsAuthenticated = () => useAuthStore(state => state.isAuthenticated)