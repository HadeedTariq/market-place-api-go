package mail

type EmailService struct {
	mailer *Mailer
}

func NewEmailService(m *Mailer) *EmailService {
	return &EmailService{mailer: m}
}

func (s *EmailService) SendOtpEmail(to string, otp string, expiryMinutes int) error {
	body := OTPTemplate("Accoswap", otp, expiryMinutes)

	return s.mailer.Send(
		to,
		"One time OTP for Accoswap",
		body,
	)
}
