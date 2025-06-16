-- Migration: Drop quota_usage table
-- Description: Remove a tabela de uso de quotas

DROP TABLE IF EXISTS quota_usage CASCADE;