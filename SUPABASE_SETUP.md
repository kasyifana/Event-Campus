# Setup Supabase Integration - Quick Guide

## ✅ Yang Sudah Dilakukan

Saya sudah mengubah kode dari **in-memory storage** ke **Supabase persistent storage**:

### File yang Diupdate:

1. **`internal/repository/user_repository.go`**
   - ✅ Replaced `map[uuid.UUID]*domain.User` dengan Supabase client
   - ✅ Semua operasi (Create, GetByEmail, GetByID, Update, UpdateRole) sekarang pakai Supabase

2. **`cmd/api/main.go`**
   - ✅ Added Supabase client initialization
   - ✅ Pass Supabase client ke repository

3. **`go.mod`**
   - ✅ Added `github.com/supabase-community/supabase-go` dependency

---

## 🚀 Cara Setup Supabase

### Step 1: Buat Project Supabase (5 menit)

1. Buka https://supabase.com
2. Sign in / Sign up
3. Klik **"New Project"**
4. Isi:
   - Name: `event-campus`
   - Database Password: **Catat password ini!**
   - Region: Southeast Asia (Singapore) - paling dekat
5. Klik **"Create new project"**
6. Tunggu ~2 menit sampai setup selesai

---

### Step 2: Execute Database Schema (2 menit)

1. Di Supabase Dashboard, buka **SQL Editor** (icon ⚡ di sidebar)
2. Klik **"New query"**
3. Copy isi file `migrations/001_initial_schema.sql`
4. Paste ke SQL Editor
5. Klik **"Run"** atau tekan `Ctrl+Enter`
6. Tunggu sampai muncul "Success. No rows returned"

✅ **Tables Created:**
- users
- events  
- registrations
- whitelist_requests
- attendances

✅ **Default Admin Created:**
- Email: `admin@uii.ac.id`
- Password: `admin123`

---

### Step 3: Get API Keys (1 menit)

1. Di Supabase Dashboard, klik **Settings** (⚙️ icon) di sidebar
2. Klik **API**
3. Copy 3 values ini:

**Project URL:**
```
https://xxxxxxxxxx.supabase.co
```

**anon public (public key):**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3...
```

**service_role (secret key):**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3...
```

⚠️ **IMPORTANT**: `service_role` key adalah SECRET! Jangan commit ke git!

---

### Step 4: Update .env File (1 menit)

Buka `.env.example` (atau buat `.env` baru) dan update:

```env
# Supabase Configuration
SUPABASE_URL=https://xxxxxxxxxx.supabase.co        # ← Update dengan Project URL
SUPABASE_ANON_KEY=eyJhbGc...                       # ← Update dengan anon key
SUPABASE_SERVICE_KEY=eyJhbGc...                    # ← Update dengan service_role key

# Yang lain tetap sama
PORT=8080
ENV=development
JWT_SECRET=event-campus-super-secret-key-change-this-in-production
JWT_EXPIRATION=24h
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
MAX_UPLOAD_SIZE=10485760
UPLOAD_PATH=./storage
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

**Save file as `.env`** (tanpa .example)

---

### Step 5: Restart Server

1. **Stop** server yang sedang running (Ctrl+C di terminal)
2. **Start** ulang:
   ```bash
   go run cmd/api/main.go
   ```

3. Lihat output, harus ada:
   ```
   🚀 Starting Event Campus API...
   📍 Environment: development
   ✅ Supabase connected    ← Ini yang penting!
   ✅ Server running on http://localhost:8080
   ```

✅ Jika ada "Supabase connected" → berhasil!  
❌ Jika error → cek credentials di `.env`

---

## 🧪 Test Connection

Test registrasi lagi:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "baru@uii.ac.id",
    "password": "password123",
    "full_name": "User Baru",
    "phone_number": "081234567890"
  }'
```

**Sekarang cek di Supabase:**

1. Buka Supabase Dashboard
2. Klik **Table Editor** (📊 icon)
3. Pilih table **`users`**
4. Lihat data baru muncul! ✅

---

## 🎯 Verifikasi Data Tersimpan

### Di Supabase Dashboard:

1. **Table Editor** → `users`
2. Harusnya ada:
   - Default admin: `admin@uii.ac.id`
   - User yang baru di-register: `baru@uii.ac.id`
   - User sebelumnya: `23523053@students.uii.ac.id`

### Test Persistence:

1. **Stop** server (Ctrl+C)
2. **Start** lagi: `go run cmd/api/main.go`
3. **Login** dengan user yang sama:
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{
       "email": "baru@uii.ac.id",
       "password": "password123"
     }'
   ```

✅ **Login berhasil** = Data persistent di Supabase!  
❌ **Login gagal** = Setup belum benar

---

## 🔍 Debugging

### Error: "Failed to initialize Supabase client"

**Penyebab:** Credentials salah di `.env`

**Solusi:**
1. Cek lagi `SUPABASE_URL` (harus https://xxx.supabase.co)
2. Cek lagi `SUPABASE_SERVICE_KEY` (bukan anon key!)
3. Pastikan no trailing spaces

### Error: "user not found" padahal data ada

**Penyebab:** Query Supabase belum sempurna

**Solusi:** Coba login dengan admin default dulu:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@uii.ac.id",
    "password": "admin123"
  }'
```

Jika berhasil → Supabase OK, masalahnya di data user tertentu

### Data tidak muncul di Supabase

**Cek:**
1. Apakah ada error di terminal saat register?
2. Response dari curl `"success": true`?
3. Refresh table editor di Supabase

---

## 📊 Monitor Database

### Di Supabase Dashboard → Database → Logs:

Bisa lihat semua query yang dijalankan real-time.

### Di Table Editor:

- **Add filter**: Filter berdasarkan email, role, etc.
- **Add column**: Customize tampilan
- **Export**: Download data sebagai CSV

---

## ✅ Checklist

Sebelum lanjut development:

- [ ] Supabase project created
- [ ] Database schema executed (5 tables created)
- [ ] `.env` configured with correct keys
- [ ] Server restart dengan "Supabase connected" log
- [ ] Test registration → data muncul di Supabase
- [ ] Test login → berhasil dengan data dari Supabase
- [ ] Server restart → data tetap ada (persistent)

---

## 🎉 Next Steps

Setelah Supabase terkoneksi, kamu bisa lanjut implement:

1. **Event Management** - Event CRUD operations
2. **Whitelist System** - Organisasi approval
3. **Registration System** - Event registration dengan capacity
4. **Attendance System** - Check-in management

Semua akan langsung tersimpan di Supabase! 🚀

---

## 💡 Tips

1. **Gunakan Table Editor** untuk debug data cepat
2. **Gunakan SQL Editor** untuk query kompleks
3. **Monitor Database Logs** untuk troubleshooting
4. **Backup Database** secara berkala (Settings → Database → Backup)

---

## 🔐 Security Notes

⚠️ **JANGAN COMMIT `.env` ke Git!**

Pastikan `.gitignore` sudah include `.env`:
```
.env
```

✅ Gunakan `.env.example` sebagai template  
✅ Share credentials via secure channel (1Password, LastPass, dll)

---

**Need Help?** Cek [API_DOCUMENTATION.md](API_DOCUMENTATION.md) untuk endpoint details!
