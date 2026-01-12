// utils/constants.go
package utils

const (
	HeaderStudentID     = "X-Student-ID"
	HeaderAuthorization = "Authorization"

	ErrInvalidEventID      = "Invalid event ID"
	ErrInvalidAttendanceID = "Invalid attendance ID"
	ErrUnauthorized        = "Unauthorized"
	ErrUserNotFound        = "User not found"
	// ReCAPTCHA
	ErrRecaptchaVerificationFailed = "recaptcha verification failed"
)

// Predefined Departments
var Departments = []string{
	"College of Education",
	"College of Engineering",
	"College of Science",
	"College of Arts and Sciences",
	"College of Business and Management",
	"College of Social Sciences",
	"College of Health Sciences",
	"College of Law",
	"College of Agriculture",
	"College of Medicine",
}

// Predefined Sections
var Sections = []string{
	"Section 1",
	"Section 2",
	"Section 3",
	"Section 4",
	"Section 5",
	"Section 6",
}
