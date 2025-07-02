'use client'

import { useState, useEffect } from 'react'
import { 
  Search, 
  Filter, 
  Plus, 
  MoreVertical, 
  Edit,
  Trash,
  UserCheck,
  UserX,
  Shield,
  Mail,
  Calendar,
  AlertTriangle
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
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
import { useAuthStore, useUserStore } from '@/store'
import { formatDate } from '@/lib/utils'
import { User, UserRole } from '@/types'
import UserModal from '@/components/users/UserModal'
import { useUsersAPI } from '@/hooks/useUsersAPI'

export default function UsersPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  
  const { user: currentUser } = useAuthStore()
  const {
    getUsersForCurrentTenant,
    addUser,
    updateUser,
    deleteUser,
    toggleUserStatus,
    checkUserQuota,
    loadInitialUsers
  } = useUserStore()

  // Load initial users if needed
  useEffect(() => {
    loadInitialUsers()
  }, [loadInitialUsers])

  // Get users for current tenant from API
  const { data: allUsers = [], isLoading: usersLoading, error: usersError } = useUsersAPI()
  const quotaInfo = checkUserQuota()

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

  // Loading state
  if (usersLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-muted-foreground">Carregando usuários...</p>
        </div>
      </div>
    )
  }

  // Error state
  if (usersError) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="max-w-md">
          <CardContent className="p-6 text-center">
            <AlertTriangle className="w-12 h-12 text-red-500 mx-auto mb-4" />
            <h2 className="text-lg font-semibold mb-2 text-red-600">Erro ao Carregar Usuários</h2>
            <p className="text-muted-foreground mb-4">
              {usersError?.message || 'Não foi possível carregar a lista de usuários'}
            </p>
            <Button onClick={() => window.location.reload()}>
              Tentar Novamente
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  const filteredUsers = allUsers.filter(user => 
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

  const handleNewUser = () => {
    setSelectedUser(null)
    setIsModalOpen(true)
  }

  const handleEditUser = (user: User) => {
    setSelectedUser(user)
    setIsModalOpen(true)
  }

  const handleDeleteUser = (userId: string) => {
    if (window.confirm('Tem certeza que deseja excluir este usuário?')) {
      deleteUser(userId)
    }
  }

  const handleToggleStatus = (userId: string) => {
    toggleUserStatus(userId)
  }

  const handleSubmitUser = (userData: { name: string; email: string; role: UserRole; isActive: boolean }) => {
    if (selectedUser) {
      updateUser(selectedUser.id, userData)
    } else {
      addUser(userData)
    }
    setIsModalOpen(false)
    setSelectedUser(null)
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
          <Button onClick={handleNewUser} disabled={!quotaInfo.canAdd}>
            <Plus className="w-4 h-4 mr-2" />
            Novo Usuário
          </Button>
          {!quotaInfo.canAdd && (
            <div className="flex items-center space-x-3">
              <div className="text-sm text-red-600 flex items-center">
                <AlertTriangle className="w-4 h-4 mr-1" />
                Limite atingido: {quotaInfo.used}/{quotaInfo.limit === -1 ? '∞' : quotaInfo.limit} usuários
              </div>
              <Button 
                variant="outline" 
                size="sm" 
                className="border-orange-300 text-orange-600 hover:bg-orange-50"
                onClick={() => window.open('/billing', '_blank')}
              >
                🚀 Upgrade Plano
              </Button>
            </div>
          )}
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
                        <DropdownMenuItem onClick={() => handleEditUser(user)}>
                          <Edit className="w-4 h-4 mr-2" />
                          Editar Usuário
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => handleToggleStatus(user.id)}>
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
                          <DropdownMenuItem 
                            className="text-red-600" 
                            onClick={() => handleDeleteUser(user.id)}
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

      {/* Statistics */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {allUsers.length}
            </div>
            <p className="text-xs text-muted-foreground">
              Total de usuários
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {allUsers.filter(u => u.isActive).length}
            </div>
            <p className="text-xs text-muted-foreground">
              Usuários ativos
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {allUsers.filter(u => u.role === 'admin').length}
            </div>
            <p className="text-xs text-muted-foreground">
              Administradores
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">
              {allUsers.filter(u => u.role === 'lawyer').length}
            </div>
            <p className="text-xs text-muted-foreground">
              Advogados
            </p>
          </CardContent>
        </Card>
      </div>

      {/* User Modal */}
      <UserModal 
        isOpen={isModalOpen}
        onClose={() => {
          setIsModalOpen(false)
          setSelectedUser(null)
        }}
        onSubmit={handleSubmitUser}
        user={selectedUser}
        quotaInfo={quotaInfo}
      />
    </div>
  )
}