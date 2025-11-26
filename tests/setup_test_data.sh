#!/bin/bash

# Setup Test Data
# Registers Admin and Organizer accounts for testing

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/helpers/color_output.sh"
source "$SCRIPT_DIR/helpers/auth_helper.sh"

BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"

print_info "============================================"
print_info "       Event Campus - Test Data Setup       "
print_info "============================================"
echo ""

# 1. Setup Admin Account
print_step "1. Setting up Admin Account"
ADMIN_EMAIL="admin@eventcampus.com"
ADMIN_PASSWORD="admin123456"

# Try to login first
ADMIN_TOKEN=$(login "$ADMIN_EMAIL" "$ADMIN_PASSWORD")

if [ -n "$ADMIN_TOKEN" ]; then
  print_success "Admin account already exists and is accessible."
else
  print_warning "Admin account not found or credentials invalid."
  print_info "Attempting to register admin user..."
  
  # Register
  REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"$ADMIN_EMAIL\",
      \"password\": \"$ADMIN_PASSWORD\",
      \"full_name\": \"System Admin\",
      \"phone_number\": \"081299999999\"
    }")
  
  SUCCESS=$(echo "$REGISTER_RESPONSE" | jq -r '.success')
  
  if [ "$SUCCESS" == "true" ]; then
    print_success "Admin user registered successfully."
    print_warning "⚠️  ACTION REQUIRED: You must manually promote this user to 'admin' in the database."
    print_info "Run this SQL command:"
    echo "UPDATE users SET role = 'admin', is_approved = true WHERE email = '$ADMIN_EMAIL';"
  else
    print_error "Failed to register admin user."
    echo "$REGISTER_RESPONSE" | jq '.'
  fi
fi

echo ""

# 2. Setup Organizer Account (for Flow 3)
print_step "2. Setting up Organizer Account"
ORG_EMAIL="organizer@eventcampus.com"
ORG_PASSWORD="orgpass123456"

# Try to login
ORG_TOKEN=$(login "$ORG_EMAIL" "$ORG_PASSWORD")

if [ -n "$ORG_TOKEN" ]; then
  print_success "Organizer account already exists."
else
  print_warning "Organizer account not found."
  print_info "Attempting to register organizer user..."
  
  # Register
  REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"$ORG_EMAIL\",
      \"password\": \"$ORG_PASSWORD\",
      \"full_name\": \"Test Organizer\",
      \"phone_number\": \"081288888888\"
    }")
  
  SUCCESS=$(echo "$REGISTER_RESPONSE" | jq -r '.success')
  
  if [ "$SUCCESS" == "true" ]; then
    print_success "Organizer user registered successfully."
    print_warning "⚠️  ACTION REQUIRED: You must manually promote this user to 'organisasi' in the database."
    print_info "Run this SQL command:"
    echo "UPDATE users SET role = 'organisasi', is_approved = true WHERE email = '$ORG_EMAIL';"
  else
    print_error "Failed to register organizer user."
    echo "$REGISTER_RESPONSE" | jq '.'
  fi
fi

echo ""
print_info "============================================"
print_info "Setup complete. Please run the SQL commands above if needed."
print_info "Then run: ./run_all_flows.sh"
print_info "============================================"
