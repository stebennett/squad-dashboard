CREATE TYPE INGESTION_TYPE AS ENUM('backfill', 'incremental');

ALTER TABLE jira_config ADD COLUMN last_ingestion_run_started TIMESTAMP DEFAULT NULL;
ALTER TABLE jira_config ADD COLUMN last_ingestion_run_completed TIMESTAMP DEFAULT NULL;
ALTER TABLE jira_config ADD COLUMN last_ingestion_type INGESTION_TYPE DEFAULT NULL;
