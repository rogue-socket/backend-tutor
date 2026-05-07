#!/bin/bash
# Creates a separate test database so integration tests don't pollute dev data.
set -euo pipefail
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE links_test;
    GRANT ALL PRIVILEGES ON DATABASE links_test TO $POSTGRES_USER;
EOSQL
