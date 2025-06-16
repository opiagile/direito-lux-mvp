# Ubiquitous Language - Direito Lux

## Core Domain Concepts

### Process Management Domain

#### **Processo (Process)**
- **Definição**: Ação judicial em curso nos tribunais brasileiros identificada por número CNJ
- **Identificador**: Número sequencial único de 20 dígitos (NNNNNNN-DD.AAAA.J.TR.OOOO)
- **Estados**: Ativo, Arquivado, Baixado, Suspenso
- **Exemplo**: "Um processo de execução fiscal está sendo monitorado desde sua distribuição"

#### **Número CNJ (CNJ Number)**
- **Definição**: Identificador único nacional padronizado pelo Conselho Nacional de Justiça
- **Formato**: NNNNNNN-DD.AAAA.J.TR.OOOO
  - NNNNNNN: Número sequencial
  - DD: Dígito verificador
  - AAAA: Ano de ajuizamento  
  - J: Segmento do Poder Judiciário
  - TR: Tribunal
  - OOOO: Órgão julgador
- **Exemplo**: "1234567-89.2023.8.26.0001"

#### **Movimentação (Movement)**
- **Definição**: Qualquer ato ou evento registrado nos autos processuais
- **Tipos**: Decisão, Despacho, Petição, Audiência, Juntada, Expedição
- **Características**: Data, descrição, responsável, documentos anexos
- **Exemplo**: "Nova movimentação: Sentença de procedência publicada em 15/03/2024"

#### **Prazo (Deadline)**
- **Definição**: Data limite para realização de ato processual
- **Tipos**: Recursal, Contestação, Manifestação, Cumprimento
- **Prioridade**: Baixa, Média, Alta, Crítica
- **Estado**: Pendente, Cumprido, Perdido, Prorrogado
- **Exemplo**: "Prazo para recurso de apelação: 15 dias úteis a partir da intimação"

#### **Monitoramento (Monitoring)**
- **Definição**: Acompanhamento automático de mudanças em processo judicial
- **Estados**: Ativo, Pausado, Parado, Suspenso
- **Frequência**: Tempo real, Diária, Semanal
- **Exemplo**: "Monitoramento ativo verificando movimentações 6x por dia"

#### **Tribunal (Court)**
- **Definição**: Órgão do Poder Judiciário responsável pelo julgamento
- **Tipos**: TJ (Tribunal de Justiça), TRF (Tribunal Regional Federal), TST (Tribunal Superior do Trabalho)
- **Instância**: 1º Grau, 2º Grau, Superior
- **Exemplo**: "TJSP - Tribunal de Justiça de São Paulo, 2ª Instância"

#### **Classe Processual (Process Class)**
- **Definição**: Categoria que define o tipo de ação judicial
- **Exemplos**: Execução Fiscal, Ação de Cobrança, Mandado de Segurança
- **Código**: Padronizado pelo CNJ (tabela unificada)

#### **Assunto Processual (Process Subject)**
- **Definição**: Matéria jurídica específica do processo
- **Área**: Civil, Criminal, Trabalhista, Tributária
- **Hierarquia**: Assunto principal e secundários
- **Exemplo**: "Direito Tributário > IPTU > Lançamento"

### Tenant Management Domain

#### **Tenant (Inquilino)**
- **Definição**: Escritório de advocacia ou departamento jurídico que utiliza o sistema
- **Isolamento**: Dados completamente separados entre tenants
- **Identificação**: ID único, CNPJ, razão social
- **Exemplo**: "Escritório Silva & Associados Advogados"

#### **Assinatura (Subscription)**
- **Definição**: Plano contratado pelo tenant determinando funcionalidades e limites
- **Planos**: Starter, Professional, Business, Enterprise
- **Ciclo**: Mensal, Anual
- **Estado**: Ativa, Suspensa, Cancelada, Trial
- **Exemplo**: "Assinatura Professional ativa até 15/04/2024"

#### **Quota (Quota)**
- **Definição**: Limite de uso de recursos por tenant
- **Tipos**: Processos, Usuários, Consultas DataJud, Armazenamento, Notificações
- **Controle**: Soft limit (aviso), Hard limit (bloqueio)
- **Exemplo**: "Quota de 200 processos, 150 utilizados"

#### **Uso de Quota (Quota Usage)**
- **Definição**: Consumo atual de recursos pelo tenant
- **Medição**: Tempo real, diária, mensal
- **Alertas**: 80%, 90%, 100% da quota
- **Reset**: Diário (DataJud), mensal (outros)

### User Management Domain

#### **Usuário (User)**
- **Definição**: Pessoa física que acessa o sistema dentro de um tenant
- **Tipos**: Advogado, Assistente, Cliente, Administrador
- **Estados**: Ativo, Inativo, Suspenso, Pendente ativação
- **Exemplo**: "Dr. João Silva, advogado responsável pelos processos tributários"

#### **Papel (Role)**
- **Definição**: Conjunto de permissões que determina o que o usuário pode fazer
- **Hierarquia**: Admin > Advogado > Assistente > Cliente
- **Escopo**: Por tenant, não global
- **Exemplo**: "Papel 'Advogado Senior' com permissão para cadastrar processos e visualizar relatórios"

#### **Cliente (Client)**
- **Definição**: Pessoa física ou jurídica representada pelo escritório
- **Dados**: Nome, documento, contatos, processos associados
- **Relacionamento**: Um cliente pode ter múltiplos processos
- **Exemplo**: "Empresa XYZ Ltda., cliente com 5 processos em andamento"

### DataJud Integration Domain

#### **DataJud**
- **Definição**: Base de dados oficial do CNJ com informações processuais
- **API**: Interface oficial para consulta de dados processuais
- **Limite**: 10.000 consultas/dia para todo o sistema
- **Atualização**: Dados atualizados em tempo quase real pelos tribunais

#### **Consulta DataJud (DataJud Query)**
- **Definição**: Requisição à API do DataJud para obter dados de processo
- **Tipos**: Completa, Movimentações, Documentos
- **Estado**: Pendente, Executada, Falhou, Cache Hit
- **Exemplo**: "Consulta completa do processo 1234567-89.2023.8.26.0001"

#### **Cache**
- **Definição**: Armazenamento temporário de dados do DataJud para reduzir consultas
- **TTL**: Tempo de vida variável por plano (7 a 90 dias)
- **Estratégia**: LRU (Least Recently Used)
- **Invalidação**: Manual ou automática

#### **Rate Limiting**
- **Definição**: Controle de frequência de consultas para respeitar limites da API
- **Algoritmo**: Token bucket por tenant
- **Janela**: Diária (00:00 às 23:59)
- **Comportamento**: Falha gradual, fallback para cache

### Notification System Domain

#### **Notificação (Notification)**
- **Definição**: Mensagem enviada ao usuário sobre eventos relevantes
- **Tipos**: Movimentação processual, Prazo se aproximando, Quota excedida
- **Canais**: WhatsApp, Email, Telegram, SMS
- **Estados**: Agendada, Enviada, Entregue, Lida, Falhou

#### **Canal de Notificação (Notification Channel)**
- **Definição**: Meio de comunicação para entrega de notificações
- **Prioridade**: WhatsApp (principal), Email (secundário), Telegram (alternativo)
- **Configuração**: Por usuário e tipo de evento
- **Fallback**: Estratégia de canais alternativos em caso de falha

#### **Template de Mensagem (Message Template)**
- **Definição**: Modelo pré-definido para formatação de notificações
- **Variáveis**: Placeholders substituídos por dados reais
- **Personalização**: Por tenant, tipo de evento e canal
- **Exemplo**: "Olá {nome}, o processo {numero} teve nova movimentação: {descricao}"

#### **Entrega (Delivery)**
- **Definição**: Confirmação de que a notificação chegou ao destinatário
- **Status**: Enviada, Entregue, Lida, Falhou
- **Tracking**: ID externo do provedor, timestamps
- **Retry**: Política de reenvio em caso de falha

### AI & Analytics Domain

#### **Resumo de IA (AI Summary)**
- **Definição**: Síntese automática gerada por IA de movimentações ou processos
- **Tipos**: Visão geral do processo, Resumo de movimentação, Análise jurídica
- **Público**: Adaptado para advogados ou clientes leigos
- **Confiança**: Score de 0-100% da qualidade do resumo

#### **Explicação de Termo (Term Explanation)**
- **Definição**: Definição simplificada de jargão jurídico para clientes
- **Contexto**: Específica para área do direito e situação
- **Linguagem**: Adaptada ao nível de conhecimento do destinatário
- **Exemplo**: "Execução fiscal = cobrança de dívida com o governo"

#### **Jurimetria**
- **Definição**: Aplicação de métodos estatísticos ao direito
- **Análises**: Probabilidade de sucesso, tempo médio de duração, padrões de decisão
- **Base**: Dados históricos de processos similares
- **Predição**: Machine learning para prever resultados

#### **Modelo de IA (AI Model)**
- **Definição**: Algoritmo treinado para tarefas específicas do domínio jurídico
- **Tipos**: Summarização, Classificação, Predição, NER (Named Entity Recognition)
- **Treinamento**: Dados gerais + dados específicos do tenant (Enterprise)
- **Versioning**: Controle de versões e rollback

### Document Management Domain

#### **Documento (Document)**
- **Definição**: Arquivo relacionado a processo ou atividade jurídica
- **Tipos**: Petição, Contrato, Parecer, Certidão, Procuração
- **Formato**: PDF, DOCX, JPG, PNG
- **Metadata**: Título, tags, processo relacionado, data de criação

#### **Template de Documento (Document Template)**
- **Definição**: Modelo pré-formatado para geração automática de documentos
- **Variáveis**: Campos substituídos automaticamente (dados do processo, cliente, etc.)
- **Personalização**: Por tenant, área jurídica
- **Versionamento**: Controle de mudanças e aprovações

#### **Geração Automática (Auto-generation)**
- **Definição**: Criação de documentos jurídicos com base em templates e dados
- **Trigger**: Manual ou automática (eventos processuais)
- **IA**: Uso de modelos para preencher campos complexos
- **Revisão**: Processo de aprovação antes da finalização

#### **Assinatura Digital (Digital Signature)**
- **Definição**: Processo de autenticação e integridade de documentos
- **Tipos**: Simples, Avançada, Qualificada (ICP-Brasil)
- **Integração**: DocuSign, Adobe Sign, Gov.br
- **Validade**: Jurídica conforme MP 2.200-2/2001

### Analytics & Reporting Domain

#### **Dashboard**
- **Definição**: Painel visual com indicadores de performance e atividade
- **Níveis**: Tenant, Usuário, Processo, Sistema
- **KPIs**: Processos ativos, prazos pendentes, taxa de entrega de notificações
- **Atualização**: Tempo real via event streaming

#### **Relatório (Report)**
- **Definição**: Documento estruturado com análise de dados
- **Tipos**: Operacional, Gerencial, Estratégico
- **Formato**: PDF, Excel, CSV, Dashboard interativo
- **Agendamento**: Automático (diário, semanal, mensal)

#### **Métrica (Metric)**
- **Definição**: Indicador quantitativo de performance
- **Categorias**: Técnicas (latência, uptime), Negócio (conversão, uso), Produto (engajamento)
- **Agregação**: Por tenant, período, dimensão
- **Alertas**: Thresholds configuráveis

## Business Rules & Invariants

### Process Management Rules

#### **Regra de Numeração CNJ**
- Todo processo deve ter número CNJ válido conforme padrão nacional
- Dígito verificador deve ser validado matematicamente
- Tribunal e órgão julgador devem existir na tabela oficial

#### **Regra de Monitoramento**
- Processo só pode ser monitorado se estiver ativo
- Monitoramento consome quota de consultas DataJud
- Frequência máxima limitada por plano de assinatura

#### **Regra de Movimentações**
- Movimentações são ordenadas cronologicamente
- Não pode haver movimentação futura (data > hoje)
- Movimentação duplicada é ignorada (idempotência)

### Tenant Management Rules

#### **Regra de Quotas**
- Hard limit bloqueia novas operações
- Soft limit apenas alerta
- Quota reseta conforme período definido (diário/mensal)

#### **Regra de Isolamento**
- Dados de um tenant nunca são visíveis para outro
- Usuário só pode acessar recursos do próprio tenant
- APIs sempre filtram por tenant do usuário autenticado

### Notification Rules

#### **Regra de Canais**
- WhatsApp é sempre tentado primeiro (disponível em todos os planos)
- Fallback automático para email em caso de falha
- Telegram apenas para planos Professional+

#### **Regra de Retry**
- Máximo 3 tentativas por canal
- Backoff exponencial entre tentativas
- Dead letter queue após todas as falhas

### DataJud Rules

#### **Regra de Rate Limiting**
- 10.000 consultas/dia compartilhadas entre todos os tenants
- Distribuição proporcional por plano de assinatura
- Cache hit não consome quota

#### **Regra de Cache**
- TTL varia por plano (7-90 dias)
- Invalidação automática em caso de nova movimentação
- Cache por processo, não por consulta

## Anti-Patterns & What to Avoid

### **❌ "Processo Inválido"**
- **Não usar**: Processo sem número CNJ válido
- **Correto**: Sempre validar formato e dígito verificador

### **❌ "Cross-tenant Data Leak"**
- **Não fazer**: Consultas que retornem dados de múltiplos tenants
- **Correto**: Sempre filtrar por tenant_id

### **❌ "API DataJud Abuse"**
- **Não fazer**: Consultas desnecessárias ou muito frequentes
- **Correto**: Usar cache inteligente e respeitar rate limits

### **❌ "Notificação Spam"**
- **Não enviar**: Múltiplas notificações para mesmo evento
- **Correto**: Deduplicação e agrupamento de eventos similares

### **❌ "Hard-coded Business Rules"**
- **Não fazer**: Regras de negócio fixas no código
- **Correto**: Configurações por tenant e feature flags

## Context Boundaries

### **Authentication Context**
- **Responsabilidade**: "Quem é o usuário e o que pode fazer?"
- **Fronteira**: Para depois da autenticação
- **Linguagem**: User, Role, Permission, Session, Token

### **Tenant Context** 
- **Responsabilidade**: "Qual escritório e que recursos pode usar?"
- **Fronteira**: Controle de acesso e quotas
- **Linguagem**: Tenant, Subscription, Quota, Billing

### **Process Context**
- **Responsabilidade**: "Quais processos estão sendo acompanhados?"
- **Fronteira**: Ciclo de vida do processo jurídico
- **Linguagem**: Process, Movement, Deadline, Monitoring

### **DataJud Context**
- **Responsabilidade**: "Como obter dados oficiais dos tribunais?"
- **Fronteira**: Integração com APIs externas
- **Linguagem**: Query, Cache, RateLimit, Integration

### **Notification Context**
- **Responsabilidade**: "Como comunicar eventos aos usuários?"
- **Fronteira**: Entrega de mensagens
- **Linguagem**: Notification, Channel, Template, Delivery

### **AI Context**
- **Responsabilidade**: "Como tornar informação jurídica acessível?"
- **Fronteira**: Processamento de linguagem natural
- **Linguagem**: Summary, Analysis, Model, Prediction

### **Document Context**
- **Responsabilidade**: "Como gerar e gerenciar documentos jurídicos?"
- **Fronteira**: Criação e armazenamento de arquivos
- **Linguagem**: Document, Template, Generation, Signature

## Glossary Extensions

### Technical Terms

#### **Event Sourcing**
- Padrão onde mudanças de estado são armazenadas como eventos
- Usado no Process Context para auditoria completa
- Permite replay e reconstrução do estado

#### **CQRS (Command Query Responsibility Segregation)**
- Separação entre operações de escrita (commands) e leitura (queries)
- Commands modificam estado, Queries retornam dados
- Usado para otimizar performance e escalabilidade

#### **Saga Pattern**
- Coordenação de transações distribuídas entre microserviços
- Usado para fluxos complexos como cadastro de processo
- Compensa operações em caso de falha

#### **Circuit Breaker**
- Padrão para proteger contra falhas em cascata
- Usado na integração com DataJud
- Estados: Closed, Open, Half-Open

#### **Multi-tenancy**
- Arquitetura onde uma instância serve múltiplos tenants
- Isolamento completo de dados entre inquilinos
- Otimização de recursos compartilhados

### Legal Domain Terms

#### **CNJ (Conselho Nacional de Justiça)**
- Órgão que supervisiona o Poder Judiciário brasileiro
- Responsável pela padronização do número único

#### **Instância**
- Grau hierárquico do tribunal
- 1ª instância: juízes singulares
- 2ª instância: tribunais colegiados
- Instâncias superiores: STJ, STF

#### **Distribuição**
- Ato de direcionamento do processo para vara competente
- Gera o número CNJ do processo
- Marco inicial do processo

#### **Intimação**
- Comunicação oficial de ato processual às partes
- Inicia contagem de prazos
- Pode ser pessoal, por edital, eletrônica

#### **Baixa Processual**
- Encerramento definitivo do processo no tribunal
- Não confundir com arquivamento (temporário)
- Processo baixado sai do acervo ativo

This ubiquitous language serves as the foundation for all development, documentation, and communication within the Direito Lux domain, ensuring consistency and clarity across all bounded contexts.