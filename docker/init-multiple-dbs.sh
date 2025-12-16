#!/bin/bash
set -e
psql -v ON_ERROR_STOP=1 --username "app" <<-EOSQL
  CREATE DATABASE auth_db;
  CREATE DATABASE user_db;
EOSQL
