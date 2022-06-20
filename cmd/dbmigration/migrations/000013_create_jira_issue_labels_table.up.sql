BEGIN;

CREATE TABLE IF NOT EXISTS jira_issue_labels(
    _id INT GENERATED ALWAYS AS IDENTITY,
    issue_key VARCHAR(16) NOT NULL,
    label VARCHAR(128) NOT NULL,
    CONSTRAINT issue_key_label_c UNIQUE (issue_key, label)
);

CREATE INDEX issue_labels_issue_key_idx ON jira_issue_labels(issue_key);
CREATE INDEX issue_labels_issue_label_idx ON jira_issue_labels(label);

COMMIT;