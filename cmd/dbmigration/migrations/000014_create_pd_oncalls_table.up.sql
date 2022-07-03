BEGIN;

CREATE TABLE IF NOT EXISTS pagerduty_oncalls(
    _id INT GENERATED ALWAYS AS IDENTITY,
    pd_user_id VARCHAR(32) NOT NULL,
    pd_user_name VARCHAR(128) NOT NULL,
    schedule_id VARCHAR(32),
    schedule_name VARCHAR(256),
    escalation_policy_id VARCHAR(32) NOT NULL,
    escalation_policy_name VARCHAR(256) NOT NULL,
    escalation_level INT NOT NULL,
    on_call_start TIMESTAMP NOT NULL,
    on_call_end TIMESTAMP NOT NULL,

    CONSTRAINT pd_on_call_cnstr UNIQUE (pd_user_id, schedule_id, escalation_policy_id, escalation_level, on_call_start, on_call_end)
);

CREATE INDEX pd_on_call_start_idx ON pagerduty_oncalls(on_call_start);
CREATE INDEX pd_on_call_end_idx ON pagerduty_oncalls(on_call_end);
CREATE INDEX pd_on_call_user_name ON pagerduty_oncalls(pd_user_name);
CREATE INDEX pd_on_call_escalation_level ON pagerduty_oncalls(escalation_level);

COMMIT;
