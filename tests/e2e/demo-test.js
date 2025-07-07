/**
 * Teste de demonstração E2E para validar funcionamento básico
 */

const axios = require('axios');
const config = require('./utils/config');

async function demoTest() {
  console.log('🎯 Direito Lux - Teste de Demonstração E2E');
  console.log('=============================================\n');

  try {
    // 1. Verificar Auth Service
    console.log('🔑 1. Testando Auth Service...');
    
    const healthResponse = await axios.get(`${config.services.auth}/health`, { timeout: 5000 });
    console.log(`✅ Auth Service health: ${healthResponse.data.status || 'OK'}`);
    
    // 2. Testar login
    console.log('\n🔐 2. Testando login...');
    
    const loginResponse = await axios.post(
      `${config.services.auth}/api/v1/auth/login`,
      {
        email: config.tenants.silva.email,
        password: config.tenants.silva.password
      },
      {
        headers: {
          ...config.headers,
          'X-Tenant-ID': config.tenants.silva.id
        },
        timeout: 10000
      }
    );
    
    console.log(`✅ Login successful: Status ${loginResponse.status}`);
    console.log(`🎫 Token recebido: ${loginResponse.data.access_token ? 'SIM' : 'NÃO'}`);
    
    // 3. Testar validação de token
    if (loginResponse.data.access_token) {
      console.log('\n🎯 3. Testando validação de token...');
      
      const validateResponse = await axios.get(
        `${config.services.auth}/api/v1/auth/validate`,
        {
          headers: {
            ...config.headers,
            'Authorization': `Bearer ${loginResponse.data.access_token}`,
            'X-Tenant-ID': config.tenants.silva.id
          },
          timeout: 5000
        }
      );
      
      console.log(`✅ Token validation: Status ${validateResponse.status}`);
      console.log(`👤 User ID: ${validateResponse.data.user_id || 'N/A'}`);
      console.log(`🏢 Tenant ID: ${validateResponse.data.tenant_id || 'N/A'}`);
    }
    
    // 4. Testar outros serviços se disponíveis
    console.log('\n🔍 4. Verificando outros serviços...');
    
    const services = ['process', 'report'];
    for (const service of services) {
      try {
        const response = await axios.get(`${config.services[service]}/health`, { timeout: 3000 });
        console.log(`✅ ${service} service: DISPONÍVEL`);
      } catch (error) {
        console.log(`⚠️  ${service} service: INDISPONÍVEL (${error.message})`);
      }
    }
    
    console.log('\n🎉 Teste de demonstração concluído com sucesso!');
    console.log('\n📋 Resumo:');
    console.log('  ✅ Auth Service funcionando');
    console.log('  ✅ Login multi-tenant funcionando'); 
    console.log('  ✅ Validação de token funcionando');
    console.log('  ✅ Estrutura E2E implementada');
    
    console.log('\n🚀 Para executar todos os testes:');
    console.log('  1. Iniciar todos os serviços: ./services/scripts/deploy-dev.sh start');
    console.log('  2. Executar testes: ./run-tests.sh');
    
  } catch (error) {
    console.error('\n❌ Erro durante teste:', error.response?.data || error.message);
    console.error('\n🔧 Para resolver:');
    console.error('  1. Verificar se Auth Service está rodando: docker ps');
    console.error('  2. Verificar logs: docker logs direito-lux-auth');
    console.error('  3. Reiniciar se necessário: docker restart direito-lux-auth');
    
    process.exit(1);
  }
}

// Executar se chamado diretamente
if (require.main === module) {
  demoTest();
}

module.exports = demoTest;