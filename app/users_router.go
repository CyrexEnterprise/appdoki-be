package app

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (a *Application) UsersRouter(router *mux.Router) {
	notifierSrv, err := newNotifier(a.firebaseApp)
	if err != nil {
		log.Fatal("could not instantiate a notifier")
	}
	usersHandler := NewUsersHandler(a.usersRepository, notifierSrv)

	router.
		Methods(http.MethodGet).
		Path("/users").
		HandlerFunc(a.JwtVerify(usersHandler.Get))

	router.
		Methods(http.MethodGet).
		Path("/users/{id}").
		HandlerFunc(a.JwtVerify(usersHandler.GetByID))

	router.
		Methods(http.MethodGet).
		Path("/users/{id}/beers").
		HandlerFunc(a.JwtVerify(usersHandler.BeersSummary))

	router.
		Methods(http.MethodPost).
		Path("/users/{id}/beers/{beers}").
		HandlerFunc(a.JwtVerify(usersHandler.GiveBeers))
}
