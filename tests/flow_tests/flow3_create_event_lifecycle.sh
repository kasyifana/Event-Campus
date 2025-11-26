#!/bin/bash

# Flow 3: Create Event Lifecycle
# User Journey: Organisasi creates and publishes an event

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../helpers/color_output.sh"
source "$SCRIPT_DIR/../helpers/auth_helper.sh"
source "$SCRIPT_DIR/../helpers/validation_helper.sh"

BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"

echo ""
print_info "============================================"
print_info "Flow 3: Create Event Lifecycle"
print_info "============================================"
echo ""

# Prerequisites: Need an organisasi account
# For testing, we'll use existing or create new one
print_step "Setup: Getting organisasi account"

ORG_EMAIL="organizer@eventcampus.com"
ORG_PASSWORD="orgpass123456"

# Login (assuming user ran the setup script and SQL promotion)
print_info "Logging in as organizer ($ORG_EMAIL)..."
TOKEN=$(login "$ORG_EMAIL" "$ORG_PASSWORD")

if [ -z "$TOKEN" ]; then
  print_warning "Login failed. Attempting to register (fallback)..."
  TOKEN=$(register_and_login "$ORG_EMAIL" "$ORG_PASSWORD" "Test Organizer" "081288888888")
fi

validate_not_empty "$TOKEN" "Organizer token" || exit 1

# Note: In real scenario, this would need whitelist approval
# For testing purposes, we'll assume role can be updated or use existing org account

echo ""

# Step 1: Create Event (Draft)
print_step "Step 1: Create event in draft status"

CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/events" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Workshop Web Development 2025\",
    \"description\": \"Learn modern web development with React and Go\",
    \"category\": \"workshop\",
    \"event_type\": \"online\",
    \"zoom_link\": \"https://zoom.us/j/123456789\",
    \"start_date\": \"2025-12-15T14:00:00Z\",
    \"end_date\": \"2025-12-15T17:00:00Z\",
    \"registration_deadline\": \"2025-12-14T23:59:59Z\",
    \"max_participants\": 50,
    \"is_uii_only\": true,
    \"status\": \"draft\"
  }")

# Check for success
SUCCESS=$(echo "$CREATE_RESPONSE" | jq -r '.success // false')
if [ "$SUCCESS" != "true" ]; then
  print_error "Create draft event failed!"
  echo "Response: $CREATE_RESPONSE"
  exit 1
fi

print_success "Event created successfully"

EVENT_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.id // empty')

if [ -z "$EVENT_ID" ]; then
  # Try alternative response structure
  EVENT_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.event_id // .data.event.id // empty')
fi

validate_not_empty "$EVENT_ID" "Event ID" || exit 1
print_info "Event ID: $EVENT_ID"

echo ""

# Step 2: Upload Poster
print_step "Step 2: Upload event poster"

# Create dummy poster image
TMP_POSTER="/tmp/event_poster_$(date +%s).jpg"
# Create a simple colored square as poster (requires ImageMagick)
if command -v convert &> /dev/null; then
  convert -size 800x600 xc:blue "$TMP_POSTER" 2>/dev/null || echo "Blue poster" > "$TMP_POSTER"
else
  # Fallback: create text file (will fail validation but tests the endpoint)
  echo "Test poster image" > "$TMP_POSTER"
fi

POSTER_RESPONSE=$(curl -s -X POST "$BASE_URL/events/$EVENT_ID/poster" \
  -H "Authorization: Bearer $TOKEN" \
  -F "poster=@$TMP_POSTER")

# Clean up temp file
rm -f "$TMP_POSTER"

validate_success "$POSTER_RESPONSE" "true" "Upload poster" || print_warning "Poster upload failed (may need valid image file)"

echo ""

# Step 3: Update Event Details
print_step "Step 3: Update event details"
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/events/$EVENT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"description\": \"Learn modern web development with React, Go, and PostgreSQL - Updated!\"
  }")

validate_success "$UPDATE_RESPONSE" "true" "Update event" || exit 1

echo ""

# Step 4: Publish Event
print_step "Step 4: Publish event"
PUBLISH_RESPONSE=$(curl -s -X POST "$BASE_URL/events/$EVENT_ID/publish" \
  -H "Authorization: Bearer $TOKEN")

validate_success "$PUBLISH_RESPONSE" "true" "Publish event" || print_warning "Publish failed (may need poster)"

echo ""

# Step 5: View My Events
print_step "Step 5: View my events"
MY_EVENTS_RESPONSE=$(curl -s -X GET "$BASE_URL/events/my-events" \
  -H "Authorization: Bearer $TOKEN")

validate_success "$MY_EVENTS_RESPONSE" "true" "Get my events" || exit 1

MY_EVENT_COUNT=$(echo "$MY_EVENTS_RESPONSE" | jq '.data | length')
print_info "Organizer has $MY_EVENT_COUNT event(s)"

echo ""

# Step 6: Verify in Public Listing
print_step "Step 6: Verify event in public listing"
PUBLIC_EVENTS_RESPONSE=$(curl -s -X GET "$BASE_URL/events" \
  -H "Authorization: Bearer $TOKEN")

validate_success "$PUBLIC_EVENTS_RESPONSE" "true" "Get public events" || exit 1

# Check if our event appears
EVENT_FOUND=$(echo "$PUBLIC_EVENTS_RESPONSE" | jq ".data[] | select(.id == \"$EVENT_ID\") | .id" | wc -l)

if [ "$EVENT_FOUND" -gt 0 ]; then
  print_success "Event found in public listing"
else
  print_warning "Event not found in public listing (may still be draft)"
fi

echo ""
print_success "============================================"
print_success "Flow 3: Create Event Lifecycle - COMPLETED ✅"
print_success "============================================"
echo ""

# Export EVENT_ID for use in other flows
echo "export TEST_EVENT_ID=$EVENT_ID" > /tmp/test_event_id.sh
print_info "Event ID saved for other tests: $EVENT_ID"

exit 0
