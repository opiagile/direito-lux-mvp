# 🎯 FLUXOS DE CONTROLE DE QUOTAS POR PLANO

## 📋 **FLUXO DETALHADO: ADICIONAR CLIENTE (PLANO STARTER)**

### **🔄 Fluxo Completo Documentado**

```yaml
CONTEXTO:
- Usuário: Dr. Silva (Advogado)
- Plano: Starter
- Limites: 50 processos, 20 clientes
- Status atual: 18 clientes cadastrados
```

### **1. 👤 SOLICITAÇÃO VIA LUXIA**

```
👤 Dr. Silva: "Luxia, preciso adicionar novo cliente"
🤖 Luxia: "Claro! Vou ajudar você a cadastrar um novo cliente.
          
          📊 Status atual:
          • Clientes: 18/20 (Starter)
          • Processos: 32/50 (Starter)
          
          Por favor, me informe os dados do cliente:"
```

### **2. 🔍 VALIDAÇÃO DE QUOTA (BACKEND)**

```go
// internal/application/client_service.go
func (s *ClientService) ValidateClientCreation(ctx context.Context, userID string) error {
    // 1. Obter dados do usuário e plano
    user, err := s.userRepository.GetUser(ctx, userID)
    if err != nil {
        return err
    }
    
    // 2. Obter limites do plano
    planLimits := s.getPlanLimits(user.Plan)
    
    // 3. Contar clientes atuais
    currentClients, err := s.clientRepository.CountClientsByTenant(ctx, user.TenantID)
    if err != nil {
        return err
    }
    
    // 4. Verificar se pode adicionar mais
    if currentClients >= planLimits.MaxClients {
        return &QuotaExceededError{
            Current: currentClients,
            Limit:   planLimits.MaxClients,
            Plan:    user.Plan,
            Resource: "clients",
        }
    }
    
    return nil
}

// getPlanLimits retorna limites por plano
func (s *ClientService) getPlanLimits(plan string) PlanLimits {
    limits := map[string]PlanLimits{
        "starter": {
            MaxClients:   20,
            MaxProcesses: 50,
            MaxUsers:     3,
        },
        "professional": {
            MaxClients:   100,
            MaxProcesses: 200,
            MaxUsers:     10,
        },
        "business": {
            MaxClients:   500,
            MaxProcesses: 500,
            MaxUsers:     50,
        },
        "enterprise": {
            MaxClients:   -1, // Ilimitado
            MaxProcesses: -1, // Ilimitado
            MaxUsers:     -1, // Ilimitado
        },
    }
    
    return limits[plan]
}
```

### **3. ✅ CENÁRIO: DENTRO DO LIMITE**

```go
// Se cliente 19/20 - PERMITIDO
func (s *ClientService) CreateClient(ctx context.Context, req *CreateClientRequest) (*Client, error) {
    // 1. Validar quota
    if err := s.ValidateClientCreation(ctx, req.UserID); err != nil {
        return nil, err
    }
    
    // 2. Criar cliente
    client := &Client{
        ID:        uuid.New().String(),
        TenantID:  req.TenantID,
        Name:      req.Name,
        Email:     req.Email,
        Phone:     req.Phone,
        Document:  req.Document,
        CreatedBy: req.UserID,
        CreatedAt: time.Now(),
    }
    
    // 3. Salvar no banco
    if err := s.clientRepository.CreateClient(ctx, client); err != nil {
        return nil, err
    }
    
    // 4. Atualizar métricas de uso
    s.usageService.UpdateClientUsage(ctx, req.TenantID, 1)
    
    // 5. Log de auditoria
    s.auditService.LogClientCreation(ctx, req.UserID, client.ID)
    
    return client, nil
}
```

### **4. 🚨 CENÁRIO: LIMITE EXCEDIDO**

```
👤 Dr. Silva: "Nome: João Silva, Email: joao@email.com"
🤖 Luxia: "❌ Limite atingido!
          
          📊 Seu plano Starter permite apenas 20 clientes
          📈 Você já possui 20 clientes cadastrados
          
          🔄 Opções disponíveis:
          1. Fazer upgrade para Professional (100 clientes)
          2. Remover cliente inativo
          3. Entrar em contato com suporte
          
          Deseja fazer upgrade do plano?"
```

### **5. 🔄 FLUXO DE UPGRADE**

```go
// internal/application/billing_service.go
func (s *BillingService) ProcessPlanUpgrade(ctx context.Context, userID, newPlan string) error {
    // 1. Obter assinatura atual
    subscription, err := s.subscriptionRepository.GetActiveSubscription(ctx, userID)
    if err != nil {
        return err
    }
    
    // 2. Calcular diferença de preço (prorate)
    priceDiff := s.calculateProratedUpgrade(subscription, newPlan)
    
    // 3. Processar pagamento da diferença
    if priceDiff > 0 {
        payment, err := s.paymentService.ProcessUpgradePayment(ctx, userID, priceDiff)
        if err != nil {
            return err
        }
        
        // Aguardar confirmação do pagamento
        if !payment.IsConfirmed {
            return fmt.Errorf("aguardando confirmação do pagamento")
        }
    }
    
    // 4. Atualizar plano
    subscription.Plan = newPlan
    subscription.UpdatedAt = time.Now()
    
    if err := s.subscriptionRepository.UpdateSubscription(ctx, subscription); err != nil {
        return err
    }
    
    // 5. Atualizar quotas do usuário
    s.quotaService.UpdateUserQuotas(ctx, userID, newPlan)
    
    // 6. Notificar usuário
    s.notificationService.NotifyPlanUpgrade(ctx, userID, newPlan)
    
    return nil
}
```

### **6. 💬 RESPOSTA PÓS-UPGRADE**

```
👤 Dr. Silva: "Sim, fazer upgrade"
🤖 Luxia: "✅ Upgrade processado com sucesso!
          
          📈 Plano atualizado: Starter → Professional
          💰 Cobrança adicional: R$ 67,00 (proporcional)
          
          📊 Novos limites:
          • Clientes: 18/100 ✨
          • Processos: 32/200 ✨
          • Usuários: 2/10 ✨
          
          Agora podemos cadastrar o cliente João Silva!
          Continuar com o cadastro?"
```

---

## 📋 **FLUXO DETALHADO: ADICIONAR PROCESSO**

### **🔄 Fluxo com Validação de Quota**

```go
// internal/application/process_service.go
func (s *ProcessService) ValidateProcessCreation(ctx context.Context, userID string) error {
    user, err := s.userRepository.GetUser(ctx, userID)
    if err != nil {
        return err
    }
    
    planLimits := s.getPlanLimits(user.Plan)
    currentProcesses, err := s.processRepository.CountProcessesByTenant(ctx, user.TenantID)
    if err != nil {
        return err
    }
    
    if currentProcesses >= planLimits.MaxProcesses {
        return &QuotaExceededError{
            Current: currentProcesses,
            Limit:   planLimits.MaxProcesses,
            Plan:    user.Plan,
            Resource: "processes",
        }
    }
    
    return nil
}

// CreateProcess com validação de quota
func (s *ProcessService) CreateProcess(ctx context.Context, req *CreateProcessRequest) (*Process, error) {
    // 1. Validar quota
    if err := s.ValidateProcessCreation(ctx, req.UserID); err != nil {
        if quotaErr, ok := err.(*QuotaExceededError); ok {
            return nil, s.handleQuotaExceeded(ctx, quotaErr)
        }
        return nil, err
    }
    
    // 2. Validar número CNJ
    if !s.validateCNJNumber(req.ProcessNumber) {
        return nil, fmt.Errorf("número CNJ inválido")
    }
    
    // 3. Verificar se processo já existe
    if exists, err := s.processRepository.ProcessExists(ctx, req.ProcessNumber); err != nil {
        return nil, err
    } else if exists {
        return nil, fmt.Errorf("processo já cadastrado")
    }
    
    // 4. Consultar dados iniciais no DataJud
    processData, err := s.dataJudService.GetProcessData(ctx, req.ProcessNumber)
    if err != nil {
        return nil, fmt.Errorf("erro ao consultar processo no DataJud: %w", err)
    }
    
    // 5. Criar processo
    process := &Process{
        ID:           uuid.New().String(),
        TenantID:     req.TenantID,
        Number:       req.ProcessNumber,
        ClientID:     req.ClientID,
        Court:        processData.Court,
        Subject:      processData.Subject,
        Status:       processData.Status,
        CreatedBy:    req.UserID,
        CreatedAt:    time.Now(),
        IsMonitored:  true,
    }
    
    // 6. Salvar processo
    if err := s.processRepository.CreateProcess(ctx, process); err != nil {
        return nil, err
    }
    
    // 7. Atualizar métricas
    s.usageService.UpdateProcessUsage(ctx, req.TenantID, 1)
    
    // 8. Iniciar monitoramento
    s.monitorService.StartProcessMonitoring(ctx, process.ID)
    
    return process, nil
}
```

### **💬 Fluxo via Luxia**

```
👤 Dr. Silva: "Adicionar processo 1001234-56.2024.8.26.0100"
🤖 Luxia: "📋 Vou adicionar este processo...
          
          🔍 Consultando dados no DataJud...
          ✅ Processo encontrado: Ação de Cobrança - TJSP
          
          📊 Status de quota:
          • Processos: 49/50 (Starter)
          
          ⚠️ ATENÇÃO: Você está próximo do limite!
          Após este processo, terá usado 100% da quota.
          
          Deseja continuar?"

👤 Dr. Silva: "Sim"
🤖 Luxia: "✅ Processo adicionado com sucesso!
          
          📋 Processo: 1001234-56.2024.8.26.0100
          🏛️ Tribunal: TJSP
          📝 Assunto: Ação de Cobrança
          🔔 Monitoramento: Ativado
          
          📊 Quota atualizada: 50/50 (100% utilizada)
          
          💡 Próximo processo exigirá upgrade do plano!"
```

---

## 📋 **FLUXO DETALHADO: ADICIONAR USUÁRIO**

### **🔄 Fluxo com Controle de Usuários**

```go
// internal/application/user_service.go
func (s *UserService) ValidateUserCreation(ctx context.Context, tenantID string) error {
    // 1. Obter dados do tenant
    tenant, err := s.tenantRepository.GetTenant(ctx, tenantID)
    if err != nil {
        return err
    }
    
    // 2. Obter limites do plano
    planLimits := s.getPlanLimits(tenant.Plan)
    
    // 3. Contar usuários atuais
    currentUsers, err := s.userRepository.CountUsersByTenant(ctx, tenantID)
    if err != nil {
        return err
    }
    
    // 4. Verificar limite
    if currentUsers >= planLimits.MaxUsers {
        return &QuotaExceededError{
            Current: currentUsers,
            Limit:   planLimits.MaxUsers,
            Plan:    tenant.Plan,
            Resource: "users",
        }
    }
    
    return nil
}

// InviteUser com validação de quota
func (s *UserService) InviteUser(ctx context.Context, req *InviteUserRequest) (*UserInvitation, error) {
    // 1. Validar quota
    if err := s.ValidateUserCreation(ctx, req.TenantID); err != nil {
        return nil, err
    }
    
    // 2. Verificar se email já existe
    if exists, err := s.userRepository.EmailExists(ctx, req.Email); err != nil {
        return nil, err
    } else if exists {
        return nil, fmt.Errorf("email já cadastrado")
    }
    
    // 3. Criar convite
    invitation := &UserInvitation{
        ID:        uuid.New().String(),
        TenantID:  req.TenantID,
        Email:     req.Email,
        Role:      req.Role,
        InvitedBy: req.InvitedBy,
        Token:     s.generateInviteToken(),
        ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 dias
        CreatedAt: time.Now(),
    }
    
    // 4. Salvar convite
    if err := s.invitationRepository.CreateInvitation(ctx, invitation); err != nil {
        return nil, err
    }
    
    // 5. Enviar email de convite
    s.emailService.SendInvitation(ctx, invitation)
    
    return invitation, nil
}
```

### **💬 Fluxo via Luxia**

```
👤 Dr. Silva: "Convidar usuário joao@escritorio.com como funcionário"
🤖 Luxia: "📋 Vou processar o convite...
          
          📊 Status atual:
          • Usuários: 2/3 (Starter)
          
          ✅ Convite enviado com sucesso!
          
          📧 Email: joao@escritorio.com
          👤 Função: Funcionário
          ⏰ Válido por: 7 dias
          
          📊 Quota atualizada: 3/3 (100% utilizada)
          
          ⚠️ Limite atingido! Próximo convite exigirá upgrade."
```

---

## 📋 **SISTEMA DE NOTIFICAÇÕES DE QUOTA**

### **🔔 Alertas Automáticos**

```go
// internal/application/quota_service.go
func (s *QuotaService) CheckQuotaAlerts(ctx context.Context, tenantID string) {
    usage := s.getCurrentUsage(ctx, tenantID)
    
    // Alertas por percentual de uso
    alerts := []QuotaAlert{
        {Threshold: 80, Message: "Você está usando 80% da quota"},
        {Threshold: 90, Message: "Você está usando 90% da quota"},
        {Threshold: 95, Message: "Você está usando 95% da quota"},
        {Threshold: 100, Message: "Quota esgotada! Upgrade necessário"},
    }
    
    for _, alert := range alerts {
        if s.shouldSendAlert(usage, alert) {
            s.sendQuotaAlert(ctx, tenantID, alert)
        }
    }
}

// sendQuotaAlert envia alerta via Luxia
func (s *QuotaService) sendQuotaAlert(ctx context.Context, tenantID string, alert QuotaAlert) {
    // Obter dados do tenant
    tenant := s.tenantRepository.GetTenant(ctx, tenantID)
    
    // Obter usuário principal (advogado)
    mainUser := s.userRepository.GetMainUserByTenant(ctx, tenantID)
    
    // Montar mensagem
    message := fmt.Sprintf(`
🚨 *Alerta de Quota*

📊 *%s*

Recursos utilizados:
• Clientes: %d/%d
• Processos: %d/%d  
• Usuários: %d/%d

💡 *Soluções:*
• Fazer upgrade do plano
• Remover recursos inativos
• Entrar em contato com suporte

Deseja fazer upgrade agora?`,
        alert.Message,
        usage.Clients, usage.MaxClients,
        usage.Processes, usage.MaxProcesses,
        usage.Users, usage.MaxUsers)
    
    // Enviar via WhatsApp
    s.notificationService.SendLuxiaMessage(ctx, mainUser.Phone, message)
}
```

---

## 📋 **FLUXO DETALHADO: FUNCIONALIDADES AVANÇADAS**

### **🔍 Busca de Jurisprudência (por plano)**

```go
// internal/application/jurisprudence_service.go
func (s *JurisprudenceService) SearchJurisprudence(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
    // 1. Verificar se plano permite
    user, err := s.userRepository.GetUser(ctx, req.UserID)
    if err != nil {
        return nil, err
    }
    
    planFeatures := s.getPlanFeatures(user.Plan)
    if !planFeatures.JurisprudenceSearch {
        return nil, &FeatureNotAvailableError{
            Feature: "jurisprudence_search",
            Plan:    user.Plan,
            RequiredPlans: []string{"professional", "business", "enterprise"},
        }
    }
    
    // 2. Verificar quota de busca
    dailySearches, err := s.usageRepository.GetDailySearches(ctx, user.TenantID)
    if err != nil {
        return nil, err
    }
    
    if dailySearches >= planFeatures.MaxDailySearches {
        return nil, &QuotaExceededError{
            Current: dailySearches,
            Limit:   planFeatures.MaxDailySearches,
            Resource: "daily_searches",
        }
    }
    
    // 3. Executar busca
    results, err := s.searchEngine.Search(ctx, req.Query)
    if err != nil {
        return nil, err
    }
    
    // 4. Filtrar resultados por plano
    filteredResults := s.filterResultsByPlan(results, planFeatures)
    
    // 5. Atualizar contador
    s.usageRepository.IncrementDailySearches(ctx, user.TenantID)
    
    return &SearchResponse{
        Results: filteredResults,
        Total:   len(filteredResults),
    }, nil
}
```

### **💬 Controle via Luxia**

```
👤 Dr. Silva: "Buscar jurisprudência sobre danos morais"
🤖 Luxia: "❌ Funcionalidade não disponível
          
          📊 Seu plano: Starter
          🔍 Busca de jurisprudência: Não incluída
          
          📈 Disponível nos planos:
          • Professional (R$ 299/mês)
          • Business (R$ 699/mês)
          • Enterprise (R$ 1.999/mês)
          
          Deseja fazer upgrade?"
```

---

## 📋 **ESTRUTURA DE CONTROLE UNIFICADA**

### **🏗️ Middleware de Validação**

```go
// internal/infrastructure/middleware/quota_middleware.go
func QuotaValidationMiddleware(quotaService QuotaService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        action := c.GetString("action") // "create_client", "create_process", etc.
        
        // Validar quota para a ação
        if err := quotaService.ValidateAction(c.Request.Context(), userID, action); err != nil {
            if quotaErr, ok := err.(*QuotaExceededError); ok {
                c.JSON(http.StatusForbidden, gin.H{
                    "error": "quota_exceeded",
                    "message": quotaErr.Message,
                    "current": quotaErr.Current,
                    "limit": quotaErr.Limit,
                    "plan": quotaErr.Plan,
                    "upgrade_options": quotaErr.UpgradeOptions,
                })
                c.Abort()
                return
            }
            
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### **📊 Dashboard de Quotas**

```go
// internal/application/dashboard_service.go
func (s *DashboardService) GetQuotaStatus(ctx context.Context, tenantID string) *QuotaStatus {
    tenant := s.tenantRepository.GetTenant(ctx, tenantID)
    usage := s.usageRepository.GetCurrentUsage(ctx, tenantID)
    limits := s.getPlanLimits(tenant.Plan)
    
    return &QuotaStatus{
        Plan: tenant.Plan,
        Usage: map[string]ResourceUsage{
            "clients": {
                Current: usage.Clients,
                Limit:   limits.MaxClients,
                Percentage: float64(usage.Clients) / float64(limits.MaxClients) * 100,
            },
            "processes": {
                Current: usage.Processes,
                Limit:   limits.MaxProcesses,
                Percentage: float64(usage.Processes) / float64(limits.MaxProcesses) * 100,
            },
            "users": {
                Current: usage.Users,
                Limit:   limits.MaxUsers,
                Percentage: float64(usage.Users) / float64(limits.MaxUsers) * 100,
            },
        },
        Alerts: s.getQuotaAlerts(usage, limits),
        UpgradeOptions: s.getUpgradeOptions(tenant.Plan),
    }
}
```

---

## ✅ **RESUMO DOS FLUXOS DOCUMENTADOS**

### **🔄 Fluxos Completos:**
1. **✅ Adicionar Cliente** - Validação quota, limites por plano, upgrade automático
2. **✅ Adicionar Processo** - Consulta DataJud, validação CNJ, controle de quota
3. **✅ Convidar Usuário** - Controle de usuários por plano, convites com expiração
4. **✅ Buscar Jurisprudência** - Funcionalidade por plano, quota diária
5. **✅ Alertas Automáticos** - Notificações 80%, 90%, 95%, 100%
6. **✅ Upgrade de Plano** - Cobrança proporcional, atualização instantânea

### **🎯 Controles Implementados:**
- **Validação prévia** - Antes de qualquer ação
- **Mensagens claras** - Usuário sabe exatamente o limite
- **Upgrade fluido** - Sem interrupção do serviço
- **Auditoria completa** - Todos os recursos são rastreados

---

## 📚 **ARQUIVO CRIADO:**
`FLUXOS_CONTROLE_QUOTAS_PLANOS.md` - **Fluxos detalhados** com código real e exemplos práticos

**🎯 Agora está claro como cada funcionalidade flui com controle de quotas por plano!**

**Algum fluxo específico precisa de mais detalhamento?**