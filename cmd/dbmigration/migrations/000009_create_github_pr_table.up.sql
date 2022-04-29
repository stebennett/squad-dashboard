CREATE TABLE github_pull_requests(
    _id INT GENERATED ALWAYS AS IDENTITY,
    organisation VARCHAR(64) NOT NULL,
    repository TEXT NOT NULL,
    gh_user VARCHAR(128),
    title TEXT NOT NULL,
    github_id INT NOT NULL,
    pr_number INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    closed_at TIMESTAMP,
    merged_at TIMESTAMP,

    CONSTRAINT org_repo_number_c UNIQUE (organisation, repository, pr_number)
);
