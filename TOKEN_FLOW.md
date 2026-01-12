# Automatic Token Flow Documentation

## Overview
The backend now automatically generates and passes tokens through the registration and password reset flows, enabling seamless frontend navigation without manual token handling.

---

## 1. REGISTRATION FLOW (reg/verify endpoint)

### Step 1: Register User
**Endpoint:** `POST /register`

**Request:**
```json
{
  "email": "student@example.com",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe",
  "student_id": "2024001"
}
```

**Response:**
```json
{
  "message": "Registration successful. Please check your email to verify your account.",
  "student_id": "2024001",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "status": "success"
}
```

**What happens:**
- User registration data is saved to `pending_users` table
- Verification code is generated and sent to email
- JWT token is automatically generated for email verification
- Token is returned to frontend

### Step 2: Verify Email with Automatic Token
**Endpoint:** `POST /reg/verify`

**Request Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Request Body:**
```json
{
  "code": "123456"
}
```

**Response:**
```json
{
  "message": "Email verified successfully",
  "user": {
    "student_id": "2024001",
    "email": "student@example.com",
    "is_verified": true
  },
  "status": "success"
}
```

**What happens:**
- Token from Step 1 is verified
- Code from email is verified against token's code
- Pending user is moved to active users table
- User is fully registered and can login

---

## 2. FORGOT PASSWORD FLOW (fgtp endpoints)

### Step 1: Request Password Reset
**Endpoint:** `POST /fgtp/forgot-password`

**Request:**
```json
{
  "email": "student@example.com"
}
```

**Response:**
```json
{
  "message": "If your email is registered, you will receive a reset code.",
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**What happens:**
- Reset code is generated and sent to email
- JWT token is automatically generated for password reset flow
- Token is returned to frontend (valid for 15 minutes)

### Step 2: Verify Reset Code with Automatic Token
**Endpoint:** `POST /fgtp/verify-reset-code`

**Request Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Request Body:**
```json
{
  "code": "123456"
}
```

**Response:**
```json
{
  "message": "Code is valid",
  "status": "success",
  "email": "student@example.com",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**What happens:**
- Token from Step 1 is verified
- Code is validated from password_resets table
- Same token is returned (can be used directly for next step)

### Step 3: Reset Password with Token
**Endpoint:** `POST /fgtp/reset-password`

**Request Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Request Body:**
```json
{
  "new_password": "NewSecurePass123!",
  "confirm_new_password": "NewSecurePass123!"
}
```

**Response:**
```json
{
  "message": "Password reset successful",
  "status": "success"
}
```

**What happens:**
- Token is verified
- New password is hashed and updated
- Reset code is marked as used
- Confirmation email is sent to user

### Step 4 (Optional): Resend Reset Code
**Endpoint:** `POST /fgtp/resend-code`

**Request:**
```json
{
  "email": "student@example.com"
}
```

**Response:**
```json
{
  "message": "New code sent if email is registered",
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**What happens:**
- Any existing unused codes are deleted
- New code is generated and sent
- New token is generated and returned

---

## Token Details

### Email Verification Token
- **Expiry:** 24 hours
- **Contains:** email, verification code
- **Used for:** `/reg/verify` endpoint
- **Purpose:** Validates email and verification code match

### Password Reset Token
- **Expiry:** 15 minutes
- **Contains:** email, reset code
- **Used for:** `/fgtp/verify-reset-code` and `/fgtp/reset-password` endpoints
- **Purpose:** Validates password reset request and code match

---

## Frontend Integration Guide

### Registration Flow
```javascript
// 1. Register
const regResponse = await fetch('http://localhost:3000/register', {
  method: 'POST',
  body: JSON.stringify(userData)
});
const { token } = await regResponse.json();

// 2. Verify with automatic token
const verifyResponse = await fetch('http://localhost:3000/reg/verify', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({ code: userCode })
});
```

### Password Reset Flow
```javascript
// 1. Forgot password
const forgotResponse = await fetch('http://localhost:3000/fgtp/forgot-password', {
  method: 'POST',
  body: JSON.stringify({ email })
});
const { token: resetToken } = await forgotResponse.json();

// 2. Verify reset code
const verifyResetResponse = await fetch('http://localhost:3000/fgtp/verify-reset-code', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${resetToken}`
  },
  body: JSON.stringify({ code: userCode })
});
const { token: finalToken } = await verifyResetResponse.json();

// 3. Reset password
const resetResponse = await fetch('http://localhost:3000/fgtp/reset-password', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${finalToken}`
  },
  body: JSON.stringify({
    new_password: newPassword,
    confirm_new_password: newPassword
  })
});
```

---

## Security Features

1. **Token Expiration:** Tokens expire automatically (24hrs for email, 15mins for password reset)
2. **Code Verification:** Codes are stored in database with expiration times
3. **One-Time Use:** Codes can only be used once, marked as used after verification
4. **Email Confirmation:** All sensitive operations require email confirmation
5. **Input Sanitization:** All inputs are sanitized before processing
6. **No Error Exposure:** Security messages don't reveal if email exists

