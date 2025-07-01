'use client'

import { useState } from 'react'
import { usePathname } from 'next/navigation'
import Link from 'next/link'
import { 
  Building,
  LayoutDashboard, 
  FileText, 
  Search, 
  BarChart3, 
  MessageSquare, 
  Settings, 
  Users,
  Bot,
  ChevronLeft,
  ChevronRight,
  CreditCard,
  Bell
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useUIStore } from '@/store'
import { usePermissions } from '@/hooks/usePermissions'
import { UserRole } from '@/types'

interface NavigationItem {
  name: string
  href: string
  icon: any
  badge?: string | number
  allowedRoles?: UserRole[]
}

const navigationItems: NavigationItem[] = [
  {
    name: 'Dashboard',
    href: '/dashboard',
    icon: LayoutDashboard,
    allowedRoles: ['admin', 'manager', 'lawyer', 'assistant'],
  },
  {
    name: 'Processos',
    href: '/processes',
    icon: FileText,
    allowedRoles: ['admin', 'manager', 'lawyer', 'assistant'],
  },
  {
    name: 'Busca',
    href: '/search',
    icon: Search,
    allowedRoles: ['admin', 'manager', 'lawyer', 'assistant'],
  },
  {
    name: 'Relatórios',
    href: '/reports',
    icon: BarChart3,
    allowedRoles: ['admin', 'manager', 'lawyer'],
  },
  {
    name: 'IA Assistant',
    href: '/ai',
    icon: Bot,
    allowedRoles: ['admin', 'manager', 'lawyer', 'assistant'],
  },
  {
    name: 'Notificações',
    href: '/notifications',
    icon: Bell,
    allowedRoles: ['admin', 'manager', 'lawyer', 'assistant'],
  },
  {
    name: 'Usuários',
    href: '/users',
    icon: Users,
    allowedRoles: ['admin'],
  },
  {
    name: 'Billing',
    href: '/billing',
    icon: CreditCard,
    allowedRoles: ['admin'],
  },
  {
    name: 'Configurações',
    href: '/settings',
    icon: Settings,
    allowedRoles: ['admin', 'manager'],
  },
]

export function Sidebar() {
  const pathname = usePathname()
  const { sidebarCollapsed, toggleSidebarCollapse } = useUIStore()
  const { currentRole, getRoleLabel } = usePermissions()

  // Filtrar itens de navegação baseado na role do usuário
  const filteredNavigationItems = navigationItems.filter(item => {
    if (!item.allowedRoles || !currentRole) return false
    return item.allowedRoles.includes(currentRole)
  })

  return (
    <div
      className={cn(
        'bg-card border-r transition-all duration-300 relative',
        sidebarCollapsed ? 'w-16' : 'w-64'
      )}
    >
      {/* Logo */}
      <div className="p-4 border-b">
        <div className="flex items-center">
          <div className="w-8 h-8 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-lg flex items-center justify-center">
            <Building className="w-5 h-5 text-white" />
          </div>
          {!sidebarCollapsed && (
            <div className="ml-3">
              <span className="text-lg font-semibold block">Direito Lux</span>
              {currentRole && (
                <span className="text-xs text-muted-foreground">
                  {getRoleLabel(currentRole)}
                </span>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Navigation */}
      <nav className="p-4 space-y-2">
        {filteredNavigationItems.map((item) => {
          const isActive = pathname === item.href
          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                'flex items-center px-3 py-2 rounded-lg transition-colors group relative',
                isActive
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:text-foreground hover:bg-accent',
                sidebarCollapsed ? 'justify-center' : ''
              )}
            >
              <item.icon className="w-5 h-5" />
              {!sidebarCollapsed && (
                <div className="ml-3 flex-1 flex items-center justify-between">
                  <span className="text-sm font-medium">{item.name}</span>
                  {item.badge && (
                    <span className="bg-primary text-primary-foreground text-xs px-2 py-0.5 rounded-full">
                      {item.badge}
                    </span>
                  )}
                </div>
              )}
              {sidebarCollapsed && (
                <span className="absolute left-16 ml-6 px-2 py-1 bg-popover text-popover-foreground text-sm rounded-md opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 whitespace-nowrap z-50">
                  {item.name}
                  {item.badge && (
                    <span className="ml-2 bg-primary text-primary-foreground text-xs px-1.5 py-0.5 rounded-full">
                      {item.badge}
                    </span>
                  )}
                </span>
              )}
            </Link>
          )
        })}
      </nav>

      {/* User Role Indicator */}
      {!sidebarCollapsed && currentRole && (
        <div className="absolute bottom-4 left-4 right-4 p-3 bg-muted/50 rounded-lg">
          <div className="text-xs text-muted-foreground">Acesso como</div>
          <div className="text-sm font-medium">{getRoleLabel(currentRole)}</div>
        </div>
      )}

      {/* Collapse Toggle */}
      <button
        onClick={toggleSidebarCollapse}
        className="absolute top-20 -right-3 w-6 h-6 bg-card border rounded-full flex items-center justify-center shadow-md hover:shadow-lg transition-shadow"
      >
        {sidebarCollapsed ? (
          <ChevronRight className="w-4 h-4" />
        ) : (
          <ChevronLeft className="w-4 h-4" />
        )}
      </button>
    </div>
  )
}