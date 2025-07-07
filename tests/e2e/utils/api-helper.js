/**
 * Helper para chamadas API nos testes E2E
 */

const axios = require('axios');
const config = require('./config');

class ApiHelper {
  constructor() {
    this.tokens = new Map(); // Cache de tokens por tenant
  }

  /**
   * Fazer login e obter token JWT
   */
  async login(tenantKey) {
    const tenant = config.tenants[tenantKey];
    if (!tenant) {
      throw new Error(`Tenant ${tenantKey} não encontrado`);
    }

    console.log(`🔑 Fazendo login para tenant: ${tenant.name}`);
    
    try {
      const response = await axios.post(
        `${config.services.auth}/api/v1/auth/login`,
        {
          email: tenant.email,
          password: tenant.password
        },
        {
          headers: {
            ...config.headers,
            'X-Tenant-ID': tenant.id
          },
          timeout: config.timeouts.auth
        }
      );

      const token = response.data.token;
      this.tokens.set(tenantKey, token);
      
      console.log(`✅ Login successful para ${tenant.name}`);
      return token;
    } catch (error) {
      console.error(`❌ Login failed para ${tenant.name}:`, error.response?.data || error.message);
      throw error;
    }
  }

  /**
   * Obter headers com autenticação
   */
  getAuthHeaders(tenantKey) {
    const tenant = config.tenants[tenantKey];
    const token = this.tokens.get(tenantKey);
    
    if (!token) {
      throw new Error(`Token não encontrado para tenant ${tenantKey}. Faça login primeiro.`);
    }

    return {
      ...config.headers,
      'Authorization': `Bearer ${token}`,
      'X-Tenant-ID': tenant.id
    };
  }

  /**
   * Chamada GET autenticada
   */
  async get(service, endpoint, tenantKey, options = {}) {
    const url = `${config.services[service]}${endpoint}`;
    const headers = this.getAuthHeaders(tenantKey);
    
    console.log(`📤 GET ${url}`);
    
    try {
      const response = await axios.get(url, {
        headers,
        timeout: config.timeouts.api,
        ...options
      });
      
      console.log(`✅ GET ${url} - Status: ${response.status}`);
      return response;
    } catch (error) {
      console.error(`❌ GET ${url} - Error:`, error.response?.status, error.response?.data);
      throw error;
    }
  }

  /**
   * Chamada POST autenticada
   */
  async post(service, endpoint, data, tenantKey, options = {}) {
    const url = `${config.services[service]}${endpoint}`;
    const headers = this.getAuthHeaders(tenantKey);
    
    console.log(`📤 POST ${url}`, { data });
    
    try {
      const response = await axios.post(url, data, {
        headers,
        timeout: config.timeouts.api,
        ...options
      });
      
      console.log(`✅ POST ${url} - Status: ${response.status}`);
      return response;
    } catch (error) {
      console.error(`❌ POST ${url} - Error:`, error.response?.status, error.response?.data);
      throw error;
    }
  }

  /**
   * Chamada PUT autenticada
   */
  async put(service, endpoint, data, tenantKey, options = {}) {
    const url = `${config.services[service]}${endpoint}`;
    const headers = this.getAuthHeaders(tenantKey);
    
    console.log(`📤 PUT ${url}`, { data });
    
    try {
      const response = await axios.put(url, data, {
        headers,
        timeout: config.timeouts.api,
        ...options
      });
      
      console.log(`✅ PUT ${url} - Status: ${response.status}`);
      return response;
    } catch (error) {
      console.error(`❌ PUT ${url} - Error:`, error.response?.status, error.response?.data);
      throw error;
    }
  }

  /**
   * Chamada DELETE autenticada
   */
  async delete(service, endpoint, tenantKey, options = {}) {
    const url = `${config.services[service]}${endpoint}`;
    const headers = this.getAuthHeaders(tenantKey);
    
    console.log(`📤 DELETE ${url}`);
    
    try {
      const response = await axios.delete(url, {
        headers,
        timeout: config.timeouts.api,
        ...options
      });
      
      console.log(`✅ DELETE ${url} - Status: ${response.status}`);
      return response;
    } catch (error) {
      console.error(`❌ DELETE ${url} - Error:`, error.response?.status, error.response?.data);
      throw error;
    }
  }

  /**
   * Verificar health de um serviço
   */
  async checkHealth(service) {
    const url = `${config.services[service]}/health`;
    
    try {
      const response = await axios.get(url, {
        timeout: 5000
      });
      
      console.log(`✅ Health check ${service}: OK`);
      return response.data;
    } catch (error) {
      console.error(`❌ Health check ${service}: FAIL`);
      throw error;
    }
  }

  /**
   * Aguardar até que um serviço esteja disponível
   */
  async waitForService(service, maxAttempts = 10, delay = 2000) {
    console.log(`⏳ Aguardando ${service} ficar disponível...`);
    
    for (let attempt = 1; attempt <= maxAttempts; attempt++) {
      try {
        await this.checkHealth(service);
        console.log(`✅ ${service} está disponível`);
        return true;
      } catch (error) {
        console.log(`⏳ Tentativa ${attempt}/${maxAttempts} para ${service}`);
        
        if (attempt === maxAttempts) {
          throw new Error(`Serviço ${service} não ficou disponível após ${maxAttempts} tentativas`);
        }
        
        await new Promise(resolve => setTimeout(resolve, delay));
      }
    }
  }

  /**
   * Limpar tokens (logout)
   */
  clearTokens() {
    this.tokens.clear();
    console.log('🧹 Tokens limpos');
  }
}

module.exports = new ApiHelper();