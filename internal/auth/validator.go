package auth

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,strong_password"`
	UserName    string `json:"user_name" validate:"required,min=3,max=30,username"`
	Gender      string `json:"gender" validate:"required,oneof=male female other"`
	CountryCode string `json:"country_code" validate:"required,len=2,uppercase"`
}

type EmailOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required,len=6,numeric"`
}

func ValidateRequest(v *validator.Validate, payload interface{}) error {
	return v.Struct(payload)
}

func InitValidator() *validator.Validate {
	v := validator.New()

	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("strong_password", validatePassword)

	return v
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return re.MatchString(username)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}
