# =============================================================================
# Dockerfile para Template Service - Direito Lux
# =============================================================================

# Build stage
FROM golang:1.21-alpine AS builder

# Instalar dependências para build
RUN apk add --no-cache git ca-certificates tzdata

# Criar usuário não-root
RUN adduser -D -s /bin/sh -u 1001 appuser

# Definir diretório de trabalho
WORKDIR /app

# Copiar go mod e sum primeiro (para cache de layers)
COPY go.mod go.sum ./

# Download das dependências
RUN go mod download && go mod verify

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o bin/template-service \
    cmd/server/main.go

# =============================================================================
# Runtime stage
FROM scratch

# Copiar certificados CA
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copiar usuário
COPY --from=builder /etc/passwd /etc/passwd

# Copiar binário
COPY --from=builder /app/bin/template-service /template-service

# Criar diretórios necessários
# Note: no scratch não podemos usar RUN, então estrutura deve vir do builder
COPY --from=builder --chown=1001:1001 /tmp /tmp

# Usar usuário não-root
USER 1001

# Expor porta
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/template-service", "healthcheck"] || exit 1

# Comando para iniciar aplicação
ENTRYPOINT ["/template-service"]