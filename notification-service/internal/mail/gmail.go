package mail

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"notification-service/internal/model"
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

func (m *Mail) SendTaskMsgEmail(subject, body string) {
	emailPurpose := []string{"example@gmail.com"}

	var response *model.Response = &model.Response{}
	err := json.Unmarshal([]byte(body), response)
	if err != nil {
		fmt.Printf("❌ failed send task message to email. err = %v", err)
	}

	// Set up authentication information
	auth := smtp.PlainAuth("", m.Email, m.Password, "smtp.gmail.com")
	sendEmail := fmt.Sprintf("Status = %s\nData = %s", response.Message, response.Data)
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, sendEmail))

	// Mengirim email via Gmail SMTP
	err = smtp.SendMail("smtp.gmail.com:587", auth, m.Email, emailPurpose, msg)
	if err != nil {
		fmt.Printf("❌ Failed to send task message to email: %s", err)
	} else {
		fmt.Printf("✔️ Task Message sent successfully to Email. with subject: %s", subject)
	}
}
