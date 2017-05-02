#!/usr/bin/env bash

# This script populates the "entities" collection of the database.
# Usage:
#   bash populate-entities.sh <JWT token> [URL] [CSV file]
#
# JWT Token: the authentication token of an admin account
# URL (optional, defaults to "http://localhost:8080"): the URL of a DAD instance
# CSV file (optional, defaults to "entities.csv"): the file containing the data to insert
#
# The CSV file is semi-colon (;) separated. The expected fields are:
# - Name: the name of the entity to insert
# - Type: the type of the entity to insert (ie. businessUnit, serviceCenter)

set -u

DAD_JWT_TOKEN="$1"
DAD_URL=${2:-'http://localhost:8080'}
FILE=${3:-entities.csv}

(
    IFS=$';\n'
    while read -r name type; do
        entity=$(printf '{"name": "%s", "type": "%s"}' "$name" "$type")
        echo "Sending: $entity"
        curl -sH "Authorization:Bearer $DAD_JWT_TOKEN" -H 'Content-Type: application/json;charset=UTF-8' -d "$entity" "$DAD_URL/api/entities/new"
        echo
    done < "$FILE"
)
