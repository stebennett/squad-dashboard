ALTER TABLE jira_issues_calculations 
    ADD COLUMN issue_created_at TIMESTAMP DEFAULT NULL,
    ADD COLUMN issue_started_at TIMESTAMP DEFAULT NULL,
    ADD COLUMN issue_completed_at TIMESTAMP DEFAULT NULL; 
