-- Drop webhook delivery trigger
DROP TRIGGER IF EXISTS update_webhooks_updated_at;

-- Drop webhook deliveries table
DROP TABLE IF EXISTS webhook_deliveries;

-- Drop webhooks table
DROP TABLE IF EXISTS webhooks;

-- Drop API keys indexes
DROP INDEX IF EXISTS idx_api_keys_is_active;
DROP INDEX IF EXISTS idx_api_keys_key_hash;

-- Drop API keys table
DROP TABLE IF EXISTS api_keys;
