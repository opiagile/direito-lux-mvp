#!/bin/bash
echo "📋 Executando Migrations..." && \
cd services/tenant-service && \
echo "1/3 Tenant Service..." && \
make migrate-up && \
echo "✅ Tenant OK" && \
cd ../auth-service && \
echo "2/3 Auth Service..." && \
make migrate-up && \
echo "✅ Auth OK" && \
cd ../process-service && \
echo "3/3 Process Service..." && \
make migrate-up && \
echo "✅ Process OK" && \
cd ../.. && \
echo "🎯 Verificando dados..." && \
PGPASSWORD=dev_password_123 psql -h localhost -U direito_lux -d direito_lux_dev -c "
SELECT 
    'Tenants: ' || COUNT(*) FROM tenants
UNION ALL
SELECT 
    'Users: ' || COUNT(*) FROM users  
UNION ALL
SELECT 
    'Processes: ' || COUNT(*) FROM processes;
" && \
echo "🎉 Setup completo!"