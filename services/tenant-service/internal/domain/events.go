package domain

import (
	"time"
	"encoding/json"
)

// DomainEvent representa um evento de domínio
type DomainEvent interface {
	EventType() string
	EventVersion() string
	AggregateID() string
	Timestamp() time.Time
	ToJSON() ([]byte, error)
}

// BaseDomainEvent implementação base para eventos de domínio
type BaseDomainEvent struct {
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	ID          string    `json:"aggregate_id"`
	OccurredAt  time.Time `json:"occurred_at"`
	TenantID    string    `json:"tenant_id"`
}

func (e BaseDomainEvent) EventType() string {
	return e.Type
}

func (e BaseDomainEvent) EventVersion() string {
	return e.Version
}

func (e BaseDomainEvent) AggregateID() string {
	return e.ID
}

func (e BaseDomainEvent) Timestamp() time.Time {
	return e.OccurredAt
}

func (e BaseDomainEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// TenantCreatedEvent evento quando um tenant é criado
type TenantCreatedEvent struct {
	BaseDomainEvent
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	PlanType  PlanType `json:"plan_type"`
	OwnerID   string   `json:"owner_id"`
}

func NewTenantCreatedEvent(tenant *Tenant) *TenantCreatedEvent {
	return &TenantCreatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "tenant.created",
			Version:    "1.0",
			ID:         tenant.ID,
			OccurredAt: time.Now(),
			TenantID:   tenant.ID,
		},
		Name:     tenant.Name,
		Email:    tenant.Email,
		PlanType: tenant.PlanType,
		OwnerID:  tenant.OwnerUserID,
	}
}

// TenantActivatedEvent evento quando um tenant é ativado
type TenantActivatedEvent struct {
	BaseDomainEvent
	ActivatedBy string `json:"activated_by"`
}

func NewTenantActivatedEvent(tenantID, activatedBy string) *TenantActivatedEvent {
	return &TenantActivatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "tenant.activated",
			Version:    "1.0",
			ID:         tenantID,
			OccurredAt: time.Now(),
			TenantID:   tenantID,
		},
		ActivatedBy: activatedBy,
	}
}

// TenantSuspendedEvent evento quando um tenant é suspenso
type TenantSuspendedEvent struct {
	BaseDomainEvent
	Reason      string `json:"reason"`
	SuspendedBy string `json:"suspended_by"`
}

func NewTenantSuspendedEvent(tenantID, reason, suspendedBy string) *TenantSuspendedEvent {
	return &TenantSuspendedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "tenant.suspended",
			Version:    "1.0",
			ID:         tenantID,
			OccurredAt: time.Now(),
			TenantID:   tenantID,
		},
		Reason:      reason,
		SuspendedBy: suspendedBy,
	}
}

// TenantCanceledEvent evento quando um tenant é cancelado
type TenantCanceledEvent struct {
	BaseDomainEvent
	Reason      string `json:"reason"`
	CanceledBy  string `json:"canceled_by"`
}

func NewTenantCanceledEvent(tenantID, reason, canceledBy string) *TenantCanceledEvent {
	return &TenantCanceledEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "tenant.canceled",
			Version:    "1.0",
			ID:         tenantID,
			OccurredAt: time.Now(),
			TenantID:   tenantID,
		},
		Reason:     reason,
		CanceledBy: canceledBy,
	}
}

// TenantUpdatedEvent evento quando um tenant é atualizado
type TenantUpdatedEvent struct {
	BaseDomainEvent
	Changes   map[string]interface{} `json:"changes"`
	UpdatedBy string                 `json:"updated_by"`
}

func NewTenantUpdatedEvent(tenantID string, changes map[string]interface{}, updatedBy string) *TenantUpdatedEvent {
	return &TenantUpdatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "tenant.updated",
			Version:    "1.0",
			ID:         tenantID,
			OccurredAt: time.Now(),
			TenantID:   tenantID,
		},
		Changes:   changes,
		UpdatedBy: updatedBy,
	}
}

// SubscriptionCreatedEvent evento quando uma assinatura é criada
type SubscriptionCreatedEvent struct {
	BaseDomainEvent
	PlanID           string            `json:"plan_id"`
	BillingInterval  BillingInterval   `json:"billing_interval"`
	TrialDays        int               `json:"trial_days"`
}

func NewSubscriptionCreatedEvent(subscription *Subscription, trialDays int) *SubscriptionCreatedEvent {
	return &SubscriptionCreatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "subscription.created",
			Version:    "1.0",
			ID:         subscription.ID,
			OccurredAt: time.Now(),
			TenantID:   subscription.TenantID,
		},
		PlanID:          subscription.PlanID,
		BillingInterval: BillingMonthly, // Default
		TrialDays:       trialDays,
	}
}

// SubscriptionActivatedEvent evento quando uma assinatura é ativada
type SubscriptionActivatedEvent struct {
	BaseDomainEvent
	PlanID string `json:"plan_id"`
}

func NewSubscriptionActivatedEvent(subscription *Subscription) *SubscriptionActivatedEvent {
	return &SubscriptionActivatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "subscription.activated",
			Version:    "1.0",
			ID:         subscription.ID,
			OccurredAt: time.Now(),
			TenantID:   subscription.TenantID,
		},
		PlanID: subscription.PlanID,
	}
}

// SubscriptionCanceledEvent evento quando uma assinatura é cancelada
type SubscriptionCanceledEvent struct {
	BaseDomainEvent
	Reason           string `json:"reason"`
	CancelAtPeriodEnd bool   `json:"cancel_at_period_end"`
}

func NewSubscriptionCanceledEvent(subscription *Subscription, reason string) *SubscriptionCanceledEvent {
	return &SubscriptionCanceledEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "subscription.canceled",
			Version:    "1.0",
			ID:         subscription.ID,
			OccurredAt: time.Now(),
			TenantID:   subscription.TenantID,
		},
		Reason:            reason,
		CancelAtPeriodEnd: subscription.CancelAtPeriodEnd,
	}
}

// SubscriptionRenewedEvent evento quando uma assinatura é renovada
type SubscriptionRenewedEvent struct {
	BaseDomainEvent
	PreviousPeriodEnd time.Time       `json:"previous_period_end"`
	NewPeriodEnd      time.Time       `json:"new_period_end"`
	BillingInterval   BillingInterval `json:"billing_interval"`
}

func NewSubscriptionRenewedEvent(subscription *Subscription, previousPeriodEnd time.Time, billingInterval BillingInterval) *SubscriptionRenewedEvent {
	return &SubscriptionRenewedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "subscription.renewed",
			Version:    "1.0",
			ID:         subscription.ID,
			OccurredAt: time.Now(),
			TenantID:   subscription.TenantID,
		},
		PreviousPeriodEnd: previousPeriodEnd,
		NewPeriodEnd:      subscription.CurrentPeriodEnd,
		BillingInterval:   billingInterval,
	}
}

// PlanChangedEvent evento quando um plano é alterado
type PlanChangedEvent struct {
	BaseDomainEvent
	FromPlanID string `json:"from_plan_id"`
	ToPlanID   string `json:"to_plan_id"`
	Reason     string `json:"reason"`
}

func NewPlanChangedEvent(subscriptionID, tenantID, fromPlanID, toPlanID, reason string) *PlanChangedEvent {
	return &PlanChangedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "subscription.plan_changed",
			Version:    "1.0",
			ID:         subscriptionID,
			OccurredAt: time.Now(),
			TenantID:   tenantID,
		},
		FromPlanID: fromPlanID,
		ToPlanID:   toPlanID,
		Reason:     reason,
	}
}

// QuotaExceededEvent evento quando uma quota é excedida
type QuotaExceededEvent struct {
	BaseDomainEvent
	QuotaType    string  `json:"quota_type"`
	CurrentUsage int     `json:"current_usage"`
	Limit        int     `json:"limit"`
	Percentage   float64 `json:"percentage"`
}

func NewQuotaExceededEvent(tenantID, quotaType string, currentUsage, limit int) *QuotaExceededEvent {
	percentage := float64(currentUsage) / float64(limit) * 100
	return &QuotaExceededEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "tenant.quota_exceeded",
			Version:    "1.0",
			ID:         tenantID,
			OccurredAt: time.Now(),
			TenantID:   tenantID,
		},
		QuotaType:    quotaType,
		CurrentUsage: currentUsage,
		Limit:        limit,
		Percentage:   percentage,
	}
}

// QuotaWarningEvent evento quando uma quota está próxima do limite
type QuotaWarningEvent struct {
	BaseDomainEvent
	QuotaType    string  `json:"quota_type"`
	CurrentUsage int     `json:"current_usage"`
	Limit        int     `json:"limit"`
	Percentage   float64 `json:"percentage"`
	Threshold    float64 `json:"threshold"`
}

func NewQuotaWarningEvent(tenantID, quotaType string, currentUsage, limit int, threshold float64) *QuotaWarningEvent {
	percentage := float64(currentUsage) / float64(limit) * 100
	return &QuotaWarningEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:       "tenant.quota_warning",
			Version:    "1.0",
			ID:         tenantID,
			OccurredAt: time.Now(),
			TenantID:   tenantID,
		},
		QuotaType:    quotaType,
		CurrentUsage: currentUsage,
		Limit:        limit,
		Percentage:   percentage,
		Threshold:    threshold,
	}
}

// EventPublisher interface para publicação de eventos
type EventPublisher interface {
	Publish(event DomainEvent) error
	PublishBatch(events []DomainEvent) error
}