#!/bin/bash

source "$(dirname "$0")/../helpers/color_output.sh"

# Validate API response success field
validate_success() {
  local response=$1
  local expected=$2
  local step_name=$3
  
  local success=$(echo "$response" | jq -r '.success // false')
  
  if [ "$success" != "$expected" ]; then
    print_error "$step_name failed!"
    echo "Expected success: $expected, got: $success"
    echo "Response: $response" | jq '.'
    return 1
  fi
  
  print_success "$step_name"
  return 0
}

# Validate field exists in response (non-strict, allows empty)
validate_field_exists() {
  local response=$1
  local field_path=$2
  local field_name=$3
  
  local value=$(echo "$response" | jq -r "$field_path // \"__NULL__\"")
  
  if [ "$value" == "__NULL__" ] || [ "$value" == "null" ]; then
    print_error "Field '$field_name' not found"
    echo "Response: $response" | jq '.'
    return 1
  fi
  
  if [ -z "$value" ]; then
    print_warning "Field '$field_name' is empty"
    return 0  # Don't fail on empty, just warn
  fi
  
  return 0
}

# Validate field value matches expected
validate_field_value() {
  local response=$1
  local field_path=$2
  local expected=$3
  local field_name=$4
  
  local actual=$(echo "$response" | jq -r "$field_path // empty")
  
  if [ "$actual" != "$expected" ]; then
    print_error "Field '$field_name' validation failed"
    echo "Expected: $expected, got: $actual"
    return 1
  fi
  
  return 0
}

# Validate HTTP status code
validate_http_status() {
  local status=$1
  local expected=$2
  local step_name=$3
  
  if [ "$status" != "$expected" ]; then
    print_error "$step_name - HTTP status mismatch"
    echo "Expected: $expected, got: $status"
    return 1
  fi
  
  return 0
}

# Validate not empty
validate_not_empty() {
  local value=$1
  local field_name=$2
  
  if [ -z "$value" ] || [ "$value" == "null" ]; then
    print_error "'$field_name' is empty or null"
    return 1
  fi
  
  return 0
}

# Validate array length
validate_array_length() {
  local response=$1
  local field_path=$2
  local expected_length=$3
  local field_name=$4
  
  local actual_length=$(echo "$response" | jq "$field_path | length")
  
  if [ "$actual_length" != "$expected_length" ]; then
    print_error "Array '$field_name' length mismatch"
    echo "Expected: $expected_length, got: $actual_length"
    return 1
  fi
  
  return 0
}
