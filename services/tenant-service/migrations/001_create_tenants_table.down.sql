-- Migration: Drop tenants table
-- Description: Remove a tabela de tenants e seus índices

DROP TABLE IF EXISTS tenants CASCADE;