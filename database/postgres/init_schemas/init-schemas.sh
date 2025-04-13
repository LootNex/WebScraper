#!/bin/bash

SCHEMAS=("auth")

for SCHEMA in "${SCHEMAS[@]}"; do
    psql "postgresql://${DBUSER}:${DBPASSWORD}@${DBHOST}:${DBPORT}/${DBNAME}" \
    -c "CREATE SCHEMA IF NOT EXISTS ${SCHEMA};"
done