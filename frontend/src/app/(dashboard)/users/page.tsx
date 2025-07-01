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
  UserCheck,
  UserX,
  Shield,
  Mail,
  Phone,
  Calendar
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
import { useAuthStore } from '@/store'
import { formatDate } from '@/lib/utils'
import { User, UserRole } from '@/types'

const mockUsers: User[] = [
  {
    id: '1',
    email: 'admin@silvaassociados.com.br',
    name: 'Carlos Silva',
    role: 'admin',
    tenantId: '11111111-1111-1111-1111-111111111111',
    isActive: true,
    createdAt: '2024-01-15T09:00:00Z',
    updatedAt: '2025-01-20T14:30:00Z'
  },
  {
    id: '2',
    email: 'gerente@silvaassociados.com.br',
    name: 'Ana Paula Santos',
    role: 'manager',
    tenantId: '11111111-1111-1111-1111-111111111111',
    isActive: true,
    createdAt: '2024-01-16T10:00:00Z',
    updatedAt: '2025-01-19T16:20:00Z'
  },
  {
    id: '3',
    email: 'advogado@silvaassociados.com.br',
    name: 'Dr. Roberto Lima',
    role: 'lawyer',
    tenantId: '11111111-1111-1111-1111-111111111111',
    isActive: true,
    createdAt: '2024-01-17T11:00:00Z',
    updatedAt: '2025-01-18T12:45:00Z'
  },
  {
    id: '4',
    email: 'cliente@silvaassociados.com.br',
    name: 'Maria José Oliveira',
    role: 'assistant',
    tenantId: '11111111-1111-1111-1111-111111111111',
    isActive: true,
    createdAt: '2024-01-18T14:00:00Z',
    updatedAt: '2025-01-17T09:15:00Z'
  },
  {
    id: '5',
    email: 'advogado2@silvaassociados.com.br',
    name: 'Dra. Fernanda Costa',
    role: 'lawyer',
    tenantId: '11111111-1111-1111-1111-111111111111',
    isActive: false,
    createdAt: '2024-02-01T08:00:00Z',
    updatedAt: '2025-01-10T11:30:00Z'
  }
]

export default function UsersPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const { user: currentUser } = useAuthStore()

  // Verificar permissões - apenas admin pode acessar
  if (currentUser?.role !== 'admin') {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="max-w-md">
          <CardContent className="p-6 text-center">
            <Shield className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
            <h2 className="text-lg font-semibold mb-2">Acesso Restrito</h2>
            <p className="text-muted-foreground">
              Apenas administradores podem acessar a gestão de usuários.
            </p>
          </CardContent>
        </Card>
      </div>
    )
  }

  const filteredUsers = mockUsers.filter(user => 
    user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    user.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
    user.role.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const getRoleColor = (role: UserRole): "default" | "secondary" | "destructive" | "outline" => {
    switch (role) {
      case 'admin': return 'destructive'
      case 'manager': return 'default'
      case 'lawyer': return 'secondary'
      case 'assistant': return 'outline'
      default: return 'outline'
    }
  }

  const getRoleLabel = (role: UserRole): string => {
    switch (role) {
      case 'admin': return 'Administrador'
      case 'manager': return 'Gerente'
      case 'lawyer': return 'Advogado'
      case 'assistant': return 'Assistente'
      default: return role
    }
  }

  const getStatusColor = (isActive: boolean): "default" | "secondary" => {
    return isActive ? 'default' : 'secondary'
  }

  const getStatusLabel = (isActive: boolean): string => {
    return isActive ? 'Ativo' : 'Inativo'
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Gestão de Usuários</h1>
          <p className="text-muted-foreground">
            Gerencie usuários, permissões e acessos do seu escritório
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Button>
            <Plus className="w-4 h-4 mr-2" />
            Novo Usuário
          </Button>
        </div>
      </div>

      {/* Filters and Search */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-2 flex-1 max-w-md">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <Input
              placeholder="Buscar usuários..."
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

      {/* Users Table */}
      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Usuário</TableHead>
                <TableHead>Email</TableHead>
                <TableHead>Role</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Último Acesso</TableHead>
                <TableHead>Criado em</TableHead>
                <TableHead>Ações</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredUsers.map((user) => (
                <TableRow key={user.id}>
                  <TableCell className="font-medium">
                    <div className="flex items-center space-x-3">
                      <div className="w-8 h-8 rounded-full bg-gradient-to-r from-blue-500 to-purple-600 flex items-center justify-center text-white text-sm font-medium">
                        {user.name.split(' ').map(n => n[0]).join('').toUpperCase()}
                      </div>
                      <div>
                        <div className="font-medium">{user.name}</div>
                        <div className="text-sm text-muted-foreground">ID: {user.id}</div>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center space-x-2">
                      <Mail className="w-4 h-4 text-muted-foreground" />
                      <span>{user.email}</span>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant={getRoleColor(user.role)}>
                      {getRoleLabel(user.role)}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center space-x-2">
                      {user.isActive ? (
                        <UserCheck className="w-4 h-4 text-green-500" />
                      ) : (
                        <UserX className="w-4 h-4 text-red-500" />
                      )}
                      <Badge variant={getStatusColor(user.isActive)}>
                        {getStatusLabel(user.isActive)}
                      </Badge>
                    </div>
                  </TableCell>
                  <TableCell>{formatDate(user.updatedAt)}</TableCell>
                  <TableCell>
                    <div className="flex items-center space-x-2">
                      <Calendar className="w-4 h-4 text-muted-foreground" />
                      <span>{formatDate(user.createdAt)}</span>
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
                          <Eye className="w-4 h-4 mr-2" />
                          Visualizar Perfil
                        </DropdownMenuItem>
                        <DropdownMenuItem>
                          <Edit className="w-4 h-4 mr-2" />
                          Editar Usuário
                        </DropdownMenuItem>
                        <DropdownMenuItem>
                          <Shield className="w-4 h-4 mr-2" />
                          Alterar Permissões
                        </DropdownMenuItem>
                        <DropdownMenuItem>
                          {user.isActive ? (
                            <>
                              <UserX className="w-4 h-4 mr-2" />
                              Desativar
                            </>
                          ) : (
                            <>
                              <UserCheck className="w-4 h-4 mr-2" />
                              Ativar
                            </>
                          )}
                        </DropdownMenuItem>
                        {user.id !== currentUser?.id && (
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
              {filteredUsers.length}
            </div>
            <p className="text-xs text-muted-foreground">
              Total de usuários
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {filteredUsers.filter(u => u.isActive).length}
            </div>
            <p className="text-xs text-muted-foreground">
              Usuários ativos
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {filteredUsers.filter(u => u.role === 'admin').length}
            </div>
            <p className="text-xs text-muted-foreground">
              Administradores
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {filteredUsers.filter(u => u.role === 'lawyer').length}
            </div>
            <p className="text-xs text-muted-foreground">
              Advogados
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}