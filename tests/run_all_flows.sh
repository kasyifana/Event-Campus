#!/bin/bash

# Master Test Runner - Run All API Flow Tests
# This script runs all flow tests in sequence and reports results

set +e  # Don't exit on error, we want to run all tests

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/helpers/color_output.sh"

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"
export BASE_URL

echo ""
print_info "╔════════════════════════════════════════════╗"
print_info "║   Event Campus - API Flow Tests Runner    ║"
print_info "╚════════════════════════════════════════════╝"
echo ""

# Check if server is running
print_step "Checking if API server is running..."
HEALTH_CHECK=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)

if [ "$HEALTH_CHECK" != "200" ]; then
  print_error "API server is not running!"
  print_info "Please start the server with: go run cmd/api/main.go"
  exit 1
fi

print_success "API server is running"
echo ""

# Check if jq is installed
if ! command -v jq &> /dev/null; then
  print_error "jq is not installed!"
  print_info "Install with: brew install jq (macOS) or apt-get install jq (Linux)"
  exit 1
fi

print_success "Dependencies OK"
echo ""

# List of flow tests to run
FLOWS=(
  "flow1_user_onboarding.sh"
  "flow2_become_organizer.sh"
  "flow3_create_event_lifecycle.sh"
  "flow4_event_registration_flow.sh"
  "flow5_attendance.sh"
  "flow6_advanced_scenarios.sh"
)

# Counters
TOTAL=${#FLOWS[@]}
PASSED=0
FAILED=0
SKIPPED=0

# Results array
declare -a RESULTS

# Run each flow
for flow in "${FLOWS[@]}"; do
  FLOW_PATH="$SCRIPT_DIR/flow_tests/$flow"
  
  if [ ! -f "$FLOW_PATH" ]; then
    print_warning "Flow test not found: $flow"
    ((SKIPPED++))
    RESULTS+=("⚠️  $flow - SKIPPED (not found)")
    continue
  fi
  
  # Make executable
  chmod +x "$FLOW_PATH"
  
  echo ""
  print_info "──────────────────────────────────────────"
  print_info "Running: $flow"
  print_info "──────────────────────────────────────────"
  echo ""
  
  # Run the flow test
  if bash "$FLOW_PATH"; then
    ((PASSED++))
    RESULTS+=("✅ $flow - PASSED")
  else
    ((FAILED++))
    RESULTS+=("❌ $flow - FAILED")
  fi
  
  echo ""
done

# Print summary
echo ""
print_info "╔════════════════════════════════════════════╗"
print_info "║            Test Results Summary            ║"
print_info "╚════════════════════════════════════════════╝"
echo ""

for result in "${RESULTS[@]}"; do
  echo "$result"
done

echo ""
print_info "──────────────────────────────────────────"
print_info "Total Tests: $TOTAL"

if [ $PASSED -gt 0 ]; then
  print_success "Passed: $PASSED"
fi

if [ $FAILED -gt 0 ]; then
  print_error "Failed: $FAILED"
fi

if [ $SKIPPED -gt 0 ]; then
  print_warning "Skipped: $SKIPPED"
fi

print_info "──────────────────────────────────────────"
echo ""

# Clean up temp files
# rm -f /tmp/test_event_id.sh /tmp/test_registration_id.sh

# Exit with appropriate code
if [ $FAILED -eq 0 ]; then
  print_success "╔════════════════════════════════════════════╗"
  print_success "║      All Tests Passed Successfully! 🎉     ║"
  print_success "╚════════════════════════════════════════════╝"
  echo ""
  exit 0
else
  print_error "╔════════════════════════════════════════════╗"
  print_error "║        Some Tests Failed ❌                 ║"
  print_error "╚════════════════════════════════════════════╝"
  echo ""
  exit 1
fi
