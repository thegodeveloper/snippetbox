#!/usr/bin/env bash
set -x
set -eo pipefail

if ! [ -x "$(command -v psql)" ]; then
  echo >&2 "Error: psql is not installed."
  exit 1
fi

DB_USER="${POSTGRES_USER:=postgres}"
DB_PASSWORD="${POSTGRES_PASSWORD:=password}"
DB_NAME="${POSTGRES_DB:=snippetbox}"
DB_PORT="${POSTGRES_PORT:=5432}"
DB_HOST="${POSTGRES_HOST:=localhost}"

# allow to skip Docker if a dockerized postgres database is already running
if [[ -z "${SKIP_DOCKER}" ]]
then
  docker run \
    -e POSTGRES_USER=${DB_USER} \
    -e POSTGRES_PASSWORD=${DB_PASSWORD} \
    -e POSTGRES_DB=${DB_NAME} \
    -p "${DB_PORT}":5432 \
    -d postgres \
    postgres -N 1000
    # ^ increased maximum number of connections for testing purposes
fi

# keep pinging postgres until it is ready to accept commands
export PGPASSWORD="${DB_PASSWORD}"
until psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d "postgres" -c '\q'; do
  >&2 echo "postgres is still unavailable - sleeping"
  sleep 1
done

>&2 echo "postgres is up and running on port ${DB_PORT} - running migrations now!"

# Create the appadmin user, grant privileges, and make it a superuser
psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d "postgres" <<-EOSQL
  CREATE USER appadmin WITH ENCRYPTED PASSWORD 'password';
  GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO appadmin;
  GRANT ALL ON SCHEMA public TO appadmin;
  ALTER ROLE appadmin SUPERUSER;
EOSQL

# Connect to the snippetbox database and create tables
psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d "${DB_NAME}" <<-EOSQL
  CREATE TABLE public.snippets (
    id serial NOT NULL,
    title varchar NOT NULL,
    "content" text NOT NULL,
    created timestamp NOT NULL,
    expires timestamp NOT NULL
  );

  CREATE INDEX idx_snippets_created ON snippets(created);

  CREATE TABLE public.users (
    id serial NOT NULL,
    name varchar NOT NULL,
    email varchar NOT NULL,
    hashed_password char(60) NOT NULL,
    created timestamp NOT NULL,
    active boolean NOT NULL DEFAULT TRUE
  );

  ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
EOSQL


>&2 echo "postgres has been migrated, ready to go!"