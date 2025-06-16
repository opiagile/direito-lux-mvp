# Processo de Documentação - Direito Lux

## 📋 Política de Documentação

**IMPORTANTE**: A documentação deve ser tratada como código - sempre atualizada, revisada e versionada.

## 🔄 Ciclo de Atualização

### Ao Finalizar Cada Módulo/Serviço

Após completar a implementação de qualquer componente, os seguintes documentos **DEVEM** ser atualizados:

#### 1. STATUS_IMPLEMENTACAO.md
- [ ] Mover item de "❌ O que Falta" para "✅ O que está Implementado"
- [ ] Atualizar percentual de progresso
- [ ] Adicionar detalhes específicos do que foi implementado
- [ ] Atualizar "Próximos Passos Recomendados"

#### 2. README.md
- [ ] Atualizar seção "📊 Status do Projeto"
- [ ] Adicionar novo serviço na arquitetura (se aplicável)
- [ ] Atualizar URLs de desenvolvimento
- [ ] Atualizar comandos úteis

#### 3. SETUP_AMBIENTE.md
- [ ] Adicionar instruções de setup do novo módulo
- [ ] Incluir novas variáveis de ambiente
- [ ] Adicionar URLs e credenciais
- [ ] Atualizar seção de troubleshooting com problemas encontrados

#### 4. Documentação Específica do Módulo
- [ ] Criar README.md no diretório do serviço
- [ ] Documentar APIs (OpenAPI/Swagger)
- [ ] Incluir exemplos de uso
- [ ] Documentar decisões arquiteturais

## 📝 Template de Atualização

### Para STATUS_IMPLEMENTACAO.md

```markdown
### N. [Nome do Serviço] (Completo)
- ✅ **services/[nome-service]/** - [Descrição breve]:
  
  **Domain Layer:**
  - `[arquivo].go` - [Descrição]
  
  **Application Layer:**
  - `[arquivo].go` - [Descrição]
  
  **Infrastructure Layer:**
  - [Componentes implementados]
  
  **Migrações:**
  - `00N_[descrição].sql`
  
  **APIs:**
  - GET/POST/PUT/DELETE /api/v1/[recurso]
  
  **Funcionalidades:**
  - [Lista de features implementadas]
```

### Para README.md

```markdown
### ✅ Implementado
- ✅ [Nome do Serviço] - [Descrição breve]

### 🔗 URLs de Desenvolvimento
| **[Nome] Service** | http://localhost:[porta] | [credenciais se houver] |

**Progresso Total**: ~[X]% completo
```

## 🎯 Checklist por Tipo de Componente

### Microserviço Go
- [ ] Atualizar docker-compose.yml
- [ ] Documentar estrutura de pastas
- [ ] Listar endpoints da API
- [ ] Documentar eventos publicados/consumidos
- [ ] Incluir exemplos de requisições
- [ ] Atualizar diagrama de arquitetura

### Serviço Python/AI
- [ ] Documentar dependências (requirements.txt)
- [ ] Explicar modelos utilizados
- [ ] Documentar API REST
- [ ] Incluir exemplos de input/output

### Frontend
- [ ] Documentar componentes principais
- [ ] Listar rotas/páginas
- [ ] Documentar integração com backend
- [ ] Incluir screenshots

### Infraestrutura
- [ ] Atualizar configurações Terraform
- [ ] Documentar recursos criados
- [ ] Atualizar variáveis de ambiente
- [ ] Documentar procedimentos de deploy

## 📊 Métricas de Documentação

Manter registro de:
- Data da última atualização
- Autor da atualização
- Versão do componente documentado
- Links para PRs relacionados

## 🚨 Regras Importantes

1. **Nunca fazer merge sem atualizar a documentação**
2. **Documentação desatualizada é pior que não ter documentação**
3. **Use exemplos reais sempre que possível**
4. **Mantenha consistência de formato**
5. **Revise a documentação em cada PR**

## 📅 Revisão Periódica

### Semanal
- Verificar se STATUS_IMPLEMENTACAO.md reflete realidade
- Atualizar progresso percentual

### Mensal
- Revisar toda documentação
- Atualizar diagramas
- Verificar links quebrados
- Atualizar screenshots

## 🔗 Documentos a Manter Atualizados

1. **STATUS_IMPLEMENTACAO.md** - Status geral do projeto
2. **README.md** - Visão geral e quick start
3. **SETUP_AMBIENTE.md** - Guia de instalação
4. **ARQUITETURA_FULLCYCLE.md** - Quando houver mudanças arquiteturais
5. **ROADMAP_IMPLEMENTACAO.md** - Ajustar prazos conforme necessário
6. **Service-specific README** - Um por serviço
7. **API Documentation** - Swagger/OpenAPI
8. **CHANGELOG.md** - Histórico de mudanças

## ⚡ Automação Futura

Considerar implementar:
- [ ] CI check para documentação atualizada
- [ ] Geração automática de changelog
- [ ] Documentação de API auto-gerada
- [ ] Status badges automáticos
- [ ] Métricas de cobertura de documentação

---

**Lembrete**: Este documento deve ser consultado ao final de CADA implementação!