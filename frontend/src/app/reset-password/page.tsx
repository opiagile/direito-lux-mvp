'use client'

import { useState, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Loader2, Lock, CheckCircle, AlertCircle, Eye, EyeOff } from 'lucide-react'
import { toast } from 'sonner'

export default function ResetPasswordPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [loading, setLoading] = useState(false)
  const [validatingToken, setValidatingToken] = useState(true)
  const [tokenValid, setTokenValid] = useState(false)
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [success, setSuccess] = useState(false)
  
  const token = searchParams.get('token')

  useEffect(() => {
    if (!token) {
      setValidatingToken(false)
      setTokenValid(false)
      return
    }

    // Validar token na API
    validateToken(token)
  }, [token])

  const validateToken = async (tokenToValidate: string) => {
    try {
      const response = await fetch(`/api/v1/auth/verify-reset-token?token=${tokenToValidate}`)
      setTokenValid(response.ok)
    } catch (error) {
      setTokenValid(false)
    } finally {
      setValidatingToken(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!password || !confirmPassword) {
      toast.error('Por favor, preencha todos os campos')
      return
    }

    if (password !== confirmPassword) {
      toast.error('As senhas não conferem')
      return
    }

    if (password.length < 8) {
      toast.error('A senha deve ter pelo menos 8 caracteres')
      return
    }

    setLoading(true)

    try {
      const response = await fetch('/api/v1/auth/reset-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          token,
          password,
        }),
      })

      if (response.ok) {
        setSuccess(true)
        toast.success('Senha alterada com sucesso!')
      } else {
        const data = await response.json()
        toast.error(data.message || 'Erro ao alterar senha')
      }
    } catch (error) {
      toast.error('Erro de conexão. Tente novamente.')
    } finally {
      setLoading(false)
    }
  }

  // Loading state
  if (validatingToken) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
        <Card className="w-full max-w-md">
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-blue-600 mb-4" />
            <p className="text-gray-600">Validando link de recuperação...</p>
          </CardContent>
        </Card>
      </div>
    )
  }

  // Invalid token
  if (!tokenValid) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <AlertCircle className="mx-auto h-12 w-12 text-red-600 mb-4" />
            <CardTitle className="text-2xl font-bold text-gray-900">
              Link Inválido
            </CardTitle>
          </CardHeader>

          <CardContent className="space-y-6">
            <div className="text-center space-y-4">
              <p className="text-gray-600">
                Este link de recuperação de senha é inválido ou expirou.
              </p>
              
              <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                <p className="text-sm text-red-800">
                  <strong>Possíveis motivos:</strong>
                  <br />
                  • O link já foi utilizado
                  <br />
                  • O link expirou (válido por 1 hora)
                  <br />
                  • O link está incompleto ou foi modificado
                </p>
              </div>
            </div>

            <div className="space-y-4">
              <Button
                onClick={() => router.push('/forgot-password')}
                className="w-full"
              >
                Solicitar Novo Link
              </Button>

              <Button
                onClick={() => router.push('/login')}
                variant="outline"
                className="w-full"
              >
                Voltar ao Login
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  // Success state
  if (success) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <CheckCircle className="mx-auto h-12 w-12 text-green-600 mb-4" />
            <CardTitle className="text-2xl font-bold text-gray-900">
              Senha Alterada!
            </CardTitle>
          </CardHeader>

          <CardContent className="space-y-6">
            <div className="text-center space-y-4">
              <p className="text-gray-600">
                Sua senha foi alterada com sucesso. Agora você pode fazer login com sua nova senha.
              </p>
              
              <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                <p className="text-sm text-green-800">
                  <strong>Próximos passos:</strong>
                  <br />
                  • Faça login com sua nova senha
                  <br />
                  • Mantenha sua senha em local seguro
                  <br />
                  • Use uma senha única e forte
                </p>
              </div>
            </div>

            <Button
              onClick={() => router.push('/login')}
              className="w-full"
              size="lg"
            >
              Fazer Login
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  // Reset password form
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <Lock className="mx-auto h-12 w-12 text-blue-600 mb-4" />
          <CardTitle className="text-2xl font-bold text-gray-900">
            Nova Senha
          </CardTitle>
          <p className="text-gray-600">
            Digite sua nova senha abaixo
          </p>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="password">Nova Senha</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  placeholder="Digite sua nova senha"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  autoFocus
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? (
                    <EyeOff className="h-4 w-4 text-gray-400" />
                  ) : (
                    <Eye className="h-4 w-4 text-gray-400" />
                  )}
                </Button>
              </div>
              <p className="text-xs text-gray-500">
                Mínimo de 8 caracteres
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="confirmPassword">Confirmar Nova Senha</Label>
              <div className="relative">
                <Input
                  id="confirmPassword"
                  type={showConfirmPassword ? 'text' : 'password'}
                  placeholder="Confirme sua nova senha"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  required
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                >
                  {showConfirmPassword ? (
                    <EyeOff className="h-4 w-4 text-gray-400" />
                  ) : (
                    <Eye className="h-4 w-4 text-gray-400" />
                  )}
                </Button>
              </div>
            </div>

            {/* Password strength indicator */}
            {password && (
              <div className="space-y-2">
                <div className="text-xs text-gray-600">Força da senha:</div>
                <div className="flex gap-1">
                  <div className={`h-1 flex-1 rounded ${password.length >= 8 ? 'bg-green-500' : 'bg-gray-200'}`} />
                  <div className={`h-1 flex-1 rounded ${password.length >= 12 ? 'bg-green-500' : 'bg-gray-200'}`} />
                  <div className={`h-1 flex-1 rounded ${/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(password) ? 'bg-green-500' : 'bg-gray-200'}`} />
                  <div className={`h-1 flex-1 rounded ${/(?=.*[!@#$%^&*])/.test(password) ? 'bg-green-500' : 'bg-gray-200'}`} />
                </div>
                <div className="text-xs text-gray-500">
                  <div className={password.length >= 8 ? 'text-green-600' : ''}>
                    ✓ Pelo menos 8 caracteres
                  </div>
                  <div className={/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(password) ? 'text-green-600' : ''}>
                    ✓ Letras maiúsculas, minúsculas e números
                  </div>
                  <div className={/(?=.*[!@#$%^&*])/.test(password) ? 'text-green-600' : ''}>
                    ✓ Caracteres especiais (!@#$%^&*)
                  </div>
                </div>
              </div>
            )}

            <Button
              type="submit"
              disabled={loading || password !== confirmPassword || password.length < 8}
              className="w-full"
              size="lg"
            >
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Alterando...
                </>
              ) : (
                <>
                  <Lock className="mr-2 h-4 w-4" />
                  Alterar Senha
                </>
              )}
            </Button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600">
              Lembrou da senha?{' '}
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