FROM boxfuse/flyway:latest

RUN apt-get update && apt-get -y install postgresql-client

COPY wait_for_postgres /scripts/wait_for_postgres

RUN chmod +x /scripts/wait_for_postgres

WORKDIR /flyway

ENTRYPOINT ["/scripts/wait_for_postgres"]

CMD ["--help"]
