'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Loader2, Mail, ArrowLeft, CheckCircle } from 'lucide-react'
import { toast } from 'sonner'

export default function ForgotPasswordPage() {
  const router = useRouter()
  const [loading, setLoading] = useState(false)
  const [email, setEmail] = useState('')
  const [sent, setSent] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!email) {
      toast.error('Por favor, informe seu email')
      return
    }

    setLoading(true)

    try {
      const response = await fetch('/api/v1/auth/forgot-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      })

      if (response.ok) {
        setSent(true)
        toast.success('Instruções enviadas para seu email!')
      } else {
        const data = await response.json()
        toast.error(data.message || 'Erro ao enviar instruções')
      }
    } catch (error) {
      toast.error('Erro de conexão. Tente novamente.')
    } finally {
      setLoading(false)
    }
  }

  if (sent) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <CheckCircle className="mx-auto h-12 w-12 text-green-600 mb-4" />
            <CardTitle className="text-2xl font-bold text-gray-900">
              Email Enviado!
            </CardTitle>
          </CardHeader>

          <CardContent className="space-y-6">
            <div className="text-center space-y-4">
              <p className="text-gray-600">
                Se o email <strong>{email}</strong> estiver cadastrado em nosso sistema, 
                você receberá instruções para recuperar sua senha.
              </p>
              
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <p className="text-sm text-blue-800">
                  <strong>Não recebeu o email?</strong>
                  <br />
                  • Verifique sua caixa de spam
                  <br />
                  • O email pode levar alguns minutos para chegar
                  <br />
                  • Verifique se o endereço está correto
                </p>
              </div>
            </div>

            <div className="space-y-4">
              <Button
                onClick={() => {
                  setSent(false)
                  setEmail('')
                }}
                variant="outline"
                className="w-full"
              >
                Tentar outro email
              </Button>

              <Button
                onClick={() => router.push('/login')}
                className="w-full"
              >
                Voltar ao Login
              </Button>
            </div>

            <div className="text-center">
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

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <Mail className="mx-auto h-12 w-12 text-blue-600 mb-4" />
          <CardTitle className="text-2xl font-bold text-gray-900">
            Esqueceu sua senha?
          </CardTitle>
          <p className="text-gray-600">
            Não se preocupe! Digite seu email e enviaremos instruções para recuperar sua senha.
          </p>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="seu@email.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                autoFocus
              />
            </div>

            <Button
              type="submit"
              disabled={loading}
              className="w-full"
              size="lg"
            >
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Enviando...
                </>
              ) : (
                <>
                  <Mail className="mr-2 h-4 w-4" />
                  Enviar Instruções
                </>
              )}
            </Button>

            <Button
              type="button"
              variant="outline"
              onClick={() => router.push('/login')}
              className="w-full"
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Voltar ao Login
            </Button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600">
              Não tem uma conta?{' '}
              <Link href="/register" className="text-blue-600 hover:underline font-semibold">
                Criar conta
              </Link>
            </p>
          </div>

          <div className="mt-6 text-center">
            <div className="bg-gray-50 border border-gray-200 rounded-lg p-4">
              <p className="text-xs text-gray-500">
                <strong>Problemas para acessar?</strong>
                <br />
                Entre em contato conosco em{' '}
                <a 
                  href="mailto:suporte@direitolux.com.br" 
                  className="text-blue-600 hover:underline"
                >
                  suporte@direitolux.com.br
                </a>
                <br />
                ou pelo WhatsApp: (11) 99999-9999
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}