# Event Campus - Backend API

Platform manajemen event kampus dengan sistem role-based access control, automated notifications, dan attendance tracking.

## 🚀 Tech Stack

- **Backend**: Go 1.21+ dengan Gin Framework
- **Database**: Supabase (PostgreSQL)
- **Authentication**: JWT-based auth
- **File Storage**: Local filesystem
- **Email**: SMTP (Gmail)
- **Scheduler**: robfig/cron
- **Deployment**: Docker + VPS via Jenkins

## 📋 Prerequisites

- Go 1.21 atau lebih tinggi
- Supabase Account
- SMTP credentials (Gmail App Password recommended)
- Make (optional)

## 🛠️ Setup Instructions

### 1. Clone & Install Dependencies

```bash
cd "Event Campus"
go mod download
```

### 2. Setup Supabase Database

1. Buat project baru di [Supabase](https://supabase.com)
2. Copy URL dan API Keys
3. Buka SQL Editor di Supabase Dashboard
4. Execute script di `migrations/001_initial_schema.sql`

### 3. Configure Environment Variables

```bash
cp .env.example .env
```

Edit `.env` dan isi dengan kredensial Anda:

```env
# Server
PORT=8080
ENV=development

# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key

# JWT
JWT_SECRET=your-strong-secret-key
JWT_EXPIRATION=24h

# Email (Gmail)
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

**Cara mendapatkan Gmail App Password:**
1. Go to https://myaccount.google.com/apppasswords
2. Generate new app password
3. Copy password ke `.env`

### 4. Run Application

```bash
# Using Make
make run

# Or directly
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

## 📚 API Documentation

### Authentication

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@uii.ac.id",
  "password": "password123",
  "full_name": "John Doe",
  "phone_number": "081234567890"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@uii.ac.id",
  "password": "password123"
}
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid",
      "email": "user@uii.ac.id",
      "full_name": "John Doe",
      "role": "mahasiswa",
      "is_uii_civitas": true
    }
  }
}
```

### Events

#### Get All Events
```http
GET /api/v1/events?category=seminar&status=published&page=1&limit=20
```

#### Get Event Detail
```http
GET /api/v1/events/:id
```

#### Create Event (Organisasi/Admin only)
```http
POST /api/v1/events
Authorization: Bearer <token>
Content-Type: multipart/form-data

{
  "title": "Workshop AI",
  "description": "Workshop tentang AI",
  "category": "workshop",
  "event_type": "online",
  "zoom_link": "https://zoom.us/j/123456",
  "start_date": "2024-01-15T10:00:00Z",
  "end_date": "2024-01-15T12:00:00Z",
  "registration_deadline": "2024-01-14T23:59:59Z",
  "max_participants": 100,
  "is_uii_only": true,
  "status": "published",
  "poster": <file>
}
```

### Registrations

#### Register to Event
```http
POST /api/v1/events/:id/register
Authorization: Bearer <token>
```

#### Cancel Registration
```http
DELETE /api/v1/registrations/:id
Authorization: Bearer <token>
```

#### My Registrations
```http
GET /api/v1/registrations/my
Authorization: Bearer <token>
```

### Whitelist (Organisasi Approval)

#### Submit Whitelist Request
```http
POST /api/v1/whitelist/request
Authorization: Bearer <token>
Content-Type: multipart/form-data

{
  "organization_name": "BEM FTI",
  "document": <pdf file>
}
```

#### Get All Requests (Admin)
```http
GET /api/v1/whitelist/requests
Authorization: Bearer <admin token>
```

#### Approve Request (Admin)
```http
PATCH /api/v1/whitelist/:id/approve
Authorization: Bearer <admin token>
Content-Type: application/json

{
  "approved": true,
  "admin_notes": "Approved"
}
```

### Attendance

#### Mark Attendance
```http
POST /api/v1/attendances
Authorization: Bearer <organisasi/admin token>
Content-Type: application/json

{
  "registration_id": "uuid",
  "notes": "Present"
}
```

## 🔒 User Roles & Permissions

### Mahasiswa (Default)
- ✅ Register & Login
- ✅ Browse & view events
- ✅ Register to events
- ✅ Cancel registration
- ✅ Submit whitelist request

### Organisasi (Approved Mahasiswa)
- ✅ All Mahasiswa permissions +
- ✅ Create, edit, delete own events
- ✅ View event registrations
- ✅ Mark attendance

### Admin (Super User)
- ✅ Full access to everything
- ✅ Approve/reject whitelist requests
- ✅ Manage all events

**Default Admin Credentials:**
- Email: `admin@uii.ac.id`
- Password: `admin123`

## 📧 Email Notifications

System akan mengirim email otomatis untuk:

1. **Registration Confirmation** - Saat berhasil daftar event
2. **Waitlist Notification** - Saat masuk waiting list
3. **Waitlist Promotion** - Saat dipromosikan dari waiting list
4. **Cancellation Confirmation** - Saat membatalkan pendaftaran
5. **H-1 Reminder** - Reminder H-1 sebelum event (includes zoom link)
6. **Whitelist Approval/Rejection** - Status pengajuan organisasi

## ⏰ Automated Schedulers

### H-1 Reminder (Daily at 00:00)
Mengirim reminder email H-1 untuk semua event yang akan berlangsung besok.

### Event Status Updater (Every hour)
Auto-update status event:
- `published` → `ongoing` (saat event mulai)
- `ongoing` → `completed` (saat event selesai)

## 🐳 Docker Deployment

### Build Image
```bash
make docker-build
```

### Run Container
```bash
make docker-run
```

### With Docker Compose
```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    env_file:
      - .env
    volumes:
      - ./storage:/app/storage
```

## 🔧 Development Commands

```bash
# Run application
make run

# Build binary
make build

# Run tests
make test

# Clean build artifacts
make clean

# Download dependencies
make deps
```

## 📁 Project Structure

```
event-campus-backend/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── config/                  # Configuration
│   ├── domain/                  # Business entities
│   ├── repository/              # Data access
│   ├── usecase/                 # Business logic
│   ├── delivery/http/           # HTTP handlers
│   ├── dto/                     # Data transfer objects
│   ├── utils/                   # Utilities
│   └── scheduler/               # Cron jobs
├── migrations/                  # Database migrations
├── storage/                     # File storage
│   ├── posters/
│   └── documents/
├── go.mod
├── Makefile
└── README.md
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/usecase/...
```

## 🚀 Deployment to VPS

### Manual Deployment

1. Build binary:
```bash
GOOS=linux GOARCH=amd64 go build -o bin/api cmd/api/main.go
```

2. Upload to VPS:
```bash
scp bin/api user@your-vps:/path/to/app/
scp .env user@your-vps:/path/to/app/
```

3. Run on VPS:
```bash
ssh user@your-vps
cd /path/to/app
./api
```

### With Jenkins (CI/CD)

Jenkins pipeline akan otomatis:
1. Pull code dari git
2. Run tests
3. Build Docker image
4. Deploy ke VPS
5. Health check

## 🔐 Security Notes

1. **Never commit `.env`** - Contains sensitive credentials
2. **Use strong JWT_SECRET** - Generate dengan `openssl rand -base64 32`
3. **Enable HTTPS** - Di production, selalu gunakan HTTPS
4. **Secure SMTP** - Gunakan App Password, bukan password asli
5. **Rate Limiting** - Consider adding rate limiting di production

## 📝 License

MIT License

## 👥 Contributors

- Your Name - Initial work

---

**Need Help?** Create an issue or contact the team.
