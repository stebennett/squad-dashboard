CREATE TABLE jira_issues_calculations(
    _id INT GENERATED ALWAYS AS IDENTITY,
    issue_key VARCHAR(16) UNIQUE NOT NULL,
    week_create INT,
    year_create INT,
    week_complete INT,
    year_complete INT,
    week_start INT,
    year_start INT,
    cycle_time INT,
    working_cycle_time INT,
    lead_time INT,
    working_lead_time INT,
    system_delay_time INT,
    working_system_delay_time INT
);

CREATE INDEX CONCURRENTLY year_week_create_idx ON jira_issues_calculations(week_create, year_create);
CREATE INDEX CONCURRENTLY year_week_complete_idx ON jira_issues_calculations(week_complete, year_complete);
CREATE INDEX CONCURRENTLY year_week_start_idx ON jira_issues_calculations(week_start, year_start);
CREATE INDEX CONCURRENTLY issue_key_idx ON jira_issues_calculations(issue_key);