DROP TRIGGER tgr_jira_status_updated ON jira_transitions;
DROP TRIGGER tgr_jira_data_lead_time ON jira_data;
DROP TRIGGER tgr_jira_data_cycle_time ON jira_data;

DROP FUNCTION fn_calculate_lead_time;
DROP FUNCTION fn_calculate_cycle_time;
DROP FUNCTION fn_calculate_dev_start_time;

ALTER TABLE jira_data DROP COLUMN jira_project_key;

DROP TABLE flow_measures;

DROP TABLE jira_config;
