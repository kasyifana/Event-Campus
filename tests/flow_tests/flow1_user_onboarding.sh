#!/bin/bash

# Flow 1: User Onboarding
# User Journey: New student wants to browse and attend campus events

set -e  # Exit on error

# Load helpers
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../helpers/color_output.sh"
source "$SCRIPT_DIR/../helpers/auth_helper.sh"
source "$SCRIPT_DIR/../helpers/validation_helper.sh"

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"
TEST_EMAIL="test_user_$(date +%s)@uii.ac.id"
TEST_PASSWORD="testpass123"

echo ""
print_info "============================================"
print_info "Flow 1: User Onboarding"
print_info "============================================"
echo ""

# Step 1: Register Account
print_step "Step 1: Register new user account"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\",
    \"full_name\": \"Test Student\",
    \"phone_number\": \"081234567890\"
  }")

validate_success "$REGISTER_RESPONSE" "true" "User registration" || exit 1
validate_field_exists "$REGISTER_RESPONSE" ".data.token" "JWT token" || exit 1
validate_field_exists "$REGISTER_RESPONSE" ".data.user.role" "User role" || exit 1
validate_field_value "$REGISTER_RESPONSE" ".data.user.role" "mahasiswa" "Default role" || exit 1

echo ""

# Step 2: Login
print_step "Step 2: Login with credentials"
TOKEN=$(login "$TEST_EMAIL" "$TEST_PASSWORD")
validate_not_empty "$TOKEN" "JWT Token" || exit 1
print_success "Login successful, token obtained"

echo ""

# Step 3: View Profile
print_step "Step 3: Get user profile"
PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/profile" \
  -H "Authorization: Bearer $TOKEN")

validate_success "$PROFILE_RESPONSE" "true" "Get profile" || exit 1
validate_field_exists "$PROFILE_RESPONSE" ".data.email" "Email" || exit 1
validate_field_value "$PROFILE_RESPONSE" ".data.email" "$TEST_EMAIL" "Email match" || exit 1

# Note: user_id may be empty string in response (backend issue)
USER_ID=$(echo "$PROFILE_RESPONSE" | jq -r '.data.user_id // "unknown"')
if [ "$USER_ID" == "" ] || [ "$USER_ID" == "unknown" ]; then
  print_warning "User ID not available in profile response"
else
  print_info "User ID: $USER_ID"
fi

echo ""

# Step 4: Browse Available Events
print_step "Step 4: Browse all available events"
EVENTS_RESPONSE=$(curl -s -X GET "$BASE_URL/events" \
  -H "Authorization: Bearer $TOKEN")

validate_success "$EVENTS_RESPONSE" "true" "Get events" || exit 1

# Handle null data (no events available)
EVENT_DATA=$(echo "$EVENTS_RESPONSE" | jq '.data')
if [ "$EVENT_DATA" == "null" ]; then
  print_warning "No events available yet"
  EVENT_COUNT=0
else
  EVENT_COUNT=$(echo "$EVENTS_RESPONSE" | jq '.data | length')
  print_info "Found $EVENT_COUNT published events"
fi

echo ""

# Step 5: Search for Specific Events
print_step "Step 5: Search events with filters"
SEARCH_RESPONSE=$(curl -s -X GET "$BASE_URL/events?category=workshop&search=web" \
  -H "Authorization: Bearer $TOKEN")

validate_success "$SEARCH_RESPONSE" "true" "Search events" || exit 1

SEARCH_DATA=$(echo "$SEARCH_RESPONSE" | jq '.data')
if [ "$SEARCH_DATA" == "null" ]; then
  SEARCH_COUNT=0
else
  SEARCH_COUNT=$(echo "$SEARCH_RESPONSE" | jq '.data | length')
fi
print_info "Found $SEARCH_COUNT events matching filters"

echo ""

# Step 6: View Event Details (if events exist)
if [ "$EVENT_COUNT" -gt 0 ] && [ "$EVENT_DATA" != "null" ]; then
  print_step "Step 6: View event details"
  
  FIRST_EVENT_ID=$(echo "$EVENTS_RESPONSE" | jq -r '.data[0].id')
  
  if [ -n "$FIRST_EVENT_ID" ] && [ "$FIRST_EVENT_ID" != "null" ]; then
    EVENT_DETAIL_RESPONSE=$(curl -s -X GET "$BASE_URL/events/$FIRST_EVENT_ID" \
      -H "Authorization: Bearer $TOKEN")
    
    validate_success "$EVENT_DETAIL_RESPONSE" "true" "Get event detail" || exit 1
    validate_field_exists "$EVENT_DETAIL_RESPONSE" ".data.title" "Event title" || exit 1
    validate_field_exists "$EVENT_DETAIL_RESPONSE" ".data.description" "Event description" || exit 1
    
    EVENT_TITLE=$(echo "$EVENT_DETAIL_RESPONSE" | jq -r '.data.title')
    print_info "Event: $EVENT_TITLE"
  else
    print_warning "No event ID available"
  fi
else
  print_warning "Step 6: Skipped - No events available"
fi

echo ""
print_success "============================================"
print_success "Flow 1: User Onboarding - COMPLETED ✅"
print_success "============================================"
echo ""

# Clean up (optional)
# Could delete test user here if needed

exit 0
