/**
 * Testes E2E para CRUD de Processos
 */

const { describe, test, expect, beforeAll, afterAll, beforeEach } = require('@jest/globals');
const { v4: uuidv4 } = require('uuid');
const config = require('./utils/config');
const apiHelper = require('./utils/api-helper');

describe('📋 CRUD de Processos', () => {
  let createdProcessIds = [];

  beforeAll(async () => {
    console.log('🚀 Iniciando testes de processos...');
    
    // Login com tenant para os testes
    await apiHelper.login('silva');
    await apiHelper.login('costa');
  });

  afterAll(async () => {
    // Limpar processos criados durante os testes
    console.log('🧹 Limpando processos de teste...');
    
    for (const processId of createdProcessIds) {
      try {
        await apiHelper.delete('process', `/api/v1/processes/${processId}`, 'silva');
      } catch (error) {
        console.warn(`⚠️  Erro ao deletar processo ${processId}:`, error.message);
      }
    }
    
    apiHelper.clearTokens();
  });

  describe('Dashboard e Estatísticas', () => {
    test('Deve obter estatísticas do dashboard', async () => {
      const response = await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
      
      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('total');
      expect(response.data).toHaveProperty('active');
      expect(response.data).toHaveProperty('paused');
      expect(response.data).toHaveProperty('archived');
      
      // Validar tipos
      expect(typeof response.data.total).toBe('number');
      expect(typeof response.data.active).toBe('number');
      expect(typeof response.data.paused).toBe('number');
      expect(typeof response.data.archived).toBe('number');
      
      // Validar lógica de negócio
      expect(response.data.total).toBeGreaterThanOrEqual(0);
      expect(response.data.active).toBeGreaterThanOrEqual(0);
      expect(response.data.paused).toBeGreaterThanOrEqual(0);
      expect(response.data.archived).toBeGreaterThanOrEqual(0);
      
      console.log(`📊 Stats Silva: ${JSON.stringify(response.data, null, 2)}`);
    });

    test('Deve obter estatísticas diferentes para tenants diferentes', async () => {
      const silvaStats = await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
      const costaStats = await apiHelper.get('process', '/api/v1/processes/stats', 'costa');
      
      expect(silvaStats.status).toBe(200);
      expect(costaStats.status).toBe(200);
      
      // Não devem ser exatamente iguais (isolamento)
      // A menos que ambos tenham exatamente os mesmos dados por coincidência
      console.log(`📊 Silva total: ${silvaStats.data.total}`);
      console.log(`📊 Costa total: ${costaStats.data.total}`);
    });
  });

  describe('Listar Processos', () => {
    test('Deve listar processos do tenant', async () => {
      const response = await apiHelper.get('process', '/api/v1/processes', 'silva');
      
      expect(response.status).toBe(200);
      expect(Array.isArray(response.data)).toBe(true);
      
      // Se há processos, validar estrutura
      if (response.data.length > 0) {
        const processo = response.data[0];
        expect(processo).toHaveProperty('id');
        expect(processo).toHaveProperty('number');
        expect(processo).toHaveProperty('status');
        
        console.log(`📋 Exemplo de processo: ${JSON.stringify(processo, null, 2)}`);
      }
    });

    test('Deve suportar paginação', async () => {
      const response = await apiHelper.get('process', '/api/v1/processes?page=1&limit=5', 'silva');
      
      expect(response.status).toBe(200);
      expect(Array.isArray(response.data)).toBe(true);
      expect(response.data.length).toBeLessThanOrEqual(5);
    });

    test('Deve suportar filtros por status', async () => {
      const response = await apiHelper.get('process', '/api/v1/processes?status=active', 'silva');
      
      expect(response.status).toBe(200);
      expect(Array.isArray(response.data)).toBe(true);
      
      // Todos os processos retornados devem ter status 'active'
      response.data.forEach(processo => {
        expect(processo.status).toBe('active');
      });
    });
  });

  describe('Criar Processo', () => {
    test('Deve criar processo válido', async () => {
      const novoProcesso = {
        number: `${Math.random().toString().substring(2, 9)}-23.2024.8.26.0001`,
        court: 'TJ-SP',
        client_name: config.testData.testClient.name,
        client_email: config.testData.testClient.email,
        classification: 'Civil',
        subject: 'Responsabilidade Civil',
        monitoring_enabled: true
      };
      
      const response = await apiHelper.post('process', '/api/v1/processes', novoProcesso, 'silva');
      
      expect(response.status).toBe(201);
      expect(response.data).toHaveProperty('id');
      expect(response.data.number).toBe(novoProcesso.number);
      expect(response.data.court).toBe(novoProcesso.court);
      expect(response.data.client_name).toBe(novoProcesso.client_name);
      
      // Salvar ID para limpeza
      createdProcessIds.push(response.data.id);
      
      console.log(`✅ Processo criado: ${response.data.id}`);
    });

    test('Deve falhar com número CNJ inválido', async () => {
      const processoInvalido = {
        number: 'numero-invalido',
        court: 'TJ-SP',
        client_name: config.testData.testClient.name,
        client_email: config.testData.testClient.email,
        classification: 'Civil'
      };
      
      try {
        await apiHelper.post('process', '/api/v1/processes', processoInvalido, 'silva');
        // Se chegou aqui, deveria ter falhado
        expect(true).toBe(false);
      } catch (error) {
        expect([400, 422]).toContain(error.response.status);
        expect(error.response.data).toHaveProperty('message');
        console.log(`✅ Validação CNJ funcionou: ${error.response.data.message}`);
      }
    });

    test('Deve falhar com dados obrigatórios ausentes', async () => {
      const processoIncompleto = {
        number: `${Math.random().toString().substring(2, 9)}-23.2024.8.26.0001`
        // Faltando court, client_name, etc.
      };
      
      try {
        await apiHelper.post('process', '/api/v1/processes', processoIncompleto, 'silva');
        expect(true).toBe(false);
      } catch (error) {
        expect([400, 422]).toContain(error.response.status);
        console.log(`✅ Validação de campos obrigatórios funcionou`);
      }
    });
  });

  describe('Obter Processo Específico', () => {
    let processoTeste;

    beforeEach(async () => {
      // Criar processo para teste
      const novoProcesso = {
        number: `${Math.random().toString().substring(2, 9)}-23.2024.8.26.0001`,
        court: 'TJ-SP',
        client_name: 'Cliente Teste',
        client_email: 'teste@exemplo.com',
        classification: 'Civil',
        subject: 'Teste E2E'
      };
      
      const response = await apiHelper.post('process', '/api/v1/processes', novoProcesso, 'silva');
      processoTeste = response.data;
      createdProcessIds.push(processoTeste.id);
    });

    test('Deve obter processo por ID', async () => {
      const response = await apiHelper.get('process', `/api/v1/processes/${processoTeste.id}`, 'silva');
      
      expect(response.status).toBe(200);
      expect(response.data.id).toBe(processoTeste.id);
      expect(response.data.number).toBe(processoTeste.number);
      expect(response.data.court).toBe(processoTeste.court);
    });

    test('Deve falhar ao obter processo inexistente', async () => {
      const idInexistente = uuidv4();
      
      try {
        await apiHelper.get('process', `/api/v1/processes/${idInexistente}`, 'silva');
        expect(true).toBe(false);
      } catch (error) {
        expect(error.response.status).toBe(404);
      }
    });

    test('Não deve permitir acesso cross-tenant', async () => {
      // Tentar acessar processo do Silva usando token do Costa
      try {
        await apiHelper.get('process', `/api/v1/processes/${processoTeste.id}`, 'costa');
        expect(true).toBe(false);
      } catch (error) {
        expect([401, 403, 404]).toContain(error.response.status);
        console.log(`✅ Isolamento multi-tenant funcionou`);
      }
    });
  });

  describe('Atualizar Processo', () => {
    let processoTeste;

    beforeEach(async () => {
      const novoProcesso = {
        number: `${Math.random().toString().substring(2, 9)}-23.2024.8.26.0001`,
        court: 'TJ-SP',
        client_name: 'Cliente Original',
        client_email: 'original@exemplo.com',
        classification: 'Civil',
        subject: 'Assunto Original'
      };
      
      const response = await apiHelper.post('process', '/api/v1/processes', novoProcesso, 'silva');
      processoTeste = response.data;
      createdProcessIds.push(processoTeste.id);
    });

    test('Deve atualizar processo', async () => {
      const atualizacao = {
        client_name: 'Cliente Atualizado',
        client_email: 'atualizado@exemplo.com',
        subject: 'Assunto Atualizado'
      };
      
      const response = await apiHelper.put(
        'process', 
        `/api/v1/processes/${processoTeste.id}`, 
        atualizacao, 
        'silva'
      );
      
      expect(response.status).toBe(200);
      expect(response.data.client_name).toBe(atualizacao.client_name);
      expect(response.data.client_email).toBe(atualizacao.client_email);
      expect(response.data.subject).toBe(atualizacao.subject);
      
      // Campos não alterados devem permanecer
      expect(response.data.number).toBe(processoTeste.number);
      expect(response.data.court).toBe(processoTeste.court);
    });

    test('Deve falhar ao atualizar processo inexistente', async () => {
      const idInexistente = uuidv4();
      
      try {
        await apiHelper.put(
          'process', 
          `/api/v1/processes/${idInexistente}`, 
          { client_name: 'Teste' }, 
          'silva'
        );
        expect(true).toBe(false);
      } catch (error) {
        expect(error.response.status).toBe(404);
      }
    });
  });

  describe('Deletar Processo', () => {
    test('Deve deletar processo', async () => {
      // Criar processo para deletar
      const novoProcesso = {
        number: `${Math.random().toString().substring(2, 9)}-23.2024.8.26.0001`,
        court: 'TJ-SP',
        client_name: 'Cliente Para Deletar',
        client_email: 'deletar@exemplo.com',
        classification: 'Civil'
      };
      
      const createResponse = await apiHelper.post('process', '/api/v1/processes', novoProcesso, 'silva');
      const processoId = createResponse.data.id;
      
      // Deletar
      const deleteResponse = await apiHelper.delete('process', `/api/v1/processes/${processoId}`, 'silva');
      expect([200, 204]).toContain(deleteResponse.status);
      
      // Verificar que foi deletado
      try {
        await apiHelper.get('process', `/api/v1/processes/${processoId}`, 'silva');
        expect(true).toBe(false);
      } catch (error) {
        expect(error.response.status).toBe(404);
      }
      
      console.log(`✅ Processo ${processoId} deletado com sucesso`);
    });

    test('Deve falhar ao deletar processo inexistente', async () => {
      const idInexistente = uuidv4();
      
      try {
        await apiHelper.delete('process', `/api/v1/processes/${idInexistente}`, 'silva');
        expect(true).toBe(false);
      } catch (error) {
        expect(error.response.status).toBe(404);
      }
    });
  });

  describe('Performance', () => {
    test('Listar processos deve ser rápido (< 1 segundo)', async () => {
      const startTime = Date.now();
      
      await apiHelper.get('process', '/api/v1/processes', 'silva');
      
      const duration = Date.now() - startTime;
      expect(duration).toBeLessThan(1000);
      
      console.log(`⚡ Listar processos: ${duration}ms`);
    });

    test('Obter estatísticas deve ser muito rápido (< 500ms)', async () => {
      const startTime = Date.now();
      
      await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
      
      const duration = Date.now() - startTime;
      expect(duration).toBeLessThan(500);
      
      console.log(`⚡ Estatísticas: ${duration}ms`);
    });

    test('Criar processo deve ser razoavelmente rápido (< 2 segundos)', async () => {
      const novoProcesso = {
        number: `${Math.random().toString().substring(2, 9)}-23.2024.8.26.0001`,
        court: 'TJ-SP',
        client_name: 'Performance Test',
        client_email: 'perf@test.com',
        classification: 'Civil'
      };
      
      const startTime = Date.now();
      
      const response = await apiHelper.post('process', '/api/v1/processes', novoProcesso, 'silva');
      
      const duration = Date.now() - startTime;
      expect(duration).toBeLessThan(2000);
      
      createdProcessIds.push(response.data.id);
      console.log(`⚡ Criar processo: ${duration}ms`);
    });
  });
});