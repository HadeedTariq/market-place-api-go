package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	utils "github.com/HadeedTariq/market-place-api-go/internal"
	repo "github.com/HadeedTariq/market-place-api-go/internal/adapters/postgresql/sqlc"
	"github.com/HadeedTariq/market-place-api-go/internal/mail"
	"github.com/HadeedTariq/market-place-api-go/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service      Service
	validator    *validator.Validate
	emailService *mail.EmailService
}

func NewHandler(authService Service, validator *validator.Validate, emailService *mail.EmailService) *handler {
	return &handler{
		service:      authService,
		validator:    validator,
		emailService: emailService,
	}
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Invalid request body",
			Status:  404,
		})
		return
	}

	err = ValidateRequest(h.validator, &req)

	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			customErrors := make(map[string]string)

			for _, e := range validationErrors {
				field := e.Field()

				switch field {
				case "Email":
					customErrors["email"] = "Invalid email format"
				case "Password":
					customErrors["password"] = "Password must be strong"
				case "UserName":
					customErrors["user_name"] = "Invalid username"
				case "Gender":
					customErrors["gender"] = "Invalid gender value"
				case "CountryCode":
					customErrors["country_code"] = "Invalid country code"
				}
			}

			utils.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
				"errors": customErrors,
			})
			return
		} else {
			utils.WriteJson(w, http.StatusBadRequest, types.ErrorResponse{
				Message: "Invalid request",
				Status:  400,
			})
			return
		}
	}

	existingUser, err := h.service.FindExistingUserByEmail(r.Context(), req.Email)

	if err != nil && err != pgx.ErrNoRows {
		log.Println(err)
		utils.WriteJson(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Internal Server Error",
			Status:  500,
		})
		return
	}

	if existingUser != nil {
		utils.WriteJson(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "User already exists with this email",
			Status:  400,
		})
		return
	}

	existingOtp, err := h.service.FindExistingOtp(r.Context(), req.Email)

	if err != nil && err != pgx.ErrNoRows {
		log.Println(err)
		utils.WriteJson(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Internal Server Error",
			Status:  500,
		})
		return
	}

	if existingOtp > 0 {
		utils.WriteJson(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Otp Already send",
			Status:  400,
		})
		return
	}
	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Internal Server Error",
			Status:  500,
		})
		return
	}

	role := "user"
	source := "general"

	err = h.service.InsertUser(r.Context(), repo.InsertUserParams{
		UserName:     &req.UserName,
		Email:        req.Email,
		PasswordHash: &hashedPassword,
		Role:         &role,
		Gender:       &req.Gender,
		Source:       &source,
	})

	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: err.Error(),
			Status:  500,
		})
		return
	}
	otp, err := utils.GenerateOTP(6)

	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: err.Error(),
			Status:  500,
		})
		return
	}
	expiresAt := time.Now().Add(5 * time.Minute)

	err = h.service.InsertEmailOtp(r.Context(), repo.InsertEmailOtpParams{
		Email: req.Email,
		Otp:   otp,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: true,
		},
	})

	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: err.Error(),
			Status:  500,
		})
		return
	}

	err = h.emailService.SendOtpEmail(req.Email, otp, 5)

	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: err.Error(),
			Status:  500,
		})
		return
	}

	utils.WriteJson(w, http.StatusCreated, types.ErrorResponse{
		Message: "Otp send on your email please verify to register",
		Status:  201,
	})
}

func (h *handler) OtpEmailChecker(w http.ResponseWriter, r *http.Request) {
	var req EmailOtpRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Invalid request body",
			Status:  404,
		})
		return
	}

	// ~ so over there have to create the validator for that
	err = ValidateRequest(h.validator, &req)

	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			customErrors := make(map[string]string)
			for _, e := range validationErrors {
				field := e.Field()
				switch field {
				case "Email":
					customErrors["email"] = "Invalid email format"
				case "Otp":
					customErrors["otp"] = "Invalid otp format"
				}
			}
			utils.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
				"errors": customErrors,
			})
			return
		} else {
			utils.WriteJson(w, http.StatusBadRequest, types.ErrorResponse{
				Message: "Invalid request",
				Status:  400,
			})
			return
		}
	}

	correctOtp, err := h.service.CheckOtp(r.Context(), repo.CheckOtpParams{
		Email: req.Email,
		Otp:   req.Otp,
	})

	if correctOtp < 1 {
		utils.WriteJson(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Otp not correct",
			Status:  400,
		})
		return
	}

}
