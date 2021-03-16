DROP TRIGGER IF EXISTS tgr_jira_status_updated ON jira_transitions;
DROP TRIGGER IF EXISTS tgr_jira_data_lead_time ON jira_data;
DROP TRIGGER IF EXISTS tgr_jira_data_cycle_time ON jira_data;

DROP FUNCTION IF EXISTS fn_calculate_lead_time;
DROP FUNCTION IF EXISTS fn_calculate_cycle_time;
DROP FUNCTION IF EXISTS fn_calculate_dev_start_time;

CREATE FUNCTION fn_calc_dev_start_time() RETURNS trigger
    LANGUAGE plpgsql
AS $$
DECLARE
    wss text;
    wsd timestamp;
BEGIN
    SELECT work_start_state INTO wss FROM jira_config
    WHERE project_key = (SELECT jira_project_key FROM jira_data WHERE _id = NEW.jira_data_id);

    SELECT jira_transition_at INTO wsd FROM jira_transitions
    WHERE jira_transition_to = wss
      AND jira_data_id = NEW.jira_data_id
    ORDER BY jira_transition_at DESC
    LIMIT 1;

    UPDATE jira_data
    SET jira_work_started_at = wsd
    WHERE _id = NEW.jira_data_id;

    RETURN NEW;
END
$$;

CREATE FUNCTION fn_calc_cycle_time() RETURNS trigger
    LANGUAGE plpgsql
AS $$
DECLARE
    cycle_time_val int;
BEGIN
    SELECT DATE_PART('day', jira_completed_at - jira_work_started_at) + 1 INTO cycle_time_val FROM jira_data
    WHERE _id = NEW._id;

    INSERT INTO flow_measures(jira_data_id, cycle_time)
    VALUES (NEW._id, cycle_time_val) ON CONFLICT(jira_data_id)
        DO UPDATE SET cycle_time = cycle_time_val WHERE flow_measures.jira_data_id = NEW._id;

    RETURN NEW;
END
$$;

CREATE FUNCTION fn_calc_lead_time() RETURNS trigger
    LANGUAGE plpgsql
AS $$
DECLARE
    lead_time_val int;
BEGIN
    SELECT DATE_PART('day', jira_completed_at - jira_created_at) + 1 INTO lead_time_val FROM jira_data
    WHERE _id = NEW._id;

    INSERT INTO flow_measures(jira_data_id, lead_time)
    VALUES (NEW._id, lead_time_val) ON CONFLICT(jira_data_id)
        DO UPDATE SET lead_time = lead_time_val WHERE flow_measures.jira_data_id = NEW._id;

    RETURN NEW;
END
$$;

CREATE TRIGGER tgr_jira_status_update
    AFTER INSERT OR UPDATE ON jira_transitions
    FOR EACH ROW EXECUTE PROCEDURE fn_calc_dev_start_time();

CREATE TRIGGER tgr_jira_data_calc_cycle_time
    AFTER INSERT OR UPDATE ON jira_data
    FOR EACH ROW
    WHEN (NEW.jira_work_started_at IS NOT NULL AND NEW.jira_completed_at IS NOT NULL)
EXECUTE PROCEDURE fn_calc_cycle_time();

CREATE TRIGGER tgr_jira_data_calc_lead_time
    AFTER INSERT OR UPDATE ON jira_data
    FOR EACH ROW
    WHEN (NEW.jira_created_at IS NOT NULL AND NEW.jira_completed_at IS NOT NULL)
EXECUTE PROCEDURE fn_calc_lead_time();
