'use client'

import Link from 'next/link'
import { ArrowLeft, Shield, Eye, Database, Lock, Users, FileText, Mail } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

export default function PrivacyPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 py-8">
      <div className="container mx-auto px-4 max-w-4xl">
        {/* Header */}
        <div className="flex items-center gap-4 mb-8">
          <Button variant="outline" size="sm" asChild>
            <Link href="/terms" className="flex items-center gap-2">
              <ArrowLeft className="h-4 w-4" />
              Voltar
            </Link>
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Política de Privacidade</h1>
            <p className="text-gray-600">Direito Lux - Proteção de Dados Pessoais</p>
          </div>
        </div>

        {/* Main Content */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Shield className="h-5 w-5 text-blue-600" />
              Política de Privacidade e Proteção de Dados
            </CardTitle>
            <p className="text-sm text-gray-600">
              Em conformidade com a LGPD (Lei 13.709/2018) • Última atualização: {new Date().toLocaleDateString('pt-BR')}
            </p>
          </CardHeader>
          
          <CardContent className="space-y-8">
            {/* Introdução */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Eye className="h-5 w-5 text-blue-600" />
                Compromisso com sua Privacidade
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>
                  O Direito Lux está comprometido com a proteção da privacidade e dos dados pessoais 
                  de nossos usuários, em total conformidade com a Lei Geral de Proteção de Dados 
                  (LGPD - Lei 13.709/2018) e demais regulamentações aplicáveis.
                </p>
                <p>
                  Esta política descreve como coletamos, usamos, armazenamos e protegemos suas 
                  informações pessoais quando você utiliza nossa plataforma.
                </p>
              </div>
            </section>

            {/* Dados Coletados */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Database className="h-5 w-5 text-blue-600" />
                Dados Pessoais Coletados
              </h2>
              <div className="space-y-4 text-gray-700">
                <div>
                  <h3 className="font-semibold text-gray-900 mb-2">Dados de Identificação:</h3>
                  <ul className="list-disc list-inside space-y-1 ml-4">
                    <li>Nome completo e CPF/CNPJ</li>
                    <li>Email e número de telefone</li>
                    <li>Endereço completo</li>
                    <li>Dados de identificação profissional (OAB, etc.)</li>
                  </ul>
                </div>
                
                <div>
                  <h3 className="font-semibold text-gray-900 mb-2">Dados de Uso:</h3>
                  <ul className="list-disc list-inside space-y-1 ml-4">
                    <li>Logs de acesso e navegação</li>
                    <li>Endereço IP e informações do dispositivo</li>
                    <li>Dados de interação com a plataforma</li>
                    <li>Preferências e configurações</li>
                  </ul>
                </div>

                <div>
                  <h3 className="font-semibold text-gray-900 mb-2">Dados Jurídicos:</h3>
                  <ul className="list-disc list-inside space-y-1 ml-4">
                    <li>Informações de processos jurídicos</li>
                    <li>Documentos e anexos carregados</li>
                    <li>Dados de clientes cadastrados</li>
                    <li>Comunicações e mensagens</li>
                  </ul>
                </div>
              </div>
            </section>

            {/* Finalidades */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <FileText className="h-5 w-5 text-blue-600" />
                Finalidades do Tratamento
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>Utilizamos seus dados pessoais para:</p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li><strong>Prestação dos serviços:</strong> Gestão de processos, notificações, relatórios</li>
                  <li><strong>Autenticação e segurança:</strong> Login, controle de acesso, prevenção de fraudes</li>
                  <li><strong>Comunicação:</strong> Suporte técnico, atualizações, notificações importantes</li>
                  <li><strong>Melhoria dos serviços:</strong> Análise de uso, desenvolvimento de funcionalidades</li>
                  <li><strong>Cumprimento legal:</strong> Obrigações fiscais, trabalhistas e regulatórias</li>
                  <li><strong>Marketing (com consentimento):</strong> Ofertas, newsletters, pesquisas</li>
                </ul>
              </div>
            </section>

            {/* Base Legal */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Lock className="h-5 w-5 text-blue-600" />
                Base Legal do Tratamento
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>O tratamento dos seus dados pessoais é fundamentado nas seguintes bases legais da LGPD:</p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li><strong>Execução de contrato:</strong> Para prestação dos serviços contratados</li>
                  <li><strong>Consentimento:</strong> Para envio de comunicações de marketing</li>
                  <li><strong>Legítimo interesse:</strong> Para segurança, prevenção de fraudes e melhoria dos serviços</li>
                  <li><strong>Cumprimento de obrigação legal:</strong> Para atender exigências fiscais e regulatórias</li>
                </ul>
              </div>
            </section>

            {/* Compartilhamento */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Users className="h-5 w-5 text-blue-600" />
                Compartilhamento de Dados
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>Seus dados podem ser compartilhados apenas nas seguintes situações:</p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li><strong>Prestadores de serviços:</strong> Empresas que nos auxiliam na operação (sempre com contrato de proteção de dados)</li>
                  <li><strong>Órgãos públicos:</strong> Quando exigido por lei ou ordem judicial</li>
                  <li><strong>Integração autorizada:</strong> APIs de tribunais e órgãos públicos para consulta de processos</li>
                  <li><strong>Consentimento específico:</strong> Outras situações previamente autorizadas por você</li>
                </ul>
                <p className="font-semibold text-gray-900">
                  Nunca vendemos seus dados pessoais para terceiros.
                </p>
              </div>
            </section>

            {/* Segurança */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Shield className="h-5 w-5 text-blue-600" />
                Medidas de Segurança
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>Implementamos rigorosas medidas de segurança técnicas e organizacionais:</p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li><strong>Criptografia:</strong> Dados em trânsito e em repouso protegidos</li>
                  <li><strong>Controle de acesso:</strong> Autenticação multi-fator e princípio do menor privilégio</li>
                  <li><strong>Monitoramento:</strong> Detecção de atividades suspeitas e logs de auditoria</li>
                  <li><strong>Backup e recuperação:</strong> Cópias de segurança regulares e plano de continuidade</li>
                  <li><strong>Treinamento:</strong> Capacitação regular da equipe em proteção de dados</li>
                  <li><strong>Certificações:</strong> Compliance com padrões internacionais de segurança</li>
                </ul>
              </div>
            </section>

            {/* Direitos do Titular */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Eye className="h-5 w-5 text-blue-600" />
                Seus Direitos como Titular dos Dados
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>Conforme a LGPD, você tem os seguintes direitos:</p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li><strong>Confirmação:</strong> Saber se tratamos seus dados pessoais</li>
                  <li><strong>Acesso:</strong> Obter cópia dos seus dados pessoais</li>
                  <li><strong>Correção:</strong> Solicitar correção de dados incompletos ou inexatos</li>
                  <li><strong>Anonimização/Bloqueio:</strong> Quando desnecessários ou excessivos</li>
                  <li><strong>Eliminação:</strong> Exclusão de dados tratados com seu consentimento</li>
                  <li><strong>Portabilidade:</strong> Receber seus dados em formato estruturado</li>
                  <li><strong>Informação:</strong> Saber com quem compartilhamos seus dados</li>
                  <li><strong>Revogação:</strong> Retirar consentimento a qualquer momento</li>
                </ul>
                <p className="mt-4 font-semibold text-gray-900">
                  Para exercer seus direitos, entre em contato conosco através do email: 
                  <span className="text-blue-600"> privacidade@direitolux.com.br</span>
                </p>
              </div>
            </section>

            {/* Retenção */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Database className="h-5 w-5 text-blue-600" />
                Período de Retenção
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>Mantemos seus dados pessoais apenas pelo tempo necessário para:</p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li>Cumprir as finalidades para as quais foram coletados</li>
                  <li>Atender obrigações legais e regulatórias</li>
                  <li>Preservar direitos em eventual litígio</li>
                </ul>
                <p>
                  Após este período, os dados são eliminados de forma segura, salvo quando a 
                  manutenção for exigida por lei.
                </p>
              </div>
            </section>

            {/* Cookies */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Eye className="h-5 w-5 text-blue-600" />
                Cookies e Tecnologias Similares
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>
                  Utilizamos cookies e tecnologias similares para melhorar sua experiência, 
                  lembrar preferências e analisar o uso da plataforma.
                </p>
                <p>
                  Você pode gerenciar suas preferências de cookies através das configurações 
                  do seu navegador.
                </p>
              </div>
            </section>

            {/* Alterações */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <FileText className="h-5 w-5 text-blue-600" />
                Alterações nesta Política
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>
                  Esta política pode ser atualizada periodicamente. Alterações substanciais 
                  serão comunicadas com antecedência de 30 dias através do email cadastrado 
                  ou notificação na plataforma.
                </p>
              </div>
            </section>

            {/* Contato */}
            <section>
              <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                <Mail className="h-5 w-5 text-blue-600" />
                Contato e Encarregado de Dados
              </h2>
              <div className="space-y-3 text-gray-700">
                <p>
                  Para questões relacionadas à proteção de dados pessoais, entre em contato 
                  com nosso Encarregado de Proteção de Dados (DPO):
                </p>
                <ul className="list-disc list-inside space-y-2 ml-4">
                  <li><strong>Email:</strong> dpo@direitolux.com.br</li>
                  <li><strong>Email geral:</strong> privacidade@direitolux.com.br</li>
                  <li><strong>Telefone:</strong> (11) 99999-9999</li>
                  <li><strong>Endereço:</strong> [Endereço da empresa]</li>
                </ul>
                <p className="mt-4">
                  <strong>Autoridade Nacional de Proteção de Dados (ANPD):</strong><br />
                  Caso não seja possível resolver sua questão conosco, você pode contatar 
                  a ANPD através do site: <span className="text-blue-600">https://www.gov.br/anpd</span>
                </p>
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
                    <Link href="/terms">Termos de Uso</Link>
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