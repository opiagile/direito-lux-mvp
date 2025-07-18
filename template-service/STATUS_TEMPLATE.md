# STATUS - [NOME DO SERVIÇO]

## 📊 **Progresso Atual**
- **Status Geral**: [ ] Não iniciado | [ ] Em desenvolvimento | [ ] Completo
- **Percentual**: 0%
- **Última Atualização**: YYYY-MM-DD HH:MM
- **Desenvolvedor**: Claude + Usuário

## ✅ **O que está Implementado**
- [ ] Estrutura base do serviço (main.go, cmd/, internal/)
- [ ] Dockerfile e docker-compose entry
- [ ] Configuração de ambiente (.env, config)
- [ ] Modelos/Domain (structs, interfaces)
- [ ] Repositórios (database layer)
- [ ] Handlers/Controllers (HTTP endpoints)
- [ ] Rotas HTTP (router setup)
- [ ] Integração com RabbitMQ (events)
- [ ] Testes unitários (>80% coverage)
- [ ] Testes de integração (E2E)
- [ ] Documentação OpenAPI/Swagger
- [ ] Health check endpoint (/health)
- [ ] Métricas e observabilidade (Prometheus)
- [ ] Logs estruturados (zap/logrus)
- [ ] Middleware (auth, cors, rate limiting)
- [ ] Migrations de banco de dados
- [ ] Validação de dados (input validation)

## 🚧 **Em Desenvolvimento**
- [Listar itens em progresso no momento]

## ❌ **O que Falta**
- [Listar funcionalidades ainda não implementadas]

## 🐛 **Problemas Conhecidos**
- [Documentar bugs ou issues conhecidos]

## 📝 **Notas de Implementação**
- [Decisões técnicas importantes tomadas]
- [Dependências externas necessárias]
- [Configurações especiais requeridas]
- [Padrões ou convenções específicas adotadas]

## 🔗 **Dependências**
- **Depende de**: [listar outros serviços necessários]
- **É dependência de**: [listar serviços que dependem deste]
- **APIs externas**: [listar integrações externas]
- **Bibliotecas principais**: [listar libs Go importantes]

## 📋 **Checklist de Conclusão**
- [ ] Todos os endpoints implementados e testados
- [ ] Testes com cobertura > 80%
- [ ] Documentação completa (README + OpenAPI)
- [ ] Integração com outros serviços testada
- [ ] Deploy em ambiente dev funcionando
- [ ] Performance validada (load testing)
- [ ] Security review realizado
- [ ] Logs e métricas implementados
- [ ] Error handling robusto
- [ ] Graceful shutdown implementado

## 🚀 **Como Testar**
```bash
# Comandos para testar o serviço
cd services/[nome-service]
docker-compose up -d
curl http://localhost:[PORT]/health
# Adicionar mais comandos de teste relevantes
```

## 📊 **Métricas Importantes**
- **Tempo de resposta médio**: N/A
- **Taxa de erro**: N/A
- **Throughput**: N/A
- **Uso de memória**: N/A
- **Uso de CPU**: N/A

## 🔧 **Configuração Específica**
```yaml
# Variáveis de ambiente necessárias
VARIABLE_1=valor
VARIABLE_2=valor
# etc...
```

## 📚 **Documentação Relacionada**
- README.md do serviço
- OpenAPI/Swagger docs
- Diagramas de arquitetura
- Fluxos de dados

---

**📝 Lembrete**: Atualizar este arquivo **sempre** que implementar/modificar algo no serviço!