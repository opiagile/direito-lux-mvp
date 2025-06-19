// User and Authentication Types
export interface User {
  id: string
  email: string
  name: string
  role: UserRole
  tenantId: string
  avatar?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export type UserRole = 'admin' | 'manager' | 'lawyer' | 'assistant'

export interface AuthSession {
  user: User
  token: string
  refreshToken: string
  expiresAt: string
}

// Tenant and Subscription Types
export interface Tenant {
  id: string
  name: string
  cnpj: string
  email: string
  phone?: string
  address?: Address
  plan: SubscriptionPlan
  subscription: Subscription
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface Address {
  street: string
  number: string
  complement?: string
  neighborhood: string
  city: string
  state: string
  zipCode: string
}

export type SubscriptionPlan = 'starter' | 'professional' | 'business' | 'enterprise'

export interface Subscription {
  id: string
  tenantId: string
  plan: SubscriptionPlan
  status: SubscriptionStatus
  startDate: string
  endDate?: string
  trial: boolean
  trialEndsAt?: string
  quotas: PlanQuotas
}

export type SubscriptionStatus = 'active' | 'trial' | 'suspended' | 'cancelled'

export interface PlanQuotas {
  processes: number
  users: number
  mcpCommands: number
  aiSummaries: number
  reports: number
  dashboards: number
  widgetsPerDashboard: number
  schedules: number
}

// Process Types
export interface Process {
  id: string
  tenantId: string
  number: string
  type: string
  subject: string
  court: string
  status: ProcessStatus
  parties: Party[]
  movements: Movement[]
  monitoring: boolean
  tags: string[]
  priority: ProcessPriority
  lawyer?: string
  estimatedValue?: number
  createdAt: string
  updatedAt: string
  lastMovement?: string
}

export type ProcessStatus = 'active' | 'suspended' | 'archived' | 'concluded'
export type ProcessPriority = 'low' | 'medium' | 'high' | 'urgent'

export interface Party {
  id: string
  name: string
  document: string
  type: PartyType
  role: PartyRole
}

export type PartyType = 'person' | 'company'
export type PartyRole = 'plaintiff' | 'defendant' | 'attorney' | 'witness' | 'other'

export interface Movement {
  id: string
  processId: string
  date: string
  description: string
  type: string
  isImportant: boolean
  aiSummary?: string
  documents: Document[]
}

export interface Document {
  id: string
  name: string
  type: string
  url: string
  size: number
  uploadedAt: string
}

// Notification Types
export interface Notification {
  id: string
  tenantId: string
  userId?: string
  type: NotificationType
  channel: NotificationChannel[]
  title: string
  message: string
  data?: Record<string, any>
  status: NotificationStatus
  priority: NotificationPriority
  scheduledAt?: string
  sentAt?: string
  readAt?: string
  createdAt: string
}

export type NotificationType = 'process_update' | 'deadline' | 'system' | 'marketing'
export type NotificationChannel = 'email' | 'whatsapp' | 'telegram' | 'push' | 'sms'
export type NotificationStatus = 'pending' | 'sent' | 'delivered' | 'failed' | 'cancelled'
export type NotificationPriority = 'low' | 'normal' | 'high' | 'critical'

// Report Types
export interface Report {
  id: string
  tenantId: string
  type: ReportType
  title: string
  format: ReportFormat
  status: ReportStatus
  filters: Record<string, any>
  fileUrl?: string
  fileSize?: number
  processingTime?: number
  scheduledBy?: string
  createdAt: string
  completedAt?: string
  expiresAt?: string
}

export type ReportType = 
  | 'executive_summary' 
  | 'process_analysis' 
  | 'productivity' 
  | 'financial' 
  | 'legal_timeline' 
  | 'jurisprudence_analysis'

export type ReportFormat = 'pdf' | 'excel' | 'csv' | 'html'
export type ReportStatus = 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled'

export interface ReportSchedule {
  id: string
  tenantId: string
  reportType: ReportType
  title: string
  format: ReportFormat
  frequency: ReportFrequency
  cronExpression?: string
  filters: Record<string, any>
  recipients: string[]
  isActive: boolean
  lastRunAt?: string
  nextRunAt?: string
  createdAt: string
}

export type ReportFrequency = 'once' | 'daily' | 'weekly' | 'monthly' | 'custom'

// Dashboard Types
export interface Dashboard {
  id: string
  tenantId: string
  title: string
  description?: string
  isPublic: boolean
  isDefault: boolean
  layout: DashboardLayout
  widgets: DashboardWidget[]
  sharedWith: string[]
  createdBy: string
  createdAt: string
  updatedAt: string
}

export interface DashboardLayout {
  columns: number
  gap: number
  padding: number
}

export interface DashboardWidget {
  id: string
  dashboardId: string
  type: WidgetType
  title: string
  dataSource: string
  chartType?: ChartType
  position: WidgetPosition
  size: WidgetSize
  config: Record<string, any>
  data?: any
}

export type WidgetType = 'kpi' | 'chart' | 'table' | 'counter' | 'gauge' | 'timeline'
export type ChartType = 'line' | 'bar' | 'pie' | 'area' | 'scatter' | 'radar'

export interface WidgetPosition {
  x: number
  y: number
}

export interface WidgetSize {
  width: number
  height: number
}

export interface KPI {
  id: string
  name: string
  displayName: string
  currentValue: number
  previousValue?: number
  target?: number
  unit: string
  trend: KPITrend
  trendPercentage?: number
  category: string
  description?: string
  updatedAt: string
}

export type KPITrend = 'up' | 'down' | 'stable'

// Search Types
export interface SearchQuery {
  query: string
  filters?: SearchFilters
  sort?: SearchSort
  pagination?: SearchPagination
}

export interface SearchFilters {
  type?: string[]
  court?: string[]
  status?: string[]
  dateRange?: {
    start: string
    end: string
  }
  tags?: string[]
}

export interface SearchSort {
  field: string
  direction: 'asc' | 'desc'
}

export interface SearchPagination {
  page: number
  limit: number
}

export interface SearchResult<T = any> {
  items: T[]
  total: number
  page: number
  limit: number
  hasMore: boolean
}

// AI Types
export interface AIAnalysis {
  id: string
  type: AIAnalysisType
  input: string
  result: Record<string, any>
  confidence: number
  processingTime: number
  model: string
  createdAt: string
}

export type AIAnalysisType = 
  | 'document_summary' 
  | 'case_similarity' 
  | 'risk_assessment' 
  | 'keyword_extraction'
  | 'sentiment_analysis'

export interface JurisprudenceSearch {
  query: string
  court?: string
  similarity?: number
  results: JurisprudenceResult[]
}

export interface JurisprudenceResult {
  id: string
  title: string
  court: string
  date: string
  summary: string
  similarity: number
  url?: string
}

// MCP Types
export interface MCPSession {
  id: string
  userId: string
  tenantId: string
  platform: MCPPlatform
  isActive: boolean
  lastActivity: string
  messageCount: number
  createdAt: string
}

export type MCPPlatform = 'whatsapp' | 'telegram' | 'claude' | 'slack'

export interface MCPMessage {
  id: string
  sessionId: string
  type: MCPMessageType
  content: string
  metadata?: Record<string, any>
  timestamp: string
}

export type MCPMessageType = 'user' | 'assistant' | 'system' | 'tool_call' | 'tool_result'

export interface MCPTool {
  name: string
  description: string
  parameters: Record<string, any>
  quota: number
  category: string
}

// API Response Types
export interface APIResponse<T = any> {
  data: T
  message?: string
  success: boolean
  timestamp: string
}

export interface APIError {
  error: string
  message: string
  code?: string
  details?: Record<string, any>
  timestamp: string
}

export interface PaginatedResponse<T = any> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    totalPages: number
    hasMore: boolean
  }
}

// Form Types
export interface FormField {
  name: string
  label: string
  type: FormFieldType
  required?: boolean
  placeholder?: string
  options?: FormOption[]
  validation?: Record<string, any>
}

export type FormFieldType = 
  | 'text' 
  | 'email' 
  | 'password' 
  | 'number' 
  | 'select' 
  | 'multiselect'
  | 'textarea' 
  | 'date' 
  | 'checkbox' 
  | 'radio'
  | 'file'

export interface FormOption {
  value: string
  label: string
  disabled?: boolean
}

// Navigation Types
export interface NavItem {
  id: string
  label: string
  icon?: string
  href?: string
  children?: NavItem[]
  badge?: string | number
  disabled?: boolean
  permissions?: string[]
}

// Theme Types
export type Theme = 'light' | 'dark' | 'system'

// Component Props Types
export interface BaseComponentProps {
  className?: string
  children?: React.ReactNode
}

export interface ButtonProps extends BaseComponentProps {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'destructive'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  onClick?: () => void
  type?: 'button' | 'submit' | 'reset'
}

export interface InputProps extends BaseComponentProps {
  type?: string
  value?: string
  defaultValue?: string
  placeholder?: string
  disabled?: boolean
  required?: boolean
  error?: string
  onChange?: (value: string) => void
}

export interface SelectProps extends BaseComponentProps {
  value?: string
  defaultValue?: string
  placeholder?: string
  options: FormOption[]
  disabled?: boolean
  required?: boolean
  error?: string
  onChange?: (value: string) => void
}

export interface ModalProps extends BaseComponentProps {
  isOpen: boolean
  onClose: () => void
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
}

export interface ToastProps {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  title: string
  message?: string
  duration?: number
  action?: {
    label: string
    onClick: () => void
  }
}