package util

import (
	"gopkg.in/gomail.v2"
)

// Gomail is a function for send email
func Gomail(host, password string, port int, from, to, subject, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	d := gomail.NewDialer(host, port, from, password)
	d.SSL = true

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
