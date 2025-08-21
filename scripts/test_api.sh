#!/bin/bash


echo "🚀 Testing Analytics Service API"
echo "=================================="

# Base URL
BASE_URL="http://localhost:8080"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' 

print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
    fi
}

# Test 1: Health check
echo -e "\n${YELLOW}1. Testing health check...${NC}"
curl -s "$BASE_URL/" | jq .
print_status $? "Health check"

# Test 2: Generate token
echo -e "\n${YELLOW}2. Generating token...${NC}"
TOKEN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth" \
    -H "Content-Type: application/json" \
    -d '{"email": "string", "password": "string"}')

echo "$TOKEN_RESPONSE" | jq .

# Extract token
TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.token')
if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    print_status 0 "Token generated successfully"
    echo "Token: $TOKEN"
else
    print_status 1 "Failed to generate token"
    exit 1
fi

# Test 3: Validate token
echo -e "\n${YELLOW}3. Validating token...${NC}"
curl -s "$BASE_URL/validate?token=$TOKEN" | jq .
print_status $? "Token validation"

# Test 4: Analytics with generated token
echo -e "\n${YELLOW}4. Testing analytics endpoint...${NC}"
ANALYTICS_RESPONSE=$(curl -s -X POST "$BASE_URL/analytics" \
    -H "Content-Type: application/json" \
    -d "{
        \"token\": \"$TOKEN\",
        \"StartDate\": \"01.01.2024\",
        \"FinishDate\": \"31.01.2024\"
    }")

echo "$ANALYTICS_RESPONSE" | jq .

# Check if analytics returned data
ITEMS_COUNT=$(echo "$ANALYTICS_RESPONSE" | jq '.items | length')
TOTAL=$(echo "$ANALYTICS_RESPONSE" | jq '.total')

if [ "$ITEMS_COUNT" -gt 0 ] && [ "$TOTAL" -gt 0 ]; then
    print_status 0 "Analytics returned data successfully"
    echo "Items count: $ITEMS_COUNT"
    echo "Total: $TOTAL"
else
    print_status 1 "Analytics returned empty data"
    echo "Items count: $ITEMS_COUNT"
    echo "Total: $TOTAL"
fi

# Test 5: Analytics with invalid token
echo -e "\n${YELLOW}5. Testing analytics with invalid token...${NC}"
curl -s -X POST "$BASE_URL/analytics" \
    -H "Content-Type: application/json" \
    -d '{
        "token": "invalid-token",
        "StartDate": "01.01.2024",
        "FinishDate": "31.01.2024"
    }' | jq .
print_status $? "Invalid token test"

echo -e "\n${GREEN}🎉 API testing completed!${NC}"
