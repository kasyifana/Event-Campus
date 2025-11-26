#!/bin/bash

# Flow 6: Advanced Scenarios
# Covers: Waitlist System, Cancellation, Event Deletion, Bulk Attendance

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../helpers/color_output.sh"
source "$SCRIPT_DIR/../helpers/auth_helper.sh"
source "$SCRIPT_DIR/../helpers/validation_helper.sh"

BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"

echo ""
print_info "============================================"
print_info "Flow 6: Advanced Scenarios"
print_info "============================================"
echo ""

# Setup: Login as Organizer
print_step "Setup: Login as Organizer"
ORG_EMAIL="organizer@eventcampus.com"
ORG_PASSWORD="orgpass123456"

ORG_TOKEN=$(login "$ORG_EMAIL" "$ORG_PASSWORD")
validate_not_empty "$ORG_TOKEN" "Organizer token" || exit 1

echo ""

# ==========================================
# SCENARIO 1: Waitlist & Cancellation
# ==========================================
print_info "--------------------------------------------"
print_info "Scenario 1: Waitlist & Cancellation System"
print_info "--------------------------------------------"

# 1. Create Event with Capacity 1
print_step "1. Create event with capacity 1"
EVENT_RESPONSE=$(curl -s -X POST "$BASE_URL/events" \
  -H "Authorization: Bearer $ORG_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Exclusive 1-on-1 Mentoring\",
    \"description\": \"Limited seat event\",
    \"category\": \"seminar\",
    \"event_type\": \"online\",
    \"zoom_link\": \"https://zoom.us/j/999888777\",
    \"start_date\": \"$(date -v+1d -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "1 day" +"%Y-%m-%dT%H:%M:%SZ")\",
    \"end_date\": \"$(date -v+1d -v+2H -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "1 day 2 hours" +"%Y-%m-%dT%H:%M:%SZ")\",
    \"registration_deadline\": \"$(date -v+23H -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "23 hours" +"%Y-%m-%dT%H:%M:%SZ")\",
    \"max_participants\": 1,
    \"is_uii_only\": false,
    \"status\": \"published\"
  }")

# Note: We create as published directly if API allows, otherwise draft then publish
# Based on previous flows, we might need to upload poster first.
# Let's try to extract ID first.
EVENT_ID=$(echo "$EVENT_RESPONSE" | jq -r '.data.id // empty')

if [ -z "$EVENT_ID" ]; then
  print_error "Failed to create event"
  echo "$EVENT_RESPONSE"
  exit 1
fi

# Upload poster (dummy)
# Create a proper dummy image file to avoid validation issues
TMP_POSTER="/tmp/dummy_poster.jpg"
echo "dummy image content" > "$TMP_POSTER"

POSTER_RESP=$(curl -s -X POST "$BASE_URL/events/$EVENT_ID/poster" \
  -H "Authorization: Bearer $ORG_TOKEN" \
  -F "poster=@$TMP_POSTER")

validate_success "$POSTER_RESP" "true" "Upload poster" || exit 1

# Publish
PUBLISH_RESP=$(curl -s -X POST "$BASE_URL/events/$EVENT_ID/publish" \
  -H "Authorization: Bearer $ORG_TOKEN")

validate_success "$PUBLISH_RESP" "true" "Publish event" || exit 1

print_success "Event created and published (ID: $EVENT_ID)"

# 2. Register User A (Should be Registered)
print_step "2. Register User A (Should get seat)"
USER_A_EMAIL="usera_$(date +%s)@example.com"
USER_A_TOKEN=$(register_and_login "$USER_A_EMAIL" "password123" "User A" "08111111111")

REG_A_RESPONSE=$(curl -s -X POST "$BASE_URL/events/$EVENT_ID/register" \
  -H "Authorization: Bearer $USER_A_TOKEN")

REG_A_ID=$(echo "$REG_A_RESPONSE" | jq -r '.data.registration_id')
REG_A_STATUS=$(echo "$REG_A_RESPONSE" | jq -r '.data.status')

if [ "$REG_A_STATUS" == "registered" ]; then
  print_success "User A registered successfully"
else
  print_error "User A failed to register. Status: $REG_A_STATUS"
  echo "Response: $REG_A_RESPONSE"
  exit 1
fi

# 3. Register User B (Should be Waitlist)
print_step "3. Register User B (Should be Waitlist)"
USER_B_EMAIL="userb_$(date +%s)@example.com"
USER_B_TOKEN=$(register_and_login "$USER_B_EMAIL" "password123" "User B" "08222222222")

REG_B_RESPONSE=$(curl -s -X POST "$BASE_URL/events/$EVENT_ID/register" \
  -H "Authorization: Bearer $USER_B_TOKEN")

REG_B_ID=$(echo "$REG_B_RESPONSE" | jq -r '.data.registration_id')
REG_B_STATUS=$(echo "$REG_B_RESPONSE" | jq -r '.data.status')

if [ "$REG_B_STATUS" == "waitlist" ]; then
  print_success "User B correctly placed in waitlist"
else
  print_error "User B status incorrect. Expected 'waitlist', got '$REG_B_STATUS'"
  exit 1
fi

# 4. User A Cancels
print_step "4. User A cancels registration"
CANCEL_RESPONSE=$(curl -s -X DELETE "$BASE_URL/registrations/$REG_A_ID" \
  -H "Authorization: Bearer $USER_A_TOKEN")

validate_success "$CANCEL_RESPONSE" "true" "Cancel registration" || exit 1

# 5. Verify User B Promoted
print_step "5. Verify User B promoted to Registered"
# User B checks their registration status
MY_REGS_B=$(curl -s -X GET "$BASE_URL/registrations/my" \
  -H "Authorization: Bearer $USER_B_TOKEN")

# Filter for this event
NEW_STATUS_B=$(echo "$MY_REGS_B" | jq -r ".data[] | select(.id == \"$REG_B_ID\") | .status")

if [ "$NEW_STATUS_B" == "registered" ]; then
  print_success "User B successfully promoted to 'registered'!"
else
  print_error "User B promotion failed. Status: $NEW_STATUS_B"
  exit 1
fi

echo ""

# ==========================================
# SCENARIO 2: Event Deletion & Cancellation
# ==========================================
print_info "--------------------------------------------"
print_info "Scenario 2: Event Deletion & Cancellation"
print_info "--------------------------------------------"

# 1. Delete Draft Event
print_step "1. Create and Delete Draft Event"
DRAFT_RESP=$(curl -s -X POST "$BASE_URL/events" \
  -H "Authorization: Bearer $ORG_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Mistake Event\",
    \"description\": \"To be deleted\",
    \"category\": \"lomba\",
    \"event_type\": \"offline\",
    \"location\": \"Nowhere\",
    \"start_date\": \"2025-12-12T10:00:00Z\",
    \"end_date\": \"2025-12-12T12:00:00Z\",
    \"registration_deadline\": \"2025-12-11T10:00:00Z\",
    \"max_participants\": 100,
    \"status\": \"draft\"
  }")
DRAFT_ID=$(echo "$DRAFT_RESP" | jq -r '.data.id')

DELETE_RESP=$(curl -s -X DELETE "$BASE_URL/events/$DRAFT_ID" \
  -H "Authorization: Bearer $ORG_TOKEN")

validate_success "$DELETE_RESP" "true" "Delete draft event" || exit 1

# Verify it's gone
GET_RESP=$(curl -s -X GET "$BASE_URL/events/$DRAFT_ID" \
  -H "Authorization: Bearer $ORG_TOKEN")
SUCCESS=$(echo "$GET_RESP" | jq -r '.success')
if [ "$SUCCESS" == "false" ]; then
  print_success "Event correctly not found after deletion"
else
  print_warning "Event still accessible after deletion (might be soft delete?)"
fi

# 2. Cancel Published Event
print_step "2. Cancel Published Event"
# Reuse the event from Scenario 1 (ID: $EVENT_ID) which is published
CANCEL_EVENT_RESP=$(curl -s -X DELETE "$BASE_URL/events/$EVENT_ID" \
  -H "Authorization: Bearer $ORG_TOKEN")

validate_success "$CANCEL_EVENT_RESP" "true" "Cancel published event" || exit 1

# Verify status is 'cancelled'
GET_EVENT_RESP=$(curl -s -X GET "$BASE_URL/events/$EVENT_ID" \
  -H "Authorization: Bearer $ORG_TOKEN")
EVENT_STATUS=$(echo "$GET_EVENT_RESP" | jq -r '.data.status')

if [ "$EVENT_STATUS" == "cancelled" ]; then
  print_success "Event status updated to 'cancelled'"
else
  print_error "Event status incorrect. Expected 'cancelled', got '$EVENT_STATUS'"
  exit 1
fi

echo ""

# ==========================================
# SCENARIO 3: Bulk Attendance
# ==========================================
print_info "--------------------------------------------"
print_info "Scenario 3: Bulk Attendance"
print_info "--------------------------------------------"

# Create new event for bulk test
print_step "1. Create event for bulk attendance"
BULK_EVENT_RESP=$(curl -s -X POST "$BASE_URL/events" \
  -H "Authorization: Bearer $ORG_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Bulk Attendance Party\",
    \"description\": \"Everyone is invited\",
    \"category\": \"konser\",
    \"event_type\": \"offline\",
    \"location\": \"Main Hall\",
    \"start_date\": \"$(date -v+1d -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "1 day" +"%Y-%m-%dT%H:%M:%SZ")\",
    \"end_date\": \"$(date -v+1d -v+2H -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "1 day 2 hours" +"%Y-%m-%dT%H:%M:%SZ")\",
    \"registration_deadline\": \"$(date -v+23H -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "23 hours" +"%Y-%m-%dT%H:%M:%SZ")\",
    \"max_participants\": 100,
    \"status\": \"published\"
  }")
BULK_EVENT_ID=$(echo "$BULK_EVENT_RESP" | jq -r '.data.id // empty')

if [ -z "$BULK_EVENT_ID" ]; then
  print_error "Failed to create bulk event"
  echo "$BULK_EVENT_RESP"
  exit 1
fi

# Upload poster & Publish
# Reuse dummy poster
curl -s -X POST "$BASE_URL/events/$BULK_EVENT_ID/poster" -H "Authorization: Bearer $ORG_TOKEN" -F "poster=@$TMP_POSTER" > /dev/null
curl -s -X POST "$BASE_URL/events/$BULK_EVENT_ID/publish" -H "Authorization: Bearer $ORG_TOKEN" > /dev/null

# Register User A and B to this new event (While event is still future)
print_step "2. Register Users A and B"
# Need to login again to get tokens if they expired, but they should be valid
# User A
REG_A_BULK=$(curl -s -X POST "$BASE_URL/events/$BULK_EVENT_ID/register" -H "Authorization: Bearer $USER_A_TOKEN")
validate_success "$REG_A_BULK" "true" "User A bulk register" || echo "Response: $REG_A_BULK"

# User B
REG_B_BULK=$(curl -s -X POST "$BASE_URL/events/$BULK_EVENT_ID/register" -H "Authorization: Bearer $USER_B_TOKEN")
validate_success "$REG_B_BULK" "true" "User B bulk register" || echo "Response: $REG_B_BULK"

# Update to Past Date (to allow attendance)
print_step "2.5 Update event to start in past"
PAST_START=$(date -v-1H -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "1 hour ago" +"%Y-%m-%dT%H:%M:%SZ")
PAST_END=$(date -v+2H -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "2 hours" +"%Y-%m-%dT%H:%M:%SZ")

UPDATE_RESP=$(curl -s -X PUT "$BASE_URL/events/$BULK_EVENT_ID" \
  -H "Authorization: Bearer $ORG_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"start_date\": \"$PAST_START\",
    \"end_date\": \"$PAST_END\"
  }")
validate_success "$UPDATE_RESP" "true" "Update bulk event time" || exit 1

# Get User IDs
REGS_RESP=$(curl -s -X GET "$BASE_URL/events/$BULK_EVENT_ID/registrations" -H "Authorization: Bearer $ORG_TOKEN")
echo "Registrations Response: $REGS_RESP"

USER_IDS=$(echo "$REGS_RESP" | jq -r '.data[].user_id')

# Format for JSON array: ["id1", "id2"]
JSON_IDS=$(echo "$USER_IDS" | jq -R . | jq -s .)

print_step "3. Mark Bulk Attendance"
BULK_ATT_RESP=$(curl -s -X POST "$BASE_URL/events/$BULK_EVENT_ID/attendance/bulk" \
  -H "Authorization: Bearer $ORG_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_ids\": $JSON_IDS
  }")

validate_success "$BULK_ATT_RESP" "true" "Bulk mark attendance" || exit 1

# Verify count
COUNT_RESP=$(curl -s -X GET "$BASE_URL/events/$BULK_EVENT_ID/attendance" -H "Authorization: Bearer $ORG_TOKEN")
COUNT=$(echo "$COUNT_RESP" | jq '.data | length')

if [ "$COUNT" -ge 2 ]; then
  print_success "Bulk attendance verified. Count: $COUNT"
else
  print_error "Bulk attendance count mismatch. Expected >= 2, got $COUNT"
  exit 1
fi

echo ""
print_success "============================================"
print_success "Flow 6: Advanced Scenarios - COMPLETED ✅"
print_success "============================================"
echo ""

exit 0
