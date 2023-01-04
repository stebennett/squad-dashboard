BEGIN;

DROP TABLE IF EXISTS jira_work_states;
CREATE TABLE IF NOT EXISTS jira_work_state_start(
    _id INT GENERATED ALWAYS AS IDENTITY,
    project VARCHAR(16) NOT NULL,
    state_type work_state_type NOT NULL,
    state_name VARCHAR(128) NOT NULL,

    CONSTRAINT jwss_project_state_name_cnstr UNIQUE (project, state_name)
);

COMMIT;