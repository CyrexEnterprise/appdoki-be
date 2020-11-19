package app

import (
	"appdoki-be/app/repositories"
	"appdoki-be/config"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Application struct {
	conf            *config.Config
	usersRepository repositories.UsersRepositoryInterface
}

func NewApplication(conf *config.Config, db *sqlx.DB) *Application {
	return &Application{
		conf:            conf,
		usersRepository: repositories.NewUsersRepository(db),
	}
}

func (a *Application) Routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)

	a.AuthRouter(router)
	a.UsersRouter(router)

	fs := http.FileServer(http.Dir("./swaggerui/"))
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", fs))

	return middlewareChain([]middleware{
		trimSuffixMiddleware,
		loggingMiddleware,
	}, router)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	type homeResponse struct {
		DocsEndpoint string
	}
	respondJSON(w, homeResponse{
		DocsEndpoint: "/docs/",
	}, http.StatusOK)
}
