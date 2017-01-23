#!/usr/bin/env bash

set -u

DAD_JWT_TOKEN="$1"
DAD_URL=${2:-'http://localhost:8080'}
FILE=${3:-entities.csv}

(
    IFS=$';\n'
    while read -r name type; do
        fs=$(printf '{"name": "%s", "type": "%s"}' "$name" "$type")
        echo "Sending: $fs"
        curl -sH "Authorization:Bearer $DAD_JWT_TOKEN" -H 'Content-Type: application/json;charset=UTF-8' -d "$fs"  "$DAD_URL/api/entities/new"
        echo
    done < "$FILE"
)
