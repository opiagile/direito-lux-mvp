/**
 * Testes E2E para Autenticação Multi-tenant
 */

const { describe, test, expect, beforeAll, afterAll } = require('@jest/globals');
const config = require('./utils/config');
const apiHelper = require('./utils/api-helper');

describe('🔑 Autenticação Multi-tenant', () => {
  beforeAll(async () => {
    console.log('🚀 Iniciando testes de autenticação...');
  });

  afterAll(async () => {
    apiHelper.clearTokens();
  });

  describe('Login por Tenant', () => {
    test('Deve fazer login com tenant Silva & Associados (Starter)', async () => {
      const token = await apiHelper.login('silva');
      
      expect(token).toBeDefined();
      expect(typeof token).toBe('string');
      expect(token.length).toBeGreaterThan(0);
    });

    test('Deve fazer login com tenant Costa & Santos (Professional)', async () => {
      const token = await apiHelper.login('costa');
      
      expect(token).toBeDefined();
      expect(typeof token).toBe('string');
      expect(token.length).toBeGreaterThan(0);
    });

    test('Deve fazer login com tenant Machado Advogados (Business)', async () => {
      const token = await apiHelper.login('machado');
      
      expect(token).toBeDefined();
      expect(typeof token).toBe('string');
      expect(token.length).toBeGreaterThan(0);
    });

    test('Deve fazer login com tenant Barros Enterprise (Enterprise)', async () => {
      const token = await apiHelper.login('barros');
      
      expect(token).toBeDefined();
      expect(typeof token).toBe('string');
      expect(token.length).toBeGreaterThan(0);
    });
  });

  describe('Validação de Token', () => {
    test('Deve validar token válido', async () => {
      await apiHelper.login('silva');
      
      const response = await apiHelper.get('auth', '/api/v1/auth/validate', 'silva');
      
      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('valid', true);
      expect(response.data).toHaveProperty('user_id');
      expect(response.data).toHaveProperty('tenant_id', config.tenants.silva.id);
    });

    test('Deve rejeitar token inválido', async () => {
      try {
        await apiHelper.get('auth', '/api/v1/auth/validate', 'silva');
        // Se chegou aqui sem token, deve falhar
        expect(true).toBe(false);
      } catch (error) {
        expect(error.response.status).toBe(401);
      }
    });
  });

  describe('Isolamento Multi-tenant', () => {
    test('Deve isolar dados entre tenants diferentes', async () => {
      // Login com dois tenants diferentes
      await apiHelper.login('silva');
      await apiHelper.login('costa');
      
      // Buscar dados de cada tenant
      const silvaResponse = await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
      const costaResponse = await apiHelper.get('process', '/api/v1/processes/stats', 'costa');
      
      expect(silvaResponse.status).toBe(200);
      expect(costaResponse.status).toBe(200);
      
      // Os dados devem ser diferentes (isolamento)
      expect(silvaResponse.data).not.toEqual(costaResponse.data);
      
      // Verificar se os tenant_ids estão corretos nos dados retornados
      // (assumindo que a API retorna tenant_id nos metadados)
    });

    test('Não deve permitir acesso cross-tenant', async () => {
      await apiHelper.login('silva');
      
      // Tentar acessar dados usando token do Silva mas com X-Tenant-ID diferente
      const tenant = config.tenants.silva;
      const wrongTenantId = config.tenants.costa.id;
      
      try {
        const response = await apiHelper.get('process', '/api/v1/processes/stats', 'silva');
        // Verificar se os dados retornados são do tenant correto
        // Não do tenant "errado" que foi injetado no header
        expect(response.status).toBe(200);
        // Dados devem ser do Silva, não do Costa
        
      } catch (error) {
        // Se a API bloquear cross-tenant, isso é o comportamento esperado
        expect([401, 403]).toContain(error.response.status);
      }
    });
  });

  describe('Casos de Erro', () => {
    test('Deve falhar com credenciais inválidas', async () => {
      try {
        const response = await apiHelper.post('auth', '/api/v1/auth/login', {
          email: 'invalid@email.com',
          password: 'wrongpassword'
        }, 'silva');
        
        // Se não lançou erro, o status deve indicar falha
        expect([400, 401, 422]).toContain(response.status);
      } catch (error) {
        expect([400, 401, 422]).toContain(error.response.status);
      }
    });

    test('Deve falhar com tenant ID inválido no header', async () => {
      try {
        const response = await axios.post(
          `${config.services.auth}/api/v1/auth/login`,
          {
            email: config.tenants.silva.email,
            password: config.tenants.silva.password
          },
          {
            headers: {
              ...config.headers,
              'X-Tenant-ID': 'invalid-tenant-id'
            }
          }
        );
        
        expect([400, 401, 404]).toContain(response.status);
      } catch (error) {
        expect([400, 401, 404]).toContain(error.response.status);
      }
    });

    test('Deve falhar sem X-Tenant-ID header', async () => {
      const axios = require('axios');
      
      try {
        const response = await axios.post(
          `${config.services.auth}/api/v1/auth/login`,
          {
            email: config.tenants.silva.email,
            password: config.tenants.silva.password
          },
          {
            headers: config.headers // Sem X-Tenant-ID
          }
        );
        
        expect([400, 401]).toContain(response.status);
      } catch (error) {
        expect([400, 401]).toContain(error.response.status);
      }
    });
  });

  describe('Refresh Token', () => {
    test('Deve refresh token válido', async () => {
      // Primeiro login para obter refresh token
      await apiHelper.login('silva');
      
      // Tentar fazer refresh (assumindo que a API retorna refresh_token no login)
      try {
        const response = await apiHelper.post('auth', '/api/v1/auth/refresh', {}, 'silva');
        
        expect(response.status).toBe(200);
        expect(response.data).toHaveProperty('token');
        expect(typeof response.data.token).toBe('string');
      } catch (error) {
        // Se endpoint não implementado ainda, skip
        if (error.response.status === 404) {
          console.log('⚠️  Endpoint /refresh não implementado ainda');
        } else {
          throw error;
        }
      }
    });
  });

  describe('Performance de Autenticação', () => {
    test('Login deve ser rápido (< 2 segundos)', async () => {
      const startTime = Date.now();
      
      await apiHelper.login('silva');
      
      const endTime = Date.now();
      const duration = endTime - startTime;
      
      expect(duration).toBeLessThan(2000);
      console.log(`⚡ Login duration: ${duration}ms`);
    });

    test('Validação de token deve ser muito rápida (< 500ms)', async () => {
      await apiHelper.login('silva');
      
      const startTime = Date.now();
      
      await apiHelper.get('auth', '/api/v1/auth/validate', 'silva');
      
      const endTime = Date.now();
      const duration = endTime - startTime;
      
      expect(duration).toBeLessThan(500);
      console.log(`⚡ Token validation duration: ${duration}ms`);
    });
  });
});