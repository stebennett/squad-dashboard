BEGIN;

CREATE TABLE IF NOT EXISTS non_working_days(
    _id INT GENERATED ALWAYS AS IDENTITY,
    project VARCHAR(16) NOT NULL,
    non_working_day DATE NOT NULL,
    CONSTRAINT project_non_working_day UNIQUE (project, non_working_day)
);

CREATE INDEX project_idx ON non_working_days(project);
CREATE INDEX project_non_working_day_idx ON non_working_days(project, non_working_day);

COMMIT;