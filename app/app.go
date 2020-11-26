package app

import (
	"appdoki-be/app/repositories"
	"appdoki-be/config"
	firebase "firebase.google.com/go/v4"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Application struct {
	conf            *config.Config
	firebaseApp     *firebase.App
	usersRepository repositories.UsersRepositoryInterface
	beersRepository repositories.BeersRepositoryInterface
}

func NewApplication(conf *config.Config, db *sqlx.DB, firebaseApp *firebase.App) *Application {
	return &Application{
		conf:            conf,
		firebaseApp:     firebaseApp,
		usersRepository: repositories.NewUsersRepository(db),
		beersRepository: repositories.NewBeersRepository(db),
	}
}

func (a *Application) Routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)

	a.AuthRouter(router)
	a.UsersRouter(router)
	a.BeersRouter(router)

	fs := http.FileServer(http.Dir("./swaggerui/"))
	router.
		PathPrefix("/docs").
		Handler(http.StripPrefix("/docs", fs))

	return middlewareChain([]middleware{
		trimSuffixMiddleware,
		loggingMiddleware,
	}, router)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	type TopicInfo struct {
		Topic       string
		Description string
	}
	type homeResponse struct {
		Version         string
		DocsEndpoint    string
		MessagingTopics []TopicInfo
	}

	res := homeResponse{
		Version:      "1.0.0",
		DocsEndpoint: "/docs/",
		MessagingTopics: []TopicInfo{
			{
				Topic:       "beers",
				Description: "Global beer transfer notifications",
			},
		},
	}

	respondJSON(w, res, http.StatusOK)
}
