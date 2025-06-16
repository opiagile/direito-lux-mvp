-- Migration: Drop quota_limits table
-- Description: Remove a tabela de limites de quotas

DROP TABLE IF EXISTS quota_limits CASCADE;