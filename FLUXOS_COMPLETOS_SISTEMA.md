# 🎯 FLUXOS COMPLETOS DO SISTEMA DIREITO LUX

## 📋 **JORNADA COMPLETA DO USUÁRIO - INÍCIO AO FIM**

---

## 🚀 **FASE 1: DESCOBERTA E AQUISIÇÃO**

### **1.1 Primeiro Contato (Landing Page)**

```yaml
CENÁRIO: Advogado pesquisa "monitoramento processos jurídicos"
├── Acessa: https://direitolux.com.br
├── Vê: Landing page com benefícios
├── Clica: "Teste Grátis 15 Dias"
└── Inicia: Processo de registro
```

### **1.2 Processo de Registro**

```yaml
PASSO 1: Dados Básicos
├── Nome: "Dr. João Silva"
├── Email: "joao@silva-advogados.com.br"
├── Telefone: "+55 11 99999-9999"
├── CNPJ: "12.345.678/0001-90"
└── Senha: [validação força]

PASSO 2: Dados do Escritório
├── Nome Escritório: "Silva & Associados"
├── Endereço: "Rua da Consolação, 1000, São Paulo"
├── Área Jurídica: "Cível, Trabalhista, Empresarial"
└── Tamanho: "5-10 advogados"

PASSO 3: Escolha do Plano
├── Starter: R$ 99/mês (20 clientes, 50 processos)
├── Professional: R$ 299/mês (100 clientes, 200 processos)
├── Business: R$ 699/mês (500 clientes, 500 processos)
└── Escolha: "Professional" (15 dias grátis)
```

### **1.3 Código de Verificação**

```go
// internal/application/onboarding_service.go
func (s *OnboardingService) ProcessRegistration(ctx context.Context, req *RegistrationRequest) (*RegistrationResponse, error) {
    // 1. Validar dados
    if err := s.validateRegistrationData(req); err != nil {
        return nil, err
    }
    
    // 2. Verificar se email já existe
    if exists, err := s.userRepository.EmailExists(ctx, req.Email); err != nil {
        return nil, err
    } else if exists {
        return nil, fmt.Errorf("email já cadastrado")
    }
    
    // 3. Criar tenant (escritório)
    tenant := &Tenant{
        ID:           uuid.New().String(),
        Name:         req.CompanyName,
        CNPJ:         req.CNPJ,
        Address:      req.Address,
        LegalAreas:   req.LegalAreas,
        Plan:         req.Plan,
        TrialEndsAt:  time.Now().Add(15 * 24 * time.Hour), // 15 dias
        CreatedAt:    time.Now(),
        Status:       "trial",
    }
    
    if err := s.tenantRepository.CreateTenant(ctx, tenant); err != nil {
        return nil, err
    }
    
    // 4. Criar usuário administrador
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    
    user := &User{
        ID:           uuid.New().String(),
        TenantID:     tenant.ID,
        Email:        req.Email,
        Name:         req.Name,
        Phone:        req.Phone,
        PasswordHash: string(hashedPassword),
        Role:         "ADVOGADO",
        IsAdmin:      true,
        CreatedAt:    time.Now(),
        PhoneVerified: false,
    }
    
    if err := s.userRepository.CreateUser(ctx, user); err != nil {
        return nil, err
    }
    
    // 5. Enviar código de verificação
    code := s.generateVerificationCode()
    if err := s.smsService.SendVerificationCode(ctx, req.Phone, code); err != nil {
        return nil, err
    }
    
    // 6. Criar trial subscription
    subscription := &Subscription{
        ID:        uuid.New().String(),
        TenantID:  tenant.ID,
        Plan:      req.Plan,
        Status:    "trial",
        StartsAt:  time.Now(),
        EndsAt:    time.Now().Add(15 * 24 * time.Hour),
        CreatedAt: time.Now(),
    }
    
    if err := s.subscriptionRepository.CreateSubscription(ctx, subscription); err != nil {
        return nil, err
    }
    
    // 7. Configurar quotas iniciais
    s.quotaService.InitializeTenantQuotas(ctx, tenant.ID, req.Plan)
    
    return &RegistrationResponse{
        TenantID:        tenant.ID,
        UserID:          user.ID,
        VerificationSent: true,
        TrialEndsAt:     tenant.TrialEndsAt,
    }, nil
}
```

### **1.4 Primeira Experiência (Onboarding)**

```yaml
VERIFICAÇÃO TELEFONE:
├── SMS: "Código: 123456"
├── Input: Usuário digita código
├── Validação: Sistema confirma
└── Sucesso: Conta ativada

TOUR GUIADO:
├── Tela 1: "Bem-vindo ao Direito Lux!"
├── Tela 2: "Vamos adicionar seu primeiro cliente"
├── Tela 3: "Agora vamos cadastrar um processo"
├── Tela 4: "Configure as notificações"
└── Tela 5: "Conheça a Luxia, sua assistente IA"
```

---

## 📋 **FASE 2: CONFIGURAÇÃO INICIAL**

### **2.1 Primeiro Cliente**

```yaml
FLUXO ONBOARDING:
├── Modal: "Vamos cadastrar seu primeiro cliente"
├── Dados: Nome, Email, Telefone, CPF/CNPJ
├── Validação: Sistema verifica CPF/CNPJ
├── Salvamento: Cliente criado
└── Sucesso: "Cliente João Silva adicionado!"

NOTIFICAÇÃO QUOTA:
├── Status: "Clientes: 1/100 (Professional)"
├── Aviso: "Você tem 99 clientes disponíveis"
└── Tip: "Importe clientes em lote (CSV)"
```

### **2.2 Primeiro Processo**

```yaml
CADASTRO PROCESSO:
├── Número: "1001234-56.2024.8.26.0100"
├── Validação: Sistema verifica formato CNJ
├── Consulta: DataJud busca dados iniciais
├── Dados: Tribunal, Assunto, Partes
├── Cliente: Vincula ao cliente cadastrado
└── Monitoramento: Inicia automaticamente

CONSULTA DATAJUD:
├── API: https://api-publica.datajud.cnj.jus.br
├── CNPJ: 12.345.678/0001-90 (escritório)
├── Query: Elasticsearch busca processo
├── Response: Dados completos do processo
└── Cache: Salva na base local
```

### **2.3 Configuração de Notificações**

```yaml
CANAIS DISPONÍVEIS:
├── WhatsApp: +55 11 99999-9999 (principal)
├── Email: joao@silva-advogados.com.br
├── Telegram: @drjoaosilva
└── SMS: +55 11 99999-9999 (backup)

TIPOS DE NOTIFICAÇÃO:
├── Movimento processual: ✅ Ativado
├── Prazos vencendo: ✅ Ativado
├── Audiências: ✅ Ativado
├── Sentenças: ✅ Ativado
├── Recursos: ✅ Ativado
└── Arquivamento: ✅ Ativado

FREQUÊNCIA:
├── Tempo real: ✅ Professional
├── Diário: ✅ Todas as 8h
├── Semanal: ✅ Segundas 9h
└── Mensal: ✅ Dia 1 às 9h
```

### **2.4 Configuração Luxia (Bot)**

```yaml
ATIVAÇÃO WHATSAPP:
├── QR Code: Sistema gera QR
├── Scan: Usuário escaneia no WhatsApp
├── Webhook: Meta configura webhook
├── Teste: "Oi Luxia" → "Olá! Sou sua assistente"
└── Sucesso: Bot ativo no WhatsApp

CONFIGURAÇÃO INICIAL:
├── Horário: 8h às 18h (horário comercial)
├── Idioma: Português brasileiro
├── Personalidade: Profissional e amigável
├── Contexto: Direito brasileiro
└── Limites: Apenas tópicos jurídicos
```

---

## 📋 **FASE 3: OPERAÇÃO DIÁRIA**

### **3.1 Monitoramento Automático**

```go
// internal/application/monitoring_service.go
func (s *MonitoringService) StartDailyMonitoring(ctx context.Context) {
    // Polling a cada 30 minutos (Professional)
    ticker := time.NewTicker(30 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if s.isWorkingHours() {
                s.checkAllProcesses(ctx)
            }
        case <-ctx.Done():
            return
        }
    }
}

func (s *MonitoringService) checkAllProcesses(ctx context.Context) {
    // Buscar todos os processos monitorados
    processes, err := s.processRepository.GetMonitoredProcesses(ctx)
    if err != nil {
        return
    }
    
    // Processar em lotes de 50
    for i := 0; i < len(processes); i += 50 {
        end := i + 50
        if end > len(processes) {
            end = len(processes)
        }
        
        batch := processes[i:end]
        s.processBatch(ctx, batch)
        
        // Delay entre lotes (rate limiting)
        time.Sleep(2 * time.Second)
    }
}

func (s *MonitoringService) processBatch(ctx context.Context, processes []Process) {
    for _, process := range processes {
        // Consultar DataJud
        updates, err := s.dataJudService.GetProcessUpdates(ctx, process.Number)
        if err != nil {
            continue
        }
        
        // Verificar mudanças
        if s.detectChanges(process, updates) {
            // Atualizar base local
            s.updateProcess(ctx, process, updates)
            
            // Notificar usuários
            s.notifyProcessUpdate(ctx, process, updates)
        }
    }
}
```

### **3.2 Detecção de Mudanças**

```yaml
VERIFICAÇÃO AUTOMÁTICA:
├── Timestamp: Última atualização no DataJud
├── Movimentos: Novos andamentos
├── Status: Mudança de fase
├── Partes: Alterações nas partes
└── Decisões: Sentenças, despachos

TIPOS DE MUDANÇA:
├── Juntada: "Juntada de petição"
├── Despacho: "Despacho do juiz"
├── Sentença: "Sentença proferida"
├── Recurso: "Recurso interposto"
├── Audiência: "Audiência designada"
└── Arquivo: "Processo arquivado"
```

### **3.3 Notificação Inteligente**

```go
// internal/application/notification_service.go
func (s *NotificationService) NotifyProcessUpdate(ctx context.Context, process Process, update ProcessUpdate) {
    // 1. Gerar resumo com IA
    summary, err := s.aiService.GenerateUpdateSummary(ctx, process, update)
    if err != nil {
        summary = s.generateBasicSummary(process, update)
    }
    
    // 2. Determinar urgência
    urgency := s.determineUrgency(update)
    
    // 3. Personalizar mensagem por destinatário
    users := s.getUsersForNotification(ctx, process.TenantID)
    
    for _, user := range users {
        message := s.personalizeMessage(user, process, update, summary, urgency)
        
        // 4. Enviar por canais configurados
        if user.WhatsAppEnabled {
            s.sendWhatsAppMessage(ctx, user.Phone, message)
        }
        
        if user.EmailEnabled {
            s.sendEmailNotification(ctx, user.Email, process, update, summary)
        }
        
        if urgency == "HIGH" && user.SMSEnabled {
            s.sendSMSAlert(ctx, user.Phone, fmt.Sprintf("URGENTE: %s", summary))
        }
    }
}

func (s *NotificationService) personalizeMessage(user User, process Process, update ProcessUpdate, summary string, urgency string) string {
    var urgencyIcon string
    switch urgency {
    case "HIGH":
        urgencyIcon = "🚨"
    case "MEDIUM":
        urgencyIcon = "⚠️"
    default:
        urgencyIcon = "📋"
    }
    
    return fmt.Sprintf(`
%s *Atualização Processual*

📋 *Processo:* %s
👤 *Cliente:* %s
🏛️ *Tribunal:* %s

📝 *Movimento:*
%s

🤖 *Resumo da Luxia:*
%s

📅 *Data:* %s

Deseja mais detalhes ou tem alguma dúvida?`,
        urgencyIcon,
        process.Number,
        process.ClientName,
        process.Court,
        update.MovementDescription,
        summary,
        update.Date.Format("02/01/2006 15:04"))
}
```

### **3.4 Exemplo de Notificação Real**

```
🚨 *Atualização Processual*

📋 *Processo:* 1001234-56.2024.8.26.0100
👤 *Cliente:* João Silva
🏛️ *Tribunal:* TJSP

📝 *Movimento:*
Sentença proferida - Julgamento Procedente

🤖 *Resumo da Luxia:*
Ótima notícia! O juiz julgou o pedido procedente. 
Isso significa que o cliente João Silva venceu a 
ação. O réu deverá pagar R$ 15.000 de indenização 
por danos morais. Prazo para recurso: 15 dias.

📅 *Data:* 15/01/2025 14:30

Deseja mais detalhes ou tem alguma dúvida?
```

---

## 📋 **FASE 4: ATENDIMENTO AO CLIENTE**

### **4.1 Cliente Acessa Luxia**

```yaml
PRIMEIRO ACESSO:
├── WhatsApp: Cliente manda "Oi" para número
├── Autenticação: Sistema pede telefone
├── Validação: Código SMS enviado
├── Confirmação: "123456"
└── Boas-vindas: "Olá João! Sou a Luxia"

AUTENTICAÇÃO TRANSPARENTE:
├── Telefone: +55 11 88888-8888
├── Busca: Sistema encontra cliente
├── Vínculo: Cliente do escritório Silva & Associados
├── Permissões: Apenas processos próprios
└── Sessão: 24h de duração
```

### **4.2 Interação Cliente-Luxia**

```yaml
EXEMPLO 1: Status do Processo
├── Cliente: "Como está meu processo?"
├── Luxia: [Busca processos do cliente]
├── Resposta: "Você tem 2 processos ativos:"
├── Lista: Processo 1, Processo 2
└── Ação: "Qual deseja ver detalhes?"

EXEMPLO 2: Explicação Jurídica
├── Cliente: "O que significa 'julgamento procedente'?"
├── Luxia: [IA processa pergunta]
├── Resposta: "Significa que você ganhou!"
├── Explicação: Detalhes em linguagem simples
└── Sugestão: "Tem mais alguma dúvida?"

EXEMPLO 3: Agendamento
├── Cliente: "Quero agendar reunião"
├── Luxia: [Consulta agenda do advogado]
├── Opções: "Horários disponíveis:"
├── Escolha: Cliente seleciona
└── Confirmação: "Reunião agendada!"
```

### **4.3 Controle de Acesso do Cliente**

```go
// internal/application/client_access_service.go
func (s *ClientAccessService) ValidateClientAccess(ctx context.Context, clientPhone, processNumber string) error {
    // 1. Buscar cliente por telefone
    client, err := s.clientRepository.GetClientByPhone(ctx, clientPhone)
    if err != nil {
        return fmt.Errorf("cliente não encontrado")
    }
    
    // 2. Buscar processo
    process, err := s.processRepository.GetProcessByNumber(ctx, processNumber)
    if err != nil {
        return fmt.Errorf("processo não encontrado")
    }
    
    // 3. Verificar se cliente é parte do processo
    if !s.isClientInProcess(client.ID, process.ID) {
        return fmt.Errorf("cliente não tem acesso a este processo")
    }
    
    // 4. Verificar se tenant é o mesmo
    if client.TenantID != process.TenantID {
        return fmt.Errorf("acesso negado")
    }
    
    return nil
}

func (s *ClientAccessService) GetClientProcesses(ctx context.Context, clientID string) ([]Process, error) {
    // Buscar apenas processos onde cliente é parte
    processes, err := s.processRepository.GetProcessesByClient(ctx, clientID)
    if err != nil {
        return nil, err
    }
    
    // Filtrar informações sensíveis
    filteredProcesses := make([]Process, len(processes))
    for i, process := range processes {
        filteredProcesses[i] = Process{
            Number:      process.Number,
            Court:       process.Court,
            Subject:     process.Subject,
            Status:      process.Status,
            LastUpdate:  process.LastUpdate,
            // Remover informações internas
            // CreatedBy, InternalNotes, etc.
        }
    }
    
    return filteredProcesses, nil
}
```

### **4.4 Diferentes Tipos de Usuário**

```yaml
ADVOGADO (Dr. João Silva):
├── Acesso: Todos os processos do escritório
├── Luxia: "Você tem 85 processos ativos"
├── Comandos: "Relatório mensal", "Buscar jurisprudência"
├── Permissões: Criar, editar, deletar
└── Notificações: Todas as atualizações

CLIENTE (João Silva):
├── Acesso: Apenas processos próprios
├── Luxia: "Você tem 2 processos ativos"
├── Comandos: "Meus processos", "Agendar reunião"
├── Permissões: Apenas visualizar
└── Notificações: Apenas processos próprios

FUNCIONÁRIO (Maria Santos):
├── Acesso: Processos atribuídos
├── Luxia: "Você tem 12 processos atribuídos"
├── Comandos: "Meus processos", "Prazos vencendo"
├── Permissões: Visualizar, editar atribuídos
└── Notificações: Apenas processos atribuídos

ESTAGIÁRIO (Pedro Costa):
├── Acesso: Supervisão obrigatória
├── Luxia: "Você tem 5 processos para estudar"
├── Comandos: "Buscar jurisprudência", "Estudar caso"
├── Permissões: Apenas visualizar
└── Notificações: Educativas e de estudo
```

---

## 📋 **FASE 5: CICLO DE VIDA DO PROCESSO**

### **5.1 Fases do Processo**

```yaml
CADASTRO INICIAL:
├── Status: "Distribuído"
├── Dados: Básicos do DataJud
├── Monitoramento: Inicia automaticamente
├── Notificações: Configuradas
└── Cliente: Vinculado e notificado

FASE INSTRUTÓRIA:
├── Movimentos: Citação, contestação, tréplica
├── Audiências: Conciliação, instrução
├── Provas: Documentos, testemunhas
├── Prazos: Manifestações, recursos
└── Notificações: Cada movimento

FASE DECISÓRIA:
├── Sentença: Procedente/improcedente
├── Análise: IA resume decisão
├── Prazo: 15 dias para recurso
├── Recomendação: Recurso ou não
└── Cliente: Notificado com explicação

FASE RECURSAL:
├── Recurso: Apelação, embargos
├── Contraminuta: Resposta da outra parte
├── Tribunal: Remessa ao 2º grau
├── Julgamento: Acórdão
└── Resultado: Confirmado ou reformado

EXECUÇÃO:
├── Título: Executivo constituído
├── Citação: Devedor citado
├── Pagamento: Voluntário ou forçado
├── Penhora: Bens penhorados
└── Satisfação: Crédito quitado

ARQUIVAMENTO:
├── Trânsito: Julgado transitado
├── Cumprimento: Obrigação cumprida
├── Arquivo: Processo arquivado
├── Baixa: Monitoramento suspenso
└── Histórico: Mantido na base
```

### **5.2 Gestão de Prazos**

```go
// internal/application/deadline_service.go
func (s *DeadlineService) MonitorDeadlines(ctx context.Context) {
    // Buscar prazos vencendo
    deadlines, err := s.deadlineRepository.GetUpcomingDeadlines(ctx, 5) // 5 dias
    if err != nil {
        return
    }
    
    for _, deadline := range deadlines {
        daysLeft := int(deadline.DueDate.Sub(time.Now()).Hours() / 24)
        
        var urgency string
        var message string
        
        switch daysLeft {
        case 0:
            urgency = "CRITICAL"
            message = "🚨 PRAZO VENCE HOJE!"
        case 1:
            urgency = "HIGH"
            message = "⚠️ PRAZO VENCE AMANHÃ!"
        case 2:
            urgency = "MEDIUM"
            message = "📅 PRAZO VENCE EM 2 DIAS"
        default:
            urgency = "LOW"
            message = fmt.Sprintf("📅 PRAZO VENCE EM %d DIAS", daysLeft)
        }
        
        notification := &DeadlineNotification{
            ProcessNumber: deadline.ProcessNumber,
            ClientName:    deadline.ClientName,
            Description:   deadline.Description,
            DueDate:       deadline.DueDate,
            Urgency:       urgency,
            Message:       message,
        }
        
        s.sendDeadlineNotification(ctx, notification)
    }
}

func (s *DeadlineService) sendDeadlineNotification(ctx context.Context, notification *DeadlineNotification) {
    message := fmt.Sprintf(`
%s

📋 *Processo:* %s
👤 *Cliente:* %s
⏰ *Prazo:* %s
📝 *Descrição:* %s

🎯 *Ação necessária:*
- Preparar petição
- Protocolar no sistema
- Confirmar entrega

Precisa de ajuda?`,
        notification.Message,
        notification.ProcessNumber,
        notification.ClientName,
        notification.DueDate.Format("02/01/2006 15:04"),
        notification.Description)
    
    // Enviar para advogado responsável
    s.notificationService.SendUrgentNotification(ctx, notification.ProcessNumber, message)
}
```

### **5.3 Análise de Jurisprudência**

```yaml
BUSCA AUTOMÁTICA:
├── Trigger: Movimento importante (sentença, acórdão)
├── Análise: IA identifica temas jurídicos
├── Busca: Precedentes similares
├── Resultado: Top 5 casos relevantes
└── Notificação: Enviada ao advogado

EXEMPLO PRÁTICO:
├── Processo: Danos morais por negativação
├── Sentença: Improcedente
├── Análise: IA encontra precedentes favoráveis
├── Sugestão: "Encontrei 3 acórdãos do STJ favoráveis"
└── Ação: "Deseja que eu prepare minuta de recurso?"
```

### **5.4 Relatórios Automáticos**

```go
// internal/application/report_service.go
func (s *ReportService) GenerateMonthlyReport(ctx context.Context, tenantID string) (*MonthlyReport, error) {
    // Dados do mês
    startDate := time.Now().AddDate(0, -1, 0)
    endDate := time.Now()
    
    // Estatísticas
    stats := s.gatherMonthlyStats(ctx, tenantID, startDate, endDate)
    
    // Processos novos
    newProcesses, _ := s.processRepository.GetNewProcesses(ctx, tenantID, startDate, endDate)
    
    // Processos concluídos
    completedProcesses, _ := s.processRepository.GetCompletedProcesses(ctx, tenantID, startDate, endDate)
    
    // Movimentos importantes
    importantMovements, _ := s.processRepository.GetImportantMovements(ctx, tenantID, startDate, endDate)
    
    // Prazos perdidos
    missedDeadlines, _ := s.deadlineRepository.GetMissedDeadlines(ctx, tenantID, startDate, endDate)
    
    report := &MonthlyReport{
        TenantID:           tenantID,
        Period:             fmt.Sprintf("%s a %s", startDate.Format("02/01/2006"), endDate.Format("02/01/2006")),
        Stats:              stats,
        NewProcesses:       newProcesses,
        CompletedProcesses: completedProcesses,
        ImportantMovements: importantMovements,
        MissedDeadlines:    missedDeadlines,
        GeneratedAt:        time.Now(),
    }
    
    return report, nil
}

func (s *ReportService) gatherMonthlyStats(ctx context.Context, tenantID string, start, end time.Time) *MonthlyStats {
    return &MonthlyStats{
        TotalProcesses:     s.processRepository.CountProcesses(ctx, tenantID),
        NewProcesses:       s.processRepository.CountNewProcesses(ctx, tenantID, start, end),
        CompletedProcesses: s.processRepository.CountCompletedProcesses(ctx, tenantID, start, end),
        ActiveProcesses:    s.processRepository.CountActiveProcesses(ctx, tenantID),
        WonCases:          s.processRepository.CountWonCases(ctx, tenantID, start, end),
        LostCases:         s.processRepository.CountLostCases(ctx, tenantID, start, end),
        PendingCases:      s.processRepository.CountPendingCases(ctx, tenantID),
        SuccessRate:       s.calculateSuccessRate(ctx, tenantID, start, end),
        AverageTime:       s.calculateAverageTime(ctx, tenantID, start, end),
        TotalRevenue:      s.calculateTotalRevenue(ctx, tenantID, start, end),
    }
}
```

---

## 📋 **FASE 6: ACESSOS POR PLANO**

### **6.1 Funcionalidades por Plano**

```yaml
STARTER (R$ 99/mês):
├── Processos: 50 máximo
├── Clientes: 20 máximo
├── Usuários: 3 máximo
├── Monitoramento: 1 hora
├── Notificações: WhatsApp + Email
├── Relatórios: Básicos mensais
├── Luxia: Comandos básicos
├── Jurisprudência: ❌ Não disponível
├── API: ❌ Não disponível
└── Suporte: Email

PROFESSIONAL (R$ 299/mês):
├── Processos: 200 máximo
├── Clientes: 100 máximo
├── Usuários: 10 máximo
├── Monitoramento: 30 minutos
├── Notificações: WhatsApp + Email + Telegram
├── Relatórios: Avançados semanais
├── Luxia: Comandos avançados + IA
├── Jurisprudência: ✅ 100 buscas/dia
├── API: ❌ Não disponível
└── Suporte: Email + Chat

BUSINESS (R$ 699/mês):
├── Processos: 500 máximo
├── Clientes: 500 máximo
├── Usuários: 50 máximo
├── Monitoramento: 15 minutos
├── Notificações: Todos os canais
├── Relatórios: Personalizados diários
├── Luxia: Comandos completos + IA avançada
├── Jurisprudência: ✅ 500 buscas/dia
├── API: ✅ 1000 calls/dia
└── Suporte: Telefone + Chat prioritário

ENTERPRISE (R$ 1.999/mês):
├── Processos: Ilimitado
├── Clientes: Ilimitado
├── Usuários: Ilimitado
├── Monitoramento: 10 minutos
├── Notificações: Todos + SMS
├── Relatórios: Personalizados tempo real
├── Luxia: Comandos customizados + IA
├── Jurisprudência: ✅ Ilimitado
├── API: ✅ Ilimitado
└── Suporte: Dedicado + WhatsApp
```

### **6.2 Controle de Acesso Dinâmico**

```go
// internal/application/access_control_service.go
func (s *AccessControlService) ValidateFeatureAccess(ctx context.Context, userID, feature string) error {
    // 1. Obter dados do usuário
    user, err := s.userRepository.GetUser(ctx, userID)
    if err != nil {
        return err
    }
    
    // 2. Obter dados do tenant
    tenant, err := s.tenantRepository.GetTenant(ctx, user.TenantID)
    if err != nil {
        return err
    }
    
    // 3. Verificar status da assinatura
    if tenant.Status != "active" && tenant.Status != "trial" {
        return fmt.Errorf("assinatura inativa")
    }
    
    // 4. Verificar funcionalidades do plano
    planFeatures := s.getPlanFeatures(tenant.Plan)
    
    switch feature {
    case "jurisprudence_search":
        if !planFeatures.JurisprudenceSearch {
            return &FeatureNotAvailableError{
                Feature: feature,
                Plan:    tenant.Plan,
                RequiredPlans: []string{"professional", "business", "enterprise"},
            }
        }
        
        // Verificar quota diária
        dailyUsage, _ := s.usageRepository.GetDailyJurisprudenceSearches(ctx, user.TenantID)
        if dailyUsage >= planFeatures.MaxJurisprudenceSearches {
            return &QuotaExceededError{
                Current:  dailyUsage,
                Limit:    planFeatures.MaxJurisprudenceSearches,
                Resource: "jurisprudence_searches",
            }
        }
        
    case "api_access":
        if !planFeatures.APIAccess {
            return &FeatureNotAvailableError{
                Feature: feature,
                Plan:    tenant.Plan,
                RequiredPlans: []string{"business", "enterprise"},
            }
        }
        
        // Verificar quota de API
        dailyAPICalls, _ := s.usageRepository.GetDailyAPICalls(ctx, user.TenantID)
        if dailyAPICalls >= planFeatures.MaxAPICalls {
            return &QuotaExceededError{
                Current:  dailyAPICalls,
                Limit:    planFeatures.MaxAPICalls,
                Resource: "api_calls",
            }
        }
        
    case "real_time_monitoring":
        if !planFeatures.RealTimeMonitoring {
            return &FeatureNotAvailableError{
                Feature: feature,
                Plan:    tenant.Plan,
                RequiredPlans: []string{"professional", "business", "enterprise"},
            }
        }
        
    case "custom_reports":
        if !planFeatures.CustomReports {
            return &FeatureNotAvailableError{
                Feature: feature,
                Plan:    tenant.Plan,
                RequiredPlans: []string{"business", "enterprise"},
            }
        }
    }
    
    return nil
}
```

---

## 📋 **FASE 7: FINALIZAÇÃO DO PROCESSO**

### **7.1 Processo Finalizado**

```yaml
DETECÇÃO AUTOMÁTICA:
├── Monitoramento: Sistema detecta status final
├── Situações: "Arquivado", "Baixado", "Cumprido"
├── Análise: IA verifica se é finalização definitiva
├── Confirmação: Advogado confirma finalização
└── Ações: Parar monitoramento, manter histórico

NOTIFICAÇÃO FINAL:
├── Advogado: "Processo finalizado com sucesso"
├── Cliente: "Seu processo foi concluído"
├── Resultado: "Procedente" ou "Improcedente"
├── Resumo: IA explica resultado final
└── Próximos: Orientações sobre recursos
```

### **7.2 Fluxo de Arquivamento**

```go
// internal/application/archival_service.go
func (s *ArchivalService) ProcessArchival(ctx context.Context, processID string, reason string) error {
    // 1. Buscar processo
    process, err := s.processRepository.GetProcess(ctx, processID)
    if err != nil {
        return err
    }
    
    // 2. Verificar se pode arquivar
    if !s.canArchive(process) {
        return fmt.Errorf("processo não pode ser arquivado")
    }
    
    // 3. Gerar relatório final
    finalReport, err := s.generateFinalReport(ctx, process)
    if err != nil {
        return err
    }
    
    // 4. Parar monitoramento
    s.monitoringService.StopMonitoring(ctx, processID)
    
    // 5. Atualizar status
    process.Status = "archived"
    process.ArchivedAt = time.Now()
    process.ArchivalReason = reason
    process.FinalReport = finalReport
    
    if err := s.processRepository.UpdateProcess(ctx, process); err != nil {
        return err
    }
    
    // 6. Notificar stakeholders
    s.notifyArchival(ctx, process)
    
    // 7. Backup de segurança
    s.backupService.BackupProcess(ctx, process)
    
    return nil
}

func (s *ArchivalService) generateFinalReport(ctx context.Context, process Process) (*FinalReport, error) {
    // Histórico completo
    movements, _ := s.processRepository.GetAllMovements(ctx, process.ID)
    
    // Documentos
    documents, _ := s.documentRepository.GetProcessDocuments(ctx, process.ID)
    
    // Timeline
    timeline := s.generateTimeline(movements)
    
    // Resultado
    result := s.analyzeResult(process, movements)
    
    // Custos
    costs := s.calculateCosts(ctx, process.ID)
    
    return &FinalReport{
        ProcessID:    process.ID,
        ProcessNumber: process.Number,
        Client:       process.ClientName,
        StartDate:    process.CreatedAt,
        EndDate:      time.Now(),
        Duration:     time.Since(process.CreatedAt),
        Result:       result,
        Timeline:     timeline,
        Documents:    documents,
        Costs:        costs,
        GeneratedAt:  time.Now(),
    }, nil
}
```

### **7.3 Notificação Final**

```
✅ *Processo Finalizado*

📋 *Processo:* 1001234-56.2024.8.26.0100
👤 *Cliente:* João Silva
🏛️ *Tribunal:* TJSP

🎯 *Resultado:* PROCEDENTE ✅

🤖 *Resumo da Luxia:*
Parabéns! Seu processo foi julgado procedente.
O cliente João Silva venceu a ação e receberá
R$ 15.000 de indenização por danos morais.

📊 *Estatísticas:*
• Duração: 8 meses e 15 dias
• Audiências: 2 realizadas
• Recursos: 1 interposto
• Resultado: Vitória integral

💰 *Valores:*
• Indenização: R$ 15.000,00
• Honorários: R$ 3.000,00
• Custas: R$ 500,00

📋 *Próximos passos:*
• Acompanhar cumprimento da sentença
• Iniciar execução se necessário
• Processo arquivado automaticamente

🏆 *Sucesso!* Mais uma vitória para o escritório!
```

### **7.4 Dados Históricos**

```yaml
MANUTENÇÃO DE DADOS:
├── Processo: Mantido na base (histórico)
├── Monitoramento: Suspenso
├── Notificações: Paradas
├── Relatórios: Disponíveis
├── Backup: Realizado
└── Acesso: Consulta apenas

REATIVAÇÃO:
├── Situação: Processo reaberto
├── Detecção: Movimento após arquivo
├── Notificação: "Processo reativado"
├── Monitoramento: Reiniciado
└── Status: Volta ao ativo
```

---

## 📋 **FASE 8: MÉTRICAS E ANALYTICS**

### **8.1 Dashboard em Tempo Real**

```yaml
VISÃO GERAL:
├── Processos ativos: 47
├── Novos este mês: 8
├── Finalizados: 12
├── Taxa de sucesso: 78%
├── Prazo médio: 6.5 meses
└── Receita gerada: R$ 125.000

KPIs PRINCIPAIS:
├── Produtividade: 15 processos/advogado
├── Eficiência: 92% prazos cumpridos
├── Satisfação: 4.8/5 (clientes)
├── Crescimento: +15% vs mês anterior
└── ROI: 340% retorno investimento
```

### **8.2 Relatórios Automatizados**

```go
// internal/application/analytics_service.go
func (s *AnalyticsService) GenerateExecutiveDashboard(ctx context.Context, tenantID string) (*ExecutiveDashboard, error) {
    // Período atual
    now := time.Now()
    currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
    lastMonth := currentMonth.AddDate(0, -1, 0)
    
    // Métricas do mês atual
    currentStats := s.getMonthStats(ctx, tenantID, currentMonth, now)
    
    // Métricas do mês anterior
    lastStats := s.getMonthStats(ctx, tenantID, lastMonth, currentMonth)
    
    // Calcular variações
    variations := s.calculateVariations(currentStats, lastStats)
    
    // Tendências
    trends := s.calculateTrends(ctx, tenantID, 6) // 6 meses
    
    // Previsões
    forecast := s.generateForecast(ctx, tenantID, trends)
    
    return &ExecutiveDashboard{
        TenantID:     tenantID,
        Period:       fmt.Sprintf("%s", currentMonth.Format("January 2006")),
        CurrentStats: currentStats,
        LastStats:    lastStats,
        Variations:   variations,
        Trends:       trends,
        Forecast:     forecast,
        GeneratedAt:  time.Now(),
    }, nil
}
```

---

## ✅ **RESUMO DOS FLUXOS COMPLETOS**

### **🎯 Jornada Mapeada:**

1. **✅ DESCOBERTA** - Landing → Registro → Verificação → Onboarding
2. **✅ CONFIGURAÇÃO** - Primeiro cliente → Primeiro processo → Notificações → Luxia
3. **✅ OPERAÇÃO** - Monitoramento → Detecção → Notificação → Interação
4. **✅ ATENDIMENTO** - Cliente acessa → Autenticação → Consultas → Suporte
5. **✅ CICLO DE VIDA** - Cadastro → Instrução → Decisão → Recurso → Execução → Arquivo
6. **✅ CONTROLE** - Quotas → Planos → Funcionalidades → Upgrades
7. **✅ FINALIZAÇÃO** - Detecção → Arquivo → Relatório → Backup → Histórico
8. **✅ ANALYTICS** - Métricas → Relatórios → Trends → Forecast

### **📊 Controles Implementados:**
- **Quotas por plano** - Recursos limitados por assinatura
- **Acessos por role** - Advogado, cliente, funcionário, estagiário
- **Funcionalidades tier** - Starter, Professional, Business, Enterprise
- **Monitoramento inteligente** - Frequência por plano
- **Notificações personalizadas** - Por usuário e urgência
- **Auditoria completa** - Todas as ações rastreadas

---

## 📚 **ARQUIVO CRIADO:**
`FLUXOS_COMPLETOS_SISTEMA.md` - **Documentação completa** da jornada do usuário

**🎯 Agora você tem TODOS os fluxos do sistema documentados do início ao fim!**

**Ficou claro ou algum fluxo específico precisa de mais detalhamento?**