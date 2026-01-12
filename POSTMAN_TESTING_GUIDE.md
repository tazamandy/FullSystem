# Attendance System - Postman Testing Guide

## Quick Start

### 1. Import Collection

1. Open **Postman**
2. Click **Import** (top left)
3. Select **Upload Files**
4. Choose `Attendance_API_Postman_Collection.json`
5. Click **Import**

### 2. Configure Environment

The collection includes default variables:
- `base_url`: http://localhost:3000
- `access_token`: (auto-populated after login)
- `user_role`: (auto-populated after login)

To use custom values:
1. Click **Environments** (top right)
2. Select or create new environment
3. Update `base_url` if needed
4. Save

---

## Complete Testing Workflow

### Test 1: Authentication Flow

#### Step 1.1 - Register User
```
POST /register
```

**Test Data:**
```json
{
  "student_id": "2024-TEST-001",
  "email": "test.student@example.com",
  "password": "TestPass123!",
  "username": "test_student_001",
  "first_name": "Test",
  "last_name": "Student",
  "middle_name": "",
  "course": "Computer Science",
  "year_level": "3",
  "section": "A",
  "department": "Engineering",
  "college": "College of Engineering",
  "contact_number": "09123456789",
  "address": "123 Test Street"
}
```

**Expected Response (201):**
```json
{
  "message": "Registration successful. Please check your email to verify your account.",
  "student_id": "2024-TEST-001",
  "status": "success"
}
```

---

#### Step 1.2 - Verify Email
```
POST /verify
```

**Test Data:**
```json
{
  "email": "test.student@example.com",
  "code": "000000"
}
```

> **Note:** Use the verification code sent to email. For testing, check email or database.

**Expected Response (200):**
```json
{
  "message": "Email verified successfully",
  "status": "success"
}
```

---

#### Step 1.3 - Login
```
POST /login
```

**Test Data:**
```json
{
  "student_id": "2024-TEST-001",
  "password": "TestPass123!"
}
```

**Expected Response (200):**
```json
{
  "message": "Login successful",
  "student_id": "2024-TEST-001",
  "role": "student",
  "first_name": "Test",
  "last_name": "Student"
}
```

> **Important:** The `access_token` will be automatically saved to the `{{access_token}}` variable through the test script.

---

### Test 2: User Profile

#### Step 2.1 - Get Profile
```
GET /profile
```

**Headers:**
```
Authorization: Bearer {{access_token}}
```

**Expected Response (200):**
```json
{
  "id": 2,
  "student_id": "2024-TEST-001",
  "email": "test.student@example.com",
  "username": "test_student_001",
  "role": "student",
  "is_verified": true,
  "first_name": "Test",
  "last_name": "Student",
  "course": "Computer Science",
  ...
}
```

---

### Test 3: Events Management

#### Step 3.1 - Get All Events
```
GET /events?page=1&limit=10
```

**Headers:**
```
Authorization: Bearer {{access_token}}
```

**Expected Response (200):**
```json
{
  "events": [
    {
      "id": 1,
      "title": "Event Title",
      "description": "Event Description",
      "event_date": "2026-01-15",
      ...
    }
  ],
  "total": 1
}
```

---

#### Step 3.2 - Create Event (Faculty/Admin Only)
```
POST /events
```

**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Test Data (Faculty account required):**
```json
{
  "title": "Test Database Lecture",
  "description": "Introduction to Database Design",
  "event_date": "2026-01-20",
  "start_time": "2026-01-20T14:00:00Z",
  "end_time": "2026-01-20T16:00:00Z",
  "location": "Computer Lab 1",
  "course": "Computer Science",
  "section": "A",
  "year_level": "3",
  "department": "Engineering",
  "college": "College of Engineering"
}
```

**Expected Response (201):**
```json
{
  "id": 2,
  "message": "Event created successfully",
  "status": "success"
}
```

> **Note:** Student accounts cannot create events. Use a faculty/admin account.

---

#### Step 3.3 - Get Event by ID
```
GET /events/1
```

**Headers:**
```
Authorization: Bearer {{access_token}}
```

**Expected Response (200):**
```json
{
  "id": 1,
  "title": "Test Database Lecture",
  "description": "Introduction to Database Design",
  ...
}
```

---

#### Step 3.4 - Get My Events
```
GET /events/my-events
```

**Headers:**
```
Authorization: Bearer {{access_token}}
```

**Expected Response (200):**
```json
{
  "events": [
    {
      "id": 2,
      "title": "Test Database Lecture",
      ...
    }
  ],
  "total": 1
}
```

---

#### Step 3.5 - Update Event (Faculty/Admin Only)
```
PUT /events/1
```

**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Test Data:**
```json
{
  "title": "Updated Event Title",
  "description": "Updated Description",
  "status": "ongoing"
}
```

**Expected Response (200):**
```json
{
  "message": "Event updated successfully",
  "status": "success"
}
```

---

#### Step 3.6 - Delete Event (Faculty/Admin Only)
```
DELETE /events/1
```

**Headers:**
```
Authorization: Bearer {{access_token}}
```

**Expected Response (200):**
```json
{
  "message": "Event deleted successfully",
  "status": "success"
}
```

---

### Test 4: Attendance Tracking

#### Step 4.1 - Mark Attendance
```
POST /attendance/mark
```

**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Test Data:**
```json
{
  "event_id": 1,
  "qr_code_data": "EVENT_2026_001"
}
```

**Expected Response (201):**
```json
{
  "message": "Attendance marked successfully",
  "event_id": 1,
  "status": "present",
  "marked_at": "2026-01-20T14:30:00Z"
}
```

---

#### Step 4.2 - Get My Attendance
```
GET /attendance/my-attendance?page=1
```

**Headers:**
```
Authorization: Bearer {{access_token}}
```

**Optional Query Parameters:**
- `event_id=1`
- `status=present|absent|late`

**Expected Response (200):**
```json
{
  "attendance_records": [
    {
      "id": 1,
      "user_id": 2,
      "event_id": 1,
      "event_title": "Test Database Lecture",
      "status": "present",
      "marked_at": "2026-01-20T14:30:00Z",
      ...
    }
  ],
  "total": 1
}
```

---

#### Step 4.3 - Get Attendance Stats
```
GET /attendance/stats
```

**Headers:**
```
Authorization: Bearer {{access_token}}
```

**Expected Response (200):**
```json
{
  "total_events": 5,
  "present_count": 4,
  "absent_count": 1,
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

#### Step 4.4 - Get Attendance by Event (Faculty/Admin Only)
```
GET /events/1/attendance?page=1
```

**Headers:**
```
Authorization: Bearer {{access_token}}
```

**Expected Response (200):**
```json
{
  "event_id": 1,
  "event_title": "Test Database Lecture",
  "attendance_records": [
    {
      "id": 1,
      "student_id": "2024-TEST-001",
      "name": "Test Student",
      "status": "present",
      "marked_at": "2026-01-20T14:30:00Z"
    }
  ],
  "total_present": 25,
  "total_absent": 5,
  "total_late": 2,
  "total_students": 32
}
```

---

#### Step 4.5 - Update Attendance Status (Faculty/Admin Only)
```
PUT /attendance/1/status
```

**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Test Data:**
```json
{
  "status": "late"
}
```

**Valid Status Values:**
- `present`
- `absent`
- `late`
- `excused`

**Expected Response (200):**
```json
{
  "message": "Attendance status updated successfully",
  "id": 1,
  "status": "late"
}
```

---

### Test 5: Admin Operations

#### Step 5.1 - Get All Users (SuperAdmin Only)
```
GET /admin/users
```

**Headers:**
```
Authorization: Bearer {{access_token}}
```

> **Requirement:** Login with superadmin account

**Expected Response (200):**
```json
{
  "users": [
    {
      "id": 1,
      "student_id": "SUPERADMIN",
      "email": "superadmin@example.com",
      "username": "superadmin",
      "role": "superadmin",
      "is_verified": true,
      "first_name": "Super",
      "last_name": "Admin",
      "created_at": "2026-01-11T15:27:16.318Z"
    }
  ],
  "total": 1
}
```

---

#### Step 5.2 - Promote User (SuperAdmin Only)
```
POST /admin/promote
```

**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Test Data:**
```json
{
  "student_id": "2024-TEST-001",
  "new_role": "faculty"
}
```

**Expected Response (200):**
```json
{
  "message": "User promoted successfully",
  "student_id": "2024-TEST-001",
  "new_role": "faculty",
  "status": "success"
}
```

---

### Test 6: Password Reset

#### Step 6.1 - Forgot Password
```
POST /forgot-password
```

**Test Data:**
```json
{
  "email": "test.student@example.com"
}
```

**Expected Response (200):**
```json
{
  "message": "If your email is registered, you will receive a reset code.",
  "status": "success"
}
```

---

#### Step 6.2 - Reset Password
```
POST /reset-password
```

**Test Data:**
```json
{
  "email": "test.student@example.com",
  "code": "000000",
  "new_password": "NewTestPass123!"
}
```

**Expected Response (200):**
```json
{
  "message": "Password reset successful",
  "status": "success"
}
```

> **Note:** Use the reset code sent to email.

---

#### Step 6.3 - Resend Reset Code
```
POST /resend-reset-code
```

**Test Data:**
```json
{
  "email": "test.student@example.com"
}
```

**Expected Response (200):**
```json
{
  "message": "New code sent if email is registered",
  "status": "success"
}
```

---

## Testing Best Practices

### ✅ Do's
- Test with actual data similar to production
- Verify response status codes match expectations
- Check response data matches schema
- Test pagination with different page numbers
- Test with different user roles
- Verify error messages are clear

### ❌ Don'ts
- Use hardcoded IDs that might not exist
- Test with invalid tokens
- Ignore error responses
- Test sensitive operations in public networks
- Forget to add Authorization header
- Use same test data repeatedly (causes duplicates)

---

## Common Testing Scenarios

### Scenario 1: Complete Student Workflow
1. Register as student
2. Verify email
3. Login
4. Get profile
5. View events
6. Mark attendance
7. View attendance stats

### Scenario 2: Faculty Operations
1. Login as faculty (or promote student to faculty)
2. Create event
3. View my events
4. View event attendance
5. Update attendance status
6. Delete event

### Scenario 3: Admin Operations
1. Login as superadmin
2. View all users
3. Promote user to faculty
4. View user list again

### Scenario 4: Password Management
1. Request password reset
2. Verify reset code received
3. Reset password
4. Login with new password

---

## Troubleshooting

### Issue: 401 Unauthorized
**Cause:** Missing or invalid token
**Solution:**
1. Login again to get fresh token
2. Verify token is in Authorization header
3. Check token format: `Bearer {token}`

### Issue: 403 Forbidden
**Cause:** Insufficient permissions
**Solution:**
1. Check user role (student, faculty, admin, superadmin)
2. Some endpoints require specific roles
3. Use appropriate account for operation

### Issue: 400 Bad Request
**Cause:** Invalid request data
**Solution:**
1. Verify JSON syntax is correct
2. Check all required fields are present
3. Validate data types (string, number, date)
4. Review field names match exactly

### Issue: 404 Not Found
**Cause:** Resource doesn't exist
**Solution:**
1. Verify ID exists in database
2. Check endpoint path is correct
3. Use valid event/attendance IDs

### Issue: 429 Too Many Requests
**Cause:** Rate limit exceeded
**Solution:**
1. Wait before making more requests
2. Check rate limit headers
3. Space out requests

---

## Tips & Tricks

### Save Token Automatically
In **Login** request → **Tests** tab, the token is auto-saved:
```javascript
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.environment.set("access_token", jsonData.access_token);
    pm.environment.set("user_role", jsonData.role);
}
```

### Create Reusable Test Data
Use environment variables for common data:
```
student_id: {{student_id}}
email: {{email}}
```

### Run Collection Tests
1. Select collection
2. Click **Run** (play icon)
3. Select requests to run
4. Click **Run Attendance API**

### Export Test Results
1. Run collection tests
2. Click **Save Results** button
3. Choose format (JSON, CSV, HTML)

---

## Contact & Support

For API documentation, see: `API_DOCUMENTATION.md`

For issues or questions, check logs: `logs/` directory

---

**Last Updated:** January 11, 2026
