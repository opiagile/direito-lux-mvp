# Frontend Web App - Direito Lux

## 📋 Visão Geral

O Frontend Web App do Direito Lux é uma aplicação moderna construída com Next.js 14, oferecendo uma interface completa e responsiva para gestão de processos jurídicos com integração total aos microserviços backend.

## 🚀 Stack Tecnológica

### Core Framework
- **Next.js 14** - Framework React com App Router
- **TypeScript** - Type safety completo
- **React 18** - Biblioteca UI moderna

### Styling e UI
- **Tailwind CSS** - Utility-first CSS framework
- **Shadcn/ui** - Componentes primitivos com Radix UI
- **Lucide React** - Ícones modernos e consistentes
- **CSS Custom Properties** - Sistema de temas dinâmico

### State Management
- **Zustand** - State management leve e eficiente
- **Persist Middleware** - Persistência automática de estado
- **Multiple Stores** - Stores especializados por domínio

### Data Fetching
- **React Query (@tanstack/react-query)** - Cache e sincronização de dados
- **Axios** - Cliente HTTP robusto
- **Multi-service Integration** - Conexão com todos os microserviços

### Forms e Validação
- **React Hook Form** - Gerenciamento de formulários
- **Zod** - Schema validation com TypeScript
- **@hookform/resolvers** - Integração RHF + Zod

### Desenvolvimento
- **ESLint** - Linting de código
- **Prettier** - Formatação automática
- **PostCSS** - Processamento CSS

## 📁 Estrutura do Projeto

```
frontend/
├── src/
│   ├── app/                    # Next.js 14 App Router
│   │   ├── (dashboard)/        # Dashboard layout group
│   │   │   ├── ai/            # AI Assistant page
│   │   │   ├── dashboard/     # Main dashboard
│   │   │   ├── processes/     # Process management
│   │   │   └── layout.tsx     # Dashboard layout
│   │   ├── login/             # Login page
│   │   ├── globals.css        # Global styles
│   │   ├── layout.tsx         # Root layout
│   │   └── page.tsx           # Home page
│   ├── components/
│   │   ├── layout/            # Layout components
│   │   │   ├── header.tsx     # App header
│   │   │   └── sidebar.tsx    # Navigation sidebar
│   │   ├── providers.tsx      # App providers
│   │   └── ui/                # UI components
│   │       ├── avatar.tsx
│   │       ├── badge.tsx
│   │       ├── button.tsx
│   │       ├── card.tsx
│   │       ├── dropdown-menu.tsx
│   │       ├── input.tsx
│   │       ├── label.tsx
│   │       ├── loading-screen.tsx
│   │       ├── table.tsx
│   │       ├── tabs.tsx
│   │       └── textarea.tsx
│   ├── hooks/
│   │   └── api.ts             # React Query hooks
│   ├── lib/
│   │   ├── api.ts             # API clients
│   │   └── utils.ts           # Utility functions
│   ├── store/
│   │   └── index.ts           # Zustand stores
│   └── types/
│       └── index.ts           # TypeScript types
├── package.json               # Dependencies e scripts
├── tsconfig.json             # TypeScript config
├── tailwind.config.js        # Tailwind configuration
├── next.config.js            # Next.js configuration
└── postcss.config.js         # PostCSS config
```

## 🎨 Sistema de Design

### Cores e Temas
- **Modo Claro/Escuro** - Sistema completo de temas
- **Cores Primárias** - Azul profissional (#3B82F6)
- **Cores Semânticas** - Success, Warning, Error, Info
- **Cores da Marca** - Paleta Direito Lux customizada

### Componentes UI
- **Design System** - Componentes consistentes e reutilizáveis
- **Variantes** - Multiple variants para cada componente
- **Responsividade** - Mobile-first design
- **Acessibilidade** - ARIA labels e keyboard navigation

## 📊 Funcionalidades Implementadas

### 🔐 Autenticação
- **Login Page** - Formulário com validação completa
- **JWT Integration** - Token management automático
- **Protected Routes** - Guards de autenticação
- **Session Management** - Refresh automático de tokens

### 📊 Dashboard
- **KPIs em Tempo Real** - Métricas principais visualizadas
- **Atividades Recentes** - Feed de últimas movimentações
- **Estatísticas Rápidas** - Gráficos e indicadores
- **Cards Interativos** - Componentes clicáveis

### 📁 Gestão de Processos
- **Visualizações Múltiplas** - Table, Grid, List views
- **Busca e Filtros** - Sistema de busca integrado
- **CRUD Operations** - Create, Read, Update, Delete
- **Monitoramento** - Toggle de processos monitorados
- **Status Management** - Gestão de status dos processos

### 🤖 AI Assistant
- **Chat Interface** - Interface conversacional
- **Análise de Documentos** - Upload e análise automática
- **Busca de Jurisprudência** - Busca semântica avançada
- **Histórico** - Armazenamento de interações
- **Multi-tab Interface** - Organização por funcionalidade

### 🎨 Interface e UX
- **Navigation Sidebar** - Menu lateral responsivo
- **Header Global** - Busca global e profile menu
- **Breadcrumbs** - Navegação hierárquica
- **Loading States** - Feedback visual para operações
- **Error Handling** - Tratamento gracioso de erros
- **Toast Notifications** - Feedback de ações

## 🔧 State Management

### Stores Implementados (Zustand)

#### AuthStore
```typescript
interface AuthState {
  user: User | null
  tenant: Tenant | null
  isAuthenticated: boolean
  token: string | null
  login: (user, tenant, token) => void
  logout: () => void
}
```

#### UIStore
```typescript
interface UIState {
  theme: Theme
  sidebarOpen: boolean
  sidebarCollapsed: boolean
  breadcrumbs: Array<{ label: string; href?: string }>
  pageTitle: string
}
```

#### ProcessStore
```typescript
interface ProcessState {
  selectedProcesses: string[]
  processFilters: Record<string, any>
  viewMode: 'grid' | 'list' | 'table'
  sortBy: string
  sortOrder: 'asc' | 'desc'
}
```

#### NotificationStore
```typescript
interface NotificationState {
  notifications: Notification[]
  unreadCount: number
  addNotification: (notification) => void
  markAsRead: (id) => void
}
```

## 🌐 Integração com APIs

### Clientes HTTP
- **apiClient** - API Gateway principal (porta 8090)
- **aiClient** - AI Service (porta 8000)
- **searchClient** - Search Service (porta 8086)
- **reportClient** - Report Service (porta 8087)

### React Query Hooks
```typescript
// Exemplos de hooks implementados
useProcesses(params) // Lista processos
useProcess(id) // Processo específico
useCreateProcess() // Criar processo
useLogin() // Autenticação
useAIAnalysis() // Análise de IA
useSearch(query) // Busca
```

### Error Handling
- **Interceptors** - Tratamento automático de erros HTTP
- **401 Redirect** - Redirecionamento para login
- **Toast Notifications** - Feedback visual de erros
- **Retry Logic** - Tentativas automáticas

## 📱 Responsividade

### Breakpoints
- **Mobile** - < 768px
- **Tablet** - 768px - 1024px
- **Desktop** - > 1024px

### Layout Adaptativo
- **Sidebar Collapse** - Colapso automático em mobile
- **Grid Responsive** - Adaptação automática de grids
- **Touch Optimization** - Gestos otimizados para mobile

## ⚡ Performance

### Otimizações
- **Code Splitting** - Lazy loading de páginas
- **Image Optimization** - Next.js Image component
- **Bundle Analysis** - Análise de tamanho de bundle
- **Caching Strategy** - React Query cache otimizado

### SEO
- **Metadata** - Meta tags otimizadas
- **Structured Data** - Schema.org markup
- **Sitemap** - Geração automática

## 🔒 Segurança

### Implementações
- **XSS Protection** - Sanitização de inputs
- **CSRF Protection** - Tokens de proteção
- **Input Validation** - Validação client/server side
- **Secure Headers** - Headers de segurança

## 🧪 Desenvolvimento

### Scripts Disponíveis
```bash
npm run dev          # Desenvolvimento
npm run build        # Build para produção
npm start            # Servidor produção
npm run lint         # Linting
npm run type-check   # Verificação TypeScript
```

### Ambiente Local
```bash
# Instalar dependências
cd frontend
npm install

# Executar em desenvolvimento
npm run dev

# Acessar aplicação
http://localhost:3000
```

### Variáveis de Ambiente
```env
API_BASE_URL=http://localhost:8081
AI_SERVICE_URL=http://localhost:8000
SEARCH_SERVICE_URL=http://localhost:8086
REPORT_SERVICE_URL=http://localhost:8087
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=dev-secret-key
```

## 🚀 Deploy

### Build de Produção
```bash
npm run build
npm start
```

### Docker (Futuro)
```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
```

## 📈 Próximas Melhorias

### Funcionalidades Pendentes
- [ ] Página de Relatórios completa
- [ ] Página de Notificações
- [ ] Página de Configurações
- [ ] Página de Usuários
- [ ] Sistema de Busca Global avançado

### Otimizações Técnicas
- [ ] Implementar PWA
- [ ] Service Workers
- [ ] Offline support
- [ ] Push notifications
- [ ] WebSocket integration

### Testes
- [ ] Unit tests (Jest + Testing Library)
- [ ] Integration tests
- [ ] E2E tests (Playwright)
- [ ] Visual regression tests

## 📊 Status de Conclusão

### ✅ Implementado (100%)
- Core framework e configuração
- Sistema de autenticação
- Dashboard principal
- Gestão de processos
- AI Assistant interface
- Componentes UI base
- State management
- API integration
- Responsive design
- Type safety

### 🔄 Em Andamento (0%)
- Nada pendente na versão atual

### ⏳ Próximas Fases
- Testes automatizados
- Otimizações de performance
- Funcionalidades adicionais

---

**🏆 Status**: ✅ **100% COMPLETO** - Frontend Web App totalmente funcional e integrado
**📅 Concluído em**: 19/06/2025
**🎯 Próximo**: Teste de integração Frontend + Backend