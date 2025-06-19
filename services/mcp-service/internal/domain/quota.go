package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// QuotaUsage representa o uso de quota de comandos MCP
type QuotaUsage struct {
	ID              uuid.UUID `json:"id" db:"id"`
	TenantID        uuid.UUID `json:"tenant_id" db:"tenant_id"`
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	Period          string    `json:"period" db:"period"` // monthly, daily
	PeriodStart     time.Time `json:"period_start" db:"period_start"`
	PeriodEnd       time.Time `json:"period_end" db:"period_end"`
	CommandCount    int       `json:"command_count" db:"command_count"`
	QuotaLimit      int       `json:"quota_limit" db:"quota_limit"`
	Plan            string    `json:"plan" db:"plan"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// PlanQuotas define os limites de comandos por plano
var PlanQuotas = map[string]int{
	"starter":      0,    // Sem acesso ao MCP
	"professional": 200,  // 200 comandos/mês
	"business":     1000, // 1000 comandos/mês
	"enterprise":   -1,   // Ilimitado
}

// NewQuotaUsage cria um novo registro de uso de quota
func NewQuotaUsage(tenantID, userID uuid.UUID, plan string) *QuotaUsage {
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	periodEnd := periodStart.AddDate(0, 1, 0).Add(-time.Second)
	
	quotaLimit, exists := PlanQuotas[plan]
	if !exists {
		quotaLimit = 0
	}
	
	return &QuotaUsage{
		ID:           uuid.New(),
		TenantID:     tenantID,
		UserID:       userID,
		Period:       "monthly",
		PeriodStart:  periodStart,
		PeriodEnd:    periodEnd,
		CommandCount: 0,
		QuotaLimit:   quotaLimit,
		Plan:         plan,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// CanExecuteCommand verifica se pode executar um comando
func (q *QuotaUsage) CanExecuteCommand() bool {
	// Plano enterprise tem quota ilimitada
	if q.QuotaLimit == -1 {
		return true
	}
	
	// Verificar se está dentro do período
	now := time.Now()
	if now.Before(q.PeriodStart) || now.After(q.PeriodEnd) {
		// Período expirou, precisa renovar
		return false
	}
	
	// Verificar se ainda tem quota disponível
	return q.CommandCount < q.QuotaLimit
}

// IncrementUsage incrementa o uso de comandos
func (q *QuotaUsage) IncrementUsage() error {
	if !q.CanExecuteCommand() {
		return fmt.Errorf("quota de comandos excedida para o plano %s", q.Plan)
	}
	
	q.CommandCount++
	q.UpdatedAt = time.Now()
	return nil
}

// GetRemainingQuota retorna quantos comandos ainda podem ser executados
func (q *QuotaUsage) GetRemainingQuota() int {
	if q.QuotaLimit == -1 {
		return -1 // Ilimitado
	}
	
	remaining := q.QuotaLimit - q.CommandCount
	if remaining < 0 {
		return 0
	}
	
	return remaining
}

// GetUsagePercentage retorna a porcentagem de uso da quota
func (q *QuotaUsage) GetUsagePercentage() float64 {
	if q.QuotaLimit == -1 || q.QuotaLimit == 0 {
		return 0
	}
	
	return float64(q.CommandCount) / float64(q.QuotaLimit) * 100
}

// IsNearLimit verifica se está próximo do limite (80% ou mais)
func (q *QuotaUsage) IsNearLimit() bool {
	if q.QuotaLimit == -1 {
		return false
	}
	
	return q.GetUsagePercentage() >= 80
}

// ResetForNewPeriod reseta o uso para um novo período
func (q *QuotaUsage) ResetForNewPeriod() {
	now := time.Now()
	q.PeriodStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	q.PeriodEnd = q.PeriodStart.AddDate(0, 1, 0).Add(-time.Second)
	q.CommandCount = 0
	q.UpdatedAt = now
}

// IsExpired verifica se o período expirou
func (q *QuotaUsage) IsExpired() bool {
	return time.Now().After(q.PeriodEnd)
}

// QuotaStats representa estatísticas de uso de quota
type QuotaStats struct {
	TenantID         uuid.UUID `json:"tenant_id"`
	Plan             string    `json:"plan"`
	CurrentUsage     int       `json:"current_usage"`
	QuotaLimit       int       `json:"quota_limit"`
	RemainingQuota   int       `json:"remaining_quota"`
	UsagePercentage  float64   `json:"usage_percentage"`
	PeriodStart      time.Time `json:"period_start"`
	PeriodEnd        time.Time `json:"period_end"`
	DaysRemaining    int       `json:"days_remaining"`
	IsNearLimit      bool      `json:"is_near_limit"`
	IsUnlimited      bool      `json:"is_unlimited"`
}

// GetStats retorna estatísticas detalhadas de uso
func (q *QuotaUsage) GetStats() *QuotaStats {
	now := time.Now()
	daysRemaining := int(q.PeriodEnd.Sub(now).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}
	
	return &QuotaStats{
		TenantID:        q.TenantID,
		Plan:            q.Plan,
		CurrentUsage:    q.CommandCount,
		QuotaLimit:      q.QuotaLimit,
		RemainingQuota:  q.GetRemainingQuota(),
		UsagePercentage: q.GetUsagePercentage(),
		PeriodStart:     q.PeriodStart,
		PeriodEnd:       q.PeriodEnd,
		DaysRemaining:   daysRemaining,
		IsNearLimit:     q.IsNearLimit(),
		IsUnlimited:     q.QuotaLimit == -1,
	}
}