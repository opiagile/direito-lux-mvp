package domain

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	// Context keys
	tenantIDKey contextKey = "tenant_id"
	userIDKey   contextKey = "user_id"
	traceIDKey  contextKey = "trace_id"
	planKey     contextKey = "plan"
)

// AuthContext representa o contexto de autenticação
type AuthContext struct {
	TenantID uuid.UUID
	UserID   uuid.UUID
	Plan     string
	TraceID  string
}

// WithAuthContext adiciona contexto de autenticação
func WithAuthContext(ctx context.Context, authCtx AuthContext) context.Context {
	ctx = context.WithValue(ctx, tenantIDKey, authCtx.TenantID)
	ctx = context.WithValue(ctx, userIDKey, authCtx.UserID)
	ctx = context.WithValue(ctx, planKey, authCtx.Plan)
	ctx = context.WithValue(ctx, traceIDKey, authCtx.TraceID)
	return ctx
}

// GetTenantID obtém o tenant ID do contexto
func GetTenantID(ctx context.Context) (uuid.UUID, bool) {
	tenantID, ok := ctx.Value(tenantIDKey).(uuid.UUID)
	return tenantID, ok
}

// GetUserID obtém o user ID do contexto
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}

// GetPlan obtém o plano do contexto
func GetPlan(ctx context.Context) (string, bool) {
	plan, ok := ctx.Value(planKey).(string)
	return plan, ok
}

// GetTraceID obtém o trace ID do contexto
func GetTraceID(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(traceIDKey).(string)
	return traceID, ok
}

// MustGetTenantID obtém o tenant ID do contexto ou panic
func MustGetTenantID(ctx context.Context) uuid.UUID {
	tenantID, ok := GetTenantID(ctx)
	if !ok {
		panic("tenant_id not found in context")
	}
	return tenantID
}

// MustGetUserID obtém o user ID do contexto ou panic
func MustGetUserID(ctx context.Context) uuid.UUID {
	userID, ok := GetUserID(ctx)
	if !ok {
		panic("user_id not found in context")
	}
	return userID
}