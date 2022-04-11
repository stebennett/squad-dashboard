CREATE TABLE IF NOT EXISTS jira_issues(
    _id INT GENERATED ALWAYS AS IDENTITY,
    issue_key VARCHAR(16) NOT NULL,
    parent_key VARCHAR(16),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
);
