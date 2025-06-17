-- Migration: Create parties table
-- Description: Tabela para armazenar as partes dos processos
-- Author: Direito Lux Team
-- Date: 2025-01-16

-- Criar tabela de partes
CREATE TABLE IF NOT EXISTS parties (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    process_id UUID NOT NULL REFERENCES processes(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('individual', 'legal')),
    name VARCHAR(500) NOT NULL,
    document VARCHAR(20),
    document_type VARCHAR(10) CHECK (document_type IN ('cpf', 'cnpj', 'other')),
    role VARCHAR(20) NOT NULL CHECK (role IN (
        'plaintiff', 'defendant', 'third_party', 'intervenor', 
        'assistant', 'guardian', 'representative', 'expert', 'witness'
    )),
    is_active BOOLEAN DEFAULT TRUE,
    lawyer JSONB,
    contact JSONB DEFAULT '{
        "email": "",
        "phone": "",
        "cell_phone": "",
        "website": ""
    }',
    address JSONB DEFAULT '{
        "street": "",
        "number": "",
        "complement": "",
        "district": "",
        "city": "",
        "state": "",
        "zip_code": "",
        "country": "Brasil"
    }',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT check_document_format CHECK (
        (document_type = 'cpf' AND LENGTH(REGEXP_REPLACE(document, '[^0-9]', '', 'g')) = 11) OR
        (document_type = 'cnpj' AND LENGTH(REGEXP_REPLACE(document, '[^0-9]', '', 'g')) = 14) OR
        (document_type = 'other') OR
        (document IS NULL)
    ),
    CONSTRAINT check_lawyer_structure CHECK (
        lawyer IS NULL OR (
            lawyer ? 'name' AND 
            (lawyer->>'name') != '' AND
            CASE 
                WHEN lawyer ? 'oab' THEN (lawyer->>'oab') ~ '^[0-9]{4,6}$'
                ELSE true
            END AND
            CASE 
                WHEN lawyer ? 'oab_state' THEN (lawyer->>'oab_state') ~ '^[A-Z]{2}$'
                ELSE true
            END
        )
    ),
    CONSTRAINT check_type_document_consistency CHECK (
        (type = 'individual' AND (document_type IS NULL OR document_type = 'cpf')) OR
        (type = 'legal' AND (document_type IS NULL OR document_type = 'cnpj')) OR
        (document_type = 'other')
    )
);

-- Comentários nas colunas
COMMENT ON TABLE parties IS 'Tabela das partes dos processos (autores, réus, etc.)';
COMMENT ON COLUMN parties.id IS 'Identificador único da parte';
COMMENT ON COLUMN parties.process_id IS 'Referência ao processo';
COMMENT ON COLUMN parties.type IS 'Tipo da parte (pessoa física ou jurídica)';
COMMENT ON COLUMN parties.name IS 'Nome/razão social da parte';
COMMENT ON COLUMN parties.document IS 'Documento de identificação (CPF/CNPJ)';
COMMENT ON COLUMN parties.document_type IS 'Tipo do documento';
COMMENT ON COLUMN parties.role IS 'Papel da parte no processo';
COMMENT ON COLUMN parties.is_active IS 'Flag indicando se a parte está ativa';
COMMENT ON COLUMN parties.lawyer IS 'Dados do advogado (JSON)';
COMMENT ON COLUMN parties.contact IS 'Informações de contato (JSON)';
COMMENT ON COLUMN parties.address IS 'Endereço da parte (JSON)';
COMMENT ON COLUMN parties.created_at IS 'Data de criação do registro';
COMMENT ON COLUMN parties.updated_at IS 'Data da última atualização';

-- Trigger para atualizar updated_at
CREATE TRIGGER update_parties_updated_at 
    BEFORE UPDATE ON parties 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Função para validar CPF (simplificada)
CREATE OR REPLACE FUNCTION validate_cpf(cpf_input TEXT)
RETURNS BOOLEAN AS $$
DECLARE
    cpf_clean TEXT;
    digit1 INT;
    digit2 INT;
    sum_value INT;
    i INT;
BEGIN
    -- Remover caracteres não numéricos
    cpf_clean := REGEXP_REPLACE(cpf_input, '[^0-9]', '', 'g');
    
    -- Verificar se tem 11 dígitos
    IF LENGTH(cpf_clean) != 11 THEN
        RETURN FALSE;
    END IF;
    
    -- Verificar se todos os dígitos são iguais
    IF cpf_clean ~ '^([0-9])\1{10}$' THEN
        RETURN FALSE;
    END IF;
    
    -- Para simplificar, aceitar qualquer CPF com 11 dígitos diferentes
    -- Em produção, implementar algoritmo completo de validação
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- Função para validar CNPJ (simplificada)
CREATE OR REPLACE FUNCTION validate_cnpj(cnpj_input TEXT)
RETURNS BOOLEAN AS $$
DECLARE
    cnpj_clean TEXT;
BEGIN
    -- Remover caracteres não numéricos
    cnpj_clean := REGEXP_REPLACE(cnpj_input, '[^0-9]', '', 'g');
    
    -- Verificar se tem 14 dígitos
    IF LENGTH(cnpj_clean) != 14 THEN
        RETURN FALSE;
    END IF;
    
    -- Verificar se todos os dígitos são iguais
    IF cnpj_clean ~ '^([0-9])\1{13}$' THEN
        RETURN FALSE;
    END IF;
    
    -- Para simplificar, aceitar qualquer CNPJ com 14 dígitos diferentes
    -- Em produção, implementar algoritmo completo de validação
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- Trigger para validar documento antes de inserir/atualizar
CREATE OR REPLACE FUNCTION validate_party_document()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.document IS NOT NULL AND NEW.document != '' THEN
        CASE NEW.document_type
            WHEN 'cpf' THEN
                IF NOT validate_cpf(NEW.document) THEN
                    RAISE EXCEPTION 'CPF inválido: %', NEW.document;
                END IF;
            WHEN 'cnpj' THEN
                IF NOT validate_cnpj(NEW.document) THEN
                    RAISE EXCEPTION 'CNPJ inválido: %', NEW.document;
                END IF;
        END CASE;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_party_document_trigger
    BEFORE INSERT OR UPDATE ON parties
    FOR EACH ROW
    EXECUTE FUNCTION validate_party_document();