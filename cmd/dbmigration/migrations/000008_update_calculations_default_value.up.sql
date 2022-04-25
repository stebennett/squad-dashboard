ALTER TABLE jira_issues_calculations
    ALTER COLUMN cycle_time SET DEFAULT -1,
    ALTER COLUMN working_cycle_time SET DEFAULT -1,
    ALTER COLUMN lead_time SET DEFAULT -1,
    ALTER COLUMN working_lead_time SET DEFAULT -1,
    ALTER COLUMN system_delay_time SET DEFAULT -1,
    ALTER COLUMN working_system_delay_time SET DEFAULT -1;

UPDATE jira_issues_calculations SET cycle_time = -1 WHERE cycle_time IS NULL;
UPDATE jira_issues_calculations SET working_cycle_time = -1 WHERE working_cycle_time IS NULL;
UPDATE jira_issues_calculations SET lead_time = -1 WHERE lead_time IS NULL;
UPDATE jira_issues_calculations SET working_lead_time = -1 WHERE working_lead_time IS NULL;
UPDATE jira_issues_calculations SET system_delay_time = -1 WHERE system_delay_time IS NULL;
UPDATE jira_issues_calculations SET working_system_delay_time = -1 WHERE working_system_delay_time IS NULL;