import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { Theme } from '@/types'

// UI Store
interface UIState {
  theme: Theme
  sidebarOpen: boolean
  sidebarCollapsed: boolean
  breadcrumbs: Array<{ label: string; href?: string }>
  pageTitle: string
  setTheme: (theme: Theme) => void
  setSidebarOpen: (open: boolean) => void
  setSidebarCollapsed: (collapsed: boolean) => void
  setBreadcrumbs: (breadcrumbs: Array<{ label: string; href?: string }>) => void
  setPageTitle: (title: string) => void
  toggleSidebar: () => void
  toggleSidebarCollapse: () => void
}

export const useUIStore = create<UIState>()(
  persist(
    (set, get) => ({
      theme: 'system',
      sidebarOpen: true,
      sidebarCollapsed: false,
      breadcrumbs: [],
      pageTitle: '',
      setTheme: (theme) => set({ theme }),
      setSidebarOpen: (sidebarOpen) => set({ sidebarOpen }),
      setSidebarCollapsed: (sidebarCollapsed) => set({ sidebarCollapsed }),
      setBreadcrumbs: (breadcrumbs) => set({ breadcrumbs }),
      setPageTitle: (pageTitle) => set({ pageTitle }),
      toggleSidebar: () => set({ sidebarOpen: !get().sidebarOpen }),
      toggleSidebarCollapse: () => set({ sidebarCollapsed: !get().sidebarCollapsed }),
    }),
    {
      name: 'ui-storage',
      partialize: (state) => ({ 
        theme: state.theme,
        sidebarCollapsed: state.sidebarCollapsed,
      }),
    }
  )
)

// Helper hooks
export const useUI = () => useUIStore()

// Selectors
export const selectTheme = () => useUIStore(state => state.theme)
export const selectSidebarOpen = () => useUIStore(state => state.sidebarOpen)