# Plummet

Plummet allows running SQL for transforming data using [DuckDB](https://duckdb.org). It uses a simple Makefile like pattern where there are targets with SQL to run that can be dependent on other targets.

## Why do this?

There are a lot of very [complex](https://docs.dagster.io/getting-started/quickstart#understanding-the-code) [pipeline](https://docs.prefect.io/latest/tutorial/flows/#run-your-first-flow) [descriptions](https://airflow.apache.org/docs/apache-airflow/stable/index.html#workflows-as-code) out there for doing ETLs. The goal is to allow a single file to define a simple pipeline that doesn't need any prior knowledge of a complex framework.

## Why DuckDB?

DuckDB has some nice features that make it well suited for this sort of process.

1. It has support for writing output to files in different [formats](https://duckdb.org/docs/sql/statements/copy#format-specific-options).
1. It has support for connecting to a large number of databases and importing from files.
1. It is optimized for fast analytic queries with zero dependencies and fits well into larger data systems.

## How does it work?

The pipeline is configured with a `plummet.yml`. Here is example.

```yml
targets:
  setup_postgres:
    sql: |
      INSTALL POSTGRES;
      LOAD POSTGRES;
  get_events:
    config:
      PG_DBURI: "postgres://timescaledb:password@localhost/userup_userservice"
    sql: "ATTACH '{{ .PG_DBURI }}' as db (TYPE postgres, READ_ONLY);"
    deps: ["setup_postgres"]
  save_events:
    sql: "COPY db.events TO 'events.parquet' (FORMAT PARQUET);"
    deps: ["get_events"]
```

The `targets` defines the different targets you can run. Each target can have a few different options.

- `sql` is the SQL code you want to run. DuckDB exposes a lot of powerful capabilities, with the simple example above showing how to grab data from a Postgres DB and write a Parquet file.
- `deps` is a list of targets that are required before this target can be run.
- `config` allows you to set variables that can be used in the SQL.

## Why "plummet"?

Pipeline are "plumbing", so "plummet" sounds like "plumb it".
