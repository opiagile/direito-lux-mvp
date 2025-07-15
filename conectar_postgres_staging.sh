#!/bin/bash

echo "🔐 CONECTANDO AO POSTGRESQL STAGING GCP"
echo "======================================="

echo "1. 🔑 Reautenticando GCloud..."
echo "Execute este comando no seu terminal:"
echo "gcloud auth application-default login"
echo ""
read -p "Pressione ENTER após completar a autenticação..."

echo "2. ⚙️ Configurando kubectl..."
gcloud container clusters get-credentials direito-lux-gke-staging --region=us-central1 --project=direito-lux-staging-2025

echo "3. 🔍 Verificando status dos pods..."
kubectl get pods -n direito-lux-staging

echo "4. 🌐 Iniciando port-forward PostgreSQL..."
echo "MANTENHA ESTE TERMINAL ABERTO!"
echo ""
echo "📋 CONFIGURAÇÃO PGADMIN:"
echo "Host: localhost"
echo "Port: 5432"
echo "Database: direito_lux_staging"
echo "Username: direito_lux"
echo "Password: dev_password_123"
echo ""
echo "🚀 Iniciando port-forward..."

kubectl port-forward -n direito-lux-staging svc/postgres-service 5432:5432