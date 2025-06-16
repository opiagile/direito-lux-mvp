# Bounded Contexts - Direito Lux

## Context Map Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        API Gateway                              │
│                    (Cross-cutting)                             │
└─────────────┬───────────────┬───────────────┬──────────────────┘
              │               │               │
    ┌─────────▼──────┐ ┌─────▼──────┐ ┌─────▼──────┐
    │  Authentication│ │   Tenant   │ │  Process   │
    │   & Identity   │ │ Management │ │ Management │
    └─────────┬──────┘ └─────┬──────┘ └─────┬──────┘
              │               │               │
              │         ┌─────▼──────┐ ┌─────▼──────┐
              │         │ Notification│ │  External  │
              │         │   System   │ │Integration │
              │         └─────┬──────┘ └─────┬──────┘
              │               │               │
         ┌────▼─────┐  ┌─────▼──────┐ ┌─────▼──────┐
         │AI & ML   │  │ Document   │ │ Analytics  │
         │Services  │  │Management  │ │& Reporting │
         └──────────┘  └────────────┘ └────────────┘
```

## 1. Authentication & Identity Context

### Domain Responsibility
Gerenciar identidade, autenticação e autorização de usuários

### Core Aggregates
```go
// User Aggregate
type User struct {
    ID        UserID        `json:"id"`
    Email     Email         `json:"email"`
    TenantID  TenantID      `json:"tenant_id"`
    Role      Role          `json:"role"`
    Profile   UserProfile   `json:"profile"`
    Status    UserStatus    `json:"status"`
    Audit     AuditInfo     `json:"audit"`
}

// Session Aggregate  
type Session struct {
    ID          SessionID     `json:"id"`
    UserID      UserID        `json:"user_id"`
    Token       JWTToken      `json:"token"`
    RefreshToken RefreshToken `json:"refresh_token"`
    ExpiresAt   time.Time     `json:"expires_at"`
    CreatedAt   time.Time     `json:"created_at"`
}

// Permission Aggregate
type Permission struct {
    ID       PermissionID   `json:"id"`
    Resource Resource       `json:"resource"`
    Action   Action         `json:"action"`
    Scope    Scope          `json:"scope"`
}
```

### Domain Events
```go
type UserRegistered struct {
    UserID   UserID    `json:"user_id"`
    Email    Email     `json:"email"`
    TenantID TenantID  `json:"tenant_id"`
    Role     Role      `json:"role"`
    OccurredAt time.Time `json:"occurred_at"`
}

type UserLoggedIn struct {
    UserID     UserID     `json:"user_id"`
    SessionID  SessionID  `json:"session_id"`
    IPAddress  string     `json:"ip_address"`
    UserAgent  string     `json:"user_agent"`
    OccurredAt time.Time  `json:"occurred_at"`
}

type UserPermissionGranted struct {
    UserID       UserID       `json:"user_id"`
    PermissionID PermissionID `json:"permission_id"`
    GrantedBy    UserID       `json:"granted_by"`
    OccurredAt   time.Time    `json:"occurred_at"`
}
```

### External Integration
- **Keycloak**: Identity provider
- **LDAP/AD**: Enterprise integration
- **OAuth2/OIDC**: Third-party authentication

### APIs Exposed
```go
// Authentication API
POST /auth/login
POST /auth/refresh  
POST /auth/logout
GET  /auth/validate

// User Management API
GET    /users
POST   /users
GET    /users/{id}
PUT    /users/{id}
DELETE /users/{id}

// Permissions API
GET  /users/{id}/permissions
POST /users/{id}/permissions
DELETE /users/{id}/permissions/{permission_id}
```

## 2. Tenant Management Context

### Domain Responsibility
Gerenciar escritórios (tenants), assinaturas, quotas e cobrança

### Core Aggregates
```go
// Tenant Aggregate
type Tenant struct {
    ID           TenantID       `json:"id"`
    Name         string         `json:"name"`
    Document     CNPJ           `json:"document"`
    Subscription Subscription   `json:"subscription"`
    Quotas       TenantQuotas   `json:"quotas"`
    BillingInfo  BillingInfo    `json:"billing_info"`
    Status       TenantStatus   `json:"status"`
    Audit        AuditInfo      `json:"audit"`
}

// Subscription Aggregate
type Subscription struct {
    ID        SubscriptionID   `json:"id"`
    TenantID  TenantID         `json:"tenant_id"`
    Plan      SubscriptionPlan `json:"plan"`
    Status    SubscriptionStatus `json:"status"`
    StartsAt  time.Time        `json:"starts_at"`
    ExpiresAt time.Time        `json:"expires_at"`
    Features  []Feature        `json:"features"`
}

// Quota Aggregate
type TenantQuota struct {
    TenantID       TenantID `json:"tenant_id"`
    ProcessLimit   int      `json:"process_limit"`
    UsersLimit     int      `json:"users_limit"`
    DataJudQuota   int      `json:"datajud_quota"`
    StorageQuota   int64    `json:"storage_quota"`
    NotificationQuota int   `json:"notification_quota"`
    CurrentUsage   QuotaUsage `json:"current_usage"`
}
```

### Domain Events
```go
type TenantCreated struct {
    TenantID    TenantID  `json:"tenant_id"`
    Name        string    `json:"name"`
    Plan        SubscriptionPlan `json:"plan"`
    CreatedBy   UserID    `json:"created_by"`
    OccurredAt  time.Time `json:"occurred_at"`
}

type TenantQuotaExceeded struct {
    TenantID    TenantID  `json:"tenant_id"`
    QuotaType   string    `json:"quota_type"`
    Limit       int       `json:"limit"`
    Current     int       `json:"current"`
    OccurredAt  time.Time `json:"occurred_at"`
}

type SubscriptionUpgraded struct {
    TenantID     TenantID         `json:"tenant_id"`
    FromPlan     SubscriptionPlan `json:"from_plan"`
    ToPlan       SubscriptionPlan `json:"to_plan"`
    EffectiveAt  time.Time        `json:"effective_at"`
    OccurredAt   time.Time        `json:"occurred_at"`
}
```

### External Integration
- **Payment Gateway**: Stripe, PagSeguro
- **Billing System**: Chargebee, interno
- **Analytics**: Mixpanel, Amplitude

### APIs Exposed
```go
// Tenant Management API
GET    /tenants
POST   /tenants
GET    /tenants/{id}
PUT    /tenants/{id}
DELETE /tenants/{id}

// Subscription API
GET  /tenants/{id}/subscription
POST /tenants/{id}/subscription/upgrade
POST /tenants/{id}/subscription/downgrade
POST /tenants/{id}/subscription/cancel

// Quota API
GET /tenants/{id}/quotas
GET /tenants/{id}/quotas/usage
```

## 3. Process Management Context

### Domain Responsibility
Gerenciar processos jurídicos, monitoramento e atualizações

### Core Aggregates
```go
// Process Aggregate (Event Sourced)
type Process struct {
    ID              ProcessID      `json:"id"`
    Number          ProcessNumber  `json:"number"`
    TenantID        TenantID       `json:"tenant_id"`
    ClientID        ClientID       `json:"client_id"`
    Court           Court          `json:"court"`
    Classification  Classification `json:"classification"`
    Status          ProcessStatus  `json:"status"`
    MonitoringState MonitoringState `json:"monitoring_state"`
    Movements       []Movement     `json:"movements"`
    Deadlines       []Deadline     `json:"deadlines"`
    Metadata        ProcessMetadata `json:"metadata"`
    Version         int            `json:"version"`
}

// Movement Value Object
type Movement struct {
    ID          MovementID `json:"id"`
    Date        time.Time  `json:"date"`
    Description string     `json:"description"`
    Type        MovementType `json:"type"`
    Document    *Document  `json:"document,omitempty"`
    Source      string     `json:"source"`
}

// Deadline Value Object
type Deadline struct {
    ID          DeadlineID   `json:"id"`
    Date        time.Time    `json:"date"`
    Description string       `json:"description"`
    Type        DeadlineType `json:"type"`
    Status      DeadlineStatus `json:"status"`
    Priority    Priority     `json:"priority"`
}
```

### Domain Events
```go
type ProcessRegistered struct {
    ProcessID   ProcessID     `json:"process_id"`
    Number      ProcessNumber `json:"number"`
    TenantID    TenantID      `json:"tenant_id"`
    ClientID    ClientID      `json:"client_id"`
    Court       Court         `json:"court"`
    RegisteredBy UserID       `json:"registered_by"`
    OccurredAt  time.Time     `json:"occurred_at"`
}

type ProcessMovementDetected struct {
    ProcessID   ProcessID   `json:"process_id"`
    MovementID  MovementID  `json:"movement_id"`
    Date        time.Time   `json:"date"`
    Description string      `json:"description"`
    Type        MovementType `json:"type"`
    Source      string      `json:"source"`
    OccurredAt  time.Time   `json:"occurred_at"`
}

type ProcessDeadlineApproaching struct {
    ProcessID    ProcessID    `json:"process_id"`
    DeadlineID   DeadlineID   `json:"deadline_id"`
    DeadlineDate time.Time    `json:"deadline_date"`
    DaysRemaining int         `json:"days_remaining"`
    Priority     Priority     `json:"priority"`
    OccurredAt   time.Time    `json:"occurred_at"`
}
```

### External Integration
- **DataJud API**: Dados processuais oficiais
- **Tribunais**: APIs específicas (quando disponível)
- **Document Storage**: S3, GCS

### APIs Exposed
```go
// Process Management API
GET    /processes
POST   /processes
GET    /processes/{id}
PUT    /processes/{id}
DELETE /processes/{id}

// Monitoring API
POST /processes/{id}/monitoring/start
POST /processes/{id}/monitoring/stop
GET  /processes/{id}/monitoring/status

// Search API
GET /processes/search
GET /processes/by-client/{client_id}
GET /processes/by-court/{court}
```

## 4. External Integration Context

### Domain Responsibility
Gerenciar integrações com APIs externas e cache

### Core Aggregates
```go
// DataJudQuery Aggregate
type DataJudQuery struct {
    ID            QueryID       `json:"id"`
    TenantID      TenantID      `json:"tenant_id"`
    ProcessNumber ProcessNumber `json:"process_number"`
    QueryType     QueryType     `json:"query_type"`
    Status        QueryStatus   `json:"status"`
    RequestData   interface{}   `json:"request_data"`
    ResponseData  interface{}   `json:"response_data"`
    RateLimit     RateLimit     `json:"rate_limit"`
    Cache         CacheInfo     `json:"cache"`
    Audit         AuditInfo     `json:"audit"`
}

// Cache Aggregate
type CacheEntry struct {
    Key        string        `json:"key"`
    Value      interface{}   `json:"value"`
    TTL        time.Duration `json:"ttl"`
    CreatedAt  time.Time     `json:"created_at"`
    AccessedAt time.Time     `json:"accessed_at"`
    TenantID   TenantID      `json:"tenant_id"`
}

// RateLimit Aggregate
type RateLimit struct {
    TenantID      TenantID  `json:"tenant_id"`
    Service       string    `json:"service"`
    Limit         int       `json:"limit"`
    WindowSize    time.Duration `json:"window_size"`
    CurrentCount  int       `json:"current_count"`
    ResetsAt      time.Time `json:"resets_at"`
}
```

### Domain Events
```go
type DataJudQueryExecuted struct {
    QueryID       QueryID       `json:"query_id"`
    TenantID      TenantID      `json:"tenant_id"`
    ProcessNumber ProcessNumber `json:"process_number"`
    Success       bool          `json:"success"`
    ResponseTime  time.Duration `json:"response_time"`
    CacheHit      bool          `json:"cache_hit"`
    OccurredAt    time.Time     `json:"occurred_at"`
}

type RateLimitExceeded struct {
    TenantID   TenantID  `json:"tenant_id"`
    Service    string    `json:"service"`
    Limit      int       `json:"limit"`
    Attempts   int       `json:"attempts"`
    OccurredAt time.Time `json:"occurred_at"`
}

type CacheEvicted struct {
    Key        string    `json:"key"`
    Reason     string    `json:"reason"`
    TenantID   TenantID  `json:"tenant_id"`
    OccurredAt time.Time `json:"occurred_at"`
}
```

### External Integration
- **DataJud CNJ**: API oficial
- **Tribunais**: TJ-SP, TRF-3, etc.
- **Cache**: Redis Cluster
- **Circuit Breaker**: Hystrix pattern

### APIs Exposed
```go
// DataJud Integration API
GET  /datajud/processes/{number}
POST /datajud/processes/batch
GET  /datajud/tribunals
GET  /datajud/cache/stats

// Rate Limiting API
GET /rate-limits/{tenant_id}
GET /rate-limits/{tenant_id}/usage
```

## 5. Notification System Context

### Domain Responsibility
Orquestrar e gerenciar envio de notificações multicanal

### Core Aggregates
```go
// Notification Aggregate
type Notification struct {
    ID          NotificationID      `json:"id"`
    TenantID    TenantID            `json:"tenant_id"`
    ProcessID   *ProcessID          `json:"process_id,omitempty"`
    RecipientID UserID              `json:"recipient_id"`
    Type        NotificationType    `json:"type"`
    Priority    Priority            `json:"priority"`
    Channel     NotificationChannel `json:"channel"`
    Template    TemplateID          `json:"template"`
    Content     NotificationContent `json:"content"`
    Status      NotificationStatus  `json:"status"`
    Attempts    int                 `json:"attempts"`
    Metadata    NotificationMetadata `json:"metadata"`
    Audit       AuditInfo           `json:"audit"`
}

// NotificationTemplate Aggregate
type NotificationTemplate struct {
    ID        TemplateID           `json:"id"`
    TenantID  TenantID             `json:"tenant_id"`
    Name      string               `json:"name"`
    Type      NotificationType     `json:"type"`
    Channel   NotificationChannel  `json:"channel"`
    Subject   string               `json:"subject"`
    Body      string               `json:"body"`
    Variables []TemplateVariable   `json:"variables"`
    IsActive  bool                 `json:"is_active"`
}

// DeliveryStrategy Aggregate
type DeliveryStrategy struct {
    TenantID       TenantID              `json:"tenant_id"`
    Type           NotificationType      `json:"type"`
    PrimaryChannel NotificationChannel   `json:"primary_channel"`
    FallbackChannels []NotificationChannel `json:"fallback_channels"`
    RetryPolicy    RetryPolicy           `json:"retry_policy"`
    Preferences    DeliveryPreferences   `json:"preferences"`
}
```

### Domain Events
```go
type NotificationRequested struct {
    NotificationID NotificationID   `json:"notification_id"`
    TenantID       TenantID         `json:"tenant_id"`
    RecipientID    UserID           `json:"recipient_id"`
    Type           NotificationType `json:"type"`
    Priority       Priority         `json:"priority"`
    Channel        NotificationChannel `json:"channel"`
    OccurredAt     time.Time        `json:"occurred_at"`
}

type NotificationSent struct {
    NotificationID NotificationID      `json:"notification_id"`
    Channel        NotificationChannel `json:"channel"`
    ExternalID     string              `json:"external_id"`
    Attempt        int                 `json:"attempt"`
    OccurredAt     time.Time           `json:"occurred_at"`
}

type NotificationDelivered struct {
    NotificationID NotificationID      `json:"notification_id"`
    Channel        NotificationChannel `json:"channel"`
    DeliveredAt    time.Time           `json:"delivered_at"`
    OccurredAt     time.Time           `json:"occurred_at"`
}

type NotificationFailed struct {
    NotificationID NotificationID      `json:"notification_id"`
    Channel        NotificationChannel `json:"channel"`
    Error          string              `json:"error"`
    Attempt        int                 `json:"attempt"`
    WillRetry      bool                `json:"will_retry"`
    OccurredAt     time.Time           `json:"occurred_at"`
}
```

### External Integration
- **WhatsApp Business API**: Canal principal
- **SMTP Provider**: SendGrid, SES
- **Telegram Bot API**: Canal alternativo
- **Push Notifications**: FCM

### APIs Exposed
```go
// Notification API
POST /notifications/send
GET  /notifications/{id}
GET  /notifications/{id}/status
POST /notifications/{id}/retry

// Template API
GET    /notifications/templates
POST   /notifications/templates
GET    /notifications/templates/{id}
PUT    /notifications/templates/{id}
DELETE /notifications/templates/{id}

// Webhook API (for delivery status)
POST /notifications/webhooks/whatsapp
POST /notifications/webhooks/email
```

## 6. AI & Analytics Context

### Domain Responsibility
Inteligência artificial, análises e jurimetria

### Core Aggregates
```go
// AISummary Aggregate
type AISummary struct {
    ID         SummaryID    `json:"id"`
    TenantID   TenantID     `json:"tenant_id"`
    ProcessID  ProcessID    `json:"process_id"`
    Type       SummaryType  `json:"type"`
    Content    string       `json:"content"`
    Confidence float64      `json:"confidence"`
    Model      string       `json:"model"`
    Version    string       `json:"version"`
    Language   string       `json:"language"`
    Metadata   SummaryMetadata `json:"metadata"`
    CreatedAt  time.Time    `json:"created_at"`
}

// AnalysisResult Aggregate
type AnalysisResult struct {
    ID         AnalysisID    `json:"id"`
    TenantID   TenantID      `json:"tenant_id"`
    Type       AnalysisType  `json:"type"`
    Subject    AnalysisSubject `json:"subject"`
    Result     AnalysisData  `json:"result"`
    Confidence float64       `json:"confidence"`
    Model      string        `json:"model"`
    CreatedAt  time.Time     `json:"created_at"`
}

// MLModel Aggregate
type MLModel struct {
    ID          ModelID     `json:"id"`
    TenantID    *TenantID   `json:"tenant_id,omitempty"` // null for global models
    Name        string      `json:"name"`
    Type        ModelType   `json:"type"`
    Version     string      `json:"version"`
    Status      ModelStatus `json:"status"`
    Accuracy    float64     `json:"accuracy"`
    TrainedAt   time.Time   `json:"trained_at"`
    DeployedAt  *time.Time  `json:"deployed_at,omitempty"`
}
```

### Domain Events
```go
type AISummarizationRequested struct {
    SummaryID  SummaryID    `json:"summary_id"`
    TenantID   TenantID     `json:"tenant_id"`
    ProcessID  ProcessID    `json:"process_id"`
    Type       SummaryType  `json:"type"`
    RequestedBy UserID      `json:"requested_by"`
    OccurredAt time.Time    `json:"occurred_at"`
}

type AISummaryGenerated struct {
    SummaryID  SummaryID    `json:"summary_id"`
    ProcessID  ProcessID    `json:"process_id"`
    Content    string       `json:"content"`
    Confidence float64      `json:"confidence"`
    Model      string       `json:"model"`
    OccurredAt time.Time    `json:"occurred_at"`
}

type AnalysisCompleted struct {
    AnalysisID AnalysisID   `json:"analysis_id"`
    TenantID   TenantID     `json:"tenant_id"`
    Type       AnalysisType `json:"type"`
    Success    bool         `json:"success"`
    Duration   time.Duration `json:"duration"`
    OccurredAt time.Time    `json:"occurred_at"`
}
```

### External Integration
- **OpenAI API**: GPT models
- **Hugging Face**: Open-source models
- **Google AI**: Vertex AI
- **Local Models**: Self-hosted

### APIs Exposed
```go
// AI Services API
POST /ai/summarize
POST /ai/explain-terms
POST /ai/classify-document
POST /ai/predict-outcome

// Analytics API
GET  /analytics/dashboard/{tenant_id}
GET  /analytics/reports/{tenant_id}
POST /analytics/custom-query

// Model Management API (Enterprise)
GET    /models
POST   /models/train
GET    /models/{id}/performance
POST   /models/{id}/deploy
```

## 7. Document Management Context

### Domain Responsibility
Geração, armazenamento e gestão de documentos jurídicos

### Core Aggregates
```go
// Document Aggregate
type Document struct {
    ID          DocumentID    `json:"id"`
    TenantID    TenantID      `json:"tenant_id"`
    ProcessID   *ProcessID    `json:"process_id,omitempty"`
    Type        DocumentType  `json:"type"`
    Name        string        `json:"name"`
    MimeType    string        `json:"mime_type"`
    Size        int64         `json:"size"`
    Hash        string        `json:"hash"`
    StoragePath string        `json:"storage_path"`
    Status      DocumentStatus `json:"status"`
    Metadata    DocumentMetadata `json:"metadata"`
    Versions    []DocumentVersion `json:"versions"`
    Audit       AuditInfo     `json:"audit"`
}

// DocumentTemplate Aggregate
type DocumentTemplate struct {
    ID          TemplateID      `json:"id"`
    TenantID    TenantID        `json:"tenant_id"`
    Name        string          `json:"name"`
    Type        DocumentType    `json:"type"`
    Category    TemplateCategory `json:"category"`
    Content     string          `json:"content"`
    Variables   []TemplateVariable `json:"variables"`
    IsActive    bool            `json:"is_active"`
    Version     int             `json:"version"`
}

// Signature Aggregate
type Signature struct {
    ID         SignatureID   `json:"id"`
    DocumentID DocumentID    `json:"document_id"`
    SignerID   UserID        `json:"signer_id"`
    Type       SignatureType `json:"type"`
    Status     SignatureStatus `json:"status"`
    SignedAt   *time.Time    `json:"signed_at,omitempty"`
    Certificate string       `json:"certificate"`
    Metadata   SignatureMetadata `json:"metadata"`
}
```

### Domain Events
```go
type DocumentGenerationRequested struct {
    DocumentID DocumentID   `json:"document_id"`
    TenantID   TenantID     `json:"tenant_id"`
    TemplateID TemplateID   `json:"template_id"`
    Variables  map[string]interface{} `json:"variables"`
    RequestedBy UserID      `json:"requested_by"`
    OccurredAt time.Time    `json:"occurred_at"`
}

type DocumentGenerated struct {
    DocumentID  DocumentID   `json:"document_id"`
    Name        string       `json:"name"`
    Type        DocumentType `json:"type"`
    Size        int64        `json:"size"`
    StoragePath string       `json:"storage_path"`
    OccurredAt  time.Time    `json:"occurred_at"`
}

type DocumentSigned struct {
    DocumentID  DocumentID    `json:"document_id"`
    SignatureID SignatureID   `json:"signature_id"`
    SignerID    UserID        `json:"signer_id"`
    Type        SignatureType `json:"type"`
    SignedAt    time.Time     `json:"signed_at"`
    OccurredAt  time.Time     `json:"occurred_at"`
}
```

### External Integration
- **Cloud Storage**: GCS, S3
- **Digital Signature**: DocuSign, Adobe Sign
- **PDF Generation**: wkhtmltopdf, Puppeteer
- **OCR Service**: Google Vision, Tesseract

### APIs Exposed
```go
// Document API
GET    /documents
POST   /documents/generate
GET    /documents/{id}
GET    /documents/{id}/download
DELETE /documents/{id}

// Template API
GET    /documents/templates
POST   /documents/templates
GET    /documents/templates/{id}
PUT    /documents/templates/{id}

// Signature API
POST /documents/{id}/signature/request
GET  /documents/{id}/signatures
POST /documents/{id}/signatures/{signature_id}/sign
```

## Context Integration Patterns

### 1. Upstream/Downstream Relationships
```
Authentication → Tenant Management (ACL)
Tenant Management → Process Management (Shared Kernel)
Process Management → Notification System (Published Language)
External Integration ← Process Management (Customer/Supplier)
```

### 2. Event Publishing Strategy
```go
// Events published to message bus
type EventBus interface {
    Publish(ctx context.Context, event DomainEvent) error
    Subscribe(eventType string, handler EventHandler) error
}

// Cross-context events
ProcessRegistered    → NotificationSystem (send welcome)
ProcessUpdated       → NotificationSystem (alert users)
QuotaExceeded        → NotificationSystem (alert admin)
NotificationFailed   → Analytics (track failures)
```

### 3. Shared Kernel Components
```go
// Shared Value Objects
type TenantID string
type UserID string  
type ProcessID string
type Money struct {
    Amount   int64
    Currency string
}

// Shared Events
type AuditInfo struct {
    CreatedAt time.Time
    CreatedBy UserID
    UpdatedAt time.Time
    UpdatedBy UserID
}
```

### 4. Anti-Corruption Layers
```go
// DataJud API responses → Domain models
type DataJudAdapter interface {
    FetchProcess(number ProcessNumber) (*Process, error)
    TransformMovement(raw DataJudMovement) Movement
}

// WhatsApp API → Notification domain
type WhatsAppAdapter interface {
    SendMessage(notification Notification) error
    HandleWebhook(payload []byte) (*DeliveryStatus, error)
}
```