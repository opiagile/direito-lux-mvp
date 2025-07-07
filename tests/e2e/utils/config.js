/**
 * Configurações para testes E2E do Direito Lux
 */

const config = {
  // URLs dos serviços
  services: {
    auth: 'http://localhost:8081',
    tenant: 'http://localhost:8082', 
    process: 'http://localhost:8080',
    datajud: 'http://localhost:8084',
    notification: 'http://localhost:8085',
    ai: 'http://localhost:8000',
    search: 'http://localhost:8086',
    report: 'http://localhost:8087',
    mcp: 'http://localhost:8088',
    frontend: 'http://localhost:3000'
  },
  
  // Tenants de teste (conforme STATUS_IMPLEMENTACAO.md)
  tenants: {
    silva: {
      id: '11111111-1111-1111-1111-111111111111',
      name: 'Silva & Associados',
      plan: 'starter',
      email: 'admin@silvaassociados.com.br',
      password: 'password'
    },
    costa: {
      id: '22222222-2222-2222-2222-222222222222', 
      name: 'Costa & Santos',
      plan: 'professional',
      email: 'admin@costasantos.com.br',
      password: 'password'
    },
    machado: {
      id: '33333333-3333-3333-3333-333333333333',
      name: 'Machado Advogados', 
      plan: 'business',
      email: 'admin@machadoadvogados.com.br',
      password: 'password'
    },
    barros: {
      id: '44444444-4444-4444-4444-444444444444',
      name: 'Barros Enterprise',
      plan: 'enterprise', 
      email: 'admin@barrosent.com.br',
      password: 'password'
    }
  },
  
  // Headers padrão
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  },
  
  // Timeouts
  timeouts: {
    api: 10000,
    auth: 5000,
    notification: 15000
  },
  
  // Dados de teste
  testData: {
    validCNJ: '0000001-23.2024.8.26.0001',
    invalidCNJ: '1234567-89.2024.8.26.0001',
    testClient: {
      name: 'João da Silva',
      email: 'joao@teste.com',
      document: '123.456.789-00'
    }
  }
};

module.exports = config;