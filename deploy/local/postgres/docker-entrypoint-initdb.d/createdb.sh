#!/bin/bash

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER default WITH PASSWORD 'secret';
    CREATE DATABASE meetings;
    GRANT ALL PRIVILEGES ON DATABASE meetings TO default;
EOSQL
