-- Drop triggers
DROP TRIGGER IF EXISTS update_projects_updated_at;
DROP TRIGGER IF EXISTS update_jobs_updated_at;
DROP TRIGGER IF EXISTS update_system_config_updated_at;

-- Drop tables in reverse order (respect foreign keys)
DROP TABLE IF EXISTS system_config;
DROP TABLE IF EXISTS job_chain_links;
DROP TABLE IF EXISTS job_chains;
DROP TABLE IF EXISTS job_tags;
DROP TABLE IF EXISTS job_executions;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS projects;
