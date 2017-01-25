#!/usr/bin/env bash

set -u

DAD_JWT_TOKEN="$1"
DAD_URL=${2:-'http://localhost:8080'}
FILE=${3:-projects.csv}  # Expected CSV format: Name;Domain;Service Center;Business Unit;Description

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
