#!/usr/bin/env bash

# This script populates the "functionalServices" collection of the database.
# Usage:
#   bash populate-functional-services.sh <JWT token> [URL] [CSV file]
#
# JWT Token: the authentication token of an admin account
# URL (optional, defaults to "http://localhost:8080"): the URL of a DAD instance
# CSV file (optional, defaults to "functional-services.csv"): the file containing the data to insert
#
# The CSV file is semi-colon (;) separated. The expected fields are:
# - Name: the name of the functional service to insert
# - Package: the package of the functional service (ie. Build, OPS, etc.)
# - Position: the relative position of the package on display (a service with position "10" will be displayed before another service with position "20")

set -u

DAD_JWT_TOKEN="$1"
DAD_URL=${2:-'http://localhost:8080'}
FILE=${3:-functional-services.csv}

(
    IFS=$';\n'
    while read -r package name position; do
        fs=$(printf '{"name": "%s", "package": "%s", "position": %s}' "$name" "$package" "$position")
        echo "Sending: $fs"
        curl -sH "Authorization:Bearer $DAD_JWT_TOKEN" -H 'Content-Type: application/json;charset=UTF-8' -d "$fs" "$DAD_URL/api/services/new"
        echo
    done < "$FILE"
)
