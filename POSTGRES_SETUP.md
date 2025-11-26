# Quick Setup: PostgreSQL Supabase Integration

## ✅ What's Done

- ✅ PostgreSQL driver installed (`github.com/lib/pq`)
- ✅ Repository implemented with SQL queries
- ✅ Main.go updated with database connection
- ✅ Build successful

## 🔧 Setup Database Credentials

### Step 1: Update `.env` File

Add these lines to your `.env` file:

```env
# PostgreSQL Configuration (from your Supabase project)
POSTGRES_HOST=aws-1-ap-northeast-1.pooler.supabase.com
POSTGRES_PORT=6543
POSTGRES_DB=postgres
POSTGRES_USER=postgres.qlrjneqefcqgukylxkuh
POSTGRES_PASSWORD=YOUR_DATABASE_PASSWORD_HERE
POSTGRES_SSLMODE=require
```

**⚠️ IMPORTANT:** Replace `YOUR_DATABASE_PASSWORD_HERE` with your actual Supabase database password!

### Step 2: Get Your Database Password

If you don't remember your password:

1. Go to Supabase Dashboard
2. Click **Settings** → **Database**
3. Look for "Connection string" section
4. Password is shown there (or reset if needed)

OR use the connection string you provided:
```
postgresql://postgres.qlrjneqefcqgukylxkuh:[YOUR-PASSWORD]@aws-1-ap-northeast-1.pooler.supabase.com:6543/postgres
```

The format is: `postgresql://USER:PASSWORD@HOST:PORT/DATABASE`

### Step 3: Start Server

```bash
go run cmd/api/main.go
```

Expected Output:
```
🚀 Starting Event Campus API...
📍 Environment: development
✅ PostgreSQL connected to Supabase    ← This confirms database connected!
✅ Server running on http://localhost:8080
```

### Step 4: Test Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@uii.ac.id",
    "password": "password123",
    "full_name": "Test User",
    "phone_number": "081234567890"
  }'
```

### Step 5: Verify in Supabase

1. Open Supabase Dashboard
2. Go to **Table Editor**
3. Select **`users`** table
4. See your new user! ✅

---

## 🔍 Connection Details

Based on your Supabase info:

```
Host: aws-1-ap-northeast-1.pooler.supabase.com
Port: 6543
Database: postgres
User: postgres.qlrjneqefcqgukylxkuh
Pool Mode: transaction
SSL Mode: require
```

---

## ❌ Troubleshooting

### Error: "POSTGRES_PASSWORD is not set"

**Solution:** Add `POSTGRES_PASSWORD=your-password` to `.env`

### Error: "Failed to connect to database"

**Possible causes:**
1. Wrong password
2. Wrong host/port
3. Firewall blocking connection
4. Supabase project paused/suspended

**Solution:** 
- Check credentials in Supabase Dashboard → Settings → Database
- Try connection pooler URL

### Error: "Failed to ping database: pq: password authentication failed"

**Solution:** Password is incorrect. Reset in Supabase Dashboard.

### Data not showing in Supabase

**Check:**
1. Registration returned `"success": true`?
2. No error in server logs?
3. Refresh Table Editor in Supabase
4. Run: `SELECT * FROM users;` in SQL Editor

---

## ✅ Verification Checklist

- [ ] `.env` file has `POSTGRES_PASSWORD` configured
- [ ] Server starts with "PostgreSQL connected to Supabase" message
- [ ] Registration curl command returns success
- [ ] Data appears in Supabase Table Editor
- [ ] Login works with registered user
- [ ] Server restart doesn't lose data (persistent!)

---

## 🎉 When Everything Works

You should be able to:

1. ✅ Register users → Saved to Supabase
2. ✅ Login with users → Retrieved from Supabase
3. ✅ Restart server → Data persists
4. ✅ View data in Supabase Dashboard
5. ✅ Query data with SQL in Supabase

**No more in-memory storage!** Everything is now persistent in PostgreSQL (Supabase).

---

## 📝 Example .env File

```env
# Server
PORT=8080
ENV=development

# Supabase
SUPABASE_URL=https://qlrjneqefcqgukylxkuh.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key

# PostgreSQL (Supabase Database)
POSTGRES_HOST=aws-1-ap-northeast-1.pooler.supabase.com
POSTGRES_PORT=6543
POSTGRES_DB=postgres
POSTGRES_USER=postgres.qlrjneqefcqgukylxkuh
POSTGRES_PASSWORD=your-actual-database-password
POSTGRES_SSLMODE=require

# JWT
JWT_SECRET=event-campus-super-secret-key
JWT_EXPIRATION=24h

# Email (optional for now)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# File Upload
MAX_UPLOAD_SIZE=10485760
UPLOAD_PATH=./storage

# CORS
ALLOWED_ORIGINS=http://localhost:3000
```

---

**Ready? Update `.env` dengan database password dan start server!** 🚀
