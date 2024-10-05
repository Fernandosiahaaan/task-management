package mail

import (
	"fmt"
	"net/smtp"
)

type Mail struct {
	Email    string
	Password string
}

func Init(email, password string) (*Mail, error) {
	var mail *Mail = &Mail{
		Email:    email,
		Password: password,
	}
	return mail, nil
}

func (m *Mail) SendEmail(subject, body string) {
	to := []string{"example@gmail.com"}

	// Set up authentication information
	auth := smtp.PlainAuth("", m.Email, m.Password, "smtp.gmail.com")

	// Pesan email
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	// Mengirim email via Gmail SMTP
	err := smtp.SendMail("smtp.gmail.com:587", auth, m.Email, to, msg)
	if err != nil {
		fmt.Printf("❌ Failed to send email: %s", err)
	} else {
		fmt.Printf("✔️ Email sent successfully with subject: %s", subject)
	}
}
