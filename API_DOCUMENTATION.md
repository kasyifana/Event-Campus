# Event Campus - REST API Documentation

Base URL: `http://localhost:3000/api/v1`

## 📑 Table of Contents

- [Authentication](#authentication)
- [Users](#users)
- [Events](#events)
- [Event Registration](#event-registration)
- [Whitelist (Organisasi Approval)](#whitelist-organisasi-approval)
- [Attendance](#attendance)
- [File Upload](#file-upload)
- [Error Responses](#error-responses)

---

## Authentication

### Register User

Creates a new user account.

**Endpoint:** `POST /auth/register`

**Access:** Public

**Request Body:**
```json
{
  "email": "mahasiswa@uii.ac.id",
  "password": "password123",
  "full_name": "Ahmad Rizki",
  "phone_number": "081234567890"
}
```

**Validation Rules:**
- `email`: Required, valid email format
- `password`: Required, minimum 8 characters
- `full_name`: Required
- `phone_number`: Required, Indonesian phone format

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "mahasiswa@uii.ac.id",
      "full_name": "Ahmad Rizki",
      "phone_number": "+6281234567890",
      "role": "mahasiswa",
      "is_uii_civitas": true,
      "is_approved": false,
      "created_at": "2024-01-15T10:00:00Z"
    }
  }
}
```

**Notes:**
- Email dengan domain `uii.ac.id` otomatis di-set `is_uii_civitas = true`
- Default role adalah `mahasiswa`
- Token JWT valid selama 24 jam

---

### Login

Authenticate user and get JWT token.

**Endpoint:** `POST /auth/login`

**Access:** Public

**Request Body:**
```json
{
  "email": "mahasiswa@uii.ac.id",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "mahasiswa@uii.ac.id",
      "full_name": "Ahmad Rizki",
      "phone_number": "+6281234567890",
      "role": "mahasiswa",
      "is_uii_civitas": true,
      "is_approved": false,
      "created_at": "2024-01-15T10:00:00Z"
    }
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "Login failed",
  "error": "invalid email or password"
}
```

---

## Users

### Get Profile

Get current user profile.

**Endpoint:** `GET /profile`

**Access:** Protected (Mahasiswa, Organisasi, Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile retrieved",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "mahasiswa@uii.ac.id",
    "role": "mahasiswa"
  }
}
```

---

### Update Profile

Update user profile information.

**Endpoint:** `PUT /profile`

**Access:** Protected (Mahasiswa, Organisasi, Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "full_name": "Ahmad Rizki Updated",
  "phone_number": "081234567891"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "mahasiswa@uii.ac.id",
    "full_name": "Ahmad Rizki Updated",
    "phone_number": "+6281234567891",
    "role": "mahasiswa"
  }
}
```

---

### Change Password

Change user password.

**Endpoint:** `POST /profile/change-password`

**Access:** Protected (Mahasiswa, Organisasi, Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "old_password": "password123",
  "new_password": "newpassword123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

---

## Events

### Get All Events

Get list of events with filtering and pagination.

**Endpoint:** `GET /events`

**Access:** Public

**Query Parameters:**
- `category` (optional): `seminar` | `workshop` | `lomba` | `konser`
- `status` (optional): `draft` | `published` | `ongoing` | `completed` | `cancelled`
- `event_type` (optional): `online` | `offline`
- `is_uii_only` (optional): `true` | `false`
- `search` (optional): Search in title and description
- `start_date` (optional): ISO 8601 date
- `end_date` (optional): ISO 8601 date
- `page` (optional): Default 1
- `limit` (optional): Default 20, max 100

**Example:**
```
GET /events?category=seminar&status=published&search=AI&page=1&limit=10
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Events retrieved successfully",
  "data": {
    "events": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "organizer_id": "550e8400-e29b-41d4-a716-446655440000",
        "organizer_name": "BEM FTI",
        "title": "Workshop AI untuk Pemula",
        "description": "Workshop pengenalan AI dan Machine Learning",
        "category": "workshop",
        "event_type": "online",
        "location": null,
        "zoom_link": "https://zoom.us/j/123456789",
        "poster_path": "posters/abc123.jpg",
        "poster_url": "http://localhost:8080/files/posters/abc123.jpg",
        "start_date": "2024-01-20T10:00:00Z",
        "end_date": "2024-01-20T12:00:00Z",
        "registration_deadline": "2024-01-19T23:59:59Z",
        "max_participants": 100,
        "current_participants": 45,
        "available_slots": 55,
        "is_uii_only": true,
        "status": "published",
        "is_full": false,
        "created_at": "2024-01-15T10:00:00Z",
        "updated_at": "2024-01-15T10:00:00Z"
      }
    ]
  },
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 45,
    "total_pages": 5
  }
}
```

---

### Get My Events

Get list of events created by current organizer.

**Endpoint:** `GET /events/my-events`

**Access:** Protected (Organisasi, Admin only)

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Events retrieved successfully",
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "Workshop AI untuk Pemula",
      "category": "workshop",
      "event_type": "online",
      "status": "published",
      "start_date": "2024-01-20T10:00:00Z",
      "end_date": "2024-01-20T12:00:00Z",
      "current_participants": 45,
      "max_participants": 100,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

**Notes:**
- Returns all events created by the authenticated organizer
- Includes events in all statuses (draft, published, etc.)

---

### Get Event Detail

Get detailed information about a specific event.

**Endpoint:** `GET /events/:id`

**Access:** Public

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Event retrieved successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "organizer_id": "550e8400-e29b-41d4-a716-446655440000",
    "organizer_name": "BEM FTI",
    "title": "Workshop AI untuk Pemula",
    "description": "Workshop lengkap tentang dasar-dasar AI...",
    "category": "workshop",
    "event_type": "online",
    "zoom_link": "https://zoom.us/j/123456789",
    "poster_url": "http://localhost:8080/files/posters/abc123.jpg",
    "start_date": "2024-01-20T10:00:00Z",
    "end_date": "2024-01-20T12:00:00Z",
    "registration_deadline": "2024-01-19T23:59:59Z",
    "max_participants": 100,
    "current_participants": 45,
    "available_slots": 55,
    "is_uii_only": true,
    "status": "published",
    "is_full": false
  }
}
```

**Notes:**
- Zoom link hanya ditampilkan jika:
  - User adalah organizer event, ATAU
  - User sudah registered dan H-1 event dimulai

---

### Create Event

Create a new event.

**Endpoint:** `POST /events`

**Access:** Protected (Organisasi, Admin only)

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request Body (multipart/form-data):**
```
title: Workshop AI untuk Pemula
description: Workshop lengkap tentang AI...
category: workshop
event_type: online
zoom_link: https://zoom.us/j/123456789
start_date: 2024-01-20T10:00:00Z
end_date: 2024-01-20T12:00:00Z
registration_deadline: 2024-01-19T23:59:59Z
max_participants: 100
is_uii_only: true
status: published
poster: <file.jpg>
```

**Validation Rules:**
- `title`: Required
- `description`: Required
- `category`: Required, one of: `seminar`, `workshop`, `lomba`, `konser`
- `event_type`: Required, `online` or `offline`
- `location`: Required if `event_type = offline`
- `zoom_link`: Required if `event_type = online`
- `start_date`: Required, must be future date
- `end_date`: Required, must be after `start_date`
- `registration_deadline`: Required, must be before `start_date`
- `max_participants`: Required, minimum 1
- `poster`: Optional, jpg/png, max 5MB

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Event created successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Workshop AI untuk Pemula",
    "status": "published",
    "poster_url": "http://localhost:8080/files/posters/abc123.jpg"
  }
}
```

---

### Upload Event Poster

Upload or update poster image for a specific event.

**Endpoint:** `POST /events/:id/poster`

**Access:** Protected (Event Owner or Admin)

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request Body (multipart/form-data):**
```
poster: <file.jpg>
```

**Validation:**
- Format: JPG, JPEG, or PNG only
- Max size: 5MB
- Event must exist and user must be the owner or admin

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Poster uploaded successfully",
  "data": {
    "poster_path": "posters/abc123.jpg"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Invalid file type. Only JPG and PNG are allowed"
}
```

**Notes:**
- Replaces existing poster if one already exists
- Old poster file is automatically deleted
- Poster URL will be updated in event details

---

### Update Event

Update an existing event.

**Endpoint:** `PUT /events/:id`

**Access:** Protected (Event Owner or Admin)

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request Body (multipart/form-data):**
All fields are optional. Only include fields to update.

```
title: Updated Title
description: Updated description
poster: <new_file.jpg>
```

**Notes:**
- Organisasi hanya bisa update event sendiri
- Admin bisa update semua event
- Tidak bisa reduce `max_participants` di bawah `current_participants`
- Tidak bisa change `is_uii_only` ke `true` jika sudah ada non-UII peserta

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Event updated successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Updated Title"
  }
}
```

---

### Publish Event

Publish an event (change status from draft to published).

**Endpoint:** `POST /events/:id/publish`

**Access:** Protected (Event Owner or Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Event published successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Failed to publish event",
  "error": "Event must have all required fields including poster"
}
```

**Notes:**
- Event must have all required fields filled
- Event must have a poster uploaded
- Published events become visible to all users
- Cannot unpublish an event once published

---

### Delete Event

Delete an event.

**Endpoint:** `DELETE /events/:id`

**Access:** Protected (Event Owner or Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Notes:**
- Hanya bisa delete jika `current_participants = 0`
- Jika ada peserta, event harus di-cancel dulu

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Event deleted successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Cannot delete event",
  "error": "Event has registered participants"
}
```

---

### Update Event Status

Update event status (admin only).

**Endpoint:** `PATCH /events/:id/status`

**Access:** Protected (Admin only)

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "status": "cancelled"
}
```

**Valid Status Values:**
- `draft`
- `published`
- `ongoing`
- `completed`
- `cancelled`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Event status updated",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "cancelled"
  }
}
```

**Notes:**
- Jika status diubah ke `cancelled`, semua peserta akan menerima email notifikasi

---

## Event Registration

### Register to Event

Register current user to an event.

**Endpoint:** `POST /events/:id/register`

**Access:** Protected (Mahasiswa, Organisasi, Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK) - Registered:**
```json
{
  "success": true,
  "message": "Successfully registered to event",
  "data": {
    "registration_id": "789e0123-e89b-12d3-a456-426614174000",
    "event_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "registered",
    "registered_at": "2024-01-15T10:00:00Z"
  }
}
```

**Response (200 OK) - Waitlist:**
```json
{
  "success": true,
  "message": "Event is full. You are added to waiting list",
  "data": {
    "registration_id": "789e0123-e89b-12d3-a456-426614174000",
    "event_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "waitlist",
    "position": 5,
    "registered_at": "2024-01-15T10:00:00Z"
  }
}
```

**Error Responses:**

*403 Forbidden (UII-Only Event):*
```json
{
  "success": false,
  "message": "Registration failed",
  "error": "This event is only for UII civitas"
}
```

*400 Bad Request (Already Registered):*
```json
{
  "success": false,
  "message": "Registration failed",
  "error": "You are already registered to this event"
}
```

*400 Bad Request (Registration Closed):*
```json
{
  "success": false,
  "message": "Registration failed",
  "error": "Registration is closed"
}
```

**Notes:**
- Otomatis ke waitlist jika `current_participants >= max_participants`
- Email konfirmasi dikirim setelah registrasi berhasil

---

### Cancel Registration

Cancel event registration.

**Endpoint:** `DELETE /registrations/:id`

**Access:** Protected (Registration Owner)

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Registration cancelled successfully",
  "data": {
    "waitlist_promoted": true,
    "promoted_user": "user@example.com"
  }
}
```

**Notes:**
- Jika ada user di waitlist, otomatis dipromosi ke `registered`
- User yang dipromosi menerima email notifikasi
- Tidak bisa cancel jika sudah attended

---

### Get My Registrations

Get all registrations of current user.

**Endpoint:** `GET /registrations/my`

**Access:** Protected (Mahasiswa, Organisasi, Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Registrations retrieved successfully",
  "data": {
    "upcoming": [
      {
        "id": "789e0123-e89b-12d3-a456-426614174000",
        "event_id": "123e4567-e89b-12d3-a456-426614174000",
        "event_title": "Workshop AI",
        "event_date": "2024-01-20T10:00:00Z",
        "status": "registered",
        "registered_at": "2024-01-15T10:00:00Z",
        "reminder_sent": false
      }
    ],
    "waitlist": [
      {
        "id": "789e0124-e89b-12d3-a456-426614174001",
        "event_id": "123e4568-e89b-12d3-a456-426614174001",
        "event_title": "Seminar Blockchain",
        "event_date": "2024-01-25T14:00:00Z",
        "status": "waitlist",
        "registered_at": "2024-01-16T10:00:00Z"
      }
    ],
    "past": [
      {
        "id": "789e0125-e89b-12d3-a456-426614174002",
        "event_id": "123e4569-e89b-12d3-a456-426614174002",
        "event_title": "Workshop Web Development",
        "event_date": "2024-01-10T10:00:00Z",
        "status": "attended",
        "registered_at": "2024-01-05T10:00:00Z"
      }
    ],
    "cancelled": []
  }
}
```

---

### Get Event Registrations

Get all registrations for a specific event.

**Endpoint:** `GET /events/:id/registrations`

**Access:** Protected (Event Owner or Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `status` (optional): `registered` | `waitlist` | `cancelled` | `attended`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Registrations retrieved successfully",
  "data": {
    "total_registered": 45,
    "total_waitlist": 10,
    "total_attended": 0,
    "registrations": [
      {
        "registration_id": "789e0123-e89b-12d3-a456-426614174000",
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Ahmad Rizki",
        "email": "mahasiswa@uii.ac.id",
        "phone_number": "+6281234567890",
        "status": "registered",
        "registered_at": "2024-01-15T10:00:00Z"
      }
    ]
  }
}
```

**Notes:**
- Organisasi hanya bisa lihat registrasi event sendiri
- Admin bisa lihat semua registrasi

---

## Whitelist (Organisasi Approval)

### Submit Whitelist Request

Submit request to become organisasi.

**Endpoint:** `POST /whitelist/request`

**Access:** Protected (Mahasiswa only)

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request Body (multipart/form-data):**
```
organization_name: BEM Fakultas Teknik Industri
document: <file.pdf>
```

**Validation:**
- `organization_name`: Required
- `document`: Required, PDF only, max 10MB
- User harus role `mahasiswa`
- Belum punya pending request

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Whitelist request submitted successfully",
  "data": {
    "id": "456e7890-e89b-12d3-a456-426614174000",
    "organization_name": "BEM Fakultas Teknik Industri",
    "status": "pending",
    "submitted_at": "2024-01-15T10:00:00Z"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Request failed",
  "error": "You already have a pending request"
}
```

---

### Get My Whitelist Request

Get current user's whitelist request status.

**Endpoint:** `GET /whitelist/my-request`

**Access:** Protected (Mahasiswa, Organisasi, Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK) - Has Request:**
```json
{
  "success": true,
  "message": "Request retrieved successfully",
  "data": {
    "id": "456e7890-e89b-12d3-a456-426614174000",
    "organization_name": "BEM Fakultas Teknik Industri",
    "document_path": "documents/abc123.pdf",
    "document_url": "http://localhost:8080/files/documents/abc123.pdf",
    "status": "pending",
    "submitted_at": "2024-01-15T10:00:00Z",
    "admin_notes": null,
    "reviewed_at": null
  }
}
```

**Response (200 OK) - No Request:**
```json
{
  "success": true,
  "message": "No request found",
  "data": null
}
```

**Notes:**
- Returns the most recent whitelist request for current user
- Shows status: `pending`, `approved`, or `rejected`
- Includes admin notes if request has been reviewed

---

### Get All Whitelist Requests

Get all whitelist requests (admin only).

**Endpoint:** `GET /whitelist/requests`

**Access:** Protected (Admin only)

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `status` (optional): `pending` | `approved` | `rejected`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Requests retrieved successfully",
  "data": [
    {
      "id": "456e7890-e89b-12d3-a456-426614174000",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "user_name": "Ahmad Rizki",
      "user_email": "mahasiswa@uii.ac.id",
      "organization_name": "BEM FTI",
      "document_path": "documents/abc123.pdf",
      "document_url": "http://localhost:8080/files/documents/abc123.pdf",
      "status": "pending",
      "submitted_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

---

### Approve Whitelist Request

Approve a whitelist request.

**Endpoint:** `PATCH /whitelist/:id/approve`

**Access:** Protected (Admin only)

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "approved": true,
  "admin_notes": "Approved. Berkas lengkap."
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Request approved successfully",
  "data": {
    "id": "456e7890-e89b-12d3-a456-426614174000",
    "status": "approved",
    "user_role": "organisasi",
    "reviewed_at": "2024-01-15T11:00:00Z"
  }
}
```

**Notes:**
- User role otomatis berubah dari `mahasiswa` → `organisasi`
- `is_approved` di-set `true`
- Email approval dikirim ke user

---

### Reject Whitelist Request

Reject a whitelist request.

**Endpoint:** `PATCH /whitelist/:id/reject`

**Access:** Protected (Admin only)

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "approved": false,
  "admin_notes": "Berkas tidak lengkap. Mohon melengkapi surat pengantar."
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Request rejected",
  "data": {
    "id": "456e7890-e89b-12d3-a456-426614174000",
    "status": "rejected",
    "reviewed_at": "2024-01-15T11:00:00Z"
  }
}
```

**Notes:**
- Email rejection dikirim ke user dengan alasan penolakan
- User bisa submit ulang request baru

---

## Attendance

### Mark Attendance

Mark user attendance at event.

**Endpoint:** `POST /events/:id/attendance`

**Access:** Protected (Event Owner or Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "notes": "Hadir tepat waktu"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Attendance marked successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Failed to mark attendance",
  "error": "User already marked as attended"
}
```

**Notes:**
- Registration status otomatis berubah ke `attended`
- Organisasi hanya bisa mark attendance untuk event sendiri
- User harus sudah registered ke event

---

### Bulk Mark Attendance

Mark attendance for multiple users at once.

**Endpoint:** `POST /events/:id/attendance/bulk`

**Access:** Protected (Event Owner or Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "user_ids": [
    "550e8400-e29b-41d4-a716-446655440000",
    "550e8400-e29b-41d4-a716-446655440001",
    "550e8400-e29b-41d4-a716-446655440002"
  ]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Bulk attendance marked successfully for 3 users"
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "No valid user IDs provided"
}
```

**Notes:**
- All users' registration status automatically change to `attended`
- Invalid user IDs are silently skipped
- Only users already registered to the event will be marked
- Organisasi can only mark attendance for their own events
- Useful for marking attendance via QR code scanner or batch import

---

### Get Event Attendance

Get all attendances for an event.

**Endpoint:** `GET /events/:id/attendance`

**Access:** Protected (Event Owner or Admin)

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Attendances retrieved successfully",
  "data": {
    "total_attended": 42,
    "total_registered": 45,
    "attendance_rate": 93.33,
    "attendances": [
      {
        "id": "321e6540-e89b-12d3-a456-426614174000",
        "registration_id": "789e0123-e89b-12d3-a456-426614174000",
        "user_name": "Ahmad Rizki",
        "user_email": "mahasiswa@uii.ac.id",
        "checked_in_at": "2024-01-20T10:05:00Z",
        "notes": "Hadir tepat waktu"
      }
    ]
  }
}
```

---

## File Upload

### Upload Poster

Upload event poster image.

**Endpoint:** `POST /upload/poster`

**Access:** Protected (Organisasi, Admin)

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request Body:**
```
poster: <file.jpg>
```

**Validation:**
- Format: JPG, JPEG, atau PNG
- Max size: 5MB

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Poster uploaded successfully",
  "data": {
    "path": "posters/abc123.jpg",
    "url": "http://localhost:8080/files/posters/abc123.jpg"
  }
}
```

---

### Upload Document

Upload whitelist document (PDF).

**Endpoint:** `POST /upload/document`

**Access:** Protected (Mahasiswa, Organisasi, Admin)

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request Body:**
```
document: <file.pdf>
```

**Validation:**
- Format: PDF only
- Max size: 10MB

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Document uploaded successfully",
  "data": {
    "path": "documents/def456.pdf",
    "url": "http://localhost:8080/files/documents/def456.pdf"
  }
}
```

---

### Serve File

Get uploaded file (poster or document).

**Endpoint:** `GET /files/:path`

**Access:** Public (for posters), Protected (for documents)

**Example:**
```
GET /files/posters/abc123.jpg
GET /files/documents/def456.pdf
```

**Response:**
Returns the file with appropriate content-type.

---

## Error Responses

### Standard Error Format

All error responses follow this format:

```json
{
  "success": false,
  "message": "Error message summary",
  "error": "Detailed error description"
}
```

### Common HTTP Status Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 200 | OK | Request berhasil |
| 201 | Created | Resource berhasil dibuat |
| 400 | Bad Request | Request invalid (validasi gagal) |
| 401 | Unauthorized | Token tidak ada atau invalid |
| 403 | Forbidden | Tidak punya permission |
| 404 | Not Found | Resource tidak ditemukan |
| 409 | Conflict | Data conflict (duplicate) |
| 500 | Internal Server Error | Server error |

### Example Error Responses

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Missing authorization header"
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "Forbidden",
  "error": "Insufficient permissions"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Not found",
  "error": "Event not found"
}
```

**400 Bad Request (Validation):**
```json
{
  "success": false,
  "message": "Invalid request",
  "error": "Field validation failed: email must be valid email format"
}
```

---

## Authentication & Authorization

### Using JWT Token

Semua protected endpoints membutuhkan JWT token di header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Role-Based Access Control

| Role | Permissions |
|------|-------------|
| **mahasiswa** | Browse events, register to events, submit whitelist request |
| **organisasi** | Semua permission mahasiswa + Create/edit/delete own events, view registrations, mark attendance |
| **admin** | Full access, approve/reject whitelist, edit/delete any event |

---

## Rate Limiting

**Current:** No rate limiting (development)

**Production Recommendation:**
- 100 requests per minute per IP
- 1000 requests per hour per user

---

## CORS

**Allowed Origins:**
Configure in `.env`:
```
ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
```

---

## Changelog

### Version 1.0.0 (Current)
- ✅ Authentication (register, login)
- ✅ User profile management
- ✅ Event CRUD (create, read, update, delete)
- ✅ Event registration and waitlist
- ✅ Whitelist system (organisasi approval)
- ✅ Attendance system (single and bulk)
- ✅ File upload (posters and documents)
- ✅ Email notifications
- ✅ Automated schedulers (H-1 reminders, status updates)

---

## Testing with cURL

### Example: Complete Registration Flow

```bash
# 1. Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@uii.ac.id",
    "password": "password123",
    "full_name": "Test User",
    "phone_number": "081234567890"
  }'

# Save the token from response

# 2. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@uii.ac.id",
    "password": "password123"
  }'

# 3. Get Profile
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## Support

**Documentation:** [README.md](README.md)  
**Implementation Plan:** [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)  
**Next Steps:** [NEXT_STEPS.md](NEXT_STEPS.md)

**Base URL (Development):** `http://localhost:8080/api/v1`  
**Base URL (Production):** `https://your-domain.com/api/v1`
