#!/usr/bin/env bash

secret=${1}
target=${2}

ids=$(
  curl -sSf \
    -X POST \
    -H "Authorization: Bearer ${secret}" \
    -H 'Content-type: application/json;charset=utf-8' \
    --data "{\"channel\": \"${target}\"}" \
    https://slack.com/api/conversations.history |
    jq '.messages[].ts'
)

for msg in ${ids}; do
  curl -sSf \
    -X POST \
    -H "Authorization: Bearer ${secret}" \
    -H 'Content-type: application/json;charset=utf-8' \
    --data "{\"channel\": \"${target}\", \"ts\": ${msg}}" \
    https://slack.com/api/chat.delete
  sleep 1
done
