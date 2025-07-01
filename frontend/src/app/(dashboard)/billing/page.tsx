'use client'

import { useEffect } from 'react'
import { 
  CreditCard, 
  Crown, 
  Package, 
  TrendingUp, 
  Calendar, 
  Download, 
  AlertTriangle,
  CheckCircle,
  Clock,
  DollarSign,
  Users,
  FileText,
  Zap,
  Bot,
  ArrowUp,
  ArrowDown,
  Shield
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useAuthStore, useBillingStore } from '@/store'
import { formatDate } from '@/lib/utils'
import { SubscriptionPlan, SubscriptionStatus } from '@/types'

interface Invoice {
  id: string
  number: string
  date: string
  period: string
  amount: number
  status: 'paid' | 'pending' | 'overdue' | 'cancelled'
  dueDate: string
  paidAt?: string
  downloadUrl?: string
}

interface Usage {
  processes: { used: number; limit: number }
  users: { used: number; limit: number }
  mcpCommands: { used: number; limit: number }
  aiSummaries: { used: number; limit: number }
  reports: { used: number; limit: number }
  datajudQueries: { used: number; limit: number }
}

// Remove os dados mockados - agora usaremos dados reais

const plans = [
  {
    name: 'Starter',
    price: 99,
    features: [
      '50 processos',
      '2 usuários',
      '10 resumos IA/mês',
      '10 relatórios/mês', 
      '100 consultas DataJud/dia',
      'WhatsApp incluído',
      'Busca manual ilimitada'
    ]
  },
  {
    name: 'Professional',
    price: 299,
    features: [
      '200 processos',
      '5 usuários',
      '50 resumos IA/mês',
      '100 relatórios/mês',
      '500 consultas DataJud/dia',
      'MCP Bot (200 comandos/mês)',
      'WhatsApp + Telegram',
      'Busca manual ilimitada'
    ]
  },
  {
    name: 'Business',
    price: 699,
    features: [
      '500 processos',
      '15 usuários',
      '200 resumos IA/mês',
      '500 relatórios/mês',
      '2.000 consultas DataJud/dia',
      'MCP Bot (1.000 comandos/mês)',
      'Recursos avançados',
      'Dashboard executivo'
    ]
  },
  {
    name: 'Enterprise',
    price: 1999,
    features: [
      'Processos ilimitados',
      'Usuários ilimitados',
      'IA ilimitada',
      'Relatórios ilimitados',
      '10.000 consultas DataJud/dia',
      'MCP Bot ilimitado',
      'White-label',
      'API completa',
      'Suporte prioritário'
    ]
  }
]

export default function BillingPage() {
  const { user: currentUser, tenant } = useAuthStore()
  const { 
    invoices, 
    currentUsage, 
    paymentMethod, 
    isLoading,
    loadBillingData,
    downloadInvoice 
  } = useBillingStore()

  // Load billing data on component mount
  useEffect(() => {
    loadBillingData()
  }, [loadBillingData])

  // Verificar permissões - apenas admin pode acessar billing
  if (currentUser?.role !== 'admin') {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="max-w-md">
          <CardContent className="p-6 text-center">
            <Shield className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
            <h2 className="text-lg font-semibold mb-2">Acesso Restrito</h2>
            <p className="text-muted-foreground">
              Apenas administradores podem acessar informações de billing.
            </p>
          </CardContent>
        </Card>
      </div>
    )
  }

  // Show loading state
  if (isLoading || !currentUsage) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Carregando dados de billing...</p>
        </div>
      </div>
    )
  }

  // Get plan info from tenant
  const planNames = {
    starter: 'Starter',
    professional: 'Professional', 
    business: 'Business',
    enterprise: 'Enterprise'
  }
  
  const planPrices = {
    starter: 99,
    professional: 299,
    business: 699,
    enterprise: 1999
  }

  const currentPlanName = planNames[tenant?.plan as keyof typeof planNames] || 'Starter'
  const currentPlanPrice = planPrices[tenant?.plan as keyof typeof planPrices] || 99

  const getUsagePercentage = (used: number, limit: number): number => {
    if (limit === 0 || limit === -1) return 0
    return Math.round((used / limit) * 100)
  }

  const getUsageColor = (percentage: number): string => {
    if (percentage >= 90) return 'bg-red-500'
    if (percentage >= 75) return 'bg-yellow-500'
    return 'bg-green-500'
  }

  const getInvoiceStatusColor = (status: Invoice['status']): "default" | "secondary" | "destructive" | "outline" => {
    switch (status) {
      case 'paid': return 'default'
      case 'pending': return 'secondary'
      case 'overdue': return 'destructive'
      case 'cancelled': return 'outline'
      default: return 'outline'
    }
  }

  const getInvoiceStatusIcon = (status: Invoice['status']) => {
    switch (status) {
      case 'paid': return <CheckCircle className="w-4 h-4" />
      case 'pending': return <Clock className="w-4 h-4" />
      case 'overdue': return <AlertTriangle className="w-4 h-4" />
      default: return <Clock className="w-4 h-4" />
    }
  }

  const formatCurrency = (amount: number): string => {
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'BRL'
    }).format(amount)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Billing & Assinaturas</h1>
          <p className="text-muted-foreground">
            Gerencie sua assinatura, faturas e uso dos recursos
          </p>
        </div>
      </div>

      {/* Current Plan Overview */}
      <Card className="border-2 border-primary">
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Crown className="w-6 h-6 text-yellow-500" />
              <span>Plano Atual: {currentPlanName}</span>
            </div>
            <Badge variant="default">Ativo</Badge>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-3">
            <div>
              <div className="text-2xl font-bold">{formatCurrency(currentPlanPrice)}</div>
              <p className="text-sm text-muted-foreground">por mês</p>
            </div>
            <div>
              <div className="text-sm text-muted-foreground">Próxima cobrança</div>
              <div className="font-medium">01 de Fevereiro, 2025</div>
            </div>
            <div className="flex justify-end">
              <Button>
                <TrendingUp className="w-4 h-4 mr-2" />
                Fazer Upgrade
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Billing Tabs */}
      <Tabs defaultValue="usage" className="space-y-6">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="usage">Uso Atual</TabsTrigger>
          <TabsTrigger value="plans">Planos</TabsTrigger>
          <TabsTrigger value="invoices">Faturas</TabsTrigger>
          <TabsTrigger value="payment">Pagamento</TabsTrigger>
        </TabsList>

        {/* Usage Tab */}
        <TabsContent value="usage" className="space-y-6">
          <div className="grid gap-6 md:grid-cols-2">
            {/* Usage Cards */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <FileText className="w-5 h-5" />
                  <span>Processos</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span>{currentUsage.processes.used} de {currentUsage.processes.limit === -1 ? '∞' : currentUsage.processes.limit}</span>
                    <span>{getUsagePercentage(currentUsage.processes.used, currentUsage.processes.limit)}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div 
                      className={`h-2 rounded-full ${getUsageColor(getUsagePercentage(currentUsage.processes.used, currentUsage.processes.limit))}`}
                      style={{ width: `${getUsagePercentage(currentUsage.processes.used, currentUsage.processes.limit)}%` }}
                    ></div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Users className="w-5 h-5" />
                  <span>Usuários</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span>{currentUsage.users.used} de {currentUsage.users.limit === -1 ? '∞' : currentUsage.users.limit}</span>
                    <span>{getUsagePercentage(currentUsage.users.used, currentUsage.users.limit)}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div 
                      className={`h-2 rounded-full ${getUsageColor(getUsagePercentage(currentUsage.users.used, currentUsage.users.limit))}`}
                      style={{ width: `${getUsagePercentage(currentUsage.users.used, currentUsage.users.limit)}%` }}
                    ></div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Zap className="w-5 h-5" />
                  <span>Resumos IA</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span>{currentUsage.aiSummaries.used} de {currentUsage.aiSummaries.limit === -1 ? '∞' : currentUsage.aiSummaries.limit}</span>
                    <span>{getUsagePercentage(currentUsage.aiSummaries.used, currentUsage.aiSummaries.limit)}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div 
                      className={`h-2 rounded-full ${getUsageColor(getUsagePercentage(currentUsage.aiSummaries.used, currentUsage.aiSummaries.limit))}`}
                      style={{ width: `${getUsagePercentage(currentUsage.aiSummaries.used, currentUsage.aiSummaries.limit)}%` }}
                    ></div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Bot className="w-5 h-5" />
                  <span>MCP Bot</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                {currentUsage.mcpCommands.limit === 0 ? (
                  <div className="text-center py-4">
                    <AlertTriangle className="w-8 h-8 text-yellow-500 mx-auto mb-2" />
                    <p className="text-sm text-muted-foreground">
                      Não disponível no plano {currentPlanName}
                    </p>
                    <Button size="sm" className="mt-2">
                      Fazer Upgrade
                    </Button>
                  </div>
                ) : (
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>{currentUsage.mcpCommands.used} de {currentUsage.mcpCommands.limit === -1 ? '∞' : currentUsage.mcpCommands.limit}</span>
                      <span>{getUsagePercentage(currentUsage.mcpCommands.used, currentUsage.mcpCommands.limit)}%</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div 
                        className={`h-2 rounded-full ${getUsageColor(getUsagePercentage(currentUsage.mcpCommands.used, currentUsage.mcpCommands.limit))}`}
                        style={{ width: `${getUsagePercentage(currentUsage.mcpCommands.used, currentUsage.mcpCommands.limit)}%` }}
                      ></div>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* DataJud Quota */}
          <Card>
            <CardHeader>
              <CardTitle>Quota DataJud CNJ</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span>Consultas hoje</span>
                  <span className="font-medium">{currentUsage.datajudQueries.used} / {currentUsage.datajudQueries.limit}</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div 
                    className={`h-3 rounded-full ${getUsageColor(getUsagePercentage(currentUsage.datajudQueries.used, currentUsage.datajudQueries.limit))}`}
                    style={{ width: `${getUsagePercentage(currentUsage.datajudQueries.used, currentUsage.datajudQueries.limit)}%` }}
                  ></div>
                </div>
                <p className="text-sm text-muted-foreground">
                  Quota será resetada às 00:00 (horário de Brasília)
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Plans Tab */}
        <TabsContent value="plans" className="space-y-6">
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
            {plans.map((plan) => {
              const isCurrent = plan.name.toLowerCase() === tenant?.plan
              return (
                <Card key={plan.name} className={`relative ${isCurrent ? 'border-2 border-primary' : ''}`}>
                  {isCurrent && (
                    <div className="absolute -top-3 left-1/2 transform -translate-x-1/2">
                      <Badge className="bg-primary">Plano Atual</Badge>
                    </div>
                  )}
                  <CardHeader>
                    <CardTitle className="flex items-center justify-between">
                      <span>{plan.name}</span>
                      {plan.name === 'Enterprise' && <Crown className="w-5 h-5 text-yellow-500" />}
                    </CardTitle>
                    <div className="text-2xl font-bold">
                      {formatCurrency(plan.price)}
                      <span className="text-sm font-normal text-muted-foreground">/mês</span>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <ul className="space-y-2 text-sm">
                      {plan.features.map((feature, index) => (
                        <li key={index} className="flex items-center space-x-2">
                          <CheckCircle className="w-4 h-4 text-green-500" />
                          <span>{feature}</span>
                        </li>
                      ))}
                    </ul>
                    <div className="mt-6">
                      {isCurrent ? (
                        <Button variant="outline" className="w-full" disabled>
                          Plano Atual
                        </Button>
                      ) : (
                        <Button className="w-full">
                          {plan.price > currentPlanPrice ? (
                            <>
                              <ArrowUp className="w-4 h-4 mr-2" />
                              Upgrade
                            </>
                          ) : (
                            <>
                              <ArrowDown className="w-4 h-4 mr-2" />
                              Downgrade
                            </>
                          )}
                        </Button>
                      )}
                    </div>
                  </CardContent>
                </Card>
              )
            })}
          </div>
        </TabsContent>

        {/* Invoices Tab */}
        <TabsContent value="invoices" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Histórico de Faturas</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Fatura</TableHead>
                    <TableHead>Período</TableHead>
                    <TableHead>Valor</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Vencimento</TableHead>
                    <TableHead>Ações</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {invoices.map((invoice) => (
                    <TableRow key={invoice.id}>
                      <TableCell className="font-medium">
                        {invoice.number}
                      </TableCell>
                      <TableCell>{invoice.period}</TableCell>
                      <TableCell className="font-medium">
                        {formatCurrency(invoice.amount)}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center space-x-2">
                          {getInvoiceStatusIcon(invoice.status)}
                          <Badge variant={getInvoiceStatusColor(invoice.status)}>
                            {invoice.status === 'paid' ? 'Pago' : 
                             invoice.status === 'pending' ? 'Pendente' :
                             invoice.status === 'overdue' ? 'Em atraso' : 'Cancelado'}
                          </Badge>
                        </div>
                      </TableCell>
                      <TableCell>
                        {formatDate(invoice.dueDate)}
                      </TableCell>
                      <TableCell>
                        {invoice.downloadUrl && (
                          <Button 
                            variant="outline" 
                            size="sm"
                            onClick={() => downloadInvoice(invoice.id)}
                          >
                            <Download className="w-4 h-4 mr-2" />
                            Download
                          </Button>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Payment Tab */}
        <TabsContent value="payment" className="space-y-6">
          <div className="grid gap-6 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <CreditCard className="w-5 h-5" />
                  <span>Método de Pagamento</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {paymentMethod ? (
                  <>
                    <div className="p-4 border rounded-lg">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <div className="w-10 h-6 bg-gradient-to-r from-blue-600 to-blue-800 rounded text-white text-xs flex items-center justify-center font-bold">
                            {paymentMethod.brand?.toUpperCase() || 'CARD'}
                          </div>
                          <div>
                            <div className="font-medium">•••• •••• •••• {paymentMethod.last4}</div>
                            <div className="text-sm text-muted-foreground">
                              Expira em {paymentMethod.expiryMonth?.toString().padStart(2, '0')}/{paymentMethod.expiryYear}
                            </div>
                          </div>
                        </div>
                        {paymentMethod.isDefault && <Badge variant="default">Principal</Badge>}
                      </div>
                    </div>
                    <Button variant="outline" className="w-full">
                      Alterar Método de Pagamento
                    </Button>
                  </>
                ) : (
                  <div className="text-center py-4">
                    <CreditCard className="w-8 h-8 text-muted-foreground mx-auto mb-2" />
                    <p className="text-sm text-muted-foreground mb-4">
                      Nenhum método de pagamento cadastrado
                    </p>
                    <Button className="w-full">
                      Adicionar Método de Pagamento
                    </Button>
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <DollarSign className="w-5 h-5" />
                  <span>Resumo de Cobrança</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span>Plano {currentPlanName}</span>
                    <span>{formatCurrency(currentPlanPrice)}</span>
                  </div>
                  <div className="flex justify-between text-sm text-muted-foreground">
                    <span>Próxima cobrança</span>
                    <span>01/02/2025</span>
                  </div>
                  <hr />
                  <div className="flex justify-between font-medium">
                    <span>Total</span>
                    <span>{formatCurrency(currentPlanPrice)}</span>
                  </div>
                </div>
                <p className="text-xs text-muted-foreground">
                  {paymentMethod ? (
                    `Cobrança automática no cartão terminado em ${paymentMethod.last4}`
                  ) : (
                    'Configure um método de pagamento para renovação automática'
                  )}
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Billing Address */}
          <Card>
            <CardHeader>
              <CardTitle>Endereço de Cobrança</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-2">
                <div>
                  <div className="font-medium">Silva & Associados</div>
                  <div className="text-sm text-muted-foreground">
                    Rua das Flores, 123<br />
                    Centro, São Paulo - SP<br />
                    CEP: 01234-567<br />
                    CNPJ: 12.345.678/0001-90
                  </div>
                </div>
                <div className="flex justify-end">
                  <Button variant="outline">
                    Alterar Endereço
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}