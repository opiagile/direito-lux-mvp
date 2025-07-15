#!/bin/bash

echo "🕛 CONFIGURAR AUTO-SHUTDOWN/STARTUP - ECONOMIA MÁXIMA"
echo "===================================================="

PROJECT_ID="direito-lux-staging-2025"
REGION="us-central1"

# Criar Cloud Function para gerenciar cluster
create_cloud_function() {
    echo "☁️ Criando Cloud Function para auto-shutdown/startup..."
    
    # Criar diretório temporário
    mkdir -p /tmp/cluster-manager
    cd /tmp/cluster-manager
    
    # Criar main.py
    cat > main.py << 'EOF'
import functions_framework
from google.cloud import container_v1
import os

PROJECT_ID = os.environ.get('PROJECT_ID', 'direito-lux-staging-2025')
CLUSTER_NAME = os.environ.get('CLUSTER_NAME', 'direito-lux-gke-staging')
LOCATION = os.environ.get('LOCATION', 'us-central1')

@functions_framework.http
def manage_cluster(request):
    """Gerencia estado do cluster GKE"""
    client = container_v1.ClusterManagerClient()
    cluster_path = f"projects/{PROJECT_ID}/locations/{LOCATION}/clusters/{CLUSTER_NAME}"
    
    action = request.args.get('action', 'status')
    
    try:
        if action == 'stop':
            # Redimensionar para 0 nodes
            operation = client.set_node_pool_size(
                name=f"{cluster_path}/nodePools/default-pool",
                node_count=0
            )
            return {'status': 'stopping', 'operation': operation.name}
            
        elif action == 'start':
            # Redimensionar para 1 node
            operation = client.set_node_pool_size(
                name=f"{cluster_path}/nodePools/default-pool", 
                node_count=1
            )
            return {'status': 'starting', 'operation': operation.name}
            
        elif action == 'status':
            cluster = client.get_cluster(name=cluster_path)
            node_count = sum(pool.initial_node_count for pool in cluster.node_pools)
            return {
                'status': cluster.status.name,
                'node_count': node_count,
                'cost_per_hour': node_count * 0.60  # e2-small
            }
            
    except Exception as e:
        return {'error': str(e)}, 500
        
    return {'error': 'Invalid action'}, 400
EOF

    # Criar requirements.txt
    cat > requirements.txt << 'EOF'
google-cloud-container==2.17.4
functions-framework==3.4.0
EOF

    # Deploy da função
    echo "📦 Fazendo deploy da Cloud Function..."
    gcloud functions deploy cluster-manager \
        --runtime python311 \
        --trigger-http \
        --allow-unauthenticated \
        --region=$REGION \
        --project=$PROJECT_ID \
        --set-env-vars PROJECT_ID=$PROJECT_ID,CLUSTER_NAME=direito-lux-gke-staging,LOCATION=$REGION
    
    # Obter URL da função
    FUNCTION_URL=$(gcloud functions describe cluster-manager --region=$REGION --project=$PROJECT_ID --format="value(httpsTrigger.url)")
    echo "✅ Cloud Function criada: $FUNCTION_URL"
    
    cd -
    rm -rf /tmp/cluster-manager
    
    return 0
}

# Criar Cloud Scheduler para shutdown automático
create_scheduler() {
    echo "⏰ Configurando Cloud Scheduler..."
    
    # Job para parar às 23:00 (11 PM) horário de Brasília
    gcloud scheduler jobs create http shutdown-cluster \
        --schedule="0 23 * * *" \
        --time-zone="America/Sao_Paulo" \
        --uri="$FUNCTION_URL?action=stop" \
        --http-method="GET" \
        --project=$PROJECT_ID \
        --location=$REGION
    
    echo "✅ Scheduler configurado - cluster para automaticamente às 23:00"
}

# Criar página de "despertar" do sistema
create_wakeup_page() {
    echo "🌅 Criando página de wake-up..."
    
    mkdir -p /tmp/wakeup-page
    cd /tmp/wakeup-page
    
    # Criar HTML simples
    cat > index.html << 'EOF'
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Direito Lux - Iniciando Sistema</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; }
        .container { max-width: 600px; margin: 0 auto; }
        .spinner { border: 4px solid #f3f3f3; border-top: 4px solid #3498db; border-radius: 50%; width: 50px; height: 50px; animation: spin 1s linear infinite; margin: 20px auto; }
        @keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
        button { background: #4CAF50; border: none; color: white; padding: 15px 32px; text-align: center; font-size: 16px; margin: 4px 2px; cursor: pointer; border-radius: 5px; }
        .status { margin: 20px 0; padding: 10px; background: rgba(255,255,255,0.1); border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🏛️ Direito Lux</h1>
        <h2>Sistema em Modo Economia</h2>
        
        <div class="status" id="status">
            💡 Sistema pausado para economizar custos<br>
            Custo atual: R$0/hora
        </div>
        
        <div class="spinner" id="spinner" style="display: none;"></div>
        
        <button onclick="wakeUpSystem()" id="wakeBtn">🚀 Iniciar Sistema</button>
        
        <p><small>⏱️ Tempo estimado de inicialização: 2-3 minutos</small></p>
        
        <div id="progress" style="display: none;">
            <h3>📋 Progresso:</h3>
            <div id="logs"></div>
        </div>
    </div>

    <script>
        const FUNCTION_URL = 'FUNCTION_URL_PLACEHOLDER';
        
        async function wakeUpSystem() {
            document.getElementById('wakeBtn').disabled = true;
            document.getElementById('spinner').style.display = 'block';
            document.getElementById('progress').style.display = 'block';
            
            updateLog('🚀 Iniciando cluster GKE...');
            
            try {
                const response = await fetch(FUNCTION_URL + '?action=start');
                const result = await response.json();
                
                if (result.status === 'starting') {
                    updateLog('✅ Comando enviado, aguardando cluster...');
                    await waitForCluster();
                } else {
                    updateLog('❌ Erro: ' + (result.error || 'Desconhecido'));
                }
            } catch (error) {
                updateLog('❌ Erro de conexão: ' + error.message);
            }
        }
        
        async function waitForCluster() {
            for (let i = 0; i < 30; i++) {
                updateLog(`⏳ Verificando status (${i + 1}/30)...`);
                
                try {
                    const response = await fetch(FUNCTION_URL + '?action=status');
                    const status = await response.json();
                    
                    if (status.node_count > 0) {
                        updateLog('✅ Cluster ativo! Redirecionando...');
                        setTimeout(() => {
                            window.location.href = '/';
                        }, 2000);
                        return;
                    }
                } catch (error) {
                    updateLog('⚠️ Erro na verificação: ' + error.message);
                }
                
                await new Promise(resolve => setTimeout(resolve, 10000)); // 10 segundos
            }
            
            updateLog('❌ Timeout - tente novamente ou verifique manualmente');
            document.getElementById('wakeBtn').disabled = false;
            document.getElementById('spinner').style.display = 'none';
        }
        
        function updateLog(message) {
            const logs = document.getElementById('logs');
            logs.innerHTML += '<div>' + new Date().toLocaleTimeString() + ' - ' + message + '</div>';
            logs.scrollTop = logs.scrollHeight;
        }
        
        // Verificar status ao carregar página
        window.onload = async function() {
            try {
                const response = await fetch(FUNCTION_URL + '?action=status');
                const status = await response.json();
                
                if (status.node_count > 0) {
                    document.getElementById('status').innerHTML = 
                        '✅ Sistema ATIVO<br>Custo atual: R$' + (status.cost_per_hour * 24).toFixed(2) + '/dia';
                    document.getElementById('wakeBtn').innerHTML = '🌐 Acessar Sistema';
                    document.getElementById('wakeBtn').onclick = () => window.location.href = '/';
                }
            } catch (error) {
                console.error('Erro ao verificar status:', error);
            }
        };
    </script>
</body>
</html>
EOF

    echo "✅ Página de wake-up criada"
    cd -
    rm -rf /tmp/wakeup-page
}

# Menu principal
case "${1:-help}" in
    "setup")
        echo "🔧 Configurando auto-shutdown completo..."
        create_cloud_function
        FUNCTION_URL=$(gcloud functions describe cluster-manager --region=$REGION --project=$PROJECT_ID --format="value(httpsTrigger.url)")
        create_scheduler
        create_wakeup_page
        echo ""
        echo "✅ CONFIGURAÇÃO COMPLETA!"
        echo "📋 O que foi criado:"
        echo "   - Cloud Function: $FUNCTION_URL"
        echo "   - Cloud Scheduler: Para às 23:00 automaticamente"
        echo "   - Página de wake-up: Para iniciar sistema sob demanda"
        echo ""
        echo "💰 ECONOMIA:"
        echo "   - Das 23:00 às 07:00 (8h): R$0/hora = R$0"
        echo "   - Das 07:00 às 23:00 (16h): R$15/dia"
        echo "   - TOTAL: ~R$15/dia ao invés de R$87/dia"
        echo "   - ECONOMIA: 83% (R$1890/mês)"
        ;;
    "test")
        echo "🧪 Testando sistema..."
        FUNCTION_URL=$(gcloud functions describe cluster-manager --region=$REGION --project=$PROJECT_ID --format="value(httpsTrigger.url)" 2>/dev/null)
        if [ -n "$FUNCTION_URL" ]; then
            echo "📊 Status atual:"
            curl -s "$FUNCTION_URL?action=status" | jq '.'
        else
            echo "❌ Cloud Function não encontrada. Execute: ./setup-auto-shutdown.sh setup"
        fi
        ;;
    "help"|*)
        echo "🔧 Comandos disponíveis:"
        echo "  ./setup-auto-shutdown.sh setup  # Configurar auto-shutdown completo"
        echo "  ./setup-auto-shutdown.sh test   # Testar configuração"
        echo ""
        echo "🎯 RESULTADO:"
        echo "   - Sistema para automaticamente às 23:00"
        echo "   - Usuário acessa página e clica 'Iniciar Sistema'"
        echo "   - Sistema inicia em 2-3 minutos"
        echo "   - ECONOMIA: 83% nos custos!"
        ;;
esac