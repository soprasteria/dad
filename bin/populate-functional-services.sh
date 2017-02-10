#!/usr/bin/env bash

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
