ALTER TABLE jira_issues_calculations
    ALTER COLUMN cycle_time SET DEFAULT NULL,
    ALTER COLUMN working_cycle_time SET DEFAULT NULL,
    ALTER COLUMN lead_time SET DEFAULT NULL,
    ALTER COLUMN working_lead_time SET DEFAULT NULL,
    ALTER COLUMN system_delay_time SET DEFAULT NULL,
    ALTER COLUMN working_system_delay_time SET DEFAULT NULL;