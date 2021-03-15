CREATE TYPE WORK_TYPE AS ENUM('story', 'task', 'bug', 'subtask');

CREATE TABLE IF NOT EXISTS jira_data (
    _id BIGINT GENERATED ALWAYS AS IDENTITY,
    jira_id INT NOT NULL UNIQUE,
    jira_key VARCHAR NOT NULL UNIQUE,
    jira_work_type WORK_TYPE NOT NULL,
    jira_created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    jira_completed_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    jira_work_started_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY(_id)
);

CREATE INDEX IF NOT EXISTS jira_data_jira_id_idx ON jira_data(jira_id);
CREATE INDEX IF NOT EXISTS jira_data_jira_key_idx ON jira_data(jira_key);
CREATE INDEX IF NOT EXISTS jira_data_jira_work_type_idx ON jira_data(jira_work_type);

CREATE TABLE IF NOT EXISTS jira_transitions(
    _id BIGINT GENERATED ALWAYS AS IDENTITY,
    jira_data_id BIGINT NOT NULL,
    jira_id INT NOT NULL UNIQUE,
    jira_transition_to VARCHAR NOT NULL,
    jira_transition_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    PRIMARY KEY(_id),
    CONSTRAINT jira_data_id_fk FOREIGN KEY(jira_data_id) REFERENCES jira_data(_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS jira_transitions_jira_id_idx ON jira_transitions(jira_id);
CREATE INDEX IF NOT EXISTS jira_transitions_jira_transition_to_idx ON jira_transitions(jira_transition_to);
CREATE INDEX IF NOT EXISTS jira_transitions_jira_transition_at_idx ON jira_transitions(jira_transition_at);
