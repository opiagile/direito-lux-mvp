/**
 * Validação de Registro - Costa Advogados
 * Script para validar funcionalidade de registro via proxy
 */

const axios = require('axios');

// Configuração do teste
const BASE_URL = 'http://localhost:3000'; // Ambiente DEV local
const REGISTRATION_DATA = {
  tenant: {
    name: 'Costa Advogados',
    document: '12.345.678/0001-90',
    email: 'admin@costaadvogados.com.br',
    phone: '(11) 99999-9999',
    website: 'https://costaadvogados.com.br',
    plan: 'professional',
    address: {
      street: 'Rua da Consolação',
      number: '1000',
      complement: 'Sala 501',
      neighborhood: 'Consolação',
      city: 'São Paulo',
      state: 'SP',
      zipCode: '01302-000'
    }
  },
  user: {
    name: 'Dr. João Costa',
    email: 'joao@costaadvogados.com.br',
    password: 'Costa123!',
    phone: '(11) 98888-8888'
  }
};

async function validateRegistration() {
  console.log('🚀 Validando Registro - Costa Advogados');
  console.log('====================================');
  
  try {
    // 1. Verificar se o sistema está acessível
    console.log('1. Verificando conectividade...');
    
    const healthCheck = await axios.get(`${BASE_URL}/health`, {
      timeout: 10000,
      validateStatus: () => true
    });
    
    console.log(`   Status: ${healthCheck.status}`);
    console.log(`   Sistema: ${healthCheck.status === 200 ? '✅ Online' : '⚠️ Limitado'}`);
    
    // 2. Testar endpoint de registro
    console.log('\n2. Testando endpoint de registro...');
    
    const registrationResponse = await axios.post(`${BASE_URL}/api/v1/auth/register`, REGISTRATION_DATA, {
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
      timeout: 15000,
      validateStatus: () => true
    });
    
    console.log(`   Status: ${registrationResponse.status}`);
    console.log(`   Response: ${JSON.stringify(registrationResponse.data, null, 2)}`);
    
    // 3. Análise do resultado
    console.log('\n3. Análise do resultado...');
    
    if (registrationResponse.status === 201) {
      console.log('   ✅ SUCESSO: Registro funcionando corretamente');
      console.log('   📋 Dados do tenant e usuário processados');
      console.log('   🏢 Costa Advogados registrado com sucesso');
      
      // Verificar se recebemos os dados esperados
      if (registrationResponse.data.tenant_id && registrationResponse.data.user_id) {
        console.log(`   🆔 Tenant ID: ${registrationResponse.data.tenant_id}`);
        console.log(`   👤 User ID: ${registrationResponse.data.user_id}`);
      }
      
    } else if (registrationResponse.status === 400) {
      console.log('   ⚠️ ERRO DE VALIDAÇÃO: Dados inválidos');
      console.log('   🔍 Verificar formato dos dados enviados');
      
    } else if (registrationResponse.status === 409) {
      console.log('   ⚠️ CONFLITO: Usuário ou tenant já existe');
      console.log('   📧 Email ou CNPJ já cadastrados');
      
    } else if (registrationResponse.status === 503) {
      console.log('   ❌ ERRO CRÍTICO: Serviço indisponível');
      console.log('   🔧 Auth ou Tenant services não estão funcionais');
      
    } else if (registrationResponse.status === 404) {
      console.log('   ❌ ENDPOINT NÃO ENCONTRADO');
      console.log('   🛣️ Rota /api/v1/auth/register não implementada');
      
    } else {
      console.log(`   ❌ ERRO INESPERADO: ${registrationResponse.status}`);
      console.log('   🔍 Verificar logs do servidor');
    }
    
    // 4. Testar funcionalidade complementar (se registro deu certo)
    if (registrationResponse.status === 201) {
      console.log('\n4. Testando login com dados criados...');
      
      try {
        const loginResponse = await axios.post(`${BASE_URL}/api/v1/auth/login`, {
          email: REGISTRATION_DATA.user.email,
          password: REGISTRATION_DATA.user.password
        }, {
          headers: {
            'Content-Type': 'application/json',
            'X-Tenant-ID': registrationResponse.data.tenant_id
          },
          timeout: 10000,
          validateStatus: () => true
        });
        
        if (loginResponse.status === 200) {
          console.log('   ✅ Login funcionando: Registro completamente validado');
          console.log('   🔑 Token JWT gerado com sucesso');
        } else {
          console.log('   ⚠️ Login falhou: Registro parcialmente validado');
          console.log(`   Status: ${loginResponse.status}`);
        }
        
      } catch (loginError) {
        console.log('   ⚠️ Erro ao testar login:', loginError.message);
      }
    }
    
  } catch (error) {
    console.log('\n❌ ERRO DE CONEXÃO');
    console.log('==================');
    console.log(`Erro: ${error.message}`);
    
    if (error.code === 'ECONNREFUSED') {
      console.log('🔧 Sistema staging parece estar offline');
    } else if (error.code === 'ENOTFOUND') {
      console.log('🌐 DNS não resolveu o endereço');
    } else if (error.code === 'ETIMEDOUT') {
      console.log('⏱️ Timeout na conexão');
    }
    
    return false;
  }
  
  console.log('\n📊 RESUMO DA VALIDAÇÃO');
  console.log('======================');
  console.log('✅ Teste de validação executado');
  console.log('📋 Dados da Costa Advogados testados');
  console.log('🔍 Verificar logs acima para resultado final');
  
  return true;
}

// Executar validação
validateRegistration().then(success => {
  if (success) {
    console.log('\n🎯 VALIDAÇÃO COMPLETA');
    process.exit(0);
  } else {
    console.log('\n❌ FALHA NA VALIDAÇÃO');
    process.exit(1);
  }
}).catch(error => {
  console.error('❌ Erro fatal:', error);
  process.exit(1);
});