# Event Campus MVP - Next Steps & Quick Start

## 🎉 What's Been Completed

✅ **Authentication System** (Fully Functional)  
✅ **Project Structure** (Clean Architecture)  
✅ **Database Schema** (Complete SQL)  
✅ **Email Templates** (7 templates ready)  
✅ **Deployment Setup** (Docker + Jenkins)  
✅ **Documentation** (Comprehensive README)

**Build Status**: ✅ Successful  
**Progress**: ~30% of MVP Complete

---

## 🚀 Quick Start Guide

### 1. Setup Database (5 minutes)

1. Go to [Supabase](https://supabase.com) and create a new project
2. Copy your project URL and keys:
   - URL: `https://xxxxx.supabase.co`
   - Anon key: Public API key
   - Service key: Secret admin key
3. Open Supabase SQL Editor
4. Copy & paste the content from `migrations/001_initial_schema.sql`
5. Click "Run" to create all tables

✅ **Default Admin Created**: `admin@uii.ac.id` / `admin123`

### 2. Configure Environment (2 minutes)

Edit `.env.example` and save as `.env`:

```env
# Update these values:
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key-here
SUPABASE_SERVICE_KEY=your-service-key-here

# Generate strong secret:
JWT_SECRET=run-this-command-to-generate: openssl rand -base64 32

# Gmail App Password (optional for now):
SMTP_USER=youremail@gmail.com
SMTP_PASSWORD=your-app-password
```

### 3. Run the Application (1 minute)

```bash
# Start the server
make run

# OR
go run cmd/api/main.go
```

✅ Server runs at: `http://localhost:8080`

### 4. Test Authentication (2 minutes)

**Register a new user:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@uii.ac.id",
    "password": "password123",
    "full_name": "Ahmad Rizki",
    "phone_number": "081234567890"
  }'
```

✅ You should receive a JWT token and user info

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@uii.ac.id",
    "password": "password123"
  }'
```

**Access Protected Route:**
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 📝 Next Implementation Steps

### Priority 1: Complete Core Repositories (4-6 hours)

You need to implement actual Supabase integration, currently using in-memory storage.

**Files to modify:**

1. **internal/repository/user_repository.go**
   - Replace in-memory maps with Supabase client
   - Use `supabase.From("users").Select()` etc.

2. **Create: internal/repository/event_repository.go**
   ```go
   - Create()
   - GetByID()
   - List() with filters
   - Update()
   - Delete()
   ```

3. **Create: internal/repository/whitelist_repository.go**
   ```go
   - Create()
   - GetByUserID()
   - ListPending()
   - Update()
   ```

4. **Create: internal/repository/registration_repository.go**
   ```go
   - Create()
   - GetByEventAndUser()
   - GetByEvent()
   - Update()
   - GetFirstWaitlist()
   ```

### Priority 2: Event Management (6-8 hours)

1. **internal/usecase/event_usecase.go**
   - CreateEvent() - With poster upload
   - GetEvents() - With filtering
   - UpdateEvent() - With ownership check
   - DeleteEvent() - With validation

2. **internal/delivery/http/handler/event_handler.go**
   - POST /api/v1/events
   - GET /api/v1/events
   - GET /api/v1/events/:id
   - PUT /api/v1/events/:id
   - DELETE /api/v1/events/:id

3. **Wire to router** - Add routes in `router.go`

### Priority 3: Whitelist System (4-6 hours)

1. **internal/usecase/whitelist_usecase.go**
   - SubmitRequest() - With PDF upload
   - GetPendingRequests()
   - ApproveRequest() - Updates user role
   - RejectRequest() - With reason

2. **internal/delivery/http/handler/whitelist_handler.go**
   - POST /api/v1/whitelist/request
   - GET /api/v1/whitelist/requests (admin)
   - PATCH /api/v1/whitelist/:id/approve (admin)
   - PATCH /api/v1/whitelist/:id/reject (admin)

3. **Integrate email notifications**
   - Call `emailSender.SendWhitelistApproval()`
   - Call `emailSender.SendWhitelistRejection()`

### Priority 4: Registration System (6-8 hours)

1. **internal/usecase/registration_usecase.go**
   - RegisterToEvent() - With capacity check
   - CancelRegistration() - With waitlist promotion
   - GetMyRegistrations()

2. **internal/delivery/http/handler/registration_handler.go**
   - POST /api/v1/events/:id/register
   - DELETE /api/v1/registrations/:id
   - GET /api/v1/registrations/my

### Priority 5: Schedulers (3-4 hours)

1. **internal/scheduler/reminder_scheduler.go**
   ```go
   func StartReminderScheduler(repo, emailSender) {
       c := cron.New()
       c.AddFunc("0 0 * * *", func() {
           // Find events starting tomorrow
           // Send reminder emails with zoom links
       })
       c.Start()
   }
   ```

2. **internal/scheduler/event_status_updater.go**
   ```go
   func StartStatusUpdater(repo) {
       c := cron.New()
       c.AddFunc("0 * * * *", func() {
           // Update event statuses
       })
       c.Start()
   }
   ```

3. **Wire to main.go**
   ```go
   go scheduler.StartReminderScheduler(...)
   go scheduler.StartStatusUpdater(...)
   ```

---

## 🔧 Development Workflow

### Adding a New Feature

1. **Domain Model** (if needed)
   - Add to `internal/domain/`

2. **Repository**
   - Interface in `internal/repository/`
   - Implementation with Supabase client

3. **Use Case**
   - Business logic in `internal/usecase/`

4. **Handler**
   - HTTP endpoints in `internal/delivery/http/handler/`

5. **Router**
   - Wire routes in `internal/delivery/http/router/router.go`

6. **Test**
   - Test with curl or Postman

### Example: Adding Attendance Feature

```go
// 1. Repository
type AttendanceRepository interface {
    Create(ctx, attendance) error
    GetByEvent(ctx, eventID) ([]Attendance, error)
}

// 2. Use Case
func (u *attendanceUsecase) MarkAttendance(registrationID) error {
    // Verify registration exists
    // Create attendance record
    // Update registration status
}

// 3. Handler
func (h *attendanceHandler) MarkAttendance(c *gin.Context) {
    // Parse request
    // Call use case
    // Return response
}

// 4. Router
protected.POST("/attendances", attendanceHandler.MarkAttendance)
```

---

## 📊 Estimated Timeline

| Phase | Task | Hours | Priority |
|-------|------|-------|----------|
| 3 | Whitelist System | 4-6 | High |
| 4 | Event Management | 6-8 | High |
| 5 | Registration System | 6-8 | High |
| 6 | Attendance System | 3-4 | Medium |
| 7-8 | Schedulers | 3-4 | Medium |
| 9 | Testing | 4-6 | High |
| DB | Supabase Integration | 4-6 | Critical |

**Total**: 30-42 hours (~1 week full-time)

---

## 🐛 Common Issues & Solutions

### Issue: "User not found" after registration
**Solution**: Currently using in-memory storage. Implement Supabase repository to persist data.

### Issue: Email not sending
**Solution**: 
1. Enable 2FA on Gmail
2. Generate App Password: https://myaccount.google.com/apppasswords
3. Use app password in `.env`, not your Gmail password

### Issue: JWT token invalid
**Solution**: Make sure JWT_SECRET in `.env` matches between registration and login.

### Issue: CORS errors
**Solution**: Add your frontend URL to `ALLOWED_ORIGINS` in `.env`

---

## 📚 Useful Commands

```bash
# Development
make run                    # Run server
make build                  # Build binary
make test                   # Run tests
make clean                  # Clean build artifacts

# Docker
make docker-build           # Build Docker image
make docker-run             # Run in container

# Database
# Execute in Supabase SQL Editor:
cat migrations/001_initial_schema.sql

# Dependencies
go mod tidy                 # Clean dependencies
go get <package>            # Add new dependency
```

---

## 🎯 Success Criteria for MVP

Before considering MVP complete, ensure:

- ✅ Authentication (register, login) **DONE**
- [ ] Whitelist approval system
- [ ] Event CRUD for organisasi
- [ ] Event registration with capacity
- [ ] Waitlist system
- [ ] Attendance marking
- [ ] H-1 reminder emails
- [ ] Event status auto-update
- [ ] All email notifications working

---

## 💡 Tips for Success

### 1. Start with Supabase Integration
Replace the in-memory repository first. Everything else depends on it.

### 2. Test Each Layer Independently
- Test repository queries in Supabase directly
- Test use cases with unit tests
- Test handlers with curl

### 3. Use Sample Data
Create test events and users directly in Supabase to speed up testing.

### 4. Follow the Pattern
Look at `auth_usecase.go` and `auth_handler.go` as templates for new features.

### 5. Commit Often
```bash
git add .
git commit -m "feat: implement event creation"
git push
```

---

## 📞 Need Help?

**Check these files:**
- [README.md](file:///Users/user/Campuss/Semester%205/BSI/Event%20Campus/README.md) - API documentation
- [Implementation Plan](file:///Users/user/.gemini/antigravity/brain/6c40377e-7991-4a35-8541-31cb26c12b60/implementation_plan.md) - Detailed specs
- [Walkthrough](file:///Users/user/.gemini/antigravity/brain/6c40377e-7991-4a35-8541-31cb26c12b60/walkthrough.md) - What's completed

**Test the APIs:**
Import into Postman or use the curl examples in README.md

---

**You're ready to continue!** The foundation is solid. Start with Supabase integration, then add features one by one. 🚀
