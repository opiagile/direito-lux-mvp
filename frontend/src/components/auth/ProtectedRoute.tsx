'use client'

import { useEffect } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { Shield, AlertTriangle, Lock } from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { usePermissions } from '@/hooks/usePermissions'
import { LoadingScreen } from '@/components/ui/loading-screen'

interface ProtectedRouteProps {
  children: React.ReactNode
  fallback?: React.ReactNode
  redirectTo?: string
}

export function ProtectedRoute({ 
  children, 
  fallback, 
  redirectTo = '/dashboard' 
}: ProtectedRouteProps) {
  const router = useRouter()
  const pathname = usePathname()
  const { canAccessPage, getPagePermission, user, currentRole, getRoleLabel } = usePermissions()

  const permission = getPagePermission(pathname)

  useEffect(() => {
    // Se não tem usuário logado, redireciona para login
    if (!user) {
      router.push('/login')
      return
    }

    // Se tem redirectTo configurado e não tem acesso, redireciona
    if (!permission.hasAccess && permission.redirectTo) {
      router.push(permission.redirectTo)
      return
    }
  }, [user, permission.hasAccess, permission.redirectTo, router])

  // Loading enquanto verifica autenticação
  if (!user) {
    return <LoadingScreen />
  }

  // Se não tem acesso, mostra fallback ou página de erro padrão
  if (!permission.hasAccess) {
    if (fallback) {
      return <>{fallback}</>
    }

    return (
      <div className="flex items-center justify-center min-h-[60vh] p-6">
        <Card className="max-w-md w-full">
          <CardContent className="p-8 text-center">
            <div className="flex justify-center mb-6">
              <div className="relative">
                <Shield className="w-16 h-16 text-muted-foreground" />
                <div className="absolute -bottom-1 -right-1 bg-red-500 rounded-full p-1">
                  <Lock className="w-4 h-4 text-white" />
                </div>
              </div>
            </div>
            
            <h2 className="text-2xl font-bold mb-2">Acesso Restrito</h2>
            
            <p className="text-muted-foreground mb-6">
              {permission.message}
            </p>

            <div className="bg-muted/50 rounded-lg p-4 mb-6">
              <div className="flex items-center justify-center space-x-2 text-sm">
                <span className="text-muted-foreground">Seu nível de acesso:</span>
                <span className="font-medium">{getRoleLabel(currentRole)}</span>
              </div>
              
              {permission.allowedRoles.length > 0 && (
                <div className="mt-2 text-xs text-muted-foreground">
                  <span>Acesso permitido para: </span>
                  <span className="font-medium">
                    {permission.allowedRoles.map(role => getRoleLabel(role)).join(', ')}
                  </span>
                </div>
              )}
            </div>

            <div className="space-y-3">
              <Button 
                onClick={() => router.push(redirectTo)}
                className="w-full"
              >
                Voltar ao Dashboard
              </Button>
              
              <Button 
                variant="outline"
                onClick={() => router.back()}
                className="w-full"
              >
                Página Anterior
              </Button>
            </div>

            {currentRole === 'assistant' && (
              <div className="mt-6 p-4 bg-blue-50 rounded-lg border border-blue-200">
                <div className="flex items-start space-x-2">
                  <AlertTriangle className="w-4 h-4 text-blue-600 mt-0.5" />
                  <div className="text-sm text-blue-800">
                    <p className="font-medium">Precisa de mais acesso?</p>
                    <p>Entre em contato com o administrador do seu escritório para solicitar permissões adicionais.</p>
                  </div>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    )
  }

  // Se tem acesso, renderiza o conteúdo
  return <>{children}</>
}

// Hook para usar proteção de rota de forma mais simples
export function useRouteProtection(requiredRoles?: string[]) {
  const { canAccessPage, getPagePermission, currentRole } = usePermissions()
  const pathname = usePathname()

  const permission = getPagePermission(pathname)
  
  // Se roles específicas foram fornecidas, verifica contra elas
  if (requiredRoles && currentRole) {
    const hasRequiredRole = requiredRoles.includes(currentRole)
    return {
      hasAccess: hasRequiredRole,
      message: hasRequiredRole ? '' : `Acesso restrito para: ${requiredRoles.join(', ')}`,
      currentRole,
      allowedRoles: requiredRoles
    }
  }

  return permission
}

// Componente para proteger ações específicas
interface ProtectedActionProps {
  resource: string
  action: string
  children: React.ReactNode
  fallback?: React.ReactNode
  hideWhenNoAccess?: boolean
}

export function ProtectedAction({ 
  resource, 
  action, 
  children, 
  fallback, 
  hideWhenNoAccess = false 
}: ProtectedActionProps) {
  const { canPerformAction } = usePermissions()
  
  const hasPermission = canPerformAction(resource, action)
  
  if (!hasPermission) {
    if (hideWhenNoAccess) return null
    if (fallback) return <>{fallback}</>
    return null
  }
  
  return <>{children}</>
}

// Componente para mostrar conteúdo baseado em role
interface RoleBasedContentProps {
  allowedRoles: string[]
  children: React.ReactNode
  fallback?: React.ReactNode
  hideWhenNoAccess?: boolean
}

export function RoleBasedContent({ 
  allowedRoles, 
  children, 
  fallback, 
  hideWhenNoAccess = false 
}: RoleBasedContentProps) {
  const { currentRole } = usePermissions()
  
  const hasAccess = currentRole && allowedRoles.includes(currentRole)
  
  if (!hasAccess) {
    if (hideWhenNoAccess) return null
    if (fallback) return <>{fallback}</>
    return null
  }
  
  return <>{children}</>
}