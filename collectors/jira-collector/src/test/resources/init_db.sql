--
-- PostgreSQL database dump
--

-- Dumped from database version 12.3 (Debian 12.3-1.pgdg100+1)
-- Dumped by pg_dump version 13.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: ingestion_type; Type: TYPE; Schema:  Owner: db_root_user
--

CREATE TYPE ingestion_type AS ENUM (
    'backfill',
    'incremental'
);


ALTER TYPE ingestion_type OWNER TO db_root_user;

--
-- Name: work_type; Type: TYPE; Schema:  Owner: db_root_user
--

CREATE TYPE work_type AS ENUM (
    'story',
    'task',
    'bug',
    'subtask'
);


ALTER TYPE work_type OWNER TO db_root_user;

--
-- Name: fn_calc_cycle_time(); Type: FUNCTION; Schema:  Owner: db_root_user
--

CREATE FUNCTION fn_calc_cycle_time() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    cycle_time_val int;
BEGIN
    SELECT DATE_PART('day', jira_completed_at - jira_work_started_at) + 1 INTO cycle_time_val FROM jira_data
    WHERE _id = NEW._id;

    INSERT INTO flow_measures(jira_data_id, cycle_time)
    VALUES (NEW._id, cycle_time_val) ON CONFLICT
        DO UPDATE SET cycle_time = cycle_time_val WHERE jira_data_id = NEW._id;

    RETURN NEW;
END
$$;


ALTER FUNCTION fn_calc_cycle_time() OWNER TO db_root_user;

--
-- Name: fn_calc_dev_start_time(); Type: FUNCTION; Schema:  Owner: db_root_user
--

CREATE FUNCTION fn_calc_dev_start_time() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    work_start_state text;
    work_start_date timestamp;
BEGIN
    SELECT work_start_state INTO work_start_state FROM jira_config
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
$$;


ALTER FUNCTION fn_calc_dev_start_time() OWNER TO db_root_user;

--
-- Name: fn_calc_lead_time(); Type: FUNCTION; Schema:  Owner: db_root_user
--

CREATE FUNCTION fn_calc_lead_time() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    lead_time_val int;
BEGIN
    SELECT DATE_PART('day', jira_completed_at - jira_created_at) + 1 INTO lead_time_val FROM jira_data
    WHERE _id = NEW._id;

    INSERT INTO flow_measures(jira_data_id, lead_time)
    VALUES (NEW._id, lead_time_val) ON CONFLICT
        DO UPDATE SET lead_time = lead_time_val WHERE jira_data_id = NEW._id;

    RETURN NEW;
END
$$;


ALTER FUNCTION fn_calc_lead_time() OWNER TO db_root_user;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: flow_measures; Type: TABLE; Schema:  Owner: db_root_user
--

CREATE TABLE flow_measures (
    _id bigint NOT NULL,
    jira_data_id integer NOT NULL,
    cycle_time integer DEFAULT '-1'::integer NOT NULL,
    lead_time integer DEFAULT '-1'::integer NOT NULL
);


ALTER TABLE flow_measures OWNER TO db_root_user;

--
-- Name: flow_measures__id_seq; Type: SEQUENCE; Schema:  Owner: db_root_user
--

ALTER TABLE flow_measures ALTER COLUMN _id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME flow_measures__id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: flyway_schema_history; Type: TABLE; Schema:  Owner: db_root_user
--

CREATE TABLE flyway_schema_history (
    installed_rank integer NOT NULL,
    version character varying(50),
    description character varying(200) NOT NULL,
    type character varying(20) NOT NULL,
    script character varying(1000) NOT NULL,
    checksum integer,
    installed_by character varying(100) NOT NULL,
    installed_on timestamp without time zone DEFAULT now() NOT NULL,
    execution_time integer NOT NULL,
    success boolean NOT NULL
);


ALTER TABLE flyway_schema_history OWNER TO db_root_user;

--
-- Name: jira_config; Type: TABLE; Schema:  Owner: db_root_user
--

CREATE TABLE jira_config (
    _id bigint NOT NULL,
    project_key character varying NOT NULL,
    work_start_state character varying NOT NULL,
    last_ingestion_run_started timestamp without time zone,
    last_ingestion_run_completed timestamp without time zone,
    last_ingestion_type ingestion_type
);


ALTER TABLE jira_config OWNER TO db_root_user;

--
-- Name: jira_config__id_seq; Type: SEQUENCE; Schema:  Owner: db_root_user
--

ALTER TABLE jira_config ALTER COLUMN _id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME jira_config__id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: jira_data; Type: TABLE; Schema:  Owner: db_root_user
--

CREATE TABLE jira_data (
    _id bigint NOT NULL,
    jira_id integer NOT NULL,
    jira_key character varying NOT NULL,
    jira_work_type work_type NOT NULL,
    jira_created_at timestamp without time zone NOT NULL,
    jira_completed_at timestamp without time zone,
    jira_work_started_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    jira_project_key character varying DEFAULT ''::character varying NOT NULL
);


ALTER TABLE jira_data OWNER TO db_root_user;

--
-- Name: jira_data__id_seq; Type: SEQUENCE; Schema:  Owner: db_root_user
--

ALTER TABLE jira_data ALTER COLUMN _id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME jira_data__id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: jira_transitions; Type: TABLE; Schema:  Owner: db_root_user
--

CREATE TABLE jira_transitions (
    _id bigint NOT NULL,
    jira_data_id bigint NOT NULL,
    jira_id integer NOT NULL,
    jira_transition_to character varying NOT NULL,
    jira_transition_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


ALTER TABLE jira_transitions OWNER TO db_root_user;

--
-- Name: jira_transitions__id_seq; Type: SEQUENCE; Schema:  Owner: db_root_user
--

ALTER TABLE jira_transitions ALTER COLUMN _id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME jira_transitions__id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: flow_measures flow_measures_pkey; Type: CONSTRAINT; Schema:  Owner: db_root_user
--

ALTER TABLE ONLY flow_measures
    ADD CONSTRAINT flow_measures_pkey PRIMARY KEY (_id);


--
-- Name: flyway_schema_history flyway_schema_history_pk; Type: CONSTRAINT; Schema:  Owner: db_root_user
--

ALTER TABLE ONLY flyway_schema_history
    ADD CONSTRAINT flyway_schema_history_pk PRIMARY KEY (installed_rank);


--
-- Name: jira_config jira_config_pkey; Type: CONSTRAINT; Schema:  Owner: db_root_user
--

ALTER TABLE ONLY jira_config
    ADD CONSTRAINT jira_config_pkey PRIMARY KEY (_id);


--
-- Name: jira_config jira_config_project_key_key; Type: CONSTRAINT; Schema:  Owner: db_root_user
--

ALTER TABLE ONLY jira_config
    ADD CONSTRAINT jira_config_project_key_key UNIQUE (project_key);


--
-- Name: jira_data jira_data_jira_id_key; Type: CONSTRAINT; Schema:  Owner: db_root_user
--

ALTER TABLE ONLY jira_data
    ADD CONSTRAINT jira_data_jira_id_key UNIQUE (jira_id);


--
-- Name: jira_data jira_data_jira_key_key; Type: CONSTRAINT; Schema:  Owner: db_root_user
--

ALTER TABLE ONLY jira_data
    ADD CONSTRAINT jira_data_jira_key_key UNIQUE (jira_key);


--
-- Name: jira_data jira_data_pkey; Type: CONSTRAINT; Schema:  Owner: db_root_user
--

ALTER TABLE ONLY jira_data
    ADD CONSTRAINT jira_data_pkey PRIMARY KEY (_id);


--
-- Name: jira_transitions jira_transitions_jira_id_key; Type: CONSTRAINT; Schema:  Owner: db_root_user
--

ALTER TABLE ONLY jira_transitions
    ADD CONSTRAINT jira_transitions_jira_id_key UNIQUE (jira_id);


--
-- Name: jira_transitions jira_transitions_pkey; Type: CONSTRAINT; Schema:  Owner: db_root_user
--

ALTER TABLE ONLY jira_transitions
    ADD CONSTRAINT jira_transitions_pkey PRIMARY KEY (_id);


--
-- Name: flyway_schema_history_s_idx; Type: INDEX; Schema:  Owner: db_root_user
--

CREATE INDEX flyway_schema_history_s_idx ON flyway_schema_history USING btree (success);


--
-- Name: jira_data_jira_id_idx; Type: INDEX; Schema:  Owner: db_root_user
--

CREATE INDEX jira_data_jira_id_idx ON jira_data USING btree (jira_id);


--
-- Name: jira_data_jira_key_idx; Type: INDEX; Schema:  Owner: db_root_user
--

CREATE INDEX jira_data_jira_key_idx ON jira_data USING btree (jira_key);


--
-- Name: jira_data_jira_work_type_idx; Type: INDEX; Schema:  Owner: db_root_user
--

CREATE INDEX jira_data_jira_work_type_idx ON jira_data USING btree (jira_work_type);


--
-- Name: jira_transitions_jira_id_idx; Type: INDEX; Schema:  Owner: db_root_user
--

CREATE INDEX jira_transitions_jira_id_idx ON jira_transitions USING btree (jira_id);


--
-- Name: jira_transitions_jira_transition_at_idx; Type: INDEX; Schema:  Owner: db_root_user
--

CREATE INDEX jira_transitions_jira_transition_at_idx ON jira_transitions USING btree (jira_transition_at);


--
-- Name: jira_transitions_jira_transition_to_idx; Type: INDEX; Schema:  Owner: db_root_user
--

CREATE INDEX jira_transitions_jira_transition_to_idx ON jira_transitions USING btree (jira_transition_to);


--
-- Name: jira_data tgr_jira_data_calc_cycle_time; Type: TRIGGER; Schema:  Owner: db_root_user
--

CREATE TRIGGER tgr_jira_data_calc_cycle_time AFTER INSERT OR UPDATE ON jira_data FOR EACH ROW WHEN (((new.jira_work_started_at IS NOT NULL) AND (new.jira_completed_at IS NOT NULL))) EXECUTE FUNCTION fn_calc_cycle_time();


--
-- Name: jira_data tgr_jira_data_calc_lead_time; Type: TRIGGER; Schema:  Owner: db_root_user
--

CREATE TRIGGER tgr_jira_data_calc_lead_time AFTER INSERT OR UPDATE ON jira_data FOR EACH ROW WHEN (((new.jira_created_at IS NOT NULL) AND (new.jira_completed_at IS NOT NULL))) EXECUTE FUNCTION fn_calc_lead_time();


--
-- Name: jira_transitions tgr_jira_status_update; Type: TRIGGER; Schema:  Owner: db_root_user
--

CREATE TRIGGER tgr_jira_status_update AFTER INSERT OR UPDATE ON jira_transitions FOR EACH ROW EXECUTE FUNCTION fn_calc_dev_start_time();


--
-- Name: jira_transitions jira_data_id_fk; Type: FK CONSTRAINT; Schema:  Owner: db_root_user
--

ALTER TABLE ONLY jira_transitions
    ADD CONSTRAINT jira_data_id_fk FOREIGN KEY (jira_data_id) REFERENCES jira_data(_id) ON DELETE CASCADE;


--
-- Name: flow_measures jira_data_id_fk; Type: FK CONSTRAINT; Schema:  Owner: db_root_user
--

ALTER TABLE ONLY flow_measures
    ADD CONSTRAINT jira_data_id_fk FOREIGN KEY (jira_data_id) REFERENCES jira_data(_id) ON DELETE CASCADE;
