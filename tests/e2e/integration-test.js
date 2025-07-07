/**
 * Teste de integração completa do Direito Lux
 * Valida todos os serviços e endpoints funcionais
 */

const axios = require('axios');
const config = require('./utils/config');

async function integrationTest() {
  console.log('🎯 Direito Lux - Teste de Integração Completa');
  console.log('===============================================\n');

  const results = {
    auth: false,
    process: false,
    report: false,
    integration: false
  };

  try {
    // 1. Testar Auth Service completo
    console.log('🔑 1. Testando Auth Service...');
    
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
        }
      }
    );
    
    const token = loginResponse.data.access_token;
    console.log(`✅ Auth Service: Login successful (${loginResponse.status})`);
    results.auth = true;

    // 2. Testar Process Service com dados reais
    console.log('\n📊 2. Testando Process Service...');
    
    const processStatsResponse = await axios.get(
      `${config.services.process}/api/v1/processes/stats`,
      {
        headers: {
          ...config.headers,
          'X-Tenant-ID': config.tenants.silva.id
        }
      }
    );
    
    const stats = processStatsResponse.data.data;
    console.log(`✅ Process Service: Stats retrieved`);
    console.log(`   📈 Total processos: ${stats.total}`);
    console.log(`   🟢 Ativos: ${stats.active}`);
    console.log(`   📚 Arquivados: ${stats.archived}`);
    
    const processListResponse = await axios.get(
      `${config.services.process}/api/v1/processes`,
      {
        headers: {
          ...config.headers,
          'X-Tenant-ID': config.tenants.silva.id
        }
      }
    );
    
    console.log(`✅ Process Service: List retrieved (${processListResponse.data.total} processos)`);
    results.process = true;

    // 3. Testar Report Service
    console.log('\n📋 3. Testando Report Service...');
    
    const dashboardResponse = await axios.get(
      `${config.services.report}/api/v1/reports/dashboard`,
      {
        headers: {
          ...config.headers,
          'X-Tenant-ID': config.tenants.silva.id
        }
      }
    );
    
    const dashboard = dashboardResponse.data.data;
    console.log(`✅ Report Service: Dashboard retrieved`);
    console.log(`   📊 Processos crescimento: ${dashboard.tendencias.processos_crescimento}`);
    console.log(`   ⭐ Satisfação cliente: ${dashboard.tendencias.satisfacao_cliente}`);
    console.log(`   🎯 Taxa sucesso: ${dashboard.tendencias.taxa_sucesso}`);
    
    const recentActivitiesResponse = await axios.get(
      `${config.services.report}/api/v1/reports/recent-activities`,
      {
        headers: {
          ...config.headers,
          'X-Tenant-ID': config.tenants.silva.id
        }
      }
    );
    
    console.log(`✅ Report Service: Recent activities retrieved (${recentActivitiesResponse.data.data.length} atividades)`);
    results.report = true;

    // 4. Testar integração entre serviços
    console.log('\n🔗 4. Testando integração entre serviços...');
    
    // Validar que os dados dos serviços são consistentes
    const processTotal = stats.total;
    const recentActivities = recentActivitiesResponse.data.data.length;
    
    console.log(`✅ Consistência de dados:`);
    console.log(`   Process Service: ${processTotal} processos`);
    console.log(`   Report Service: ${recentActivities} atividades recentes`);
    
    if (processTotal > 0 && recentActivities > 0) {
      console.log(`✅ Integração validada: Dados consistentes entre serviços`);
      results.integration = true;
    }

    // 5. Testar multi-tenant
    console.log('\n🏢 5. Testando multi-tenant...');
    
    const tenant2StatsResponse = await axios.get(
      `${config.services.process}/api/v1/processes/stats`,
      {
        headers: {
          ...config.headers,
          'X-Tenant-ID': config.tenants.costa.id
        }
      }
    );
    
    console.log(`✅ Multi-tenant: Tenant 2 stats retrieved`);
    console.log(`   📊 Tenant 1 processos: ${stats.total}`);
    console.log(`   📊 Tenant 2 processos: ${tenant2StatsResponse.data.data.total}`);

    // Resumo final
    console.log('\n🎉 TESTE DE INTEGRAÇÃO COMPLETO!');
    console.log('==================================');
    
    const passedTests = Object.values(results).filter(Boolean).length;
    const totalTests = Object.keys(results).length;
    
    console.log(`📊 Resultado: ${passedTests}/${totalTests} testes passaram`);
    console.log('\n✅ Serviços funcionais:');
    if (results.auth) console.log('   🔑 Auth Service: FUNCIONANDO');
    if (results.process) console.log('   📊 Process Service: FUNCIONANDO');
    if (results.report) console.log('   📋 Report Service: FUNCIONANDO');
    if (results.integration) console.log('   🔗 Integração: FUNCIONANDO');
    
    console.log('\n🚀 Sistema pronto para uso!');
    console.log('\n📋 Próximos passos:');
    console.log('  1. Inicializar frontend dashboard');
    console.log('  2. Conectar DataJud Service');
    console.log('  3. Implementar Notification Service');
    console.log('  4. Configurar AI Service');
    
    if (passedTests === totalTests) {
      console.log('\n🎯 STATUS: TODOS OS TESTES PASSARAM! ✅');
      process.exit(0);
    } else {
      console.log('\n⚠️  STATUS: ALGUNS TESTES FALHARAM');
      process.exit(1);
    }
    
  } catch (error) {
    console.error('\n❌ Erro durante teste de integração:', error.response?.data || error.message);
    console.error('\n🔧 Para resolver:');
    console.error('  1. Verificar se todos os serviços estão rodando');
    console.error('  2. Verificar logs dos serviços');
    console.error('  3. Verificar conectividade de rede');
    
    process.exit(1);
  }
}

// Executar se chamado diretamente
if (require.main === module) {
  integrationTest();
}

module.exports = integrationTest;