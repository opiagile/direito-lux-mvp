# AI Service - Inteligência Artificial para Análise Jurisprudencial

## 📋 Visão Geral

O **AI Service** é o microserviço de inteligência artificial do Direito Lux, responsável por análise jurisprudencial, busca semântica, geração de documentos e análise de texto legal. Implementado em **Python 3.11 + FastAPI**, oferece APIs robustas para processamento inteligente de dados jurídicos.

## 🚀 Funcionalidades Principais

### 🔍 Busca Semântica em Jurisprudência
- Busca por similaridade usando embeddings vetoriais
- Suporte a múltiplos tribunais (STF, STJ, TJ, TRT, TST, TRF)
- Filtros por tipo de decisão, data, área jurídica
- Cache Redis para performance otimizada

### 📊 Análise de Similaridade Multi-dimensional
- **Semântica**: Análise de conteúdo usando embeddings
- **Legal**: Comparação de fundamentos jurídicos 
- **Factual**: Análise de fatos e circunstâncias
- **Procedimental**: Comparação de ritos e procedimentos
- **Contextual**: Análise temporal e de contexto

### 📄 Análise de Documentos Legais
- Extração de entidades jurídicas (leis, artigos, valores, datas)
- Classificação automática por área jurídica
- Análise de risco jurídico
- Verificação de conformidade regulatória
- Processamento de texto jurídico brasileiro

### ✍️ Geração de Documentos
- Templates para contratos, petições, pareceres
- Personalização por variáveis dinâmicas
- Múltiplos níveis de qualidade (Draft, Standard, Professional, Premium)
- Validação de completude e precisão jurídica

## 🏗️ Arquitetura Técnica

### Stack Tecnológica
- **Framework**: FastAPI + Python 3.11
- **Validação**: Pydantic para modelos e configuração
- **Database**: SQLAlchemy + PostgreSQL com pgvector
- **Cache**: Redis com TTL configurável
- **ML/AI**: OpenAI + HuggingFace (com fallbacks)
- **Vector Store**: FAISS + pgvector para busca semântica
- **Containerização**: Docker multi-stage build

### Estrutura do Projeto
```
ai-service/
├── app/
│   ├── api/                    # Endpoints FastAPI
│   │   ├── health.py          # Health checks
│   │   ├── jurisprudence.py   # APIs de jurisprudência
│   │   ├── analysis.py        # APIs de análise
│   │   └── generation.py      # APIs de geração
│   ├── core/                  # Configuração e exceções
│   │   ├── config.py          # Settings com Pydantic
│   │   ├── exceptions.py      # Exceções customizadas
│   │   └── logging.py         # Logging estruturado
│   ├── db/                    # Database e modelos
│   │   ├── database.py        # Setup SQLAlchemy
│   │   └── models.py          # Modelos de dados
│   ├── models/                # Pydantic models
│   │   ├── jurisprudence.py   # Modelos de jurisprudência
│   │   ├── analysis.py        # Modelos de análise
│   │   └── generation.py      # Modelos de geração
│   ├── services/              # Business logic
│   │   ├── embeddings.py      # Geração de embeddings
│   │   ├── vector_store.py    # Busca vetorial
│   │   ├── text_processing.py # Processamento de texto
│   │   ├── cache.py          # Cache Redis
│   │   └── jurisprudence_service.py # Lógica de negócio
│   └── main.py               # Aplicação FastAPI
├── requirements.txt          # Dependências Python
├── pyproject.toml           # Configuração do projeto
├── Dockerfile               # Container Docker
└── .env.example            # Variáveis de ambiente
```

## 🌐 APIs Disponíveis

### Health Endpoints
- `GET /health` - Health check básico
- `GET /ready` - Readiness check com dependências

### Jurisprudence API (`/api/v1/jurisprudence/`)
- `POST /search` - Busca semântica em decisões judiciais
- `POST /similarity` - Análise de similaridade entre casos
- `POST /find-precedents` - Busca precedentes jurídicos
- `GET /courts` - Lista tipos de tribunais
- `GET /stats` - Estatísticas da base

### Analysis API (`/api/v1/analysis/`)
- `POST /analyze-document` - Análise completa de documentos
- `POST /analyze-process` - Análise de processos
- `GET /analysis-types` - Tipos de análise disponíveis

### Generation API (`/api/v1/generation/`)
- `POST /generate-document` - Geração de documentos
- `GET /document-types` - Tipos de documentos
- `GET /templates` - Templates disponíveis

## 🔧 Configuração e Deploy

### Variáveis de Ambiente

```bash
# Serviço
SERVICE_NAME=ai-service
VERSION=1.0.0
ENVIRONMENT=development
PORT=8000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=direito_lux_ai
DB_USER=postgres
DB_PASSWORD=your_password

# OpenAI (Opcional)
OPENAI_API_KEY=your_openai_key
OPENAI_MODEL=gpt-3.5-turbo
OPENAI_EMBEDDING_MODEL=text-embedding-ada-002

# HuggingFace (Opcional)
HUGGINGFACE_MODEL_NAME=sentence-transformers/all-MiniLM-L6-v2

# Vector Store
VECTOR_STORE_TYPE=pgvector  # ou faiss
EMBEDDING_DIMENSION=1536

# Cache
REDIS_URL=redis://localhost:6379
CACHE_TTL=3600

# Autenticação
JWT_SECRET_KEY=your_jwt_secret
```

### Docker Setup

```bash
# Build da imagem
docker build -t direito-lux/ai-service .

# Executar com docker-compose
docker-compose up ai-service

# Ou executar standalone
docker run -p 8000:8000 \
  -e DB_PASSWORD=password \
  -e OPENAI_API_KEY=key \
  -e JWT_SECRET_KEY=secret \
  direito-lux/ai-service
```

### Desenvolvimento Local

```bash
# Instalar dependências
pip install -r requirements.txt

# Configurar environment
cp .env.example .env
# Editar .env com suas configurações

# Executar aplicação
uvicorn app.main:app --reload --port 8000

# Testar aplicação
curl http://localhost:8000/health
```

## 🤖 Machine Learning & IA

### Modelos de Embeddings

**OpenAI (Opcional)**:
- `text-embedding-ada-002` - Embeddings de alta qualidade
- Suporte a grandes volumes de texto
- API robusta com rate limiting automático

**HuggingFace (Opcional)**:
- `sentence-transformers/all-MiniLM-L6-v2` - Modelo local
- Execução offline, sem custos de API
- Otimizado para texto em português

### Vector Stores

**pgvector (Recomendado)**:
- Extensão PostgreSQL para vetores
- Busca nativa no banco de dados
- Escalonabilidade e consistência

**FAISS (Alternativo)**:
- Biblioteca Facebook para busca vetorial
- Performance otimizada para grandes volumes
- Persistência em disco

### Processamento de Texto

**Funcionalidades**:
- Limpeza de texto jurídico brasileiro
- Extração de entidades (CPF, CNPJ, processos, valores)
- Classificação por área jurídica
- Extração de frases-chave
- Sumarização automática

## 📊 Integração com Planos

### Funcionalidades por Plano

| Funcionalidade | Starter | Professional | Business | Enterprise |
|---------------|---------|--------------|----------|------------|
| Busca Semântica | 10/mês | 50/mês | 200/mês | Ilimitado |
| Análise de Documentos | 5/mês | 25/mês | 100/mês | Ilimitado |
| Geração de Documentos | 3/mês | 15/mês | 50/mês | Ilimitado |
| Precedentes Jurídicos | STF/STJ | + TJs | + TRTs | Completo |
| Qualidade IA | Standard | Standard | Premium | Premium |

### Limites de Rate

```python
# Configuração por plano
RATE_LIMITS = {
    "starter": {"requests_per_minute": 10, "concurrent": 2},
    "professional": {"requests_per_minute": 30, "concurrent": 5},
    "business": {"requests_per_minute": 100, "concurrent": 10},
    "enterprise": {"requests_per_minute": 500, "concurrent": 20}
}
```

## 🔍 Monitoramento e Observabilidade

### Métricas Prometheus
- Request rate e latency por endpoint
- Cache hit/miss ratios
- Embedding generation metrics
- Vector search performance
- Error rates e tipos

### Health Checks
- Database connectivity
- Redis connectivity  
- Vector store status
- External APIs health

### Logging Estruturado
- Request/response tracing
- Performance metrics
- Error tracking
- Business events

## 🧪 Testes e Qualidade

### Testes Implementados
```bash
# Testar aplicação completa
DB_PASSWORD=test OPENAI_API_KEY=test JWT_SECRET_KEY=test \
python -c "from app.main import app; print('✅ AI Service OK')"

# Testar endpoints
curl -X POST http://localhost:8000/api/v1/jurisprudence/search \
  -H "Content-Type: application/json" \
  -d '{"query": "responsabilidade civil", "max_results": 5}'

# Testar health
curl http://localhost:8000/health
curl http://localhost:8000/ready
```

### Fallbacks e Resiliência
- Funciona sem bibliotecas ML (modo degradado)
- Cache para reduzir dependência de APIs externas
- Circuit breakers para APIs externas
- Graceful degradation por plano

## 🚀 Deploy e Produção

### Preparação para Produção

1. **Configurar Secrets**:
   ```bash
   # OpenAI API Key
   kubectl create secret generic ai-service-secrets \
     --from-literal=openai-api-key=your-key
   ```

2. **Database Migration**:
   ```bash
   # Executar migrações
   alembic upgrade head
   
   # Instalar extensão pgvector
   psql -c "CREATE EXTENSION IF NOT EXISTS vector;"
   ```

3. **Vector Store Setup**:
   ```bash
   # Inicializar FAISS index (se usado)
   python -c "from app.services.vector_store import init_vector_store; \
              import asyncio; asyncio.run(init_vector_store())"
   ```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ai-service
  template:
    metadata:
      labels:
        app: ai-service
    spec:
      containers:
      - name: ai-service
        image: direito-lux/ai-service:latest
        ports:
        - containerPort: 8000
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: ai-service-secrets
              key: db-password
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: ai-service-secrets
              key: openai-api-key
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
```

## 📚 Referências e Documentação

### Links Úteis
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [OpenAI API Reference](https://platform.openai.com/docs/api-reference)
- [HuggingFace Sentence Transformers](https://www.sbert.net/)
- [pgvector Documentation](https://github.com/pgvector/pgvector)
- [FAISS Documentation](https://faiss.ai/)

### Próximas Funcionalidades
- [ ] Análise de sentimento em decisões
- [ ] Predição de resultados processuais
- [ ] Integração com APIs de tribunais
- [ ] Fine-tuning de modelos para direito brasileiro
- [ ] Análise de jurisprudência por área específica
- [ ] Detecção de tendências jurisprudenciais

---

**🔄 Última Atualização**: 18/06/2025  
**📈 Status**: ✅ Implementado e Funcional  
**👨‍💻 Responsável**: Full Cycle Developer  
**🎯 Próximo**: Deploy em ambiente DEV