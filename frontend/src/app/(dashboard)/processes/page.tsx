'use client'

import { useState } from 'react'
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
import { useProcesses } from '@/hooks/api'
import { useProcessStore } from '@/store'
import { formatDate, getStatusColor, getStatusLabel } from '@/lib/utils'

const mockProcesses = [
  {
    id: '1',
    number: '5001234-20.2023.4.03.6109',
    subject: 'Ação de Cobrança',
    court: 'TJSP - 1ª Vara Cível',
    status: 'active',
    priority: 'high',
    monitoring: true,
    lastMovement: '2025-01-18T10:30:00Z',
    parties: ['João Silva', 'Maria Santos'],
    lawyer: 'Dr. Carlos Oliveira',
    estimatedValue: 50000,
  },
  {
    id: '2',
    number: '5009876-15.2023.4.03.6109',
    subject: 'Ação Trabalhista',
    court: 'TRT - 2ª Região',
    status: 'active',
    priority: 'medium',
    monitoring: false,
    lastMovement: '2025-01-17T14:15:00Z',
    parties: ['Pedro Costa', 'Empresa ABC Ltda'],
    lawyer: 'Dra. Ana Paula',
    estimatedValue: 25000,
  },
  {
    id: '3',
    number: '5005555-30.2023.4.03.6109',
    subject: 'Divórcio Consensual',
    court: 'TJSP - Vara de Família',
    status: 'concluded',
    priority: 'low',
    monitoring: true,
    lastMovement: '2025-01-15T09:00:00Z',
    parties: ['Roberto Lima', 'Sandra Lima'],
    lawyer: 'Dr. Carlos Oliveira',
    estimatedValue: 0,
  },
]

export default function ProcessesPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const { viewMode, setViewMode } = useProcessStore()
  const { data: processes, isLoading } = useProcesses()

  const filteredProcesses = mockProcesses.filter(process => 
    process.number.toLowerCase().includes(searchQuery.toLowerCase()) ||
    process.subject.toLowerCase().includes(searchQuery.toLowerCase()) ||
    process.court.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high': return 'destructive'
      case 'medium': return 'default'
      case 'low': return 'secondary'
      default: return 'outline'
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
                    {process.priority}
                  </Badge>
                </TableCell>
                <TableCell>{formatDate(process.lastMovement)}</TableCell>
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
                      <DropdownMenuItem>
                        <Edit className="w-4 h-4 mr-2" />
                        Editar
                      </DropdownMenuItem>
                      <DropdownMenuItem>
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
                <Badge variant={getStatusColor(process.status)} size="sm">
                  {getStatusLabel(process.status)}
                </Badge>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-xs text-muted-foreground">Prioridade:</span>
                <Badge variant={getPriorityColor(process.priority)} size="sm">
                  {process.priority}
                </Badge>
              </div>
              <div className="text-xs text-muted-foreground">
                <p>Tribunal: {process.court}</p>
                <p>Última mov.: {formatDate(process.lastMovement)}</p>
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
                    <DropdownMenuItem>
                      <Edit className="w-4 h-4 mr-2" />
                      Editar
                    </DropdownMenuItem>
                    <DropdownMenuItem>
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
          <Button>
            <Plus className="w-4 h-4 mr-2" />
            Novo Processo
          </Button>
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
      {viewMode === 'table' ? renderTableView() : renderGridView()}

      {/* Statistics */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {filteredProcesses.length}
            </div>
            <p className="text-xs text-muted-foreground">
              Total de processos
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {filteredProcesses.filter(p => p.status === 'active').length}
            </div>
            <p className="text-xs text-muted-foreground">
              Processos ativos
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {filteredProcesses.filter(p => p.monitoring).length}
            </div>
            <p className="text-xs text-muted-foreground">
              Monitorados
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {filteredProcesses.filter(p => p.priority === 'high').length}
            </div>
            <p className="text-xs text-muted-foreground">
              Alta prioridade
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}