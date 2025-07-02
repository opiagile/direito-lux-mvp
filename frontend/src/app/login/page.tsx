'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Eye, EyeOff, Lock, Mail, Building } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useLogin } from '@/hooks/api'
import { useAuthStore } from '@/store'
import { toast } from 'sonner'
import { Tenant } from '@/types'

const loginSchema = z.object({
  email: z.string().email('Email inválido'),
  password: z.string().min(6, 'Senha deve ter pelo menos 6 caracteres'),
})

type LoginForm = z.infer<typeof loginSchema>

export default function LoginPage() {
  const router = useRouter()
  const [showPassword, setShowPassword] = useState(false)
  const login = useLogin()
  const { login: loginStore } = useAuthStore()

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
  })

  const onSubmit = async (data: LoginForm) => {
    try {
      console.log('🔐 Autenticando com banco real:', data.email)
      
      // Make direct API call to auth service - NO MOCKS
      const result = await login.mutateAsync(data)
      console.log('✅ Resposta do auth-service:', result)
      
      // Validate response structure from real API
      if (!result.user || !result.access_token) {
        console.error('❌ Resposta inválida do auth-service:', result)
        toast.error('Erro na resposta do servidor de autenticação')
        return
      }
      
      const { user, access_token } = result
      console.log('👤 Usuário autenticado:', user.email, user.role)
      console.log('🔍 Estrutura completa do user:', JSON.stringify(user, null, 2))
      console.log('🔍 Tenant ID do usuário:', user.tenant_id || user.tenantId || user.tenant_id)
      
      // Store token for authenticated requests
      localStorage.setItem('auth_token', access_token)
      
      // Extract tenant ID from user response
      const tenantId = user.tenant_id
      if (!tenantId) {
        console.error('❌ Tenant ID não encontrado na resposta do usuário')
        console.error('🔍 Dados do usuário recebidos:', user)
        toast.error('❌ Dados incompletos do usuário. Contate o administrador.')
        return
      }
      
      console.log('✅ Tenant ID extraído:', tenantId)
      
      // Fetch tenant data from tenant-service - REQUIRED ONLINE
      let tenant: Tenant
      try {
        console.log('🏢 Buscando dados do tenant:', tenantId)
        const tenantResponse = await fetch(`http://localhost:8082/api/v1/tenants/${tenantId}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${access_token}`,
            'X-Tenant-ID': tenantId,
            'Content-Type': 'application/json'
          },
          signal: AbortSignal.timeout(10000)  // 10 second timeout
        })
        
        if (tenantResponse.ok) {
          const tenantData = await tenantResponse.json()
          tenant = tenantData.data || tenantData
          console.log('✅ Dados do tenant recuperados:', tenant.name, tenant.plan)
        } else if (tenantResponse.status === 404) {
          console.error('❌ Tenant não encontrado no sistema')
          toast.error('❌ Dados do escritório não encontrados. Contate o administrador.')
          return
        } else {
          console.error('❌ Erro no tenant-service:', tenantResponse.status, tenantResponse.statusText)
          toast.error(`❌ Erro no serviço do escritório (${tenantResponse.status}). Contate o suporte.`)
          return
        }
      } catch (tenantError: any) {
        console.error('❌ Erro crítico ao buscar tenant:', tenantError)
        
        // NO FALLBACK - System must be online
        if (tenantError.name === 'AbortError') {
          toast.error('❌ Timeout ao conectar com serviços. Verifique sua conexão.')
        } else if (tenantError.message?.includes('fetch') || tenantError.code === 'ECONNREFUSED') {
          toast.error('❌ Serviços indisponíveis. Contate o administrador do sistema.')
        } else {
          toast.error('❌ Erro ao carregar dados do escritório. Contate o suporte.')
        }
        return
      }
      
      // Store authentication data
      loginStore(user, tenant, access_token)
      console.log('🚀 Login completo, redirecionando...')
      toast.success(`Bem-vindo, ${user.first_name} ${user.last_name} (${user.role})`)
      router.push('/dashboard')
      
    } catch (error: any) {
      console.error('❌ Erro de autenticação:', error)
      
      // Handle specific error cases with clear user messages
      if (error.response?.status === 400) {
        const errorMsg = error.response?.data?.message || error.response?.data?.error
        if (errorMsg?.includes('credenciais inválidas') || errorMsg?.includes('senha')) {
          toast.error('❌ Email ou senha incorretos. Verifique suas credenciais.')
        } else if (errorMsg?.includes('usuário não encontrado') || errorMsg?.includes('not found')) {
          toast.error(`❌ Email "${data.email}" não encontrado. Verifique o email digitado.`)
        } else {
          toast.error(`❌ Erro: ${errorMsg}`)
        }
        console.error('Detalhes do erro 400:', error.response?.data)
      } else if (error.response?.status === 401) {
        toast.error('❌ Email ou senha incorretos. Verifique suas credenciais.')
      } else if (error.response?.status === 403) {
        toast.error('❌ Usuário inativo ou sem permissão de acesso.')
      } else if (error.response?.status === 404) {
        toast.error(`❌ Email "${data.email}" não encontrado no sistema.`)
      } else if (error.code === 'ECONNREFUSED' || error.message.includes('Network Error')) {
        toast.error('🔌 Serviço de autenticação indisponível. Contate o suporte.')
      } else if (error.message?.includes('não está cadastrado no sistema')) {
        toast.error(`❌ Email "${data.email}" não está cadastrado. Entre em contato com o administrador.`)
      } else {
        // Generic error - show user-friendly message
        toast.error('❌ Erro ao fazer login. Verifique suas credenciais ou contate o suporte.')
        console.error('Erro completo:', error)
      }
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800 px-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <div className="mx-auto w-24 h-24 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-2xl flex items-center justify-center mb-4">
            <Building className="w-12 h-12 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            Direito Lux
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Gestão Jurídica Inteligente
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Entrar na sua conta</CardTitle>
            <CardDescription>
              Digite suas credenciais para acessar a plataforma
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <div className="relative">
                  <Mail className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
                  <Input
                    id="email"
                    type="email"
                    placeholder="seu@email.com"
                    className="pl-10"
                    {...register('email')}
                  />
                </div>
                {errors.email && (
                  <p className="text-sm text-red-600">{errors.email.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="password">Senha</Label>
                <div className="relative">
                  <Lock className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
                  <Input
                    id="password"
                    type={showPassword ? 'text' : 'password'}
                    placeholder="Sua senha"
                    className="pl-10 pr-10"
                    {...register('password')}
                  />
                  <button
                    type="button"
                    className="absolute right-3 top-3 h-4 w-4 text-gray-400 hover:text-gray-600"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? <EyeOff /> : <Eye />}
                  </button>
                </div>
                {errors.password && (
                  <p className="text-sm text-red-600">{errors.password.message}</p>
                )}
              </div>

              <Button
                type="submit"
                className="w-full"
                disabled={login.isPending}
              >
                {login.isPending ? 'Entrando...' : 'Entrar'}
              </Button>
            </form>

            <div className="mt-6 text-center">
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Não tem uma conta?{' '}
                <button
                  onClick={() => router.push('/register')}
                  className="text-blue-600 hover:underline font-medium"
                >
                  Criar conta
                </button>
              </p>
            </div>
          </CardContent>
        </Card>

        <div className="mt-8 text-center text-xs text-gray-500">
          © 2025 Direito Lux. Todos os direitos reservados.
        </div>
      </div>
    </div>
  )
}