#!/bin/bash

# Base URL for API
export BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"

# Login and return JWT token
login() {
  local email=$1
  local password=$2
  
  local response=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$email\", \"password\": \"$password\"}")
  
  local token=$(echo "$response" | jq -r '.data.token // empty')
  
  if [ -z "$token" ] || [ "$token" == "null" ]; then
    return 1
  fi
  
  echo "$token"
  return 0
}

# Register new user and return JWT token
register_and_login() {
  local email=$1
  local password=$2
  local full_name=$3
  local phone=$4
  
  # Register
  local reg_response=$(curl -s -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"$email\",
      \"password\": \"$password\",
      \"full_name\": \"$full_name\",
      \"phone_number\": \"$phone\"
    }")
  
  local success=$(echo "$reg_response" | jq -r '.success // false')
  
  if [ "$success" != "true" ]; then
    echo "$reg_response" | jq -r '.error // "Registration failed"' >&2
    return 1
  fi
  
  # Login
  login "$email" "$password"
}

# Get user ID from token
get_user_id() {
  local token=$1
  
  local response=$(curl -s -X GET "$BASE_URL/profile" \
    -H "Authorization: Bearer $token")
  
  echo "$response" | jq -r '.data.user_id // empty'
}

# Get user role from token
get_user_role() {
  local token=$1
  
  local response=$(curl -s -X GET "$BASE_URL/profile" \
    -H "Authorization: Bearer $token")
  
  echo "$response" | jq -r '.data.role // empty'
}
