package auth

import (
	"encoding/json"
	"log"
	"net/http"

	utils "github.com/HadeedTariq/market-place-api-go/internal"
	"github.com/HadeedTariq/market-place-api-go/internal/types"
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

	}
}
