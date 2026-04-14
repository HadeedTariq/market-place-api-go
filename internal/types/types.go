package types

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}
