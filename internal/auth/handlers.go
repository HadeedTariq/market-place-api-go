package auth

import (
	"encoding/json"
	"log"
	"net/http"

	utils "github.com/HadeedTariq/market-place-api-go/internal"
	"github.com/HadeedTariq/market-place-api-go/internal/types"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	service Service
}

func NewHandler(authService Service) *handler {
	return &handler{
		service: authService,
	}
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// ~ so first of all have to integrate some validation process in to that
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Invalid request body",
			Status:  404,
		})
	}

	err = validate.Struct(req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, e := range validationErrors {
			field := e.Field()

			switch field {
			case "Email":
				errors["email"] = "Invalid email format"
			case "Password":
				errors["password"] = "Password must be strong"
			case "UserName":
				errors["user_name"] = "Invalid username"
			case "Gender":
				errors["gender"] = "Invalid gender value"
			case "CountryCode":
				errors["country_code"] = "Invalid country code"
			}
		}
		utils.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"errors": errors,
		})
	}

	// ~ so over there the validation related stuff is done
	// ~ so now sqlc will comes because have to generate the query that with this email the user already exist
	existingUser, err := h.service.FindExistingUserByEmail(r.Context(), req.Email)
	if err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: err.Error(),
			Status:  404,
		})
	}

	if existingUser > 0 {
		utils.WriteJson(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Already exist user with this email",
			Status:  404,
		})
	}

	// ~ so now for the simplicity have to create the user with in the database

}
