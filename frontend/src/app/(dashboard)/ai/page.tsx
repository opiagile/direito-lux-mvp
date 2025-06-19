'use client'

import { useState } from 'react'
import { 
  Bot, 
  Send, 
  FileText, 
  Search, 
  Download,
  Upload,
  History,
  Sparkles,
  MessageSquare
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Textarea } from '@/components/ui/textarea'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const aiTools = [
  {
    name: 'Análise de Documento',
    description: 'Análise inteligente de petições, contratos e documentos jurídicos',
    icon: FileText,
    category: 'Análise',
  },
  {
    name: 'Busca de Jurisprudência',
    description: 'Encontre precedentes relevantes com IA avançada',
    icon: Search,
    category: 'Pesquisa',
  },
  {
    name: 'Geração de Documentos',
    description: 'Crie petições, contratos e outros documentos automaticamente',
    icon: Download,
    category: 'Geração',
  },
  {
    name: 'Resumo de Processos',
    description: 'Resumos automáticos de movimentações e decisões',
    icon: MessageSquare,
    category: 'Análise',
  },
]

const chatHistory = [
  {
    id: 1,
    type: 'user',
    content: 'Analise este contrato de locação e identifique possíveis cláusulas abusivas',
    timestamp: '14:30',
  },
  {
    id: 2,
    type: 'assistant',
    content: 'Analisei o contrato de locação fornecido. Identifiquei 3 pontos que merecem atenção:\n\n1. **Cláusula de Reajuste** (Cláusula 8.2): O reajuste anual está vinculado ao IGP-M sem limitação, o que pode ser abusivo conforme o CDC.\n\n2. **Multa por Rescisão** (Cláusula 12.1): A multa de 3 aluguéis para rescisão antecipada pode ser considerada excessiva.\n\n3. **Responsabilidade por Benfeitorias** (Cláusula 15): O locatário não tem direito a indenização por benfeitorias necessárias, o que contraria o CC/2002.\n\nRecomendo revisar essas cláusulas para adequação à legislação vigente.',
    timestamp: '14:32',
  },
  {
    id: 3,
    type: 'user',
    content: 'Pode gerar uma cláusula alternativa para o item 2?',
    timestamp: '14:35',
  },
  {
    id: 4,
    type: 'assistant',
    content: 'Claro! Aqui está uma sugestão de cláusula mais equilibrada:\n\n**Cláusula 12.1 - Rescisão Antecipada (Proposta Revisada)**\n\n"Em caso de rescisão antecipada do contrato pelo LOCATÁRIO, este deverá comunicar por escrito com antecedência mínima de 30 (trinta) dias, pagando multa correspondente a 1 (um) aluguel vigente, proporcional ao tempo restante de contrato, limitada ao máximo de 2 (dois) aluguéis."\n\nEsta versão:\n- Reduz a multa de 3 para máximo de 2 aluguéis\n- Torna proporcional ao tempo restante\n- Exige aviso prévio, dando segurança ao locador\n- Está em conformidade com jurisprudência do STJ',
    timestamp: '14:36',
  },
]

export default function AIPage() {
  const [selectedTool, setSelectedTool] = useState('chat')
  const [message, setMessage] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const handleSendMessage = async () => {
    if (!message.trim()) return
    
    setIsLoading(true)
    // Simulate API call
    setTimeout(() => {
      setIsLoading(false)
      setMessage('')
    }, 2000)
  }

  const renderChatInterface = () => (
    <div className="flex flex-col h-[600px]">
      {/* Chat Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {chatHistory.map((message) => (
          <div
            key={message.id}
            className={`flex ${message.type === 'user' ? 'justify-end' : 'justify-start'}`}
          >
            <div
              className={`max-w-[80%] rounded-lg p-3 ${
                message.type === 'user'
                  ? 'bg-primary text-primary-foreground'
                  : 'bg-muted'
              }`}
            >
              <div className="flex items-start space-x-2">
                {message.type === 'assistant' && (
                  <Bot className="w-5 h-5 mt-0.5 text-blue-500" />
                )}
                <div className="flex-1">
                  <p className="text-sm whitespace-pre-wrap">
                    {message.content}
                  </p>
                  <p className="text-xs opacity-70 mt-2">
                    {message.timestamp}
                  </p>
                </div>
              </div>
            </div>
          </div>
        ))}
        {isLoading && (
          <div className="flex justify-start">
            <div className="bg-muted rounded-lg p-3">
              <div className="flex items-center space-x-2">
                <Bot className="w-5 h-5 text-blue-500" />
                <div className="flex space-x-1">
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Message Input */}
      <div className="border-t p-4">
        <div className="flex space-x-2">
          <Button variant="outline" size="icon">
            <Upload className="w-4 h-4" />
          </Button>
          <Input
            placeholder="Digite sua pergunta..."
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSendMessage()}
            disabled={isLoading}
          />
          <Button onClick={handleSendMessage} disabled={isLoading || !message.trim()}>
            <Send className="w-4 h-4" />
          </Button>
        </div>
      </div>
    </div>
  )

  const renderDocumentAnalysis = () => (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Upload de Documento</CardTitle>
          <CardDescription>
            Faça upload de um documento para análise com IA
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="border-2 border-dashed border-muted-foreground/25 rounded-lg p-8 text-center">
            <Upload className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
            <p className="text-sm text-muted-foreground mb-2">
              Arraste e solte um arquivo aqui ou clique para selecionar
            </p>
            <p className="text-xs text-muted-foreground mb-4">
              PDF, DOC, DOCX até 10MB
            </p>
            <Button>Selecionar Arquivo</Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Tipos de Análise</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2">
            <div className="p-4 border rounded-lg cursor-pointer hover:bg-accent">
              <h4 className="font-medium">Análise de Cláusulas</h4>
              <p className="text-sm text-muted-foreground">
                Identifica cláusulas potencialmente abusivas
              </p>
            </div>
            <div className="p-4 border rounded-lg cursor-pointer hover:bg-accent">
              <h4 className="font-medium">Extração de Dados</h4>
              <p className="text-sm text-muted-foreground">
                Extrai informações estruturadas do documento
              </p>
            </div>
            <div className="p-4 border rounded-lg cursor-pointer hover:bg-accent">
              <h4 className="font-medium">Resumo Executivo</h4>
              <p className="text-sm text-muted-foreground">
                Cria um resumo dos pontos principais
              </p>
            </div>
            <div className="p-4 border rounded-lg cursor-pointer hover:bg-accent">
              <h4 className="font-medium">Análise de Riscos</h4>
              <p className="text-sm text-muted-foreground">
                Avalia riscos legais e contratuais
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )

  const renderJurisprudenceSearch = () => (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Busca Inteligente de Jurisprudência</CardTitle>
          <CardDescription>
            Use IA para encontrar precedentes relevantes
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <Textarea
              placeholder="Descreva seu caso ou questão jurídica..."
              className="min-h-[100px]"
            />
            <div className="flex space-x-2">
              <Button>
                <Search className="w-4 h-4 mr-2" />
                Buscar Jurisprudência
              </Button>
              <Button variant="outline">
                Busca Avançada
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Resultados Recentes</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="p-4 border rounded-lg">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h4 className="font-medium">STJ - REsp 1.234.567</h4>
                  <p className="text-sm text-muted-foreground mt-1">
                    Cobrança de taxa de corretagem em financiamento imobiliário...
                  </p>
                  <div className="flex items-center space-x-2 mt-2">
                    <Badge variant="secondary">95% similaridade</Badge>
                    <Badge variant="outline">Direito Civil</Badge>
                  </div>
                </div>
                <Button variant="ghost" size="sm">
                  <Download className="w-4 h-4" />
                </Button>
              </div>
            </div>

            <div className="p-4 border rounded-lg">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h4 className="font-medium">TJSP - Apelação 5555.666</h4>
                  <p className="text-sm text-muted-foreground mt-1">
                    Rescisão contratual por inadimplemento do devedor...
                  </p>
                  <div className="flex items-center space-x-2 mt-2">
                    <Badge variant="secondary">87% similaridade</Badge>
                    <Badge variant="outline">Contratos</Badge>
                  </div>
                </div>
                <Button variant="ghost" size="sm">
                  <Download className="w-4 h-4" />
                </Button>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">IA Assistant</h1>
          <p className="text-muted-foreground">
            Assistente jurídico inteligente com IA avançada
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Badge variant="secondary" className="bg-blue-100 text-blue-800">
            <Sparkles className="w-3 h-3 mr-1" />
            GPT-4 Powered
          </Badge>
        </div>
      </div>

      {/* AI Tools Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {aiTools.map((tool) => (
          <Card key={tool.name} className="cursor-pointer hover:shadow-md transition-shadow">
            <CardContent className="p-4">
              <div className="flex items-start space-x-3">
                <div className="p-2 bg-blue-100 rounded-lg">
                  <tool.icon className="w-5 h-5 text-blue-600" />
                </div>
                <div className="flex-1">
                  <h3 className="font-medium text-sm">{tool.name}</h3>
                  <p className="text-xs text-muted-foreground mt-1">
                    {tool.description}
                  </p>
                  <Badge variant="outline" className="mt-2 text-xs">
                    {tool.category}
                  </Badge>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Main Interface */}
      <Tabs value={selectedTool} onValueChange={setSelectedTool}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="chat">Chat IA</TabsTrigger>
          <TabsTrigger value="analysis">Análise de Docs</TabsTrigger>
          <TabsTrigger value="jurisprudence">Jurisprudência</TabsTrigger>
          <TabsTrigger value="history">Histórico</TabsTrigger>
        </TabsList>

        <TabsContent value="chat" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <MessageSquare className="w-5 h-5 mr-2" />
                Chat com IA
              </CardTitle>
              <CardDescription>
                Converse com o assistente jurídico inteligente
              </CardDescription>
            </CardHeader>
            <CardContent>
              {renderChatInterface()}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="analysis" className="mt-6">
          {renderDocumentAnalysis()}
        </TabsContent>

        <TabsContent value="jurisprudence" className="mt-6">
          {renderJurisprudenceSearch()}
        </TabsContent>

        <TabsContent value="history" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <History className="w-5 h-5 mr-2" />
                Histórico de Interações
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-muted-foreground">
                Seu histórico de interações com a IA aparecerá aqui.
              </p>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}