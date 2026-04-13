package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	utils "github.com/HadeedTariq/market-place-api-go/internal"
	repo "github.com/HadeedTariq/market-place-api-go/internal/adapters/postgresql/sqlc"
	"github.com/HadeedTariq/market-place-api-go/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

type handler struct {
	service   Service
	validator *validator.Validate
}

func NewHandler(authService Service, validator *validator.Validate) *handler {
	return &handler{
		service:   authService,
		validator: validator,
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

	err = h.validator.Struct(req)

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

	if err != nil {
		if err == pgx.ErrNoRows {
			existingUser = 0
		} else {
			log.Println(err)
			utils.WriteJson(w, http.StatusInternalServerError, types.ErrorResponse{
				Message: "Internal Server Error",
				Status:  500,
			})
			return
		}
	}

	if existingUser > 0 {
		utils.WriteJson(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Already exist user with this email",
			Status:  400,
		})
		return
	}

	// ~ so now for the simplicity have to create the user with in the database
	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Internal Server Error",
			Status:  500,
		})
		return

	}

	err = h.service.InsertUser(r.Context(), repo.InsertUserParams{
		UserName:     req.UserName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "user",
		Gender:       req.Gender,
		Source:       "general",
	})

	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: err.Error(),
			Status:  500,
		})
		return
	}

	utils.WriteJson(w, http.StatusCreated, types.ErrorResponse{
		Message: "User registered successfully",
		Status:  201,
	})
}
