import { useAuthStore } from '@/store'
import { UserRole } from '@/types'

interface PermissionConfig {
  allowedRoles: UserRole[]
  message?: string
  redirectTo?: string
}

interface PagePermissions {
  [key: string]: PermissionConfig
}

// Definição de permissões por página
const pagePermissions: PagePermissions = {
  // Gestão de Usuários - Apenas admins
  '/users': {
    allowedRoles: ['admin'],
    message: 'Apenas administradores podem acessar a gestão de usuários.'
  },
  
  // Configurações - Admin e Manager
  '/settings': {
    allowedRoles: ['admin', 'manager'],
    message: 'Apenas administradores e gerentes podem acessar as configurações.'
  },
  
  // Billing - Apenas admins
  '/billing': {
    allowedRoles: ['admin'],
    message: 'Apenas administradores podem acessar informações de billing.'
  },
  
  // Relatórios - Admin, Manager e Lawyer
  '/reports': {
    allowedRoles: ['admin', 'manager', 'lawyer'],
    message: 'Você não tem permissão para acessar os relatórios.'
  },
  
  // Busca - Todos podem acessar
  '/search': {
    allowedRoles: ['admin', 'manager', 'lawyer', 'assistant'],
    message: 'Você não tem permissão para acessar a busca.'
  },
  
  // Notificações - Todos podem acessar
  '/notifications': {
    allowedRoles: ['admin', 'manager', 'lawyer', 'assistant'],
    message: 'Você não tem permissão para acessar as notificações.'
  },
  
  // Dashboard - Todos podem acessar
  '/dashboard': {
    allowedRoles: ['admin', 'manager', 'lawyer', 'assistant'],
    message: 'Você não tem permissão para acessar o dashboard.'
  },
  
  // Processos - Todos podem acessar (com diferentes níveis)
  '/processes': {
    allowedRoles: ['admin', 'manager', 'lawyer', 'assistant'],
    message: 'Você não tem permissão para acessar os processos.'
  },
  
  // IA Assistant - Todos podem acessar
  '/ai': {
    allowedRoles: ['admin', 'manager', 'lawyer', 'assistant'],
    message: 'Você não tem permissão para acessar o AI Assistant.'
  },
  
  // Perfil - Todos podem acessar seu próprio perfil
  '/profile': {
    allowedRoles: ['admin', 'manager', 'lawyer', 'assistant'],
    message: 'Você não tem permissão para acessar o perfil.'
  }
}

// Definição de ações específicas
interface ActionPermissions {
  [key: string]: {
    [action: string]: UserRole[]
  }
}

const actionPermissions: ActionPermissions = {
  users: {
    create: ['admin'],
    edit: ['admin'],
    delete: ['admin'],
    view: ['admin']
  },
  
  processes: {
    create: ['admin', 'manager', 'lawyer'],
    edit: ['admin', 'manager', 'lawyer'],
    delete: ['admin', 'manager'],
    view: ['admin', 'manager', 'lawyer', 'assistant']
  },
  
  reports: {
    create: ['admin', 'manager'],
    schedule: ['admin', 'manager'],
    download: ['admin', 'manager', 'lawyer'],
    view: ['admin', 'manager', 'lawyer']
  },
  
  settings: {
    edit: ['admin'],
    view: ['admin', 'manager']
  },
  
  billing: {
    edit: ['admin'],
    view: ['admin']
  },
  
  notifications: {
    markAsRead: ['admin', 'manager', 'lawyer', 'assistant'],
    delete: ['admin', 'manager', 'lawyer', 'assistant'],
    configure: ['admin', 'manager']
  },
  
  profile: {
    view: ['admin', 'manager', 'lawyer', 'assistant'],
    edit: ['admin', 'manager', 'lawyer', 'assistant'],
    changePassword: ['admin', 'manager', 'lawyer', 'assistant'],
    manageSessions: ['admin', 'manager', 'lawyer', 'assistant']
  }
}

export function usePermissions() {
  const { user } = useAuthStore()
  const currentRole = user?.role

  /**
   * Verifica se o usuário tem permissão para acessar uma página
   */
  const canAccessPage = (page: string): boolean => {
    if (!currentRole) return false
    
    const permission = pagePermissions[page]
    if (!permission) return true // Se não está definido, permite acesso
    
    return permission.allowedRoles.includes(currentRole)
  }

  /**
   * Retorna informações sobre permissão de página
   */
  const getPagePermission = (page: string) => {
    const permission = pagePermissions[page]
    const hasAccess = canAccessPage(page)
    
    return {
      hasAccess,
      message: permission?.message || 'Você não tem permissão para acessar esta página.',
      allowedRoles: permission?.allowedRoles || [],
      redirectTo: permission?.redirectTo
    }
  }

  /**
   * Verifica se o usuário pode executar uma ação específica
   */
  const canPerformAction = (resource: string, action: string): boolean => {
    if (!currentRole) return false
    
    const permissions = actionPermissions[resource]
    if (!permissions) return false
    
    const actionRoles = permissions[action]
    if (!actionRoles) return false
    
    return actionRoles.includes(currentRole)
  }

  /**
   * Verifica múltiplas permissões de uma vez
   */
  const hasAnyPermission = (checks: Array<{resource: string, action: string}>): boolean => {
    return checks.some(check => canPerformAction(check.resource, check.action))
  }

  /**
   * Verifica todas as permissões de uma lista
   */
  const hasAllPermissions = (checks: Array<{resource: string, action: string}>): boolean => {
    return checks.every(check => canPerformAction(check.resource, check.action))
  }

  /**
   * Retorna o nível de acesso do usuário
   */
  const getAccessLevel = (): 'full' | 'limited' | 'read-only' | 'none' => {
    if (!currentRole) return 'none'
    
    switch (currentRole) {
      case 'admin': return 'full'
      case 'manager': return 'limited'
      case 'lawyer': return 'limited'
      case 'assistant': return 'read-only'
      default: return 'none'
    }
  }

  /**
   * Verifica se é admin
   */
  const isAdmin = (): boolean => currentRole === 'admin'

  /**
   * Verifica se é manager ou superior
   */
  const isManagerOrAbove = (): boolean => ['admin', 'manager'].includes(currentRole || '')

  /**
   * Verifica se é lawyer ou superior
   */
  const isLawyerOrAbove = (): boolean => ['admin', 'manager', 'lawyer'].includes(currentRole || '')

  /**
   * Retorna rótulo amigável para role
   */
  const getRoleLabel = (role?: UserRole): string => {
    switch (role || currentRole) {
      case 'admin': return 'Administrador'
      case 'manager': return 'Gerente'
      case 'lawyer': return 'Advogado'
      case 'assistant': return 'Assistente'
      default: return 'Usuário'
    }
  }

  /**
   * Filtra itens de menu baseado em permissões
   */
  const filterMenuItems = (items: Array<{href: string, [key: string]: any}>) => {
    return items.filter(item => canAccessPage(item.href))
  }

  return {
    // Estado atual
    currentRole,
    user,
    
    // Verificações de página
    canAccessPage,
    getPagePermission,
    
    // Verificações de ação
    canPerformAction,
    hasAnyPermission,
    hasAllPermissions,
    
    // Helpers de role
    isAdmin,
    isManagerOrAbove,
    isLawyerOrAbove,
    getAccessLevel,
    getRoleLabel,
    
    // Utilitários
    filterMenuItems
  }
}