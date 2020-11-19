package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (a *Application) UsersRouter(router *mux.Router) {
	usersHandler := NewUsersHandler(a.usersRepository)

	router.
		Methods(http.MethodGet).
		Path("/users").
		HandlerFunc(a.JwtVerify(usersHandler.Get))

	//router.
	//	Methods(http.MethodPost).
	//	Path("/users").
	//	HandlerFunc(usersHandler.Create)
	//
	//router.
	//	Methods(http.MethodPut).
	//	Path("/{id}").
	//	HandlerFunc(usersHandler.Update)
}
