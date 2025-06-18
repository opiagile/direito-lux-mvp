package domain

import "errors"

// Erros de notificação
var (
	ErrNotificationNotFound   = errors.New("notification not found")
	ErrNotificationCancelled  = errors.New("notification cancelled")
	ErrNotificationExpired    = errors.New("notification expired")
	ErrInvalidNotification    = errors.New("invalid notification")
	ErrNotificationExists     = errors.New("notification already exists")
)

// Erros de template
var (
	ErrTemplateNotFound    = errors.New("template not found")
	ErrTemplateInvalid     = errors.New("template invalid")
	ErrTemplateExists      = errors.New("template already exists")
	ErrSystemTemplate      = errors.New("cannot modify system template")
)

// Erros de preferência
var (
	ErrPreferenceNotFound = errors.New("preference not found")
	ErrPreferenceInvalid  = errors.New("preference invalid")
	ErrPreferenceExists   = errors.New("preference already exists")
)

// Erros de provedor
var (
	ErrProviderNotFound    = errors.New("provider not found")
	ErrProviderUnavailable = errors.New("provider unavailable")
	ErrProviderConfigInvalid = errors.New("provider configuration invalid")
	ErrRateLimitExceeded   = errors.New("rate limit exceeded")
)

// Erros de fila
var (
	ErrQueueEmpty       = errors.New("queue is empty")
	ErrQueueFull        = errors.New("queue is full")
	ErrQueueUnavailable = errors.New("queue unavailable")
)

// Erros de validação
var (
	ErrValidationFailed   = errors.New("validation failed")
	ErrRequiredField      = errors.New("required field missing")
	ErrInvalidFormat      = errors.New("invalid format")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrInvalidURL         = errors.New("invalid URL")
)

// Erros de autorização
var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
)

// Erros de configuração
var (
	ErrConfigMissing  = errors.New("configuration missing")
	ErrConfigInvalid  = errors.New("configuration invalid")
	ErrSecretMissing  = errors.New("secret missing")
)