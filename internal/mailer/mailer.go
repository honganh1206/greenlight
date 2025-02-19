package mailer

import (
	"bytes"
	"embed"
	"text/template"
	"time"
)

// Embed templates directly into the compiled binary
// These templates will be read during compile time

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	Dialer *Dialer
	Sender string
}

type MailerConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

func New(host string, port int, username, password, sender string) *Mailer {
	dialer := NewDialer(host, port, username, password)

	dialer.Timeout = 5 * time.Second
	return &Mailer{
		Dialer: dialer,
		Sender: sender,
	}
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)

	if err != nil {
		return err
	}

	// Execute the "subject" template
	// Passing the dynamic data and store inside a buffer variable
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Execute the "plainBody" template
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	// fmt.Println(plainBody)

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.Sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// FIXME: Remove hard-coded retry times
	for i := 1; i <= 3; i++ {
		err = m.Dialer.DialAndSend(msg)
		if err == nil {
			return nil
		}
		// Sleep for a short time and retry
		time.Sleep(500 * time.Millisecond)
	}

	return err
}
