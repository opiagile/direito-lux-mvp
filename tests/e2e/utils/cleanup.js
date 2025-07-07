/**
 * Cleanup após testes E2E
 */

const apiHelper = require('./api-helper');

async function cleanup() {
  console.log('\n🧹 Iniciando cleanup após testes...');
  
  try {
    // Limpar tokens de autenticação
    apiHelper.clearTokens();
    
    // Aqui podemos adicionar outras limpezas se necessário:
    // - Remover dados de teste criados
    // - Resetar estados
    // - Limpar caches
    
    console.log('✅ Cleanup concluído');
    
  } catch (error) {
    console.error('❌ Erro durante cleanup:', error.message);
    // Não falhar os testes por erro de cleanup
  }
}

// Executar cleanup se chamado diretamente
if (require.main === module) {
  cleanup();
}

module.exports = cleanup;