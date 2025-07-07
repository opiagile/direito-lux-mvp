/**
 * Setup para testes E2E
 */

const config = require('./config');
const apiHelper = require('./api-helper');

async function setupTests() {
  console.log('🚀 Iniciando setup dos testes E2E...');
  
  try {
    // 1. Verificar se todos os serviços estão disponíveis
    console.log('\n📡 Verificando disponibilidade dos serviços...');
    
    const services = Object.keys(config.services);
    const healthChecks = [];
    
    for (const service of services) {
      try {
        await apiHelper.waitForService(service, 3, 1000);
        healthChecks.push({ service, status: 'OK' });
      } catch (error) {
        healthChecks.push({ service, status: 'FAIL', error: error.message });
        console.warn(`⚠️  ${service} não está disponível: ${error.message}`);
      }
    }
    
    // Mostrar resultado dos health checks
    console.log('\n📊 Status dos serviços:');
    healthChecks.forEach(({ service, status, error }) => {
      const icon = status === 'OK' ? '✅' : '❌';
      console.log(`${icon} ${service}: ${status}${error ? ` (${error})` : ''}`);
    });
    
    // Verificar se serviços críticos estão funcionando
    const criticalServices = ['auth', 'process', 'tenant'];
    const failedCritical = healthChecks.filter(hc => 
      criticalServices.includes(hc.service) && hc.status === 'FAIL'
    );
    
    if (failedCritical.length > 0) {
      console.error('\n❌ Serviços críticos não estão disponíveis:', failedCritical.map(s => s.service));
      console.error('💡 Execute: ./scripts/deploy-dev.sh start');
      throw new Error('Serviços críticos indisponíveis');
    }
    
    // 2. Fazer login inicial para validar autenticação
    console.log('\n🔑 Testando autenticação...');
    
    try {
      await apiHelper.login('silva');
      console.log('✅ Autenticação funcionando');
    } catch (error) {
      console.error('❌ Falha na autenticação:', error.message);
      throw error;
    }
    
    console.log('\n✅ Setup concluído com sucesso!');
    console.log('🎯 Iniciando execução dos testes...\n');
    
  } catch (error) {
    console.error('\n💥 Falha no setup:', error.message);
    console.error('\n📋 Passos para resolver:');
    console.error('1. Verificar se os serviços estão rodando: ./scripts/deploy-dev.sh status');
    console.error('2. Iniciar serviços se necessário: ./scripts/deploy-dev.sh start');
    console.error('3. Verificar logs: ./scripts/deploy-dev.sh logs');
    
    process.exit(1);
  }
}

// Executar setup se chamado diretamente
if (require.main === module) {
  setupTests();
}

module.exports = setupTests;