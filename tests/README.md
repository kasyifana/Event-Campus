# API Flow Tests

Automated flow-based testing for Event Campus API based on realistic user journeys.

## 📁 Structure

```
tests/
├── flow_tests/          # Flow-based test scripts
│   ├── flow1_user_onboarding.sh
│   ├── flow2_become_organizer.sh
│   ├── flow3_create_event_lifecycle.sh
│   └── flow4_event_registration_flow.sh
├── helpers/             # Reusable helper functions
│   ├── auth_helper.sh
│   ├── validation_helper.sh
│   └── color_output.sh
└── run_all_flows.sh     # Master test runner
```

## 🚀 Quick Start

### Prerequisites

1. **API Server Running**
   ```bash
   go run cmd/api/main.go
   # Server should be running on http://localhost:8080
   ```

2. **Install jq** (JSON processor)
   ```bash
   # macOS
   brew install jq
   
   # Linux
   sudo apt-get install jq
   ```

3. **Database Setup**
   - PostgreSQL running
   - Migrations applied
   - Admin account seeded (for Flow 2)

### Run All Tests

```bash
cd tests
./run_all_flows.sh
```

### Run Individual Flow

```bash
cd tests/flow_tests
./flow1_user_onboarding.sh
```

## 📋 Test Flows

### Flow 1: User Onboarding
**Journey:** New student browses and explores events

**Steps:**
1. Register account with UII email
2. Login and get JWT token
3. View user profile
4. Browse available events
5. Search for specific events  
6. View event details

**Validates:**
- User registration works
- JWT token generation
- Profile retrieval
- Event listing and filtering

---

### Flow 2: Become Organizer
**Journey:** Student organization gets approval to create events

**Steps:**
1. Register as mahasiswa
2. Submit whitelist request with document
3. Check request status
4. Admin reviews and approves
5. User role upgraded to organisasi
6. Verify new event creation permission

**Validates:**
- Whitelist request submission
- Document upload
- Admin approval workflow
- Role upgrade
- Permission changes

---

### Flow 3: Create Event Lifecycle
**Journey:** Organisasi creates and publishes event

**Steps:**
1. Create event in draft status
2. Upload poster image
3. Update event details
4. Publish event
5. View my events
6. Verify in public event listing

**Validates:**
- Event creation
- File upload
- Event updates
- Draft → Published workflow
- Event visibility

---

### Flow 4: Event Registration
**Journey:** Students register for events

**Steps:**
1. Register new mahasiswa account
2. Browse available events
3. Register for event
4. Receive confirmation
5. View my registrations
6. Verify duplicate prevention

**Validates:**
- Event registration
- Email notifications
- Registration listing
- Duplicate prevention

## 🎯 Usage Examples

### Set Custom Base URL

```bash
BASE_URL=http://your-server.com/api/v1 ./run_all_flows.sh
```

### Run with Admin Credentials

```bash
ADMIN_EMAIL=your-admin@example.com \
ADMIN_PASSWORD=yourpassword \
./run_all_flows.sh
```

### Debug Individual Flow

```bash
# Add set -x for verbose output
bash -x flow_tests/flow1_user_onboarding.sh
```

## ✅ Success Criteria

All tests pass when:
- ✅ All API endpoints respond correctly
- ✅ Database operations succeed
- ✅ JWT tokens are valid
- ✅ File uploads work
- ✅ Email notifications sent (check SMTP logs)
- ✅ Role-based access control enforced

## 🐛 Troubleshooting

### Server Not Running

```
❌ API server is not running!
```

**Solution:** Start the server
```bash
go run cmd/api/main.go
```

### jq Not Found

```
jq: command not found
```

**Solution:** Install jq
```bash
brew install jq  # macOS
```

### Admin Account Missing

```
⚠️ Admin account not found, skipping admin approval steps
```

**Solution:** Seed admin account or create manually
```sql
INSERT INTO users (email, password_hash, role, is_approved)
VALUES ('admin@eventcampus.com', '$hashed_password', 'admin', true);
```

### File Upload Fails

```
⚠️ Poster upload failed
```

**Solution:** Ensure ImageMagick installed or use valid JPG/PNG file
```bash
brew install imagemagick
```

## 📊 Test Results

Example output:

```
╔════════════════════════════════════════════╗
║            Test Results Summary            ║
╚════════════════════════════════════════════╝

✅ flow1_user_onboarding.sh - PASSED
✅ flow2_become_organizer.sh - PASSED
✅ flow3_create_event_lifecycle.sh - PASSED
✅ flow4_event_registration_flow.sh - PASSED

──────────────────────────────────────────
Total Tests: 4
Passed: 4
──────────────────────────────────────────

╔════════════════════════════════════════════╗
║      All Tests Passed Successfully! 🎉     ║
╚════════════════════════════════════════════╝
```

## 🔧 Adding New Flow Tests

1. Create new script in `flow_tests/`
2. Follow naming convention: `flowN_descriptive_name.sh`
3. Use helper functions from `helpers/`
4. Add to `FLOWS` array in `run_all_flows.sh`

Example template:

```bash
#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../helpers/color_output.sh"
source "$SCRIPT_DIR/../helpers/auth_helper.sh"
source "$SCRIPT_DIR/../helpers/validation_helper.sh"

print_info "Flow N: Your Flow Name"

# Your test steps here

print_success "Flow N Completed ✅"
exit 0
```

## 📝 Notes

- Tests create temporary users with unique emails (timestamp-based)
- Some tests depend on previous flows (e.g., Flow 4 needs event from Flow 3)
- Cleanup is automatic for temp files
- Database records persist (good for manual inspection)

## 🔗 Related Documentation

- [API Documentation](../API_DOCUMENTATION.md)
- [Testing Guide](../TESTING_GUIDE.md)
- [Implementation Plan](../.gemini/antigravity/brain/*/implementation_plan.md)
