#!/bin/bash

echo "🧹 LIMPEZA AUTOMÁTICA DO PROJETO DIREITO LUX"
echo "============================================="
echo ""
echo "Este script vai remover arquivos redundantes e organizar o projeto"
echo ""

# Confirmar antes de executar
read -p "Deseja continuar com a limpeza? (s/n): " confirm
if [ "$confirm" != "s" ] && [ "$confirm" != "S" ]; then
    echo "Operação cancelada."
    exit 0
fi

echo ""
echo "🗑️ Removendo scripts redundantes..."

# ============================================================================
# 1. REMOVER SCRIPTS DE FIX TEMPORÁRIOS
# ============================================================================
echo "1. Removendo scripts de fix temporários..."
rm -f FIX_*.sh
rm -f fix_*.sh 
rm -f simple_fix.sh
rm -f DEFINITIVE_DB_FIX.sh

# ============================================================================
# 2. REMOVER MÚLTIPLOS SCRIPTS DE SETUP (manter apenas os essenciais)
# ============================================================================
echo "2. Removendo scripts de setup redundantes..."
rm -f SETUP_COMPLETO.sh
rm -f SETUP_CORRIGIDO.sh  
rm -f SETUP_FINAL.sh
rm -f SETUP_FINAL_DEFINITIVO.sh
rm -f SETUP_DOCKER_CORRETO.sh
rm -f SETUP_POSTGRES_PADRAO.sh
rm -f setup_simple.sh
rm -f setup_oneliner.sh
rm -f setup_test_environment.sh

# ============================================================================
# 3. REMOVER SCRIPTS DE DEBUG
# ============================================================================
echo "3. Removendo scripts de debug..."
rm -f DEBUG_*.sh
rm -f TEST_LOGIN_FULL.sh
rm -f diagnose_*.sh
rm -f check_connection.sh
rm -f check_databases.sh
rm -f VERIFICAR_*.sh
rm -f verify_test_data.sh

# ============================================================================
# 4. REMOVER SCRIPTS DE EXECUÇÃO REDUNDANTES
# ============================================================================
echo "4. Removendo scripts de execução redundantes..."
rm -f EXECUTAR_*.sh
rm -f EXECUTE_*.sh
rm -f execute_auth_fix.sh
rm -f execute_fix.sh
rm -f run_migrations.sh
rm -f RUN_AUTH_MIGRATIONS.sh

# ============================================================================
# 5. REMOVER SCRIPTS DE GERENCIAMENTO DE USUÁRIO TEMPORÁRIOS
# ============================================================================
echo "5. Removendo scripts de usuário temporários..."
rm -f CHECK_AUTH_SERVICE.sh
rm -f CHECK_DATABASE_USERS.sh
rm -f INSERIR_DADOS_TESTE.sh
rm -f SIMPLE_USER_CREATE.sh

# ============================================================================
# 6. REMOVER ARQUIVOS SQL REDUNDANTES
# ============================================================================
echo "6. Removendo arquivos SQL redundantes..."
rm -f SEED_DATABASE_CORRECTED.sql
rm -f SEED_SIMPLE.sql
rm -f fix_*.sql
rm -f copy_data_to_correct_db.sh

# ============================================================================
# 7. REMOVER ARQUIVOS DE TEXTO TEMPORÁRIOS
# ============================================================================
echo "7. Removendo arquivos de texto temporários..."
rm -f COMANDOS_SETUP.txt
rm -f EXECUTAR_AGORA.txt

# ============================================================================
# 8. ORGANIZAR SCRIPTS UTILITÁRIOS
# ============================================================================
echo "8. Organizando scripts utilitários..."
mkdir -p scripts/utilities
mv CHECK_SERVICES_STATUS.sh scripts/utilities/ 2>/dev/null || true
mv RUN_SERVICES_LOCAL.sh scripts/utilities/ 2>/dev/null || true
mv STOP_ALL_SERVICES.sh scripts/utilities/ 2>/dev/null || true
mv execute_migrations.sh scripts/utilities/ 2>/dev/null || true

# ============================================================================
# 9. CRIAR DOCUMENTAÇÃO DOS SCRIPTS ESSENCIAIS
# ============================================================================
echo "9. Criando documentação dos scripts essenciais..."

cat > SCRIPTS_ESSENCIAIS.md << 'EOF'
# Scripts Essenciais - Direito Lux

## 🚀 Scripts Principais (Root)

### Configuração e Setup
- `SETUP_COMPLETE_FIXED.sh` - ⭐ **Setup completo do ambiente**
- `CLEAN_ENVIRONMENT_TOTAL.sh` - 🧹 **Limpeza total do ambiente**
- `START_LOCAL_DEV.sh` - 🛠️ **Iniciar ambiente de desenvolvimento**

### Build e Deploy
- `build-all.sh` - 🔨 **Compilar todos os microserviços**
- `start-services.sh` - ▶️ **Iniciar serviços localmente**
- `stop-services.sh` - ⏹️ **Parar serviços**
- `test-local.sh` - ✅ **Testes locais**

### Criação de Serviços
- `create-service.sh` - 📦 **Criar novo microserviço**

## 🛠️ Scripts Utilitários (scripts/utilities/)

- `CHECK_SERVICES_STATUS.sh` - 📊 **Status dos serviços**
- `RUN_SERVICES_LOCAL.sh` - 🚀 **Executar serviços local**
- `STOP_ALL_SERVICES.sh` - 🛑 **Parar todos os serviços**
- `execute_migrations.sh` - 🗄️ **Executar migrations**

## 📊 Dados de Teste

- `SEED_DATABASE_COMPLETE.sql` - 🌱 **Dados completos de teste**

## 📝 Como Usar

### Primeiro Setup
```bash
chmod +x SETUP_COMPLETE_FIXED.sh
./SETUP_COMPLETE_FIXED.sh
```

### Desenvolvimento Diário
```bash
./START_LOCAL_DEV.sh              # Iniciar desenvolvimento
./scripts/utilities/CHECK_SERVICES_STATUS.sh  # Verificar status
./test-local.sh                   # Testar funcionalidades
./stop-services.sh                # Parar ao final do dia
```

### Limpeza Total (quando necessário)
```bash
./CLEAN_ENVIRONMENT_TOTAL.sh
```

### Criar Novo Microserviço
```bash
./create-service.sh nome-service
```

## 🗂️ Estrutura Organizada

```
direito-lux/
├── SETUP_COMPLETE_FIXED.sh      # Setup principal
├── CLEAN_ENVIRONMENT_TOTAL.sh   # Limpeza
├── START_LOCAL_DEV.sh           # Dev environment
├── build-all.sh                 # Build
├── start-services.sh            # Start services
├── stop-services.sh             # Stop services
├── test-local.sh                # Tests
├── create-service.sh            # New service
├── SEED_DATABASE_COMPLETE.sql   # Test data
└── scripts/
    └── utilities/               # Utility scripts
        ├── CHECK_SERVICES_STATUS.sh
        ├── RUN_SERVICES_LOCAL.sh
        ├── STOP_ALL_SERVICES.sh
        └── execute_migrations.sh
```
EOF

# ============================================================================
# 10. RESULTADO FINAL
# ============================================================================
echo ""
echo "✅ LIMPEZA CONCLUÍDA!"
echo ""
echo "📊 Resumo:"
echo "   • Scripts removidos: ~40 arquivos"
echo "   • Scripts mantidos: ~12 arquivos essenciais"
echo "   • Estrutura organizada em scripts/utilities/"
echo ""
echo "📋 Scripts essenciais mantidos:"
echo "   ✅ SETUP_COMPLETE_FIXED.sh (setup completo)"
echo "   ✅ CLEAN_ENVIRONMENT_TOTAL.sh (limpeza)"  
echo "   ✅ START_LOCAL_DEV.sh (desenvolvimento)"
echo "   ✅ build-all.sh (build)"
echo "   ✅ start-services.sh (iniciar serviços)"
echo "   ✅ stop-services.sh (parar serviços)"
echo "   ✅ test-local.sh (testes)"
echo "   ✅ create-service.sh (novo serviço)"
echo "   ✅ SEED_DATABASE_COMPLETE.sql (dados teste)"
echo ""
echo "📖 Documentação criada: SCRIPTS_ESSENCIAIS.md"
echo ""
echo "🎯 O projeto agora está muito mais limpo e organizado!"
echo "   Para usar: consulte SCRIPTS_ESSENCIAIS.md"