package authhandler

import "net/http"

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthenticationRestApiProvider interface {
	Register(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}
