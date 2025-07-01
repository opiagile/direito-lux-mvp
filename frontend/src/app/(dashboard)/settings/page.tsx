'use client'

import { useState } from 'react'
import { 
  Building, 
  Settings as SettingsIcon, 
  Bell, 
  Shield, 
  Smartphone,
  Mail,
  Phone,
  MapPin,
  Save,
  Edit,
  Eye,
  EyeOff,
  Copy,
  CheckCircle,
  AlertCircle,
  Webhook,
  Key,
  Bot
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useAuthStore } from '@/store'

export default function SettingsPage() {
  const { user: currentUser, tenant } = useAuthStore()
  const [isEditing, setIsEditing] = useState(false)
  const [showApiKey, setShowApiKey] = useState(false)

  // Mock data for settings
  const [tenantData, setTenantData] = useState({
    name: 'Silva & Associados',
    cnpj: '12.345.678/0001-90',
    email: 'contato@silvaassociados.com.br',
    phone: '(11) 98765-4321',
    address: {
      street: 'Rua das Flores, 123',
      neighborhood: 'Centro',
      city: 'São Paulo',
      state: 'SP',
      zipCode: '01234-567'
    }
  })

  const [integrations, setIntegrations] = useState({
    whatsapp: {
      enabled: true,
      phone: '+5511987654321',
      apiKey: 'wa_test_key_abc123...',
      status: 'connected'
    },
    email: {
      enabled: true,
      smtpHost: 'smtp.gmail.com',
      smtpPort: 587,
      username: 'noreply@silvaassociados.com.br',
      status: 'connected'
    },
    telegram: {
      enabled: false,
      botToken: '',
      status: 'disconnected'
    },
    datajud: {
      enabled: true,
      cnpjPool: ['12.345.678/0001-90', '98.765.432/0001-10'],
      quotaLimit: 500,
      quotaUsed: 127,
      status: 'connected'
    }
  })

  const [preferences, setPreferences] = useState({
    notifications: {
      emailAlerts: true,
      whatsappAlerts: true,
      systemMaintenance: true,
      weeklyReports: true
    },
    security: {
      twoFactor: false,
      sessionTimeout: 30,
      passwordPolicy: 'strong',
      ipWhitelist: false
    }
  })

  // Verificar permissões - apenas admin e manager podem acessar
  if (!currentUser || !['admin', 'manager'].includes(currentUser.role)) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="max-w-md">
          <CardContent className="p-6 text-center">
            <Shield className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
            <h2 className="text-lg font-semibold mb-2">Acesso Restrito</h2>
            <p className="text-muted-foreground">
              Apenas administradores e gerentes podem acessar as configurações.
            </p>
          </CardContent>
        </Card>
      </div>
    )
  }

  const handleSave = () => {
    setIsEditing(false)
    // Here you would normally save to backend
    console.log('Saving settings...', { tenantData, integrations, preferences })
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'connected': return 'default'
      case 'disconnected': return 'secondary'
      case 'error': return 'destructive'
      default: return 'outline'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'connected': return <CheckCircle className="w-4 h-4" />
      case 'disconnected': return <AlertCircle className="w-4 h-4" />
      case 'error': return <AlertCircle className="w-4 h-4" />
      default: return null
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Configurações</h1>
          <p className="text-muted-foreground">
            Gerencie as configurações do seu escritório e integrações
          </p>
        </div>
        <div className="flex items-center space-x-2">
          {isEditing ? (
            <>
              <Button variant="outline" onClick={() => setIsEditing(false)}>
                Cancelar
              </Button>
              <Button onClick={handleSave}>
                <Save className="w-4 h-4 mr-2" />
                Salvar Alterações
              </Button>
            </>
          ) : (
            <Button onClick={() => setIsEditing(true)}>
              <Edit className="w-4 h-4 mr-2" />
              Editar
            </Button>
          )}
        </div>
      </div>

      {/* Settings Tabs */}
      <Tabs defaultValue="general" className="space-y-6">
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="general">Geral</TabsTrigger>
          <TabsTrigger value="integrations">Integrações</TabsTrigger>
          <TabsTrigger value="notifications">Notificações</TabsTrigger>
          <TabsTrigger value="security">Segurança</TabsTrigger>
          <TabsTrigger value="api">API</TabsTrigger>
        </TabsList>

        {/* General Settings */}
        <TabsContent value="general" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Building className="w-5 h-5" />
                <span>Dados do Escritório</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Nome do Escritório</Label>
                  <Input
                    id="name"
                    value={tenantData.name}
                    onChange={(e) => setTenantData({...tenantData, name: e.target.value})}
                    disabled={!isEditing}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="cnpj">CNPJ</Label>
                  <Input
                    id="cnpj"
                    value={tenantData.cnpj}
                    onChange={(e) => setTenantData({...tenantData, cnpj: e.target.value})}
                    disabled={!isEditing}
                  />
                </div>
              </div>
              
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="email">Email Principal</Label>
                  <Input
                    id="email"
                    type="email"
                    value={tenantData.email}
                    onChange={(e) => setTenantData({...tenantData, email: e.target.value})}
                    disabled={!isEditing}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="phone">Telefone</Label>
                  <Input
                    id="phone"
                    value={tenantData.phone}
                    onChange={(e) => setTenantData({...tenantData, phone: e.target.value})}
                    disabled={!isEditing}
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="address">Endereço Completo</Label>
                <Textarea
                  id="address"
                  value={`${tenantData.address.street}, ${tenantData.address.neighborhood}, ${tenantData.address.city} - ${tenantData.address.state}, CEP: ${tenantData.address.zipCode}`}
                  disabled={!isEditing}
                  rows={3}
                />
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Integrations */}
        <TabsContent value="integrations" className="space-y-6">
          <div className="grid gap-6 md:grid-cols-2">
            {/* WhatsApp Integration */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Smartphone className="w-5 h-5" />
                    <span>WhatsApp Business</span>
                  </div>
                  <Badge variant={getStatusColor(integrations.whatsapp.status)}>
                    {getStatusIcon(integrations.whatsapp.status)}
                    <span className="ml-1">{integrations.whatsapp.status}</span>
                  </Badge>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label>Número do WhatsApp</Label>
                  <Input value={integrations.whatsapp.phone} disabled={!isEditing} />
                </div>
                <div className="space-y-2">
                  <Label>API Key</Label>
                  <div className="flex space-x-2">
                    <Input 
                      type={showApiKey ? "text" : "password"}
                      value={integrations.whatsapp.apiKey}
                      disabled={!isEditing}
                    />
                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => setShowApiKey(!showApiKey)}
                    >
                      {showApiKey ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                    </Button>
                  </div>
                </div>
                <Button variant="outline" className="w-full">
                  <Webhook className="w-4 h-4 mr-2" />
                  Testar Conexão
                </Button>
              </CardContent>
            </Card>

            {/* Email Integration */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Mail className="w-5 h-5" />
                    <span>Email SMTP</span>
                  </div>
                  <Badge variant={getStatusColor(integrations.email.status)}>
                    {getStatusIcon(integrations.email.status)}
                    <span className="ml-1">{integrations.email.status}</span>
                  </Badge>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label>Servidor SMTP</Label>
                  <Input value={integrations.email.smtpHost} disabled={!isEditing} />
                </div>
                <div className="grid grid-cols-2 gap-2">
                  <div className="space-y-2">
                    <Label>Porta</Label>
                    <Input value={integrations.email.smtpPort} disabled={!isEditing} />
                  </div>
                  <div className="space-y-2">
                    <Label>Usuário</Label>
                    <Input value={integrations.email.username} disabled={!isEditing} />
                  </div>
                </div>
                <Button variant="outline" className="w-full">
                  <Mail className="w-4 h-4 mr-2" />
                  Testar Email
                </Button>
              </CardContent>
            </Card>

            {/* Telegram Integration */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Bot className="w-5 h-5" />
                    <span>Telegram Bot</span>
                  </div>
                  <Badge variant={getStatusColor(integrations.telegram.status)}>
                    {getStatusIcon(integrations.telegram.status)}
                    <span className="ml-1">{integrations.telegram.status}</span>
                  </Badge>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label>Bot Token</Label>
                  <Input 
                    type="password"
                    placeholder="Cole o token do seu bot aqui"
                    value={integrations.telegram.botToken}
                    disabled={!isEditing}
                  />
                </div>
                <Button variant="outline" className="w-full">
                  <Bot className="w-4 h-4 mr-2" />
                  Conectar Bot
                </Button>
              </CardContent>
            </Card>

            {/* DataJud Integration */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Building className="w-5 h-5" />
                    <span>DataJud CNJ</span>
                  </div>
                  <Badge variant={getStatusColor(integrations.datajud.status)}>
                    {getStatusIcon(integrations.datajud.status)}
                    <span className="ml-1">{integrations.datajud.status}</span>
                  </Badge>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label>Pool de CNPJs</Label>
                  <Textarea 
                    value={integrations.datajud.cnpjPool.join('\n')}
                    placeholder="Um CNPJ por linha"
                    disabled={!isEditing}
                    rows={3}
                  />
                </div>
                <div className="space-y-2">
                  <Label>Quota Utilizada</Label>
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-blue-600 h-2 rounded-full" 
                        style={{ width: `${(integrations.datajud.quotaUsed / integrations.datajud.quotaLimit) * 100}%` }}
                      ></div>
                    </div>
                    <span className="text-sm text-muted-foreground">
                      {integrations.datajud.quotaUsed}/{integrations.datajud.quotaLimit}
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Notifications */}
        <TabsContent value="notifications" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Bell className="w-5 h-5" />
                <span>Preferências de Notificação</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <Label className="text-base">Alertas por Email</Label>
                    <p className="text-sm text-muted-foreground">Receber notificações importantes por email</p>
                  </div>
                  <Button
                    variant={preferences.notifications.emailAlerts ? "default" : "outline"}
                    size="sm"
                    onClick={() => setPreferences({
                      ...preferences,
                      notifications: {
                        ...preferences.notifications,
                        emailAlerts: !preferences.notifications.emailAlerts
                      }
                    })}
                    disabled={!isEditing}
                  >
                    {preferences.notifications.emailAlerts ? "Ativado" : "Desativado"}
                  </Button>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <Label className="text-base">Alertas via WhatsApp</Label>
                    <p className="text-sm text-muted-foreground">Receber notificações urgentes via WhatsApp</p>
                  </div>
                  <Button
                    variant={preferences.notifications.whatsappAlerts ? "default" : "outline"}
                    size="sm"
                    onClick={() => setPreferences({
                      ...preferences,
                      notifications: {
                        ...preferences.notifications,
                        whatsappAlerts: !preferences.notifications.whatsappAlerts
                      }
                    })}
                    disabled={!isEditing}
                  >
                    {preferences.notifications.whatsappAlerts ? "Ativado" : "Desativado"}
                  </Button>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <Label className="text-base">Relatórios Semanais</Label>
                    <p className="text-sm text-muted-foreground">Receber resumo semanal por email</p>
                  </div>
                  <Button
                    variant={preferences.notifications.weeklyReports ? "default" : "outline"}
                    size="sm"
                    onClick={() => setPreferences({
                      ...preferences,
                      notifications: {
                        ...preferences.notifications,
                        weeklyReports: !preferences.notifications.weeklyReports
                      }
                    })}
                    disabled={!isEditing}
                  >
                    {preferences.notifications.weeklyReports ? "Ativado" : "Desativado"}
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Security */}
        <TabsContent value="security" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Shield className="w-5 h-5" />
                <span>Configurações de Segurança</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <Label className="text-base">Autenticação de Dois Fatores</Label>
                    <p className="text-sm text-muted-foreground">Adicionar camada extra de segurança</p>
                  </div>
                  <Button
                    variant={preferences.security.twoFactor ? "default" : "outline"}
                    size="sm"
                    disabled={!isEditing}
                  >
                    {preferences.security.twoFactor ? "Ativado" : "Desativado"}
                  </Button>
                </div>

                <div className="space-y-2">
                  <Label>Timeout de Sessão (minutos)</Label>
                  <Input 
                    type="number"
                    value={preferences.security.sessionTimeout}
                    disabled={!isEditing}
                  />
                </div>

                <div className="space-y-2">
                  <Label>Política de Senhas</Label>
                  <select className="w-full p-2 border rounded" disabled={!isEditing}>
                    <option value="basic">Básica (8 caracteres)</option>
                    <option value="strong" selected>Forte (12 caracteres + especiais)</option>
                    <option value="enterprise">Enterprise (16 caracteres + complexidade)</option>
                  </select>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* API Settings */}
        <TabsContent value="api" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Key className="w-5 h-5" />
                <span>API e Integrações</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-4">
                <div className="space-y-2">
                  <Label>API Key</Label>
                  <div className="flex space-x-2">
                    <Input 
                      type="password"
                      value="sk_live_abcd1234567890..."
                      disabled
                    />
                    <Button variant="outline" size="icon">
                      <Copy className="w-4 h-4" />
                    </Button>
                  </div>
                </div>

                <div className="space-y-2">
                  <Label>Webhook URL</Label>
                  <Input 
                    value="https://webhook.silvaassociados.com.br/api/webhooks"
                    disabled={!isEditing}
                  />
                </div>

                <div className="p-4 bg-muted rounded-lg">
                  <h4 className="font-medium mb-2">Limites de API</h4>
                  <div className="grid grid-cols-3 gap-4 text-sm">
                    <div>
                      <p className="text-muted-foreground">Requests/hora</p>
                      <p className="font-medium">1,000 / 1,000</p>
                    </div>
                    <div>
                      <p className="text-muted-foreground">Requests/dia</p>
                      <p className="font-medium">15,000 / 25,000</p>
                    </div>
                    <div>
                      <p className="text-muted-foreground">Rate Limit</p>
                      <p className="font-medium">50 req/min</p>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}