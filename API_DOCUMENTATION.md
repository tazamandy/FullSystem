# Attendance System API Documentation

## Base URL
```
http://localhost:3000
```

## Authentication
Most endpoints require Bearer Token authentication. Include the token in the `Authorization` header:
```
Authorization: Bearer {access_token}
```

---

## Table of Contents
1. [Authentication Endpoints](#authentication-endpoints)
2. [User Endpoints](#user-endpoints)
3. [Event Endpoints](#event-endpoints)
4. [Attendance Endpoints](#attendance-endpoints)

---

## Authentication Endpoints

### 1. Register User
**POST** `/register`

Create a new user account.

**Authentication:** No

**Request Body:**
```json
{
  "student_id": "2024-001",
  "email": "student@example.com",
  "password": "SecurePass123!",
  "username": "student_username",
  "first_name": "John",
  "last_name": "Doe",
  "middle_name": "M",
  "course": "Computer Science",
  "year_level": "3",
  "section": "A",
  "department": "Engineering",
  "college": "College of Engineering",
  "contact_number": "09123456789",
  "address": "123 Main Street"
}
```

**Response (201):**
```json
{
  "message": "Registration successful. Please check your email to verify your account.",
  "student_id": "2024-001",
  "status": "success"
}
```

**Postman Steps:**
1. Create new request → POST
2. URL: `http://localhost:3000/register`
3. Body → raw → JSON
4. Paste request JSON above
5. Send

---

### 2. Verify Email
**POST** `/verify`

Verify user email with verification code.

**Authentication:** No

**Request Body:**
```json
{
  "email": "student@example.com",
  "code": "123456"
}
```

**Response (200):**
```json
{
  "message": "Email verified successfully",
  "status": "success"
}
```

---

### 3. Login
**POST** `/login`

Login with student ID or email and password.

**Authentication:** No

**Request Body:**
```json
{
  "student_id": "2024-001",
  "password": "SecurePass123!"
}
```

**Or with email:**
```json
{
  "student_id": "student@example.com",
  "password": "SecurePass123!"
}
```

**Response (200):**
```json
{
  "message": "Login successful",
  "student_id": "2024-001",
  "role": "student",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Postman Steps:**
1. Create new request → POST
2. URL: `http://localhost:3000/login`
3. Body → raw → JSON
4. Send

---

### 4. Refresh Token
**POST** `/refresh-token`

Refresh expired access token.

**Authentication:** Yes (Bearer Token)

**Request Body:**
```json
{
  "refresh_token": "your_refresh_token_here"
}
```

**Response (200):**
```json
{
  "access_token": "new_access_token",
  "refresh_token": "new_refresh_token",
  "expires_in": 3600
}
```

---

### 5. Forgot Password
**POST** `/forgot-password`

Request password reset code.

**Authentication:** No

**Request Body:**
```json
{
  "email": "student@example.com"
}
```

**Response (200):**
```json
{
  "message": "If your email is registered, you will receive a reset code.",
  "status": "success"
}
```

---

### 6. Reset Password
**POST** `/reset-password`

Reset password with verification code.

**Authentication:** No

**Request Body:**
```json
{
  "email": "student@example.com",
  "code": "123456",
  "new_password": "NewSecurePass123!"
}
```

**Response (200):**
```json
{
  "message": "Password reset successful",
  "status": "success"
}
```

---

### 7. Resend Reset Code
**POST** `/resend-reset-code`

Resend password reset code.

**Authentication:** No

**Request Body:**
```json
{
  "email": "student@example.com"
}
```

**Response (200):**
```json
{
  "message": "New code sent if email is registered",
  "status": "success"
}
```

---

## User Endpoints

### 1. Get Profile
**GET** `/profile`

Get current logged-in user profile.

**Authentication:** Yes (Bearer Token)

**Response (200):**
```json
{
  "id": 1,
  "student_id": "2024-001",
  "email": "student@example.com",
  "username": "student_username",
  "role": "student",
  "is_verified": true,
  "first_name": "John",
  "last_name": "Doe",
  "middle_name": "M",
  "course": "Computer Science",
  "year_level": "3",
  "section": "A",
  "department": "Engineering",
  "college": "College of Engineering",
  "contact_number": "09123456789",
  "address": "123 Main Street",
  "created_at": "2026-01-11T15:27:16.318Z",
  "verified_at": "2026-01-11T15:27:16.318Z"
}
```

**Postman Steps:**
1. Create new request → GET
2. URL: `http://localhost:3000/profile`
3. Headers → Add `Authorization: Bearer {your_token}`
4. Send

---

### 2. Get All Users (Admin Only)
**GET** `/admin/users`

Get all users in the system.

**Authentication:** Yes (Bearer Token - SuperAdmin role required)

**Response (200):**
```json
{
  "users": [
    {
      "id": 1,
      "student_id": "2024-001",
      "email": "student@example.com",
      "username": "student_username",
      "role": "student",
      "is_verified": true,
      "first_name": "John",
      "last_name": "Doe",
      "created_at": "2026-01-11T15:27:16.318Z"
    }
  ],
  "total": 1
}
```

---

### 3. Promote User (Admin Only)
**POST** `/admin/promote`

Promote student to faculty/admin.

**Authentication:** Yes (Bearer Token - SuperAdmin role required)

**Request Body:**
```json
{
  "student_id": "2024-001",
  "new_role": "faculty"
}
```

**Valid Roles:** `student`, `faculty`, `admin`, `superadmin`

**Response (200):**
```json
{
  "message": "User promoted successfully",
  "student_id": "2024-001",
  "new_role": "faculty",
  "status": "success"
}
```

---

## Event Endpoints

### 1. Get All Events
**GET** `/events`

Get all events.

**Authentication:** Yes (Bearer Token)

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `status` (optional): Filter by status (scheduled, ongoing, completed)

**Response (200):**
```json
{
  "events": [
    {
      "id": 1,
      "title": "Operating Systems Lecture",
      "description": "Guest lecture on OS Concepts",
      "event_date": "2026-01-15",
      "start_time": "2026-01-15T10:00:00Z",
      "end_time": "2026-01-15T12:00:00Z",
      "location": "Room 101",
      "course": "Computer Science",
      "section": "A",
      "year_level": "3",
      "department": "Engineering",
      "college": "College of Engineering",
      "created_by": "faculty@example.com",
      "created_by_role": "faculty",
      "status": "scheduled",
      "is_active": true,
      "qr_code_data": "EVENT_2026_001",
      "created_at": "2026-01-11T15:27:16.318Z"
    }
  ],
  "total": 1
}
```

---

### 2. Get My Events
**GET** `/events/my-events`

Get events created by current user.

**Authentication:** Yes (Bearer Token)

**Response (200):**
```json
{
  "events": [
    {
      "id": 1,
      "title": "Operating Systems Lecture",
      ...
    }
  ],
  "total": 1
}
```

---

### 3. Get Event by ID
**GET** `/events/:id`

Get specific event details.

**Authentication:** Yes (Bearer Token)

**URL Parameters:**
- `id`: Event ID

**Response (200):**
```json
{
  "id": 1,
  "title": "Operating Systems Lecture",
  "description": "Guest lecture on OS Concepts",
  ...
}
```

---

### 4. Create Event (Faculty/Admin Only)
**POST** `/events`

Create new event.

**Authentication:** Yes (Bearer Token - Faculty or Admin role required)

**Request Body:**
```json
{
  "title": "Operating Systems Lecture",
  "description": "Guest lecture on OS Concepts",
  "event_date": "2026-01-15",
  "start_time": "2026-01-15T10:00:00Z",
  "end_time": "2026-01-15T12:00:00Z",
  "location": "Room 101",
  "course": "Computer Science",
  "section": "A",
  "year_level": "3",
  "department": "Engineering",
  "college": "College of Engineering"
}
```

**Response (201):**
```json
{
  "id": 1,
  "message": "Event created successfully",
  "status": "success"
}
```

**Postman Steps:**
1. Create new request → POST
2. URL: `http://localhost:3000/events`
3. Headers → Add `Authorization: Bearer {your_token}`
4. Body → raw → JSON → Paste request JSON
5. Send

---

### 5. Update Event (Faculty/Admin Only)
**PUT** `/events/:id`

Update event details.

**Authentication:** Yes (Bearer Token - Faculty or Admin role required)

**URL Parameters:**
- `id`: Event ID

**Request Body:**
```json
{
  "title": "Updated Event Title",
  "description": "Updated description",
  "status": "ongoing"
}
```

**Response (200):**
```json
{
  "message": "Event updated successfully",
  "status": "success"
}
```

---

### 6. Delete Event (Faculty/Admin Only)
**DELETE** `/events/:id`

Delete event.

**Authentication:** Yes (Bearer Token - Faculty or Admin role required)

**URL Parameters:**
- `id`: Event ID

**Response (200):**
```json
{
  "message": "Event deleted successfully",
  "status": "success"
}
```

---

## Attendance Endpoints

### 1. Mark Attendance
**POST** `/attendance/mark`

Mark attendance for an event.

**Authentication:** Yes (Bearer Token)

**Request Body:**
```json
{
  "event_id": 1,
  "qr_code_data": "EVENT_2026_001"
}
```

**Response (201):**
```json
{
  "message": "Attendance marked successfully",
  "event_id": 1,
  "status": "present",
  "marked_at": "2026-01-15T10:30:00Z"
}
```

---

### 2. Get My Attendance
**GET** `/attendance/my-attendance`

Get attendance records for current user.

**Authentication:** Yes (Bearer Token)

**Query Parameters:**
- `event_id` (optional): Filter by event ID
- `status` (optional): Filter by status (present, absent, late)
- `page` (optional): Page number

**Response (200):**
```json
{
  "attendance_records": [
    {
      "id": 1,
      "user_id": 1,
      "event_id": 1,
      "event_title": "Operating Systems Lecture",
      "status": "present",
      "marked_at": "2026-01-15T10:30:00Z",
      "event_date": "2026-01-15",
      "created_at": "2026-01-15T10:30:00Z"
    }
  ],
  "total": 1
}
```

---

### 3. Get Attendance Stats
**GET** `/attendance/stats`

Get attendance statistics for current user.

**Authentication:** Yes (Bearer Token)

**Response (200):**
```json
{
  "total_events": 10,
  "present_count": 8,
  "absent_count": 2,
  "late_count": 0,
  "attendance_rate": 80.0,
  "courses": {
    "Computer Science": {
      "total": 5,
      "present": 4,
      "rate": 80.0
    }
  }
}
```

---

### 4. Get Attendance by Event (Faculty/Admin Only)
**GET** `/events/:event_id/attendance`

Get all attendance records for an event.

**Authentication:** Yes (Bearer Token - Faculty or Admin role required)

**URL Parameters:**
- `event_id`: Event ID

**Query Parameters:**
- `status` (optional): Filter by status (present, absent, late)
- `page` (optional): Page number

**Response (200):**
```json
{
  "event_id": 1,
  "event_title": "Operating Systems Lecture",
  "attendance_records": [
    {
      "id": 1,
      "student_id": "2024-001",
      "name": "John Doe",
      "status": "present",
      "marked_at": "2026-01-15T10:30:00Z"
    }
  ],
  "total_present": 25,
  "total_absent": 5,
  "total_late": 2,
  "total_students": 32
}
```

---

### 5. Update Attendance Status (Faculty/Admin Only)
**PUT** `/attendance/:id/status`

Update attendance status for a record.

**Authentication:** Yes (Bearer Token - Faculty or Admin role required)

**URL Parameters:**
- `id`: Attendance record ID

**Request Body:**
```json
{
  "status": "late"
}
```

**Valid Status Values:** `present`, `absent`, `late`, `excused`

**Response (200):**
```json
{
  "message": "Attendance status updated successfully",
  "id": 1,
  "status": "late"
}
```

---

## System Endpoints

### 1. Health Check
**GET** `/health`

Check system status.

**Authentication:** No

**Response (200):**
```json
{
  "status": "ok",
  "service": "attendance-backend",
  "port": "3000"
}
```

---

### 2. Root Info
**GET** `/`

Get request information.

**Authentication:** No

**Response (200):**
```json
{
  "remote_ip": "127.0.0.1",
  "user_agent": "PostmanRuntime/7.32.1",
  "host": "localhost:3000"
}
```

---

## Error Responses

All endpoints follow standard HTTP status codes:

| Code | Meaning |
|------|---------|
| 200 | OK - Request successful |
| 201 | Created - Resource created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing/invalid token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |

**Error Response Format:**
```json
{
  "error": "Error message description"
}
```

---

## Testing in Postman

### Setup

1. **Download Postman** from https://www.postman.com/downloads/

2. **Create Collection**
   - Open Postman
   - Click "Create" → "Collection" → Name: "Attendance API"

3. **Set Base URL (Environment Variable)**
   - Click "Environments" → "Create New Environment"
   - Add variable: `base_url` = `http://localhost:3000`
   - Add variable: `token` = (will be populated after login)

### Test Flow

#### 1. Register User
```
POST {{base_url}}/register
```

#### 2. Verify Email
```
POST {{base_url}}/verify
```

#### 3. Login
```
POST {{base_url}}/login
```
- **Important:** Copy the response and store the token in the `token` environment variable

#### 4. Get Profile
```
GET {{base_url}}/profile
```
- Add Header: `Authorization: Bearer {{token}}`

#### 5. Create Event (Faculty only)
```
POST {{base_url}}/events
```

#### 6. Mark Attendance
```
POST {{base_url}}/attendance/mark
```

#### 7. Get Attendance Stats
```
GET {{base_url}}/attendance/stats
```

### Pre-request Script (for storing token)

In Login request, go to **Tests** tab and add:
```javascript
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.environment.set("token", jsonData.access_token);
}
```

---

## Postman Collection JSON (Import)

You can import this JSON into Postman:

```json
{
  "info": {
    "name": "Attendance API",
    "description": "Attendance Management System API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Auth",
      "item": [
        {
          "name": "Register",
          "request": {
            "method": "POST",
            "url": "{{base_url}}/register"
          }
        },
        {
          "name": "Login",
          "request": {
            "method": "POST",
            "url": "{{base_url}}/login"
          }
        }
      ]
    }
  ]
}
```

---

## Common Issues & Solutions

### Issue: 401 Unauthorized
**Solution:** 
- Ensure token is included in Authorization header
- Token might be expired, login again
- Check if token format is correct: `Bearer {token}`

### Issue: 403 Forbidden
**Solution:**
- User doesn't have required role
- Faculty endpoints require faculty/admin role
- Admin endpoints require superadmin role

### Issue: 429 Too Many Requests
**Solution:**
- Rate limit exceeded
- Wait before making more requests
- Check rate limit headers in response

### Issue: Database Column Not Found
**Solution:**
- Run migration: `psql -U postgres -d attendance < attendance.sql`
- Ensure database schema is up to date

---

## Notes

- All timestamps are in UTC format (ISO 8601)
- Passwords must meet security requirements
- Email verification is required for account activation
- Rate limiting applies to all endpoints
- CORS is enabled for development

