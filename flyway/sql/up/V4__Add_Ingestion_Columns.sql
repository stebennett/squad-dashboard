CREATE TYPE INGESTION_TYPE AS ENUM('Backfill', 'Incremental');

ALTER TABLE jira_config ADD COLUMN last_ingestion_run_started TIMESTAMP WITH TIME ZONE DEFAULT NULL;
ALTER TABLE jira_config ADD COLUMN last_ingestion_run_completed TIMESTAMP WITH TIME ZONE DEFAULT NULL;
ALTER TABLE jira_config ADD COLUMN last_ingestion_type INGESTION_TYPE DEFAULT NULL;
