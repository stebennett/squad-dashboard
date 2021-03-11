CREATE TABLE IF NOT EXISTS jira_config (
    _id BIGINT GENERATED ALWAYS AS IDENTITY,
    project_key VARCHAR NOT NULL UNIQUE,
    work_start_state VARCHAR NOT NULL,
    PRIMARY KEY(_id)
);

CREATE TABLE IF NOT EXISTS flow_measures (
    _id BIGINT GENERATED ALWAYS AS IDENTITY,
    jira_data_id INT NOT NULL,
    cycle_time INT NOT NULL DEFAULT -1,
    lead_time INT NOT NULL DEFAULT -1,

    PRIMARY KEY(_id),
    CONSTRAINT jira_data_id_fk FOREIGN KEY(jira_data_id) REFERENCES jira_data(_id) ON DELETE CASCADE
);

ALTER TABLE jira_data ADD COLUMN IF NOT EXISTS
    jira_project_key VARCHAR NOT NULL DEFAULT '';

CREATE FUNCTION fn_calculate_dev_start_time() RETURNS trigger AS $$
    DECLARE
        work_start_date timestamp;
    BEGIN
        SELECT work_start_state INTO work_start_date FROM jira_config
        WHERE project_key = (SELECT jira_project_key FROM jira_data WHERE _id = NEW.jira_data_id);

        SELECT jira_transition_at INTO work_start_date FROM jira_transitions
        WHERE jira_transition_to = work_start_state
        AND jira_data_id = NEW.jira_data_id
        ORDER BY jira_transition_at DESC
        LIMIT 1;

        UPDATE jira_data
        SET jira_work_started_at = work_start_date
        WHERE _id = NEW.jira_data_id;

        RETURN NEW;
    END
$$ LANGUAGE plpgsql;

CREATE FUNCTION fn_calculate_cycle_time() RETURNS trigger AS $$
    DECLARE
        cycle_time_val int;
    BEGIN
        SELECT DATE_PART('day', jira_completed_at - jira_work_start_at) + 1 INTO cycle_time_val FROM jira_data
        WHERE _id = NEW._id;

        INSERT INTO flow_measures(jira_data_id, cycle_time)
        VALUES (NEW._id, cycle_time_val) ON CONFLICT
        DO UPDATE SET cycle_time = cycle_time_val WHERE jira_data_id = NEW._id;

        RETURN NEW;
    END
$$ LANGUAGE plpgsql;

CREATE FUNCTION fn_calculate_lead_time() RETURNS trigger AS $$
    DECLARE
        lead_time_val int;
    BEGIN
        SELECT DATE_PART('day', jira_completed_at - jira_created_at) + 1 INTO lead_time_val FROM jira_data
        WHERE _id = NEW._id;

        INSERT INTO flow_measures(jira_data_id, lead_time_val)
        VALUES (NEW._id, lead_time_val) ON CONFLICT
        DO UPDATE SET lead_time = lead_time_val WHERE jira_data_id = NEW._id;

        RETURN NEW;
    END
$$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_jira_data_cycle_time
    AFTER INSERT OR UPDATE ON jira_data
    FOR EACH ROW
    WHEN (NEW.jira_work_started_at IS NOT NULL AND NEW.jira_completed_at IS NOT NULL)
    EXECUTE PROCEDURE fn_calculate_cycle_time();

CREATE TRIGGER tgr_jira_data_lead_time
    AFTER INSERT OR UPDATE ON jira_data
    FOR EACH ROW
    WHEN (NEW.jira_created_at IS NOT NULL AND NEW.jira_completed_at IS NOT NULL)
    EXECUTE PROCEDURE fn_calculate_lead_time();

CREATE TRIGGER tgr_jira_status_updated
    AFTER INSERT OR UPDATE ON jira_transitions
    FOR EACH ROW EXECUTE PROCEDURE fn_calculate_dev_start_time()
