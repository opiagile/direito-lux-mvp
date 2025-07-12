#!/bin/bash

# 🔐 Setup Profissional de Gerenciamento de Segredos
# Direito Lux - Solução Enterprise-Grade

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Banner
show_banner() {
    echo -e "${CYAN}"
    echo "╔═══════════════════════════════════════════════════╗"
    echo "║    🔐 SETUP PROFISSIONAL DE SEGREDOS             ║"
    echo "║    Enterprise-Grade Secret Management             ║"
    echo "╚═══════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

log_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }
log_success() { echo -e "${GREEN}✅ $1${NC}"; }
log_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
log_error() { echo -e "${RED}❌ $1${NC}"; }
log_step() { echo -e "${PURPLE}📋 $1${NC}"; }

# Verificar dependências
check_dependencies() {
    log_step "Verificando dependências..."
    
    # GitHub CLI
    if ! command -v gh &> /dev/null; then
        log_error "GitHub CLI não encontrado. Instale: brew install gh"
        exit 1
    fi
    
    # kubectl
    if ! command -v kubectl &> /dev/null; then
        log_warning "kubectl não encontrado. Algumas funcionalidades serão limitadas."
    fi
    
    # Helm
    if ! command -v helm &> /dev/null; then
        log_warning "Helm não encontrado. Algumas funcionalidades serão limitadas."
    fi
    
    log_success "Dependências verificadas"
}

# Menu de opções
show_menu() {
    echo ""
    log_step "Escolha a solução de segredos:"
    echo ""
    echo "1. 🚀 QUICK START - GitHub Secrets (Recomendado para começar)"
    echo "2. 🏢 ENTERPRISE - External Secrets + Google Secret Manager"
    echo "3. 🔐 MAXIMUM - HashiCorp Vault completo"
    echo "4. 🔄 GITOPS - Sealed Secrets"
    echo "5. 📊 COMPARAR - Ver todas as opções"
    echo ""
    read -p "Opção (1-5): " choice
}

# Opção 1: GitHub Secrets
setup_github_secrets() {
    log_step "Configurando GitHub Secrets..."
    
    # Login no GitHub se necessário
    if ! gh auth status &> /dev/null; then
        log_warning "Fazendo login no GitHub..."
        gh auth login
    fi
    
    echo ""
    log_step "Configure os seguintes secrets no GitHub:"
    echo ""
    
    echo -e "${YELLOW}1. Telegram Bot Token:${NC}"
    read -s -p "TELEGRAM_BOT_TOKEN: " TELEGRAM_TOKEN
    echo ""
    
    echo -e "${YELLOW}2. WhatsApp Access Token:${NC}"
    read -s -p "WHATSAPP_ACCESS_TOKEN: " WHATSAPP_TOKEN
    echo ""
    
    echo -e "${YELLOW}3. OpenAI API Key (opcional):${NC}"
    read -s -p "OPENAI_API_KEY: " OPENAI_KEY
    echo ""
    
    # Definir secrets
    log_step "Definindo secrets no GitHub..."
    
    if [ ! -z "$TELEGRAM_TOKEN" ]; then
        gh secret set TELEGRAM_BOT_TOKEN --body "$TELEGRAM_TOKEN"
        log_success "TELEGRAM_BOT_TOKEN definido"
    fi
    
    if [ ! -z "$WHATSAPP_TOKEN" ]; then
        gh secret set WHATSAPP_ACCESS_TOKEN --body "$WHATSAPP_TOKEN"
        log_success "WHATSAPP_ACCESS_TOKEN definido"
    fi
    
    if [ ! -z "$OPENAI_KEY" ]; then
        gh secret set OPENAI_API_KEY --body "$OPENAI_KEY"
        log_success "OPENAI_API_KEY definido"
    fi
    
    # Criar workflow que usa os secrets
    create_github_workflow
    
    log_success "GitHub Secrets configurado com sucesso!"
}

# Criar workflow do GitHub Actions
create_github_workflow() {
    log_step "Criando workflow de deploy seguro..."
    
    mkdir -p .github/workflows
    
    cat > .github/workflows/deploy-with-secrets.yml << EOF
name: Deploy with Secure Secrets

on:
  push:
    branches: [main]
  workflow_dispatch:

env:
  TELEGRAM_BOT_TOKEN: \${{ secrets.TELEGRAM_BOT_TOKEN }}
  WHATSAPP_ACCESS_TOKEN: \${{ secrets.WHATSAPP_ACCESS_TOKEN }}
  OPENAI_API_KEY: \${{ secrets.OPENAI_API_KEY }}

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Environment
      run: |
        echo "Setting up secure environment..."
        # Secrets são automaticamente mascarados nos logs
        
    - name: Deploy to Staging
      run: |
        echo "Deploying with secure secrets..."
        # Usar secrets aqui sem exposição
        
    - name: Test Deployment
      run: |
        echo "Testing deployment..."
        # Verificações sem expor tokens
EOF

    log_success "Workflow criado: .github/workflows/deploy-with-secrets.yml"
}

# Opção 2: External Secrets + Google Secret Manager
setup_external_secrets() {
    log_step "Configurando External Secrets Operator + Google Secret Manager..."
    
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl necessário para esta opção"
        exit 1
    fi
    
    if ! command -v helm &> /dev/null; then
        log_error "Helm necessário para esta opção"
        exit 1
    fi
    
    # Instalar External Secrets Operator
    log_step "Instalando External Secrets Operator..."
    
    helm repo add external-secrets https://charts.external-secrets.io
    helm repo update
    
    kubectl create namespace external-secrets-system --dry-run=client -o yaml | kubectl apply -f -
    
    helm upgrade --install external-secrets \
        external-secrets/external-secrets \
        -n external-secrets-system \
        --create-namespace
    
    # Aguardar deployment
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=external-secrets -n external-secrets-system --timeout=300s
    
    log_success "External Secrets Operator instalado"
    
    # Configurar Google Secret Manager
    setup_google_secret_manager
    
    # Criar ExternalSecret
    create_external_secret_manifests
    
    log_success "External Secrets configurado com sucesso!"
}

# Setup Google Secret Manager
setup_google_secret_manager() {
    log_step "Configurando Google Secret Manager..."
    
    echo ""
    log_warning "Para configurar o Google Secret Manager:"
    echo "1. Acesse: https://console.cloud.google.com/security/secret-manager"
    echo "2. Ative a API Secret Manager"
    echo "3. Crie os seguintes secrets:"
    echo "   - telegram-bot-token"
    echo "   - whatsapp-access-token"
    echo "   - openai-api-key"
    echo ""
    
    read -p "Pressione Enter quando terminar de configurar no GCP..."
    
    # Service Account Key
    echo ""
    log_step "Configure a Service Account:"
    echo "1. IAM & Admin → Service Accounts"
    echo "2. Create Service Account: 'external-secrets-sa'"
    echo "3. Role: Secret Manager Secret Accessor"
    echo "4. Download JSON key"
    echo ""
    
    read -p "Caminho para o JSON da Service Account: " SA_KEY_PATH
    
    if [ -f "$SA_KEY_PATH" ]; then
        kubectl create secret generic gcpsm-secret \
            --from-file=secret-access-credentials="$SA_KEY_PATH" \
            -n external-secrets-system
        log_success "Service Account configurada"
    else
        log_error "Arquivo de Service Account não encontrado"
        exit 1
    fi
}

# Criar manifests do External Secret
create_external_secret_manifests() {
    log_step "Criando manifests do External Secret..."
    
    mkdir -p k8s/secrets
    
    # SecretStore
    cat > k8s/secrets/secret-store.yaml << EOF
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: gcpsm-secret-store
  namespace: default
spec:
  provider:
    gcpsm:
      projectId: "seu-project-id"
      auth:
        secretRef:
          secretAccessKey:
            name: gcpsm-secret
            key: secret-access-credentials
            namespace: external-secrets-system
EOF

    # ExternalSecret
    cat > k8s/secrets/notification-external-secret.yaml << EOF
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: notification-secrets
  namespace: default
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: gcpsm-secret-store
    kind: SecretStore
  target:
    name: notification-secrets
    creationPolicy: Owner
  data:
  - secretKey: telegram-token
    remoteRef:
      key: telegram-bot-token
  - secretKey: whatsapp-token
    remoteRef:
      key: whatsapp-access-token
  - secretKey: openai-key
    remoteRef:
      key: openai-api-key
EOF

    log_success "Manifests criados em k8s/secrets/"
    
    echo ""
    log_warning "Para aplicar:"
    echo "1. Edite k8s/secrets/secret-store.yaml com seu Project ID"
    echo "2. kubectl apply -f k8s/secrets/"
}

# Opção 3: HashiCorp Vault
setup_vault() {
    log_step "Configurando HashiCorp Vault..."
    
    if ! command -v helm &> /dev/null; then
        log_error "Helm necessário para esta opção"
        exit 1
    fi
    
    # Instalar Vault
    log_step "Instalando HashiCorp Vault..."
    
    helm repo add hashicorp https://helm.releases.hashicorp.com
    helm repo update
    
    # Valores para desenvolvimento
    cat > vault-values.yaml << EOF
server:
  dev:
    enabled: true
    devRootToken: "root"
  
  dataStorage:
    enabled: true
    size: 10Gi
    
ui:
  enabled: true
  serviceType: "LoadBalancer"
EOF

    kubectl create namespace vault --dry-run=client -o yaml | kubectl apply -f -
    
    helm upgrade --install vault hashicorp/vault \
        -n vault \
        --create-namespace \
        -f vault-values.yaml
    
    # Aguardar deployment
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=vault -n vault --timeout=300s
    
    log_success "HashiCorp Vault instalado"
    
    # Configurar secrets no Vault
    configure_vault_secrets
    
    log_success "HashiCorp Vault configurado com sucesso!"
}

# Configurar secrets no Vault
configure_vault_secrets() {
    log_step "Configurando secrets no Vault..."
    
    # Port forward para acessar Vault
    kubectl port-forward svc/vault -n vault 8200:8200 &
    PF_PID=$!
    
    sleep 5
    
    export VAULT_ADDR="http://localhost:8200"
    export VAULT_TOKEN="root"
    
    # Ativar KV engine
    kubectl exec -n vault vault-0 -- vault secrets enable -path=secret kv-v2
    
    echo ""
    log_step "Digite os secrets para armazenar no Vault:"
    
    read -s -p "Telegram Bot Token: " TELEGRAM_TOKEN
    echo ""
    read -s -p "WhatsApp Access Token: " WHATSAPP_TOKEN
    echo ""
    
    # Armazenar secrets
    kubectl exec -n vault vault-0 -- vault kv put secret/notification \
        telegram-token="$TELEGRAM_TOKEN" \
        whatsapp-token="$WHATSAPP_TOKEN"
    
    # Parar port forward
    kill $PF_PID
    
    log_success "Secrets armazenados no Vault"
}

# Opção 4: Sealed Secrets
setup_sealed_secrets() {
    log_step "Configurando Sealed Secrets..."
    
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl necessário para esta opção"
        exit 1
    fi
    
    # Instalar Sealed Secrets Controller
    kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml
    
    # Aguardar deployment
    kubectl wait --for=condition=ready pod -l name=sealed-secrets-controller -n kube-system --timeout=300s
    
    # Instalar kubeseal CLI
    if ! command -v kubeseal &> /dev/null; then
        log_step "Instalando kubeseal CLI..."
        
        KUBESEAL_VERSION='0.24.0'
        wget "https://github.com/bitnami-labs/sealed-secrets/releases/download/v${KUBESEAL_VERSION}/kubeseal-${KUBESEAL_VERSION}-linux-amd64.tar.gz"
        tar xfz "kubeseal-${KUBESEAL_VERSION}-linux-amd64.tar.gz"
        sudo install -m 755 kubeseal /usr/local/bin/kubeseal
        rm kubeseal*
        
        log_success "kubeseal instalado"
    fi
    
    # Criar sealed secrets
    create_sealed_secrets
    
    log_success "Sealed Secrets configurado com sucesso!"
}

# Criar sealed secrets
create_sealed_secrets() {
    log_step "Criando Sealed Secrets..."
    
    mkdir -p k8s/sealed-secrets
    
    echo ""
    read -s -p "Telegram Bot Token: " TELEGRAM_TOKEN
    echo ""
    read -s -p "WhatsApp Access Token: " WHATSAPP_TOKEN
    echo ""
    
    # Criar secret regular e selá-lo
    echo -n "$TELEGRAM_TOKEN" | kubectl create secret generic telegram-bot-secret \
        --dry-run=client --from-file=token=/dev/stdin -o yaml | \
        kubeseal -o yaml > k8s/sealed-secrets/telegram-bot-sealed.yaml
    
    echo -n "$WHATSAPP_TOKEN" | kubectl create secret generic whatsapp-api-secret \
        --dry-run=client --from-file=token=/dev/stdin -o yaml | \
        kubeseal -o yaml > k8s/sealed-secrets/whatsapp-api-sealed.yaml
    
    log_success "Sealed Secrets criados em k8s/sealed-secrets/"
    
    echo ""
    log_warning "Para aplicar:"
    echo "kubectl apply -f k8s/sealed-secrets/"
}

# Comparar opções
compare_options() {
    echo ""
    log_step "Comparação de Soluções de Segredos:"
    echo ""
    
    echo "┌─────────────────┬─────────────┬──────────┬───────────┬─────────────┐"
    echo "│ Solução         │ Complexidade│ Custo    │ Segurança │ Compliance  │"
    echo "├─────────────────┼─────────────┼──────────┼───────────┼─────────────┤"
    echo "│ GitHub Secrets  │ ⭐⭐         │ Free     │ ⭐⭐⭐      │ ⭐⭐         │"
    echo "│ External Secrets│ ⭐⭐⭐        │ Low      │ ⭐⭐⭐⭐     │ ⭐⭐⭐        │"
    echo "│ HashiCorp Vault │ ⭐⭐⭐⭐⭐     │ Medium   │ ⭐⭐⭐⭐⭐    │ ⭐⭐⭐⭐⭐     │"
    echo "│ Sealed Secrets  │ ⭐⭐⭐        │ Free     │ ⭐⭐⭐⭐     │ ⭐⭐⭐⭐      │"
    echo "└─────────────────┴─────────────┴──────────┴───────────┴─────────────┘"
    echo ""
    
    echo -e "${GREEN}Recomendações:${NC}"
    echo "🚀 Começar: GitHub Secrets"
    echo "🏢 Crescer: External Secrets + Google Secret Manager"
    echo "🔐 Enterprise: HashiCorp Vault"
    echo ""
}

# Cleanup
cleanup() {
    log_step "Limpando arquivos temporários..."
    rm -f vault-values.yaml
    rm -f kubeseal*
}

# Main
main() {
    show_banner
    check_dependencies
    
    while true; do
        show_menu
        
        case $choice in
            1)
                setup_github_secrets
                break
                ;;
            2)
                setup_external_secrets
                break
                ;;
            3)
                setup_vault
                break
                ;;
            4)
                setup_sealed_secrets
                break
                ;;
            5)
                compare_options
                ;;
            *)
                log_error "Opção inválida"
                ;;
        esac
    done
    
    cleanup
    
    echo ""
    log_success "🎉 Setup de segredos concluído!"
    echo ""
    log_info "📚 Documentação completa: SECRETS_MANAGEMENT_OPTIONS.md"
}

# Executar se chamado diretamente
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi