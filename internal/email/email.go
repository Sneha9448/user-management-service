package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"user-management-service/internal/config"
)

var smtpCfg *config.Config

// Init initializes the email package with the required configuration
func Init(cfg *config.Config) {
	smtpCfg = cfg
}

// SendOTPEmail sends a 6-digit OTP code to the specified email address.
// It uses the configuration provided during Init().
func SendOTPEmail(to string, otp string) error {
	if smtpCfg == nil {
		return fmt.Errorf("email package not initialized")
	}

	// For development: If SMTPEmail is not set, log the OTP to a file
	if smtpCfg.SMTPEmail == "" {
		log.Printf("DEMO MODE: Sending OTP %s to %s\n", otp, to)
		if err := os.WriteFile("otp_debug.log", []byte(otp), 0644); err != nil {
			log.Printf("Failed to write otp_debug.log: %v", err)
		}
		return nil
	}

	from := smtpCfg.SMTPEmail
	password := smtpCfg.SMTPPassword
	smtpHost := smtpCfg.SMTPHost
	smtpPort := smtpCfg.SMTPPort

	// Message composition
	subject := "Your OTP Code"
	body := fmt.Sprintf("Your 6-digit verification code is: %s\nThis code expires in 10 minutes.", otp)
	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("OTP email successfully sent to %s", to)
	return nil
}
