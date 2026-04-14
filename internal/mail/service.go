package mail

type EmailService struct {
	mailer *Mailer
}

func NewEmailService(m *Mailer) *EmailService {
	return &EmailService{mailer: m}
}

func (s *EmailService) SendWelcomeEmail(to string, name string) error {
	body := WelcomeTemplate(name)

	return s.mailer.Send(
		to,
		"Welcome to Our Platform",
		body,
	)
}
