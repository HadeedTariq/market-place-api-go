package mail

import "gopkg.in/gomail.v2"

type Mailer struct {
	dialer *gomail.Dialer
	from   string
}

func NewMailer(host string, port int, username, password, from string) *Mailer {
	d := gomail.NewDialer(host, port, username, password)

	return &Mailer{
		dialer: d,
		from:   from,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	msg.SetBody("text/html", body)

	return m.dialer.DialAndSend(msg)
}
