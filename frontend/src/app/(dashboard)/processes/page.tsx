'use client'

import { useState, useMemo } from 'react'
import { 
  Search, 
  Filter, 
  Plus, 
  MoreVertical, 
  Eye,
  Edit,
  Trash,
  Bell,
  BellOff,
  Grid,
  List,
  Table as TableIcon
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
import { useProcessDataStore, useProcessUIStore } from '@/store'
import { usePermissions } from '@/hooks/usePermissions'
import { ProcessModal } from '@/components/processes/ProcessModal'
import { formatDate, getStatusColor, getStatusLabel } from '@/lib/utils'

export default function ProcessesPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingProcessId, setEditingProcessId] = useState<string | undefined>()
  
  const { viewMode, setViewMode } = useProcessUIStore()
  const { 
    processes,
    getProcessesByFilter, 
    deleteProcess, 
    toggleMonitoring, 
    getStats 
  } = useProcessDataStore()
  const { canPerformAction } = usePermissions()

  const filteredProcesses = useMemo(() => {
    return getProcessesByFilter({
      search: searchQuery
    })
  }, [getProcessesByFilter, searchQuery, processes]) // Add processes to dependency array

  const stats = getStats()

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high': return 'destructive'
      case 'urgent': return 'destructive'
      case 'medium': return 'default'
      case 'low': return 'secondary'
      default: return 'outline'
    }
  }

  const getPriorityLabel = (priority: string) => {
    switch (priority) {
      case 'high': return 'Alta'
      case 'urgent': return 'Urgente'
      case 'medium': return 'Média'
      case 'low': return 'Baixa'
      default: return priority
    }
  }

  const renderTableView = () => (
    <Card>
      <CardContent className="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Número</TableHead>
              <TableHead>Assunto</TableHead>
              <TableHead>Tribunal</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Prioridade</TableHead>
              <TableHead>Última Movimentação</TableHead>
              <TableHead>Ações</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredProcesses.map((process) => (
              <TableRow key={process.id}>
                <TableCell className="font-medium">
                  <div className="flex items-center space-x-2">
                    <span>{process.number}</span>
                    {process.monitoring && (
                      <Bell className="w-4 h-4 text-blue-500" />
                    )}
                  </div>
                </TableCell>
                <TableCell>{process.subject}</TableCell>
                <TableCell>{process.court}</TableCell>
                <TableCell>
                  <Badge variant={getStatusColor(process.status)}>
                    {getStatusLabel(process.status)}
                  </Badge>
                </TableCell>
                <TableCell>
                  <Badge variant={getPriorityColor(process.priority)}>
                    {getPriorityLabel(process.priority)}
                  </Badge>
                </TableCell>
                <TableCell>{process.lastMovement ? formatDate(process.lastMovement) : 'N/A'}</TableCell>
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
                      {canPerformAction('processes', 'edit') && (
                        <DropdownMenuItem onClick={() => {
                          setEditingProcessId(process.id)
                          setIsModalOpen(true)
                        }}>
                          <Edit className="w-4 h-4 mr-2" />
                          Editar
                        </DropdownMenuItem>
                      )}
                      {canPerformAction('processes', 'edit') && (
                        <DropdownMenuItem onClick={() => toggleMonitoring(process.id)}>
                          {process.monitoring ? (
                            <>
                              <BellOff className="w-4 h-4 mr-2" />
                              Parar Monitoramento
                            </>
                          ) : (
                            <>
                              <Bell className="w-4 h-4 mr-2" />
                              Monitorar
                            </>
                          )}
                        </DropdownMenuItem>
                      )}
                      {canPerformAction('processes', 'delete') && (
                        <DropdownMenuItem 
                          className="text-red-600"
                          onClick={() => {
                            if (confirm(`Tem certeza que deseja excluir o processo ${process.number}?`)) {
                              deleteProcess(process.id)
                            }
                          }}
                        >
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
  )

  const renderListView = () => (
    <div className="space-y-3">
      {filteredProcesses.map((process) => (
        <Card key={process.id} className="hover:shadow-md transition-shadow">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <div className="flex items-center space-x-3 mb-2">
                  <h3 className="font-semibold text-lg">{process.number}</h3>
                  {process.monitoring && (
                    <Bell className="w-4 h-4 text-blue-500" />
                  )}
                </div>
                <p className="text-muted-foreground mb-2">{process.subject}</p>
                <div className="flex items-center space-x-4 text-sm">
                  <span className="text-muted-foreground">Tribunal: {process.court}</span>
                  <Badge variant={getStatusColor(process.status)}>
                    {getStatusLabel(process.status)}
                  </Badge>
                  <Badge variant={getPriorityColor(process.priority)}>
                    {getPriorityLabel(process.priority)}
                  </Badge>
                  <span className="text-muted-foreground">
                    Última mov.: {process.lastMovement ? formatDate(process.lastMovement) : 'N/A'}
                  </span>
                </div>
              </div>
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
                  {canPerformAction('processes', 'edit') && (
                    <DropdownMenuItem onClick={() => {
                      setEditingProcessId(process.id)
                      setIsModalOpen(true)
                    }}>
                      <Edit className="w-4 h-4 mr-2" />
                      Editar
                    </DropdownMenuItem>
                  )}
                  {canPerformAction('processes', 'edit') && (
                    <DropdownMenuItem onClick={() => toggleMonitoring(process.id)}>
                      {process.monitoring ? (
                        <>
                          <BellOff className="w-4 h-4 mr-2" />
                          Parar Monitoramento
                        </>
                      ) : (
                        <>
                          <Bell className="w-4 h-4 mr-2" />
                          Monitorar
                        </>
                      )}
                    </DropdownMenuItem>
                  )}
                  {canPerformAction('processes', 'delete') && (
                    <DropdownMenuItem 
                      className="text-red-600"
                      onClick={() => {
                        if (confirm(`Tem certeza que deseja excluir o processo ${process.number}?`)) {
                          deleteProcess(process.id)
                        }
                      }}
                    >
                      <Trash className="w-4 h-4 mr-2" />
                      Excluir
                    </DropdownMenuItem>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )

  const renderGridView = () => (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {filteredProcesses.map((process) => (
        <Card key={process.id} className="hover:shadow-md transition-shadow">
          <CardHeader className="pb-3">
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <CardTitle className="text-sm font-medium">
                  {process.number}
                </CardTitle>
                <p className="text-sm text-muted-foreground mt-1">
                  {process.subject}
                </p>
              </div>
              {process.monitoring && (
                <Bell className="w-4 h-4 text-blue-500" />
              )}
            </div>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-xs text-muted-foreground">Status:</span>
                <Badge variant={getStatusColor(process.status)}>
                  {getStatusLabel(process.status)}
                </Badge>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-xs text-muted-foreground">Prioridade:</span>
                <Badge variant={getPriorityColor(process.priority)}>
                  {getPriorityLabel(process.priority)}
                </Badge>
              </div>
              <div className="text-xs text-muted-foreground">
                <p>Tribunal: {process.court}</p>
                <p>Última mov.: {process.lastMovement ? formatDate(process.lastMovement) : 'N/A'}</p>
              </div>
              <div className="flex justify-end pt-2">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="sm">
                      <MoreVertical className="w-4 h-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem>
                      <Eye className="w-4 h-4 mr-2" />
                      Visualizar
                    </DropdownMenuItem>
                    {canPerformAction('processes', 'edit') && (
                      <DropdownMenuItem onClick={() => {
                        setEditingProcessId(process.id)
                        setIsModalOpen(true)
                      }}>
                        <Edit className="w-4 h-4 mr-2" />
                        Editar
                      </DropdownMenuItem>
                    )}
                    {canPerformAction('processes', 'edit') && (
                      <DropdownMenuItem onClick={() => toggleMonitoring(process.id)}>
                        {process.monitoring ? (
                          <>
                            <BellOff className="w-4 h-4 mr-2" />
                            Parar Monitoramento
                          </>
                        ) : (
                          <>
                            <Bell className="w-4 h-4 mr-2" />
                            Monitorar
                          </>
                        )}
                      </DropdownMenuItem>
                    )}
                    {canPerformAction('processes', 'delete') && (
                      <DropdownMenuItem 
                        className="text-red-600"
                        onClick={() => {
                          if (confirm(`Tem certeza que deseja excluir o processo ${process.number}?`)) {
                            deleteProcess(process.id)
                          }
                        }}
                      >
                        <Trash className="w-4 h-4 mr-2" />
                        Excluir
                      </DropdownMenuItem>
                    )}
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Processos</h1>
          <p className="text-muted-foreground">
            Gerencie todos os seus processos jurídicos
          </p>
        </div>
        <div className="flex items-center space-x-2">
          {canPerformAction('processes', 'create') && (
            <Button onClick={() => setIsModalOpen(true)}>
              <Plus className="w-4 h-4 mr-2" />
              Novo Processo
            </Button>
          )}
        </div>
      </div>

      {/* Filters and Search */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-2 flex-1 max-w-md">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <Input
              placeholder="Buscar processos..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>
          <Button variant="outline" size="icon">
            <Filter className="w-4 h-4" />
          </Button>
        </div>

        {/* View Mode Toggle */}
        <div className="flex items-center space-x-1 border rounded-lg p-1">
          <Button
            variant={viewMode === 'table' ? 'default' : 'ghost'}
            size="sm"
            onClick={() => setViewMode('table')}
          >
            <TableIcon className="w-4 h-4" />
          </Button>
          <Button
            variant={viewMode === 'grid' ? 'default' : 'ghost'}
            size="sm"
            onClick={() => setViewMode('grid')}
          >
            <Grid className="w-4 h-4" />
          </Button>
          <Button
            variant={viewMode === 'list' ? 'default' : 'ghost'}
            size="sm"
            onClick={() => setViewMode('list')}
          >
            <List className="w-4 h-4" />
          </Button>
        </div>
      </div>

      {/* Content */}
      {viewMode === 'table' && renderTableView()}
      {viewMode === 'grid' && renderGridView()}
      {viewMode === 'list' && renderListView()}

      {/* Statistics */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {stats.total}
            </div>
            <p className="text-xs text-muted-foreground">
              Total de processos
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {stats.active}
            </div>
            <p className="text-xs text-muted-foreground">
              Processos ativos
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {stats.monitoring}
            </div>
            <p className="text-xs text-muted-foreground">
              Monitorados
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {stats.highPriority}
            </div>
            <p className="text-xs text-muted-foreground">
              Alta prioridade
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Process Modal */}
      <ProcessModal
        isOpen={isModalOpen}
        onClose={() => {
          setIsModalOpen(false)
          setEditingProcessId(undefined)
          // Force re-render by updating a dummy state
          setSearchQuery(searchQuery)
        }}
        processId={editingProcessId}
      />
    </div>
  )
}