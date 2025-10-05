#!/bin/bash

# Wrapper script to run Newman tests with automatic token generation

newman run collection/fitapi.postman_collection.json \
  --env-var "auth_token=$(./collection/generate-token.sh)" \
  "$@"
