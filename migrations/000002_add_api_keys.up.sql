-- API Keys table for authentication
CREATE TABLE IF NOT EXISTS api_keys (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT NOT NULL,
    key_hash TEXT NOT NULL UNIQUE,
    key_prefix TEXT NOT NULL,
    scopes TEXT NOT NULL DEFAULT 'read,write',
    last_used_at DATETIME,
    expires_at DATETIME,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc'))
);

-- Index for fast key lookup during authentication
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_is_active ON api_keys(is_active);

-- Webhooks table for notifications
CREATE TABLE IF NOT EXISTS webhooks (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    secret TEXT,
    events TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc'))
);

-- Webhook deliveries for tracking delivery attempts
CREATE TABLE IF NOT EXISTS webhook_deliveries (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    webhook_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'success', 'failed')),
    response_code INTEGER,
    response_body TEXT,
    error TEXT,
    attempts INTEGER NOT NULL DEFAULT 0,
    next_retry_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc')),
    FOREIGN KEY (webhook_id) REFERENCES webhooks(id) ON DELETE CASCADE
);

-- Indexes for webhook delivery processing
CREATE INDEX idx_webhook_deliveries_webhook_id ON webhook_deliveries(webhook_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status, next_retry_at);

-- Trigger for webhook updated_at
CREATE TRIGGER update_webhooks_updated_at AFTER UPDATE ON webhooks
BEGIN
    UPDATE webhooks SET updated_at = datetime('now', 'utc') WHERE id = NEW.id;
END;
