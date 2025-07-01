'use client'

import { useState } from 'react'
import { 
  FileText, 
  Download, 
  Calendar, 
  Filter, 
  Plus, 
  MoreVertical, 
  Eye,
  Trash,
  Clock,
  CheckCircle,
  AlertCircle,
  BarChart3,
  PieChart,
  TrendingUp,
  FileSpreadsheet,
  Search,
  Settings,
  Play,
  Pause
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useAuthStore } from '@/store'
import { formatDate } from '@/lib/utils'
import { Report, ReportType, ReportStatus, ReportSchedule } from '@/types'

const mockReports: Report[] = [
  {
    id: '1',
    tenantId: '11111111-1111-1111-1111-111111111111',
    type: 'executive_summary',
    title: 'Relatório Executivo - Janeiro 2025',
    format: 'pdf',
    status: 'completed',
    filters: { period: '2025-01' },
    fileUrl: '/reports/exec_jan_2025.pdf',
    fileSize: 2048576,
    processingTime: 5.2,
    scheduledBy: 'admin@silvaassociados.com.br',
    createdAt: '2025-01-20T10:30:00Z',
    completedAt: '2025-01-20T10:35:00Z',
    expiresAt: '2025-02-20T10:35:00Z'
  },
  {
    id: '2',
    tenantId: '11111111-1111-1111-1111-111111111111',
    type: 'process_analysis',
    title: 'Análise de Processos - Dezembro 2024',
    format: 'excel',
    status: 'completed',
    filters: { period: '2024-12', court: 'TJSP' },
    fileUrl: '/reports/proc_analysis_dec_2024.xlsx',
    fileSize: 1536000,
    processingTime: 8.7,
    scheduledBy: 'gerente@silvaassociados.com.br',
    createdAt: '2025-01-18T14:15:00Z',
    completedAt: '2025-01-18T14:24:00Z',
    expiresAt: '2025-02-18T14:24:00Z'
  },
  {
    id: '3',
    tenantId: '11111111-1111-1111-1111-111111111111',
    type: 'productivity',
    title: 'Relatório de Produtividade - Q4 2024',
    format: 'pdf',
    status: 'processing',
    filters: { period: 'Q4-2024' },
    scheduledBy: 'admin@silvaassociados.com.br',
    createdAt: '2025-01-20T16:00:00Z'
  },
  {
    id: '4',
    tenantId: '11111111-1111-1111-1111-111111111111',
    type: 'financial',
    title: 'Relatório Financeiro - 2024',
    format: 'excel',
    status: 'failed',
    filters: { year: '2024' },
    scheduledBy: 'admin@silvaassociados.com.br',
    createdAt: '2025-01-19T09:00:00Z'
  }
]

const mockSchedules: ReportSchedule[] = [
  {
    id: '1',
    tenantId: '11111111-1111-1111-1111-111111111111',
    reportType: 'executive_summary',
    title: 'Relatório Executivo Mensal',
    format: 'pdf',
    frequency: 'monthly',
    filters: { automated: true },
    recipients: ['admin@silvaassociados.com.br', 'gerente@silvaassociados.com.br'],
    isActive: true,
    lastRunAt: '2025-01-20T10:30:00Z',
    nextRunAt: '2025-02-20T10:30:00Z',
    createdAt: '2024-12-01T10:00:00Z'
  },
  {
    id: '2',
    tenantId: '11111111-1111-1111-1111-111111111111',
    reportType: 'productivity',
    title: 'Relatório de Produtividade Semanal',
    format: 'excel',
    frequency: 'weekly',
    filters: { includeGraphs: true },
    recipients: ['gerente@silvaassociados.com.br'],
    isActive: true,
    lastRunAt: '2025-01-17T08:00:00Z',
    nextRunAt: '2025-01-24T08:00:00Z',
    createdAt: '2024-11-15T14:00:00Z'
  },
  {
    id: '3',
    tenantId: '11111111-1111-1111-1111-111111111111',
    reportType: 'process_analysis',
    title: 'Análise Trimestral de Processos',
    format: 'pdf',
    frequency: 'custom',
    cronExpression: '0 0 1 */3 *',
    filters: { detailed: true },
    recipients: ['admin@silvaassociados.com.br'],
    isActive: false,
    lastRunAt: '2024-10-01T00:00:00Z',
    nextRunAt: '2025-01-01T00:00:00Z',
    createdAt: '2024-09-01T10:00:00Z'
  }
]

export default function ReportsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const { user: currentUser } = useAuthStore()

  // Verificar permissões - admin e manager têm acesso total, lawyer acesso limitado
  if (!currentUser || !['admin', 'manager', 'lawyer'].includes(currentUser.role)) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="max-w-md">
          <CardContent className="p-6 text-center">
            <FileText className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
            <h2 className="text-lg font-semibold mb-2">Acesso Restrito</h2>
            <p className="text-muted-foreground">
              Você não tem permissão para acessar os relatórios.
            </p>
          </CardContent>
        </Card>
      </div>
    )
  }

  const filteredReports = mockReports.filter(report => 
    report.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    report.type.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const getStatusColor = (status: ReportStatus): "default" | "secondary" | "destructive" | "outline" => {
    switch (status) {
      case 'completed': return 'default'
      case 'processing': return 'secondary'
      case 'failed': return 'destructive'
      case 'cancelled': return 'outline'
      default: return 'outline'
    }
  }

  const getStatusIcon = (status: ReportStatus) => {
    switch (status) {
      case 'completed': return <CheckCircle className="w-4 h-4" />
      case 'processing': return <Clock className="w-4 h-4" />
      case 'failed': return <AlertCircle className="w-4 h-4" />
      default: return <Clock className="w-4 h-4" />
    }
  }

  const getTypeIcon = (type: ReportType) => {
    switch (type) {
      case 'executive_summary': return <TrendingUp className="w-4 h-4" />
      case 'process_analysis': return <BarChart3 className="w-4 h-4" />
      case 'productivity': return <PieChart className="w-4 h-4" />
      case 'financial': return <FileSpreadsheet className="w-4 h-4" />
      default: return <FileText className="w-4 h-4" />
    }
  }

  const getTypeLabel = (type: ReportType): string => {
    switch (type) {
      case 'executive_summary': return 'Resumo Executivo'
      case 'process_analysis': return 'Análise de Processos'
      case 'productivity': return 'Produtividade'
      case 'financial': return 'Financeiro'
      case 'legal_timeline': return 'Timeline Jurídico'
      case 'jurisprudence_analysis': return 'Análise Jurisprudencial'
      default: return type
    }
  }

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Relatórios</h1>
          <p className="text-muted-foreground">
            Gerencie e analise relatórios de performance e dados jurídicos
          </p>
        </div>
        <div className="flex items-center space-x-2">
          {['admin', 'manager'].includes(currentUser?.role || '') && (
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Novo Relatório
            </Button>
          )}
        </div>
      </div>

      {/* Reports Tabs */}
      <Tabs defaultValue="reports" className="space-y-6">
        <TabsList>
          <TabsTrigger value="reports">Relatórios</TabsTrigger>
          {['admin', 'manager'].includes(currentUser?.role || '') && (
            <TabsTrigger value="schedules">Agendamentos</TabsTrigger>
          )}
          <TabsTrigger value="templates">Templates</TabsTrigger>
        </TabsList>

        {/* Reports List */}
        <TabsContent value="reports" className="space-y-6">
          {/* Filters and Search */}
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2 flex-1 max-w-md">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <Input
                  placeholder="Buscar relatórios..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Button variant="outline" size="icon">
                <Filter className="w-4 h-4" />
              </Button>
            </div>
          </div>

          {/* Reports Table */}
          <Card>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Relatório</TableHead>
                    <TableHead>Tipo</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Formato</TableHead>
                    <TableHead>Tamanho</TableHead>
                    <TableHead>Criado em</TableHead>
                    <TableHead>Ações</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredReports.map((report) => (
                    <TableRow key={report.id}>
                      <TableCell className="font-medium">
                        <div className="flex items-center space-x-3">
                          {getTypeIcon(report.type)}
                          <div>
                            <div className="font-medium">{report.title}</div>
                            <div className="text-sm text-muted-foreground">
                              Por: {report.scheduledBy}
                            </div>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">
                          {getTypeLabel(report.type)}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center space-x-2">
                          {getStatusIcon(report.status)}
                          <Badge variant={getStatusColor(report.status)}>
                            {report.status}
                          </Badge>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant="secondary">
                          {report.format.toUpperCase()}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        {report.fileSize ? formatFileSize(report.fileSize) : '-'}
                      </TableCell>
                      <TableCell>{formatDate(report.createdAt)}</TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon">
                              <MoreVertical className="w-4 h-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem>
                              <Eye className="w-4 h-4 mr-2" />
                              Visualizar
                            </DropdownMenuItem>
                            {report.status === 'completed' && (
                              <DropdownMenuItem>
                                <Download className="w-4 h-4 mr-2" />
                                Download
                              </DropdownMenuItem>
                            )}
                            {['admin', 'manager'].includes(currentUser?.role || '') && (
                              <DropdownMenuItem className="text-red-600">
                                <Trash className="w-4 h-4 mr-2" />
                                Excluir
                              </DropdownMenuItem>
                            )}
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          {/* Statistics */}
          <div className="grid gap-4 md:grid-cols-4">
            <Card>
              <CardContent className="p-4">
                <div className="text-2xl font-bold">
                  {filteredReports.length}
                </div>
                <p className="text-xs text-muted-foreground">
                  Total de relatórios
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4">
                <div className="text-2xl font-bold">
                  {filteredReports.filter(r => r.status === 'completed').length}
                </div>
                <p className="text-xs text-muted-foreground">
                  Concluídos
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4">
                <div className="text-2xl font-bold">
                  {filteredReports.filter(r => r.status === 'processing').length}
                </div>
                <p className="text-xs text-muted-foreground">
                  Em processamento
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4">
                <div className="text-2xl font-bold">
                  {mockSchedules.filter(s => s.isActive).length}
                </div>
                <p className="text-xs text-muted-foreground">
                  Agendamentos ativos
                </p>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Schedules */}
        {['admin', 'manager'].includes(currentUser?.role || '') && (
          <TabsContent value="schedules" className="space-y-6">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold">Agendamentos de Relatórios</h2>
              <Button>
                <Calendar className="w-4 h-4 mr-2" />
                Novo Agendamento
              </Button>
            </div>

            <Card>
              <CardContent className="p-0">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Relatório</TableHead>
                      <TableHead>Frequência</TableHead>
                      <TableHead>Destinatários</TableHead>
                      <TableHead>Próxima Execução</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Ações</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {mockSchedules.map((schedule) => (
                      <TableRow key={schedule.id}>
                        <TableCell className="font-medium">
                          <div className="flex items-center space-x-3">
                            {getTypeIcon(schedule.reportType)}
                            <div>
                              <div className="font-medium">{schedule.title}</div>
                              <div className="text-sm text-muted-foreground">
                                {getTypeLabel(schedule.reportType)} - {schedule.format.toUpperCase()}
                              </div>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline">
                            {schedule.frequency}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <div className="text-sm">
                            {schedule.recipients.length} destinatários
                          </div>
                        </TableCell>
                        <TableCell>
                          {schedule.nextRunAt ? formatDate(schedule.nextRunAt) : '-'}
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center space-x-2">
                            {schedule.isActive ? (
                              <Play className="w-4 h-4 text-green-500" />
                            ) : (
                              <Pause className="w-4 h-4 text-gray-500" />
                            )}
                            <Badge variant={schedule.isActive ? 'default' : 'secondary'}>
                              {schedule.isActive ? 'Ativo' : 'Pausado'}
                            </Badge>
                          </div>
                        </TableCell>
                        <TableCell>
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="icon">
                                <MoreVertical className="w-4 h-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem>
                                <Settings className="w-4 h-4 mr-2" />
                                Configurar
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                {schedule.isActive ? (
                                  <>
                                    <Pause className="w-4 h-4 mr-2" />
                                    Pausar
                                  </>
                                ) : (
                                  <>
                                    <Play className="w-4 h-4 mr-2" />
                                    Ativar
                                  </>
                                )}
                              </DropdownMenuItem>
                              <DropdownMenuItem className="text-red-600">
                                <Trash className="w-4 h-4 mr-2" />
                                Excluir
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </TabsContent>
        )}

        {/* Templates */}
        <TabsContent value="templates" className="space-y-6">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold">Templates de Relatórios</h2>
            {['admin', 'manager'].includes(currentUser?.role || '') && (
              <Button>
                <Plus className="w-4 h-4 mr-2" />
                Novo Template
              </Button>
            )}
          </div>

          <div className="grid gap-4 md:grid-cols-3">
            {[
              { type: 'executive_summary', title: 'Resumo Executivo', description: 'Visão geral do escritório com KPIs principais' },
              { type: 'process_analysis', title: 'Análise de Processos', description: 'Análise detalhada de processos e movimentações' },
              { type: 'productivity', title: 'Produtividade', description: 'Métricas de produtividade da equipe' },
              { type: 'financial', title: 'Financeiro', description: 'Relatórios financeiros e de faturamento' },
              { type: 'legal_timeline', title: 'Timeline Jurídico', description: 'Cronologia de eventos jurídicos importantes' },
              { type: 'jurisprudence_analysis', title: 'Análise Jurisprudencial', description: 'Análise de jurisprudência e tendências' }
            ].map((template) => (
              <Card key={template.type} className="hover:shadow-md transition-shadow cursor-pointer">
                <CardHeader className="pb-3">
                  <div className="flex items-center space-x-3">
                    {getTypeIcon(template.type as ReportType)}
                    <CardTitle className="text-base">{template.title}</CardTitle>
                  </div>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground mb-4">
                    {template.description}
                  </p>
                  <Button size="sm" className="w-full">
                    Gerar Relatório
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}