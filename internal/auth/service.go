package auth

type Service interface {
}

type svc struct {
}

func NewService() Service {
	return &svc{}
}
