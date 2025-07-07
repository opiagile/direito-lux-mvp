-- Inserir dados de teste para CNPJ Providers
-- Usando tenant_id do Silva & Associados como exemplo

INSERT INTO cnpj_providers (
    tenant_id,
    cnpj, 
    name, 
    email,
    api_key, 
    daily_limit, 
    priority, 
    is_active
) VALUES 
(
    '11111111-1111-1111-1111-111111111111', -- Silva & Associados
    '11.222.333/0001-81',
    'Escritório Principal - Silva & Associados',
    'datajud@silvaassociados.com.br',
    'sk_test_4eC39HqLyjWDarjtT1zdp7dc',
    10000,
    1,
    true
),
(
    '11111111-1111-1111-1111-111111111111', -- Silva & Associados
    '22.333.444/0001-92',
    'Escritório Backup - Costa & Santos',
    'datajud@costasantos.com.br',
    'sk_test_5dD40IrMzxWEbkjtS2aeq8ed',
    10000,
    2,
    true
),
(
    '11111111-1111-1111-1111-111111111111', -- Silva & Associados
    '33.444.555/0001-03',
    'Escritório Reserva - Barros Entidades',
    'datajud@barrosent.com.br',
    'sk_test_6eE51JsNayXFclkuT3bfr9fe',
    10000,
    3,
    true
)
ON CONFLICT (cnpj) 
DO UPDATE SET
    name = EXCLUDED.name,
    email = EXCLUDED.email,
    api_key = EXCLUDED.api_key,
    daily_limit = EXCLUDED.daily_limit,
    priority = EXCLUDED.priority,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Atualizar uso diário para simular uso
UPDATE cnpj_providers 
SET 
    daily_usage = 4500,
    last_used_at = NOW() - INTERVAL '30 minutes'
WHERE cnpj = '11.222.333/0001-81';

UPDATE cnpj_providers 
SET 
    daily_usage = 2000,
    last_used_at = NOW() - INTERVAL '2 hours'
WHERE cnpj = '22.333.444/0001-92';

-- Verificar dados inseridos
SELECT 
    cnpj,
    name,
    email,
    api_key,
    daily_limit,
    daily_usage,
    priority,
    is_active,
    daily_limit - daily_usage as quota_disponivel,
    round((daily_usage::numeric / daily_limit * 100)::numeric, 2) as percentual_uso
FROM cnpj_providers
WHERE tenant_id = '11111111-1111-1111-1111-111111111111'
ORDER BY priority;