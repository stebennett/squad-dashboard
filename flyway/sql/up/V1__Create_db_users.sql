CREATE USER ${user_writer_name} WITH ENCRYPTED PASSWORD '${user_writer_password}';

GRANT CONNECT ON DATABASE ${database_name} TO ${user_writer_name};
GRANT USAGE ON SCHEMA public TO ${user_writer_name};

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO ${user_writer_name};
GRANT SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO ${user_writer_name};

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO ${user_writer_name};
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, UPDATE ON SEQUENCES TO ${user_writer_name};
