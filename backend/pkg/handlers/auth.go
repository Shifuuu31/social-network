package handlers

import (
	"net/http"
)

func (rt *Root) NewUserHandler() (userMux *http.ServeMux) {
	userMux = http.NewServeMux()

	userMux.HandleFunc("POST /signup", rt.SignUp)
	userMux.HandleFunc("POST /signin", rt.SignIn)
	userMux.HandleFunc("DELETE /signout", rt.SignOut)

	return userMux
}

func (rt *Root) SignUp(w http.ResponseWriter, r *http.Request) {
	// TODO verify user input

	// TODO hash user password

	// TODO insert user into db
}

func (rt *Root) SignIn(w http.ResponseWriter, r *http.Request) {
	// TODO verify user credentials

	// TODO set user session
}

func (rt *Root) SignOut(w http.ResponseWriter, r *http.Request) {
	// TODO unset user session
}
