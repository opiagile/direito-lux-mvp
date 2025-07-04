-- Criar tabela de processos jurídicos
CREATE TABLE IF NOT EXISTS processes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    number VARCHAR(255) NOT NULL,
    court VARCHAR(255) NOT NULL,
    subject TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused', 'archived')),
    monitoring BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Garantir número único por tenant
    CONSTRAINT unique_process_number_per_tenant UNIQUE (tenant_id, number)
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_processes_tenant_id ON processes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_processes_status ON processes(status);
CREATE INDEX IF NOT EXISTS idx_processes_number ON processes(number);
CREATE INDEX IF NOT EXISTS idx_processes_monitoring ON processes(monitoring);
CREATE INDEX IF NOT EXISTS idx_processes_created_at ON processes(created_at);

-- Trigger para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_processes_updated_at 
    BEFORE UPDATE ON processes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Adicionar comentários para documentação
COMMENT ON TABLE processes IS 'Tabela principal de processos jurídicos monitorados';
COMMENT ON COLUMN processes.id IS 'Identificador único do processo';
COMMENT ON COLUMN processes.tenant_id IS 'ID do escritório/tenant dono do processo';
COMMENT ON COLUMN processes.number IS 'Número do processo no tribunal (ex: 5001234-12.2024.8.26.0100)';
COMMENT ON COLUMN processes.court IS 'Tribunal onde tramita o processo (ex: TJSP, TRT2)';
COMMENT ON COLUMN processes.subject IS 'Assunto/descrição do processo';
COMMENT ON COLUMN processes.status IS 'Status do processo: active (ativo), paused (pausado), archived (arquivado)';
COMMENT ON COLUMN processes.monitoring IS 'Se o processo está sendo monitorado automaticamente';
COMMENT ON COLUMN processes.created_at IS 'Data/hora de criação do registro';
COMMENT ON COLUMN processes.updated_at IS 'Data/hora da última atualização';