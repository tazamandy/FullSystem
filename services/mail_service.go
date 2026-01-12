// services/mail.go

package services

import (
	"attendance-system/logging"
	"fmt"
	"html"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

// SendEmail sends a multipart email with both plain-text and HTML bodies.
// The caller should pass a fully-formed HTML body; the function will generate
// a plain-text fallback by stripping tags.
func SendEmail(to string, subject string, htmlBody string) error {
	plain := htmlToPlain(htmlBody)

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// Set plain text as primary body and add HTML as alternative
	m.SetBody("text/plain", plain)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		getSMTPPort(),
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)

	// Set timeout for SMTP connection
	d.LocalName = "attendance-system"

	err := d.DialAndSend(m)
	if err != nil {
		logging.Logger.Sugar().Errorf("Failed to send email to %s: %v", to, err)
		return err
	}
	return nil
}

func getSMTPPort() int {
	port := os.Getenv("SMTP_PORT")
	if portInt, err := strconv.Atoi(port); err == nil {
		return portInt
	}
	return 587
}

// htmlToPlain produces a simple plain-text fallback from HTML by removing tags
// and unescaping HTML entities. This is intentionally simple; for production
// consider a library for better conversions.
func htmlToPlain(s string) string {
	// Remove script/style blocks first (basic)
	reScripts := regexp.MustCompile(`(?s)<(script|style)[^>]*>.*?</(script|style)>`)
	s = reScripts.ReplaceAllString(s, "")

	// Strip remaining tags
	re := regexp.MustCompile(`<[^>]+>`)
	noTags := re.ReplaceAllString(s, " ")

	// Collapse whitespace
	collapsed := strings.Join(strings.Fields(noTags), " ")
	return html.UnescapeString(strings.TrimSpace(collapsed))
}

// BuildHTMLEmail wraps a fragment of HTML content into a modern, responsive
// email template with inline styles. Keep content minimal and mobile-friendly.
func BuildHTMLEmail(preheader, heading, contentHTML, footerHTML string) string {
	template := `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial; background:#f5f7fb; margin:0; padding:20px; }
		.container { max-width:600px; margin:0 auto; background:#ffffff; border-radius:8px; overflow:hidden; box-shadow:0 2px 8px rgba(0,0,0,0.08); }
		.header { background:linear-gradient(90deg,#4f46e5,#06b6d4); color:#fff; padding:20px; }
		.title { margin:0; font-size:20px; font-weight:600; }
		.content { padding:24px; color:#111827; line-height:1.5; font-size:15px; }
		.cta { display:inline-block; background:#4f46e5; color:#fff; padding:12px 20px; border-radius:6px; text-decoration:none; }
		.muted { color:#6b7280; font-size:13px; }
		.footer { padding:16px 24px; background:#fafafa; color:#6b7280; font-size:13px; }
		@media (max-width:420px) { .content { padding:16px } .header { padding:16px } }
	</style>
</head>
<body>
	<span style="display:none!important;visibility:hidden;mso-hide:all;">%s</span>
	<div class="container">
		<div class="header">
			<h1 class="title">%s</h1>
		</div>
		<div class="content">
			%s
		</div>
		<div class="footer">
			%s
		</div>
	</div>
</body>
</html>`

	return fmt.Sprintf(template, heading, preheader, heading, contentHTML, footerHTML)
}
