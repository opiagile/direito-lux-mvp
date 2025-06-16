package domain

import (
	"time"
	"errors"
)

// QuotaUsage representa o uso atual de quotas de um tenant
type QuotaUsage struct {
	ID                    string    `json:"id" db:"id"`
	TenantID              string    `json:"tenant_id" db:"tenant_id"`
	ProcessesCount        int       `json:"processes_count" db:"processes_count"`
	UsersCount            int       `json:"users_count" db:"users_count"`
	ClientsCount          int       `json:"clients_count" db:"clients_count"`
	DataJudQueriesDaily   int       `json:"datajud_queries_daily" db:"datajud_queries_daily"`
	DataJudQueriesMonth   int       `json:"datajud_queries_month" db:"datajud_queries_month"`
	AIQueriesMonthly      int       `json:"ai_queries_monthly" db:"ai_queries_monthly"`
	StorageUsedGB         float64   `json:"storage_used_gb" db:"storage_used_gb"`
	WebhooksCount         int       `json:"webhooks_count" db:"webhooks_count"`
	APICallsDaily         int       `json:"api_calls_daily" db:"api_calls_daily"`
	APICallsMonthly       int       `json:"api_calls_monthly" db:"api_calls_monthly"`
	LastUpdated           time.Time `json:"last_updated" db:"last_updated"`
	LastResetDaily        time.Time `json:"last_reset_daily" db:"last_reset_daily"`
	LastResetMonthly      time.Time `json:"last_reset_monthly" db:"last_reset_monthly"`
}

// QuotaLimit representa os limites de quota para um tenant
type QuotaLimit struct {
	TenantID              string    `json:"tenant_id"`
	MaxProcesses          int       `json:"max_processes"`
	MaxUsers              int       `json:"max_users"`
	MaxClients            int       `json:"max_clients"`
	DataJudQueriesDaily   int       `json:"datajud_queries_daily"`
	AIQueriesMonthly      int       `json:"ai_queries_monthly"`
	StorageGB             int       `json:"storage_gb"`
	MaxWebhooks           int       `json:"max_webhooks"`
	MaxAPICallsDaily      int       `json:"max_api_calls_daily"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// QuotaCheck resultado de verificação de quota
type QuotaCheck struct {
	QuotaType    string  `json:"quota_type"`
	Current      int     `json:"current"`
	Limit        int     `json:"limit"`
	Available    int     `json:"available"`
	Percentage   float64 `json:"percentage"`
	IsExceeded   bool    `json:"is_exceeded"`
	IsWarning    bool    `json:"is_warning"`
}

// QuotaRepository interface para persistência de quotas
type QuotaRepository interface {
	GetUsage(tenantID string) (*QuotaUsage, error)
	UpdateUsage(usage *QuotaUsage) error
	IncrementCounter(tenantID, counterType string, amount int) error
	ResetDailyCounters(tenantID string) error
	ResetMonthlyCounters(tenantID string) error
	GetLimits(tenantID string) (*QuotaLimit, error)
	UpdateLimits(limits *QuotaLimit) error
}

// Erros de domínio para quotas
var (
	ErrQuotaExceeded        = errors.New("quota excedida")
	ErrQuotaNotFound        = errors.New("quota não encontrada")
	ErrInvalidQuotaType     = errors.New("tipo de quota inválido")
	ErrNegativeUsage        = errors.New("uso não pode ser negativo")
	ErrQuotaLimitExceeded   = errors.New("limite de quota excedido")
)

// CheckProcessesQuota verifica quota de processos
func (qu *QuotaUsage) CheckProcessesQuota(limit int) *QuotaCheck {
	return qu.checkQuota("processes", qu.ProcessesCount, limit)
}

// CheckUsersQuota verifica quota de usuários
func (qu *QuotaUsage) CheckUsersQuota(limit int) *QuotaCheck {
	return qu.checkQuota("users", qu.UsersCount, limit)
}

// CheckClientsQuota verifica quota de clientes
func (qu *QuotaUsage) CheckClientsQuota(limit int) *QuotaCheck {
	return qu.checkQuota("clients", qu.ClientsCount, limit)
}

// CheckDataJudDailyQuota verifica quota diária do DataJud
func (qu *QuotaUsage) CheckDataJudDailyQuota(limit int) *QuotaCheck {
	return qu.checkQuota("datajud_daily", qu.DataJudQueriesDaily, limit)
}

// CheckAIMonthlyQuota verifica quota mensal de IA
func (qu *QuotaUsage) CheckAIMonthlyQuota(limit int) *QuotaCheck {
	return qu.checkQuota("ai_monthly", qu.AIQueriesMonthly, limit)
}

// CheckStorageQuota verifica quota de armazenamento
func (qu *QuotaUsage) CheckStorageQuota(limitGB int) *QuotaCheck {
	current := int(qu.StorageUsedGB * 100) // Converter para centésimos de GB
	limit := limitGB * 100
	return qu.checkQuota("storage", current, limit)
}

// CheckWebhooksQuota verifica quota de webhooks
func (qu *QuotaUsage) CheckWebhooksQuota(limit int) *QuotaCheck {
	return qu.checkQuota("webhooks", qu.WebhooksCount, limit)
}

// CheckAPIDailyQuota verifica quota diária de API
func (qu *QuotaUsage) CheckAPIDailyQuota(limit int) *QuotaCheck {
	return qu.checkQuota("api_daily", qu.APICallsDaily, limit)
}

// checkQuota método auxiliar para verificação de quotas
func (qu *QuotaUsage) checkQuota(quotaType string, current, limit int) *QuotaCheck {
	if limit <= 0 { // Quota ilimitada
		return &QuotaCheck{
			QuotaType:    quotaType,
			Current:      current,
			Limit:        -1,
			Available:    -1,
			Percentage:   0,
			IsExceeded:   false,
			IsWarning:    false,
		}
	}

	available := limit - current
	if available < 0 {
		available = 0
	}

	percentage := float64(current) / float64(limit) * 100
	isExceeded := current >= limit
	isWarning := percentage >= 80.0 && !isExceeded

	return &QuotaCheck{
		QuotaType:    quotaType,
		Current:      current,
		Limit:        limit,
		Available:    available,
		Percentage:   percentage,
		IsExceeded:   isExceeded,
		IsWarning:    isWarning,
	}
}

// CanIncrementProcesses verifica se pode incrementar processos
func (qu *QuotaUsage) CanIncrementProcesses(limit, increment int) bool {
	if limit <= 0 { // Ilimitado
		return true
	}
	return (qu.ProcessesCount + increment) <= limit
}

// CanIncrementUsers verifica se pode incrementar usuários
func (qu *QuotaUsage) CanIncrementUsers(limit, increment int) bool {
	if limit <= 0 { // Ilimitado
		return true
	}
	return (qu.UsersCount + increment) <= limit
}

// CanIncrementClients verifica se pode incrementar clientes
func (qu *QuotaUsage) CanIncrementClients(limit, increment int) bool {
	if limit <= 0 { // Ilimitado
		return true
	}
	return (qu.ClientsCount + increment) <= limit
}

// CanIncrementDataJudDaily verifica se pode fazer consulta DataJud
func (qu *QuotaUsage) CanIncrementDataJudDaily(limit, increment int) bool {
	if limit <= 0 { // Ilimitado
		return true
	}
	return (qu.DataJudQueriesDaily + increment) <= limit
}

// CanIncrementAIMonthly verifica se pode fazer consulta IA
func (qu *QuotaUsage) CanIncrementAIMonthly(limit, increment int) bool {
	if limit <= 0 { // Ilimitado
		return true
	}
	return (qu.AIQueriesMonthly + increment) <= limit
}

// CanIncrementStorage verifica se pode usar mais armazenamento
func (qu *QuotaUsage) CanIncrementStorage(limitGB int, incrementGB float64) bool {
	if limitGB <= 0 { // Ilimitado
		return true
	}
	return (qu.StorageUsedGB + incrementGB) <= float64(limitGB)
}

// CanIncrementWebhooks verifica se pode adicionar webhook
func (qu *QuotaUsage) CanIncrementWebhooks(limit, increment int) bool {
	if limit <= 0 { // Ilimitado
		return true
	}
	return (qu.WebhooksCount + increment) <= limit
}

// CanIncrementAPIDaily verifica se pode fazer chamada API
func (qu *QuotaUsage) CanIncrementAPIDaily(limit, increment int) bool {
	if limit <= 0 { // Ilimitado
		return true
	}
	return (qu.APICallsDaily + increment) <= limit
}

// IncrementProcesses incrementa contador de processos
func (qu *QuotaUsage) IncrementProcesses(amount int) error {
	if amount < 0 {
		return ErrNegativeUsage
	}
	qu.ProcessesCount += amount
	qu.LastUpdated = time.Now()
	return nil
}

// DecrementProcesses decrementa contador de processos
func (qu *QuotaUsage) DecrementProcesses(amount int) error {
	if amount < 0 {
		return ErrNegativeUsage
	}
	qu.ProcessesCount -= amount
	if qu.ProcessesCount < 0 {
		qu.ProcessesCount = 0
	}
	qu.LastUpdated = time.Now()
	return nil
}

// IncrementUsers incrementa contador de usuários
func (qu *QuotaUsage) IncrementUsers(amount int) error {
	if amount < 0 {
		return ErrNegativeUsage
	}
	qu.UsersCount += amount
	qu.LastUpdated = time.Now()
	return nil
}

// DecrementUsers decrementa contador de usuários
func (qu *QuotaUsage) DecrementUsers(amount int) error {
	if amount < 0 {
		return ErrNegativeUsage
	}
	qu.UsersCount -= amount
	if qu.UsersCount < 0 {
		qu.UsersCount = 0
	}
	qu.LastUpdated = time.Now()
	return nil
}

// IncrementClients incrementa contador de clientes
func (qu *QuotaUsage) IncrementClients(amount int) error {
	if amount < 0 {
		return ErrNegativeUsage
	}
	qu.ClientsCount += amount
	qu.LastUpdated = time.Now()
	return nil
}

// DecrementClients decrementa contador de clientes
func (qu *QuotaUsage) DecrementClients(amount int) error {
	if amount < 0 {
		return ErrNegativeUsage
	}
	qu.ClientsCount -= amount
	if qu.ClientsCount < 0 {
		qu.ClientsCount = 0
	}
	qu.LastUpdated = time.Now()
	return nil
}

// IncrementDataJudDaily incrementa contador DataJud diário
func (qu *QuotaUsage) IncrementDataJudDaily(amount int) error {
	if amount < 0 {
		return ErrNegativeUsage
	}
	qu.DataJudQueriesDaily += amount
	qu.DataJudQueriesMonth += amount
	qu.LastUpdated = time.Now()
	return nil
}

// IncrementAIMonthly incrementa contador IA mensal
func (qu *QuotaUsage) IncrementAIMonthly(amount int) error {
	if amount < 0 {
		return ErrNegativeUsage
	}
	qu.AIQueriesMonthly += amount
	qu.LastUpdated = time.Now()
	return nil
}

// UpdateStorageUsage atualiza uso de armazenamento
func (qu *QuotaUsage) UpdateStorageUsage(usageGB float64) error {
	if usageGB < 0 {
		return ErrNegativeUsage
	}
	qu.StorageUsedGB = usageGB
	qu.LastUpdated = time.Now()
	return nil
}

// IncrementAPIDaily incrementa contador API diário
func (qu *QuotaUsage) IncrementAPIDaily(amount int) error {
	if amount < 0 {
		return ErrNegativeUsage
	}
	qu.APICallsDaily += amount
	qu.APICallsMonthly += amount
	qu.LastUpdated = time.Now()
	return nil
}

// ShouldResetDaily verifica se deve resetar contadores diários
func (qu *QuotaUsage) ShouldResetDaily() bool {
	now := time.Now()
	lastReset := qu.LastResetDaily
	
	// Se nunca foi resetado ou se passou de 00:00 do dia atual
	return lastReset.IsZero() || 
		   (now.Year() > lastReset.Year() || 
		    now.YearDay() > lastReset.YearDay())
}

// ShouldResetMonthly verifica se deve resetar contadores mensais
func (qu *QuotaUsage) ShouldResetMonthly() bool {
	now := time.Now()
	lastReset := qu.LastResetMonthly
	
	// Se nunca foi resetado ou se passou para um novo mês
	return lastReset.IsZero() || 
		   (now.Year() > lastReset.Year() || 
		    now.Month() > lastReset.Month())
}

// ResetDailyCounters reseta contadores diários
func (qu *QuotaUsage) ResetDailyCounters() {
	qu.DataJudQueriesDaily = 0
	qu.APICallsDaily = 0
	qu.LastResetDaily = time.Now()
	qu.LastUpdated = time.Now()
}

// ResetMonthlyCounters reseta contadores mensais
func (qu *QuotaUsage) ResetMonthlyCounters() {
	qu.DataJudQueriesMonth = 0
	qu.AIQueriesMonthly = 0
	qu.APICallsMonthly = 0
	qu.LastResetMonthly = time.Now()
	qu.LastUpdated = time.Now()
}

// GetQuotaLimitFromPlan cria QuotaLimit a partir das quotas do plano
func GetQuotaLimitFromPlan(tenantID string, quotas PlanQuotas) *QuotaLimit {
	return &QuotaLimit{
		TenantID:              tenantID,
		MaxProcesses:          quotas.MaxProcesses,
		MaxUsers:              quotas.MaxUsers,
		MaxClients:            quotas.MaxClients,
		DataJudQueriesDaily:   quotas.DataJudQueriesDaily,
		AIQueriesMonthly:      quotas.AIQueriesMonthly,
		StorageGB:             quotas.StorageGB,
		MaxWebhooks:           quotas.MaxWebhooks,
		MaxAPICallsDaily:      quotas.MaxAPICallsDaily,
		UpdatedAt:             time.Now(),
	}
}

// NewQuotaUsage cria nova instância de QuotaUsage
func NewQuotaUsage(tenantID string) *QuotaUsage {
	now := time.Now()
	return &QuotaUsage{
		TenantID:              tenantID,
		ProcessesCount:        0,
		UsersCount:            0,
		ClientsCount:          0,
		DataJudQueriesDaily:   0,
		DataJudQueriesMonth:   0,
		AIQueriesMonthly:      0,
		StorageUsedGB:         0,
		WebhooksCount:         0,
		APICallsDaily:         0,
		APICallsMonthly:       0,
		LastUpdated:           now,
		LastResetDaily:        now,
		LastResetMonthly:      now,
	}
}