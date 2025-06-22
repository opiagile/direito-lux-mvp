-- =============================================================================
-- DIREITO LUX - Criação de Usuário e Banco de Dados
-- =============================================================================
-- Este script é executado pelo superusuário postgres durante a inicialização

-- Criar o role direito_lux
CREATE ROLE direito_lux WITH 
    LOGIN 
    PASSWORD 'dev_password_123' 
    CREATEDB 
    SUPERUSER;

-- Criar o banco de dados
CREATE DATABASE direito_lux_dev 
    OWNER direito_lux 
    ENCODING 'UTF8';

-- Conectar ao novo banco para criar extensões
\c direito_lux_dev

-- Criar extensões necessárias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Verificação
SELECT 'Usuario e banco direito_lux_dev criados com sucesso!' as status;