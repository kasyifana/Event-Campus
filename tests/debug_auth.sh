#!/bin/bash

# Debug Auth & Roles
# Logs in as Admin and Organizer to check tokens and roles

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/helpers/color_output.sh"
source "$SCRIPT_DIR/helpers/auth_helper.sh"

BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"

print_info "============================================"
print_info "       Debug Auth & Roles                   "
print_info "============================================"
echo ""

# 1. Check Admin
print_step "1. Checking Admin Account (admin@eventcampus.com)"

# Try password 'admin123456'
print_info "Attempting login with password 'admin123456'..."
ADMIN_TOKEN=$(login "admin@eventcampus.com" "admin123456")

if [ -n "$ADMIN_TOKEN" ]; then
  print_success "Login SUCCESS with 'admin123456'"
  print_info "Token: ${ADMIN_TOKEN:0:20}..."
  
  ROLE=$(get_user_role "$ADMIN_TOKEN")
  print_info "Role in DB: $ROLE"
  
  if [ "$ROLE" == "admin" ]; then
    print_success "Role is correct (admin)"
  else
    print_error "Role is INCORRECT (expected 'admin', got '$ROLE')"
  fi
else
  print_warning "Login FAILED with 'admin123456'"
  
  # Try password 'admin123'
  print_info "Attempting login with password 'admin123'..."
  ADMIN_TOKEN=$(login "admin@eventcampus.com" "admin123")
  
  if [ -n "$ADMIN_TOKEN" ]; then
    print_success "Login SUCCESS with 'admin123'"
    print_info "Token: ${ADMIN_TOKEN:0:20}..."
    
    ROLE=$(get_user_role "$ADMIN_TOKEN")
    print_info "Role in DB: $ROLE"
     if [ "$ROLE" == "admin" ]; then
      print_success "Role is correct (admin)"
    else
      print_error "Role is INCORRECT (expected 'admin', got '$ROLE')"
    fi
  else
    print_error "Login FAILED with 'admin123' too."
  fi
fi

echo ""

# 2. Check Organizer
print_step "2. Checking Organizer Account (organizer@eventcampus.com)"

print_info "Attempting login with password 'orgpass123456'..."
ORG_TOKEN=$(login "organizer@eventcampus.com" "orgpass123456")

if [ -n "$ORG_TOKEN" ]; then
  print_success "Login SUCCESS"
  print_info "Token: ${ORG_TOKEN:0:20}..."
  
  ROLE=$(get_user_role "$ORG_TOKEN")
  print_info "Role in DB: $ROLE"
  
  if [ "$ROLE" == "organisasi" ]; then
    print_success "Role is correct (organisasi)"
  else
    print_error "Role is INCORRECT (expected 'organisasi', got '$ROLE')"
    print_info "This is why Flow 3 fails with 'Forbidden'!"
  fi
else
  print_error "Login FAILED"
fi

echo ""
