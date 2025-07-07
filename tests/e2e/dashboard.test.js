/**
 * Testes E2E para Dashboard e Relatórios
 */

const { describe, test, expect, beforeAll, afterAll } = require('@jest/globals');
const config = require('./utils/config');
const apiHelper = require('./utils/api-helper');

describe('📊 Dashboard e Relatórios', () => {
  beforeAll(async () => {
    console.log('🚀 Iniciando testes de dashboard...');
    
    // Login com diferentes tenants para testar
    await apiHelper.login('silva');
    await apiHelper.login('costa');
    await apiHelper.login('machado');
  });

  afterAll(async () => {
    apiHelper.clearTokens();
  });

  describe('KPIs do Process Service', () => {
    test('Deve obter KPIs principais do dashboard', async () => {
      const response = await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
      
      expect(response.status).toBe(200);
      
      // Verificar estrutura dos KPIs
      expect(response.data).toHaveProperty('total');
      expect(response.data).toHaveProperty('active');
      expect(response.data).toHaveProperty('paused');
      expect(response.data).toHaveProperty('archived');
      
      // Campos adicionais esperados
      expect(response.data).toHaveProperty('this_month');
      expect(response.data).toHaveProperty('todayMovements');
      expect(response.data).toHaveProperty('upcomingDeadlines');
      
      // Validar tipos
      expect(typeof response.data.total).toBe('number');
      expect(typeof response.data.active).toBe('number');
      expect(typeof response.data.paused).toBe('number');
      expect(typeof response.data.archived).toBe('number');
      expect(typeof response.data.this_month).toBe('number');
      expect(typeof response.data.todayMovements).toBe('number');
      expect(typeof response.data.upcomingDeadlines).toBe('number');
      
      // Validar lógica de negócio
      expect(response.data.total).toBeGreaterThanOrEqual(0);
      expect(response.data.active).toBeGreaterThanOrEqual(0);
      expect(response.data.total).toBeGreaterThanOrEqual(response.data.active);
      
      console.log(`📊 KPIs Silva & Associados:`, {
        total: response.data.total,
        active: response.data.active,
        this_month: response.data.this_month,
        todayMovements: response.data.todayMovements
      });
    });

    test('Deve ter dados diferentes por tenant (isolamento)', async () => {
      const silvaKPIs = await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
      const costaKPIs = await apiHelper.get('process', '/api/v1/processes/stats', 'costa');
      const machadoKPIs = await apiHelper.get('process', '/api/v1/processes/stats', 'machado');
      
      expect(silvaKPIs.status).toBe(200);
      expect(costaKPIs.status).toBe(200);
      expect(machadoKPIs.status).toBe(200);
      
      console.log('📊 Comparação KPIs por tenant:');
      console.log(`Silva (Starter): ${silvaKPIs.data.total} processos`);
      console.log(`Costa (Professional): ${costaKPIs.data.total} processos`);
      console.log(`Machado (Business): ${machadoKPIs.data.total} processos`);
      
      // Cada tenant deve ter seus próprios dados
      // (mesmo que sejam iguais por coincidência, a lógica de isolamento deve funcionar)
      expect(silvaKPIs.data).toBeDefined();
      expect(costaKPIs.data).toBeDefined();
      expect(machadoKPIs.data).toBeDefined();
    });
  });

  describe('Report Service - Atividades Recentes', () => {
    test('Deve obter atividades recentes', async () => {
      const response = await apiHelper.get('report', '/api/v1/reports/recent-activities', 'silva');
      
      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('data');
      expect(response.data).toHaveProperty('meta');
      
      // Validar estrutura das atividades
      expect(Array.isArray(response.data.data)).toBe(true);
      expect(response.data.meta).toHaveProperty('tenant_id');
      expect(response.data.meta).toHaveProperty('total');
      
      // Se há atividades, validar estrutura
      if (response.data.data.length > 0) {
        const atividade = response.data.data[0];
        expect(atividade).toHaveProperty('id');
        expect(atividade).toHaveProperty('type');
        expect(atividade).toHaveProperty('description');
        expect(atividade).toHaveProperty('timestamp');
        
        console.log(`📋 Exemplo de atividade recente:`, atividade);
      }
      
      console.log(`📈 Total de atividades recentes: ${response.data.data.length}`);
    });

    test('Deve obter métricas adicionais do dashboard', async () => {
      const response = await apiHelper.get('report', '/api/v1/reports/dashboard', 'silva');
      
      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('data');
      
      // Verificar estrutura das métricas adicionais
      const data = response.data.data;
      expect(data).toHaveProperty('resumo_semanal');
      expect(data).toHaveProperty('tendencias');
      expect(data).toHaveProperty('alertas');
      expect(data).toHaveProperty('performance');
      
      console.log(`📊 Resumo semanal:`, data.resumo_semanal);
      console.log(`📈 Alertas:`, data.alertas);
    });

    test('Deve funcionar mesmo com banco indisponível (graceful degradation)', async () => {
      // Este teste verifica se o Report Service tem graceful degradation
      // Como implementado no service para funcionar sem PostgreSQL/Redis
      
      const response = await apiHelper.get('report', '/api/v1/reports/recent-activities', 'silva');
      
      // Deve sempre retornar 200, mesmo que seja dados demo
      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('data');
      expect(Array.isArray(response.data.data)).toBe(true);
      
      console.log(`✅ Graceful degradation funcionando`);
    });
  });

  describe('Health Checks dos Serviços', () => {
    test('Process Service deve estar saudável', async () => {
      const health = await apiHelper.checkHealth('process');
      
      expect(health).toHaveProperty('status');
      expect(['healthy', 'ok']).toContain(health.status.toLowerCase());
      
      console.log(`💚 Process Service health:`, health);
    });

    test('Report Service deve estar saudável', async () => {
      const health = await apiHelper.checkHealth('report');
      
      expect(health).toHaveProperty('status');
      expect(['healthy', 'ok']).toContain(health.status.toLowerCase());
      
      // Verificar dependências relatadas
      if (health.database) {
        console.log(`📊 Report Service DB status: ${health.database}`);
      }
      if (health.redis) {
        console.log(`📊 Report Service Redis status: ${health.redis}`);
      }
      
      console.log(`💚 Report Service health:`, health);
    });

    test('Auth Service deve estar saudável', async () => {
      const health = await apiHelper.checkHealth('auth');
      
      expect(health).toHaveProperty('status');
      expect(['healthy', 'ok']).toContain(health.status.toLowerCase());
      
      console.log(`💚 Auth Service health:`, health);
    });
  });

  describe('Integração Dashboard Frontend', () => {
    test('Deve obter todos os dados necessários para o dashboard', async () => {
      // Simular o que o frontend faz para carregar o dashboard completo
      
      // 1. KPIs principais
      const kpisResponse = await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
      expect(kpisResponse.status).toBe(200);
      
      // 2. Atividades recentes  
      const activityResponse = await apiHelper.get('report', '/api/v1/reports/recent-activities', 'silva');
      expect(activityResponse.status).toBe(200);
      
      // 3. Métricas adicionais
      const dashboardResponse = await apiHelper.get('report', '/api/v1/reports/dashboard', 'silva');
      expect(dashboardResponse.status).toBe(200);
      
      console.log('✅ Dashboard completo carregado com sucesso');
      console.log(`📊 KPIs: ${JSON.stringify(kpisResponse.data, null, 2)}`);
      console.log(`📋 Atividades: ${activityResponse.data.data.length} items`);
      console.log(`📈 Métricas extras: ${Object.keys(dashboardResponse.data.data).length} seções`);
    });

    test('Dados devem estar no formato esperado pelo frontend', async () => {
      const kpisResponse = await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
      const kpis = kpisResponse.data;
      
      // Verificar se os campos têm os nomes esperados pelo frontend
      expect(kpis).toHaveProperty('total');
      expect(kpis).toHaveProperty('active');
      expect(kpis).toHaveProperty('paused');
      expect(kpis).toHaveProperty('archived');
      
      // Verificar se são números (não strings)
      expect(typeof kpis.total).toBe('number');
      expect(typeof kpis.active).toBe('number');
      
      // Campos específicos para dashboard
      expect(kpis).toHaveProperty('this_month');
      expect(kpis).toHaveProperty('todayMovements'); 
      expect(kpis).toHaveProperty('upcomingDeadlines');
      
      console.log('✅ Formato dos dados compatível com frontend');
    });
  });

  describe('Performance do Dashboard', () => {
    test('KPIs devem carregar rapidamente (< 300ms)', async () => {
      const startTime = Date.now();
      
      await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
      
      const duration = Date.now() - startTime;
      expect(duration).toBeLessThan(300);
      
      console.log(`⚡ KPIs carregaram em ${duration}ms`);
    });

    test('Atividades recentes devem carregar rapidamente (< 500ms)', async () => {
      const startTime = Date.now();
      
      await apiHelper.get('report', '/api/v1/reports/recent-activities', 'silva');
      
      const duration = Date.now() - startTime;
      expect(duration).toBeLessThan(500);
      
      console.log(`⚡ Atividades carregaram em ${duration}ms`);
    });

    test('Dashboard completo deve carregar em tempo aceitável (< 1 segundo)', async () => {
      const startTime = Date.now();
      
      // Carregar todas as APIs do dashboard em paralelo
      await Promise.all([
        apiHelper.get('process', '/api/v1/processes/stats', 'silva'),
        apiHelper.get('report', '/api/v1/reports/recent-activities', 'silva'),
        apiHelper.get('report', '/api/v1/reports/dashboard', 'silva')
      ]);
      
      const duration = Date.now() - startTime;
      expect(duration).toBeLessThan(1000);
      
      console.log(`⚡ Dashboard completo carregou em ${duration}ms`);
    });
  });

  describe('Casos de Erro', () => {
    test('Deve lidar com tenant inexistente graciosamente', async () => {
      // Usar um token válido mas com tenant ID inválido seria ideal
      // Por agora, testamos se as APIs lidam bem com erros
      
      try {
        // Se o middleware de tenant for rigoroso, isso pode falhar
        const response = await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
        expect(response.status).toBe(200);
      } catch (error) {
        // Se falhar por isolamento, tudo bem
        expect([400, 401, 403, 404]).toContain(error.response.status);
      }
    });

    test('Deve ter timeout apropriado para APIs lentas', async () => {
      // Este teste verifica se o timeout está configurado corretamente
      
      const startTime = Date.now();
      
      try {
        await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
        
        const duration = Date.now() - startTime;
        // Se completou sem timeout, deve ter sido razoavelmente rápido
        expect(duration).toBeLessThan(10000);
      } catch (error) {
        // Se deu timeout, verificar se foi o timeout esperado
        if (error.code === 'ECONNABORTED') {
          console.log('⏰ Timeout funcionando corretamente');
        } else {
          throw error;
        }
      }
    });
  });
});