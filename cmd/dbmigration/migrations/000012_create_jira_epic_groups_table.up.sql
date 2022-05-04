BEGIN;

CREATE TABLE IF NOT EXISTS jira_epic_groups(
    _id INT GENERATED ALWAYS AS IDENTITY,
    epic_key VARCHAR(16) NOT NULL,
    group_name VARCHAR(16) NOT NULL,
    CONSTRAINT epic_group_name_c UNIQUE (epic_key, group_name)
);

CREATE INDEX group_name_idx ON jira_epic_groups(group_name);

COMMIT;