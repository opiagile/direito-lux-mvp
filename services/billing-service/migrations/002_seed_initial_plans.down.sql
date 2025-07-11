-- Migration: Remove initial plans
-- Created: 2025-07-11

-- Drop view
DROP VIEW IF EXISTS active_plans;

-- Drop functions
DROP FUNCTION IF EXISTS get_plan_by_name(VARCHAR(100));
DROP FUNCTION IF EXISTS plan_has_feature(UUID, VARCHAR(50));
DROP FUNCTION IF EXISTS get_plan_price(UUID, VARCHAR(20));

-- Delete all plans
DELETE FROM plans;