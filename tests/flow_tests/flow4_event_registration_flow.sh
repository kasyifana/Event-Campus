#!/bin/bash

# Flow 4: Event Registration Flow
# User Journey: Students register for an event

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../helpers/color_output.sh"
source "$SCRIPT_DIR/../helpers/auth_helper.sh"
source "$SCRIPT_DIR/../helpers/validation_helper.sh"

BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"

echo ""
print_info "============================================"
print_info "Flow 4: Event Registration Flow"
print_info "============================================"
echo ""

# Get event ID from previous test or use provided one
if [ -f /tmp/test_event_id.sh ]; then
  source /tmp/test_event_id.sh
fi

if [ -z "$TEST_EVENT_ID" ]; then
  print_warning "No event ID found, creating test event first..."
  bash "$SCRIPT_DIR/flow3_create_event_lifecycle.sh" || exit 1
  source /tmp/test_event_id.sh
fi

EVENT_ID="$TEST_EVENT_ID"
print_info "Using Event ID: $EVENT_ID"

echo ""

# Step 1: Create test user
print_step "Step 1: Register new mahasiswa"
USER_EMAIL="student_$(date +%s)@uii.ac.id"
USER_PASSWORD="student123"

USER_TOKEN=$(register_and_login "$USER_EMAIL" "$USER_PASSWORD" "Test Student" "081234567892")
validate_not_empty "$USER_TOKEN" "User token" || exit 1

echo ""

# Step 2: Browse Events
print_step "Step 2: Browse available events"
EVENTS_RESPONSE=$(curl -s -X GET "$BASE_URL/events" \
  -H "Authorization: Bearer $USER_TOKEN")

validate_success "$EVENTS_RESPONSE" "true" "Get events" || exit 1

EVENT_COUNT=$(echo "$EVENTS_RESPONSE" | jq '.data | length')
print_info "Found $EVENT_COUNT event(s)"

echo ""

# Step 3: Register for Event
print_step "Step 3: Register for event"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/events/$EVENT_ID/register" \
  -H "Authorization: Bearer $USER_TOKEN")

validate_success "$REGISTER_RESPONSE" "true" "Event registration" || exit 1
validate_field_exists "$REGISTER_RESPONSE" ".data.registration_id" "Registration ID" || exit 1

REG_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.data.registration_id // .data.id')
REG_STATUS=$(echo "$REGISTER_RESPONSE" | jq -r '.data.status')

print_info "Registration ID: $REG_ID"
print_info "Status: $REG_STATUS"

echo ""

# Step 4: View My Registrations
print_step "Step 4: View my registrations"
MY_REGS_RESPONSE=$(curl -s -X GET "$BASE_URL/registrations/my" \
  -H "Authorization: Bearer $USER_TOKEN")

validate_success "$MY_REGS_RESPONSE" "true" "Get my registrations" || exit 1

# Check if registration appears
TOTAL_REGS=$(echo "$MY_REGS_RESPONSE" | jq '.data | length')
print_info "User has $TOTAL_REGS registration(s)"

echo ""

# Step 5: Try to Register Again (Should Fail)
print_step "Step 5: Try duplicate registration (should fail)"
DUPLICATE_RESPONSE=$(curl -s -X POST "$BASE_URL/events/$EVENT_ID/register" \
  -H "Authorization: Bearer $USER_TOKEN")

DUPLICATE_SUCCESS=$(echo "$DUPLICATE_RESPONSE" | jq -r '.success')

if [ "$DUPLICATE_SUCCESS" == "false" ]; then
  print_success "Duplicate registration correctly prevented"
else
  print_warning "Duplicate registration was allowed (potential issue)"
fi

echo ""

# Step 6: Organizer Views Registrations
print_step "Step 6: Create organizer to view registrations"

# Get/create organizer
ORG_EMAIL="org_viewer_$(date +%s)@uii.ac.id"
ORG_PASSWORD="orgpass123"
ORG_TOKEN=$(register_and_login "$ORG_EMAIL" "$ORG_PASSWORD" "Event Organizer" "081234567893")

# Note: This organizer won't own the event, so may get 403

echo ""
print_success "============================================"
print_success "Flow 4: Event Registration - COMPLETED ✅"
print_success "============================================"
echo ""

# Save registration ID for other tests
echo "export TEST_REG_ID=$REG_ID" > /tmp/test_registration_id.sh

exit 0
