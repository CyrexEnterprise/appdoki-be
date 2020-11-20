package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (a *Application) BeersRouter(router *mux.Router) {
	beersHandler := NewBeersHandler(a.beersRepository)

	router.
		Methods(http.MethodGet).
		Path("/beers").
		HandlerFunc(a.JwtVerify(beersHandler.Get))
}
