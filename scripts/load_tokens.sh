#!/bin/bash


set -e

DEFAULT_FILE="routes/LogPas.txt"
TOKEN_FILE="${1:-$DEFAULT_FILE}"

echo "Loading tokens from: $TOKEN_FILE"

if [ ! -f "$TOKEN_FILE" ]; then
    echo "Error: File $TOKEN_FILE not found"
    exit 1
fi

echo "Tokens found:"
echo "============="

line_number=1
while IFS= read -r line; do
    if [ -n "$line" ]; then
        token=$(echo "$line" | awk '{print $1}')
        user_id=$(echo "$line" | awk '{print $2}')
        
        echo "Line $line_number:"
        echo "  Token: $token"
        echo "  User ID: $user_id"
        echo ""
    fi
    ((line_number++))
done < "$TOKEN_FILE"

echo "Total tokens loaded: $((line_number - 1))"
echo "============="
echo "Note: Tokens are automatically loaded when the service starts."
echo "You can also use the /auth endpoint to generate new tokens."
