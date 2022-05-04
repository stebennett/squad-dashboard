CREATE TABLE jira_issues_reports(
    _id INT GENERATED ALWAYS AS IDENTITY,
    project VARCHAR(16) NOT NULL,
    week_start TIMESTAMP NOT NULL,
    number_of_items_started INT NOT NULL DEFAULT -1,
    number_of_items_completed INT NOT NULL DEFAULT -1,

    CONSTRAINT project_week_start_c UNIQUE (project, week_start)
);
