CREATE TABLE IF NOT EXISTS jira_transitions(
    _id INT GENERATED ALWAYS AS IDENTITY,
    issue_key VARCHAR(16) NOT NULL,
    from_state VARCHAR(128) NOT NULL,
    to_state VARCHAR(128) NOT NULL,
    created_at TIMESTAMP NOT NULL
);
