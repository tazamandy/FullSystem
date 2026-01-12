# Postman Testing Guide - Automatic Token Flow

## Testing the Complete Registration & Password Reset Flows

### 1. REGISTRATION FLOW TEST

#### Request 1: Register New User
```
POST http://localhost:3000/register

Content-Type: application/json

{
  "email": "testuser@example.com",
  "password": "TestPassword123!",
  "first_name": "Test",
  "last_name": "User",
  "course": "CS101",
  "year_level": "1st Year"
}
```

**Expected Response:**
```json
{
  "message": "Registration successful. Please check your email to verify your account.",
  "student_id": "2024XXXX",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "status": "success"
}
```

**⚠️ IMPORTANT:** 
- Copy the `token` value from response
- Check email for verification code (in test, check console/logs)

---

#### Request 2: Verify Email with Token from Request 1
```
POST http://localhost:3000/reg/verify

Authorization: Bearer {TOKEN_FROM_REQUEST_1}
Content-Type: application/json

{
  "code": "123456"  // Code from email
}
```

**Expected Response:**
```json
{
  "message": "Email verified successfully",
  "user": {
    "student_id": "2024XXXX",
    "email": "testuser@example.com",
    "is_verified": true
  },
  "status": "success"
}
```

**✅ At this point:**
- User account is now active
- User can login with email and password
- User has been assigned a QR code

---

### 2. FORGOT PASSWORD FLOW TEST

#### Request 3: Request Password Reset
```
POST http://localhost:3000/fgtp/forgot-password

Content-Type: application/json

{
  "email": "testuser@example.com"
}
```

**Expected Response:**
```json
{
  "message": "If your email is registered, you will receive a reset code.",
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**⚠️ IMPORTANT:**
- Copy the `token` value from response
- Check email for reset code (in test, check console/logs)
- Token is valid for 15 minutes only

---

#### Request 4: Verify Reset Code with Token from Request 3
```
POST http://localhost:3000/fgtp/verify-reset-code

Authorization: Bearer {TOKEN_FROM_REQUEST_3}
Content-Type: application/json

{
  "code": "123456"  // Code from email
}
```

**Expected Response:**
```json
{
  "message": "Code is valid",
  "status": "success",
  "email": "testuser@example.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**✅ The response token is the SAME token (can be reused)**

---

#### Request 5: Reset Password with Token from Request 4
```
POST http://localhost:3000/fgtp/reset-password

Authorization: Bearer {TOKEN_FROM_REQUEST_4}
Content-Type: application/json

{
  "new_password": "NewPassword456!",
  "confirm_new_password": "NewPassword456!"
}
```

**Expected Response:**
```json
{
  "message": "Password reset successful",
  "status": "success"
}
```

**✅ At this point:**
- Password has been changed
- User can login with new password
- Confirmation email has been sent
- Old reset code is marked as used (cannot be reused)

---

### 3. OPTIONAL: Resend Reset Code

#### Request 6: Resend Reset Code
```
POST http://localhost:3000/fgtp/resend-code

Content-Type: application/json

{
  "email": "testuser@example.com"
}
```

**Expected Response:**
```json
{
  "message": "New code sent if email is registered",
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**What happens:**
- Previous reset codes are deleted
- New reset code is generated and sent
- New token is provided (for fresh attempt)

---

## Postman Environment Variables Setup

Create these environment variables for easier testing:

```
{
  "base_url": "http://localhost:3000",
  "register_token": "",
  "reset_token": "",
  "test_email": "testuser@example.com",
  "test_password": "TestPassword123!",
  "new_password": "NewPassword456!"
}
```

Then update your requests to use variables:

```
POST {{base_url}}/register
Authorization: Bearer {{register_token}}
```

---

## Postman Pre-request & Tests Scripts

### For Register Request - Extract Token
```javascript
// Tests tab
var jsonData = pm.response.json();
pm.environment.set("register_token", jsonData.token);
console.log("Register Token: " + jsonData.token);
```

### For Verify Email Request
```javascript
// Pre-request Script
var token = pm.environment.get("register_token");
pm.request.headers.add({
  key: "Authorization",
  value: "Bearer " + token
});
```

### For Forgot Password - Extract Token
```javascript
// Tests tab
var jsonData = pm.response.json();
pm.environment.set("reset_token", jsonData.token);
console.log("Reset Token: " + jsonData.token);
```

### For Verify Reset Code - Extract Token
```javascript
// Tests tab
var jsonData = pm.response.json();
pm.environment.set("reset_token", jsonData.token);
console.log("Reset Token (Updated): " + jsonData.token);
```

---

## Troubleshooting

### Token Expired Error
- **Issue:** "Invalid or expired token"
- **Solution:** Get a fresh token by making the initial request again

### Code Not Matching
- **Issue:** "Invalid verification code"
- **Cause:** Code from email doesn't match code in request
- **Solution:** Make sure you're using the exact code sent to email

### User Not Found
- **Issue:** "user not found" or "Registration not found"
- **Cause:** Email doesn't exist or pending registration was deleted
- **Solution:** Register the user again with `/register`

### Password Requirements
- **Minimum 8 characters**
- **At least 1 uppercase letter** (A-Z)
- **At least 1 lowercase letter** (a-z)
- **At least 1 number** (0-9)
- **At least 1 special character** (!@#$%^&*() etc)

Example valid password: `SecurePass123!`

---

## Complete Test Scenario

1. ✅ Call `/register` → Get registration token
2. ✅ Check email for verification code
3. ✅ Call `/reg/verify` with token + code → User verified
4. ✅ Call `/fgtp/forgot-password` → Get reset token
5. ✅ Check email for reset code
6. ✅ Call `/fgtp/verify-reset-code` with reset token + code → Get same token back
7. ✅ Call `/fgtp/reset-password` with reset token + new password → Password changed
8. ✅ User can now login with new password

