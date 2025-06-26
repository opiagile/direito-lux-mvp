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
