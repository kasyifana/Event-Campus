# 🧪 Testing Guide - Event Campus API

Panduan lengkap untuk testing semua fitur Event Campus Backend API.

## 📋 Prerequisites

1. **Setup Database PostgreSQL**
   ```bash
   # Pastikan PostgreSQL sudah running di localhost:5432
   # Database sudah dibuat sesuai .env
   ```

2. **Setup Environment Variables**
   ```bash
   # Copy .env.example ke .env
   cp .env.example .env
   
   # Edit .env dengan konfigurasi yang benar
   ```

3. **Run Server**
   ```bash
   go run cmd/api/main.go
   # Server akan running di http://localhost:8080
   ```

---

## 🔐 Phase 2: Authentication Testing

### 1. Register User (Mahasiswa)
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mahasiswa@uii.ac.id",
    "password": "password123",
    "full_name": "Ahmad Mahasiswa",
    "phone_number": "081234567890",
    "nim": "21523001",
    "is_uii_civitas": true
  }'
```

**Expected:** Status 201, success message

### 2. Register User (Organisasi - akan ditolak karena perlu whitelist)
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "organisasi@example.com",
    "password": "password123",
    "full_name": "Organisasi Test",
    "phone_number": "081234567891",
    "role": "organisasi"
  }'
```

**Expected:** Status 400 (harus melalui whitelist)

### 3. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mahasiswa@uii.ac.id",
    "password": "password123"
  }'
```

**Expected:** Status 200, JWT token

**💡 Save TOKEN untuk request selanjutnya:**
```bash
export TOKEN="<jwt_token_dari_response>"
```

### 4. Get Profile
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $TOKEN"
```

**Expected:** Status 200, user profile data

---

## 📝 Phase 3: Whitelist System Testing

### 1. Submit Whitelist Request (as Mahasiswa)
```bash
# Siapkan file PDF dummy
echo "Test Document" > /tmp/test_document.pdf

curl -X POST http://localhost:8080/api/v1/whitelist/request \
  -H "Authorization: Bearer $TOKEN" \
  -F "organization_name=HMTI UII" \
  -F "organization_type=Himpunan Mahasiswa" \
  -F "pic_name=Ahmad" \
  -F "pic_position=Ketua" \
  -F "pic_phone=081234567890" \
  -F "document=@/tmp/test_document.pdf"
```

**Expected:** Status 201, whitelist request created

### 2. Get My Whitelist Request
```bash
curl -X GET http://localhost:8080/api/v1/whitelist/my-request \
  -H "Authorization: Bearer $TOKEN"
```

**Expected:** Status 200, request data with status "pending"

### 3. Admin: Get All Whitelist Requests
```bash
# Login sebagai admin terlebih dahulu
# Default admin: admin@eventcampus.com / admin123

curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@eventcampus.com",
    "password": "admin123"
  }'

export ADMIN_TOKEN="<admin_jwt_token>"

curl -X GET http://localhost:8080/api/v1/whitelist/requests \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**Expected:** Status 200, list of all requests

### 4. Admin: Approve Whitelist Request
```bash
curl -X PATCH http://localhost:8080/api/v1/whitelist/<REQUEST_ID>/review \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved",
    "admin_notes": "Approved - valid organization"
  }'
```

**Expected:** Status 200, user role upgraded to "organisasi", email sent

### 5. Verify Role Change
```bash
# Login ulang dengan user yang di-approve
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mahasiswa@uii.ac.id",
    "password": "password123"
  }'
```

**Expected:** Token dengan role "organisasi"

---

## 🎉 Phase 4: Event Management Testing

### 1. Create Event (as Organisasi)
```bash
export ORG_TOKEN="<organisasi_jwt_token>"

curl -X POST http://localhost:8080/api/v1/events \
  -H "Authorization: Bearer $ORG_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Workshop Web Development",
    "description": "Belajar membuat website modern dengan React dan Go",
    "category": "workshop",
    "event_type": "online",
    "zoom_link": "https://zoom.us/j/123456789",
    "start_date": "2025-12-01T14:00:00Z",
    "end_date": "2025-12-01T17:00:00Z",
    "registration_deadline": "2025-11-30T23:59:59Z",
    "max_participants": 50,
    "is_uii_only": true,
    "status": "draft"
  }'
```

**Expected:** Status 201, event created

**💡 Save EVENT_ID:**
```bash
export EVENT_ID="<event_id_from_response>"
```

### 2. Upload Event Poster
```bash
# Buat dummy image
convert -size 800x600 xc:blue /tmp/poster.jpg
# Atau gunakan image yang ada

curl -X POST http://localhost:8080/api/v1/events/$EVENT_ID/poster \
  -H "Authorization: Bearer $ORG_TOKEN" \
  -F "poster=@/tmp/poster.jpg"
```

**Expected:** Status 200, poster uploaded

### 3. Publish Event
```bash
curl -X POST http://localhost:8080/api/v1/events/$EVENT_ID/publish \
  -H "Authorization: Bearer $ORG_TOKEN"
```

**Expected:** Status 200, event published

### 4. Get All Events (Public)
```bash
curl -X GET "http://localhost:8080/api/v1/events" \
  -H "Authorization: Bearer $TOKEN"
```

**Expected:** Status 200, list of published events

### 5. Filter Events
```bash
curl -X GET "http://localhost:8080/api/v1/events?category=workshop&search=web" \
  -H "Authorization: Bearer $TOKEN"
```

**Expected:** Status 200, filtered events

### 6. Get Event Detail
```bash
curl -X GET http://localhost:8080/api/v1/events/$EVENT_ID \
  -H "Authorization: Bearer $TOKEN"
```

**Expected:** Status 200, event detail dengan poster URL

---

## 🎫 Phase 5: Registration System Testing

### 1. Register for Event
```bash
curl -X POST http://localhost:8080/api/v1/events/$EVENT_ID/register \
  -H "Authorization: Bearer $TOKEN"
```

**Expected:** 
- Status 201 jika ada slot → email konfirmasi dikirim
- Status with waitlist message jika penuh

### 2. Get My Registrations
```bash
curl -X GET http://localhost:8080/api/v1/registrations/my \
  -H "Authorization: Bearer $TOKEN"
```

**Expected:** Status 200, list of user's registrations

### 3. Test Waitlist (Register banyak user hingga penuh)
```bash
# Register 50+ users untuk mengisi kapasitas
# User ke-51 akan masuk waitlist

for i in {1..51}; do
  curl -X POST http://localhost:8080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"user$i@test.com\",
      \"password\": \"pass123\",
      \"full_name\": \"User $i\",
      \"phone_number\": \"08123456789$i\",
      \"is_uii_civitas\": true
    }"
  
  # Login
  LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"user$i@test.com\", \"password\": \"pass123\"}")
  
  USER_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
  
  # Register for event
  curl -X POST http://localhost:8080/api/v1/events/$EVENT_ID/register \
    -H "Authorization: Bearer $USER_TOKEN"
done
```

### 4. Cancel Registration (Test Waitlist Promotion)
```bash
# Ambil satu registration ID
export REG_ID="<registration_id>"

curl -X DELETE http://localhost:8080/api/v1/registrations/$REG_ID \
  -H "Authorization: Bearer $TOKEN"
```

**Expected:** 
- Status 200
- User pertama di waitlist otomatis dipromosikan
- Email promotion dikirim

### 5. Organizer: View Event Registrations
```bash
curl -X GET http://localhost:8080/api/v1/events/$EVENT_ID/registrations \
  -H "Authorization: Bearer $ORG_TOKEN"
```

**Expected:** Status 200, list of all registrations

---

## ✅ Phase 6: Attendance System Testing

### 1. Mark Single Attendance
```bash
# Ambil user_id dari registrations
export PARTICIPANT_ID="<user_id>"

curl -X POST http://localhost:8080/api/v1/events/$EVENT_ID/attendance \
  -H "Authorization: Bearer $ORG_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "'$PARTICIPANT_ID'",
    "notes": "Hadir tepat waktu"
  }'
```

**Expected:** Status 200, attendance marked

### 2. Bulk Mark Attendance
```bash
curl -X POST http://localhost:8080/api/v1/events/$EVENT_ID/attendance/bulk \
  -H "Authorization: Bearer $ORG_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_ids": ["<user_id_1>", "<user_id_2>", "<user_id_3>"]
  }'
```

**Expected:** Status 200, multiple attendance marked

### 3. Get Event Attendance List
```bash
curl -X GET http://localhost:8080/api/v1/events/$EVENT_ID/attendance \
  -H "Authorization: Bearer $ORG_TOKEN"
```

**Expected:** Status 200, list of attendees

---

## 📧 Phase 7: Email Testing

### Verify Email Delivery

1. **Registration Confirmation Email**
   - Register untuk event
   - Check email inbox untuk konfirmasi
   - Verify: subject, event details, registration ID

2. **Waitlist Notification Email**
   - Register saat event full
   - Check email untuk notifikasi waitlist
   - Verify: posisi di waitlist

3. **Waitlist Promotion Email**
   - Cancel satu registration
   - Check email user pertama di waitlist
   - Verify: promotion message, zoom link

4. **Whitelist Approval Email**
   - Admin approve whitelist request
   - Check email applicant
   - Verify: approval message, new permissions

5. **Whitelist Rejection Email**
   - Admin reject whitelist request
   - Check email applicant
   - Verify: rejection reason

---

## ⏰ Phase 8: Scheduler Testing

### 1. Test H-1 Reminder (Manual Trigger)
```bash
# Buat event dengan start_date besok
# Wait atau trigger scheduler manually

# Check logs untuk:
# - "🔔 Running H-1 reminder job..."
# - "✅ H-1 reminders sent: X"

# Check email inbox untuk reminder
# Verify: zoom link included
```

### 2. Test Event Status Auto-Updater
```bash
# Create events dengan different timestamps:
# - Event A: start_date = now - 1 hour (should be ongoing)
# - Event B: end_date = now - 1 hour (should be completed)

# Wait untuk scheduler run atau trigger manually

# Verify status changes:
curl -X GET http://localhost:8080/api/v1/events/<EVENT_A_ID> \
  -H "Authorization: Bearer $TOKEN"

# Expected: status = "ongoing" atau "completed"
```

---

## 🔄 End-to-End Testing Scenarios

### Scenario 1: Complete Event Lifecycle
1. ✅ Mahasiswa register account
2. ✅ Submit whitelist request
3. ✅ Admin approve whitelist (role → organisasi)
4. ✅ Organisasi create event (draft)
5. ✅ Upload poster
6. ✅ Publish event
7. ✅ Multiple users register (some waitlisted)
8. ✅ One user cancels (waitlist promotion)
9. ✅ H-1 reminder sent
10. ✅ Event status auto-updates (ongoing → completed)
11. ✅ Organizer marks attendance

### Scenario 2: Capacity & Waitlist Flow
1. ✅ Create event with max_participants = 3
2. ✅ 5 users register (3 confirmed, 2 waitlist)
3. ✅ User 1 cancels
4. ✅ User 4 (first in waitlist) gets promoted
5. ✅ User 5 still in waitlist
6. ✅ Verify emails sent correctly

### Scenario 3: Permission Testing
1. ✅ Mahasiswa tries to create event → 403
2. ✅ Organisasi tries to edit other's event → 403
3. ✅ Non-admin tries to approve whitelist → 403
4. ✅ Non-organizer tries to mark attendance → 403

---

## ✅ Testing Checklist

### Authentication & Authorization
- [ ] ✅ User registration works
- [ ] ✅ Login returns valid JWT
- [ ] ✅ Protected routes require token
- [ ] ✅ Role-based access control works
- [ ] ✅ Token expiration handled

### Whitelist System
- [ ] ✅ Mahasiswa can submit request
- [ ] ✅ File upload works (PDF validation)
- [ ] ✅ Admin can view all requests
- [ ] ✅ Approval upgrades role to organisasi
- [ ] ✅ Rejection sends email with reason
- [ ] ✅ Duplicate requests prevented

### Event Management
- [ ] ✅ Organisasi can create events
- [ ] ✅ Poster upload works (JPG/PNG)
- [ ] ✅ Draft → Published workflow
- [ ] ✅ Filtering & search works
- [ ] ✅ Only organizer can edit own events
- [ ] ✅ Poster URL generated correctly

### Registration System
- [ ] ✅ Registration creates record
- [ ] ✅ Capacity check works
- [ ] ✅ Auto-waitlist when full
- [ ] ✅ Cancellation triggers promotion (FIFO)
- [ ] ✅ UII-only validation works
- [ ] ✅ Registration deadline enforced
- [ ] ✅ Duplicate registration prevented

### Attendance System
- [ ] ✅ Single attendance marking works
- [ ] ✅ Bulk operations work (transaction-safe)
- [ ] ✅ Registration status updated to 'attended'
- [ ] ✅ Only organizer can mark
- [ ] ✅ Event must have started

### Email Notifications
- [ ] ✅ Registration confirmation sent
- [ ] ✅ Waitlist notification sent
- [ ] ✅ Promotion email sent
- [ ] ✅ Cancellation confirmation sent
- [ ] ✅ Whitelist approval sent
- [ ] ✅ Whitelist rejection sent
- [ ] ✅ H-1 reminder sent with zoom link

### Automated Schedulers
- [ ] ✅ H-1 reminder runs daily at 9AM
- [ ] ✅ Status updater runs hourly
- [ ] ✅ published → ongoing transition
- [ ] ✅ ongoing → completed transition
- [ ] ✅ Reminder marked as sent

---

## 🐛 Common Issues & Solutions

### Issue: Database connection failed
**Solution:**
```bash
# Check PostgreSQL is running
pg_isready

# Check credentials in .env
psql -U postgres -d event_campus
```

### Issue: Email not sending
**Solution:**
```bash
# Check SMTP configuration in .env
# Test with Mailtrap.io or Gmail SMTP
# Check application logs for email errors
```

### Issue: File upload fails
**Solution:**
```bash
# Check upload directory exists and writable
mkdir -p uploads/posters uploads/documents
chmod 755 uploads
```

### Issue: JWT token invalid
**Solution:**
```bash
# Token might be expired (check JWT_EXPIRATION in .env)
# Login again to get fresh token
# Verify JWT_SECRET matches
```

---

## 📊 Performance Testing

### Load Test Registration (Optional)
```bash
# Install Apache Bench
brew install ab  # macOS
# atau apt-get install apache2-utils  # Linux

# Test 100 concurrent requests
ab -n 100 -c 10 -H "Authorization: Bearer $TOKEN" \
  -p event_reg.json -T application/json \
  http://localhost:8080/api/v1/events/$EVENT_ID/register
```

---

## ✅ Final Verification

Setelah semua testing selesai, verify:

1. **Database Integrity**
   ```sql
   -- Check data counts
   SELECT 
     (SELECT COUNT(*) FROM users) as users,
     (SELECT COUNT(*) FROM events) as events,
     (SELECT COUNT(*) FROM registrations) as registrations,
     (SELECT COUNT(*) FROM attendances) as attendances;
   ```

2. **Application Logs**
   - No error messages
   - Scheduler running correctly
   - All emails sent successfully

3. **API Response Times**
   - All endpoints respond < 500ms
   - Database queries optimized

**Congratulations! MVP Testing Complete! 🎉**
