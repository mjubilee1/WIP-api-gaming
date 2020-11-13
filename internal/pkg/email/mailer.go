package email

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"api-gaming/internal/util"
)

type Request struct {
	From    string
	TO      []string
	Subject string
	Body    string
}

// MIME - formatting
var (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	server = util.ViperEnvVariable("SMTP_SERVER")
	port = util.ViperEnvVariable("SMTP_PORT")
	email = util.ViperEnvVariable("SMTP_EMAIL")
	password = util.ViperEnvVariable("SMTP_PASSWORD")
	t *template.Template
)

func init() {
	fileName := "web/templates/verification-email.html"
	t = template.Must(template.ParseFiles(fileName))
}

// NewRequest - sets up the email subject line and who to send to.
func NewRequest(to []string, subject string) *Request {
	return &Request{
		TO:      to,
		Subject: subject,
	}
}

func (r *Request) parseTemplate(data interface{}) error {
	buffer := new(bytes.Buffer)
	if err := t.Execute(buffer, data); err != nil {
		log.Println("Failed to execute template: ", err)
	}

	r.Body = buffer.String()
	return nil
}

func (r *Request) sendMail() bool {
	addr := "smtp.gmail.com:587"
	body := "To: " + r.TO[0] + "\r\nSubject: " + r.Subject + "\r\n" + MIME + "\r\n" + r.Body

	if err := smtp.SendMail(addr, smtp.PlainAuth("", email, password, server), email, r.TO, []byte(body)); err != nil {
		return false
	}
	return true
}

// Send - takes in the html email template and sends a request to the SMTP server
func (r *Request) Send(items interface{}) {
	err := r.parseTemplate(items)
	if err != nil {
		log.Fatal(err)
	}

	if ok := r.sendMail(); ok {
		log.Printf("Email has been sent to %s\n", r.TO)
	} else {
		log.Printf("Failed to send the email to %s\n", r.TO)
	}
}