import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { User, Tenant, Theme, Notification } from '@/types'

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
      isLoading: true,
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

// Notification Store
interface NotificationState {
  notifications: Notification[]
  unreadCount: number
  addNotification: (notification: Notification) => void
  markAsRead: (id: string) => void
  markAllAsRead: () => void
  removeNotification: (id: string) => void
  clearNotifications: () => void
  setNotifications: (notifications: Notification[]) => void
}

export const useNotificationStore = create<NotificationState>((set, get) => ({
  notifications: [],
  unreadCount: 0,
  addNotification: (notification) => {
    const notifications = [notification, ...get().notifications].slice(0, 50) // Keep only last 50
    const unreadCount = notifications.filter(n => !n.readAt).length
    set({ notifications, unreadCount })
  },
  markAsRead: (id) => {
    const notifications = get().notifications.map(n => 
      n.id === id ? { ...n, readAt: new Date().toISOString() } : n
    )
    const unreadCount = notifications.filter(n => !n.readAt).length
    set({ notifications, unreadCount })
  },
  markAllAsRead: () => {
    const notifications = get().notifications.map(n => ({
      ...n,
      readAt: n.readAt || new Date().toISOString()
    }))
    set({ notifications, unreadCount: 0 })
  },
  removeNotification: (id) => {
    const notifications = get().notifications.filter(n => n.id !== id)
    const unreadCount = notifications.filter(n => !n.readAt).length
    set({ notifications, unreadCount })
  },
  clearNotifications: () => set({ notifications: [], unreadCount: 0 }),
  setNotifications: (notifications) => {
    const unreadCount = notifications.filter(n => !n.readAt).length
    set({ notifications, unreadCount })
  },
}))

// Dashboard Store
interface DashboardState {
  selectedDashboard: string | null
  dashboardFilters: Record<string, any>
  refreshInterval: number
  autoRefresh: boolean
  setSelectedDashboard: (dashboardId: string | null) => void
  setDashboardFilters: (filters: Record<string, any>) => void
  setRefreshInterval: (interval: number) => void
  setAutoRefresh: (autoRefresh: boolean) => void
  updateFilter: (key: string, value: any) => void
  clearFilters: () => void
}

export const useDashboardStore = create<DashboardState>()(
  persist(
    (set, get) => ({
      selectedDashboard: null,
      dashboardFilters: {},
      refreshInterval: 30000, // 30 seconds
      autoRefresh: true,
      setSelectedDashboard: (selectedDashboard) => set({ selectedDashboard }),
      setDashboardFilters: (dashboardFilters) => set({ dashboardFilters }),
      setRefreshInterval: (refreshInterval) => set({ refreshInterval }),
      setAutoRefresh: (autoRefresh) => set({ autoRefresh }),
      updateFilter: (key, value) => 
        set({ dashboardFilters: { ...get().dashboardFilters, [key]: value } }),
      clearFilters: () => set({ dashboardFilters: {} }),
    }),
    {
      name: 'dashboard-storage',
      partialize: (state) => ({
        refreshInterval: state.refreshInterval,
        autoRefresh: state.autoRefresh,
      }),
    }
  )
)

// Search Store
interface SearchState {
  recentSearches: string[]
  savedSearches: Array<{ id: string; name: string; query: any }>
  searchFilters: Record<string, any>
  addRecentSearch: (query: string) => void
  saveSearch: (id: string, name: string, query: any) => void
  removeSavedSearch: (id: string) => void
  setSearchFilters: (filters: Record<string, any>) => void
  updateSearchFilter: (key: string, value: any) => void
  clearSearchFilters: () => void
  clearRecentSearches: () => void
}

export const useSearchStore = create<SearchState>()(
  persist(
    (set, get) => ({
      recentSearches: [],
      savedSearches: [],
      searchFilters: {},
      addRecentSearch: (query) => {
        const recentSearches = [
          query,
          ...get().recentSearches.filter(q => q !== query)
        ].slice(0, 10) // Keep only last 10
        set({ recentSearches })
      },
      saveSearch: (id, name, query) => {
        const savedSearches = [...get().savedSearches, { id, name, query }]
        set({ savedSearches })
      },
      removeSavedSearch: (id) => {
        const savedSearches = get().savedSearches.filter(s => s.id !== id)
        set({ savedSearches })
      },
      setSearchFilters: (searchFilters) => set({ searchFilters }),
      updateSearchFilter: (key, value) => 
        set({ searchFilters: { ...get().searchFilters, [key]: value } }),
      clearSearchFilters: () => set({ searchFilters: {} }),
      clearRecentSearches: () => set({ recentSearches: [] }),
    }),
    {
      name: 'search-storage',
    }
  )
)

// Process Store
interface ProcessState {
  selectedProcesses: string[]
  processFilters: Record<string, any>
  viewMode: 'grid' | 'list' | 'table'
  sortBy: string
  sortOrder: 'asc' | 'desc'
  setSelectedProcesses: (processes: string[]) => void
  addSelectedProcess: (processId: string) => void
  removeSelectedProcess: (processId: string) => void
  clearSelectedProcesses: () => void
  setProcessFilters: (filters: Record<string, any>) => void
  updateProcessFilter: (key: string, value: any) => void
  clearProcessFilters: () => void
  setViewMode: (mode: 'grid' | 'list' | 'table') => void
  setSortBy: (sortBy: string) => void
  setSortOrder: (order: 'asc' | 'desc') => void
}

export const useProcessStore = create<ProcessState>()(
  persist(
    (set, get) => ({
      selectedProcesses: [],
      processFilters: {},
      viewMode: 'table',
      sortBy: 'updatedAt',
      sortOrder: 'desc',
      setSelectedProcesses: (selectedProcesses) => set({ selectedProcesses }),
      addSelectedProcess: (processId) => {
        const selectedProcesses = [...get().selectedProcesses, processId]
        set({ selectedProcesses })
      },
      removeSelectedProcess: (processId) => {
        const selectedProcesses = get().selectedProcesses.filter(id => id !== processId)
        set({ selectedProcesses })
      },
      clearSelectedProcesses: () => set({ selectedProcesses: [] }),
      setProcessFilters: (processFilters) => set({ processFilters }),
      updateProcessFilter: (key, value) => 
        set({ processFilters: { ...get().processFilters, [key]: value } }),
      clearProcessFilters: () => set({ processFilters: {} }),
      setViewMode: (viewMode) => set({ viewMode }),
      setSortBy: (sortBy) => set({ sortBy }),
      setSortOrder: (sortOrder) => set({ sortOrder }),
    }),
    {
      name: 'process-storage',
      partialize: (state) => ({
        viewMode: state.viewMode,
        sortBy: state.sortBy,
        sortOrder: state.sortOrder,
      }),
    }
  )
)

// Settings Store
interface SettingsState {
  preferences: {
    language: string
    timezone: string
    dateFormat: string
    currency: string
    pageSize: number
    emailNotifications: boolean
    pushNotifications: boolean
    whatsappNotifications: boolean
    autoSave: boolean
    compactMode: boolean
  }
  updatePreference: <K extends keyof SettingsState['preferences']>(
    key: K, 
    value: SettingsState['preferences'][K]
  ) => void
  resetPreferences: () => void
}

const defaultPreferences = {
  language: 'pt-BR',
  timezone: 'America/Sao_Paulo',
  dateFormat: 'dd/MM/yyyy',
  currency: 'BRL',
  pageSize: 20,
  emailNotifications: true,
  pushNotifications: true,
  whatsappNotifications: true,
  autoSave: true,
  compactMode: false,
}

export const useSettingsStore = create<SettingsState>()(
  persist(
    (set, get) => ({
      preferences: defaultPreferences,
      updatePreference: (key, value) => 
        set({ preferences: { ...get().preferences, [key]: value } }),
      resetPreferences: () => set({ preferences: defaultPreferences }),
    }),
    {
      name: 'settings-storage',
    }
  )
)

// Global store combining all stores
export const useStore = () => ({
  auth: useAuthStore(),
  ui: useUIStore(),
  notifications: useNotificationStore(),
  dashboard: useDashboardStore(),
  search: useSearchStore(),
  process: useProcessStore(),
  settings: useSettingsStore(),
})

// Helper hooks
export const useAuth = () => useAuthStore()
export const useUI = () => useUIStore()
export const useNotifications = () => useNotificationStore()
export const useDashboard = () => useDashboardStore()
export const useSearch = () => useSearchStore()
export const useProcess = () => useProcessStore()
export const useSettings = () => useSettingsStore()

// Selectors
export const selectUser = () => useAuthStore(state => state.user)
export const selectTenant = () => useAuthStore(state => state.tenant)
export const selectIsAuthenticated = () => useAuthStore(state => state.isAuthenticated)
export const selectTheme = () => useUIStore(state => state.theme)
export const selectSidebarOpen = () => useUIStore(state => state.sidebarOpen)
export const selectUnreadCount = () => useNotificationStore(state => state.unreadCount)