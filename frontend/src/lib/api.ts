import axios, { AxiosError, AxiosResponse } from 'axios'
import { APIResponse, APIError } from '@/types'

// API Configuration
const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8081'
const PROCESS_SERVICE_URL = process.env.PROCESS_SERVICE_URL || 'http://127.0.0.1:8083'
const AI_SERVICE_URL = process.env.AI_SERVICE_URL || 'http://localhost:8000'
const SEARCH_SERVICE_URL = process.env.SEARCH_SERVICE_URL || 'http://localhost:8086'
const REPORT_SERVICE_URL = process.env.REPORT_SERVICE_URL || 'http://localhost:8087'

// Create axios instances for different services
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const aiClient = axios.create({
  baseURL: AI_SERVICE_URL,
  timeout: 60000,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const searchClient = axios.create({
  baseURL: SEARCH_SERVICE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const processClient = axios.create({
  baseURL: PROCESS_SERVICE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const reportClient = axios.create({
  baseURL: REPORT_SERVICE_URL,
  timeout: 60000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token and tenant ID
const addAuthInterceptor = (client: typeof apiClient) => {
  client.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('auth_token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      
      // Add tenant ID from auth store if available
      try {
        const authData = localStorage.getItem('auth-storage')
        if (authData) {
          const parsed = JSON.parse(authData)
          if (parsed.state?.tenant?.id) {
            config.headers['X-Tenant-ID'] = parsed.state.tenant.id
          }
        }
      } catch (error) {
        console.warn('Could not extract tenant ID from auth storage:', error)
      }
      
      return config
    },
    (error) => {
      return Promise.reject(error)
    }
  )
}

// Response interceptor for error handling
const addResponseInterceptor = (client: typeof apiClient) => {
  client.interceptors.response.use(
    (response: AxiosResponse<APIResponse>) => {
      return response
    },
    (error: AxiosError<APIError>) => {
      if (error.response?.status === 401) {
        // Handle unauthorized access
        localStorage.removeItem('auth_token')
        window.location.href = '/login'
      }
      
      return Promise.reject(error)
    }
  )
}

// Apply interceptors to all clients
;[apiClient, processClient, aiClient, searchClient, reportClient].forEach((client) => {
  addAuthInterceptor(client)
  addResponseInterceptor(client)
})

// API endpoints
export const endpoints = {
  // Auth endpoints
  auth: {
    login: '/api/v1/auth/login',
    logout: '/api/v1/auth/logout',
    refresh: '/api/v1/auth/refresh',
    validate: '/api/v1/auth/validate',
    register: '/api/v1/auth/register',
  },
  
  // User endpoints
  users: {
    profile: '/api/v1/users/profile',
    list: '/api/v1/users',
    create: '/api/v1/users',
    update: (id: string) => `/api/v1/users/${id}`,
    delete: (id: string) => `/api/v1/users/${id}`,
  },
  
  // Tenant endpoints
  tenants: {
    current: '/api/v1/tenants/current',
    list: '/api/v1/tenants',
    create: '/api/v1/tenants',
    update: (id: string) => `/api/v1/tenants/${id}`,
    subscription: '/api/v1/tenants/subscription',
    quotas: '/api/v1/tenants/quotas',
  },
  
  // Process endpoints
  processes: {
    list: '/api/v1/processes',
    create: '/api/v1/processes',
    get: (id: string) => `/api/v1/processes/${id}`,
    update: (id: string) => `/api/v1/processes/${id}`,
    delete: (id: string) => `/api/v1/processes/${id}`,
    movements: (id: string) => `/api/v1/processes/${id}/movements`,
    monitor: (id: string) => `/api/v1/processes/${id}/monitor`,
    unmonitor: (id: string) => `/api/v1/processes/${id}/unmonitor`,
    stats: '/api/v1/processes/stats',
  },
  
  // DataJud endpoints
  datajud: {
    search: '/api/v1/datajud/search',
    process: (number: string) => `/api/v1/datajud/process/${number}`,
    movements: (number: string) => `/api/v1/datajud/process/${number}/movements`,
    stats: '/api/v1/datajud/stats',
  },
  
  // Notification endpoints
  notifications: {
    list: '/api/v1/notifications',
    create: '/api/v1/notifications',
    get: (id: string) => `/api/v1/notifications/${id}`,
    markRead: (id: string) => `/api/v1/notifications/${id}/read`,
    preferences: '/api/v1/notifications/preferences',
    templates: '/api/v1/notifications/templates',
    stats: '/api/v1/notifications/stats',
  },
  
  // AI endpoints (AI Service)
  ai: {
    analyze: '/api/v1/analysis/document',
    jurisprudence: '/api/v1/jurisprudence/search',
    similarity: '/api/v1/jurisprudence/similarity',
    generate: '/api/v1/generation/document',
    history: '/api/v1/analysis/history',
    types: '/api/v1/analysis/types',
  },
  
  // Search endpoints (Search Service)
  search: {
    basic: '/api/v1/search',
    advanced: '/api/v1/search/advanced',
    suggestions: '/api/v1/search/suggestions',
    aggregations: '/api/v1/search/aggregate',
    index: '/api/v1/index',
  },
  
  // Report endpoints (Report Service)
  reports: {
    list: '/api/v1/reports',
    create: '/api/v1/reports',
    get: (id: string) => `/api/v1/reports/${id}`,
    download: (id: string) => `/api/v1/reports/${id}/download`,
    delete: (id: string) => `/api/v1/reports/${id}`,
    stats: '/api/v1/reports/stats',
  },
  
  // Dashboard endpoints (Report Service)
  dashboards: {
    list: '/api/v1/dashboards',
    create: '/api/v1/dashboards',
    get: (id: string) => `/api/v1/dashboards/${id}`,
    update: (id: string) => `/api/v1/dashboards/${id}`,
    delete: (id: string) => `/api/v1/dashboards/${id}`,
    data: (id: string) => `/api/v1/dashboards/${id}/data`,
    widgets: (id: string) => `/api/v1/dashboards/${id}/widgets`,
    addWidget: (id: string) => `/api/v1/dashboards/${id}/widgets`,
    updateWidget: (id: string, widgetId: string) => `/api/v1/dashboards/${id}/widgets/${widgetId}`,
    removeWidget: (id: string, widgetId: string) => `/api/v1/dashboards/${id}/widgets/${widgetId}`,
  },
  
  // Schedule endpoints (Report Service)
  schedules: {
    list: '/api/v1/schedules',
    create: '/api/v1/schedules',
    get: (id: string) => `/api/v1/schedules/${id}`,
    update: (id: string) => `/api/v1/schedules/${id}`,
    delete: (id: string) => `/api/v1/schedules/${id}`,
  },
  
  // KPI endpoints (Report Service)
  kpis: {
    list: '/api/v1/kpis',
    calculate: '/api/v1/kpis/calculate',
  },
  
  // MCP endpoints
  mcp: {
    sessions: '/api/v1/mcp/sessions',
    messages: (sessionId: string) => `/api/v1/mcp/sessions/${sessionId}/messages`,
    tools: '/api/v1/mcp/tools',
    execute: '/api/v1/mcp/execute',
    stats: '/api/v1/mcp/stats',
  },
}

// Generic API methods
export const api = {
  // Generic GET request
  get: async <T = any>(url: string, client = apiClient): Promise<T> => {
    const response = await client.get<APIResponse<T>>(url)
    return response.data.data
  },
  
  // Generic POST request
  post: async <T = any>(url: string, data?: any, client = apiClient): Promise<T> => {
    const response = await client.post<APIResponse<T>>(url, data)
    return response.data.data
  },
  
  // Generic PUT request
  put: async <T = any>(url: string, data?: any, client = apiClient): Promise<T> => {
    const response = await client.put<APIResponse<T>>(url, data)
    return response.data.data
  },
  
  // Generic PATCH request
  patch: async <T = any>(url: string, data?: any, client = apiClient): Promise<T> => {
    const response = await client.patch<APIResponse<T>>(url, data)
    return response.data.data
  },
  
  // Generic DELETE request
  delete: async <T = any>(url: string, client = apiClient): Promise<T> => {
    const response = await client.delete<APIResponse<T>>(url)
    return response.data.data
  },
  
  // File upload
  upload: async <T = any>(url: string, file: File, client = apiClient): Promise<T> => {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await client.post<APIResponse<T>>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    
    return response.data.data
  },
  
  // Download file
  download: async (url: string, filename?: string, client = apiClient): Promise<void> => {
    const response = await client.get(url, {
      responseType: 'blob',
    })
    
    const blob = new Blob([response.data])
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = filename || 'download'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)
  },
}

// Specific API methods for services
export const authAPI = {
  login: async (credentials: { email: string; password: string }) => {
    console.log(`🔐 authAPI.login chamado com email: ${credentials.email}`)
    
    try {
      // Make login request - auth service will find tenant from user email
      console.log(`📤 Fazendo POST para: ${endpoints.auth.login}`)
      const response = await apiClient.post(endpoints.auth.login, credentials)
      console.log(`📨 Resposta recebida:`, response.status, response.data)
      return response.data // Return data directly since auth doesn't wrap in APIResponse
    } catch (error: any) {
      console.error(`❌ Erro no authAPI.login:`, error)
      console.error(`❌ Response status:`, error.response?.status)
      console.error(`❌ Response data:`, error.response?.data)
      throw error
    }
  },
    
  logout: () =>
    api.post(endpoints.auth.logout),
    
  refresh: (refreshToken: string) =>
    api.post(endpoints.auth.refresh, { refreshToken }),
    
  validate: () =>
    api.get(endpoints.auth.validate),
}

export const processAPI = {
  list: (params?: any) =>
    api.get(endpoints.processes.list + (params ? `?${new URLSearchParams(params)}` : ''), processClient),
    
  get: (id: string) =>
    api.get(endpoints.processes.get(id), processClient),
    
  create: (data: any) =>
    api.post(endpoints.processes.create, data, processClient),
    
  update: (id: string, data: any) =>
    api.put(endpoints.processes.update(id), data, processClient),
    
  delete: (id: string) =>
    api.delete(endpoints.processes.delete(id), processClient),
    
  getMovements: (id: string) =>
    api.get(endpoints.processes.movements(id), processClient),
    
  monitor: (id: string) =>
    api.post(endpoints.processes.monitor(id), processClient),
    
  unmonitor: (id: string) =>
    api.delete(endpoints.processes.unmonitor(id), processClient),
    
  getStats: () =>
    api.get(endpoints.processes.stats, processClient),
}

export const reportAPI = {
  list: (params?: any) =>
    api.get(endpoints.reports.list + (params ? `?${new URLSearchParams(params)}` : ''), reportClient),
    
  create: (data: any) =>
    api.post(endpoints.reports.create, data, reportClient),
    
  get: (id: string) =>
    api.get(endpoints.reports.get(id), reportClient),
    
  download: (id: string, filename?: string) =>
    api.download(endpoints.reports.download(id), filename, reportClient),
    
  delete: (id: string) =>
    api.delete(endpoints.reports.delete(id), reportClient),
    
  getStats: () =>
    api.get(endpoints.reports.stats, reportClient),
}

export const dashboardAPI = {
  list: () =>
    api.get(endpoints.dashboards.list, reportClient),
    
  create: (data: any) =>
    api.post(endpoints.dashboards.create, data, reportClient),
    
  get: (id: string) =>
    api.get(endpoints.dashboards.get(id), reportClient),
    
  update: (id: string, data: any) =>
    api.put(endpoints.dashboards.update(id), data, reportClient),
    
  delete: (id: string) =>
    api.delete(endpoints.dashboards.delete(id), reportClient),
    
  getData: (id: string) =>
    api.get(endpoints.dashboards.data(id), reportClient),
    
  addWidget: (id: string, widget: any) =>
    api.post(endpoints.dashboards.addWidget(id), widget, reportClient),
    
  updateWidget: (id: string, widgetId: string, widget: any) =>
    api.put(endpoints.dashboards.updateWidget(id, widgetId), widget, reportClient),
    
  removeWidget: (id: string, widgetId: string) =>
    api.delete(endpoints.dashboards.removeWidget(id, widgetId), reportClient),
}

export const searchAPI = {
  basic: (query: any) =>
    api.post(endpoints.search.basic, query, searchClient),
    
  advanced: (query: any) =>
    api.post(endpoints.search.advanced, query, searchClient),
    
  suggestions: (q: string) =>
    api.get(endpoints.search.suggestions + `?q=${encodeURIComponent(q)}`, searchClient),
    
  aggregations: (query: any) =>
    api.post(endpoints.search.aggregations, query, searchClient),
}

export const aiAPI = {
  analyze: (data: any) =>
    api.post(endpoints.ai.analyze, data, aiClient),
    
  searchJurisprudence: (query: any) =>
    api.post(endpoints.ai.jurisprudence, query, aiClient),
    
  checkSimilarity: (data: any) =>
    api.post(endpoints.ai.similarity, data, aiClient),
    
  generateDocument: (data: any) =>
    api.post(endpoints.ai.generate, data, aiClient),
    
  getHistory: (params?: any) =>
    api.get(endpoints.ai.history + (params ? `?${new URLSearchParams(params)}` : ''), aiClient),
}

export default api