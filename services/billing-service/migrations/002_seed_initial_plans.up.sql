-- Migration: Seed initial plans
-- Created: 2025-07-11

-- Insert default plans
INSERT INTO plans (
    id,
    name,
    display_name,
    description,
    price_monthly,
    price_yearly,
    trial_days,
    active,
    max_processes,
    max_users,
    max_ai_requests,
    max_bot_commands,
    has_whatsapp,
    has_telegram,
    has_mcp_bot,
    has_unlimited_search,
    has_ai_analysis,
    has_predictions,
    has_doc_generation,
    has_white_label,
    has_api_access,
    has_priority_support,
    created_at,
    updated_at
) VALUES 
(
    uuid_generate_v4(),
    'starter',
    'Starter',
    'Ideal para advogados autônomos',
    9900,  -- R$ 99,00
    99000, -- R$ 990,00 (2 meses grátis)
    15,    -- 15 dias de trial
    true,
    50,    -- 50 processos
    2,     -- 2 usuários
    10,    -- 10 IA requests/mês
    0,     -- 0 comandos bot
    true,  -- WhatsApp
    true,  -- Telegram
    false, -- MCP Bot
    true,  -- Busca ilimitada
    true,  -- IA Analysis
    false, -- Predictions
    false, -- Doc Generation
    false, -- White Label
    false, -- API Access
    false, -- Priority Support
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    uuid_generate_v4(),
    'professional',
    'Professional',
    'Para pequenos escritórios',
    29900, -- R$ 299,00
    299000, -- R$ 2.990,00 (2 meses grátis)
    15,    -- 15 dias de trial
    true,
    200,   -- 200 processos
    5,     -- 5 usuários
    50,    -- 50 IA requests/mês
    200,   -- 200 comandos bot/mês
    true,  -- WhatsApp
    true,  -- Telegram
    true,  -- MCP Bot
    true,  -- Busca ilimitada
    true,  -- IA Analysis
    true,  -- Predictions
    true,  -- Doc Generation
    false, -- White Label
    false, -- API Access
    false, -- Priority Support
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    uuid_generate_v4(),
    'business',
    'Business',
    'Para escritórios médios',
    69900, -- R$ 699,00
    699000, -- R$ 6.990,00 (2 meses grátis)
    15,    -- 15 dias de trial
    true,
    500,   -- 500 processos
    15,    -- 15 usuários
    200,   -- 200 IA requests/mês
    1000,  -- 1000 comandos bot/mês
    true,  -- WhatsApp
    true,  -- Telegram
    true,  -- MCP Bot
    true,  -- Busca ilimitada
    true,  -- IA Analysis
    true,  -- Predictions
    true,  -- Doc Generation
    false, -- White Label
    true,  -- API Access
    true,  -- Priority Support
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    uuid_generate_v4(),
    'enterprise',
    'Enterprise',
    'Para grandes escritórios',
    0,     -- Sob consulta
    0,     -- Sob consulta
    30,    -- 30 dias de trial
    true,
    -1,    -- Ilimitado
    -1,    -- Ilimitado
    -1,    -- Ilimitado
    -1,    -- Ilimitado
    true,  -- WhatsApp
    true,  -- Telegram
    true,  -- MCP Bot
    true,  -- Busca ilimitada
    true,  -- IA Analysis
    true,  -- Predictions
    true,  -- Doc Generation
    true,  -- White Label
    true,  -- API Access
    true,  -- Priority Support
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Create a function to get plan by name
CREATE OR REPLACE FUNCTION get_plan_by_name(plan_name VARCHAR(100))
RETURNS TABLE (
    id UUID,
    name VARCHAR(100),
    display_name VARCHAR(200),
    description TEXT,
    price_monthly BIGINT,
    price_yearly BIGINT,
    trial_days INTEGER,
    active BOOLEAN,
    max_processes INTEGER,
    max_users INTEGER,
    max_ai_requests INTEGER,
    max_bot_commands INTEGER,
    has_whatsapp BOOLEAN,
    has_telegram BOOLEAN,
    has_mcp_bot BOOLEAN,
    has_unlimited_search BOOLEAN,
    has_ai_analysis BOOLEAN,
    has_predictions BOOLEAN,
    has_doc_generation BOOLEAN,
    has_white_label BOOLEAN,
    has_api_access BOOLEAN,
    has_priority_support BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.name, p.display_name, p.description, p.price_monthly, p.price_yearly, p.trial_days, p.active,
           p.max_processes, p.max_users, p.max_ai_requests, p.max_bot_commands,
           p.has_whatsapp, p.has_telegram, p.has_mcp_bot, p.has_unlimited_search,
           p.has_ai_analysis, p.has_predictions, p.has_doc_generation, p.has_white_label,
           p.has_api_access, p.has_priority_support, p.created_at, p.updated_at
    FROM plans p
    WHERE p.name = plan_name AND p.active = true;
END;
$$ LANGUAGE plpgsql;

-- Create a function to check if plan has feature
CREATE OR REPLACE FUNCTION plan_has_feature(plan_id UUID, feature_name VARCHAR(50))
RETURNS BOOLEAN AS $$
DECLARE
    result BOOLEAN := false;
BEGIN
    SELECT CASE feature_name
        WHEN 'whatsapp' THEN has_whatsapp
        WHEN 'telegram' THEN has_telegram
        WHEN 'mcp_bot' THEN has_mcp_bot
        WHEN 'unlimited_search' THEN has_unlimited_search
        WHEN 'ai_analysis' THEN has_ai_analysis
        WHEN 'predictions' THEN has_predictions
        WHEN 'doc_generation' THEN has_doc_generation
        WHEN 'white_label' THEN has_white_label
        WHEN 'api_access' THEN has_api_access
        WHEN 'priority_support' THEN has_priority_support
        ELSE false
    END INTO result
    FROM plans
    WHERE id = plan_id;
    
    RETURN COALESCE(result, false);
END;
$$ LANGUAGE plpgsql;

-- Create a function to get plan price by cycle
CREATE OR REPLACE FUNCTION get_plan_price(plan_id UUID, billing_cycle VARCHAR(20))
RETURNS BIGINT AS $$
DECLARE
    price BIGINT := 0;
BEGIN
    SELECT CASE billing_cycle
        WHEN 'yearly' THEN price_yearly
        ELSE price_monthly
    END INTO price
    FROM plans
    WHERE id = plan_id;
    
    RETURN COALESCE(price, 0);
END;
$$ LANGUAGE plpgsql;

-- Create a view for active plans
CREATE OR REPLACE VIEW active_plans AS
SELECT 
    id,
    name,
    display_name,
    description,
    price_monthly,
    price_yearly,
    trial_days,
    max_processes,
    max_users,
    max_ai_requests,
    max_bot_commands,
    has_whatsapp,
    has_telegram,
    has_mcp_bot,
    has_unlimited_search,
    has_ai_analysis,
    has_predictions,
    has_doc_generation,
    has_white_label,
    has_api_access,
    has_priority_support,
    created_at,
    updated_at
FROM plans
WHERE active = true
ORDER BY 
    CASE name
        WHEN 'starter' THEN 1
        WHEN 'professional' THEN 2
        WHEN 'business' THEN 3
        WHEN 'enterprise' THEN 4
        ELSE 5
    END;