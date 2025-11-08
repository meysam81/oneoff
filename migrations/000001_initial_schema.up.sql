-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    color TEXT,
    icon TEXT,
    is_archived BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc'))
);

-- Create default project
INSERT INTO projects (id, name, description, color, icon)
VALUES ('default', 'Default Project', 'Default project for all jobs', '#6366f1', 'folder');

-- Tags table
CREATE TABLE IF NOT EXISTS tags (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT NOT NULL UNIQUE,
    color TEXT,
    is_default BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc'))
);

-- Insert default tags
INSERT INTO tags (name, color, is_default) VALUES
    ('maintenance', '#8b5cf6', 1),
    ('deployment', '#10b981', 1),
    ('backup', '#3b82f6', 1),
    ('cleanup', '#f59e0b', 1),
    ('migration', '#ec4899', 1),
    ('testing', '#14b8a6', 1),
    ('production', '#ef4444', 1),
    ('staging', '#f97316', 1);

-- Jobs table
CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    config TEXT NOT NULL, -- JSON blob for job-specific configuration
    scheduled_at DATETIME NOT NULL, -- UTC timestamp
    priority INTEGER NOT NULL DEFAULT 5 CHECK(priority >= 1 AND priority <= 10),
    project_id TEXT NOT NULL DEFAULT 'default',
    timezone TEXT NOT NULL DEFAULT 'UTC',
    status TEXT NOT NULL DEFAULT 'scheduled' CHECK(status IN ('scheduled', 'running', 'completed', 'failed', 'cancelled')),
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc')),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Index for faster queries
CREATE INDEX idx_jobs_scheduled_at ON jobs(scheduled_at);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_project_id ON jobs(project_id);
CREATE INDEX idx_jobs_priority ON jobs(priority);

-- Job executions table
CREATE TABLE IF NOT EXISTS job_executions (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    job_id TEXT NOT NULL,
    started_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc')),
    completed_at DATETIME,
    status TEXT NOT NULL DEFAULT 'running' CHECK(status IN ('running', 'completed', 'failed', 'cancelled')),
    output TEXT, -- Execution output/logs
    exit_code INTEGER,
    error TEXT,
    duration_ms INTEGER,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc')),
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
);

-- Index for execution queries
CREATE INDEX idx_job_executions_job_id ON job_executions(job_id);
CREATE INDEX idx_job_executions_status ON job_executions(status);
CREATE INDEX idx_job_executions_started_at ON job_executions(started_at);

-- Job tags many-to-many relationship
CREATE TABLE IF NOT EXISTS job_tags (
    job_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc')),
    PRIMARY KEY (job_id, tag_id),
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Job chains table
CREATE TABLE IF NOT EXISTS job_chains (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT NOT NULL,
    project_id TEXT NOT NULL DEFAULT 'default',
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc')),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Job chain links (defines the sequence)
CREATE TABLE IF NOT EXISTS job_chain_links (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    chain_id TEXT NOT NULL,
    job_id TEXT NOT NULL,
    sequence_order INTEGER NOT NULL,
    stop_on_failure BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc')),
    FOREIGN KEY (chain_id) REFERENCES job_chains(id) ON DELETE CASCADE,
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
    UNIQUE(chain_id, sequence_order)
);

-- System configuration table
CREATE TABLE IF NOT EXISTS system_config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL, -- JSON for complex configs
    updated_at DATETIME NOT NULL DEFAULT (datetime('now', 'utc'))
);

-- Insert default system config
INSERT INTO system_config (key, value) VALUES
    ('workers_count', '0'), -- 0 means auto-detect (N/2 cores)
    ('default_timezone', 'UTC'),
    ('log_retention_days', '90'),
    ('default_priority', '5');

-- Triggers for updated_at timestamps
CREATE TRIGGER update_projects_updated_at AFTER UPDATE ON projects
BEGIN
    UPDATE projects SET updated_at = datetime('now', 'utc') WHERE id = NEW.id;
END;

CREATE TRIGGER update_jobs_updated_at AFTER UPDATE ON jobs
BEGIN
    UPDATE jobs SET updated_at = datetime('now', 'utc') WHERE id = NEW.id;
END;

CREATE TRIGGER update_system_config_updated_at AFTER UPDATE ON system_config
BEGIN
    UPDATE system_config SET updated_at = datetime('now', 'utc') WHERE key = NEW.key;
END;
