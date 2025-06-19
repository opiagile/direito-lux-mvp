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
import { formatNumber, formatPercentage } from '@/lib/utils'

const mockKPIData = [
  {
    title: 'Total de Processos',
    value: 1247,
    change: 12.5,
    icon: FileText,
    trend: 'up' as const,
  },
  {
    title: 'Processos Ativos',
    value: 892,
    change: 8.2,
    icon: Activity,
    trend: 'up' as const,
  },
  {
    title: 'Movimentações Hoje',
    value: 23,
    change: -4.1,
    icon: TrendingUp,
    trend: 'down' as const,
  },
  {
    title: 'Prazos Próximos',
    value: 8,
    change: 0,
    icon: AlertTriangle,
    trend: 'stable' as const,
  },
]

const recentActivities = [
  {
    id: 1,
    process: '5001234-20.2023.4.03.6109',
    action: 'Nova movimentação',
    description: 'Sentença publicada',
    time: '2 horas atrás',
    priority: 'high' as const,
  },
  {
    id: 2,
    process: '5009876-15.2023.4.03.6109',
    action: 'Prazo vencendo',
    description: 'Contestação em 2 dias',
    time: '4 horas atrás',
    priority: 'medium' as const,
  },
  {
    id: 3,
    process: '5005555-30.2023.4.03.6109',
    action: 'Documento anexado',
    description: 'Petição inicial',
    time: '6 horas atrás',
    priority: 'low' as const,
  },
]

export default function DashboardPage() {
  const { data: processStats, isLoading } = useProcessStats()

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
          <Button>
            <Plus className="w-4 h-4 mr-2" />
            Novo Processo
          </Button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {mockKPIData.map((kpi) => (
          <Card key={kpi.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {kpi.title}
              </CardTitle>
              <kpi.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {formatNumber(kpi.value)}
              </div>
              <div className="flex items-center space-x-2 text-xs text-muted-foreground">
                <span className={getTrendColor(kpi.trend)}>
                  {getTrendIcon(kpi.trend)} {Math.abs(kpi.change)}%
                </span>
                <span>em relação ao mês anterior</span>
              </div>
            </CardContent>
          </Card>
        ))}
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
              {recentActivities.map((activity) => (
                <div
                  key={activity.id}
                  className="flex items-start space-x-3 p-3 rounded-lg border"
                >
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center justify-between">
                      <p className="text-sm font-medium truncate">
                        {activity.process}
                      </p>
                      <Badge variant={getPriorityColor(activity.priority)}>
                        {activity.action}
                      </Badge>
                    </div>
                    <p className="text-sm text-muted-foreground mt-1">
                      {activity.description}
                    </p>
                    <p className="text-xs text-muted-foreground mt-1">
                      {activity.time}
                    </p>
                  </div>
                </div>
              ))}
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