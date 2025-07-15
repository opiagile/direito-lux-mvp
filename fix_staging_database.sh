#!/bin/bash

echo "🔧 CORRIGINDO DATABASE STAGING - CARREGANDO DADOS DE TESTE"
echo "=============================================================="

# Reautenticar GCloud
echo "1. 🔐 Reautenticando GCloud..."
gcloud auth application-default login

# Configurar kubectl
echo "2. ⚙️  Configurando kubectl..."
gcloud container clusters get-credentials direito-lux-gke-staging --region=us-central1 --project=direito-lux-staging-2025

# Verificar se PostgreSQL está rodando
echo "3. 🔍 Verificando status PostgreSQL..."
kubectl get pods -n direito-lux-staging -l app=postgres

# Port-forward PostgreSQL em background
echo "4. 🌐 Criando port-forward PostgreSQL..."
kubectl port-forward -n direito-lux-staging svc/postgres-service 5433:5432 &
PORT_FORWARD_PID=$!
sleep 5

# Verificar conexão
echo "5. 🧪 Testando conexão PostgreSQL..."
PGPASSWORD=dev_password_123 psql -h localhost -p 5433 -U direito_lux -d direito_lux_staging -c "SELECT 1;" 2>/dev/null

if [ $? -eq 0 ]; then
    echo "✅ PostgreSQL acessível!"
    
    # Verificar se tabelas existem
    echo "6. 📋 Verificando tabelas existentes..."
    PGPASSWORD=dev_password_123 psql -h localhost -p 5433 -U direito_lux -d direito_lux_staging -c "\dt" 
    
    # Verificar dados na tabela users
    echo "7. 👥 Verificando usuários cadastrados..."
    PGPASSWORD=dev_password_123 psql -h localhost -p 5433 -U direito_lux -d direito_lux_staging -c "SELECT email, created_at FROM users LIMIT 5;" 2>/dev/null
    
    if [ $? -ne 0 ]; then
        echo "❌ Tabela users não existe ou sem dados!"
        echo "8. 🔄 Executando migrations e seed data..."
        
        # Executar migrations via auth-service pod
        echo "Executando migrations no auth-service..."
        AUTH_POD=$(kubectl get pods -n direito-lux-staging -l app=auth-service -o jsonpath="{.items[0].metadata.name}")
        kubectl exec -n direito-lux-staging $AUTH_POD -- /bin/sh -c "cd /app && migrate -path migrations -database 'postgres://direito_lux:dev_password_123@postgres-service:5432/direito_lux_staging?sslmode=disable' up"
        
        # Verificar novamente
        echo "9. ✅ Verificando dados após migrations..."
        PGPASSWORD=dev_password_123 psql -h localhost -p 5433 -U direito_lux -d direito_lux_staging -c "SELECT email, created_at FROM users LIMIT 5;"
    else
        echo "✅ Dados encontrados na tabela users!"
    fi
    
else
    echo "❌ Erro ao conectar PostgreSQL!"
    echo "Verificando logs do PostgreSQL..."
    kubectl logs -n direito-lux-staging -l app=postgres --tail=20
fi

# Cleanup
echo "10. 🧹 Finalizando..."
kill $PORT_FORWARD_PID 2>/dev/null

echo ""
echo "🎯 PRÓXIMOS PASSOS:"
echo "1. Execute: kubectl port-forward -n direito-lux-staging svc/postgres-service 5432:5432"
echo "2. Conecte pgAdmin em localhost:5432"
echo "3. Use credenciais: direito_lux/dev_password_123/direito_lux_staging"
echo "4. Teste login no sistema: admin@silvaassociados.com.br/password"