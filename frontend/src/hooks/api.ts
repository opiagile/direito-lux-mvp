import React from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { 
  processAPI, 
  reportAPI, 
  dashboardAPI, 
  searchAPI, 
  aiAPI,
  authAPI 
} from '@/lib/api'
import { 
  Process, 
  Report, 
  Dashboard, 
  SearchQuery, 
  AIAnalysis,
  User,
  PaginatedResponse 
} from '@/types'
import { toast } from 'sonner'

// Query Keys
export const queryKeys = {
  auth: ['auth'] as const,
  user: ['user'] as const,
  processes: ['processes'] as const,
  process: (id: string) => ['process', id] as const,
  processStats: ['processes', 'stats'] as const,
  reports: ['reports'] as const,
  report: (id: string) => ['report', id] as const,
  reportStats: ['reports', 'stats'] as const,
  dashboards: ['dashboards'] as const,
  dashboard: (id: string) => ['dashboard', id] as const,
  dashboardData: (id: string) => ['dashboard', id, 'data'] as const,
  kpis: ['kpis'] as const,
  schedules: ['schedules'] as const,
  search: (query: string) => ['search', query] as const,
  searchSuggestions: (query: string) => ['search', 'suggestions', query] as const,
  aiHistory: ['ai', 'history'] as const,
  notifications: ['notifications'] as const,
}

// Auth Hooks
export const useAuth = () => {
  return useQuery({
    queryKey: queryKeys.auth,
    queryFn: authAPI.validate,
    retry: false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

export const useLogin = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: authAPI.login,
    onSuccess: (data) => {
      // Store access_token from response
      if (data.access_token) {
        localStorage.setItem('auth_token', data.access_token)
      }
      queryClient.setQueryData(queryKeys.auth, data)
      // Remove automatic success toast - component will handle it
    },
    onError: (error: any) => {
      // Log error for debugging but don't show toast here - component will handle it
      console.log('🔍 useLogin onError capturado:', error)
      console.log('🔍 Error response:', error.response)
      console.log('🔍 Error status:', error.response?.status)
      // Error is automatically propagated to component by react-query
    },
    throwOnError: true // Ensure error is thrown to component catch block
  })
}

export const useLogout = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: authAPI.logout,
    onSuccess: () => {
      localStorage.removeItem('auth_token')
      queryClient.clear()
      toast.success('Logout realizado com sucesso!')
    },
  })
}

// Process Hooks
export const useProcesses = (params?: any) => {
  return useQuery({
    queryKey: [...queryKeys.processes, params],
    queryFn: () => processAPI.list(params),
    staleTime: 30 * 1000, // 30 seconds
  })
}

export const useProcess = (id: string) => {
  return useQuery({
    queryKey: queryKeys.process(id),
    queryFn: () => processAPI.get(id),
    enabled: !!id,
    staleTime: 30 * 1000,
  })
}

export const useProcessStats = () => {
  return useQuery({
    queryKey: queryKeys.processStats,
    queryFn: processAPI.getStats,
    staleTime: 60 * 1000, // 1 minute
    retry: false, // Don't retry 404s
    meta: {
      errorHandler: 'silent' // Don't show toast errors for missing endpoints
    }
  })
}

export const useCreateProcess = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: processAPI.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.processes })
      queryClient.invalidateQueries({ queryKey: queryKeys.processStats })
      toast.success('Processo criado com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao criar processo')
    },
  })
}

export const useUpdateProcess = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      processAPI.update(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.processes })
      queryClient.invalidateQueries({ queryKey: queryKeys.process(id) })
      queryClient.invalidateQueries({ queryKey: queryKeys.processStats })
      toast.success('Processo atualizado com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao atualizar processo')
    },
  })
}

export const useDeleteProcess = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: processAPI.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.processes })
      queryClient.invalidateQueries({ queryKey: queryKeys.processStats })
      toast.success('Processo excluído com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao excluir processo')
    },
  })
}

export const useMonitorProcess = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: processAPI.monitor,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.process(id) })
      queryClient.invalidateQueries({ queryKey: queryKeys.processes })
      toast.success('Processo adicionado ao monitoramento!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao monitorar processo')
    },
  })
}

// Report Hooks
export const useReports = (params?: any) => {
  return useQuery({
    queryKey: [...queryKeys.reports, params],
    queryFn: () => reportAPI.list(params),
    staleTime: 30 * 1000,
  })
}

export const useReport = (id: string) => {
  return useQuery({
    queryKey: queryKeys.report(id),
    queryFn: () => reportAPI.get(id),
    enabled: !!id,
    staleTime: 30 * 1000,
  })
}

export const useReportStats = () => {
  return useQuery({
    queryKey: queryKeys.reportStats,
    queryFn: reportAPI.getStats,
    staleTime: 60 * 1000,
  })
}

export const useCreateReport = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: reportAPI.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.reports })
      queryClient.invalidateQueries({ queryKey: queryKeys.reportStats })
      toast.success('Relatório criado com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao criar relatório')
    },
  })
}

export const useDeleteReport = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: reportAPI.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.reports })
      queryClient.invalidateQueries({ queryKey: queryKeys.reportStats })
      toast.success('Relatório excluído com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao excluir relatório')
    },
  })
}

// Dashboard Hooks
export const useDashboards = () => {
  return useQuery({
    queryKey: queryKeys.dashboards,
    queryFn: dashboardAPI.list,
    staleTime: 60 * 1000,
  })
}

export const useDashboard = (id: string) => {
  return useQuery({
    queryKey: queryKeys.dashboard(id),
    queryFn: () => dashboardAPI.get(id),
    enabled: !!id,
    staleTime: 30 * 1000,
  })
}

export const useDashboardData = (id: string) => {
  return useQuery({
    queryKey: queryKeys.dashboardData(id),
    queryFn: () => dashboardAPI.getData(id),
    enabled: !!id,
    staleTime: 15 * 1000, // 15 seconds for real-time data
    refetchInterval: 30 * 1000, // Auto refresh every 30 seconds
  })
}

export const useCreateDashboard = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: dashboardAPI.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.dashboards })
      toast.success('Dashboard criado com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao criar dashboard')
    },
  })
}

export const useUpdateDashboard = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      dashboardAPI.update(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.dashboards })
      queryClient.invalidateQueries({ queryKey: queryKeys.dashboard(id) })
      toast.success('Dashboard atualizado com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao atualizar dashboard')
    },
  })
}

export const useAddWidget = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ dashboardId, widget }: { dashboardId: string; widget: any }) => 
      dashboardAPI.addWidget(dashboardId, widget),
    onSuccess: (_, { dashboardId }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.dashboard(dashboardId) })
      queryClient.invalidateQueries({ queryKey: queryKeys.dashboardData(dashboardId) })
      toast.success('Widget adicionado com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao adicionar widget')
    },
  })
}

export const useUpdateWidget = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ dashboardId, widgetId, widget }: { 
      dashboardId: string; 
      widgetId: string; 
      widget: any 
    }) => dashboardAPI.updateWidget(dashboardId, widgetId, widget),
    onSuccess: (_, { dashboardId }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.dashboard(dashboardId) })
      queryClient.invalidateQueries({ queryKey: queryKeys.dashboardData(dashboardId) })
      toast.success('Widget atualizado com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao atualizar widget')
    },
  })
}

export const useRemoveWidget = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ dashboardId, widgetId }: { 
      dashboardId: string; 
      widgetId: string 
    }) => dashboardAPI.removeWidget(dashboardId, widgetId),
    onSuccess: (_, { dashboardId }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.dashboard(dashboardId) })
      queryClient.invalidateQueries({ queryKey: queryKeys.dashboardData(dashboardId) })
      toast.success('Widget removido com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao remover widget')
    },
  })
}

// Search Hooks
export const useSearch = (query: SearchQuery) => {
  return useQuery({
    queryKey: queryKeys.search(JSON.stringify(query)),
    queryFn: () => searchAPI.basic(query),
    enabled: !!query.query && query.query.length > 0,
    staleTime: 60 * 1000,
  })
}

export const useAdvancedSearch = () => {
  return useMutation({
    mutationFn: searchAPI.advanced,
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro na busca avançada')
    },
  })
}

export const useSearchSuggestions = (query: string) => {
  return useQuery({
    queryKey: queryKeys.searchSuggestions(query),
    queryFn: () => searchAPI.suggestions(query),
    enabled: !!query && query.length > 2,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

// AI Hooks
export const useAIAnalysis = () => {
  return useMutation({
    mutationFn: aiAPI.analyze,
    onSuccess: () => {
      toast.success('Análise realizada com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro na análise de IA')
    },
  })
}

export const useJurisprudenceSearch = () => {
  return useMutation({
    mutationFn: aiAPI.searchJurisprudence,
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro na busca de jurisprudência')
    },
  })
}

export const useGenerateDocument = () => {
  return useMutation({
    mutationFn: aiAPI.generateDocument,
    onSuccess: () => {
      toast.success('Documento gerado com sucesso!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao gerar documento')
    },
  })
}

export const useAIHistory = (params?: any) => {
  return useQuery({
    queryKey: [...queryKeys.aiHistory, params],
    queryFn: () => aiAPI.getHistory(params),
    staleTime: 60 * 1000,
  })
}

// Custom hooks for common patterns
export const useDebounce = <T>(value: T, delay: number): T => {
  const [debouncedValue, setDebouncedValue] = React.useState<T>(value)

  React.useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value)
    }, delay)

    return () => {
      clearTimeout(handler)
    }
  }, [value, delay])

  return debouncedValue
}

export const usePagination = (initialPage = 1, initialLimit = 10) => {
  const [page, setPage] = React.useState(initialPage)
  const [limit, setLimit] = React.useState(initialLimit)

  const nextPage = () => setPage(prev => prev + 1)
  const prevPage = () => setPage(prev => Math.max(1, prev - 1))
  const goToPage = (newPage: number) => setPage(Math.max(1, newPage))
  const changeLimit = (newLimit: number) => {
    setLimit(newLimit)
    setPage(1) // Reset to first page when changing limit
  }

  return {
    page,
    limit,
    setPage,
    setLimit,
    nextPage,
    prevPage,
    goToPage,
    changeLimit,
  }
}

export const useLocalStorage = <T>(key: string, initialValue: T) => {
  const [storedValue, setStoredValue] = React.useState<T>(() => {
    try {
      const item = window.localStorage.getItem(key)
      return item ? JSON.parse(item) : initialValue
    } catch (error) {
      return initialValue
    }
  })

  const setValue = (value: T | ((val: T) => T)) => {
    try {
      const valueToStore = value instanceof Function ? value(storedValue) : value
      setStoredValue(valueToStore)
      window.localStorage.setItem(key, JSON.stringify(valueToStore))
    } catch (error) {
      console.error('Error setting localStorage:', error)
    }
  }

  return [storedValue, setValue] as const
}