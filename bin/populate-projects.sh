#!/usr/bin/env bash

# This script populates the "projects" collection of the database.
# Usage:
#   bash populate-projects.sh <JWT token> [URL] [CSV file]
#
# JWT Token: the authentication token of an admin account
# URL (optional, defaults to "http://localhost:8080"): the URL of a DAD instance
# CSV file (optional, defaults to "projects.csv"): the file containing the data to insert
#
# The CSV file is semi-colon (;) separated. The expected fields are:
# - Name: the name of the project to insert
# - Domain: the domain of the project
# - Service Center: the service center of the project
# - Business Unit: the business unit of the project
# - Description: a description of the project

set -u

DAD_JWT_TOKEN="$1"
DAD_URL=${2:-'http://localhost:8080'}
FILE=${3:-projects.csv}

entities=$(curl -sH "Authorization:Bearer $DAD_JWT_TOKEN" "$DAD_URL/api/entities")
(
    IFS=$';\n'
    while read -r name domain serviceCenter businessUnit description; do
        serviceCenterID=$(echo "$entities" | jq -M -r ".[] | select(.name==\"$serviceCenter\") | .id")
        businessUnitID=$(echo "$entities" | jq -M -r ".[] | select(.name==\"$businessUnit\") | .id")

        project=$(printf '{"name": "%s", "domain": "%s", "serviceCenter": "%s", "businessUnit": "%s", "description": "%s"}' "$name" "$domain" "$serviceCenterID" "$businessUnitID" "$description")
        echo "Sending: $project"
        curl -sH "Authorization:Bearer $DAD_JWT_TOKEN" -H 'Content-Type: application/json;charset=UTF-8' -d "$project" "$DAD_URL/api/projects/new"
        echo
    done < "$FILE"
)
