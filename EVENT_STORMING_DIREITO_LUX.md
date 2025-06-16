# Event Storming - Direito Lux

## Workshop Overview
**Objetivo**: Mapear o domínio do Direito Lux e identificar bounded contexts, agregados e eventos

**Participantes**: Product Owner, Tech Lead, Stakeholders jurídicos

**Duração**: 4-6 horas (pode ser dividido em 2 sessões)

## 🟠 Domain Events (Eventos de Domínio)

### Cliente/User Journey
```
🟠 UserRegistered
🟠 UserActivated  
🟠 UserLoggedIn
🟠 UserSubscriptionChanged
🟠 UserDeactivated
```

### Processo Jurídico Lifecycle
```
🟠 ProcessRequested          // Cliente solicita monitoramento
🟠 ProcessValidated          // Número do processo validado
🟠 ProcessRegistered         // Processo cadastrado no sistema
🟠 ProcessMonitoringStarted  // Monitoramento ativo iniciado
🟠 ProcessDataFetched        // Dados obtidos do DataJud
🟠 ProcessUpdated            // Nova movimentação detectada
🟠 ProcessMovementDetected   // Específico para movimentações
🟠 ProcessDeadlineDetected   // Prazo identificado
🟠 ProcessStatusChanged      // Mudança de status processual
🟠 ProcessConcluded          // Processo finalizado
🟠 ProcessMonitoringStopped  // Monitoramento pausado/parado
🟠 ProcessArchived           // Processo arquivado
```

### Notificações
```
🟠 NotificationRequested     // Sistema solicita envio
🟠 NotificationScheduled     // Agendada para envio
🟠 NotificationSent          // Enviada com sucesso
🟠 NotificationFailed        // Falha no envio
🟠 NotificationDelivered     // Confirmação de entrega
🟠 NotificationRead          // Lida pelo destinatário
🟠 NotificationRetryScheduled // Reagendada para retry
```

### Integração DataJud
```
🟠 DataJudQueryRequested     // Consulta solicitada
🟠 DataJudQueryExecuted      // Consulta executada
🟠 DataJudDataReceived       // Dados recebidos
🟠 DataJudRateLimitReached   // Limite atingido
🟠 DataJudErrorOccurred      // Erro na consulta
🟠 DataJudCacheHit           // Cache utilizado
🟠 DataJudCacheMiss          // Cache miss
```

### Tenant Management
```
🟠 TenantCreated            // Novo escritório cadastrado
🟠 TenantActivated          // Escritório ativado
🟠 TenantQuotaUpdated       // Cota alterada
🟠 TenantQuotaExceeded      // Cota excedida
🟠 TenantBillingCycleStarted // Novo ciclo de cobrança
🟠 TenantSuspended          // Escritório suspenso
🟠 TenantReactivated        // Escritório reativado
```

### Inteligência Artificial
```
🟠 AISummarizationRequested  // Resumo solicitado
🟠 AISummaryGenerated        // Resumo gerado
🟠 AITermExplanationRequested // Explicação solicitada
🟠 AITermExplained           // Termo explicado
🟠 AIAnalysisCompleted       // Análise concluída
🟠 AIModelTrainingStarted    // Treinamento iniciado
```

### Documentos
```
🟠 DocumentGenerationRequested // Geração solicitada
🟠 DocumentGenerated           // Documento gerado
🟠 DocumentSigned              // Documento assinado
🟠 DocumentSent                // Documento enviado
🟠 DocumentArchived            // Documento arquivado
```

## 🔵 Commands (Comandos)

### User Commands
```
🔵 RegisterUser
🔵 ActivateUser
🔵 LoginUser
🔵 ChangeUserSubscription
🔵 DeactivateUser
```

### Process Commands
```
🔵 RequestProcessMonitoring
🔵 ValidateProcessNumber
🔵 RegisterProcess
🔵 StartMonitoring
🔵 StopMonitoring
🔵 UpdateProcessData
🔵 ArchiveProcess
🔵 SetProcessDeadline
```

### Notification Commands
```
🔵 SendNotification
🔵 ScheduleNotification
🔵 RetryNotification
🔵 CancelNotification
🔵 UpdateNotificationTemplate
```

### DataJud Commands
```
🔵 QueryDataJud
🔵 FetchProcessData
🔵 RefreshCache
🔵 UpdateRateLimit
```

### AI Commands
```
🔵 SummarizeProcess
🔵 ExplainLegalTerm
🔵 AnalyzeDocument
🔵 TrainModel
```

## 🟡 Aggregates (Agregados)

### 1. User Aggregate
```go
type User struct {
    ID           UserID
    Email        string
    TenantID     TenantID
    Role         UserRole
    Subscription SubscriptionType
    Status       UserStatus
    CreatedAt    time.Time
    LastLoginAt  time.Time
}
```

### 2. Process Aggregate
```go
type Process struct {
    ID              ProcessID
    Number          ProcessNumber
    TenantID        TenantID
    ClientID        ClientID
    Court           Court
    Status          ProcessStatus
    MonitoringState MonitoringState
    Movements       []ProcessMovement
    Deadlines       []Deadline
    CreatedAt       time.Time
    LastUpdatedAt   time.Time
}
```

### 3. Notification Aggregate
```go
type Notification struct {
    ID          NotificationID
    TenantID    TenantID
    ProcessID   ProcessID
    RecipientID UserID
    Channel     NotificationChannel
    Template    MessageTemplate
    Status      NotificationStatus
    Attempts    int
    CreatedAt   time.Time
    SentAt      *time.Time
}
```

### 4. Tenant Aggregate
```go
type Tenant struct {
    ID           TenantID
    Name         string
    Subscription SubscriptionPlan
    Quotas       TenantQuotas
    Status       TenantStatus
    BillingInfo  BillingInfo
    CreatedAt    time.Time
}
```

### 5. DataJudQuery Aggregate
```go
type DataJudQuery struct {
    ID           QueryID
    TenantID     TenantID
    ProcessNumber ProcessNumber
    QueryType    QueryType
    Status       QueryStatus
    RateLimit    RateLimit
    CacheStatus  CacheStatus
    ExecutedAt   time.Time
}
```

## 🟢 Views/Read Models (Projeções)

### Process Dashboard View
```go
type ProcessDashboardView struct {
    TenantID        TenantID
    TotalProcesses  int
    ActiveProcesses int
    RecentUpdates   []ProcessUpdate
    PendingDeadlines []Deadline
    QuotaUsage      QuotaUsage
}
```

### Notification History View
```go
type NotificationHistoryView struct {
    TenantID         TenantID
    TotalSent        int
    DeliveryRate     float64
    FailedByChannel  map[Channel]int
    RecentFailures   []FailedNotification
}
```

### Analytics View
```go
type TenantAnalyticsView struct {
    TenantID           TenantID
    ProcessesPerMonth  map[string]int
    NotificationsPerDay map[string]int
    TopClients         []ClientStats
    ApiUsage           ApiUsageStats
}
```

## 🔴 Bounded Contexts

### 1. Authentication & Authorization Context
**Responsabilidade**: Gestão de usuários, autenticação e autorização
```
Aggregates: User, Role, Permission
Events: UserRegistered, UserLoggedIn, RoleAssigned
Services: AuthService, PermissionService
```

### 2. Tenant Management Context  
**Responsabilidade**: Gestão de tenants, quotas e billing
```
Aggregates: Tenant, Subscription, Quota
Events: TenantCreated, QuotaExceeded, BillingCycleStarted
Services: TenantService, QuotaService, BillingService
```

### 3. Process Management Context
**Responsabilidade**: Gestão do ciclo de vida dos processos
```
Aggregates: Process, ProcessMovement, Deadline
Events: ProcessRegistered, ProcessUpdated, DeadlineDetected
Services: ProcessService, MonitoringService
```

### 4. External Integration Context (DataJud)
**Responsabilidade**: Integração com APIs externas
```
Aggregates: DataJudQuery, Cache, RateLimit
Events: DataJudQueryExecuted, CacheHit, RateLimitReached
Services: DataJudService, CacheService, RateLimitService
```

### 5. Notification Context
**Responsabilidade**: Orquestração e envio de notificações
```
Aggregates: Notification, NotificationTemplate, Channel
Events: NotificationSent, NotificationFailed, NotificationDelivered
Services: NotificationService, WhatsAppService, EmailService
```

### 6. AI & Analytics Context
**Responsabilidade**: Inteligência artificial e análises
```
Aggregates: AISummary, AnalysisResult, Model
Events: SummaryGenerated, AnalysisCompleted, ModelTrained
Services: AIService, SummarizationService, AnalyticsService
```

### 7. Document Management Context
**Responsabilidade**: Geração e gestão de documentos
```
Aggregates: Document, Template, Signature
Events: DocumentGenerated, DocumentSigned, DocumentArchived
Services: DocumentService, TemplateService, SignatureService
```

## 🔀 Context Map

```
[Authentication] --> [Tenant Management] : User belongs to Tenant
[Tenant Management] --> [Process Management] : Tenant owns Processes
[Process Management] --> [External Integration] : Process data from DataJud
[Process Management] --> [Notification] : Process changes trigger notifications
[Process Management] --> [AI & Analytics] : Process data for AI analysis
[Notification] --> [Authentication] : User preferences for notifications
[AI & Analytics] --> [Document Management] : AI generates documents
[Tenant Management] --> [All Contexts] : Quota enforcement
```

## 📋 Saga Patterns Identificados

### 1. Process Registration Saga
```
1. ValidateProcessNumber
2. CheckTenantQuota
3. RegisterProcess
4. StartMonitoring
5. SendConfirmationNotification

Compensations:
- RemoveProcess if monitoring fails
- RestoreQuota if process invalid
```

### 2. Notification Delivery Saga
```
1. PrepareNotification
2. SelectChannel (WhatsApp -> Email -> SMS)
3. SendNotification
4. ConfirmDelivery
5. UpdateStatistics

Compensations:
- Retry with different channel
- Mark as failed after max attempts
```

### 3. Tenant Onboarding Saga
```
1. CreateTenant
2. SetupDefaultQuotas
3. CreateAdminUser
4. SendWelcomeNotification
5. ActivateTenant

Compensations:
- Cleanup if any step fails
- Rollback tenant creation
```

## 🎯 Key Insights do Event Storming

### 1. Bounded Contexts bem definidos
- **Authentication**: Isolado e reutilizável
- **Process Management**: Core domain 
- **Notification**: Complex orchestration needed
- **AI**: Separate context for ML concerns

### 2. Critical Integration Points
- **DataJud API**: Rate limiting critical
- **WhatsApp Business**: Primary channel
- **Multi-tenant**: Isolation everywhere

### 3. Event-driven Opportunities
- **Process monitoring**: Event-sourcing ideal
- **Notifications**: Pub/Sub pattern
- **Analytics**: Event streaming for real-time

### 4. Consistency Boundaries
- **Strong consistency**: Within aggregates
- **Eventual consistency**: Between contexts
- **Sagas**: For cross-context transactions

## 📖 Ubiquitous Language

### Core Terms
- **Processo**: Processo jurídico a ser monitorado
- **Movimentação**: Atualização/mudança no processo
- **Prazo**: Data limite para ação processual
- **Tenant**: Escritório de advocacia (inquilino)
- **Monitoramento**: Verificação automática de mudanças
- **Notificação**: Comunicação enviada ao usuário
- **Quota**: Limite de uso por tenant
- **DataJud**: API oficial do CNJ para dados processuais

### Domain-specific Terms
- **Número CNJ**: Identificador único nacional do processo
- **Tribunal**: Órgão julgador (TJ, TRF, etc.)
- **Grau**: Instância processual (1º, 2º grau)
- **Classe**: Tipo de ação judicial
- **Assunto**: Matéria jurídica do processo

## 🔧 Technical Decisions from Event Storming

### 1. Event Store Strategy
- **Process Context**: Full event sourcing
- **Other Contexts**: Event-driven with projections
- **Snapshots**: Every 50 events for performance

### 2. Consistency Strategy
- **Command side**: Strong consistency
- **Query side**: Eventual consistency acceptable
- **SLA**: 5 seconds max for notifications

### 3. Data Partitioning
- **Tenant-based**: All data partitioned by TenantID
- **Time-based**: Events partitioned by month
- **Geographic**: Future multi-region consideration

### 4. Integration Patterns
- **DataJud**: Publish-subscribe with cache
- **WhatsApp**: Request-response with retry
- **Internal**: Event-driven with message bus