-- Migration: Create movements table
-- Description: Tabela para armazenar movimentações/andamentos dos processos
-- Author: Direito Lux Team
-- Date: 2025-01-16

-- Criar tabela de movimentações
CREATE TABLE IF NOT EXISTS movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    process_id UUID NOT NULL REFERENCES processes(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    sequence INTEGER NOT NULL,
    external_id VARCHAR(100),
    date TIMESTAMP NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN (
        'decision', 'order', 'act', 'filing', 'hearing', 
        'publication', 'remittance', 'return', 'citation', 
        'intimation', 'archiving', 'reactivation', 'appeal', 
        'execution', 'other'
    )),
    code VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    content TEXT,
    judge VARCHAR(200),
    responsible VARCHAR(200),
    attachments JSONB DEFAULT '[]',
    related_parties TEXT[] DEFAULT '{}',
    is_important BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT TRUE,
    notification_sent BOOLEAN DEFAULT FALSE,
    tags TEXT[] DEFAULT '{}',
    metadata JSONB DEFAULT '{
        "original_source": "",
        "datajud_id": "",
        "import_batch": "",
        "keywords": [],
        "analysis": {
            "sentiment": "",
            "importance": 0,
            "category": "",
            "has_deadline": false,
            "deadline_date": null,
            "requires_action": false,
            "action_type": "",
            "confidence": 0.0,
            "processed_by": "",
            "processed_at": null
        },
        "custom_fields": {}
    }',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    synced_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT unique_process_sequence UNIQUE(process_id, sequence),
    CONSTRAINT unique_external_id UNIQUE(external_id) DEFERRABLE,
    CONSTRAINT check_sequence_positive CHECK (sequence > 0),
    CONSTRAINT check_date_not_future CHECK (date <= CURRENT_TIMESTAMP + INTERVAL '1 day'),
    CONSTRAINT check_importance_notification CHECK (
        (is_important = false AND notification_sent = false) OR 
        (is_important = true)
    )
);

-- Comentários nas colunas
COMMENT ON TABLE movements IS 'Tabela de movimentações/andamentos dos processos';
COMMENT ON COLUMN movements.id IS 'Identificador único da movimentação';
COMMENT ON COLUMN movements.process_id IS 'Referência ao processo';
COMMENT ON COLUMN movements.tenant_id IS 'Identificador do tenant';
COMMENT ON COLUMN movements.sequence IS 'Número sequencial da movimentação no processo';
COMMENT ON COLUMN movements.external_id IS 'Identificador externo (DataJud)';
COMMENT ON COLUMN movements.date IS 'Data/hora da movimentação';
COMMENT ON COLUMN movements.type IS 'Tipo da movimentação';
COMMENT ON COLUMN movements.code IS 'Código CNJ da movimentação';
COMMENT ON COLUMN movements.title IS 'Título da movimentação';
COMMENT ON COLUMN movements.description IS 'Descrição da movimentação';
COMMENT ON COLUMN movements.content IS 'Conteúdo completo da movimentação';
COMMENT ON COLUMN movements.judge IS 'Juiz responsável pela movimentação';
COMMENT ON COLUMN movements.responsible IS 'Responsável pela movimentação';
COMMENT ON COLUMN movements.attachments IS 'Anexos da movimentação (JSON array)';
COMMENT ON COLUMN movements.related_parties IS 'IDs das partes relacionadas';
COMMENT ON COLUMN movements.is_important IS 'Flag indicando se é uma movimentação importante';
COMMENT ON COLUMN movements.is_public IS 'Flag indicando se é pública';
COMMENT ON COLUMN movements.notification_sent IS 'Flag indicando se notificação foi enviada';
COMMENT ON COLUMN movements.tags IS 'Tags para categorização';
COMMENT ON COLUMN movements.metadata IS 'Metadados adicionais incluindo análise IA (JSON)';
COMMENT ON COLUMN movements.created_at IS 'Data de criação do registro';
COMMENT ON COLUMN movements.updated_at IS 'Data da última atualização';
COMMENT ON COLUMN movements.synced_at IS 'Data da última sincronização';

-- Trigger para atualizar updated_at
CREATE TRIGGER update_movements_updated_at 
    BEFORE UPDATE ON movements 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Função para auto-incrementar sequence por processo
CREATE OR REPLACE FUNCTION set_movement_sequence()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.sequence IS NULL OR NEW.sequence = 0 THEN
        SELECT COALESCE(MAX(sequence), 0) + 1 
        INTO NEW.sequence 
        FROM movements 
        WHERE process_id = NEW.process_id;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER set_movement_sequence_trigger
    BEFORE INSERT ON movements
    FOR EACH ROW
    EXECUTE FUNCTION set_movement_sequence();