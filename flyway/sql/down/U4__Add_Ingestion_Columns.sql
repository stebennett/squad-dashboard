CREATE TYPE INGESTION_TYPE AS ENUM('backfill', 'incremental');

ALTER TABLE jira_config DROP COLUMN last_ingestion_run_started;
ALTER TABLE jira_config DROP COLUMN last_ingestion_run_completed;
ALTER TABLE jira_config DROP COLUMN last_ingestion_type;
