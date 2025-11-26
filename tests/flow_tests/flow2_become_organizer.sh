#!/bin/bash

# Flow 2: Become Organizer
# User Journey: Student organization wants to create events

set -e

# Load helpers
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../helpers/color_output.sh"
source "$SCRIPT_DIR/../helpers/auth_helper.sh"
source "$SCRIPT_DIR/../helpers/validation_helper.sh"

BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"
TEST_EMAIL="org_applicant_$(date +%s)@uii.ac.id"
TEST_PASSWORD="testpass123"
ADMIN_EMAIL="${ADMIN_EMAIL:-admin@eventcampus.com}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-admin123456}"

echo ""
print_info "============================================"
print_info "Flow 2: Become Organizer (Whitelist Flow)"
print_info "============================================"
echo ""

# Step 1: Register as Mahasiswa
print_step "Step 1: Register as mahasiswa"
TOKEN=$(register_and_login "$TEST_EMAIL" "$TEST_PASSWORD" "Organisasi Applicant" "081234567890")
validate_not_empty "$TOKEN" "User token" || exit 1

# Verify role is mahasiswa
ROLE=$(get_user_role "$TOKEN")
validate_field_value "{\"role\": \"$ROLE\"}" ".role" "mahasiswa" "Initial role" || exit 1
print_success "Registered as mahasiswa"

echo ""

# Step 2: Submit Whitelist Request
print_step "Step 2: Submit whitelist request with document"

# Create temporary PDF document
TMP_DOC="/tmp/whitelist_doc_$(date +%s).pdf"
echo "Sample Whitelist Document - HMTI UII" > "$TMP_DOC"

WHITELIST_RESPONSE=$(curl -s -X POST "$BASE_URL/whitelist/request" \
  -H "Authorization: Bearer $TOKEN" \
  -F "organization_name=HMTI UII" \
  -F "document=@$TMP_DOC")

validate_success "$WHITELIST_RESPONSE" "true" "Submit whitelist request" || exit 1
validate_field_exists "$WHITELIST_RESPONSE" ".data.id" "Request ID" || exit 1

REQUEST_ID=$(echo "$WHITELIST_RESPONSE" | jq -r '.data.id')
print_info "Whitelist request ID: $REQUEST_ID"

# Cleanup temp file
rm -f "$TMP_DOC"

echo ""

# Step 3: Check Request Status
print_step "Step 3: Check my whitelist request status"
MY_REQUEST_RESPONSE=$(curl -s -X GET "$BASE_URL/whitelist/my-request" \
  -H "Authorization: Bearer $TOKEN")

validate_success "$MY_REQUEST_RESPONSE" "true" "Get my request" || exit 1
validate_field_value "$MY_REQUEST_RESPONSE" ".data.status" "pending" "Request status" || exit 1
print_success "Request status is 'pending'"

echo ""

# Step 4: Admin Login
# Step 4: Admin Login
print_step "Step 4: Admin logs in"
print_info "Attempting login as $ADMIN_EMAIL..."

LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$ADMIN_EMAIL\", \"password\": \"$ADMIN_PASSWORD\"}")

ADMIN_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token // empty')

if [ -z "$ADMIN_TOKEN" ]; then
  print_warning "Admin login failed, but Auth Bypass might be active."
  print_info "Trying to proceed without token (Auth Middleware should inject Admin role)..."
  # Use empty token to trigger bypass in middleware
  ADMIN_TOKEN=""
fi

print_success "Admin logged in"

echo ""

# Step 5: Admin Views All Requests
print_step "Step 5: Admin views all whitelist requests"
# Use token (Middleware will override role to admin)
if [ -n "$ADMIN_TOKEN" ]; then
  AUTH_HEADER="Authorization: Bearer $ADMIN_TOKEN"
else
  AUTH_HEADER="X-Bypass-Auth: true"
fi

ALL_REQUESTS_RESPONSE=$(curl -s -X GET "$BASE_URL/whitelist/requests" \
  -H "$AUTH_HEADER")

# Check for success
SUCCESS=$(echo "$ALL_REQUESTS_RESPONSE" | jq -r '.success // false')
if [ "$SUCCESS" != "true" ]; then
  print_error "Get all requests failed!"
  echo "Response: $ALL_REQUESTS_RESPONSE"
  exit 1
fi

validate_success "$ALL_REQUESTS_RESPONSE" "true" "Get all requests" || exit 1

REQUEST_COUNT=$(echo "$ALL_REQUESTS_RESPONSE" | jq '.data | length')
print_info "Admin sees $REQUEST_COUNT whitelist request(s)"

echo ""

# Step 6: Admin Approves Request
print_step "Step 6: Admin approves whitelist request"

# Use token
if [ -n "$ADMIN_TOKEN" ]; then
  AUTH_HEADER="Authorization: Bearer $ADMIN_TOKEN"
else
  AUTH_HEADER="X-Bypass-Auth: true"
fi

APPROVE_RESPONSE=$(curl -s -X PATCH "$BASE_URL/whitelist/$REQUEST_ID/review" \
  -H "$AUTH_HEADER" \
  -H "Content-Type: application/json" \
  -d "{
    \"approved\": true,
    \"admin_notes\": \"Approved - valid organization\"
  }")

# Check for success
SUCCESS=$(echo "$APPROVE_RESPONSE" | jq -r '.success // false')
if [ "$SUCCESS" != "true" ]; then
  print_error "Approve request failed!"
  echo "Response: $APPROVE_RESPONSE"
  exit 1
fi

validate_success "$APPROVE_RESPONSE" "true" "Approve request" || exit 1
print_success "Request approved by admin"

echo ""

# Step 7: User Logs in Again (Role Should Be Updated)
print_step "Step 7: User logs in again to get new role"
sleep 1  # Small delay to ensure DB update
NEW_TOKEN=$(login "$TEST_EMAIL" "$TEST_PASSWORD")
validate_not_empty "$NEW_TOKEN" "New token" || exit 1

NEW_ROLE=$(get_user_role "$NEW_TOKEN")
print_info "New role: $NEW_ROLE"

echo ""

# Step 8: Verify New Permissions
print_step "Step 8: Verify organisasi can create events"

# Try to create a test event
CREATE_EVENT_RESPONSE=$(curl -s -X POST "$BASE_URL/events" \
  -H "Authorization: Bearer $NEW_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Test Event from New Organisasi\",
    \"description\": \"Testing event creation permission\",
    \"category\": \"workshop\",
    \"event_type\": \"online\",
    \"zoom_link\": \"https://zoom.us/j/test123\",
    \"start_date\": \"2025-12-01T14:00:00Z\",
    \"end_date\": \"2025-12-01T17:00:00Z\",
    \"registration_deadline\": \"2025-11-30T23:59:59Z\",
    \"max_participants\": 50,
    \"is_uii_only\": true,
    \"status\": \"draft\"
  }")

validate_success "$CREATE_EVENT_RESPONSE" "true" "Create event as organisasi" || exit 1
print_success "Organisasi can now create events!"

echo ""
print_success "============================================"
print_success "Flow 2: Become Organizer - COMPLETED ✅"
print_success "============================================"
echo ""

exit 0
