'use client'

import { useState } from 'react'
import { 
  Bell, 
  Check, 
  Search, 
  Calendar, 
  Clock, 
  AlertTriangle, 
  Info, 
  CheckCircle, 
  XCircle,
  Mail,
  MessageCircle,
  Smartphone,
  MoreVertical,
  Settings,
  Archive,
  Trash,
  MailCheck,
  Eye,
  EyeOff
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
import { formatDate } from '@/lib/utils'
import { Notification, NotificationType, NotificationChannel, NotificationStatus, NotificationPriority } from '@/types'

const mockNotifications: Notification[] = [
  {
    id: '1',
    tenantId: '11111111-1111-1111-1111-111111111111',
    userId: 'current-user',
    type: 'process_update',
    channel: ['whatsapp', 'email'],
    title: 'Nova Movimentação - Processo 5001234-20.2023.4.03.6109',
    message: 'Foi registrada uma nova movimentação no processo de Ação de Cobrança. Tipo: Sentença de Procedência.',
    data: { processNumber: '5001234-20.2023.4.03.6109', court: 'TJSP' },
    status: 'delivered',
    priority: 'high',
    sentAt: '2025-01-20T10:30:00Z',
    readAt: undefined,
    createdAt: '2025-01-20T10:30:00Z'
  },
  {
    id: '2',
    tenantId: '11111111-1111-1111-1111-111111111111',
    type: 'deadline',
    channel: ['email', 'push'],
    title: 'Prazo se Aproximando',
    message: 'O prazo para contestação no processo 5009876-15.2023.4.03.6109 vence em 2 dias (22/01/2025).',
    data: { processNumber: '5009876-15.2023.4.03.6109', deadline: '2025-01-22T23:59:59Z' },
    status: 'delivered',
    priority: 'critical',
    sentAt: '2025-01-20T08:00:00Z',
    readAt: '2025-01-20T09:15:00Z',
    createdAt: '2025-01-20T08:00:00Z'
  },
  {
    id: '3',
    tenantId: '11111111-1111-1111-1111-111111111111',
    userId: 'current-user',
    type: 'system',
    channel: ['email'],
    title: 'Backup Realizado com Sucesso',
    message: 'O backup automático dos dados foi concluído com sucesso. Todos os dados estão seguros.',
    status: 'delivered',
    priority: 'low',
    sentAt: '2025-01-20T03:00:00Z',
    readAt: '2025-01-20T07:30:00Z',
    createdAt: '2025-01-20T03:00:00Z'
  },
  {
    id: '4',
    tenantId: '11111111-1111-1111-1111-111111111111',
    type: 'process_update',
    channel: ['whatsapp'],
    title: 'Audiência Marcada',
    message: 'Foi marcada audiência de conciliação para o processo 5005555-30.2023.4.03.6109 em 25/01/2025 às 14:00.',
    data: { processNumber: '5005555-30.2023.4.03.6109', hearingDate: '2025-01-25T14:00:00Z' },
    status: 'sent',
    priority: 'high',
    sentAt: '2025-01-19T16:45:00Z',
    readAt: undefined,
    createdAt: '2025-01-19T16:45:00Z'
  },
  {
    id: '5',
    tenantId: '11111111-1111-1111-1111-111111111111',
    type: 'marketing',
    channel: ['email'],
    title: 'Nova Funcionalidade: AI Assistant',
    message: 'Descubra a nova funcionalidade de AI Assistant para análise automática de documentos jurídicos.',
    status: 'delivered',
    priority: 'normal',
    sentAt: '2025-01-18T12:00:00Z',
    readAt: undefined,
    createdAt: '2025-01-18T12:00:00Z'
  },
  {
    id: '6',
    tenantId: '11111111-1111-1111-1111-111111111111',
    type: 'deadline',
    channel: ['whatsapp', 'email', 'push'],
    title: 'Prazo Vencido',
    message: 'ATENÇÃO: O prazo para recurso no processo 5007890-25.2023.4.03.6109 venceu ontem.',
    data: { processNumber: '5007890-25.2023.4.03.6109', deadline: '2025-01-19T23:59:59Z' },
    status: 'delivered',
    priority: 'critical',
    sentAt: '2025-01-20T00:01:00Z',
    readAt: undefined,
    createdAt: '2025-01-20T00:01:00Z'
  }
]

export default function NotificationsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedStatus, setSelectedStatus] = useState<'all' | 'unread' | 'read'>('all')
  const [selectedType, setSelectedType] = useState<'all' | NotificationType>('all')

  const filteredNotifications = mockNotifications.filter(notification => {
    // Filter by search query
    if (searchQuery && !notification.title.toLowerCase().includes(searchQuery.toLowerCase()) && 
        !notification.message.toLowerCase().includes(searchQuery.toLowerCase())) {
      return false
    }

    // Filter by status
    if (selectedStatus === 'unread' && notification.readAt) return false
    if (selectedStatus === 'read' && !notification.readAt) return false

    // Filter by type
    if (selectedType !== 'all' && notification.type !== selectedType) return false

    return true
  })

  const unreadCount = mockNotifications.filter(n => !n.readAt).length

  const getTypeIcon = (type: NotificationType) => {
    switch (type) {
      case 'process_update': return <Bell className="w-5 h-5 text-blue-500" />
      case 'deadline': return <AlertTriangle className="w-5 h-5 text-red-500" />
      case 'system': return <Settings className="w-5 h-5 text-gray-500" />
      case 'marketing': return <Info className="w-5 h-5 text-green-500" />
      default: return <Bell className="w-5 h-5" />
    }
  }

  const getTypeLabel = (type: NotificationType): string => {
    switch (type) {
      case 'process_update': return 'Processo'
      case 'deadline': return 'Prazo'
      case 'system': return 'Sistema'
      case 'marketing': return 'Novidades'
      default: return type
    }
  }

  const getPriorityColor = (priority: NotificationPriority): "default" | "secondary" | "destructive" | "outline" => {
    switch (priority) {
      case 'critical': return 'destructive'
      case 'high': return 'default'
      case 'normal': return 'secondary'
      case 'low': return 'outline'
      default: return 'outline'
    }
  }

  const getPriorityLabel = (priority: NotificationPriority): string => {
    switch (priority) {
      case 'critical': return 'Crítica'
      case 'high': return 'Alta'
      case 'normal': return 'Normal'
      case 'low': return 'Baixa'
      default: return priority
    }
  }

  const getChannelIcon = (channel: NotificationChannel) => {
    switch (channel) {
      case 'email': return <Mail className="w-4 h-4" />
      case 'whatsapp': return <MessageCircle className="w-4 h-4" />
      case 'push': return <Smartphone className="w-4 h-4" />
      case 'telegram': return <MessageCircle className="w-4 h-4" />
      case 'sms': return <Smartphone className="w-4 h-4" />
      default: return <Bell className="w-4 h-4" />
    }
  }

  const getStatusIcon = (status: NotificationStatus) => {
    switch (status) {
      case 'sent': return <Clock className="w-4 h-4 text-blue-500" />
      case 'delivered': return <CheckCircle className="w-4 h-4 text-green-500" />
      case 'failed': return <XCircle className="w-4 h-4 text-red-500" />
      default: return <Clock className="w-4 h-4" />
    }
  }

  const markAsRead = (id: string) => {
    console.log('Marking as read:', id)
    // Here you would update the notification in the backend
  }

  const markAllAsRead = () => {
    console.log('Marking all as read')
    // Here you would update all notifications in the backend
  }

  const deleteNotification = (id: string) => {
    console.log('Deleting notification:', id)
    // Here you would delete the notification in the backend
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Notificações</h1>
          <p className="text-muted-foreground">
            Gerencie todas as suas notificações e alertas
          </p>
        </div>
        <div className="flex items-center space-x-2">
          {unreadCount > 0 && (
            <Button variant="outline" onClick={markAllAsRead}>
              <MailCheck className="w-4 h-4 mr-2" />
              Marcar Todas como Lidas
            </Button>
          )}
          <Button variant="outline">
            <Settings className="w-4 h-4 mr-2" />
            Configurações
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold">{mockNotifications.length}</div>
            <p className="text-xs text-muted-foreground">Total de notificações</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold text-red-600">{unreadCount}</div>
            <p className="text-xs text-muted-foreground">Não lidas</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold text-orange-600">
              {mockNotifications.filter(n => n.priority === 'critical').length}
            </div>
            <p className="text-xs text-muted-foreground">Críticas</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-2xl font-bold text-green-600">
              {mockNotifications.filter(n => n.status === 'delivered').length}
            </div>
            <p className="text-xs text-muted-foreground">Entregues</p>
          </CardContent>
        </Card>
      </div>

      {/* Filters and Search */}
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center justify-between space-x-4">
            <div className="flex items-center space-x-2 flex-1 max-w-md">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <Input
                  placeholder="Buscar notificações..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>

            <div className="flex space-x-2">
              <select 
                className="px-3 py-2 border rounded"
                value={selectedStatus}
                onChange={(e) => setSelectedStatus(e.target.value as any)}
              >
                <option value="all">Todas</option>
                <option value="unread">Não lidas</option>
                <option value="read">Lidas</option>
              </select>

              <select 
                className="px-3 py-2 border rounded"
                value={selectedType}
                onChange={(e) => setSelectedType(e.target.value as any)}
              >
                <option value="all">Todos os tipos</option>
                <option value="process_update">Processos</option>
                <option value="deadline">Prazos</option>
                <option value="system">Sistema</option>
                <option value="marketing">Novidades</option>
              </select>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Notifications List */}
      <div className="space-y-2">
        {filteredNotifications.map((notification) => (
          <Card 
            key={notification.id} 
            className={`transition-all hover:shadow-md ${!notification.readAt ? 'border-l-4 border-l-primary bg-primary/5' : ''}`}
          >
            <CardContent className="p-4">
              <div className="flex items-start space-x-4">
                <div className="flex-shrink-0 pt-1">
                  {getTypeIcon(notification.type)}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center space-x-2 mb-1">
                        <h3 className={`font-medium ${!notification.readAt ? 'font-semibold' : ''}`}>
                          {notification.title}
                        </h3>
                        {!notification.readAt && (
                          <div className="w-2 h-2 bg-primary rounded-full"></div>
                        )}
                      </div>
                      
                      <div className="flex items-center space-x-2 mb-2">
                        <Badge variant="outline" className="text-xs">
                          {getTypeLabel(notification.type)}
                        </Badge>
                        <Badge variant={getPriorityColor(notification.priority)} className="text-xs">
                          {getPriorityLabel(notification.priority)}
                        </Badge>
                        <div className="flex items-center space-x-1">
                          {notification.channel.map((channel, index) => (
                            <div key={index} className="text-muted-foreground">
                              {getChannelIcon(channel)}
                            </div>
                          ))}
                        </div>
                        <div className="flex items-center space-x-1 text-muted-foreground">
                          {getStatusIcon(notification.status)}
                        </div>
                      </div>
                      
                      <p className="text-muted-foreground text-sm mb-2">
                        {notification.message}
                      </p>
                      
                      <div className="flex items-center space-x-4 text-xs text-muted-foreground">
                        <span className="flex items-center space-x-1">
                          <Calendar className="w-3 h-3" />
                          <span>{formatDate(notification.createdAt)}</span>
                        </span>
                        {notification.readAt && (
                          <span className="flex items-center space-x-1">
                            <Eye className="w-3 h-3" />
                            <span>Lida em {formatDate(notification.readAt)}</span>
                          </span>
                        )}
                      </div>
                    </div>
                    
                    <div className="flex items-center space-x-1">
                      {!notification.readAt && (
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => markAsRead(notification.id)}
                        >
                          <Check className="w-4 h-4" />
                        </Button>
                      )}
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <MoreVertical className="w-4 h-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          {!notification.readAt ? (
                            <DropdownMenuItem onClick={() => markAsRead(notification.id)}>
                              <Eye className="w-4 h-4 mr-2" />
                              Marcar como lida
                            </DropdownMenuItem>
                          ) : (
                            <DropdownMenuItem>
                              <EyeOff className="w-4 h-4 mr-2" />
                              Marcar como não lida
                            </DropdownMenuItem>
                          )}
                          <DropdownMenuItem>
                            <Archive className="w-4 h-4 mr-2" />
                            Arquivar
                          </DropdownMenuItem>
                          <DropdownMenuItem 
                            className="text-red-600"
                            onClick={() => deleteNotification(notification.id)}
                          >
                            <Trash className="w-4 h-4 mr-2" />
                            Excluir
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Empty State */}
      {filteredNotifications.length === 0 && (
        <Card>
          <CardContent className="p-12 text-center">
            <Bell className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-lg font-semibold mb-2">Nenhuma notificação encontrada</h3>
            <p className="text-muted-foreground">
              {searchQuery ? 'Tente ajustar os filtros de busca.' : 'Você está em dia com todas as notificações!'}
            </p>
          </CardContent>
        </Card>
      )}
    </div>
  )
}