package auth

import "net/http"

type handler struct {
	service Service
}

func NewHandler(authService Service) *handler {
	return &handler{
		service: authService,
	}
}

func (h *handler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	// ~ ok so what this handler is going to do is based on the token will return the user by taking from the middleware
}
