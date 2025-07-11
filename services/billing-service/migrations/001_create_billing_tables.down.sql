-- Migration: Drop billing tables
-- Created: 2025-07-11

-- Drop triggers
DROP TRIGGER IF EXISTS update_plans_updated_at ON plans;
DROP TRIGGER IF EXISTS update_customers_updated_at ON customers;
DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP TRIGGER IF EXISTS update_invoices_updated_at ON invoices;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_subscriptions_tenant_id;
DROP INDEX IF EXISTS idx_subscriptions_plan_id;
DROP INDEX IF EXISTS idx_subscriptions_status;
DROP INDEX IF EXISTS idx_subscriptions_next_billing_date;
DROP INDEX IF EXISTS idx_subscriptions_asaas_id;

DROP INDEX IF EXISTS idx_payments_subscription_id;
DROP INDEX IF EXISTS idx_payments_tenant_id;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_due_date;
DROP INDEX IF EXISTS idx_payments_asaas_payment_id;
DROP INDEX IF EXISTS idx_payments_now_payment_id;

DROP INDEX IF EXISTS idx_invoices_subscription_id;
DROP INDEX IF EXISTS idx_invoices_tenant_id;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_due_date;
DROP INDEX IF EXISTS idx_invoices_number;
DROP INDEX IF EXISTS idx_invoices_nfe_number;

DROP INDEX IF EXISTS idx_customers_tenant_id;
DROP INDEX IF EXISTS idx_customers_document;
DROP INDEX IF EXISTS idx_customers_email;
DROP INDEX IF EXISTS idx_customers_asaas_id;

DROP INDEX IF EXISTS idx_plans_name;
DROP INDEX IF EXISTS idx_plans_active;

-- Drop tables (reverse order due to foreign keys)
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS plans;