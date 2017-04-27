#!/usr/bin/env bash

# This script populates the "technologies" collection of the database.
# Usage:
#   bash populate-technologies.sh <JWT token> [URL] [CSV file]
#
# JWT Token: the authentication token of an admin account
# URL (optional, defaults to "http://localhost:8080"): the URL of a DAD instance
# CSV file (optional, defaults to "technologies.csv"): the file containing the data to insert
#
# The CSV file is semi-colon (;) separated. The expected fields are:
# - Name: the name of the technology to insert

set -u

DAD_JWT_TOKEN="$1"
DAD_URL=${2:-'http://localhost:8080'}
FILE=${3:-technologies.csv}

(
    IFS=$';\n'
    while read -r name; do
        technology=$(printf '{"name": "%s"}' "$name")
        echo "Sending: $technology"
        curl -sH "Authorization:Bearer $DAD_JWT_TOKEN" -H 'Content-Type: application/json;charset=UTF-8' -d "$technology" "$DAD_URL/api/projects/new"
        echo
    done < "$FILE"
)
