#!/bin/bash

echo "🔧 CONFIGURANDO ACESSO LOCAL AO STAGING"
echo "========================================"

# Verificar se já existe
if grep -q "staging.direitolux.com.br" /etc/hosts 2>/dev/null; then
    echo "⚠️  Entradas já existem em /etc/hosts"
    echo "🔍 Removendo entradas antigas..."
    sudo sed -i '' '/staging.direitolux.com.br/d' /etc/hosts 2>/dev/null || sudo sed -i '/staging.direitolux.com.br/d' /etc/hosts
    sudo sed -i '' '/api-staging.direitolux.com.br/d' /etc/hosts 2>/dev/null || sudo sed -i '/api-staging.direitolux.com.br/d' /etc/hosts
fi

echo ""
echo "📝 Adicionando entradas no /etc/hosts..."
echo ""
echo "35.188.198.87 staging.direitolux.com.br" | sudo tee -a /etc/hosts
echo "35.188.198.87 api-staging.direitolux.com.br" | sudo tee -a /etc/hosts

echo ""
echo "✅ Configuração concluída!"
echo ""
echo "🌐 AGORA VOCÊ PODE ACESSAR:"
echo "   https://staging.direitolux.com.br"
echo ""
echo "🔓 Credenciais de teste:"
echo "   Email: admin@silvaassociados.com.br"
echo "   Senha: password"
echo ""
echo "⚠️  IMPORTANTE: Isso é temporário até configurar o DNS real"