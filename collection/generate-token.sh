#!/bin/bash

# Script to generate and cache authentication token for Newman tests

CACHE_FILE=".token_cache.json"
BUFFER_SECONDS=300  # 5 minutes

# Check if cached token exists and is still valid
if [ -f "$CACHE_FILE" ]; then
    EXPIRES_AT=$(jq -r '.expires_at' "$CACHE_FILE")
    NOW=$(date +%s)
    TIME_LEFT=$((EXPIRES_AT - NOW - BUFFER_SECONDS))

    if [ $TIME_LEFT -gt 0 ]; then
        MINUTES_LEFT=$((TIME_LEFT / 60))
        echo "✅ Using cached token (expires in $MINUTES_LEFT minutes)" >&2
        jq -r '.access_token' "$CACHE_FILE"
        exit 0
    else
        echo "⏰ Cached token expired, generating new one..." >&2
    fi
fi

# Generate new token with --json flag for machine-readable output
echo "🔄 Generating new authentication token..." >&2
TOKEN_JSON=$(go run cmd/gettoken/main.go --json 2>&1)

if [ $? -ne 0 ]; then
    echo "❌ Failed to generate token" >&2
    exit 1
fi

# Save to cache
echo "$TOKEN_JSON" > "$CACHE_FILE"

# Extract and return just the access token
ACCESS_TOKEN=$(echo "$TOKEN_JSON" | jq -r '.access_token')
EXPIRES_IN=$(echo "$TOKEN_JSON" | jq -r '.expires_in')

echo "✅ New token generated and cached (expires in ${EXPIRES_IN}s)" >&2
echo "$ACCESS_TOKEN"
