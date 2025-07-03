'use client'

import { 
  FileText, 
  Users, 
  TrendingUp, 
  AlertTriangle,
  Activity,
  Clock,
  BarChart3,
  Plus
} from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { useProcessStats } from '@/hooks/api'
import { usePermissions } from '@/hooks/usePermissions'
import { formatNumber, formatPercentage } from '@/lib/utils'

// ❌ MOCK REMOVIDO - Usar dados reais do useProcessStats hook
// KPIs devem vir de: GET /api/v1/reports/dashboard

// ❌ MOCK REMOVIDO - Usar dados reais das atividades recentes
// Atividades devem vir de: GET /api/v1/reports/recent-activities

export default function DashboardPage() {
  const { data: processStats, isLoading, error } = useProcessStats()
  const { canPerformAction } = usePermissions()

  // Handle missing endpoints gracefully
  const hasProcessStats = processStats && !error

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high': return 'destructive'
      case 'medium': return 'default'
      case 'low': return 'secondary'
      default: return 'outline'
    }
  }

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up': return '↗'
      case 'down': return '↘'
      default: return '→'
    }
  }

  const getTrendColor = (trend: string) => {
    switch (trend) {
      case 'up': return 'text-green-600'
      case 'down': return 'text-red-600'
      default: return 'text-gray-600'
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground">
            Visão geral dos seus processos e atividades
          </p>
        </div>
        <div className="flex items-center space-x-2">
          {canPerformAction('processes', 'create') && (
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Novo Processo
            </Button>
          )}
        </div>
      </div>

      {/* KPI Cards - USANDO DADOS REAIS OU PLACEHOLDERS */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Processos</CardTitle>
            <FileText className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {hasProcessStats ? formatNumber(processStats.total || 0) : '--'}
            </div>
            <div className="text-xs text-muted-foreground">
              {hasProcessStats ? 
                <span>Dados atualizados do sistema</span> : 
                <span className="text-orange-600">Aguardando API /processes/stats</span>
              }
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Processos Ativos</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {hasProcessStats ? formatNumber(processStats.active || 0) : '--'}
            </div>
            <div className="text-xs text-muted-foreground">
              {hasProcessStats ? 
                <span>Em andamento no momento</span> : 
                <span className="text-orange-600">Aguardando API /processes/stats</span>
              }
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Movimentações Hoje</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {hasProcessStats ? formatNumber(processStats.todayMovements || 0) : '--'}
            </div>
            <div className="text-xs text-muted-foreground">
              {hasProcessStats ? 
                <span>Novidades de hoje</span> : 
                <span className="text-orange-600">Aguardando API /processes/stats</span>
              }
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Prazos Próximos</CardTitle>
            <AlertTriangle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {hasProcessStats ? formatNumber(processStats.upcomingDeadlines || 0) : '--'}
            </div>
            <div className="text-xs text-muted-foreground">
              {hasProcessStats ? 
                <span>Requerem atenção</span> : 
                <span className="text-orange-600">Aguardando API /processes/stats</span>
              }
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Recent Activities */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <Clock className="w-5 h-5 mr-2" />
              Atividades Recentes
            </CardTitle>
            <CardDescription>
              Últimas movimentações dos seus processos
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="text-center py-8 text-muted-foreground">
                <Clock className="w-8 h-8 mx-auto mb-2" />
                <p>❌ Mock removido - Implementar busca real de atividades</p>
                <p className="text-xs mt-1">TODO: Conectar a GET /api/v1/reports/recent-activities</p>
              </div>
            </div>
            <Button variant="outline" className="w-full mt-4">
              Ver todas as atividades
            </Button>
          </CardContent>
        </Card>

        {/* Quick Stats */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <BarChart3 className="w-5 h-5 mr-2" />
              Estatísticas Rápidas
            </CardTitle>
            <CardDescription>
              Resumo do seu escritório hoje
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">
                  Taxa de sucesso
                </span>
                <span className="text-sm font-medium">87.2%</span>
              </div>
              <div className="w-full bg-secondary rounded-full h-2">
                <div 
                  className="bg-primary h-2 rounded-full" 
                  style={{ width: '87.2%' }}
                />
              </div>

              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">
                  Tempo médio de resolução
                </span>
                <span className="text-sm font-medium">180 dias</span>
              </div>
              <div className="w-full bg-secondary rounded-full h-2">
                <div 
                  className="bg-blue-500 h-2 rounded-full" 
                  style={{ width: '65%' }}
                />
              </div>

              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">
                  Processos monitorados
                </span>
                <span className="text-sm font-medium">234 de 892</span>
              </div>
              <div className="w-full bg-secondary rounded-full h-2">
                <div 
                  className="bg-green-500 h-2 rounded-full" 
                  style={{ width: '26.2%' }}
                />
              </div>

              <div className="pt-4 border-t">
                <Button variant="outline" className="w-full">
                  Ver relatório completo
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}