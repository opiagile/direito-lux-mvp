# Domain Events - Direito Lux

## Event-driven Architecture Overview

O Direito Lux utiliza uma arquitetura orientada a eventos para garantir baixo acoplamento entre bounded contexts e alta escalabilidade. Todos os eventos seguem um formato padronizado e são publicados via Pub/Sub.

## Event Schema Padrão

```go
type DomainEvent interface {
    GetEventID() string
    GetEventType() string
    GetAggregateID() string
    GetTenantID() string
    GetOccurredAt() time.Time
    GetVersion() int
    GetMetadata() EventMetadata
}

type BaseEvent struct {
    EventID     string        `json:"event_id"`
    EventType   string        `json:"event_type"`
    AggregateID string        `json:"aggregate_id"`
    TenantID    string        `json:"tenant_id"`
    OccurredAt  time.Time     `json:"occurred_at"`
    Version     int           `json:"version"`
    Metadata    EventMetadata `json:"metadata"`
}

type EventMetadata struct {
    CausationID   string            `json:"causation_id,omitempty"`
    CorrelationID string            `json:"correlation_id,omitempty"`
    UserID        string            `json:"user_id,omitempty"`
    Source        string            `json:"source"`
    Headers       map[string]string `json:"headers,omitempty"`
}
```

## 1. Authentication Context Events

### UserRegistered
```go
type UserRegistered struct {
    BaseEvent
    Payload UserRegisteredPayload `json:"payload"`
}

type UserRegisteredPayload struct {
    UserID      string    `json:"user_id"`
    Email       string    `json:"email"`
    TenantID    string    `json:"tenant_id"`
    Role        string    `json:"role"`
    FirstName   string    `json:"first_name"`
    LastName    string    `json:"last_name"`
    IsActive    bool      `json:"is_active"`
    ActivationToken string `json:"activation_token,omitempty"`
}
```

**Consumers:**
- Tenant Service: Atualizar contagem de usuários
- Notification Service: Enviar email de boas-vindas
- Analytics Service: Registrar novo usuário

### UserLoggedIn
```go
type UserLoggedIn struct {
    BaseEvent
    Payload UserLoggedInPayload `json:"payload"`
}

type UserLoggedInPayload struct {
    UserID     string `json:"user_id"`
    SessionID  string `json:"session_id"`
    IPAddress  string `json:"ip_address"`
    UserAgent  string `json:"user_agent"`
    LoginMethod string `json:"login_method"` // password, oauth, sso
}
```

**Consumers:**
- Analytics Service: Registrar atividade de login
- Security Service: Detectar logins suspeitos

### UserPermissionChanged
```go
type UserPermissionChanged struct {
    BaseEvent
    Payload UserPermissionChangedPayload `json:"payload"`
}

type UserPermissionChangedPayload struct {
    UserID      string   `json:"user_id"`
    OldRole     string   `json:"old_role"`
    NewRole     string   `json:"new_role"`
    Permissions []string `json:"permissions"`
    ChangedBy   string   `json:"changed_by"`
}
```

## 2. Tenant Management Events

### TenantCreated
```go
type TenantCreated struct {
    BaseEvent
    Payload TenantCreatedPayload `json:"payload"`
}

type TenantCreatedPayload struct {
    TenantID     string    `json:"tenant_id"`
    Name         string    `json:"name"`
    Document     string    `json:"document"` // CNPJ
    Plan         string    `json:"plan"`
    AdminUserID  string    `json:"admin_user_id"`
    BillingEmail string    `json:"billing_email"`
    Address      Address   `json:"address"`
    TrialEndsAt  time.Time `json:"trial_ends_at"`
}

type Address struct {
    Street     string `json:"street"`
    Number     string `json:"number"`
    City       string `json:"city"`
    State      string `json:"state"`
    PostalCode string `json:"postal_code"`
    Country    string `json:"country"`
}
```

**Consumers:**
- Process Service: Configurar quotas iniciais
- Notification Service: Enviar welcome package
- Document Service: Criar templates padrão
- Billing Service: Setup billing account

### TenantSubscriptionUpgraded
```go
type TenantSubscriptionUpgraded struct {
    BaseEvent
    Payload TenantSubscriptionUpgradedPayload `json:"payload"`
}

type TenantSubscriptionUpgradedPayload struct {
    TenantID      string    `json:"tenant_id"`
    FromPlan      string    `json:"from_plan"`
    ToPlan        string    `json:"to_plan"`
    EffectiveAt   time.Time `json:"effective_at"`
    ProrationAmount int64   `json:"proration_amount"` // centavos
    NewQuotas     TenantQuotas `json:"new_quotas"`
}

type TenantQuotas struct {
    ProcessLimit      int `json:"process_limit"`
    UsersLimit        int `json:"users_limit"`
    DataJudDailyQuota int `json:"datajud_daily_quota"`
    StorageQuotaGB    int `json:"storage_quota_gb"`
    NotificationsMonth int `json:"notifications_month"`
}
```

**Consumers:**
- Process Service: Atualizar quotas
- Billing Service: Processar cobrança pro-rata
- Analytics Service: Track upgrade
- Notification Service: Confirmar upgrade

### TenantQuotaExceeded
```go
type TenantQuotaExceeded struct {
    BaseEvent
    Payload TenantQuotaExceededPayload `json:"payload"`
}

type TenantQuotaExceededPayload struct {
    TenantID    string `json:"tenant_id"`
    QuotaType   string `json:"quota_type"` // processes, users, datajud_calls, storage
    Limit       int    `json:"limit"`
    CurrentUsage int   `json:"current_usage"`
    Percentage   float64 `json:"percentage"`
    Action      string `json:"action"` // warn, block, upgrade_suggest
}
```

**Consumers:**
- Notification Service: Alertar administrador
- Process Service: Bloquear novos cadastros se necessário
- Analytics Service: Track quota usage patterns

## 3. Process Management Events

### ProcessRegistered
```go
type ProcessRegistered struct {
    BaseEvent
    Payload ProcessRegisteredPayload `json:"payload"`
}

type ProcessRegisteredPayload struct {
    ProcessID     string    `json:"process_id"`
    Number        string    `json:"number"`
    TenantID      string    `json:"tenant_id"`
    ClientID      string    `json:"client_id"`
    Court         Court     `json:"court"`
    Classification Classification `json:"classification"`
    RegisteredBy  string    `json:"registered_by"`
    MonitoringEnabled bool  `json:"monitoring_enabled"`
}

type Court struct {
    Code     string `json:"code"`
    Name     string `json:"name"`
    Type     string `json:"type"` // tj, trf, tst, etc
    State    string `json:"state"`
    Instance int    `json:"instance"` // 1, 2
}

type Classification struct {
    Class   string `json:"class"`
    Subject string `json:"subject"`
    Area    string `json:"area"` // civil, criminal, trabalhista
}
```

**Consumers:**
- DataJud Service: Agendar primeira consulta
- Notification Service: Confirmar cadastro
- Analytics Service: Registrar novo processo
- Tenant Service: Incrementar contador de processos

### ProcessMovementDetected
```go
type ProcessMovementDetected struct {
    BaseEvent
    Payload ProcessMovementDetectedPayload `json:"payload"`
}

type ProcessMovementDetectedPayload struct {
    ProcessID    string    `json:"process_id"`
    MovementID   string    `json:"movement_id"`
    Date         time.Time `json:"date"`
    Description  string    `json:"description"`
    Type         string    `json:"type"`
    IsDecision   bool      `json:"is_decision"`
    HasDocument  bool      `json:"has_document"`
    DocumentURL  string    `json:"document_url,omitempty"`
    Source       string    `json:"source"` // datajud, tribunal_api, manual
    PreviousMovementID string `json:"previous_movement_id,omitempty"`
}
```

**Consumers:**
- Notification Service: Enviar alertas para interessados
- AI Service: Analisar e resumir movimentação
- Document Service: Baixar documentos anexos
- Analytics Service: Registrar atividade processual

### ProcessDeadlineDetected
```go
type ProcessDeadlineDetected struct {
    BaseEvent
    Payload ProcessDeadlineDetectedPayload `json:"payload"`
}

type ProcessDeadlineDetectedPayload struct {
    ProcessID     string    `json:"process_id"`
    DeadlineID    string    `json:"deadline_id"`
    Date          time.Time `json:"date"`
    Description   string    `json:"description"`
    Type          string    `json:"type"` // recurso, contestacao, manifestacao
    Priority      string    `json:"priority"` // low, medium, high, critical
    DaysRemaining int       `json:"days_remaining"`
    ResponsibleUserID string `json:"responsible_user_id,omitempty"`
    MovementID    string    `json:"movement_id"` // movimentação que gerou o prazo
}
```

**Consumers:**
- Notification Service: Agendar lembretes
- Calendar Service: Criar eventos no calendário
- Analytics Service: Track compliance com prazos

### ProcessStatusChanged
```go
type ProcessStatusChanged struct {
    BaseEvent
    Payload ProcessStatusChangedPayload `json:"payload"`
}

type ProcessStatusChangedPayload struct {
    ProcessID   string `json:"process_id"`
    OldStatus   string `json:"old_status"`
    NewStatus   string `json:"new_status"`
    Reason      string `json:"reason"`
    MovementID  string `json:"movement_id,omitempty"`
    ChangedBy   string `json:"changed_by"` // system, user_id
}
```

## 4. DataJud Integration Events

### DataJudQueryExecuted
```go
type DataJudQueryExecuted struct {
    BaseEvent
    Payload DataJudQueryExecutedPayload `json:"payload"`
}

type DataJudQueryExecutedPayload struct {
    QueryID       string        `json:"query_id"`
    ProcessNumber string        `json:"process_number"`
    QueryType     string        `json:"query_type"` // full, movements, documents
    Success       bool          `json:"success"`
    ResponseTime  time.Duration `json:"response_time"`
    CacheHit      bool          `json:"cache_hit"`
    RateLimitInfo RateLimitInfo `json:"rate_limit_info"`
    ErrorMessage  string        `json:"error_message,omitempty"`
    ResultSize    int           `json:"result_size"` // bytes
}

type RateLimitInfo struct {
    Remaining int       `json:"remaining"`
    Limit     int       `json:"limit"`
    ResetsAt  time.Time `json:"resets_at"`
}
```

**Consumers:**
- Process Service: Atualizar dados do processo
- Analytics Service: Monitor API usage
- Notification Service: Alertar sobre falhas
- Cache Service: Atualizar estratégias

### DataJudRateLimitExceeded
```go
type DataJudRateLimitExceeded struct {
    BaseEvent
    Payload DataJudRateLimitExceededPayload `json:"payload"`
}

type DataJudRateLimitExceededPayload struct {
    TenantID      string    `json:"tenant_id"`
    DailyLimit    int       `json:"daily_limit"`
    CurrentUsage  int       `json:"current_usage"`
    ResetsAt      time.Time `json:"resets_at"`
    BlockedQueries int      `json:"blocked_queries"`
    SuggestedAction string  `json:"suggested_action"`
}
```

**Consumers:**
- Tenant Service: Bloquear novas consultas
- Notification Service: Alertar administrador
- Process Service: Pausar monitoramento automático

## 5. Notification System Events

### NotificationRequested
```go
type NotificationRequested struct {
    BaseEvent
    Payload NotificationRequestedPayload `json:"payload"`
}

type NotificationRequestedPayload struct {
    NotificationID string             `json:"notification_id"`
    Type          string             `json:"type"` // process_update, deadline_reminder, quota_warning
    Priority      string             `json:"priority"` // low, normal, high, urgent
    RecipientID   string             `json:"recipient_id"`
    Channel       string             `json:"channel"` // whatsapp, email, telegram, sms
    TemplateID    string             `json:"template_id"`
    Variables     map[string]interface{} `json:"variables"`
    ScheduledFor  *time.Time         `json:"scheduled_for,omitempty"`
    ProcessID     string             `json:"process_id,omitempty"`
}
```

### NotificationSent
```go
type NotificationSent struct {
    BaseEvent
    Payload NotificationSentPayload `json:"payload"`
}

type NotificationSentPayload struct {
    NotificationID string        `json:"notification_id"`
    Channel       string        `json:"channel"`
    ExternalID    string        `json:"external_id"` // ID do provedor externo
    Attempt       int           `json:"attempt"`
    ResponseTime  time.Duration `json:"response_time"`
    Cost          int64         `json:"cost"` // centavos
}
```

### NotificationDelivered
```go
type NotificationDelivered struct {
    BaseEvent
    Payload NotificationDeliveredPayload `json:"payload"`
}

type NotificationDeliveredPayload struct {
    NotificationID string    `json:"notification_id"`
    Channel       string    `json:"channel"`
    DeliveredAt   time.Time `json:"delivered_at"`
    ReadAt        *time.Time `json:"read_at,omitempty"`
    ExternalID    string    `json:"external_id"`
}
```

### NotificationFailed
```go
type NotificationFailed struct {
    BaseEvent
    Payload NotificationFailedPayload `json:"payload"`
}

type NotificationFailedPayload struct {
    NotificationID string `json:"notification_id"`
    Channel       string `json:"channel"`
    Attempt       int    `json:"attempt"`
    ErrorCode     string `json:"error_code"`
    ErrorMessage  string `json:"error_message"`
    WillRetry     bool   `json:"will_retry"`
    NextRetryAt   *time.Time `json:"next_retry_at,omitempty"`
    FallbackChannel string `json:"fallback_channel,omitempty"`
}
```

## 6. AI & Analytics Events

### AISummarizationRequested
```go
type AISummarizationRequested struct {
    BaseEvent
    Payload AISummarizationRequestedPayload `json:"payload"`
}

type AISummarizationRequestedPayload struct {
    SummaryID    string `json:"summary_id"`
    ProcessID    string `json:"process_id"`
    Type         string `json:"type"` // process_overview, movement_summary, legal_analysis
    Language     string `json:"language"`
    TargetAudience string `json:"target_audience"` // lawyer, client
    RequestedBy  string `json:"requested_by"`
    Priority     string `json:"priority"`
}
```

### AISummaryGenerated
```go
type AISummaryGenerated struct {
    BaseEvent
    Payload AISummaryGeneratedPayload `json:"payload"`
}

type AISummaryGeneratedPayload struct {
    SummaryID    string  `json:"summary_id"`
    ProcessID    string  `json:"process_id"`
    Content      string  `json:"content"`
    Confidence   float64 `json:"confidence"`
    Model        string  `json:"model"`
    TokensUsed   int     `json:"tokens_used"`
    ProcessingTime time.Duration `json:"processing_time"`
    CostUSD      float64 `json:"cost_usd"`
}
```

### AnalysisCompleted
```go
type AnalysisCompleted struct {
    BaseEvent
    Payload AnalysisCompletedPayload `json:"payload"`
}

type AnalysisCompletedPayload struct {
    AnalysisID   string        `json:"analysis_id"`
    Type         string        `json:"type"` // jurimetria, sentiment, risk_assessment
    SubjectID    string        `json:"subject_id"` // process_id, tenant_id
    SubjectType  string        `json:"subject_type"`
    Result       interface{}   `json:"result"`
    Confidence   float64       `json:"confidence"`
    Duration     time.Duration `json:"duration"`
    ModelUsed    string        `json:"model_used"`
}
```

## 7. Document Management Events

### DocumentGenerationRequested
```go
type DocumentGenerationRequested struct {
    BaseEvent
    Payload DocumentGenerationRequestedPayload `json:"payload"`
}

type DocumentGenerationRequestedPayload struct {
    DocumentID   string                 `json:"document_id"`
    TemplateID   string                 `json:"template_id"`
    ProcessID    string                 `json:"process_id,omitempty"`
    Type         string                 `json:"type"` // peticao, contrato, parecer
    Format       string                 `json:"format"` // pdf, docx, html
    Variables    map[string]interface{} `json:"variables"`
    RequestedBy  string                 `json:"requested_by"`
    Priority     string                 `json:"priority"`
}
```

### DocumentGenerated
```go
type DocumentGenerated struct {
    BaseEvent
    Payload DocumentGeneratedPayload `json:"payload"`
}

type DocumentGeneratedPayload struct {
    DocumentID   string        `json:"document_id"`
    Name         string        `json:"name"`
    Type         string        `json:"type"`
    Format       string        `json:"format"`
    Size         int64         `json:"size"`
    StoragePath  string        `json:"storage_path"`
    Hash         string        `json:"hash"`
    GenerationTime time.Duration `json:"generation_time"`
    TemplateUsed string        `json:"template_used"`
}
```

## Event Publishing Strategy

### 1. Event Bus Configuration
```go
// Pub/Sub Topics by Context
const (
    TopicAuth         = "auth-events"
    TopicTenant       = "tenant-events"  
    TopicProcess      = "process-events"
    TopicDataJud      = "datajud-events"
    TopicNotification = "notification-events"
    TopicAI           = "ai-events"
    TopicDocument     = "document-events"
    TopicAnalytics    = "analytics-events"
)

// Dead Letter Queue para eventos com falha
const (
    TopicDeadLetter = "dead-letter-events"
)
```

### 2. Event Ordering
```go
// Eventos ordenados por Aggregate ID
type OrderedEvent struct {
    Event       DomainEvent `json:"event"`
    OrderingKey string      `json:"ordering_key"` // tenant_id:aggregate_id
    Partition   int         `json:"partition"`
}
```

### 3. Event Replay
```go
// Capacidade de replay para recovery
type EventReplay struct {
    FromTimestamp time.Time `json:"from_timestamp"`
    ToTimestamp   time.Time `json:"to_timestamp"`
    EventTypes    []string  `json:"event_types"`
    TenantID      string    `json:"tenant_id,omitempty"`
    AggregateID   string    `json:"aggregate_id,omitempty"`
}
```

### 4. Monitoring & Observability
```go
// Métricas por evento
type EventMetrics struct {
    EventType     string        `json:"event_type"`
    TenantID      string        `json:"tenant_id"`
    PublishedAt   time.Time     `json:"published_at"`
    ProcessedAt   time.Time     `json:"processed_at"`
    Latency       time.Duration `json:"latency"`
    Success       bool          `json:"success"`
    ConsumerCount int           `json:"consumer_count"`
}
```

## Event Subscription Patterns

### 1. Process Saga Subscriptions
```go
// Process Registration Saga
ProcessRegistered → ValidateWithDataJud
DataJudQueryExecuted → StartMonitoring  
ProcessMonitoringStarted → SendConfirmation

// Notification Delivery Saga
NotificationRequested → AttemptWhatsApp
NotificationFailed → TryEmailFallback
NotificationSent → TrackDelivery
```

### 2. Analytics Projections
```go
// Real-time dashboards
ProcessRegistered → UpdateProcessCount
NotificationSent → UpdateDeliveryStats
DataJudQueryExecuted → UpdateAPIUsage
TenantQuotaExceeded → UpdateQuotaAlerts
```

### 3. Cross-Context Integration
```go
// Tenant → All contexts
TenantQuotaExceeded → [ProcessService, NotificationService]
TenantSuspended → [AuthService, ProcessService, AIService]

// Process → Multiple consumers  
ProcessMovementDetected → [NotificationService, AIService, DocumentService]
```