// models/errors_model.go
// Centralized error messages and constants

package models

// HTTP Error Messages
const (
	// General Errors
	ErrInvalidRequest  = "Invalid request"
	ErrUnauthorized    = "Unauthorized"
	ErrForbidden       = "Forbidden"
	ErrNotFound        = "Not found"
	ErrInternalServer  = "Internal server error"
	ErrTooManyRequests = "Too many requests"

	// Authentication Errors
	ErrInvalidCredentials         = "Invalid email or password"
	ErrInvalidStudentID           = "Invalid student ID or password"
	ErrEmailNotVerified           = "Email not verified. Please check your email."
	ErrRefreshTokenExpired        = "Invalid or expired refresh token"
	ErrTokenRequired              = "refresh_token is required"
	ErrFailedGenerateAccessToken  = "Failed to generate access token"
	ErrFailedGenerateRefreshToken = "Failed to generate refresh token"

	// Validation Errors
	ErrFieldsRequired        = "Student ID/Email and password are required"
	ErrEmailPasswordRequired = "Email and password are required"
	ErrEmailRequired         = "Email is required"
	ErrPasswordRequired      = "Password is required"
	ErrCodeRequired          = "Code is required"
	ErrAllFieldsRequired     = "All fields are required"

	// User Errors
	ErrUserNotFound           = "User not found"
	ErrFailedFetchUserProfile = "Failed to fetch user profile"
	ErrEmailAlreadyExists     = "Email already exists"
	ErrStudentIDAlreadyExists = "Student ID already exists"
	ErrUsernameTaken          = "Username is already taken"

	// Verification Errors
	ErrInvalidVerificationCode  = "Invalid verification code"
	ErrVerificationCodeExpired  = "Verification code has expired. Please register again."
	ErrInvalidResetCode         = "Invalid or expired reset code"
	ErrPasswordResetCodeExpired = "Password reset code has expired"

	// Event Errors
	ErrEventNotFound      = "Event not found"
	ErrFailedCreateEvent  = "Failed to create event"
	ErrFailedUpdateEvent  = "Failed to update event"
	ErrFailedDeleteEvent  = "Failed to delete event"
	ErrInvalidEventDate   = "Invalid event date"
	ErrEventAlreadyActive = "Event is already active"
	ErrEventNotActive     = "Event is not active"

	// Attendance Errors
	ErrFailedMarkAttendance    = "Failed to mark attendance"
	ErrAttendanceAlreadyMarked = "Attendance already marked for this event"
	ErrInvalidQRCode           = "Invalid QR code"
	ErrQRCodeExpired           = "QR code has expired"

	// Database Errors
	ErrDatabaseConnection = "Database connection error"
	ErrFailedCreateRecord = "Failed to create record"
	ErrFailedUpdateRecord = "Failed to update record"
	ErrFailedDeleteRecord = "Failed to delete record"

	// Email Errors
	ErrFailedSendEmail = "Failed to send email"
)

// Success Messages
const (
	SuccessRegistration      = "Registration successful. Please check your email to verify your account."
	SuccessEmailVerified     = "Email verified successfully"
	SuccessLoginFacultyAdmin = "Login successful"
	SuccessTokenRefreshed    = "Token refreshed successfully"
	SuccessResetCodeSent     = "If your email is registered, you will receive a reset code."
	SuccessCodeValid         = "Code is valid"
	SuccessPasswordReset     = "Password reset successful"
	SuccessNewCodeSent       = "New code sent if email is registered"
	SuccessProfileFetched    = "Profile fetched successfully"
	SuccessUserPromoted      = "User promoted successfully"
	SuccessEventCreated      = "Event created successfully"
	SuccessEventUpdated      = "Event updated successfully"
	SuccessEventDeleted      = "Event deleted successfully"
	SuccessAttendanceMarked  = "Attendance marked successfully"
	SuccessAttendanceUpdated = "Attendance status updated successfully"
	SuccessUsersFetched      = "Users fetched successfully"
)

// Status Constants
const (
	// User Roles
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleFaculty    = "faculty"
	RoleStaff      = "staff"
	RoleStudent    = "student"

	// Event Status
	EventStatusScheduled = "scheduled"
	EventStatusOngoing   = "ongoing"
	EventStatusCompleted = "completed"
	EventStatusCancelled = "cancelled"

	// Attendance Status
	AttendanceStatusPresent = "present"
	AttendanceStatusAbsent  = "absent"
	AttendanceStatusLate    = "late"
	AttendanceStatusExcused = "excused"

	// QR Type
	QRTypeStudentID = "student_id"
	QRTypeEvent     = "event"
)
