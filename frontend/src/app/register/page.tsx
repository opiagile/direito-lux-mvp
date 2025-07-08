'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import { Loader2, Building, User, CreditCard, Shield } from 'lucide-react'
import { toast } from 'sonner'

interface RegisterFormData {
  // Dados do Tenant
  tenant: {
    name: string
    document: string
    email: string
    phone: string
    website?: string
    plan: string
    // Endereço
    address: {
      street: string
      number: string
      complement?: string
      neighborhood: string
      city: string
      state: string
      zipCode: string
    }
  }
  // Dados do Usuário Admin
  user: {
    name: string
    email: string
    password: string
    confirmPassword: string
    phone: string
  }
  // Aceites
  terms: boolean
  privacy: boolean
  marketing: boolean
}

const plans = [
  {
    id: 'starter',
    name: 'Starter',
    price: 'R$ 99/mês',
    features: ['50 processos', '20 clientes', '100 consultas CNJ/dia', 'WhatsApp incluído'],
    badge: 'Ideal para começar'
  },
  {
    id: 'professional', 
    name: 'Professional',
    price: 'R$ 299/mês',
    features: ['200 processos', '100 clientes', '500 consultas CNJ/dia', 'Bot MCP incluído'],
    badge: 'Mais popular',
    popular: true
  },
  {
    id: 'business',
    name: 'Business', 
    price: 'R$ 699/mês',
    features: ['500 processos', '500 clientes', '2000 consultas CNJ/dia', 'Todas as funcionalidades'],
    badge: 'Para escritórios grandes'
  }
]

export default function RegisterPage() {
  const router = useRouter()
  const [loading, setLoading] = useState(false)
  const [step, setStep] = useState(1)
  const [formData, setFormData] = useState<RegisterFormData>({
    tenant: {
      name: '',
      document: '',
      email: '',
      phone: '',
      website: '',
      plan: 'professional',
      address: {
        street: '',
        number: '',
        complement: '',
        neighborhood: '',
        city: '',
        state: '',
        zipCode: ''
      }
    },
    user: {
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
      phone: ''
    },
    terms: false,
    privacy: false,
    marketing: false
  })

  const handleInputChange = (section: 'tenant' | 'user', field: string, value: string) => {
    if (section === 'tenant' && field.startsWith('address.')) {
      const addressField = field.replace('address.', '')
      setFormData(prev => ({
        ...prev,
        tenant: {
          ...prev.tenant,
          address: {
            ...prev.tenant.address,
            [addressField]: value
          }
        }
      }))
    } else {
      setFormData(prev => ({
        ...prev,
        [section]: {
          ...prev[section],
          [field]: value
        }
      }))
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    // Validações
    if (formData.user.password !== formData.user.confirmPassword) {
      toast.error('Senhas não conferem')
      return
    }

    if (!formData.terms || !formData.privacy) {
      toast.error('Você deve aceitar os termos e política de privacidade')
      return
    }

    setLoading(true)

    try {
      const response = await fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          tenant: formData.tenant,
          user: {
            name: formData.user.name,
            email: formData.user.email,
            password: formData.user.password,
            phone: formData.user.phone
          }
        }),
      })

      if (response.ok) {
        toast.success('Conta criada com sucesso! Verifique seu email.')
        router.push('/login?message=registration-success')
      } else {
        const data = await response.json()
        toast.error(data.message || 'Erro ao criar conta')
      }
    } catch (error) {
      toast.error('Erro de conexão. Tente novamente.')
    } finally {
      setLoading(false)
    }
  }

  const renderStep1 = () => (
    <div className="space-y-6">
      <div className="text-center">
        <Building className="mx-auto h-12 w-12 text-blue-600 mb-4" />
        <h2 className="text-2xl font-bold">Dados do Escritório</h2>
        <p className="text-gray-600">Informações do seu escritório de advocacia</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <Label htmlFor="tenant-name">Nome do Escritório *</Label>
          <Input
            id="tenant-name"
            placeholder="Silva & Associados"
            value={formData.tenant.name}
            onChange={(e) => handleInputChange('tenant', 'name', e.target.value)}
            required
          />
        </div>

        <div>
          <Label htmlFor="tenant-document">CNPJ *</Label>
          <Input
            id="tenant-document"
            placeholder="00.000.000/0001-00"
            value={formData.tenant.document}
            onChange={(e) => handleInputChange('tenant', 'document', e.target.value)}
            required
          />
        </div>

        <div>
          <Label htmlFor="tenant-email">Email do Escritório *</Label>
          <Input
            id="tenant-email"
            type="email"
            placeholder="contato@silvaassociados.com.br"
            value={formData.tenant.email}
            onChange={(e) => handleInputChange('tenant', 'email', e.target.value)}
            required
          />
        </div>

        <div>
          <Label htmlFor="tenant-phone">Telefone *</Label>
          <Input
            id="tenant-phone"
            placeholder="(11) 99999-9999"
            value={formData.tenant.phone}
            onChange={(e) => handleInputChange('tenant', 'phone', e.target.value)}
            required
          />
        </div>

        <div className="md:col-span-2">
          <Label htmlFor="tenant-website">Website (opcional)</Label>
          <Input
            id="tenant-website"
            placeholder="https://www.silvaassociados.com.br"
            value={formData.tenant.website}
            onChange={(e) => handleInputChange('tenant', 'website', e.target.value)}
          />
        </div>
      </div>

      {/* Endereço */}
      <div>
        <h3 className="text-lg font-semibold mb-4">Endereço</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <Label htmlFor="address-street">Rua *</Label>
            <Input
              id="address-street"
              placeholder="Av. Paulista"
              value={formData.tenant.address.street}
              onChange={(e) => handleInputChange('tenant', 'address.street', e.target.value)}
              required
            />
          </div>

          <div>
            <Label htmlFor="address-number">Número *</Label>
            <Input
              id="address-number"
              placeholder="1000"
              value={formData.tenant.address.number}
              onChange={(e) => handleInputChange('tenant', 'address.number', e.target.value)}
              required
            />
          </div>

          <div>
            <Label htmlFor="address-neighborhood">Bairro *</Label>
            <Input
              id="address-neighborhood"
              placeholder="Bela Vista"
              value={formData.tenant.address.neighborhood}
              onChange={(e) => handleInputChange('tenant', 'address.neighborhood', e.target.value)}
              required
            />
          </div>

          <div>
            <Label htmlFor="address-city">Cidade *</Label>
            <Input
              id="address-city"
              placeholder="São Paulo"
              value={formData.tenant.address.city}
              onChange={(e) => handleInputChange('tenant', 'address.city', e.target.value)}
              required
            />
          </div>

          <div>
            <Label htmlFor="address-state">Estado *</Label>
            <Input
              id="address-state"
              placeholder="SP"
              value={formData.tenant.address.state}
              onChange={(e) => handleInputChange('tenant', 'address.state', e.target.value)}
              required
            />
          </div>

          <div>
            <Label htmlFor="address-zipcode">CEP *</Label>
            <Input
              id="address-zipcode"
              placeholder="01310-100"
              value={formData.tenant.address.zipCode}
              onChange={(e) => handleInputChange('tenant', 'address.zipCode', e.target.value)}
              required
            />
          </div>
        </div>
      </div>

      <Button 
        onClick={() => setStep(2)}
        className="w-full"
        size="lg"
      >
        Continuar
      </Button>
    </div>
  )

  const renderStep2 = () => (
    <div className="space-y-6">
      <div className="text-center">
        <User className="mx-auto h-12 w-12 text-blue-600 mb-4" />
        <h2 className="text-2xl font-bold">Dados do Administrador</h2>
        <p className="text-gray-600">Informações do usuário administrador</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <Label htmlFor="user-name">Nome Completo *</Label>
          <Input
            id="user-name"
            placeholder="João Silva"
            value={formData.user.name}
            onChange={(e) => handleInputChange('user', 'name', e.target.value)}
            required
          />
        </div>

        <div>
          <Label htmlFor="user-email">Email *</Label>
          <Input
            id="user-email"
            type="email"
            placeholder="joao@silvaassociados.com.br"
            value={formData.user.email}
            onChange={(e) => handleInputChange('user', 'email', e.target.value)}
            required
          />
        </div>

        <div>
          <Label htmlFor="user-phone">Telefone *</Label>
          <Input
            id="user-phone"
            placeholder="(11) 99999-9999"
            value={formData.user.phone}
            onChange={(e) => handleInputChange('user', 'phone', e.target.value)}
            required
          />
        </div>

        <div>
          <Label htmlFor="user-password">Senha *</Label>
          <Input
            id="user-password"
            type="password"
            placeholder="********"
            value={formData.user.password}
            onChange={(e) => handleInputChange('user', 'password', e.target.value)}
            required
          />
        </div>

        <div className="md:col-span-2">
          <Label htmlFor="user-confirm-password">Confirmar Senha *</Label>
          <Input
            id="user-confirm-password"
            type="password"
            placeholder="********"
            value={formData.user.confirmPassword}
            onChange={(e) => handleInputChange('user', 'confirmPassword', e.target.value)}
            required
          />
        </div>
      </div>

      <div className="flex gap-4">
        <Button 
          variant="outline"
          onClick={() => setStep(1)}
          className="flex-1"
        >
          Voltar
        </Button>
        <Button 
          onClick={() => setStep(3)}
          className="flex-1"
        >
          Continuar
        </Button>
      </div>
    </div>
  )

  const renderStep3 = () => (
    <div className="space-y-6">
      <div className="text-center">
        <CreditCard className="mx-auto h-12 w-12 text-blue-600 mb-4" />
        <h2 className="text-2xl font-bold">Escolha seu Plano</h2>
        <p className="text-gray-600">Selecione o plano ideal para seu escritório</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {plans.map((plan) => (
          <Card 
            key={plan.id}
            className={`cursor-pointer transition-all ${
              formData.tenant.plan === plan.id 
                ? 'ring-2 ring-blue-500 border-blue-500' 
                : 'hover:shadow-lg'
            } ${plan.popular ? 'border-blue-200' : ''}`}
            onClick={() => handleInputChange('tenant', 'plan', plan.id)}
          >
            <CardHeader className="text-center">
              {plan.popular && (
                <Badge className="mx-auto mb-2 bg-blue-500">Mais Popular</Badge>
              )}
              <CardTitle className="text-xl">{plan.name}</CardTitle>
              <div className="text-2xl font-bold text-blue-600">{plan.price}</div>
              <p className="text-sm text-gray-500">{plan.badge}</p>
            </CardHeader>
            <CardContent>
              <ul className="space-y-2 text-sm">
                {plan.features.map((feature, index) => (
                  <li key={index} className="flex items-center">
                    <Shield className="h-4 w-4 text-green-500 mr-2" />
                    {feature}
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="space-y-4">
        <div className="flex items-center space-x-2">
          <Checkbox
            id="terms"
            checked={formData.terms}
            onCheckedChange={(checked) => setFormData(prev => ({ ...prev, terms: !!checked }))}
          />
          <Label htmlFor="terms" className="text-sm">
            Aceito os <Link href="/terms" className="text-blue-600 hover:underline">Termos de Uso</Link> *
          </Label>
        </div>

        <div className="flex items-center space-x-2">
          <Checkbox
            id="privacy"
            checked={formData.privacy}
            onCheckedChange={(checked) => setFormData(prev => ({ ...prev, privacy: !!checked }))}
          />
          <Label htmlFor="privacy" className="text-sm">
            Aceito a <Link href="/privacy" className="text-blue-600 hover:underline">Política de Privacidade</Link> *
          </Label>
        </div>

        <div className="flex items-center space-x-2">
          <Checkbox
            id="marketing"
            checked={formData.marketing}
            onCheckedChange={(checked) => setFormData(prev => ({ ...prev, marketing: !!checked }))}
          />
          <Label htmlFor="marketing" className="text-sm">
            Aceito receber comunicações de marketing (opcional)
          </Label>
        </div>
      </div>

      <div className="flex gap-4">
        <Button 
          variant="outline"
          onClick={() => setStep(2)}
          className="flex-1"
        >
          Voltar
        </Button>
        <Button 
          onClick={handleSubmit}
          disabled={loading || !formData.terms || !formData.privacy}
          className="flex-1"
        >
          {loading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Criando conta...
            </>
          ) : (
            'Criar Conta'
          )}
        </Button>
      </div>
    </div>
  )

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <Card className="w-full max-w-4xl">
        <CardHeader className="text-center">
          <CardTitle className="text-3xl font-bold text-gray-900">
            Direito Lux
          </CardTitle>
          <p className="text-gray-600">
            Crie sua conta e revolucione a gestão jurídica
          </p>
          
          {/* Progress indicator */}
          <div className="flex justify-center mt-6">
            <div className="flex space-x-4">
              {[1, 2, 3].map((stepNumber) => (
                <div key={stepNumber} className="flex items-center">
                  <div className={`
                    w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold
                    ${step >= stepNumber 
                      ? 'bg-blue-600 text-white' 
                      : 'bg-gray-200 text-gray-600'
                    }
                  `}>
                    {stepNumber}
                  </div>
                  {stepNumber < 3 && (
                    <div className={`
                      w-12 h-1 mx-2
                      ${step > stepNumber ? 'bg-blue-600' : 'bg-gray-200'}
                    `} />
                  )}
                </div>
              ))}
            </div>
          </div>
        </CardHeader>

        <CardContent className="p-8">
          <form onSubmit={handleSubmit}>
            {step === 1 && renderStep1()}
            {step === 2 && renderStep2()}
            {step === 3 && renderStep3()}
          </form>

          <div className="mt-8 text-center">
            <p className="text-sm text-gray-600">
              Já tem uma conta?{' '}
              <Link href="/login" className="text-blue-600 hover:underline font-semibold">
                Fazer login
              </Link>
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}