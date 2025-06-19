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
  ChevronRight
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useUIStore } from '@/store'

const navigationItems = [
  {
    name: 'Dashboard',
    href: '/dashboard',
    icon: LayoutDashboard,
  },
  {
    name: 'Processos',
    href: '/processes',
    icon: FileText,
  },
  {
    name: 'Busca',
    href: '/search',
    icon: Search,
  },
  {
    name: 'Relatórios',
    href: '/reports',
    icon: BarChart3,
  },
  {
    name: 'IA Assistant',
    href: '/ai',
    icon: Bot,
  },
  {
    name: 'Notificações',
    href: '/notifications',
    icon: MessageSquare,
  },
  {
    name: 'Usuários',
    href: '/users',
    icon: Users,
  },
  {
    name: 'Configurações',
    href: '/settings',
    icon: Settings,
  },
]

export function Sidebar() {
  const pathname = usePathname()
  const { sidebarCollapsed, toggleSidebarCollapse } = useUIStore()

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
            <span className="ml-3 text-lg font-semibold">Direito Lux</span>
          )}
        </div>
      </div>

      {/* Navigation */}
      <nav className="p-4 space-y-2">
        {navigationItems.map((item) => {
          const isActive = pathname === item.href
          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                'flex items-center px-3 py-2 rounded-lg transition-colors group',
                isActive
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:text-foreground hover:bg-accent',
                sidebarCollapsed ? 'justify-center' : ''
              )}
            >
              <item.icon className="w-5 h-5" />
              {!sidebarCollapsed && (
                <span className="ml-3 text-sm font-medium">{item.name}</span>
              )}
              {sidebarCollapsed && (
                <span className="absolute left-16 ml-6 px-2 py-1 bg-popover text-popover-foreground text-sm rounded-md opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 whitespace-nowrap z-50">
                  {item.name}
                </span>
              )}
            </Link>
          )
        })}
      </nav>

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