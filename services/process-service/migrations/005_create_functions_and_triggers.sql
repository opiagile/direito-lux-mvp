-- Migration: Create functions and triggers
-- Description: Funções e triggers específicos do negócio
-- Author: Direito Lux Team
-- Date: 2025-01-16

-- ========================================
-- FUNÇÕES UTILITÁRIAS
-- ========================================

-- Função para limpar e formatar número CNJ
CREATE OR REPLACE FUNCTION format_cnj_number(input_number TEXT)
RETURNS TEXT AS $$
DECLARE
    clean_number TEXT;
BEGIN
    -- Remove todos os caracteres não numéricos
    clean_number := REGEXP_REPLACE(input_number, '[^0-9]', '', 'g');
    
    -- Verifica se tem 20 dígitos
    IF LENGTH(clean_number) != 20 THEN
        RETURN input_number; -- Retorna original se inválido
    END IF;
    
    -- Formata no padrão CNJ: NNNNNNN-DD.AAAA.J.TR.OOOO
    RETURN SUBSTRING(clean_number, 1, 7) || '-' || 
           SUBSTRING(clean_number, 8, 2) || '.' ||
           SUBSTRING(clean_number, 10, 4) || '.' ||
           SUBSTRING(clean_number, 14, 1) || '.' ||
           SUBSTRING(clean_number, 15, 2) || '.' ||
           SUBSTRING(clean_number, 17, 4);
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Função para validar número CNJ
CREATE OR REPLACE FUNCTION validate_cnj_number(cnj_number TEXT)
RETURNS BOOLEAN AS $$
DECLARE
    clean_number TEXT;
    verification_digits TEXT;
    calculated_digits TEXT;
BEGIN
    -- Remove caracteres não numéricos
    clean_number := REGEXP_REPLACE(cnj_number, '[^0-9]', '', 'g');
    
    -- Verifica se tem 20 dígitos
    IF LENGTH(clean_number) != 20 THEN
        RETURN FALSE;
    END IF;
    
    -- Para simplificar, aceita qualquer número com 20 dígitos
    -- Em produção, implementar algoritmo completo de validação CNJ
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Função para calcular dias úteis entre duas datas
CREATE OR REPLACE FUNCTION calculate_business_days(start_date DATE, end_date DATE)
RETURNS INTEGER AS $$
DECLARE
    days INTEGER := 0;
    current_date DATE := start_date;
BEGIN
    WHILE current_date <= end_date LOOP
        -- Conta apenas dias úteis (segunda a sexta)
        IF EXTRACT(dow FROM current_date) BETWEEN 1 AND 5 THEN
            days := days + 1;
        END IF;
        current_date := current_date + 1;
    END LOOP;
    
    RETURN days;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Função para extrair palavras-chave de texto
CREATE OR REPLACE FUNCTION extract_keywords(input_text TEXT, min_length INTEGER DEFAULT 4)
RETURNS TEXT[] AS $$
DECLARE
    words TEXT[];
    cleaned_text TEXT;
    stop_words TEXT[] := ARRAY['para', 'pela', 'pelo', 'esta', 'este', 'essa', 'esse', 'aquela', 'aquele', 'desde', 'entre', 'sobre', 'durante', 'antes', 'depois', 'contra', 'conforme', 'segundo'];
    word TEXT;
    result TEXT[] := ARRAY[]::TEXT[];
BEGIN
    -- Limpa o texto e converte para minúsculo
    cleaned_text := LOWER(REGEXP_REPLACE(input_text, '[^a-záàâãéèêíìîóòôõúùûç\s]', ' ', 'g'));
    
    -- Divide em palavras
    words := string_to_array(cleaned_text, ' ');
    
    -- Filtra palavras relevantes
    FOREACH word IN ARRAY words LOOP
        word := TRIM(word);
        IF LENGTH(word) >= min_length AND NOT (word = ANY(stop_words)) THEN
            result := array_append(result, word);
        END IF;
    END LOOP;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- ========================================
-- TRIGGERS DE NEGÓCIO
-- ========================================

-- Trigger para atualizar last_movement_at no processo quando há nova movimentação
CREATE OR REPLACE FUNCTION update_process_last_movement()
RETURNS TRIGGER AS $$
BEGIN
    -- Atualiza a data da última movimentação no processo
    UPDATE processes 
    SET last_movement_at = NEW.date,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.process_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_process_last_movement_trigger
    AFTER INSERT ON movements
    FOR EACH ROW
    EXECUTE FUNCTION update_process_last_movement();

-- Trigger para extrair palavras-chave automaticamente de movimentações
CREATE OR REPLACE FUNCTION auto_extract_movement_keywords()
RETURNS TRIGGER AS $$
DECLARE
    keywords TEXT[];
BEGIN
    -- Extrai palavras-chave do título e descrição
    keywords := extract_keywords(NEW.title || ' ' || NEW.description);
    
    -- Atualiza o metadata com as palavras-chave
    NEW.metadata := jsonb_set(
        COALESCE(NEW.metadata, '{}'::jsonb),
        '{keywords}',
        to_jsonb(keywords)
    );
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER auto_extract_movement_keywords_trigger
    BEFORE INSERT OR UPDATE ON movements
    FOR EACH ROW
    EXECUTE FUNCTION auto_extract_movement_keywords();

-- Trigger para detectar movimentações importantes automaticamente
CREATE OR REPLACE FUNCTION auto_detect_important_movements()
RETURNS TRIGGER AS $$
DECLARE
    important_keywords TEXT[] := ARRAY[
        'sentença', 'acórdão', 'decisão', 'despacho', 'citação', 'intimação',
        'audiência', 'recurso', 'agravo', 'apelação', 'embargos', 'execução',
        'penhora', 'arrematação', 'adjudicação', 'bloqueio', 'extinção'
    ];
    keyword TEXT;
    text_to_check TEXT;
BEGIN
    text_to_check := LOWER(NEW.title || ' ' || NEW.description);
    
    -- Verifica se contém palavras-chave importantes
    FOREACH keyword IN ARRAY important_keywords LOOP
        IF text_to_check LIKE '%' || keyword || '%' THEN
            NEW.is_important := TRUE;
            EXIT;
        END IF;
    END LOOP;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER auto_detect_important_movements_trigger
    BEFORE INSERT OR UPDATE ON movements
    FOR EACH ROW
    EXECUTE FUNCTION auto_detect_important_movements();

-- ========================================
-- FUNÇÕES DE ESTATÍSTICAS
-- ========================================

-- Função para calcular estatísticas de um processo
CREATE OR REPLACE FUNCTION get_process_statistics(process_uuid UUID)
RETURNS TABLE(
    total_movements INTEGER,
    important_movements INTEGER,
    last_movement_date TIMESTAMP,
    average_days_between_movements NUMERIC,
    most_common_movement_type TEXT
) AS $$
BEGIN
    RETURN QUERY
    WITH movement_stats AS (
        SELECT 
            COUNT(*)::INTEGER as total,
            COUNT(CASE WHEN is_important THEN 1 END)::INTEGER as important,
            MAX(date) as last_date,
            mode() WITHIN GROUP (ORDER BY type) as common_type
        FROM movements 
        WHERE process_id = process_uuid
    ),
    date_diffs AS (
        SELECT 
            AVG(EXTRACT(days FROM date - LAG(date) OVER (ORDER BY date)))::NUMERIC as avg_days
        FROM movements 
        WHERE process_id = process_uuid
    )
    SELECT 
        ms.total,
        ms.important,
        ms.last_date,
        dd.avg_days,
        ms.common_type::TEXT
    FROM movement_stats ms
    CROSS JOIN date_diffs dd;
END;
$$ LANGUAGE plpgsql;

-- Função para buscar processos similares
CREATE OR REPLACE FUNCTION find_similar_processes(
    reference_process_uuid UUID,
    similarity_threshold REAL DEFAULT 0.3,
    max_results INTEGER DEFAULT 10
)
RETURNS TABLE(
    process_id UUID,
    process_number TEXT,
    process_title TEXT,
    similarity_score REAL
) AS $$
BEGIN
    RETURN QUERY
    WITH reference_process AS (
        SELECT title, description, subject, tags
        FROM processes 
        WHERE id = reference_process_uuid
    )
    SELECT 
        p.id,
        p.number,
        p.title,
        (
            similarity(p.title, rp.title) * 0.4 +
            similarity(p.description, rp.description) * 0.3 +
            CASE WHEN p.tags && rp.tags THEN 0.3 ELSE 0 END
        ) as score
    FROM processes p, reference_process rp
    WHERE p.id != reference_process_uuid
      AND p.tenant_id = (SELECT tenant_id FROM processes WHERE id = reference_process_uuid)
    ORDER BY score DESC
    LIMIT max_results;
END;
$$ LANGUAGE plpgsql;

-- ========================================
-- FUNÇÕES DE LIMPEZA E MANUTENÇÃO
-- ========================================

-- Função para limpar dados antigos (soft delete)
CREATE OR REPLACE FUNCTION cleanup_old_data(days_to_keep INTEGER DEFAULT 2555) -- ~7 anos
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER := 0;
    cutoff_date DATE := CURRENT_DATE - INTERVAL '1 day' * days_to_keep;
BEGIN
    -- Arquiva processos muito antigos sem atividade
    UPDATE processes 
    SET status = 'archived',
        archived_at = CURRENT_TIMESTAMP,
        updated_at = CURRENT_TIMESTAMP
    WHERE status != 'archived'
      AND (last_movement_at < cutoff_date OR 
           (last_movement_at IS NULL AND created_at < cutoff_date));
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Função para reindexar e otimizar tabelas
CREATE OR REPLACE FUNCTION maintain_database()
RETURNS TEXT AS $$
DECLARE
    result TEXT := '';
BEGIN
    -- Atualiza estatísticas
    ANALYZE processes;
    ANALYZE movements;
    ANALYZE parties;
    
    result := result || 'Estatísticas atualizadas. ';
    
    -- Limpa índices não utilizados (seria implementado com mais critério)
    result := result || 'Manutenção concluída.';
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- ========================================
-- COMENTÁRIOS NAS FUNÇÕES
-- ========================================

COMMENT ON FUNCTION format_cnj_number(TEXT) IS 'Formata número CNJ no padrão oficial';
COMMENT ON FUNCTION validate_cnj_number(TEXT) IS 'Valida formato do número CNJ';
COMMENT ON FUNCTION calculate_business_days(DATE, DATE) IS 'Calcula dias úteis entre duas datas';
COMMENT ON FUNCTION extract_keywords(TEXT, INTEGER) IS 'Extrai palavras-chave relevantes de um texto';
COMMENT ON FUNCTION get_process_statistics(UUID) IS 'Retorna estatísticas completas de um processo';
COMMENT ON FUNCTION find_similar_processes(UUID, REAL, INTEGER) IS 'Encontra processos similares baseado em texto e tags';
COMMENT ON FUNCTION cleanup_old_data(INTEGER) IS 'Arquiva dados antigos para otimização';
COMMENT ON FUNCTION maintain_database() IS 'Executa rotinas de manutenção do banco';