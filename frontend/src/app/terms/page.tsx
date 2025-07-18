'use client'

import Link from 'next/link'
import { ArrowLeft, Scale, Shield, Users, FileText, Clock, Mail } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

export default function TermsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 py-8">
      <div className="container mx-auto px-4 max-w-4xl">
        {/* Header */}
        <div className="flex items-center gap-4 mb-8">
          <Button variant="outline" size="sm" asChild>
            <Link href="/register" className="flex items-center gap-2">
              <ArrowLeft className="h-4 w-4" />
              Voltar
            </Link>
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Termos de Uso</h1>
            <p className="text-gray-600">Direito Lux - Plataforma de Gestão Jurídica</p>
          </div>
        </div>

        {/* Main Content */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Scale className="h-5 w-5 text-blue-600" />
              Termos e Condições de Uso
            </CardTitle>
            <p className="text-sm text-gray-600">
              Última atualização: {new Date().toLocaleDateString('pt-BR')}
            </p>
          </CardHeader>
          
          <CardContent className="space-y-8">
            {/* Seção 1 - Aceitação */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <FileText className="h-5 w-5 text-blue-600" />
                1. Aceitação dos Termos
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>
                  Ao acessar e utilizar a plataforma Direito Lux, você concorda em cumprir e estar 
                  vinculado aos termos e condições de uso descritos neste documento.
                </p>
                <p>
                  Se você não concordar com qualquer parte destes termos, não deve usar nossos serviços.
                </p>
              </div>
            </section>

            {/* Seção 2 - Descrição dos Serviços */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Users className="h-5 w-5 text-blue-600" />
                2. Descrição dos Serviços
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>
                  O Direito Lux é uma plataforma SaaS (Software as a Service) para gestão de 
                  processos jurídicos que oferece:
                </p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li>Gestão e monitoramento de processos jurídicos</li>
                  <li>Integração com APIs do sistema judiciário brasileiro</li>
                  <li>Notificações automáticas via WhatsApp, email e Telegram</li>
                  <li>Análise de dados com inteligência artificial</li>
                  <li>Relatórios e dashboards personalizados</li>
                  <li>Gestão de clientes e documentos</li>
                </ul>
              </div>
            </section>

            {/* Seção 3 - Responsabilidades do Usuário */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Shield className="h-5 w-5 text-blue-600" />
                3. Responsabilidades do Usuário
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>Você se compromete a:</p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li>Fornecer informações verdadeiras, precisas e atualizadas</li>
                  <li>Manter a confidencialidade de suas credenciais de acesso</li>
                  <li>Usar a plataforma apenas para fins legais e éticos</li>
                  <li>Respeitar os direitos de propriedade intelectual</li>
                  <li>Não tentar acessar sistemas ou dados não autorizados</li>
                  <li>Cumprir todas as leis aplicáveis, incluindo o Marco Civil da Internet</li>
                </ul>
              </div>
            </section>

            {/* Seção 4 - Proteção de Dados */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Shield className="h-5 w-5 text-blue-600" />
                4. Proteção de Dados e Privacidade
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>
                  Estamos comprometidos com a proteção dos seus dados pessoais em conformidade 
                  com a Lei Geral de Proteção de Dados (LGPD - Lei 13.709/2018).
                </p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li>Coletamos apenas dados necessários para a prestação dos serviços</li>
                  <li>Implementamos medidas de segurança técnicas e organizacionais</li>
                  <li>Não compartilhamos dados com terceiros sem consentimento</li>
                  <li>Você pode solicitar acesso, correção ou exclusão dos seus dados</li>
                </ul>
                <p>
                  Para mais detalhes, consulte nossa{' '}
                  <Link href="/privacy" className="text-blue-600 hover:underline">
                    Política de Privacidade
                  </Link>.
                </p>
              </div>
            </section>

            {/* Seção 5 - Planos e Pagamentos */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Clock className="h-5 w-5 text-blue-600" />
                5. Planos e Pagamentos
              </h2>
              <div className="space-y-3 text-gray-700">
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li>Os planos são cobrados mensalmente ou anualmente conforme contratado</li>
                  <li>Oferecemos período de teste gratuito de 15 dias</li>
                  <li>O cancelamento pode ser feito a qualquer momento</li>
                  <li>Reembolsos seguem nossa política específica</li>
                  <li>Preços podem ser alterados com aviso prévio de 30 dias</li>
                </ul>
              </div>
            </section>

            {/* Seção 6 - Limitações */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <FileText className="h-5 w-5 text-blue-600" />
                6. Limitações de Responsabilidade
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>
                  O Direito Lux não se responsabiliza por:
                </p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li>Interrupções temporárias dos serviços para manutenção</li>
                  <li>Indisponibilidade de APIs externas (tribunais, governo)</li>
                  <li>Decisões tomadas com base nas informações da plataforma</li>
                  <li>Uso inadequado da plataforma pelos usuários</li>
                </ul>
              </div>
            </section>

            {/* Seção 7 - Alterações */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Clock className="h-5 w-5 text-blue-600" />
                7. Alterações nos Termos
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>
                  Reservamo-nos o direito de modificar estes termos a qualquer momento. 
                  As alterações serão comunicadas com antecedência de 30 dias e entrarão 
                  em vigor na data especificada.
                </p>
                <p>
                  O uso continuado da plataforma após as alterações constitui aceitação 
                  dos novos termos.
                </p>
              </div>
            </section>

            {/* Seção 8 - Contato */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Mail className="h-5 w-5 text-blue-600" />
                8. Contato
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>
                  Para dúvidas sobre estes termos ou nossos serviços, entre em contato:
                </p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li>Email: contato@direitolux.com.br</li>
                  <li>Telefone: (11) 99999-9999</li>
                  <li>Website: https://direitolux.com.br</li>
                </ul>
              </div>
            </section>

            {/* Footer */}
            <div className="pt-8 border-t">
              <div className="flex flex-col sm:flex-row justify-between items-center gap-4">
                <p className="text-sm text-gray-500">
                  © 2025 Direito Lux. Todos os direitos reservados.
                </p>
                <div className="flex gap-4">
                  <Button variant="outline" size="sm" asChild>
                    <Link href="/privacy">Política de Privacidade</Link>
                  </Button>
                  <Button variant="outline" size="sm" asChild>
                    <Link href="/register">Aceitar e Continuar</Link>
                  </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}