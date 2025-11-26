#!/bin/bash

# Flow 5: Attendance Flow
# User Journey: Organizer marks attendance for participants

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../helpers/color_output.sh"
source "$SCRIPT_DIR/../helpers/auth_helper.sh"
source "$SCRIPT_DIR/../helpers/validation_helper.sh"

BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"

echo ""
print_info "============================================"
print_info "Flow 5: Attendance Flow"
print_info "============================================"
echo ""

# Prerequisites: Need Event ID and Registration ID
if [ -f /tmp/test_event_id.sh ]; then
  source /tmp/test_event_id.sh
fi
if [ -f /tmp/test_registration_id.sh ]; then
  source /tmp/test_registration_id.sh
fi

if [ -z "$TEST_EVENT_ID" ] || [ -z "$TEST_REG_ID" ]; then
  print_warning "Missing Event ID or Registration ID from previous flows."
  print_info "Skipping Flow 5 (Dependencies not met)"
  exit 0
fi

EVENT_ID="$TEST_EVENT_ID"
REG_ID="$TEST_REG_ID"

print_info "Using Event ID: $EVENT_ID"
print_info "Using Registration ID: $REG_ID"

echo ""

# Step 1: Login as Organizer
print_step "Step 1: Login as Organizer"
ORG_EMAIL="organizer@eventcampus.com"
ORG_PASSWORD="orgpass123456"

TOKEN=$(login "$ORG_EMAIL" "$ORG_PASSWORD")
validate_not_empty "$TOKEN" "Organizer token" || exit 1

echo ""

# Step 2: Get Event Registrations
print_step "Step 2: Get event registrations"
REGS_RESPONSE=$(curl -s -X GET "$BASE_URL/events/$EVENT_ID/registrations" \
  -H "Authorization: Bearer $TOKEN")

validate_success "$REGS_RESPONSE" "true" "Get registrations" || exit 1

# Verify our registration is in the list
USER_ID=$(echo "$REGS_RESPONSE" | jq -r ".data[] | select(.id == \"$REG_ID\") | .user_id")

if [ -n "$USER_ID" ] && [ "$USER_ID" != "null" ]; then
  print_success "Target registration found. User ID: $USER_ID"
else
  print_error "Target registration NOT found in list"
  exit 1
fi

echo ""

# Step 2.5: Update Event to Start Now (so we can mark attendance)
print_step "Step 2.5: Update event start time to allow attendance"
# Get current date - 1 hour in ISO 8601 format
# MacOS date command is different from Linux
if [[ "$OSTYPE" == "darwin"* ]]; then
  START_DATE=$(date -v-1H -u +"%Y-%m-%dT%H:%M:%SZ")
  END_DATE=$(date -v+2H -u +"%Y-%m-%dT%H:%M:%SZ")
else
  START_DATE=$(date -u -d "1 hour ago" +"%Y-%m-%dT%H:%M:%SZ")
  END_DATE=$(date -u -d "2 hours" +"%Y-%m-%dT%H:%M:%SZ")
fi

UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/events/$EVENT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"start_date\": \"$START_DATE\",
    \"end_date\": \"$END_DATE\"
  }")

# Check for success
SUCCESS=$(echo "$UPDATE_RESPONSE" | jq -r '.success // false')
if [ "$SUCCESS" != "true" ]; then
  print_error "Update event time failed!"
  echo "Response: $UPDATE_RESPONSE"
  exit 1
fi

validate_success "$UPDATE_RESPONSE" "true" "Update event time" || exit 1

echo ""

# Step 3: Mark Attendance
print_step "Step 3: Mark attendance"
ATTENDANCE_RESPONSE=$(curl -s -X POST "$BASE_URL/events/$EVENT_ID/attendance" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"$USER_ID\",
    \"status\": \"attended\"
  }")

# Check for success
SUCCESS=$(echo "$ATTENDANCE_RESPONSE" | jq -r '.success // false')
if [ "$SUCCESS" != "true" ]; then
  print_error "Mark attendance failed!"
  echo "Response: $ATTENDANCE_RESPONSE"
  exit 1
fi

validate_success "$ATTENDANCE_RESPONSE" "true" "Mark attendance" || exit 1

echo ""

# Step 4: Verify Attendance Status
print_step "Step 4: Verify attendance status"
# We can verify by getting registrations again
VERIFY_RESPONSE=$(curl -s -X GET "$BASE_URL/events/$EVENT_ID/registrations" \
  -H "Authorization: Bearer $TOKEN")

STATUS=$(echo "$VERIFY_RESPONSE" | jq -r ".data[] | select(.id == \"$REG_ID\") | .status")

if [ "$STATUS" == "attended" ]; then
  print_success "Status updated to 'attended'"
else
  print_error "Status update failed. Expected 'attended', got '$STATUS'"
  exit 1
fi

echo ""
print_success "============================================"
print_success "Flow 5: Attendance Flow - COMPLETED ✅"
print_success "============================================"
echo ""

exit 0
